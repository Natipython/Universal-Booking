import 'dart:convert';

import 'package:sqflite_common/sqlite_api.dart';

import '../features/bookings/domain/booking.dart';

/// V0 only: unencrypted synthetic SQLite store, never production data/auth.
class DemoRepository {
  DemoRepository._(this.database, this.fixtures);
  final Database database;
  final DemoFixtures fixtures;

  static Future<DemoRepository> open({
    required DatabaseFactory factory,
    required String path,
    required DemoFixtures fixtures,
  }) async {
    final db = await factory.openDatabase(
      path,
      options: OpenDatabaseOptions(
        version: 1,
        singleInstance: false,
        onConfigure: (db) async {
          await db.execute('PRAGMA foreign_keys = ON');
          await db.execute('PRAGMA synchronous = FULL');
        },
        onCreate: (db, version) async {
          await db.execute('''CREATE TABLE operations (
          organization_id TEXT NOT NULL, id TEXT NOT NULL,
          PRIMARY KEY (organization_id, id))''');
          await db.execute('''CREATE TABLE services (
          organization_id TEXT NOT NULL, operation_id TEXT NOT NULL,
          id TEXT NOT NULL, PRIMARY KEY (organization_id, operation_id, id),
          FOREIGN KEY (organization_id, operation_id)
            REFERENCES operations (organization_id, id))''');
          await db.execute('''CREATE TABLE bookings (
          organization_id TEXT NOT NULL, operation_id TEXT NOT NULL, id TEXT NOT NULL,
          service_id TEXT NOT NULL, service_name TEXT NOT NULL,
          resource_id TEXT NOT NULL, customer_label TEXT NOT NULL, source TEXT NOT NULL,
          start_ms INTEGER NOT NULL, end_ms INTEGER NOT NULL CHECK (end_ms > start_ms),
          price_minor INTEGER NOT NULL CHECK (price_minor >= 0),
          status TEXT NOT NULL CHECK (status IN ('pending','simulatedConflict','cancelled')),
          revision INTEGER NOT NULL CHECK (revision >= 1),
          PRIMARY KEY (organization_id, id), UNIQUE (organization_id, operation_id, id),
          FOREIGN KEY (organization_id, operation_id, service_id)
            REFERENCES services (organization_id, operation_id, id))''');
          await db.execute('''CREATE INDEX booking_scope_time
          ON bookings (organization_id, operation_id, start_ms)''');
          await db.execute('''CREATE TABLE outbox (
          organization_id TEXT NOT NULL, operation_id TEXT NOT NULL,
          command_id TEXT NOT NULL, booking_id TEXT NOT NULL,
          payload TEXT NOT NULL, state TEXT NOT NULL
            CHECK (state IN ('pending','simulatedConflict','cancelled')),
          PRIMARY KEY (organization_id, command_id),
          FOREIGN KEY (organization_id, operation_id, booking_id)
            REFERENCES bookings (organization_id, operation_id, id))''');
          await db.execute('''CREATE TABLE local_audit (
          event_id INTEGER PRIMARY KEY AUTOINCREMENT,
          organization_id TEXT NOT NULL, operation_id TEXT NOT NULL,
          booking_id TEXT NOT NULL, action TEXT NOT NULL,
          before_json TEXT, after_json TEXT NOT NULL,
          FOREIGN KEY (organization_id, operation_id, booking_id)
            REFERENCES bookings (organization_id, operation_id, id))''');
          await db.execute('''CREATE TABLE drafts (
          organization_id TEXT NOT NULL, operation_id TEXT NOT NULL,
          payload TEXT NOT NULL, PRIMARY KEY (organization_id, operation_id),
          FOREIGN KEY (organization_id, operation_id)
            REFERENCES operations (organization_id, id))''');
          for (final operation in fixtures.operations) {
            await db.insert('operations', {
              'organization_id': operation.organizationId,
              'id': operation.id,
            });
          }
          for (final service in fixtures.services) {
            await db.insert('services', {
              'organization_id': service.organizationId,
              'operation_id': service.operationId,
              'id': service.id,
            });
          }
        },
      ),
    );
    return DemoRepository._(db, fixtures);
  }

  void _checkScope(DemoScope scope) {
    if (!fixtures.operations.any(
      (o) =>
          o.organizationId == scope.organizationId && o.id == scope.operationId,
    )) {
      throw ArgumentError('Unknown demo organization/operation.');
    }
  }

  DemoService _validate(BookingIntent intent) {
    _checkScope(intent.scope);
    final service = fixtures.serviceFor(intent.scope, intent.serviceId);
    const labels = {'Demo Guest 01', 'Demo Guest 02', 'Demo Group 03'};
    const sources = {'Call', 'Walk-in', 'WhatsApp', 'Telegram', 'Other'};
    final local = fixtures.utcToLocal(intent.startUtc);
    final end = local.add(Duration(minutes: service.durationMinutes));
    if (!labels.contains(intent.customerLabel) ||
        !sources.contains(intent.source) ||
        local.hour < 9 ||
        local.minute % 30 != 0 ||
        local.second != 0 ||
        local.millisecond != 0 ||
        local.microsecond != 0 ||
        end.hour > 18 ||
        (end.hour == 18 && end.minute != 0) ||
        end.day != local.day ||
        local.year < fixtures.referenceDate.year - 1 ||
        local.year > fixtures.referenceDate.year + 1) {
      throw ArgumentError(
        'Use the demo choices and a half-hour start within 09:00–18:00.',
      );
    }
    return service;
  }

  Future<List<LocalBooking>> list(
    DemoScope scope, {
    bool includeCancelled = false,
  }) async {
    _checkScope(scope);
    final rows = await database.query(
      'bookings',
      where:
          'organization_id = ? AND operation_id = ?${includeCancelled ? '' : ' AND status != ?'}',
      whereArgs: [
        scope.organizationId,
        scope.operationId,
        if (!includeCancelled) LocalStatus.cancelled.name,
      ],
      orderBy: 'start_ms ASC',
    );
    return rows.map(LocalBooking.fromRow).toList();
  }

  Future<LocalBooking> get(DemoScope scope, String id) async {
    _checkScope(scope);
    final rows = await database.query(
      'bookings',
      where: 'organization_id = ? AND operation_id = ? AND id = ?',
      whereArgs: [scope.organizationId, scope.operationId, id],
      limit: 1,
    );
    if (rows.isEmpty) throw StateError('Demo item is not in this scope.');
    return LocalBooking.fromRow(rows.single);
  }

  Future<void> _checkLocalOverlap(
    Transaction tx,
    BookingIntent intent,
    DemoService service, {
    String? exceptId,
  }) async {
    final start = intent.startUtc.toUtc().millisecondsSinceEpoch;
    final end = start + service.durationMinutes * 60000;
    final rows = await tx.query(
      'bookings',
      columns: ['id'],
      where:
          'organization_id = ? AND resource_id = ? AND status != ? '
          'AND start_ms < ? AND end_ms > ?${exceptId == null ? '' : ' AND id != ?'}',
      whereArgs: [
        intent.scope.organizationId,
        service.resourceId,
        LocalStatus.cancelled.name,
        end,
        start,
        ?exceptId,
      ],
      limit: 1,
    );
    if (rows.isNotEmpty) throw const LocalBookingConflict();
  }

  Future<LocalBooking> save(
    BookingIntent intent, {
    required String commandId,
  }) async {
    final service = _validate(intent);
    if (commandId.isEmpty || commandId.length > 128) {
      throw ArgumentError('Invalid command ID.');
    }
    final payload = jsonEncode({
      'type': 'createLocalIntent',
      ...intent.toJson(),
    });
    final id = await database.transaction((tx) async {
      final old = await tx.query(
        'outbox',
        where: 'organization_id = ? AND command_id = ?',
        whereArgs: [intent.scope.organizationId, commandId],
      );
      if (old.isNotEmpty) {
        if (old.single['payload'] != payload ||
            old.single['operation_id'] != intent.scope.operationId) {
          throw StateError(
            'Command ID already belongs to different input. Nothing was overwritten.',
          );
        }
        return old.single['booking_id'] as String;
      }
      await _checkLocalOverlap(tx, intent, service);
      final id = newDemoId();
      final row = {
        'organization_id': intent.scope.organizationId,
        'operation_id': intent.scope.operationId,
        'id': id,
        'service_id': service.id,
        'service_name': service.name,
        'resource_id': service.resourceId,
        'customer_label': intent.customerLabel,
        'source': intent.source,
        'start_ms': intent.startUtc.toUtc().millisecondsSinceEpoch,
        'end_ms':
            intent.startUtc.toUtc().millisecondsSinceEpoch +
            service.durationMinutes * 60000,
        'price_minor': service.priceMinor,
        'status': LocalStatus.pending.name,
        'revision': 1,
      };
      await tx.insert('bookings', row);
      await tx.insert('outbox', {
        'organization_id': intent.scope.organizationId,
        'operation_id': intent.scope.operationId,
        'command_id': commandId,
        'booking_id': id,
        'payload': payload,
        'state': LocalStatus.pending.name,
      });
      await tx.insert('local_audit', {
        'organization_id': intent.scope.organizationId,
        'operation_id': intent.scope.operationId,
        'booking_id': id,
        'action': 'local_intent_saved',
        'after_json': jsonEncode(row),
      });
      await tx.delete(
        'drafts',
        where: 'organization_id = ? AND operation_id = ?',
        whereArgs: [intent.scope.organizationId, intent.scope.operationId],
      );
      return id;
    });
    return get(intent.scope, id);
  }

  Future<void> saveDraft(BookingIntent intent) async {
    _validate(intent);
    await database.insert('drafts', {
      'organization_id': intent.scope.organizationId,
      'operation_id': intent.scope.operationId,
      'payload': jsonEncode(intent.toJson()),
    }, conflictAlgorithm: ConflictAlgorithm.replace);
  }

  Future<BookingIntent?> readDraft(DemoScope scope) async {
    _checkScope(scope);
    final rows = await database.query(
      'drafts',
      where: 'organization_id = ? AND operation_id = ?',
      whereArgs: [scope.organizationId, scope.operationId],
    );
    return rows.isEmpty
        ? null
        : BookingIntent.fromJson(
            jsonDecode(rows.single['payload'] as String)
                as Map<String, dynamic>,
          );
  }

  Future<void> changeStatus(
    DemoScope scope,
    String id,
    LocalStatus status,
  ) async {
    _checkScope(scope);
    if (status == LocalStatus.pending) {
      throw ArgumentError('Use a new reviewed local intent to retry.');
    }
    await database.transaction((tx) async {
      final rows = await tx.query(
        'bookings',
        where: 'organization_id = ? AND operation_id = ? AND id = ?',
        whereArgs: [scope.organizationId, scope.operationId, id],
      );
      if (rows.isEmpty) throw StateError('Demo item is outside this scope.');
      final before = rows.single;
      if (before['status'] == status.name) return;
      if (before['status'] == LocalStatus.cancelled.name) {
        throw StateError('Cancelled local item cannot change.');
      }
      final next = {
        ...before,
        'status': status.name,
        'revision': (before['revision'] as int) + 1,
      };
      await tx.update(
        'bookings',
        {'status': next['status'], 'revision': next['revision']},
        where: 'organization_id = ? AND operation_id = ? AND id = ?',
        whereArgs: [scope.organizationId, scope.operationId, id],
      );
      await tx.update(
        'outbox',
        {'state': status.name},
        where: 'organization_id = ? AND operation_id = ? AND booking_id = ?',
        whereArgs: [scope.organizationId, scope.operationId, id],
      );
      await tx.insert('local_audit', {
        'organization_id': scope.organizationId,
        'operation_id': scope.operationId,
        'booking_id': id,
        'action': status == LocalStatus.cancelled
            ? 'local_intent_cancelled'
            : 'simulated_conflict',
        'before_json': jsonEncode(before),
        'after_json': jsonEncode(next),
      });
    });
  }

  Future<LocalBooking> revise(
    LocalBooking expected,
    BookingIntent intent, {
    required String commandId,
  }) async {
    final service = _validate(intent);
    if (expected.scope.key != intent.scope.key) {
      throw ArgumentError('Cannot move a booking between demo scopes.');
    }
    if (commandId.isEmpty || commandId.length > 128) {
      throw ArgumentError('Invalid command ID.');
    }
    final payload = jsonEncode({
      'type': 'reviseLocalIntent',
      'bookingId': expected.id,
      'expectedRevision': expected.revision,
      ...intent.toJson(),
    });
    await database.transaction((tx) async {
      final prior = await tx.query(
        'outbox',
        where: 'organization_id = ? AND command_id = ?',
        whereArgs: [intent.scope.organizationId, commandId],
      );
      if (prior.isNotEmpty) {
        if (prior.single['payload'] != payload) {
          throw StateError('Command ID input differs.');
        }
        return;
      }
      final rows = await tx.query(
        'bookings',
        where: 'organization_id = ? AND operation_id = ? AND id = ?',
        whereArgs: [
          intent.scope.organizationId,
          intent.scope.operationId,
          expected.id,
        ],
      );
      if (rows.isEmpty ||
          rows.single['revision'] != expected.revision ||
          rows.single['status'] == LocalStatus.cancelled.name) {
        throw StateError(
          'This demo item changed. Reopen it before editing. Your input is kept.',
        );
      }
      await _checkLocalOverlap(tx, intent, service, exceptId: expected.id);
      final before = rows.single;
      final next = {
        ...before,
        'service_id': service.id,
        'service_name': service.name,
        'resource_id': service.resourceId,
        'customer_label': intent.customerLabel,
        'source': intent.source,
        'start_ms': intent.startUtc.toUtc().millisecondsSinceEpoch,
        'end_ms':
            intent.startUtc.toUtc().millisecondsSinceEpoch +
            service.durationMinutes * 60000,
        'price_minor': service.priceMinor,
        'status': LocalStatus.pending.name,
        'revision': expected.revision + 1,
      };
      // Old commands retain immutable payloads, but cannot be dispatched as current work.
      await tx.update(
        'outbox',
        {'state': LocalStatus.cancelled.name},
        where: 'organization_id = ? AND operation_id = ? AND booking_id = ? AND state = ?',
        whereArgs: [
          intent.scope.organizationId,
          intent.scope.operationId,
          expected.id,
          LocalStatus.pending.name,
        ],
      );
      await tx.update(
        'bookings',
        next,
        where: 'organization_id = ? AND operation_id = ? AND id = ?',
        whereArgs: [
          intent.scope.organizationId,
          intent.scope.operationId,
          expected.id,
        ],
      );
      await tx.insert('outbox', {
        'organization_id': intent.scope.organizationId,
        'operation_id': intent.scope.operationId,
        'command_id': commandId,
        'booking_id': expected.id,
        'payload': payload,
        'state': LocalStatus.pending.name,
      });
      await tx.insert('local_audit', {
        'organization_id': intent.scope.organizationId,
        'operation_id': intent.scope.operationId,
        'booking_id': expected.id,
        'action': 'local_intent_revised',
        'before_json': jsonEncode(before),
        'after_json': jsonEncode(next),
      });
      await tx.delete(
        'drafts',
        where: 'organization_id = ? AND operation_id = ?',
        whereArgs: [intent.scope.organizationId, intent.scope.operationId],
      );
    });
    return get(intent.scope, expected.id);
  }

  Future<List<Map<String, Object?>>> history(DemoScope scope, String id) async {
    await get(scope, id);
    return database.query(
      'local_audit',
      where: 'organization_id = ? AND operation_id = ? AND booking_id = ?',
      whereArgs: [scope.organizationId, scope.operationId, id],
      orderBy: 'event_id ASC',
    );
  }

  Future<void> close() => database.close();
}

import 'dart:convert';
import 'dart:io';

import 'package:flutter_test/flutter_test.dart';
import 'package:path/path.dart' as path;
import 'package:sqflite_common_ffi/sqflite_ffi.dart';
import 'package:ubi_business_v0/features/bookings/domain/booking.dart';
import 'package:ubi_business_v0/local/demo_repository.dart';

void main() {
  late DemoFixtures fixtures;
  late DemoRepository repository;
  late Directory temp;
  const a = DemoScope('demo-org-a', 'studio');
  const b = DemoScope('demo-org-b', 'studio');
  const event = DemoScope('demo-org-a', 'events');
  BookingIntent intent({
    DemoScope scope = a,
    String service = 'portrait',
    int hour = 9,
    int minute = 0,
    String customer = 'Demo Guest 01',
  }) => BookingIntent(
    scope: scope,
    serviceId: service,
    customerLabel: customer,
    source: 'Call',
    startUtc: fixtures.localToUtc(fixtures.referenceDate, hour, minute),
  );

  setUpAll(() {
    sqfliteFfiInit();
  });
  setUp(() async {
    fixtures = DemoFixtures.parse(
      await File('../fixtures/booking-v0.json').readAsString(),
    );
    temp = await Directory.systemTemp.createTemp('ubi-v0-store-');
    repository = await DemoRepository.open(
      factory: databaseFactoryFfi,
      path: path.join(temp.path, 'test.sqlite'),
      fixtures: fixtures,
    );
  });
  tearDown(() async {
    await repository.close();
    await temp.delete(recursive: true);
  });

  test('fixtures are explicitly synthetic, two tenants, two operations in one tenant', () {
    expect(fixtures.organizations.length, 2);
    expect(
      fixtures.operations
          .where((o) => o.organizationId == a.organizationId)
          .length,
      2,
    );
    final json = jsonDecode(
      File('../fixtures/booking-v0.json').readAsStringSync(),
    ) as Map<String, dynamic>;
    json['syntheticOnly'] = false;
    expect(() => DemoFixtures.parse(jsonEncode(json)), throwsFormatException);
  });

  test('atomic local save produces pending booking, command and history; never confirmed', () async {
    final saved = await repository.save(intent(), commandId: 'test-command');
    expect(saved.status, LocalStatus.pending);
    expect(saved.priceMinor, 200000);
    expect(saved.endUtc.difference(saved.startUtc).inMinutes, 60);
    expect((await repository.database.query('outbox')).length, 1);
    expect(
      (await repository.history(a, saved.id)).single['action'],
      'local_intent_saved',
    );
  });

  test('same command and input return one logical item', () async {
    final first = await repository.save(intent(), commandId: 'same-command');
    final retry = await repository.save(intent(), commandId: 'same-command');
    expect(first.id, retry.id);
    expect((await repository.list(a)).length, 1);
    expect((await repository.history(a, first.id)).length, 1);
  });

  test(
    'same command with changed input fails without rewriting old intent',
    () async {
      final original = await repository.save(
        intent(),
        commandId: 'same-command',
      );
      await expectLater(
        repository.save(intent(hour: 11), commandId: 'same-command'),
        throwsStateError,
      );
      expect(
        (await repository.get(a, original.id)).startUtc,
        original.startUtc,
      );
    },
  );

  test(
    'outbox insert failure rolls back booking and audit; draft remains',
    () async {
      await repository.saveDraft(intent());
      await repository.database.execute(
        "CREATE TRIGGER reject_test_command BEFORE INSERT ON outbox WHEN NEW.command_id = 'fail-command' BEGIN SELECT RAISE(ABORT, 'synthetic outbox failure'); END",
      );
      await expectLater(
        repository.save(intent(), commandId: 'fail-command'),
        throwsA(isA<DatabaseException>()),
      );
      expect(await repository.list(a), isEmpty);
      expect(await repository.database.query('local_audit'), isEmpty);
      expect(await repository.database.query('outbox'), isEmpty);
      expect(await repository.readDraft(a), isNotNull);
    },
  );

  test(
    'local overlap rejects without losing draft; adjacent interval succeeds',
    () async {
      await repository.save(intent(), commandId: 'first');
      await repository.saveDraft(intent(minute: 30));
      await expectLater(
        repository.save(intent(minute: 30), commandId: 'second'),
        throwsA(isA<LocalBookingConflict>()),
      );
      expect(
        (await repository.readDraft(a))!.startUtc,
        intent(minute: 30).startUtc,
      );
      await repository.save(intent(hour: 10), commandId: 'adjacent');
      expect((await repository.list(a)).length, 2);
    },
  );

  test(
    'concurrent local contenders serialize and create one pending item',
    () async {
      final results = await Future.wait(
        List.generate(12, (i) async {
          try {
            await repository.save(intent(), commandId: 'contender-$i');
            return true;
          } on LocalBookingConflict {
            return false;
          }
        }),
      );
      expect(results.where((s) => s).length, 1);
      expect((await repository.list(a)).length, 1);
    },
  );

  test(
    'tenant scope prevents reads, history and changes of another tenant item',
    () async {
      final one = await repository.save(intent(), commandId: 'shared-command');
      final two = await repository.save(
        intent(scope: b),
        commandId: 'shared-command',
      );
      expect((await repository.list(a)).single.id, one.id);
      expect((await repository.list(b)).single.id, two.id);
      await expectLater(repository.get(b, one.id), throwsStateError);
      await expectLater(repository.history(b, one.id), throwsStateError);
      await expectLater(
        repository.changeStatus(b, one.id, LocalStatus.cancelled),
        throwsStateError,
      );
      expect((await repository.get(a, one.id)).status, LocalStatus.pending);
    },
  );

  test('operation scope is independent inside the same organization', () async {
    final studio = await repository.save(intent(), commandId: 'studio');
    final consultation = await repository.save(
      intent(scope: event, service: 'consultation'),
      commandId: 'event',
    );
    expect((await repository.list(a)).single.id, studio.id);
    expect((await repository.list(event)).single.id, consultation.id);
    await expectLater(repository.get(event, studio.id), throwsStateError);
  });

  test('scoped composite FK refuses a cross-operation service', () async {
    final saved = await repository.save(intent(), commandId: 'first');
    await expectLater(
      repository.database.update(
        'bookings',
        {'operation_id': 'events'},
        where: 'id = ?',
        whereArgs: [saved.id],
      ),
      throwsA(isA<DatabaseException>()),
    );
  });

  test('draft persists across close/reopen, scoped to one operation', () async {
    await repository.saveDraft(intent(hour: 11, customer: 'Demo Guest 02'));
    final dbPath = repository.database.path;
    await repository.close();
    repository = await DemoRepository.open(
      factory: databaseFactoryFfi,
      path: dbPath,
      fixtures: fixtures,
    );
    expect((await repository.readDraft(a))!.customerLabel, 'Demo Guest 02');
    expect(
      (await repository.readDraft(a))!.startUtc,
      intent(hour: 11).startUtc,
    );
    expect(await repository.readDraft(b), isNull);
    expect(await repository.readDraft(event), isNull);
  });

  test('simulated conflict retains details and is idempotent', () async {
    final saved = await repository.save(intent(), commandId: 'first');
    await repository.changeStatus(a, saved.id, LocalStatus.simulatedConflict);
    await repository.changeStatus(a, saved.id, LocalStatus.simulatedConflict);
    final conflicted = await repository.get(a, saved.id);
    expect(conflicted.customerLabel, saved.customerLabel);
    expect(conflicted.startUtc, saved.startUtc);
    expect(conflicted.status, LocalStatus.simulatedConflict);
    expect((await repository.history(a, saved.id)).length, 2);
  });

  test(
    'local revision preserves command payload/history and one booking ID',
    () async {
      final saved = await repository.save(intent(), commandId: 'first');
      await repository.changeStatus(a, saved.id, LocalStatus.simulatedConflict);
      final conflicted = await repository.get(a, saved.id);
      final revised = await repository.revise(
        conflicted,
        intent(hour: 11),
        commandId: 'revision',
      );
      final retried = await repository.revise(
        conflicted,
        intent(hour: 11),
        commandId: 'revision',
      );
      expect(revised.id, saved.id);
      expect(retried.id, saved.id);
      expect(revised.revision, 3);
      expect(revised.status, LocalStatus.pending);
      final old = (await repository.database.query(
        'outbox',
        where: 'command_id = ?',
        whereArgs: ['first'],
      )).single;
      expect(
        jsonDecode(old['payload'] as String)['startUtcMs'],
        saved.startUtc.millisecondsSinceEpoch,
      );
      expect((await repository.history(a, saved.id)).length, 3);
    },
  );

  test(
    'failed revision keeps the old schedule; stale revision is rejected',
    () async {
      final first = await repository.save(intent(), commandId: 'first');
      await repository.save(intent(hour: 11), commandId: 'other');
      await expectLater(
        repository.revise(first, intent(hour: 11), commandId: 'failed'),
        throwsA(isA<LocalBookingConflict>()),
      );
      expect((await repository.get(a, first.id)).startUtc, first.startUtc);
      await repository.changeStatus(a, first.id, LocalStatus.simulatedConflict);
      await expectLater(
        repository.revise(first, intent(hour: 13), commandId: 'stale'),
        throwsStateError,
      );
    },
  );

  test('cancellation retains history and removes pending work; cannot be resurrected', () async {
    final first = await repository.save(intent(), commandId: 'first');
    await repository.changeStatus(a, first.id, LocalStatus.cancelled);
    expect(await repository.list(a), isEmpty);
    expect(
      (await repository.list(a, includeCancelled: true)).single.status,
      LocalStatus.cancelled,
    );
    expect((await repository.history(a, first.id)).length, 2);
    await expectLater(
      repository.changeStatus(a, first.id, LocalStatus.simulatedConflict),
      throwsStateError,
    );
    await repository.save(intent(), commandId: 'new-local-item');
    expect((await repository.list(a)).length, 1);
  });

  test(
    'invalid scope, real-looking label and closing-hours overflow cannot save',
    () async {
      await expectLater(
        repository.save(
          intent(customer: 'Actual Person'),
          commandId: 'bad-label',
        ),
        throwsArgumentError,
      );
      await expectLater(
        repository.save(intent(hour: 17, minute: 30), commandId: 'too-late'),
        throwsArgumentError,
      );
      await expectLater(
        repository.save(
          intent(scope: const DemoScope('unknown', 'studio')),
          commandId: 'bad-scope',
        ),
        throwsArgumentError,
      );
      expect(await repository.list(a), isEmpty);
    },
  );

  test('exact money and fixed Ethiopia prototype timezone boundaries', () {
    expect(demoMoney(123456), 'ETB 1234.56');
    final utc = fixtures.localToUtc(DateTime(2030, 1, 1), 0, 0);
    expect(utc, DateTime.utc(2029, 12, 31, 21));
    expect(fixtures.utcToLocal(utc), DateTime.utc(2030, 1, 1));
  });

  test('committed SQLite item survives child-process SIGKILL without graceful close', () async {
    final db = path.join(temp.path, 'killed.sqlite');
    final dart = path.normalize(
      path.join(
        Directory.current.path,
        '../../.tooling/flutter/bin/cache/dart-sdk/bin/dart',
      ),
    );
    final child = await Process.start(dart, [
      'run',
      'tool/persistence_probe.dart',
      db,
      path.normalize(
        path.join(Directory.current.path, '../fixtures/booking-v0.json'),
      ),
    ]);
    final stderr = child.stderr.transform(utf8.decoder).join();
    try {
      final ready = await child.stdout
          .transform(utf8.decoder)
          .transform(const LineSplitter())
          .first
          .timeout(const Duration(seconds: 45));
      expect(ready, 'COMMITTED');
      child.kill(ProcessSignal.sigkill);
      expect(await child.exitCode, isNot(0));
      final reopened = await DemoRepository.open(
        factory: databaseFactoryFfi,
        path: db,
        fixtures: fixtures,
      );
      try {
        expect((await reopened.list(a)).single.status, LocalStatus.pending);
        expect((await reopened.database.query('outbox')).length, 1);
        expect((await reopened.database.query('local_audit')).length, 1);
      } finally {
        await reopened.close();
      }
    } finally {
      child.kill(ProcessSignal.sigkill);
      await child.exitCode;
      final diagnostic = await stderr;
      if (diagnostic.isNotEmpty) {
        // Sanitized synthetic test tool output only.
        // ignore: avoid_print
        print(diagnostic);
      }
    }
  }, timeout: const Timeout(Duration(seconds: 60)));
}

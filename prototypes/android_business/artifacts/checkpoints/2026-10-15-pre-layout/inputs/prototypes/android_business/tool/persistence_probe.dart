import 'dart:io';

import 'package:sqflite_common_ffi/sqflite_ffi.dart';
import 'package:ubi_business_v0/features/bookings/domain/booking.dart';
import 'package:ubi_business_v0/local/demo_repository.dart';

// Parent test kills this process only AFTER SQLite acknowledges the transaction.
Future<void> main(List<String> args) async {
  if (args.length != 2) {
    throw ArgumentError('Expected isolated DB path and fixture path.');
  }
  sqfliteFfiInit();
  final fixtures = DemoFixtures.parse(await File(args[1]).readAsString());
  final repository = await DemoRepository.open(
    factory: databaseFactoryFfi,
    path: args[0],
    fixtures: fixtures,
  );
  await repository.save(
    BookingIntent(
      scope: const DemoScope('demo-org-a', 'studio'),
      serviceId: 'portrait',
      customerLabel: 'Demo Guest 01',
      source: 'Call',
      startUtc: fixtures.localToUtc(fixtures.referenceDate, 9, 0),
    ),
    commandId: 'probe-command',
  );
  stdout.writeln('COMMITTED');
  await stdout.flush();
  await stdin.first; // No graceful DB close; parent issues SIGKILL.
}

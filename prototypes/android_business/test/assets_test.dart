import 'dart:io';

import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:ubi_business_v0/app/demo_assets.dart';
import 'package:ubi_business_v0/app/theme.dart';
import 'package:ubi_business_v0/features/bookings/domain/booking.dart';

void main() {
  TestWidgetsFlutterBinding.ensureInitialized();

  test(
    'startup asset keys load exact shared synthetic fixture and palette',
    () async {
      final booking = await rootBundle.loadString(DemoAssets.bookings);
      final palette = await rootBundle.loadString(DemoAssets.palette);
      expect(booking, await File('../fixtures/booking-v0.json').readAsString());
      expect(palette, await File('../../design/tokens.json').readAsString());
      expect(DemoFixtures.parse(booking).services, isNotEmpty);
      expect(UbiPalette.parse(palette).primary.a, 1);
    },
  );

  test('asset manifest declares safe project-local startup keys', () async {
    final manifest = await AssetManifest.loadFromAssetBundle(rootBundle);
    final keys = manifest.listAssets();
    expect(keys, containsAll([DemoAssets.bookings, DemoAssets.palette]));
    expect(keys.any((key) => key.split('/').contains('..')), isFalse);
  });
}

import 'dart:io';
import 'dart:ui' as ui;

import 'package:flutter/material.dart';
import 'package:flutter/rendering.dart';
import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:sqflite_common_ffi/sqflite_ffi.dart';
import 'package:ubi_business_v0/app/theme.dart';
import 'package:ubi_business_v0/features/bookings/domain/booking.dart';
import 'package:ubi_business_v0/local/demo_repository.dart';
import 'package:ubi_business_v0/main.dart';

Future<void> waitFor(WidgetTester tester, bool Function() ready) async {
  for (var i = 0; i < 100; i++) {
    await tester.runAsync(
      () => Future<void>.delayed(const Duration(milliseconds: 10)),
    );
    await tester.pump(const Duration(milliseconds: 20));
    if (ready()) {
      await tester.pumpAndSettle();
      return;
    }
  }
  fail('UI did not reach the expected durable result.');
}

void main() {
  TestWidgetsFlutterBinding.ensureInitialized();
  late DemoFixtures fixtures;
  late UbiPalette palette;
  sqfliteFfiInit();
  setUpAll(() async {
    for (final family in ['Roboto', 'MaterialIcons']) {
      final loader = FontLoader(family);
      final files = family == 'Roboto'
          ? ['Roboto-Regular.ttf', 'Roboto-Medium.ttf', 'Roboto-Bold.ttf']
          : ['MaterialIcons-Regular.otf'];
      for (final file in files) {
        loader.addFont(
          File(
            '../../.tooling/flutter/bin/cache/artifacts/material_fonts/$file',
          ).readAsBytes().then(ByteData.sublistView),
        );
      }
      await loader.load();
    }
    fixtures = DemoFixtures.parse(
      await File('../fixtures/booking-v0.json').readAsString(),
    );
    palette = UbiPalette.parse(
      await File('../../design/tokens.json').readAsString(),
    );
  });

  Future<DemoRepository> launch(
    WidgetTester tester, {
    TextScaler? scale,
    GlobalKey? capture,
    bool seed = false,
  }) async {
    tester.view.physicalSize = const Size(1080, 2400);
    tester.view.devicePixelRatio = 3;
    addTearDown(tester.view.resetPhysicalSize);
    addTearDown(tester.view.resetDevicePixelRatio);
    final repository = await tester.runAsync(() async {
      final repo = await DemoRepository.open(
        factory: databaseFactoryFfi,
        path: inMemoryDatabasePath,
        fixtures: fixtures,
      );
      if (seed) {
        await repo.save(
          BookingIntent(
            scope: const DemoScope('demo-org-a', 'studio'),
            serviceId: 'portrait',
            customerLabel: 'Demo Guest 01',
            source: 'Call',
            startUtc: fixtures.localToUtc(fixtures.referenceDate, 9, 0),
          ),
          commandId: 'screenshot-first',
        );
      }
      return repo;
    });
    addTearDown(() => repository!.close());
    await tester.pumpWidget(
      UbiPrototype(
        repository: repository!,
        palette: palette,
        textScaler: scale,
        repaintKey: capture,
      ),
    );
    await waitFor(
      tester,
      () => find.byType(LinearProgressIndicator).evaluate().isEmpty,
    );
    return repository;
  }

  Future<void> openForm(WidgetTester tester) async {
    await tester.tap(find.text('Add demo booking'));
    await tester.pump(const Duration(milliseconds: 400));
    await waitFor(
      tester,
      () => find.byType(CircularProgressIndicator).evaluate().isEmpty,
    );
  }

  Future<void> saveDefault(WidgetTester tester) async {
    await tester.tap(find.text('Review details'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Save on this phone'));
    await waitFor(
      tester,
      () => find.text('Your demo schedule').evaluate().isNotEmpty,
    );
  }

  testWidgets(
    'customer-safe demo boundary and working manual save/details path',
    (tester) async {
      final repo = await launch(tester);
      expect(find.textContaining('SYNTHETIC DEMO'), findsOneWidget);
      expect(find.text('Local demo · no server'), findsOneWidget);
      await openForm(tester);
      expect(find.text('Made-up customer'), findsOneWidget);
      await saveDefault(tester);
      expect(
        (await tester.runAsync(
          () => repo.list(const DemoScope('demo-org-a', 'studio')),
        ))!.length,
        1,
      );
      await tester.scrollUntilVisible(find.text('Open saved details'), 180);
      await tester.pumpAndSettle();
      await tester.tap(find.text('Open saved details'));
      await tester.pumpAndSettle();
      await waitFor(
        tester,
        () => find.text('Saved demo details').evaluate().isNotEmpty,
      );
      expect(find.text('Saved on this phone · not confirmed'), findsOneWidget);
      await tester.scrollUntilVisible(
        find.textContaining(
          'No server, customer message or payment is connected',
        ),
        180,
      );
      await tester.pumpAndSettle();
      expect(
        find.textContaining(
          'No server, customer message or payment is connected',
        ),
        findsOneWidget,
      );
      await tester.scrollUntilVisible(
        find.text('• Created on this phone'),
        180,
      );
      await tester.pumpAndSettle();
      expect(find.text('• Created on this phone'), findsOneWidget);
      expect(tester.takeException(), isNull);
    },
  );

  testWidgets('double tap saves one logical local intent', (tester) async {
    final repo = await launch(tester);
    await openForm(tester);
    await tester.tap(find.text('Review details'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Save on this phone'));
    await tester.tap(find.text('Save on this phone'));
    await waitFor(
      tester,
      () => find.text('Your demo schedule').evaluate().isNotEmpty,
    );
    expect(
      (await tester.runAsync(
        () => repo.list(const DemoScope('demo-org-a', 'studio')),
      ))!.length,
      1,
    );
  });

  testWidgets(
    'failed transaction keeps review/input and never shows a saved item',
    (tester) async {
      final repo = await launch(tester);
      await tester.runAsync(
        () => repo.database.execute(
          "CREATE TRIGGER reject_all BEFORE INSERT ON outbox BEGIN SELECT RAISE(ABORT, 'synthetic storage failure'); END",
        ),
      );
      await openForm(tester);
      await tester.tap(find.text('Review details'));
      await tester.pumpAndSettle();
      await tester.tap(find.text('Save on this phone'));
      await waitFor(
        tester,
        () => find
            .textContaining('The save could not finish')
            .evaluate()
            .isNotEmpty,
      );
      expect(find.text('Review before saving'), findsOneWidget);
      expect(
        await tester.runAsync(
          () => repo.list(const DemoScope('demo-org-a', 'studio')),
        ),
        isEmpty,
      );
      expect(find.text('Demo Guest 01'), findsOneWidget);
    },
  );

  testWidgets('simulated conflict is explicit and retains saved details', (
    tester,
  ) async {
    final repo = await launch(tester, seed: true);
    await tester.scrollUntilVisible(find.text('Open saved details'), 180);
    await tester.pumpAndSettle();
    await tester.tap(find.text('Open saved details'));
    await tester.pumpAndSettle();
    await tester.ensureVisible(find.text('Simulate a conflict'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Simulate a conflict'));
    await waitFor(
      tester,
      () => find.text('Choose another time').evaluate().isNotEmpty,
    );
    expect(find.text('Choose another time'), findsOneWidget);
    await tester.drag(find.byType(ListView), const Offset(0, 700));
    await tester.pumpAndSettle();
    expect(find.text('Simulated conflict · review needed'), findsOneWidget);
    final saved = (await tester.runAsync(
      () => repo.list(const DemoScope('demo-org-a', 'studio')),
    ))!.single;
    expect(saved.status, LocalStatus.simulatedConflict);
    expect(saved.customerLabel, 'Demo Guest 01');
  });

  testWidgets('switching demo tenant shows only that tenant schedule', (
    tester,
  ) async {
    await launch(tester, seed: true);
    await tester.tap(find.text('Demo lab'));
    await tester.pumpAndSettle();
    await tester.ensureVisible(find.text('Demo River Studio\nStudio sessions'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Demo River Studio\nStudio sessions'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Today'));
    await waitFor(
      tester,
      () => find.byType(LinearProgressIndicator).evaluate().isEmpty,
    );
    expect(find.text('Demo River Studio'), findsOneWidget);
    expect(find.text('Open saved details'), findsNothing);
    await tester.scrollUntilVisible(find.text('Start with one booking'), 180);
    await tester.pumpAndSettle();
    expect(find.text('Start with one booking'), findsOneWidget);
  });

  testWidgets('home and form remain operable at 200 percent font scale', (
    tester,
  ) async {
    await launch(tester, scale: const TextScaler.linear(2));
    expect(tester.takeException(), isNull);
    await openForm(tester);
    expect(tester.takeException(), isNull);
    await tester.tap(find.text('Review details'));
    await tester.pumpAndSettle();
    expect(tester.takeException(), isNull);
    expect(find.text('Save on this phone'), findsOneWidget);
  });

  testWidgets('primary save action has a labelled minimum touch target', (
    tester,
  ) async {
    await launch(tester);
    final semantics = tester.ensureSemantics();
    try {
      await openForm(tester);
      final button = find.widgetWithText(FilledButton, 'Review details');
      final size = tester.getSize(button);
      expect(size.height, greaterThanOrEqualTo(48));
      expect(size.width, greaterThanOrEqualTo(48));
      expect(
        tester.getSemantics(button),
        matchesSemantics(
          label: 'Review details',
          isButton: true,
          hasEnabledState: true,
          isEnabled: true,
          isFocusable: true,
          hasTapAction: true,
          hasFocusAction: true,
        ),
      );
    } finally {
      semantics.dispose();
    }
  });

  testWidgets(
    'capture actual Flutter home/review rendering and verify palette pairs',
    (tester) async {
      final originalPainting = debugDisableShadows;
      try {
        debugDisableShadows = false;
        final key = GlobalKey();
        await launch(tester, capture: key);
        double ratio(Color a, Color b) {
          final x = a.computeLuminance(), y = b.computeLuminance();
          return ((x > y ? x : y) + .05) / ((x < y ? x : y) + .05);
        }

        expect(ratio(palette.primary, Colors.white), greaterThanOrEqualTo(4.5));
        expect(
          ratio(palette.info, palette.infoSoft),
          greaterThanOrEqualTo(4.5),
        );
        expect(
          ratio(palette.warning, palette.warningSoft),
          greaterThanOrEqualTo(4.5),
        );
        expect(ratio(palette.muted, Colors.white), greaterThanOrEqualTo(4.5));
        Future<void> capture(String name) async {
          final originalShadows = debugDisableShadows;
          try {
            debugDisableShadows = false;
            await tester.pumpAndSettle();
            final boundary =
                key.currentContext!.findRenderObject()!
                    as RenderRepaintBoundary;
            boundary.markNeedsPaint();
            await tester.pump();
            await tester.runAsync(() async {
              final image = await boundary.toImage(pixelRatio: 2);
              final data = await image.toByteData(
                format: ui.ImageByteFormat.png,
              );
              await Directory('artifacts').create(recursive: true);
              await File('artifacts/$name.png')
                  .writeAsBytes(data!.buffer.asUint8List());
              image.dispose();
            });
          } finally {
            debugDisableShadows = originalShadows;
          }
        }

        await capture('home');
        await openForm(tester);
        await capture('add-booking');
        await tester.tap(find.text('Review details'));
        await tester.pumpAndSettle();
        await capture('review');
        expect(tester.takeException(), isNull);
      } finally {
        debugDisableShadows = originalPainting;
      }
    },
  );
}

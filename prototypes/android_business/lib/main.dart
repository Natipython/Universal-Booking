import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:path/path.dart' as path;
import 'package:sqflite/sqflite.dart';

import 'app/demo_assets.dart';
import 'app/theme.dart';
import 'features/bookings/domain/booking.dart';
import 'features/bookings/presentation/business_home.dart';
import 'local/demo_repository.dart';

Future<void> main() async {
  WidgetsFlutterBinding.ensureInitialized();
  try {
    final fixtures = DemoFixtures.parse(
      await rootBundle.loadString(DemoAssets.bookings),
    );
    final palette = UbiPalette.parse(
      await rootBundle.loadString(DemoAssets.palette),
    );
    final repository = await DemoRepository.open(
      factory: databaseFactory,
      path: path.join(await getDatabasesPath(), 'ubi_synthetic_v0.sqlite'),
      fixtures: fixtures,
    );
    runApp(UbiPrototype(repository: repository, palette: palette));
  } catch (_) {
    runApp(
      const MaterialApp(
        home: Scaffold(
          body: SafeArea(
            child: Padding(
              padding: EdgeInsets.all(24),
              child: Center(
                child: Text(
                  'Ubi demo storage could not open. No new work was saved. Close the demo and retry. Do not enter real customer information.',
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }
}

class UbiPrototype extends StatelessWidget {
  const UbiPrototype({
    super.key,
    required this.repository,
    required this.palette,
    this.textScaler,
    this.repaintKey,
  });
  final DemoRepository repository;
  final UbiPalette palette;
  final TextScaler? textScaler;
  final GlobalKey? repaintKey;

  @override
  Widget build(BuildContext context) => MaterialApp(
    title: 'Ubi Business · Demo',
    debugShowCheckedModeBanner: false,
    theme: palette.theme,
    builder: (context, child) {
      final content = MediaQuery(
        data: MediaQuery.of(
          context,
        ).copyWith(textScaler: textScaler ?? MediaQuery.textScalerOf(context)),
        child: child!,
      );
      return repaintKey == null
          ? content
          : RepaintBoundary(key: repaintKey, child: content);
    },
    home: BusinessHome(repository: repository, palette: palette),
  );
}

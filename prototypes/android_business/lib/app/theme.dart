import 'dart:convert';

import 'package:flutter/material.dart';

class UbiPalette {
  UbiPalette.parse(String source) {
    final colors =
        (jsonDecode(source) as Map<String, dynamic>)['color']
            as Map<String, dynamic>;
    Color value(String group, String name) {
      final hex = colors[group][name]['\$value'] as String;
      return Color(int.parse('FF${hex.substring(1)}', radix: 16));
    }

    primary = value('brand', 'primary');
    primarySoft = value('brand', 'primarySoft');
    info = value('semantic', 'info');
    infoSoft = value('semantic', 'infoSoft');
    warning = value('semantic', 'warning');
    warningSoft = value('semantic', 'warningSoft');
    success = value('semantic', 'success');
    danger = value('semantic', 'danger');
    text = value('neutral', 'text');
    muted = value('neutral', 'muted');
    canvas = value('neutral', 'canvas');
    border = value('neutral', 'border');
  }
  late final Color primary,
      primarySoft,
      info,
      infoSoft,
      warning,
      warningSoft,
      success,
      danger,
      text,
      muted,
      canvas,
      border;

  ThemeData get theme => ThemeData(
    useMaterial3: true,
    colorScheme: ColorScheme.fromSeed(seedColor: primary).copyWith(
      primary: primary,
      onPrimary: Colors.white,
      secondary: info,
      onSecondary: Colors.white,
      surface: Colors.white,
      onSurface: text,
      error: danger,
      outline: muted,
    ),
    scaffoldBackgroundColor: canvas,
    appBarTheme: const AppBarTheme(
      backgroundColor: Colors.white,
      surfaceTintColor: Colors.transparent,
      centerTitle: false,
    ),
    textTheme: TextTheme(
      headlineMedium: TextStyle(
        fontSize: 28,
        fontWeight: FontWeight.w700,
        color: text,
      ),
      titleLarge: TextStyle(
        fontSize: 22,
        fontWeight: FontWeight.w700,
        color: text,
      ),
      titleMedium: TextStyle(
        fontSize: 18,
        fontWeight: FontWeight.w700,
        color: text,
      ),
      bodyLarge: TextStyle(fontSize: 16, height: 1.45, color: text),
      bodyMedium: TextStyle(fontSize: 14, height: 1.45, color: muted),
    ),
    filledButtonTheme: FilledButtonThemeData(
      style: FilledButton.styleFrom(
        minimumSize: const Size(48, 56),
        padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 16),
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(14)),
      ),
    ),
    outlinedButtonTheme: OutlinedButtonThemeData(
      style: OutlinedButton.styleFrom(
        minimumSize: const Size(48, 52),
        side: BorderSide(color: muted),
        padding: const EdgeInsets.symmetric(horizontal: 18, vertical: 14),
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(14)),
      ),
    ),
    inputDecorationTheme: InputDecorationTheme(
      filled: true,
      fillColor: Colors.white,
      contentPadding: const EdgeInsets.all(16),
      border: OutlineInputBorder(borderRadius: BorderRadius.circular(14)),
      enabledBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(14),
        borderSide: BorderSide(color: muted),
      ),
      focusedBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(14),
        borderSide: BorderSide(color: primary, width: 2),
      ),
    ),
    navigationBarTheme: NavigationBarThemeData(
      backgroundColor: Colors.white,
      indicatorColor: primarySoft,
    ),
  );
}

# Generated synthetic Android assets

`demo/booking-v0.json` and `demo/tokens.json` are byte-identical generated copies of the canonical shared files in `prototypes/fixtures/` and `design/`. Do not edit these copies or add real customer data.

From the repository root, run `node .tooling/sync-android-assets-v0.mjs` after changing canonical fixtures/tokens. `--check` verifies without repairing or writing a receipt. Both host-validation and APK-build runners synchronize before fingerprinting. Host tests exercise the same root-bundle keys used at startup; the APK runner separately checks the actual archived asset bytes. None of these checks proves execution on Android.

# WB-01/02 synthetic prototype evidence

Status: **PARTIAL — historical host-test evidence; subsequent bounded native result recorded separately; broader device/field checks pending**
Updated: 2026-10-11 (Africa/Addis_Ababa)
Related: specifications 21, 23, 35, 61, 62, 71–73, 85; DEC-022; EV-03/05

## Outcome and exact boundary

Subsequent 2026-10-15 checkpoint: [bounded native checks on an authorized Samsung SM-G975W/Android 12](WB-02-physical-android-smoke-evidence.md) now establish actual native SQLite execution and one separately approved force-stop/draft-recovery sample. This supersedes the dated absence of all Android execution below, not the missing minimum-device, performance, accessibility or full-scenario evidence. The October 11 host receipt itself was not rerun or relabelled.

Later 2026-10-15 layout work has **41/41 current host tests** and a separate candidate; see [the layout follow-up](WB-02-home-action-layout-evidence.md). This document retains the October 11 result and now links frozen copies so rebuilding cannot substitute new evidence for old claims.

The isolated [Ubi Business V0 project](../../../prototypes/android_business/README.md) implements a synthetic manual-booking/local-state feasibility slice. The October 11 machine run passed **28/28 tests**, with analyzer exit 0, no skipped tests and no covered-source change during the run. This is Linux-host Flutter/domain/SQLite FFI/startup-asset evidence, not a physical Android run or a working production service. A separate [debug APK receipt](WB-02-synthetic-android-apk-evidence.md) establishes compilation and static inspection, not native execution.

Historical machine receipt: [validation-receipt.json](../../../prototypes/android_business/artifacts/checkpoints/2026-10-15-pre-layout/inputs/prototypes/android_business/artifacts/validation-receipt.json). Rechecked at `2026-10-11T14:36:10.457Z` / `2026-10-11 17:36:10.457 EAT`. Covered-source SHA-256 is `c35b013fd3607d5c6c13bff44a71f7c4fb7dee262507987126c4075385833348`; all 28 tests pass. The two added tests cover the corrected startup bundle keys; canonical fixture/theme values and booking logic are unchanged. The previous 26-test receipt is preserved under `artifacts/host-before-asset-fix/`. This host-runner digest does not cover every Android wrapper/build-tool input. Re-running the runner regenerates the current receipt, not this archived copy; this prose describes that dated run, not a permanent guarantee. No independent human reviewer sign-off is recorded.

Scope isolation matters: this code lives under `prototypes/`, not production `apps/` or `services/`. There are no real customers, credentials, identity verification, server sync, external notices, payments or native calls/messages. Local storage is unencrypted and is restricted to invented fixture data. No public domain, provider/company account, production database, commit or push was changed.

## What was implemented

- Two invented organizations, three operations (two in one organization), four services, three made-up customer choices and a deliberately future Gregorian fixture date.
- Today/day-selection list, dummy manual entry, explicit service/date/clock/customer/source/price review and local saved details.
- A transaction that saves a tentative item, immutable local command payload and local history together. “Saved on this phone” never means server-confirmed.
- Persisted, operation-scoped draft; retry/double-tap deduplication; local same-resource overlap checks and retained input on failure.
- Explicit conflict simulation, reviewed revision with original payload/history retained, and local cancellation without deleting history.
- Demo scope switching and scoped repository reads/writes. These are fixture isolation checks, not real authentication, staff permissions or PostgreSQL RLS.
- Founder palette: white/purple surfaces, accessible dark-blue information text with sky-blue surfaces, amber review state. No green confirmation appears because there is no authoritative confirmation.

The “Calendar” tab is a simple date-selectable list in V0, not a finished resource/month/drag calendar. Source choices such as WhatsApp/Telegram record a dummy origin label only; they do not connect to inboxes. In-app messaging and voice calls remain required future WB-13/WB-16 work, not working features.

## Executed checks

| Boundary | Result | What the result does not establish |
|---|---|---|
| 18 repository/domain tests | PASS: transaction rollback, stable-command retry, changed-input rejection, local overlap/adjacency/concurrency, scoped composite keys/reads/writes/drafts/history, revision/cancellation, exact integer money and fixture time | Multi-device/server capacity truth, real identity/authorization, provider correctness or complete SCN-002–006/013/014 |
| Child-process persistence, included in the 18 | PASS: a committed SQLite item/outbox/audit survives SIGKILL and reopening without graceful close | Android OS process death, filesystem durability across power loss, background restrictions, plugin behavior or low storage |
| 8 widget tests | PASS: manual save/details, double tap, failed-storage review retention, explicit conflict, tenant switch, 200% text scale, one labelled minimum touch target, selected contrast pairs and captures | All screens/controls, TalkBack, Amharic/Ethiopian calendar, participant comprehension, device startup/frames/battery or network |
| 2 startup-asset tests | PASS: the exact root-bundle keys used by main load canonical synthetic fixture/palette contents; manifest uses project-local keys | Native launch or SQLite plugin execution; actual APK contents are separately inspected |
| Flutter analyzer | PASS, no issues | Runtime completeness, security certification or production readiness |
| Documentation-checker tests | 9/9 PASS, including external prototype Markdown links and missing-link rejection | Application correctness or gate approval |

These exercise **subsets**, not completion, of SCN-002/003/005/006/013/018. The 28 full scenario definitions remain NOT RUN. Actual roles, approved server operations, device/network conditions and research observations are absent.

## UI artifacts

Actual Flutter widget-test renders at a 360 × 800 logical viewport, captured at 2× with SDK Roboto/MaterialIcons fonts and normal shadow rendering:

- [Home](../../../prototypes/android_business/artifacts/checkpoints/2026-10-15-pre-layout/inputs/prototypes/android_business/artifacts/home.png)
- [Manual entry](../../../prototypes/android_business/artifacts/checkpoints/2026-10-15-pre-layout/inputs/prototypes/android_business/artifacts/add-booking.png)
- [Review before local save](../../../prototypes/android_business/artifacts/checkpoints/2026-10-15-pre-layout/inputs/prototypes/android_business/artifacts/review.png)

These are not physical Android screenshots or golden-regression baselines. Visual inspection found readable main text and an explicit demo/local boundary at the standard scale. The automated large-text test checks the limited home/form/review path for framework layout errors; it does not certify every label or accessible interaction. Home currently has substantial introductory/demo content before the list, so field work must check whether the first useful schedule/action is sufficiently immediate. Production copy/layout is not frozen.

## Supported device decision and observed tooling

The founder specified **Android 11 or later and at least 4 GB RAM** on 2026-10-04 (DEC-022). It supersedes the provisional Android 10/2 GB business D1 baseline, not the separate guest-web support policy. The prototype Gradle configuration sets `minSdk = 30`; RAM is a physical-device support/test baseline, not a free-memory threshold or startup denial. [Android's SDK documentation](https://developer.android.com/guide/topics/manifest/uses-sdk-element) supports the OS/API mapping and minimum-SDK distinction. Compile/target SDK and store compliance require separate release checks.

Observed installed SDK: Flutter 3.47.1, Dart 3.13.1. Dependencies and checksums are in the prototype lockfile. Host tests use `sqflite_common_ffi`; the Android app path uses `sqflite`, which has not been exercised on Android. This records an installed feasibility toolchain, not approval of a final production stack or a latest-version claim.

The earlier `flutter doctor -v` found no configured Android SDK or JDK. Subsequently, explicit approval of both SDK agreement revisions allowed isolated JDK/Android tooling and required stable packages to be installed and verified; the [setup checkpoint](WB-02-android-build-preflight.md) records that preparation. The corrected [debug APK](WB-02-synthetic-android-apk-evidence.md) now compiles and passes static manifest/signature/adapter/asset checks, but no Android phone or emulator has been installed/tested. Generated debug signing is not production signing. The phone answer supplied a minimum target, not a concrete model or installation approval. Go was initially absent; the subsequent [WB-03 experiment](WB-03-synthetic-integrity-evidence.md) installed a local Go runtime and ran a separate newly created disposable PostgreSQL cluster. That adds server-component evidence without implementing Android/server synchronization.

## Reproduction

From the repository root:

```sh
node .tooling/validate-v0.mjs
node --test .tooling/verify-blueprint.test.mjs
node .tooling/verify-blueprint.mjs --report
```

The V0 runner invokes the installed local Flutter SDK, analyzer and machine-format tests, captures toolchain/test details, fingerprints the covered inputs, and regenerates receipt/logs/screenshots. The isolated [task protocol](../../../prototypes/validation/WB-01-task-protocol.md) describes field attempts without treating unfamiliarity as a person's defect. Do not use real personal information in this prototype.

## Remaining work, in order

1. SDK/JDK setup and the inspected synthetic debug APK are complete within the [bounded build receipt](WB-02-synthetic-android-apk-evidence.md). Identify an authorized test device and obtain installation direction before native execution. Do not distribute a debug-signed build as a production release.
2. Obtain a concrete physical Android 11/4 GB device and newer supported device/version coverage. Check the native SQLite path, process death, low storage/background behavior, startup/frames/network/data/battery and accessibility. No benchmark result has yet been measured.
3. WB-04 now has [local customer-link component evidence](WB-04-synthetic-customer-link-evidence.md). Complete representative guest compatibility/network/localization/comprehension checks and review persistent recovery/identity; keep guest web support separate from the merchant Android floor.
4. Review the recorded WB-03 component result and resolve/disposition its spec-linter limitation; real identity, client/server integration, migration/restore and broader allocation rules remain unverified. Do not replace database truth with local overlap checks.
5. Complete appropriate consenting research, localization and review; use evidence to review WB-05 contracts/frameworks. EV-03/04/05 remain PARTIAL, other launch evidence remains missing, and G7/G8/G9 stay OPEN.

This progress does not confirm Addis Ababa/photo studios, pilot size/duration, entity/payment roles, providers, counsel approval, funded team or a launch date. Those proposals retain their decision-register status.

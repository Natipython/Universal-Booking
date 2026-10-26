# WB-02 first physical Android smoke-test evidence

Status: **PARTIAL — bounded native checks pass on one Android 12 phone; performance/accessibility/release evidence pending**  
Tested: 2026-10-15 (Africa/Addis_Ababa)  
Related: DEC-022; partial EV-03/05; WB-02 and G7/G8/G9 remain open

## Outcome and safety boundary

The founder connected a Samsung **SM-G975W**, authorized USB debugging and directed careful Ubi-only testing. The founder separately answered **“Yes, force-stop Ubi only”** for one lifecycle check. The verified synthetic debug APK was freshly installed, opened and exercised through Android's native SQLite adapter. No existing Ubi package was overwritten. No contacts, messages, photos, other app stores or real customer records were accessed. No rooting, bootloader change, app-data clearing, uninstall, network/security/display/accessibility setting change or storage-filling test was performed.

The APK actually installed on the phone has SHA-256 `d42e02e366b255213b54cb1866b9127f8a4171f8d09f3c48e19dcd95d47b07f4`, matching the [October 11 inspected build](WB-02-synthetic-android-apk-evidence.md). Its covered application/build sources were unchanged at that phone run. It remains a debug, synthetic-only application with unencrypted local storage, not a production release or Android/server integration.

Subsequent approved continuation on the same date fixed the action layout **in source/host tests only** and produced a separate 0.1.1 (2) candidate; see [layout evidence](WB-02-home-action-layout-evidence.md). The phone was not reconnected/updated. This original 0.1.0 (1) APK, covered inputs and receipts are [preserved separately](../../../prototypes/android_business/artifacts/checkpoints/2026-10-15-pre-layout/checkpoint.json). Offline rechecks now use those frozen inputs and explicitly report that the working source differs; they are not new device runs.

Machine evidence: [native receipt](../../../prototypes/android_business/artifacts/android-native/2026-10-15-sm-g975w/validation-receipt.json), [device metadata](../../../prototypes/android_business/artifacts/android-native/2026-10-15-sm-g975w/device-info.json), [bounded authorization](../../../prototypes/android_business/artifacts/android-native/2026-10-15-sm-g975w/session-authorization.json) and [session actions](../../../prototypes/android_business/artifacts/android-native/2026-10-15-sm-g975w/session-events.json). The offline verifier passed **11 evidence cross-checks**; these are assertions over one manually driven session, not 11 independent devices or full acceptance scenarios.

## Actual device and build

| Item | Observed value |
|---|---|
| Phone | Samsung SM-G975W; serial/IMEI not retained |
| Android | 12, API 31 |
| CPU compatibility | `arm64-v8a,armeabi-v7a,armeabi`; APK includes compatible ABIs |
| Memory | `/proc/meminfo` MemTotal 7,640,540 KiB, about 7.3 GiB OS-usable RAM; not an independently verified physical-RAM specification |
| Available data storage after installation | 41,884,180 KiB, about 40 GiB |
| Display | 1440 × 3040; reported physical density 420, existing override 560; unchanged |
| Installed package | `dev.ubi.validation.ubi_business_v0`, version 0.1.0 (1), min/target API 30/36 |
| Artifact | Existing 175,588,986-byte universal debug APK; no rebuild or UI/domain change |

This is an above-minimum Android 12 sample, **not verification of the Android 11/4 GB minimum configuration** or newer Android coverage. The red “60” indicator and edge overlay visible on the phone were already present; they were not configured or treated as an app frame-rate measurement. Debug launch diagnostics in the action record are not accepted startup benchmarks.

## Executed native paths

| Path | Result and evidence | Important limit |
|---|---|---|
| Install/launch | Fresh installation succeeded; native store created; first screen showed the correct fixture and explicit demo boundary | One device/debug build; no release compatibility certification |
| Review | Correct service, date, 24-hour time, customer, source and invented ETB amount shown; no booking written before save | English/Gregorian/international clock only |
| Save/repeated taps | Two consecutive device save taps produced exactly one booking, one outbox command and one history event | Inter-tap interval not measured; not a stress/reconnect/idempotent-server test |
| Conflict simulation/revision | Original details retained; reviewed change from 09:00 to 10:00 kept the booking ID and immutable original command payload; history appended | Deliberately local simulated conflict, not a server response |
| Draft/process restart | Guest 02/10:00 unfinished draft plus the existing booking survived the **one authorized force-stop**; before/after SQLite SHA-256 identical; reopening Add restored the draft | Force-stop is not low-memory kill, reboot, power loss or arbitrary background restriction |
| Actual local overlap | Saving the recovered 10:00 draft against the existing 10:00 item was rejected; visible input and entire database snapshot unchanged | Only local exclusive fixture resource rules |
| Cancellation safety/history | “Keep item” changed nothing. Explicit cancellation retained the original record and four history events, cancelled queued commands and freed the time | Synthetic local cancellation; no money, notice or remote booking |
| Operation/business draft separation | Events operation showed its own 30-minute/ETB 500 service and default Guest 01; other business showed 45-minute/ETB 1200 and default Guest 01; original studio restored Guest 02/10:00; database unchanged while switching | Fixture UI/read isolation only, not real staff permissions, authentication or every cross-scope write case |
| Reuse after cancellation | Recovered draft then saved successfully at the freed time, retaining original cancellation/history; final state: two stored records (one pending, one cancelled), three commands, five history events, no draft | Not an authoritative confirmed booking |
| Limited logs | App-PID logs sampled around launch/end contained no matched fatal, Flutter-error or SQLite-exception signature; local debug-service addresses redacted | Not complete lifetime logs, zero crashes guarantee, security audit or crash-free rate |

SQLite snapshots were copied only from the newly installed demo's known private database while it was idle and cross-checked read-only. The snapshots and JSON summaries agree and pass `PRAGMA quick_check`. They do not establish filesystem/power-loss durability. Screenshots were taken only after confirming Ubi was focused before and after capture; status/navigation bars can include system indicators. No notification shade or unrelated app was opened.

Selected actual phone captures: [review](../../../prototypes/android_business/artifacts/android-native/2026-10-15-sm-g975w/03-review-before-save-screen.png), [recovered draft](../../../prototypes/android_business/artifacts/android-native/2026-10-15-sm-g975w/12-recovered-draft-screen.png), [overlap rejection](../../../prototypes/android_business/artifacts/android-native/2026-10-15-sm-g975w/13-overlap-rejected-screen.png), [retained cancellation history](../../../prototypes/android_business/artifacts/android-native/2026-10-15-sm-g975w/15-cancelled-details-screen.png) and [final schedule after scrolling](../../../prototypes/android_business/artifacts/android-native/2026-10-15-sm-g975w/20-final-schedule-screen.png). Visual observations were made by the agent, not an independent human reviewer or automatic image-semantic validator.

## Finding in the tested 0.1.0 build — NATIVE-UX-01

At the initial Home scroll position, the floating **Add demo booking** button overlaps part of a booking card and its pending-status text. [The initial saved screen](../../../prototypes/android_business/artifacts/android-native/2026-10-15-sm-g975w/04-after-save-screen.png) shows this. Scrolling makes the whole card accessible, as the final schedule capture demonstrates. Introductory demo content also pushes useful schedule content down.

**Recorded, not fixed during this phone session.** That session authorized testing, not a UI redesign. The later [local follow-up](WB-02-home-action-layout-evidence.md) fixes the overlap in source and small/large-text host viewports; actual phone verification of the fix is still pending. The phone's installed APK remains unchanged. The historical native receipt retains `OPEN_NOT_FIXED` for the old tested artifact, not the later source.

## Remaining checks and handoff

Airplane-mode transitions, weak network, TalkBack, actual OS text/display scaling, low storage, background restrictions, low-memory process death, reboot/power loss, battery and release/profile performance were **not run**. These would need appropriate devices/builds and any relevant setting/lifecycle authorization. Do not claim Telegram-like smoothness, Android 11/4 GB coverage, Amharic/Ethiopian-calendar comprehension or full acceptance from this sample.

The isolated testing ADB server was stopped after the session. Ubi and its synthetic test records remain installed; nothing was erased. The founder was asked to turn USB debugging off. No automatic reconnect or further phone interaction is authorized by this receipt.

The offline evidence check can be repeated without touching the phone:

```sh
node .tooling/verify-android-native-evidence-v0.mjs
```

`.tooling/android-device-v0.mjs` is a session-specific Ubi-only interaction helper, **not an unattended test runner or standing device permission**. The saved session is now CLOSED and the helper rejects commands before ADB access. Do not reopen/reuse it without renewed session authorization and target checks. Its recorded one-time force-stop cannot be repeated under this session authorization.

Next: obtain renewed device-update/recheck authorization for the separate layout candidate; then separately prepare and authorize profile/release measurement and the remaining [native checklist](WB-02-android-native-test-checklist.md), Android 11/4 GB/newer-device coverage and consenting comprehension/localization research. All 28 full application scenarios and release gates remain open. Desktop, mandatory messaging/calling and production identity/synchronization remain in the [ordered backlog](../specifications/85-release-roadmap-and-implementation-backlog.md).

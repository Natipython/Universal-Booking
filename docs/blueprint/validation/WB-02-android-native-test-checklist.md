# WB-02 native Android follow-up checklist

Status: **PARTIAL — first Android 12 physical-phone sample executed; baseline/performance/accessibility checks remain**
Prepared: 2026-10-11 (Africa/Addis_Ababa)  
Updated: 2026-10-15; [actual device receipt and limits](WB-02-physical-android-smoke-evidence.md)  
Scope: Synthetic business V0 only; partial EV-03/05, not complete acceptance scenarios

The PASS/PARTIAL native rows below describe the **0.1.0 (1)** phone session, not the later build. [The layout follow-up](WB-02-home-action-layout-evidence.md) fixes NATIVE-UX-01 in source/host tests and prepares 0.1.1 (2); no phone update/recheck has been authorized or performed. The old APK/evidence is preserved separately.

## Entry conditions

1. Identify an authorized test phone: model, actual Android version/API, physical RAM and free storage. The founder's support floor remains **Android 11+/4 GB physical RAM**, not a free-RAM threshold. Include a newer supported Android version separately.
2. Explicitly authorize installation on that identified test device. Local SDK/build approval does **not** authorize phone installation, developer-mode changes, device erasure or access to personal data.
3. Record the APK path, SHA-256, package/version, minimum/target API, ABI and debug signing from the actual build receipt. A debug APK is not production-signed or suitable for speed acceptance.
4. Use invented fixture records only. Storage is unencrypted; there is no real authentication, authoritative server, payment, inbox, message or call integration. Do not import contacts, real bookings or payment screenshots.
5. Preserve other applications and existing device data. Record a result and evidence for each attempt; do not mark a missing/blocked attempt as passed.

## Ordered native smoke checks

| Check | Expected observation | Current result |
|---|---|---|
| Updated Home action layout | Add stays clear of booking content at initial/mid/end scroll, with native insets/text scale; prior synthetic records preserved when updating | HOST PASS only — 12 viewport/state checks and one semantics/navigation check; new-build device update/recheck NOT RUN |
| Launch and scope | Clearly labelled synthetic demo; expected fixture organization/operation; no permission request for contacts, microphone, camera or location | PASS bounded sample — SM-G975W/Android 12 |
| Review before save | Select a dummy service/customer/time; exact operation, service, time, duration and price are reviewable; back/edit preserves input | PASS bounded reviewed create/revision paths |
| Native SQLite save | Save through the Android `sqflite` adapter; Today/details, local tentative label, command/outbox and local history agree | PASS bounded sample; actual private synthetic database inspected |
| Repeat tap/retry | Repeated save does not create a duplicate local booking or command; changed input is not silently reused under the old command | PARTIAL — two consecutive taps made one item; retry stress/changed-command cases not exercised natively |
| Draft reopening | Leave an unfinished form, close/reopen and inspect the scoped recovered draft; no cross-operation customer/form leakage | PASS bounded sample; draft recovered after authorized force-stop and scope return |
| Process recreation | Background the app and use a specifically authorized process-termination method; reopen and verify committed local item, pending command, history and recoverable draft | PARTIAL — one explicitly authorized force-stop passed; not low-memory/background/power-loss coverage |
| Conflict and revision | Simulated conflict is explicit and retains input; reviewed revision preserves previous payload/history; no green server-confirmed claim | PASS bounded local simulation/revision and actual local-overlap rejection |
| Local cancellation | Cancellation is clearly local/tentative where applicable; schedule/details/history remain consistent and no audit history is deleted | PASS bounded Keep/Cancel/history and freed-time reuse |
| Scope switching | Two fixture organizations and all three operations remain isolated across list/details/draft/history/revision/cancellation paths | PARTIAL — correct scoped forms and draft recovery observed; all cross-scope mutation paths not run |
| Airplane mode | Already-local schedule and manual intent remain usable; no fake synchronization, delivered-message or payment success | NOT RUN |
| Android accessibility | TalkBack focus/labels/actions, Android font/display scaling, keyboard and scrolling remain usable on actual screens | NOT RUN |
| Low storage/background | On a designated safe test device or controlled environment, a storage/background failure leaves truthful status and recoverable input rather than false save success | NOT RUN |

A force-stop, process kill, activity recreation and power loss are different conditions. Record the exact method; host SIGKILL tests or a force-stop alone do not establish every Android lifecycle/durability property. Do not fill a personal phone's storage or clear application data to simulate a failure without separate explicit approval.

## Performance and user comprehension — separate evidence

After native correctness, use a separately inspected profile/release-mode build for startup and frame measurements. Record model, OS, physical/free RAM, storage, build mode/hash, cold/warm definitions, repetitions and distribution of results; do not report one run as a reliable target or debug timings as launch performance. Assess battery/background and weak-network behavior only against features that actually exist; this local prototype has no server synchronization to benchmark.

Conduct consenting first-time-user tasks and language/calendar checks separately. Flutter host text-scale/layout tests do not establish Android TalkBack, Amharic copy, Ethiopian-calendar/clock understanding or an unfamiliar user's ability to complete booking unaided.

## Evidence and completion boundary

Keep a dated device receipt with each check's result, failure/recovery steps and minimal redacted screenshots/logs. Review failures before adding features. APK compilation/static inspection, this checklist and host tests do not close WB-02, WB-05, all 28 full scenarios or G7/G8/G9. See [inspected APK evidence](WB-02-synthetic-android-apk-evidence.md), [build setup](WB-02-android-build-preflight.md), [host evidence](WB-01-02-synthetic-prototype-evidence.md) and [ordered backlog](../specifications/85-release-roadmap-and-implementation-backlog.md).

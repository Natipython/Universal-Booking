# WB-02 Android build preflight and next actions

Status: **Android toolchain verified; corrected synthetic debug APK built and statically inspected**
Checked: 2026-10-11 (Africa/Addis_Ababa)

## Verified checkpoint

Subsequent 2026-10-15 checkpoint: [bounded checks on an authorized Samsung SM-G975W/Android 12](WB-02-physical-android-smoke-evidence.md) now add actual native SQLite and one explicitly approved force-stop/draft-recovery sample for the unchanged APK. This setup document's October 11 statements remain historical; minimum-device, broader native/device and performance evidence are still pending.

Later that date, [the Home action follow-up](WB-02-home-action-layout-evidence.md) fixes the overlap locally and builds a separate 0.1.1 candidate, without reopening the phone session. This setup history is not a new-build phone result; the original APK/receipts are preserved separately.

WB-04's local customer-link component receipt exists; do not repeat that implementation or infer that its browser checks prove Android behavior. The next dependency is a native build of the existing synthetic Android business prototype, not a new production application or framework rewrite.

The initial read-only preflight found no `java`, `javac`, `adb`, `sdkmanager` or `emulator` on PATH. `JAVA_HOME`, `ANDROID_HOME` and `ANDROID_SDK_ROOT` were not configured; the common SDK directories and proposed repository-local JDK/SDK directories were absent. Those observations meant no usable toolchain was detected, not proof that no archive existed anywhere on the computer. No installation or agreement acceptance was performed at that initial checkpoint. The subsequently authorized setup and [successful debug APK inspection](WB-02-synthetic-android-apk-evidence.md) supersede those tooling/build gaps, not the missing native/device evidence.

The existing Gradle wrapper JAR/script and Flutter SDK path are present. The installed Flutter source and project files give these build inputs:

| Input | Existing value / setup target | Evidence |
|---|---|---|
| Flutter / Dart | 3.47.1 / 3.13.1 in the saved host receipt | `prototypes/android_business/artifacts/validation-receipt.json`; not a latest-version claim |
| Minimum phone API | 30 — Android 11 | `android/app/build.gradle.kts`; founder floor remains Android 11+/4 GB physical RAM |
| Compile / target API | 36 / 36 | Installed Flutter `FlutterExtension.kt`; inherited by the app |
| Android Gradle Plugin | 9.1.0 | `android/settings.gradle.kts` |
| Gradle wrapper | 9.3.1 | `android/gradle/wrapper/gradle-wrapper.properties` |
| Kotlin plugin | 2.4.0 | `android/settings.gradle.kts`; bounded debug build passed, broader compatibility untested |
| Java compilation target | 17 | App compile/Kotlin settings; installed Temurin 17.0.20.1+1 runs successfully |
| SDK Build Tools | 36.0.0 installed | Installed `source.properties`, SDK-manager list and setup receipt |
| NDK | 28.2.13676358 | Installed Flutter extension, inherited by the app |

The [official AGP 9.1-family notes](https://developer.android.com/build/releases/agp-9-1-0-release-notes), checked 2026-10-11, list Gradle 9.3.1, JDK 17, Build Tools 36.0.0 and NDK 28.2.13676358. The current page is headed 9.1.1 while this project pins 9.1.0; that family-level guidance is not proof that the complete Flutter/Kotlin/plugin combination builds. Keep the saved versions initially and test them; do not silently upgrade to whichever version is newest. The downloaded command-line tool revision's Java compatibility must also be checked during authorized setup.

Minimum phone API is **not** the compile/target API. Building with SDK 36 does not change the decided Android 11 minimum or turn RAM into a startup exclusion check.

## Resource and isolation plan

The development host reports about 7.1 GiB physical memory and 329 GiB free on the project drive. These are host observations, not an Android device benchmark. The existing template permits an 8 GiB Gradle heap plus 4 GiB metaspace; these are ceilings rather than immediate allocations. An isolated `.tooling/gradle-home/gradle.properties` now sets a 2 GiB heap, 768 MiB metaspace, two workers, no parallel project execution and in-process Kotlin compilation for the first build. These are conservative host limits, not a benchmarked optimum or Android performance tuning.

Use repository-local tool/cache directories after authorization, rather than system-wide installs or global shell/Flutter configuration changes:

- `.tooling/jdk-17/` — a verified official OpenJDK distribution; pin the downloaded revision/checksum.
- `.tooling/android-sdk/` — stable command-line tools and only required platform/build/NDK packages.
- `.tooling/gradle-home/` — isolated Gradle downloads/cache.
- `.tooling/android-downloads/` — verified archives on the project drive; do not unpack large SDK/NDK archives in the smaller RAM-backed `/tmp`.

The [official command-line tools documentation](https://developer.android.com/tools) supports an SDK-manager-only setup without installing the full Android Studio IDE. No emulator/system image is needed for the initial APK build. An emulator may later help native smoke testing but cannot replace the agreed physical Android 11+/4 GB performance checks. Do not start one or install onto a personal phone without the relevant direction.

## Agreement boundary

The [Android SDK download page](https://developer.android.com/studio) requires agreement before downloading the command-line tools, not merely before running `sdkmanager --licenses`. Its terms describe the SDK agreement as binding and require authority when accepting for an organization. [The sdkmanager documentation](https://developer.android.com/tools/sdkmanager#accept-licenses) separately describes package-license prompts. This is a setup/authority finding, not legal advice about Ubi's launch.

**Explicit human acceptance/authorization received on 2026-10-11.** In direct response to the question asking whether the founder accepts the linked Android SDK agreement and authorizes local build-tool installation and a synthetic APK build, the founder answered: “YES I APPROVE JUST BE VERY CAREFUL AND CONSCIOUS!” This authorizes that bounded setup/build, not publication, personal-device installation or additional/different agreements. Do not pipe blanket yes responses into license prompts, fabricate accepted-license hashes or accept SDK telemetry. No Google/Play account or purchase is needed.

### Subsequent verified setup and agreement mismatch

- Downloaded official Eclipse Temurin **17.0.20.1+1** and Google Android command-line tools **22.0**, verifying their published SHA-256 values before extraction. Their local `java -version` and `sdkmanager --version` commands succeed. The latter prints a deprecation warning recommending Android CLI; it still reports 22.0. No new framework or CLI migration is inferred from that warning.
- All tools and archives are under `.tooling/`; no global shell/Flutter configuration or system package installation was changed. Record: [setup receipt](../../../prototypes/android_business/artifacts/android-build/setup-receipt.json), [source/checksum pins](../../../.tooling/android-toolchain-pins.json).
- The template wrapper JAR identifies itself as Gradle 2.10. Replaced that launcher with the official checksum-verified **9.3.1 wrapper JAR**, preserving the old binary in `.tooling/android-downloads/gradle-wrapper-template-2.10.jar`. Added the official Gradle 9.3.1 distribution SHA-256 and download timeout/URL validation to the existing wrapper properties. The wrapper subsequently downloaded/verified the distribution and `gradlew --version` reports **9.3.1 on Temurin 17.0.20.1+1**. No AGP/Kotlin/Flutter version was changed; a version command is not APK compilation.
- Google's download page presents SDK terms dated **April 28, 2026**, covered by the founder's explicit approval. However, the official [stable package metadata](https://dl.google.com/android/repository/repository2-3.xml) references `android-sdk-license` dated **January 16, 2019** for all four requested packages. Its extracted [review text](../../../.tooling/android-downloads/android-sdk-license-2019.txt) has SHA-256 `1f8729233617b193fd619213792ae16a41b95d2bbbf525dfe66998252ba68b16`.
- **The founder explicitly authorized the January 2019 stable-package agreement on 2026-10-11:** “Yes, accept that stable SDK agreement and continue the local test build.” The actual package-manager prompt matched the reviewed text after whitespace normalization; that one prompt was accepted interactively. No acceptance file was fabricated and no blanket `--licenses`/yes command was used. No preview license, telemetry, account, publication or device installation is authorized.
- The four packages installed successfully: platform-tools **37.0.1**, Android 36 platform **revision 2**, Build Tools **36.0.0**, NDK **28.2.13676358/r28c**. The setup runner verifies installed properties against the required API/build/NDK versions and records the SDK manager's actual license-file hash. It reports `PASS_REQUIRED_TOOL_SETUP`, not Android runtime success. The optional independent NDK-archive hash attempt found that SDK manager had already removed its temporary archive; [that check](../../../prototypes/android_business/artifacts/android-build/ndk-archive-checksum.json) is explicitly **not captured**, with no independent SHA-256 claim.
- `node .tooling/build-android-v0.mjs` passed locked dependency resolution, corrected synthetic debug compilation and actual manifest/permissions/native plugin/signature/hash/bundled-asset inspection at `2026-10-11T14:38:26.712Z`. [APK evidence](WB-02-synthetic-android-apk-evidence.md) records exact size/hash, preserved failed/superseded attempts and the startup-asset correction. The supplied `apkanalyzer` shell launcher mishandles spaces in the SDK path; the runner invokes its unchanged official Java entry point with an argument array. SDK binaries were not patched. The runner never accepts agreements or installs onto devices itself.

## Next steps, in order

1. Download-page SDK agreement and the January 2019 stable-package revision explicitly authorized on 2026-10-11. The latter was accepted through its actual package-manager prompt. No company registration or Play Console account is needed for this local test build.
2. JDK and command-line tools installed and checksum/version verified. Continue using per-command environment settings, not global configuration changes.
3. Required stable packages installed and their versions verified. No emulator/system image or CMake package was installed. Any genuinely required new package/agreement must still be handled explicitly.
4. Java/Gradle/Flutter/SDK diagnostics and locked Dart resolution passed. Preserve the recorded warnings/failures and pins; conservative host limits are not a benchmarked optimum.
5. **Synthetic debug APK built and statically inspected.** The separate receipt records commands, covered inputs, artifact SHA-256, actual merged manifest/package/API/permissions/signing identity, adapter code and exact bundled fixture/theme bytes. The prototype's release configuration currently uses a debug key; a future release-mode build would still not be a production-signed release. Do not create/upload production keys or distribute the build publicly.
6. With an identified authorized test device, exercise the native sqflite adapter, restart/process-death recovery, tentative-state/review/conflict/revision/cancellation paths and scope isolation. Keep real personal information out of this unencrypted fixture prototype.
7. Measure release/profile behavior separately on representative Android 11+/4 GB and newer devices: startup, frames, low storage/background, weak/offline network, battery, TalkBack, text scaling and comprehension/localization. A debug APK or emulator result alone cannot establish these outcomes.

## What has not changed

Booking/domain/database behavior, canonical fixture/theme values and supported phone floor were not changed. Alongside isolated tooling/wrapper integrity changes, startup now loads generated project-local asset keys because the former parent-relative files were absent from the actual APK. The host runner was rerun at `2026-10-11T14:36:10.457Z`: **28/28 pass**, analyzer exit 0, no covered-source changes during that run; two tests were added for startup assets. The separate successful APK receipt has broader/different covered build inputs. WB-03's known linter disagreement remains open; WB-04 remains a local synthetic component. All 28 full acceptance scenarios and G7/G8/G9 remain open. Desktop, in-app messaging and calling remain in their existing dependency order.

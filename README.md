# Ubi

Ubi is booking infrastructure for East African service businesses that currently manage customers through calls, messaging apps, notebooks, and payment screenshots.

The first business application is **Ubi Business for Android**, followed by a desktop business workspace. A lightweight customer web experience supports shared booking links without requiring installation. Ethiopia is the decided launch country; Addis Ababa and photography studios remain proposed first-city/vertical choices.

## Product promise

> Manage every booking from your phone. Customers tap, choose, and book.

Ubi is designed for people with limited experience using booking software, unreliable connectivity, low-cost Android phones, and local payment habits.

## Repository

```text
apps/                Reserved; no Ubi application implementation yet
services/            Reserved; no production Ubi API implementation yet
packages/contracts/  Reserved; no executable shared contracts yet
prototypes/
  android_business/  Synthetic Flutter/SQLite V0; host-tested, debug APK inspected
  booking_integrity/ Isolated Go/PostgreSQL/HTTP experiment; disposable DB only
  customer_web/     Scoped, editable guest-link browser experiment; synthetic only
  fixtures/         Invented organizations, operations, services and customers
  validation/       Moderated task protocol; field research not conducted
design/
  tokens.json        Provisional visual/semantic tokens; not contrast-approved
docs/
  blueprint/         Active controlled planning, validation and build backlog
  product-foundation.md, delivery-roadmap.md, architecture.md
                     Archived preliminary inputs, not build authority
.tooling/
  verify-blueprint.mjs Documentation-only structural/traceability checker
  validate-v0.mjs    Reproducible host analyzer/test/evidence runner
  validate-integrity-v0.mjs Synthetic database/API checks and explicit lint gaps
  validate-customer-v0.mjs Scoped form, retry, browser and low-byte component checks
```

## Current programme state

Latest checkpoint — **2026-10-15**: the Home Add-button/card overlap is fixed in source and host tests by reserving an action area below the schedule. The [layout follow-up](docs/blueprint/validation/WB-02-home-action-layout-evidence.md) separates the new 0.1.1 (2) local candidate from the old 0.1.0 (1) APK still on the phone. No phone was reconnected or updated for this fix. The earlier [physical-phone receipt](docs/blueprint/validation/WB-02-physical-android-smoke-evidence.md) records bounded native SQLite/save/review/conflict/revision/cancellation and one separately approved force-stop/draft-recovery path on an authorized Samsung SM-G975W/Android 12, with 11 offline evidence cross-checks. That historical evidence is preserved against its original APK, not transferred to the new build. The testing connection remains stopped. This is not Android 11/4 GB coverage, TalkBack, offline-radio transitions, release performance, or WB-02 closure.

Work has reached **Stage 9 synthetic feasibility validation**. The isolated Android-targeted V0 implements local manual entry/review, tentative SQLite save, draft recovery, conflict simulation, revision/cancellation history and fixture scope switching. It has 41 passing Linux-host repository/widget/startup-asset tests and a clean analyzer, including 13 new action-layout/navigation checks. The separate Go/PostgreSQL experiment exercises booking contention, held/committed capacity, atomic changes, exact quotation conversion and lost-response retry, with saved HTTP contract assertions. Its Postman specification-linter disagreement is recorded explicitly; no all-linters pass is claimed. Neither slice is a production application or an integrated Android/server flow. The 28 full scenario definitions are **not run**, and release gates G7/G8/G9 remain **open**. See the historical [Android host evidence](docs/blueprint/validation/WB-01-02-synthetic-prototype-evidence.md) and [APK evidence](docs/blueprint/validation/WB-02-synthetic-android-apk-evidence.md), the latest [layout follow-up](docs/blueprint/validation/WB-02-home-action-layout-evidence.md), [integrity evidence](docs/blueprint/validation/WB-03-synthetic-integrity-evidence.md) and [active planning index](docs/blueprint/README.md).

The founder-supported business Android floor is **Android 11+ and at least 4 GB RAM**, recorded as DEC-022 on 2026-10-04. Actual phone performance remains unmeasured. Guest web compatibility is separate; the business app's minimum does not automatically exclude guests with older devices.

WB-04 now implements the local customer-link → editable details → review → truthful instant/custom outcome flow, with disposable PostgreSQL behind instant demo bookings. Its [evidence record](docs/blueprint/validation/WB-04-synthetic-customer-link-evidence.md) records 22 top-level tests plus 45 nested cases, repeated three times: 201/201 case results, including saved local Postman and JavaScript-disabled browser checks. This is not a complete public customer app, verified identity or Android/server integration.

Next dependencies: obtain renewed authorization to install/recheck the layout candidate on a phone, remaining native/device coverage and physical Android 11/4 GB release/profile measurements, together with web compatibility/localization, weak-network and comprehension research. The first authorized Android 12 correctness sample is recorded above; further device/setting changes need appropriate direction. Explicitly approved, project-isolated JDK/Android tooling and stable packages are verified; universal debug APKs remain large unoptimized three-ABI artifacts, not production distributions or speed results. See the [setup checkpoint](docs/blueprint/validation/WB-02-android-build-preflight.md) and [partially executed native checklist](docs/blueprint/validation/WB-02-android-native-test-checklist.md). Resolve/disposition the WB-03 linter limitation and obtain relevant evidence/review before WB-05 contract closure. See the [ordered backlog](docs/blueprint/specifications/85-release-roadmap-and-implementation-backlog.md) and [candidate release/evidence contract](docs/blueprint/specifications/72-first-release-contract-and-evidence-gates.md). In-app messaging and voice calling remain required future workstreams, not working features or silently removed scope.

## Intended first build target

The first usable vertical slice is:

1. Set up a photography business.
2. Add or choose a service package.
3. Set availability.
4. Add a phone, walk-in, or message-sourced booking in seconds.
5. See the booking on Today and Calendar.
6. Preserve eligible offline intent; confirm/complete/cancel/reschedule only under the applicable server authority, with explicit conflict recovery.

## Non-negotiable quality targets

- Fast first useful screen on budget Android hardware.
- Offline-readable schedule and safe queued changes.
- One clear primary action per screen.
- Minimum 48 dp touch targets and readable contrast.
- Phone number first; email is optional.
- English and Amharic-ready from the beginning.
- Green only communicates success or approval.
- Human call/message fallback is always understandable.

Run `node .tooling/verify-blueprint.mjs` to check documentation structure and requirement coverage. Passing this check does not establish application correctness, performance, legal approval or launch readiness.

Run `node .tooling/validate-v0.mjs` for the separate synthetic prototype checks. It writes the [machine receipt](prototypes/android_business/artifacts/validation-receipt.json) and actual Flutter-rendered test screenshots; these are host renders, not physical Android screenshots.

Run `node .tooling/validate-integrity-v0.mjs` for disposable PostgreSQL, Go race-detector and local HTTP checks. Its [machine receipt](prototypes/booking_integrity/artifacts/validation-receipt.json) preserves the specification-linter failure as a tooling limitation; exit 2 means this recorded partial result, not a full pass. No existing database URL, real credentials or live data may be used.

Run `node .tooling/validate-customer-v0.mjs` for WB-04 synthetic form/grant/review/retry and local browser checks. See the [demo instructions](prototypes/customer_web/README.md) and [machine receipt](prototypes/customer_web/artifacts/validation-receipt.json). No public service, social inbox or payment is connected.

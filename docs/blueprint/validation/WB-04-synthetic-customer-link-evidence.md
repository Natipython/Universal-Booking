# WB-04 synthetic customer-link evidence

Status: **PARTIAL programme evidence — local component checks pass; package/release gates remain open**  
Updated: 2026-10-11 (Africa/Addis_Ababa)  
Related: specifications 17, 21, 22, 25, 34, 44, 51, 71–73, 85; partial EV-03/04/05

## Outcome and exact boundary

The [local customer-web experiment](../../../prototypes/customer_web/README.md) now implements prepared link → explicit continue → edit → review → truthful outcome. Instant demo bookings use the actual WB-03 transaction service in a newly created private PostgreSQL cluster. Custom enquiries are saved only in demo memory and do **not** reserve capacity. No customer app installation is needed for the browser path; no native app association is implemented.

The recorded run passed **22 top-level tests plus 45 nested cases**, repeated three times with Go's race detector: **201/201 case results**, zero failures/skips. These are 67 distinct component case names, including parent cases, not 201 independent user scenarios. Each repetition executes 21 saved local Postman requests and 109 assertions: **63 requests / 327 assertions** total. Each also runs the JavaScript-disabled Chromium browser harness with **26 bounded checks**.

Machine receipt: [validation-receipt.json](../../../prototypes/customer_web/artifacts/validation-receipt.json). Checked at `2026-10-11T12:38:21.270Z` / `2026-10-11 15:38:21.270 EAT`. Covered-source SHA-256: `f12634708055d2daa1546ad207771d8aa0f1799439ca94ab82d88be86a019d8e`. Sources were unchanged during validation. The receipt includes the exact covered source paths and imported WB-03 runtime/schema inputs. Artifacts are regenerated on rerun; this is a dated component result, not permanent approval or an independent human review.

## What was built and cross-checked

| Boundary | Recorded local evidence | Important limit |
|---|---|---|
| Link preview and redemption | GET/HEAD/crawler preview shows no private draft and does not claim the grant or book capacity; only explicit POST redeems; same-browser retry recovers the same session | An unused copied bearer link may be claimed by another browser. Binding is not verified customer identity |
| Scope and expiry | Hashed 32-byte opaque grants; purpose, mode, organization/operation/service, expiry and revocation checked again after redemption; forwarded used links fail generically | Grant/session/review state is ephemeral; no durable production identity, cross-device or post-install continuation |
| Input and browser privacy | Exact Host; Origin/navigation checks; cookie-bound CSRF; duplicate/unknown/oversized fields, forged authority, query data, invalid calendar/clock/date and unsafe HTML rejected | Not a complete public HTTPS/auth/rate-limit/penetration assessment; compatibility without required browser metadata needs review |
| Review and edit | Changed date/time/customer persists through back-to-edit; review shows selected customer and duration/price or event size before sending; stale/foreign/expired reviews cannot commit | No real minimum-contact flow, actual hours/availability picker, Ethiopian-calendar/clock or Amharic acceptance |
| Booking truth | Actual SQL adapter creates one scoped `confirmed_demo` reservation; capacity and changed terms preserve the draft; response scope/version/price/currency/customer/start/end must match review | Exclusive invented services only; no live business, payment, notice or Android synchronization |
| Retry and unknown outcomes | Eight concurrent duplicate submissions produce one booking/command/outbox; actual HTTP connection dropped after commit then retried without duplicate; committed/absent/mismatched-price unknown results tested | No host/power-loss or distributed restart guarantee. After review expiry, unknown attempts use read-only receipt lookup; no visible receipt means keep checking, not start a new booking |
| Custom enquiry | Selected event size is reviewed; duplicate submit saves one in-memory request, zero bookings/reservations; amber/written pending expectation | Full clarification, versioned quotation, acceptance and guest manage UI are not implemented |
| Browser and byte scope | Core path works with JavaScript disabled; 390px viewport plus 320px/200% root-font simulation; skip-link first focus, no horizontal page overflow, primary controls ≥48 CSS px, no third-party requests | Headless Linux Chromium, not a physical phone, TalkBack, real OS font scaling, full WCAG audit, field comprehension, weak-network or Core Web Vitals proof |

Privacy headers are `no-store`, `no-referrer`, restrictive CSP, frame denial, nosniff and noindex. Cookies are host-only/HttpOnly/SameSite=Lax/short-lived, with Secure enabled for configured HTTPS and explicitly demo-only HTTP names. Generated logs redact opaque secrets; no cookie export, HAR, raw request bodies, real contacts or transcripts are retained.

## Recheck findings and repairs

1. **Expired unknown attempts could reach a creation path.** Recovery after review expiry is now read-only. A missing committed receipt remains checking; absence does not authorize another commitment.
2. **Review omitted selected customer/event details.** The review now exposes those choices and instant duration before explicit submit; UI assertions check them.
3. **A real browser submission failed despite direct HTTP tests.** Chromium sent `Origin: null` under `no-referrer`. The [Fetch Standard's Origin algorithm](https://fetch.spec.whatwg.org/#append-a-request-origin-header), checked 2026-10-11, supports that behavior. The repair retains no-referrer and requires same-origin navigation/document metadata plus the issued cookie and secret CSRF. Null-origin cross-site/same-site/missing-metadata/wrong-mode/wrong-destination/bad-CSRF cases are rejected. This is not a blanket null-origin allowance.
4. **The Postman harness duplicated browser cookies.** Explicit cookie copies were removed; its own cookie jar carries issued cookies. The invalid-session negative remains explicit. No server check was weakened to accommodate the harness.
5. **Display and protocol gaps.** HEAD body suppression, query rejection, exact date/time validation, buffered template rendering, complete response matching, after-outcome routing and large-text decorative logo geometry were checked/repaired.

Earlier failed exploratory runs were not counted as successful evidence. The saved final receipt follows these repairs and a fresh three-run validation.

## Contract and tooling results

Go format/vet/module verification pass. Independent `@apidevtools/swagger-parser` validation accepts OpenAPI 3.1.0. Postman specification lint returns exit 0 with zero violations; collection lint scans all 21 request files with zero errors/warnings. All 63 per-request loopback guards execute. No Postman account login, cloud workspace, spec publication or authenticated governance review was performed; the CLI's unauthenticated warning is retained.

The customer experiment returns HTML/forms, not WB-03's problem-JSON API. **Its lint pass does not resolve WB-03's `application/problem+json` disagreement**, which remains recorded separately. The production API contract and stack are not frozen by either experiment.

Local tools recorded: Go 1.27.1 linux/amd64, PostgreSQL 18.6, Postman CLI 1.69.0, Chromium 155.0.8059.39 and the bundled Playwright harness. These are observed host versions, not production dependency approval.

Uncompressed sampled response bodies:

| Resource | Bytes |
|---|---:|
| Public example HTML | 2138 |
| Generic preview HTML | 1680 |
| Editable details HTML | 2855 |
| Review HTML | 2086 |
| Outcome HTML | 1792 |
| CSS | 4129 |

Required application JavaScript: **0 bytes**. Required third-party assets: **0**. These samples are below the local 16 KiB-per-body experiment guard; they are not compressed transfer totals, production budgets, network latency or speed benchmarks. See the [payload receipt](../../../prototypes/customer_web/artifacts/payload-receipt.json) and [browser receipt](../../../prototypes/customer_web/artifacts/browser-receipt.json). Screenshots are actual synthetic renders, including [review](../../../prototypes/customer_web/artifacts/review-390.png) and [large-text simulation](../../../prototypes/customer_web/artifacts/large-text-320.png), not accepted product designs.

## Safety, reproduction and next dependencies

The [cluster receipt](../../../prototypes/customer_web/artifacts/cluster-receipt.json) confirms a newly created private synthetic database, no existing connection accepted, no TCP listener, server stopped and owned temporary files removed. No persistent database needs recovery; fixtures/schema recreate it. No real customers, credentials, social inboxes, messages, calls, payments, domain/provider/company account, deployment, commit or push changed.

From the repository root:

```sh
node .tooling/validate-customer-v0.mjs
node .tooling/validate-v0.mjs
node .tooling/validate-integrity-v0.mjs
node --test .tooling/verify-blueprint.test.mjs
node .tooling/verify-blueprint.mjs --report
```

The WB-03 runner currently exits 2 for its precisely recorded lint limitation, not an all-checks pass. The earlier 2026-10-11 reruns passed 26/26 Android-host tests and 75/75 integrity case results with unchanged covered source digests. The subsequent [Android APK checkpoint](WB-02-synthetic-android-apk-evidence.md) adds corrected startup-asset packaging, 28/28 host tests and debug compilation/static inspection. These remain separate component results, not native adapter execution, Android/server integration or phone performance.

Next, in order: inspect Android toolchain/license/device prerequisites; build/test the Android native slice once those prerequisites are authorized and available; gather representative device/network/localization/comprehension evidence and review guest compatibility/durable recovery; resolve/disposition remaining technical decisions and the WB-03 lint issue before WB-05 closure. Do not silently accept vendor licenses, invent a physical test phone, turn dummy labels into real identity or advance to production to satisfy a continuation request.

WB-04 is **component-validated, not package/gate-complete**. Selected synthetic subsets of SCN-007/009/010/011 are exercised; all 28 full acceptance scenarios remain NOT RUN and G7/G8/G9 remain OPEN. Required messaging/calling remain WB-13/WB-16, not removed scope or working features. City/vertical/pilot/entity/payment proposals and licensed review are still distinct from decided launch country Ethiopia.

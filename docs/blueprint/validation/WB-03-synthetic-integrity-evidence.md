# WB-03 synthetic Go/PostgreSQL integrity evidence

Status: **PARTIAL — local component checks pass; spec-linter limitation and release evidence remain**  
Updated: 2026-10-11 (Africa/Addis_Ababa)
Related: specifications 33–35, 40–44, 61, 71–73, 85; EV-03/04

## Outcome and exact boundary

The isolated [WB-03 experiment](../../../prototypes/booking_integrity/README.md) now executes an invented exclusive-resource booking schema, Go transaction commands and loopback HTTP contract. The recorded run has **23 top-level tests plus two nested race cases**, repeated three times with Go's race detector: **75/75 case results**, zero failures/skips. Each repetition runs 25 saved Postman requests and 172 response assertions: **75 requests / 516 assertions** in total. These counts describe component checks, not 75 independent product scenarios or a production load benchmark.

Machine receipt: [validation-receipt.json](../../../prototypes/booking_integrity/artifacts/validation-receipt.json). Rechecked at `2026-10-11T12:36:05.794Z` / `2026-10-11 15:36:05.794 EAT`. Covered-source SHA-256 remains `45a97237cea1be90cb0df16591a765cfc144ba6c2ac1e0782cf5809befda0993`. Covered inputs did not change during the run. Receipts/logs are regenerated on rerun; this prose describes this dated build, not a permanent guarantee. No independent human reviewer or release sign-off is recorded.

Overall receipt status is **PARTIAL_TOOLING_LIMITATION**, not an all-checks pass. Core tests, static checks, independent OpenAPI validation and collection lint pass. Postman **specification lint fails**, with the precise disagreement described below. EV-03/04 remain partial and G7/G8/G9 stay OPEN.

## Implemented component slice

- Shared invented organizations/operations/services/customer labels from the Android fixtures; a studio and events operation share one organization rather than unrelated merchant ownership.
- One tenant-qualified reservation table for active holds and committed bookings; composite scope keys, buffered half-open intervals and a GiST exclusion constraint.
- Transactional create, reschedule, cancel, hold consumption and versioned quote acceptance with successful command receipt, audit and outbox in the same commit.
- Canonical-input idempotency, exact integer ETB snapshots, strict bounded JSON, safe error codes and explicit synthetic UTC-offset times.
- Forced tenant RLS and a non-owner/non-superuser/non-bypass runtime role, with transaction-local tenant context. Audit/outbox rows cannot be updated/deleted by that role.
- An experimental [OpenAPI 3.1 contract](../../../prototypes/booking_integrity/openapi.json) and saved [Postman Collection v3](../../../prototypes/booking_integrity/postman/collections/ubi-synthetic-v0/.resources/definition.yaml), developed contract-first using the Postman skill workflow.

`X-Ubi-Demo-Principal` chooses a fixed invented owner/staff/revoked fixture and **is not authentication**. `confirmed_demo` only describes an invented database commitment. The Android client is not connected to this server; quotations have no customer UI here. Nothing sends a customer notice, calls a provider or handles money.

## Executed checks and limits

| Boundary | Recorded outcome | Not established |
|---|---|---|
| Competing capacity | 100 synchronized Go contenders produce one committed booking and 99 capacity conflicts; duplicate successful keys produce one logical effect; direct SQL overlap is rejected without service advisory locks | 100 simultaneous database connections, sustained throughput, many resources/tenants or production double-booking guarantee; runtime pool is 32 connections |
| Hold/commit ownership | Shared collision boundary, same-row hold conversion, expiry reconciliation, active/expired hold-versus-booking races, adjacency and setup/cleanup buffers pass | Autonomous expiry worker, pooled/nightly or multi-resource allocation, real payment-triggered conversion |
| Atomic change and failure | Failed occupied-time/stale/outbox-error moves preserve the original allocation/version/history; cancellation releases capacity without resurrecting old commitment | Real client edit/race UX, cutoff policy, approval workflows, refunds or notice delivery |
| Retry and lost response | An actual loopback handler commits, then its TCP response is deliberately dropped; same-key retry recovers the result with one booking/audit/outbox; changed successful-key input is rejected | Android reconnect integration, multi-process crashes, power-loss durability, retry policy for failed commands or queue delivery |
| Quote terms/conversion | Superseded/stale versions rejected; accepted version replay and concurrent distinct-key acceptance produce one booking; catalogue changes do not silently rewrite issued terms; conflicting conversion leaves no partial booking | Full custom-request/negotiation/guest-consent flow, quote lines/taxes, legal signature or real customer acceptance |
| Scope and privileges | Two-tenant/two-operation negatives, staff/revoked fixture denials, composite foreign keys, missing-context RLS denial, same-connection commit/rollback context reset, non-owner role and append-only privilege checks pass | Real login/current membership, resistance to an application role forging tenant context, complete penetration testing or a production RLS policy |
| HTTP contract | 25 requests / 172 assertions per repetition, all pass; own/foreign lookup, invalid input, missing principal/key, wrong content type, large body, money/schema/state/header/retry checks; all 25 per-request loopback scripts execute | Public HTTPS/auth/rate limits, all protocol edge cases/default 404/405 errors, complete endpoint/security fuzzing or Postman cloud governance review |
| Static and independent validation | `gofmt -l` empty; Go vet and module verification pass; independent OpenAPI 3.1 validation passes; collection lint scans 25 files with zero errors/warnings | Human schema/API/architecture approval, production migrations or dependency vulnerability assessment |

These exercise selected transaction/HTTP/fixture-isolation subsets of SCN-002/004/005/006/011/012/014. All **28 full scenario definitions remain NOT RUN**. Real actors, end-to-end client synchronization, guest grants, representative devices/networks and appropriate research/provider/review evidence are absent.

## Specification-linter disagreement

Postman CLI 1.69.0 specification lint exits 1 with **144 diagnostics**, all media-range/property-name errors on response content whose sole media type is `application/problem+json`. [RFC 9457](https://www.rfc-editor.org/rfc/rfc9457.html) specifies that media type. The separate `@apidevtools/swagger-parser` validator accepts OpenAPI 3.1.0, and actual HTTP schema/content-type assertions pass.

The inference is a validator/tooling disagreement; no vendor-confirmed defect or approved waiver is claimed. The contract was not weakened to suppress diagnostics. Raw [lint stdout](../../../prototypes/booking_integrity/artifacts/postman-spec-lint.stdout.txt) and [stderr](../../../prototypes/booking_integrity/artifacts/postman-spec-lint.stderr.txt) are retained. The unauthenticated CLI also reports an authentication warning; no account login or authenticated governance review was performed. Any diagnostic outside the narrowly recorded media-type case makes the runner fail rather than silently accepting it. Resolve/disposition this limitation with review before WB-05 contract freeze.

Collection execution needs no account here. CLI `--output` required login, so output is captured by the local harness instead. Every request has its own loopback guard because collection-level script inheritance was not observed in this CLI's run. Both limitations were checked against actual CLI output; no Postman workspace, published spec/collection or vendor-feedback submission was created.

## Database safety and tooling

The harness accepts **no existing database URL**. It creates a new private `ubi-wb03-*` temporary directory and PostgreSQL cluster, rejects host/TCP authentication and binds PostgreSQL only to its private Unix socket. Tests reseed invented fixtures in that cluster. HTTP uses a temporary `127.0.0.1` test listener only. The [cleanup receipt](../../../prototypes/booking_integrity/artifacts/cluster-receipt.json) confirms the server stopped and its owned temporary files were removed. There is no retained database to recover; the checked-in schema/fixtures recreate it. If shutdown fails on a future run, the harness preserves the directory and reports the failed cleanup.

Observed local toolchain: Go 1.27.1 linux/amd64; PostgreSQL 18.6; pgx v5.11.0; Postman CLI 1.69.0. Go/Postman/validator dependencies were installed in local `.tooling/` paths, not production environments. Go's archive checksum was verified against [official download metadata](https://go.dev/dl/). The installed versions are feasibility-toolchain facts, not a final-stack approval. [PostgreSQL 18 constraints](https://www.postgresql.org/docs/18/ddl-constraints.html) and [RLS documentation](https://www.postgresql.org/docs/18/ddl-rowsecurity.html) informed the schema/role checks; tests, not documentation alone, establish this bounded result.

The 11-table synthetic schema is deliberately narrower than catalogue 42: quotations combine header/version snapshots, dummy request IDs have no full request lifecycle table, reservation ownership is limited to one exclusive resource, and fixture principals replace real identity/membership tables. It is not an executable production migration or a substitute for the complete logical model. Reconcile and review those differences before approving a production schema.

No public DNS/domain purchase, production/managed database, real data, credential, customer message, payment, provider/company account, commit or push changed. At the earlier integrity checkpoint the separate Android slice had 26 host tests and no SDK/JDK or device result. The subsequent [Android APK checkpoint](WB-02-synthetic-android-apk-evidence.md) records explicitly approved isolated tool setup, 28 passing host tests and debug compilation/static inspection. Native execution and real-device evidence remain absent; no Android/server integration result is inferred.

## Reproduction and next work

From the repository root:

```sh
node .tooling/validate-integrity-v0.mjs
node .tooling/validate-v0.mjs
node --test .tooling/verify-blueprint.test.mjs
node .tooling/verify-blueprint.mjs --report
```

Integrity runner exit **2** preserves the recorded partial/linter-limitation status; it is not an all-linters success. Exit 1 denotes a test/unclassified failure, and 0 requires the local checks and both linters to pass. The machine receipt distinguishes results, covered-source digest, raw test events and unverified scope.

WB-04 now has a separate [local customer-link evidence record](WB-04-synthetic-customer-link-evidence.md), not a public or fully approved guest product. Next dependencies include Android SDK/JDK/license setup and real Android 11+/4 GB/native-adapter/accessibility/performance checks without treating host timings as phone evidence. Obtain relevant web/network/research and technical review, resolve/disposition the tooling issue and review WB-05 bounded contracts. Messaging/calling remain required future workstreams, not working demo buttons or removed scope. Founder city/vertical/pilot/entity/payment facts and legal/provider/production approval remain unresolved.

# WB-03 — synthetic booking-integrity contract V0

Status: **Experimental contract; no production API/authentication or release approval**  
Date: 2026-10-04 (Africa/Addis_Ababa)

This disposable Go/PostgreSQL slice tests exclusive photography-studio resources. It deliberately shares the invented organizations/operations/services/customer labels in `../fixtures/booking-v0.json`. It is separate from production `services/`, cannot send notices or take payment, and does not connect the Android app yet.

## Authority and safety boundary

Only a new, private temporary PostgreSQL cluster may be used. Its runtime role cannot own tables, bypass RLS, perform DDL, delete history or update audit/outbox rows. Transactions set tenant context locally; absent context denies tenant reads/writes. The experiment does not accept an existing/managed database URL.

HTTP is test-only loopback. `X-Ubi-Demo-Principal` selects a server-defined invented owner/staff/revoked fixture; it is **not authentication** and must never become a production login or published endpoint. Organization/operation IDs in JSON are validated against that fixture rather than treated as authority. Production requires real identity, current memberships/capabilities, HTTPS, rate limits and reviewed privacy/security controls.

Every success includes `demo_only: true`; `confirmed_demo` means only that a synthetic reservation committed in this disposable database. It does not confirm any real business booking. Error responses use `application/problem+json`, stable `code`, safe `detail`, status, request ID and retryability. Sensitive cross-scope lookups return 404. Responses use `Cache-Control: no-store`.

## HTTP contract

Machine spec: `openapi.json`, generated from the checked-in contract generator. Versioned experimental paths are `/v0`, not final Ubi `/v1` commitments.

| Method / path | Required input | Success and safety |
|---|---|---|
| GET `/v0/health` | None | 200 synthetic marker; never mutates state |
| POST `/v0/bookings` | Scope, service/version, explicit-offset start, dummy customer/source; Idempotency-Key | 201 synthetic booking, authoritative snapshot, Location; rejects occupied capacity or stale service version |
| GET `/v0/bookings/{id}` | Authorized fixture principal | 200 own scoped item; 404 outside fixture scope; no mutation |
| POST `/v0/bookings/{id}/reschedules` | Expected booking version and explicit-offset new start; Idempotency-Key | 200 same booking ID/new version; original interval survives any failed move |
| POST `/v0/bookings/{id}/cancellations` | Expected booking version; Idempotency-Key | 200 cancelled_demo plus released capacity; history preserved; no refund/notice |
| POST `/v0/holds` | Same service/time scope plus bounded TTL; Idempotency-Key | 201 active hold, not confirmation; held/committed capacity shares one collision boundary |
| POST `/v0/holds/{id}/bookings` | Expected hold version; Idempotency-Key | 201 synthetic booking; converts the same locked reservation, never release-and-reinsert |
| POST `/v0/quotations` | Dummy request ID, service/time scope, optional integer demo price; Idempotency-Key | 201 immutable quote version; supersedes previous unaccepted version, does not allocate capacity |
| POST `/v0/quotations/{id}/acceptances` | Expected quote version; Idempotency-Key | 201 one synthetic booking/quote conversion; stale, expired or conflicting acceptance never partially succeeds |

POST JSON is bounded to 32 KiB, rejects unknown fields/trailing JSON and requires a stable 8–128-character idempotency key. Successful commands replay the original status/representation under the same principal/tenant/action/key and canonical input. Changed input with a previously successful key returns 409. Failed transactions do not retain a successful-command record; callers must not infer confirmation from a timeout or 202.

Money is integer ETB minor units, never float. Service duration/resource/buffers/current price come from the database, not client-supplied end times or a “confirmed” flag. V0 time inputs require an explicit RFC3339 offset and are normalized to UTC; the fixture date/clock is synthetic, not Ethiopian-calendar support. Only fixture labels/sources and years 2029–2031 are accepted.

## Transaction invariants

1. One tenant-qualified `resource_reservations` table holds both active holds and committed bookings. GiST exclusion rejects overlapping `[start-buffer, end+buffer)` intervals per organization/resource. Exact adjacency is allowed; release retains history.
2. Composite foreign keys prevent cross-tenant/operation service/resource links. Runtime RLS is tested separately from application filters. Transaction-local context must not leak through reused connections.
3. Mutation, successful idempotency receipt, outbox and audit commit together. Forced outbox failure rolls everything back. Audit/outbox cannot be rewritten by the runtime role.
4. Resource locks serialize service commands; the exclusion constraint remains a database backstop. Concurrent contenders must yield one capacity owner, not merely one HTTP success.
5. Holds are checked against database wall time after locks, not device clocks. Resource commands reconcile expired holds. There is no autonomous expiry worker in V0; idle expired rows may remain until reconciliation.
6. Rescheduling updates allocation/booking/version atomically. Stale edits, cancellation and unavailable alternatives cannot discard the old valid schedule.
7. Quote acceptance locks/version-checks the quotation and allocates capacity atomically. Duplicate conversion cannot produce another booking or event. Quotes preserve their original price/service snapshot; later catalogue edits do not silently rewrite accepted terms.

## Test plan and non-goals

Required tests: 100 synchronized contenders, held-versus-committed contention, exact adjacency/buffers, tenant and operation negatives, absent-context/pool reset/least-privilege checks, duplicated commands/lost response, atomic move/failure, hold conversion/expiry race, stale and duplicate quotation acceptance, forced rollback, exact money/offset validation, safe HTTP errors and contract assertions.

This does not implement pooled/nightly/multi-resource booking, hotel inventory, real role administration, customer tokens, providers, payment, production migrations/restore, encryption, network/device performance, queue delivery or complete SCN-002–006/011/012/023. SQL/model/API remain candidates until relevant review and WB-05. No legal/entity/payment/provider approval follows from a pass.

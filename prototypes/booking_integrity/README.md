# Ubi WB-03 synthetic booking-integrity experiment

Status: **Component-tested candidate; not a production service or approved release**  
Updated: 2026-10-04 (Africa/Addis_Ababa)

This Go/PostgreSQL experiment tests booking truth before connecting the Android client or customer web page. It uses the same invented two organizations, three operations and studio service fixtures as the Android V0. One organization can contain multiple operations; scope is not represented by separate unrelated merchant identities.

The [contract](CONTRACT.md) and [OpenAPI spec](openapi.json) describe the experimental interface. `confirmed_demo` is a committed **invented** reservation in a disposable database, never a real confirmation. The loopback-only HTTP test server is started by tests, not a publishable application command. Selecting `X-Ubi-Demo-Principal` is a fixture convenience, **not login**.

## What is exercised

- Exclusive held/committed capacity in one reservation table, with database exclusion constraints, buffers and half-open adjacency.
- Atomic create, move, cancel, hold conversion and exact quotation-version conversion; retained booking/history on failure.
- Tenant/operation negative cases, forced RLS, missing-context denial, transaction-local pool reset and a non-owner/non-bypass runtime role.
- Idempotent command/result replay, changed-input rejection and recovery after an actual committed HTTP response is dropped.
- Successful receipt, reservation, audit and outbox in the same transaction; injected outbox failures cannot partially confirm or move a booking.
- Immutable quote snapshots, superseded versions, catalogue-change rejection and concurrent duplicate acceptance.
- Strict bounded HTTP JSON, safe typed errors, integer ETB money and explicit-offset synthetic times.

These are selected component checks, not the complete 28-scenario acceptance suite. Holds expire by database time when resource commands reconcile them; no autonomous timer worker exists. The outbox is durable local evidence, not a delivery worker or a sent message. Runtime tenant settings are trusted server context, not a defense against a compromised application role inventing another context.

## Reproduce locally

Use a non-root Linux account, the installed local Go/Postman runtimes and PostgreSQL 18 binaries at `/usr/lib/postgresql/18/bin`. PostgreSQL's `btree_gist` extension and a C compiler are required. There is **no database URL input**. The test harness creates a new private `ubi-wb03-*` temporary directory, starts a new PostgreSQL cluster without TCP, then stops it and removes only that owned temporary directory. If shutdown fails, it preserves the directory and records a failed cleanup. Do not replace this safety harness with a connection to an existing database.

From the repository root:

```sh
node .tooling/validate-integrity-v0.mjs
```

The runner captures static checks, independent OpenAPI validation, Postman spec/collection diagnostics and three repetitions of race-detector-enabled Go tests. It fingerprints covered sources before/after, saves [validation-receipt.json](artifacts/validation-receipt.json), raw test events and cleanup evidence. Generated logs/receipts may be overwritten on a rerun; they contain dummy data only. Exact versions/dependencies are recorded in the receipt, `go.mod`, `go.sum` and the local tool lockfiles.

For one focused run from this directory:

```sh
UBI_RUN_SYNTHETIC_PG=1 ../../.tooling/go/bin/go test -race -count=1 -timeout=180s ./...
```

Contract/collection generators are `.tooling/build-integrity-contract.mjs` and `.tooling/build-integrity-collection.mjs` at the repository root. Regenerate and review their changes before revalidation. Every saved request has a loopback guard; the default target is deliberately unusable until the harness supplies its temporary loopback address.

## Postman workflow and recorded limitation

The Postman API-engineering workflow guided contract-first input/output/error definition and saved Collection v3 assertions. Bootstrap was local (`--no-cloud --no-fetch`), no workspace/account was linked and no collection/spec was published. Runs use `--no-report-events`, with local report capture; no authenticated governance review is claimed. CLI `--output` required login here, so the harness captures terminal output rather than logging in or exporting to a workspace.

The recorded CLI spec-lint run rejects `application/problem+json` with media-range/property-name diagnostics. This media type is specified by [RFC 9457](https://www.rfc-editor.org/rfc/rfc9457.html); the independent OpenAPI 3.1 validator accepts the contract, and actual response-schema/content-type assertions pass. This is a **recorded validator disagreement**, not a spec-lint pass, vendor-confirmed defect or approved release waiver. The runner reports `PARTIAL_TOOLING_LIMITATION` and exits **2** for this known case (0 only for all checks/linters passing, 1 for a test/unclassified failure). Any new diagnostic fails the bounded check. The limitation must be resolved/dispositioned with review before contract freeze. No credentials or vendor-feedback submission are needed for this local slice.

## Not implemented or approved

No real authentication/membership administration, HTTPS public API, rate limits, merchant working-hours configuration, Ethiopian calendar/clock UI, guest tokens, Android/server sync, multi-resource/pool/nightly inventory, migration/restore, worker delivery, external messages/payments or native messaging/calling exists here. No phone-speed or production-throughput claim follows from localhost timings. Founder/reviewer/provider/legal decisions and G7/G8/G9 remain open.

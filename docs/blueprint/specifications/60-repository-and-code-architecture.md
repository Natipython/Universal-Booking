# Repository, code and language architecture

Status: **Draft — proposed implementation architecture, not permission to start production build**  
Depends on: domain/platform/data/security blueprints, ADR-002, quality attributes and team/hosting decisions  
Blocks: repository scaffolding, build tooling, dependency selection, CI/CD and implementation slices

## Purpose

This document turns Ubi's logical architecture into maintainable source-code boundaries. It defines where Go, Flutter/Dart, React/TypeScript, Kotlin/Java, Swift, Rust, C, C++ and Assembly may be used without creating nine competing implementations of the same business rules.

The primary constraint is a small early team serving low-confidence users on constrained devices and networks. Operational simplicity and correctness outrank language enthusiasm.

Diagram: `diagrams/60-code-and-language-boundaries.mmd`.

## Initial deployable set

```text
1. Go application/API deployable
2. Go worker/scheduler deployable (may share modules and release)
3. Flutter/Dart Ubi Business Android application
4. React/TypeScript web application serving customer, business and operations route groups
5. PostgreSQL migrations and generated contracts
```

Logical modules are numerous; independently deployed services are deliberately few. Native libraries are added only by ADR-002 gates.

## Proposed monorepo

```text
/
├── apps/
│   ├── business_mobile/          # Flutter/Dart; Android first, iOS later
│   └── web/                      # React/TypeScript SSR: customer/business/ops boundaries
├── services/
│   ├── platform/                 # Go modular application/API
│   └── worker/                   # Go worker entrypoint; shared domain/application packages
├── modules/                      # Go bounded-context packages
│   ├── identity/
│   ├── organization/
│   ├── catalog/
│   ├── scheduling/
│   ├── booking/
│   ├── quotation/
│   ├── conversation/
│   ├── integration/
│   ├── payment/
│   ├── notification/
│   ├── trust/
│   ├── privacy/
│   └── audit/
├── packages/
│   ├── contracts/                # OpenAPI/event/schema sources + generated clients
│   ├── design_tokens/            # Semantic design data, generated per client
│   ├── localization/             # Source messages/glossary, not secrets
│   └── test_fixtures/            # Synthetic builders and provider contract fixtures
├── native/
│   ├── include/                  # Versioned C ABI headers
│   ├── rust/                     # Optional approved Rust component
│   └── cpp/                      # Optional approved C/C++ component/wrappers
├── database/
│   ├── migrations/
│   ├── policies/                 # RLS/roles/grants
│   └── verification/
├── infra/                        # Provider-neutral modules then selected IaC
├── tests/                        # Cross-system/e2e/performance/security/recovery suites
├── tools/                        # Pinned developer/release helpers
└── docs/                         # Controlled blueprint, ADRs, runbooks and API docs
```

Exact framework-generated folders may differ, but ownership and dependency direction must remain.

## Dependency direction

```text
transport/UI/provider adapter
        ↓
application use case / command/query
        ↓
domain model and policy interfaces
        ↓
repository/provider ports
        ↓
PostgreSQL, HTTP, object, notification and OS adapters
```

- Domain packages do not import HTTP, Flutter, React, SQL driver, provider SDK or cloud SDK.
- Application code coordinates transactions/idempotency/authorization through explicit ports.
- Adapters translate external schemas into internal typed commands and preserve raw evidence only where approved.
- No module reads another module's tables directly for mutations. Queries use documented read models/interfaces.
- Cross-module events/contracts are versioned; in-process calls do not excuse undocumented coupling.
- Circular imports and shared `utils` dumping grounds are prohibited.

## Go server structure

Each bounded context contains, where needed:

```text
domain/          entities, values, invariants, state transitions
application/     commands, queries, authorization requirements, transactions
ports/           repository/provider interfaces owned by the consumer
adapters/        PostgreSQL/provider implementations
transport/       HTTP/webhook/job input-output mapping
observability/   safe domain metrics and event descriptors
tests/           examples, properties, concurrency and contract fixtures
```

Rules:

- `context.Context` deadline/cancellation crosses every I/O boundary; it is not stored in domain objects.
- Goroutines are created only by owned bounded supervisors/pools, never unbounded per item.
- Database pool, provider concurrency and queue leases have explicit budgets.
- Monetary values, IDs, time instants, timezones and policy versions use domain types, not loose strings/floats.
- Errors are typed/classified as validation, authorization, conflict, retryable dependency, unknown external outcome or internal fault.
- Panics are defects, recovered only at process boundaries with safe failure/telemetry; they do not become normal control flow.
- Build output is reproducible as far as practical, includes version/commit/schema compatibility and supports graceful shutdown.

## Flutter/Dart mobile structure

```text
lib/
├── app/                    bootstrap, navigation, themes, localization
├── core/                   typed errors, clocks, IDs, telemetry contracts
├── features/<capability>/
│   ├── presentation/
│   ├── application/
│   ├── domain/
│   └── data/
├── local/                  relational schema, outbox, sync metadata, secure references
├── sync/                   bootstrap/pull/push/conflict coordinator
├── contracts/              generated API/event types
└── platform/               narrow plugin interfaces
```

- Screens render durable repository state, not raw request futures.
- Eligible actions are transactionally stored before showing “saved on this phone.”
- Widget code does not contain pricing, permission, payment-verification or booking-allocation truth.
- State survives rotation/process death where required; background work follows OS policy.
- Platform plugins are behind Dart interfaces and fakeable in tests.
- Android build contains only required ABIs/resources; native symbols and mapping files are retained securely.

## React/TypeScript web structure

Route groups separate public customer, authenticated merchant and workforce operations policy. They may share components/tokens but never authentication assumptions.

- Server rendering supplies business/service/price method/next action before client hydration.
- Progressive enhancement keeps critical customer forms usable with delayed JavaScript where practical.
- TypeScript strict mode and generated server contracts are required; runtime validation still applies at trust boundaries.
- Server-only secrets/provider calls cannot enter client bundles.
- Authenticated/one-time-token/private responses are never placed in public shared cache.
- Accessible HTML semantics precede custom interaction; keyboard/screen-reader behavior is tested.
- Bundle budgets are route-specific; operations code is not shipped to public customer routes.

## Native module architecture

Native code is a separate high-risk dependency boundary:

```text
Dart / Kotlin / Swift / Go safe wrapper
              ↓
versioned C header and generated/tested bindings
              ↓
Rust or C++ implementation
              ↓
OS/hardware or audited native dependency
```

Requirements:

- no domain/database/network credentials inside native module;
- byte/CPU/memory/time bounds for untrusted input;
- no cross-boundary exception/unwind;
- explicit allocator and thread ownership;
- cancellation and timeout behavior;
- deterministic error codes plus safe diagnostics;
- sanitizer/fuzzer/static-analysis builds and corpus regression;
- per-ABI artifact checksum, SBOM, license, signature and symbol archive;
- safe fallback or feature disablement if the module cannot load.

Go uses cgo only for an approved dependency that cannot be safely isolated as a process or replaced. Cgo changes cross-compilation, static linking, crash/debug and supply-chain behavior and must not be introduced casually.

## Contract generation

One reviewed source defines each external contract. Generated Go/Dart/TypeScript bindings are outputs, not hand-edited truth.

- OpenAPI/JSON Schema or another approved schema describes HTTP payloads/errors/idempotency.
- Event and sync schemas include compatibility/version policy.
- Database schema is not exposed as an API contract.
- Code generation is pinned, reproducible and diff-reviewed.
- Generated code has ownership and vulnerability/update path.
- Consumer/provider contract tests verify old supported clients and provider fixtures.

## Configuration and secrets

- Typed configuration loads at startup from environment/config service; invalid required configuration fails closed.
- Secrets are references delivered by workload identity/secret manager, not repository files, images, mobile apps or web bundles.
- Safe defaults differ by environment only through controlled configuration; debug bypasses cannot compile into release unintentionally.
- Feature flags have owner, purpose, scope, default, expiry and removal issue.
- Country/provider differences use policy/config/adapter interfaces, not scattered conditionals.

## Dependency governance

Before adding a dependency:

- prove necessity versus platform/standard library/small owned code;
- review maintainer/activity/license/security/transitives/binary size/data behavior;
- pin with lockfile/checksum and generate SBOM;
- record owner, update policy and replacement/removal path;
- sandbox/minimize permissions for code generators/build scripts;
- run platform/device/build compatibility tests.

High-risk SDKs, native binaries, analytics/crash tools and provider SDKs also pass document 55 vendor/data review.

## Code ownership and review

- Every top-level area has accountable and backup owners.
- Booking, money, authorization, sync, native/unsafe code, migrations, security/privacy and production infrastructure require specialist review.
- Author cannot be sole approver of a high-risk change.
- Review checks behavior, failure, data, operations and rollback—not only style.
- Temporary exceptions have owner, risk, expiry and removal plan.

## Language addition/removal gate

A new production language requires an ADR covering:

1. measured problem and rejected existing-language solution;
2. team owner/on-call/hiring feasibility;
3. build/test/security/SBOM/signing/patch pipeline;
4. observability/debug/profiling/crash support;
5. data/contract/deployment boundaries;
6. total lifecycle cost and exit plan.

The same review can retire a language/component when its value no longer exceeds its cost.

## Release gates

- ADR-002 prototype evidence exists and DEC-031/032/033/044 are approved before full scaffolding.
- Repository dependency rules are automatically checked; no cyclic/domain-to-adapter violation.
- All client/server/provider contracts are schema-defined and compatibility-tested.
- Go, Flutter and React critical paths meet the applicable quality budgets on real D1/W1/W2 conditions.
- Native code, if any, has measured need, safe wrapper, per-ABI tests, fuzz/sanitizer evidence, symbols, SBOM and fallback.
- No booking/payment/authorization truth exists only in UI, provider adapter or native module.
- Secrets, debug bypasses and production data are absent from repository/build artifacts.
- Owners, runbooks, version/rollback and dependency patch paths exist for every deployable.

## Open decisions

- Final Flutter versus native Android result after prototype.
- Exact React server framework/runtime and hosting after deployment/provider decision.
- Go router/data-access/migration/codegen libraries after minimal dependency evaluation.
- Whether any C/C++/Rust component is needed in first release; default is none.
- Supported Android ABIs/OS versions and later iOS/desktop minimums.
- Monorepo CI/build orchestration tool only if plain language-native tooling becomes insufficient.


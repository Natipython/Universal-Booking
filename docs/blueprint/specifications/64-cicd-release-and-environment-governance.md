# CI/CD, release and environment governance

Status: **Draft — provider/store/tooling not selected**  
Depends on: repository/test/observability/deployment/security/vendor architecture  
Blocks: production credentials, automated deployment, app-store release and pilot promotion

## Purpose

Ubi needs repeatable, reviewable releases across Go services, Flutter Android, React web, PostgreSQL and any approved native libraries. CI/CD must reduce change risk without giving automation broad uncontrolled access or treating rollback as a substitute for data compatibility.

Diagram: `diagrams/64-build-promote-release-flow.mmd`.

## Source-control model

- One protected mainline is continuously releasable; short-lived branches and small reviewed changes are preferred.
- Direct production/mainline pushes are prohibited except documented emergency mechanism with retrospective review.
- Required review and checks scale with risk; code owner cannot approve their own high-risk change alone.
- Signed commits are optional policy; release artifacts/provenance and protected identity are mandatory concerns.
- Generated/vendor/lock/schema/native binary changes are visible and reviewed, not hidden by automation.
- Repository permissions are least privilege, regularly reviewed and separate from production runtime access.

## Change classification

| Class | Examples | Minimum additional control |
|---|---|---|
| C0 — documentation/no runtime | Clarification, non-executable diagram | Review, link/structure check |
| C1 — low-risk compatible | Copy, additive optional UI, internal metric | Normal tests/review/flag where useful |
| C2 — domain/data/client behavior | Booking flow, permission, sync, provider adapter | Specialist review, contract/integration/e2e, staged rollout |
| C3 — high-risk | Money, ownership/recovery, RLS, migration, retention, native/unsafe, auth/security | Two-person/specialist approvals, adversarial and rollback/forward evidence |
| C4 — emergency | Active incident containment | Named incident authority, minimal scope, live verification, retrospective within fixed period |

## Pipeline stages

```text
source review
→ deterministic dependency/tool resolution
→ static/schema/security/license checks
→ unit/property/component builds and tests
→ integration/PostgreSQL/contract tests
→ immutable artifacts + SBOM + provenance + signatures
→ staging deployment and migration rehearsal
→ end-to-end/accessibility/performance/security/recovery gates
→ approval against one exact artifact set
→ progressive production promotion
→ automated/user/domain health observation
→ complete, continue, disable, roll forward or rollback safely
```

Build once and promote the same digest. Rebuilding “the same version” for production creates unreviewed output.

## Language-specific build controls

### Go

- pinned supported toolchain and modules/checksums;
- format/vet/static analysis, unit/race/integration/security tests as applicable;
- reproducible version metadata and minimal runtime image/non-root identity;
- no compiler/debug/admin endpoint accidentally exposed;
- graceful shutdown and compatibility smoke test.

### Flutter/Dart and Android

- pinned Flutter/Dart/Gradle/Android toolchain and lockfiles;
- analyze/test/build release per supported ABI/flavor;
- package/binary size, permissions/manifest, SDK/data and native-library diff;
- signed bundle through protected signing service/key access;
- symbols/mapping/native debug artifacts retained privately;
- internal track → controlled pilot → broader release; staged rollout and kill switch where store permits.

### React/TypeScript web

- pinned runtime/package manager/lock and install integrity;
- strict type/lint/test/SSR/build/bundle/security checks;
- public/private route cache and secret-in-bundle scanning;
- content-addressed static assets and source maps protected appropriately;
- server/web asset compatibility across rolling deployment.

### C/C++/Rust/native

- pinned compiler/target/build scripts/dependencies;
- per-platform/ABI matrix, fuzz/sanitizer/static analysis and wrapper tests;
- export-symbol/ABI compatibility and forbidden API/architecture checks;
- signed checksummed artifacts, SBOM/license/source/provenance and symbols;
- no opaque prebuilt binary without vendor/security/exit review.

## Artifact inventory

Each release manifest binds:

- source commit and reviewed change set;
- Go/web/mobile/native artifact digests and signatures;
- database migration set/schema compatibility range;
- API/event/sync/provider contract versions;
- configuration schema and safe defaults;
- feature-flag definitions;
- SBOM, dependency/license/security results;
- symbols/source maps and evidence bundle;
- deployment/IaC/module/provider capability versions;
- rollback/roll-forward and minimum client constraints.

## CI identity and secret protection

- CI uses short-lived workload identity where available; no long-lived cloud/store keys in repository variables when avoidable.
- Fork/untrusted pull-request code cannot access protected secrets, signing or production networks.
- Build jobs receive minimum scoped secret only at required step; logs/artifacts are scrubbed.
- Signing keys live in protected KMS/HSM/store mechanism; access is approved/audited and supports recovery/rotation.
- Production deploy identity cannot alter source-review rules; source admins cannot silently use production data.
- Runners are ephemeral/isolated, patched and do not reuse contaminated workspace/caches across trust levels.

## Environment governance

The environment model in specification 39 applies. Additional rules:

- environment configuration is versioned and reviewed without committing secrets;
- production/non-production cloud, database, provider apps/webhooks, domains and signing are separated;
- production data never seeds lower environments by default;
- staging tests real topology/contracts with synthetic data and external sandbox accounts;
- ephemeral environments auto-expire with DNS, data, credential and cost cleanup;
- drift detection compares actual infrastructure/configuration to approved source;
- console/manual changes are emergency/diagnostic only and reconciled back to source.

## Database release sequence

1. Validate backup/PITR and restore readiness proportional to change.
2. Run compatible **expand** migration with lock/statement timeout and observable progress.
3. Deploy code that supports old and new schema.
4. Backfill in bounded resumable idempotent batches.
5. Verify constraints/counts/quality and all old readers/writers retired.
6. Apply **contract** migration in a later reviewed release.

Application rollback cannot undo already written incompatible data. Prefer forward correction or feature disablement; whole-database PITR is not ordinary deployment rollback.

## Progressive delivery

### Server/web

- deploy small canary or one replica/traffic cohort first where platform supports;
- readiness only after configuration/schema/dependency validation;
- observe task/domain/SLO/queue/database/provider/cost signals, not only HTTP health;
- pause automatically/manual on gate breach;
- drain gracefully and preserve accepted work during replacement.

### Android

- internal engineering track → field-test/pilot track → staged percentage rollout;
- server remains compatible with supported previous app/protocol versions;
- remote feature/capability flags can disable risky optional flows but cannot weaken authorization/migration;
- minimum safe version only for proven security/integrity necessity with understandable upgrade/fallback path;
- rollback means stopping rollout and shipping fixed version; installed clients cannot be instantly replaced.

## Feature flags

Every flag records owner, purpose, type, scope, default, dependencies, created/expiry date, metrics and removal issue.

- default off/fail-safe for high-risk incomplete capability;
- server enforces authorization/domain rule regardless of client flag;
- mutually inconsistent combinations are tested/prevented;
- flag evaluation and material changes are observable/audited;
- stale flags are removed; flags are not permanent configuration dumping ground.

## Release decision record

Before promotion, release owner records:

- included changes/artifacts/migrations/flags/providers;
- requirement/risk/control evidence and known issues;
- customer/merchant/support/data/legal impact;
- performance/capacity/cost/SLO headroom;
- deployment/verification/stop/rollback/roll-forward plan;
- responsible on-call/support and communication;
- approvals appropriate to change class.

## Rollback and roll-forward

| Failure | Preferred response |
|---|---|
| Optional feature defect | Disable scoped flag; preserve evidence/state |
| Stateless server regression | Stop promotion/restore compatible prior artifact |
| Web asset/server mismatch | Restore matched manifest/digests; never mix arbitrary builds |
| Compatible DB code bug | Roll back code or roll forward fix while schema remains expanded |
| Data-writing/migration defect | Stop writers, preserve evidence, run reviewed correction/restore decision |
| Provider policy/schema break | Disable capability, queue/quarantine safely, link/manual fallback |
| Android regression | Stop staged rollout, server-side disable optional path, publish fixed build |
| Security compromise | Incident runbook; revoke/rotate/contain, not blind redeploy only |

Rollback is rehearsed. The command existing in documentation is not evidence it works.

## Dependency and vulnerability releases

- Automated update proposals are reviewed/tested; no blind merge by severity label.
- Critical exploitable findings receive immediate triage against actual exposure and compensating control.
- Runtime/base/toolchain/provider SDK support/EOL dates are inventoried and alerted.
- Emergency patch still produces signed artifact/evidence and compatibility checks.
- Native/mobile store review lead time is included in vulnerability response planning.

## Release verification

After promotion verify:

- correct artifact/config/schema/flag/provider capability active;
- synthetic and real-safe critical flows;
- booking/payment/tenant/offline invariants and audit events;
- error budget, latency, crash/ANR, queue/database/provider health;
- no telemetry gap, cost spike or support surge;
- old replicas/jobs drained and migration/backfill state correct.

The release is not complete until the observation window and handoff are complete.

## Release gates

- Protected review/check rules and CODEOWNERS/approvals match risk classes.
- CI secrets/signing/production identity are isolated and tested for recovery/revocation.
- One immutable manifest/digest set is promoted across environments.
- SBOM/provenance/signature/vulnerability/license evidence exists for every artifact, including native.
- Database expand/backfill/contract and old-client compatibility are rehearsed.
- Canary/staged rollout, stop, rollback/roll-forward and feature disablement are exercised.
- On-call/support/runbooks/dashboards/status communication are ready before pilot users receive release.
- No production deployment depends on one developer laptop or undocumented console clicks.

## Open decisions

- Source/CI/artifact/signing/cloud/app-store providers.
- Branch/review tooling and exact change-class approvers.
- Android signing/account ownership and store-track governance.
- Canary platform and release observation windows.
- Supported client compatibility/minimum-version policy.
- Infrastructure-as-code tool after cloud selection.


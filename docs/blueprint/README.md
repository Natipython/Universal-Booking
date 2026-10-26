# Ubi architecture and operating blueprint

Status: **Working blueprint — synthetic validation permitted; release contract not frozen**  
Owner: Ubi product and technology leadership  
Last updated: 2026-10-15 (Africa/Addis_Ababa)

This directory is the controlled source of truth for designing, validating, building, operating, and governing Ubi. It covers the product, platform, applications, data, security, privacy, permissions, integrations, payments, operations, and commercial system.

No screen, table, API, vendor, or implementation technology is considered approved merely because it appears in a draft. Decisions become buildable only after their dependencies and acceptance gates are satisfied.

## Document states

- **Exploration** — alternatives and questions; not a decision.
- **Draft** — coherent proposal; incomplete or unvalidated.
- **Proposed** — ready for stakeholder review.
- **Approved** — authorized source for implementation.
- **Implemented** — verified in a released system.
- **Superseded** — retained for history but no longer authoritative.

Every substantial decision must record its state, owner, date, rationale, alternatives, consequences, legal/data impact, migration or reversal cost, and validation evidence.

## Blueprint hierarchy

```text
Vision and market boundary
└── Stakeholders, actors and responsibilities
    └── Business capabilities and end-to-end workflows
        └── Domain model and lifecycle rules
            ├── Product experience architecture
            ├── Platform and integration architecture
            ├── Database and data architecture
            ├── Security, privacy and legal architecture
            └── Operating and reliability architecture
                └── Release contract and implementation backlog
```

The dependency direction is intentional. A database schema must implement the domain model; the domain model must implement validated workflows; workflows must solve a verified business problem for defined actors.

## Current programme position

- **2026-10-15 layout follow-up:** the observed Home Add-button/card overlap is fixed in source/host tests; [the follow-up evidence](validation/WB-02-home-action-layout-evidence.md) records 41 passing host tests and a separate 0.1.1 (2) candidate awaiting authorized phone recheck. No phone was reconnected or updated. The original 0.1.0 (1) APK and its covered sources/receipts are preserved separately.
- **2026-10-15 native checkpoint:** the earlier inspected 0.1.0 (1) APK was freshly installed on an authorized Samsung SM-G975W/Android 12 and bounded native SQLite/review/save/conflict/revision/cancellation/scope/draft checks passed, including one separately approved force-stop. [Physical-phone evidence](validation/WB-02-physical-android-smoke-evidence.md) records 11 offline cross-checks, the originally observed overlap and the stopped testing connection. This supersedes the earlier absence of all native evidence, not missing Android 11/4 GB, accessibility, performance, broader lifecycle or research evidence. No full scenario/gate is closed.
- Stage 7 security, privacy, policy and legal **architecture drafts** (`50`–`57`) and their diagrams now exist and have passed structural/consistency review.
- Ethiopia is the decided first-launch country. Gate G7 is **not passed**: the Ethiopian launch city, contracting entity, data roles, final payment arrangement and other listed decisions remain open, and licensed Ethiopian counsel has not reviewed the final release boundary.
- Stage 8 reliability, engineering and operations **architecture drafts** (`60`–`70`) and their diagrams now exist and have passed structural review. ADR-002 proposes Go + Flutter/Dart + React/TypeScript with tightly gated Kotlin/Java, Swift, Rust, C/C++ and Assembly boundaries.
- Gate G8 is **not passed**: tooling/providers are not selected, representative-device/pilot measurements do not exist, and critical runbooks have not yet been exercised.
- Stage 9 validation **definitions and release-candidate drafts** (`71`–`73`) now exist, together with the previously missing experience specifications (`23`–`26`) and business execution drafts (`80`–`85`). All 79 ledger requirements have planned scenario coverage; none of the 28 full scenarios has been executed. An isolated synthetic Android-targeted V0 has 41 passing host repository/widget/startup-asset tests including the layout follow-up; historical [host evidence](validation/WB-01-02-synthetic-prototype-evidence.md) and [inspected debug APK](validation/WB-02-synthetic-android-apk-evidence.md) retain their original 28-test/0.1.0 boundaries. A separate Go/PostgreSQL/loopback HTTP experiment now has component correctness checks and [its own evidence record](validation/WB-03-synthetic-integrity-evidence.md), including a specification-linter disagreement, not an all-linters pass.
- Gate G9 is **not passed** and the release contract is **not frozen**. Field research, founder release decisions, actual team/budget/provider facts, physical-device/application tests, specialist review and operating exercises remain missing.
- WB-01 fixtures/task protocol are prepared; WB-02 is implemented and host-tested, with debug compilation/static inspection complete and one bounded Android 12 physical sample recorded; remaining native/device checks are pending. Its explicitly approved isolated [Android/JDK setup](validation/WB-02-android-build-preflight.md) and required stable packages are verified. WB-03 has a disposable database/HTTP implementation and synthetic test evidence; real identity, migration/restore, sustained load, review and its linter limitation remain pending. WB-04 now has [local customer-link component evidence](validation/WB-04-synthetic-customer-link-evidence.md): editable prepared context, exact review, distinct instant/custom outcomes, privacy/retry checks and JavaScript-disabled Chromium paths. It has no verified customer identity or persistent production continuation. Remaining Android native/device checks, web/network/localization/comprehension evidence and review remain dependencies before WB-05 closure. The founder's Android 11+/4 GB business baseline is DEC-022; the above-minimum sample is not that baseline's verification. Synthetic prototypes do not authorize real-data processing, customer messages, real money or production release. Drafting later stages does not waive G7 or G8.
- Recheck scope, repaired contradictions and unresolved decisions are recorded in specification 73. `node .tooling/verify-blueprint.mjs` checks document structure/identifier coverage only, not rendered diagrams or application/legal readiness.

## Catalogue

The catalogue is also the intended document set. A listed path that does not yet exist is planned work, not an implied completed specification.

### Programme control

- `00-program-plan.md` — complete order of work, deliverables, reviews, and gates.
- `01-requirements-ledger.md` — traceable requirements already stated or inferred.
- `02-decision-register.md` — unresolved and approved product/technical/business decisions.
- `03-scope-and-principles.md` — product promise, boundaries, and non-negotiable rules.
- `04-risk-register.md` — product, market, engineering, data, legal, vendor, and operational risks.
- `05-glossary.md` — one precise meaning for every important Ubi term.

### Business and product blueprint

- `specifications/10-stakeholders-and-actors.md`
- `specifications/11-identity-organization-and-tenancy.md`
- `specifications/12-roles-permissions-and-approvals.md`
- `specifications/13-business-taxonomy.md`
- `specifications/14-booking-archetypes.md`
- `specifications/15-standard-booking-workflows.md`
- `specifications/16-custom-request-and-quotation-workflows.md`
- `specifications/17-conversation-handoff-and-human-takeover.md`
- `specifications/18-payments-deposits-refunds-and-reconciliation.md`
- `specifications/19-notifications-reminders-messaging-and-calls.md`
- `specifications/20-functional-requirements.md`
- `specifications/21-quality-attributes.md`
- `specifications/22-accessibility-localization-and-calendar.md`
- `specifications/23-android-business-experience.md`
- `specifications/24-desktop-business-experience.md`
- `specifications/25-customer-web-and-app-experience.md`
- `specifications/26-ubi-operations-console.md`

### System and platform blueprint

- `specifications/30-system-context.md`
- `specifications/31-domain-and-bounded-contexts.md`
- `specifications/32-platform-components.md`
- `specifications/33-api-and-contract-governance.md`
- `specifications/34-event-and-workflow-architecture.md`
- `specifications/35-offline-sync-and-conflict-resolution.md`
- `specifications/36-channel-integration-architecture.md`
- `specifications/37-payment-integration-architecture.md`
- `specifications/38-media-search-notification-and-call-architecture.md`
- `specifications/39-deployment-network-and-environment-architecture.md`

### Database and data-system blueprint

- `specifications/40-data-domain-model.md`
- `specifications/41-operational-database-architecture.md`
- `specifications/42-logical-schema-and-table-catalogue.md`
- `specifications/43-keys-constraints-indexes-and-partitioning.md`
- `specifications/44-multitenancy-authorization-and-data-isolation.md`
- `specifications/45-cache-search-object-and-secret-storage.md`
- `specifications/46-data-classification-retention-and-deletion.md`
- `specifications/47-analytics-event-model-and-metric-dictionary.md`
- `specifications/48-data-quality-lineage-and-governance.md`
- `specifications/49-migrations-backups-restore-and-disaster-recovery.md`

### Trust, policy, privacy and legal blueprint

- `policies/50-threat-model-and-security-controls.md`
- `policies/51-privacy-and-consent-architecture.md`
- `policies/52-jurisdiction-and-regulatory-matrix.md`
- `policies/53-platform-permissions-and-oauth-governance.md`
- `policies/54-fraud-abuse-moderation-and-safety.md`
- `policies/55-vendor-subprocessor-and-data-transfer-governance.md`
- `policies/56-legal-document-register.md`
- `policies/57-records-audit-and-law-enforcement-process.md`

### Engineering and operations blueprint

- `specifications/60-repository-and-code-architecture.md`
- `specifications/61-testing-and-quality-engineering.md`
- `specifications/62-performance-capacity-and-cost-model.md`
- `specifications/63-observability-slo-and-reliability.md`
- `specifications/64-cicd-release-and-environment-governance.md`
- `specifications/65-support-onboarding-training-and-service-operations.md`
- `runbooks/66-incident-response.md`
- `runbooks/67-security-and-privacy-incident-response.md`
- `runbooks/68-integration-outage-and-message-recovery.md`
- `runbooks/69-payment-reconciliation-and-dispute-response.md`
- `runbooks/70-backup-restore-and-business-continuity.md`

### Business execution blueprint

- `specifications/80-market-validation-and-pilot.md`
- `specifications/81-commercial-model-pricing-and-unit-economics.md`
- `specifications/82-partnership-and-vendor-strategy.md`
- `specifications/83-kpi-scorecard-and-experiment-governance.md`
- `specifications/84-team-ownership-raci-and-hiring-sequence.md`
- `specifications/85-release-roadmap-and-implementation-backlog.md`

### End-to-end validation and release control

- `specifications/71-end-to-end-validation-and-traceability.md`
- `specifications/72-first-release-contract-and-evidence-gates.md`
- `specifications/73-blueprint-audit-and-decision-closure.md`

### Diagrams and decision records

- `diagrams/` — versioned Mermaid source for context, container, data, sequence, lifecycle, permission, deployment, and recovery diagrams.
- `decisions/` — architecture decision records (ADRs), one decision per file.

## Traceability rule

Every implementation item must trace through this chain:

```text
Business objective
→ actor and scenario
→ requirement ID
→ workflow and domain rule
→ security/privacy/legal obligations
→ API/event/data contract
→ acceptance test and operational metric
→ implementation item
```

An item without that chain is either missing analysis or unnecessary scope.

## Blueprint completion definition

The blueprint is complete enough to authorize the first release only when:

1. The launch country, city, first business vertical, and commercial operator are decided.
2. Every first-release actor and permission is explicit.
3. Every first-release workflow has happy paths, exceptions, state transitions, and audit requirements.
4. Functional and non-functional requirements are measurable.
5. The logical data model and retention rules are reviewed together.
6. Threat, privacy, payment, messaging, and integration risks have controls and owners.
7. Android, desktop, customer, and operations responsibilities are separated clearly.
8. Failure, offline, retry, conflict, backup, and recovery behaviour is testable.
9. Country-specific legal counsel has reviewed the applicable policy and regulatory matrix.
10. Pilot scenarios demonstrate that the release boundary can create value without unbuilt future features.

## Relationship to earlier documents

The product foundation, architecture and delivery-roadmap documents directly under `docs/` are archived preliminary inputs. Their old scope, thresholds and lifecycle simplifications are not implementation authority. The controlled blueprint is the active planning record, but a draft in this directory still does not authorize a live release merely by existing.

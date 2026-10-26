# Ubi blueprint and delivery programme

Status: **Draft**  
Purpose: define the complete order of investigation, decision, specification, validation, and implementation.

## Programme rule

Ubi will be built in verified vertical slices, but planned as one coherent platform. “Plan everything” means discovering cross-system obligations and future boundaries now; it does not mean implementing every future feature in the first release.

## Stage 0 — Programme control and evidence

### Questions

- What has the founder explicitly requested?
- Which statements are decisions, preferences, assumptions, aspirations, or unresolved questions?
- What evidence is needed before each assumption becomes a decision?
- Who owns product, technical, data, legal, security, and operational approval?

### Work

1. Consolidate the complete conversation and existing notes.
2. Create requirement IDs and source references.
3. Create the glossary, decision register, assumption log, dependency map, and risk register.
4. Establish document states, version control, approval authority, and change-control rules.
5. Select the launch country, city, first vertical, pilot size, and legal/commercial operator.
6. Create a research plan for merchants, staff, customers, payment providers, and channel-platform prerequisites.

### Gate G0

- Requirements are traceable and contradictions are visible.
- Critical launch decisions have owners and due dates.
- No implementation assumption is disguised as an approved decision.

## Stage 1 — Market, business and operating model

### Questions

- Whose problem is Ubi solving first?
- Who pays, who uses, who administers, who supports, and who bears risk?
- What must happen outside software for onboarding and adoption to succeed?

### Work

1. Interview representative owners, receptionists, service staff, and customers.
2. Observe real bookings from calls, walk-ins, social messages, notebooks, and payment evidence.
3. Quantify lost enquiries, double-booking, response time, no-shows, deposit disputes, and training cost.
4. Define the Ubi value proposition and first vertical’s operating model.
5. Define assisted onboarding, merchant verification, customer support, escalation, and account recovery.
6. Model pricing, messaging/payment costs, support cost, acquisition cost, margins, and vendor exposure.
7. Define success and explicit stop/pivot criteria for the pilot.

### Gate G1

- At least one narrow market has a verified painful problem and willingness to adopt.
- The business can be supported operationally at a plausible cost.
- The first release has a measurable value hypothesis.

## Stage 2 — Actors, identity, organization and authority

### Questions

- What is a person, account, customer, organization, brand, branch, booking operation, staff membership, and public profile?
- Can one person manage several businesses? Can one business operate several booking models?
- Who may view, create, approve, cancel, refund, export, delete, connect, or impersonate?

### Work

1. Model guest and verified customer identities.
2. Model business organizations, brands, branches, departments, and booking operations.
3. Model staff membership and capability-based permissions.
4. Model Ubi support access, break-glass access, customer consent, and audit.
5. Define phone verification, account recovery, device/session management, and high-risk step-up verification.
6. Define stable IDs, public usernames, mutable slugs, and business subdomains.
7. Create the responsibility and data-controller/processor matrix.

### Gate G2

- Every first-release action has an authorized actor and audit rule.
- Multi-business and multi-operation ownership is unambiguous.
- Account recovery cannot silently transfer business ownership.

## Stage 3 — Business taxonomy and workflow architecture

### Questions

- Which booking archetypes cover the target businesses?
- Where must workflows differ, and where can they share components?
- What is a booking, request, quote, hold, reservation, appointment, rental, stay, or class?

### Work

1. Map every target business to reusable booking archetypes.
2. Document standard, request-to-confirm, and custom quotation journeys.
3. Model availability, capacity, staff, venue, asset, room, service-area, and multi-resource constraints.
4. Model pricing, packages, add-ons, taxes, discounts, deposits, payment evidence, refunds, and balances.
5. Model calls, walk-ins, connected-channel enquiries, Ubi messages, and customer web bookings.
6. Define holds, expiry, waitlists, recurrence, rescheduling, cancellation, no-show, completion, and dispute behaviour.
7. Create state machines, sequence diagrams, invariants, and exception catalogues.

### Gate G3

- Every first-release scenario can be described without UI or database terminology.
- No two states have ambiguous financial or scheduling meaning.
- Double-booking and false-confirmation protections are explicit.

## Stage 4 — Requirements and experience architecture

### Work

1. Write functional requirements by actor and capability.
2. Write measurable quality attributes: performance, offline behaviour, availability, durability, accessibility, localization, security, privacy, and cost.
3. Create task inventories and information architecture for Android, desktop, customer web/app, and Ubi operations.
4. Create low-fidelity flows before visual styling.
5. Test language, comprehension, touch targets, error recovery, calendar presentation, and low-bandwidth behaviour.
6. Create a shared design system with white surfaces, purple primary actions, sky-blue informational actions, and green restricted to success/approval.
7. Define content design and translations for English and the selected launch languages.

### Gate G4

- Representative users complete critical prototype tasks with acceptable error/help rates.
- Every requirement has an observable acceptance condition.
- Platform responsibilities do not depend on a specific screen layout.

## Stage 5 — Domain and platform architecture

### Work

1. Define bounded contexts and ownership of identity, organizations, catalog, availability, booking, conversation, payment, notification, audit, and analytics rules.
2. Produce system-context, container, component, deployment, and trust-boundary diagrams.
3. Select modular-monolith boundaries for the first release and rules for later extraction.
4. Define synchronous APIs, asynchronous events, idempotency, versioning, compatibility, rate limits, and error contracts.
5. Define external adapters for messaging, payments, SMS/push, media, search, and calls.
6. Define environment, configuration, feature-flag, secret, and tenant-isolation strategies.
7. Record all build-versus-buy and vendor-exit decisions.

### Gate G5

- Component responsibilities and data ownership do not overlap ambiguously.
- Failure and timeout behaviour is defined across every boundary.
- The architecture supports one release without premature distributed-system complexity.

## Stage 6 — Database and data-system architecture

### Work

1. Create the conceptual data model from the approved domain model.
2. Create logical schemas, entities, relationships, lifecycle fields, audit fields, and tenant ownership.
3. Define primary/foreign keys, unique/check/exclusion constraints, indexes, locking, isolation, concurrency, and idempotency records.
4. Define temporal data, timezones, Ethiopian/Gregorian presentation, currency precision, phone identity, and multilingual fields.
5. Define local Android storage, sync commands, server versions, tombstones, conflicts, and recovery.
6. Define PostgreSQL topology, connection pooling, migrations, backups, point-in-time recovery, replicas, maintenance, and restore tests.
7. Define cache, rate-limit, search, object/media, secret, and ephemeral storage boundaries.
8. Define event collection, analytics model, metric dictionary, data quality, lineage, warehouse/lake timing, and access controls.
9. Define classification, minimization, retention, deletion, export, legal hold, and anonymization/pseudonymization.

### Gate G6

- Scheduling and financial invariants are enforced below the UI.
- Restore, deletion, export, and tenant-isolation tests are designed.
- Operational and analytical workloads cannot silently corrupt each other.

## Stage 7 — Security, privacy, policy and legal architecture

### Work

1. Threat-model accounts, APIs, mobile/desktop/web clients, integrations, payments, support tools, calls, media, and supply chain.
2. Define authentication, authorization, encryption, key/secret management, session/device security, abuse controls, and audit.
3. Create privacy data-flow maps, purposes, lawful-basis questions, consent records, retention, access/correction/deletion/export processes, and child/minor handling.
4. Map launch-country company, consumer, electronic transaction, data protection, communications, payment, tax, and recordkeeping obligations.
5. Define platform permissions, OAuth scopes, credential rotation, webhook authenticity, and merchant disconnect/deletion.
6. Define business/customer terms, privacy notices, data-processing terms, acceptable-use policy, messaging/call consent, cookie notice, refund/cancellation responsibilities, subprocessor list, and law-enforcement process.
7. Obtain licensed local counsel review before production launch.

### Gate G7

- High-risk threats have controls, owners, tests, and residual-risk decisions.
- Data processing has documented purpose, authority, retention, and deletion.
- No production integration relies on a merchant sharing a social-media password.
- Country-specific legal review is complete for the release boundary.

## Stage 8 — Reliability, operations and engineering system

### Work

1. Define service-level indicators/objectives, error budgets, capacity assumptions, and cost budgets.
2. Define logs, metrics, traces, crash reporting, privacy filtering, alerting, dashboards, and audit review.
3. Define unit, domain, property, integration, contract, migration, UI, accessibility, performance, security, offline, chaos/failure, and recovery testing.
4. Define repository structure, coding standards, dependency governance, code review, CI/CD, signing, environments, release channels, rollback, and feature flags.
5. Define merchant support, incident response, privacy breach response, integration outage, payment reconciliation, backup restore, and business continuity runbooks.
6. Define documentation, training, sandbox/demo accounts, support SLAs, and escalation routes.

### Gate G8

- The release can be observed, supported, rolled back, restored, and explained.
- Critical runbooks have been exercised, not merely written.
- Performance and cost can be measured on representative low-cost devices and weak networks.

## Stage 9 — End-to-end validation and release contract

Current drafting position (2026-10-04): experience 23–26, validation/release control 71–73 and business execution 80–85 exist. Scenarios are not run; the candidate contract is not frozen. Specification 72 distinguishes synthetic feasibility/research from real-data/live release, so necessary prototypes can produce evidence without bypassing privacy/legal/production gates.

### Work

1. Walk representative scenarios across people, software, providers, money, data, and failure states.
2. Prototype and field-test the first release with realistic merchants and customers.
3. Resolve contradictions and update architecture decision records.
4. Freeze the first-release actor/capability matrix, workflows, contracts, schemas, quality budgets, policy documents, and operational obligations.
5. Create the implementation backlog with dependency order, estimates, owners, acceptance tests, risks, and release gates.

### Gate G9

- The first release contract is approved by product, technology, data, security/privacy, operations, and business leadership.
- Remaining unknowns are explicitly accepted, deferred, or blocking.

## Stage 10 — Work & Build

Implementation proceeds as thin, deployable vertical slices:

1. Engineering foundations, CI, environments, observability, and contract-test harness.
2. Identity, organization, membership, permission, and audit slice.
3. Catalog, availability, and manual booking slice.
4. Android offline schedule and safe synchronization slice.
5. Customer booking and manage-booking link slice.
6. Desktop operations slice.
7. Notifications and reminders slice.
8. Connected-conversation and human-handoff slice.
9. Deposit/payment and reconciliation slice.
10. Production hardening, migration, support, and pilot release.

Each slice must include domain rules, database migration, APIs/events, client behaviour, permissions, audit, privacy, observability, documentation, tests, rollback, and operational ownership. A backend endpoint or attractive screen alone is not a completed slice.

## Stage 11 — Pilot and controlled expansion

1. Assisted onboarding and shadow operation.
2. Gradual migration of real bookings.
3. Daily review of failed bookings, sync conflicts, support requests, integration errors, and payment discrepancies.
4. Weekly product and operating review against the pilot scorecard.
5. Expansion only after explicit retention, reliability, comprehension, and unit-economic gates.

## Stage 12 — Future-platform governance

Native customer applications, general username discovery, in-app calls, marketplace discovery, hotels, healthcare, multi-country payments, and AI automation each require a separate approval package. In-app messaging and calling are required platform direction, not discarded optional ideas; specification 85 retains explicit workstreams and validation gates for them. Future data fields or interface placeholders do not authorize premature implementation.

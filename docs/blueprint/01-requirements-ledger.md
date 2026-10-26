# Ubi requirements ledger

Status: **Draft — founder input consolidated; validation pending**  
Notation: `MUST` is mandatory, `SHOULD` is a strong default, and `MAY` is optional. These terms indicate priority, not legal enforceability.

## Requirement states

- **Stated** — explicitly requested by the founder.
- **Derived** — necessary consequence of a stated requirement.
- **Proposed** — expert recommendation requiring approval.
- **Validated** — supported by research or testing.
- **Approved** — authorized for a release.

## Product and market

| ID | Requirement | State | Release implication |
|---|---|---|---|
| PRD-001 | Ubi MUST serve East African operating realities rather than copy a foreign booking product unchanged. | Stated | Requires country overlays, local research, local payments, language and network testing. |
| PRD-002 | Ubi MUST remain understandable to people with below-average experience using booking applications. | Stated | Requires comprehension testing, large targets, ordinary language, progressive disclosure and assisted recovery. |
| PRD-003 | Ubi MUST support small businesses and have an architecture that can later support larger organizations. | Stated | First release remains narrow; organization/branch/operation boundaries must not block growth. |
| PRD-004 | Ubi MUST begin with the business-side Android application. | Stated | Android is the first production client. |
| PRD-005 | Ubi SHOULD provide a desktop business workspace after Android for larger operations. | Stated | Desktop workflow research is required before choosing browser, installable web or native form. |
| PRD-006 | The first launch MUST be one country, one city and a deliberately selected vertical before regional expansion. | Proposed | Country/city/vertical are blocking decisions for final legal, payment and localization design. |
| PRD-007 | Ubi MUST start as business booking infrastructure, not a consumer marketplace. | Derived | Merchant-supplied demand provides immediate value and avoids marketplace cold start. |
| PRD-008 | Hotels and regulated healthcare SHOULD remain outside the initial release. | Proposed | Both require specialized regulatory/inventory/reliability analysis. |

## Business identity and organization

| ID | Requirement | State | Release implication |
|---|---|---|---|
| ORG-001 | A human MUST have one login identity that can hold memberships in one or more organizations. | Derived | Login identity cannot be the tenant itself. |
| ORG-002 | One organization MUST support multiple brands, branches and booking operations. | Stated | Example: hotel rooms and event halls share an organization but use distinct booking rules. |
| ORG-003 | Each booking operation MUST have its own workflow, catalog, resources, availability, policies and public entry point where required. | Stated | Operations are configurable bounded units within the organization. |
| ORG-004 | Businesses MUST have a stable internal ID and MAY have a changeable public username/slug. | Stated/Derived | Foreign keys never depend on a username or domain slug. |
| ORG-005 | Business public pages SHOULD use business subdomains. | Stated | Certificate, slug collision, reserved-name, abuse and domain-change policies are required. |
| ORG-006 | Staff access MUST be membership- and permission-based, not shared credentials. | Derived | Requires invites, roles, revocation, session control and audit. |
| ORG-007 | An organization MUST be able to constrain staff to a branch, operation or assigned resource. | Derived | Authorization scope is more granular than organization-wide access. |

## Booking engine

| ID | Requirement | State | Release implication |
|---|---|---|---|
| BKG-001 | Ubi MUST support regular/standard and custom booking workflows. | Stated | They require separate lifecycle states and confirmation semantics. |
| BKG-002 | Standard booking flows MUST be tailored to the booking archetype and business configuration. | Stated | One universal form is prohibited. |
| BKG-003 | Custom requests MUST support requirements gathering, business review, quotation, customer acceptance, holds and deposit rules. | Derived | A custom request is not automatically a booking. |
| BKG-004 | Ubi MUST support appointments, service events, venues/resources, stays, classes, rentals, field services, viewings, consultations and quotations as platform archetypes over time. | Derived | First release selects only required archetypes. |
| BKG-005 | Businesses MUST be able to create bookings arising from phone calls and walk-ins quickly. | Stated | Manual booking is a first-release capability. |
| BKG-006 | Booking source MUST be recorded without changing the core booking lifecycle. | Derived | Sources include Ubi, web, call, walk-in, staff and connected channels. |
| BKG-007 | Ubi MUST prevent or explicitly resolve conflicting exclusive resource reservations. | Derived | Requires transactional constraints and conflict semantics. |
| BKG-008 | Requested, held, accepted, confirmed, completed, cancelled, declined, expired, no-show and rescheduled meanings MUST be unambiguous. | Derived | State machines and invariants precede schema/UI design. |
| BKG-009 | A business MUST choose whether eligible bookings auto-confirm or require staff acceptance. | Derived | Confirmation policy belongs to operation/service configuration. |
| BKG-010 | Customers and staff MUST be able to edit prefilled data before final confirmation. | Stated | Conversation extraction is assistance, not unquestionable truth. |

## Conversation, channel and handoff

| ID | Requirement | State | Release implication |
|---|---|---|---|
| MSG-001 | Social-channel connection and message management MUST remain optional. | Stated | Direct Ubi link and manual booking work without channel integration. |
| MSG-002 | Ubi SHOULD unify messages received through officially supported business integrations. | Stated | Integration capability differs by platform and account eligibility. |
| MSG-003 | Ubi MUST use official authorization and MUST NOT request a merchant’s social-media password. | Derived | OAuth/embedded signup, minimum scopes, token protection and revocation are mandatory. |
| MSG-004 | Staff MUST be able to see new supported conversations, reply, assign, pause automation and take over. | Stated | Conversation ownership and bot/human state are first-class. |
| MSG-005 | Automation MUST escalate uncertainty and MUST NOT silently confirm unavailable capacity, unverifiable payment or irreversible changes. | Derived | Policy engine, confidence boundaries and audit are required. |
| MSG-006 | A secure contextual link SHOULD continue a conversation into Ubi or the customer web experience with editable prefill. | Stated | Use opaque, scoped, expiring, revocable server-side handoff tokens. |
| MSG-007 | The link SHOULD open the installed app and MUST fall back to web when unavailable. | Stated | Verified HTTPS App Links/universal-link strategy is required. |
| MSG-008 | Ubi MUST distinguish external-channel identity from Ubi identity and support safe linking/merging. | Derived | Phone numbers, platform-scoped IDs and guest/customer profiles cannot be treated as identical. |
| MSG-009 | In-app messaging MUST be included in the long-term system and SHOULD begin as booking-linked conversations. | Stated/Proposed | Messaging data and safety boundaries are planned now; broad social discovery is later. |
| MSG-010 | In-app calls MUST be included in long-term architecture but released only after identity, messaging, safety and network quality are proven. | Stated/Proposed | Call signaling/media/vendor/storage/consent remain separate approval scope. |

## Customer experience

| ID | Requirement | State | Release implication |
|---|---|---|---|
| CUX-001 | A customer MUST be able to complete an ordinary web booking without installing an app. | Proposed from prior analysis | Native customer app is not a prerequisite for merchant value. |
| CUX-002 | Ordinary booking MUST NOT require email or password creation. | Derived | Phone-first guest flow with later identity promotion. |
| CUX-003 | The customer MUST see service, date/time, price/deposit, location and policy before confirming. | Derived | Review step is mandatory. |
| CUX-004 | Customers SHOULD manage cancellation/rescheduling through a private link without a password. | Proposed | Requires scoped token security, expiry and step-up verification for risk. |
| CUX-005 | Human call or message fallback MUST be understandable wherever automation cannot complete the task. | Stated/Derived | Contact policy and business availability are required. |

## Business experience

| ID | Requirement | State | Release implication |
|---|---|---|---|
| BUX-001 | The first-use experience MUST use familiar business-type templates and suggested services. | Stated | Avoid blank enterprise configuration. |
| BUX-002 | Assisted onboarding and staff training SHOULD be supported for early merchants. | Stated | Requires onboarding console, checklist, demo/test booking and audit of support access. |
| BUX-003 | Owners MUST be able to configure services, prices, duration, working hours, exceptions, staff and booking policy. | Derived | First-release scope is selected by vertical. |
| BUX-004 | Daily operations MUST prioritize today’s work, upcoming bookings, requests requiring action, sync state and one fast add-booking action. | Derived | Information architecture follows frequency and urgency. |
| BUX-005 | Larger organizations MUST be able to operate through an efficient desktop workspace. | Stated | Keyboard, density, multiple columns, bulk actions and multi-branch workflows require separate UX. |

## Performance, offline and reliability

| ID | Requirement | State | Release implication |
|---|---|---|---|
| QLT-001 | Ubi MUST feel very fast and smooth on Android, iOS and desktop over time. | Stated | Concrete device/network/startup/frame/API budgets must be approved. |
| QLT-002 | Android MUST expose locally cached useful content without waiting for the network. | Derived | Local-first reads and durable local writes are required. |
| QLT-003 | A dropped connection MUST NOT silently lose an accepted user change. | Stated | Durable outbox, visible sync state and retry semantics are required. |
| QLT-004 | Retried writes MUST be idempotent. | Derived | Idempotency keys and server command records are required. |
| QLT-005 | Financial and scheduling conflicts MUST NOT use silent last-write-wins resolution. | Derived | Explicit merge or authoritative rejection paths are required. |
| QLT-006 | Representative low-cost devices and constrained networks MUST be part of performance testing. | Derived | Simulator-only performance approval is prohibited. |
| QLT-007 | Reliability MUST include observable SLOs, tested backups, restore targets and incident runbooks. | Derived | “Cloud hosted” is not a recovery plan. |

## Accessibility, localization and visual system

| ID | Requirement | State | Release implication |
|---|---|---|---|
| A11Y-001 | Workflows MUST use ordinary, locally understandable language and one primary decision at a time. | Stated | Content design is part of functional correctness. |
| A11Y-002 | Interactive targets MUST be at least 48 dp on touch clients unless a stricter platform rule applies. | Proposed | Component tests and design tokens enforce the minimum. |
| A11Y-003 | Meaning MUST NOT depend on color alone. | Derived | Text/icon/state labels accompany color. |
| A11Y-004 | White is the main surface, purple is the primary brand/action accent, sky blue is secondary/informational, and green is reserved for genuine approval/success. | Stated | Semantic design tokens must preserve accessible contrast. |
| A11Y-005 | The product MUST support font scaling, screen readers and keyboard use where applicable. | Derived | Accessibility testing is required on every client. |
| LOC-001 | Language and calendar architecture MUST be localizable from the beginning. | Stated/Derived | Strings, layouts, sorting and date formatting cannot be English-hard-coded. |
| LOC-002 | Timed records MUST store unambiguous instants/timezones while permitting Gregorian and Ethiopian presentation. | Derived | Calendar conversion is an edge/presentation concern with tested rules. |

## Payments and money

| ID | Requirement | State | Release implication |
|---|---|---|---|
| PAY-001 | Businesses MUST be able to configure no deposit, fixed deposit, percentage deposit or full prepayment when supported. | Stated/Derived | Amount snapshots and policy versioning are required. |
| PAY-002 | Cash and pay-at-location MUST remain supported. | Stated | Digital payment is not mandatory for every booking. |
| PAY-003 | Customer-reported payment and provider-verified payment MUST be distinct. | Derived | Screenshots cannot directly produce verified state. |
| PAY-004 | Payment, booking and fulfilment state machines MUST be separate. | Derived | Payment failure does not erase the booking history. |
| PAY-005 | Currency values MUST use exact minor-unit or currency-aware decimal representation, never binary floating point. | Derived | Currency metadata and provider reconciliation rules are required. |
| PAY-006 | Refunds, reversals, disputes, reconciliation and webhook retries MUST be designed before real-money release. | Derived | Payment integration includes operations, not just checkout. |

## Security, privacy, permission and policy

| ID | Requirement | State | Release implication |
|---|---|---|---|
| SEC-001 | Access MUST be least-privilege and tenant-scoped on the server. | Derived | Client routing or subdomains never authorize access. |
| SEC-002 | Sensitive administrative, ownership and financial actions MUST use step-up verification where risk warrants it. | Proposed | Risk matrix determines exact operations. |
| SEC-003 | Secrets and provider tokens MUST be encrypted, access-controlled, rotated and auditable. | Derived | Tokens are never exposed to ordinary client code or logs. |
| SEC-004 | Booking, permission, payment and support-access changes MUST produce tamper-evident audit events appropriate to risk. | Derived | Audit access and retention must also be controlled. |
| PRV-001 | Ubi MUST collect the minimum personal data needed for stated purposes. | Derived | Every data field requires purpose, owner and retention. |
| PRV-002 | Private handoff and manage-booking links MUST be scoped, expiring and revocable. | Derived | Personal/conversation data is not embedded in URL query parameters. |
| PRV-003 | Users MUST have documented access, correction, export and deletion processes subject to jurisdiction and legitimate retention. | Derived | Product, support, backups and analytics must implement the process together. |
| PRV-004 | Support access to merchant/customer data MUST be purpose-limited, time-bounded where possible and audited. | Derived | No invisible universal support access. |
| LEG-001 | Final terms, privacy, payment, communications and data-processing documents MUST be reviewed by licensed counsel for each launch jurisdiction. | Proposed | Product architecture work is not a substitute for legal advice. |

## Data and analytics

| ID | Requirement | State | Release implication |
|---|---|---|---|
| DAT-001 | PostgreSQL SHOULD be the initial authoritative operational database unless evidence disproves the fit. | Proposed | Final selection follows workload and operational review. |
| DAT-002 | Operational, audit, analytical, cache, search, object and secret data MUST have explicit ownership and purpose. | Derived | “Database” is not one undifferentiated store. |
| DAT-003 | Analytics MUST use a governed event and metric dictionary. | Derived | Reports cannot redefine booking or revenue inconsistently. |
| DAT-004 | Data lineage, quality checks and deletion propagation MUST cover derived analytical data. | Derived | A warehouse cannot become an ungoverned permanent copy. |
| DAT-005 | Schema changes MUST be versioned, reviewed, backward-compatible during deployment and rehearsed for recovery. | Derived | Migrations are release artifacts. |

## Items deliberately not yet approved

- Launch city and first vertical. Launch country is decided: Ethiopia, recorded in DEC-001/ADR-001.
- Legal entity and contracting structure.
- Exact first-release applications beyond Android-first and desktop-next preference.
- Technology stack, cloud, database vendor, payment provider and communication vendors.
- Exact role list and permission boundaries.
- Exact booking states and database schema.
- Exact retention periods and data-residency commitments.
- Languages, currencies and calendar defaults for the launch.
- Pricing, subscription tiers and transaction fees.
- AI/automation behaviour and model/vendor selection.
- Native voice-call architecture and recording policy.

These are inputs to the decision programme, not omissions to be guessed during coding.

## Current validation linkage

Updated 2026-10-04: all 79 ledger requirements have planned coverage in `specifications/71-end-to-end-validation-and-traceability.md`. This is traceability, not passed tests or requirement approval. The proposed release boundary and missing evidence are specification 72; dependency-ordered work is specification 85. Native messaging and calling remain required platform direction, with explicit delivery workstreams and separately reviewed release timing.

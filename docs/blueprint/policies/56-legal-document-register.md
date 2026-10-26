# Legal document register and product-policy governance

Status: **Draft — structure only; not legal advice and not launch-approved wording**  
Depends on: launch country, contracting entity, product scope, data roles, payment role, channel/call/AI scope and licensed local counsel  
Blocks: public launch, paid subscriptions, real payments, connected channels, native calls and production personal-data processing

## Purpose

Ubi needs a controlled set of contracts, notices, consents, public rules and internal procedures. A generic terms page or copied privacy policy is not sufficient. Each document must match the actual product, data flow, responsible legal entity, country, language and operational capability.

This register defines what must exist, who it governs, where it appears, what evidence Ubi retains and what decisions must be made before counsel can finalize it.

## Document types are not interchangeable

| Type | Function | User action | Example |
|---|---|---|---|
| Contract | Creates enforceable obligations between identified parties | Acceptance/signature when legally required | Merchant service agreement |
| Notice | Explains facts, rights and practices | Must be available; acknowledgment is not automatically consent | Privacy notice |
| Consent | Freely given, specific, informed and revocable permission where consent is the chosen lawful basis | Clear affirmative choice | Optional marketing or call recording |
| Policy/rule | Defines permitted conduct and consequences | Incorporated into contract where enforceability is required | Acceptable-use policy |
| Operational procedure | Tells authorized staff how to perform controlled work | Training and access approval | Data-rights procedure |
| Evidence record | Proves what happened without replacing the governing document | System-generated | Terms version acceptance event |

Ubi must never label a forced contractual necessity as optional consent, or use one checkbox for unrelated purposes.

## Controlled legal-document register

### Public and user-facing documents

| ID | Document | Principal audience | Product surface/evidence | Minimum dependency |
|---|---|---|---|---|
| LEG-DOC-001 | Ubi Business service agreement/terms | Owner and authorized merchant representative | Signup/contract flow; signer, authority representation, version, language, time and evidence | Entity, country, service/pricing/liability model |
| LEG-DOC-002 | Customer booking terms | Guest and account customer | Before final booking/request; version attached to booking | Platform-versus-merchant role, cancellation/refund/dispute allocation |
| LEG-DOC-003 | Merchant-specific service, cancellation, rescheduling, no-show and refund terms | Customer | Service page and review step; policy-version snapshot | Merchant configuration and consumer-law review |
| LEG-DOC-004 | Account privacy notice | Owners, staff and account customers | Signup/settings; notice version presented | Data roles, inventory, retention, transfer and rights process |
| LEG-DOC-005 | Booking/contact privacy notice | Guest/customer whose data is entered by self, merchant or connected channel | First practical collection point plus durable link | Controller/processor allocation and indirect-collection duties |
| LEG-DOC-006 | Cookie/local-storage notice and preference control | Web users | Web entry/settings; preference record where required | Actual cookie/SDK inventory and lawful basis |
| LEG-DOC-007 | Acceptable-use and prohibited-business policy | Merchants, staff and customers | Onboarding/help/report flow; enforcement evidence | Launch law, platform rules and risk appetite |
| LEG-DOC-008 | Messaging, notification and connected-channel notice | Message participants | Thread/channel connection and notification settings | Channel capability, retention and opt-out design |
| LEG-DOC-009 | Calling and recording notice/consent | Call participants | Before call/recording; separate proof for recording | DEC-056, telecom/privacy law and vendor architecture |
| LEG-DOC-010 | Payment, deposit and payment-status disclosure | Payer and merchant | Checkout/payment record; provider and policy version | Payment role, provider, verification and refund model |
| LEG-DOC-011 | Safety, reporting, blocking, moderation and appeal rules | All participants | Messaging/call/profile/report surfaces | Trust-and-safety operating model |
| LEG-DOC-012 | Accessibility and support statement | All users | Public help and app settings | Supported channels/languages and response capability |
| LEG-DOC-013 | Subprocessor/provider list and material-change process | Contracting merchants and other entitled parties | Public/contract portal; dated versions and notifications | Approved vendor register and transfer assessment |
| LEG-DOC-014 | Data-rights request information | Data subjects/authorized agents | Privacy centre/support | Identity, authorization, exceptions and response process |
| LEG-DOC-015 | Law-enforcement/request guidelines and transparency approach | Authorities and public | Controlled public page when operationally ready | Entity, service address and procedure in document 57 |

### Business-to-business and provider documents

| ID | Document | Parties/purpose | Minimum dependency |
|---|---|---|
| LEG-DOC-020 | Data-processing agreement and schedules | Ubi and merchant where one processes for the other | Per-activity role analysis, instructions, security, subprocessors, transfers, deletion and assistance |
| LEG-DOC-021 | Commercial order form/subscription schedule | Ubi and merchant | Currency, tax, tier, limits, renewal, suspension and termination |
| LEG-DOC-022 | Service-level/support schedule | Ubi and eligible merchant | Measurable service boundaries, exclusions, channels and remedies |
| LEG-DOC-023 | Vendor/subprocessor DPA and security schedule | Ubi and provider | Document 55 review, locations, subprocessors, incident terms and exit |
| LEG-DOC-024 | Channel/provider supplemental terms | Ubi, merchant and/or official provider | Provider relationship, permissions, restrictions and disconnect consequences |
| LEG-DOC-025 | Payment-provider/merchant supplemental agreement | Applicable parties | Settlement/custody/fees/disputes/KYB allocation confirmed by counsel/provider |
| LEG-DOC-026 | API/developer agreement | Future external integrator | Approved public API, credentials, data use, rate limits and revocation |
| LEG-DOC-027 | Pilot agreement | Pilot merchant | Pilot scope, support, data, pricing, feedback, exit and no-production assumptions |

### Workforce and internal control documents

| ID | Document/procedure | Minimum content |
|---|---|---|
| LEG-DOC-030 | Employee/contractor confidentiality and data-handling terms | Confidentiality, least privilege, approved systems, return/deletion and survival |
| LEG-DOC-031 | Security and acceptable-use policy | Identity, devices, secrets, access, reporting and disciplinary path |
| LEG-DOC-032 | Privacy/data-governance policy | Roles, inventory, purpose, minimization, retention, rights and transfers |
| LEG-DOC-033 | Incident and breach procedure | Triage, evidence, decision authority, notification clocks and communication |
| LEG-DOC-034 | Data-rights operating procedure | Intake, verification, search, exemptions, approval, delivery and audit |
| LEG-DOC-035 | Retention/deletion/legal-hold schedule | Record class, trigger, period, authority, deletion and hold override |
| LEG-DOC-036 | Law-enforcement/civil-request procedure | Authenticity, jurisdiction, legal basis, narrowness, notice and disclosure approval |
| LEG-DOC-037 | Support and privileged-access procedure | Just-in-time scope, approval, monitoring, emergency access and revocation |
| LEG-DOC-038 | Moderation/enforcement procedure | Evidence, proportional action, escalation, appeal and vulnerable-user safety |
| LEG-DOC-039 | Vendor procurement/renewal/exit procedure | Document 55 lifecycle, re-review and deletion evidence |
| LEG-DOC-040 | Records and audit procedure | Document 57 events, access, export, integrity and disposal |

## Product-to-document mapping

Every relevant product change must answer:

1. Which party provides the service and which party promises what?
2. What data is collected, from whom, for what purpose and under whose control?
3. Which contract, notice, choice or policy applies at that exact moment?
4. Can a low-literacy user understand the immediate consequence before acting?
5. What version and evidence must be attached to the booking/payment/account record?
6. Can the user later find the document, change an optional choice or exercise a right?
7. What happens offline, after a link expires, after policy change or after account termination?

### Required evidence model

An acceptance/notice/consent event should record only what is necessary, such as:

- actor or defensible guest/session reference;
- organization/role when acting for a business;
- document ID, immutable version, language and content checksum;
- action (`presented`, `accepted`, `declined`, `withdrawn`, `superseded`);
- time, product surface and relevant transaction reference;
- collection method and legally justified supporting metadata;
- parent/guardian or representative authority evidence where applicable.

Do not store full access tokens, payment credentials, message bodies or unnecessary device fingerprinting in legal evidence.

## Plain-language and mistake-resistant presentation

- Put the immediate operational consequence next to the action; do not hide it only in a long policy.
- Use short sentences, familiar local wording, large touch targets and one primary decision per screen.
- Translate professionally and review meaning in context; a machine translation is not legal approval.
- Explain essential versus optional choices separately.
- Do not preselect optional marketing, recording or unrelated data-sharing choices.
- Provide a simple summary with access to the complete authoritative text; the summary must not contradict it.
- Use audio/assisted explanation only as a supplement, with privacy-safe replay and equivalent text.
- Do not force account creation merely to read applicable terms or privacy information.
- Display merchant policy changes before they apply to an existing booking unless law/contract permits otherwise.

## Versioning, hierarchy and change control

Each controlled document needs:

- unique document/version ID, owner, status, jurisdiction, language and audience;
- drafting basis, counsel reviewer, approval date, effective date and next review;
- canonical signed/approved source and checksum;
- change summary and affected product/data/workflows;
- prior versions retained for disputes and evidence under the retention schedule;
- delivery/notification and re-acceptance rule;
- rollback/correction process for publication error.

Order of precedence must be explicit in the merchant agreement. A help article, chatbot response, salesperson statement or UI tooltip cannot silently amend the contract.

Material changes to purpose, price, liability, dispute route, provider/data transfer or user rights require legal review and proportionate notice. Re-acceptance is used only when required or genuinely appropriate, not as a substitute for fair terms.

## Merchant-authored policies

Ubi may provide reviewed templates and guided configuration, but the merchant remains responsible for accurate service descriptions, availability, eligibility, cancellation/refund rules and lawful content where allocated by contract.

Controls:

- structured choices before unrestricted free text;
- plain-language preview of what the customer will see;
- effective date and immutable snapshot attached to each booking;
- warnings for internally contradictory settings;
- restricted terms and review/escalation for risky categories;
- change history, role permission and approval for material edits;
- no claim that a template guarantees legal compliance in every country/business type.

## Special scope gates

### Minors and sensitive services

Ubi must not infer one regional age threshold. Before any child-directed or likely-minor service, counsel must define age/authority, parental or guardian evidence, marketing limits, photography/media permissions, rights handling and safe support. Sensitive healthcare is outside the first-release recommendation.

### Messaging, automation and AI

Users must know when they interact with automation, how to reach a human, what connected channel is involved and which party sent a message. AI-specific terms cannot authorize undisclosed model training or override privacy-purpose limits. Deterministic booking facts remain the authoritative source.

### Calls and recording

Calls do not imply recording. Recording is off unless separately approved. If enabled, the product must provide lawful notice/consent, visible state, stop control, retention, access, export/deletion handling and vendor/transfer analysis.

### Payments

Copy must distinguish customer-reported, pending/provider-processing and provider-verified payment. Ubi must not promise regulated custody, guaranteed settlement, refunds or buyer protection unless the approved provider/contract/operating model actually delivers them.

## Ownership and approval

| Responsibility | Accountable role |
|---|---|
| Product/workflow accuracy | Product owner |
| Data map and privacy architecture | Privacy/data lead |
| Security and evidence controls | Security/CTO |
| Pricing/payment/settlement facts | CEO/finance/payment owner |
| Merchant operations/support feasibility | Operations/support owner |
| Local legal interpretation and final text | Licensed counsel in each launch jurisdiction |
| Entity approval and risk acceptance | Authorized company leadership |
| Publication/version integrity | Designated legal-operations owner |

Counsel approval is necessary but not sufficient: product, engineering and operations must verify that published statements match runtime behavior.

## Release gates

- DEC-001, DEC-005, DEC-006, DEC-007 and applicable trust decisions are resolved.
- Every first-release workflow maps to the applicable contract, notice, choice, evidence and internal procedure.
- Licensed launch-country counsel reviews the complete product/data/payment/channel facts and final localized texts.
- Terms identify the correct entity, address/contact, governing framework, merchant/platform responsibility and dispute route.
- Privacy notices match the verified data inventory, vendor register, transfers, retention and user-rights implementation.
- Consent is used only when appropriate, granular and withdrawable without deceptive friction.
- Customer-facing terms and merchant policy snapshots are attached to bookings in a retrievable form.
- Accessibility, language comprehension and assisted-onboarding tests pass with representative users.
- Publication, checksum, version, notice, acceptance, withdrawal and rollback paths are tested.
- No marketing, salesperson, support bot or UI claim exceeds the approved documents and actual capability.

## Open decisions

- Launch country/city, legal entity and service address.
- First vertical, pilot arrangement, merchant/customer contracting relationships and dispute forum.
- Controller/processor allocation per activity and any representative/DPO obligations.
- Payment custody/settlement/provider role and tax treatment.
- First channels, message/call/recording/AI scope and notification lawful basis.
- Minor/guardian policy and excluded business categories.
- Languages, authoritative-text hierarchy and accessible alternatives.
- Retention periods, user-notice exceptions, arbitration/class-action/limitation clauses where lawful and appropriate.


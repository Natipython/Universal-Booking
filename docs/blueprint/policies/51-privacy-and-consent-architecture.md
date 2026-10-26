# Privacy and consent architecture

Status: **Draft — product architecture; legal roles, bases, notices and deadlines require launch-country counsel**  
Evidence review date: 2026-09-05  
Depends on: identity/organization, workflows, data inventory/lifecycle, tenant isolation and threat model  
Blocks: production collection, analytics/AI, channel/payment providers, final notices/agreements and launch

## Purpose

This document defines how Ubi makes privacy understandable and operational across Accounts, merchant Customer records, guest/scoped-link users, staff, connected channels, payments, messages, calls, devices, analytics and support.

It is not a privacy notice and does not select a legal basis. Those depend on the launch country, legal entity, exact feature and Ubi/merchant/provider responsibility.

## Privacy promise

Ubi should be able to answer, in plain language:

1. What information are we asking for?
2. Why is it needed now?
3. Is it required to complete the booking?
4. Which business and providers will receive/use it?
5. How long is it kept and what can the person do about it?
6. What happens if they decline or withdraw an optional choice?

Low digital literacy requires clearer explanation and fewer choices—not hidden consent or weaker protection.

## Foundational distinctions

- A global **Ubi Account** is not the same as a merchant’s **Customer record**.
- A phone number may be shared, recycled, mistyped or controlled by someone other than the service recipient.
- A Telegram/WhatsApp/Instagram/Facebook/TikTok identifier is provider/connection scoped and is not automatically a Ubi Account.
- A booker, payer, attendee, child, organizer and service recipient may be different people.
- A merchant may need a transaction record after an Account/customer profile closes, subject to law and purpose.
- Publicly visible social/profile information is not automatically free for unrelated collection or reuse.
- Consent to receive one booking update is not consent to marketing, recording, AI training or cross-merchant profiling.

## Responsibility-by-processing model

Do not assign one controller/processor label to the entire platform. Create a reviewed row for each processing activity.

Provisional questions—not legal conclusions:

| Activity | Likely responsibility question |
|---|---|
| Ubi Account authentication/security | Does Ubi independently determine essential purpose/means and act as controller? |
| Merchant customer/booking management | Is merchant controller and Ubi processor/service provider, and where does Ubi independently determine purposes? |
| Public business profile | Does merchant decide publication while Ubi separately controls platform abuse/security? |
| Ubi software billing/support/security | Ubi likely has independent purposes; exact role and notice need review |
| Connected social messages | Merchant, Ubi and platform may have distinct/overlapping roles under separate terms |
| Payment provider transaction | Merchant/Ubi/provider roles depend on who contracts, initiates, receives or controls funds/data |
| Product analytics and AI improvement | Separate Ubi purpose requiring explicit necessity, transparency and possibly choice/consent |
| Law-enforcement response | Ubi independently evaluates legal validity and disclosure duties for data it controls |

The responsibility register names legal entity, jurisdiction, data subjects, purpose, decision maker, processor instructions, providers/recipients, rights route, incident duties and contract/notice.

## Privacy data-flow method

For every release workflow, diagram:

```text
person/merchant/staff/provider source
→ client/form/chat/import collection point and notice
→ Ubi API/domain owner
→ database/object/device/search/cache/audit/analytics copies
→ authorized merchant/Ubi workforce recipients
→ messaging/payment/call/cloud/analytics subprocessors or independent recipients
→ country/region transfer
→ retention trigger and deletion/export/correction path
```

Document data fields, classification, purpose, legal-basis question, authority, encryption/access, retention, automated decision, child/sensitive data and failure behavior at each arrow.

Diagram: `diagrams/51-privacy-data-and-choice-flow.mmd`.

## Processing-purpose register

Every collection/use links to one approved stable purpose, such as:

- authenticate, secure and recover an Account;
- establish/administer merchant authority;
- answer an enquiry and create/fulfil/support a booking/request/quote;
- communicate a requested transactional update;
- verify/reconcile payment evidence;
- operate and troubleshoot a chosen integration;
- prevent/investigate fraud, abuse and security incidents;
- satisfy approved contractual/legal record duties;
- provide Ubi support under scoped access;
- measure reliability/usability using minimized telemetry;
- conduct optional merchant marketing where lawful and chosen.

A new purpose requires inventory, necessity/alternative review, role/basis/notice assessment, retention, rights and vendor update before collection or reuse.

## Lawful-basis decision record

Consent is not the default answer to all processing. For each activity, counsel determines applicable local bases and conditions, potentially including necessity for contract/pre-contract steps, legal obligation, vital/public interests, legitimate interests or consent depending on law.

The record contains:

```text
activity/purpose and responsible legal entity
jurisdiction and applicable law/version
candidate/approved basis and necessity test
data/subject/sensitive/child categories
balancing/impact assessment where relevant
notice and choice/withdrawal behavior
consequence if data is not provided
retention/recipients/transfers
approver/effective/review date
```

No engineer invents a basis in code. If the selected law does not recognize a proposed basis for that activity, product scope changes.

## Consent architecture

Where consent is required or chosen as the proper basis, it must be demonstrable and usable.

### Consent qualities

- affirmative, specific, informed and unambiguous under the applicable standard;
- separate from unrelated terms/choices;
- understandable in the selected language and accessible format;
- no pre-checked boxes, disguised button or unnecessary service denial;
- granular by purpose/channel where choices genuinely differ;
- easy to withdraw through a route comparable to giving it;
- withdrawal stops future processing and propagates without rewriting valid past evidence;
- parental/representative authority verified where required.

### Consent evidence

Store subject/representative, purpose, allowed/denied, notice/version/language, exact presentation context, affirmative action, source/device/session, time, authority evidence where necessary, withdrawal/supersession and propagation result. Do not store secret/biometric proof unnecessarily.

### Never bundle automatically

- transactional booking updates with marketing;
- Ubi marketing with merchant marketing;
- one messaging channel with all channels;
- call participation with recording;
- operating automation with training/model improvement;
- public business publication with public staff/customer data;
- one merchant relationship with cross-merchant identity/profile reuse.

## Notice and comprehension design

Use layered notices:

1. one short statement at collection explaining immediate purpose and recipient;
2. required versus optional fields visibly marked;
3. a clear action summary before confirmation;
4. expandable details for retention, rights, providers and responsibility;
5. accessible full notice and contact route.

Requirements:

- ordinary local language, tested with representative users;
- no legal jargon as the only explanation;
- icons never replace text;
- audio/help may supplement but not expose private content or become mandatory;
- user can review/correct before submitting;
- important changes are actively notified where required, not hidden in a footer;
- acceptance of Terms and privacy acknowledgement/consent are recorded distinctly.

## Data minimization by journey

### Public discovery

No login required for public service/business browsing. Collect only necessary security/rate telemetry. Do not force phone access/contacts/location when manual area selection works.

### Enquiry and availability

Ask service/date/branch/party/capacity facts required to answer. Avoid identity until a response/hold requires it.

### Booking/hold

Collect the minimum reliable contact and fulfilment facts. Additional questions must come from a versioned merchant form with purpose/classification/retention and prohibition of unsafe free text.

### Account/profile reuse

Offer convenience transparently. A guest booking does not require permanent profile creation unless the release contract proves necessity. Claiming an existing merchant Customer record requires verified linking, not contact similarity alone.

### Connected conversation

Ingest only events/content available through the approved business connection and needed for the configured workflow. Convert accepted booking facts into structured snapshots; do not keep full raw conversation indefinitely because it is convenient.

### Payment

Prefer provider-hosted collection and tokens. Store exact operational/reconciliation evidence, not payment credentials Ubi does not need.

## Sensitive and high-impact data

The pilot should exclude health, biometric, government-ID, precise sensitive location and similarly high-risk data unless the first vertical cannot function without it and receives a separate approval package.

Photography/wedding/event workflows require special caution:

- a booking may describe an event involving children or family details;
- booking permission does not grant Ubi or merchant broad image-publication rights;
- private galleries/media delivery, public portfolio use and face recognition are separate purposes;
- biometric/face recognition is out of initial scope;
- public testimonials/portfolio require separate evidence and withdrawal/contract rules.

## Children and minors

Age thresholds and parental authority vary by jurisdiction. Ubi therefore cannot ship one universal “under 18” checkbox as compliance.

Proposed safe pilot boundary:

- Ubi is not marketed/directed as a child social/booking service;
- Account creation and contracting/booking actions use an approved minimum age/authority model for launch country;
- an adult booker may make a service booking involving a child only when the merchant workflow legitimately needs that fact;
- collect minimal child information, normally no separate child Account/contact;
- child photos/health/school/location data and direct messaging/calling require additional review and safeguards;
- block profiling/marketing/AI training using child data by default;
- provide parent/guardian rights and safety escalation according to local law;
- age assurance must be proportionate and not collect more identity data than the risk justifies.

Document 52 and DEC-054 must be approved before a child-directed or sensitive vertical.

## Data-subject and merchant rights operations

Support access, correction, export/portability, deletion/erasure, restriction, objection/withdrawal and complaint/appeal as applicable.

Workflow requirements:

- accessible route for people without a Ubi Account;
- proportional identity/authority verification that does not disclose whether another account exists;
- route to responsible merchant versus Ubi while Ubi tracks its own duties;
- search all governed systems through data inventory/lineage;
- protect other persons, businesses, security/provider secrets and privileged material;
- secure expiring export and common machine-readable format where required/useful;
- exact jurisdictional deadline and exception calendar;
- task-level processor/provider propagation;
- reasoned partial/refusal decision and redress;
- deletion/restriction tombstone prevents delayed provider/offline resurrection;
- minimal completion evidence without retaining removed content.

## Account, Customer and transaction closure

- Account closure revokes sessions/authenticators and removes/minimizes Ubi Account data under schedule.
- Each merchant receives/acts on the relevant Customer-link/profile request according to responsibility and active duties.
- Booking/payment/audit evidence is not blindly deleted; retain/restrict/anonymize only as approved.
- Customer-account claim is severed when appropriate while historical party role remains understandable.
- Organization closure invokes export, provider disconnect, staff revoke, public unpublish, retention and deletion workflow.
- “Delete app” and “delete account/data” are clearly different.

## Cookies, SDKs, device permissions and telemetry

- Customer web uses necessary storage only until a reviewed cookie/analytics decision; non-essential tracking waits for applicable notice/choice.
- No advertising SDK/device fingerprinting in pilot by default.
- Android asks minimum runtime permissions at point of use; denial leaves a usable alternative where feasible.
- Crash/performance data is allowlisted/scrubbed/sampled and never includes message/contact/token payload by default.
- Push notifications hide sensitive booking/customer text on lock screens by default and offer granular controls.
- Clipboard/screenshots/OS backup are restricted only where risk justifies usability cost.

## Messaging, calls and marketing

- Transactional updates use the contact/channel requested or reasonably tied to the booking under the approved basis/rules.
- Marketing has separate merchant/Ubi identity, evidence, frequency and opt-out/suppression list.
- Opt-out does not necessarily suppress essential security/contract messages; copy explains the distinction.
- Connected-platform rules and customer initiation/template windows are enforced per current provider capability.
- In-app block/report/mute and merchant/staff takeover are privacy/safety controls.
- Call recording is off by default; if ever enabled, define participants/jurisdictions, notice/consent, purpose, access, retention, transcript/AI and refusal alternative.

## Automated assistance and AI

- Tell users when they are interacting with automation where material.
- Automation gathers only approved booking facts and exposes human takeover.
- Do not make solely automated high-impact decisions about payment, cancellation/refund, fraud restriction, ownership or sensitive eligibility without lawful design, explanation/review and human process.
- Store source/provenance/confidence for extracted suggestions.
- Prompt/model/provider logs, training use, retention and transfers are separate inventory/vendor decisions.
- Merchant/customer content is not used to train/improve general models by default.
- Prompt injection content cannot access secrets, other tenants or unapproved tools.

## Privacy engineering controls

- tenant and purpose-bound authorization with field-level DTOs;
- classification, encryption and least-privilege/JIT support access;
- data-protection defaults: unpublished/private, minimal retention, no marketing/recording;
- pseudonymous analytics IDs and minimum cohort controls;
- exact deletion/export/lineage orchestration across stores/providers;
- privacy-safe logs/errors/test data;
- DPIA/privacy-impact assessment before high-risk/sensitive/large-scale/AI/call features;
- records of processing and processor instructions/contracts;
- incident detection, triage and jurisdictional notification clock;
- periodic access/retention/vendor/consent review.

## Privacy incident handling

Document 67 will implement the runbook. Architecture must:

1. detect/contain without destroying evidence;
2. identify data, people/merchants/countries, sensitivity, volume, access and likely harm;
3. engage security/privacy/legal and responsible merchants/providers;
4. track each jurisdiction’s regulator/person/merchant notice threshold and deadline;
5. communicate accurately in understandable language;
6. remediate credentials/access/copies and monitor misuse;
7. record decision/evidence and improve controls.

No universal breach deadline is hardcoded before jurisdiction approval.

## Release gates

- Processing responsibility and basis register is approved for every first-release activity.
- Every collected field/event/object/SDK has purpose, classification, notice, owner, retention, rights and recipients.
- Guest/customer/account/external-channel identities cannot auto-merge.
- Notices/choices are usability-tested in launch languages with below-average digital-literacy users.
- Consent evidence/withdrawal and transactional-versus-marketing suppression pass end-to-end tests.
- Rights requests work for Account and non-Account people across all stores/providers/backups.
- Minor/child and sensitive-data boundary is approved and enforced.
- No call recording, advertising tracking or general model training ships through an unrelated toggle.
- Local counsel approves notices, roles, deadlines, transfers and exceptions.

## Open decisions

- Launch country, Ubi legal entity and Ubi/merchant/provider responsibility per activity.
- Exact lawful bases, rights/deadlines and regulator registration/DPO requirements.
- Guest booking versus required Account boundary.
- Age/minor/guardian rules and prohibited first-release verticals/fields.
- Merchant versus Ubi marketing architecture.
- Analytics/crash/AI vendors and user choices.
- Call recording remains deferred and off.
- Privacy contact/DPO/representative organization and complaint escalation.

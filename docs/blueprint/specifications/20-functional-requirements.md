# Functional requirements specification

Status: **Draft**  
Depends on: actors, identity, permissions, business taxonomy and booking workflows  
Blocks: domain architecture, APIs/events, database design, client information architecture and acceptance testing

## Purpose and use

This specification states what the complete Ubi platform must be capable of doing. It does not assign every capability to Version 1.

Each requirement must later trace to:

```text
actor and scenario
→ permission and policy
→ domain command/event
→ API/data contract
→ client behaviour
→ audit/observability
→ acceptance test
→ release milestone
```

## Requirement conventions

- `MUST` — mandatory when the capability is included in an approved release.
- `SHOULD` — expected direction; deviation needs a recorded reason.
- `MAY` — supported option or future extension.
- **Foundation** — architecture/domain capability required before dependent features.
- **Pilot candidate** — likely first-vertical capability, subject to final pilot decision.
- **Expansion** — later product capability.
- **Future gated** — requires a separate approval package.

Release labels do not replace the formal release contract.

## Identity, authentication and account lifecycle

| ID | Requirement | Stage |
|---|---|---|
| FR-IAM-001 | Ubi MUST let a person authenticate with an approved phone-first method without requiring email. | Foundation |
| FR-IAM-002 | Ubi MUST represent verified contact methods separately from the account and retain verification/revocation history. | Foundation |
| FR-IAM-003 | Ubi MUST assign immutable opaque account identifiers that are not phone numbers, emails or usernames. | Foundation |
| FR-IAM-004 | Ubi MUST support multiple active/revoked sessions and allow the person to review and terminate sessions/devices. | Pilot candidate |
| FR-IAM-005 | Ubi MUST require recent/stronger authentication for risk-classified ownership, settlement, export, integration and security actions. | Foundation |
| FR-IAM-006 | Ubi MUST provide risk-based account/contact recovery that does not silently transfer organization authority. | Foundation |
| FR-IAM-007 | Ubi MUST notify affected contacts/sessions of material security changes when safe and possible. | Pilot candidate |
| FR-IAM-008 | Ubi MUST support suspension, restriction, closure and redress states without deleting required records immediately. | Foundation |
| FR-IAM-009 | Ubi SHOULD support stronger device-bound or passkey authentication for privileged business users when target devices permit. | Expansion |
| FR-IAM-010 | Ubi MUST rate-limit, monitor and safely explain verification failures without disclosing whether unrelated accounts exist. | Foundation |

## Organization, brand, branch and operation

| ID | Requirement | Stage |
|---|---|---|
| FR-ORG-001 | An authenticated account MUST be able to create or join multiple organizations through separate controlled workflows. | Foundation |
| FR-ORG-002 | One organization MUST support multiple brands, branches and booking operations. | Foundation |
| FR-ORG-003 | Each operation MUST independently select a booking archetype/version, catalog, resources, availability and policies. | Foundation |
| FR-ORG-004 | Ubi MUST support one-person businesses without forcing unnecessary branch/staff complexity. | Pilot candidate |
| FR-ORG-005 | Ubi MUST maintain lifecycle and verification state separately for organization, public profile and individual operations. | Foundation |
| FR-ORG-006 | Authorized users MUST be able to activate, pause and archive operations without rewriting historical bookings. | Pilot candidate |
| FR-ORG-007 | Ubi MUST provide a public business slug/subdomain policy with uniqueness, reservation, rename and abuse handling. | Expansion |
| FR-ORG-008 | Ubi MUST provide a controlled organization closure workflow covering future bookings, balances, integrations, export and retention. | Foundation |
| FR-ORG-009 | Ubi SHOULD support test/demo operations that cannot be confused with production bookings or money. | Pilot candidate |
| FR-ORG-010 | Ubi MUST display organization/branch/operation context clearly whenever an account can act in more than one scope. | Foundation |

## Membership, roles, permissions and approvals

| ID | Requirement | Stage |
|---|---|---|
| FR-AUTHZ-001 | Ubi MUST authorize every protected server operation using active subject, tenant, capability, scope, resource and context. | Foundation |
| FR-AUTHZ-002 | Authorization MUST deny by default when any required rule is absent, stale or ambiguous. | Foundation |
| FR-AUTHZ-003 | Authorized users MUST be able to invite staff with an understandable role, explicit scope and optional expiry. | Pilot candidate |
| FR-AUTHZ-004 | An invite MUST be expiring, single-purpose, revocable and inactive until authenticated acceptance. | Foundation |
| FR-AUTHZ-005 | A delegating user MUST NOT grant capabilities or scopes beyond their delegated authority. | Foundation |
| FR-AUTHZ-006 | Authorized users MUST be able to suspend/revoke membership and promptly invalidate affected authority. | Foundation |
| FR-AUTHZ-007 | Offboarding MUST preserve authored operational history and support reassignment of open work. | Pilot candidate |
| FR-AUTHZ-008 | Sensitive actions MUST support step-up, amount/condition limits, delay or approval as configured by risk policy. | Foundation |
| FR-AUTHZ-009 | Approval requests MUST bind an exact action, target and before/after values and MUST expire or invalidate when context changes. | Expansion |
| FR-AUTHZ-010 | Ubi MUST audit membership, permission, ownership, privileged support and approval changes. | Foundation |

## Catalog, templates and configuration

| ID | Requirement | Stage |
|---|---|---|
| FR-CAT-001 | Ubi MUST offer familiar business-type templates that create appropriate operation/service defaults. | Pilot candidate |
| FR-CAT-002 | Authorized users MUST be able to create/edit/archive services, packages and add-ons within the operation’s archetype. | Pilot candidate |
| FR-CAT-003 | Services MUST support localized names/descriptions, duration, price method, fulfilment location and confirmation policy. | Foundation |
| FR-CAT-004 | Ubi MUST validate required archetype configuration before a service/operation can be publicly booked. | Foundation |
| FR-CAT-005 | Configuration and templates MUST be versioned so edits do not alter submitted/confirmed commitments. | Foundation |
| FR-CAT-006 | Ubi MUST support fixed, starting, per-unit/duration/person/night and quotation price methods only where the archetype permits. | Expansion |
| FR-CAT-007 | Ubi MUST use exact currency-aware amounts and reject unsupported currency/provider combinations. | Foundation |
| FR-CAT-008 | Ubi MUST support typed custom intake fields and bounded conditional questions without merchant executable code. | Expansion |
| FR-CAT-009 | Ubi MUST prevent forms/templates from collecting prohibited or unjustified sensitive data through policy controls. | Foundation |
| FR-CAT-010 | Authorized users SHOULD be able to test configuration with a clearly labeled test booking before publication. | Pilot candidate |

## Resources and availability

| ID | Requirement | Stage |
|---|---|---|
| FR-AVL-001 | Ubi MUST represent eligible staff, spaces, assets, pooled quantity and capacity resources required by approved archetypes. | Foundation |
| FR-AVL-002 | Authorized users MUST be able to define recurring hours, exceptions, breaks, blocks and setup/travel/cleanup buffers. | Pilot candidate |
| FR-AVL-003 | Availability computation MUST include operation/service rules, required resources, holds, confirmed commitments, capacity and timezone. | Foundation |
| FR-AVL-004 | Public availability MUST reveal only safe candidate choices, not private resource calendars or customer information. | Foundation |
| FR-AVL-005 | Final confirmation MUST revalidate authoritative availability transactionally. | Foundation |
| FR-AVL-006 | Ubi MUST support exact-resource choice and “any eligible resource” assignment where configured. | Pilot candidate |
| FR-AVL-007 | Ubi MUST support full-day/date, interval, session, date-range and time-window calculations as archetypes are released. | Expansion |
| FR-AVL-008 | Holds MUST consume configured capacity and expire/release/consume through authoritative idempotent transitions. | Foundation |
| FR-AVL-009 | Availability/configuration changes affecting future commitments MUST produce a conflict/impact review rather than silently invalidating them. | Foundation |
| FR-AVL-010 | Authorized overrides MUST require capability, reason, conflict disclosure and audit. | Pilot candidate |

## Standard booking and fulfilment

| ID | Requirement | Stage |
|---|---|---|
| FR-BKG-001 | Customers MUST be able to review and submit a standard booking/request without mandatory email/password creation. | Pilot candidate |
| FR-BKG-002 | Ubi MUST support instant, approval-required, payment-conditioned and approval-then-payment confirmation modes. | Foundation |
| FR-BKG-003 | Customer/staff MUST see truthful outcomes: draft, submitted/requested, held, payment processing, confirmed, cancelled or completed. | Foundation |
| FR-BKG-004 | Booking submission/confirmation MUST be idempotent and return the existing result for a duplicate logical attempt. | Foundation |
| FR-BKG-005 | A confirmed booking MUST snapshot schedule, service, price, policy, parties, source and confirmation evidence. | Foundation |
| FR-BKG-006 | Staff MUST be able to create call, walk-in, message and other authorized manual bookings quickly. | Pilot candidate |
| FR-BKG-007 | Ubi MUST record booking source/provenance separately from lifecycle state. | Foundation |
| FR-BKG-008 | Authorized actors MUST be able to accept, decline, reschedule, cancel, check in, start, complete and record no-show according to policy/state. | Pilot candidate |
| FR-BKG-009 | Rescheduling MUST preserve prior schedule history and atomically activate the new valid schedule. | Foundation |
| FR-BKG-010 | Cancellation MUST preview capacity, customer, payment/refund and policy consequences before commitment. | Pilot candidate |
| FR-BKG-011 | Ubi MUST distinguish customer, business/provider and system cancellation/no-show responsibility. | Foundation |
| FR-BKG-012 | Customers MUST receive a safe return path and understandable booking code/status after submission. | Pilot candidate |
| FR-BKG-013 | Customers SHOULD be able to reschedule/cancel through a scoped link or account when merchant policy permits. | Pilot candidate |
| FR-BKG-014 | Ubi MUST provide conflict resolution when capacity, permission, price or policy changed between selection and commit. | Foundation |
| FR-BKG-015 | Ubi MUST retain history when a service/resource is archived after bookings were created. | Foundation |

## Custom request, quotation and amendment

| ID | Requirement | Stage |
|---|---|---|
| FR-REQ-001 | Customer/staff MUST be able to create an editable structured custom-request draft from direct entry or conversation. | Expansion |
| FR-REQ-002 | Ubi MUST record answer provenance and permit correction of extracted/suggested values before commitment. | Expansion |
| FR-REQ-003 | Staff MUST be able to assign, triage, request clarification, schedule an assessment, prepare a proposal, decline or close a request. | Expansion |
| FR-REQ-004 | Every request MUST identify the responsible next party and action/deadline when applicable. | Expansion |
| FR-REQ-005 | Ubi MUST issue immutable, versioned quotations/options with scope, exclusions, schedule, price, deposit, policy and validity. | Expansion |
| FR-REQ-006 | Quotation revision MUST create a superseding version and MUST NOT rewrite an issued/accepted version. | Foundation |
| FR-REQ-007 | Customer acceptance MUST identify the exact quotation/option version and authentication/token evidence. | Foundation |
| FR-REQ-008 | Ubi MUST support an optional expiring capacity hold linked to a request/quotation. | Expansion |
| FR-REQ-009 | Conversion MUST atomically validate accepted version, conditions, capacity and idempotency before creating a confirmed booking. | Foundation |
| FR-REQ-010 | Scope/price changes after confirmation MUST use versioned amendments/change orders. | Expansion |

## Customer relationship records

| ID | Requirement | Stage |
|---|---|---|
| FR-CST-001 | Each organization MUST maintain tenant-private customer relationship records. | Foundation |
| FR-CST-002 | Staff MUST be able to find/reuse a customer using privacy-safe identifiers within authorized scope. | Pilot candidate |
| FR-CST-003 | Ubi SHOULD suggest possible duplicates but MUST require authorized review before merge. | Pilot candidate |
| FR-CST-004 | Merge/unmerge MUST preserve source, bookings, history and audit. | Expansion |
| FR-CST-005 | A customer MAY link merchant bookings/records to their Ubi account only after appropriate verification. | Pilot candidate |
| FR-CST-006 | Linking MUST expose only authorized customer-facing data and MUST NOT expose merchant-internal notes. | Foundation |
| FR-CST-007 | Ubi MUST distinguish booker, payer, attendee/beneficiary and organizer where relevant. | Expansion |
| FR-CST-008 | Staff notes MUST have purpose and visibility controls; regulated/health notes require separately approved products. | Foundation |

## Conversations and connected channels

| ID | Requirement | Stage |
|---|---|---|
| FR-CON-001 | Ubi MUST support booking/request-linked native conversations between authorized participants. | Expansion |
| FR-CON-002 | Ubi MUST normalize supported provider messages while retaining source provider IDs, connection and capabilities. | Expansion |
| FR-CON-003 | A merchant MUST connect external business channels only through official authorization with disclosed scopes. | Foundation |
| FR-CON-004 | Ubi MUST track channel connection lifecycle, token expiry/revocation, capability changes and disconnect effects. | Expansion |
| FR-CON-005 | Inbox users MUST be able to view assigned queues, claim/assign conversations, reply, close and reopen according to policy. | Expansion |
| FR-CON-006 | Staff MUST be able to pause automation/take over; substantive automated replies MUST stop until explicitly resumed. | Expansion |
| FR-CON-007 | Ubi MUST show source channel, delivery status and limitations without implying unsupported inbox history/control. | Expansion |
| FR-CON-008 | A conversation MUST be linkable to enquiries, requests, quotations, bookings and customers without merging their lifecycles. | Foundation |
| FR-CON-009 | Ubi MUST support block/report/rate-limit/retention controls before broad native messaging release. | Future gated |
| FR-CON-010 | Ubi MUST support fallback/action-needed handling when a provider/channel is unavailable or disconnected. | Expansion |

## Contextual handoff and deep links

| ID | Requirement | Stage |
|---|---|---|
| FR-LNK-001 | Ubi MUST create opaque, purpose-scoped, expiring and revocable handoff/manage-booking tokens. | Foundation |
| FR-LNK-002 | Token URLs MUST NOT contain raw personal, payment or conversation content. | Foundation |
| FR-LNK-003 | The same verified HTTPS link SHOULD open the authorized native destination when installed and a safe web fallback otherwise. | Expansion |
| FR-LNK-004 | Handoff-prefilled fields MUST show source and remain editable unless a verified immutable value is required. | Expansion |
| FR-LNK-005 | Token redemption MUST enforce tenant/resource/action scope and risk-based re-verification. | Foundation |
| FR-LNK-006 | Token forwarding, replay, expiry and revocation MUST fail safely without revealing unrelated data. | Foundation |

## Payments, deposits and reconciliation

| ID | Requirement | Stage |
|---|---|---|
| FR-PAY-001 | An operation/service MUST be able to configure no deposit, fixed deposit, percentage deposit or full prepayment where legally/provider supported. | Expansion |
| FR-PAY-002 | Cash/pay-at-location and manual payment evidence MUST remain valid methods with distinct states. | Pilot candidate |
| FR-PAY-003 | Ubi MUST separate payment intent, attempt/transaction, customer-reported evidence, provider verification, settlement, refund and dispute. | Foundation |
| FR-PAY-004 | Provider requests/webhooks MUST be authenticated, idempotent and linked through provider/Ubi references. | Foundation |
| FR-PAY-005 | Client/provider return pages MUST NOT alone mark a payment verified. | Foundation |
| FR-PAY-006 | Ubi MUST reconcile expected amounts/currency against provider events and settlement records. | Expansion |
| FR-PAY-007 | Amount/currency/late-payment/duplicate/reversal discrepancies MUST enter an exception queue without creating false booking state. | Expansion |
| FR-PAY-008 | Authorized users MUST be able to request/approve/execute/refuse refunds within policy/provider limits and audit. | Expansion |
| FR-PAY-009 | Booking deposits and refundable security/damage deposits MUST be modeled separately. | Future gated |
| FR-PAY-010 | Real-money release MUST include dispute, reversal, refund, settlement and provider-outage operations. | Foundation |

## Notifications and reminders

| ID | Requirement | Stage |
|---|---|---|
| FR-NTF-001 | Ubi MUST generate notification intent from authoritative domain events, not client assumptions. | Foundation |
| FR-NTF-002 | Ubi MUST distinguish transactional booking/request/payment messages from optional promotional communication. | Foundation |
| FR-NTF-003 | Users/businesses MUST be able to configure permitted channels/preferences within safety/legal/provider constraints. | Expansion |
| FR-NTF-004 | Notification delivery failure MUST NOT roll back a successful booking; it MUST retry/escalate according to importance. | Foundation |
| FR-NTF-005 | Ubi MUST record provider message ID/status, attempts and final outcome without treating delivery as customer comprehension. | Expansion |
| FR-NTF-006 | Reminder scheduling MUST use booking timezone/policy and remain idempotent after reschedule/cancellation. | Pilot candidate |
| FR-NTF-007 | Messages MUST use truthful state-specific wording and approved localized templates. | Foundation |
| FR-NTF-008 | Businesses MUST be able to see operationally important failed notifications requiring human follow-up. | Pilot candidate |

## Android local-first and synchronization

| ID | Requirement | Stage |
|---|---|---|
| FR-SYNC-001 | Android MUST read critical schedule, catalog, customer and pending-action views from durable local storage. | Foundation |
| FR-SYNC-002 | Android MUST display cached data without waiting for a network attempt to complete. | Foundation |
| FR-SYNC-003 | Eligible offline edits/bookings MUST commit locally with a durable command/outbox record before success is shown. | Pilot candidate |
| FR-SYNC-004 | Each queued command MUST carry stable identity, tenant, actor/membership version, device/time and idempotency information. | Foundation |
| FR-SYNC-005 | Ubi MUST visibly distinguish local/pending, synchronizing, synchronized, rejected and conflict states. | Foundation |
| FR-SYNC-006 | The server MUST reauthenticate/reauthorize and validate every synchronized command. | Foundation |
| FR-SYNC-007 | Conflicting scheduling, money, permission or customer-visible terms MUST NOT use silent last-write-wins. | Foundation |
| FR-SYNC-008 | Rejected/conflicting commands MUST preserve user input and provide resolution/cancel/escalation choices. | Pilot candidate |
| FR-SYNC-009 | Synchronization MUST support incremental pull, retry/backoff, deduplication, deletion tombstones and resumability. | Foundation |
| FR-SYNC-010 | High-risk online-only actions MUST explain connectivity need before the user enters unsupported work. | Pilot candidate |

## Reporting, analytics and exports

| ID | Requirement | Stage |
|---|---|---|
| FR-RPT-001 | Ubi MUST define metrics from governed operational states and a versioned metric dictionary. | Foundation |
| FR-RPT-002 | Business reports MUST enforce organization/branch/operation/permission scope. | Foundation |
| FR-RPT-003 | Ubi MUST distinguish requests, holds, confirmed bookings, fulfilment and verified/settled money in reports. | Foundation |
| FR-RPT-004 | Pilot reporting SHOULD include volume/source, response time, completion, cancellation/no-show, sync conflicts and payment exceptions. | Pilot candidate |
| FR-RPT-005 | Personal/financial row-level export MUST require separate capability, audit and secure expiring delivery. | Expansion |
| FR-RPT-006 | Analytics collection MUST respect minimization, consent/purpose and deletion/retention propagation. | Foundation |
| FR-RPT-007 | Ubi operations MUST be able to monitor aggregate reliability/adoption without routine access to message/customer content. | Foundation |

## Ubi operations, support, trust and administration

| ID | Requirement | Stage |
|---|---|---|
| FR-OPS-001 | Ubi MUST provide separate workforce identity and permissions from merchant membership. | Foundation |
| FR-OPS-002 | Support MUST default to redacted metadata and require purpose-scoped audited elevation for content/action access. | Foundation |
| FR-OPS-003 | Assisted onboarding MUST record merchant authorization and changes made on the merchant’s behalf. | Pilot candidate |
| FR-OPS-004 | Ubi MUST provide controlled queues for account recovery, integration failure, payment reconciliation, abuse and privacy requests. | Expansion |
| FR-OPS-005 | Break-glass access MUST require strong authentication, reason, short expiry, alert and mandatory review. | Foundation |
| FR-OPS-006 | Authorized trust/safety staff MUST be able to restrict/suspend actors/resources proportionately with reason, notice and redress controls. | Expansion |
| FR-OPS-007 | Operators MUST NOT directly rewrite domain/financial history outside controlled commands and reconciliation procedures. | Foundation |
| FR-OPS-008 | Production administration MUST be auditable and separated from development/test environments. | Foundation |

## Privacy and data-subject operations

| ID | Requirement | Stage |
|---|---|---|
| FR-PRV-001 | Ubi MUST maintain a data inventory mapping field/event/object to purpose, owner, classification and retention. | Foundation |
| FR-PRV-002 | Ubi MUST capture and version required notices, agreements and consent/communication preferences. | Foundation |
| FR-PRV-003 | Authorized users MUST be able to request/access/correct/export/delete/restrict data according to verified identity, responsibility and law. | Expansion |
| FR-PRV-004 | Deletion/restriction MUST propagate to operational, search, object, cache and analytical systems subject to documented exceptions. | Foundation |
| FR-PRV-005 | Legal/contractual retention and legal hold MUST be explicit and MUST NOT be presented as completed deletion. | Foundation |
| FR-PRV-006 | Merchant closure/disconnect MUST revoke provider access and follow export/retention/deletion procedures. | Foundation |
| FR-PRV-007 | Ubi MUST maintain a current subprocessor/vendor/data-location register for production processing. | Foundation |

## Future native calling

| ID | Requirement | Stage |
|---|---|---|
| FR-CALL-001 | Calls MUST be limited to verified/scoped Ubi identities or booking relationships before general discovery is considered. | Future gated |
| FR-CALL-002 | Ubi MUST support ringing, answer/decline, missed/ended/failed states and poor-network adaptation. | Future gated |
| FR-CALL-003 | Calls MUST include block/report/rate-limit/consent/privacy controls before release. | Future gated |
| FR-CALL-004 | Call signaling, metadata, media relay and optional recording MUST have separate retention/access purposes. | Future gated |
| FR-CALL-005 | Recording MUST be disabled unless a separately approved lawful, consented and operational use is established. | Future gated |
| FR-CALL-006 | Ubi MUST disclose that it is not an emergency service and cannot guarantee carrier-grade reachability. | Future gated |

## Functional acceptance rule

A requirement is implementation-ready only when it has:

- Approved release and actor.
- Happy path and exception cases.
- Permission and audit decision.
- Data ownership, classification and retention.
- Domain/API/event contract.
- Offline/provider failure behaviour where applicable.
- Localized customer/staff wording.
- Automated and human acceptance tests.
- Observability and operational owner.
- Rollback/migration effect.

## Open decisions

- Which subset constitutes the first pilot release.
- Launch market, first vertical and its required archetype modules.
- Authentication/provider/payment/channel capabilities available in that market.
- Exact notification, refund, cancellation and support policies.
- Which customer-web features precede desktop versus follow it.
- Whether any custom-request capability is needed in the first pilot.
- Which future requirements require separate repositories/services versus modules in one platform.

# Conversation, contextual handoff and human-takeover workflows

Status: **Draft**  
Depends on: actors, identity/tenancy, permissions, booking archetypes and standard/custom workflows  
Blocks: integration architecture, native messaging, automation, contextual links, client inboxes, data model, privacy and operations

## Purpose

Ubi must help a business turn enquiries received by call, walk-in, Ubi messaging or an officially connected channel into accurate structured work. Conversation assists the booking process; it is not the authoritative booking, availability or payment record.

This document defines a provider-neutral model. A named platform is supported only after its official business API, account eligibility, permissions, review process, regional availability, pricing, message windows and operational limits are verified at implementation time.

## Product boundary

Ubi MAY support three independent entry modes:

1. **Direct Ubi entry** — the business shares its public Ubi page/link and the customer books or requests directly.
2. **Contextual handoff** — staff or approved automation collects limited information, creates an editable structured draft and sends a secure Ubi link to continue.
3. **Connected conversation** — a business explicitly connects an eligible business channel through its official authorization process so authorized staff can handle supported messages in Ubi.

A merchant may use any supported mode without enabling the others. Social-channel connection is never required to use Ubi booking.

Ubi MUST NOT:

- ask for a merchant’s social-media password;
- scrape a consumer/user inbox or automate a personal account;
- imply access to history, delivery/read status, templates, calling or message types a provider does not expose;
- make a booking look confirmed because a message was received or replied to;
- send an irreversible automated answer based only on uncertain extracted text;
- combine unrelated customers merely because they share a display name or phone.

## Core concepts

| Concept | Meaning |
|---|---|
| Channel provider | Ubi, an approved external messaging platform, phone/walk-in/manual source or another approved origin. |
| Channel connection | One organization-authorized link between Ubi and an eligible external business account/inbox. |
| Provider identity | Provider-specific external account/contact ID. It is not automatically a Ubi account. |
| Conversation | A provider-scoped thread/container with participants, source and operational links. |
| Message | An immutable inbound/outbound/system communication item with provider/source identity and delivery evidence. |
| Inbox work item | Ubi’s operational assignment/status wrapper for handling a conversation. |
| Automation mode | Whether automation is disabled, suggesting, responding within policy, paused or requiring takeover. |
| Structured draft | Editable booking/request/customer information proposed from a conversation or entered by staff. |
| Contextual handoff | A secure link that lets a customer continue a narrowly scoped draft/action in Ubi web or app. |
| Human takeover | Explicit transfer from automated response to an authorized staff member/queue. |

Conversation, work assignment, automation, booking/request, notification and delivery each have separate state.

## Channel connection lifecycle

```text
not_connected
→ authorization_started
→ active
↔ degraded
→ reauthorization_required
→ active
→ disconnected
→ deletion_pending
→ deleted

authorization_started → failed/expired
active/degraded/reauthorization_required → revoked_by_provider
```

### Rules

1. The initiating account must hold `integration.connect` for the exact organization and intended scope.
2. Ubi redirects through or embeds only an official provider authorization flow.
3. The merchant must see the business account, requested permissions, data categories, intended Ubi functions and disconnect effect before completion.
4. Ubi validates the returned account and capabilities; successful OAuth alone does not prove that messaging/webhooks are available.
5. Provider tokens are environment- and connection-scoped, encrypted, inaccessible to ordinary staff and rotatable/revocable.
6. The connection stores granted scopes and discovered capabilities separately from the capabilities Ubi has enabled.
7. Webhook/event authenticity and account/connection matching must be verified before accepting content.
8. Scope loss, token expiry, provider policy change and webhook failure move the connection to a truthful degraded/reauthorization state.
9. Disconnect immediately stops new provider access and outgoing work where possible; retained Ubi booking/audit records follow the approved retention schedule.
10. Reconnection must not silently join a different external business account to old conversations.

## Capability discovery and truthful degradation

Every connection has a versioned capability profile, for example:

- receive/send text;
- receive/send supported media;
- reply window and template restrictions;
- delivery/read status availability;
- participant/contact fields;
- conversation history range;
- interactive actions/buttons;
- webhook/event availability;
- initiation rules;
- calling availability;
- rate/throughput limits.

The client and automation policy consume this profile. A missing/unknown capability is unavailable by default. Provider capability changes generate operational alerts, impacted-work queries and adapter contract tests.

## Conversation operational model

### Work status

```text
new
→ open
↔ waiting_for_customer
↔ waiting_for_business
→ resolved
→ closed

new/open/waiting_* → spam_or_blocked
resolved → open                (new relevant activity or authorized reopen)
closed → open                  (authorized/manual/provider event where policy permits)
```

### Assignment

Assignment is independent of work status:

- unassigned queue;
- assigned to a staff member;
- assigned to a team/branch/operation queue;
- claimed temporarily;
- escalated to manager/specialist;
- released/reassigned.

An assignment does not grant access beyond the staff member’s existing membership/capability/scope.

### Automation mode

```text
disabled
suggest_only
active_within_policy
paused_by_staff
paused_by_risk
handoff_required
```

Human takeover changes automation mode; it does not change booking or conversation status automatically.

### Message delivery state

Inbound messages record `received` after Ubi accepts verified provider evidence. Outbound messages may progress through:

```text
draft → queued → submitted_to_provider → sent/delivered/read
                            ↘ failed
                            ↘ unknown
```

States are provider-capability dependent. Ubi must not infer `delivered` or `read` from `submitted_to_provider`. Retrying an outbound logical message must be idempotent and must avoid duplicate customer messages where the provider permits.

## End-to-end inbound enquiry workflow

### Preconditions

- The channel is Ubi-native/manual or has an active official connection.
- The organization/branch/operation routing policy is known or a safe triage queue exists.
- Retention, notice, permission and automation policy for the source are active.

### Workflow

1. Ubi receives an authenticated event or authorized staff creates a call/walk-in enquiry.
2. Ubi deduplicates the provider event/logical message using stable provider identifiers and payload evidence.
3. The adapter stores the original provider reference and normalizes only supported fields.
4. Ubi resolves the connection/tenant before any customer matching or routing.
5. Ubi creates/updates the provider-scoped participant and conversation. It may suggest a tenant-private customer match but cannot auto-merge uncertain identities.
6. Routing selects the appropriate branch/operation/queue from explicit provider mapping, customer choice or staff triage.
7. Policy classifies the message as ordinary enquiry, actionable booking/request input, sensitive/unsupported content, abuse/spam or human-review required.
8. Authorized automation may answer low-risk factual questions from approved business configuration, request missing non-sensitive information or prepare a suggestion/draft.
9. Availability, price, policy and payment answers come from authoritative versioned Ubi data and include truthful limitations.
10. If information is sufficient, Ubi creates an editable structured draft linked to message evidence; it does not confirm the booking.
11. Ubi either lets authorized staff review/create the record or sends a purpose-scoped handoff link for customer review.
12. The booking/request follows its own authoritative workflow. Conversation displays its result only after observing the domain event.
13. Work status and next responsible actor are updated explicitly.

## Information extraction and provenance

Each proposed field from a conversation MUST record:

- destination field and normalized candidate value;
- source message/provider/time and relevant excerpt reference;
- entry method: customer selection, staff entry, deterministic parse or model suggestion;
- confidence/validation result where applicable;
- who reviewed/edited/accepted it;
- whether the source may be retained or must be deleted/redacted earlier;
- the final value committed to a request/booking snapshot.

Low confidence, contradictory dates, ambiguous calendars, uncertain customer identity, price changes, payment claims and exclusive-resource selections require review. Ubi may suggest; it cannot silently make a commitment.

## Contextual handoff workflow

The target experience is: a customer starts in conversation, receives a link, sees a familiar Ubi page already prepared for the correct business/service, corrects it if needed and safely completes the next step.

### Token requirements

A contextual handoff token must be:

- cryptographically random/opaque;
- purpose-, tenant-, resource- and action-scoped;
- short-lived according to risk;
- single-use when the action is irreversible, or bounded-use for a manage flow;
- revocable and auditable;
- stored hashed or otherwise protected so database/log disclosure does not reveal a usable secret;
- excluded from analytics, referrer leakage, screenshots/logs and third-party page resources where feasible;
- verified again for sensitive data/action according to risk.

The URL MUST NOT contain raw name, phone, message, service answer, price, payment detail or conversation transcript.

### Handoff states

```text
draft_prepared
→ link_issued
→ opened
→ reviewed
→ submitted/accepted

link_issued/opened/reviewed → expired
link_issued/opened/reviewed → revoked
link_issued/opened/reviewed → superseded
```

Opening a token is not acceptance. Link delivery is not opening. Submission still invokes the request/booking command with server-side permission, configuration, availability, price and idempotency validation.

### Web/app continuation

1. Ubi issues one verified HTTPS URL under an approved Ubi domain.
2. If a compatible customer app is installed, verified app/universal-link association may open the native destination.
3. Otherwise the same URL opens a lightweight responsive web flow.
4. Authentication state never travels as raw conversation data in the URL.
5. The destination displays the exact business/branch/operation, purpose and fields proposed from conversation.
6. Suggested fields are visually identified and editable; immutable verified fields explain why they cannot be edited.
7. The customer reviews the authoritative price/availability/policy and submits.
8. Ubi returns a booking/request code and truthful status, then posts an allowed status message back to the source conversation if configured.

Deep linking is an enhancement. The web fallback must complete the supported action without installing an app.

## Human takeover workflow

### Takeover triggers

Takeover is required for:

- explicit customer request for a person;
- confidence or policy threshold failure;
- complaint, dispute, cancellation exception or refund;
- ambiguous calendar/date, price, identity, capacity or payment;
- sensitive/prohibited data;
- abusive or threatening communication;
- provider limitation or automation failure;
- defined high-value/VIP/regulated scenarios;
- automation loop or repeated non-comprehension.

### Takeover process

1. Automation stops substantive outbound replies immediately and records the trigger.
2. The conversation enters `handoff_required` and an authorized queue with urgency/age.
3. The assigned staff member sees a concise provenance-linked summary plus the original authorized content, not invented conclusions.
4. Staff claim/accept assignment and the customer receives a truthful acknowledgement if permitted.
5. Staff communicate and create/change Ubi records only through their own capabilities.
6. Automation may still offer private suggestions, clearly labelled, if policy allows.
7. Resume requires an explicit authorized action and visible automation scope; new sensitive context may keep it paused.
8. Takeover/resume, assignment and material actions are audited.

No user is told a person has taken over until a person/queue has actually accepted responsibility or the wording clearly says the message is waiting.

## Staff inbox experience rules

The eventual business inbox MUST prioritize:

- new/unassigned and overdue customer work;
- source channel and connection health;
- branch/operation context;
- customer/request/booking links without lifecycle merging;
- last message, next responsible actor and time waiting;
- automation/human-control state;
- delivery failures and provider limitations;
- explicit `Create booking`, `Create request`, `Send secure link`, `Assign`, `Pause automation` and `Close` actions where authorized.

The inbox cannot be the only way to find a confirmed booking. Booking operations remain usable when every external channel is disconnected.

## Manual call and walk-in capture

Phone calls and walk-ins do not produce provider messages by default. Staff can:

1. select `Phone call` or `Walk-in` as source;
2. find/create the tenant-private customer with minimal justified details;
3. choose standard booking or custom request;
4. record structured answers and optional short operational note;
5. review/confirm under the normal booking rules;
6. send a customer receipt/manage link through an approved channel if consent and contact permit.

Call recording or transcript ingestion is excluded unless separately approved for law, consent, security, retention and vendor processing.

## Automation authority boundary

### Automation may, when approved

- answer a configured address, hours or service-description question;
- present safe candidate services/times/prices;
- ask bounded intake questions;
- prepare editable drafts and summaries;
- issue a scoped handoff link;
- send state-based acknowledgements/reminders;
- route/escalate to staff.

### Automation may not independently

- override unavailable capacity or permission;
- accept a custom scope/quotation on behalf of a party;
- mark unverified payment as verified;
- grant refund/discount/ownership/integration access;
- make a legal, medical or safety-critical determination;
- reveal tenant/customer information to an unverified participant;
- resume after human takeover without explicit policy/action;
- train a third-party model on customer/business content without separate authority and disclosure.

Every substantive automated outbound message records the policy/rule/model version, source facts, result and takeover path necessary for audit and investigation.

## Identity linking and customer privacy

- A provider participant remains provider-connection scoped until verified linking occurs.
- A phone number displayed by a provider may be a matching clue, not proof of ownership or a global identity.
- Customer/account linking requires purpose-appropriate verification and cannot expose merchant-internal notes or unrelated bookings.
- Cross-merchant conversation/customer matching is prohibited for ordinary business use.
- Staff search and content access are tenant/branch/operation/case scoped.
- Sensitive content accidentally received must support restricted access, redaction/quarantine and deletion according to policy/provider obligations.
- Participants need appropriate notice that the business uses Ubi and may use permitted assistance/automation.

## Retention, deletion and evidence

Separate retention schedules are required for:

- provider event envelope and deduplication key;
- message content and supported media;
- provider participant/contact fields;
- normalized conversation metadata;
- structured extraction/provenance;
- booking/request snapshots derived from messages;
- automation decision evidence;
- delivery status and operational logs;
- connection authorization/audit and revoked tokens;
- abuse/legal-hold evidence.

Disconnecting a provider connection does not automatically delete Ubi bookings or legally required audit/financial evidence. It does stop future access and starts the approved connection-data retention/deletion process. A merchant/customer deletion request must report provider-owned copies outside Ubi’s control truthfully.

## Failure and exception catalogue

| Failure | Required behaviour |
|---|---|
| Duplicate provider event | Deduplicate; no duplicate message, reply, draft or booking. |
| Event arrives out of order | Preserve provider time/reference and causal ordering where known; do not rewrite booking truth. |
| Provider omits delivery/read events | Show the highest evidenced state only. |
| Token/scope expires | Stop affected operations, mark degraded/reauthorization required and notify authorized business users. |
| Connection mapped to wrong operation | Quarantine/unroute until authorized correction; never leak to another tenant. |
| Provider outage/rate limit | Queue only policy-eligible outbound work with expiry; show delay/failure and offer direct Ubi/manual fallback. |
| Customer repeats a message | Do not generate repeated substantive automation blindly; detect loop and hand over. |
| Staff and automation reply together | Human takeover/lease/version check prevents automated send; if already submitted, show both evidence and pause. |
| Conversation says a slot is available but it is later taken | Final authoritative availability check rejects/conflicts and preserves the draft. |
| Link is forwarded | Scope/re-verification and expiry limit exposure; do not reveal unrelated conversation data. |
| Link expires during editing | Preserve non-sensitive entered work locally where safe, explain expiry and offer reissue/reverification. |
| Customer blocks/opts out | Respect provider/Ubi restrictions immediately except legally permitted essential notices through appropriate means. |
| Unsupported media/content | Preserve safe provider reference/metadata, explain limitation and provide a supported alternate path. |
| Integration disconnected | Ubi bookings/manual operations continue; queued sends stop/fail safely; reconciliation identifies uncertain submissions. |

## Invariants

1. A conversation or message never confirms a booking by itself.
2. Provider delivery/read status never proves booking acceptance, attendance or payment.
3. A channel connection belongs to exactly one organization connection boundary, even if routed to many branches/operations.
4. Every external message is attributable to provider, connection, provider ID and direction.
5. Every automated reply is bounded by current connection capability and Ubi policy.
6. Human takeover prevents further substantive automated replies until explicit resume.
7. Extracted data remains proposed until reviewed/committed by an authorized actor.
8. Handoff URLs contain opaque tokens, not personal or conversation data.
9. Link redemption cannot exceed its original tenant/resource/action scope.
10. Disconnect/revocation prevents new provider access without deleting authoritative Ubi history silently.
11. Cross-tenant customer/conversation discovery is denied by default.
12. A provider outage cannot make a committed Ubi booking disappear or become unconfirmed.

## Permissions and audit

Separate capabilities are required for:

- view/reply to conversations;
- view media/sensitive content;
- assign/manage inbox work;
- create booking/request from conversation;
- issue/revoke contextual links;
- pause/resume automation;
- configure automation knowledge/policy;
- connect/reauthorize/disconnect a provider;
- export/delete conversation data;
- review quarantined/abuse content;
- privileged Ubi support access.

Audit must record connection changes, routing/configuration changes, assignment, takeover/resume, link issue/revocation/redemption, automated sends, staff sends, content access where risk requires and every resulting domain command.

## Release sequence

The realistic sequence is:

1. **Foundation:** direct Ubi business link, manual call/walk-in source, tenant-safe provenance and scoped manage/handoff token service.
2. **Pilot candidate:** staff sends a secure prefilled draft/review link through an already approved contact method; no universal inbox required.
3. **Expansion:** Ubi-native messaging and staff inbox with assignment/human takeover/safety controls.
4. **Expansion by provider:** one official connected channel at a time after technical, commercial, privacy and operational validation.
5. **Later:** bounded automated responses/suggestions after real enquiry data and human escalation performance are proven.
6. **Future gated:** open username messaging, general calling, recording/transcription and multi-provider AI automation.

This order provides booking value even if no external platform approves an integration.

## Acceptance scenarios

1. A receptionist creates a confirmed walk-in booking while offline; source remains `walk_in`, the record later synchronizes once and no conversation is invented.
2. A customer messages an officially connected account, provides service/date information and receives a link showing editable proposed answers; opening it does not confirm anything.
3. A forwarded handoff link cannot display another conversation or perform a broader action than intended.
4. A provider sends the same webhook five times; Ubi records one logical inbound message and one automation outcome.
5. Staff takes over while automation has a reply queued; version/lease policy prevents the automated substantive reply from being sent afterward.
6. Provider delivery status is unavailable; Ubi says submitted/sent only, never read.
7. The provider token is revoked; the business sees reauthorization required, Ubi bookings remain accessible and new outbound messages stop safely.
8. A suggested time is taken before final confirmation; the booking command returns a clear conflict and preserves all customer input.
9. Two provider contacts share a display name; Ubi does not merge them or reveal one person’s booking to the other.
10. A business chooses never to connect social channels and still completes every core booking workflow through its Ubi link and manual capture.

## Open decisions

- First launch channel and official business-account eligibility in the selected country.
- Whether any connected inbox capability belongs in the pilot after direct-link validation.
- Retention schedule and controller/processor responsibilities for connected content.
- Automation technology, approved knowledge source, confidence policy and vendor/data-processing terms.
- Handoff token duration and re-verification by action risk.
- Provider-specific customer-notice/template/message-window rules and cost.
- Ubi-native messaging identity discovery boundary and block/report operations.

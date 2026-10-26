# Notifications, reminders, native messaging and future calls

Status: **Draft**  
Depends on: actors, identity, permissions, booking/payment/conversation workflows, quality and localization rules  
Blocks: communications architecture, client experiences, provider selection, safety/privacy/legal design and operations

## Purpose

Ubi needs to tell the right person what actually happened, through an available and permitted channel, without confusing a sent message with a confirmed booking or understood instruction. Ubi-native messaging and calls are long-term requirements, but their identity, abuse, privacy, weak-network and operating systems must exist before broad release.

This document separates:

1. **Notifications** — system-generated, event-driven information or action prompts.
2. **Reminders** — time/policy-driven notifications about an existing obligation or booking.
3. **Native messages** — participant-authored or approved automated conversation content inside Ubi.
4. **Calls** — real-time Ubi audio communication, future-gated.

An external connected-channel reply is governed primarily by the conversation specification, even if the delivery infrastructure is shared.

## Communication principles

1. Domain truth produces communication intent; a client screen does not invent it.
2. Communication failure does not undo a committed booking/payment action.
3. `queued`, `submitted`, `sent`, `delivered`, `read` and `understood` are different.
4. A recipient receives only data necessary for that communication purpose and channel risk.
5. Transactional/service communication and promotional marketing are separately classified and governed.
6. User/business preferences apply within safety, contract, law and provider limitations; they do not authorize spam.
7. Critical status must be retrievable in Ubi/web even when push, SMS or external messaging fails.
8. Templates are versioned, localized and tied to exact domain states.
9. Retries are idempotent and bounded; stale messages expire rather than arrive misleadingly late.
10. Channel cost, latency, data usage, deliverability and accessibility are product constraints.

## Core communication objects

| Object | Meaning |
|---|---|
| Notification policy | Versioned rule mapping a domain event/recipient/purpose to eligible content, urgency and channels. |
| Notification intent | Durable first-party instruction to inform a recipient about an exact event/action. |
| Recipient endpoint | Verified or approved push token, phone/SMS endpoint, email, Ubi inbox or connected participant address. |
| Delivery attempt | One provider/channel submission for an intent with outcome evidence. |
| Template version | Approved localized subject/body/action structure with data classification. |
| Reminder schedule | Versioned plan derived from booking/obligation time, timezone and policy. |
| Preference/consent record | Purpose/channel choice, source, time, version and withdrawal/override conditions. |
| Native thread/message | Ubi participant conversation and immutable message item. |
| Call session | Authorized signaling/media session with participant, booking relationship and lifecycle metadata. |

## Notification intent lifecycle

```text
created
→ scheduled/ready
→ dispatching
→ completed

scheduled → cancelled/superseded/expired
ready/dispatching → partially_completed
ready/dispatching → failed_action_required
```

Each delivery attempt independently records:

```text
queued → submitted_to_provider → sent/delivered/read
                           ↘ retryable_failure
                           ↘ permanent_failure
                           ↘ unknown
                           ↘ expired
```

Only states supported by provider evidence are used. Intent `completed` means the configured dispatch policy completed; it does not prove reading or comprehension.

## Notification generation workflow

1. An authoritative domain transaction records a versioned event and durable outbox item.
2. Communication policy evaluates event type/version, tenant/operation, recipient role, language, purpose, sensitivity, urgency, preferences, legal/provider rules and available endpoints.
3. Ubi creates one idempotent notification intent per logical recipient/purpose/event.
4. The renderer resolves an approved localized template and minimum safe fields from the authoritative snapshot.
5. Scheduled items wait on an authoritative server time; immediate items enter a priority-aware queue.
6. A channel adapter checks current endpoint validity, provider capability, template/window/rate limits and cost policy.
7. The adapter submits with a stable Ubi/provider reference and bounded timeout.
8. Provider callbacks/status checks update delivery attempts idempotently.
9. Retry policy considers error class, urgency, age, cost and duplicate risk. Permanent/expired failures do not loop.
10. If a critical notice cannot be delivered, Ubi exposes it in the business action queue and, where approved, attempts a fallback channel.
11. Opening an action link revalidates its token and current domain state; an old notification cannot perform a stale action.

## Notification classes

| Class | Examples | Default treatment |
|---|---|---|
| Security | OTP, recovery, new privileged device, ownership/provider-account change | Highest protection; minimal content; no marketing consent dependency; strict rate/abuse controls |
| Transactional booking | submitted, accepted/declined, confirmed, rescheduled, cancelled, action/payment needed | Required service communication subject to lawful channel/purpose; truthful state-specific wording |
| Payment/reconciliation | payment reported, processing, verified, failed, refund, exception | Sensitive/minimized; never include full credentials/reference; booking state shown separately |
| Operational staff | new request, conflict, overdue response, sync rejection, failed critical customer notice | Role/scope restricted; app queue primary, push as alert |
| Reminder | appointment/event/hold/payment deadline | Scheduled from current domain state; cancel/reschedule safe |
| Product/service | onboarding help, feature notice | Preference/necessity assessed; not disguised as booking notice |
| Promotional | offers, campaigns, cross-sell | Explicit separate lawful/consent basis, opt-out and merchant/platform policy required |

Promotional content MUST NOT be inserted into OTP, security, receipt, cancellation or essential booking messages merely to bypass preference/provider rules.

## Channel strategy

### In-app notification center

- Primary durable place for authenticated business operational notices.
- Stores safe title, summary, created/read/archived status, action link and current-state resolution.
- Read state is per account/device synchronization, not evidence that an external staff process is complete.
- High-volume informational notices should be grouped without hiding critical items.

### Push notifications

- Used as a wake-up/summary channel, not a secret-bearing database.
- Lock-screen text defaults to privacy-safe wording for sensitive matters.
- Push tokens are device/app-instance endpoints, rotatable/revocable and removed after invalidation.
- Tap opens an authenticated, authorization-checked destination; payload data is not trusted as current state.
- Failure or disabled push is normal and does not block core operation.

### SMS

- Reserved for approved security, critical booking/customer communication and selected fallback where cost/consent/law permit.
- Messages identify Ubi/business/purpose clearly enough to resist phishing and avoid sensitive detail.
- Short links use an approved Ubi domain and scoped tokens; no third-party public shortener.
- Delivery receipts are provider evidence only and may be unavailable/ambiguous.
- Phone recycling, shared phones, sender-ID rules, country routing, cost and fraud are assessed per launch market.

### Email

- Optional where users/businesses have a usable verified address; it cannot be mandatory for phone-first core use.
- Appropriate for longer receipts/exports but still uses secure links for sensitive records.

### Connected external channel

- Only through an active official business connection and provider-allowed message/template/window.
- Conversation identity/capability and opt-out rules apply.
- Failure falls back to Ubi/manual handling according to policy; Ubi does not route around provider restrictions.

## Reminder workflow

### Schedule creation

1. A confirmed/requested booking, hold, quote or payment obligation emits a scheduling event.
2. Reminder policy calculates candidate times using the booking/recipient timezone, quiet hours, urgency and channel availability.
3. Ubi stores reminder schedule, domain version, template/policy version and idempotency identity.
4. A reschedule, cancellation, confirmation-state change, payment resolution or opt-out invalidates/recalculates affected reminders.
5. At dispatch, Ubi rechecks current state and relevance before rendering/sending.

### Rules

- “Tomorrow” must include the exact local date for consequential reminders.
- The business/branch/service and required next action must be unambiguous.
- Appointment reminders cannot expose a sensitive service category on a shared lock screen/SMS without approved purpose.
- A reminder arriving after cancellation/completion is a correctness incident, not a harmless cosmetic error.
- The same reminder/event cannot be re-sent by every retry; recurrence must be explicit.
- No-show reduction experiments need consent/purpose, metric definition and fairness review; do not penalize customers solely because a provider says delivery failed.

## Preferences, consent and endpoint ownership

Preferences are purpose- and channel-specific. The system records:

- subject/account/merchant relationship;
- communication purpose/category;
- endpoint/channel;
- allowed/denied/required status and reason;
- capture source, notice/version, time and actor;
- withdrawal/change time;
- lawful/contractual override and expiry where applicable.

Endpoint verification proves control at that time, not permanent identity. Changed/recycled numbers and devices require session/contact lifecycle handling.

The business cannot use Ubi transactional permissions as unlimited marketing permission. Ubi and merchant controller/processor responsibilities, sender identity, opt-out, suppression retention and marketing rules require jurisdiction/provider analysis.

## Template and content governance

Every template version defines:

- stable purpose/event/state;
- allowed recipient roles and channels;
- locale and fallback locale;
- approved variables with classification and maximum length;
- link/action type and token risk;
- provider template/reference/version where required;
- reviewer/approval and effective/expiry time;
- accessibility/plain-language checks;
- test snapshots for important variants;
- whether the text is transactional, security, operational or promotional.

Templates use immutable booking/price/policy evidence where a historical commitment is described and query current state when presenting a current action. Mixing these must be explicit.

## Failure, retry and fallback

| Failure | Required behaviour |
|---|---|
| Provider timeout/unknown | Record unknown, query/reconcile if supported and avoid blind duplicate send. |
| Temporary provider/rate failure | Bounded exponential backoff with jitter, priority and expiry. |
| Invalid/revoked endpoint | Mark unusable, stop retry, prompt safe re-verification through another authenticated path. |
| Template rejected/provider window closed | Record permanent/action-required outcome; choose only an approved alternate template/channel. |
| Language template missing | Use explicit reviewed fallback if policy permits; block legal/payment/safety send if unsafe. |
| Notification becomes stale | Recheck domain version/state and cancel/supersede before send. |
| All customer channels fail | Booking remains committed; business sees follow-up task and printable/call/manual option where permitted. |
| Push tap from logged-out/unauthorized device | Authenticate and authorize; disclose no sensitive destination data before success. |
| Duplicate domain event | Reuse notification intent identity; do not send duplicate logical notice. |
| Provider falsely/temporarily reports status | Preserve evidence history; no business/payment state is derived from read/delivery. |
| Notification queue delay | Measure age by priority, alert before expiry and suppress misleading late reminders. |

## Ubi-native messaging boundary

Native messaging is required long-term for business/customer communication, but the first safe release SHOULD limit who can contact whom.

### Recommended relationship model

A customer may message a business when at least one is true:

- the customer opened the business’s verified Ubi public profile/link and starts an enquiry permitted by that business;
- an existing request/booking creates a scoped conversation relationship;
- the business sent a customer-authorized invitation/handoff;
- a later discovery model explicitly permits contact with anti-abuse controls.

Knowing a username/unique ID must not automatically grant unlimited contact or reveal private activity. Business/public usernames and immutable internal IDs remain separate.

### Native conversation/message lifecycle

Thread work/assignment follows the conversation specification. A native message may be:

```text
local_draft → queued → accepted_by_ubi → delivered_to_recipient_device/account → read
                         ↘ rejected/failed
                         ↘ removed_or_restricted
```

`read` is a synchronized Ubi user action/receipt where enabled, not proof of comprehension. The sender must see when a message remains only on their device or is rejected.

### Minimum safety system before release

- verified/scoped identity and relationship rules;
- block, mute and report;
- per-account/device/network/business rate limits;
- spam/scam/impersonation and malicious-link/media controls;
- staff/queue access and offboarding;
- user-visible automation/bot identity;
- retention, deletion, export and legal-hold process;
- content/media access and moderation policy;
- notification privacy and quiet controls;
- abuse investigation tooling with audited restricted access;
- appeal/redress for restrictions where appropriate;
- emergency/law-enforcement process reviewed for the launch country.

### Message integrity and editing

- Each message has an immutable ID, author/source, server acceptance time and conversation scope.
- If editing is offered, revisions or edit history/evidence rules must prevent rewriting accepted quotations, policies or abuse evidence.
- Deletion must distinguish local hide, participant-visible removal and system retention obligations.
- Attachments are separate protected objects with type/size/scan/access/expiry controls.
- Quotation, booking, payment and identity actions use structured Ubi objects/commands; free text cannot impersonate them.
- End-to-end encryption must not be claimed unless the implemented key, multi-device, backup, reporting and metadata model actually provides it and has been independently reviewed.

## Future Ubi call boundary

In-app calling is an explicit long-term product requirement, but it is **future gated** because reliable real-time media, abuse response, network adaptation, privacy and operating cost are materially different from booking/message delivery.

### Initial allowed call relationships

- verified Ubi business staff to a customer linked to an active enquiry/request/booking, subject to contact preference and business policy;
- customer to a verified business through its profile or active booking relationship;
- no unrestricted person-to-person discovery or arbitrary username calling in the first call release.

### Call lifecycle

```text
created/authorized
→ signaling
→ ringing
→ connecting
→ connected
→ ended

created/signaling/ringing/connecting → declined
created/signaling/ringing/connecting → cancelled
created/signaling/ringing/connecting → missed
created/signaling/ringing/connecting → failed
connected → reconnecting → connected/failed/ended
```

Call, notification and booking states stay separate. A missed call does not cancel or decline a booking.

### Required technical/product behaviour

- authorization rechecks relationship, block/restriction, staff membership/scope and rate limits before signaling;
- ringing identifies verified business/display context and automation/recording status;
- poor-network adaptation prioritizes understandable audio, lowers bitrate and communicates reconnecting/ended truthfully;
- permission denial for microphone/notifications has a clear settings/retry path and does not block text booking;
- calls support speaker/earpiece/Bluetooth routing, mute, interruption and background behaviour appropriate to Android/iOS/desktop releases;
- signaling is authenticated, replay resistant and bounded by session expiry;
- relay/media provider regions, encryption, metadata, IP exposure, cost and subprocessor terms are assessed;
- emergency calls and guaranteed carrier reachability are explicitly excluded;
- missed-call notification exposes only safe context.

### Recording and transcription

Recording is off by default and not part of the initial call release. Enabling recording/transcription requires a separately approved purpose, jurisdiction/party-consent analysis, prominent in-call indication, participant controls, provider/data location, encryption/access, retention/deletion/export, staff policy, quality/error review and dispute process.

No hidden recording is permitted. A booking note or AI summary cannot be presented as a verbatim record and must preserve provenance/review.

### Call safety

Before calling launches, Ubi must operate:

- block/report and contact-preference enforcement;
- frequency/time-of-day/relationship limits;
- spam/harassment/impersonation investigation;
- staff identity/offboarding and business accountability;
- minimum metadata needed for security/support with approved retention;
- escalation/redress and emergency limitations;
- capacity/cost/quality monitoring by network class without collecting unnecessary content.

## Accessibility and localization

- Notification copy follows the recipient’s approved language independent of business default where a reviewed template exists.
- A notification must make sense when read aloud without relying on color/icon/visual layout.
- Push/SMS action links have descriptive text and do not require rapid transcription.
- Critical in-app notices receive screen-reader announcements once and remain in a navigable notification center.
- Read/unread cannot be color alone.
- Dates show exact calendar/date/time/timezone according to the regional specification.
- Calls expose accessible labelled controls, focus order, large targets and non-audio state indicators.
- Captions for live calling require separate feasibility/quality review; availability cannot be implied before implementation.

## Privacy and data minimization

- Lock-screen, SMS and external-channel content is minimized according to shared-device and interception risk.
- Provider payloads/endpoints, template variables and delivery logs each have purpose/classification/retention.
- Push providers receive opaque routing and minimal display data where architecture permits.
- Recipient endpoints are not sold/shared or reused for unrelated promotion.
- Analytics records purpose/state/latency/error with privacy-safe IDs, not full message text.
- Message/call metadata and content have separate access and retention.
- Blocking/suppression may require retaining a minimal protected identifier so deletion does not re-enable unwanted contact; policy/legal review defines it.

## Permissions and audit

Separate capabilities cover:

- view/act on notifications;
- configure operation reminder policy;
- edit/approve localized transactional templates;
- send manual one-to-one service communication;
- create promotional campaign (future/controlled);
- manage endpoints/preferences/consent;
- view delivery status and failures;
- message/reply/assign/moderate/report;
- initiate/receive calls;
- configure call availability;
- access message/call metadata/content;
- configure providers/credentials;
- privileged Ubi support or abuse investigation.

Audit records policy/template/provider changes, manual sends, preference/consent changes, suppression, security communications, sensitive content access, block/report outcomes and call authorization/lifecycle where required. Ordinary high-volume delivery telemetry need not become privileged audit content but must remain operationally traceable.

## Invariants

1. Notifications report authoritative states; they never create those states by wording alone.
2. Delivery/read never proves customer acceptance, service completion or payment.
3. Notification failure cannot roll back a committed booking.
4. Duplicate event/retry cannot create duplicate logical communication beyond explicit recurrence/escalation.
5. Stale reminders are cancelled/superseded before dispatch whenever the current state is known.
6. Transactional communication is not a loophole for promotion.
7. Opening a notification/action link never bypasses current authentication, authorization or domain validation.
8. Native free text cannot impersonate a structured quote, confirmed booking or verified payment.
9. Blocking/restriction is enforced before native message/call delivery.
10. Calls are limited to an approved relationship and cannot silently record.
11. Provider status is represented only to its evidenced capability.
12. Critical status remains retrievable inside Ubi/web even if every external delivery attempt fails.

## Release sequence

1. **Foundation:** domain-event notification intent/outbox, in-app business action queue, versioned localized templates and delivery evidence model.
2. **Pilot candidate:** push for business operations plus one approved customer channel for receipts/reminders; failed critical delivery becomes a staff task.
3. **Expansion:** multi-channel preference/fallback and merchant-configurable bounded reminder policy.
4. **Expansion:** Ubi-native booking-scoped messaging with assignment, block/report, abuse and retention operations.
5. **Expansion:** selected official external-channel delivery/inbox capabilities.
6. **Future gated:** Ubi audio calls after relationship, safety, network, cost and provider validation.
7. **Future separately gated:** calling recording/transcription, open username discovery and marketing automation.

## Acceptance scenarios

1. A confirmed booking emits one customer receipt intent; five worker retries produce at most one logical send per channel/provider idempotency ability.
2. Push fails but the booking remains confirmed and the business sees a critical follow-up task if no customer fallback succeeds.
3. A booking is rescheduled just before an old reminder runs; dispatch rechecks version and suppresses the stale date.
4. A shared lock screen shows “You have a Ubi update,” not sensitive service/payment detail.
5. An SMS link forwarded to another person cannot reveal or change unrelated booking data.
6. Provider offers only `submitted`, not delivered/read; Ubi does not display unsupported states.
7. A customer disables promotion but still receives a legally/permissibly necessary cancellation update; no offer is inserted.
8. A blocked native participant cannot start a new thread/call through username or an old booking link.
9. An offline native message remains clearly local/queued and survives process death; it is not shown as sent.
10. A poor-network call moves through reconnecting/failed honestly and offers the booking conversation/call-back path; the booking is unchanged.

## Open decisions

- Launch notification channels/providers, country coverage, cost, sender identity and deliverability.
- Which customer transactional notices are essential versus configurable.
- Quiet hours, reminder timing and escalation for the pilot vertical.
- Controller/processor/sender responsibility and consent/opt-out/suppression retention.
- Ubi-native message release timing, allowed relationships and content/moderation scope.
- Call provider/build approach, media region, codec/network targets, unit cost and mobile background behaviour.
- Whether read receipts, message editing/deletion and business calling hours are offered.
- Recording/transcription remains unapproved pending separate legal/privacy/safety review.

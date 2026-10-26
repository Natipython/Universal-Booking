# Customer web, contextual links and future application

Status: **Draft — account-free web candidate; native customer app later**  
Updated: 2026-10-04 (Africa/Addis_Ababa)  
Depends on: specifications 11–22, 33–38, 44, policies 50–57

## Customer promise

A person can understand a specific business's service and submit a standard booking or custom enquiry without installing Ubi or creating a password. They must know whether the business has confirmed it, what remains to be done and how to contact a human.

## Entry and destination

- Public profile: `<business>.<primary-domain>`; the domain has not been bought or approved.
- Service link: resolves the correct published service/operation.
- Contextual link: `go.<primary-domain>/h/<opaque-token>`; carries only an opaque scoped grant, not chat text or phone/date parameters.
- Manage link: scoped to one booking and allowed customer actions; distinct from a public business slug or display code.

Use a central verified HTTPS link host as the initial candidate. Merchant subdomains resolve public content; never assume a wildcard automatically associates every subdomain with every app. Domain ownership and per-platform association are implementation dependencies. [Android's App Links documentation](https://developer.android.com/training/app-links), rechecked 2026-10-04, describes verified HTTPS links and web destinations.

When the authorized customer app exists and is installed, a supported verified link may open that app. Otherwise web must complete the task. The business app is not a substitute customer destination. Retaining context through a future installation requires an explicit secure resume flow; do not promise automatic post-install continuation.

## Public business page

Above the fold: business name, operation/location, service choices with meaningful prices, exact confirmation expectation and human contact. A small optional portfolio may show relevant work; never make photos/video load before booking can begin.

Verification indicators must explain what was checked. A Ubi profile is not a guarantee of merchant quality, ownership, licensing or payment safety. No invented reviews, stars, booking counts or urgency claims.

## Standard booking, shortest complete flow

| Step | Customer sees/does | Safety rule |
|---|---|---|
| Service | Choose a package or review one already discussed | “From”/estimated prices clearly differ from a final amount |
| Schedule | Choose eligible date/time and necessary options | Exact date, calendar/time convention and location; no private staff calendar |
| Contact | Enter minimum name/contact where needed | No required email/password; phone verification proportional to risk |
| Review | Service, price/deposit, date/time, place, policy and confirmation condition | Prefill editable; changed price/terms need fresh review |
| Submit | One clear **Send request** or **Confirm booking** action | Label follows actual configured mode; idempotent command |
| Outcome | Request sent / confirmed / checking / action needed | Show next responsible party and safe return access |

Several logical steps may fit on one screen when the chosen service is already known. Do not use a fixed multi-step wizard that forces re-entry of chat information. Back navigation preserves drafts.

## Custom event enquiry

Ask only useful first facts: event type, preferred date or “not sure,” location/service area, approximate scope and contact. Ask guest count or budget only when the merchant genuinely needs it; permit unknown where appropriate. Do not request ID documents, children's images or unnecessary intimate details.

Submission means **Request sent**, not reserved/confirmed. Business clarification and an issued quote show exact scope, exclusions, schedule, validity, price and deposit condition. Customer acceptance identifies a precise quote version. An amended quote requires new review; expired acceptance cannot confirm capacity.

## Conversation handoff: example

1. A customer asks the studio about birthday photography in a supported channel or by phone.
2. Staff/authorized assistance records a draft: requested service, preferred date and approximate group size, with provenance.
3. Ubi issues a scoped contextual link and displays relevant service information.
4. The customer reviews “Details from your conversation,” changes an incorrect date and completes the missing contact/review step.
5. Server revalidates capacity, terms and permission before producing the truthful outcome.

No integration is required: staff may create/send a draft link manually under policy, or the business simply shares its public link. Permission to message via a platform does not imply unrestricted access to all inbox history.

## Link and browser safety

- GET/page preview must not submit a booking, accept a quote, cancel, or irreversibly consume a grant. Link-preview crawlers must not perform customer actions.
- Grant redemption/authenticated POST and CSRF protections must match the actual cookie/token architecture.
- Bearer links can be forwarded. Limit exposed data and actions; require stronger verification for sensitive disclosure or consequential changes. Expiry cannot make a leaked active link harmless.
- Separate display code from authorization secret. Phone similarity alone never merges identities or grants history access.
- Private pages use appropriate no-store/referrer protections; no third-party analytics, ads or untrusted media requests receive the link token.
- Redact tokens from application/edge logs, support views and screenshots.
- Stale, revoked, expired or wrong-scope links fail with an understandable route to the business, without revealing another booking.
- If the app is installed but the wrong account is active, preserve the destination safely and reauthorize; never expose data through navigation alone.

## Payment and recovery

The initial candidate supports no deposit/pay-at-location and separately labelled merchant-reported payment evidence. Digital payments are feature-gated until provider, commercial-role, reconciliation and legal reviews are complete.

A provider return screen or screenshot never produces verified payment. A timeout displays “Checking payment” with a return path; do not suggest paying again until the existing attempt is reconciled. Cancellation and refund remain separate outcomes.

Offline public pages may show previously cached non-private service information with age indicators. Final availability and booking commit need the server. Network failure preserves safe draft input where possible and retries the same logical attempt; local submission is not a confirmed reservation.

## Native messages and calling

Required future native experience: booking-linked thread, notifications, permitted attachments, delivery evidence, report/block, business hours/response expectation, scoped identity and username privacy. A guest web thread, if included, needs its own authenticated/link-authorized participant session and safety review; it is not an open public chat box.

Required calling work: authorized booking/identity relationship, call offer/consent, ringing/answer/decline, mute/speaker/end, missed/failed states, privacy, abuse limits and weak-network tests. Phone/WhatsApp/Telegram contact buttons are external fallbacks, not Ubi-native calls. Calling is not an emergency service. No video or recording is automatically included.

## Local comprehension and acceptance

Apply readable English/Amharic source/translation review, large touch targets, text-plus-color statuses, TalkBack/keyboard paths, optional media, and specified web byte/performance budgets.

Specifically test Ethiopian versus Gregorian dates and Ethiopian versus international clock conventions. Do not infer a time convention from interface language. An ambiguous “2 o'clock” from a conversation requires clarification and an explicit equivalent before commitment; a display conversion cannot change the stored instant.

Required observations: new customer completes standard booking, edits a prefill, submits a custom enquiry, distinguishes request/confirmation/payment, recovers after interruption, understands cancellation, uses human fallback and handles a forwarded/expired link safely. Thresholds and evidence are controlled in specifications 71–72 and 83.

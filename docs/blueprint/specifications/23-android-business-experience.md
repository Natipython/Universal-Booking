# Ubi Business Android experience

Status: **Draft — screen/workflow contract, not a built or user-tested interface**  
Updated: 2026-10-04 (Africa/Addis_Ababa)  
Depends on: specifications 10–22, 33–35, 44, 60–65  
Owners: Product/Design, Android, Backend/Data, Security, Operations

## Scope and first-use model

Android is the first business production client. Ethiopia is decided; Addis Ababa and photo studios remain proposed. The default setup presents one business and one booking operation. Additional branches/operations appear only when needed; internal organization concepts must not become a first-use vocabulary test.

The public customer experience is separate. A customer link must not open the business app into a staff workflow or require business membership.

## Navigation and screen inventory

| Area | First view and primary action | Secondary work |
|---|---|---|
| Today | Date, upcoming work, attention items; **Add booking** | Explicit last-sync time, local pending count, next-day preview |
| Calendar | List-first daily schedule; **Add booking** | Day/week views, resource filters, unavailable blocks |
| Requests | Items waiting for this staff member/business; **Review request** | Standard approvals, custom requests, expiring quotes |
| Messages | Booking-linked threads; **Open conversation** | Source badge, assignment, human takeover, delivery evidence |
| More | Plain labelled settings/tasks | Services, hours, customers, staff, business link, help, account |

This is a candidate five-destination layout. Test the labels and discoverability before approving it. A disabled or unreleased module is hidden with a truthful feature-availability explanation where needed; never display a fake unread counter or working call button.

Core screens: language selection, sign-in/recovery, organization selection, guided setup, Today, Calendar, Add booking, review, booking details/history, conflict resolution, custom intake, quotation review/issue, thread, services, hours/blocks, staff invitation, connection settings, sync details, help, sessions/logout.

## Assisted onboarding, in order

1. Choose interface language without assuming literacy or English familiarity.
2. Authenticate with the approved phone-first method; support paste and resend limits.
3. Create a business or accept a scoped staff invitation. A phone OTP alone does not prove business ownership.
4. Choose a familiar template. For the candidate photo pilot, offer a standard studio session and a separate custom event enquiry.
5. Enter business name, service location/contact and timezone. Show what becomes public.
6. Review suggested services, exact ETB prices, duration and required resource. Suggestions are examples, not market prices.
7. Set hours, breaks, buffers and known unavailable dates. Record existing bookings before publishing availability.
8. Choose request-to-confirm or instant confirmation for eligible services; explain the consequence in one sentence.
9. Complete a clearly marked test booking, find it and cancel it. Test data cannot occupy live capacity.
10. Owner reviews publication and staff permissions. Social connections are optional and occur afterward.

Permit save-and-resume. A support worker uses authorized onboarding tools, never a shared owner password. Do not collect business documents until an identified verification/legal purpose requires them.

## Manual booking: minimum complete journey

1. Tap **Add booking**; active business/operation stays visible.
2. Choose a service; recent choices reduce searching.
3. Choose exact date/time and required resource; availability includes buffers and existing commitments.
4. Find or enter the minimum customer/contact information. A walk-in without a phone follows the approved no-contact policy; do not invent a phone number.
5. Select source: call, walk-in, external message, or other. No social account connection is necessary.
6. Review service, schedule, location, price, payment condition and customer-notification choice.
7. Save one durable intent. Online confirmation needs a committed server result; offline saving creates a local tentative item.
8. Show the actual outcome and a clear route back to Today or details.

Customer communication after staff entry needs an approved endpoint/purpose. Staff must not silently enroll a walk-in for promotional messages.

## Status language and visual semantics

| Evidence/state | Candidate plain-language text | Treatment |
|---|---|---|
| Durable local command only | Saved on this phone. Waiting to send. | Neutral/sky-blue; not green |
| Server request accepted, not confirmed | Request sent. Waiting for the business to approve. | Informational; identify next actor |
| Capacity/version conflict | This time is no longer available. Your details are saved. | Attention; alternatives and review |
| Authoritative confirmed booking | Booking confirmed. | Green with text/icon and exact receipt |
| Customer payment claim | Payment reported. Not verified yet. | Informational; no verified receipt |
| Unknown server outcome after timeout | Checking whether this was saved. Do not add it again. | Reconcile the same command ID |

The English text is source copy for translation, not approved Amharic. Each phrase needs comprehension testing.

## Booking details and changes

Details show service, customer contact where authorized, exact schedule, location, booking status, separate payment status, source and history. Present only actions allowed in the current lifecycle and permission scope.

Reschedule previews both old and proposed dates. The old committed allocation remains until an atomic successful move. Cancellation previews policy, customer notice and any separate refund obligation. No refund is implied by cancellation. Completing/no-show marking records actor and reason where required; undo uses an audited domain action, not history deletion.

Custom enquiries remain outside the confirmed schedule unless an explicit hold exists. The owner reviews requirements and an immutable quote version before conversion. Staff cannot replace a custom quote with an unreviewed package price.

## Messages and calls

Booking-linked native messaging and in-app calls remain required product work. Messaging screens must distinguish native and external sources and prevent a reply from going to the wrong channel. Human takeover persists server-side; provider disconnect cannot erase booking history.

Call UI is specified for the later calling slice: booking/thread context, caller identity, consent/availability, answer/decline, mute/speaker/end, missed/failed state, report/block. A **Call business phone** link is an external phone fallback, never represented as a completed Ubi call. Recording is off by default; a separate approval is required to add it.

## Offline, privacy and shared phones

- Cached Today loads without a network request. Show age of data and pending work.
- Only approved offline commands are queued; ownership, refunds, social authorization and publication require current online authority.
- Outbox writes and local projection updates commit atomically. Process death preserves committed commands and recoverable drafts.
- Revocation cannot be detected instantly while disconnected. Define a bounded offline-access lease and sensitive cache horizon before real-data use; cached permission never authorizes server execution.
- Organization/account switching must isolate cached records, cursors and pending commands.
- Logout explains pending work and safe recovery options. It cannot upload another person's commands or leave readable private records to the next user.
- Lock-screen notices minimize customer details; contacts/media permissions are requested only for the selected optional task.

## Visual/performance acceptance

White surfaces, purple primary actions, sky-blue information and green authoritative success follow the semantic tokens. Use at least 48 dp targets, labels with icons, readable dates and scalable text. Price/status/date must not be truncated. Test TalkBack, large fonts, keyboard obstruction, reduced motion and Ethiopic-script rendering.

The current sky-blue token `#218ACB` has about 3.78:1 contrast against white: do not use normal-sized white text on that background. Use a tested darker informational/action variant or tested dark text while retaining sky-blue accents. The subtle border token alone is not sufficient for an essential control boundary on white. Tokens are provisional; every component/hover/focus/error/disabled state needs its own review.

Use specification 21/62 budgets; test the release/profile build on physical D1 devices and N1 networks. Lazy lists, bounded rebuilds and simple effects are implementation candidates supported by [Flutter's performance guidance](https://docs.flutter.dev/perf/best-practices), rechecked 2026-10-04. They do not establish Ubi performance before measurement.

Required proof: manual booking and custom enquiry comprehension, cold/warm startup, process-death survival, retry convergence, cross-account cache isolation, conflict explanation and no false green success. Evidence belongs in specification 72.

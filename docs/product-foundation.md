# Ubi product foundation

Status: **Archived preliminary product input — not approved release scope**  
Archived: 2026-10-04 (Africa/Addis_Ababa)

See [the active blueprint](blueprint/README.md), [decisions](blueprint/02-decision-register.md) and [candidate release contract](blueprint/specifications/72-first-release-contract-and-evidence-gates.md). Geography/vertical, language, provider sequence and calling timing stated below are historical proposals unless explicitly decided in the current register. Ethiopia and Android-first are decided; native messaging/calling remain required direction. Supported social inbox access must be verified per official provider/account; no unconditional access to all incoming chats is promised.

## Historical product notes

## 1. Problem

Small and medium service businesses lose time and trust when bookings are split across calls, WhatsApp, Telegram, Instagram, Facebook, TikTok, notebooks, and memory. Customers repeatedly ask the same questions, businesses can double-book, deposits are hard to verify, and neither side has one reliable record.

Ubi turns those fragmented conversations into one clear booking record while preserving the human ways people already communicate.

## 2. Initial market boundary

Ubi is not launching as a marketplace or a hotel property-management system.

- Initial operating model: business-to-customer booking infrastructure.
- Initial geography: one city and one country before regional expansion.
- Initial vertical: photography studios.
- Next validation vertical: salons and barbers.
- Later vertical: event halls and venues.
- Hotels and regulated healthcare stay outside the early product.

Photography is the first template because it validates packages, add-ons, deposits, staff, locations, event dates, variable durations, and manual bookings without introducing hotel-grade inventory complexity.

## 3. The two products

### Ubi Business

An Android-first operating tool for the owner and staff. It owns the calendar, availability, customers, services, payment state, and all incoming booking conversations.

### Ubi customer experience

A lightweight web app opened from a business subdomain or secure handoff link. It requires no installation and no account before booking. A native customer app comes later when repeated use justifies it.

## 4. Core journeys

### Business owner

1. Choose a familiar business type.
2. Accept suggested services and edit prices.
3. Set ordinary working days and exceptions.
4. Share a business link or connect an approved messaging channel.
5. Add calls, walk-ins, and chat requests to the same calendar.
6. Accept, reschedule, complete, or cancel a booking.
7. Track whether a deposit is due, reported, or verified.

### Customer

1. Open a business link or contextual handoff link.
2. Choose a service or confirm the service already discussed.
3. Choose a date and an available time.
4. Enter a name and phone number.
5. Review price, deposit, place, and cancellation rules.
6. Confirm and receive a simple booking code.
7. Return through the private link to message, reschedule, or cancel.

## 5. Ubi Connect

Social message management is optional, not a requirement for receiving value from Ubi.

- Telegram is the first bot integration because it offers the most practical official path.
- WhatsApp follows through the official business platform and approved templates.
- Instagram and Facebook follow through official Meta business APIs where account eligibility permits.
- TikTok begins as a profile/link handoff until official messaging access is commercially and technically dependable.
- No integration may impersonate a human or silently make irreversible booking decisions.

A conversation handoff creates a short-lived opaque token. The token opens the native app when installed and otherwise opens the web app. Prefilled details remain editable, and sensitive conversation text is not placed directly in the URL.

## 6. Messaging and calls

In-app messaging and calling are part of the long-term product and data model.

Delivery sequence:

1. Booking-linked text messages, attachments, delivery state, and human handoff.
2. Push notifications, moderation/reporting, retention controls, and reliable offline outbox.
3. Voice calls associated with a booking or verified business/customer identity.
4. General username discovery only after abuse prevention and privacy controls prove reliable.

Calls are not placed in the first MVP because signaling, media quality, background execution, safety, and telecom costs would distract from proving the booking engine.

## 7. Usability contract

Every flow must satisfy these rules:

- Use ordinary words such as “Add booking,” never internal system terminology.
- Show one main decision at a time.
- Prefer choices, recent values, and sensible defaults over typing.
- Keep important actions reachable with one hand.
- Use a minimum 48 dp interactive target and never rely on color alone.
- State consequences before destructive or financial actions.
- Let users go back without losing entered information.
- Keep dates human-readable: Today, Tomorrow, Saturday, followed by the exact date.
- Support Ethiopian and Gregorian calendar presentation without storing ambiguous dates.
- Show sync state in plain language and never claim data is saved before it is durable locally.
- Keep a visible human call/message option when automation cannot answer confidently.
- Never require an email address for an ordinary customer booking.

## 8. Language and accessibility

- English and Amharic are the first supported UI languages.
- Oromo and Swahili are next based on pilot geography.
- Store message keys rather than embedding business logic in translated text.
- Support system font scaling without clipped controls.
- Screen-reader labels describe action and state.
- Icons accompany text; icons do not replace unfamiliar words.
- Error messages explain what happened and what the person can do next.

## 9. Performance contract

“Telegram-fast” is treated as a measurable engineering requirement, not a visual effect.

- App shell usable from local data immediately after launch.
- No network request blocks access to the cached Today schedule.
- Lists remain smooth on representative low-cost Android devices.
- Images are optional, compressed, sized to display, and never block booking steps.
- Changes are written locally first and synchronized through an observable outbox.
- Duplicate retries are safe through idempotency keys.
- Customer pages ship minimal JavaScript and work on slow or interrupted connections.
- Startup, frame timing, API latency, failed syncs, and crash-free sessions are monitored.

## 10. Trust and safety

- Business and customer phone numbers are verified before privileged actions.
- Staff have explicit roles; not every staff member can change money or business settings.
- Booking and payment-state changes produce an audit record.
- Deposit state distinguishes customer-reported payment from provider-verified payment.
- Private links expire, can be revoked, and reveal the minimum necessary information.
- Calls and messages include block, report, consent, and retention controls before broad release.
- Ubi does not describe a booking as confirmed until the configured business rule is satisfied.

## 11. Product language

- Brand: **Ubi**.
- Business app: **Ubi Business**.
- Customer app, when justified: **Ubi**.
- Business pages: `business-name.<primary-domain>`.
- Secure handoff links: `go.<primary-domain>/h/<opaque-token>`.

The domain remains configurable until ownership and trademark checks are complete. Product code must not hard-code `ubi.click`.

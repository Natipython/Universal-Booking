# WB-04 — synthetic customer handoff contract V0

Status: **Experimental, local-only; no real customer identity, live bookings or production approval**  
Date: 2026-10-11 (Africa/Addis_Ababa)

## Goal and boundary

A customer can continue a prepared conversation draft without installing an application: open → explicitly continue → correct details → review → submit → understand the outcome. Use the invented WB-01 records only. Business Android 11+/4 GB is not a minimum for this browser flow. No domain, social account, phone contact, payment, analytics or native-app association is enabled.

Server-rendered HTML/forms are a bounded low-byte feasibility choice, not replacement of the proposed production React/TypeScript stack. The core path has no JavaScript, third-party assets or external fonts. English/Gregorian/international-24-hour fixture input is explicit; no Amharic/Ethiopian-calendar/clock completeness is implied.

## HTTP and state contract

All endpoints are served only on a configured loopback origin. Host/Origin checks use that configured origin, never forwarded headers. An exact Origin is accepted; `Origin: null` is accepted **only** with `Sec-Fetch-Site: same-origin`, `Sec-Fetch-Mode: navigate` and `Sec-Fetch-Dest: document`, plus the same issued cookie and secret CSRF checks. Missing Origin, other origins, same-site/cross-site metadata and incomplete null-origin metadata fail closed. The [Fetch Standard, Origin header algorithm](https://fetch.spec.whatwg.org/#append-a-request-origin-header), checked 2026-10-11, explains why `no-referrer` HTML form navigation can send `null`; this exception preserves no-referrer without accepting arbitrary null-origin submissions. Older browsers without the required metadata still need a separately reviewed compatibility fallback.

POST bodies are bounded URL-encoded forms; unknown/duplicate fields, cross-origin submissions and invalid/missing CSRF fail without mutation. Cookies are HttpOnly, SameSite=Lax, host-only and short-lived. `Secure` is on for HTTPS; HTTP cookies use explicitly demo-only names on loopback, never production `__Host-` names.

| Method / route | Input and authority | Result |
|---|---|---|
| GET/HEAD `/` | None | Public invented studio/service facts, no private context; demo entry buttons |
| POST `/demo/starts` | Bootstrap CSRF, allowlisted fixture scenario | Mint one opaque draft grant and 303 to `/h/{token}`; not a merchant/staff API |
| GET/HEAD `/h/{token}` | Opaque locator only | Generic safe preview + Continue form, even for unknown/revoked/expired tokens; no grant redemption, private draft or booking effect |
| POST `/v0/handoff-sessions` | Token, bootstrap cookie and matching CSRF | Bind an unexpired/unrevoked `prepare_booking` grant to this browser, 303 `/continue`; same-browser redemption retry recovers the same session, different browser fails generically |
| GET/HEAD `/continue` | Session cookie; grant checked again | Editable bounded draft or existing outcome; private no-store projection only |
| POST `/v0/reviews` | Session/CSRF, expected draft version, fixed service, explicit calendar/clock, date/time, dummy customer/event size | Save edited draft and immutable current-price review, 303 `/review`; validation failure retains entered safe values; no capacity reservation |
| GET/HEAD `/review` | Session + current review | Exact business/operation/service/time/price or enquiry expectation and one explicit submit button; no mutation |
| POST `/v0/submissions` | Session/CSRF, owned current unexpired review ID | Instant fixture calls WB-03 transactional booking with one server-generated stable command key; custom fixture saves an in-memory demo enquiry without capacity; 303 `/outcome`; same review retry recovers one effect |
| GET/HEAD `/outcome` | Session; scoped result only | Confirmed-demo / request-saved-demo / checking / action-needed truth; no account history or manage actions |

Protected form routes must not accept tenant/operation/principal/price/end-time from the client as authority. The grant fixes organization, operation, service, purpose and confirmation mode. An untrusted proposed date is editable; submitted date/time/customer label is validated. Service changes outside grant scope require a new public flow. Prices/durations/version come from the server, never a posted “confirmed” flag.

## Grant and session rules

- Grant secrets are 32 random bytes encoded as base64url. Only SHA-256 digests are stored. URLs contain no raw name, phone, date, message, price or transcript.
- Preview GET does not validate/consume the grant or reveal validity/private context. Bootstrap cookie/CSRF issuance is technical browser state, not a business action.
- First explicit redemption records browser binding and a server-derived opaque session. Session cookie raw values are not retained in the store; retry derives the same value using a private per-process HMAC key and the same bootstrap binding.
- A copied link before redemption can still be claimed by another person. Browser binding is **not customer identity verification**; only the minimal invented service/date draft is exposed, never contact/chat/history/manage authority. Stronger identity verification is required for consequential/private production flows.
- Expiry, revocation, scope, purpose and browser/session binding are checked again on every protected read/write. Redemption does not grant permanent access or refresh grant expiry.
- Grant/session/review records are memory-only in V0; a process restart loses continuation. No production persistence, contact OTP, token recovery, cross-device or post-install handoff is claimed.
- No-store, no-referrer, restrictive CSP, frame denial, nosniff and noindex apply. No token/raw URL/request body logging, third-party requests, open redirects or cross-origin permissions are enabled.

## Outcome and retry truth

The instant fixture is connected through a narrow adapter to the existing WB-03 Go/PostgreSQL service, in a newly created private temporary cluster only. Its successful result is `confirmed_demo`: an invented reservation, no real business action or message. The quote/custom lifecycle beyond a minimal enquiry remains in WB-03/future work; this browser slice does not pretend that an enquiry is accepted or paid.

Each reviewed attempt gets a stable backend command key. A network/storage-unknown result must preserve that review/key and show **Checking demo booking**, not a green success or a new command. Known capacity/price/version rejection keeps the editable draft and requires a new review. Duplicate submissions and lost-success retries must return the same outcome with one booking effect. Editing after a confirmed/request-saved outcome cannot implicitly create another booking. Old reviews cannot submit after a newer draft review.

After a review expires, an unknown attempt can perform **read-only receipt recovery**, never create a new booking. If no committed receipt is visible, it remains checking for review/support; absence at one instant is not permission to make another attempt. Grant expiry/revocation still blocks the session. Production must define durable recovery beyond these ephemeral demo sessions. Confirmation must match the reviewed scope, service version, price, currency, customer, start and duration—not merely a successful status.

## Required checks

Preview/crawler/HEAD safety; token digests; same-browser redemption retry and forwarded-link negatives; wrong purpose/scope; expiry/revocation after redemption; CSRF/origin/Host/duplicate/oversized-field negatives; HTML escaping; preserved edits; stale/expired reviews and changed terms; no forged tenant/price/customer fields; instant/custom outcome distinction; duplicate submission; lost committed outcome retry; actual WB-03 adapter; private cache/referrer/CSP/cookie headers; keyboard/large-text/narrow-view/no-JS browser path; measured document/CSS byte sizes and no third-party requests.

No full SCN-007–011 pass, weak-network/phone speed, real identity, production cookie/TLS review, domain/app link verification, legal/provider readiness or G7/G8/G9 passage follows from this prototype.

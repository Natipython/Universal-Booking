# Ubi architecture

Status: **Archived preliminary sketch — do not implement these simplified lifecycles**  
Archived: 2026-10-04 (Africa/Addis_Ababa)

Active planning authority is [the controlled blueprint](blueprint/README.md). In particular, [standard workflows](blueprint/specifications/15-standard-booking-workflows.md), [custom quotations](blueprint/specifications/16-custom-request-and-quotation-workflows.md), [offline authority](blueprint/specifications/35-offline-sync-and-conflict-resolution.md) and [the proposed language ADR](blueprint/decisions/ADR-002-language-runtime-and-native-code-strategy.md) supersede this sketch for planning.

The historical `accepted`/green rule below is not current: approval does not mean confirmed until all configured capacity/terms/payment conditions succeed. A local saved intent is tentative, not confirmed. The old TypeScript-or-Go backend choice is superseded by proposed Go, still subject to feasibility evidence. Nothing here proves an app exists.

## Historical sketch

## Guiding decisions

1. The phone remains useful while offline.
2. Server time and storage use unambiguous UTC instants; the UI localizes dates and calendars.
3. Every write is idempotent and auditable.
4. Booking state, payment state, and message state are separate state machines.
5. External messaging and payment providers connect through adapters, not core business logic.
6. Customer handoff URLs contain opaque tokens, not personal or conversation data.

## Planned system shape

```text
Ubi Business Android ─┐
Customer Web App ─────┼── API gateway ── Booking service ── PostgreSQL
Ubi customer clients ─┘       │               │
                              │               ├── Availability engine
                              │               ├── Outbox/event workers
                              │               └── Audit log
                              ├── Identity/OTP
                              ├── Messaging adapters
                              ├── Payment adapters
                              └── Object storage/media processing
```

## Android application layers

```text
presentation  Screens, accessible controls, localized copy
application   Booking use-cases and validation
domain        Entities, value objects, state transitions
data          Local database, sync outbox, API adapters
platform      Notifications, deep links, secure storage, connectivity
```

The first vertical slice uses a replaceable local repository so product behaviour can be tested before the backend exists. The repository contract will later be backed by an encrypted local database plus synchronization.

## Core booking states

`draft → requested → accepted → completed`

Side paths:

- `requested → declined`
- `requested|accepted → cancelled`
- `requested|accepted → rescheduled` creates an audited schedule change rather than hiding history.

“Approved” and green success treatment are reserved for the accepted state. A requested booking uses informational styling.

## Payment states

`not_required | due | customer_reported | processing | verified | failed | refunded | disputed`

The application must never turn a screenshot or customer claim directly into `verified`.

## Sync model

- The UI reads from local storage.
- User actions commit a local transaction and append an outbox command.
- A background synchronizer retries with stable idempotency keys.
- The server returns a version and canonical state.
- Conflicts affecting time or money are shown for explicit resolution; silent last-write-wins is prohibited.

## Multi-tenancy

Every business-owned record carries a tenant identifier. Authorization is enforced on the server, not inferred from a subdomain or client-supplied identifier. Public business pages resolve a verified slug to the tenant.

## Date and calendar rules

- Persist UTC instant plus IANA timezone for timed bookings.
- Persist a local service date for all-day events where appropriate.
- Format Gregorian or Ethiopian calendar presentation at the edge.
- Availability calculations use the business timezone and retain daylight/timezone history.

## Deep links

- Android App Links use verified HTTPS domains.
- The same URL works in a browser when the app is absent.
- Handoff tokens are short-lived, single-purpose, revocable, and exchanged server-side.
- The receiving client previews prefilled fields and requires customer confirmation.

## Initial technology direction

- Ubi Business: Flutter/Dart, optimized and tested on Android first.
- Customer web: server-rendered TypeScript web application with minimal client JavaScript.
- Data: PostgreSQL as the source of truth; Redis only for bounded cache/rate-limit/ephemeral jobs.
- API: modular TypeScript or Go service selected after the Android domain contract stabilizes.
- Media: object storage plus asynchronous resizing and malware/type validation.

The architecture deliberately avoids microservices at the beginning. A modular monolith and reliable worker process are easier to operate and can be split only when real load or team boundaries demand it.

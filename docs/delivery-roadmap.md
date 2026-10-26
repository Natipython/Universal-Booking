# Ubi delivery roadmap

Status: **Archived preliminary roadmap — not the current build/release contract**  
Archived: 2026-10-04 (Africa/Addis_Ababa)

Use [the dependency-ordered backlog](blueprint/specifications/85-release-roadmap-and-implementation-backlog.md), [candidate release contract](blueprint/specifications/72-first-release-contract-and-evidence-gates.md) and [current scorecard](blueprint/specifications/83-kpi-scorecard-and-experiment-governance.md). The five-business/four-week and 70% task-completion values below are historical proposals, not current requirements or passed gates. The active proposed pilot is ten businesses over 8–12 weeks, pending founder/research confirmation; Ethiopia is decided, city/vertical are not. Native messaging and voice calling remain required platform direction with explicit workstreams.

## Historical notes

The word “entire” describes the destination, not one risky release. Each stage must create something a real business can use and must earn the next stage through evidence.

## Stage 0 — Foundation (current)

- Product scope, design tokens, accessibility, performance budgets, and architecture decisions.
- Ubi Business Android shell and first end-to-end local booking flow.
- Automated analysis and widget tests.

Exit: a person unfamiliar with booking software can add a booking and find it again without instruction.

## Stage 1 — Pilot-ready Ubi Business

- Assisted onboarding for photography studios.
- Services, packages, add-ons, duration, price, staff, and working hours.
- Today, calendar, customers, and manual booking.
- Offline local database, sync outbox, audit trail, and conflict handling.
- Phone OTP, owner/staff roles, push notifications, English and Amharic.

Exit: 5 pilot studios run every booking through Ubi for four continuous weeks.

## Stage 2 — Customer booking web app

- Business subdomains and mobile-first service pages.
- Guest booking, private manage-booking links, availability, cancellation, and rescheduling.
- Low-bandwidth rendering and install-free fallback.
- SMS confirmation and configurable cash/deposit instructions.

Exit: at least 70% of observed customers complete a standard booking without staff help.

## Stage 3 — Connected conversations

- Booking-linked in-app messaging and staff inbox.
- Telegram bot and secure contextual handoff.
- Human takeover, message status, attachments, abuse controls, and notification settings.

Exit: connected conversations reduce repetitive staff messages without increasing incorrect bookings.

## Stage 4 — Verified payments

- One launch-country payment provider, webhooks, reconciliation, receipts, refunds, and disputes.
- Fixed, percentage, full, or no-deposit policies.
- Cash and pay-at-location remain first-class options.

Exit: payment state is independently reconcilable and support staff can resolve exceptions.

## Stage 5 — Market expansion

- Salon/barber template and repeat-booking features.
- Event hall template with tentative holds, viewings, capacity, and packages.
- WhatsApp, Instagram, and Facebook integrations where official APIs permit.
- Multi-branch, permissions, reports, and localized payment providers.

## Stage 6 — Native customer app and calls

- Permanent customer identity promoted from verified guest history.
- Cross-business conversations, favourites, and repeat booking.
- Booking-linked voice calls with block/report controls and poor-network adaptation.

## Stage 7 — Discovery marketplace

Only after supply density and merchant retention are proven: location-based discovery, reviews with verified booking context, ranking integrity, and marketplace commercial rules.

## Pilot scorecard

- Median time for staff to add a manual booking.
- Booking completion rate and median completion time.
- Double-booking and lost-change rate.
- Weekly active business rate and four-week retention.
- Percentage of bookings captured from calls/walk-ins/messages.
- No-show rate before and after reminders/deposits.
- Crash-free sessions, cold/warm startup, and failed-sync rate.
- Help requests per completed booking, segmented by language and device tier.

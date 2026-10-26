# Ubi scope and architecture principles

Status: **Draft**

## Product definition

Ubi is a service-booking and business-conversation platform for organizations that currently coordinate customers through phone calls, walk-ins, messaging applications, notebooks, and disconnected payment evidence.

It is designed to convert an enquiry into a reliable operational record while preserving human conversation and local business practice.

## Product promise

### To a business

Manage every booking from one trustworthy place, even when customers call, walk in, use a booking link, or contact a supported connected channel.

### To a customer

Understand the service, choose or request what you need, confirm it safely, and know what happens next.

## Platform boundary

Ubi owns:

- Ubi identities, organizations, staff membership and permissions.
- Business catalog, availability, booking/request/quote records and audit.
- Ubi messages, notifications and future Ubi call coordination.
- Integration adapters, normalized channel conversations and handoff records.
- Payment intent/state records and reconciliation evidence for approved providers.
- Merchant-facing and customer-facing Ubi experiences.

Ubi does not own:

- External social-network accounts, policies, delivery guarantees or historical inbox access.
- Merchant service quality or legal licensing unless Ubi explicitly assumes that role later.
- A payment provider’s settlement network or truth beyond verified provider evidence.
- The customer’s device, network availability or mobile-money account.
- Country-specific professional regulation for merchants.

These boundaries must be stated truthfully in product behaviour, support processes and legal documents.

## Non-negotiable product principles

1. **One human identity, explicit memberships.** People do not share owner credentials.
2. **One organization, many operations.** A real business may contain brands, branches and different booking models.
3. **Universal through building blocks.** Reusable archetypes provide flexibility; arbitrary custom code is not the default.
4. **Standard and custom are different.** A request or quotation never looks confirmed merely because it was submitted.
5. **Conversation is a source, not the database.** Calls and messages produce structured drafts/records with source provenance.
6. **Human control remains visible.** Staff can review, correct, take over and audit automation.
7. **Official integration only.** No password collection, scraping or unofficial user-account automation.
8. **Local-first business operations.** Today’s schedule remains useful offline; changes are durably queued and visibly synchronized.
9. **Server-enforced truth.** Clients assist; authoritative permissions, constraints and money rules are enforced server-side.
10. **Separate state machines.** Booking, payment, conversation, notification and fulfilment states do not impersonate each other.
11. **No silent loss or conflict.** Retries are idempotent and material conflicts are explicitly resolved.
12. **Plain language is correctness.** A technically correct action users cannot understand is a product defect.
13. **Minimum necessary data.** Each field, event and copy has a stated purpose and retention rule.
14. **Audit powerful actions.** Ownership, access, booking, money, integration and support changes are attributable.
15. **Measure speed.** “Fast like Telegram” becomes device, frame, startup, latency, payload and failure budgets.
16. **Design for weak conditions.** Low bandwidth, interruptions, low-cost devices and limited storage are normal conditions.
17. **Accessibility is a release requirement.** Color, language, target size, screen readers, font scaling and keyboard behaviour are tested.
18. **Regional core, country overlays.** Shared platform concepts coexist with country-specific law, payments, language, currency and operations.
19. **Modular before distributed.** Clear boundaries do not require premature microservices.
20. **Operational completeness.** A feature includes support, observability, recovery, policy, cost and ownership—not only code.

## Experience principles for low learning burden

- One primary action or decision per step.
- Familiar business vocabulary, never architecture terminology.
- Suggested templates and recent values instead of blank configuration.
- Large labelled controls; icons reinforce words rather than replace them.
- Phone-first identity; email only where justified.
- Human-readable dates followed by an exact date.
- Prices, deposits, confirmation conditions and cancellation consequences shown before commitment.
- Back navigation preserves work.
- Errors explain what happened, what is saved, and the next safe action.
- Progress is meaningful; unnecessary step counters and animations are avoided.
- Green is reserved for completed success/approval; pending/requested states use neutral or informational treatment.
- Training supplements clarity but never excuses a confusing primary workflow.

## First-release exclusion discipline

Unless later approved through the decision process, the first release excludes:

- Consumer discovery marketplace and public ranking.
- Hotels’ full property-management/channel-manager functions.
- Regulated clinical records and healthcare workflows.
- General-purpose social network or open username messaging.
- Voice/video calls and recording.
- Loyalty, payroll, accounting suite, inventory ERP and marketing automation.
- Unrestricted no-code workflow programming.
- AI that independently approves money, availability, identity or irreversible actions.
- Multi-country launch.

Exclusion from Version 1 does not mean the architecture ignores future boundaries; it means no premature implementation.

## Visual semantics

- White and near-white surfaces keep the interface calm and readable.
- Purple identifies Ubi and the principal action/selection.
- Sky blue identifies helpful information, connection and secondary action.
- Green communicates verified success or approval only.
- Amber communicates attention or a reversible pending risk.
- Red communicates error, danger or destructive consequence.
- Text and iconography always accompany semantic color.

Exact tokens require contrast and device testing before approval.

## Architecture fitness tests

Any proposed architecture must answer yes to these questions:

1. Can a business add a call/walk-in booking while disconnected and later synchronize safely?
2. Can one organization operate two fundamentally different booking operations without duplicate ownership accounts?
3. Can staff permissions prevent a receptionist from changing ownership or issuing refunds?
4. Can an external integration be disconnected without losing Ubi booking history or exposing its token?
5. Can payment evidence remain “reported” until independently verified?
6. Can a customer continue through web when the native app is absent?
7. Can we delete/export a person’s data consistently across operational and analytical systems?
8. Can a duplicate webhook or mobile retry be processed without duplicate money or bookings?
9. Can we restore the service and verify tenant isolation after a disaster?
10. Can we explain every automated action to the business and customer?

If a proposal fails one of these tests, it must change or document an explicitly accepted exception.

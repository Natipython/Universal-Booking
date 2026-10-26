# Ubi controlled glossary

Status: **Draft**

These terms are deliberately separated. Product copy may use simpler words, but architecture, database, policy and analytics must not collapse different concepts.

## Identity and organization

**Person** — a real human being; not necessarily an authenticated Ubi user.

**Account** — a person’s authentication identity, sessions and recovery controls.

**Guest customer** — a customer represented for a booking without a permanent authenticated Ubi account.

**Customer profile** — customer information used by an organization or linked to a verified Ubi account according to privacy rules.

**Organization** — the tenant-level real-world business or operating entity using Ubi.

**Brand** — a public commercial identity operated by an organization.

**Branch** — an organizational operating location or management scope.

**Booking operation** — a configurable unit that owns a coherent booking model, catalog, resources, availability and policy; for example hotel rooms or a wedding hall.

**Membership** — the relationship granting an account scoped capabilities in an organization, branch or operation.

**Role** — a named bundle of capabilities used for convenient assignment.

**Permission/capability** — a server-enforced right to perform a specific action on a defined scope.

**Business public profile** — publicly safe brand/location/contact information presented to customers.

**Slug/username** — a changeable human-readable public identifier; never the authoritative database identity.

## Catalog, capacity and scheduling

**Service** — something a business offers, with rules such as duration, price, location and required resources.

**Package** — a sellable configuration of service components and inclusions.

**Add-on** — an optional priced or capacity-affecting component selected with a service/package.

**Resource** — a schedulable person, room, venue, asset, vehicle or equipment unit.

**Availability rule** — recurring or exceptional rule describing when an operation/resource may be considered.

**Capacity** — the maximum simultaneous quantity or participants allowed by a resource/session.

**Slot** — a computed candidate interval; not itself proof that a booking is confirmed.

**Hold** — temporary, expiring reservation of capacity while a request, quotation or payment completes.

## Commercial lifecycle

**Enquiry** — an expression of interest that may not yet contain enough structured information to create a request.

**Booking draft** — structured but unsubmitted booking information.

**Booking request** — submitted request awaiting a configured decision or condition.

**Custom request** — request requiring business review, scope clarification or quotation before confirmation.

**Quotation** — versioned commercial proposal with scope, price, validity and acceptance state.

**Booking** — authoritative record coordinating a service commitment, schedule, parties and policies.

**Reservation** — booking archetype centered on exclusive/capacity-controlled resource or inventory.

**Appointment** — booking archetype centered on a scheduled service interval, commonly involving staff.

**Fulfilment** — delivery of the promised service; separate from booking and payment.

**Reschedule** — audited change to the promised schedule; not deletion of prior history.

**Cancellation** — termination before fulfilment according to actor, reason, policy and financial consequences.

**No-show** — a confirmed party did not attend/perform as recorded under defined policy.

## Money

**Price snapshot** — the exact commercial terms captured for a request/booking, protected from later catalog edits.

**Deposit policy** — rule determining whether and how much must be paid before confirmation or fulfilment.

**Payment intent** — Ubi record of an attempted collection request; not proof of payment.

**Customer-reported payment** — unverified claim or evidence provided by a customer/staff member.

**Provider-verified payment** — payment state supported by authenticated provider data and reconciliation rules.

**Settlement** — provider movement of funds to the entitled party.

**Reconciliation** — comparison of Ubi records, provider records and settlement outcomes to identify differences.

**Refund** — return of captured/settled value according to provider and policy.

**Reversal** — provider/network cancellation of a previous transaction outcome.

**Dispute** — contested payment or service-related financial claim.

## Conversation and integration

**Conversation** — ordered communication context between participants, potentially associated with a booking/request.

**Channel** — communication medium such as Ubi, Telegram, WhatsApp, Instagram or Facebook Page messaging.

**Channel connection** — merchant-authorized integration between an external business asset and Ubi.

**External identity** — platform-scoped user/business identifier received from an integration; it is not automatically a Ubi account.

**Normalized message** — Ubi representation of a provider message retaining source IDs, capabilities and provenance.

**Human takeover** — explicit transition that stops automated replies and assigns conversation control to staff.

**Handoff token** — opaque, scoped, expiring and revocable credential used to continue context in another Ubi surface.

**Webhook** — provider-initiated authenticated notification to Ubi; it may be delayed, duplicated or delivered out of order and must be handled accordingly.

## Data and systems

**System of record** — authoritative source for a specific category of truth.

**Operational database** — transactional store serving live business operations.

**Local store** — device database enabling fast/offline reads and queued changes; not independently authoritative for synchronized server truth.

**Outbox command** — durable record of a local/server action awaiting processing or synchronization.

**Idempotency key** — stable identifier ensuring safe retries do not repeat a logical effect.

**Audit event** — attributable, purpose-limited record of a security/business-significant action.

**Domain event** — statement that an approved domain state change occurred; not identical to analytics telemetry.

**Analytics event** — governed measurement record used to understand product/operational behaviour.

**Tenant** — authorization and data-isolation boundary, initially corresponding to an organization unless an approved decision changes it.

**Tombstone** — retained deletion marker used for synchronization and propagation without retaining unnecessary deleted content.

**RPO** — maximum acceptable data loss measured in time after a disaster.

**RTO** — maximum acceptable time to restore a service after a disaster.

**SLO** — measurable reliability target for a service or user journey.

**Error budget** — tolerated unreliability within an SLO period used to balance reliability and change.

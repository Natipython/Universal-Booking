# Booking archetype blueprint

Status: **Draft**  
Depends on: business taxonomy  
Blocks: booking workflows, availability engine, domain model, database schema, onboarding and UI

## Purpose

An archetype defines the scheduling and commitment invariants shared by many business categories. A merchant operation selects one primary archetype and a controlled set of modules/policies. Complex businesses create several operations or use an approved composite instead of forcing contradictory behaviour into one form.

## Common domain objects

The following concepts are deliberately separate:

- **Enquiry** — unstructured or partly structured interest.
- **Booking draft** — client/local working data not yet submitted.
- **Booking request** — submitted proposal awaiting a decision or condition.
- **Quotation** — versioned scope and commercial proposal.
- **Hold** — temporary claim on capacity/resources with explicit expiry.
- **Booking** — confirmed service commitment between defined parties.
- **Schedule occurrence** — exact planned time/location/resources for a booking.
- **Fulfilment record** — arrival/start/completion/no-show/service outcome.
- **Payment intent/transaction** — separate money lifecycle.
- **Conversation** — communication context linked to any of the above.

An auto-confirm flow may create a booking directly from a validated draft. A request-to-confirm flow creates a booking request first. A custom flow may require enquiry, request, quotation and hold before creating a booking.

## Shared configuration modules

### Schedule modules

- Exact start with fixed/calculated duration.
- Customer-selected duration from allowed increments.
- Date-only/full-day commitment.
- Time window.
- Multi-day date range.
- Published session/timetable.
- Recurrence rule with occurrence exceptions.
- Setup/travel/cleanup buffers.

### Resource modules

- Named staff preference or required staff.
- Any eligible staff/resource automatic assignment.
- Room/venue/court/station.
- Discrete asset/vehicle/equipment.
- Quantity inventory.
- Several resources required simultaneously.
- Substitutable resource group.

### Capacity modules

- Exclusive capacity of one.
- Participant-count capacity.
- Unit-count inventory.
- Pooled capacity.
- Shared versus exclusive use.
- Waitlist when full.

### Catalog and price modules

- Service.
- Package.
- Add-on.
- Fixed or duration/person/unit/night-based price.
- Starting price or estimated range.
- Quotation.
- Price adjustments/discount authority.
- Tax/fee presentation after country design.

### Location modules

- Fixed branch.
- One of several branches.
- Customer address/service area.
- Event/custom location.
- Remote meeting.
- Pickup/drop-off or origin/destination.

### Commitment modules

- Instant confirmation.
- Staff approval.
- Customer acceptance of proposed option.
- Deposit/payment condition.
- Document/eligibility condition.
- Temporary hold with expiry.

### Policy modules

- Cancellation/rescheduling deadline and consequence.
- No-show rule.
- Minimum/maximum advance booking.
- Blackout/exception dates.
- Booking frequency/quantity limits.
- Late arrival/grace period.
- Refund/deposit rules.
- Customer eligibility/age/guardian rules where approved.

## A01 — Appointment

### Best for

Barbers, salons, consultations, tutors, trainers, inspections, simple photography sessions and other time-bound services.

### Commitment

One customer/party receives a service during an exact interval, usually involving one or more staff/resources.

### Required facts

- Service/package.
- Branch/location mode.
- Start time and calculated/selected duration.
- Customer/party contact.
- Eligible or selected staff/resource.
- Price snapshot and confirmation policy.

### Availability rule

The complete service interval plus buffers must fit business hours, service rules and all required resource calendars without prohibited overlap.

### Variants

- Customer chooses staff.
- “Any available” staff assigned automatically.
- Staff and room/chair/bay required simultaneously.
- Repeating series.
- Request-to-confirm for variable duration.

### Not sufficient for

Weddings, multi-stage projects, hotel stays or repair jobs whose final scope is unknown after intake.

## A02 — Resource slot/reservation

### Best for

Halls, studios, meeting rooms, sports courts, coworking rooms, gaming stations and other bookable spaces/resources.

### Commitment

Exclusive or capacity-limited use of a resource for an interval or event date.

### Required facts

- Resource or resource class.
- Date/start/duration or full-day block.
- Party/participant count when capacity matters.
- Usage/package/add-ons.
- Setup/cleanup/access buffers.
- Hold/confirmation/deposit policy.

### Availability rule

Resource capacity over the complete blocked interval cannot exceed allowed capacity. Combined/divisible resources require explicit compatibility rules.

### Variants

- Instant short room/court booking.
- Full-day venue request with manual approval.
- Tentative event-date hold.
- Shared session versus exclusive rental.
- Room plus equipment/staff.

### Not sufficient for

Accommodation nights, asset pickup/return condition, or an event service whose main constraint is a service team rather than a venue.

## A03 — Event/project service

### Best for

Wedding photography, videography, DJs, MCs, planners, decorators, catering and other services committed to a specific event/project.

### Commitment

A business reserves service capacity/team for an event date/time and agreed scope. Delivery may include preparation and post-event work outside the visible event interval.

### Required facts

- Event type/date/time/location.
- Scope/package/deliverables.
- Expected participants/scale where relevant.
- Team/resources and travel/setup requirements.
- Price/quotation and deposit/confirmation conditions.

### Availability rule

Business/team capacity and travel/setup/other project obligations must be reviewed. An open calendar slot alone may not prove feasibility.

### Variants

- Fixed package with manual confirmation.
- Package plus add-ons.
- Fully custom quotation.
- Multi-stage project with consultation/rehearsal/delivery milestones.

### Not sufficient for

Simple room rental or a service with a known fixed-duration appointment.

## A04 — Stay/unit-night reservation

### Best for

Hotels, guesthouses, resorts, lodges, hostels and vacation properties.

### Commitment

One or more units or beds across a continuous range of nights with check-in/check-out rules and occupancy.

### Required facts

- Arrival/check-in and departure/check-out dates.
- Unit type or unit assignment policy.
- Unit quantity, adults/children/occupants.
- Rate plan and included services.
- Restrictions, guarantee/payment and cancellation policy.

### Availability rule

Requested unit-night inventory must remain available for every night in the stay. Maintenance/out-of-service/overbooking policies are specialized concerns.

### Variants

- Room-type inventory with later unit assignment.
- Exact property/unit booking.
- Bed-level hostel inventory.
- Package with activities/transport as related operations.

### Release restriction

This archetype is documented for platform boundaries but is not approved for the initial release. A full hospitality programme must cover rate plans, taxes, occupancy, housekeeping, channel distribution and overbooking.

## A05 — Class/session/capacity booking

### Best for

Training centres, sports academies, gym classes, pool sessions, group lessons, tours and go-kart sessions.

### Commitment

One or more seats/participants are reserved in a predefined session or repeating cohort with capacity.

### Required facts

- Session/cohort/timetable occurrence.
- Participant count and identities only where necessary.
- Capacity/eligibility.
- Instructor/location/resources.
- Per-seat/package/membership price rule.

### Availability rule

Confirmed plus held seats must not exceed session capacity; instructor/location availability defines whether the session itself can run.

### Variants

- Single session.
- Recurring cohort enrollment.
- Waitlist.
- Membership entitlement.
- Private group booking.

### Not sufficient for

One-to-one appointments or exclusive venue rental unless combined deliberately.

## A06 — Asset rental

### Best for

Vehicles, sound/lighting systems, tents, chairs, tables, tools and equipment.

### Commitment

Specific units or quantities are reserved between pickup/delivery and return/collection, potentially with setup/operator services.

### Required facts

- Asset class/unit/quantity.
- Rental interval.
- Pickup/delivery and return/collection details.
- Eligibility/document requirements.
- Price, booking deposit and potentially separate security/damage deposit.
- Condition/handover process.

### Availability rule

Every required unit/quantity is available across the blocked interval including delivery, cleaning, inspection and maintenance buffers.

### Variants

- Exact serialized vehicle/equipment.
- Quantity pool such as chairs.
- Rental plus operator/setup crew.
- One-way pickup/return locations.

### Important distinction

Booking deposit, payment balance and refundable security/damage deposit are different financial obligations.

## A07 — Field service/dispatch/transport

### Best for

Cleaning, plumbing, electrical work, repairs, chauffeur service and other mobile services.

### Commitment

A worker/team/vehicle visits a customer/event location during an exact time or arrival window, with travel/service-area constraints.

### Required facts

- Service/request type.
- Address/location and service-area eligibility.
- Preferred exact time or time window.
- Estimated duration/urgency.
- Required team/vehicle/tools.
- Price, estimate or call-out fee.

### Availability rule

Worker/resource capacity must include travel time and geographic feasibility, not only service duration.

### Variants

- Customer chooses a window; business assigns arrival.
- Immediate/on-demand request, later programme.
- Recurring cleaning/service.
- Origin/destination journey with passenger/cargo details.

### Not sufficient for

Emergency dispatch promises without 24/7 operational capability and appropriate safety/process design.

## A08 — Viewing/site visit/tour

### Best for

Property viewing, venue viewing, office tours, site surveys and pre-service assessments.

### Commitment

A customer and representative meet at a location for evaluation before a later booking/project decision.

### Required facts

- Property/site/venue.
- Agent/staff.
- Date/time/duration.
- Visitor/party count.
- Safe directions/access and qualification information where justified.

### Availability rule

The location and required representative must be available; occupant/security/privacy restrictions may limit public detail.

### Variants

- Public open-house session.
- Private viewing.
- Remote/video consultation.
- Site survey feeding a custom request/quotation.

### Outcome

May complete without another booking or produce a linked custom request/quotation. It must not automatically commit the main service.

## A09 — Queue/walk-in

### Best for

Barbers, car washes, clinics only after regulated review, tire service and businesses that mix appointments with walk-ins.

### Commitment

A customer joins an ordered or estimated queue rather than reserving an exact guaranteed interval.

### Required facts

- Service/category.
- Arrival/check-in time.
- Customer/party/vehicle identifier sufficient for service.
- Assigned/eligible resource if applicable.
- Estimated wait with clear non-guarantee wording.

### Capacity rule

Queue estimates use active workload and throughput but remain estimates. Appointments may reserve capacity that walk-in estimates must respect.

### Variants

- Physical walk-in entered by staff.
- Remote join/check-in if the business can operationally honor it.
- Number/token display.

### Caution

Do not present an estimated service time as a confirmed appointment.

## A10 — Custom request and quotation

### Best for

Weddings, events, decoration, catering, moving, repair diagnosis, construction/design and any service with unknown scope, price or feasibility.

### Commitment

No booking commitment exists until requirements are reviewed, a scope/option is proposed, the customer accepts, conditions are satisfied and capacity is confirmed.

### Required facts

- Request type and guided intake answers.
- Customer/event/location/timing preferences.
- Attachments only where needed.
- Business review/assignee.
- Versioned quotation or proposed options.
- Hold, acceptance and deposit deadline where applicable.

### Availability rule

Availability is evaluated by staff/rules against the agreed scope. A preferred date can be recorded without falsely reserving it. A separate hold reserves capacity when authorized.

### Variants

- Quote only, no date yet.
- Several proposed schedule options.
- Site visit before quotation.
- Fixed package customized with controlled modules.
- Quote revision/amendment after customer feedback.

Detailed lifecycle is defined in the custom-request blueprint.

## Composite-operation rules

An operation may combine an archetype with subordinate modules when the commitment remains coherent.

Valid examples:

- Appointment requiring staff plus treatment room.
- Venue reservation plus equipment add-ons.
- Event service plus equipment rental.
- Field service preceded by a viewing/site survey.
- Class session with optional equipment.

Prefer separate linked operations/records when commitments have independent lifecycle or inventory:

- Hotel stay and wedding hall.
- Repair intake appointment and later repair quotation.
- Venue viewing and venue reservation.
- Wedding photography event and studio portrait appointment.
- Rental asset and optional staffed event production when either may change independently.

## Configuration safety levels

### Level 1 — Template defaults

Business type produces understandable preconfigured services, durations, policies and role suggestions.

### Level 2 — Safe edits

Owner edits names, prices, durations, availability, staff, add-ons and approved policy choices.

### Level 3 — Advanced modules

Owner enables holds, multiple resources, capacity, recurrence, service areas, quotations or approvals through guided configuration.

### Level 4 — Specialist configuration

Ubi-assisted setup for complex multi-resource/event/rental/hospitality/regulated operations. It remains typed/versioned configuration, not arbitrary code.

## Configuration invariants

1. Every active service belongs to one operation and one primary archetype/version.
2. A service cannot publish until required configuration is complete and validated.
3. Duration, buffers, capacity and required resource rules are snapshotted or version-referenced for existing commitments.
4. Catalog price edits do not rewrite submitted quotations or confirmed bookings.
5. Policy edits do not silently remove customer rights or change existing commitments unless an explicit lawful migration is accepted.
6. Resources from another tenant cannot satisfy availability.
7. A held resource consumes configured capacity until consumed, released or expired.
8. A request/preferred time does not consume capacity unless an explicit hold exists.
9. Auto-confirmation is allowed only when the system can authoritatively validate all required capacity/conditions.
10. Custom fields cannot silently request prohibited/unnecessary sensitive data.
11. Deactivated services remain represented in historical bookings.
12. Archetype/template migrations are versioned and reversible/tested.

## Availability engine implications

The engine must evaluate a request against:

```text
operation publication/eligibility
∩ service booking rules
∩ business/branch hours and exceptions
∩ required resource calendars and buffers
∩ holds and confirmed commitments
∩ capacity/unit inventory
∩ advance/cutoff/timezone rules
∩ location/service-area/travel feasibility when applicable
∩ staff/provider eligibility
```

The output is a candidate offer with an expiry/version, not a permanent guarantee. Final commitment occurs transactionally against current authoritative constraints.

## Archetype evolution

- Archetype IDs and semantic versions are stable.
- New optional fields do not alter old meaning.
- Breaking invariant changes require a new version and merchant migration plan.
- A merchant tests changes before publishing.
- Existing future bookings retain their applicable snapshots/rules.
- Analytics records both business category and archetype/version.

## Approval criteria

- Each selected pilot service fits one archetype without contradictory exceptions.
- Availability and confirmation invariants are implementable transactionally.
- The archetype is understandable through ordinary business-language templates.
- Policy and price changes preserve historical commitments.
- Customization remains bounded and testable.
- Specialized stay, healthcare and complex regulated overlays remain gated.

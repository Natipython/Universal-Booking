# Business taxonomy blueprint

Status: **Draft**  
Depends on: actors, organization and tenancy  
Blocks: booking archetypes, workflow configuration, onboarding templates, catalog, availability, database design and pilot selection

## Purpose

Ubi needs a universal classification system without pretending every business works the same way. The taxonomy classifies an operation by how it commits time, resources, capacity, location and money—not merely by its marketing category.

A business category supplies familiar onboarding language and defaults. The underlying booking behaviour comes from one or more reusable archetypes and modules.

## Classification dimensions

Every booking operation is described using these dimensions.

### 1. Customer intent

- **Immediate selection** — customer can choose a known service/package.
- **Availability request** — business must approve a proposed time.
- **Custom scope request** — requirements and price need discussion/quotation.
- **Viewing/consultation first** — preliminary appointment precedes the main commitment.
- **Walk-in/queue** — service is delivered by order or estimated wait rather than a reserved interval.

### 2. Time commitment

- No scheduled time yet.
- Exact start time.
- Start time plus calculated duration.
- Customer-selected time window.
- Business-assigned time window.
- Full-day or event date.
- Multi-day date range.
- Repeating schedule.
- Session from a published timetable.
- Queue position/estimated wait.

### 3. Capacity model

- One exclusive booking per interval.
- Counted participant capacity.
- Multiple interchangeable units.
- Named resource capacity.
- Pooled capacity across resources.
- Queue/service throughput.
- Unscheduled/custom capacity reviewed manually.

### 4. Resource model

- No named resource required.
- One staff member.
- One room/venue/court/studio.
- One asset/vehicle/equipment unit.
- Several simultaneous resources.
- Alternative eligible resources chosen automatically.
- Customer-owned resource or location.

### 5. Fulfilment location

- At a fixed business branch.
- At a selected business location.
- At the customer’s address.
- At an event/custom location.
- Remote/online.
- Mobile route/journey.
- Pickup and return locations.

### 6. Commercial model

- Free.
- Fixed price.
- Starting/from price.
- Per duration/unit/person/night/distance.
- Package plus add-ons.
- Tiered/seasonal/rate-plan price.
- Estimated range.
- Custom quotation.
- Membership/included entitlement, later module.

### 7. Confirmation model

- Instant confirmation.
- Business approval required.
- Customer confirmation of a business proposal.
- Deposit/payment required before confirmation.
- Verification/documentation required.
- Multi-party approval.
- Tentative hold followed by a deadline.

### 8. Fulfilment pattern

- One-time service.
- Repeating appointments.
- Multi-stage project.
- Check-in/use/check-out.
- Pickup/use/return.
- Arrival/queue/service.
- Class/cohort attendance.

### 9. Risk and regulatory overlay

- Ordinary service.
- High-value deposit/dispute risk.
- Safety-sensitive location/service.
- Minor/guardian involvement.
- Professional licensing.
- Health/sensitive personal information.
- Aviation/transport or other regulated operation.

The risk overlay changes verification, consent, retention, insurance and support obligations without creating a completely unrelated scheduling engine.

## Business family map

The identifiers `A01`–`A10` refer to the archetypes defined in the next blueprint.

### Photography and media

| Business | Primary archetype(s) | Important modules/variations | Practical caution |
|---|---|---|---|
| Photography studio | A01 Appointment, A02 Resource slot, A03 Event/project, A10 Custom request | Packages, photographer, studio room, add-ons, location, deposit | Passport sessions and weddings need very different defaults. |
| Wedding photographer | A03 Event/project, A10 Custom request | Event date, coverage hours, travel, team, deliverables, quotation, deposit | Date may need a temporary hold while scope/contract is agreed. |
| Videographer | A03, A10 | Crew/equipment, event schedule, deliverables, revisions, travel | A single “time slot” cannot represent production scope. |
| Drone operator | A03, A10 | Location, permits, weather contingency, safety review | Regulatory/airspace approval cannot be represented as ordinary availability. |
| Photo/video production company | A03, A10 | Multi-stage project, crew, equipment, locations, quotation milestones | Requires project phases later; first release can support request/quote/booking commitment. |
| Podcast studio | A02 Resource slot, A01 Appointment | Room, engineer, equipment, duration, add-ons | Setup/cleanup buffers affect resource availability. |
| Recording studio | A02, A01, A10 | Studio/engineer, hourly blocks, equipment, production package | Long/custom sessions may need approval; short studio rental may be instant. |
| Content studio | A02, A01 | Studio set, lighting/equipment, creator assistance | Resource combinations and buffers matter. |

### Beauty and grooming

| Business | Primary archetype(s) | Important modules/variations | Practical caution |
|---|---|---|---|
| Barber | A01 Appointment, A09 Queue/walk-in | Barber preference, fixed duration/price, repeat booking | Walk-ins and appointments must share staff capacity. |
| Hair salon | A01, A10 for complex work | Stylist, variable duration, service combination, add-ons | Service combinations and hair assessment can change duration. |
| Nail salon | A01 | Technician, service type, station, add-ons | Staff and station may both be required. |
| Makeup artist | A01, A03, A10 | Location, event date, travel, package, trial session | On-location travel and bridal work differ from in-studio appointments. |
| Braiding salon | A01, A10 | Style choice/reference image, duration range, stylist | Exact duration may depend on style/hair assessment. |
| Spa | A01, A02 | Therapist, treatment room, package, buffers | Requires simultaneous staff and room resources. |
| Massage business | A01 | Therapist, room/location, duration | Privacy, safety and practitioner rules require market review. |
| Beauty clinic | A01 with regulated overlay | Practitioner, room, intake/consent | Do not launch as ordinary salon flow until health/privacy obligations are reviewed. |
| Bridal makeup service | A03, A10 | Event schedule, party size, assistants, travel, package, deposit | Multi-person schedule and location logistics need custom scope. |

### Events and venues

| Business | Primary archetype(s) | Important modules/variations | Practical caution |
|---|---|---|---|
| Wedding hall | A02 Resource/event-date, A08 Viewing, A10 | Date hold, capacity, packages, add-ons, deposit, viewing | A hold is not confirmation; cleanup/setup blocks must be modeled. |
| Conference hall | A02, A08 | Time block, layout, capacity, equipment, catering | Half-day/full-day and room-combination rules may apply. |
| Banquet hall | A02, A08, A10 | Guest count, catering/decoration packages, date hold | Capacity and commercial scope must both be confirmed. |
| Party venue | A02, A08 | Date/time, capacity, package, noise/age policy | Damage/security deposits may differ from booking deposits. |
| Meeting room | A02 | Room, duration, equipment, participant capacity | Often suitable for instant booking. |
| Wedding planner | A03, A10 | Event date, planning scope, milestones, suppliers | Booking reserves service capacity, not all third-party suppliers. |
| Event planner | A03, A10 | Event type, scope, date, location, suppliers | Quotation and change control are essential. |
| DJ | A03, A10 | Date, duration, location, equipment, travel | Event availability plus equipment/transport. |
| MC | A03, A10 | Date, schedule, language, location, rehearsal | Custom event brief. |
| Decorator | A03, A10 | Venue, theme, quantity, setup/teardown, materials | Images and quotation versions are common. |
| Catering company | A03, A10 | Guest count, menu, location, delivery/service, dietary needs | Capacity is production/logistics-based, not a simple calendar slot. |
| Sound-system rental | A06 Rental, A03/A10 with operation | Equipment bundle, operator, delivery, setup, deposit | Inventory and service crew may both be required. |
| Lighting rental | A06, A03/A10 | Equipment, venue, setup/operation, return | Asset condition and delivery/return are separate events. |
| Tent/chair/table rental | A06, A10 | Item quantities, delivery, setup, event date, return | Quantity inventory replaces one-resource exclusivity. |

### Hospitality

| Business | Primary archetype(s) | Important modules/variations | Practical caution |
|---|---|---|---|
| Hotel | A04 Stay/unit-night; separate A02/A08 for venues | Unit type, occupancy, rate plan, check-in/out, taxes, restrictions | Full PMS/channel-manager/housekeeping scope is excluded initially. |
| Guesthouse | A04 | Small unit inventory, stay dates, occupancy, payment | May be simpler but still needs night-based inventory. |
| Resort | A04 plus A01/A02/A05 activities | Unit inventory, packages, activities, transfers | Multiple operations should remain under one organization. |
| Lodge | A04 | Unit inventory, meals/activities, transport | Remote connectivity and packaged stays may matter. |
| Vacation property | A04 | Property/unit, date range, occupancy, cleaning, deposit | Host/property management introduces separate roles. |
| Hostel | A04 with bed capacity | Bed/room type, occupancy, date range | Bed-level versus room-level inventory must be explicit. |

Hospitality remains a specialized future programme. Its inclusion in the taxonomy protects the organization model; it does not authorize building a hotel system in the first release.

### Sports and recreation

| Business | Primary archetype(s) | Important modules/variations | Practical caution |
|---|---|---|---|
| Football pitch | A02 Resource slot | Pitch, duration, team/contact, lighting, equipment | Back-to-back buffer and partial-pitch variants may exist. |
| Futsal court | A02 | Court, fixed slot, equipment | High repeat usage favours fast rebooking. |
| Basketball court | A02 | Court, slot, shared/exclusive policy | Public/shared sessions differ from exclusive rental. |
| Tennis court | A02 | Court, slot, optional coach/equipment | Court and coach are separate resources. |
| Gym | A05 Class/session; later membership | Class capacity, instructor, membership entitlement | General gym access is not necessarily a booking. |
| Personal trainer | A01, A05 | Trainer, location, recurring sessions, package | Recurrence and package usage matter. |
| Swimming pool | A02/A05 | Lane/private slot or public session capacity | Safety/age rules and session capacity differ. |
| Sports academy | A05 | Cohort/session, coach, age group, capacity, recurrence | Enrollment and attendance extend beyond one booking. |
| Gaming lounge | A02 | Station/room, duration, players, equipment | Multiple interchangeable stations may be pooled. |
| Go-kart location | A05/A02 | Session, kart capacity, participant/age rules | Safety waivers and arrival/check-in matter. |

### Automotive and transport

| Business | Primary archetype(s) | Important modules/variations | Practical caution |
|---|---|---|---|
| Car wash | A01, A09 | Bay/staff, vehicle type, service package, queue | Queue and appointment throughput may coexist. |
| Auto detailing | A01, A10 | Vehicle type/condition, package, duration range | Assessment may change duration/price. |
| Mechanic | A01 intake, A10 estimate/project | Vehicle, symptoms, inspection, estimate approval, parts | Appointment is intake; repair commitment follows diagnosis. |
| Tire service | A01, A09 | Vehicle, tire/service, bay, stock note | Parts inventory is separate future capability. |
| Car inspection | A01 | Bay/inspector, vehicle details, fixed procedure | Regulatory inspection data requires jurisdiction review. |
| Vehicle rental | A06 Rental | Vehicle class/unit, pickup/return, driver eligibility, deposit | License/insurance/condition/damage processes are specialized. |
| Chauffeur service | A07 Field/transport, A10 | Pickup, destination, time, passengers, vehicle/driver, distance | Route and travel time affect availability and price. |

### Professional services and education

| Business | Primary archetype(s) | Important modules/variations | Practical caution |
|---|---|---|---|
| Lawyer | A01 Consultation, A10 engagement | Matter type, conflict check, confidentiality, consultation | Legal-client confidentiality and engagement terms require special treatment. |
| Consultant | A01, A03/A10 | Topic, meeting, scope, deliverables, quotation | Consultation and project are different commitments. |
| Accountant | A01, A10 | Service type, period, documents, deadline | Sensitive financial documents require stronger controls. |
| Architect | A01, A03/A10 | Site, project stage, scope, deliverables | Project workflow exceeds an appointment. |
| Designer | A01, A03/A10 | Brief, deliverables, revisions, quotation | Versioned scope and milestones matter. |
| Tutor | A01, A05 | Subject, learner, recurring sessions, location | Guardian/minor handling may apply. |
| Language teacher | A01/A05 | Level, one-to-one/class, recurrence, capacity | Assessment may precede placement. |
| Driving school | A05/A01 plus resources | Instructor, vehicle, class/lesson, learner eligibility | Instructor and vehicle must be available together; regulation applies. |
| Training centre | A05 | Course/cohort, timetable, capacity, enrollment | Course enrollment and attendance extend beyond ordinary appointments. |

### Property and workspace

| Business | Primary archetype(s) | Important modules/variations | Practical caution |
|---|---|---|---|
| Real-estate property viewing | A08 Viewing | Property, agent, visitor, time, directions | Occupant privacy and exact-address disclosure controls may apply. |
| Apartment viewing | A08 | Unit/property, agent, visitor, time | Avoid publishing resident/customer data. |
| Office-space tour | A08 | Building/unit, agent, visitor/company | Group visitor and security check-in may apply. |
| Coworking room | A02 | Room/desk, duration, capacity, equipment | Desk memberships differ from room bookings. |
| Meeting room | A02 | Room, slot, capacity, equipment | Same resource archetype as event meeting rooms. |

### Healthcare — deferred regulated overlay

Clinics, dentists, optometrists, physiotherapists and similar providers primarily use A01 Appointment and sometimes A05 repeated sessions. They are not approved as an early vertical because they may require:

- Health/sensitive-data controls and restricted staff visibility.
- Clinical/professional licensing and record-separation decisions.
- Emergency disclaimers/escalation and patient-safety reliability.
- Consent, guardianship and longer statutory retention.
- Jurisdiction-specific telehealth/prescription/referral rules.

Ubi must not call ordinary booking notes a clinical record or accidentally invite collection of health data without a dedicated programme.

### Home and field services

| Business | Primary archetype(s) | Important modules/variations | Practical caution |
|---|---|---|---|
| Cleaning service | A07 Field service, A10 | Address/area, property size, frequency, team, supplies | Recurrence, travel zones and access instructions. |
| Plumber | A07, A10 | Address, issue, urgency, visit window, estimate | Emergency claims and arrival windows require careful wording. |
| Electrician | A07, A10 | Address, issue, safety screening, estimate | Safety/licensing and urgent hazard messaging. |
| Handyman | A07, A10 | Task list, photos, location, estimate | Several tasks can change duration and tools. |
| Painter | A08 site visit, A10 | Site, dimensions, surfaces, materials, quotation | Site assessment normally precedes final scope. |
| Appliance repair | A07, A10 | Appliance/model, symptoms, address, parts/estimate | Diagnosis and parts availability affect confirmation. |
| Moving company | A08 survey, A10 | Origin/destination, inventory, access, date, vehicle/team | Quote, route, team and vehicle form one commitment. |

## Additional candidate families not in the original list

These are taxonomy checks, not launch recommendations:

- Restaurant/table reservations — A02 capacity/resource, waiting list and party-size rules.
- Tours/guides/activities — A05 session/capacity or A03 private event, location and safety modules.
- Electronics/phone repair — A01 intake plus A10 diagnosis/quotation.
- Pet grooming — A01 with animal/profile/safety considerations; veterinary care is regulated.
- Tool/equipment rental — A06 quantity/unit rental with deposit and condition processes.
- Religious/community facilities — A02 venue or A03 event with governance/eligibility policy.
- Funeral/memorial services — A03/A10 with highly sensitive content and urgent support requirements.
- Agricultural equipment/service providers — A06 rental or A07 field service with location/logistics.

These candidates test whether the archetypes are general enough. They should not expand the first release.

## Operation composition examples

### Barber

```text
Archetype: appointment
Time: exact start + service duration
Resources: one barber, optional chair
Price: fixed service or simple add-ons
Confirmation: instant
Payment: pay at location
Other: walk-in queue and repeat booking
```

### Photography studio

```text
Operation 1: passport/portrait appointments
  exact slot + studio + photographer + fixed package + instant/request confirmation

Operation 2: event photography
  event date + package/custom requirements + team/location + quotation + deposit

Operation 3: studio rental
  studio resource + duration + equipment add-ons + deposit/approval
```

### Hotel organization

```text
Operation 1: rooms — unit-night stay inventory
Operation 2: wedding hall — venue date/slot + hold + package + deposit
Operation 3: conference rooms — room time slot + equipment/catering
Operation 4: airport transport — journey/field-service booking
```

These operations share one organization and may share customers, staff identity and reporting, but they do not share one booking form or availability algorithm.

## Onboarding taxonomy

The owner should not see the dimensions above as a technical questionnaire. Proposed onboarding:

1. Choose familiar business family.
2. Choose the specific business type or “another service.”
3. Ubi proposes one or more operations and common services.
4. Owner confirms how customers usually book: exact time, date/event, rental period, class or custom request.
5. Owner answers only configuration questions that change actual behaviour.
6. Ubi creates a test booking and explains the result.
7. Advanced modules remain hidden until enabled or needed.

If the owner chooses a contradictory setup—such as instant confirmation with an unavailable manual-only custom quote—Ubi should explain and resolve the contradiction rather than expose raw configuration flags.

## Taxonomy governance

- Business marketing category and booking archetype are versioned separately.
- New categories should reuse existing archetypes unless a genuine invariant cannot be represented.
- Adding fields is not sufficient evidence for a new archetype.
- Regulated overlays require legal/security review before activation.
- Template changes do not silently alter existing merchant operations.
- Merchant customizations remain within typed modules and policies; arbitrary executable workflow code is excluded initially.
- Analytics use stable taxonomy IDs, not translated display labels.

## First-release selection criteria

Score candidate verticals using evidence for:

- Frequency and cost of current booking problems.
- Merchant willingness to change behaviour/pay.
- Booking volume and customer repeat rate.
- Importance of deposits/no-show reduction.
- Workflow complexity relative to team capacity.
- Availability of reachable pilot merchants.
- Channel/payment prerequisites.
- Regulatory and safety burden.
- Ability to prove value in four to eight weeks.
- Transferability of the archetype to the next vertical.

Photography remains a strong complexity test and salons/barbers remain a strong frequency/retention test, but no first vertical is approved until field evidence and founder decision.

## Approval criteria

- Every originally listed business maps to at least one archetype and explains its exceptions.
- Pilot merchants recognize the proposed operations and terminology.
- Two business types that look similar but schedule differently are not forced together.
- Regulated/high-risk families are clearly gated.
- Configuration modules can be expressed without arbitrary per-merchant code.
- First vertical and its exact operation subset are selected through the decision register.

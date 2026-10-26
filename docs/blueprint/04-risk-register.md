# Ubi initial risk register

Status: **Draft**  
Scale: likelihood and impact use `Low`, `Medium`, `High`, or `Critical`. Scores are not yet quantitative because launch evidence is incomplete.

| ID | Risk | Likelihood | Impact | Early mitigation | Decision owner |
|---|---|---|---|---|---|
| RSK-001 | Building for too many business types before proving one workflow. | High | Critical | One-country/one-city/one-vertical gate; archetype architecture without broad first-release scope. | CEO/Product |
| RSK-002 | Product is simple for customers but too complex for owners/staff. | High | Critical | Observational research, assisted onboarding, task-based prototypes and comprehension metrics. | Product/Design |
| RSK-003 | Businesses continue using notebooks, leaving Ubi incomplete and untrusted. | High | Critical | Fast manual booking, walk-in/call capture, migration assistance and daily operational value. | Product/Operations |
| RSK-004 | External channels do not permit the inbox access Ubi promises. | High | Critical | Official API capability matrix, truthful product copy, optional adapters and link fallback. | CTO/Legal/Product |
| RSK-005 | Meta onboarding/app review makes WhatsApp/Instagram adoption difficult. | High | High | Prototype embedded onboarding, assisted support, prerequisite checks and Telegram/direct-link fallback. | Integrations/Operations |
| RSK-006 | Weak connectivity causes missing or duplicated bookings. | High | Critical | Local-first storage, idempotent commands, conflict handling, outbox monitoring and device/network tests. | CTO/Data |
| RSK-007 | Concurrent bookings create double allocation. | Medium | Critical | Database-enforced invariants, holds with expiry, isolation/locking tests and authoritative availability. | Data/Backend |
| RSK-008 | Payment screenshots or callbacks create false payment state. | High | Critical | Separate reported/processing/verified states, signed webhooks, reconciliation and manual exception queue. | Finance/Backend/Security |
| RSK-009 | Account recovery enables business takeover. | Medium | Critical | Verified ownership model, step-up checks, recovery delays/evidence, alerts and audit. | Security/Operations |
| RSK-010 | Staff or support access exposes unrelated tenant/customer data. | Medium | Critical | Server-side tenant authorization, scoped memberships, just-in-time support access and monitoring. | Security/Privacy |
| RSK-011 | A regional launch accidentally violates country-specific privacy/payment/consumer law. | Medium | Critical | Country selection before final design, jurisdiction matrix and licensed local counsel review. | CEO/Legal |
| RSK-012 | Data accumulates without purpose or deletion capability. | High | High | Data inventory, classification, retention schedule, deletion propagation and lineage tests. | Privacy/Data |
| RSK-013 | In-app messaging/calling creates harassment, fraud or safety incidents. | Medium | Critical | Booking-linked release, identity, block/report, moderation, rate limits, consent and later separate call approval. | Trust & Safety |
| RSK-014 | “Telegram-like” performance remains a slogan rather than an engineering property. | High | High | Representative device/network budgets, profiling, observability and release gates. | Mobile/SRE |
| RSK-015 | Premature microservices increase cost and unreliability. | Medium | High | Modular monolith initially; extraction only with measured scale/team evidence. | CTO |
| RSK-016 | Vendor lock-in or regional service gaps threaten operations. | Medium | High | Adapter boundaries, export/delete tests, multi-provider contingency for critical channels and vendor register. | CTO/Procurement |
| RSK-017 | SMS, messaging, storage or payment fees destroy small-merchant unit economics. | High | High | Per-booking cost model, configurable notification policy, pricing experiments and usage limits. | CEO/Finance |
| RSK-018 | Translation/calendar errors create wrong dates or inaccessible flows. | Medium | Critical | UTC/timezone model, country calendar tests, professional translation and field validation. | Product/Data |
| RSK-019 | Analytics definitions conflict with operational truth. | Medium | High | Metric dictionary, source lineage, versioned events and reconciliation against operational records. | Data/Product |
| RSK-020 | The platform becomes responsible for merchant service disputes without operational/legal preparation. | Medium | High | Clear platform role, terms, evidence/audit, escalation policy and local legal review. | CEO/Legal/Support |
| RSK-021 | Draft document coverage is mistaken for implemented/tested/legal release readiness. | High | Critical | Separate synthetic/research/live boundaries, missing-evidence register, truthful gate status and reproducible document-only checks. | CEO/CTO/Release owner |
| RSK-022 | Ethiopian/international clock ambiguity changes a real booking time despite correct calendar conversion. | Medium | Critical | Explicit clock convention, clarify ambiguous chat/voice time, repeat equivalent/date/timezone and test comprehension. | Product/Data |
| RSK-023 | Founder-only staffing cannot sustain promised support, privacy, financial or incident responsibilities. | High | High | Actual capacity/budget, named coverage, scoped external reviewers, no fictional independent approval or 24/7 promise. | CEO/Operations/Finance |

## Risk treatment workflow

1. Identify the affected actors, requirements, data and external dependencies.
2. Assign an accountable owner and a review date.
3. Select avoid, reduce, transfer or accept.
4. Define preventive, detective and recovery controls.
5. Define measurable control tests and residual risk.
6. Obtain the required approval for high/critical residual risks.
7. Link incidents and pilot evidence back to this register.

No critical risk may be accepted implicitly through schedule pressure.

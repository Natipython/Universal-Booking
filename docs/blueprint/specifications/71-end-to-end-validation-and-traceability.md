# End-to-end validation and requirement traceability

Status: **Draft — full scenario definitions not run; isolated V0 component evidence recorded separately**  
Updated: 2026-10-11 (Africa/Addis_Ababa)
Depends on: requirements ledger, specifications 10–70, policies 50–57, release contract 72  
Owner: Product + Engineering/QA; accountable review roles in specification 84

## Evidence rule

Coverage here means a requirement has a planned scenario. It does not mean an API, database migration, screen, legal document or automated test exists. Each release item also needs an executable requirement-level test and linked artifact before approval. Postponed capabilities keep their own future acceptance tests.

Use synthetic fixtures before approved field research/live data. Test at least two tenants, two operations in one tenant, owner/receptionist/assigned-worker/revoked-worker identities, a guest and a verified customer, D1/W1/W2 devices and N0–N3 networks. Distinguish simulated providers from real approved sandbox evidence.

## Founder-requirement crosswalk

| Validation group | Ledger requirements | Main scenario set |
|---|---|---|
| Local problem and onboarding | PRD-001, PRD-002, PRD-006, BUX-001, BUX-002, A11Y-001, LOC-001 | SCN-001, SCN-018 |
| Organization and staff scope | ORG-001, ORG-002, ORG-003, ORG-006, ORG-007 | SCN-002, SCN-012 |
| Manual work and durable sync | BKG-005, BKG-006, BUX-004, QLT-002, QLT-003, QLT-004 | SCN-003, SCN-004, SCN-013 |
| Capacity and truthful lifecycle | BKG-007, BKG-008, QLT-005 | SCN-005, SCN-006 |
| Standard guest flow | BKG-001, BKG-002, BKG-009, CUX-001, CUX-002, CUX-003 | SCN-007, SCN-008 |
| Editable contextual continuation | BKG-010, MSG-006, MSG-007, MSG-008, PRV-002 | SCN-009, SCN-010 |
| Custom requests and quotation | BKG-003 | SCN-011 |
| Money and evidence separation | PAY-001, PAY-002, PAY-003, PAY-004, PAY-005, PAY-006 | SCN-014, SCN-015 |
| Optional official channels and takeover | MSG-001, MSG-002, MSG-003, MSG-004, MSG-005, SEC-003 | SCN-016, SCN-017 |
| Required native messaging/calling and fallback | MSG-009, MSG-010, CUX-005 | SCN-019, SCN-020 |
| Android speed, localization and accessibility | PRD-004, QLT-001, QLT-006, A11Y-002, A11Y-003, A11Y-004, A11Y-005, LOC-002 | SCN-018, SCN-021 |
| Desktop work | PRD-005, BUX-005 | SCN-022 |
| Identity, isolation, authority and recovery | ORG-004, ORG-005, SEC-001, SEC-002, SEC-004, DAT-001, DAT-002, QLT-007 | SCN-002, SCN-012, SCN-023 |
| Privacy, legal, analytics and change | PRV-001, PRV-003, PRV-004, LEG-001, DAT-003, DAT-004, DAT-005 | SCN-024, SCN-025, SCN-026 |
| Controlled growth and customer return | PRD-003, PRD-007, PRD-008, BKG-004, BUX-003, CUX-004 | SCN-002, SCN-008, SCN-027, SCN-028 |

## Scenario catalogue

All full scenario rows have execution state **NOT RUN**. “Owner” names a role, not a staffed individual. A receipt must identify the build/schema/policy versions, fixtures, device/network, result, failures and reviewer. The [WB-01/02 synthetic host receipt](../validation/WB-01-02-synthetic-prototype-evidence.md) exercises local subsets of SCN-002/003/005/006/013/018. The separate [WB-03 component receipt](../validation/WB-03-synthetic-integrity-evidence.md) exercises selected transaction/HTTP/fixture-isolation subsets of SCN-002/004/005/006/011/012/014; it does not pass real identity, complete client/server, guest, device/network, provider or research scenarios. No full scenario is upgraded by these component results.

The [WB-04 customer-link receipt](../validation/WB-04-synthetic-customer-link-evidence.md) adds selected synthetic subsets of SCN-007/009/010/011: browser edits/review, ephemeral grant/crawler/expiry safety, retry and instant/enquiry distinction. SCN-008 manage actions, minimum real contact, persistent recovery, full quotation UI, app association, representative devices/networks and research remain absent. The 201 repeated case results are component checks, not 201 product scenarios; all 28 full rows remain NOT RUN.

| ID | Scenario and failure injection | Required outcome/evidence | Owner |
|---|---|---|---|
| SCN-001 | New owner chooses a photo template, edits prices/hours, records notebook bookings, creates test booking; interrupt setup halfway | Resume safely; demo/live unmistakable; owner understands publication and no social connection needed | Product/Operations |
| SCN-002 | One identity belongs to two tenants; one tenant has studio and event operations; switch scopes, rename slug | Shared organization ownership, separate configuration/resource rules, immutable IDs; no cross-tenant projection or authority | Backend/Data |
| SCN-003 | Receptionist enters call/walk-in booking on N0, then immediately kill/restart the app | Committed local command persists; displayed as local tentative, not confirmed; source/customer input retained | Android/QA |
| SCN-004 | Reconnect SCN-003 while duplicating/reordering requests and losing successful response | One logical booking effect; same ID/result recovered; audit/outbox complete; no repeat customer notice | Backend/Android |
| SCN-005 | Two devices/web users select the same exclusive slot; include active hold and setup/cleanup buffers | Exactly one authorized conflicting allocation succeeds; losers retain input/alternatives; same physical reservation boundary for hold/commit | Backend/Data |
| SCN-006 | Reschedule to occupied time, cancel during a stale edit, change service price after selection | Old schedule retained until valid atomic move; typed version/price conflict; snapshots/history intact | Backend/Product |
| SCN-007 | Guest uses public service link, skips email/account, chooses instant or approval mode; disable optional media/JS enhancement | Server-backed submit with minimum contact; readable review; outcome matches actual mode; no private calendar details | Web/Product |
| SCN-008 | Guest revisits/manage link and cancels/reschedules near cutoff; forward link to another person | Token limits/step-up policy enforced; consequences reviewed; display code alone grants nothing; no unrelated history exposed | Web/Security |
| SCN-009 | Chat-assisted birthday enquiry has wrong date; edit before submit; link opens web with no customer app | Provenance visible, edit survives, correct operation/service, fresh terms/capacity validation; no raw chat/phone in URL | Web/Integrations |
| SCN-010 | Replay expired/revoked/wrong-scope link; crawler previews; customer app absent/wrong account | Safe fallback/re-auth; GET does not commit or consume irreversibly; forwarded bearer risk documented; no unauthorized disclosure | Security/Mobile |
| SCN-011 | Wedding enquiry → clarify → quote v1 → v2 → accept stale v1 → accept v2 twice; expire hold | Exact accepted quote only; request/hold/booking distinct; one atomic conversion; amended terms never overwrite commitments | Backend/Product |
| SCN-012 | Staff queued booking then owner revokes membership; staff tries refund/export/provider/ownership action | Current server permission denies; local intent remains recoverable under policy; no high-risk offline commitment; audited denial | Security/Backend |
| SCN-013 | Log out/switch account with pending work; low disk, schema update, cursor expiry, device clock drift | No silent loss or cross-account upload; explicit local-access policy; safe resync/migration/expiry handling | Android/Data |
| SCN-014 | Cash/no deposit, merchant cash receipt, customer screenshot and provider verification paths | Evidence classes visibly distinct; exact ETB arithmetic; booking/payment separate; customer claim never creates provider-verified status | Backend/Finance |
| SCN-015 | Provider simulator sends duplicate/out-of-order/forged callbacks, late payment after hold expiry, refund timeout/reversal | Authenticated idempotent processing; exceptions queued; no false booking/refund success; merchant/provider funds path retained | Payments/Finance |
| SCN-016 | Business declines social integration, then eligible sandbox connects/revokes an official channel | Link/manual work unaffected; minimum disclosed scope, encrypted credentials, capability-bound history access, disconnect preserves records | Integrations/Security |
| SCN-017 | Bot collects draft; customer asks for human; two staff claim; disconnect provider mid-reply | Single ownership/takeover state, automation stops, source/delivery states truthful, reply cannot bypass domain truth | Integrations/Product |
| SCN-018 | Limited-digital-experience participants use large fonts/TalkBack, Amharic, calendar/year boundaries and ambiguous local clock | Correct task and explanation, not memorized taps; no date/time shift; translated status distinguished without color | Design/QA |
| SCN-019 | Booking-linked native message: retry, unauthorized participant, blocked sender, revoked staff, failed attachment | One accepted message, correct participant scope, safety controls, truthful delivery state; no quote/payment authority from free text | Messaging/Safety |
| SCN-020 | Future in-app voice call: consent/ring/answer/decline/end, network switch/backgrounding, block, relay unavailable | Accurate lifecycle, bounded reconnect/failure, no emergency/reachability guarantee, recording off, metadata/media cost and privacy evidence | Calling/Safety |
| SCN-021 | D1 cold/warm startup, scrolling, N1 submit, process death, low storage/thermal load | Raw traces against specification 21/62; no fake success or missing accessibility to meet timing; byte/battery budgets measured | Android/QA |
| SCN-022 | Receptionist works desktop queue/calendar with keyboard/zoom; two tabs and lost response | Scope/focus retained, list alternative, idempotent writes, atomic reschedule, safe shared-PC sign-out | Web/QA |
| SCN-023 | Restore isolated database/object fixtures; rebuild workers, revoke keys and reconcile simulated provider outcomes | Constraints/RLS, reservation integrity, idempotency/outbox/timers/deletion state and application usability verified; measured RPO/RTO | Operations/Data |
| SCN-024 | Support case needs content elevation; expired/wrong-tenant grant; exercise privacy request/legal hold | Redacted-first, scoped/JIT access, audit, verified requester, honest partial deletion and lawful exception process | Privacy/Security |
| SCN-025 | Metric job counts requested/confirmed/cancelled and reported/verified money; delete/restrict source data | Versioned denominators, no inflated revenue, purpose/minimization, propagation/lineage verified | Data/Product |
| SCN-026 | Roll forward/back-compatible schema/client change, rollback app with older offline commands | No corrupt pending work, supported-version behavior, rehearsed expand/contract plan and recovery receipt | Engineering/Data |
| SCN-027 | New service configuration, block hours, archive resource, variant unsupported archetype | Existing commitments unchanged; conflict preview; unsupported hotel/health workflow cannot silently activate | Product/Backend |
| SCN-028 | Rehearse launch packet, merchant closure, incident pause and expansion decision with missing approver/vendor/counsel | Missing mandatory evidence blocks relevant live capability; no claim of legal approval, provider approval or completed release | CEO/Operations |

## Test-design detail: capacity concurrency

Use the same tenant/resource with overlapping intervals, exact adjacent boundaries, buffers, held and committed ownership, hold expiry/conversion, cancelled allocation, multi-resource booking and a second tenant using different resources.

For each test, start multiple contenders behind a barrier. Record command/result IDs and query authoritative reservation/booking/outbox/audit state afterward. Success is one allowable effect, not one HTTP success or one green screen. Retry after timeout must retrieve that effect. Multi-resource partial allocation must roll back as one transaction.

The production SQL design remains proposed. WB-03 executes a narrower disposable exclusive-resource schema; pooled/nightly/multi-resource rules and production migration/restore remain untested. [PostgreSQL's exclusion-constraint documentation](https://www.postgresql.org/docs/18/ddl-constraints.html), rechecked 2026-10-04, supports exclusion rules but does not prove the eventual Ubi schema correct.

## Moderated comprehension protocol

Recruit target owners, receptionists, customers and service staff separately. Present tasks in familiar words, avoid coaching during the scored attempt, record intervention and critical misunderstanding, then debrief. Test errors/recovery as well as happy paths. Offer accessible alternatives and obtain research permission; do not record real customer content by default.

Pilot thresholds in specification 83 are provisional. Report per-task/per-language/per-role sample counts and raw outcomes. A small round finds problems; it cannot establish a statistically reliable regional success rate.

## Release mapping

The candidate first live boundary and gated extensions are in specification 72. A deferred scenario remains required before its later capability ships. Cases 15–17 and 20 need approved provider/calling evidence if activated; a simulator pass never substitutes for provider eligibility or legal review.

Create requirement-level automated/manual case IDs in the implementation backlog. Each failure has severity, owner, next safe action, retest and evidence. Mandatory integrity/privacy controls cannot be waived by a good usability average.

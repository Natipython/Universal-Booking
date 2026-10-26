# Pilot scorecard and experiment governance

Status: **Draft — definitions/provisional thresholds; no measured values**  
Updated: 2026-10-04 (Africa/Addis_Ababa)  
Depends on: metric dictionary 47, data governance 48, quality 21/62/63, pilot 80

## Measurement discipline

Every metric records numerator, denominator, eligible population, time window, data source, timezone, segment, missing data, owner and privacy purpose. Report counts with percentages and medians/percentiles with sample sizes. Downloads, page visits, message delivery and polite interest are not successful bookings, merchant retention or revenue.

Definitions below are pilot interpretation of the governed dictionary, not separate lifecycle meanings. Changes are versioned and preserve comparability or explain breaks.

## Scorecard

| Metric | Definition/evidence | Proposed decision use |
|---|---|---|
| Manual booking task time | Start to correct durable entry, including corrections; local save and server confirmation measured separately | Candidate median ≤30 s for a trained repeated standard task; first-time results separately; validate against merchant baseline |
| Critical unassisted task completion | Correct task and correct status/date understanding, without staff coaching, divided by attempted eligible tasks | Candidate ≥90% for each primary task/subgroup, not blended average; insufficient sample is unverified |
| Wrong confirmation/payment belief | Participant thinks request/local save/payment report is confirmed/verified | Zero observed critical misunderstanding before releasing that flow; investigate every occurrence |
| Booking capture | Ubi recorded in-scope bookings divided by reconciled actual merchant bookings, including calls/walk-ins/messages | Candidate ≥90% after assisted transition; unknown denominator means not measurable |
| Weekly active merchant | Enrolled merchant performs a defined meaningful operational action in local week | Observe by cohort; logging in alone does not qualify |
| Four-week retained usage | Enrolled eligible cohort still captures agreed workflow in week four | Candidate ≥8 of 10 if ten are enrolled and have eligible work; publish reasons and small-sample uncertainty |
| Paid continuation | Merchants actually pay and continue using approved scope after disclosed trial | Evidence for pricing; no fixed universal conversion claim before experiment |
| Request response time | Valid customer request to first substantive authorized response, segmented by declared business hours | Compare baseline; an automated acknowledgment is not substantive resolution |
| No-show/cancellation | Defined fulfilled/no-show/cancelled outcomes per eligible confirmed bookings | Baseline comparison; no causal claim without controlling service mix/seasonality |
| Scheduling/data integrity | Proven prohibited double allocation, lost durable intent, cross-tenant disclosure, duplicate logical effect | Zero known unresolved release-critical defect; incidents trigger pause/runbook |
| Android quality | D1 startup/frame/byte/battery traces and crash/ANR rates, segmented by build/network/device | Specification 21/62 targets; short pilot cannot prove long-window reliability |
| Customer web quality | Text task under N1 plus real/lab LCP/INP/CLS and transferred bytes | Separate field targets from constrained-network lab evidence |
| Sync health | Age/outcome of pending commands, rejection/conflict rates, reconciliation completion | Watch p95/p99 and worst queue; averages cannot hide stranded work |
| Notice delivery/action-needed | Per-provider attempts and final evidence, stale reminder suppression, unresolved critical notices | Preserve current status access even if notices fail |
| Support effort | Staff minutes/cases per active merchant and correctly completed task | Budget against actual staffing; learning improvements should reduce repeated burden |
| Contribution/cash | Collected software revenue less attributable costs; actual cash timing/runway | Specification 81; no merchant GMV counted as Ubi revenue |
| Safety/privacy | Reports, restrictions/appeals, access elevations, privacy request outcomes and overdue cases | Quality of response and evidence, not “zero complaints means safe” |

These thresholds are hypotheses for founder/product review. The earlier 70% standard-booking criterion is not the current quality gate; an average 70% could leave many customers unable to finish safely. Do not label a pilot passed by selecting the easiest subgroup or omitting interrupted attempts.

## Usability sample and inference limits

Small moderated rounds find critical defects and guide redesign. A 9/10 result is an observed 90%, not proof that 90% of Ethiopia's population will succeed. Plan larger validation samples or state uncertainty before formal population/reliability claims. Supplement SUS, if used, with direct tasks and comprehension; candidate median SUS ≥80 is not proof of correctness or accessibility.

Use limited-digital-experience owners/staff/customers, approved languages, low-cost devices, shared-device scenarios and accessibility needs. Track interventions without shaming participants. Revise the product when people misunderstand it.

## Experiment controls

Before launch, record hypothesis, primary metric, eligibility, duration/sample rationale, baseline, segmentation, privacy purpose, exposure/rollback, harm guardrails, budget, owner and stop condition.

No experiment can weaken authorization, hide fees, confuse confirmation, override capacity, fabricate urgency or enroll marketing without appropriate basis. Pricing experiments must disclose the actual offered terms and honor commitments. Do not optimize conversion by removing necessary review or human help.

## Review cadence and decisions

- Supervised live transition: daily integrity/sync/support and booking-reconciliation review.
- Weekly: adoption, task problems, capture, support costs and open incidents.
- Pilot exit: reviewed cohort results, economics and safety; explicit continue/rework/pivot/pause.
- Expansion: proof first workflow is retained and supportable; separate country/vertical capability packet.

Any critical privacy/integrity/false-success event pauses the affected feature; invoke 66–70 and communicate actual impact. Good growth does not compensate for a known critical defect.

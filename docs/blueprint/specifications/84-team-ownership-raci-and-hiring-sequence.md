# Ownership, practical staffing and review responsibilities

Status: **Draft — roles defined; actual team, founder composition and budget unknown**  
Updated: 2026-10-04 (Africa/Addis_Ababa)  
Depends on: decisions 5–7, release 72, support 65, finance 81, backlog 85

## Realistic company operating model

CEO, CTO, CFO, lawyer and data engineer are accountability perspectives, not evidence that Ubi employs five executives. If there is one founder, that person may initially own product/business and engineering coordination; obtain qualified external review where skill, law or separation of duties requires it. Do not invent cofounders, shareholders, capital or registered status.

The private-limited-company structure in DEC-005 is a proposal for qualified Ethiopian review using actual ownership/operating facts. Company name, objects, registration, tax and employment arrangements remain unverified. A planning assistant is not licensed Ethiopian counsel or a statutory officer.

## RACI-style responsibility map

R = performs/co-ordinates the work; A = accountable for the decision; C = consulted for necessary expertise. Each accountable role must eventually be assigned to a named person. No person currently assigned by this document.

| Work | R | A | C / independent review |
|---|---|---|---|
| Market/pilot/release scope | Product/merchant operations | CEO/founder | Merchants, staff/customers, finance, CTO |
| Architecture/framework/runtime | Engineering | CTO/technical lead | Android/web/data/security reviewers |
| Booking/database integrity | Backend/data engineer | CTO | Independent concurrency/tenant-isolation review |
| UX/accessibility/local wording | Design/content/local researcher | Product lead | Low-confidence users, human translators, accessibility reviewer |
| Security/privacy architecture | Security/privacy owner | Accountable company operator | Qualified independent specialists and counsel as applicable |
| Entity/terms/regulatory interpretation | Licensed counsel with actual founder facts | Company operator for business commitments | Accountant, provider, product/security |
| Cash plan/pricing/subscriptions | Finance/bookkeeping | CEO/finance owner | Qualified Ethiopian accountant, CTO cost inputs |
| Payment role/provider/reconciliation | Payments/finance operations | Company operator/finance owner | Licensed provider and qualified counsel; technical security review |
| Merchant onboarding/support | Operations/support | Operations owner | Product, engineering, merchant-authorized representative |
| Native messages/calls/safety | Communications engineering + safety | CTO/product accountable owner | Privacy, abuse response, finance/provider reviewers |
| Production release/incident/restore | Release/on-call operations | Release owner/CTO | Security/data, finance for financial incident, independent checks |

## Staffing sequence driven by constraints

1. **Discovery/planning:** founder/product coordination, strong technical ownership, locally capable researcher/translator and scoped legal/accounting advice. Not all full-time roles.
2. **Synthetic feasibility:** experienced Android/Flutter and Go/SQL capability; web support for the minimal customer link and contract tests; competent independent correctness/security review.
3. **Controlled pilot:** named merchant-onboarding/support coverage, release/incident owner and finance/privacy exception handling; one person can wear several hats only within a funded workload.
4. **Connected communication/payment:** add actual integration/reconciliation/safety ownership before activation, not after a customer loses trust.
5. **Calling/expansion:** add media/network expertise and sustained abuse/operating coverage based on demonstrated volume; avoid promising a globally available communications service on founder-only staffing.

Hire/contract for missing capabilities, not impressive titles. Rates, headcount and schedule require available budget and actual skills. If only one person can build the pilot, reduce parallel features and commercial promises.

## Separation of duties without fictitious staff

Ownership transfer, high-risk export, production credential access, refunds/financial reconciliation and security exceptions need documented review requirements. If a second approver is required and none exists, use a qualified authorized external reviewer or keep the action unavailable; do not have the founder click two roles and call it independent approval.

Incident emergency access is controlled break-glass, not permanent administrator freedom. Retrospective review remains mandatory under the relevant runbook. Support/onboarding cannot ask merchants for passwords or OTPs.

## Capacity and coverage

Record each person's available weekly hours, concurrent responsibilities, absence cover and escalation contact before committing delivery dates. Scope buffers for research, rework, provider review, translation, testing and incidents.

Customer-visible support hours match staffed coverage. There is no invented 24/7 help desk. Production alert routing must reach a real person and safely pause affected optional capabilities when response coverage is unavailable.

## Assignment gate

Before a backlog item becomes scheduled, it needs a named accountable owner and accepted capacity. Before live release, all critical operating/privacy/financial responsibilities and specialist review gaps must be resolved or the affected capability remains disabled. Role tables alone cannot pass EV-09/EV-10.

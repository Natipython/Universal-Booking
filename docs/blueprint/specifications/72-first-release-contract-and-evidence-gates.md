# First-release candidate contract and evidence gates

Status: **Proposed candidate — NOT frozen, NOT launch-approved**  
Updated: 2026-10-15 (Africa/Addis_Ababa)
Owners: accountable roles in specification 84  
Inputs: decision register, specifications 10–71, policies 50–57, business execution 80–85

## What is decided and what is not

Decided: Ethiopia is the first launch country; Ubi Business Android is the first business application; desktop follows; native in-app messages and calls are required in the platform direction; the founder's visual palette is retained.

Planning assumptions only: Addis Ababa, photo studios/photographers, ten businesses over 8–12 weeks, merchant SaaS, a suitable Ethiopian private-limited company and non-custodial merchant/provider payments. Founder composition, capital, ownership and actual legal/tax facts are unknown. No domain is purchased, company registered, provider accepted, counsel approval obtained or technology selection validated by this document.

Keep these assumptions editable. Planning may continue with them; irreversible/live commitments may not treat silence as approval.

## Separate four boundaries

| Boundary | Allowed purpose | Required before progressing |
|---|---|---|
| V0 synthetic feasibility prototypes | Compare frameworks, test invariants/UI with fake identities/data/providers | Recorded scope, isolated environments, no external customer messages or money |
| V1 moderated research | Observe consenting participants with dummy bookings | Research notice/permission, minimal data, test devices, no misleading live availability |
| R1 controlled live pilot candidate | Real approved merchants/bookings in the narrow release | Founder release decisions, actual accountable operator, applicable legal/privacy/security/provider reviews, exercised operations and acceptance evidence |
| R2 gated extensions/expansion | Native messaging, official connectors, digital payments, calls and additional markets as separately scoped | Capability-specific tests, policy/cost/vendor/operational proof and approved amended contract |

This avoids a circular demand for a completed production app before validating its framework. V0 prototypes are disposable validation work, not authorization to expose real customers or bypass launch gates. The founder's request for counsel review after planning/build preparation is accommodated by a complete scoped review packet before real-data/public/regulated release; no legal advice or readiness is implied here.

## R1 capability and actor matrix

| Capability | Owner/authorized staff | Customer | Ubi operations | Boundary |
|---|---|---|---|---|
| Identity, org/membership, roles/audit | Own session; create/join org under authority; scoped tasks | Guest contact/link grant; optional verified identity later | Separate workforce identity and audited support | Required |
| Photo template/services/hours/resources | Configure within capability; snapshot existing terms | View published service/price/location | Assist setup with authorization | Required |
| Manual call/walk-in/message booking | Fast entry; authoritative online result or offline local tentative | Receive notice only through approved permitted endpoint | Diagnose command state | Required |
| Android Today/calendar/offline sync | Cached authorized schedule, durable commands, conflicts | Not a staff client | Redacted diagnostics | Required |
| Standard link booking/manage link | Instant or request-to-confirm eligibility/policy | Account-free completion/review/return | Support token expiry without seeing raw secret | Required supporting customer surface |
| Custom enquiry and minimal versioned quote | Triage/clarify, issue quote, convert atomically | Submit enquiry, review exact quote, accept conditionally | Diagnose lifecycle without editing history | Candidate narrow subset; validate pilot need |
| Desktop business workspace | Scoped calendar/queue, keyboard operation | Not a business workspace | Separate console | Next business surface; prove shared contracts before rollout |
| Cash/no deposit/pay-at-location | Merchant receipt/evidence with provenance | Clear obligation, no fake provider verification | Reconcile evidence classes | Candidate default |
| Essential status/reminders | Action-needed queue and failed-notice handling | Approved channel, truthful current-state notice | Delivery diagnostics/incident route | Minimum release-dependent notification method required |
| Booking-linked native messaging | Authorized participant/queue, takeover, safety | Scoped thread/session | Audited safety/support | Required platform milestone; separately gated activation |
| Official social connectors | Optional connection through authorized business account | Use existing supported channel | Capability/connection diagnostics | Optional; no live connector dependency for basic bookings |
| Provider-verified deposits/payments | Merchant/provider settlement path and reconciliation | Approved checkout, exact money, unknown recovery | Exception handling; no custody | Separate real-money gate |
| Native in-app voice calling | Authorized call relationship, identity/privacy | Consent, ring/answer/decline/end, block/report | Metadata-limited support | Required platform milestone after safe messaging; separate activation gate |

This is a proposal for delivery order, not founder approval to drop or indefinitely forget messaging/calls. Backlog 85 names both workstreams and their acceptance work. If the first contracted live release must include them, revise staffing, scope and timing explicitly rather than showing fake functionality.

## Exclusions from this candidate live boundary

Open consumer discovery/ranking; general person-to-person messaging; hotel property-management/inventory/channel-manager functionality; healthcare records; multiple-country launch; arbitrary executable merchant workflows; autonomous AI confirmation/refund; wallets/pooled/split customer funds; video and recording. A future architecture entry is not a shipped feature.

## Domain/data/API contract to freeze

Before R1 is build/release-authoritative, identify actual versions and receipts for:

1. Account/contact, organization/branch/operation/membership and public slug ownership.
2. Service/resource/hours/price/policy versions and commitment snapshots.
3. Booking/request/quote/hold/allocation/payment/notification/conversation lifecycles, with permitted actor/action matrix.
4. Physical reservation constraints covering holds and committed exclusive intervals together; composite tenant keys, RLS/authorization and transactional command/outbox rules.
5. Schema migrations, typed command/result/error/event contracts, supported client/sync versions and conflict cases.
6. Link scope, expiry/revocation, participant/session authorization, required step-up and crawler-safe redemption.
7. Local dataset/quota/offline lease, logout/revocation/resync and queued-command compatibility rules.
8. Field/event purposes, classifications, retention, deletion/legal hold, data location and vendor roles.

The documents define candidates. No executable schema, approved OpenAPI contract or database migration exists yet. A freeze references immutable versions/checksums and an explicit change procedure; “use the latest draft” is not a freeze.

## Quality and integrity acceptance

Use the full requirements in 21/22/61/62; do not replace them with this summary.

- Zero accepted Ubi-induced prohibited double allocations, cross-tenant disclosure or lost durable user intent in the release test suite.
- Idempotent retries, conflict preservation, truthful booking/payment/message statuses and atomic reschedule/quote conversion.
- D1 cached cold TTID ≤2 s p75, usable ≤3 s p75, warm usable ≤1 s p75, local acknowledgment ≤100 ms; device/network traces required.
- W1 text flow usable on N1; field p75 web targets and byte budgets measured separately. A 250 KiB byte ceiling alone cannot prove 2.5 s LCP on a 400 kbit/s network.
- Critical low-confidence user tasks meet agreed comprehension/accessibility criteria, with subgroup results and reviewed English/Amharic/date/time wording.
- Recovery, release rollback, incident/support/privacy workflows exercised; actual measured RPO/RTO, costs and named coverage recorded.

Proposed 99.8% crash-free sessions, 99.9% general-launch booking-path availability, RPO ≤5 min and RTO ≤4 h remain targets. A tiny pilot or short prototype cannot establish monthly reliability or a contractual SLA.

## Evidence register

2026-10-15 additions: [one physical Android 12 sample](../validation/WB-02-physical-android-smoke-evidence.md) adds bounded native SQLite and explicitly authorized force-stop/draft-recovery evidence to partial EV-03/05. The later [Home action follow-up](../validation/WB-02-home-action-layout-evidence.md) fixes the observed overlap in source/host tests, with 41 passing tests and a separate candidate awaiting authorized phone recheck. The original phone-tested APK/inputs remain archived. Neither addition verifies Android 11/4 GB, release/profile performance, full native acceptance, accessibility/localization, Android/server integration or any release gate. Earlier dated evidence below remains historical.

Every entry is **MISSING** at creation. Owner assignment is also pending until a real person accepts it.

2026-10-11 update: EV-03/05 have **PARTIAL — synthetic host tests and debug compilation/static APK inspection** evidence in the [WB-01/02 host receipt](../validation/WB-01-02-synthetic-prototype-evidence.md) and [APK receipt](../validation/WB-02-synthetic-android-apk-evidence.md). The [WB-03 receipt](../validation/WB-03-synthetic-integrity-evidence.md) adds **PARTIAL EV-03/04** evidence for disposable Go/PostgreSQL invariants and loopback HTTP contract assertions, with a recorded spec-linter disagreement. The [WB-04 receipt](../validation/WB-04-synthetic-customer-link-evidence.md) adds **PARTIAL EV-03/04/05** form/grant/retry and headless-browser evidence, not a complete guest product. Android adapter execution/device measurements, weak-network, localization/accessibility/research acceptance, real identity, durable web continuation, production migration/restore/load and independent review remain missing. Other entries remain missing; no named reviewer approval or gate passage is inferred.

| Evidence ID | Receipt needed | Gate and accountable role |
|---|---|---|
| EV-01 | Founder city/vertical/pilot/scope decisions and ownership facts | Market/release; CEO |
| EV-02 | Consenting merchant/staff/customer observations, baseline and comprehension results | Market/product; Product |
| EV-03 | Framework/runtime/native-boundary ADR evidence: D1 prototype, Go/SQL invariants, low-byte web | Technical feasibility; CTO |
| EV-04 | Executable schema/contracts, tenant-negative/concurrency/idempotency/migration results | Data/correctness; Data/CTO |
| EV-05 | Release-device/offline/process-death/security/UX/translation/date tests | Client quality; QA/Product |
| EV-06 | Current entity/activity/privacy/retention/terms/communications review for actual release | Trust/legal; accountable operator + licensed Ethiopian counsel |
| EV-07 | Vendor scope/data-location/exit/cost review and applicable official sandbox/production access | Provider/data; CTO/Privacy |
| EV-08 | Backup/restore, rollout/rollback, incident, support and privacy exercise receipts | Operations; Operations/CTO |
| EV-09 | Funded cost envelope, pricing hypothesis, support capacity and named responsibility | Commercial; CEO/Finance |
| EV-10 | Pilot-specific release artifact/signing/rollback, known issues, gate decisions and operator sign-off | Controlled live pilot; release owner |
| EV-11 | Native messaging identity/safety/delivery/offline/privacy/support test results | Messaging activation; CTO/Safety |
| EV-12 | Official channel eligibility/scopes/review/onboarding/webhook/takeover evidence, per provider | Connector activation; Integrations/Privacy |
| EV-13 | Approved non-custodial arrangement, provider truth/refund/settlement/reconciliation/counsel proof | Real-money activation; Finance/CTO/counsel |
| EV-14 | Call feasibility, background/weak-network/abuse/consent/media-cost/vendor review | Calling activation; CTO/Safety/Finance |

Each receipt records date, source, artifact/build, test boundary, result, limitations, reviewer and expiry/recheck trigger. No numeric target, checklist tick or test simulator is its own evidence.

## Gate status and freeze procedure

G7 **OPEN** (legal/privacy/entity/payment/provider conclusions); G8 **OPEN** (measured reliability and exercised operations); G9 **OPEN** (validation and release approval). Market G1 and experience G4 also have missing field evidence; existence of documents does not close them.

Review sequence: complete V0/V1 evidence → resolve/explicitly disposition decisions and contradictions → version candidate contracts → named role review → licensed/qualified review of relevant live scope → authorize bounded implementation/release activities → rehearse/verify live prerequisites → controlled pilot. Production gates apply to the relevant activated capability; optional connectors/payments/calls can remain disabled without blocking a verified basic booking pilot.

A freeze records each requirement as included, excluded for this release with a future milestone, or blocking; exact evidence IDs; accepted risks with accountable approver and expiry; release hashes; stop/rollback conditions. Safety-critical missing evidence blocks release, not merely a line in “known issues.”

## Immediate next work

The synthetic Android manual-booking/local-state prototype, field task protocol, separate Go/PostgreSQL experiment and WB-04 customer-link component checks now exist. The approved isolated Android setup and inspected debug APK are recorded separately. A [bounded Android 12 physical-phone session](../validation/WB-02-physical-android-smoke-evidence.md) adds actual native storage/recovery evidence for the old 0.1.0 APK. The [layout fix](../validation/WB-02-home-action-layout-evidence.md) is verified locally; next obtain renewed authorization for candidate installation/native recheck, then remaining appropriately authorized native/device/network/localization/comprehension checks and reviewed web compatibility/recovery design. Minimum-device and profile/release performance evidence remain missing; debug size/timings do not establish production targets. WB-04's in-memory sessions and dummy customer labels are not approved production shortcuts. The WB-03 linter disagreement and relevant technical/research reviews must be resolved/dispositioned before WB-05 contract freeze. Resolve actual founder choices without repeating decisions already made, including the decided Android 11+/4 GB floor. Prepare real-data/counsel/vendor review before inviting live merchants. Build sequence and dependencies are specification 85.

# Dependency-ordered release roadmap and implementation backlog

Status: **Draft — build preparation; estimates, named owners and release freeze pending**  
Updated: 2026-10-15 (Africa/Addis_Ababa)
Depends on: validation 71, candidate contract/evidence 72, decision closure 73, commercial/operating work 80–84

## Recovered next step

Continue Stage 9; do not repeat completed Stage 8 drafting. The first concrete validation slice is Android manual booking/offline state plus a Go/PostgreSQL integrity experiment using synthetic records. The minimal customer link page supports that shared contract. Desktop remains the next business product surface after Android.

Stages below are delivery increments, not replacements for programme stages 0–12. Scope may be reduced or changed by evidence. A new language/framework does not get its own product rewrite without a measured benefit.

## Execution checkpoint — 2026-10-15

| Package | Observed state | Still required |
|---|---|---|
| WB-01 | Synthetic two-tenant/three-operation fixtures and moderated task protocol prepared | Real research permission/recruitment, founder city/vertical/pilot decisions and field outcomes |
| WB-02 | 41 host tests/analyzer pass; Home action overlap fixed locally with 13 regression checks; separate 0.1.1 (2) candidate. Old 0.1.0 authorized Android 12/SM-G975W save/review/conflict/revision/cancellation/draft sample passed, including one authorized force-stop; 11 historical offline cross-checks pass | Renewed authorization to update/recheck layout on a phone; remaining lifecycle/scope/retry checks, Android 11/4 GB and newer devices, release/profile performance, network/TalkBack/localization/comprehension evidence |
| WB-03 | Isolated Go/PostgreSQL/HTTP experiment implemented; transaction/race/contract checks recorded in a disposable cluster | Independent review, actual identity/migrations/restore/load scope and recorded spec-linter disagreement; not production or Android integration |
| WB-04 | Local synthetic handoff/edit/review/outcome implemented; 201/201 repeated component case results including Postman and no-JavaScript Chromium paths | Weak-network/device/compatibility/localization/comprehension evidence, persistent continuation/identity, full enquiry/quote/manage flow, app association and review; component pass is not package closure |
| WB-05 | Not ready | Relevant evidence, review and versioned bounded build-contract decision |

The historical [WB-01/02 host receipt](../validation/WB-01-02-synthetic-prototype-evidence.md), [debug APK receipt](../validation/WB-02-synthetic-android-apk-evidence.md) and latest [layout follow-up](../validation/WB-02-home-action-layout-evidence.md) provide partial EV-03/05 evidence; the [WB-03 receipt](../validation/WB-03-synthetic-integrity-evidence.md) adds partial EV-03/04 component evidence; the [WB-04 receipt](../validation/WB-04-synthetic-customer-link-evidence.md) adds partial EV-03/04/05 form/browser evidence. None is package/gate completion. Next dependencies are authorized Android native/device work, representative web/network/localization/comprehension checks and relevant reviews, not automatic WB-05 closure. DEC-022 sets the business-app support floor to Android 11+/4 GB; RAM/OS do not themselves prove speed. Desktop and required native messages/calls remain in the ordered work below.

The [Android build setup](../validation/WB-02-android-build-preflight.md) records the explicitly authorized SDK agreements and verified local tooling. The corrected debug APK passes static identity/API/backup/signature/adapter/asset checks; failed/superseded attempts are preserved. Its 175.6 MB debug size is not a production distribution target. The subsequent [2026-10-15 physical-phone sample](../validation/WB-02-physical-android-smoke-evidence.md) records scoped installation/test permission and a separately authorized single force-stop on Android 12. The [native checklist](../validation/WB-02-android-native-test-checklist.md) is now partially executed, not complete. The testing connection was stopped; future device access/settings tests need appropriate direction. No general release approval or minimum-device performance result is inferred.

## Work packages and dependency graph

| Item | Scope / responsible role | Depends on | Acceptance and evidence | Size/uncertainty |
|---|---|---|---|---|
| WB-01 | Research fixtures, release assumptions, draft screens/copy and task protocol; Product/Operations | Current 71–84 drafts | Safe synthetic fixtures for two tenants/operations; bounded task script; founder decisions recorded rather than guessed; EV-01/02 scope | S–M; field recruitment uncertain |
| WB-02 | Android Flutter/Dart feasibility prototype; Android | WB-01 | Manual add → local Today/details → process-death recovery → pending/conflict labels; D1 release/profile timings and accessibility evidence; EV-03/05 | M; physical devices/toolchain required |
| WB-03 | Go/PostgreSQL transactional experiment; Backend/Data | WB-01 | Tenant-negative tests, overlapping held/committed reservation race, atomic reschedule/quote conversion, idempotency and outbox receipt; EV-03/04 | M; DB/contract selections not frozen |
| WB-04 | Thin customer link and contextual resume prototype; Web/Security | WB-01; aligns with WB-03 contracts | Standard/custom review, editable prefill, private token/crawler/expiry tests, web fallback, small-byte N1 measurements; SCN-007–011 | M; domain/association prototype can be local |
| WB-05 | Architecture/experience closure and versioned build contract; CTO/Product/Data | WB-02/03/04 and reviewed research | Close/disposition stack/API/localDB/isolation decisions with evidence; approved bounded build scope; no live approval inferred; EV-03–05 | S–M; reviewers unavailable = blocker |
| WB-06 | Repository, dependency/toolchain pins, CI, environments, security/telemetry foundations; Engineering | WB-05 | Verified runtime versions, reproducible tests/builds, secret scan, SBOM/signing plan, failure-safe observability, isolated synthetic/staging environments | M; provider selection/review |
| WB-07 | Identity, org/operation/membership/permission/audit vertical slice; Backend/Android | WB-06 | Actual migrations/contracts/UI; phone-first recovery and tenant-negative/revoked-session cases; SCN-002/012/013; no fake live OTP | L; identity provider/local-access lease |
| WB-08 | Catalog/hours/resource/manual-booking/offline sync vertical slice; Backend/Data/Android | WB-07 | End-to-end durable command, committed reservation, conflict recovery, cached Today, versioned price and history; SCN-003–006/013/021/026 | L; core integrity and schema risk |
| WB-09 | Customer booking/manage link and minimum custom enquiry/quote slice; Backend/Web/Product | WB-08 | Exact review/commit/return permissions, guest privacy, immutable quote, no false success; SCN-007–011/014/018 | L; custom subset needs market evidence |
| WB-10 | Desktop business workspace; Web/Product | Shared WB-08/09 APIs and Android contract proven | Queue/calendar/keyboard and concurrency acceptance; shared-PC policy; SCN-022; DEC-018 closure | M–L; packaging stays gated |
| WB-11 | Essential notifications/support/privacy/operating hardening; Operations/Engineering | Relevant WB-07–10 data/events; may develop alongside them | Safe notices/reminders, redacted support, privacy propagation, rollout/restore exercises; SCN-023–026/028; EV-06–10 | L; actual reviews/provider/support coverage |
| WB-12 | Controlled first live pilot; Product/Operations/CEO | Candidate included slices complete; mandatory live EV-01–10, applicable G7/G8/G9 decisions | Phased merchant onboarding, imported booking reconciliation, staffed support, stop/rollback plan and weekly scorecard | Unknown calendar duration until approvals; proposed 8–12 weeks |
| WB-13 | Required booking-linked native messaging slice; Messaging/Android/Web/Safety | WB-07/08 authority/data + safety/privacy/support portions of WB-11 | Participant authorization, queue/outbox/delivery, report/block, retention, staff takeover and guest web-session policy if included; SCN-019, EV-11 | L; native endpoint/safety scope |
| WB-14 | Optional first official social connector; Integrations/Operations | Conversation contracts/ownership proven; approved provider account/scopes | Eligibility, webhook replay/rate limits, takeover, editable handoff, disconnect and direct-link fallback; SCN-016/017, EV-12 | M–L; external review timing not controllable |
| WB-15 | Optional verified digital-deposit/payment slice; Payments/Finance/Data | WB-08/09 and actual provider/entity/role agreement; operational WB-11 | Signed evidence, exact amounts, late/duplicate/refund/reversal/settlement exceptions; SCN-014/015, EV-13 | L; legal/provider/finance gating |
| WB-16 | Required in-app voice calling feasibility and delivery; Calling/Safety/Android/Web | WB-13 safe relationship/identity model; call-specific EV-14 | Ring/answer/decline/end/reconnect, background/poor-network/block/report, media-cost/privacy proof; SCN-020; recording off | L–XL; specialized expertise and field network tests |
| WB-17 | iOS and justified packaged Windows/macOS clients; Mobile/Desktop | Stable domain contracts and demonstrated demand/support capacity | Platform links/notifications/calls/accessibility and low-device performance; no copied Android assumptions; native adapters justified | L–XL; per-platform certification/maintenance |
| WB-18 | Adjacent verticals/marketplace/country expansion; Product/CEO/CTO | Retained/supportable core workflow, funded unit economics, new review packets | New archetype/privacy/payment/operational evidence; first market performance not degraded | XL; separate business/release decisions |

WB-13 is not optional platform direction; WB-16 is not forgotten. Their timing is deliberately tied to prerequisites and is not an invented commitment to ship later without ownership. WB-14 and WB-15 may activate before/after WB-12 according to an amended approved pilot contract; base booking cannot depend on them being approved.

## Estimate policy

S/M/L/XL are relative planning sizes, not time promises: one bounded task, one module/surface, a cross-system slice, or a multi-platform/review-heavy programme. Re-estimate only after an accountable implementer inspects concrete contracts and capacity. Each item needs an optimistic/expected/risk estimate and confidence; external approvals are tracked separately from engineering effort.

No credible total delivery date can be calculated from unknown team/budget/provider/counsel/device facts. Do not sum relative sizes or claim the “entire project” takes a few weeks.

## Minimum ticket contract

Every implementation ticket includes:

1. Business objective, actor, scenario and ledger/functional/non-functional requirement IDs.
2. Current scope/release, accountable owner, dependencies and concrete exit tests.
3. Domain states/invariants, API/event/schema version and migration plan.
4. Permission, audit, data purpose/classification/retention and localized copy.
5. Offline/network/provider/unknown/conflict behavior and security-negative cases.
6. UI/accessibility/performance/byte budgets where relevant.
7. Observability, support explanation, cost, failure/rollback and evidence receipt.
8. Reviewed implementation/test artifacts and remaining risk/limitations.

WB packages must be split into such tickets before work is scheduled. References to documents alone are not automated tests.

## First prototype acceptance packet

WB-02: fake studio/service/resource/customer records; persistent local tentative booking/outbox; recover after process death; show Today/details/conflict labels; no real authentication, external notices, live server, money or production claim.

WB-03: disposable isolated database, tenant keys and proposed exclusion/transaction model; simultaneous contenders, hold consumption and retry recovery; no managed production database or provider writes. Explain any mismatch with logical catalogue before approving schema.

WB-04: fake-context opaque grant service; customer edits prefill and reviews distinct request/confirmed examples; no public DNS/domain purchase, social connection or customer data.

The V0 prototypes can be thrown away after lessons are recorded. Do not let demo shortcuts become live authorization or storage behavior. Real D1/network/comprehension evidence cannot be invented from a headless simulation.

## Language sequence and safety

Go serves the candidate API/worker core; SQL enforces proposed scheduling/tenant invariants; Flutter/Dart serves candidate Android; React/TypeScript serves candidate web/desktop. Kotlin/Java and Swift bridge platform requirements. Consider C/C++/Rust only for measured native/media/hotspot needs with safe interfaces, testing/SBOM/ownership; no blanket C/C++ rewrite and no default handwritten assembly. ADR-002 remains the review authority for those choices.

## Definition of done and stop rule

An item is done only with its verified relevant tests/evidence, reviewed artifacts, operational owner and honest known limitations. A release is done only when its exact contract/gates are satisfied and deployed behavior is independently verified.

Pause the affected feature for critical tenant/integrity/false-success defects, missing live authority or uncontrolled spend. Preserve recoverable work, communicate actual impact and invoke runbooks. Document completion may start the next validation task, never silently approve production.

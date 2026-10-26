# Blueprint recheck, contradictions and decision closure

Status: **Draft — verified local inventory; semantic and release audit remains partial**  
Updated: 2026-10-11 (Africa/Addis_Ababa)
Owner: Product/CTO; named sign-off pending

## Recovered position

At the start of this resumption, the saved planning index stopped after Stage 8 drafts dated 2026-09-18. Stage 9 had no documents or executed tests. The project is not a Git repository; `apps/` and `services/` contain no production application implementation. `.tooling/flutter/` is a bundled SDK checkout, not a Ubi application. Its presence alone did not prove a functioning toolchain or an approved framework.

The drive was unmounted when work resumed and was mounted through the normal desktop storage service. Files were then inspected; no app build, deployment, external message, account change, payment or company registration was performed.

## Scope of this recheck

Full document inventory, catalogue/fence/reference/identifier checks; substantive reread of requirements, decisions, programme, scope, key booking/custom/conversation, experience and quality definitions; targeted search of cross-document state/capacity/payment rules. This is not a line-by-line semantic approval of every older document, rendered Mermaid validation, executable-schema audit, threat penetration test or legal review.

Run the reproducible structural check from the repository root:

```sh
node .tooling/verify-blueprint.mjs
node --test .tooling/verify-blueprint.test.mjs
```

The check validates catalogue existence, unique numbered blueprint documents, declared status, balanced Markdown fences, referenced diagram/local Markdown paths, requirement/decision definition uniqueness, and coverage/identifier validity in the Stage 9 crosswalk. Read its limitations as well as its result. It does not execute application tests or validate Mermaid syntax/rendering.

## Findings and disposition

| Finding | Evidence at resume | Disposition |
|---|---|---|
| AUD-01 Missing product experience documents | Catalogue listed 23–26; files absent | Four screen/workflow drafts added; field/design validation still missing |
| AUD-02 Missing business execution documents | Catalogue listed 80–85; files absent | Market, finance, vendor, scorecard, responsibility and backlog drafts added; no commercial/field proof implied |
| AUD-03 Stage 9 absent | No crosswalk, release/evidence contract or closure record | Specifications 71–73 added; G9 remains open |
| AUD-04 Stale country statement | Ledger said launch country unapproved despite DEC-001 Ethiopia | Corrected to Ethiopia decided; city/vertical still proposed |
| AUD-05 Older roadmap contradicted active proposals | Five businesses/four weeks and 70% usability contrasted with ten/8–12 weeks and proposed stricter quality gate | Preliminary notes explicitly archived; active roadmap points to 85, decisions and scorecard |
| AUD-06 Older lifecycle conflated approval/confirmation | Preliminary architecture treated accepted/green as enough and omitted current conditions | Archived/replaced summary refers to current workflow truth; no implementation inherited from it |
| AUD-07 Mandatory native messaging/calls could be mistaken for discarded scope | Founder insists both are inevitable; older documents label future | Requirement retained; explicit backlog workstreams and activation gates; first-release timing remains a review item |
| AUD-08 Calendar support omitted explicit local clock choice | Calendar/timezone rules existed, but ambiguous Ethiopian/international clock wording lacked specific acceptance cases | Specification 22, experience 25 and SCN-018 explicitly separate calendar, language and clock convention |
| AUD-09 Feasibility/release gates risk circularity | “Not approved for implementation” alongside demand for device/provider test results | Contract 72 permits isolated synthetic validation without authorizing real-data/production release |
| AUD-10 Documents and toolchain could be mistaken for built software | Root README showed a planned app tree; directories absent | Root README labels planned paths and actual planning-only state |
| AUD-11 Sky-blue pairing is not automatically accessible | Token `#218ACB` against white computes about 3.78:1; border `#D9DAE5` against white about 1.39:1 | Not approved for normal-sized white button text or sole essential boundary respectively; retain palette with tested dark text/darker semantic variants and redundant cues; full component/state audit still needed |
| AUD-12 Contract/permission vocabulary remains illustrative | Capability/API/schema examples exist across drafts, but no executable registry/contracts | WB-05/06 must reconcile examples into canonical versioned contracts; document-only traceability does not prove compatibility |

These are documentation repairs, not proven repairs to running software. Historical proposals remain distinguishable from founder decisions.

## Decision queue by boundary

| Boundary | Decisions to close or explicitly disposition | Safe work while open |
|---|---|---|
| Research/product | DEC-002, DEC-003, DEC-004, DEC-013, DEC-015, DEC-018 | Draft protocols and dummy prototypes under labelled assumptions |
| Technical prototype | DEC-030, DEC-031, DEC-032, DEC-033, DEC-034, DEC-035, DEC-036, DEC-037, DEC-044 | Compare implementation candidates; no production selection claim |
| Real-data/live pilot | DEC-005, DEC-006, DEC-038, DEC-039, DEC-040, DEC-050, DEC-051, DEC-052, DEC-053, DEC-054, DEC-055, DEC-058 | Assemble review packets; no live processing before relevant approvals |
| Native messaging | DEC-055 plus relationship, safety, identity and support controls | Test scoped threads with fake participants |
| Official connector | DEC-016, DEC-017, provider-specific permissions and capability evidence | Simulator/approved isolated sandbox, direct-link fallback |
| Digital payments | DEC-007, DEC-059, actual provider/commercial/refund arrangements | Synthetic payments; retain no-deposit/cash candidate |
| Native calling | Call identity/consent/media/vendor/cost decisions; DEC-056 recording separate | Synthetic feasibility and weak-network tests; no implied call service |

Missing decisions may be intentionally postponed if their capability is disabled. Essential security/privacy/identity/operating obligations cannot be postponed merely because a release is called a pilot.

## Remaining evidence and honesty boundary

No merchant interview, physical-device result, translated UI review, approved provider account, production migration, restore exercise or licensed counsel sign-off has been observed in this workspace. They are **missing/unverified**, not “failed.” Synthetic Flutter component tests and a disposable Go/PostgreSQL schema/HTTP experiment now exist; they do not replace those missing items. The planned full scenario suite is **not run**. The observed Postman specification-linter failure is recorded separately, not mislabeled as missing or passing.

Current drafting completeness is not full blueprint approval. The next verifiable progress is the bounded validation programme and dependency-ordered build preparation in 85, with receipts tracked in 72. Do not restart Stage 0 or silently mark gates passed just because the catalogue now has files.

The generated `validation/document-structure-report.json` is a dated machine receipt produced by `node .tooling/verify-blueprint.mjs --report`. Its source digest fingerprints the blueprint Markdown and Mermaid content, not an approved release. Re-run after edits; this file does not certify anything beyond the listed document checks. The checker now has nine synthetic positive/negative tests, including missing files, invalid references, duplicate IDs, false prefix coverage and prototype links outside the blueprint directory. These are tooling tests, not the 28 Ubi application scenarios.

## Subsequent synthetic implementation checkpoint

2026-10-15 additions: [one authorized physical Android 12 session](../validation/WB-02-physical-android-smoke-evidence.md) adds bounded native SQLite and single-force-stop recovery checks for the original 0.1.0 APK. The later [Home action follow-up](../validation/WB-02-home-action-layout-evidence.md) fixes the observed overlap in source/host tests (41 passing tests) and prepares a separate 0.1.1 candidate awaiting renewed device-update/recheck authorization. The phone remains unchanged; original APK/inputs/evidence are archived. Minimum-device, profile/release performance, accessibility, broader native acceptance and real integration remain unverified. Earlier dated checkpoints below are historical, not claims that no phone has ever been tested.

At the earlier host checkpoint, WB-01 fixtures/protocol and the WB-02 isolated Flutter/SQLite prototype existed under `prototypes/`; Flutter 3.47.1 / Dart 3.13.1, analyzer and 26 local repository/widget tests passed, with no APK or physical-device result at that time. The subsequent 2026-10-11 checkpoint adds two startup-asset tests (**28/28**, analyzer clean), explicitly approved isolated JDK/Android setup and a corrected [debug APK with static inspection](../validation/WB-02-synthetic-android-apk-evidence.md). Failed/superseded attempts and the asset-packaging/parser corrections are recorded; native execution and physical-device measurements remain missing. [The host evidence record](../validation/WB-01-02-synthetic-prototype-evidence.md) preserves the narrower scope and gaps. No production API or live account/data exists. The founder decided Android 11+/4 GB for the business app (DEC-022), superseding the old provisional Android 10/2 GB D1 baseline. Framework/database/release proposals are not thereby approved.

WB-03 now has a separate Go/PostgreSQL/loopback HTTP implementation and repeatable component checks using only a newly created private synthetic cluster. [Its receipt](../validation/WB-03-synthetic-integrity-evidence.md) records integrity, retry, RLS/privilege and contract evidence plus the Postman spec-linter limitation. Temporary database cleanup is checked; no existing database or public endpoint was used. The 2026-10-11 recheck again passed 75/75 component case results; the lint disagreement remains open.

WB-04 now has [customer-link component evidence](../validation/WB-04-synthetic-customer-link-evidence.md) for scoped ephemeral handoff/edit/review, distinct enquiry/instant outcomes, actual lost-response retry and no-JavaScript browser paths. Recheck repaired expired-unknown retry authority, incomplete review details and a real browser no-referrer/Origin incompatibility. The customer experiment's own HTML contract lint passes; it does not resolve WB-03's problem-JSON linter disagreement. Android/server integration and persistent/verified public guest flows remain unimplemented; all full acceptance scenarios remain NOT RUN. No receipt closes G7/G8/G9 or freezes the stack/schema/API.

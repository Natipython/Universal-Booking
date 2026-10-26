# Ubi decision register

Status: **Draft**  
Purpose: prevent important choices from being hidden in code, vendor accounts, database migrations, or casual conversation.

## Decision levels

- **Founder/board** — market, ownership, brand, capital, risk appetite and irreversible commercial commitments.
- **Product/business** — customer, value proposition, release scope, pricing and operating policy.
- **Technology/data/security** — architecture and control decisions within approved business constraints.
- **Legal/compliance** — jurisdiction-specific interpretation reviewed by licensed counsel.
- **Operational** — runbook and service decisions inside approved policy and risk boundaries.

## Blocking launch decisions

| ID | Decision | State | Why it blocks downstream work | Required evidence/owner |
|---|---|---|---|---|
| DEC-001 | Launch country | Decided — Ethiopia (founder, 2026-09-05) | Ethiopia is the first-launch jurisdiction; expansion countries remain separate gated decisions. | Founder decision recorded; Ethiopian legal/regulatory verification continues |
| DEC-002 | Launch city | Proposed — Addis Ababa | Founder proximity, merchant density, connectivity/provider access and field support make Addis the practical first validation city; this still needs founder confirmation. | Founder/product |
| DEC-003 | First vertical | Proposed — photo studios/photographers | Matches the original use case, supports appointments and custom event quotations, and avoids hotel/health complexity; validate with real merchant interviews before approval. | Founder/product |
| DEC-004 | Pilot size and duration | Proposed — 10 businesses for 8–12 weeks | Small enough for founder-led onboarding/support and large enough to expose workflow differences; expand only against agreed success/stop metrics. | Founder/product/operations |
| DEC-005 | Contracting legal entity | Proposed — Ethiopian private limited company | Recommended working structure: `Ubi Technologies` as a one-person private limited company if there is one founder, or a private limited company if there are multiple founders; exact name, ownership, capital, objects, registration and tax treatment require founder facts and Ethiopian counsel/accounting. | Founder ownership decision + CEO + Ethiopian corporate counsel/accountant |
| DEC-006 | Data-controller/processor roles | Proposed — role varies by activity | Ubi is likely controller for Ubi Accounts, subscriptions, security, platform abuse and its own operations; each merchant is likely controller for its customer/service relationship, with Ubi processing booking/customer records on documented instructions. Providers may have independent roles. Final classification is per activity, not one label for the company. | Privacy architecture + Ethiopian counsel |
| DEC-007 | First-release payment model | Proposed — non-custodial licensed-provider orchestration | Customer pays the merchant through an NBE-licensed provider; Ubi records initiation/provider status and reconciliation but does not hold, pool, store, split or settle customer funds. Start with separate merchant SaaS fees. | Founder confirmation + provider due diligence + CEO/CTO/finance + Ethiopian payments counsel |

## Product decisions

| ID | Decision | State | Current direction | Reversal cost |
|---|---|---|---|---|
| DEC-010 | Brand name | Working | Ubi; conduct trademark/domain/company-name checks before commitment. | Medium/high after launch |
| DEC-011 | Product surfaces | Direction stated | Ubi Business Android first; desktop business workspace next; customer web/app sequence requires workflow review. | Medium |
| DEC-012 | Marketplace timing | Proposed | Begin as merchant booking infrastructure; marketplace only after supply density and retention. | Low before marketplace investment |
| DEC-013 | Customer account requirement | Proposed | Guest booking first; phone-verified identity promotion later. | Medium if identity model is wrong |
| DEC-014 | Universal configuration model | Proposed | Reusable booking archetypes and modules, not unrestricted custom code/workflows. | High after merchant configuration data exists |
| DEC-015 | Confirmation policy | Open | Support service/operation configuration for auto-confirm versus request-to-confirm. | Medium |
| DEC-016 | Initial automation boundary | Open | Start with structured assistance and deterministic rules; human takeover mandatory. | High if customers experience incorrect automation |
| DEC-017 | First connected channel | Proposed | Validate channel usage commercially; Telegram is technically tractable, WhatsApp may have stronger market priority. | Medium |
| DEC-018 | Desktop form | Open | Compare responsive web console/PWA with native desktop requirements through research. | High if two client stacks are committed too early |
| DEC-019 | Native in-app messaging release boundary | Direction stated; release timing open | Required platform capability, initially booking-linked/scoped. WB-13 and EV-11 define safety/delivery/privacy validation; first-live-pilot inclusion requires an explicit release decision, not silent omission. | High after participant/message data exists |
| DEC-020 | Native in-app voice calling release boundary | Direction stated; release timing open | Required platform capability. WB-16 and EV-14 cover feasibility, identity, poor-network/background behavior, safety, consent and cost; external phone buttons are not Ubi-native calling. Recording remains separate DEC-056 scope. | High after real-time media/vendor commitments |
| DEC-021 | Calendar and clock presentation convention | Open | Separate language, Gregorian/Ethiopian calendar, international/Ethiopian clock presentation and timezone. Ambiguous conversation times require clarification; no automatic inference from language. | High if accepted instants are misinterpreted |
| DEC-022 | Minimum supported business Android baseline | Decided — founder, 2026-10-04 | Android 11 or later and at least 4 GB RAM. Replaces provisional Android 10/2 GB D1. Actual locally representative hardware, accessibility and performance still need measurement; RAM is a support/test baseline, not a startup denial. Public guest-web compatibility remains a separate decision. | Medium; excludes older/lower-memory merchant devices from supported scope |

## Architecture decisions requiring records

| ID | Topic | State | Options to evaluate |
|---|---|---|---|
| DEC-030 | Initial architecture style | Proposed | Modular monolith plus workers versus distributed services. |
| DEC-031 | Mobile framework | Proposed | Flutter/Dart for Ubi Business mobile, Android first; Kotlin/Java Android and Swift Apple adapters only where platform APIs/plugins require them. Approval requires D1 offline/performance/accessibility prototype evidence. |
| DEC-032 | Desktop/web framework | Proposed | Server-rendered, progressively enhanced React/TypeScript web/PWA for customer, business desktop and operations surfaces; Flutter desktop remains a later measured option, not a pilot dependency. |
| DEC-033 | Backend language/runtime | Proposed | Go for the modular application/API and workers. PostgreSQL/SQL remains data truth; Rust/C/C++ are isolated native/tooling options, not alternate booking backends. |
| DEC-034 | Authoritative operational database | Proposed | PostgreSQL versus alternatives after workload review. |
| DEC-035 | Tenant isolation model | Open | Shared schema with tenant keys/RLS, schema isolation, database isolation, or tiers. |
| DEC-036 | Offline local database and sync protocol | Open | SQLite-family choices, command/outbox protocol and conflict model. |
| DEC-037 | API style | Open | REST/JSON, GraphQL, gRPC internally, event contracts, or justified combination. |
| DEC-038 | Cloud/region/deployment provider | Open | Region availability, data transfer, cost, managed services, exit strategy and skills. |
| DEC-039 | Identity/OTP provider strategy | Open | Build core identity plus providers versus managed identity service. |
| DEC-040 | Notification providers | Open | Push, SMS and channel notification vendors per country. |
| DEC-041 | Media storage/processing | Open | Object storage, encryption, lifecycle, moderation and transformations. |
| DEC-042 | Search architecture | Open | Database-native search initially versus dedicated engine. |
| DEC-043 | Analytics architecture timing | Open | Operational reporting first; warehouse/lake only at justified volume/complexity. |
| DEC-044 | Native/low-level language boundary | Proposed | New C/C++/Rust code needs measured CPU/memory/latency need, stable C ABI where shared, fuzz/sanitizer/security review and ownership. Prefer Rust for new memory-sensitive code where supportable; use C/C++ for proven libraries/interop. Handwritten assembly is prohibited by default. |

## Trust, legal and data decisions

| ID | Topic | State | Required analysis |
|---|---|---|---|
| DEC-050 | Data classification | Open | Public, internal, confidential, sensitive personal, credential, financial and audit categories. |
| DEC-051 | Retention schedule | Open | Purpose, law, merchant needs, disputes, deletion, backups and analytics. |
| DEC-052 | Data residency/transfers | Open | Launch law, provider regions, safeguards, subprocessors and customer commitments. |
| DEC-053 | Support access model | Open | Just-in-time access, merchant consent, break-glass, audit and monitoring. |
| DEC-054 | Age/minor policy | Open | Booking context, parental authority, photography/health sensitivity and country law. |
| DEC-055 | Messaging consent and retention | Open | Ubi messaging, connected channels, notification templates and opt-out. |
| DEC-056 | Call consent/recording policy | Deferred | Recording default should be off unless a separately justified lawful use exists. |
| DEC-057 | AI processing and vendor use | Deferred/open | Purpose limits, disclosure, human control, training use, retention, cross-border processing and model safety. |
| DEC-058 | Merchant/customer dispute allocation | Open | Platform role, merchant responsibility, cancellation/refund rules and escalation. |
| DEC-059 | Payment custody/settlement role | Proposed — Ubi is not custodian or payment-system operator | Merchant/provider remains payee/settlement path; Ubi is booking software and technical integration layer. Any later wallet, pooled collection, split settlement, money movement or gateway role is a new regulated-scope decision. |

## Decision record template

Each approved decision receives `decisions/ADR-NNN-short-title.md`:

```text
Title
Status and decision date
Decision owner and required approvers
Problem/context
Constraints and requirements affected
Options considered
Decision
Rationale and evidence
Security/privacy/legal/data consequences
Operational and cost consequences
Failure modes and mitigations
Migration/reversal plan
Validation and review date
Related requirements, diagrams and implementation items
```

## Change-control rules

The 2026-10-04 Stage 9 package proposes release scope and a decision queue in specifications 72/73. It does not approve existing technology, entity, city, pilot or payment proposals. Founder silence is not approval. Review each decision against the exact synthetic, real-data or production boundary affected.

1. An approved decision may be challenged with evidence; it is not silently overwritten.
2. Security, privacy, payment or data-retention impact triggers specialist review.
3. Breaking API, event, schema or policy changes require compatibility/migration plans.
4. Vendor changes require data export, deletion, credential revocation and rollback analysis.
5. Emergency operational changes are documented retrospectively within the incident process.
6. A decision is reviewed when its evidence, jurisdiction, scale, vendor terms or risk assumptions materially change.

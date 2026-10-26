# Partnerships and vendor selection strategy

Status: **Draft — no vendor shortlist winner, account approval or contract**  
Updated: 2026-10-04 (Africa/Addis_Ababa)  
Depends on: decisions 38–43, policies 52–55, contracts 72, cost model 81

## Select by demonstrated launch need

The first partnerships are willing pilot merchants, locally competent support/research people, translation reviewers and qualified Ethiopian legal/accounting reviewers. A cloud/payment brand name does not solve merchant adoption or establish local permissions.

Run technology procurement after a bounded requirement and cost/data envelope exist. Do not install every plugin, create vendor accounts or sign contracts just because the blueprint lists a future capability.

## Selection matrix

| Need | Evidence required before use | Initial fallback |
|---|---|---|
| Hosting/managed PostgreSQL | Region/data transfer, latency, backup/PITR/restore, costs, availability, exit export | Isolated local synthetic development; no production resilience claim |
| Identity/OTP | Ethiopia delivery/eligibility, abuse spend limits, recovery, privacy, actual routing tests | Demo identities for synthetic tests; never accept fake OTP in live auth |
| Push/essential notices | Supported devices, endpoint lifecycle, consent/purpose, failure evidence/cost | In-app queue and approved human follow-up |
| Domain/DNS/TLS | Name/trademark/renewal facts, ownership, wildcard/association behavior, credential isolation | Local named test hosts; no claim domain purchased |
| Object/media | Access scope, upload limits/scan, region, lifecycle/egress/deletion | Text-only critical flows; optional media disabled |
| Official channel APIs | Merchant eligibility, scopes, review, capability, support, revocation, rates/windows and terms | Public/private Ubi links plus manual capture |
| Licensed payment provider | Actual license/current services, merchant onboarding, non-custodial role, payee/settlement, signed status, refunds/reversal/reconciliation and costs | No deposit/cash/pay-at-location candidate |
| Call signaling/media | Identity/privacy, local network feasibility, relay/egress price, background support, abuse and exit | Explicit external phone fallback, not a Ubi call |
| Observability/support | Redaction, location/retention, security, sampling/cardinality/cost and export | Bounded first-party metadata plus practiced runbooks |

Fallbacks themselves need applicable live authorization and operational support. A direct phone fallback does not resolve a failed payment or establish capacity.

## Vendor scorecard

Pass/fail controls first: country/account eligibility, actual permissions, security/privacy requirements, required function, acceptable exit and funded cost. Then score usability/onboarding, observed latency/reliability, support, total recurring cost, integration complexity and recovery. Document weights before selecting; a cheap plan cannot compensate for failed essential legal/data controls.

Keep evidence dated: official documentation, terms, provider responses, test account IDs without secrets, test receipts and signed agreements. Provider policies/capabilities may change; recheck before public claims, submission or production activation.

## Language/tooling boundary

Go + Flutter/Dart + React/TypeScript + PostgreSQL remains ADR-002's proposal. C/C++ and Go are considered intentionally, not forgotten. Java/Kotlin are Android interoperability options; Swift is Apple platform work; Rust is a supported-team candidate for new memory-sensitive native components. React is a UI library; Flutter a framework; neither is a low-level language replacement.

Use a measured hotspot or mature required library to justify C/C++ rather than rewriting business CRUD in multiple languages. Go/native integration needs the ownership/lifetime/pointer rules described in [the official cgo documentation](https://pkg.go.dev/cmd/cgo), rechecked 2026-10-04; direct or ABI-wrapped foreign calls add build and memory-boundary obligations. Prefer simple service boundaries when justified by measurements and ownership. No native code or handwritten assembly is approved merely by listing it.

## Lifecycle and exit

1. Assign accountable owner and exact purpose/data fields.
2. Complete policy 55 due diligence and role/location/contract review.
3. Use isolated synthetic/approved sandbox credentials; store secrets in the designated secret system.
4. Prove failure, rate-limit, replay, revocation and reconciliation paths.
5. Record funded prices, quotas, alerts, renewal/review date and live activation gate.
6. Exercise export/replacement and deletion/revocation, including downstream data and callback shutdown.

No production credentials, social passwords, session cookies or OTP codes belong in planning documents or support messages.

## Partnership economics

Merchant associations, photography networks and service providers may later reduce acquisition/onboarding cost. Test a small opt-in partnership before exclusivity, reseller access or referral commissions. Share only approved business contact information; referral permission is not consent to unrestricted customer data. Quantify support burden and reputation risk alongside acquisition benefit.

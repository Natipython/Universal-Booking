# ADR-001 — Ethiopia as Ubi's first-launch country

Status: **Approved founder decision**  
Decision date: 2026-09-05 (Africa/Addis_Ababa)  
Decision owner: Founder  
Required downstream reviewers: CEO/product, CTO/security/data, CFO/accounting, operations and licensed Ethiopian counsel

## Problem and context

Ubi cannot responsibly design one undifferentiated “East Africa” launch. Country selection affects company formation, privacy, payments, tax/records, telecom and messaging, languages/calendars, providers, contracts, support and merchant operations.

## Decision

Ubi's first-launch country is **Ethiopia**.

This decision does not yet select:

- launch city;
- first merchant vertical;
- contracting legal entity or ownership structure;
- tax/business-license classification;
- payment provider or regulated payment role;
- cloud/data location or cross-border transfer arrangement;
- launch languages, pilot size or public-launch date.

It does not authorize expansion to another country. Each additional country repeats the legal, commercial, provider and operational gate.

## Rationale

- It narrows an otherwise unbuildable multi-jurisdiction design.
- The founder explicitly selected Ethiopia.
- Ubi is designed around problems visible in Ethiopian/East-African merchant booking: calls, walk-ins, messaging, cash/digital payment evidence, constrained devices/networks and assisted onboarding.
- Country-first decisions allow realistic provider, legal, language, payment and support validation.

## Consequences

- Ethiopian law and regulator/provider reality become the first-release baseline.
- Ethiopia's Personal Data Protection Proclamation No. 1321/2024 must be included; “no general data-protection law” is not a valid assumption.
- NBE-licensed payment providers and current NBE directives are the initial payment environment; Ubi's proposed model remains non-custodial.
- Product testing must use representative Ethiopian merchants, staff, customers, Android devices, networks, languages, calendars and payment behavior.
- The remaining candidate-country matrix is expansion research, not launch scope.

## Risks and mitigations

| Risk | Mitigation |
|---|---|
| Treating all Ethiopia as one identical market | Select one city and vertical; conduct observed workflow research |
| Late legal redesign | Complete internal architecture first, obtain Ethiopian design review before irreversible implementation/contracts, and final counsel approval before production launch |
| Ubi accidentally becomes a payment operator/custodian | Keep merchant/provider settlement direct; prohibit Ubi-held/pooled/stored/split funds unless separately approved/licensed |
| Provider/API not available for entity/use case | Verify account eligibility and production capability before promising it |
| Product too complex for target users | Assisted onboarding, plain local language, low-bandwidth design and measured field testing |

## Validation and review

- Founder confirms city, first vertical, ownership/founders and pilot model.
- Ethiopian corporate/accounting advice confirms entity, business objects, registrations and tax/invoicing setup.
- Ethiopian privacy/payments/consumer/telecom counsel reviews the completed release boundary.
- Selected payment/channel/cloud providers confirm eligibility in writing or in verified production onboarding.
- Review this ADR if the first launch moves outside Ethiopia.

## Related artifacts

- `02-decision-register.md` — DEC-001 through DEC-007
- `policies/52-jurisdiction-and-regulatory-matrix.md`
- `policies/56-legal-document-register.md`
- `policies/57-records-audit-and-law-enforcement-process.md`


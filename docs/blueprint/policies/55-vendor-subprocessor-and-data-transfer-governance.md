# Vendor, subprocessor and data-transfer governance

Status: **Draft — governance framework; no cloud/provider approved**  
Depends on: architecture, data inventory, threat/privacy model, jurisdiction matrix, recovery and cost model  
Blocks: vendor accounts/contracts, production data, country launch and provider integrations

## Purpose

Ubi will depend on cloud hosting, identity, object storage, notifications, SMS, messaging, payments, calls, analytics/crash reporting, support and possibly AI. Every provider changes security, privacy, availability, legal, unit-economic and exit risk.

A familiar brand, certification badge or cheap introductory price is not approval. The exact product, legal entity, region, configuration, data flow and contract must be reviewed.

## Principles

- Buy commodity capability when it improves safety/time-to-market; keep Ubi domain truth portable.
- Use the minimum vendor and SDK set for pilot.
- No vendor receives data without an inventory purpose, role, classification, retention and deletion path.
- “Region selected” does not prove all support, logs, backups, subprocessors or control-plane access stay there.
- Encryption at rest does not solve provider access or unlawful transfer by itself.
- A vendor’s compliance report supports due diligence; it does not transfer Ubi/merchant accountability.
- Critical providers need tested degradation, reconciliation and exit—not only a contract.
- Provider price includes messages, calls, egress, storage, transformations, support, currency/tax and engineering operations.
- Merchant-connected platforms may be independent parties as well as processors; classify per activity.
- Policy/terms/API changes are production changes even when Ubi deploys no code.

## Vendor inventory classes

- managed PostgreSQL, compute/network/DNS/CDN/WAF/object/KMS/secrets/backups;
- identity/OTP/email/SMS/push and telecom sender services;
- WhatsApp/Meta, Telegram, TikTok and future channel/solution providers;
- payment gateways/mobile money/banks/reconciliation/KYB providers;
- call signaling/media/TURN/recording/transcription;
- analytics, crash/performance, logs/traces/security monitoring;
- customer support/ticketing/moderation/communications;
- AI/model/inference/vector/evaluation services;
- source control/CI/CD/build signing/dependency registries;
- accountants, lawyers, support contractors and other people-processors;
- open-source dependencies and hosted APIs with material supply-chain risk.

## Criticality tiers

| Tier | Consequence of failure/compromise | Examples | Minimum direction |
|---|---|---|---|
| V0 critical | Booking/payment/identity truth, keys or broad confidential data unavailable/compromised | Primary database/cloud/KMS/identity/payment | Executive/security/privacy/legal approval; tested recovery/exit; strong contract/monitoring |
| V1 high | Core workflow blocked or significant content exposed | Object storage, channel provider, SMS/OTP, operations/support | Degraded fallback, incident/data/exit terms and periodic review |
| V2 medium | Important but owner truth remains available | Analytics/crash/search/notification add-on | Minimized data, kill switch, export/delete and acceptable outage |
| V3 low | Easily replaceable/non-production with minimal data | Design/research utility | Basic security/privacy/license/cost review; no production secrets/data by default |

Criticality is based on actual configuration/data, not vendor category. One analytics SDK collecting message bodies would become high risk.

## Vendor register

Each product/account has:

```text
vendor and contracting/service legal entity
product/API/version/account owner/environment
business capability and accountable Ubi owner
criticality, data inventory/purposes/classification/subjects
Ubi/merchant/vendor legal roles and instructions
countries/regions: primary, replica, backup, logs, support, subprocessors
security controls/reports/pen-test/incident history and shared responsibility
authentication/admin/SCIM/MFA/JIT/API credentials
retention/deletion/export/portability and backup behavior
subprocessors and change-notice/objection process
transfer mechanism/assessment and government-access analysis
availability/SLA/support/RPO/RTO/status/incident contacts
pricing, currency, tax, minimums, egress and cost limits
terms/AUP/platform policies/licensing/IP/AI training
fallback, reconciliation, exit trigger/test and deletion certificate/evidence
approval/effective/renewal/review dates and open risks
```

No personal spreadsheet becomes the sole register; controlled documents/database records and reminders have named backups/owners.

## Selection process

1. Define required capability, data, regions, SLO, volume, cost and exit before looking at vendors.
2. Consider build, managed provider, official platform directly and licensed intermediary.
3. Shortlist only products available for the launch country/legal entity/use case.
4. Run security/privacy/legal/financial/operational due diligence.
5. Prototype latency, Android/network behavior, API correctness, limits, onboarding and failure.
6. Model realistic pilot and growth cost per active merchant/completed booking/message/payment/media GB.
7. Review contract/DPA/terms/subprocessors/transfer/incident/deletion/exit.
8. Record alternatives, risk, reversibility and approval in ADR/vendor register.
9. Use non-production/synthetic data until approved.
10. Configure least privilege, monitoring, budget, backup/fallback and evidence before production.

## Security due diligence

Review exact evidence and scope:

- security governance/ownership and independent assurance scope/date/exceptions;
- encryption/key options, tenant isolation, network/private access and secrets;
- authentication/MFA/admin roles/service accounts/audit/export;
- secure development, vulnerability disclosure, patching and dependency controls;
- incident detection, notification commitments, evidence/cooperation;
- employee/contractor access, background checks where lawful, JIT and support access;
- data deletion/media sanitization and backup retention;
- availability architecture, capacity, DDoS, backup/PITR/restore and disaster tests;
- penetration-test/customer testing rules and recent material incidents;
- API/webhook/OAuth/idempotency/replay/version/deprecation behavior;
- customer-managed safeguards Ubi must configure.

Certifications are evidence only if the exact service/entity/region/control period and exceptions apply.

## Privacy and contractual due diligence

- processing role and written instructions/purpose restrictions;
- confidentiality/security and assistance with rights/DPIA/incidents/regulators;
- data categories/subjects, sensitive/child data restrictions;
- exact retention, deletion/return, backups and verification;
- all subprocessors, locations and advance material-change notice;
- cross-border transfer basis/safeguards/localization and remote support;
- government/law-enforcement request policy, challenge and transparency where lawful;
- secondary use, advertising, sale/sharing, AI/model training and service-improvement terms;
- audit/information rights proportionate to risk;
- liability, indemnity, insurance and order-of-precedence;
- termination, data export format/time/cost and continuing obligations.

Contract promises must match technical configuration. A DPA stating short deletion is useless if the enabled product retains debug payloads longer.

## Data-location and transfer assessment

For each path, record:

```text
exporting entity/country and subjects
receiving vendor/subprocessor legal entity/country
service/primary/replica/backup/log/support/control-plane locations
data/purpose/classification/volume/frequency and encryption/key holder
applicable localization/transfer rules and approved mechanism
destination law/government-access and practical risk assessment
technical/contractual/organizational supplementary controls
rights/incident/deletion ability
residual risk, approver and review trigger
```

Do not describe data as “stored in Africa” or “never leaves country” unless all copies/access/support/logs/subprocessors and contract make it true and monitored.

## Provider access and credentials

- Separate production/staging/development accounts/projects/keys.
- Named workforce access, phishing-resistant MFA where practical, least privilege and no shared vendor admin.
- SSO/SCIM/JIT where risk/scale justifies; quarterly or tighter critical-access review.
- Service identity uses minimum scope and short-lived/federated credentials where supported.
- Provider tokens/keys enter managed secret/KMS design and never merchant/support chat.
- Alert on admin, billing, region, webhook, key, retention, export and subprocessor configuration changes.
- Preserve break-glass access under dual control outside one failed identity path.

## SDK and mobile/web supply-chain gate

An SDK can execute inside Ubi and collect data outside API expectations.

- Prefer server-side API over client SDK when capability permits.
- Inventory exact package/version/transitive dependencies/permissions/network endpoints/data collection.
- Verify publisher/signature/provenance/license/update history and vulnerability process.
- Disable automatic broad data collection, session replay, advertising IDs and content capture.
- Wrap SDK behind adapter/feature flag/consent and kill switch.
- Test Android startup, memory, battery, bytes, crash and low-end device impact.
- Pin/lock and review updates; no remote code/config enabling unreviewed collection.
- Remove SDK and verify outbound traffic/data deletion during exit test.

## AI/model provider gate

Before any merchant/customer content is sent:

- exact use case and minimum input/output;
- whether provider stores, trains, human-reviews or reuses data, including abuse logs;
- region/subprocessors/transfer and deletion/retention controls;
- tenant isolation, prompt/output logging and secret filtering;
- model/version change and availability/cost controls;
- safety/evaluation for launch languages and high-impact error classes;
- contractual IP/confidentiality/liability and opt-out;
- non-AI/human fallback and rapid disable;
- prohibition on high-risk autonomous booking/payment/ownership/restriction action.

AI is deferred until this gate passes; a zero-retention marketing statement is verified contractually and technically.

## Concentration and dependency risk

Map shared failure domains: one cloud for API/database/object/KMS/backup; one Meta account for multiple channels; one phone/SMS provider for login and notifications; one identity administrator; one payment rail; one DNS/registrar account.

For critical concentration:

- decide accepted risk versus separate control/account/region/provider;
- maintain status/contact/manual fallback and recovery priority;
- avoid fake multi-provider complexity before pilot, but preserve adapter/export boundary;
- test loss/suspension of vendor and billing/admin account;
- cap spend/traffic without silently dropping security/transactional work.

## Ongoing monitoring and changes

Monitor uptime/latency/errors/limits/cost/security advisories, terms/policy/API versions, subprocessors/locations, certificates/reports, incidents and financial viability.

Material vendor change triggers impact review, test and possibly merchant/person notice or reauthorization. Examples: new country/subprocessor, training/secondary-use term, shortened API deprecation, price increase, lost certification, breach, account-model/identifier change or acquisition.

## Exit and replacement

Every V0/V1 provider has an exit plan and test:

1. inventory data/config/identifiers/templates/webhooks/keys/dependencies;
2. export in usable documented format with checksums/counts;
3. dual-read/capability switch only when consistency allows;
4. migrate/reconcile owner state and external unknown outcomes;
5. update DNS/routes/apps/providers without split brain;
6. revoke credentials/webhooks/admins and stop billing;
7. delete/return data including backups on contractual schedule;
8. obtain/record deletion evidence and monitor late callbacks;
9. preserve only justified legal/audit facts;
10. communicate realistic service changes.

Provider-scoped external IDs remain namespaced; migration cannot pretend they are portable identities.

## Release gates

- Every production vendor/SDK/contractor exists in register with owner, tier, purpose, data, location, terms and cost.
- Launch-country availability and legal entity/account eligibility are verified.
- V0/V1 security/privacy/legal/financial review and signed contracts/DPA are complete.
- Required transfer/localization assessment and regulator/merchant approvals exist.
- Least-privilege accounts/keys, monitoring, budget alerts, incident contacts and kill switch are configured.
- Backup/recovery/fallback/reconciliation and exit/export/delete are tested proportionate to tier.
- No vendor uses merchant/customer data for advertising/general model training outside approved purpose.
- Android SDK impact/data/permissions/outbound traffic are verified on low-end devices.
- Marketing/data-location/SLA claims match evidence.

## Open decisions

- Launch country/legal entity and selected cloud region/provider.
- Identity/OTP, payment, notification/SMS and first channel providers.
- Whether an intermediary/solution provider improves Meta/TikTok onboarding enough for its cost/data dependency.
- Observability/crash/support/security tooling minimum set.
- AI and native call vendors remain deferred.
- Required separate-account/region/off-provider backup level.
- Vendor review staffing, renewal calendar and acceptable concentration risk.

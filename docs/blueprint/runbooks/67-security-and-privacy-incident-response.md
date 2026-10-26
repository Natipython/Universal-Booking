# Runbook 67 — Security and privacy incident response

Status: **Draft — Ethiopian counsel/regulatory procedure pending**  
Owners: Security incident lead and privacy lead  
Applies to: suspected/confirmed unauthorized access, disclosure, alteration, destruction, credential compromise, abuse of privilege or personal-data incident

## Safety rule

Contain quickly but preserve evidence and affected-person rights. Do not wipe devices/logs, contact an attacker, notify externally, promise “no breach,” or broadly search production without incident authority.

Diagram: `diagrams/67-security-privacy-incident-flow.mmd`.

## Triggers

- cross-tenant/customer/merchant data visible to wrong party;
- stolen account/session/owner recovery or privileged-support misuse;
- provider/OAuth/payment/cloud/signing/database credential exposure;
- suspicious export, bulk access, malware, dependency/build compromise;
- lost/stolen device with Ubi production data/access;
- misdirected notification/export/backup/log/analytics data;
- vulnerability exploitation or credible researcher report;
- vendor reports security/privacy incident;
- integrity/destruction/ransomware or unlawful processing/retention.

When uncertain, preserve and classify; do not wait for proof before stopping active harm.

## Specialized roles

- Security IC/technical containment lead.
- Privacy lead for processing/data-subject impact, merchant roles and rights.
- Evidence custodian/forensics with chain of custody.
- Legal/counsel for privilege, Ethiopian/current-country duties, law enforcement and notification.
- Affected domain/tenant/provider owner.
- Communications/support/executive leads.

Only authorized counsel/decision owner determines legal notification classification/deadline. No universal breach deadline is hardcoded before review.

## Immediate actions

1. Open restricted incident/case; use out-of-band channel if identity/collaboration may be compromised.
2. Name roles, severity and exact systems/accounts/data suspected.
3. Preserve volatile and durable evidence: identity/admin/cloud/audit/deploy/access/provider logs, artifact/config versions and time source.
4. Stop active access/exfiltration/unsafe processing with narrowest effective containment.
5. Revoke sessions/tokens/keys or isolate workload/account only through tracked action; retain metadata/evidence.
6. Protect unaffected tenants and critical booking continuity; do not expose incident details through broad banners prematurely.
7. Notify counsel/privacy/executive/vendor/merchant controller according to role and credible scope.
8. Establish next assessment/communication time.

## Evidence and forensics

- Record who collected what, source, time/timezone, method, checksum, storage and every transfer/access.
- Prefer snapshots/export copies and provider preservation over modifying originals.
- Document gaps, log retention/sampling/clock drift and whether attacker could alter evidence.
- Limit evidence to incident scope; unrelated customer content is not investigative entitlement.
- Malware/native binaries are analyzed in isolated environment.
- Legal hold is scoped/reviewed; incident does not justify indefinite retention.

## Assessment matrix

Determine and update:

- confidentiality, integrity and availability impact;
- data categories/sensitivity/content versus metadata;
- people/merchants/tenants/countries/children or vulnerable groups;
- approximate records and exposure window;
- access/recipient/actor intent and likelihood of acquisition/use;
- encryption/key/access state and effective protections;
- booking/payment/safety/identity/financial/reputational harm;
- whether data was Ubi-controller, merchant-controller/Ubi-processor, provider-held or mixed;
- recoverability and continuing risk;
- regulator, merchant, affected-person, provider, insurer/contract/law-enforcement duties.

Unknown values remain `unknown`; they are not counted as zero.

## Containment playbooks

### Account/session takeover

- revoke affected sessions/refresh credentials and high-risk pending changes;
- secure primary recovery factors and freeze ownership/payment/provider connection changes;
- review membership/export/customer/booking/payment/provider actions;
- notify legitimate owner through independently verified channel;
- restore access through enhanced recovery with cool-down/redress.

### Provider/OAuth/API credential

- identify exact app/environment/tenant/scopes and last valid use;
- revoke/rotate with provider and disable callbacks/sends as needed;
- search for unauthorized actions and reconcile provider truth;
- repair secret distribution and remove exposed copies/history/logs/artifacts;
- reauthorize merchants transparently where required.

### Cross-tenant authorization/data leak

- disable affected route/query/cache/export/projection immediately;
- preserve request/auth/RLS/cache/release evidence;
- enumerate affected tenant-object pairs using independently reviewed query;
- test adjacent routes/stores/analytics/backups;
- correct authorization/data safely and verify negative isolation tests.

### Build/dependency/native compromise

- stop releases/distribution, revoke signing/CI/cloud credentials at risk;
- identify artifact digests/versions/devices/environments;
- reproduce/compare provenance/SBOM/source and isolate malicious component;
- ship signed clean replacement and server-side containment;
- assess data/actions possible from compromised client/native/server context.

### Malicious insider/support access

- suspend scoped access without alerting in a way that destroys evidence;
- preserve workforce identity/JIT/case/query/export/device records;
- involve executive/legal/HR/security under need-to-know;
- protect reporters and affected merchants/users;
- review shared access/control failures, not only individual action.

### Ransomware/destructive change

- isolate affected identities/systems and protect backup controls;
- invoke runbook 70; do not connect clean backups to compromised control plane;
- determine exfiltration as well as encryption/destruction;
- rotate trust roots and validate recovery before traffic.

## Eradication and recovery

- close root access path and adjacent variants; patch/configure/rebuild from trusted source;
- rotate credentials in dependency order and invalidate old sessions/links/keys;
- restore/correct only from validated clean source; preserve incident copies;
- validate tenant isolation, permissions, audit, deletion/hold and provider outcomes;
- monitor attacker indicators and affected cohorts with proportionate duration;
- remove temporary access/tools/data and release scoped legal holds when approved.

Reopening requires security/privacy/domain/reliability approval appropriate to severity.

## Notification decision

Counsel/privacy creates a per-jurisdiction/contract record:

```text
applicable entity/role/law/contract
trigger threshold and known facts
regulator/merchant/person/provider recipient
deadline start/event and time remaining
required content/method/language
delay/prohibition/exception authority
approver, sent evidence and updates/correction
```

Notifications are accurate, understandable and actionable: what happened, data/risk, what Ubi did, what recipient should do, contact/complaint path and updates. Do not overstate certainty or hide known material facts.

For merchant-controlled booking data, Ubi promptly provides the merchant the facts/evidence/support required by the approved DPA and law; roles do not excuse delay or duplicate contradictory notices.

## Researcher/vulnerability intake

- Acknowledge safely without demanding sensitive proof in public channels.
- Give secure submission/reference and scope guidance.
- Reproduce without retaliating against good-faith research; legal/abuse cases are reviewed individually.
- Coordinate fix/disclosure and credit where policy permits.
- Never authorize production data access casually.

## Post-incident requirements

- complete affected-record/action inventory and corrections;
- fulfill notification, rights, support, monitoring and regulator/merchant follow-up;
- update threat/privacy/data maps, risk/vendor/retention and tests;
- review why prevention/detection/JIT/RLS/telemetry/response failed;
- validate credential/data copies removed and temporary containment retired;
- executive accepts residual risk; critical gap blocks launch/release.

## Exercise gates

Exercise at least account takeover, cross-tenant response, leaked provider token, malicious release and vendor breach scenarios. Verify counsel/contact path, evidence chain, credential rotation, merchant communication and restored negative isolation tests.


# Fraud, abuse, moderation and safety

Status: **Draft — control and operating model; launch-country policy, staffing and thresholds unresolved**  
Depends on: actors, payments, conversations/calls, permissions, privacy, threat model and operations console  
Blocks: public discovery, messaging/calls, automated restrictions, payment operations and production pilot

## Purpose

Ubi can be abused by fake businesses, fraudulent customers, compromised accounts, payment screenshot scams, refund abuse, spam, harassment, unsafe content, false reports and malicious automation. This document defines prevention, detection, proportionate action, human review and appeal without treating ordinary East African connectivity/device patterns as guilt.

## Safety principles

- Protect people and platform integrity while minimizing unnecessary identity collection.
- A risk signal is not proof. Explainable combinations, domain evidence and human review govern high-impact restrictions.
- Shared phones, recycled numbers, SIM changes, carrier-grade NAT, cybercafé/shared Wi-Fi, low-cost devices, cash and repeated offline retries are normal possibilities—not automatic fraud.
- Separate merchant/customer misconduct from provider outage, app defect and misunderstanding.
- Friction is proportional: low-risk users should not face government-ID demands for an ordinary booking.
- Never let revenue or merchant status override serious safety evidence.
- Protect reporters from retaliation and accused parties from malicious/duplicate reporting.
- Automation can triage; irreversible/high-impact action has defined human authority and appeal.
- Evidence is purpose-limited, access-controlled and retained only as justified.
- Ubi is initially merchant infrastructure, not an open marketplace/social network; discovery/calls expand only after safety capability exists.

## Abuse surfaces

### Account and organization

- automated account/OTP creation, credential stuffing and takeover;
- fake/impersonated business, stolen ownership claim or hostile staff invitation;
- compromised former staff/owner, shared credentials and recovery social engineering;
- malicious subdomain/business name/logo/link impersonation;
- merchant using Ubi for prohibited/unlicensed/deceptive services.

### Booking and customer

- fake bookings, slot hoarding, no-show campaigns or competitor sabotage;
- abusive cancellations/reschedules and repeated scoped-link guessing/sharing;
- customer/staff manipulating price, availability, completion or dispute evidence;
- collection of unnecessary sensitive data through custom forms/notes;
- unsafe field-service address or event circumstances.

### Payments

- fake screenshot/reference, duplicate/reversed transfer or wrong-recipient claim;
- stolen payment method/provider account;
- over-refund, refund to different route, collusion and chargeback abuse;
- merchant claiming payment/Ubi custody inaccurately;
- laundering/regulated money movement if Ubi accidentally controls funds.

### Messaging, content and calls

- spam/unsolicited marketing, phishing, malicious links/files and impersonation;
- harassment, threats, hate, sexual exploitation, grooming and extortion;
- staff/customer stalking through contact/location/schedule information;
- automated misinformation, fake urgency or deceptive human impersonation;
- report/block evasion, repeated calls and caller-ID abuse;
- illegal/copyright/privacy-violating media, especially child/event images.

### Platform and operations

- scraping, enumeration, denial/cost exhaustion and referral/promotion gaming;
- false reports intended to suspend a competitor/customer;
- Ubi support/merchant moderators abusing access or evidence;
- compromised provider/automation/vendor producing harmful actions;
- retaliation against reporter or leaking case outcome/evidence.

## Trust and verification levels

Verification is attribute/purpose specific, not one universal “verified” badge.

| Level | Meaning | Does not mean |
|---|---|---|
| Contact controlled | Recent successful phone/email/provider challenge | Real/legal identity, ownership or trustworthiness |
| Account secured | Authenticator/session/device controls meet baseline | Business authority |
| Organization claimed | Account asserted business relationship | Legal business verified |
| Representative verified | Evidence links person to organization under approved process | Every service/license/claim is valid |
| Provider asset verified | Official provider confirms admin/control of exact asset | Provider endorses Ubi/merchant |
| Sector credential verified | Specific credential checked with source/expiry | General service-quality guarantee |
| Payment evidence verified | Authenticated provider/manual approved evidence for exact transaction | Settlement finality beyond provider semantics |

UI states exact verification and date/source; no generic green check that customers can misunderstand.

## Merchant onboarding controls

Initial assisted pilot:

1. invitation/referral/source and contact control;
2. organization name, representative authority and official/business evidence proportionate to country/provider/vertical;
3. exact branch/operation/service and public claims review;
4. payment/provider account belongs to approved party/model;
5. prohibited/restricted services and licensing screen;
6. owner/admin security setup and recovery contacts;
7. terms, privacy/data-role and customer-dispute responsibility acceptance;
8. test booking/refund/message/support scenarios;
9. limited pilot activation, monitoring and escalation route.

Verification documents are high-risk data: collect only when needed, segregate, restrict support access and delete under schedule. Ubi does not invent a “business license verification” claim if it only reviewed an uploaded image.

## Customer friction model

Public browsing remains open with rate/abuse controls. Booking may use guest/contact verification. Escalate friction only for specific risk/action:

- resend/rate delay or challenge for automated abuse;
- contact re-verification for suspicious change or risky scoped-link action;
- deposit/provider verification under merchant policy;
- stronger account/representative review for repeated harmful behavior;
- block/restrict only affected merchant/action where possible.

Do not require national ID merely because a user lacks a modern phone, changes SIM or shares a household number.

## Signals and risk decisions

Allowed signal families, subject to privacy/legal review:

- authenticated account/session/device/recovery changes;
- booking velocity, hold expiry, no-show/cancellation/dispute patterns with denominator/context;
- idempotency/replay/automation/resource consumption anomalies;
- verified payment reversals/disputes/reference reuse;
- provider account/webhook/security/reputation state;
- content/link/file safety classifier plus human report;
- merchant verification/licensing/complaint history;
- confirmed relation between accounts/devices/payment sources only where lawfully and accurately established.

Weak signals requiring caution:

- IP address/location mismatch, because mobile networks/VPNs/shared access are common;
- cheap/old/rooted device or failed integrity signal;
- language, name, neighborhood, disability, typing speed or literacy proxy;
- contact shared/recycled/SIM-changed;
- offline retries, delayed events or cash/manual payment;
- one uncorroborated report.

Restricted/protected traits and their proxies are not used for risk scoring without exceptional lawful necessity/fairness review.

Every automated risk decision records model/rule version, material signals, confidence, action, expiry, evidence and review/appeal route. Avoid opaque vendor scores as sole high-impact decision.

## Proportionate action ladder

1. Observe/log minimized signal.
2. Give understandable warning/help.
3. Add rate/cooldown/size limit.
4. Require contact re-verification or step-up.
5. Delay/hold the risky action while preserving unrelated access.
6. Require merchant/customer confirmation or second approval.
7. Route to human case review.
8. Restrict exact feature, contact, conversation, payment or operation temporarily.
9. Suspend account/organization for severe/current risk.
10. Terminate/remove/report only under approved policy/law/provider duty.

Each action has scope, reason category, owner, start, review/expiry, affected workflows, notification rule and appeal/redress. Indefinite “under review” is prohibited.

## Booking and no-show abuse

- Short, capacity-aware holds with expiry and per-principal/contact/merchant risk limits.
- Idempotency prevents retries from looking like distinct bookings.
- Merchant-configured deposits/cancellation policy must be visible before confirmation and versioned.
- Repeated no-show/cancel signals are merchant-contextual; one merchant cannot silently blacklist a customer across all businesses.
- Merchant may block a customer for its organization with reason/expiry/review; global Ubi restriction requires separate platform evidence.
- Waitlist/hold priority cannot be purchased/manipulated invisibly.
- Staff override records reason/approval and does not hide double booking.

## Payment fraud and dispute safety

- Screenshot/message/reference is reported evidence only.
- Verify exact provider account, amount, currency, reference, payer where available, merchant and transaction state.
- Reused reference and mismatched recipient/amount create a case; never mark verified by text similarity.
- Refund amount/source/destination follows verified allocation and approved policy; change of destination is high risk.
- Provider timeout/unknown status triggers reconcile, not duplicate charge/refund.
- Manual verification/refund needs capability, step-up, limits and second approval at defined threshold.
- Customer and merchant can dispute/report with evidence; Ubi states whether it facilitates communication or bears a contractual decision role.
- Ubi avoids unapproved fund custody/money transmission and uses licensed providers under document 52.

## Messaging and call safety controls

Before native messaging/calls:

- clear business/customer/automation identity and booking/business context;
- customer-initiation/contact rules and consent/transactional purpose;
- block, report, mute and stop-call controls visible during/after interaction;
- rate/frequency/quiet-hours and repeated-unanswered-call controls;
- link/file/media scanning and safe previews/downloads;
- staff takeover and automation pause;
- conversation-level tenant and participant authorization;
- abuse category, imminent-danger escalation and evidence preservation policy;
- moderator/support tools with masking/JIT/audit;
- minors/sensitive vertical/recording restrictions;
- device/network failure does not repeatedly call/message;
- provider/platform complaint/quality limits and opt-out sync.

Blocking stops future direct contact through Ubi for the scoped relationship and prevents automation bypass. It does not promise to block contact through independent phone/social channels Ubi does not control.

Call recording remains off until separate country/consent/storage/access/retention approval.

## Content and service policy structure

Final policy is country/vertical/provider specific. Categories requiring prohibition or restriction review include:

- sexual exploitation/abuse, especially children;
- credible threats, violence, harassment, hate and targeted privacy abuse;
- scams, phishing, malware, impersonation and deceptive business claims;
- illegal goods/services or regulated activity without required authorization;
- non-consensual intimate/private imagery and unlawful surveillance;
- doxxing/identity/payment credentials;
- copyright/trademark/privacy-violating media;
- medical/legal/financial claims beyond approved first-release scope;
- spam/unsolicited bulk messaging and provider-policy violations;
- unsafe physical-service/event circumstances.

Policy gives definitions, examples, exceptions/public-interest process where relevant, enforcement range, evidence, notice and appeal. It is written in usable launch languages; moderators receive contextual training.

## Report, block and emergency flow

1. Immediate block/mute/leave/stop-call works without writing a long report.
2. Optional report captures category, relationship/message/call references and safe description/evidence.
3. Ubi acknowledges without promising an outcome it cannot deliver.
4. Triage prioritizes imminent physical/child/financial/account compromise risk.
5. Preserve exact minimum evidence under case/legal-hold rules; restrict access.
6. Apply temporary scoped protection if necessary.
7. Investigator reviews content/context/history/provider evidence and conflicts.
8. Take proportionate action and notify parties without exposing reporter/private evidence.
9. Offer appeal/redress and monitor evasion/retaliation.
10. Escalate emergency/law enforcement only under approved local protocol/document 57.

Ubi is not an emergency service. UI must provide locally verified emergency/support resources only where accurate and maintainable.

## Moderation operations

- Queues separated by severity, language, country, data sensitivity and conflict of interest.
- Reviewer access is case-bound, minimum, time-limited and audited.
- Two-person review for permanent/high-impact action where feasible.
- No moderator downloads/copies evidence outside controlled tools.
- Decision reasons use controlled categories plus bounded notes.
- Quality sampling measures consistency, false positives/negatives, timeliness and user harm.
- Staff wellbeing, exposure minimization and escalation support are part of the operating design.
- Vendor moderation requires subprocessor, transfer, security, training, audit and deletion approval.

Pilot scope must match available trained coverage. If Ubi cannot staff safety in a language/time window, it cannot safely launch open messaging/calls there.

## Appeals and redress

- Notice states affected feature/account/organization, general reason, duration and appeal route unless prohibited by safety/law.
- Appeal is accessible without relying on the suspended session and supports low-literacy assistance.
- Different qualified reviewer handles serious appeal where possible.
- Review considers mistaken identity, shared/recycled contact/device/network and new evidence.
- Restore access/data/configuration safely after reversal and prevent automated immediate re-restriction by the same stale signal.
- Track outcome and policy/model correction.
- Merchant/customer dispute is not always a Ubi policy violation; route contractual dispute separately.

## Abuse of reporting and merchant power

- Rate/deduplicate coordinated reports without suppressing genuine repeated harm.
- One merchant cannot create a secret platform-wide customer blacklist.
- Staff cannot label customer fraud to avoid refund/complaint review without evidence.
- Reporters/subjects do not see private internal/provider/security information.
- Conflicts involving Ubi staff/major merchants escalate independently.
- Public review/rating marketplace is deferred; moderation design precedes it.

## Privacy, retention and sharing

- Risk data has explicit purpose/classification/retention and access.
- Do not retain full conversations/files indefinitely for possible abuse.
- Hash/device/network/linkage data still may be personal/sensitive and is governed.
- Cross-merchant linking is exceptional and platform-safety justified, not commercial profiling.
- Vendor/model training on reports/content is off unless separately approved.
- Share with provider/regulator/law enforcement only under contract/law/policy and document 57.
- Delete unsupported reports/signals and review stale restrictions.

## Metrics and guardrails

Measure by country/language/vertical/provider with minimum-group privacy:

- report/block/call-stop rate and severity;
- time to protective action and final decision;
- upheld/reversed appeal rate and reason;
- repeat/evasion/retaliation rate;
- false positive/negative findings from quality review;
- verification completion/abandonment/help burden;
- payment discrepancy/loss/refund abuse;
- spam/provider complaint/account restriction;
- moderator queue age/exposure and access violations;
- impact on legitimate bookings, especially shared phone/network/low-cost device cohorts.

Do not reward teams for high suspension volume. Safety, fairness, recovery and legitimate task success are joint guardrails.

## Release gates

- First-release prohibited/restricted service/content policy and merchant verification are locally reviewed.
- Report/block/appeal and emergency/escalation processes are usable in launch languages.
- Payment reported/verified/refund/dispute controls and licensed-provider boundary pass adversarial tests.
- Risk system documents weak signals and cannot auto-ban solely on device/IP/SIM/cash/offline patterns.
- Staff/moderator/support access is case/JIT/audited with training and capacity.
- Serious actions have expiry/review, notice and independent appeal.
- Messaging/calls do not launch before block/report/rate/takeover/minor/recording controls.
- Automation cannot independently perform high-impact restriction, payment, ownership or booking decisions.
- Legal/privacy/provider reporting, evidence and retention rules are approved.

## Open decisions

- First vertical/country prohibited-service and merchant-verification requirements.
- Pilot safety staffing, languages, coverage hours and escalation partners.
- Risk rules/vendor versus simple reviewed first-party signals.
- Deposit/no-show/block policies and cross-merchant platform action threshold.
- Emergency resource display and law-enforcement escalation by country.
- Native messaging/call launch and recording remain future gated.
- Public discovery/reviews/marketplace moderation remains deferred.

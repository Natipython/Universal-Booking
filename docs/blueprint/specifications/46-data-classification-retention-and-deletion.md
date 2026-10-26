# Data classification, retention and deletion

Status: **Draft — governance framework proposed; jurisdiction-specific periods and legal bases unresolved**  
Depends on: data catalogue, isolation/storage architecture, workflows and launch-market decision  
Blocks: privacy notices/agreements, executable retention rules, deletion workers, support procedures and production data

## Purpose and legal caution

This document defines how Ubi classifies, minimizes, retains, exports, corrects and deletes data across databases, devices, objects, search, cache, integrations, logs and backups.

It deliberately does **not** declare one retention period or legal basis for all East African countries. Applicable privacy, consumer, tax, communications, employment, payment and record-keeping obligations depend on the launch jurisdiction, Ubi’s local entity/role, the merchant’s role, the exact feature and provider contracts. Country-qualified counsel must approve the executable schedule and documents before launch.

## Non-negotiable principles

- Collect only data needed for a named product, security, contractual or legally reviewed purpose.
- A field is not justified because it may be useful later.
- Classification follows the most sensitive meaningful content, not the table name.
- Retention has a start trigger, duration/condition, disposal action, owner and exception—not “keep forever.”
- Soft deletion is lifecycle state, not proof of erasure.
- Legal hold suspends only the covered disposal, not all privacy/security controls.
- Deletion propagates to projections, objects, devices and processors and produces evidence without retaining deleted content.
- Backups have bounded expiry and restore-time deletion replay; instant selective removal from immutable backups may be impossible.
- Merchant/customer closure never silently erases records still required for bookings, disputes, security or legal obligations.
- Ubi never claims that it can delete copies held independently by WhatsApp, Telegram, Meta, TikTok, telecom/payment providers or message recipients; responsibility and user guidance must be accurate.

## Classification levels

| Level | Name | Examples | Baseline handling |
|---|---|---|---|
| C0 | Approved public | Published business name, public service, branch hours, approved public image | Integrity/abuse controls; only explicit publication; cache/CDN allowed |
| C1 | Internal operational | Non-sensitive feature/config IDs, minimized aggregate metrics, runbook state | Authenticated workforce/service access; no public exposure |
| C2 | Confidential business/customer | Booking details, customer contacts, staff schedules, private business configuration, ordinary messages | Tenant/purpose isolation, least privilege, encryption in transit/at rest, masked logs/exports |
| C3 | Restricted/high impact | Authenticators/tokens, identity evidence, precise sensitive location, private attachments, payment credential/evidence, security events, high-risk notes, key material | Narrow services/fields, step-up/JIT where appropriate, stronger encryption/key separation, access audit, minimal retention |

Classification is independent of “personal data.” A confidential merchant price rule may not identify a person; a public profile may still include personal contact information whose publication needs authority and notice.

### Prohibited by default

- plaintext passwords, OTPs, recovery codes, private keys or provider tokens in ordinary tables/logs;
- full card/mobile-money credentials unless a separately approved compliant payment design requires them;
- arbitrary government/health/biometric information for normal booking;
- secret questions based on discoverable personal facts;
- unrestricted free-text internal notes about protected/sensitive traits;
- full social-message content copied indefinitely “for AI training”;
- production data in developer laptops, source repositories, fixtures or casual spreadsheets.

## Data inventory contract

Every field/object/event/log/derived document has an inventory entry containing:

```text
stable inventory ID and schema/event/object path
business/domain owner and data steward
subjects and tenant/global scope
classification and sensitivity tags
purpose(s) and prohibited secondary uses
collection/source and recipients/processors
required/optional/derived status
countries/regions and transfer path
authoritative owner and copies/projections
retention trigger/rule and deletion mechanism
access roles/capabilities and audit requirement
encryption/masking/tokenization requirement
notice/consent/contract reference where applicable
quality/lineage/schema version
```

CI/schema review should reject an implementation field without a registered purpose/classification/retention owner.

## Purpose model and responsibility

Processing purposes use controlled IDs, for example:

- authenticate and secure an account;
- administer merchant membership;
- create/fulfil/support a booking or custom request;
- communicate a requested reminder/update;
- verify/reconcile a payment claim;
- operate/fix an integration;
- prevent fraud/abuse and investigate incidents;
- meet approved record-keeping/dispute duties;
- improve product reliability through minimized telemetry;
- optional merchant marketing under its own rules.

Purposes are not interchangeable. Contact provided to complete a wedding booking does not automatically authorize unrelated Ubi or merchant advertising, model training or sale/sharing.

Ubi’s legal/data role may differ by feature and country—controller, joint controller, processor/service provider, independent controller or another local concept. The product architecture records the responsibility matrix but local counsel approves the terminology and contractual allocation.

## Collection and minimization

### Progressive collection

Ask only what is necessary at the current step:

- discovery can show public services without login;
- availability search may need service/date/location but not full identity;
- a hold needs the minimum contact/anti-abuse facts;
- confirmation collects required party/fulfilment/payment facts;
- optional profile reuse is separate from completing one booking;
- identity evidence is requested only for a documented high-risk use case.

Below-average digital literacy is not permission to hide collection in long text. Use short layered notices in the user’s language, clear required/optional labels, confirmation summaries and accessible full policy links.

### Free text and attachments

Free text is high-risk because users/staff may enter unnecessary data. Use structured choices where practical, prompt users not to include secrets, classify visibility, limit size, scan attachments and provide controlled correction/removal. Do not parse/copy every message into every downstream system.

## Retention rule structure

An executable rule contains:

```text
rule/version/jurisdiction/data class and exact object scope
purpose/legal/contractual justification reference
retention trigger event
active operational period
archive/restricted period if any
terminal action: delete, anonymize, aggregate, return or review
legal-hold/incident/dispute exception
downstream systems/processors and propagation objective
backup expiry/restore replay behavior
owner, approver, effective/review dates
```

Do not calculate expiry only from row creation. Useful triggers include account closure, booking completion/cancellation, payment settlement/dispute closure, provider disconnection, upload abandonment, invite/grant expiry, incident closure and contract termination.

### Data-class schedule requiring country review

| Record family | Normal trigger | Retention decision inputs | Terminal action direction |
|---|---|---|---|
| Account/auth/session | Session/recovery expiry or account closure | Security, account recovery, incident evidence | Revoke immediately; delete/minimize bounded metadata after reviewed period |
| Membership/authority | Membership/organization closure | Ownership disputes, audit, employment/contract context | Remove live authority immediately; retain minimal authority history as justified |
| Customer CRM | Last relationship/merchant instruction | Open bookings, disputes, merchant instructions, subject rights | Delete/anonymize non-required profile; preserve necessary transaction snapshot |
| Booking/request/quote | Completion/cancellation/expiry | Consumer/tax/contract/dispute obligations by market | Restrict then delete/anonymize where allowed; preserve exact required evidence only |
| Conversation/message | Closure/last activity | Booking evidence, provider policy, safety reports, merchant responsibility | Shortest justified content retention; structured commercial snapshots may outlive chat |
| Payment/reconciliation | Settlement/dispute closure | Payment/tax/AML/provider/local obligations | Retain exact necessary evidence; avoid credentials/full unnecessary payload |
| Media/upload | Abandonment, unlink, booking closure | Evidence need, merchant/customer expectation, legal hold | Delete original/variants/CDN; keep minimal tombstone/checksum only if justified |
| Scoped links/invites | Expiry/redemption/revocation | Abuse/security evidence | Delete secret/verifier promptly; retain bounded minimal event evidence |
| Provider tokens/events | Disconnection/event completion | Reconciliation/debug/security and provider contract | Revoke upstream; delete credentials; minimize/delete raw payload earlier than canonical facts |
| Logs/traces | Collection time/incident closure | Reliability/security need and content risk | Short bounded hot retention; aggregate/minimize; delete raw telemetry |
| Audit/security | Event/incident/account/contract trigger | Integrity, investigation, legal duty | Protected bounded retention with access audit; never “forever just in case” |
| Analytics | Event/experiment end | Metric purpose, consent/notice, identifiability | Aggregate/de-identify early; delete subject-level event on schedule |
| Backups | Backup creation | Approved RPO/RTO, investigation hold, key availability | Automated expiry/destruction; deletion replay after restore |

Exact durations remain `TBD-JURISDICTION` until a signed register exists. Code must not embed guessed universal years.

## Lifecycle states versus deletion

Use precise words:

- **Active**: used for current approved purpose.
- **Archived/restricted**: retained but excluded from ordinary work and narrowly accessible.
- **De-identified/aggregated**: direct identifiers removed and re-identification risk assessed; not automatically anonymous.
- **Anonymized**: irreversibly no longer attributable under the applicable standard; requires evidence, not a label.
- **Tombstoned**: minimal marker prevents resurrection/replay while content is removed.
- **Deleted**: removed from active systems and scheduled/propagated through all governed copies.
- **Expired backup copy**: disappeared through documented backup lifecycle; before then protected and replay-controlled.

Closing, hiding, archiving, encrypting or removing from search alone is not deletion.

## Deletion orchestration

Deletion is an idempotent multi-system process managed through `privacy.privacy_cases`, `privacy_case_tasks`, retention work and `deletion_tombstones`.

1. Resolve verified subject/tenant/resource and responsibility.
2. Freeze accidental mutation/merge as needed.
3. Discover authoritative rows and copies from inventory/lineage.
4. Check active booking/payment/dispute/security duties and legal holds.
5. Decide delete, correct, restrict, anonymize, retain or partially fulfill with reason.
6. Execute owner-domain commands; never a blind graph cascade.
7. Delete/rewrite search, cache, object variants, device sync projections and provider-held copies where Ubi controls/instructs them.
8. Issue tombstones so delayed events/offline clients do not resurrect removed content.
9. Produce a minimized manifest of systems/actions/results/errors without copying deleted content.
10. Retry/escalate failures and communicate the accurate outcome/limitations.
11. Ensure any later backup restore replays tombstones/deletion ledger before service resumes.

### Relationship-specific handling

- Deleting an Account does not automatically delete each merchant’s independently required Customer record.
- Deleting a Customer profile does not rewrite immutable booking/payment facts beyond approved minimization/anonymization.
- Merchant organization closure triggers export/contract/retention workflow, not `DELETE organization CASCADE`.
- Customer-account claims can be severed while merchant transactional records remain.
- Merged customer records retain reversible/protected merge evidence until its own retention ends.
- Public profile unpublication is fast but is followed by cache/search/object invalidation; third-party search-engine copies are handled through available processes without impossible guarantees.

## Correction and provenance

Correct current factual data without rewriting what was actually accepted/issued/paid at an earlier time.

- Current customer name/contact can change with provenance.
- A booking snapshot retains the value shown/accepted then; an amendment/correction event explains later change.
- Provider evidence is preserved as received plus canonical interpretation/revision.
- Disputed staff notes can be restricted/corrected with authorship/history rather than silently overwritten.
- Derived search/analytics data is rebuilt from corrected owner state subject to legitimate historical metric definitions.

## Legal holds and investigations

A legal/incident hold requires exact subject/data scope, authority/request source, purpose, issuer, start, review/expiry, access restriction and release evidence.

- Holds do not grant content access.
- Scope is the minimum defensible set, not “entire database.”
- Conflicting deletion tasks pause only affected objects and record the reason.
- Hold inventory is reviewed regularly; expired/unsupported holds release and deferred disposal resumes.
- Merchant requests cannot force Ubi to retain prohibited/unnecessary data indefinitely.

## Subject/merchant request operations

The privacy-case workflow supports access, portability/export, correction, deletion, restriction, objection/withdrawal and notice/consent evidence as applicable.

Controls include:

- identity/authority verification proportional to disclosure risk;
- safe handling for a customer who has no Ubi Account;
- merchant-versus-Ubi responsibility routing;
- no disclosure of another person/business/provider secrets;
- structured machine-readable/common format where required/useful;
- secure expiring export delivery;
- statutory/contractual deadline calendar configured by jurisdiction;
- decision, exception and appeal/escalation evidence;
- subprocessors/provider task tracking.

Support cannot deny or fulfill a request from intuition; it follows an approved jurisdiction/playbook and escalates uncertain cases.

## Backup, device and third-party behavior

### Backups

Active deletion completes in production stores while immutable backups expire on a bounded schedule. Backups are isolated/encrypted/access-controlled and not restored for ordinary lookup. Any restore runs deletion/revocation replay and reconciliation before customer traffic.

### Android devices

Deletion/revocation emits sync tombstones and purge commands. The app removes affected rows/files from its accessible store and prevents pending commands from recreating them. Ubi cannot guarantee erasure from a rooted device, screenshots, manual exports or external OS backups outside its control; settings and notices reduce those paths.

### Providers and recipients

For social, messaging, calling and payments, Ubi records which copy it controls, which it can request the provider to remove and which remains under the business/customer/provider. Provider contract/API limitations appear in the data inventory and user-facing explanation.

## Logs, analytics and AI/automation

- Logs contain identifiers/status/error codes, not message bodies/tokens/payment details by default.
- Debug payload capture is exceptional, sampled, redacted, time-bounded and access-controlled.
- Analytics events use controlled properties and reject free-text/high-risk fields.
- Automation extraction stores source references and structured suggestions; raw content retention follows the conversation purpose, not model convenience.
- Training/fine-tuning/model improvement with merchant/customer content is a separate product/legal decision, default off until approved, transparent and controllable as required.
- Vendor prompts/responses, abuse logs and model telemetry are inventoried as subprocessors/copies with their own retention and deletion capabilities.

## Governance and change control

Data owner approves purpose and operational need. Privacy/legal approves applicable basis/notice/rights/retention. Security approves classification/control. Engineering/data stewardship approves implementation/lineage/deletion proof.

New data or secondary use requires:

1. inventory/purpose proposal;
2. minimization/alternative review;
3. classification and threat/privacy impact review;
4. jurisdiction/contract/vendor assessment;
5. notice/choice/consent changes where required;
6. retention/deletion/export implementation;
7. tests/metrics/access review;
8. approval before production collection.

## Release gates

- Launch country, Ubi legal entity and merchant/customer responsibility model are decided.
- Local counsel approves the legal/regulatory matrix, processing purposes, notice/agreements and exact retention schedule.
- Every first-release field, object, event and log has inventory/classification/purpose/retention/owner.
- Prohibited data rules are represented in schemas, UI copy, logs and tests.
- Privacy-case identity verification, deadlines, responsibility routing and secure export are rehearsed.
- Delete/correct/restrict workflows pass end-to-end across PostgreSQL, cache, search, object, device sync, analytics and chosen providers.
- Legal hold prevents only covered disposal and access remains restricted.
- Backup restore replays deletion/revocation before reopening.
- Vendor/subprocessor deletion and data-location behavior is contractually/technically verified.
- Staff training and audited JIT access exist before real support data.

## Open decisions

- Launch jurisdiction and applicable sector/payment/communications rules.
- Ubi/merchant legal roles per workflow and contract.
- Exact retention trigger/duration/action for every first-release inventory item.
- Whether any health/identity/sensitive-category booking fields are excluded entirely from pilot.
- Customer request path when no Account exists and merchant cannot be reached.
- Backup expiry and restore deletion-replay implementation.
- Provider deletion/export limits and contractual commitments.
- Definition/testing standard for de-identification/anonymization.
- AI/automation provider data-use and zero/short-retention requirements.

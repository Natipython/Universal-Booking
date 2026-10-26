# Logical schema and table catalogue

Status: **Draft — table inventory, not executable DDL**  
Depends on: conceptual model and operational database architecture  
Blocks: columns/types, keys/constraints/indexes, RLS, retention, migrations and implementation

## Purpose

This catalogue names the logical PostgreSQL schemas/tables, their write owner and role. It prevents missing lifecycle/evidence/support data and prevents developers from creating duplicate concepts under different names.

## Release legend

- **F** — foundation required by dependent architecture.
- **P** — likely pilot capability, subject to first vertical/release decision.
- **E** — expansion.
- **G** — separately future gated.

A listed table is not authorization to implement it in Version 1.

## Common logical columns

Unless the table’s nature differs, owner tables use:

```text
id                      immutable opaque identifier
organization_id         immutable tenant key for tenant-owned rows
version                 optimistic/domain version where mutable
created_at              server persistence time
created_by_*            attributable principal/source where required
updated_at              current-row update time where meaningful
effective_from/until    version/rule applicability where meaningful
lifecycle/status        context-specific constrained state
```

Append-only evidence rows use `recorded_at`, `occurred_at`, actor/service/correlation/causation/source and do not pretend to be mutable current records. Every foreign key, lifecycle, deletion and audit field is finalized in later documents.

## `iam` — identity and sessions

| Table | Stage | Owner/purpose | Important relationships |
|---|:---:|---|---|
| `accounts` | F | Ubi authentication identity/lifecycle, not tenant/member/customer | Optional person link; referenced by sessions/memberships/claims |
| `persons` | E | Separately justified real-human identity record | Optional one-to-one/many evidence relationship with account; not auto-created for every customer |
| `account_contacts` | F | Normalized/display contact method and lifecycle | Account; verification episodes |
| `contact_verifications` | F | Challenge/point-in-time control evidence, purpose and result | Account contact; provider reference; no OTP secret |
| `authenticators` | F | Credential/passkey/password/approved authenticator metadata and lifecycle | Account; protected verifier/secret reference outside ordinary output |
| `sessions` | F | Login session/refresh family/assurance/revocation | Account, installation; security events |
| `session_events` | P | Created/refreshed/stepped-up/revoked/security history | Session/account |
| `device_installations` | F | App-install/platform/push/sync/security association | Account sessions, push endpoints, sync scopes |
| `account_recovery_cases` | F | Recovery request/risk/evidence/decision; cannot grant membership itself | Account/contact/session; operations case |
| `service_identities` | F | Workload/service principal metadata and environment scope | Runtime roles/audit |

## `org` — organization, hierarchy and authority

| Table | Stage | Owner/purpose | Important relationships |
|---|:---:|---|---|
| `organizations` | F | Primary tenant/lifecycle/default country/timezone/currency/language/calendar | Brands, branches, operations, memberships |
| `organization_verifications` | F | Business/representative verification evidence and lifecycle | Organization, account/person, media/evidence reference |
| `ownership_claims` | F | Account/legal representative authority and transfer history | Organization, account, verification/approval |
| `brands` | F | Public commercial identity under organization | Organization, public profile, operations |
| `branches` | F | Physical/operational scope/location/timezone | Organization; operations, membership scopes |
| `booking_operations` | F | Booking configuration/workflow boundary | Organization; branch/brand links; catalog/resources/bookings |
| `operation_brand_links` | F | Explicit brand-operation presentation relationship | Same-tenant operation/brand |
| `operation_branch_links` | F | Explicit branch-operation relationship/effective history | Same-tenant operation/branch |
| `memberships` | F | Account-to-organization authority lifecycle/version | Account, organization; roles/scopes/grants |
| `membership_scopes` | F | Branch/operation/resource/queue/assigned boundary | Membership and exactly typed scope target |
| `roles` | F | Tenant/platform role bundle identity/version | Organization/system template; role capabilities |
| `capabilities` | F | Controlled action registry | Role/grant/policy/test catalogue |
| `role_capabilities` | F | Capability included in a role bundle | Role, capability; optional conditions |
| `role_assignments` | F | Role assigned to membership with scope/effective range | Membership, role, scope |
| `capability_grants` | F | Explicit grant/deny/limit beyond role bundle | Membership, capability, scope, delegator |
| `membership_invites` | P | Expiring single-purpose invitation and acceptance | Organization, inviter, contact/account, proposed role/scope, access grant |
| `approval_requests` | F | Exact privileged action/before-after/expiry/current status | Organization, requester, resource/action |
| `approval_decisions` | F | Approver allow/refuse evidence | Approval request, membership/account |
| `authority_versions` | F | Optional tenant/member policy epoch for cache/sync invalidation | Organization/membership |
| `organization_restrictions` | F | Proportionate restriction/suspension/reason/redress | Organization; operations case/audit |

## `public_presence` — public publication and URLs

| Table | Stage | Owner/purpose | Important relationships |
|---|:---:|---|---|
| `public_profiles` | P | Published safe organization/brand/branch/operation representation | Exact owner/version/publication state |
| `public_slugs` | P | Normalized slug/subdomain target/lifecycle/reservation | Public profile/organization/brand |
| `slug_redirects` | E | Rename redirect/source/target/effective expiry | Old/new slug; abuse/suspension policy |
| `reserved_slugs` | F | Platform/protected/abusive/confusable names | Normalized slug/reason |
| `public_service_projections` | P | Rebuildable safe service/version presentation | Catalog publication/profile |
| `public_location_projections` | P | Safe public address/landmark/map presentation | Branch/location owner |

## `catalog` — services, forms, price and policy versions

| Table | Stage | Owner/purpose | Important relationships |
|---|:---:|---|---|
| `operation_config_versions` | F | Immutable archetype/config publication/draft/effective version | Booking operation, publication |
| `services` | F | Stable service identity/lifecycle | Booking operation; versions |
| `service_versions` | F | Immutable duration/fulfilment/confirmation/price-method configuration | Service, operation config, form/policy |
| `service_localizations` | P | Locale-specific service name/description/version | Service version, locale |
| `packages` | E | Stable package identity | Operation; versions |
| `package_versions` | E | Immutable inclusions/price/configuration | Package; service/add-on components |
| `add_ons` | E | Stable optional component identity | Operation; versions |
| `add_on_versions` | E | Immutable add-on configuration/price/capacity effect | Add-on |
| `service_add_on_rules` | E | Which add-on versions/options apply to service version | Same operation/config version |
| `service_resource_requirements` | F | Required resource roles/quantity/assignment/buffers | Service version, resource type |
| `price_components` | F | Exact fixed/rate/percentage/per-unit component definition | Service/package/add-on version, currency |
| `intake_forms` | P | Stable form identity/purpose | Operation/service |
| `intake_form_versions` | P | Immutable ordered form/conditional schema publication | Intake form |
| `intake_field_definitions` | P | Typed question/validation/classification/purpose | Form version |
| `intake_field_options` | P | Versioned localized allowed choices | Field definition |
| `policy_documents` | F | Stable policy identity/type/owner | Organization/operation/service |
| `policy_versions` | F | Immutable customer-visible terms/effective publication | Policy document; localized text/media reference |
| `catalog_publications` | F | Atomic set of configuration versions published/tested/retired | Operation config and included versions |

## `scheduling` — resources, rules, holds and commitments

| Table | Stage | Owner/purpose | Important relationships |
|---|:---:|---|---|
| `resource_types` | F | Operation/archetype resource classification | Booking operation |
| `resources` | F | Staff/space/asset/vehicle/equipment/pool identity/lifecycle | Organization/operation/branch/resource type |
| `resource_account_links` | P | Staff resource ↔ account/membership relationship | Same-tenant resource/membership; effective dates |
| `resource_eligibilities` | F | Resource may satisfy service requirement/role | Resource, service/resource requirement |
| `availability_rules` | F | Effective recurring local schedule/timezone/capacity | Resource/operation |
| `availability_rule_intervals` | F | Weekday/local start/end/break portions | Availability rule |
| `availability_exceptions` | P | Dated open/closed/block/capacity override | Resource/operation; reason/actor |
| `capacity_pools` | E | Pooled quantity/session/group inventory identity | Operation/resource/service |
| `capacity_buckets` | E | Bounded session/date/interval capacity/locked current summary | Capacity pool; ledger entries |
| `capacity_ledger_entries` | E | Append-only reserve/release/adjust evidence | Pool/bucket, hold/booking/reason |
| `holds` | F | Expiring purpose-bound capacity reservation aggregate | Operation/request/quotation/payment process |
| `hold_allocations` | F | Exact resource/pool interval/quantity held | Hold, resource/pool/bucket |
| `hold_transitions` | F | Placed/extended/consumed/released/expired history | Hold |
| `capacity_commitments` | F | Active/released booking allocation over resource/time/quantity | Booking, resource/pool; exclusion/quantity rules |
| `resource_reservations` | F | Shared physical collision boundary for exclusive held/committed intervals | Same-tenant resource plus exactly one hold or booking commitment owner |
| `capacity_commitment_history` | F | Commit/move/release/override evidence | Commitment, booking schedule version |
| `schedule_overrides` | P | Authorized conflict/after-hours/capacity override reason/approval | Resource/booking/approval/audit |

## `crm` — tenant customer relationships

| Table | Stage | Owner/purpose | Important relationships |
|---|:---:|---|---|
| `customer_records` | F | Tenant-private stable customer relationship/lifecycle | Organization; created from call/walk-in/web/channel |
| `customer_names` | F | Preferred/display/local-script name versions/provenance | Customer record |
| `customer_contacts` | F | Tenant-purpose contact methods/preferences/verification source | Customer record; not Account contact authority automatically |
| `customer_addresses` | E | Purpose-scoped local address/landmark/location | Customer record; protected location/media source |
| `customer_account_claims` | P | Verified link to Ubi account/booking-visible scope | Customer record, account, verification/grant |
| `customer_duplicate_candidates` | P | Similarity evidence/review state | Two same-tenant customer records |
| `customer_merges` | E | Surviving/source relationship and reversible audit | Same-tenant customer records |
| `customer_notes` | P | Purpose/visibility/classified operational note | Customer, author membership; no regulated free-for-all |
| `customer_communication_preferences` | P | Merchant-specific purpose/channel choice/evidence | Customer/contact, notice/version |

## `booking` — authoritative booking and fulfilment

| Table | Stage | Owner/purpose | Important relationships |
|---|:---:|---|---|
| `bookings` | F | Stable booking/current lifecycle/operation/archetype/source/next action | Organization/operation/customer references |
| `booking_parties` | F | Booker/payer/organizer/attendee/provider roles and visibility | Booking; customer/account/membership/external reference via typed links |
| `booking_sources` | F | Source/provenance instance such as call/walk-in/web/channel/staff | Booking; conversation/provider/client command reference |
| `booking_snapshots` | F | Immutable committed service/price/policy/location/intake evidence version | Booking; catalog/policy versions |
| `booking_snapshot_line_items` | F | Exact price/service/add-on/quantity/rate/currency components | Booking snapshot |
| `booking_snapshot_answers` | P | Typed customer answers against field/version/provenance | Booking snapshot, intake field |
| `booking_schedule_versions` | F | Immutable local/UTC/date/timezone/calendar/resource schedule history | Booking; current marker/version |
| `booking_state_transitions` | F | Append-only lifecycle transition/actor/reason/source | Booking |
| `booking_responsibilities` | P | Next actor/action/deadline/history | Booking/customer/business/system |
| `booking_cancellations` | P | Cancellation actor/responsibility/reason/policy/financial consequence | Booking/state transition |
| `fulfilment_events` | P | Check-in/start/complete/no-show/operational evidence | Booking, actor/resource |
| `booking_relations` | E | Approved parent/rebook/replacement/group relation | Same-tenant bookings; constrained relation type |

## `request` — custom request, quotation and amendment

| Table | Stage | Owner/purpose | Important relationships |
|---|:---:|---|---|
| `custom_requests` | E | Stable request/current lifecycle/operation/customer/next action | Organization/operation/customer/conversation |
| `request_state_transitions` | E | Append-only request lifecycle evidence | Custom request |
| `request_assignments` | E | Queue/staff/assessment ownership history | Request, membership/queue |
| `request_answers` | E | Typed requirement value/provenance/reviewer | Request, field definition/source message |
| `request_assessments` | E | Site visit/consultation/feasibility findings | Request; schedule/author/media |
| `quotations` | E | Stable quotation lineage/current lifecycle/version set | Request, issuer/operation |
| `quotation_versions` | E | Immutable issued scope/exclusions/validity/terms | Quotation, policy/price/version |
| `quotation_options` | E | Alternative option within exact quotation version | Quotation version |
| `quotation_line_items` | E | Exact item/rate/quantity/currency components | Quotation option/version |
| `quotation_acceptances` | E | Exact actor/token/auth/version/option acceptance/refusal | Quotation version/option/customer/grant |
| `request_booking_conversions` | E | Idempotent accepted quote → resulting booking outcome | Request/quote/booking/hold/payment process |
| `amendments` | E | Stable post-confirmation change proposal lineage | Booking/request |
| `amendment_versions` | E | Immutable scope/price/schedule/term proposal | Amendment |
| `amendment_acceptances` | E | Exact party/version acceptance | Amendment version/customer/account/grant |

## `conversation` — native/normalized messages and inbox

| Table | Stage | Owner/purpose | Important relationships |
|---|:---:|---|---|
| `conversations` | E | Tenant thread/current work/routing/automation summary | Channel connection/operation/queue/customer |
| `conversation_participants` | E | Customer/account/membership/external role in conversation | Typed constrained participant reference |
| `external_participants` | E | Provider-connection-scoped participant identity/profile evidence | Integration connection; optional reviewed customer link |
| `messages` | E | Immutable normalized/native message source/content/state/version | Conversation, author, provider reference |
| `message_revisions` | E | Provider/native edit/delete/tombstone history | Message |
| `message_delivery_evidence` | E | Supported submitted/sent/delivered/read/failure evidence | Message/provider event |
| `conversation_assignments` | E | Queue/staff claim/assign/release/escalation history | Conversation, membership/queue |
| `conversation_work_transitions` | E | New/open/waiting/resolved/closed/spam status history | Conversation |
| `automation_control_events` | E | Disabled/suggest/active/pause/takeover/resume history | Conversation/policy/actor |
| `automation_decisions` | E | Bounded rule/model/source facts/suggestion/send/escalation evidence | Message/conversation; no secrets/raw prompt dump |
| `extraction_suggestions` | E | Proposed structured field/source/confidence/reviewer/result | Message, request/booking draft field |
| `conversation_domain_links` | E | Explicit link to request/booking/quotation/customer | Same-tenant constrained link type |
| `conversation_blocks_reports` | E | Participant block/report/restriction/case evidence | Conversation/participant/operations |

## `access_grant` — secure scoped links

| Table | Stage | Owner/purpose | Important relationships |
|---|:---:|---|---|
| `access_grants` | F | Purpose/tenant/resource/action/expiry/use/revocation and protected token verifier | Booking/request/invite/customer action |
| `access_grant_redemptions` | F | Open/exchange/verification/action result evidence | Grant, anonymous/account/contact session |
| `access_grant_supersessions` | P | Old/new grant relationship/reason | Grant pair |
| `link_sessions` | P | Short-lived browser/app grant after URL-token exchange | Grant, session/device; no broader authority |

## `integration` — provider connections and raw/canonical transport

| Table | Stage | Owner/purpose | Important relationships |
|---|:---:|---|---|
| `providers` | F | Approved provider registry/category/owner/status | Adapter versions/capabilities |
| `provider_applications` | E | Ubi provider app/client/environment/review configuration metadata | Provider; secret references |
| `connections` | E | Tenant/provider business asset/account/scopes/lifecycle/routing | Organization/provider application |
| `connection_capabilities` | E | Versioned discovered/enabled provider feature/limit | Connection |
| `connection_credentials` | E | Secret-manager key reference/expiry/scope/rotation metadata | Connection; never plaintext credential |
| `connection_routes` | E | Provider asset to operation/queue mapping | Connection, operation/queue |
| `provider_event_receipts` | E | Authenticated raw/hash event/dedup/ack/process state | Provider app/connection; partition candidate |
| `provider_outbound_requests` | E | Stable request/reference/attempt/unknown/reconcile evidence | Connection, domain intent |
| `provider_health_checks` | E | Capability/token/webhook/latency validation results | Connection |
| `adapter_versions` | F | Supported provider schema/canonical mapping/release status | Provider/application/events |

Channel and payment-specific connection extensions live in their owner schema referencing `integration.connections`, not by adding every provider field here.

## `payment` — money evidence and reconciliation

| Table | Stage | Owner/purpose | Important relationships |
|---|:---:|---|---|
| `payment_connections` | E | Merchant payment product/account/KYB/currency/routing extension | Integration connection, organization/operation |
| `payment_obligations` | F | Exact amount/currency/purpose/deadline/condition | Booking/quotation/amendment |
| `payment_intents` | E | Logical attempt to satisfy obligation/current summary | Obligation, method/connection |
| `payment_attempts` | E | One customer/provider/manual interaction | Intent, provider request/reference |
| `reported_payment_evidence` | P | Customer/staff claim/reference/receipt/unverified status | Obligation/attempt/customer/media |
| `manual_payment_verifications` | P | Authorized merchant verification/rejection/source/reviewer | Reported evidence, approval/audit |
| `provider_transactions` | E | Authenticated provider transaction/current evidence summary | Payment connection/provider references |
| `provider_transaction_events` | E | Append-only provider state evidence | Provider transaction/event receipt |
| `payment_allocations` | E | Exact verified value applied to obligation | Provider/manual verified transaction, obligation |
| `refund_intents` | E | Exact controlled refund request/approval/current state | Transaction/allocation/booking |
| `refund_attempts` | E | Provider/manual refund interaction/evidence | Refund intent/provider request |
| `payment_reversals` | E | Provider reversal exact evidence/current handling | Transaction/allocation |
| `payment_disputes` | E | Dispute/chargeback lifecycle/deadline/evidence | Transaction/booking/operations case |
| `settlement_records` | E | Provider payout/settlement header/control totals | Payment connection/currency/period |
| `settlement_lines` | E | Transaction/refund/reversal/fee line evidence | Settlement/transaction/refund |
| `reconciliation_cases` | E | Typed discrepancy/current owner/deadline/resolution | Organization/payment connection |
| `reconciliation_case_links` | E | Exact transaction/intent/booking/refund/settlement evidence links | Reconciliation case |
| `reconciliation_runs` | E | Provider range/checkpoint/control totals/outcome | Connection/currency/period |

## `communications` — notifications and reminders

| Table | Stage | Owner/purpose | Important relationships |
|---|:---:|---|---|
| `recipient_endpoints` | F | Purpose-approved push/phone/email/channel endpoint projection/lifecycle | Account/customer/installation/provider participant |
| `communication_preferences` | P | Purpose/channel allow/deny/required/source/version | Recipient/contact/organization |
| `message_templates` | F | Stable purpose/event/channel template identity | Tenant/platform owner |
| `message_template_versions` | F | Immutable localized content/variables/provider mapping | Template, locale |
| `notification_policies` | F | Event/recipient/purpose/channel/fallback/urgency rule version | Organization/operation/platform |
| `notification_intents` | F | One logical event/recipient/purpose/current result | Domain event/resource/recipient |
| `delivery_attempts` | P | One provider/channel submission/status/retry/cost evidence | Notification intent/endpoint/provider request |
| `delivery_events` | P | Provider submitted/delivered/read/failure evidence history | Delivery attempt/event receipt |
| `reminder_schedules` | P | Due time/timezone/domain version/current state | Booking/obligation/quote/notification policy |
| `in_app_notifications` | P | Durable account/business action-center projection/read/archive | Notification intent/account/membership |
| `communication_suppressions` | P | Opt-out/block/bounce/invalid endpoint minimal protected entry | Endpoint/purpose/source |

## `media` — object metadata and processing

| Table | Stage | Owner/purpose | Important relationships |
|---|:---:|---|---|
| `media_objects` | P | Tenant/purpose/classification/storage/checksum/type/size/current lifecycle | Uploader, object provider key reference |
| `upload_sessions` | P | Expiring authorized upload destination/limits/result | Media object, principal/resource |
| `media_processing_jobs` | P | Scan/transform/OCR-approved job/attempt/result | Media object, worker version |
| `media_variants` | P | Derived safe format/size/checksum/storage reference | Parent media object |
| `catalog_media_links` | P | Public/private catalog/service relationship/publication order | Media, service/profile version |
| `conversation_media_links` | E | Message attachment relationship | Media, message |
| `payment_media_links` | P | Receipt/evidence relationship | Media, reported evidence/case |
| `request_media_links` | E | Request/assessment/reference relationship | Media, request/answer/assessment |
| `export_media_links` | E | Privacy/report export package relationship/expiry | Media, export job/case |
| `media_access_events` | E | Risk-classified access/signed issuance/download evidence | Media, principal/case |

Typed link tables avoid an unconstrained generic polymorphic foreign key.

## `sync` — Android commands and convergence

| Table | Stage | Owner/purpose | Important relationships |
|---|:---:|---|---|
| `sync_scopes` | F | Installation/account/tenant/filter/protocol/schema/current cursor state | Device installation, membership/operation |
| `client_commands` | F | Immutable client command envelope/hash/base versions/current processing | Tenant/account/install; typed payload version |
| `client_command_dependencies` | F | Directed acyclic dependency graph | Parent/child client commands |
| `client_resource_mappings` | F | Client-created ID → canonical server ID/type | Tenant/install/accepted command |
| `client_command_results` | F | Stable accepted/duplicate/conflict/rejected/retry outcome | Client command, resource/version/problem |
| `sync_conflicts` | P | Typed preserved input/current context/resolution status | Client command/resource |
| `sync_conflict_resolutions` | P | User/manager selected resolution/new command outcome | Conflict, actor, command |
| `sync_cursors` | F | Opaque scope/filter position/expiry/protocol | Sync scope |
| `sync_changes` | F | Bounded projection change sequence/version/reference | Tenant/scope/resource; retention candidate |
| `sync_tombstones` | F | Minimal delete/archive/revoke/merge propagation marker | Tenant/resource/version/reason |
| `sync_bootstrap_manifests` | F | Snapshot projection/version/page/checksum/activation evidence | Sync scope |
| `sync_runs` | P | Diagnostic push/pull/result/bytes/lag/outcome | Scope/install; minimized metadata |

## `platform` — durable cross-cutting mechanics

| Table | Stage | Owner/purpose | Important relationships |
|---|:---:|---|---|
| `idempotency_records` | F | Logical key/scope/request hash/processing/terminal replay result | Principal/tenant/action/resource |
| `outbox_items` | F | Committed integration event/work pending worker delivery | Producer tenant/aggregate/event schema |
| `consumer_inbox_receipts` | F | Consumer-event dedup/current result | Consumer/event/work |
| `scheduled_tasks` | F | Due/lease/idempotent timer/expiry/reminder/retention work | Owner context/resource/version |
| `process_instances` | F | Long-running workflow current step/deadline/outcome | Tenant/type/version/correlation |
| `process_steps` | F | Append-only step/command/event/attempt/result evidence | Process instance |
| `worker_leases` | F | Optional logical lease/heartbeat where not stored on work row | Work type/resource |
| `feature_flags` | F | Controlled flag definition/default/owner/expiry | Environment/country/tenant eligibility references |
| `feature_flag_overrides` | F | Scoped override/effective/audit | Flag, organization/account/cohort |
| `schema_versions` | F | Application/local/event/provider schema compatibility registry | Migrations/releases |

## `operations` — Ubi internal case work

| Table | Stage | Owner/purpose | Important relationships |
|---|:---:|---|---|
| `cases` | F | Onboarding/support/integration/payment/privacy/trust/security case | Category, tenant/subject, severity, status/deadline |
| `case_subjects` | F | Typed resource/account/organization/provider references | Case; constrained target registry |
| `case_assignments` | F | Workforce team/person claim/escalation history | Case, workforce identity |
| `case_notes` | F | Purpose/classified internal note with author | Case; not domain truth |
| `case_actions` | F | Requested/executed domain command/outcome | Case, workforce principal, audit |
| `privileged_access_requests` | F | Exact purpose/resource/data/expiry/elevation request | Case, requester/approver |
| `privileged_access_grants` | F | Approved scoped JIT/break-glass grant and revocation | Request, workforce session |
| `privileged_access_usage` | F | Resource/action/content access evidence/review | Grant/case/audit |
| `platform_restrictions` | E | Account/business/resource restriction/rule/reason/redress | Case, target owner context |
| `redress_appeals` | E | Appeal/evidence/reviewer/decision | Restriction/case/account/organization |

## `audit` — protected evidence

| Table | Stage | Owner/purpose | Important relationships |
|---|:---:|---|---|
| `audit_events` | F | Append-only significant action/access/decision/outcome | Actor/service/session/tenant/resource/correlation |
| `audit_event_changes` | F | Bounded before/after field/value hash/reference evidence | Audit event; classified/retained |
| `audit_integrity_checkpoints` | F | Sequence/hash/signature/storage checkpoint if approved | Audit event range |
| `audit_exports` | E | Scoped immutable review/export package and requester | Case/organization/date range |
| `audit_access_events` | F | Who queried/exported protected audit | Workforce/account/case |

Audit implementation must not copy secrets/full content unnecessarily and cannot be writable by ordinary application roles except controlled append interface.

## `privacy` — data governance and rights/retention

| Table | Stage | Owner/purpose | Important relationships |
|---|:---:|---|---|
| `data_inventory_items` | F | Field/object/event/system purpose/classification/owner/location | Schema/table/field/provider/version |
| `processing_purposes` | F | Stable business/legal processing purpose registry | Inventory/notices/consents |
| `notice_documents` | F | Stable notice/agreement identity/type/jurisdiction/audience | Versions |
| `notice_versions` | F | Immutable localized published notice/text/effective range | Notice document |
| `notice_acceptances` | F | Actor/subject/version/presentation/auth acceptance evidence | Notice version/account/organization/customer |
| `consent_records` | P | Purpose/channel/allowed/withdrawn/source/version evidence where consent applies | Subject/contact/merchant |
| `retention_rules` | F | Data class/purpose/jurisdiction/trigger/duration/action/owner | Inventory objects |
| `legal_holds` | F | Exact subject/data scope/reason/authority/effective/release | Case/resource/retention task |
| `privacy_cases` | F | Access/correct/export/delete/restrict/objection request and verified decision | Subject/organization/responsibility/jurisdiction |
| `privacy_case_tasks` | F | Per owner/system/provider propagation/action/result/error | Privacy case; data inventory item |
| `privacy_exports` | E | Export build/manifest/access/expiry/deletion evidence | Privacy case/media object |
| `deletion_tombstones` | F | Minimal protected deletion propagation/idempotency proof | Subject/resource/system |
| `subprocessors` | F | Vendor/legal entity/purpose/status/contract/incident contact | Services/data locations |
| `data_locations` | F | Service/region/replica/backup/support access mapping | Subprocessor/inventory |
| `transfer_assessments` | F | Jurisdiction/data transfer basis/risk/review | Data location/subprocessor/country |

## `reporting` — rebuildable projections and metrics

| Table | Stage | Owner/purpose | Important relationships |
|---|:---:|---|---|
| `metric_definitions` | F | Stable metric name/version/formula/source/owner/exclusions | Reports/KPI scorecard |
| `reporting_checkpoints` | F | Projection consumer cursor/version/freshness/rebuild | Event source/projection |
| `operation_daily_metrics` | P | Aggregate booking/request/fulfilment/cancellation/source counts | Tenant/operation/local date/metric version |
| `operation_response_metrics` | P | Request/conversation response distributions | Tenant/operation/date |
| `sync_quality_metrics` | P | Aggregate pending/conflict/rejection/full-resync/performance | Tenant/app version/date; minimized |
| `payment_reconciliation_metrics` | E | Exact aggregate verified/refund/exception/settlement values | Tenant/connection/currency/date |
| `provider_delivery_metrics` | E | Aggregate messaging/notification/provider outcomes/cost | Provider/tenant/date/purpose |
| `experiment_assignments` | E | Approved experiment/cohort/version/subject pseudonym | Experiment; privacy/governance |
| `experiment_outcomes` | E | Governed metric observation | Assignment/metric/version |

Operational reports may initially be views/materialized/query projections rather than all physical tables. They never write owner state.

## `billing` — Ubi merchant subscription (later)

| Table | Stage | Owner/purpose | Important relationships |
|---|:---:|---|---|
| `plans` | E | Ubi software plan/version/entitlements/pricing | Country/currency/effective range |
| `subscriptions` | E | Organization’s Ubi service subscription lifecycle | Organization/plan |
| `usage_records` | E | Idempotent billable Ubi usage measure | Organization/metric/source event |
| `invoices` | E | Ubi-issued commercial invoice lifecycle | Organization/subscription/currency |
| `invoice_lines` | E | Exact plan/usage/credit/fee line | Invoice |
| `subscription_payments` | E | Payment for Ubi software, separate from merchant customer payment | Invoice/provider connection |

## Cross-schema relationship rules

1. Owner table can be written only through its module/runtime role/interface.
2. Tenant-owned child carries `organization_id`; composite FK proves same tenant.
3. Cross-schema FKs are allowed only when they improve integrity without blocking domain lifecycle/deletion/extraction; otherwise immutable ID plus application/event contract and reconciliation is documented.
4. Global Account links to tenant data through Membership/Claim/Party, never a direct “account owns all customer rows” shortcut.
5. Provider/resource polymorphism uses typed link tables or a constrained target registry with validation—not free-form type/ID.
6. Historical snapshot rows retain source version IDs plus copied accepted values; current catalog/customer records are not required to render history.
7. Deleting/closing parent tenant/account/customer does not cascade through booking/payment/audit blindly.
8. Reporting/search/public/sync projections are rebuildable and cannot be mutation parents.

## DDL dependency order

Proposed creation sequence:

1. database roles, schemas, safe extensions and common types/functions;
2. IAM global identity/service tables;
3. Organization/hierarchy/authority;
4. Catalog/configuration and CRM;
5. Scheduling resources/rules;
6. Booking/request core;
7. Access grants, platform idempotency/outbox/timers/sync;
8. Communications/media/integration/payment/conversation expansion;
9. Operations/audit/privacy governance;
10. Reporting/billing projections;
11. composite FKs, RLS policies, indexes and validation in lock-safe migration sequence;
12. seed controlled registries/capabilities/problem/event/schema versions.

## Table ownership/release acceptance

Before a table is implementation-ready it needs:

- bounded-context owner and accountable team/person;
- release/actor/workflow/requirement trace;
- exact columns/types/null/default/semantic definitions;
- PK/FK/unique/check/exclusion/concurrency constraints;
- tenant/RLS/privilege/read/write policy and negative tests;
- approved queries/indexes/cardinality/volume/growth;
- lifecycle/history/version/idempotency rules;
- classification/purpose/retention/deletion/legal hold/export;
- audit/event/analytics obligations;
- migration/backfill/rollback/restore;
- fixtures/factories and domain/concurrency/property tests.

## Open decisions

- Reduce catalogue to the selected first-release vertical.
- Exact cross-schema FK and typed-link target implementation.
- Account/Person necessity and identity-provider schema boundary.
- Snapshot/form/provider JSONB versus typed child-table details.
- Pooled-capacity and availability recurrence tables by first archetype.
- RLS policy/schema ownership and runtime-role mapping.
- Audit/privacy/reporting physical separation/partition timing.
- Billing/commercial subscription scope and provider.

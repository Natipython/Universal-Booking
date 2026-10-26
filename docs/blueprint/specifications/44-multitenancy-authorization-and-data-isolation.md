# Multitenancy, authorization and data isolation

Status: **Draft — defense model proposed; threat-model, PostgreSQL and load prototypes required**  
Depends on: identity/organization model, roles and permissions, database architecture and constraints  
Blocks: RLS DDL, repository interfaces, operations-console access, security testing and release authorization

## Purpose

Ubi will hold bookings, contacts, conversations, schedules and payment evidence for many unrelated businesses. This document establishes how one shared platform keeps those businesses separate while allowing a person to belong to several businesses deliberately.

The central rule is:

> Authentication says who is present. Membership, capability, scope, resource state and purpose together decide what that principal may do. A record ID or selected organization is never authority by itself.

## Initial tenancy decision

Use a shared application and shared managed PostgreSQL database with row-level organization isolation for the initial release.

Do not create one database/schema per small business initially. That model would multiply migrations, backups, pooling, support and reporting work before Ubi has evidence that dedicated tenancy is required.

This decision does not mean “one weak shared bucket.” Isolation is layered across API authorization, database keys/RLS, cache/search/object prefixes and ACLs, queues/events, exports, operations access and continuous negative tests.

A future dedicated database/region may be offered only for measured regulatory, risk, contract or scale needs. Moving a tenant then requires an approved isolation tier, routing directory, migration/rollback, export, restore and key-management design.

## Tenant and actor meanings

- **Organization** is the primary commercial/security tenant.
- **Brand** is a public commercial identity; it is not an authorization tenant.
- **Branch** is a physical/operational scope within an organization.
- **Booking operation** is a workflow/catalog/schedule boundary within an organization.
- **Account** is a global authentication identity and may have zero or many memberships.
- **Membership** is the account-to-organization authority relationship.
- **Customer record** is tenant-private CRM data and is not automatically an Account.
- **External participant** is provider-connection-scoped and is not automatically an Account or Customer.
- **Workforce principal** is a controlled Ubi operator identity, separate from merchant membership.
- **Service identity** is a workload, never a human account disguised as a bot.

Organization is the database tenant key even when a row is scoped more narrowly to a branch, operation, queue or resource.

## Authorization decision model

Every protected request is evaluated in this order:

```text
authenticated principal and valid session
→ active account/workload lifecycle
→ explicitly selected organization
→ active membership or approved non-member path
→ current authority version
→ capability for exact action
→ scope covers exact branch/operation/resource/queue/assignment
→ resource lifecycle/state permits the action
→ required approval or step-up is satisfied
→ field/purpose/channel restrictions are satisfied
→ command executes with tenant-bound repository/database context
```

Failure at any step denies by default and returns a safe problem response. It never silently falls back to another organization, broader role or “first membership.”

### Authorization input contract

The server constructs an immutable request context containing at least:

```text
request_id / trace_id
principal_type and principal_id
account_id or service_identity_id
session_id and assurance level
organization_id
membership_id and authority_version, when applicable
capability being exercised
resolved scope/resource IDs
client/app/version/channel
operations case/grant, only for approved workforce access
```

Client-supplied IDs are lookup inputs only. The server resolves and verifies every relation.

## Capability and scope rules

Role names are user-friendly bundles, not authorization code. Code checks stable capabilities such as `booking.confirm`, `schedule.override` or `customer.contact.view`, combined with scope and state.

### Scope hierarchy

Proposed scope types:

```text
organization
branch
booking_operation
resource
queue
assigned_only
```

A broad scope contains narrower targets only through verified same-tenant hierarchy. No string-prefix or UI-screen assumption establishes containment.

### Conflict precedence

1. Platform/legal/security restriction denies.
2. Explicit active scoped deny denies.
3. Required approval/step-up not satisfied denies.
4. A matching scoped capability grant or role capability allows.
5. Otherwise deny.

Owners are not silently superusers. Certain actions—ownership transfer, payment refund, schedule override, sensitive export, integration credential change, support impersonation—can require step-up, separation of duties or a second approval even for an owner.

### Resource-state authorization

Possessing a capability does not make every transition legal. For example:

- `booking.confirm` cannot confirm a cancelled booking or consume expired capacity;
- `refund.execute` cannot exceed verified refundable value;
- `message.send` cannot bypass a block, channel policy or automation takeover state;
- `customer.delete` cannot erase protected commercial or legal evidence;
- `membership.manage` cannot remove the last verified owner without transfer/recovery rules.

## Data classes by isolation path

| Data class | Examples | Access path |
|---|---|---|
| Tenant operational | Bookings, customers, services, resources, conversations | Tenant-bound server module + RLS |
| Global controlled registry | Capabilities, provider registry, schema versions | Read-only approved view/function; platform owner write |
| Global account security | Accounts, authenticators, sessions | IAM service/module only; no tenant browse |
| Public projection | Published profile/service/location | Anonymous public API/view with publication/abuse controls |
| Ubi workforce | Cases, JIT grants, internal notes | Operations service with workforce authorization and audit |
| Protected audit/privacy | Audit events, legal holds, privacy cases | Narrow append/query interfaces; no ordinary merchant mutation |
| Rebuildable projection | Search, reporting, sync view | Source-authorized projector and tenant-filtered reader |

Avoid tables where `organization_id NULL` ambiguously means “global.” Prefer distinct global and tenant tables/views so fail-closed behavior is obvious.

## Layered isolation architecture

### Layer 1 — application authorization

The application owns fine-grained policy: capability, scope, state, field, approval, purpose and channel. Domain modules expose authorized commands/queries rather than generic cross-schema CRUD.

Repository signatures require trusted tenant context. A method such as `getBooking(id)` is rejected in favor of a tenant/resource-aware contract such as `getBooking(context, bookingId)`.

Batch APIs authorize every item or guarantee one identical resolved scope. Partial authorization outcomes are explicit; they never leak existence through different timing/messages unnecessarily.

### Layer 2 — tenant-qualified relational integrity

Tenant-owned rows contain immutable `organization_id`. Composite FKs prove parent and child belong to the same tenant. This catches accidental/malicious cross-tenant ID pairing even before an RLS read.

### Layer 3 — PostgreSQL privileges and RLS

Runtime roles:

- do not own tables or schemas;
- do not have superuser or `BYPASSRLS`;
- cannot change DDL;
- receive only exact schema/table/function privileges;
- use RLS on all tenant-owner/base tables, with `FORCE ROW LEVEL SECURITY` where compatible;
- cannot access base tables through an accidentally permissive `public` grant.

RLS provides a coarse tenant boundary and selected special paths. It does not replace domain capability/state authorization.

### Layer 4 — downstream propagation

Every tenant-bearing event, queue item, object, search document, cache key, sync cursor/change and export carries verified organization scope. Consumers validate it against the referenced owner rather than trusting an arbitrary header.

### Layer 5 — detection and proof

Negative tests, canary tenants, audit, anomaly alerts, access reviews and restore/export verification continuously test isolation. “No incident reported” is not proof.

## PostgreSQL request context

Tenant queries run only inside an explicit transaction. After application authentication/authorization, the server sets transaction-local context such as:

```text
app.current_principal_type
app.current_account_id
app.current_organization_id
app.current_membership_id
app.current_authority_version
app.current_service_identity_id
app.current_workforce_grant_id
app.current_request_id
```

Requirements:

- use `SET LOCAL` or a reviewed transaction-entry function, never session-persistent state;
- context helper functions return `NULL`/deny on absent or malformed values, not a broad default;
- every success, exception and cancellation commits/rolls back before returning the pooled connection;
- pool checkout/reset tests prove no prior tenant context remains;
- transactions with tenant base-table queries cannot run in accidental autocommit before context;
- read replicas, jobs and migrations have separately designed context/roles.

Conceptual tenant policy:

```sql
using (
  organization_id = auth_context.current_organization_id()
)
with check (
  organization_id = auth_context.current_organization_id()
)
```

The helper is small, stable and safely configured. Policy predicates use indexed columns and avoid per-row network calls or large joins.

### Honest limitation

If the whole API runtime is compromised and its database credential may establish arbitrary transaction context, a simple custom setting cannot prove which end user was authenticated. RLS still protects against missing filters and many coding mistakes, but not every full-process compromise.

Mitigations include least-privilege module/workload roles, short-lived workload credentials, network isolation, protected session validation, narrow database functions for especially sensitive actions, field encryption/key separation, anomaly detection and rapid credential rotation. Do not market RLS as magical cryptographic tenant isolation.

## RLS policy families

### Tenant current-state tables

- `SELECT`: current tenant rows, further filtered by application scope.
- `INSERT`: new row `organization_id` must equal current tenant.
- `UPDATE`: both old visibility and new row check equal current tenant; tenant key cannot change.
- `DELETE`: generally unavailable to ordinary runtime; lifecycle command used instead.

### Append-only history/evidence

- Runtime can append only through owner command/function.
- Ordinary update/delete privilege is absent.
- Merchant reads use approved tenant-filtered views that omit internal/security fields.
- Audit/privacy access itself is audited.

### Global Account/IAM

Do not expose Account rows through tenant RLS merely because an account has Membership. IAM functions return only purpose-limited identity/profile facts. A business cannot browse another business’s members or an account’s other memberships.

### Public data

Anonymous role reads only explicitly published projection rows/columns. It never relies on “tenant row marked public” across arbitrary base tables. Unpublication, restriction and slug reassignment invalidate caches/search promptly.

### Background workers

Preferred worker patterns:

1. claim a bounded work item through a restricted queue interface;
2. obtain its immutable tenant/resource scope;
3. open a tenant-context transaction;
4. process one tenant/work unit;
5. commit evidence/result;
6. release context before the next item.

Cross-tenant jobs do not casually disable RLS. A separately reviewed service role/function may enumerate due work metadata, but domain reads/writes re-enter the exact tenant. Maintenance/migration roles are not reused by workers.

### Reporting

Merchant reports are tenant-bound projections/views. Platform aggregate reporting uses minimized/de-identified governed projections and a separate role. Raw cross-tenant joins are not available to general analysts or customer support.

### Ubi operations access

Support cannot impersonate a merchant silently or browse production tables freely.

Required path:

1. a case exists with category, tenant/subject, purpose and urgency;
2. workforce principal authenticates with required assurance;
3. an exact read/action scope and expiry is requested;
4. approval is obtained where risk requires;
5. a short-lived JIT grant is issued;
6. every viewed resource/field and command is recorded;
7. sensitive values are masked unless separately necessary;
8. grant expires/revokes and usage is reviewed.

Break-glass is narrowly available for real emergencies, alerts reviewers immediately and requires post-event review. It does not mean permanent administrator access.

## Connection and database-role separation

| Connection role | May do | Must not do |
|---|---|---|
| API tenant runtime | Tenant-context owner-module queries/commands | DDL, bypass RLS, arbitrary global/workforce reads |
| IAM runtime | Approved identity/session functions | Browse tenant CRM/bookings |
| Worker runtime | Claim work and process scoped domain effects | Use migration/owner privileges |
| Public API runtime | Read safe public projections | Read tenant base/private rows |
| Operations runtime | Case/JIT-controlled interfaces | General unlogged tenant access |
| Reporting runtime | Governed views/projections | Write owner data or expose raw PII broadly |
| Migration role | Reviewed DDL/migrations | Serve application traffic |
| Schema owner | Own objects, normally no-login | Authenticate applications |

If one deployment process initially carries several runtime responsibilities for cost reasons, logical grants/interfaces still remain separate and the consolidation is a recorded temporary risk with a split trigger.

## Field- and purpose-level protection

Row access does not imply access to every column.

- Sensitive contact, address, payment, credential and private note fields use purpose-specific views/DTOs.
- Ordinary booking lists show only the minimum customer details needed to perform work.
- Search results reveal no more than the caller could open.
- Exports require an explicit capability, scope, reason and audit record; large exports may require approval and asynchronous secure delivery.
- Secrets, OTPs and payment credentials never appear in merchant/customer/support DTOs or logs.
- Data used for analytics/training/model improvement requires separate documented purpose and governance; operating Ubi is not blanket permission for every secondary use.

## Revocation and consistency

Membership suspension, role/scope change, organization restriction and session revocation update an `authority_version` or equivalent epoch.

- New requests must use current authority.
- Access tokens are short-lived and checked against revocation/epoch for risky actions.
- Offline clients may display cached data within approved policy but cannot make newly unauthorized server changes.
- Sync rejects revoked scopes and propagates authorization tombstones/required local purge.
- Long-running jobs re-check authority/policy before a privileged irreversible step.
- Cache keys include authority-sensitive scope/version where user-specific results are cached.

Revocation objectives are measurable by risk class; “eventually” is not sufficient for ownership/payment/export/security actions.

## Cross-system isolation rules

### Cache

Keys include environment, data class, organization and resource/query version. Never cache tenant data under a user-supplied unverified organization. Public and private namespaces are separate. Authorization decisions are short-lived/versioned and fail closed.

### Search

Every private document has immutable organization and ACL/scope fields set by a trusted projector. Query filters are server-generated and mandatory. Search hit IDs are re-authorized before sensitive hydration/action.

### Object storage/media

Private bucket/object metadata records authoritative organization/owner/purpose. Object keys are unguessable locators, not authority. Downloads/uploads use short-lived scoped grants; callbacks validate key, tenant, checksum and expected processing job.

### Queues/events

Producer writes tenant from owner state. Consumer validates event tenant against fetched aggregate before mutation. Dead-letter and operations tools mask/minimize payloads and enforce the same tenant/case controls.

### Logs/traces/metrics

Telemetry minimizes personal content, uses pseudonymous IDs where possible and separates tenant/business support views from platform operations. Trace IDs are not authorization. Log search access is governed and audited.

### Export/backup/restore

Tenant export selects through a verified manifest and reconciliation counts; no guessed prefix alone. A tenant-level restore is a controlled reconstruction/import unless the storage architecture proves isolated physical restore. Full backup access is highly privileged and cannot be given to ordinary support.

## Preventing enumeration and confused-deputy failures

- Protected not-found versus forbidden responses follow a consistent disclosure policy.
- Public references/slugs are rate-limited and reveal only published facts.
- IDs remain opaque but opacity is not treated as security.
- Server ignores or rejects body/header organization that conflicts with the resolved membership/resource.
- Provider callbacks bind provider application, environment, connection and external asset before accepting tenant routing.
- Scoped links bind purpose, resource, allowed action, expiry/use limit and optional verification; opening one cannot upgrade to normal tenant membership.
- A service acting for another service passes authenticated workload identity and immutable causation, not a forged end-user header.

## Security and isolation test matrix

Every tenant-owned endpoint/query/command gets positive and negative tests. Required classes include:

1. Account A in Organization A cannot read/update/delete Organization B row by B’s valid ID.
2. Account in A and B sees only the explicitly selected tenant per request.
3. Suspended membership fails despite a previously valid token.
4. Branch-scoped staff cannot access another branch in the same organization.
5. Operation/queue/assigned-only boundaries are enforced.
6. Capability without legal resource state still fails.
7. Missing database context returns zero/error, never all tenants.
8. Wrong tenant on child insert fails composite FK and RLS.
9. Attempt to change `organization_id` fails.
10. Runtime role cannot disable RLS, change DDL or query owner-only tables.
11. Table owner/migration role is not the runtime role; forced RLS behavior is verified.
12. Pooled connection reused after success/error/cancellation has no old context.
13. Worker cannot use one claimed tenant to mutate another referenced ID.
14. Search/cache/object/export paths reproduce negative isolation behavior.
15. Public role cannot reach an unpublished/private base row or column.
16. Scoped-link holder can perform only its allowed action and cannot enumerate neighbors.
17. Operations access without active case/JIT grant fails and is alerted.
18. Field masking/export limits match capability and purpose.
19. RLS policies retain index-supported query plans at representative tenant sizes.
20. Backup/restore and analytics fixtures contain synthetic canary tenants to detect mixing.

Property-based tests generate arbitrary tenant/resource combinations. CI runs schema privilege/RLS linting and migration tests. Pre-production uses realistic cardinality. Production synthetic isolation probes must not expose real customer data.

Diagram: `diagrams/44-authorization-and-tenant-isolation.mmd`.

## Review and release gates

Before first production data:

- every table is classified as tenant/global/public/workforce/protected/projection;
- every tenant base table has RLS and composite tenant FK coverage, or a signed exception;
- grants/default privileges are generated and inspected from a clean database;
- every route maps to capability, scope, state, fields and audit duty;
- all runtime roles fail attempted DDL/BYPASSRLS/cross-schema access;
- pooling/context cleanup tests pass under errors/timeouts;
- API, worker, search, object, sync and export negative tests pass;
- JIT operations path exists before support receives production access;
- authority revocation SLOs are measured;
- threat model and external security review cover IDOR/BOLA, confused deputy, SSRF/provider callbacks, token/link leakage and backup access;
- country-specific privacy/legal review confirms responsibilities and access controls.

## Open decisions

- Exact first-release scope hierarchy and capability registry.
- Access-token lifetime, authority-version verification and high-risk step-up mechanism.
- Whether database request context uses direct `SET LOCAL` or one reviewed entry function.
- Which modules justify distinct PostgreSQL runtime roles at pilot scale.
- Field protection/encryption/tokenized lookup strategy by launch-country risks.
- Operations JIT provider/build-versus-buy decision.
- Dedicated-tenant/region eligibility and migration trigger.
- Measurable revocation targets for ordinary versus high-risk actions.

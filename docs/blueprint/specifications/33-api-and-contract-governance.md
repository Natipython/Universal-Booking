# API and contract governance

Status: **Draft**  
Depends on: system context, bounded contexts and platform components  
Blocks: client implementation, provider adapters, event schemas, testing, security review and support policy

## Purpose

This document defines how Ubi clients, modules, workers and providers communicate without sharing hidden assumptions. It establishes stable identifiers, authorization context, idempotency, concurrency, errors, versioning and change-control rules before endpoints are implemented.

## Contract families

Ubi has separate contract families:

| Contract | Participants | Initial form | Authority |
|---|---|---|---|
| Client/public HTTP API | Android, customer web, business web, operations console ↔ Ubi | HTTPS + JSON, resource/query and explicit command endpoints | Owning domain module |
| Internal module contract | Application workflow ↔ domain modules | Typed in-process commands/queries/results | Owning module |
| Synchronization contract | Android sync engine ↔ Ubi | Bounded HTTPS batch/change-feed protocol | Domain owner through Sync envelope |
| Domain/integration events | Modules/workers/process managers | Versioned durable event envelope | Event-producing context |
| Provider adapter contract | Provider-specific ingress/egress ↔ canonical Ubi domain adapter | Provider protocol outside; typed canonical contract inside | Integration + consuming domain |
| Public/scoped links | Browser/app resolver ↔ Access Grant and target domain | Verified HTTPS URL with opaque token | Link context plus target domain |
| Media contract | Client/domain ↔ object/media service | Initiate/upload/finalize/access sequence | Media plus linking domain |
| Realtime contract | Later native messaging/calls | Authenticated WebSocket/signaling plus durable HTTP recovery | Conversation/Call owner |

An HTTP DTO, provider payload, database row, domain aggregate and analytical event are not the same object and must not be reused as if they were.

## Initial client API style

### Direction

The initial external API SHOULD use:

- HTTPS;
- JSON request/response representations;
- ordinary HTTP method/status/cache semantics;
- `/v1/` major-version route namespace;
- OpenAPI as the machine-readable client contract;
- resource-oriented queries plus explicit command endpoints for lifecycle actions;
- RFC 9457 `application/problem+json` for structured HTTP problems;
- cursor pagination;
- opaque immutable identifiers;
- stable idempotency and conditional/version controls.

GraphQL is not an initial requirement. gRPC MAY later serve controlled internal/high-throughput boundaries, but putting a network protocol between in-process modules is prohibited.

### Resource and command examples

```text
GET    /v1/me
GET    /v1/organizations
POST   /v1/organizations
GET    /v1/organizations/{organization_id}
GET    /v1/organizations/{organization_id}/operations
POST   /v1/organizations/{organization_id}/operations
GET    /v1/operations/{operation_id}/services
POST   /v1/operations/{operation_id}/bookings
GET    /v1/bookings/{booking_id}
POST   /v1/bookings/{booking_id}/actions/confirm
POST   /v1/bookings/{booking_id}/actions/reschedule
POST   /v1/bookings/{booking_id}/actions/cancel
GET    /v1/operations/{operation_id}/availability
POST   /v1/operations/{operation_id}/holds
POST   /v1/requests/{request_id}/quotations
POST   /v1/quotations/{quotation_id}/actions/accept
POST   /v1/sync/commands
GET    /v1/sync/changes?cursor={opaque_cursor}
```

Exact resource names remain subject to API review. The rules are:

- `GET`/`HEAD` never causes a business mutation.
- Lifecycle actions use `POST` to a clearly named command/action endpoint when CRUD semantics would hide the state transition and consequences.
- `PUT` replaces a complete client-owned representation only when that is genuinely the model.
- `PATCH` is used only with an approved patch format and field-authorization semantics; generic arbitrary-object patching is prohibited.
- `DELETE` means the documented resource lifecycle action, not guaranteed physical erasure.
- Provider callbacks have separate provider-specific routes and authentication, not public client schemas.

## Addressing and tenant context

Identifiers in paths are opaque locators. They do not authorize access.

For every protected request, the server resolves:

1. authenticated account/service/link/provider principal;
2. assurance/session/grant state;
3. target resource and its immutable tenant;
4. current membership/capability/scope or customer/token relationship;
5. state-specific domain policy.

Rules:

- The server MUST NOT trust a client `organization_id`, role, branch or username as authority.
- Resource tenant equality is validated below routing.
- Where organization appears in a path, it improves clarity/routing and is checked against the child resource; it is not the security boundary by itself.
- Merchant APIs never expose cross-tenant lookup by customer contact, public slug or predictable ID.
- Public endpoints return deliberately public projections only.
- Disclosure-sensitive authorization MAY return a not-found response instead of confirming resource existence.

## Request metadata

### Standard/request headers

| Header/concept | Use | Rule |
|---|---|---|
| `Authorization` or approved secure session cookie | Principal/session | Never accepted in URL; short/revocable policy defined by IAM |
| `Idempotency-Key` | Logical mutation retry identity | Required for designated POST commands; not a secret |
| `If-Match` / representation version | Conditional mutation | Used for online resource edits where HTTP conditional semantics fit |
| `Accept-Language` | Preferred response human text/representation | Language cannot change currency/calendar/domain truth |
| `Content-Type` / `Accept` | Representation negotiation | Strict supported media types; reject ambiguous content |
| Client build/platform metadata | Compatibility and diagnostics | Bounded/minimized; cannot authorize |

The edge generates or validates a safe correlation identifier. A malicious client value cannot control logs, traces or another request’s identity.

Device ID, app-install ID and client-command ID are distinct. None is a hardware fingerprint or authentication credential by itself.

## Response conventions

### Successful single resource/command

Return a typed representation containing the minimum data needed for the caller, including where relevant:

- immutable resource ID;
- context-qualified state;
- server resource version;
- authoritative timestamps/timezone;
- `pending/action_required/next_action` semantics;
- links/actions currently available as convenience only—the server still reauthorizes later;
- source/provenance and confirmation/payment truth needed for comprehension.

HTTP headers SHOULD provide correlation ID, `ETag`/cache controls and `Location` for newly addressable resources where appropriate.

### Collections

Proposed shape:

```json
{
  "items": [],
  "next_cursor": null,
  "has_more": false
}
```

Rules:

- cursor is opaque, integrity-protected or server-resolved and scoped to query/tenant/principal as required;
- ordering is stable and explicitly documented, normally `(business_time, immutable_id)` or another deterministic key;
- page size has server maximum and a conservative default;
- filtering/sorting fields are allow-listed and indexed before release;
- total counts are omitted when expensive, sensitive or misleading;
- cursor expiry/restart returns a typed recoverable problem.

### Empty and asynchronous outcomes

- `204` is used only when the client genuinely needs no representation.
- A long-running accepted operation may return `202` with durable operation/process ID and status route.
- “Queued”/`202` never means the domain outcome succeeded.
- Provider-dependent work returns current state and polling/subscription/recovery guidance, not an indefinite open request.

## Error contract

Ubi adopts RFC 9457 Problem Details for HTTP APIs.

```json
{
  "type": "https://api.ubi.example/problems/booking-capacity-conflict",
  "title": "That time is no longer available",
  "status": 409,
  "detail": "Choose another available time or ask a manager to review the conflict.",
  "instance": "urn:ubi:problem:opaque-occurrence-id",
  "code": "booking.capacity_conflict",
  "correlation_id": "opaque-support-id",
  "retryable": false,
  "current_version": 17,
  "field_errors": []
}
```

Rules:

- `type` and `code` are stable machine semantics; client logic MUST NOT parse `title` or `detail`.
- `title`/`detail` are human-readable/localizable and never contain internal stack traces, SQL, secrets or unrelated resource existence.
- `instance`/correlation identifiers are safe for support and do not expose record counts.
- `retryable` describes whether the same logical request may be retried, not whether a new business action is safe.
- Field problems use stable field paths/codes and preserve valid input.
- One primary problem is returned; structured extensions may include related field conflicts/options.

### Status mapping direction

| HTTP status | Ubi meaning direction |
|---|---|
| `400` | Malformed request/protocol/unsupported representation |
| `401` | Missing/invalid/expired authentication; reauthentication may be needed |
| `403` | Authenticated/granted principal cannot perform this action; resource disclosure is safe |
| `404` | Resource absent or existence intentionally undisclosed |
| `409` | Domain state, capacity, idempotency or synchronization conflict |
| `410` | Public/scoped resource intentionally no longer available where disclosure is safe |
| `412` | Explicit conditional request/precondition version failed |
| `422` | Well-formed request violates field/domain validation before a state conflict |
| `423` | Use only if an explicitly documented resource lock/approval state fits; not a generic busy error |
| `429` | Rate/abuse/quota limit; include safe retry guidance where appropriate |
| `503` | Ubi/provider dependency temporarily unavailable for that operation |

Do not return `200` with `success: false` for a failed HTTP operation. Batch/sync commands may return HTTP success for a valid batch envelope with per-command typed outcomes.

## Idempotency contract

### Required operations

Idempotency is mandatory for:

- booking/request/quotation submission and lifecycle commands;
- hold placement/consumption;
- offline synchronized commands;
- payment intent/allocation/refund/provider operations;
- invitation/link issuance where retry could duplicate effects;
- notification/manual message sends;
- provider callbacks and worker/process-manager steps;
- privacy/organization closure actions.

### Scope and behaviour

An idempotency record binds:

- logical key/client command ID;
- principal/tenant/action route and API/command version;
- canonical request fingerprint excluding unstable transport fields;
- first-received time and processing/terminal state;
- resulting resource IDs/version/status/problem necessary to replay;
- expiry/retention class and security/audit linkage.

Rules:

1. Same key + same logical request returns the same committed/terminal result without repeating effects.
2. Same key + different fingerprint returns `409 idempotency.key_reused`.
3. Concurrent duplicate waits briefly or returns a typed `processing` result with the same operation identity; it cannot start twice.
4. A transport timeout does not authorize a new key for the same intended action.
5. Retention covers the maximum offline/provider/retry/reconciliation window. Exact periods are selected by action; financial and long-offline commands may require longer evidence/tombstones.
6. Keys are unpredictable enough to avoid accidental collision but are not treated as authorization secrets.
7. Side effects after the authoritative transaction inherit causation/idempotency identity and deduplicate independently.

HTTP method idempotency alone is insufficient for POST commands and external side effects.

## Optimistic concurrency and versions

Every mutable aggregate exposed to clients has a monotonically comparable opaque/integer `version` owned by the server.

### Online editing

- `GET` returns representation version/`ETag`.
- A material edit sends `If-Match` or typed `expected_version`.
- Version mismatch returns current version and a domain-safe conflict/reload/merge path.
- Server may auto-merge only explicitly commutative/non-material fields.

### Domain commands

Commands include the evidence versions they rely on:

- booking version;
- selected catalog/price/policy version;
- availability/hold identity and expiry;
- accepted quotation version;
- authorization/approval reference where relevant.

Final command revalidates current authoritative invariants even when all versions match.

### Offline commands

Offline uses `base_version`, stable client command/resource IDs and domain intent rather than a generic object replacement. Conflict rules are defined in the sync specification.

## Identifier contract

- Internal IDs are immutable, opaque, non-reassigned and safe to expose only where authorized.
- Public booking codes/slugs are separate and may be human-friendly/changeable under policy.
- Token secrets are high entropy and never used as ordinary resource IDs.
- Provider IDs are namespaced by provider/environment/connection/account.
- Offline-created resources use collision-resistant client-generated IDs or a stable client-to-server ID mapping; exact format awaits database decision.
- IDs MUST NOT embed phone, email, tenant name, date, region or sequential volume.
- JSON represents identifiers as strings.

## Time, calendar, money and number contract

### Instants and dates

- Machine interchange uses ISO 8601/RFC 3339-compatible unambiguous timestamps with UTC offset; authoritative instants are normalized to UTC while retaining relevant IANA timezone/original local context in explicit fields.
- Date-only fields use `YYYY-MM-DD` plus calendar/timezone context where required; they are not midnight timestamps.
- Recurrence rules/timezone/effective range are structured fields, not display strings.
- Ethiopian calendar is an input/display representation with explicit calendar identifier and verified conversion; canonical machine dates remain explicit.
- Durations are numeric unit/value or documented ISO duration, not localized text.

### Money

Proposed transport representation:

```json
{
  "amount_minor": "12500",
  "currency": "ETB"
}
```

The integer may be a JSON string to avoid cross-language integer limits. Currency minor-unit/rounding policy is server-owned. For currencies/uses requiring a decimal representation, the contract defines it explicitly; binary float is prohibited.

### Phone and locale

- Normalized phone contact uses country-aware parsing/E.164 where valid plus protected entered/display context as needed.
- Locale uses BCP 47; timezone uses IANA identifiers; currency uses ISO 4217.
- Localized display strings are separate from structured machine values.

## Create/update input policy

- Request schemas allow-list every writable field.
- Server-owned fields such as tenant, owner, state, totals, verification, audit, permission and timestamps are never mass-assignable.
- Unknown fields are rejected for commands where silent ignore would hide client/server incompatibility; explicit extensible metadata areas are tightly typed/limited.
- `null`, absent and empty have documented distinct meaning where applicable.
- Trimming/Unicode normalization is field-specific and never silently changes legal/accepted text.
- Validation runs syntactic → authorization → referential/tenant → domain/state → concurrency order while avoiding sensitive disclosure.
- Business-custom intake uses typed approved field definitions; it cannot inject arbitrary JSON into protected domain objects.

## Mutation response and receipts

Every consequential mutation result identifies:

- logical command ID/idempotency key reference;
- resulting resource and version;
- actual domain state, not generic `success`;
- committed server time;
- pending external/local effects separately;
- conflict/rejected field input recovery where applicable;
- current next action/actor/deadline;
- audit/support reference appropriate to the caller.

A request interrupted after commit can recover the result through the idempotency contract or resource/operation query.

## Batch and synchronization APIs

The sync endpoint is not a generic bulk CRUD API.

Rules:

- bounded command count/compressed byte size and processing time;
- each command has client command ID, type/version, tenant/actor/device context, creation/base versions, dependencies and typed payload;
- batch envelope authentication does not skip per-command current authorization;
- each result is `accepted`, `duplicate`, `conflict`, `rejected`, `retry_later` or `dependency_failed` with stable server resource/version/problem;
- one invalid command does not necessarily roll back unrelated commands; declared dependency groups may be atomic;
- commands for different tenants are not mixed;
- changes feed is cursor/version based, contains projections/tombstones and can demand full resync;
- compression is allowed with decompressed-size/bomb protections.

## Public and scoped-link API

- Public endpoints have separate schemas, rate limits, cache policy and enumeration threat tests.
- Public availability returns safe candidate intervals/capacity bands only, never raw staff/customer occupancy.
- Scoped tokens are presented in URL only for initial navigation when necessary, then exchanged for a short-lived protected session/grant and removed from visible/referrer-prone navigation.
- Token response pages forbid third-party resources until token stripping/exchange and apply restrictive referrer/cache policies.
- Link open does not mutate acceptance/confirmation.
- Sensitive actions may require phone/account step-up even with a valid grant.

## Media API

Preferred sequence:

1. Client requests upload intent with purpose, tenant/linking resource, declared type/size/checksum.
2. Server authorizes and returns an expiring, single-purpose upload destination/credentials.
3. Client uploads directly where approved.
4. Client finalizes with upload/object ID.
5. Media worker validates actual size/type/checksum, scans/processes and marks ready/quarantined/rejected.
6. Domain links only approved/appropriate object states.
7. Downloads use short authorized URLs or protected proxy; public derivatives are separately published.

Clients never choose permanent storage keys or ACLs.

## Provider ingress contracts

Each provider endpoint defines:

- exact environment/provider/account route and authentication/signature validation;
- raw-body/headers required for verification;
- strict accepted size/content/time skew/replay limits;
- provider event/deduplication identity;
- acknowledgement deadline/semantics;
- out-of-order and unknown-event handling;
- canonical adapter version;
- tenant/connection resolution;
- sensitive payload storage/redaction/retention;
- reconciliation when callbacks are missing/ambiguous;
- sandbox/test/live separation.

Never expose a single unauthenticated “generic webhook” that trusts provider/type/tenant fields in the body.

## Authentication, browser and mobile transport controls

- All production traffic uses approved modern TLS; cleartext is prohibited.
- Browser cookie sessions, if used, are Secure/HttpOnly/SameSite as appropriate and all state-changing routes have CSRF/origin defenses.
- Token-bearing browser APIs use strict CORS allow-lists; CORS is not authorization.
- Mobile app distribution/signing/app-attestation signals may increase risk confidence but cannot replace account/session authorization.
- Access tokens are short-lived and refresh/session families are revocable; exact token format is an IAM decision.
- Secrets/tokens are never put in ordinary query parameters, analytics or logs except the constrained initial opaque-link case with immediate safe exchange.
- Responses set cache policy by public/private/sensitive class; authenticated/customer/payment responses are not shared-cacheable.
- Request decompression, JSON depth/count/number/string limits prevent resource exhaustion.

## Rate limits, quotas and overload

Limits are layered by route risk and available keys:

- IP/network/device/app-install;
- account/session/token;
- organization/connection/provider account;
- destination phone/contact;
- resource/action and global safety.

Rules:

- OTP, public availability, link redemption, search, message/call and provider ingress have dedicated abuse policy.
- One tenant/provider cannot exhaust all booking capacity.
- Accepted durable work is not silently discarded during overload.
- `429`/`503` returns safe retry guidance and correlation; client applies jitter/backoff and command idempotency.
- Limits account for shared carrier NAT and weak networks; IP-only permanent blocking is unsafe.
- Commercial quotas and protective rate limits are distinct and explained truthfully.

## Compatibility and versioning

### Client API

- `/v1` is a major compatibility boundary, not the application release number.
- Additive optional fields are normally backward-compatible; clients ignore unknown response fields safely.
- Changing meaning, required input, state semantics, enum exhaustiveness or authorization is reviewed as potentially breaking.
- Enum clients implement unknown/future-safe handling unless the contract explicitly closes the set.
- Fields are deprecated with replacement, telemetry, supported-client window and removal approval.
- Server may require a minimum client only for security/data-integrity incompatibility and must provide an understandable upgrade path.
- Offline queued commands remain interpretable or receive a recoverable version rejection across the supported offline window.

### Contract-first change process

1. Requirement/workflow/owner identified.
2. Threat/privacy/data/cost impact classified.
3. Schema and examples changed in machine-readable source.
4. Compatibility diff reviewed.
5. Consumer/provider contract tests updated.
6. Server supports old/new deployment window.
7. Observability and rollback/roll-forward defined.
8. Documentation/SDK/client generated artifacts updated.
9. Deprecation/removal tracked to completion.

Provider adapters pin/test provider API versions and treat provider changes as external compatibility events.

## Schema and documentation governance

The repository will contain:

- OpenAPI documents split by approved public/customer/business/operations surface or tags with strict access review;
- shared primitive schemas with small ownership;
- problem-type registry;
- command/event catalog and examples;
- provider canonical contracts and fixtures;
- data classification/field owner annotations;
- compatibility lint rules;
- generated client models where safe;
- contract tests and published internal developer reference.

Generated code never replaces domain validation. API docs exposed publicly must not reveal operations-only routes/schemas or security internals unnecessarily.

## Contract observability

Measure by endpoint/command/problem code/version—not raw path IDs:

- requests, latency, status/problem outcome;
- payload compressed/decompressed size;
- idempotency hit/reuse/processing duration;
- concurrency conflict and user resolution;
- client version/compatibility failures;
- authz deny/step-up without sensitive resource detail;
- provider timeout/retry/reconciliation;
- sync batch size/command outcome/lag;
- rate/overload result;
- schema/unknown-event errors.

High-cardinality personal/provider identifiers remain out of ordinary metric labels.

## Contract acceptance tests

1. Submit the same booking command concurrently 20 times with one key; one booking/capacity effect and one stable result exist.
2. Reuse that key with changed time/customer; receive typed conflict and no second effect.
3. Edit a booking using a stale version; preserve input and return current version/domain resolution.
4. Attempt a child resource through another organization path; disclose nothing and deny below routing.
5. Send an unknown field trying to set `status`, `organization_id`, `verified` or `total`; reject and preserve server truth.
6. Verify every non-success response follows approved HTTP/problem semantics and leaks no stack/tenant data.
7. Upgrade server while an older supported Android client has queued commands; commands succeed or recover through explicit compatibility result.
8. Replay duplicate/out-of-order provider callbacks; canonical domain state remains correct.
9. Forward/expire/revoke a scoped link; token cannot broaden access or perform a stale action.
10. Simulate timeouts after commit; client recovers exact result using the original idempotency identity.

## Authoritative references

- [RFC 9110: HTTP Semantics](https://www.rfc-editor.org/rfc/rfc9110.html) defines HTTP methods, status and conditional semantics, including idempotent methods.
- [RFC 9457: Problem Details for HTTP APIs](https://www.rfc-editor.org/rfc/rfc9457.html) defines the structured HTTP problem format and obsoletes RFC 7807.
- [OpenAPI Specification](https://spec.openapis.org/oas/) is the machine-readable HTTP API description baseline; the exact supported release is pinned in the toolchain decision.
- [OWASP API Security Top 10](https://owasp.org/www-project-api-security/) informs API threat coverage but is not a complete verification standard.

## Open decisions

- Exact HTTP resource/action naming after first-release command inventory.
- Backend/runtime/OpenAPI tooling and compatibility-diff tools.
- Session/access-token format and browser cookie versus bearer usage by surface.
- Idempotency retention by action and maximum supported offline period.
- Identifier format/client-generated ID strategy.
- Money transport representation across selected languages/providers.
- Public/business/operations API separation and hostname policy.
- Whether/when to add an API gateway product, GraphQL, gRPC or generated SDKs.
- Realtime protocol is deferred to native messaging/call architecture.


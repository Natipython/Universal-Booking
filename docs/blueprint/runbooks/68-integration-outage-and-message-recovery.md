# Runbook 68 — Integration outage and message recovery

Status: **Draft — provider-specific annexes require approved integrations**  
Owner: Integration operations lead  
Applies to: Telegram, WhatsApp/Meta, Instagram, Facebook, TikTok, SMS/email/push and future official channel/provider adapters

## Core rule

A provider outage must not take down manual booking or cause Ubi to claim a message was delivered, read or acted on without provider evidence. Official APIs only; never recover by scraping, shared passwords, copied sessions, WhatsApp Web or AccessibilityService automation.

Diagram: `diagrams/68-integration-outage-recovery.mmd`.

## Trigger categories

| Category | Examples |
|---|---|
| Provider availability | API/DNS/TLS/regional outage, elevated errors/latency |
| Authentication/authorization | token expired/revoked, merchant asset removed, scope/review changed |
| Ingress/webhook | callbacks missing/delayed/duplicate/invalid signature/schema |
| Egress/send | rate limit, template/policy rejection, timeout/unknown result, delivery failure |
| Capability/policy | message window changed, account restricted, product/version deprecated |
| Ubi adapter | bad deploy/parser/routing/queue/credential/cache issue |
| Tenant-specific | merchant disconnected/misconfigured/ineligible/unpaid provider account |
| Abuse/security | spam complaints, compromised token, forged callback, provider suspension |

## Message/effect states

Use evidence-specific states such as:

```text
draft → queued → sending → provider_accepted → delivered → read
                  ↘ failed
                  ↘ unknown_requires_reconciliation
queued → blocked_by_policy / authorization_required / cancelled
```

Provider support varies. Never synthesize `delivered`/`read`, and never resend an unknown non-idempotent effect blindly.

## Operating modes

- **Normal:** capability verified and health within objective.
- **Degraded:** delays/errors; bounded queue/retry, user-visible warning and fallback.
- **Paused:** unsafe/unauthorized/policy-broken; no new sends, preserve intake/work.
- **Recovering:** probe, reconcile unknowns, fair controlled backlog replay.
- **Disconnected:** credential/asset removed; direct Ubi/manual flow only until authorized reconnect.

Mode is provider + connection + capability scoped; one provider/merchant does not disable all Ubi.

## Initial response

1. Declare incident/case and identify provider/product/app/environment/capability/tenants/time.
2. Verify Ubi telemetry pipeline, adapter release/config/credential status and provider official status independently.
3. Stop unsafe retry/resend/token-refresh storms; preserve queues/webhooks/raw approved evidence.
4. Classify already accepted, definitely failed, unknown and not-attempted effects.
5. Switch affected capability to degraded/paused and surface direct-link/manual/phone fallback.
6. Protect core booking; notification failure never rolls back committed booking.
7. Contact provider through approved support with sanitized correlation/account evidence.
8. Set merchant/support/status update cadence.

## Ingress/webhook recovery

- Validate signature/secret, timestamp/replay, app/environment/asset/connection before parsing.
- Store accepted raw envelope/checksum/headers only per classification/retention, then inbox-deduplicate.
- Quarantine unknown schema/account/tenant rather than route by untrusted payload label.
- Compare provider event/list API where officially available; respect retention/window/rate limits.
- Replay stored/quarantined events idempotently in original provider ordering context where known.
- Detect gaps using cursors/sequence/time/provider queries; gaps remain visible if provider offers no recovery.
- Never manufacture missing message content from notification previews/device scraping.

## Egress/send recovery

- If no provider request occurred: keep queued and retry with backoff/jitter within validity.
- If provider explicitly rejected: mark failed/policy/auth action required; do not loop.
- If provider accepted with ID: store evidence and wait/query supported status; no duplicate send.
- If timeout/connection loss after possible acceptance: mark unknown and reconcile using provider idempotency/reference/query before retry.
- Expired reminders/templates are cancelled, not delivered late without product policy.
- Apply tenant/provider fairness so one backlog cannot starve current critical messages.
- Review customer harm before replaying old conversation/marketing content.

## Authentication and merchant reconnect

- Automatic refresh only under provider rules, locking and token-reuse protections.
- Invalid grant/revoked scope/asset removal pauses capability and tells authorized merchant what changed.
- Reconnect starts a new bound OAuth/official authorization transaction; never ask for password/session.
- Re-discover verified assets/scopes/capabilities after authorization; do not assume old IDs/permissions.
- Old credential is revoked/deleted under retention/evidence rules.
- Merchant chooses routing before sends resume.

## Policy/capability change

1. Capture current official evidence/date and affected app/account/country/use case.
2. Disable unsupported claim/action in capability registry, UI and automation.
3. Determine queued/in-flight/data retention/deletion/merchant notice obligations.
4. Update adapter/contracts/tests/onboarding/terms/cost model.
5. Re-run provider review/prototype before re-enable.

Do not build unofficial workaround to preserve a sales promise.

## Automation and human takeover

- During uncertainty, automation stops high-impact replies/booking/payment claims.
- Human sees provider/Ubi/booking state and whether replying is possible.
- Direct Ubi booking link or phone/manual instruction remains available where safe.
- Duplicate replies are prevented through conversation/action idempotency.
- Messages received while disconnected are not claimed as imported unless official API actually provides them.

## Backlog replay

Before replay:

- provider and token/capability probe pass;
- current policy/window/template still permits each item;
- item remains useful/not superseded/expired/cancelled;
- unknown prior effects reconciled;
- rate/concurrency/cost and merchant fairness limits set;
- monitoring/support ready.

Replay in bounded batches; pause on error/burn/cost/complaint threshold.

## Reconciliation report

For each affected connection/window:

- ingress count received/deduplicated/quarantined/gap/late;
- egress queued/attempted/accepted/delivered-if-supported/failed/unknown/cancelled;
- bookings/conversations/notifications requiring merchant/customer action;
- duplicate/late/misrouted evidence;
- provider cost/rate/account restriction and credential status;
- data deletion/retention changes;
- final unresolved items and owner.

## Release/readiness gates

- Kill switch and per-provider/tenant/capability degraded mode are exercised.
- Signature/replay/routing/schema and unknown-outcome tests pass.
- Queue is durable, bounded, fair and observable; replay cannot duplicate logical effect.
- Official reconnect/revocation/delete and direct-link/manual fallback are tested.
- Provider support/contact, capability evidence, rate/cost limits and merchant script exist.
- No integration is required for core manual booking continuity.


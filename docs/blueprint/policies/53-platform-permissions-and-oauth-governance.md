# Platform permissions and OAuth governance

Status: **Draft — provider-neutral control architecture with time-stamped capability screening**  
Evidence checked: 2026-09-05; provider documentation/dashboard/app-review must be reverified before implementation and release  
Depends on: optional channel architecture, privacy/security, organization authority and vendor governance  
Blocks: provider app creation/review, merchant onboarding, message automation, production tokens and marketing claims

## Direct conclusion

Ubi may connect only through official provider products and authorization. It must never ask a merchant for a WhatsApp/Telegram/Instagram/Facebook/TikTok password, session cookie, device backup, QR-session takeover or unofficial “personal account automation.”

Channel connection remains optional. Manual call/walk-in booking and a Ubi link must remain useful if a provider rejects Ubi, a merchant is ineligible, a capability disappears or the integration is disconnected.

## Current capability screening

This is directional evidence, not an implementation contract.

| Provider/product | Verified direction | Important limitation/gate | Ubi decision |
|---|---|---|---|
| Telegram Bot API | Official [Telegram bot documentation](https://core.telegram.org/bots) says bots use an HTTPS API, users must contact/add a bot before it can start a conversation, and bot tokens give full control; [Bot API](https://core.telegram.org/bots/api/) supports secret webhook headers | Standard bot, group privacy and Telegram Business bot behavior differ; official Bot API does not provide secret chats; token custody/merchant bot ownership model must be chosen | Strong prototype candidate, but distinguish Ubi bot, merchant-owned bot and Business bot connection |
| Telegram Business bots | Telegram documents [Business bots](https://core.telegram.org/api/bots) that can process/answer messages on a user’s behalf | Eligibility/account plan, rights and exact Bot API updates/actions need live prototype and current terms | Capability-gated, not assumed for every merchant |
| WhatsApp Business Platform | Meta’s current developer material describes [unified Embedded Signup v4](https://developers.meta.com/resources/videos/unified-onboarding-whatsapp/) and a changing WhatsApp account model; official platform app/dashboard authorization is required | Business/app verification, country/account/product eligibility, messaging rules, pricing and permissions change; never promise access before onboarding probe | High market-priority research/prototype; direct-link fallback required |
| WhatsApp identity evolution | Meta documents [WhatsApp Usernames/Business Scoped User IDs](https://developers.meta.com/resources/videos/whatsapp-usernames/) and account-model evolution rather than permanent phone-only identity | Ubi cannot use phone number as immutable universal person identity; webhook/parser/API mapping must be versioned | Store provider-scoped IDs plus separately verified contact claims |
| Instagram professional messaging | Meta’s official Postman workspace documents the [Instagram Conversations API](https://www.postman.com/meta/instagram/folder/23987686-6a91368f-1fa8-4614-9ed6-7d1e08c21e62): professional account, access token from a person able to manage messages, webhooks, permissions and Advanced Access for third-party accounts | Requests/history/media/access semantics and permissions are constrained; app review required | Prototype only with eligible test/professional accounts and live capability checks |
| Facebook Page Messenger | Official Messenger/Graph products exist, but exact 2026 Page permissions, messaging windows, review and regional availability were not sufficiently verified in this screening | Must be verified in Meta’s current app dashboard/docs for Ubi’s legal entity and use case | No build promise until evidence package passes |
| TikTok Business Messaging API | TikTok for Business currently documents a [Business Messaging API](https://business-api.tiktok.com/gateway/docs/index?doc_id=1739585600931842&identify_key=c0138ffadd90a955c1f0670a56fe348d1d40680b3c89461e09f78ed26785164b&language=ENGLISH) for business-account direct messaging/automation/thread management | Access, security/privacy review, market/account eligibility, message limits and partner status must be confirmed | Future adapter behind eligibility/cost evidence; not pilot assumption |
| TikTok Data Portability | TikTok’s [Data Portability API](https://developers.tiktok.com/docs/en/data-portability-api-get-started) includes direct-message export scopes | Export/portability is not a live inbox/send API and must not be misrepresented as one | Never use portability scope as a messaging workaround |

Provider changes observed in 2026 reinforce the architecture rule: capability is runtime/versioned data, not a permanent enum or sales promise.

## Connection models

### Merchant-owned provider asset

Preferred for business inboxes. Merchant admin authorizes Ubi to a specific business/Page/professional account/phone/bot asset through the provider’s official flow. Ubi stores external asset identity, scopes/capabilities and protected token—not the password.

### Ubi-owned shared bot/application

May suit a Telegram discovery/booking bot or notification sender. Customers interact with clearly identified Ubi/merchant context. Tenant routing is explicit and one merchant cannot control another’s bot conversations.

### Merchant-created bot token

Telegram may require the merchant to create/own a bot and transfer its token into a secure one-time secret-entry flow. The token never appears in support chat/logs/screens after acceptance; connection verification and rotation/revoke are mandatory. This model has higher onboarding burden and should be usability-tested.

### Provider/solution partner

A licensed/approved provider may simplify WhatsApp/SMS/payment/call onboarding. It becomes a subprocessor/critical dependency with contract, pricing, export/deletion, credential ownership and exit analysis. Partner convenience does not authorize password/session scraping.

## Capability registry

Each provider application/connection records a verified snapshot:

```text
provider/product/API and schema version
country/region and account/business asset type
Ubi app identity/review/access level
merchant asset IDs and verified administrator
granted OAuth scopes/permissions and expiry
read/send/reply/history/media/template/automation/call capabilities
customer-initiation and message-window rules
rate/size/type/retention/edit/delete/delivery semantics
webhook authenticity and replay mechanism
pricing/billing payer and currency
policy/terms/data-use/deletion requirements
last verified time, evidence URL/test and next review
degraded fallback and disable/kill switch
```

UI and automation use this registry. Unsupported action is hidden/explained and server-rejected; it is never attempted optimistically against provider policy.

## Authorization prerequisites

Only an active Ubi member with `integration.connect` capability and exact organization/operation scope may start connection. High-risk providers may require owner/admin role, current step-up and approval.

Before redirect, Ubi explains in plain language:

- which provider/business asset will connect;
- exact message/profile/media actions requested;
- why each permission is needed;
- what Ubi stores and who can access it;
- automation/human takeover behavior;
- provider/Ubi costs and prerequisites;
- how to pause/disconnect/delete;
- which Ubi functionality remains without connection.

Do not request all future permissions “just in case.” New scope requires a new reviewed authorization/change flow.

## OAuth flow

Follow provider requirements and [OAuth Security BCP RFC 9700](https://www.rfc-editor.org/info/rfc9700/).

1. Authorized merchant selects provider and intended operation/queue.
2. Server creates a single-use short-lived authorization transaction bound to Ubi account, membership, organization, provider application, exact redirect and requested scopes.
3. Client opens provider authorization in system browser/approved flow; no embedded credential capture.
4. Use authorization code with PKCE where supported/required, exact redirect matching, CSRF `state`, issuer/mix-up protection and OIDC `nonce` when applicable.
5. Callback validates transaction, state/issuer/code/redirect/expiry/use and current merchant authority.
6. Confidential token exchange occurs server-side; no provider client secret/token enters Android/web logs or URL.
7. Server queries provider to discover authenticated user/business assets/scopes rather than trusting callback labels.
8. Merchant selects exact eligible asset and routing; Ubi verifies their provider authority.
9. Store encrypted token/secret reference and non-secret metadata, test webhook/send/read capability safely and show result.
10. Append connection/audit/notice evidence; enable only capabilities that passed.

Forbidden OAuth patterns:

- resource-owner password grant/password collection;
- implicit token in browser URL;
- wildcard/unvalidated redirect URI;
- reusable/unbound state;
- accepting access token for wrong app/audience/environment;
- long-lived privileged token on client/device;
- using one merchant’s callback transaction for another tenant;
- logging authorization code/token/full callback query.

## Token and credential governance

- Encrypt per-connection provider credentials with KMS/secret-manager design; plaintext only inside owning integration workload memory.
- Store provider, app, environment, external account/asset, scopes, issue/expiry/refresh/key version and last validation.
- Refresh uses locking/idempotency; rotation supports overlap and detects invalidation/reuse where provider supports.
- No ordinary support/merchant/reporting read of tokens.
- Provider app secrets and merchant tokens are different classes/keys.
- Reauthorization is explicit when scope/admin/account changes; do not silently broaden.
- Merchant/admin removal or provider revocation transitions connection to truthful limited/expired state.
- Suspected leak disables actions, revokes upstream, rotates and investigates all affected connections.

## Webhook governance

1. Separate provider/application/environment endpoints or cryptographically bind them.
2. Verify current official signature/secret/token mechanism over required original bytes.
3. Validate timestamp/replay window, method/content type/size and TLS.
4. Durably record receipt ID/digest/provider/application before acknowledging within provider timeout.
5. Deduplicate by provider event/update ID and namespaced fallback only when documented.
6. Resolve connection/business asset/tenant from server records; never accept arbitrary tenant header/field.
7. Validate schema/version and quarantine unknown events/statuses.
8. Process asynchronously/idempotently; no long provider/database operation in webhook response.
9. Preserve edits/deletes/retries/out-of-order semantics as provider evidence.
10. Monitor signature failures, duplicate rate, lag, unknown routing and provider delivery gap.

Telegram’s current Bot API supports a `secret_token` delivered as `X-Telegram-Bot-Api-Secret-Token`; use it plus random endpoint path/TLS and normal controls, not IP lists alone.

## Message identity and linking

- Namespace participant ID by provider, application/connection and external asset/context.
- Store provider-scoped stable identifier even if a phone/display name is supplied.
- Phone/usernames/display names are mutable claims, not authority or merge key.
- Link to a merchant Customer only through reviewed/verified evidence; show provenance and allow unlink/correction.
- Link to a Ubi Account only through explicit verified flow; never because contacts match.
- A user messaging two merchant assets has two merchant relationship contexts; Ubi does not create an invisible cross-business profile.

## Automation and human takeover

- Merchant enables automation per connection/operation and sees its exact capability/limits.
- Automation identifies itself where material and provides human fallback.
- Untrusted message text cannot change permissions, reveal secrets or invoke arbitrary tools.
- Structured extraction remains suggestion until rule/confidence/human review permits.
- Never auto-confirm unavailable capacity, payment, refund, cancellation, ownership or unbounded quote.
- Human takeover atomically pauses automation; queued sends re-check state before dispatch.
- Every automated suggestion/send/escalation retains rule/model/source/version and actor evidence without secret prompt dumps.

## Provider-policy and data-use governance

For each product, legal/privacy/integration owners review current:

- developer/platform/business terms and acceptable use;
- app/business verification and review evidence;
- user initiation, reply windows, templates and categories;
- automation/AI, marketing, prohibited content and opt-out;
- data minimization, storage/retention, transfer, service-provider and training restrictions;
- security, incident and vulnerability obligations;
- user data access/deletion/revocation callbacks or instructions;
- pricing, billing, rate limits and account quality/suspension;
- logo/trademark/UX disclosure rules;
- audit/certification and policy-change notice;
- termination/export/delete and continuing obligations.

Maintain source URL, captured version/date, reviewer, applicability, controls and next review. Legal copy never says Ubi is “official partner” without actual permission.

## Pause, disconnect and deletion

### Pause automation

Stops automated sends immediately while preserving the authorized connection/inbound handling according to merchant choice/provider rules.

### Disconnect

Requires capability, confirmation and impact preview. Ubi:

- disables new provider actions/webhook routing;
- revokes/deauthorizes token upstream where supported;
- deletes protected credential after necessary revoke evidence;
- marks capabilities unavailable and cancels/holds queued sends safely;
- retains or deletes canonical booking/evidence/message content under approved purpose/retention;
- provides provider-side cleanup instructions and verifies where possible;
- updates merchant/customer UI and routes to direct link/manual work.

### Organization/account deletion

Enumerates every connection, provider app asset, token, webhook/subscription, template/sender and retained copy. Deletion is task-driven/reconciled; removing a local token alone does not prove provider deauthorization/deletion.

## Failure and policy-change behavior

- Expired/revoked token: stop actions, preserve drafts, show owner/admin recovery steps.
- Missing scope/capability: disable only affected feature; do not request broad reauth silently.
- Provider outage/rate limit: queue bounded allowed work, show delay, prevent duplicate retry and offer Ubi/manual fallback.
- Unknown send result: reconcile by stable request/provider ID before sending again.
- App review rejected/country unsupported: do not workaround unofficially; revise scope/provider/fallback.
- Policy/price change: impact assessment, merchant notice where material, feature flag and explicit reauthorization/acceptance if required.
- Provider account suspension: show truthful status and evidence/support path; Ubi cannot promise reinstatement.

## Test and approval package per provider

- official product, app, legal entity, business verification and test assets;
- exact scopes/permissions mapped to screens/actions and justification;
- complete OAuth state/PKCE/redirect/issuer/tenant tests;
- webhook signature/replay/duplicate/out-of-order/unknown/rotation tests;
- participant/asset/tenant-routing negative tests;
- supported media/message/edit/delete/history/window/rate semantics;
- automation takeover/block/report/opt-out and prohibited action tests;
- token leak/revoke/reauthorize/admin-removal/disconnect/delete exercise;
- provider outage/unknown result/cost-limit/fallback;
- current terms/privacy/data-use/transfer/subprocessor/retention review;
- app-review demonstration and evidence archive;
- merchant onboarding usability test on low-cost Android and weak network.

## Release gates

- Only official API/provider connection is used; no password, session-cookie, QR takeover or scraping.
- Merchant can use core Ubi without any social connection.
- Runtime capability registry prevents unsupported promises/actions.
- Exact scopes are minimum and pass provider review for Ubi’s real use case/legal entity.
- OAuth/webhook/token/tenant negative tests pass.
- Participant identity cannot auto-merge across providers/merchants/accounts.
- Pause/takeover/disconnect/revoke/delete and provider outage are rehearsed.
- Provider terms, data use, pricing, country eligibility and policy monitoring have owners.
- Customer/merchant disclosures are understandable and accurate.

## Open decisions

- First connected channel: Telegram technical pilot versus WhatsApp market-priority pilot after merchant research.
- Ubi-owned versus merchant-owned Telegram bot/Business bot model.
- Meta legal entity/app/tech-provider/partner path and exact 2026 permissions.
- TikTok Business Messaging market/partner access and commercial value.
- Provider message/content retention and deletion by country/responsibility.
- Whether any provider calling capability is relevant versus future native Ubi calls.
- Provider cost ceilings and merchant pricing/allocation.

# Cache, search, object, device and secret storage

Status: **Draft — product- and provider-neutral storage responsibilities proposed**  
Depends on: platform architecture, offline sync, database architecture, tenant isolation and data catalogue  
Blocks: provider selection, storage schemas, media pipeline, performance/cost tests and production controls

## Purpose

Ubi needs to feel fast on low-cost Android devices and inconsistent East African networks without letting caches, search indexes, files or offline copies become a second source of truth. This document assigns one purpose to each storage class and defines its security, freshness, failure and recovery behavior.

## Executive direction

Start with the smallest dependable storage set:

1. managed PostgreSQL for authoritative operational data, durable jobs and early search;
2. private S3-compatible object storage for files/media, with metadata and ownership in PostgreSQL;
3. platform secret manager/KMS for service secrets and cryptographic key control;
4. encrypted, bounded Android local database/file cache for offline-tolerant business work;
5. CDN/edge cache only for explicitly public, versioned assets/projections;
6. optional small in-process cache for safe immutable/reference data.

Do not require Redis, a dedicated search cluster, data lake or streaming platform for the pilot unless measurement proves PostgreSQL/object storage cannot meet the approved workload. “Telegram-like speed” comes first from small payloads, local-first screens, bounded queries, good indexes, background sync and close-enough regional infrastructure—not from naming more databases.

## Storage responsibility matrix

| Store | Authoritative for | Must never be authoritative for |
|---|---|---|
| PostgreSQL | Booking/capacity, configuration versions, authority, customers, payment evidence, workflows, idempotency/outbox and object metadata | Large media bytes, plaintext secrets, arbitrary telemetry blobs |
| Private object storage | Exact immutable file bytes/versions after verified upload | Ownership, permission, booking state, payment state |
| CDN/public edge | Published public assets/projections within freshness rules | Private customer/business data or authorization decisions |
| Android local database | Approved offline projection, command queue, sync state and drafts | Server-confirmed availability, membership authority or final payment status |
| In-process/shared cache | Recomputable, bounded acceleration | Booking lock, durable queue, final state, sole session/revocation truth |
| PostgreSQL search | Early operational/public search projection/index | Mutation of owner records |
| Future search service | Rebuildable documents and ranking after measured need | Owner data or final ACL decision |
| Secret manager/KMS | Service/provider secrets, key-encryption keys, rotation metadata/control | Ordinary queryable business data |
| Observability store | Minimized logs/traces/metrics | Customer content archive or audit system of record |

Diagram: `diagrams/45-data-storage-topology.mmd`.

## Source-of-truth rules

- One fact has one authoritative owner.
- Cache/search/device/document projections include source version and can be discarded/rebuilt.
- The booking database transaction alone decides capacity and confirmation.
- Search hits and cached IDs are re-authorized before private hydration/action.
- Object possession or URL knowledge grants nothing beyond a valid short-lived access decision.
- A cache outage may make Ubi slower; it must not make availability, permissions or payment truth incorrect.
- Reconciliation detects divergence for every derived store that can affect a user-visible workflow.

## Caching architecture

### L1 process cache

An API process may cache small, stable and non-sensitive reference data such as capability registry, country metadata, schema/feature metadata and immutable published configuration by version.

Requirements:

- bounded memory and item size;
- short/explicit TTL;
- versioned keys;
- no correctness dependency;
- no unbounded tenant-cardinality cache;
- metrics for hit rate, eviction and stale reads;
- safe behavior after deployment/config version change.

### Public edge/CDN cache

Eligible:

- published business/service/location projection;
- transformed public business images;
- static application assets;
- localized public configuration that contains no private data.

Not eligible:

- private booking pages or scoped-link responses;
- customer details, staff calendars, unpublished drafts;
- authorization/session responses;
- payment/provider callbacks;
- media whose audience may change without a versioned public projection.

Public responses use versioned cache keys/ETags and deliberate `Cache-Control`. Publication, unpublication, restriction and slug changes emit invalidation. A short bounded TTL remains the safety net because purge is not perfectly instantaneous.

### Private/API response cache

Avoid general shared caching of private dynamic responses in the first release. If a proven expensive read needs caching:

- the key includes environment, organization, authorization scope/epoch, query/version and locale;
- the result contains no broader fields than its reader may see;
- revocation/invalidation requirements are measurable;
- `private`/`no-store` headers prevent browser/intermediary mistakes where appropriate;
- reauthorization occurs before action;
- a penetration test attempts cross-account and cross-tenant cache poisoning.

### Shared cache/Redis introduction gate

A shared key-value system is added only for a measured use case such as high-volume distributed rate limiting, short-lived presence/typing indicators, hot derived reads or coordination that PostgreSQL/provider infrastructure cannot meet affordably.

It must not be used as:

- booking/capacity lock or final inventory source;
- only copy of commands, notifications or payment work;
- permanent membership/revocation truth;
- substitute for database constraints/transactions;
- uncontrolled message history.

If introduced, use managed TLS/auth, private networking, separate environments, ACLs, bounded TTL/eviction policy, memory alarms and tenant-prefixed versioned keys. Decide explicitly whether persistence is off/on and prove that total cache loss is safe.

## Cache key and invalidation contract

Conceptual private key:

```text
environment:data_class:vN:organization_id:scope_or_authority_epoch:resource_or_query_hash:source_version
```

Keys must not contain raw phone numbers, emails, tokens or sensitive names. Hashing a low-entropy contact alone does not make it safe.

Invalidation is emitted transactionally through outbox after owner state commits. Consumers are idempotent. Each projection/cache documents:

- source owner and event/version;
- maximum tolerated staleness;
- invalidation/rebuild path;
- behavior when event delivery is late/duplicated/out of order;
- security-critical purge target;
- reconciliation query and alarm.

## Search architecture

### Pilot decision: PostgreSQL first

Use tenant-qualified indexed queries plus PostgreSQL full-text/trigram only for proven search fields. This reduces operations and keeps consistency understandable.

Likely pilot searches:

- merchant staff: booking reference, customer display name/contact, service and upcoming work;
- public customer: published business/service/branch/location;
- operations: exact case/reference/provider IDs through controlled tooling.

Search is not permission. Every private result is tenant/scoped at query time and sensitive detail is hydrated through the owner module.

### East African usability and normalization

Search design must be tested with the launch market rather than assuming English spelling and Western addresses.

- preserve original script/display text;
- support configured local languages/scripts and common transliterations;
- normalize Unicode deliberately without destroying meaningful characters;
- make accent/case/punctuation tolerance language aware;
- handle names with different ordering/spelling;
- normalize phone contacts using selected country/numbering rules while retaining display form;
- support landmark/area descriptions where formal street addresses are weak;
- use forgiving prefix/error tolerance for human search but exact normalized matching for controlled identifiers;
- never auto-merge customers solely because fuzzy search thinks they match.

Search relevance needs real anonymized/synthetic local query sets and merchant user testing. A technically fast result that people cannot recognize is not successful.

### Search document rules

Every derived document carries:

```text
document_id and source_type/source_id/source_version
organization_id for private content
scope/visibility fields
publication state for public content
locale/script fields
created/updated/deleted projection positions
minimal searchable/display fields
```

Owner changes emit upsert/delete/tombstone events. The projector rejects older versions. A periodic reconciler compares counts/checkpoints/samples with PostgreSQL. Rebuild occurs into a new version/index and changes the read alias only after parity checks.

### Dedicated search service gate

Introduce a managed search engine only when evidence shows PostgreSQL cannot meet approved latency/relevance/index size/write load or when capabilities such as robust multilingual typo tolerance materially improve conversion/work speed.

Before selection measure:

- expected documents and update/query rate;
- launch languages/scripts and tokenizer quality;
- tenant/field ACL filtering correctness/performance;
- region/data location and subprocessor implications;
- outage/rebuild time and operational expertise;
- minimum monthly cost and growth curve;
- export/deletion propagation and backup expectations;
- safe fallback search behavior.

An external engine remains disposable and cannot accept direct client mutation.

## Object and media storage

### Ownership split

PostgreSQL `media.objects` owns metadata and access relationship:

- immutable object ID, organization/owner/purpose/classification;
- storage provider/bucket/key/version reference;
- expected and verified byte size/content type/checksum;
- upload/scan/process/publication/lifecycle state;
- encryption/key version where relevant;
- retention/legal hold/deletion state;
- creator/source/correlation and timestamps.

Private object storage owns the exact bytes. Business/domain rows link through controlled `media.object_links`; they do not store public URLs as permanent authority.

### Object key pattern

Use opaque server-generated keys, conceptually:

```text
environment/private/organization_id/object_id/version/randomized_name
environment/public/publication_version/object_id/variant
```

The tenant segment aids operations; it is not authorization. Original filenames are metadata after sanitization and are not trusted as keys/paths/content types.

### Upload flow

1. Client requests an upload session for exact tenant, purpose, owner, media class, size and type.
2. Server authorizes and creates pending metadata plus a short-lived, size/type-limited upload grant.
3. Client uploads directly/resumably where supported, useful on unstable networks.
4. Storage event/client completion queues verification.
5. Worker verifies expected key, tenant, size, checksum/signature and actual content type.
6. Malware/unsafe-content and decompression/image processing controls run as risk requires.
7. Safe variants are created without overwriting originals; metadata advances state.
8. Domain owner may link/publish only an approved ready object.
9. Orphan pending uploads expire and are reconciled/deleted.

The API does not proxy large file bytes by default. Upload retry resumes or safely creates a new version; it never silently attaches someone else’s object.

### Download flow

- Server authorizes current resource/field/purpose.
- It returns a short-lived audience/action/object/version-scoped URL or streams through a controlled gateway for higher-risk content.
- Response uses safe content disposition/type, anti-sniffing and cache policy.
- Revocation-sensitive content uses shorter expiry and no public CDN cache.
- Access to classified content is auditable where required.

### Media processing isolation

- Processors read/write only designated input/output prefixes/buckets.
- Untrusted formats run in constrained workers with CPU/memory/time/size limits.
- Metadata/EXIF is removed from public variants unless explicitly needed; location metadata is protected.
- Image dimensions/pixels, archive expansion and document parsing have hard limits.
- A failed scan/transform cannot become published through fallback.
- Provider callback/event cannot choose an arbitrary storage key.

### Retention, deletion and recovery

- Lifecycle policy is driven by approved retention class, domain state and legal hold—not storage age alone.
- Deletion is a recorded idempotent workflow across original, variants, CDN and backup expiry.
- Object versioning/soft-delete may protect against accidental deletion, but retention/privacy cost and final purge must be defined.
- Cross-region replication is not automatic; select it from RPO/RTO, data-location law and cost.
- Database backup without object manifest/checksum reconciliation is not a complete restore.
- Periodic restore drills sample file readability and metadata linkage.

## Android device storage

The Business app needs local speed and limited offline work, but a lost low-cost phone is a realistic threat.

### May store, within release policy

- bounded upcoming bookings/work projection;
- service/resource reference projection needed for offline screens;
- drafts and idempotent pending commands;
- sync cursor, source versions and tombstones;
- thumbnails/media explicitly needed and expiring;
- minimal display contact necessary for assigned work.

### Must not store casually

- integration/provider secrets or server credentials;
- full organization customer/message/payment history;
- reusable scoped-link secrets;
- raw OTP/password/authenticator values;
- permanent authoritative availability/payment/permission state;
- unnecessary high-risk identity documents/media.

### Controls

- app-private sandbox and encrypted database/files with keys protected by Android Keystore where supported;
- no sensitive data in logs, notifications, clipboard, screenshots/recents where risk demands controls;
- account/organization partitioning inside local schema;
- auto-lock and step-up for risky actions;
- bounded time/row/byte cache with deterministic eviction;
- remote/session revocation propagates local scope purge;
- uninstall/logout/ownership transfer behavior tested;
- backup/export flags configured so OS/cloud backups do not leak unsupported app data;
- rooted/compromised device limitations stated honestly.

Cached local data may be shown with clear last-synced/pending state. Any capacity/payment/authority-sensitive action returns final server status.

## Secret storage and key management

### Secret classes

| Secret | Storage direction |
|---|---|
| Service/database credentials | Managed secret manager, short-lived workload identity where available |
| OAuth/provider tokens | Envelope-encrypted ciphertext or per-connection managed secret reference; metadata only in ordinary tables |
| Signing/encryption key material | Managed KMS/HSM-backed service; application receives use permission, not raw master key |
| Scoped-link/API token | Random secret shown/transmitted as needed; verifier/hash + metadata in PostgreSQL |
| Password authenticator | Approved slow password hash with per-record salt/parameters; never decryptable plaintext |
| OTP/recovery code | Short-lived verifier/hash, attempt/expiry controls; never logged |
| Mobile app | No embedded privileged backend secret; platform attestation is signal, not sole trust |

### Envelope encryption for high-cardinality tenant credentials

Storing every provider token as an independent secret-manager object may become expensive/slow. A reviewed envelope design may store ciphertext in a restricted table while KMS protects versioned key-encryption keys.

Requirements:

- authenticated encryption with immutable associated context such as environment/provider/connection/credential ID;
- unique nonce generated by approved library;
- key version retained for decrypt/rotation;
- decrypt only in the owning integration workload;
- plaintext exists only briefly in memory and never in logs/errors/traces;
- no general support/reporting/database role can decrypt;
- rotation and compromised-provider revocation are independently supported;
- backup restore retains access to required historical key versions until data rewrap/deletion is verified.

Exact algorithms/libraries/KMS are approved through security ADR, not improvised in application code.

### Secret lifecycle

Every secret has owner, environment, purpose, issuer/provider, scope, creation, last use where available, rotation/expiry, revocation and incident procedure.

- Production/staging/development keys are separate.
- Secrets never enter source control, images, tickets, chat, analytics or crash reports.
- CI uses workload federation/short-lived credentials where available.
- Access is least privilege and auditable.
- Rotation supports overlapping old/new validity when providers require safe cutover.
- Provider disconnection revokes upstream token and deletes ciphertext after evidence/retention rules.
- Secret scanning runs in pre-commit/CI and repository/history incident response is documented.

## Data encryption direction

- TLS for all network paths, with certificate validation and approved modern settings.
- Managed storage encryption at rest is baseline, not a substitute for authorization.
- Application/field encryption is selected for data whose threat model requires separation from database/support access.
- Encryption keys are environment/purpose separated and versioned.
- Searchable encryption/tokenization is introduced only for exact approved lookup needs; it leaks patterns and needs a documented threat model.
- End-to-end encryption is not claimed for business workflows where Ubi servers must route, search, support or automate content. Marketing language must match reality.

## Failure and degradation behavior

| Failure | Required behavior |
|---|---|
| Process/shared cache unavailable | Fall back to bounded PostgreSQL query; shed optional load rather than return wrong truth |
| Public CDN stale/unavailable | Origin/version fallback where safe; never expose unpublished object |
| Search projection late/down | Show explicit limited/direct lookup or unavailable state; booking owner data remains safe |
| Object store/upload slow | Resume/retry and keep pending state; do not confirm attachment/publication prematurely |
| Media scanner/processor down | Queue pending; do not bypass safety gate |
| Secret manager/KMS unavailable | Fail protected provider action safely; do not use stale secret beyond policy |
| Android offline | Show cached last-sync status, queue approved commands, block/qualify server-truth actions |
| Projection event missed | Reconciliation/backfill/rebuild from owner database and checkpoints |

## Provisional performance and freshness budgets

These are design hypotheses to validate in `62-performance-capacity-and-cost-model.md`, not promises:

- cached Android screen paints from local data without waiting for network;
- critical interactive server responses use bounded/indexed database reads and small payloads;
- public publication/unpublication reaches normal edge/search paths within a defined short target, with urgent restriction purge faster;
- private operational search freshness is visible and normally seconds, not hours;
- upload is resumable for realistic mobile network interruption and progress survives screen/process restart where feasible;
- no first screen requires original full-resolution media;
- cache/search/object failure never permits double booking or permission bypass.

Test with affordable Android hardware, constrained CPU/memory/storage, packet loss, high latency, low bandwidth and intermittent connectivity—not only a developer laptop on Wi-Fi.

## Cost controls

- Image/media size, variants, egress and retention quotas are product decisions.
- Prefer thumbnails and adaptive sizes; never download originals for list screens.
- Bound cache/search document cardinality and payloads.
- Lifecycle orphan uploads and obsolete non-held variants.
- Track storage, request, transformation, CDN egress, KMS and search cost per active organization/workflow.
- Provider selection models an East African latency region plus legal location and egress, not headline storage price alone.
- Adding a managed datastore requires an on-call/backup/security/cost owner.

## Validation and release gates

- Every field/object/projection is assigned one authoritative owner and classification.
- Cache/search/object/device schemas include tenant and source version.
- Cross-tenant cache poisoning/search/object URL negative tests pass.
- Search relevance/usability is tested with launch-language names, numbers, landmarks and misspellings.
- Upload/download limits, malware/parser isolation and orphan reconciliation pass adversarial tests.
- Device storage is inspected on logout/revocation/loss scenarios and low-end hardware.
- Secret inventory, rotation and simulated leak/revocation exercises pass.
- Projection rebuild from a clean store is rehearsed and timed.
- Database plus object restore is reconciled by manifest/checksum.
- Failure injection proves cache/search/storage/KMS degradation does not corrupt owner truth.
- Performance and recurring cost are measured on the candidate region/provider before commitment.

## Open decisions

- Launch country/city/languages and realistic search corpus.
- Managed PostgreSQL/object/KMS region/provider and data location.
- Exact Android local database/encryption library and supported OS/device floor.
- Whether pilot actually needs a shared cache beyond provider/CDN controls.
- PostgreSQL trigram/full-text configuration and protected contact lookup method.
- Media types, maximum sizes, transformations, malware tooling and retention.
- Envelope encryption implementation versus managed per-connection secrets at pilot volume.
- CDN purge/freshness targets for ordinary unpublish versus emergency restriction.
- Object versioning/replication and tested RPO/RTO based on cost/risk.

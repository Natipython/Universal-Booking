# Payments, deposits, refunds and reconciliation workflows

Status: **Draft — provider and jurisdiction neutral**  
Depends on: actors, permissions, booking/custom workflows and quality attributes  
Blocks: payment architecture, data model, security/privacy/legal review, business/client experience and operations

## Purpose

Ubi must represent money truthfully in a market where customers may pay cash, transfer/mobile money directly to a business, show a receipt/reference, or use an integrated provider. A customer saying “I paid,” a provider accepting a request, a provider confirming a transaction, money settling to a business and a booking becoming confirmed are different events.

This specification defines the lifecycle and safeguards without selecting a provider or causing Ubi to hold/route money. The launch-country payment role, licensing, tax, settlement and consumer obligations require provider due diligence and licensed legal/accounting advice before real-money release.

## Recommended initial commercial boundary

Ubi SHOULD begin as booking software that records payment requirements and evidence while the merchant receives payment through approved methods. If integrated payment is later enabled, funds should flow through a licensed provider structure approved for Ubi’s exact role and country.

Ubi MUST NOT accidentally become a wallet, escrow service, marketplace settlement operator, money transmitter or merchant of record through interface wording or implementation. Any movement from “software record/provider integration” toward custody, split settlement, stored value or payout requires a new legal, financial, risk and operations approval package.

## Money-domain separation

| Concept | Meaning |
|---|---|
| Price component | Exact charge/discount/tax/fee/deposit line with currency, quantity/rate and reason. |
| Price snapshot | Immutable version of customer-visible components/total used for a request, quote or booking commitment. |
| Payment obligation | Amount/currency a party is expected to pay by a condition/deadline for a booking/quote. |
| Payment intent | Ubi’s request to satisfy one or more obligations using one approved method/provider. |
| Payment attempt | One customer/provider interaction under an intent. Several attempts may fail; they do not multiply the amount owed. |
| Customer-reported payment | A claim/reference/receipt supplied by customer or staff, not yet independently verified. |
| Provider transaction | Provider-identified payment event/status with authenticated evidence. |
| Verified payment allocation | Exact verified amount applied to one or more obligations. |
| Settlement | Provider/bank evidence that funds were included in a merchant payout/accounting settlement. |
| Refund | A controlled return request/attempt/verified outcome against eligible verified payment. |
| Reversal/chargeback | Provider/network reversal or dispute that may occur after apparent payment success. |
| Reconciliation case | A discrepancy requiring automated or human resolution without fabricating money/booking state. |
| Security/damage deposit | Refundable guarantee distinct from a booking prepayment; future-gated unless explicitly designed. |

Booking, payment, refund, settlement and dispute state machines remain separate and link by explicit policy/event.

## Exact amount and currency rules

1. Every amount has an ISO 4217 currency and exact integer-minor-unit or approved decimal representation.
2. Binary floating point is prohibited for monetary calculations.
3. Currency precision, rounding, fees, tax inclusion and minimum/maximum values are server-controlled and versioned.
4. One payment attempt has one expected currency; conversion is a separate provider/financial event and cannot be inferred from display.
5. Price snapshots itemize subtotal, discounts, tax/fees, required deposit/prepayment, paid/verified amount and remaining balance.
6. “Starting from,” estimate, quotation and fixed total must have distinct labels.
7. Client-calculated totals are previews. The server recomputes and returns authoritative totals before commitment.
8. A changed price creates a new version/amendment; issued/accepted evidence is never overwritten.

## Payment methods

### No payment required

The booking may confirm through its availability/approval policy without a payment obligation.

### Cash/pay at location

- A future obligation may be recorded with method `cash_at_location`.
- Selecting cash does not mark it paid.
- Authorized staff records collection after receipt with amount, time, collector and optional receipt/reference.
- Cash cannot be independently provider-verified, so the state/evidence source must say staff-recorded.
- Cash reconciliation follows merchant shift/register procedures if that product scope is approved.

### Customer-reported transfer/mobile-money payment

- Customer/staff may submit a bounded reference, amount, time and approved receipt image.
- The claim enters `reported_unverified`.
- Authorized business staff or an integrated provider verifies/rejects it under policy.
- Receipt images/references are sensitive evidence with restricted access and retention.
- Visual resemblance, message text or customer statement alone cannot create `provider_verified`.

### Integrated provider payment

- Ubi creates an idempotent payment intent with exact amount/currency and stable references.
- Sensitive entry SHOULD occur in provider-hosted/approved components so Ubi minimizes payment credential scope.
- Customer/provider redirects are navigation evidence only.
- Authenticated provider API/webhook status is reconciled before Ubi marks provider verification.
- Provider initiation/processing/success may be asynchronous, reversed or duplicated; the UI reflects that reality.

## Payment-intent lifecycle

```text
draft
→ ready
→ initiated
→ processing
→ partially_verified
→ verified

ready/initiated/processing → failed
ready/initiated/processing → expired
ready/initiated/processing → cancelled
processing/verified → exception
verified → partially_refunded → refunded
verified/partially_refunded/refunded → reversed_or_disputed
```

This is a Ubi intent summary. Individual attempts, provider transactions, allocations, refunds, settlements and disputes retain their own evidence and lifecycle.

### Logical state meanings

| State | Meaning |
|---|---|
| draft | Obligation exists but payment initiation is not ready/published. |
| ready | Exact payable amount/method/deadline is valid for initiation. |
| initiated | Customer/provider interaction began; no provider verification yet. |
| processing | Provider evidence says processing/unknown-final. |
| partially_verified | Verified allocations are less than the required amount. |
| verified | Required amount is independently verified and allocated, subject to later reversal. |
| failed | Latest/terminal eligible attempt failed; another attempt may be permitted. |
| expired | Deadline/intent validity ended; late provider success still requires reconciliation. |
| cancelled | Ubi/authorized actor stopped further attempts before final verification. |
| exception | Evidence conflicts, amount/currency/reference is wrong/duplicate, or human review is required. |
| refunded | Approved refund total equals the eligible verified amount according to policy. |
| reversed_or_disputed | Provider/network challenged or removed value after an earlier state. |

“Paid” MAY be a user-facing summary only when the exact underlying meaning and source are available; “reported,” “being checked,” “verified,” “settled,” and “refunded” must not collapse into one green state.

## Payment-conditioned booking workflow

### Preconditions

- Booking/service configuration has a versioned payment policy.
- Exact price/deposit obligation and expiry/hold rules are known.
- The method/provider supports the currency, amount, customer/merchant eligibility and country.
- Availability policy defines whether a hold exists while payment is attempted.

### Workflow

1. Customer/staff selects a service/time or accepts an exact quotation version.
2. The server recomputes price and authoritative availability.
3. Ubi creates a booking request/conditional record and, if configured, an expiring capacity hold.
4. Ubi creates/reuses an idempotent payment intent for the exact obligation.
5. The customer sees amount, currency, purpose, business recipient/provider, method, expiry, booking condition and failure/recovery before starting.
6. Payment is initiated through the approved method.
7. Redirect/client result may update the experience to `being checked`; it cannot verify payment.
8. Ubi receives/fetches authenticated provider evidence or authorized manual verification.
9. Reconciliation validates provider account, transaction ID, amount, currency, status, intended merchant/obligation, duplicates and prior reversals.
10. A verified allocation triggers an idempotent booking policy command.
11. Booking confirmation revalidates hold/availability, quote/policy version and payment condition transactionally.
12. If payment verifies but capacity/booking cannot confirm, Ubi creates an urgent reconciliation/refund/rebooking case; it never hides the money event.
13. The customer/business receives separate truthful payment and booking receipts/status.

## Holds, expiry and late payment

- A payment hold reserves only the configured capacity/resources and has an authoritative server expiry.
- The UI shows the hold deadline and what happens after expiry.
- Payment initiation does not extend a hold unless an explicit idempotent policy does so.
- On expiry, availability releases independently; the payment intent may expire but later provider events are still ingested.
- A provider success after hold/intent expiry enters `late_payment` reconciliation. It does not recreate or confirm the expired booking automatically.
- Staff options may include offer alternative capacity, authorized new booking allocation or refund, each with customer communication and audit.
- Provider delay/retry behaviour must be included when setting hold duration; an unrealistically short hold causes paid-without-booking incidents.

## Manual evidence and verification workflow

1. Customer or authorized staff reports payment against a specific obligation.
2. Ubi checks basic format/amount/currency/reference duplication but labels it unverified.
3. Evidence enters a tenant-scoped verification queue with deadline and assigned capability.
4. Verifier compares an approved source (provider portal/API/account evidence), not merely the uploaded screenshot.
5. Verifier records `verified`, `partially verified`, `rejected`, `cannot verify` or `escalated` with reason and source/time.
6. A second approval is required when amount/risk threshold or self-verification conflict applies.
7. Verified allocation invokes the same booking policy command as integrated verification.
8. Customer/business is notified of the real outcome.

The person who reports/collects manual payment SHOULD NOT verify their own high-risk transaction where staffing permits. Small one-person businesses may require a disclosed reduced-control mode and periodic reconciliation rather than fake separation of duties.

## Webhook and provider event processing

For every provider:

1. Receive raw bytes and required headers through the provider-specific endpoint.
2. Authenticate signature/source according to current official provider rules before trusting content.
3. Resolve environment, provider account and Ubi connection/merchant.
4. Persist a bounded immutable receipt/deduplication record before acknowledging when architecture permits.
5. Deduplicate by provider event ID and logical transaction transition.
6. Parse into a versioned provider adapter contract; unknown event types are retained safely/observed, not guessed.
7. Acquire concurrency control for the Ubi intent/transaction/allocation.
8. Validate legal state transition, amount, currency, provider merchant, reference and prior outcomes.
9. Commit provider evidence, normalized event, allocations/state changes and outbox events atomically.
10. Acknowledge promptly; downstream notification/booking work runs idempotently.

Webhook ordering is not assumed. Polling/reconciliation may be required when events are missing or ambiguous. Provider secrets/signatures and sensitive raw payload fields are excluded from ordinary logs.

## Refund workflow

### Preconditions

- An eligible verified payment/allocation exists.
- Booking cancellation/amendment policy has calculated the candidate refundable amount.
- The actor has `refund.request`; amount/risk may require separate `refund.approve`.
- Provider/method supports the proposed path or an approved manual-return process exists.

### Workflow

1. Ubi shows verified paid amount, prior refunds/reversals, policy calculation, provider constraints and remaining refundable maximum.
2. Requester selects amount/reason and confirms customer destination/details only through an approved safe method.
3. Approval binds exact payment/allocation, amount, currency, reason and expiry.
4. Ubi creates an immutable idempotent refund intent.
5. Executor/provider attempts the refund; client success alone is not final.
6. Authenticated provider/manual evidence marks processing, verified, failed, unknown or exception.
7. Allocation and booking financial summary update from verified outcome.
8. Customer receives truthful status and reference/receipt where allowed.
9. Settlement reconciliation later confirms financial inclusion; discrepancy reopens an exception without rewriting prior evidence.

Cancellation and refund are separate: a booking may be cancelled without a refund, and a refund may be pending/failed after cancellation.

## Reversals, disputes and chargebacks

- Later provider reversal/dispute events are always ingested even after booking completion or refund.
- Ubi creates a case linked to the original transaction/allocation/booking and preserves all earlier evidence.
- The booking is not silently deleted or automatically marked customer-at-fault.
- Access to evidence and response actions is capability- and deadline-controlled.
- Merchant/customer communication states that the payment is disputed/reversed, not that fraud is proven.
- Accounting/settlement effects remain separate from service fulfilment decisions.
- Evidence export, representment or liability allocation is provider/country/contract specific and requires approved operations/legal procedure.

## Reconciliation model

### Layers

1. **Real-time transaction reconciliation** — Ubi intent/expected amount versus provider event/transaction.
2. **Daily operational reconciliation** — provider transaction list versus Ubi verified/failed/unknown records.
3. **Settlement reconciliation** — verified provider transactions/refunds/fees/reversals versus payout/settlement evidence.
4. **Merchant/manual reconciliation** — staff-recorded cash/transfers versus merchant-controlled source evidence.

### Exception types

- success at provider, absent/unknown in Ubi;
- success in Ubi, not supported by provider evidence;
- wrong amount/currency/merchant/account;
- duplicate transaction or reused reference;
- partial/overpayment;
- late payment after expiry;
- payment verified but booking failed;
- booking confirmed while required payment unresolved;
- refund/reversal/dispute mismatch;
- settlement missing/short/fee mismatch;
- orphan provider transaction with no valid obligation.

### Case controls

Each case has type, severity, amount/currency, provider evidence, linked records, owner, due date, permitted resolution commands, customer/merchant communication, decision/reason and closure review. Staff cannot resolve a case by directly editing balances or deleting provider evidence.

## Booking financial summary

The booking may expose a derived summary such as:

- total agreed price;
- amount required now;
- verified allocated payment;
- reported/unverified amount;
- verified refunds;
- reversed/disputed amount;
- balance due;
- next payment deadline/method;
- exception indicator.

The summary is recalculated from immutable price versions, obligations and verified allocations. It is not a manually editable balance column.

## Offline behaviour

### Allowed offline

- view previously synchronized price/payment summary with a visible “as of” time;
- record a customer-reported reference/receipt or cash-collection draft as a pending local command;
- prepare a refund request draft;
- show the user that server/provider verification requires connection.

### Online-only by default

- provider payment initiation when no offline provider mechanism exists;
- provider verification/reconciliation;
- final confirmation conditioned on payment;
- refund execution/approval above risk policy;
- settlement/dispute action;
- changing payment provider/account.

An offline record MUST say “Saved on this phone—not verified” and can be rejected/conflicted after server authorization, duplicate and state checks. It can never display provider-verified green success.

## Permissions, approvals and separation of duties

Distinct capabilities include:

- view booking financial summary;
- view restricted payment evidence/reference;
- report/record payment;
- verify/reject manual evidence;
- configure method/deposit/policy;
- request refund;
- approve refund by threshold;
- execute/retry refund;
- view settlement/fees;
- resolve reconciliation exceptions;
- connect/change payment provider/account;
- export financial records;
- privileged Ubi support investigation.

Changing the merchant settlement/provider account is ownership-risk and requires recent strong authentication, approval/notice, cooling/risk checks and audit. Ordinary reception staff cannot perform it.

## Security and privacy controls

- Use provider-hosted tokenized flows where practical; Ubi should not store raw card/mobile-money PIN or reusable financial credentials.
- Payment references/receipts, payer contact and provider payloads are confidential and field-level access/redaction is required.
- Secrets and signing keys use managed secret/key controls, rotation and environment isolation.
- All create/confirm/refund/provider-event operations require idempotency and replay protection.
- Logs and analytics use Ubi IDs and coarse outcome/latency, not full references, tokens or receipt media.
- Receipt uploads require media type/size validation, malware/content safeguards, protected object access and retention/deletion policy.
- Risk signals must not become hidden discriminatory decisions; denials/reviews need explainable reasons and redress where applicable.
- Payment-security standard scope, provider responsibilities and compliance evidence must be determined before selecting an integrated method.

## Customer and business receipts

Receipts/status views MUST clearly show:

- business/legal or trading identity as required;
- booking/request reference and payment purpose;
- amount and currency;
- payment method/source category without exposing secrets;
- reported, processing, verified, failed, refunded or disputed status;
- provider/Ubi reference appropriate for support;
- date/time/timezone;
- whether the booking itself is confirmed;
- balance/next action;
- help/dispute path and applicable policy/version.

Ubi cannot label a document a legally valid tax receipt/invoice until launch-country requirements and responsible issuer are approved.

## Failure and exception catalogue

| Failure | Required behaviour |
|---|---|
| Customer closes provider page | Retain/recover intent; query/wait for provider truth; show not verified. |
| Redirect says success, webhook/API absent | Remain processing/being checked; reconcile; never verify from redirect. |
| Duplicate tap/retry/webhook | Return/reuse same logical intent/effect; no duplicate charge/allocation/booking. |
| Provider event out of order | Apply legal transition/evidence rules; retain event; never regress verified state blindly. |
| Amount/currency mismatch | Exception queue; no verified allocation or false booking confirmation. |
| Partial payment | Allocate exact verified amount; apply configured policy; explain remaining amount. |
| Overpayment | Verify actual amount, prevent excess allocation, open refund/credit decision case; no hidden balance. |
| Paid after hold expired | Late-payment case with rebook/refund resolution; no automatic unavailable booking. |
| Payment verifies, booking commit fails | Urgent paid-without-booking case and customer acknowledgement; capacity/refund/rebooking path. |
| Booking confirms, payment later reverses | Preserve booking history, mark financial exception/dispute and invoke merchant policy; do not erase fulfilment. |
| Refund provider timeout | Keep processing/unknown and reconcile with same idempotency key; do not submit a new refund blindly. |
| Staff uploads fake/duplicate receipt | Unverified/rejected with audit and abuse escalation; no provider-verified state. |
| Provider account disconnected | Stop new initiation, retain/reconcile existing transactions through approved access/export path, show business action needed. |

## Invariants

1. A client redirect, screenshot, customer statement or message never establishes provider-verified payment.
2. One provider transaction value cannot be allocated more than once beyond its verified available amount.
3. Refunds cannot exceed eligible verified value after prior refunds/reversals unless an explicit separately accounted adjustment is approved.
4. Amount/currency/provider merchant must match before automatic allocation.
5. Booking confirmation and payment verification remain separate events connected by versioned policy.
6. Every provider event, attempt, allocation, refund and reconciliation decision is attributable and immutable in evidence.
7. Retrying a logical command/event cannot create duplicate money movement or booking effects.
8. Expired/late payment never silently consumes newly available capacity.
9. A price/quote/policy change does not rewrite an accepted snapshot.
10. Offline/manual evidence cannot masquerade as independently verified.
11. Settlement is not inferred from payment verification.
12. Staff and support cannot repair money history with direct database edits.

## Audit requirements

Audit includes:

- payment configuration/provider account changes;
- price/deposit/policy version publication;
- intent/attempt creation and idempotency key scope;
- customer/staff manual report and evidence access;
- manual verification/rejection and approval;
- provider event receipt/authentication outcome/normalized transition;
- allocations and booking policy results;
- refund request/approval/execution/outcome;
- reconciliation assignment/decision/closure;
- reversal/dispute/settlement updates;
- privileged support/export actions.

Audit must protect secrets and minimize personal/payment content while retaining sufficient evidence.

## Release sequence

1. **Foundation:** exact prices/obligations, cash/pay-later, explicit `unpaid/reported/verified` source semantics, booking financial summary and audit.
2. **Pilot candidate:** merchant-recorded cash and customer-reported transfer reference/evidence with controlled manual verification, if validated by the pilot.
3. **Expansion:** one licensed integrated payment method/provider after country, role, security and reconciliation approval.
4. **Expansion:** provider refunds, daily transaction and settlement reconciliation, disputes/reversals.
5. **Future gated:** security/damage deposits, split payouts, stored balance/credits, multi-currency conversion or any Ubi custody/escrow-like role.

The first pilot can create value without integrated payments; integrating payment before reconciliation/support is dangerous.

## Acceptance scenarios

1. Staff marks a mobile-money transfer as reported; every screen says unverified until independent evidence is recorded.
2. Five identical provider webhooks create one verified transaction allocation and one booking-policy command.
3. A customer refreshes/taps Pay twice; Ubi reuses the intent/idempotency scope and does not request duplicate payment.
4. Provider redirect says success but authenticated status is processing; the customer sees “being checked,” not “paid/confirmed.”
5. Payment verifies after the capacity hold expired; Ubi opens a paid-without-booking resolution rather than confirming an unavailable slot.
6. A refund times out; retry uses the same provider/idempotency reference and reconciles before another attempt.
7. A later reversal preserves the original verified evidence and booking history while creating a financial exception.
8. A receptionist can record cash collection but cannot change settlement accounts or approve a high-value refund.
9. Offline cash/evidence entry survives process death but is labelled local/unverified and reauthorized on synchronization.
10. Reports distinguish quoted total, verified payment, reported evidence, refund, settlement and balance.

## Open decisions

- Launch-country operator/payment role and licensed legal/accounting analysis.
- Pilot payment behaviour: pay later, cash, direct transfer evidence, integrated provider or combination.
- Merchant versus Ubi responsibilities for verification, refund, disputes, receipts and customer support.
- First provider, currencies, merchant/customer eligibility, fees, settlement timing and API/webhook reliability.
- Deposit/hold duration and late-payment policy for the first booking archetype.
- Manual verification thresholds and feasible separation of duties for one-person businesses.
- Financial record/evidence/receipt retention and tax-document obligations.
- Whether Ubi charges subscription, transaction fee or other pricing and how provider costs appear.

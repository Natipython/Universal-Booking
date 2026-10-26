# Runbook 69 — Payment reconciliation and dispute response

Status: **Draft — non-custodial Ethiopia-first model; provider and counsel approval pending**  
Owners: Payment/finance operations and payment engineering  
Applies to: customer-reported payment, licensed-provider payment initiation/status, cash/pay-at-location records, refunds/reversals where supported, discrepancies and merchant/customer disputes

## Boundary

Ubi is proposed as booking software and a technical integration layer. A licensed provider settles customer funds directly to the merchant/provider-designated account. Ubi does not hold, pool, store, split or settle customer funds.

Any wallet, pooled collection, split settlement, gateway/operator or custody behavior is outside this runbook and requires a new legal/regulatory/architecture decision.

Diagram: `diagrams/69-payment-reconciliation-flow.mmd`.

## Non-negotiable states

```text
not_required
reported_by_customer_or_merchant
provider_initiated
provider_processing
provider_verified
provider_failed
unknown_requires_reconciliation
reversal_pending / reversed
refund_requested / refund_processing / refunded / refund_failed
cash_expected / cash_recorded_by_merchant
```

- Screenshot, SMS text, chat message or merchant assertion may support `reported`; it never directly creates `provider_verified`.
- Provider verification includes provider, merchant/payee, amount, currency, stable reference and authoritative status evidence.
- Booking, fulfillment and payment lifecycles remain separate.

## Invoke this runbook for

- provider callback/polling discrepancy or missing/duplicate/out-of-order event;
- Ubi/provider/merchant totals or references do not reconcile;
- customer claims paid but provider is absent/failed/different;
- verified payment mapped to wrong booking/merchant/amount/currency;
- duplicate initiation/charge concern;
- refund/reversal/chargeback/dispute request;
- provider outage, settlement concern or account restriction;
- suspected forged evidence/fraud or unauthorized refund;
- cash/deposit record disagreement affecting booking.

False verification or wrong-tenant financial exposure is SEV-0/SEV-1 and invokes runbooks 66/67.

## Roles and separation

- Payment operations owns case, provider/merchant/customer evidence and reconciliation.
- Payment engineering owns adapter/idempotency/webhook/query defects.
- Finance owner approves financial corrections/fees/write-offs within policy.
- Merchant owns service/cancellation/refund decision where contract allocates it.
- Security/trust handles fraud/account compromise.
- Counsel/support handles regulated/consumer/legal communications.

No person should request and solely approve a high-risk refund/manual verification/export where staffing allows separation.

## Immediate response

1. Open case/incident; verify requester and merchant/booking/provider scope.
2. Freeze unsafe automated transition/refund/retry for affected payment/reference/connection.
3. Preserve Ubi ledger/audit/idempotency/webhook/provider request-response evidence.
4. Classify current Ubi state and what evidence established it.
5. Query/contact provider through official authenticated path; do not rely on customer screenshot alone.
6. Prevent duplicate customer action: show processing/unknown and safe next step, not “try again” blindly.
7. Check booking hold/confirmation/customer communication separately.
8. Set owner, next update and any time-sensitive service/refund obligation.

## Reconciliation dimensions

Match on:

- Ubi organization/operation/booking/payment attempt;
- provider/account/payee/merchant identifier;
- provider transaction/idempotency/reference and parent/refund/reversal reference;
- exact amount/currency/fee where provided;
- initiation/authorization/settlement/reversal timestamps and timezone/source;
- current provider status and event history;
- callback signature/app/environment/connection;
- booking/payment policy snapshot and responsible party.

Phone/name/screenshot alone is not a stable payment key.

## Daily/event reconciliation

For every approved provider/account/period:

1. Ingest signed callbacks idempotently into immutable provider-event evidence.
2. Query/import official transaction/settlement report when available.
3. Compare opening + successful/failed/reversed/refunded amounts/counts and stable references.
4. Classify matches and exceptions.
5. Auto-resolve only deterministic, policy-approved matches.
6. Queue ambiguous/wrong-amount/wrong-payee/duplicate/missing items for human review.
7. Record corrections as new ledger/audit events; never overwrite history.
8. Obtain merchant/provider acknowledgment where responsibility requires.
9. Close period only when material exceptions have owners and due dates.

Because Ubi is non-custodial, provider settlement to merchant may not equal Ubi booking payment status timing. Both are visible and separately reconciled where provider data allows.

## Discrepancy playbooks

### Customer reports paid; provider not verified

- record report/evidence safely and say “reported, being checked”;
- search exact provider/payee/reference/amount/time through official path;
- check wrong merchant/account/environment/reference and delayed provider status;
- merchant decides service hold/continuation under policy without Ubi inventing payment truth;
- if no verification, retain reported/failed/unknown outcome and explain dispute/provider path.

### Provider verified; Ubi missing/stuck

- validate payee/amount/currency/reference and callback/query authenticity;
- idempotently apply missing provider evidence to correct payment attempt;
- update booking only through approved payment/booking policy command;
- notify merchant/customer of corrected status and duplicate-action prevention;
- fix ingestion/reconciliation gap and test replay.

### Ubi shows verified; provider does not

- immediately stop dependent automated confirmation/refund and classify severe;
- identify whether forged callback, wrong environment/account, bug or provider reversal;
- correct via append-only event, communicate affected merchant/customer and review all same-path records;
- invoke security/privacy/legal as needed.

### Duplicate concern

- determine duplicate UI tap/request versus distinct provider charges;
- idempotency prevents duplicate Ubi initiation, but provider evidence controls charge claim;
- never refund one record until both charges/payee/settlement status are verified;
- use provider/merchant refund process and preserve linked references.

### Amount/currency/payee mismatch

- quarantine; do not auto-apply partial/over/foreign/wrong-payee amount;
- merchant/provider determines top-up/refund/reallocation under contract/law;
- Ubi never moves value internally as if it were custodian.

### Callback timeout/unknown outcome

- do not resend non-idempotent initiation blindly;
- poll/query by stable provider/idempotency reference where supported;
- if provider cannot resolve, escalate and keep action-required with merchant/customer guidance;
- expire booking hold only under authoritative booking policy, not device/provider guess.

## Refund, reversal and dispute

Before action verify:

- authorized merchant/staff role and step-up/approval;
- original provider-verified transaction/payee/amount/currency/status;
- booking/service/cancellation/refund policy snapshot and customer request;
- provider capability/deadline/fees/partial support;
- no existing refund/reversal/dispute in flight;
- consumer/legal exception and merchant/Ubi responsibility.

Refund is sent to provider/merchant official flow; state becomes processing until provider verifies outcome. A submitted refund is not a completed refund. Reversal/chargeback/provider clawback is linked without erasing original success.

Ubi may facilitate evidence/communication but does not claim final adjudication authority beyond contract and approved process.

## Cash and pay-at-location

- Cash status is merchant-recorded, not provider-verified.
- Record collector/staff, amount/currency/time and receipt/reference as policy permits.
- Shared staff accounts are prohibited; corrections are append-only/audited.
- Cash disagreement follows merchant/customer dispute process; Ubi does not imply bank verification.
- Digital provider outage never disables cash/manual booking if merchant policy allows.

## Fraud and account compromise

Indicators include changed payee details, unusual refunds/exports, forged proof, provider token compromise and social-engineered ownership change.

- Freeze high-risk payment/provider/ownership action, not unrelated service by default.
- Verify legitimate owner/payee through independent channel and review recent changes.
- Revoke/rotate credentials and invoke runbook 67.
- Preserve evidence and avoid warning suspected insider/attacker in a destructive way.
- Risk signals alone do not permanently punish shared-phone/cash/offline users; review and appeal apply.

## Communication

Use precise language:

- “You reported payment” versus “Provider verified payment.”
- “Refund requested/processing” versus “Refund completed.”
- identify responsible merchant/provider/Ubi action and next update;
- warn against paying again while outcome is unknown;
- never request PIN/OTP/password/full credential;
- explain whether booking remains held, confirmed, action-required or unaffected.

## Closure evidence

- original and final provider/Ubi/booking states;
- matching fields and official source/checksum/report version;
- all corrections/refund/reversal linked events;
- merchant/customer/provider communications and obligations;
- financial amount/count exception resolution and approvers;
- root/control/test/runbook/product action;
- retained/deleted evidence under approved schedule.

## Readiness gates

- Provider sandbox/certification proves signatures, replay, idempotency, query and reconciliation exports.
- Reported/processing/verified/unknown/refund/reversal states appear correctly in UI/audit.
- Daily reconciliation and exception queue run with synthetic mismatch cases.
- Duplicate/timeout/wrong amount/currency/payee tests pass.
- Cash/manual flow remains available independently.
- Step-up/separation/JIT and merchant/customer scripts are exercised.
- Licensed Ethiopian provider/counsel confirms Ubi's non-custodial role and exact responsibilities before real money.


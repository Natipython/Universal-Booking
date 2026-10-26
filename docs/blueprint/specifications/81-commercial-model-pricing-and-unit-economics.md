# Commercial model, pricing experiments and unit economics

Status: **Draft — no approved prices, budget, tax treatment or revenue forecast**  
Updated: 2026-10-04 (Africa/Addis_Ababa)  
Owners: CEO/Finance; qualified Ethiopian accounting/legal review for applicable obligations  
Depends on: market evidence 80, costs 62, vendor facts 82, scorecard 83

## Candidate commercial role

Start as merchant booking software. A separate recurring merchant subscription is the candidate core revenue; optional provider-heavy capabilities may have transparent usage limits/add-ons. Avoid making an initial product dependent on consumer marketplace commissions or holding customer deposits.

The proposed non-custodial arrangement is not a legal exemption. Actual fund flow, contracts, payee, refunds, fees and settlement must be verified with the licensed provider and applicable Ethiopian reviewers before real-money activation. Do not claim Ubi has a payment license or is a payment-system operator.

## Revenue versus money observed

| Amount | Accounting/measurement boundary |
|---|---|
| Merchant SaaS subscription collected | Candidate Ubi software revenue; recognition/tax treatment requires accounting review |
| Customer pays merchant for photography | Merchant service/payment volume, not automatically Ubi revenue or cash |
| Provider fee charged to merchant/customer | Third-party charge with disclosed payer; not silently a Ubi charge |
| Ubi usage/add-on fee | Separate approved price/contract, collection and accounting records |
| Payment reported but not verified | Evidence state, not proven settled volume or revenue |

Subscription status and merchant customer payment status are distinct. A software subscription billing outage cannot quietly delete customer bookings or financial history; commercial restriction and export access require explicit policy.

## Pricing discovery before setting tiers

1. Measure booking volume, number of staff/devices, current administration time and affordability with actual merchants.
2. Obtain dated provider/cloud/communications quotes, including Ethiopia routing, taxes/fees, FX, minimum commitments and failure/retry billing.
3. Compare a simple core plan with a clearly bounded usage option. Test comprehension and willingness to pay using honest scope, not invented competitor prices.
4. Separate one-off assisted onboarding from recurring support cost. Disclose whether onboarding is free, subsidized or charged.
5. Test actual paid continuation after a clearly defined trial. “Would you pay?” is weaker evidence than retained paid use.
6. Approve collection method, billing/invoice/refund/cancellation terms and applicable tax/accounting processes before selling.

No ETB amount is approved here. Avoid a complicated five-tier grid before proving the first merchant task. A free trial needs an end date, capacity/usage boundaries and truthful feature limitations; it is not lifetime free unlimited messaging/calls.

## Cost worksheet and formulas

Record each input with source/date, currency, owner, confidence and low/base/high case.

```text
merchant_variable_cost = attributable compute/database/storage/egress
                       + notification/channel/payment integration units
                       + media relay/processing units, if enabled
                       + support_minutes × fully_loaded_cost_per_minute
                       + allocated fraud/refund/exception handling

merchant_contribution = collected_net_software_revenue - merchant_variable_cost
contribution_margin = merchant_contribution / collected_net_software_revenue

monthly_cash_burn = actual_cash_outflows - actual_cash_inflows
runway_months = unrestricted_operating_cash / positive_monthly_cash_burn

break_even_merchants = fixed_monthly_cost / positive_mean_merchant_contribution
```

Use zero/negative guards: undefined contribution margin with zero revenue is not 0%; no finite break-even exists when contribution is non-positive. Separate cash timing from accrual/profit and merchant/provider funds from Ubi operating cash.

Include founder time at a transparent replacement-cost assumption even when unpaid, to expose operational unsustainability. Do not invent salaries, market rates, budget or foreign-currency costs as facts.

## Cost categories often missed

- Physical low-cost Android devices, test SIM/data, field travel and participant research.
- Translation/content review, accessibility, security review and legal/accounting work.
- Signing/release infrastructure, backups, restore drills and monitored environments.
- SMS OTP fraud, retries, sender/routing limitations and failed delivery.
- Provider onboarding/support and per-message/media/relay costs.
- Customer/merchant support, disputes, privacy requests and account recovery.
- Taxes, FX movement, bank/payment fees and vendor minimum commitments.
- Incident response, emergency external expertise and exit/export/deletion costs.

## Commercial safeguards

No unlimited messaging/call promise without an evidenced affordable envelope. Quotas must be understandable and cannot discard already accepted work. Budget alerts can stop new optional workloads but must preserve safety, booking visibility and reconciliation.

Publish support coverage actually staffed. A founder-led pilot does not have a 24/7 support team. Track merchant acquisition cost and churn; avoid lifetime-value claims from a few trial weeks.

## Finance evidence gate

Before paid/live launch: actual operator/ownership facts, approved role/contracts, dated cost quotes, funded low/base/high cash plan, collection/invoice/tax process, user-understood price/limits, exception reserve and assigned finance ownership. These are EV-09/EV-13 inputs, not verified outcomes.

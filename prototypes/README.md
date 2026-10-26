# Ubi isolated feasibility prototypes

Status: **V0 synthetic validation only — not a production scaffold or launch approval**

This directory implements the bounded feasibility work in blueprint 72/85. It is separate from future production `apps/` and `services/` so demo shortcuts cannot quietly become live authorization, booking truth or data policy.

Current slices:

- [Android business V0](android_business/README.md): Flutter Android project with synthetic SQLite booking intents and outbox, working manual entry/review/details, draft persistence, simulated conflict/revision/cancellation and scope switching. Host tests pass; APK/device evidence is absent.
- [Booking integrity V0](booking_integrity/README.md): separate Go/PostgreSQL transaction and loopback HTTP experiment, with held/committed capacity, retries, atomic moves/quote conversion and contract tests. Only a newly created disposable database is allowed. A spec-linter disagreement remains recorded, not waived.
- [Customer link V0](customer_web/README.md): no-JavaScript prepared-link/edit/review/outcome experiment with scoped ephemeral grants and an adapter to the disposable integrity service. Synthetic browser/negative/retry checks pass; real identity, persistent continuation and native app links are absent.

Canonical dummy records live in `fixtures/booking-v0.json`; prices and future dates are invented. Android/server integration has not been implemented.

No real customers, contact numbers, credentials, identity verification, client/server synchronization, external messages, payments, calls or company/provider accounts are used. Storage is deliberately unencrypted synthetic-only storage. A fixture principal header is not authentication. Do not put real merchant/customer data here.

Protocol: `validation/WB-01-task-protocol.md`. Technical test results and screenshots are generated, not design mockups or market-validation evidence. WB-04 component checks now exist; physical Android, weak-network, localization/comprehension evidence and review remain pending before WB-05 closure.

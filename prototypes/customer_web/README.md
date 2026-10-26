# Ubi — WB-04 customer-link feasibility

Status: **Local synthetic component, not a production customer app.**

This slice continues a prepared conversation draft in a browser without installing Ubi. It uses the invented WB-01 businesses and the disposable WB-03 PostgreSQL booking service. No existing database, real contact, live inbox, payment or public domain is accepted.

## Try the local demo

From the repository root:

```sh
cd prototypes/customer_web
UBI_RUN_SYNTHETIC_WEB=1 ../../.tooling/go/bin/go run ./cmd/demo
```

Open the printed `http://127.0.0.1:<port>` address. Stop with Ctrl+C. The process also stops after 30 minutes. Its newly created private PostgreSQL cluster has no TCP listener and is stopped/removed on normal shutdown; cleanup failure returns a nonzero exit and retains files for diagnosis. Abrupt process/host termination is not a tested recovery path.

Requirements on this development host: repository Go runtime, PostgreSQL 18 binaries and `btree_gist`. This is not an Android build or production installation procedure.

## What to check

1. Choose a prepared instant booking or a custom enquiry.
2. Open the generic preview; no private draft is disclosed and no time is booked.
3. Explicitly continue, correct the invented date/time/customer choices, and review them.
4. Submit once. Instant booking commits a `confirmed_demo` reservation; custom enquiry saves `request_saved_demo` in memory and does **not** hold capacity.
5. Repeat the same submission safely. Unknown results remain checking; after review expiry, recovery is read-only and cannot create a booking.

The review shows business, service, date/time, practice customer, duration/price or event size/enquiry expectations. White is the surface, purple the primary action, sky blue the informational accent, green verified demo confirmation and amber an enquiry/warning. Status is also written in words, never colour alone.

## Reproduce checks

```sh
node .tooling/validate-customer-v0.mjs
```

Run that command from the repository root. It performs format/vet/module checks, independent OpenAPI validation, local Postman specification/collection checks, and three race-detector test runs including Chromium with JavaScript disabled. It writes a source-bound [machine receipt](artifacts/validation-receipt.json), private-cluster cleanup receipt, redacted local Postman output, document byte samples, browser receipt and rendered screenshots. The browser harness uses the configured Codex dependency runtime's Playwright package and `/usr/bin/google-chrome`; it never uses a personal browser profile.

See [CONTRACT.md](CONTRACT.md) for authority, privacy, expiry and retry rules. The saved [OpenAPI](openapi.json) and local Postman collection describe this experiment only; they are not the canonical production API.

## Deliberate limits

- Grants, sessions, reviews and custom enquiries are memory-only. Restart loses continuation; persisted booking receipts are disposable test data, not a complete production recovery system.
- An unused forwarded bearer link may be claimed by another browser. Browser binding is not verified customer identity; there is no contact/history/manage permission.
- English, Gregorian dates and the international clock only. Ethiopian calendar/clock, Amharic and real-user comprehension must still be designed and validated.
- Browsers sending `Origin: null` must supply same-origin navigation metadata as well as the issued cookie and CSRF token. Older browsers without that metadata fail closed; their compatibility fallback is not implemented.
- No native app installation/deep-link association, production auth, merchant Android integration, verified public business page, real service catalogue, hours/availability grid, full quotation workflow, messages, calls, payment or notification.
- Narrow viewport/root-font simulation is not physical-device, TalkBack, full accessibility or weak-network evidence. Small payloads do not prove fast user experience.
- All full acceptance scenarios and G7/G8/G9 remain open. Do not invite real customers on the basis of this receipt.

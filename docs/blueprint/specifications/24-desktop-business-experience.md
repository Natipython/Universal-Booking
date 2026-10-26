# Desktop business workspace

Status: **Draft — responsive web candidate; packaged desktop not approved**  
Updated: 2026-10-04 (Africa/Addis_Ababa)  
Depends on: specifications 12, 15–22, 23, 33–35, 44, 60–65

## Role and delivery order

Desktop is the next business workspace after Android. A lightweight customer completion page is supporting infrastructure for booking links, not a second business application that replaces this order. The desktop form remains DEC-018 open; responsive browser delivery is the candidate to test on actual Windows/Mac machines before separate packaging.

It shares domain/API/authorization contracts with Android, not Android's small-screen layout. No independent desktop definition of confirmation, price or customer identity is allowed.

## Workspace map

| Area | Layout | Common task |
|---|---|---|
| Today/calendar | Navigation, scoped schedule, details pane | Find/add/move bookings with exact dates |
| Requests/quotes | Action queue, request, versioned quotation | Assign, clarify, quote and review acceptance |
| Messages | Queue, conversation, linked booking context | Claim, reply, pause automation, hand back |
| Services/hours | Configuration list and editor | Preview effect of changes on existing commitments |
| Customers | Authorized search and relationship summary | Find contact/history without cross-business leakage |
| Staff/settings | Role/scope preview and configuration | Invite/offboard, publish links, manage connections |

Use progressive disclosure for one-person businesses. Persistent business/branch/operation labels prevent editing the wrong scope. A hotel example may demonstrate multiple operations in synthetic testing, but hotel inventory is not pilot functionality.

## Interaction rules

1. Every task works with keyboard and visible focus. Shortcuts are optional, discoverable and cannot fire destructive actions in text fields.
2. A schedule list/table alternative complements the visual calendar. Dragging is never the only reschedule route.
3. Drag/drop creates a proposal with old/new times and consequence review; confirmation still needs a server transaction.
4. Filters apply explicit scope, persist predictably and show when items are hidden. Search cannot cross authorization boundaries.
5. Details open without losing queue position. Browser back/refresh preserves eligible draft state safely.
6. Bulk actions, when separately approved, show record count, individual eligibility, skipped failures and per-record outcomes. “All successful” cannot mask partial execution.
7. No silent customer-data download into local browser storage. Browser caching/offline scope requires separate security analysis.

## Online and failure semantics

The first desktop candidate is online for writes. Lost connectivity keeps entered information in a purpose-limited draft where safe; it does not imply a committed booking. Reconnection reconciles a timed-out command with the same idempotency key. Disable repeat submit while checking, but provide a recoverable failure state instead of an endless spinner.

Browser tab competition and two staff users must pass the same version/conflict rules as Android. Real-time updates improve freshness but never prove that a selected slot remains available at commit.

## Permissions and shared computers

- Staff see only authorized branch/operation/assigned records; hiding a button is insufficient enforcement.
- Owner actions require the appropriate capability and step-up controls. The workspace must not contain an unrestricted “approve everything” action.
- Shared-PC policy defines session locking, expiry and sign-out. Avoid persistent customer content on a public/shared computer.
- Exports require their own capability, purpose, audit and secure delivery. Screen printouts must not expose private notes by default.
- Tokens and provider credentials never enter URLs, client-readable configuration or support screenshots.

## Messages/calls and accessibility

Native booking-linked threads share history and safety rules with Android. External threads show source and reply constraints. A future in-app voice call needs microphone permission on explicit action, visible active-call controls, reconnect/failure state and block/report; desktop browser calling is not a promised current feature.

Apply specification 22 to keyboard, screen-reader tables, zoom/reflow, contrast, reduced motion and localized dates. Low-bandwidth operation avoids giant client bundles and auto-loaded portfolio media.

## Acceptance and open items

Test receptionist, owner and assigned photographer tasks on the proposed minimum hardware/browser matrix, including keyboard-only, 200% zoom, Amharic layout, expired sessions, wrong tenant, double-submit, lost network and concurrent reschedule.

Open: supported browsers/hardware, offline scope, packaging need, print workflow, multi-monitor demand, bulk-action necessity and call media provider. No usability/performance results exist yet.

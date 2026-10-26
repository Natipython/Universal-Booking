# Accessibility, localization, calendar and regional-format specification

Status: **Draft**  
Depends on: actors, product principles, functional requirements and quality attributes  
Blocks: content design, design system, Android/desktop/web experience, domain model, APIs and acceptance testing

## Purpose

Ubi must be understandable and operable by people with limited booking-app experience, people using assistive technology, people under time pressure and people using languages, calendars, names, phones and addresses that generic global software often handles badly.

This specification treats comprehension and accessibility as correctness. It does not describe target users as unintelligent. Difficulty is a product/system failure when the interface assumes technical knowledge, literacy, memory, vision, dexterity, connectivity or conventions that the task does not require.

## Standards and target

- Customer web, desktop web and public content target **WCAG 2.2 Level AA**.
- Native Android targets equivalent perceivable, operable, understandable and robust behaviour using platform accessibility semantics and testing.
- Ubi SHOULD exceed minimum standards where a larger target, clearer wording or simpler choice materially reduces error.
- Conformance cannot be claimed from an automated scanner alone. Manual keyboard, screen-reader, magnification/font scaling, contrast and task tests are mandatory.
- Country-specific disability/accessibility duties remain part of the later jurisdiction review.

## Inclusive product-language policy

### Required voice

Content MUST be:

- direct, respectful and calm;
- written in familiar business/customer words;
- explicit about who will do what and what happens next;
- short without removing price, policy, timing or risk;
- translated for meaning and local use, not word-for-word from technical English.

Ubi MUST NOT expose terms such as *tenant*, *entity*, *mutation*, *sync conflict*, *state machine*, *OAuth scope* or *webhook* in ordinary workflows.

Preferred transformations include:

| Internal concept | User-facing direction |
|---|---|
| organization/operation context | business, branch and service area/name |
| command queued in outbox | “Saved on this phone. We’ll send it when you’re connected.” |
| synchronization conflict | “This time was taken before your booking reached Ubi. Choose another time or ask a manager.” |
| payment intent pending provider evidence | “Payment is being checked. Your booking is not confirmed yet.” |
| authorization failure | “You don’t have permission to do this. Ask the owner or a manager.” |
| destructive archive | “Hide this service from new bookings” plus effect on existing bookings |

Final wording requires comprehension tests and professional translation.

## Cognitive accessibility and error prevention

| ID | Requirement |
|---|---|
| A11Y-COG-001 | Each step MUST present one clear primary decision/action; secondary and destructive actions must not compete visually. |
| A11Y-COG-002 | The interface MUST show a useful example, template, recent value or safe default instead of an unexplained blank setup whenever possible. |
| A11Y-COG-003 | Labels MUST remain visible after input; placeholder-only labels are prohibited. |
| A11Y-COG-004 | Required fields MUST be explained before submission and marked with text/semantics, not color alone. Optional questions SHOULD be identified as optional. |
| A11Y-COG-005 | Related choices MUST be grouped in the order people naturally decide them. Branching questions must hide irrelevant work without changing earlier answers silently. |
| A11Y-COG-006 | Long or consequential flows MUST provide a plain-language review showing business/service, customer, date/time, price/deposit, confirmation condition and policy consequence before commitment. |
| A11Y-COG-007 | Back navigation, interruption, screen rotation and safe app restart MUST preserve entered work unless security requires otherwise; any loss must be disclosed before it occurs. |
| A11Y-COG-008 | Destructive and high-impact actions MUST name the object and consequence. Confirmation wording such as “Are you sure?” by itself is insufficient. |
| A11Y-COG-009 | Success MUST state what was actually completed and the next expected action. Submission, payment reported, payment verified and booking confirmed must never share the same success copy. |
| A11Y-COG-010 | Error recovery MUST put focus near the problem, preserve valid input, explain the correction and never erase the whole form for one invalid field. |
| A11Y-COG-011 | Technical identifiers may be offered for support/copying but MUST NOT be the only explanation a person receives. |
| A11Y-COG-012 | Time-limited actions MUST show the time remaining and provide extension/recovery where safety permits. Important completion must not depend on reading speed alone. |
| A11Y-COG-013 | Ubi MUST avoid memory tests, puzzle-style CAPTCHA and blocking copy/paste into OTP/password fields. Authentication must support accessible alternatives. |

## Visual and interaction accessibility

### Touch and pointer targets

- Android interactive targets MUST provide at least **48 × 48 dp** of effective touch area, including icon buttons.
- Web interactive targets SHOULD provide at least **44 × 44 CSS px** where practical and MUST meet WCAG 2.2 minimum target-size rules and spacing.
- Closely spaced destructive/confirm actions are prohibited.
- Gesture-only actions MUST have a labelled button/menu alternative. Swipe may accelerate a task but cannot be its only path.

### Typography, scaling and layout

| ID | Requirement |
|---|---|
| A11Y-VIS-001 | Body text MUST use a legible size and line height; the approved design-token values require testing on D1 screens and target scripts. |
| A11Y-VIS-002 | Android MUST remain usable with platform font/display scaling through the supported accessibility range; web MUST support text zoom to 200% without loss of content or operation. |
| A11Y-VIS-003 | Customer web MUST reflow at a 320 CSS-pixel viewport without two-dimensional scrolling except content that intrinsically requires it. |
| A11Y-VIS-004 | Text MUST achieve at least 4.5:1 contrast for normal text and 3:1 for large text; meaningful non-text UI boundaries/states MUST meet 3:1 where WCAG applies. |
| A11Y-VIS-005 | Content MUST not be obscured by keyboards, sticky actions, system bars, zoom or orientation changes. |
| A11Y-VIS-006 | Truncation MUST NOT hide price, status, date/time, customer identity disambiguation or destructive consequences. A reliable expansion/access path is required. |
| A11Y-VIS-007 | Horizontal screen orientation and larger layouts MUST adapt without creating a second, inconsistent task model. |

### Color semantics

Ubi’s provisional palette is governed as follows:

- white/near-white: calm surfaces and reading areas;
- purple: brand, selected state and principal action;
- sky blue: neutral information, connection and secondary assistance;
- green: verified completion/approval only;
- amber: pending attention or reversible risk;
- red: error, danger or destructive consequence.

Every use must pass contrast testing in its foreground/background/state combination. Color MUST be accompanied by text, shape, icon and/or accessibility semantics. “Requested” or “saved only on this phone” cannot be green because neither is confirmed success.

## Assistive technology semantics

| ID | Requirement |
|---|---|
| A11Y-AT-001 | Every control MUST expose an accurate accessible name, role, value/state and action. Decorative images/icons MUST be hidden from assistive technology. |
| A11Y-AT-002 | Focus/read order MUST follow the visual and task order. Opening a dialog, error or new step MUST place focus predictably; closing it MUST return focus safely. |
| A11Y-AT-003 | Screen changes, pending work, validation errors and confirmed success MUST be announced with appropriate live-region/platform semantics without repeatedly interrupting the user. |
| A11Y-AT-004 | Custom controls MUST reproduce native keyboard, focus, state and screen-reader behaviour or be replaced with native controls. |
| A11Y-AT-005 | Images that communicate a service, venue, document or status MUST have concise meaningful alternatives; complex information needs a text equivalent. |
| A11Y-AT-006 | Android releases MUST be tested with TalkBack and switch/keyboard navigation as applicable. Web/desktop releases MUST be tested with keyboard and at least the approved screen-reader/browser combinations. |
| A11Y-AT-007 | Tables, calendars and schedules MUST expose headings/relationships and offer a linear/list alternative on small screens or assistive technology. |

## Keyboard, voice, motion and media

- Every desktop/web function MUST be operable with a keyboard, with visible focus and no keyboard trap.
- Keyboard order MUST remain logical across responsive layout changes.
- Information cannot depend on hover. Hover help needs focus/touch access.
- Motion MUST respect the platform/browser reduced-motion preference. Non-essential parallax, looping and celebratory animation are excluded from critical workflows.
- Animation must not delay completion or communicate state alone. Flashing content is prohibited beyond safe accessibility thresholds.
- Voice input may be supported, but users MUST review interpreted dates, names, prices and services before commitment.
- Prerecorded instructional video MUST have captions and a text alternative. Important audio information needs a transcript or equivalent visual message.
- Audio/call releases require separate volume, device, caption/transcript feasibility and consent work; calling remains future-gated.

## Form, date/time and status presentation

| ID | Requirement |
|---|---|
| A11Y-FRM-001 | Text fields MUST identify the expected value, local example and permitted format before error occurs. |
| A11Y-FRM-002 | Appropriate input modes/autocomplete MAY reduce typing, but auto-filled values must remain reviewable and editable. |
| A11Y-FRM-003 | Date selection MUST offer a calendar/list appropriate to the task and a typed/accessibility path. The selected result MUST be repeated as weekday, day, month and year. |
| A11Y-FRM-004 | Relative words such as “today” and “tomorrow” MUST be followed by the explicit local date when used for a commitment. |
| A11Y-FRM-005 | Time MUST include an unambiguous local format and timezone when parties/branches may differ. 12-hour time includes AM/PM translated appropriately; 24-hour time has no AM/PM. |
| A11Y-FRM-006 | Status labels MUST use stable verbs/phrases and explain the responsible next actor: “Waiting for the business to approve,” not merely “Pending.” |
| A11Y-FRM-007 | Validation MUST happen at a useful time; it must not announce errors on every keystroke or reject reasonable localized punctuation/spacing. |
| A11Y-FRM-008 | OTP and recovery flows MUST support paste, password-manager/OS code assistance and sufficient time. Users cannot be forced to transcribe inaccessible visual content. |

## Language architecture

### Locale model

Ubi MUST separate:

- interface language;
- content language(s) supplied by the business;
- number/currency format preference;
- calendar display preference;
- time zone;
- clock presentation convention, including international versus Ethiopian-clock wording when supported;
- country/legal/payment configuration.

Changing language must not silently change currency, calendar, branch, country, price or legal agreement.

### Engineering requirements

| ID | Requirement |
|---|---|
| L10N-001 | Language/locales MUST use standard BCP 47 identifiers; time zones MUST use IANA identifiers and currencies ISO 4217 codes. |
| L10N-002 | User-visible strings MUST live in localization resources, not be embedded in business/domain code. Server-generated customer text must use versioned localized templates. |
| L10N-003 | Sentence fragments MUST NOT be concatenated. Plural, gender/case and select rules must use ICU/CLDR-style message formatting or an approved equivalent. |
| L10N-004 | Layouts MUST tolerate at least 50% text expansion in testing and support right-to-left mirroring before an RTL language is released. |
| L10N-005 | Fonts MUST render the complete approved script set, including Ethiopic, digits and punctuation, at all required weights without fallback instability. |
| L10N-006 | Search, sorting, comparison, casing and name validation MUST be locale-aware without assuming Latin characters. User names and business names are stored in entered Unicode form with safe normalization rules. |
| L10N-007 | Translation keys MUST carry context, screenshots/character constraints where needed, source version, reviewer and approval state. |
| L10N-008 | Machine translation MAY create an internal draft but MUST NOT publish legal, payment, policy, safety or commitment text without qualified human review. |
| L10N-009 | Pseudo-localization, missing-string detection, text-expansion, script-rendering and mixed-language tests MUST run before each localized release. |
| L10N-010 | When a business has not provided translated catalog content, Ubi MUST label the available content language rather than fabricate a translation or show an empty critical field. |

### Proposed language sequence

The launch language set cannot be approved before the first country/city and user research are chosen. If Ethiopia is the first market, the practical starting proposal is:

1. English as the engineering/source language.
2. Amharic as a professionally translated and locally tested first customer/business language.
3. Afaan Oromo and/or other languages according to the pilot city and actual users.
4. Swahili when entering markets where it is operationally dominant, with country-specific terminology rather than one generic East African translation.

This is a planning proposal, not an assumption that every East African user speaks or prefers the same language.

## Calendar, date and time architecture

Calendar support is a data-integrity problem, not only a visual preference.

### Canonical storage rules

| ID | Requirement |
|---|---|
| CAL-001 | An actual occurrence/instant MUST be stored as an unambiguous UTC instant plus its relevant IANA timezone and original local scheduling context when required. |
| CAL-002 | A date-only commitment (for example a local service day) MUST be stored as a date with its business calendar/timezone context, not converted blindly to midnight UTC. |
| CAL-003 | Recurring availability MUST store local rules, timezone and version/effective range; precomputing UTC forever is prohibited because timezone rules can change. |
| CAL-004 | Gregorian/ISO dates form the canonical calculation/interchange basis unless an approved domain requires otherwise. Ethiopian and other calendar dates are supported as verified input/display conversions. |
| CAL-005 | The originally presented/accepted local date, time, calendar and timezone SHOULD be captured in commitment evidence so a future conversion/library change cannot rewrite what the parties saw. |
| CAL-006 | Device time is not authoritative for expiry, ordering, payment, audit or permission. Clock offset may assist presentation only. |

### Ethiopian and Gregorian calendar experience

- A user or organization MAY select Ethiopian or Gregorian display after implementation is approved.
- Critical commitments SHOULD show the selected calendar prominently and the equivalent other-calendar date where this materially prevents mistakes.
- The calendar name must be visible; a bare numeric date such as `02/03/26` is prohibited.
- The final review/receipt MUST spell out weekday, month, day and full year in the selected language/calendar, plus time and timezone when relevant.
- Conversion MUST use an approved, thoroughly tested library/algorithm. Boundary tests must cover Ethiopian new year, Pagumen length/leap years, Gregorian leap years, year transitions and timezone day boundaries.
- Calendar conversion happens at presentation/input boundaries; booking and availability invariants use the canonical representation.
- Import/export/API contracts MUST declare calendar/format explicitly and default to ISO 8601 for machine exchange.

The initial engineering design must not promise Ethiopian-calendar editing until library, accessibility and field usability tests pass. Accurate dual-date display is preferable to a rushed calendar picker.

## Timezone rules

### Clock-convention clarity

Calendar preference, timezone and clock convention are separate. For the Ethiopia launch, ambiguous conversation/voice input such as “2 o'clock” must prompt clarification of the intended convention and a readable equivalent before commitment. Do not infer clock convention from Amharic/English, phone settings or calendar choice. A display conversion must not change the actual stored instant.

Ethiopian-clock input/display remains proposed until its conversion rules, translated AM/PM/day-period wording, boundary cases and locally observed comprehension are verified. Use explicitly labelled international times plus human clarification where a tested alternative is not yet available. Include these cases in SCN-018 and release evidence; do not silently apply an unreviewed fixed offset to user input.

- Every branch/fulfilment location has an explicit IANA timezone; it does not inherit the phone timezone silently.
- A new booking defaults to the service/branch timezone and labels it when the actor may be elsewhere.
- Changing a branch timezone is a sensitive migration requiring preview of future schedules; it cannot reinterpret existing instants silently.
- Notifications use the booking/recipient policy and show the appointment timezone where ambiguity is possible.
- Daylight-saving and government timezone-rule changes must be covered by library updates and regression tests even if the first launch location currently uses a stable offset.
- `Africa/Addis_Ababa` is a likely Ethiopia launch value, not a global or hard-coded default.

## Numbers, money and pricing

| ID | Requirement |
|---|---|
| REG-MNY-001 | Monetary values MUST store an exact amount and ISO 4217 currency; display uses the approved locale/currency rules. |
| REG-MNY-002 | Currency and language are independent. A language change cannot convert or relabel the amount. |
| REG-MNY-003 | Ubi MUST display currency code/symbol clearly where a symbol may be ambiguous and show the charge/deposit/balance meaning. |
| REG-MNY-004 | Decimal/grouping separators may be localized for presentation but parsing must reject ambiguous amounts before commitment and repeat the normalized amount for review. |
| REG-MNY-005 | Percentage, per-person, per-night, per-hour, estimated and starting prices require explicit labels and calculation explanations. |
| REG-MNY-006 | Rounding, minor units, tax/fee inclusion and exchange rules require country/provider approval; clients cannot invent rounding. |

## Phone, identity and names

- Phone contacts SHOULD be normalized to E.164 after country-aware parsing while preserving the user-entered/display form when useful.
- Country calling code selection must use country names/codes and searchable text; flags cannot be the sole representation.
- Shared, recycled or changed phone numbers are expected. A phone number is a contact method, not a permanent person identifier.
- Verification proves control at a time; it does not prove legal identity, ownership of a business or continued ownership forever.
- Names MUST support Unicode, multiple words, single names, local scripts, transliteration and ordering differences. “First name/last name” is not a universal legal model.
- A preferred/display name and a verified legal name, if ever required, MUST be separate fields with separate purposes and access.
- Customer matching must not use name spelling as a unique key.

## Address, landmark, map and service-area model

Ubi MUST NOT require a Western street-address structure where it does not fit local reality.

The regional address model SHOULD support:

- country and administrative area/city/subcity/district where useful;
- free-form local address/directions;
- familiar landmark and optional nearby-place description;
- map pin/coordinates with stated accuracy/source;
- business branch or saved service location;
- service area/radius/polygon or named areas for field services;
- contact instructions that are disclosed only when operationally necessary.

A map pin is not automatically a verified address. Exact home/customer coordinates are sensitive and must be purpose-limited, access-controlled and retained only as necessary. Users need a text/manual alternative when maps, GPS or data are unavailable.

## Content and catalog localization

| ID | Requirement |
|---|---|
| CONTENT-001 | A service/option/venue may have multiple localized names and descriptions, each with locale and publication state. |
| CONTENT-002 | The business remains responsible for merchant-authored claims and translations, while Ubi supplies safe templates, previews and reporting/correction paths. |
| CONTENT-003 | A booking snapshot MUST preserve the exact localized customer-visible service/price/policy version used at commitment. |
| CONTENT-004 | Search SHOULD match approved transliterations/aliases without replacing the canonical business-entered name. |
| CONTENT-005 | Generated notification and policy text MUST fall back through an explicit locale chain and MUST NOT mix languages unpredictably within one sentence. |

## Low-literacy assistance without loss of control

Ubi MAY provide:

- illustrated service templates using real locally meaningful examples;
- optional read-aloud instructions where language quality is verified;
- guided setup with progress based on business outcomes;
- safe suggested answers from an earlier conversation, visibly marked and editable;
- phone-number/customer lookup within authorized tenant scope;
- staff-assisted onboarding with recorded authorization and review.

Ubi MUST NOT:

- hide contractual price or policy behind pictures/voice only;
- auto-submit an inferred answer;
- describe automation as a human;
- require literacy-heavy training to recover from a normal error;
- use shame, childish language or “idiot-proof” wording in product/research materials.

## Research and acceptance testing

### Participant coverage

Testing MUST include people from the intended city/vertical across:

- first-time or infrequent booking-app users;
- business owners and staff with different authority levels;
- Android devices at the D1/D2 boundary;
- preferred launch languages and calendar conventions;
- people who use screen readers, larger text, keyboard/switch or other relevant access methods;
- realistic N0/N1 interruptions and shared/noisy business environments.

Compensation, consent, recording, data handling and recruitment must follow the approved research/privacy process.

### Critical task protocol

At minimum, the pilot protocol measures whether a person can:

1. identify the current business/branch/operation;
2. add a booking received by phone or walk-in;
3. distinguish request from confirmed booking;
4. find today’s schedule while offline;
5. understand and recover a synchronization/capacity conflict;
6. reschedule or cancel with consequences understood;
7. identify reported versus verified payment;
8. invite/restrict staff without granting ownership accidentally;
9. complete a customer booking from a shared link;
10. find help and communicate the booking/reference safely.

Record completion, time, hesitation, wrong turns, assistance, comprehension, recovery and confidence—not only taps. Results must be segmented by language, device, prior experience and access need where sample size permits.

### Accessibility release gate

A critical flow is not releasable when it has:

- a keyboard or screen-reader dead end;
- inaccessible verification/authentication;
- clipped/hidden commitment information at supported scaling;
- color-only status or unlabeled controls;
- an ambiguous calendar/currency/price result;
- an error that erases valid work;
- an untranslated or unreviewed legal/payment/safety message in a released locale.

## Authoritative external baselines

- [Web Content Accessibility Guidelines (WCAG) 2.2](https://www.w3.org/TR/WCAG22/) is the web conformance baseline.
- [Android Compose accessibility defaults](https://developer.android.com/develop/ui/compose/accessibility/api-defaults) documents platform semantics and the 48 dp minimum touch-target convention.

Platform guidance is a baseline; locally observed comprehension and access barriers still require product changes.

## Open decisions

- Launch country, city, vertical and participant profile.
- First released interface/content languages and professional translation reviewers.
- Ethiopian-calendar input/editing scope versus accurate dual-date display only.
- Approved fonts, typography scale, color tokens and dark-mode scope.
- Supported screen-reader/browser/device combinations.
- Whether read-aloud/voice input adds enough benefit to justify language, privacy, download and support cost.
- Map/geocoding provider, offline fallback and location-data retention.
- Jurisdiction-specific accessibility, consumer-language and records obligations.

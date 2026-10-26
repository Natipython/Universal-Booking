# WB-01 synthetic fixtures and task protocol

Status: **Prepared — not recruited, field-tested or commercially validated**  
Date: 2026-10-04 (Africa/Addis_Ababa)

## Safe test boundary

Use only the two invented businesses, three operations, four service templates and three made-up customer choices. No real phone number/name, chat history, identity document, credentials, money or publicly bookable capacity. Fixture date is 15 January 2030, Gregorian, Addis Ababa international 24-hour time. Do not describe it as today's real schedule. Fixed timezone conversion is a V0 fixture constraint, not a regional calendar implementation.

The scenario script describes merchant owner/receptionist work but V0 has no authentication or role enforcement. Scope switching checks local repository isolation; it does not prove real identity/server authorization. Founder city/vertical/team/payment choices remain proposed/unknown where recorded in the decision register.

## Tasks, in order

| Task | Prompt without coaching | Correct observable outcome | Blueprint linkage |
|---|---|---|---|
| T01 | Explain whether this app can make a real reservation today | Participant recognizes synthetic-only, no server, no money; do not proceed if confused | SCN-001/028 |
| T02 | A caller requests the dummy portrait session at 09:00 on the fixture day. Add it and review the choices | Correct service, date, explicitly labelled clock, dummy customer/source/price; review before local save | SCN-003/007/018 subset |
| T03 | Find the saved item and explain whether the business/customer has confirmed it | Understands local tentative, not confirmed; can find details and next-state explanation | SCN-003/018 subset |
| T04 | Close the demo and reopen it. Find the saved work | Local item/command/history persist; manual device test still required | SCN-003/013 subset |
| T05 | Add another dummy item overlapping the first | Clear local-conflict message, input retained; choose a valid alternative | SCN-005/006 local-only subset |
| T06 | Simulate a conflict on an existing item and change its time | Recognizes simulation; original evidence retained; new local intent remains not confirmed | SCN-006/013 subset |
| T07 | Cancel the local item and explain the consequence | Preview/confirmation, no real customer/refund action; retained local history | SCN-008 local-only subset |
| T08 | Switch between two businesses and studio/event operations | No unrelated items shown; full active business/operation understandable | SCN-002 local-only subset |
| T09 | Leave an edited form, return and continue | Restored authorized synthetic draft; distinguish draft from submitted local item | SCN-013 subset |
| T10 | Repeat key tasks with large fonts and accessibility support | Operable controls/readable date/status and correct comprehension; actual TalkBack test pending | SCN-018/021 subset |

## Moderation/recording

Before involving people, explain research purpose, voluntary participation, whether recording is proposed, retention/access and right to stop; use an appropriate consent/research notice. Do not recruit or contact participants from this script alone. No recorded voice/video or personal participant data by default.

Score first attempt without coaching: correct outcome, correct understanding, elapsed time, intervention, recoverable mistakes, accessibility barrier and participant's own explanation. Then debrief and help. Avoid shame or treating unfamiliarity as a person’s defect. A technically successful tap with wrong confirmation/date understanding is a failed task.

This protocol can prepare V1 but cannot substitute for field evidence or legal/privacy review. Actual groups/languages/devices, approval/permission, recruitment and sample size remain to be established. The founder set the business-app support floor to Android 11+ and 4 GB RAM on 2026-10-04; this is not an actual phone model or a performance test. Guest web compatibility is separate.

## Technical versus field receipts

Automated tests check local persistence/atomicity/deduplication, fixture isolation, limited UI paths, 200% font scaling and selected color pairs. A subprocess SIGKILL verifies recovery of a committed SQLite write on this Linux host. It is not Android OS process-death/background/low-storage testing.

Physical D1 release/profile build testing still needs Android SDK/JDK, a real low-cost phone, test data/SIM/network setup and measured startup/frame/data/battery traces. Field comprehension, Amharic review, Ethiopian-calendar/clock entry, real authorization and server sync are not claimed.

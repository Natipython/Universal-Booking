# Ubi operations and support console

Status: **Draft — least-privilege workforce workflow, not production tooling**  
Updated: 2026-10-04 (Africa/Addis_Ababa)  
Depends on: specifications 10–12, 33–49, 61–65, policies 50–57, runbooks 66–70

## Boundary

This is the internal Ubi workforce console, not an owner account with access to every business. Separate workforce identities, strong authentication and scoped capabilities govern it. A sole founder may perform several accountable roles; that does not justify permanent unrestricted production access.

## Work areas and permitted starting view

| Queue | Default evidence | Controlled action |
|---|---|---|
| Onboarding | Business setup checklist, consent/authorization, test status | Assist within approved merchant scope; merchant reviews publication |
| Support | Case ID, redacted command/status timeline, device/release, sync error | Diagnose/reconcile without customer-content access first |
| Identity/recovery | Verification history and account-risk metadata | Approved recovery process, not manual owner substitution |
| Integration health | Connection/capability, token expiry metadata, retry/queue status | Reconnect/disconnect via official merchant authority; never show token |
| Payment exception | References, amount/currency, provider evidence and booking state | Reconcile/route merchant-provider case; never invent settlement |
| Privacy | Verified request, activity roles, affected stores/holds | Tracked export/correction/deletion/restriction workflow |
| Safety/appeals | Report category, necessary evidence, restriction history | Proportionate restriction and reasoned appeal path |
| Reliability | SLOs, failures, queue age, backup/exercise state | Follow runbooks, feature gates and controlled incident actions |

Only implemented and authorized queues are enabled. A deployment dashboard or direct database console is not a substitute support workflow.

## Case lifecycle

1. Record purpose, requested action, affected actor/organization and severity.
2. Authenticate/verify the requester proportionately. A screenshot or caller's knowledge of a business name is insufficient.
3. Diagnose from redacted metadata and scoped event/command history.
4. Request time-bounded content/action access only when necessary and authorized by the correct policy/role.
5. Perform a controlled domain command or escalation; record actor, before/after state and evidence.
6. Give a truthful explanation: known, unknown, retained, next action and responsible party.
7. Close with user-visible result where appropriate, expiry of access and review of follow-up work.

Unknown outcomes remain open for reconciliation; neither “case closed” nor a screenshot proves booking/payment success.

## High-risk controls

- No generic invisible “log in as merchant.” Any delegated session must be scoped, purpose-bound, visible/audited and unable to impersonate unrestricted consent.
- No direct rewrite of booking, ledger, quote or audit history. Recovery commands preserve attribution and original evidence.
- Exports and content views have independent authority and retention controls.
- Break-glass access has reason, recent strong authentication, short expiry, alert and retrospective review. It cannot remove financial, legal or tenant-isolation constraints.
- Until independent staffing exists, high-risk peer review uses a qualified external reviewer or action remains unavailable; do not invent an imaginary second approver.
- Production/test environments are unmistakable. Synthetic support practice never contacts real customers or charges money.
- Support credentials, provider tokens, OTPs and raw private links are never requested from users or pasted into cases.

## Merchant point of view

The merchant sees why help is needed, what data/action access is requested, how long it lasts, the case reference and what was changed. Access expires on completion/revocation. Whether merchant authorization is sufficient for a particular data subject or legal activity is reviewed under policy 51/52; a checkbox is not universal legal authority.

## Console acceptance

Negative authorization tests cover wrong tenant, stale/revoked grant, unrelated case, expired delegation, export without permission, direct-history-edit attempt and misuse of break-glass. Usability tests cover assisting a low-confidence owner without taking their password, explaining an offline conflict and preserving an unknown payment result.

Operators need keyboard navigation, accessible tables, filtered queues, exact dates and redacted copyable references. Every live queue has an actual owner and escalation path; gaps block the relevant production capability rather than being hidden in a staffing plan.

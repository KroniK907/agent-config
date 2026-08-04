---
name: research
description: Investigate a wayfinder:research ticket, post structured findings as an issue comment, and propose non-binding tracker updates. Use when a wayfinder:research ticket is ready, wayfinder Route suggests research, or the user wants fact-gathering on a map research ticket.
---

# Research

Load a **`wayfinder:research`** ticket, investigate per [behavior rules](REFERENCE.md#behavior-rules), post a **structured findings comment** on the ticket, and hand off **non-binding Proposed tracker updates**. Does **not** append decision-log rows, post Reconcile approval phrases, or edit the map without human review.

**v1 is human-initiated HITL only** — cloud AFK pickup deferred to [#10](https://github.com/KroniK907/skills/issues/10).

## Not this skill

| Skill | When instead |
|-------|----------------|
| [wayfinder](../wayfinder/SKILL.md) | Chart, Materialize, Reconcile, Route only |
| [grill-me](../grill-me/SKILL.md) | Depth-first Q&A → binding `{MAP-SLUG}-GM-*` rows |
| [strategic-ideation](../strategic-ideation/SKILL.md) | Scope/strategy expand → tension → prune |
| [feature-discovery](../feature-discovery/SKILL.md) | Breadth-first zone triage before tickets exist |

## Prerequisites

- Open **`wayfinder:research`** ticket (`wayfinder:todo` + `wayfinder:hitl` in v1) with **Question**, **Done when**, and **Map** sections
- `gh` authenticated on the target repo (for posting findings comment)
- Parent map + decision log links available from ticket **Map** section

## Workflow

### 1. Load context

```text
gh issue view <research-num> --json body,title,url,labels
gh issue view <map-num> --json body,title,url
```

From the ticket: **Question**, **Done when** bullets, **Map** link, optional **Source hints** and **Perspectives**.

From the map: slug, decision log link, **Notes**, relevant **Not yet specified** fog.

**Gate:** Ticket must be labelled `wayfinder:research`. If the **Question** is scope/strategy shape rather than fact-gathering, hand off to [strategic-ideation](../strategic-ideation/SKILL.md) or [grill-me](../grill-me/SKILL.md) per [wayfinder routing](../wayfinder/REFERENCE.md#routing-table).

### 2. Investigate

Follow [behavior rules](REFERENCE.md#behavior-rules):

- **Primary sources first** — docs, specs, code, official APIs
- **Secondary sources** only when labeled; lower weight in Findings
- Seek **≥1 alternate viewpoint** for Viewpoints/alternatives
- If Coverage fails on a **valid premise**, run **one scope-expansion pass** then stop

Use web search, codebase exploration, and ticket **Source hints** as appropriate. Do not treat "couldn't find sources" as invalid premise.

### 3. Draft findings comment

Build comment per [output template](REFERENCE.md#findings-comment-template):

Summary → Findings → Gaps & follow-ups → Viewpoints/alternatives → Coverage → (optional **Invalid premise**) → Proposed tracker updates.

**Coverage:** Map each **Done when** bullet to `satisfied` / `partial` / `not satisfied` with brief evidence.

**Invalid premise:** Include only for logical/category ticket errors (wrong question type, impossible ask, mis-scoped ticket) — not missing data.

### 4. Post comment

Post as a **new top-level comment** on the research ticket:

```powershell
gh issue comment <research-num> --body-file path\to\findings.md
```

Each session posts a **new** comment; prior runs stay in the thread.

**Default:** Do not post Reconcile approval phrases or edit map/decision log from this skill.

### 5. Hand off

Tell the user:

- Findings are **non-binding** — review **Proposed tracker updates** before any map or log edits
- **Binding decisions:** run [grill-me](../grill-me/SKILL.md) on findings, or invoke wayfinder **Reconcile** after explicit human approval
- **Follow-up research:** start a new session on the same ticket (new comment); or close via wayfinder **Reconcile** when **Done when** is satisfied and user approves

End with: *Review the findings — invoke wayfinder **Reconcile** when ready to sync approved outcomes to the map.*

## Interaction rules

1. **Non-binding only** — never append decision-log rows; never post **Approved — reconcile and close**
2. **Comment, not file** — findings live on the ticket thread, not repo Markdown
3. **One scope-expansion pass** — if Coverage still fails after expansion, note gaps and stop
4. **HITL v1** — wait for human direction before follow-up research or Reconcile
5. **Invalid premise is rare** — distinguish from insufficient evidence

## Quick start

User: "Research ticket #8 on map #12."

Load ticket #8 + map #12 → investigate per behavior rules → post structured findings comment → hand off for human review.

See [REFERENCE.md](REFERENCE.md) for ticket template, output sections, behavior rules, and design defaults.

## Bootstrap

This skill was first implemented manually from [Task: implement research skill](https://github.com/KroniK907/skills/issues/21) (bundle [#20](https://github.com/KroniK907/skills/issues/20)). That task is the reference example for the workflow above.

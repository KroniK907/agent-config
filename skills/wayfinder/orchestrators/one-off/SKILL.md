---
name: one-off
description: one-off, wayfinder Route one-off, map To Do implementation, To Do ticket repo deliverables, small map-scoped slice, without define-bundle, one-off intent, materialize To Do ticket
agent-config-sync: true
---

# One-off

**HITL-only** entry for map **To Do** tickets that ship repo deliverables without the **define-bundle, create-tasks, Implementing** pipeline. Draft (or load) ticket, materialize with **`wf:approved`**, then the **implement-task** tail (task worktree and pull request) with documented gate waivers.

Detail: [REFERENCE.md](REFERENCE.md)

## Not this skill

| Skill | When instead |
|-------|----------------|
| [wayfinder](../../SKILL.md) | Chart, Materialize, Reconcile, Route only |
| [define-bundle](../../actions/define-bundle/SKILL.md) | GM cluster into an approved bundle |
| [create-tasks](../../actions/create-tasks/SKILL.md) | Split approved bundle → **Implementing** tasks |
| [implement-task](../implement-task/SKILL.md) | Direct pickup on **`wf:approved`** bundle tasks - stops on To Do tickets unless entered via **one-off** |
| Agent checklist or human | Trivial map errands with no repo deliverables (rename label, post comment, update tracker text) |

## Prerequisites

- Parent map (`wf:map`) with **To Do** table
- `gh` authenticated on the target repo
- Human declares one-off intent in chat

## Orchestration checklist

Run in order. **Stop at first gate failure** - post **Blocked** resolution per [REFERENCE § Resolution](REFERENCE.md#resolution-comment-one-off-variant); do not edit the repo.

### 1. Load context

```text
gh issue view <map-num> --json body,title,url
gh issue view <ticket-num> --json body,title,url,labels   # when ticket exists
```

From the map: slug, decision log link, **To Do**, optional **Dev branch:** line.

From the ticket (if loading existing): **Question**, **Done when**, **## Method**, **Blocked by**, **Status**, **PR:**

**Wrong entry:** If the user invoked [implement-task](../implement-task/SKILL.md) on a To Do ticket with no bundle parent, stop and redirect here - see [REFERENCE § Wrong-entry redirect](REFERENCE.md#wrong-entry-redirect).

### 2. Draft or skip

| Situation | Action |
|-----------|--------|
| New work | Chat-only draft per [REFERENCE § Ticket draft](REFERENCE.md#ticket-draft-chat-only); pause for human review |
| Complete existing ticket | Skip draft - verify **Question**, **Done when**, **## Method**, labels |
| Incomplete existing ticket | Fill gaps in chat; human confirms before materialize |

Optional: short [grill-me](../../ideation/grill-me/SKILL.md) pass on draft when scope is fuzzy.

### 3. Materialize

On human **`draft approved`** (or **`ticket approved`**):

1. Create or update ticket - labels `wf:todo` + `wf:task` + `wf:hitl` + **`wf:approved`**
2. Set body **Status:** `ready`; link map **Parent:** only (no **Parent bundle:**)
3. Append **To Do** row on map if new ticket
4. Proceed to build in the **same session** - do not wait for a separate pickup

See [REFERENCE § Ticket template](REFERENCE.md#ticket-template).

### 4. Git

Run [implement-task task worktree](../implement-task/REFERENCE.md#4-task-worktree) for this ticket. Same integration branch, same `task/{issue-num}-{slug}` worktree, same pull request. There is no separate one-off branch pattern.

### 5. Implementation tail

Follow the [implement-task orchestration checklist](../implement-task/SKILL.md#orchestration-checklist) with [gate waivers](REFERENCE.md#implement-task-gate-waivers) documented in REFERENCE only - **do not edit implement-task skill files**.

Includes: Method dispatch → [code-review](../../actions/code-review/SKILL.md) → push → resolution comment → **Status:** `awaiting-reconcile` + **`wf:needs-review`**.

Ticket stays on map **To Do** throughout - never **Implementing**.

### 6. Hand off

Tell the human:

- Task is **`awaiting-reconcile`** - review the resolution comment and the pull request on **PR:**
- Invoke wayfinder **Reconcile** with **`Approved - reconcile and close`** when accepted
- **Reconcile** moves map **To Do → Completed** gist and closes ticket - no Decision coverage **`implemented`** updates unless the ticket body explicitly references GM rows

## Interaction rules

1. **HITL only** - no AFK path; never add **`wf:afk`** or **`wf:afk-running`**
2. **To Do persistence** - ticket never moves to **Implementing**
3. **implement-task unchanged** - all gate waivers live in **one-off** REFERENCE
4. **Blocked by stops the run** - open blockers halt like implement-task
5. **Human closes the ticket** - requires **`Approved - reconcile and close`**

## Quick start

User: "One-off: ship {deliverable} on map #N."

Load map, draft ticket (or load #ticket), **`draft approved`**, materialize, task worktree, implement-task tail with waivers, resolution, Reconcile.

See [REFERENCE.md](REFERENCE.md) for templates, gate waivers, and worked example.

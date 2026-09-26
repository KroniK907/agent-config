---
name: one-off
description: one-off, wayfinder Route one-off, map To Do implementation, To Do ticket repo deliverables, small map-scoped slice, without define-bundle, one-off intent, materialize To Do ticket
agent-config-sync: true
---

# One-off

**HITL-only** path for a map **To Do** ticket that ships repo deliverables without the define-bundle → create-tasks → Implementing pipeline. Draft (or load) the ticket, materialize it with **`wf:approved`**, then run the implement-task tail - worktree, Method, code-review, PR - under the [gate waivers](REFERENCE.md#implement-task-gate-waivers).

| Instead | When |
|---------|------|
| [define-bundle](../../actions/define-bundle/SKILL.md) | The work belongs to a GM cluster |
| [implement-task](../implement-task/SKILL.md) | **Implementing** tasks from an approved bundle |
| Just do it | Trivial map errands with no repo deliverables (retitle, label, comment, map prose) |

Needs: a parent `wf:map` with **To Do**, `gh` auth.

## Checklist

If a gate fails, post the **Blocked** resolution ([one-off variant](REFERENCE.md#resolution-comment-one-off-variant)) and stop before touching the repo.

1. **Load** - the map (slug, log link, **To Do**) and the ticket if it exists (**Question**, **Done when**, **## Method**, **Blocked by**, **Status**, **PR:**).
2. **Draft** - for new work, show a [chat draft](REFERENCE.md#ticket-draft-chat-only). For an existing ticket, fill any gaps. A short [grill-me](../../ideation/grill-me/SKILL.md) pass helps when scope is fuzzy. When the user's request already pins down the deliverable, keep the draft brief and move on.
3. **Materialize** - once the user is happy with the draft: create or update the ticket with labels `wf:todo` `wf:task` `wf:hitl` `wf:approved`, **Status:** `ready`, map **Parent:** link only ([template](REFERENCE.md#ticket-template)); add a **To Do** row for a new ticket. Carry straight on to the build.
4. **Worktree** - [implement-task task worktree](../implement-task/REFERENCE.md#4-task-worktree).
5. **Implementation tail** - the [implement-task checklist](../implement-task/SKILL.md#checklist) from Method dispatch onward, with the waivers. Ends with **Status:** `awaiting-reconcile` + **`wf:needs-review`**.
6. **Hand off** - point the user at the PR. Once they're happy, Reconcile moves the **To Do** row to **Completed** and closes the ticket.

## Rules

1. **HITL only** - never `wf:afk` or `wf:afk-running`.
2. **Stays on To Do** - the ticket never moves to **Implementing**.
3. **Waivers live here** - implement-task files stay unchanged.
4. **Open blockers stop the run.**

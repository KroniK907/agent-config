# Create tasks reference

## Task issue template

Use for GitHub issue body. Title: `Task: {short name}`.

```markdown
# Task: {short name}

**Status:** draft | ready | awaiting-reconcile

## Parent bundle

[Bundle: {name}](bundle-url) - map [{FeatureName}:Map](map-url) - decision log [#N](log-url)

## What to build

{Concrete deliverables for this vertical slice; 1-2 paragraphs + bullet list when helpful.}

## Method

**Skill:** {skill-name}

{One line - why this Method for this slice. Propose from `wayfinder/` or `wayfinder/actions/`; repo-root one-offs only when human sets explicitly.}

## Decisions

Decision log: [#N](log-url). All `[global]` rows in the log apply. Load prose with `wf log --ids` or `--global`.

- **{MAP-SLUG}-GM-NNN** - one-line summary

## Outcomes / stories covered

N/A - meta/infra

- …

## Done when

- [ ] …

## Blocked by

 - 
```

For product tasks with user stories, replace the `N/A` line with story bullets and use **Done when** as acceptance criteria.

### Status lifecycle

| Status | Set by | Meaning |
|--------|--------|---------|
| `draft` | create-tasks | Split not yet accepted |
| `ready` | create-tasks on promote | Eligible for [implement-task](../../orchestrators/implement-task/SKILL.md) once it has **`wf:approved`** |
| `awaiting-reconcile` | implement-task only | PR open; resolution posted; waiting on review and Reconcile |

### Method field

Required at draft. Pick from **`wayfinder/**/<name>/SKILL.md`** in the pinned pack: **`write-code`** by default, **`prototype`** for throwaway exploration. Repo-root skills (`tdd`, `commit`, `writing-for-agents`, …) only when the user sets them. AFK tasks need a valid Method before **`wf:approved`** - see [Method validation](../../orchestrators/implement-task/REFERENCE.md#3-method-validation).

---

## Pickup order

Add **`wf:approved`** to one task at a time - serial pickup (especially AFK) expects one frontier task. A task is eligible when **Blocked by** is empty or every blocker is closed or `awaiting-reconcile`. Among eligible tasks pick:

1. The first in the bundle's stated rollout order
2. Foundations (infra, shared contracts) before consumers
3. The only remaining task

Tasks with open blockers stay `ready` without the label; [implement-task](../../orchestrators/implement-task/REFERENCE.md#unblock-and-handoff) adds it (plus the AFK pickup comment) when blockers clear.

---

## Split heuristics

| Signal | Split |
|--------|-------|
| Meta/infra bundle (single skill folder + doc updates) | **One task** - skip lengthy split debate |
| Distinct subsystems or deploy order | **2-3 tasks** - one vertical slice each |
| Task would exceed one focused agent session | Propose smaller slices |
| Hard dependency between slices | Earlier slice first; **Blocked by** on downstream task |

Default cap: 3 tasks per bundle unless the user asks for more. When unsure, propose fewer, larger slices. A single-task bundle still gets a task issue so it is tracked on **Implementing**.

---

## Decision coverage updates

### On promote

For each GM ID listed in bundle **Decisions** (not **Constraints**):

```markdown
| {MAP-SLUG}-GM-NNN | assigned | [#task](task-url) |
```

Linked issue moves from bundle (`scoped`) to task when work is assigned.

When multiple tasks cover one bundle, all share the same GM row until implementation Reconcile marks **`implemented`**. Link the task that will complete that GM, or the first task in build order.

### On implementation Reconcile

When wayfinder **Reconcile** closes a shipped task:

```markdown
| {MAP-SLUG}-GM-NNN | implemented | [#task](task-url) |
```

| Status | Meaning | Linked issue |
|--------|---------|--------------|
| `scoped` | In approved bundle | Bundle issue |
| `assigned` | Implementation task exists | Task issue |
| `implemented` | Shipped | Task issue (closed) |

---

## Map Implementing table

Add on promote:

```markdown
| [Task: {name}](task-url) | [#bundle](bundle-url) | HITL / AFK | ready | - |
```

Remove row and add **Completed** gist on implementation Reconcile.

---

## Parent linkage

Default (matches [define-bundle](../define-bundle/REFERENCE.md#bundle-issue-template)):

- **Required:** **Parent bundle** section in task body with markdown links to bundle, map, decision log
- **Optional:** native GitHub sub-issues (bundle → task); not required for workflow

---

## Route heuristics (for wayfinder)

Suggest **create-tasks** when:

- User asks to split or implement from an approved bundle
- Map has an approved **`wf:bundle`** with Status `approved` and empty or stale **Implementing** frontier
- define-bundle just approved a bundle

Prefer **define-bundle** when GM rows are still **`open`** and need bundling first.

After promote, suggest implement-task on the **`wf:approved`** task.

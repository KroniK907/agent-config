---
name: create-tasks
description: create-tasks, split bundle into tasks, implementation tasks from bundle, approved wf:bundle, approve tasks, wayfinder Route create-tasks, wf:task draft
agent-config-sync: true
---

# Create tasks

Split an **approved** bundle (`wf:bundle`, **Status:** `approved`) into **thin vertical slices** - each one session or one PR - as task issues, then put them on the map's **Implementing** frontier and mark the first one ready for pickup. No product code unless the task itself is the work.

| Instead | When |
|---------|------|
| [wayfinder](../../SKILL.md) | Chart, Materialize, Reconcile, Route |
| [define-bundle](../define-bundle/SKILL.md) | Bundle is still `draft`, or rows need bundling |
| [design-modules](../design-modules/SKILL.md) | Module shape is still open - recommended before splitting |

Needs: the approved bundle, a parent map with **Implementing** and **Decision coverage**, `gh` auth, labels `wf:task` / `wf:prototype`, `wf:hitl` / `wf:afk`, `wf:approved`.

## Workflow

### 1. Load context

```text
gh issue view <bundle-num> --json body,title,url
gh issue view <map-num> --json body,title,url
```

From the bundle: decision IDs, scope, boundaries, outcomes, and any module-design comments. Load binding prose with `wf log <log-num> --ids ...` / `--global`. From the map: slug, log link, **Implementing**, coverage for the bundle's GMs.

### 2. Propose the split

Give the user the task titles, what each ships, which outcomes it covers, HITL vs AFK, blocked-by, and why the boundaries fall there. Default to **one task** for a meta/infra bundle that fits one session; split only when subsystems, dependencies, or session size warrant it ([split heuristics](REFERENCE.md#split-heuristics)). Adjust from their feedback.

### 3. Draft task issues

`gh issue create` (or update drafts in place): title `Task: {short name}`, labels type + mode, **Status:** `draft`, body per the [task template](REFERENCE.md#task-issue-template). Link the parent bundle in the body. List the GM IDs each task honors with one-line summaries; all `[global]` rows apply.

Every task gets a **## Method** at draft - default **`write-code`**, **`prototype`** for throwaway exploration ([Method field](REFERENCE.md#method-field)). AFK tasks cannot be picked up without one.

Hand the drafts to the user and add **`wf:needs-review`**. If they already agreed to the split in step 2 and the drafts match it, go straight to step 4.

### 4. Promote

Once the user is happy with the split:

1. **Map Implementing** - one row per task (Ticket, Bundle, Mode, Status `ready`, Blocked by), via `wf map-edit`.
2. **Decision coverage** - each bundle-scoped GM → `assigned`, linked to the task that completes it.
3. **Tasks** - **Status:** `ready`; remove **`wf:needs-review`**.
4. **`wf:approved`** - add to **one** unblocked task, the next in build order ([pickup order](REFERENCE.md#pickup-order)). Tasks with open blockers wait; implement-task labels them when blockers clear.
5. **AFK pickup** - if that task is `wf:afk`, post the [AFK pickup comment](../../orchestrators/implement-task/references/afk-pickup-comment.md).

**Done when:** every task is on **Implementing** as `ready`, its GMs are `assigned`, and exactly one unblocked task carries **`wf:approved`**.

### 5. Hand off

Point the user at the **`wf:approved`** task: [implement-task](../../orchestrators/implement-task/SKILL.md) opens a worktree and a PR into `integrationBranch`. When it ships, wayfinder **Reconcile** closes it and marks its GMs `implemented`.

## Rules

1. **Approved bundles only** - a `draft` bundle goes back to define-bundle.
2. **Implementing is separate** - planning **To Do** tickets stay where they are.
3. **Globals are inherited** - never mark constraint-only GMs `assigned`.
4. **Leave the bundle alone** - no Status changes or re-scoping here.

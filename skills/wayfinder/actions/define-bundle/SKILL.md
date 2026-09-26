---
name: define-bundle
description: define-bundle, build bundles, approve bundle, GM cluster ready, group decisions into bundle, wayfinder Route define-bundle, wf:bundle draft
agent-config-sync: true
---

# Define bundle

Group **`open`** decision-log rows into a bundle issue (`wf:bundle`), and once the user is happy with it, mark it approved and claim its rows in map **Decision coverage**. Maps implement incrementally - don't wait for an empty To Do, a PRD, or cleared fog.

| Instead | When |
|---------|------|
| [wayfinder](../../SKILL.md) | Chart, Materialize, Reconcile, Route |
| [grill-me](../../ideation/grill-me/SKILL.md) | Unknowns still need binding GM rows |
| [design-modules](../design-modules/SKILL.md) | Shape module interfaces from an approved bundle |
| [create-tasks](../create-tasks/SKILL.md) | Split an approved bundle into tasks |

Needs: a `wf:map` with **Decision coverage** and a linked log, `gh` auth, label `wf:bundle`.

## Workflow

### 1. Load context

```text
go run <wayfinder>/utilities/wf/wf.go section <map-num> "Decision coverage" "To Do" "Implementing" "Notes"
go run <wayfinder>/utilities/wf/wf.go log <log-num> --list
```

Load full text for candidate IDs with `wf log <log-num> --ids A,B`. Classify each row by [global vs bundle-scoped](REFERENCE.md#global-vs-bundle-scoped-rows): coverage `open` without `[global]` are candidates; globals are constraints only; `scoped` / `assigned` / `implemented` are already claimed.

### 2. Propose the cluster

Give the user a build-oriented **name**, the **covered GM IDs** (one vertical slice - same subsystem, shared deliverables), the **rationale**, and **excluded rows**. Adjust from their feedback.

### 3. Draft the bundle issue

`gh issue create` (or update an existing draft in place) with title `Bundle: {short name}`, label `wf:bundle`, **Status:** `draft`, body per the [bundle template](REFERENCE.md#bundle-issue-template).

- **Decisions** and **Constraints** list IDs with one-line summaries and a log link - agents load prose with `wf log --ids` / `--global`.
- Link the parent map in the body; no sub-issues, no **Branch:** line. Each task gets its own worktree and PR later.

Hand the draft to the user and add **`wf:needs-review`**. If they already agreed to the cluster in step 2 and the draft matches it, go straight to step 4.

### 4. Approve

Once the user is happy with the bundle:

1. Set **Status:** `approved`; remove **`wf:needs-review`**.
2. `wf map-edit --coverage` each covered row to `scoped`, linked to the bundle. That row is the claim; the decision log is not edited.
3. Optionally add a Notes line linking the bundle, and comment on the bundle with what changed.

**Done when:** the bundle is `approved` and every covered row is `scoped` with a link to it.

Rows move to **Implementing** in [create-tasks](../create-tasks/SKILL.md), not here.

### 5. Hand off

Suggest [design-modules](../design-modules/SKILL.md) when module shape is still open, otherwise [create-tasks](../create-tasks/SKILL.md) - or implement straight from the bundle when one session covers it. Planning **To Do** may stay open.

## Rules

1. **One bundle per row** - a bundle-scoped GM is claimed by one bundle (`scoped`).
2. **Globals are inherited** - listed under Constraints, never `scoped`.
3. **The log is binding** - reference it; don't copy prose into the bundle.
4. **Bundle, don't build** - product code belongs to tasks, unless the bundle itself is meta/infra.

---
name: create-tasks
description: Split an approved wayfinder bundle into draft implementation tasks and sync map Implementing rows on scope and task approval. Use when a wayfinder:bundle is approved, wayfinder Route suggests create-tasks, or the user wants implementation tasks from a bundle.
---

# Create tasks

Split an **approved** bundle issue (`wayfinder:bundle`, Status `approved`) into agent-run-sized **thin vertical slices** as draft implementation issues, then on human approval sync map **Implementing** + **Decision coverage**. Does **not** implement product code unless the task itself is the work.

Runs on an approved bundle + parent map + decision log. After **`tasks approved`**, hand off to implementation; on ship, invoke [wayfinder](../wayfinder/SKILL.md) **Reconcile** to close the task and move coverage to **`implemented`**.

## Not this skill

| Skill | When instead |
|-------|----------------|
| [wayfinder](../wayfinder/SKILL.md) | Chart, Materialize, Reconcile, Route only |
| [define-bundle](../define-bundle/SKILL.md) | Coalesce GM rows into draft/approved bundles |
| [grill-me](../grill-me/SKILL.md) | Resolve unknowns → new binding GM rows |
| [write-a-prd](../write-a-prd/SKILL.md) | Small map-free scope only |

## Prerequisites

- Approved bundle issue (`wayfinder:bundle`, **Status:** `approved`)
- Parent map with **Implementing** table and **Decision coverage**
- `gh` authenticated on the target repo
- Labels `wayfinder:task` or `wayfinder:prototype`, `wayfinder:hitl` or `wayfinder:afk`, `wayfinder:approved`

## Workflow

### 1. Load context

```text
gh issue view <bundle-num> --json body,title,url
gh issue view <map-num> --json body,title,url
```

From the bundle: **Decisions** (covered GM rows), **Constraints**, scope summary, boundaries, outcomes.

From the map: slug, decision log link, **Implementing**, **Decision coverage** for bundle-scoped GM IDs.

**Gate:** Bundle **Status** must be `approved`. If `draft`, hand off to [define-bundle](../define-bundle/SKILL.md) for **`bundle approved`**.

### 2. Propose split

Present to the user:

- **Task count** and titles (1–3 typical; more only when deliverables are clearly separable)
- **Per task:** what ships, which bundle outcomes it covers, HITL vs AFK, blocked-by if any
- **Rationale** — why this slice boundaries

**Split heuristics:** See [REFERENCE.md](REFERENCE.md#split-heuristics). Default **one task** when the bundle is meta/infra and fits one session (e.g. a single skill folder). Split when subsystems, dependencies, or session size clearly warrant it.

Wait for confirmation or edits before creating/updating task issues.

### 3. Draft task issues

Create early with `gh issue create` or update drafts in place.

| Field | Value |
|-------|--------|
| Title | `Task: {short name}` |
| Labels | `wayfinder:task` or `wayfinder:prototype` + `wayfinder:hitl` or `wayfinder:afk` |
| Body | Per [REFERENCE.md](REFERENCE.md#task-issue-template) |
| **Status** | `draft` |

**Parent bundle:** link in **Parent bundle** section (body link only; native sub-issues optional).

**Decisions:** copy bundle **Decisions** rows **verbatim** + relevant **Constraints** the task must honor.

Fill **What to build**, **Outcomes/stories covered**, **Done when**, **Blocked by**.

Post or narrate drafts; end with: *Review the tasks — reply **scope approved** when the split is accepted, or request edits.*

**Default:** do not run **`scope approved`** writes without the explicit phrase.

### 4. On `scope approved`

When the user says **`scope approved`** (optionally naming task issues):

1. **Map Implementing** — one row per draft task (Ticket, Bundle link, Mode, Status `draft`, Blocked by)
2. **Map Decision coverage** — each bundle-scoped GM in bundle **Decisions** → **`assigned`**, **Linked issue** → task URL (when multiple tasks cover one GM, link the primary task or the task that completes that GM)
3. **Comment** on each task summarizing executed updates

Use `gh issue edit --body-file` for full map body replacements.

**Do not** add `wayfinder:approved` or set Status `ready` yet.

### 5. On `tasks approved`

When the user says **`tasks approved`**, **`task approved`**, or issue comment **`approved`** (per task or all):

1. **Task issue(s)** — set **Status:** `ready` in body; add label **`wayfinder:approved`**
2. **Map Implementing** — update Status column to `ready` for approved tasks
3. **Comment** on each task — ready for implementation

**Default:** do not start implementation without **`wayfinder:approved`** on the task.

### 6. Hand off

Tell the user:

- **Next:** implement from the ready task(s) — agent or human worker picks up **`wayfinder:approved`** tasks
- On completion: invoke wayfinder **Reconcile** with **`Approved — reconcile and close`** on the task issue

## Implementation Reconcile (after ship)

Owned by [wayfinder](../wayfinder/SKILL.md) **Reconcile**, not create-tasks. On **`Approved — reconcile and close`** for an implementation task:

1. Close task issue; remove **`wayfinder:approved`** label
2. **Map Implementing** — move row gist to **Completed**
3. **Decision coverage** — bundle-scoped GMs fully shipped by this task → **`implemented`**, linked issue stays task URL

See [REFERENCE.md](REFERENCE.md#implementation-reconcile).

## Interaction rules

1. **Approved bundle only** — never split draft bundles
2. **Planning To Do stays separate** — **Implementing** is the implementation frontier; do not move planning tickets
3. **Globals inherited** — copy bundle **Constraints** into each task **Decisions**; never mark constraint-only GMs **`assigned`**
4. **Draft early** — create task issues while split is still being refined; update in place
5. **No bundle edits** — create-tasks does not change bundle Status or re-scope GM rows

## Quick start

User: "Split bundle #17 into tasks on map #12."

Load bundle #17 + map #12 → propose split → create draft `wayfinder:task` issue(s) → user says **`scope approved`** → sync Implementing + coverage → user says **`tasks approved`** → add `wayfinder:approved` → implement.

See [REFERENCE.md](REFERENCE.md) for task template, approval phrases, and coverage updates.

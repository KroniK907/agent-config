# Create tasks reference

## Task issue template

Use for GitHub issue body. Title: `Task: {short name}`.

```markdown
# Task: {short name}

**Status:** draft | ready

## Parent bundle

[Bundle: {name}](bundle-url) · map [{FeatureName}:Map](map-url) · decision log [#N](log-url)

## What to build

{Concrete deliverables for this vertical slice; 1–2 paragraphs + bullet list when helpful.}

## Decisions

*(Bundle Decisions verbatim + inherited Constraints this task must honor)*

**{MAP-SLUG}-GM-NNN** — …

## Outcomes / stories covered

N/A — meta/infra

- …

## Done when

- [ ] …

## Blocked by

—
```

For product tasks with user stories, replace the `N/A` line with story bullets and use **Done when** as acceptance criteria.

---

## Approval phrases

| User says | Agent may |
|-----------|-----------|
| **scope approved** | Add map **Implementing** rows; set Decision coverage **`assigned`** + task links for bundle-scoped GMs |
| **tasks approved** / **task approved** / issue comment **approved** | Set task Status `ready`; add **`wayfinder:approved`**; update Implementing Status → `ready` |
| (edits requested) | Update draft task bodies in place; keep Status `draft`; no `wayfinder:approved` |
| (no approval) | Narrate or post drafts only; **do not** write Implementing or coverage |

Synonyms accepted if unambiguous: "approve the tasks", "approve task #N", "approved" on a specific task thread.

**Separate from Reconcile:** `scope approved` / `tasks approved` are owned by **create-tasks**. wayfinder **Reconcile** owns implementation task close + coverage **`implemented`**.

---

## Split heuristics

| Signal | Split |
|--------|-------|
| Meta/infra bundle (single skill folder + doc updates) | **One task** — skip lengthy split debate |
| Distinct subsystems or deploy order | **2–3 tasks** — one vertical slice each |
| Task would exceed one focused agent session | Propose smaller slices |
| Hard dependency between slices | Earlier slice first; **Blocked by** on downstream task |

Default cap: **3 tasks per bundle session** unless user asks for more. When unsure, propose fewer larger slices and let the user split further.

**Single-task bundles:** Still create a task issue — implements approval gates and **Implementing** tracking even when split is obvious.

---

## Decision coverage updates

### On `scope approved`

For each GM ID listed in bundle **Decisions** (not **Constraints**):

```markdown
| {MAP-SLUG}-GM-NNN | assigned | [#task](task-url) |
```

Linked issue moves from bundle (`scoped`) to task when work is assigned.

When multiple tasks cover one bundle, all share the same GM row until implementation Reconcile marks **`implemented`**. Link the task that will complete that GM, or the first task in build order.

### On implementation Reconcile

When wayfinder **Reconcile** closes a shipped task (`Approved — reconcile and close`):

```markdown
| {MAP-SLUG}-GM-NNN | implemented | [#task](task-url) |
```

Remove **`wayfinder:approved`** from the closed task. Move **Implementing** row gist to **Completed**.

| Status | Meaning | Linked issue |
|--------|---------|--------------|
| `scoped` | In approved bundle | Bundle issue |
| `assigned` | Implementation task exists | Task issue |
| `implemented` | Shipped | Task issue (closed) |

---

## Map Implementing table

Add on **`scope approved`**:

```markdown
| [Task: {name}](task-url) | [#bundle](bundle-url) | HITL / AFK | draft | — |
```

Update Status to **`ready`** on **`tasks approved`**.

Remove row and add **Completed** gist on implementation Reconcile.

---

## Parent linkage

Default (matches [define-bundle](../define-bundle/REFERENCE.md#bundle-issue-template)):

- **Required:** **Parent bundle** section in task body with markdown links to bundle, map, decision log
- **Optional:** native GitHub sub-issues (bundle → task); not required for workflow

---

## Design defaults

| Topic | Default |
|-------|---------|
| Split heuristics | 1–3 tasks; single task for meta/infra one-session bundles |
| Reconcile ownership | create-tasks documents gates; wayfinder Reconcile owns post-ship close + **`implemented`** |
| Sub-issues vs body links | Body **Parent bundle** link required; sub-issues optional |
| Single-task bundles | Still mint a task issue; skip split debate when obvious |
| Global re-tag pass | **Deferred** — separate Reconcile pass; not part of create-tasks |
| Map-free PRD workflow | **Constraint only** — write-a-prd Route stays in wayfinder REFERENCE; not bundled here |

---

## Implementation Reconcile

After an agent or human ships a **`wayfinder:approved`** task:

1. Post resolution comment on the task issue (what shipped, checklist against **Done when**)
2. Human: **`Approved — reconcile and close`**
3. wayfinder **Reconcile** executes:
   - Close task; remove **`wayfinder:approved`**
   - **Implementing** → **Completed** gist
   - Decision coverage **`implemented`** for GMs fully delivered by this task

create-tasks does **not** run implementation Reconcile — remind the user to invoke wayfinder when ready.

---

## Route heuristics (for wayfinder)

Suggest **create-tasks** when:

- User explicitly asks to split or implement from an approved bundle
- Map has an approved **`wayfinder:bundle`** with Status `approved` and empty or stale **Implementing** frontier
- **`define-bundle`** just completed **`bundle approved`** handoff

Prefer **define-bundle** when GM rows are still **`open`** and need bundling first.

After **`tasks approved`**, suggest picking up a **`wayfinder:approved`** task for implementation.

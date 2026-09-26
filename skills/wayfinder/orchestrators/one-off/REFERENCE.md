# One-off reference

---

## When to use

| Path | Use when |
|------|----------|
| **one-off** | Map **To Do** ticket with **repo deliverables** (skill folder, code, docs in target repo) - one session or thin vertical slice |
| **define-bundle, create-tasks, implement-task** | Multiple GM rows, **Implementing** frontier, AFK eligibility, one pull request per task |
| **Agent checklist or human** | Trivial map errands - no repo deliverables (retitle ticket, add label, post comment, edit map prose) |

**Scope gate:** is this map-scoped repo work that fits one session?

---

## Wrong-entry redirect

**implement-task** stops on To Do tickets without an approved bundle parent:

| Symptom | Fix |
|---------|-----|
| User says "implement task #N" on a **To Do** row | Redirect to **one-off** |
| Startup fails: missing bundle parent, bundle Status not `approved` | Re-enter via **one-off** - do not patch implement-task gates |
| User wants bundle pipeline | Hand off to [define-bundle](../../actions/define-bundle/SKILL.md) |

Narrate: *This ticket is on map **To Do** with no bundle - use **one-off** for map-scoped implementation, or **define-bundle** when the work belongs in an approved bundle.*

---

## Ticket draft (chat-only)

Show before `gh issue create`.

```markdown
## Draft one-off ticket

**Map:** [{FeatureName}:Map](map-url)
**Title:** Task: {short name}

## Question

<One line - what this session ships.>

## Done when

- [ ] <verifiable deliverable 1>
- [ ] <verifiable deliverable 2>

## Method

{skill-name}

<One line - why this Method (from `wayfinder/` or `wayfinder/actions/`; repo-root one-offs only when human sets explicitly).>

## Map

Parent: [{FeatureName}:Map](map-url)

**Type:** task | **Mode:** HITL

## Blocked by

<!-- Omit when none. Open blocker issues stop the run at startup. -->

- [#N Blocker title](url)

## Labels (on materialize)

`wf:todo` - `wf:task` - `wf:hitl` - `wf:approved`

## Status (on materialize)

`ready`

## Proposed git

Task worktree `task/{issue-num}-{slug}` and a pull request into `integrationBranch`, per [implement-task](../implement-task/REFERENCE.md#4-task-worktree). Slug from the ticket title, kebab-case, at most four words.
```

Ask for any changes, then materialize and build once the user is happy.

---

## Ticket template

GitHub issue body after materialize. Title: `Task: {short name}`.

```markdown
**Status:** ready

**PR:** {url after implement-task opens it}

## Question

<What this session ships.>

## Done when

- [ ] …

## Method

{skill-name}

## Map

Parent: [{FeatureName}:Map](map-url)

## Blocked by

<!-- Optional - open issues stop the run at startup -->
```

**No Parent bundle:** section. Ticket stays on map **To Do** for the entire run.

---

## implement-task gate waivers

When entered via **one-off**, apply these overrides to [implement-task startup gates](../implement-task/REFERENCE.md#startup-gates). **Waivers apply only on the one-off path** - implement-task skill files are not edited.

| Gate | Standard implement-task | One-off waiver |
|------|-------------------------|----------------|
| Bundle parent | Required; Status `approved` | **Waived** - no bundle. Still run [task worktree](../implement-task/REFERENCE.md#4-task-worktree) and open the pull request |
| Label **`wf:approved`** | Required at pickup | **Expected at materialize** - one-off adds it before build tail |
| **Blocked by** | Stops on open blockers | **Kept** - no waiver |
| AFK serial | AFK only | **N/A** - HITL only |
| Method validation | Per implement-task | **Kept** |
| Code review | After Method | **Kept** |
| Close task / remove **`wf:approved`** | Left to Reconcile | **Kept** |

Run the full implement-task tail: Method → code-review → push → resolution → **`awaiting-reconcile`**.

---

## Git

Use [implement-task task worktree](../implement-task/REFERENCE.md#4-task-worktree). One-off does not define a second checkout. Commit only in the task worktree. The pull request targets `integrationBranch`.

---

## Resolution comment (one-off variant)

Post on the one-off ticket at end-of-run. Same structure as [implement-task success template](../implement-task/references/resolution-comment.md#success-template) with these deltas:

- **Task** line only - omit **Bundle:** (no parent bundle)
- **Next** - review the pull request on **PR:**
- All other sections unchanged: Summary, Method, Code review, Commits, Done when, Reconcile

After posting: set **Status:** `awaiting-reconcile`; add **`wf:needs-review`**; keep **`wf:approved`**.

**Blocked runs:** Use implement-task [Blocked template](../implement-task/references/resolution-comment.md#blocked-template) - omit bundle URL; do not set **`awaiting-reconcile`**.

---

## Reconcile handoff

Wayfinder **Reconcile** closes the ticket once the user is happy with the PR: **To Do** row → **Completed** gist, **`wf:approved`** removed, coverage untouched unless the ticket body names GM rows. After `awaiting-reconcile`, suggest Reconcile; for rework, reset **Status** to `ready`.

---

## Worked example (generalized)

**Situation:** Map `{FeatureName}:Map` has a **To Do** row - ship a new sibling skill folder while planning tickets remain open.

1. **Declare** - Human: "One-off: add `{skill-name}` skill on `{FeatureName}:Map`."
2. **Draft** - Agent posts chat draft (Question, Done when, **## Method** `{skill-name}` or `writing-for-agents` for meta skills).
3. **Agree** - User signs off on the draft.
4. **Materialize** - `gh issue create` with labels; **Status:** `ready`; **To Do** row on map; **`wf:approved`** added.
5. **Git** - task worktree per implement-task; branch `task/{N}-{skill-name}`.
6. **Build** - Follow task **## Method** skill; record pre-Method SHA.
7. **Code review** - implement-task mode on diff since pre-Method SHA.
8. **Push** - commit in the worktree; push; open the pull request; write **PR:** on the ticket.
9. **Resolve** - post resolution comment; **Status:** `awaiting-reconcile`; **`wf:needs-review`**.
10. **Reconcile** - after the user reviews the PR → Completed gist; close ticket.

**Not this example:** Trivial "add a Note line to the map" - use agent checklist, not one-off.

---

## Route trigger (for wayfinder)

Suggest **one-off** when:

- Map **To Do** row is type **`task`** with repo deliverables
- User declares one-off intent or says "ship X on this map" without bundling
- User hit implement-task gates on a To Do ticket - redirect here

Do **not** suggest one-off for:

- **`Implementing`** rows - use [implement-task](../implement-task/SKILL.md)
- Trivial checklist errands
- AFK pickup (one-off is HITL-only permanently)

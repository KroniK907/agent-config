# Implement task reference

---

## Split from create-tasks

| Phase | Owner | Human gates | Agent may |
|-------|-------|-------------|-----------|
| Split bundle → draft tasks | [create-tasks](../create-tasks/SKILL.md) | **`scope approved`** | Create task issues; map **Implementing**; coverage **`assigned`** |
| Promote to pickup | create-tasks | **`tasks approved`** / **`approved`** | **Status:** `ready`; label **`wf:approved`** |
| Implementation run | **implement-task** | **`Approved — reconcile and close`** (wayfinder Reconcile) | Build via Method; push; resolution comment; **Status:** `awaiting-reconcile` |
| Close + coverage | [wayfinder](../SKILL.md) Reconcile | **`Approved — reconcile and close`** | Close task; remove **`wf:approved`**; **Implementing** → **Completed**; coverage **`implemented`** |

create-tasks never runs Method playbooks or posts implementation resolution comments. implement-task never splits bundles or adds **`wf:approved`**.

---

## Startup gates

Run **before any repository edit**. First failure → **Blocked** resolution ([resolution-comment.md](references/resolution-comment.md)); stop.

### 1. Load context

```text
gh issue view <task-num> --json body,title,url,labels
gh issue view <bundle-num> --json body,title,url,labels
```

From the task: **Status**, **Parent bundle**, **What to build**, **Decisions**, **Done when**, **Blocked by**, **## Method** (required for AFK).

From the bundle: map link, **Branch:** line, covered **Decisions**, **Status:** `approved`.

Optional: load map **Implementing** row for mode (HITL / AFK).

### 2. Task gates

| Gate | Fail when |
|------|-----------|
| Label | Missing **`wf:approved`** |
| Status | Not **`ready`** (re-runs: also accept **`awaiting-reconcile`** only when human explicitly restarted implementation in chat) |
| Blockers | Any **Blocked by** issue still **open** |
| Bundle parent | Bundle **Status** not **`approved`** |
| Mode label | Missing **`wf:hitl`** or **`wf:afk`** |

### 3. Method validation

Resolve Method skill path:

- Default pool: **`wayfinder/**/<name>/SKILL.md`** in the pinned skills pack
- Repo-root one-offs (`tdd`, `commit`, `writing-for-agents`, …) valid **only** when task **## Method** explicitly names them (human-set)

| Mode | Rule |
|------|------|
| **AFK** | **## Method** required; skill file must exist — fail closed if missing or invalid |
| **HITL** | Use task **## Method** by default; operator may **session-only override** in chat; persist override to task body only on explicit human request |

Record resolved Method name for resolution comment **Method** section.

### 4. Bundle branch

From bundle body **`Branch:`** line — pattern **`afk/bundle-{issue-num}-{slug}`**.

```powershell
git fetch origin
git checkout <branch>   # or git checkout -b <branch> when branch first needed
git pull --rebase origin <branch>
```

**Agents never open PRs.** Human opens PR when bundle is complete.

Gate failure: branch missing and cannot be created, checkout conflict, or pull failure → **Blocked**.

### 5. AFK serial gate — AFK only

Before repo edits on an AFK task:

1. If repo already has **`wf:afk-running`** on another open issue → **Blocked** (serial queue)
2. Else add **`wf:afk-running`** to **this** task

**Bypass:** Issue comment containing **`@cursor`** on the AFK task skips the serial gate for that pickup (still run other gates).

HITL tasks **never** add or remove **`wf:afk-running`**.

---

## Invariants

Throughout the run — violations are workflow bugs:

1. **Never close** the implementation task issue
2. **Never remove** **`wf:approved`** from the task you are implementing
3. **Never post** Reconcile approval phrases (**`Approved — reconcile and close`**, **`Approved — reconcile, keep open`**) on the task
4. **Never open PRs** — push to bundle branch only
5. **Orchestration vs build** — git push, code-review invoke, resolution comment, status `awaiting-reconcile`, unblock scan, AFK handoff stay in implement-task; Method skill owns product/doc deliverables; [code-review](../code-review/SKILL.md) owns review + obvious auto-fix

---

## Method dispatch

After all startup gates pass:

1. Record **`pre-method-sha`**: `git rev-parse HEAD`
2. Read resolved Method skill (frontmatter **`name`** must match task **## Method**)
3. Follow that skill's REFERENCE workflow for build work only
4. Honor task **Decisions**, **What to build**, and **Done when**
5. Return artifacts to orchestrator for resolution **Done when** table

Action skills live under **`wayfinder/actions/<name>/`** per [PATTERNS.md](../actions/PATTERNS.md). Map-frontier skills (`research`, `grill-me`, …) are **not** default Methods for implementation tasks unless explicitly set.

---

## Code review

After Method build work completes, **before commit/push**:

1. If `git diff pre-method-sha...HEAD` is empty → skip (no file changes)
2. Invoke [code-review](../code-review/SKILL.md) in **implement-task mode**:
   - Fixed point: **`pre-method-sha`**
   - Spec: task issue + bundle **Decisions** (already loaded)
3. Apply [auto-fix policy](../code-review/REFERENCE.md#auto-fix-policy) — code-review fixes obvious items in-repo
4. Capture [return artifact](../code-review/REFERENCE.md#implement-task-return-artifact) for resolution **Code review** section

Code-review does **not** replace human Reconcile — remaining Standards/Spec findings are for reviewer attention, not blockers unless the run cannot proceed (e.g. unfixable build break — narrate and stop before push).

HITL and AFK both run code-review automatically. No task **## Method** override.

---

## End-of-run sequence

After code-review completes:

### 1. Commit and push

- Commit on bundle branch — Method deliverables + code-review auto-fixes; messages reference task `#N` when helpful
- **`git push origin <bundle-branch>`**
- If push fails → narrate blocker; do not set **`awaiting-reconcile`** until push succeeds (or human directs otherwise)

### 2. Resolution comment

Post **Success** template from [references/resolution-comment.md](references/resolution-comment.md):

Sections: **Summary**, **Method**, **Code review**, **Commits**, **Done when**, **Next**, **Reconcile**

Paste code-review [return artifact](../code-review/REFERENCE.md#implement-task-return-artifact) under **Code review**. Map each task **Done when** bullet in the table with evidence.

### 3. Task status

```powershell
gh issue edit <task-num> --body-file path\to\updated-body.md
```

Set **Status:** `awaiting-reconcile`. Keep **`wf:approved`**. Add **`wf:needs-review`**. Do not remove **`wf:approved`**.

### 4. Unblock and handoff

See [Unblock and handoff](#unblock-and-handoff) below.

---

## HITL vs AFK

| Topic | HITL | AFK |
|-------|------|-----|
| Pickup | Human starts chat with task link / `#N` | Automation on **`wf:approved`** label add |
| **`wf:afk-running`** | Never | Acquire at startup; remove at end-of-run |
| Serial queue | N/A | One AFK run per repo; handoff after end-of-run |
| **`@cursor` bypass** | N/A | Comment on task skips serial gate |
| Method | Default from task; session override OK | **## Method** required; no override |
| Resolution + **`awaiting-reconcile`** | Same | Same |
| Unblock scan | Same | Same |
| Reconcile close | Human **`Approved — reconcile and close`** | Same — automation never closes task |

Task bodies are **identical** for HITL and AFK. Mode is label-only.

---

## Unblock and handoff

After success resolution and **`awaiting-reconcile`**:

### Dependent unblock (HITL + AFK)

Scan implementation tasks (map **Implementing** or bundle siblings) that list this task in **Blocked by**:

- When **all** blockers for a dependent are **closed** or **`awaiting-reconcile`** / reconciled as shipped, and dependent **Status** is **`ready`** with **`scope approved`** already applied:
  - Add **`wf:approved`** to the dependent (create-tasks deferred label when blocked; implement-task restores when unblocked)
- Do **not** start the dependent automatically in HITL unless the human asks

create-tasks owns deferring **`wf:approved`** while blockers exist; implement-task performs the **add** when blockers clear.

### AFK serial handoff (AFK only)

1. Remove **`wf:afk-running`** from the current task
2. Find next eligible AFK task: **`wf:approved`**, **`ready`**, unblocked, no **`afk-running`** elsewhere
3. Hand off to automation (next pickup) — do **not** implement the next task in the same orchestration run unless explicitly configured

If no eligible task, queue idle.

---

## Status lifecycle

| Status | Set by | Meaning |
|--------|--------|---------|
| `draft` | create-tasks | Task not approved for pickup |
| `ready` | create-tasks on **`tasks approved`** | Eligible for implement-task |
| `awaiting-reconcile` | **implement-task** end-of-run | Work pushed; resolution posted; awaits human Reconcile |
| (closed) | wayfinder Reconcile | Human **`Approved — reconcile and close`** |

Only implement-task sets **`awaiting-reconcile`**.

---

## Route heuristics (for wayfinder)

Suggest **implement-task** when:

- Map **Implementing** row has **`wf:approved`** and **Status:** `ready`
- User says "implement task #N" or "pick up #N" on an approved implementation task
- AFK automation triggers on **`wf:approved`** add

After **`awaiting-reconcile`**, suggest wayfinder **Reconcile** — not another implement-task pass unless human requests rework (reset **Status** to **`ready`** explicitly before re-run).

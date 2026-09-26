# Implement task reference

## Startup gates

Run **before any repository edit**. First failure → **Blocked** resolution ([resolution-comment.md](references/resolution-comment.md)); stop.

### 1. Load context

```text
gh issue view <task-num> --json body,title,url,labels
gh issue view <bundle-num> --json body,title,url,labels
```

From the task: **Status**, **Parent bundle**, **What to build**, **Decisions**, **Done when**, **Blocked by**, **## Method** (required for AFK).

From the bundle: map link, covered **Decisions**, **Status:** `approved`.

Optional: load map **Implementing** row for mode (HITL / AFK).

### 2. Task gates

| Gate | Stop when |
|------|-----------|
| Label | Missing **`wf:approved`** |
| Status | Not **`ready`** (or **`awaiting-reconcile`** when the user asked for rework) |
| Blockers | Any **Blocked by** issue still **open** |
| Bundle parent | Bundle **Status** not **`approved`** |
| Mode label | Missing **`wf:hitl`** or **`wf:afk`** |

### 3. Method validation

Resolve Method skill path:

- Default pool: **`wayfinder/**/<name>/SKILL.md`** in the pinned skills pack
- Repo-root one-offs (`tdd`, `commit`, `writing-for-agents`, …) valid **only** when task **## Method** explicitly names them (human-set)

| Mode | Rule |
|------|------|
| **AFK** | **## Method** required; skill file must exist - stop if missing or invalid |
| **HITL** | Use task **## Method**; the user may override it for the session. Write the override to the task body if they want it kept |

Record resolved Method name for resolution comment **Method** section.

### 4. Task worktree

**Integration branch** is `.cursor/agent-manifest.json` field `integrationBranch` (for example `dev`, `staging`, or `main`). Read it before any git write.

**Done when:** the field names a branch that exists on `origin`. If it is missing, ask the user (HITL) or post **Blocked** (AFK) rather than guessing.

Each task gets its own worktree. The pull request into `integrationBranch` is the review artifact.

**Done when:** the shell cwd is a linked worktree whose current branch is not `integrationBranch`.

Stay in the current worktree when `git branch --show-current` is already not `integrationBranch`. A host worktree with another branch name counts. Do not add a second worktree.

Create a worktree only when the checkout is `integrationBranch`. Branch name `task/{issue-num}-{slug}`. Slug is kebab-case from the task title, at most four words.

```powershell
git fetch origin
$base = (Get-Content .cursor/agent-manifest.json -Raw | ConvertFrom-Json).integrationBranch
$branch = "task/<issue-num>-<slug>"
$wt = Join-Path (Split-Path (git rev-parse --show-toplevel) -Parent) ("task-<issue-num>-<slug>")
git worktree add -b $branch $wt "origin/$base"
Set-Location $wt
```

Gate failure: `integrationBranch` missing, `origin/$base` missing, or `git worktree add` fails → **Blocked**. Do not commit on the integration branch.

### 5. AFK serial gate - AFK only

Before repo edits on an AFK task:

1. If repo already has **`wf:afk-running`** on another open issue → **Blocked** (serial queue)
2. Else add **`wf:afk-running`** to **this** task

**Bypass:** An issue comment containing **`afk-serial-bypass`** on the AFK task skips the serial gate for that pickup (still run other gates). A comment containing **`@cursor`** is the same bypass.

HITL tasks **never** add or remove **`wf:afk-running`**.

---

## Invariants

1. The task stays open with **`wf:approved`** - Reconcile closes it after the user reviews the PR.
2. Commit only in the task worktree; ship through a PR against `integrationBranch`.
3. implement-task owns push, code-review, resolution, status, unblock, and AFK handoff; the Method skill owns the deliverables.

---

## Method dispatch

After all startup gates pass:

1. Record **`pre-method-sha`**: `git rev-parse HEAD`
2. Read resolved Method skill (frontmatter **`name`** must match task **## Method**)
3. Follow that skill's REFERENCE workflow for build work only
4. Honor task **Decisions**, **What to build**, and **Done when**
5. Return artifacts to orchestrator for resolution **Done when** table

Action skills live under **`wayfinder/actions/<name>/`** per [PATTERNS.md](../../actions/PATTERNS.md). Default **`wf:task`** Method is **`write-code`**; **`prototype`** for throwaway demos. Ideation skills (`research`, `grill-me`, `design-modules`, …) may also live under **`wayfinder/actions/`** or **`wayfinder/ideation/`** - valid as **## Method** when explicitly set.

---

## Code review

After Method build work completes, **before commit/push**:

1. If `git diff pre-method-sha...HEAD` is empty → skip (no file changes)
2. Invoke [code-review](../../actions/code-review/SKILL.md) in **implement-task mode**:
 - Fixed point: **`pre-method-sha`**
 - Spec: task issue + bundle **Decisions** (already loaded)
3. Apply [auto-fix policy](../../actions/code-review/REFERENCE.md#auto-fix-policy) - code-review fixes obvious items in-repo
4. Capture [return artifact](../../actions/code-review/REFERENCE.md#implement-task-return-artifact) for resolution **Code review** section

Remaining Standards/Spec findings go to the reviewer in the resolution comment. They block only when the run cannot proceed (e.g. an unfixable build break - say so and stop before push).

HITL and AFK both run code-review automatically. No task **## Method** override.

---

## End-of-run sequence

After code-review completes:

### 1. Commit, push, and open the pull request

- If `git diff origin/<integrationBranch>...HEAD` is empty, skip push and the pull request. Post the resolution comment and say no files changed.
- Otherwise commit in the task worktree. Method deliverables and code-review auto-fixes. Messages reference task `#N` when helpful.
- `git push -u origin HEAD`
- `gh pr create --base <integrationBranch>` when no open PR exists for this head. If `gh pr view` already returns one, use that URL. The head is the current branch, including a host worktree branch that is not `task/{issue-num}-{slug}`.
- Write **PR:** and the URL on the task body.
- **Done when:** `gh pr view --json url,baseRefName` shows that URL and `baseRefName` equals `integrationBranch`, or the diff against `integrationBranch` was empty.
- If push or PR creation fails, say why and hold off on **`awaiting-reconcile`** until both succeed.

### 2. Resolution comment

Post **Success** template from [references/resolution-comment.md](references/resolution-comment.md):

Sections: **Summary**, **Method**, **Code review**, **Commits**, **Done when**, **Next**

Paste code-review [return artifact](../../actions/code-review/REFERENCE.md#implement-task-return-artifact) under **Code review**. Map each task **Done when** bullet in the table with evidence.

### 3. Task status

```powershell
gh issue edit <task-num> --body-file path\to\updated-body.md
```

Set **Status:** `awaiting-reconcile`, keep **`wf:approved`**, and add **`wf:needs-review`**.

### 4. Unblock and handoff

See [Unblock and handoff](#unblock-and-handoff) below.

---

## HITL vs AFK

| Topic | HITL | AFK |
|-------|------|-----|
| Pickup | Human starts chat with task link / `#N` | Repo automation on issue comment **`Approved - AFK implement`** ([afk-pickup-comment.md](references/afk-pickup-comment.md)) |
| **`wf:approved`** | Startup gate; added by create-tasks or unblock | Same; the pickup comment is what triggers automation |
| **`wf:afk-running`** | Never | Acquire at startup; remove at end-of-run |
| Serial queue | N/A | One AFK run per repo; handoff after end-of-run |
| **Serial bypass** | N/A | Comment **`afk-serial-bypass`** (or **`@cursor`**) on the task skips the serial gate |
| Method | Default from task; session override OK | **## Method** required; no override |
| Resolution + **`awaiting-reconcile`** | Same | Same |
| Unblock scan | Label **`wf:approved`** only | Label + AFK pickup comment |
| Close | Reconcile, after the user reviews the PR | Same - automation never closes the task |

Task bodies are **identical** for HITL and AFK. Mode is label-only.

---

## Unblock and handoff

After success resolution and **`awaiting-reconcile`**:

### Dependent unblock (HITL + AFK)

Scan implementation tasks (map **Implementing** or bundle siblings) that list this task in **Blocked by**:

- When **all** blockers for a dependent are **closed** or **`awaiting-reconcile`**, and the dependent is **`ready`**:
 - Add **`wf:approved`** to the dependent
 - **AFK only** (`wf:afk`): post pickup comment per [references/afk-pickup-comment.md](references/afk-pickup-comment.md)
- In HITL, leave starting the dependent to the user

### AFK serial handoff (AFK only)

1. Remove **`wf:afk-running`** from the current task
2. Find next eligible AFK task: **`wf:approved`** or ready to receive it, **`ready`**, unblocked, no **`afk-running`** elsewhere, **`wf:afk`** label
3. If next task lacks **`wf:approved`**: add label, then post AFK pickup comment per [references/afk-pickup-comment.md](references/afk-pickup-comment.md)
4. Hand off to automation (comment trigger) - do **not** implement the next task in the same orchestration run unless explicitly configured

If no eligible task, queue idle.

---

## Status lifecycle

See [create-tasks status lifecycle](../../actions/create-tasks/REFERENCE.md#status-lifecycle). Only implement-task sets `awaiting-reconcile`; Reconcile closes the task.

---

## Route heuristics (for wayfinder)

Suggest **implement-task** when:

- Map **Implementing** row has **`wf:approved`** and **Status:** `ready`
- User says "implement task #N" or "pick up #N" on an approved implementation task
- AFK automation triggers on issue comment **`Approved - AFK implement`** (see [references/afk-pickup-comment.md](references/afk-pickup-comment.md))

After **`awaiting-reconcile`**, suggest wayfinder **Reconcile**. For rework, reset **Status** to `ready` and re-run.

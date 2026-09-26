# AFK pickup comment

Post when an AFK implementation task becomes eligible for **implement-task** pickup - by [create-tasks](../../../actions/create-tasks/SKILL.md) on promote, and by **implement-task** when unblocking dependents or handing off to the next AFK task.

This is a machine trigger the agent posts, not a human sign-off. Repo automation filters issue comments on this exact line (case-sensitive):

```text
Approved - AFK implement
```

Label **`wf:approved`** is the implement-task startup gate. If the automation host gains a label-added trigger, the comment can go away once every app repo migrates.

## When to post

Post **only** when **all** of the following hold:

1. Task has label **`wf:afk`** (not HITL)
2. Task **Status:** `ready`
3. **Blocked by** cleared (empty, or every blocker **closed** / **`awaiting-reconcile`**)
4. Agent adds (or confirms) label **`wf:approved`** on the same pickup decision

| Owner | Moment |
|-------|--------|
| **create-tasks** | On promote, when adding **`wf:approved`** to the one eligible AFK task |
| **implement-task** | After success resolution when a dependent becomes unblocked |
| **implement-task** | AFK serial handoff - next eligible task in queue after removing **`wf:afk-running`** from the finished task |

Post once per pickup decision; repost only when the user resets pickup.

## Comment template

Post as a **new top-level comment** on the AFK task. The trigger phrase must appear **verbatim on its own line** so automation filters match reliably.

```markdown
## AFK implement pickup

**Approved - AFK implement**

Task **Status:** `ready` - label **`wf:approved`** added for reviewer visibility.
**Method:** `{method-name}` - task worktree and pull request per implement-task.

<!-- Automation: issue comment trigger v1. Future: wf:approved label add may replace comment trigger. -->
```

Replace `{method-name}` with the task **## Method** value. Omit the HTML comment when posting via `gh issue comment` if your team prefers a clean thread - automation needs only the trigger line.

## gh example

```powershell
gh issue comment <task-num> --body-file path\to\afk-pickup.md
```

```bash
gh issue comment <task-num> --body-file path/to/afk-pickup.md
```

## Pair with label

Always add **`wf:approved`** in the same pickup decision (before or after the comment):

```powershell
gh issue edit <task-num> --add-label "wf:approved"
```

HITL tasks (**`wf:hitl`**) never receive this comment - human starts implement-task in chat instead.

---
name: implement-task
description: implement-task, wf:approved, wf:approved task ready, pick up task, implement bundle task, wayfinder Route implement-task, HITL task pickup, AFK task pickup, awaiting-reconcile
agent-config-sync: true
---

# Implement task

Orchestrates a **`wf:approved`** implementation task: startup gates, task worktree, **Method** dispatch, [code-review](../../actions/code-review/SKILL.md), push, pull request, [resolution comment](references/resolution-comment.md), **Status:** `awaiting-reconcile`. The build itself belongs to the Method skill.

Detail: [REFERENCE.md](REFERENCE.md). Templates: [resolution comment](references/resolution-comment.md), [AFK pickup comment](references/afk-pickup-comment.md).

Splitting bundles and first labelling tasks is [create-tasks](../../actions/create-tasks/SKILL.md). Closing a shipped task is wayfinder **Reconcile**. A map **To Do** ticket with no bundle goes through [one-off](../one-off/SKILL.md).

## Checklist

Run in order. If a startup gate fails, post the **Blocked** resolution and stop before touching the repo.

1. **Load** - task issue + parent bundle; fetch decision prose with `wf log <log-num> --ids ...` and `--global`.
2. **Startup gates** - [REFERENCE § Startup gates](REFERENCE.md#startup-gates).
3. **Worktree** - [task worktree](REFERENCE.md#4-task-worktree) off `integrationBranch`.
4. **Method** - record `pre-method-sha`, then follow the task's **## Method** skill.
5. **Code review** - [code-review](../../actions/code-review/SKILL.md) in implement-task mode on `<pre-method-sha>...HEAD` ([REFERENCE § Code review](REFERENCE.md#code-review)).
6. **Push** - commit in the worktree, push, open a PR into `integrationBranch`, write **PR:** on the task body.
7. **Resolve** - post the success resolution; set **Status:** `awaiting-reconcile`; add **`wf:needs-review`**.
8. **Unblock** - label dependents whose blockers are now clear ([REFERENCE § Unblock and handoff](REFERENCE.md#unblock-and-handoff)).
9. **AFK only** - release **`wf:afk-running`** and hand off to the next eligible AFK task.

The task stays open with **`wf:approved`** until the user has reviewed the PR and Reconcile closes it - the implementing agent doesn't close its own work.

## Hand off

Tell the user the PR is ready for review. Once they're happy with it, wayfinder **Reconcile** closes the task.

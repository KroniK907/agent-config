# Cursor automation prompt - wayfinder AFK implementation pickup

Copy this prompt into **one repo-scoped Cursor automation** per implementation repository. The agent opens one pull request per task into `integrationBranch`.

## Trigger

| Setting | Value |
|---------|--------|
| Event | Issue comment created |
| Match | Comment body contains **`Approved - AFK implement`** (exact phrase, case-sensitive) |
| Scope | This repository only |

**Why comment, not label:** Cursor automations v1 support PR label triggers reliably; **issue label added** for **`wf:approved`** is not available yet. Skills still add **`wf:approved`** for human reviewers and implement-task startup gates; the comment phrase is the automation trigger. When issue-label triggers ship, app repos may switch automation to **`wf:approved`** label add and optionally stop posting pickup comments - see [afk-pickup-comment.md](../../orchestrators/implement-task/references/afk-pickup-comment.md).

Optional bypass: an issue comment containing **`@cursor`** on an AFK task skips the serial queue gate for that pickup ([implement-task REFERENCE](../../orchestrators/implement-task/REFERENCE.md#5-afk-serial-gate--afk-only)).

## Prompt

```text
You are picking up a wayfinder implementation task in AFK (unattended) mode.

1. Read the triggered issue number from the automation context (the issue that received the pickup comment).
2. Invoke the **implement-task** skill on that issue (`/implement-task` or load wayfinder/orchestrators/implement-task/SKILL.md from the installed skills directory; Cursor Cloud copies that directory to ~/.cursor/skills/).
3. Follow implement-task exactly - orchestration only:
 - Run startup gates; stop on first failure (Status, labels, Method, integration branch, AFK serial lock)
 - Create the task worktree from integrationBranch
 - Method dispatch from task **## Method** (required for AFK - no session override)
 - Code review after Method (implement-task mode)
 - Commit, push, and open a pull request into integrationBranch
 - Post resolution comment; set task **Status:** awaiting-reconcile
 - Keep the task open with wf:approved until wayfinder Reconcile closes it
4. On startup gate failure, post Blocked resolution per implement-task references/resolution-comment.md - no repo edits.
5. At end-of-run: remove wf:afk-running from this task; hand off to next eligible AFK task (wf:approved + pickup comment Approved - AFK implement) if any.

Contract references:
- implement-task: wayfinder/orchestrators/implement-task/SKILL.md + REFERENCE.md
- AFK pickup comment: wayfinder/orchestrators/implement-task/references/afk-pickup-comment.md
- Task body template (identical for HITL and AFK): wayfinder/actions/create-tasks/REFERENCE.md
- AFK bootstrap checklist: wayfinder/utilities/AFK-BOOTSTRAP.md

GH_TOKEN: use the Cursor dashboard secret by default; optional per-repo override via `environment.json` `env.GH_TOKEN` when dashboard secret is not set.
```

## Human steps after automation runs

1. Review the **Implementation resolution** comment on the task issue.
2. Review the pull request linked from the task **PR:** line.
3. Invoke wayfinder **Reconcile** with **`Approved - reconcile and close`** when accepted.

## PR creation

**ON.** implement-task opens one pull request per task. The base is `integrationBranch` in `.cursor/agent-manifest.json`.

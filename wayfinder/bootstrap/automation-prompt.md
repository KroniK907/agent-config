# Cursor automation prompt — wayfinder AFK implementation pickup

Copy this prompt into **one repo-scoped Cursor automation** per implementation repository. Do **not** enable PR creation in the automation — agents push to the bundle branch only ([WF-ECO-GM-027](https://github.com/KroniK907/skills/issues/11)).

## Trigger

| Setting | Value |
|---------|--------|
| Event | Issue label added |
| Label | `wayfinder:approved` |
| Scope | This repository only |

Optional bypass: an issue comment containing `@cursor` on an AFK task skips the serial queue gate for that pickup ([implement-task REFERENCE](../implement-task/REFERENCE.md#5-afk-serial-gate--afk-only)).

## Prompt

```text
You are picking up a wayfinder implementation task in AFK (unattended) mode.

1. Read the triggered issue number from the automation context.
2. Invoke the **implement-task** skill on that issue (`/implement-task` or load wayfinder/implement-task/SKILL.md from ~/.cursor/skills/).
3. Follow implement-task exactly — orchestration only:
   - Fail-closed startup gates (Status, labels, Method, bundle branch, AFK serial lock)
   - Checkout/pull bundle branch from parent bundle **Branch:** line
   - Method dispatch from task **## Method** (required for AFK — no session override)
   - Code review after Method (implement-task mode)
   - Commit + push to bundle branch — **never open PRs**
   - Post resolution comment; set task **Status:** awaiting-reconcile
   - Never close the task; never remove wayfinder:approved; never post Reconcile approval phrases
4. On startup gate failure, post Blocked resolution per implement-task references/resolution-comment.md — no repo edits.
5. At end-of-run: remove wayfinder:afk-running from this task; hand off to next eligible AFK task if any.

Contract references:
- implement-task ([#29](https://github.com/KroniK907/skills/issues/29)): wayfinder/implement-task/SKILL.md + REFERENCE.md
- Task body template (identical for HITL and AFK): wayfinder/create-tasks/REFERENCE.md
- AFK bootstrap checklist: wayfinder/AFK-BOOTSTRAP.md

GH_TOKEN: use the Cursor dashboard secret by default; optional per-repo override via `environment.json` `env.GH_TOKEN` when dashboard secret is not set.
```

## Human steps after automation runs

1. Review the **Implementation resolution** comment on the task issue.
2. Review diff on the bundle branch named in the parent bundle issue.
3. Invoke wayfinder **Reconcile** with **`Approved — reconcile and close`** when accepted.

## PR creation

**OFF.** Bundle work accumulates on `afk/bundle-{N}-{slug}` until a human opens one PR for the complete bundle.

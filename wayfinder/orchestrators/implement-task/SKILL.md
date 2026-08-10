---
name: implement-task
description: Orchestrate pickup and completion of wf:approved implementation tasks - startup gates, bundle-branch git, Method dispatch, code-review (auto-fix obvious), resolution comment, awaiting-reconcile, and AFK queue handoff. Use when a wf:approved task is ready to implement (HITL or AFK), or wayfinder Route suggests implement-task.
---

# Implement task

**Orchestration-only** entry for **`wf:approved`** implementation tasks. Fail-closed startup → bundle-branch git → **Method dispatch** → **[code-review](../../actions/code-review/SKILL.md)** (auto-fix obvious; defer rest) → push → [resolution comment](references/resolution-comment.md) → **`Status: awaiting-reconcile`**. Does **not** close the task, remove **`wf:approved`**, or post Reconcile approval phrases.

Detail: [REFERENCE.md](REFERENCE.md) - resolution templates: [references/resolution-comment.md](references/resolution-comment.md) - AFK pickup comment: [references/afk-pickup-comment.md](references/afk-pickup-comment.md)

## Not this skill

| Skill | When instead |
|-------|----------------|
| [create-tasks](../../actions/create-tasks/SKILL.md) | Split bundle, **`scope approved`**, **`tasks approved`**, add **`wf:approved`** |
| [wayfinder](../../SKILL.md) | Chart, Materialize, **Reconcile** (close task after human **`Approved - reconcile and close`**) |
| [actions/*](../actions/PATTERNS.md) | Build work delegated via task **## Method** |

## Orchestration checklist

Run in order. **Stop at first gate failure** - post **Blocked** resolution per [references/resolution-comment.md](references/resolution-comment.md); do not edit the repo.

1. **Load** - task issue + parent bundle (map link, decision log, **Branch:**, **Decisions**)
2. **Startup gates** - [REFERENCE § Startup gates](REFERENCE.md#startup-gates) (Status, labels, Method, bundle branch, AFK serial)
3. **Git** - checkout/pull bundle branch from bundle **Branch:** line; create if missing
4. **Method dispatch** - record pre-Method `HEAD`; load and follow task **## Method** skill (HITL session override allowed; AFK requires valid Method)
5. **Build** - action skill owns deliverables; orchestrator does not duplicate build steps
6. **Code review** - [code-review](../../actions/code-review/SKILL.md) in **implement-task mode** on `<pre-method-sha>...HEAD`; auto-fix obvious mistakes; capture [return artifact](../../actions/code-review/REFERENCE.md#implement-task-return-artifact) ([REFERENCE § Code review](REFERENCE.md#code-review))
7. **Push** - commit Method + auto-fixes on bundle branch; push to remote
8. **Resolve** - post success resolution comment (include **Code review** section); set body **Status:** `awaiting-reconcile` (keep **`wf:approved`**); add label **`wf:needs-review`**
9. **Unblock** - scan dependents; for each cleared AFK dependent: add **`wf:approved`** + post AFK pickup comment ([references/afk-pickup-comment.md](references/afk-pickup-comment.md)); HITL dependents: label only ([REFERENCE § Unblock and handoff](REFERENCE.md#unblock-and-handoff))
10. **AFK only** - remove **`wf:afk-running`**; serial handoff to next eligible AFK task (**`wf:approved`** + pickup comment per [afk-pickup-comment.md](references/afk-pickup-comment.md))

**Invariants:** Never close task - never remove own **`wf:approved`** - never post **`Approved - reconcile and close`**

## Hand off

Tell the human (or AFK queue): task is **`awaiting-reconcile`** - review resolution comment, then invoke wayfinder **Reconcile** with **`Approved - reconcile and close`** when accepted.

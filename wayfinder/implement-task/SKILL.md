---
name: implement-task
description: Orchestrate pickup and completion of wayfinder:approved implementation tasks — startup gates, bundle-branch git, Method dispatch, resolution comment, awaiting-reconcile, and AFK queue handoff. Use when a wayfinder:approved task is ready to implement (HITL or AFK), or wayfinder Route suggests implement-task.
---

# Implement task

**Orchestration-only** entry for **`wayfinder:approved`** implementation tasks. Fail-closed startup → bundle-branch git → **Method dispatch** to the task's action skill → push → [resolution comment](references/resolution-comment.md) → **`Status: awaiting-reconcile`**. Does **not** close the task, remove **`wayfinder:approved`**, or post Reconcile approval phrases.

Detail: [REFERENCE.md](REFERENCE.md) · resolution templates: [references/resolution-comment.md](references/resolution-comment.md)

## Not this skill

| Skill | When instead |
|-------|----------------|
| [create-tasks](../create-tasks/SKILL.md) | Split bundle, **`scope approved`**, **`tasks approved`**, add **`wayfinder:approved`** |
| [wayfinder](../SKILL.md) | Chart, Materialize, **Reconcile** (close task after human **`Approved — reconcile and close`**) |
| [actions/*](../actions/PATTERNS.md) | Build work delegated via task **## Method** |

## Orchestration checklist

Run in order. **Stop at first gate failure** — post **Blocked** resolution per [references/resolution-comment.md](references/resolution-comment.md); do not edit the repo.

1. **Load** — task issue + parent bundle (map link, decision log, **Branch:**, **Decisions**)
2. **Startup gates** — [REFERENCE § Startup gates](REFERENCE.md#startup-gates) (Status, labels, Method, bundle branch, AFK serial)
3. **Git** — checkout/pull bundle branch; create if missing per GM-027
4. **Method dispatch** — load and follow task **## Method** skill (HITL session override allowed; AFK requires valid Method)
5. **Build** — action skill owns deliverables; orchestrator does not duplicate build steps
6. **Push** — commit on bundle branch; push to remote
7. **Resolve** — post success resolution comment; set body **Status:** `awaiting-reconcile` (keep **`wayfinder:approved`**)
8. **Unblock** — scan dependents; add **`wayfinder:approved`** where **Blocked by** cleared ([REFERENCE § Unblock and handoff](REFERENCE.md#unblock-and-handoff))
9. **AFK only** — remove **`wayfinder:afk-running`**; serial handoff to next eligible AFK task

**Invariants:** Never close task · never remove own **`wayfinder:approved`** · never post **`Approved — reconcile and close`**

## Hand off

Tell the human (or AFK queue): task is **`awaiting-reconcile`** — review resolution comment, then invoke wayfinder **Reconcile** with **`Approved — reconcile and close`** when accepted.

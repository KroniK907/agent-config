---
name: design-module
description: Shape deep modules from bundle decisions or planning tickets - vocabulary, deepening, and design-it-twice exploration. HITL only. Posts module-design artifact as issue comment. Use after bundle approved before create-tasks, on wf:prototype interface exploration, or when user wants module interface design.
---

# Design module

Turn bundled **Decisions** (or a planning ticket **Question**) into a **deep module** shape: small interface, clear **seam**, dependency category, and a recommended design with rationale. Output is a structured comment on the bundle or ticket - not repo edits unless the human explicitly requests spike code.

Vocabulary and deepening: [REFERENCE.md](REFERENCE.md). Parallel exploration: [DESIGN-IT-TWICE.md](DESIGN-IT-TWICE.md). Dependency categories: [DEEPENING.md](DEEPENING.md).

Adapted from [mattpocock/skills - engineering/codebase-design](https://github.com/mattpocock/skills/tree/main/skills/engineering/codebase-design). Replaces **design-an-interface**.

**v1 is human-initiated HITL only** - no AFK path.

## Not this skill

| Skill | When instead |
|-------|----------------|
| [wayfinder](../SKILL.md) | Chart, Materialize, Reconcile, Route only |
| [define-bundle](../define-bundle/SKILL.md) | Coalesce GM rows into draft/approved bundles |
| [create-tasks](../create-tasks/SKILL.md) | Split approved bundle into implementation tasks |
| [prototype](../actions/prototype/SKILL.md) | Throwaway code on bundle branch when **## Method:** `prototype` |
| [improve-codebase-architecture](../../improve-codebase-architecture/SKILL.md) | Codebase-wide exploration → GitHub RFC issues |
| [grill-me](../grill-me/SKILL.md) | Binding decisions via Q&A → GM rows |

## Entry paths

| Entry | Input | Output |
|-------|-------|--------|
| **Bundle** (primary) | Approved `wf:bundle` + **Decisions** / **Constraints** | Comment on bundle issue → hand off to create-tasks |
| **Planning** | Map **To Do** `wf:prototype` or ad-hoc module question | Comment on ticket or map issue |

## Quick start

**Bundle:** User says "design module for bundle #N" after **`bundle approved`**.

Load bundle + map → frame module from **Decisions** → [DESIGN-IT-TWICE.md](DESIGN-IT-TWICE.md) when interface shape is open → post artifact comment → suggest [create-tasks](../create-tasks/SKILL.md).

**Planning:** User invokes on a prototype or interface ticket on map **To Do**.

Load ticket → gather requirements → design-it-twice → post comment → human Reconcile when ticket complete.

See [REFERENCE.md](REFERENCE.md) for workflow, artifact template, and bundle vs planning deltas.

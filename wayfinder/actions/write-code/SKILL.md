---
name: write-code
description: Implement task deliverables with TDD at pre-agreed seams - typecheck and targeted tests during work, full suite at end. Method playbook for wf:task bundle work; dispatched by implement-task. Excludes code review, commit, and push.
---

# Write code

Build the task **What to build** on the bundle branch using **vertical-slice TDD** at pre-agreed seams, plus regular typecheck and test runs. Dispatched by [implement-task](../../implement-task/SKILL.md) after startup gates.

Detail: [REFERENCE.md](REFERENCE.md). TDD reference: [tests.md](tests.md) - [mocking.md](mocking.md).

Adapted from [mattpocock/skills - engineering/implement](https://github.com/mattpocock/skills/tree/main/skills/engineering/implement) and [engineering/tdd](https://github.com/mattpocock/skills/tree/main/skills/engineering/tdd).

## When to use

- **`wf:task`** implementation task with **## Method:** `write-code` (default for normal bundle build work)
- Task **What to build** ships repo code, tests, or docs in the target project

## Not this skill

| Skill | When instead |
|-------|----------------|
| [implement-task](../../implement-task/SKILL.md) | Gates, bundle branch, code-review, commit, push, resolution comment |
| [code-review](../../code-review/SKILL.md) | Standards + Spec review after Method (always implement-task) |
| [prototype](../prototype/SKILL.md) | Throwaway demos when **## Method:** `prototype` |
| [tdd](../../../tdd/SKILL.md) | Ad-hoc TDD sessions outside implement-task Method dispatch |
| [design-modules](../../design-modules/SKILL.md) | Modules interface shaping before tasks exist |

## Rules

1. **Seams first** - agree test seams with the human (or task spec) before writing tests; no tests at unconfirmed seams.
2. **Vertical slices** - one failing test → minimal pass → repeat; no horizontal "all tests then all code."
3. **Verify often** - typecheck regularly; run the focused test file after each slice; run the full suite once before returning to implement-task.
4. **Honor Decisions** - task **Decisions** and bundle **Constraints** are binding; scope to **Done when**.
5. **Stop at build** - do not code-review, commit, push, or post resolution; return artifacts to implement-task.

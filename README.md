# Agent Skills

A collection of agent skills that extend capabilities across planning, development, and delivery.

Install paths below use this repository (`KroniK907/skills`). If you use a fork, substitute your GitHub `owner/repo` prefix.

## Planning & design

Skills for shaping work before or alongside implementation.

- **wayfinder** — Bootstrap and maintain `FeatureName:Map` GitHub trackers: map skeleton, materialize tickets from map-discovery comment, reconcile after approval, suggest next skill. Use when a feature is too big for one session.

  ```
  npx skills@latest add KroniK907/skills/wayfinder
  ```

- **define-bundle** — Coalesce decision-log clusters into draft `wayfinder:bundle` issues; promote on `bundle approved` while planning To Do or fog stay open.

  ```
  npx skills@latest add KroniK907/skills/define-bundle
  ```

- **feature-discovery** — Breadth-first five-zone interview; posts a map-discovery artifact as a comment on the map issue for wayfinder Materialize.

  ```
  npx skills@latest add KroniK907/skills/feature-discovery
  ```

- **strategic-ideation** — Expand/tension/prune at idea level for scope and strategy; hand off to grill-me or PRD (renamed from feature-ideation).

  ```
  npx skills@latest add KroniK907/skills/strategic-ideation
  ```

- **write-a-prd** — Turn an existing long design discussion or decision artifact into a PRD, with codebase exploration and module sketching, then submit as a GitHub issue (after decisions exist; use grill-me separately if you need depth-first Q&A first).

  ```
  npx skills@latest add KroniK907/skills/write-a-prd
  ```

- **prd-to-plan** — Turn a PRD into a multi-phase implementation plan using tracer-bullet vertical slices (saved under `./plans/`).

  ```
  npx skills@latest add KroniK907/skills/prd-to-plan
  ```

- **prd-to-issues** — Break a PRD into independently-grabbable GitHub issues using vertical slices.

  ```
  npx skills@latest add KroniK907/skills/prd-to-issues
  ```

- **grill-me** — Stress-test a plan or design through sequential Q&A until open branches are resolved.

  ```
  npx skills@latest add KroniK907/skills/grill-me
  ```

- **design-an-interface** — Generate multiple radically different interface designs for a module using parallel sub-agents.

  ```
  npx skills@latest add KroniK907/skills/design-an-interface
  ```

- **request-refactor-plan** — Create a detailed refactor plan with tiny commits via user interview, then file it as a GitHub issue.

  ```
  npx skills@latest add KroniK907/skills/request-refactor-plan
  ```

## Development

Skills for building, fixing, and evolving code.

- **tdd** — Test-driven development with a red-green-refactor loop (features, fixes, integration tests).

  ```
  npx skills@latest add KroniK907/skills/tdd
  ```

- **triage-issue** — Investigate a bug by exploring the codebase, identify root cause, and file a GitHub issue with a TDD-based fix plan.

  ```
  npx skills@latest add KroniK907/skills/triage-issue
  ```

- **improve-codebase-architecture** — Explore a codebase for architectural improvement opportunities, focusing on deepening shallow modules and testability.

  ```
  npx skills@latest add KroniK907/skills/improve-codebase-architecture
  ```

- **commit** — Stage and commit only changes attributable to the current agent chat (split into logical commits when appropriate; uses `git` and `gh`).

  ```
  npx skills@latest add KroniK907/skills/commit
  ```

## Writing & knowledge

- **write-a-skill** — Create new agent skills with proper structure, progressive disclosure, and bundled resources.

  ```
  npx skills@latest add KroniK907/skills/write-a-skill
  ```

- **ubiquitous-language** — Extract a DDD-style ubiquitous language glossary from the current conversation; saves to `UBIQUITOUS_LANGUAGE.md`.

  ```
  npx skills@latest add KroniK907/skills/ubiquitous-language
  ```

## Skills in this repo

| Skill | Folder |
|-------|--------|
| wayfinder | `wayfinder/` |
| define-bundle | `define-bundle/` |
| feature-discovery | `feature-discovery/` |
| strategic-ideation | `strategic-ideation/` |
| feature-ideation | `feature-ideation/` (stub → strategic-ideation) |
| write-a-prd | `write-a-prd/` |
| prd-to-plan | `prd-to-plan/` |
| prd-to-issues | `prd-to-issues/` |
| grill-me | `grill-me/` |
| design-an-interface | `design-an-interface/` |
| request-refactor-plan | `request-refactor-plan/` |
| tdd | `tdd/` |
| triage-issue | `triage-issue/` |
| improve-codebase-architecture | `improve-codebase-architecture/` |
| commit | `commit/` |
| write-a-skill | `write-a-skill/` |
| ubiquitous-language | `ubiquitous-language/` |

Related Cursor-focused skills (hooks, rules, canvas, SDK, CLI status line, and so on) may live in a separate `skills-cursor` tree alongside this repo on your machine; they are not bundled here.

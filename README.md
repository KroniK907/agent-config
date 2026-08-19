# Agent Skills

A collection of agent skills that extend capabilities across planning, development, and delivery.

Install paths below use this repository (`KroniK907/skills`). If you use a fork, substitute your GitHub `owner/repo` prefix.

## Wayfinder ecosystem

Skills for large-feature planning and incremental implementation via GitHub map trackers. The pack lives under `wayfinder/` with subfolders `actions/`, `ideation/`, `orchestrators/`, and `utilities/`. Hub skill: `wayfinder/SKILL.md`.

- **wayfinder** - Bootstrap and maintain `FeatureName:Map` GitHub trackers: map skeleton, materialize tickets from map-discovery comment, reconcile after approval, suggest next skill. Use when a feature is too big for one session.

  ```
  npx skills@latest add KroniK907/skills/wayfinder
  ```

- **define-bundle** (action) - Group decision-log clusters into draft `wf:bundle` issues; promote on `bundle approved` while planning To Do or fog stay open.

  ```
  npx skills@latest add KroniK907/skills/wayfinder/actions/define-bundle
  ```

- **feature-discovery** (ideation) - Breadth-first five-zone interview; posts a map-discovery artifact as a comment on the map issue for wayfinder Materialize.

  ```
  npx skills@latest add KroniK907/skills/wayfinder/ideation/feature-discovery
  ```

- **strategic-ideation** (ideation) - Expand/tension/prune at idea level for scope and strategy; hand off to grill-me or PRD (renamed from feature-ideation).

  ```
  npx skills@latest add KroniK907/skills/wayfinder/ideation/strategic-ideation
  ```

- **grill-me** (ideation) - Stress-test a plan or design through sequential Q&A until open branches are resolved.

  ```
  npx skills@latest add KroniK907/skills/wayfinder/ideation/grill-me
  ```

- **design-modules** (action) - Shape one or more deep modules from bundle decisions or planning tickets; seam discovery and design-it-twice exploration; HITL only. Replaces design-an-interface.

  ```
  npx skills@latest add KroniK907/skills/wayfinder/actions/design-modules
  ```

- **write-code** (action) - Default bundle **`wf:task`** Method: TDD at pre-agreed seams via implement-task.

  ```
  npx skills@latest add KroniK907/skills/wayfinder/actions/write-code
  ```

- **create-tasks** (action) - Split an approved bundle into implementation tasks on the map **Implementing** frontier.

  ```
  npx skills@latest add KroniK907/skills/wayfinder/actions/create-tasks
  ```

- **research** (action) - Investigate `wf:research` tickets; post structured findings and non-binding tracker updates.

  ```
  npx skills@latest add KroniK907/skills/wayfinder/actions/research
  ```

## Planning utilities (map-free)

Skills for shaping work without a wayfinder map - small scope, PRDs, or standalone design.

- **write-a-prd** - Turn an existing long design discussion or decision artifact into a PRD, with codebase exploration and module sketching, then submit as a GitHub issue (after decisions exist; use grill-me separately if you need depth-first Q&A first).

  ```
  npx skills@latest add KroniK907/skills/write-a-prd
  ```

- **prd-to-plan** - Turn a PRD into a multi-phase implementation plan using tracer-bullet vertical slices (saved under `./plans/`).

  ```
  npx skills@latest add KroniK907/skills/prd-to-plan
  ```

- **prd-to-issues** - Break a PRD into independently-grabbable GitHub issues using vertical slices.

  ```
  npx skills@latest add KroniK907/skills/prd-to-issues
  ```

- **request-refactor-plan** - Create a detailed refactor plan with tiny commits via user interview, then file it as a GitHub issue.

  ```
  npx skills@latest add KroniK907/skills/request-refactor-plan
  ```

## Development

Skills for building, fixing, and evolving code.

- **triage-issue** - Investigate a bug by exploring the codebase, identify root cause, and file a GitHub issue with a TDD-based fix plan.

  ```
  npx skills@latest add KroniK907/skills/triage-issue
  ```

- **improve-codebase-architecture** - Explore a codebase for architectural improvement opportunities, focusing on deepening shallow modules and testability.

  ```
  npx skills@latest add KroniK907/skills/improve-codebase-architecture
  ```

- **commit** - Stage and commit only changes attributable to the current agent chat (split into logical commits when appropriate; uses `git` and `gh`).

  ```
  npx skills@latest add KroniK907/skills/commit
  ```

## Writing & knowledge

- **writing-for-agents** - Write documents agents consume (skills, AGENTS.md, Cursor rules): context pointers, information hierarchy, completion criteria, leading words. Adapted from [mattpocock/skills](https://github.com/mattpocock/skills/tree/main/skills/productivity/writing-for-agents).

  ```
  npx skills@latest add KroniK907/skills/writing-for-agents
  ```

- **write-a-skill** - Router to `writing-for-agents` for backward-compatible installs when creating a new skill.

  ```
  npx skills@latest add KroniK907/skills/write-a-skill
  ```

- **ubiquitous-language** - Extract a DDD-style ubiquitous language glossary from the current conversation; saves to `UBIQUITOUS_LANGUAGE.md`.

  ```
  npx skills@latest add KroniK907/skills/ubiquitous-language
  ```

- **unslop** - Cut AI tells from any writing; rewrite for plain human voice. **Always on** via description + optional rule. Adapted from [cursor/plugins/pstack](https://github.com/cursor/plugins/tree/main/pstack/skills/unslop).

  ```
  npx skills@latest add KroniK907/skills/unslop
  ```

  For global always-on behavior in every project, copy `.cursor/rules/unslop.mdc` to `~/.cursor/rules/unslop.mdc`.

## Skills in this repo

| Skill | Folder |
|-------|--------|
| wayfinder (hub) | `wayfinder/` |
| define-bundle | `wayfinder/actions/define-bundle/` |
| feature-discovery | `wayfinder/ideation/feature-discovery/` |
| strategic-ideation | `wayfinder/ideation/strategic-ideation/` |
| grill-me | `wayfinder/ideation/grill-me/` |
| design-modules | `wayfinder/actions/design-modules/` |
| write-code (action) | `wayfinder/actions/write-code/` |
| create-tasks | `wayfinder/actions/create-tasks/` |
| implement-task | `wayfinder/orchestrators/implement-task/` |
| one-off | `wayfinder/orchestrators/one-off/` |
| prototype (action) | `wayfinder/actions/prototype/` |
| research | `wayfinder/actions/research/` |
| write-a-prd | `write-a-prd/` |
| prd-to-plan | `prd-to-plan/` |
| prd-to-issues | `prd-to-issues/` |
| request-refactor-plan | `request-refactor-plan/` |
| triage-issue | `triage-issue/` |
| improve-codebase-architecture | `improve-codebase-architecture/` |
| commit | `commit/` |
| writing-for-agents | `writing-for-agents/` |
| write-a-skill | `write-a-skill/` (router → writing-for-agents) |
| ubiquitous-language | `ubiquitous-language/` |
| unslop | `unslop/` |

**Layout:** Ecosystem skills under `wayfinder/actions/`, `wayfinder/ideation/`, and `wayfinder/orchestrators/`; utilities under `wayfinder/utilities/`. One-off repo-root utilities stay at repo root.

Related Cursor-focused skills (hooks, rules, canvas, SDK, CLI status line, and so on) may live in a separate `skills-cursor` tree alongside this repo on your machine; they are not bundled here.

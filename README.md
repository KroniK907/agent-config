# Agent config hub

Team skills, Cursor rules, and scripts for NUS agent tooling. v1 layout per **AgentConfigHub** map.

**Repo:** `KroniK907/agent-config` (GitHub rename from `KroniK907/skills` pending). **Catalog:** [`catalog.json`](catalog.json) lists every installable skill, rule, and script with paths and labels.

## Layout

| Path | Purpose |
|------|---------|
| `skills/` | Agent skills - flat folders + `skills/wayfinder/` tree |
| `rules/` | Team Cursor rules pack (`*.mdc`) |
| `scripts/` | Apply, bootstrap, and validation scripts |
| `AGENTS.md` | What agents should know about this repo |
| `.cursor/` | Example project templates only |

## Install

```text
npx skills@latest add KroniK907/agent-config/skills/<skill-path>
```

Substitute your fork if needed. Paths are in [`catalog.json`](catalog.json). For global always-on **unslop**, copy `rules/unslop.mdc` to `~/.cursor/rules/unslop.mdc`.

## Wayfinder ecosystem

Skills for large-feature planning and incremental implementation via GitHub map trackers. Hub skill: `skills/wayfinder/SKILL.md`.

- **wayfinder** - Bootstrap and maintain `FeatureName:Map` GitHub trackers: map skeleton, materialize tickets from map-discovery comment, reconcile after approval, suggest next skill. Use when a feature is too big for one session.
- **define-bundle** (action) - Group decision-log clusters into draft `wf:bundle` issues; promote on `bundle approved` while planning To Do or fog stay open.
- **feature-discovery** (ideation) - Breadth-first five-zone interview; posts a map-discovery artifact as a comment on the map issue for wayfinder Materialize.
- **strategic-ideation** (ideation) - Expand/tension/prune at idea level for scope and strategy; hand off to grill-me or PRD.
- **grill-me** (ideation) - Stress-test a plan or design through sequential Q&A until open branches are resolved.
- **design-modules** (action) - Shape one or more deep modules from bundle decisions; seam discovery and design-it-twice exploration; HITL only.
- **write-code** (action) - Default bundle **`wf:task`** Method: TDD at pre-agreed seams via implement-task.
- **create-tasks** (action) - Split an approved bundle into implementation tasks on the map **Implementing** frontier.
- **research** (action) - Investigate `wf:research` tickets; post structured findings and non-binding tracker updates.
- **implement-task** (orchestrator) - Pick up **`wf:approved`** tasks on a bundle branch; code review, push, resolution comment.
- **one-off** (orchestrator) - Map **To Do** tickets that ship repo work without the define-bundle pipeline.
- **prototype** (action) - Throwaway demos when **## Method:** `prototype`.
- **code-review** (action) - Standards + spec review; auto-fix obvious issues in implement-task mode.
- **constrain-fog** (ideation) - Groom **Not yet specified** fog on a map into reconcile-ready tickets.

## Planning utilities (map-free)

Skills for shaping work without a wayfinder map - small scope, PRDs, or standalone design.

- **write-a-prd** - Turn an existing long design discussion or decision artifact into a PRD, with codebase exploration and module sketching, then submit as a GitHub issue.
- **prd-to-plan** - Turn a PRD into a multi-phase implementation plan using tracer-bullet vertical slices (saved under `./plans/`).
- **prd-to-issues** - Break a PRD into independently-grabbable GitHub issues using vertical slices.
- **request-refactor-plan** - Create a detailed refactor plan with tiny commits via user interview, then file it as a GitHub issue.

## Development

- **triage-issue** - Investigate a bug by exploring the codebase, identify root cause, and file a GitHub issue with a TDD-based fix plan.
- **improve-codebase-architecture** - Explore a codebase for architectural improvement opportunities, focusing on deepening shallow modules and testability.
- **commit** - Stage and commit only changes attributable to the current agent chat (split into logical commits when appropriate).

## Writing and knowledge

- **writing-for-agents** - Write documents agents consume (skills, AGENTS.md, Cursor rules): context pointers, information hierarchy, completion criteria, leading words.
- **write-a-skill** - Router to `writing-for-agents` for backward-compatible installs when creating a new skill.
- **ubiquitous-language** - Extract a DDD-style ubiquitous language glossary from the current conversation; saves to `UBIQUITOUS_LANGUAGE.md`.
- **unslop** - Cut AI tells from any writing; rewrite for plain human voice. **Always on** via description + optional rule.
- **ccr-summary** - Generate the opening summary for a CCR from a full report PDF per 40 CFR § 141.156.

Cursor product skills (hooks, canvas, SDK) live in a separate `skills-cursor` tree - not bundled here.

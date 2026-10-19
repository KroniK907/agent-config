# Agent config hub

Team skills, project rules, and scripts for NUS agent tooling. v1 layout per **AgentConfigHub** map. Claude Code and Codex load this pack, usually through T3 Code.

**Repo:** [`KroniK907/agent-config`](https://github.com/KroniK907/agent-config). **Catalog:** [`catalog.json`](catalog.json) lists every installable skill, rule, and script with paths and labels.

## Layout

| Path | Purpose |
|------|---------|
| `skills/` | Agent skills - flat folders + `skills/wayfinder/` tree |
| `rules/` | Team rules pack (`*.md`) |
| `scripts/` | `validate-catalog`, run in CI |
| `skills/wayfinder/utilities/wf/` | `wf` helper for decision-log comments and surgical map edits (`go run wf.go`) |
| `AGENTS.md` | What agents should know about this repo |

## Install on a machine (primary)

Clone this repo to a stable path. Symlink each skill into `~/.claude/skills/` and each rule into `~/.claude/rules/`. [AGENTS.md](AGENTS.md#working-in-this-repo) has the commands and the rule-loading rules Claude applies.

## One-off skill install (secondary)

```text
npx skills@latest add KroniK907/agent-config/skills/<skill-path>
```

Use for a single skill without cloning the repo. Substitute your fork if needed. Paths are in [`catalog.json`](catalog.json).

## Wayfinder ecosystem

Skills for large-feature planning and incremental implementation via GitHub map trackers. Hub skill: `skills/wayfinder/SKILL.md`.

- **wayfinder** - Bootstrap and maintain `FeatureName:Map` GitHub trackers: map skeleton, materialize tickets from map-discovery comment, reconcile after approval, suggest next skill. Use when a feature is too big for one session.
- **define-bundle** (action) - Group decision-log clusters into draft `wf:bundle` issues; approve them while planning To Do or fog stay open.
- **feature-discovery** (ideation) - Breadth-first five-zone interview; posts a map-discovery artifact as a comment on the map issue for wayfinder Materialize.
- **strategic-ideation** (ideation) - Expand/tension/prune at idea level for scope and strategy; hand off to grill-me or PRD.
- **grill-me** (ideation) - Stress-test a plan or design through sequential Q&A until open branches are resolved.
- **design-modules** (action) - Shape one or more deep modules from bundle decisions; seam discovery and design-it-twice exploration; HITL only.
- **write-code** (action) - Default bundle **`wf:task`** Method: TDD at pre-agreed seams via implement-task.
- **create-tasks** (action) - Split an approved bundle into implementation tasks on the map **Implementing** frontier.
- **research** (action) - Investigate `wf:research` tickets; post structured findings and non-binding tracker updates.
- **implement-task** (orchestrator) - Pick up **`wf:approved`** tasks in a worktree; code review, push, pull request, resolution comment.
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
- **backlog-sweep** - Sweep open GitHub issues against merged work: close shipped and stale issues, narrow half-done ones, and label small quick work `priority: high`. User-invoked; schedule it for a weekly run.
- **commit** - Stage and commit only changes attributable to the current agent chat, then open a pull request into `integrationBranch` unless the user says not to.

## Writing and knowledge

- **writing-for-agents** - Write documents agents consume (skills, AGENTS.md, CLAUDE.md, project rules): context pointers, information hierarchy, completion criteria, leading words.
- **write-a-skill** - Router to `writing-for-agents` for backward-compatible installs when creating a new skill.
- **ubiquitous-language** - Extract a DDD-style ubiquitous language glossary from the current conversation; saves to `UBIQUITOUS_LANGUAGE.md`.
- **unslop** - Cut AI tells from any writing; rewrite for plain human voice. **Always on** via description + optional rule.
- **ccr-summary** - Generate the opening summary for a CCR from a full report PDF per 40 CFR § 141.156.

Host-bundled product skills (editor hooks, canvas, SDK) stay outside this repo.

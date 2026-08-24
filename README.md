# Agent config hub

Team skills, Cursor rules, and scripts for NUS agent tooling. v1 layout per **AgentConfigHub** map.

**Repo:** [`KroniK907/agent-config`](https://github.com/KroniK907/agent-config). **Catalog:** [`catalog.json`](catalog.json) lists every installable skill, rule, and script with paths and labels.

## Layout

| Path | Purpose |
|------|---------|
| `skills/` | Agent skills - flat folders + `skills/wayfinder/` tree |
| `rules/` | Team Cursor rules pack (`*.mdc`) |
| `scripts/` | Bootstrap, validation, and [agent-config-wizard](scripts/wizard/README.md) apply TUI |
| `AGENTS.md` | What agents should know about this repo |
| `.cursor/` | Example project templates only |

## Apply to a project (primary)

1. Clone this repo to a stable path on your machine (team root).
2. From your **project** repo root, run the wizard binary or `go run`.

```powershell
# Windows - built binary from project cwd
C:\path\to\agent-config\scripts\wizard\bin\agent-config-wizard.exe
```

```bash
# Linux
/path/to/agent-config/scripts/wizard/bin/agent-config-wizard
```

Dev without a binary:

```powershell
go run -C C:\path\to\agent-config\scripts\wizard\src . -project C:\path\to\your-project
```

The team clone does **not** need to live inside the project repo. Full run/build/test docs: [scripts/wizard/README.md](scripts/wizard/README.md).

**First run:** opt in to skills and rules in the TUI, press **a** to apply. Creates `.cursor/agent-config.local.json` (gitignore) and copies selected items into `.cursor/skills/` and `.cursor/rules/`. See [.cursor/README.md](.cursor/README.md) for manifest, gitignore, and re-run behavior.

**Cloud agents:** commit `.cursor/agent-manifest.json` (subset of desktop opt-in). Template: [.cursor/examples/agent-manifest.json.example](.cursor/examples/agent-manifest.json.example).

## One-off skill install (secondary)

```text
npx skills@latest add KroniK907/agent-config/skills/<skill-path>
```

Use for a single skill without full project apply. Substitute your fork if needed. Paths are in [`catalog.json`](catalog.json).

## Legacy global install

[`skills/wayfinder/utilities/bootstrap/install-skills.ps1`](skills/wayfinder/utilities/bootstrap/install-skills.ps1) copies the full wayfinder pack into `~/.cursor/skills/` for **Cloud AFK bootstrap** and pre-apply smoke tests. Prefer **agent-config-wizard** for day-to-day desktop project setup. See [AFK-BOOTSTRAP.md](skills/wayfinder/utilities/AFK-BOOTSTRAP.md) for cloud agent setup.

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

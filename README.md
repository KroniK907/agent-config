# Agent config hub

Team skills, Cursor rules, and scripts for NUS agent tooling. v1 layout per **AgentConfigHub** map.

Install paths below use **`KroniK907/agent-config`** (GitHub rename from `KroniK907/skills` after the v1 layout PR merges). If you use a fork, substitute your `owner/repo` prefix. Catalog authority is root [`catalog.json`](catalog.json).

## Layout

| Path | Purpose |
|------|---------|
| `skills/` | Agent skills (flat + `skills/wayfinder/` tree) |
| `rules/` | Team Cursor rules pack (`*.mdc`) |
| `scripts/` | Apply, bootstrap, and validation scripts |
| `AGENTS.md` | Team baseline for agents |
| `.cursor/` | Example project templates only |

## Wayfinder ecosystem

Skills for large-feature planning and incremental implementation via GitHub map trackers. Hub skill: `skills/wayfinder/SKILL.md`.

- **wayfinder** - Bootstrap and maintain `FeatureName:Map` GitHub trackers.

  ```
  npx skills@latest add KroniK907/agent-config/skills/wayfinder
  ```

- **define-bundle** (action) - Group decision-log clusters into draft `wf:bundle` issues.

  ```
  npx skills@latest add KroniK907/agent-config/skills/wayfinder/actions/define-bundle
  ```

- **feature-discovery** (ideation) - Breadth-first five-zone interview.

  ```
  npx skills@latest add KroniK907/agent-config/skills/wayfinder/ideation/feature-discovery
  ```

- **strategic-ideation** (ideation) - Expand/tension/prune at idea level.

  ```
  npx skills@latest add KroniK907/agent-config/skills/wayfinder/ideation/strategic-ideation
  ```

- **grill-me** (ideation) - Stress-test a plan through sequential Q&A.

  ```
  npx skills@latest add KroniK907/agent-config/skills/wayfinder/ideation/grill-me
  ```

- **design-modules** (action) - Shape deep modules from bundle decisions.

  ```
  npx skills@latest add KroniK907/agent-config/skills/wayfinder/actions/design-modules
  ```

- **write-code** (action) - Default bundle **`wf:task`** Method.

  ```
  npx skills@latest add KroniK907/agent-config/skills/wayfinder/actions/write-code
  ```

- **create-tasks** (action) - Split an approved bundle into implementation tasks.

  ```
  npx skills@latest add KroniK907/agent-config/skills/wayfinder/actions/create-tasks
  ```

- **research** (action) - Investigate `wf:research` tickets.

  ```
  npx skills@latest add KroniK907/agent-config/skills/wayfinder/actions/research
  ```

## Planning utilities (map-free)

- **write-a-prd** - Turn a design discussion into a PRD GitHub issue.

  ```
  npx skills@latest add KroniK907/agent-config/skills/write-a-prd
  ```

- **prd-to-plan** - Turn a PRD into a multi-phase plan under `./plans/`.

  ```
  npx skills@latest add KroniK907/agent-config/skills/prd-to-plan
  ```

- **prd-to-issues** - Break a PRD into vertical-slice GitHub issues.

  ```
  npx skills@latest add KroniK907/agent-config/skills/prd-to-issues
  ```

- **request-refactor-plan** - Create a refactor plan and file as a GitHub issue.

  ```
  npx skills@latest add KroniK907/agent-config/skills/request-refactor-plan
  ```

## Development

- **triage-issue** - Investigate a bug and file a GitHub issue with a fix plan.

  ```
  npx skills@latest add KroniK907/agent-config/skills/triage-issue
  ```

- **improve-codebase-architecture** - Find architectural improvement opportunities.

  ```
  npx skills@latest add KroniK907/agent-config/skills/improve-codebase-architecture
  ```

- **commit** - Commit only changes from the current agent chat.

  ```
  npx skills@latest add KroniK907/agent-config/skills/commit
  ```

## Writing and knowledge

- **writing-for-agents** - Write documents agents consume (skills, AGENTS.md, rules).

  ```
  npx skills@latest add KroniK907/agent-config/skills/writing-for-agents
  ```

- **write-a-skill** - Router to `writing-for-agents`.

  ```
  npx skills@latest add KroniK907/agent-config/skills/write-a-skill
  ```

- **ubiquitous-language** - Build a DDD glossary from conversation.

  ```
  npx skills@latest add KroniK907/agent-config/skills/ubiquitous-language
  ```

- **unslop** - Cut AI tells from any writing. **Always on** via description + rule.

  ```
  npx skills@latest add KroniK907/agent-config/skills/unslop
  ```

  For global always-on behavior, copy `rules/unslop.mdc` to `~/.cursor/rules/unslop.mdc`.

- **ccr-summary** - Summarize CCR contact records.

  ```
  npx skills@latest add KroniK907/agent-config/skills/ccr-summary
  ```

## Skills in this repo

| Skill | Folder |
|-------|--------|
| wayfinder (hub) | `skills/wayfinder/` |
| define-bundle | `skills/wayfinder/actions/define-bundle/` |
| feature-discovery | `skills/wayfinder/ideation/feature-discovery/` |
| strategic-ideation | `skills/wayfinder/ideation/strategic-ideation/` |
| grill-me | `skills/wayfinder/ideation/grill-me/` |
| design-modules | `skills/wayfinder/actions/design-modules/` |
| write-code (action) | `skills/wayfinder/actions/write-code/` |
| create-tasks | `skills/wayfinder/actions/create-tasks/` |
| implement-task | `skills/wayfinder/orchestrators/implement-task/` |
| one-off | `skills/wayfinder/orchestrators/one-off/` |
| prototype (action) | `skills/wayfinder/actions/prototype/` |
| research | `skills/wayfinder/actions/research/` |
| write-a-prd | `skills/write-a-prd/` |
| prd-to-plan | `skills/prd-to-plan/` |
| prd-to-issues | `skills/prd-to-issues/` |
| request-refactor-plan | `skills/request-refactor-plan/` |
| triage-issue | `skills/triage-issue/` |
| improve-codebase-architecture | `skills/improve-codebase-architecture/` |
| commit | `skills/commit/` |
| ccr-summary | `skills/ccr-summary/` |
| writing-for-agents | `skills/writing-for-agents/` |
| write-a-skill | `skills/write-a-skill/` |
| ubiquitous-language | `skills/ubiquitous-language/` |
| unslop | `skills/unslop/` |

Related Cursor-focused skills (hooks, canvas, SDK, and so on) may live in a separate `skills-cursor` tree; they are not bundled here.

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

Path pattern (substitute your fork if needed):

```text
npx skills@latest add KroniK907/agent-config/skills/<skill-path>
```

Examples: `skills/wayfinder`, `skills/wayfinder/actions/research`, `skills/unslop`. See [`catalog.json`](catalog.json) for every `path` entry.

For global always-on **unslop**, copy `rules/unslop.mdc` to `~/.cursor/rules/unslop.mdc`.

## Skills

| Skill | Path | Role |
|-------|------|------|
| wayfinder (hub) | `skills/wayfinder/` | Map trackers - Chart, Materialize, Reconcile, Route |
| define-bundle | `skills/wayfinder/actions/define-bundle/` | Group GMs into bundles |
| create-tasks | `skills/wayfinder/actions/create-tasks/` | Split bundles into tasks |
| implement-task | `skills/wayfinder/orchestrators/implement-task/` | Pick up `wf:approved` tasks |
| one-off | `skills/wayfinder/orchestrators/one-off/` | Map To Do without bundle pipeline |
| write-code | `skills/wayfinder/actions/write-code/` | Default bundle task Method |
| code-review | `skills/wayfinder/actions/code-review/` | Review + auto-fix pass |
| design-modules | `skills/wayfinder/actions/design-modules/` | Module interface shaping |
| prototype | `skills/wayfinder/actions/prototype/` | Throwaway demos |
| research | `skills/wayfinder/actions/research/` | `wf:research` tickets |
| feature-discovery | `skills/wayfinder/ideation/feature-discovery/` | Breadth-first map discovery |
| strategic-ideation | `skills/wayfinder/ideation/strategic-ideation/` | Idea-level scope work |
| grill-me | `skills/wayfinder/ideation/grill-me/` | Depth-first Q&A |
| constrain-fog | `skills/wayfinder/ideation/constrain-fog/` | Resolve map fog |
| write-a-prd | `skills/write-a-prd/` | PRD from decisions |
| prd-to-plan | `skills/prd-to-plan/` | PRD to phased plan |
| prd-to-issues | `skills/prd-to-issues/` | PRD to GitHub issues |
| request-refactor-plan | `skills/request-refactor-plan/` | Refactor RFC |
| triage-issue | `skills/triage-issue/` | Bug investigation + issue |
| improve-codebase-architecture | `skills/improve-codebase-architecture/` | Architecture review |
| commit | `skills/commit/` | Commit current chat only |
| writing-for-agents | `skills/writing-for-agents/` | Docs agents consume |
| write-a-skill | `skills/write-a-skill/` | Router to writing-for-agents |
| ubiquitous-language | `skills/ubiquitous-language/` | DDD glossary |
| unslop | `skills/unslop/` | Cut AI tells from writing |
| ccr-summary | `skills/ccr-summary/` | CCR PDF summary |

Cursor product skills (hooks, canvas, SDK) live in a separate `skills-cursor` tree - not bundled here.

# Agent config hub

Team skills, project rules, and scripts for NUS agent tooling. v1 layout per **AgentConfigHub** map. Claude Code and Codex load this pack, usually through T3 Code.

**Repo:** [`KroniK907/agent-config`](https://github.com/KroniK907/agent-config). **Catalog:** [`catalog.json`](catalog.json) lists every installable skill, rule, and script with paths and labels.

## Layout

| Path | Purpose |
|------|---------|
| `skills/` | Agent skills - flat folders + `skills/wayfinder/` tree |
| `rules/` | Team rules pack (`*.md`) |
| `scripts/` | `validate-catalog`, run in CI |
| `skills/wayfinder/utilities/wf/` | `wf` helper that reads map sections and uploads map bodies (`go run wf.go`) |
| `AGENTS.md` | What agents should know about this repo |

## Install on a machine (primary)

Clone this repo to a stable path. Symlink each skill into `~/.claude/skills/` and each rule into `~/.claude/rules/`. [AGENTS.md](AGENTS.md#working-in-this-repo) has the commands and the rule-loading rules Claude applies.

## One-off skill install (secondary)

```text
npx skills@latest add KroniK907/agent-config/skills/<skill-path>
```

Use for a single skill without cloning the repo. Substitute your fork if needed. Paths are in [`catalog.json`](catalog.json).

## Wayfinder

Large features live in one GitHub map issue: the plan, its decisions and its build order. Hub: `skills/wayfinder/SKILL.md`; map format: `skills/wayfinder/MAP.md`.

- **map** - Research a feature and file its `{FeatureName}:Map` issue, or post a research or cross-map note on one.
- **grill** - Settle a map in waves: a Proposed/Ask bag, then one question at a time. Writes GM decision comments and the build order (**Slices**).
- **build** - Orchestrate a map's slices as a stack of PRs: planner contract, one builder and one babysitter per layer, a cross-model stack review, merge, and an HTML report.
- **release** - Ship `dev` to `main` with a version bump and user-facing changelog, then back-merge.
- **research** - Answer a map question from primary sources and post the findings on the map.
- **design-modules** - Shape deep-module interfaces and seams before a build.

## Development

- **code-review** - Two-axis review (standards and spec) of a PR, branch or diff; can post findings on the PR.
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

---
name: subagent
description: "Pick the model and reasoning effort for a subagent before you spawn one. Use when delegating work to a subagent, fanning out parallel workers, or a skill tells you to hand a task off."
agent-config-sync: true
---

Every subagent gets the cheapest **tier** that will finish its task on the first try. Cost here means API list price per completed task, not per token. A cheap model that needs a retry, or a slow one that stalls the main loop, costs more than the next tier up.

## Steps

### 1. Decide whether to delegate

A subagent starts cold. It re-reads what you already know, and its brief has to carry everything it needs. Do the task yourself when you already hold the context and it takes about three tool calls or fewer. Delegate when the task is self-contained, when it would flood your context with file or log dumps, or when several independent tasks can run in parallel.

**Done when:** you can say in one sentence why this task goes to a subagent.

### 2. Pick the tier

Match the task to the first row whose markers fit. When a task straddles two rows, take the higher one.

| Tier | Markers | Examples |
|------|---------|----------|
| **scout** | Read-only. The answer is a fact you can check. | Find files or symbols, list call sites, pull the failing test out of a CI log, summarize a file or thread, look up a documented value |
| **edit** | Writes code, but the brief names the files and the change. A pattern to copy exists. No design choice is left. | Apply a review nit, rename, add a test that mirrors its neighbour, fix a lint or type error with an obvious fix, add a catalog entry |
| **build** | Ordinary engineering judgement. Several files, or a bug with a reproduction, or a review finding you must trace through call paths. | Multi-file feature slice, debug a failing test, review a diff, write tests for new behaviour |
| **hard** | Ambiguous spec, cross-cutting design, security, concurrency, data migrations, or a task that already failed at a lower tier. | Architecture change, auth or permissions logic, race condition, schema migration |

**Done when:** the task has one tier and you can name the marker that put it there.

### 3. Pick the model and effort

Use the family your host runs: Claude models under Claude Code, GPT models under Codex.

| Tier | Claude | GPT |
|------|--------|-----|
| scout | Haiku 5.5, `low` (`medium` for wide searches) | GPT-6 Luna, `medium` (`low` for single lookups) |
| edit | Haiku 5.5, `high` | GPT-6 Luna, `high` |
| build | Opus 5.5, `medium` | GPT-6.1 Sol, `medium` |
| hard | Opus 5.5, `high` | GPT-6.1 Sol, `high` |

Effort rules:

- `max` is off the table. Over `xhigh` it buys zero to four index points for 1.3 to 2.7 times the cost, and two to twelve minutes per task.
- `xhigh` is for a **hard** task that already failed at `high`, when the failure was shallow reasoning rather than missing context. Fix missing context by improving the brief instead.
- When the main loop is waiting on the result and seconds matter, swap a Haiku **edit** for Sonnet 5.5 at `low`. Haiku at `high` writes long reasoning and averages about 30 seconds per task. Sonnet at `low` starts answering in about one second.

Leave Claude Fable 5.1 and GPT-6 Astra unpicked unless the user names them. Opus 5.5 at `high` scores above Fable 5.1 at `xhigh` for under a third of the cost, and GPT-6.1 Sol at `xhigh` matches Astra at `high` for under a quarter of the cost. Skip GPT-5.6 Terra too: Sol beats it at every price point.

**Done when:** you have a model and an effort level, and neither is `max`.

### 4. Write the brief

The subagent sees only the brief. Give it:

- **Goal:** one sentence on what to produce.
- **Where:** the files, directories, PR, or URLs to start from.
- **Constraints:** what to leave alone, conventions to follow, and whether it may edit files.
- **Done when:** a checkable condition, for example "the test passes" or "every call site listed with `file:line`".
- **Return:** the shape of the answer. Ask for conclusions with `file:line` references, not file contents.

Keep **scout** tasks narrow enough to stay under about 100K tokens of context. Haiku 5.5 bills five times its base rate past 100K prompt tokens. GPT-6 Luna's surcharge starts at 272K.

**Done when:** a reader with no other context could do the task from the brief alone.

### 5. Spawn and check

- **Claude Code:** the Agent tool takes `model` (`haiku`, `sonnet`, `opus`) and `effort` (`low`, `medium`, `high`, `xhigh`). This skill is the instruction that permits setting `effort`.
- **Codex:** name the model and reasoning effort when you spawn. Explicit spawn values override `[agents] default_subagent_model` and `default_subagent_reasoning_effort` in `config.toml`. A custom agent file in `.codex/agents/` pins its own `model` and `model_reasoning_effort` instead.
- **Other hosts:** pick the agent definition whose model matches the row.

Launch independent subagents in one batch so they run in parallel. Check each result against its done-when before you use it, and read the diff of any subagent that edited files. If a result fails the check, rerun it one tier up with a sharper brief. Retrying at the same tier rarely helps.

**Done when:** every result passed its check or was rerun one tier up.

## Data behind the tables

The tiers rest on API list prices and Artificial Analysis cost-per-task figures from early October 2026. [benchmarks.md](benchmarks.md) has the numbers, sources, and the date they were checked. Refresh it when a new model ships in either family.

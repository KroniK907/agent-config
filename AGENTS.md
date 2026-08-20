# Team agent baseline

Northern Utility Services internal agent defaults. Project repos opt in via agent-config apply (see map **AgentConfigHub**).

## Environment

- **OS:** Windows 10/11
- **Shell:** PowerShell for terminal commands
- **Go:** Installed (`go` on PATH) - prefer Go for small scripts and utilities
- **Node.js / Python:** Not assumed on staff machines unless a project documents them

## Coding

- Keep changes small and obvious. Match existing conventions in the repo you are in.
- Comments explain non-obvious business logic, not every line.
- Tests should cover real behavior, not ceremony.
- Do not touch production, live databases, or deploy channels unless explicitly asked.

## Writing

- Apply **unslop** to all user-facing text.
- Use ASCII hyphen `-`, not em dash or en dash.
- Plain language; short sentences.

## Tracker vs portable docs

- GitHub issues hold map, bundle, and task trackers with concrete links and GM IDs.
- Committed skill markdown in this repo uses placeholders (`#N`, `{MAP-SLUG}-GM-001`) per **portable-skill-docs** rule.

## Install paths

Skills and rules ship from **`KroniK907/agent-config`** (rename pending after v1 layout PR merges). Catalog authority is root `catalog.json`.

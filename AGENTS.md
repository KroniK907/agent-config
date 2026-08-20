# Agent config hub

Northern Utility Services team source for agent skills, Cursor rules, and bootstrap scripts. Other repos opt in via agent-config apply; this repo is the canonical copy those tools pull from.

## What lives here

| Path | Role |
|------|------|
| `skills/` | Installable agent skills - flat folders plus `skills/wayfinder/` sub-tree |
| `rules/` | Team `.mdc` rules copied into project `.cursor/rules/` |
| `scripts/` | Validation, apply, and bootstrap tooling |
| `catalog.json` | Sole catalog - every skill, rule, and script with `path` and `label` |
| `.cursor/` | Example project templates only - not the live rules pack |

## Wayfinder tree

`skills/wayfinder/` is the hierarchical opt-in example:

- `SKILL.md` - hub (Chart, Materialize, Reconcile, Route)
- `actions/` - build playbooks (`write-code`, `create-tasks`, `research`, …)
- `ideation/` - planning interviews (`grill-me`, `feature-discovery`, …)
- `orchestrators/` - `implement-task`, `one-off`
- `utilities/` - bootstrap scripts, map validators - not installable skills

Repo-root skills under `skills/<name>/` are map-free utilities (PRD tools, `commit`, `unslop`, etc.).

## Working in this repo

**Catalog is authoritative.** If you add or move a skill folder, update `catalog.json`. CI runs `go run scripts/validate-catalog/main.go` on PR - drift fails the build.

**Portable markdown.** Committed `.md` files use placeholders (`#N`, `{MAP-SLUG}-GM-001`, `{FeatureName}:Map`) - not live issue URLs or map-specific GM rows. Concrete tracker links belong in GitHub issue bodies and comments. See `rules/portable-skill-docs.mdc`.

**Sync frontmatter.** Team skill sources include `agent-config-sync: true` in YAML frontmatter so apply tooling knows they are managed copies.

**Install paths.** Docs and examples reference `KroniK907/agent-config/skills/...` (post-rename). Until GitHub rename lands, the remote may still be `KroniK907/skills`.

**Do not treat `.cursor/` as the rules source.** Rules ship from `rules/`. Project manifests and gitignore patterns are documented under `.cursor/examples/`.

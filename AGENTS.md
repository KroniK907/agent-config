# Agent config hub

Northern Utility Services team source for agent skills, project rules, and bootstrap scripts. Other repos opt in via agent-config apply; this repo is the canonical copy those tools pull from. Claude CLI and Cursor CLI both load these skills and rules.

## What lives here

| Path | Role |
|------|------|
| `skills/` | Installable agent skills - flat folders plus `skills/wayfinder/` sub-tree |
| `rules/` | Team rules (`*.mdc`) applied into the project for whichever CLI loads them |
| `scripts/` | Validation, bootstrap, and **agent-config-wizard** ([scripts/wizard/](scripts/wizard/README.md)) |
| `catalog.json` | Sole catalog - every skill, rule, and script with `path` and `label` |
| `.cursor/` | Example project templates only - not the live rules pack |

## Wayfinder tree

`skills/wayfinder/` is the hierarchical opt-in example:

- `SKILL.md` - hub (Chart, Materialize, Reconcile, Route)
- `actions/` - build playbooks (`write-code`, `create-tasks`, `research`, …)
- `ideation/` - planning interviews (`grill-me`, `feature-discovery`, …)
- `orchestrators/` - `implement-task`, `one-off`
- `utilities/` - bootstrap scripts and the `wf` helper (`go run utilities/wf/wf.go`) - not installable skills

Repo-root skills under `skills/<name>/` are map-free utilities (PRD tools, `commit`, `unslop`, etc.).

## Working in this repo

**Catalog is authoritative.** If you add or move a skill folder, update `catalog.json`. CI runs `go run scripts/validate-catalog/main.go` on PR - drift fails the build.

**Local CLI symlinks.** On a machine where this checkout is the personal skill source, `~/.cursor/skills` is a directory symlink to `skills/` and `~/.cursor/rules` is a directory symlink to `rules/`. New skills and rules show up for the Cursor CLI with no extra link. The Claude CLI cannot use one directory symlink: `~/.claude/skills/` also holds other content, so each top-level skill is its own symlink. After you add `skills/<name>/`, link that folder into Claude's skills directory, following the existing links (stable clone, not a temporary worktree):

```bash
ln -sfn "$(readlink -f ~/.cursor/skills)/<name>" ~/.claude/skills/<name>
```

Claude loads a skill only when `SKILL.md` sits directly in `~/.claude/skills/<name>/`. A link to `skills/wayfinder` exposes the hub and hides every skill under `actions/`, `ideation/`, and `orchestrators/`. Link each of those folders by its skill name (same stable clone, skip the hub because `wayfinder` is already linked):

```bash
src="$(readlink -f ~/.cursor/skills)"
find "$src/wayfinder" -name SKILL.md -printf '%h\n' | while read -r dir; do
  name="$(basename "$dir")"
  [ "$name" = wayfinder ] && continue
  ln -sfn "$dir" "$HOME/.claude/skills/$name"
done
```

**Portable markdown.** Committed `.md` files use placeholders (`#N`, `{MAP-SLUG}-GM-001`, `{FeatureName}:Map`) - not live issue URLs or map-specific GM rows. Concrete tracker links belong in GitHub issue bodies and comments. See `rules/portable-skill-docs.mdc`.

**Sync frontmatter.** Team skill sources include `agent-config-sync: true` in YAML frontmatter so apply tooling knows they are managed copies.

**Install paths.** Docs and examples use `KroniK907/agent-config/skills/...`.

**Do not treat `.cursor/` as the rules source.** Rules ship from `rules/`. Project manifests and gitignore patterns are documented under `.cursor/examples/`.

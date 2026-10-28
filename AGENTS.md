# Agent config hub

Northern Utility Services team source for agent skills and rules. Claude Code and Codex load them from a local checkout of this repo, usually through T3 Code.

## What lives here

| Path | Role |
|------|------|
| `skills/` | Installable agent skills - flat folders plus `skills/wayfinder/` sub-tree |
| `rules/` | Team rules (`*.md`), linked into `~/.claude/rules/` |
| `scripts/` | `validate-catalog` (CI) |
| `catalog.json` | Sole catalog - every skill, rule, and script with `path` and `label` |

## Wayfinder tree

`skills/wayfinder/` holds the map lifecycle:

- `SKILL.md` - hub overview; `MAP.md` - the map issue format every wayfinder skill shares
- `map/`, `grill/`, `build/`, `release/` - one skill per lifecycle step (`sol-build/` is the sol-loop variant of `build/`); `research/` and `design-modules/` are helpers
- `utilities/` - the `wf` helper (`go run utilities/wf/wf.go`) and the label bootstrap - not installable skills

Wayfinder skills reach `MAP.md` as `../MAP.md`, so resolve a skill's folder with `readlink -f` first. Repo-root skills under `skills/<name>/` are map-free utilities (`code-review`, `commit`, `unslop`, etc.).

## Working in this repo

**Catalog is authoritative.** If you add or move a skill folder, update `catalog.json`. CI runs `go run scripts/validate-catalog/main.go` on PR - drift fails the build.

**Local CLI symlinks.** On a machine where this checkout is the personal skill source, each skill and rule is symlinked from the stable clone (not a temporary worktree) into `~/.claude/`. `~/.claude/skills/` also holds other content, so link each top-level skill on its own:

```bash
ln -sfn /path/to/agent-config/skills/<name> ~/.claude/skills/<name>
```

Claude loads a skill only when `SKILL.md` sits directly in `~/.claude/skills/<name>/`. A link to `skills/wayfinder` exposes only the hub. Link each skill folder under it by its skill name, and skip the hub because `wayfinder` is already linked:

```bash
src=/path/to/agent-config/skills
find "$src/wayfinder" -name SKILL.md -printf '%h\n' | while read -r dir; do
  name="$(basename "$dir")"
  [ "$name" = wayfinder ] && continue
  ln -sfn "$dir" "$HOME/.claude/skills/$name"
done
```

Link rules the same way: `ln -sfn /path/to/agent-config/rules/<name>.md ~/.claude/rules/<name>.md`. Claude judges a linked rule by the target's extension and skips targets that do not end in `.md`, with no warning. Scope a rule with a `paths:` list in its frontmatter. Claude ignores Cursor's `globs` and `alwaysApply`, so a rule without `paths:` loads in every session.

**Portable markdown.** Committed `.md` files use placeholders (`#N`, `{MAP-SLUG}-GM-001`, `{FeatureName}:Map`) - not live issue URLs or map-specific GM rows. Concrete tracker links belong in GitHub issue bodies and comments. See `rules/portable-skill-docs.md`.

**Sync frontmatter.** Team skill sources include `agent-config-sync: true` in YAML frontmatter to mark them as managed team copies.

**Install paths.** Docs and examples use `KroniK907/agent-config/skills/...`.


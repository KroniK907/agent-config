# Example project `.cursor/` files

Templates for project repos using agent-config. Copy or generate via **agent-config-wizard** ([scripts/wizard](../scripts/wizard/README.md)).

| File | Purpose |
|------|---------|
| [examples/agent-config.local.json.example](examples/agent-config.local.json.example) | Desktop state (gitignore in real projects) |
| [examples/agent-manifest.json.example](examples/agent-manifest.json.example) | Committed cloud bootstrap subset |
| [examples/gitignore.snippet](examples/gitignore.snippet) | Paste into project `.gitignore` |
| [examples/environment.json.example](examples/environment.json.example) | Cloud agent Build `install` hook |

Team rules and skills live in repo-root `rules/` and `skills/`, not here.

## First run

1. Paste [examples/gitignore.snippet](examples/gitignore.snippet) into the project `.gitignore` so desktop state stays local.
2. Clone the team repo (`KroniK907/agent-config`) somewhere stable.
3. From the **project** root, run the wizard (see [scripts/wizard/README.md](../scripts/wizard/README.md)).
4. Toggle skills and rules in the opt-in tree. Groups cascade to children. Rules appear above skills.
5. Optional: enable **Environment details rule** at the top of the tree - writes `.cursor/rules/environment.mdc` on apply.
6. Press **a** to apply. Wizard writes `.cursor/agent-config.local.json`, copies enabled items, and refreshes env-details when enabled.

## Re-run

Re-run the wizard from the project root any time the catalog changes or you want to adjust opt-in.

| Behavior | Detail |
|----------|--------|
| New catalog entries | Marked **NEW** in the tree; default off until you opt in |
| Deselected managed copies | Removed on apply (`agent-config-sync: true` in source frontmatter) |
| Project override | Local file with `agent-config-sync: false` shows **(project override)**; apply skips and deselect does not delete |
| Env details | Silent refresh on every apply when enabled in manifest |
| State pane | Tab to view/edit manifest JSON; apply saves selections |

When desktop opt-in is stable, copy the skill/rule list into `.cursor/agent-manifest.json` for cloud agents. Pin `source.ref` to a semver tag matching [`catalog.json`](../catalog.json) `catalog.version`.

## Cloud agents

Commit `.cursor/agent-manifest.json` with the subset cloud agents need. Build hook template: [examples/environment.json.example](examples/environment.json.example). Cloud copy-only bootstrap is separate from the desktop wizard - see map decision log for the two-script model.

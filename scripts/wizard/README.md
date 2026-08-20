# agent-config-wizard

Desktop TUI for applying team agent-config to a project repo: hierarchical skill/rule opt-in, manifest persistence, managed copy lifecycle.

Validated in prototype [#51](https://github.com/KroniK907/agent-config/issues/51). Promotion to production apply is a follow-on write-code task.

## Layout

| Path | Purpose |
|------|---------|
| [src/main.go](src/main.go) | Wizard source (single package) |
| [src/go.mod](src/go.mod) | Go module for bubbletea deps |
| [agent-config-wizard.exe](agent-config-wizard.exe) | Windows binary (rebuild after src changes) |

## Two roots

| Root | Resolves from |
|------|----------------|
| Team | Walk up from binary or `src/` until `catalog.json` is found |
| Project | cwd, or `-project` flag (manifest + `.cursor/` destination) |

The agent-config clone does **not** need to live inside the project repo.

## Run

Go resolves modules from the `-C` directory. Use `-C` when cwd is another project (especially one with its own `go.mod`).

### Separate project (typical)

Built binary from project cwd:

```powershell
cd C:\Users\DanielKranich\Documents\Projects\jrdev
C:\Users\DanielKranich\.cursor\skills\scripts\wizard\agent-config-wizard.exe
```

Or dev with explicit project:

```powershell
go run -C C:\Users\DanielKranich\.cursor\skills\scripts\wizard\src . -project C:\Users\DanielKranich\Documents\Projects\jrdev
```

### Same repo (team = project)

```powershell
cd C:\Users\DanielKranich\.cursor\skills
go run -C .\scripts\wizard\src .
```

## Build

```powershell
go build -C src -o ../agent-config-wizard.exe .
```

## Tests

From the repo root:

```powershell
go test -C scripts/wizard/src ./...
```

From `scripts/wizard/src`:

```powershell
go test ./...
```

Coverage spans manifest load/save, catalog tree and group cascade, project override detection, apply copy/remove/skip, and new-entry defaults.

## Keys

| Key | Action |
|-----|--------|
| Up/Down | Move selection (list) or scroll one line (state) |
| PgUp/PgDn | Page scroll in current pane |
| Home/End | Jump to top/bottom of current pane |
| Tab | Switch between opt-in tree and state JSON |
| Space | Toggle opt-in (list pane); groups cascade to children |
| a | Apply - copy enabled, remove deselected managed copies |
| s | Jump to state pane |
| q | Save manifest and quit |

## Behavior

- Rules above skills; nested groups (`Wayfinder`, `actions/`, etc.) with cascade select
- Manifest at `.cursor/agent-config.local.json` (see `.cursor/examples/agent-config.local.json.example`)
- `agent-config-sync: true` frontmatter on managed copies; `false` = detached override
- Re-run: `lastCatalogPaths` tracks seen entries; new catalog paths show **NEW** and default off
- Project-native files (no sync marker, never applied) show **(project override)**
- Copy delivery to `.cursor/rules/` and `.cursor/skills/`; framework copy to `.cursor/agent-config/`
- Optional **Environment details rule** at top of tree - toggles `envDetails` in manifest; writes `.cursor/rules/environment.mdc` on apply (`alwaysApply: true`, `generated-by: agent-config-wizard`); silent refresh on every apply when enabled; removes only wizard-generated copies when disabled

## Collision demo

1. Run apply with `unslop` rule enabled.
2. Edit `.cursor/rules/unslop.mdc` locally and set `agent-config-sync: false` (simulate detach).
3. Re-run - rule shows **(project override)**; apply skips it; deselect does not delete detached copy.

## Cleanup

Run from the **project** directory:

```powershell
Remove-Item -Recurse -Force .cursor/agent-config.local.json, .cursor/agent-config, .cursor/skills, .cursor/rules -ErrorAction SilentlyContinue
```

Only remove paths created for wizard testing.

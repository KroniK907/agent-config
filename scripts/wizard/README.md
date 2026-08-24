# agent-config-wizard

Desktop TUI for applying team agent-config to a project repo: hierarchical skill/rule opt-in, manifest persistence, managed copy lifecycle.

Primary desktop apply path - see repo [README](../../README.md#apply-to-a-project-primary).

## Layout

| Path | Purpose |
|------|---------|
| [src/main.go](src/main.go) | Wizard source - single file (TUI, apply, env-details probe) |
| [src/main_test.go](src/main_test.go) | All wizard tests |
| [src/go.mod](src/go.mod) | Go module for bubbletea deps |
| [bin/agent-config-wizard](bin/agent-config-wizard) | Linux binary (built by CI on `main`) |
| [bin/agent-config-wizard.exe](bin/agent-config-wizard.exe) | Windows binary (built by CI on `main`) |

Per **AGENT-CFG-GM-010**, wizard Go code stays in one source file plus one test file - no extra `.go` modules under `src/`.

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
C:\Users\DanielKranich\.cursor\skills\scripts\wizard\bin\agent-config-wizard.exe
```

On Linux:

```bash
cd ~/projects/jrdev
/path/to/agent-config/scripts/wizard/bin/agent-config-wizard
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

CI builds tracked copies under `bin/` on every push to `main`. Download the latest run artifact from the [Validate catalog workflow](https://github.com/KroniK907/agent-config/actions/workflows/validate-catalog.yml) if you need binaries before the bot commit lands.

For local testing only (gitignored):

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
| a | Apply desktop config - copy enabled, remove deselected managed copies, save manifest |
| c | Enter/exit **cloud bootstrap** mode |
| l | (cloud mode) Sync cloud opt-in from desktop selections; shows diff in state pane |
| [ / ] | (cloud mode) Cycle `source.ref` through GitHub release tags |
| w | (cloud mode) Write `.cursor/agent-manifest.json` and `.cursor/environment.json` |
| s | Jump to state pane |
| q | Quit without saving TUI selections |

## Behavior

- Rules above skills; nested groups (`Wayfinder`, `actions/`, etc.) with cascade select
- Manifest at `.cursor/agent-config.local.json` (see `.cursor/examples/agent-config.local.json.example`)
- `agent-config-sync: true` frontmatter on managed copies; `false` = detached override
- Re-run: `lastCatalogPaths` tracks seen entries; new catalog paths show **NEW** and default off
- Project-native files (no sync marker, never applied) show **(project override)**
- Copy delivery to `.cursor/rules/` and `.cursor/skills/`; framework copy to `.cursor/agent-config/`
- Optional **Environment details rule** at top of tree - toggles `envDetails` in manifest; writes `.cursor/rules/environment.mdc` on apply (`alwaysApply: true`, `generated-by: agent-config-wizard`); silent refresh on every apply when enabled; removes only wizard-generated copies when disabled
- **Cloud bootstrap** (`c`) - optional committed cloud config: separate skills/rules opt-in, GitHub release tag picker, path validation at pinned ref, writes `.cursor/agent-manifest.json` + `.cursor/environment.json` (`w`). Desktop apply (`a`) does not write cloud files.

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

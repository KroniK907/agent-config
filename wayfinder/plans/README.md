# Wayfinder plans

**GitHub issues are the canonical tracker** for wayfinder maps on this repo.

## Wayfinder-Ecosystem (bootstrap)

| Artifact | Issue |
|----------|-------|
| **Wayfinder-Ecosystem:Map** | https://github.com/KroniK907/skills/issues/12 |
| **Wayfinder-Ecosystem:Decision-Log** | https://github.com/KroniK907/skills/issues/11 |

To Do tickets `#4`–`#10` are sub-issues of the map. Native **blocked-by** dependencies are wired on the repo.

### Labels to add (manual)

The bootstrap token could not create custom labels. Add these on GitHub when convenient:

`wayfinder:map`, `wayfinder:decision-log`, `wayfinder:todo`, `wayfinder:research`, `wayfinder:prototype`, `wayfinder:grilling`, `wayfinder:task`, `wayfinder:hitl`, `wayfinder:afk`

Then apply them to issues `#4`–`#12`.

### Cleanup

Close spurious test issues `#2` and `#3` if still open (created during API permission probing).

## Local files

Use `wayfinder/plans/{FeatureName}.Map.md` only when GitHub issues are unavailable (offline) or as a Windows-friendly export. Do not duplicate the canonical map in the repo while GitHub is enabled.

Re-bootstrap script: [`../scripts/bootstrap-github-map.sh`](../scripts/bootstrap-github-map.sh) (creates new issues — do not re-run on an existing map).

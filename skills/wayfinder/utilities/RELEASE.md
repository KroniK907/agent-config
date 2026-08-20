# Skills repo releases

Semver tags on **`KroniK907/agent-config`** are the pin target for app-repo [AFK bootstrap](AFK-BOOTSTRAP.md) (`WAYFINDER_SKILLS_TAG` in `.cursor/environment.json`).

---

## Versioning

Use [Semantic Versioning](https://semver.org/):

| Bump | When |
|------|------|
| **MAJOR** | Breaking changes to wayfinder task contracts, implement-task gates, or skill paths that require app-repo config updates |
| **MINOR** | New skills, new wayfinder labels, backward-compatible features |
| **PATCH** | Doc fixes, non-breaking corrections |

First release after AFK v1 bundle lands: **`v0.1.0`**.

---

## Release checklist (maintainer)

Run on **`main`** after the AFK v1 bundle PR merges and all Implementing tasks are reconciled:

1. **Verify main** - `orchestrators/implement-task`, `actions/` playbooks, `utilities/bootstrap/`, and hub REFERENCE rows present on `main`.
2. **Tag** - annotated tag preferred:

   ```bash
   git checkout main
   git pull origin main
   git tag -a v0.1.0 -m "AFK v1 - implement-task, bootstrap pack, wayfinder layout"
   git push origin v0.1.0
   ```

3. **GitHub Release** - publish official release for the tag:

   ```bash
   gh release create v0.1.0 --title "v0.1.0 - AFK v1" --notes-file release-notes.md
   ```

   Summarize: implement-task orchestration, bootstrap pack, `wf:*` labels, four-folder layout (`actions/`, `ideation/`, `orchestrators/`, `utilities/`), write-code + design-modules, one-off entry path.

4. **Announce** - note tag in bundle / map Completed gist; app repos bump `WAYFINDER_SKILLS_TAG`.
5. **App repos** - follow [AFK-BOOTSTRAP.md § Pin skills pack](AFK-BOOTSTRAP.md#2-pin-skills-pack-environmentjson); trigger Cloud Agent environment rebuild.

Subsequent releases repeat with the next semver (`v0.1.1`, `v0.2.0`, …).

---

## Pre-release validation

Before tagging:

- [ ] Bundle branch merged to `main` (or release branch policy satisfied)
- [ ] `wayfinder/orchestrators/implement-task/` and `wayfinder/utilities/bootstrap/` paths stable
- [ ] [labels-manifest.json](bootstrap/labels-manifest.json) includes all labels referenced in REFERENCE
- [ ] Smoke: `install-skills.sh` with `WAYFINDER_SKILLS_TAG=<candidate>` populates `~/.cursor/skills/wayfinder/orchestrators/implement-task/SKILL.md`

---

## Consumer pin

App repo `.cursor/environment.json`:

```json
{
  "env": {
    "WAYFINDER_SKILLS_REPO": "KroniK907/agent-config",
    "WAYFINDER_SKILLS_TAG": "v0.1.3"
  }
}
```

Never pin to a moving branch (`main`) in production AFK repos - always an exact tag.

---

## Not in scope

- npm / `npx skills@latest` publish flow (orthogonal install path for HITL local dev)

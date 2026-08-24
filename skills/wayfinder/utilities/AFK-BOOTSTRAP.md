# AFK bootstrap checklist

Cross-repo setup for **wayfinder AFK v1** unattended implementation pickup. Complete every step in an **app implementation repo** before adding `wf:afk` tasks or enabling comment-trigger automation.

**Desktop project skills:** use **agent-config-wizard** ([scripts/wizard/README.md](../../../scripts/wizard/README.md)) - not this checklist. This doc is for Cloud AFK orchestration only.

**Binding contract:** tracker lives in each app repo; skills come from [`KroniK907/agent-config`](https://github.com/KroniK907/agent-config) pinned to a **semver tag**; one automation per repo; agents **never open PRs** - bundle branch + resolution comment only.

---

## Prerequisites

- [ ] GitHub issues enabled on the app repo (wayfinder map tracker)
- [ ] `gh` CLI authenticated with issue write on the app repo
- [ ] `jq` available for Unix label bootstrap script (`bootstrap-labels.sh`)
- [ ] Cursor Cloud Agents / Automations access for the org or user
- [ ] Skills repo **semver tag** published (see [RELEASE.md](RELEASE.md) - pin `v0.1.0` or later)

---

## 1. Bootstrap GitHub labels

Run once per app repo from any checkout that includes this skills pack:

```powershell
# Windows (repo root = app repo)
.\wayfinder\utilities\bootstrap\bootstrap-labels.ps1
```

```bash
# macOS / Linux / Cloud Agent shell
bash wayfinder/utilities/bootstrap/bootstrap-labels.sh
```

Manifest: [bootstrap/labels-manifest.json](bootstrap/labels-manifest.json) - includes **`wf:afk-running`** (serial queue lock) and **`wf:needs-review`** (approval gate pending).

Verify:

```powershell
gh label list --limit 100 | Select-String wf:
```

---

## 2. Pin skills pack (manifest + environment.json)

**Pin target:** `.cursor/agent-manifest.json` `source.ref` (semver tag). Do not use an `env` block in `environment.json` for repo/tag pin - that pattern is deprecated per AgentConfigHub map decisions.

1. Copy [`.cursor/examples/agent-manifest.json.example`](../../../.cursor/examples/agent-manifest.json.example) to **`.cursor/agent-manifest.json`** in the app repo. Set **`source.ref`** to the exact semver tag you released (e.g. `v1.0.0`). List enabled **`skills`** and **`rules`** paths from [`catalog.json`](../../../catalog.json) at that tag.
2. Copy [`.cursor/examples/environment.json.example`](../../../.cursor/examples/environment.json.example) to **`.cursor/environment.json`** in the app repo root. Update the tag in the `install` curl URL to match **`source.ref`**.
3. Commit both files under `.cursor/`.

The **`install`** command runs on Cloud Agent Build creation. It invokes [`scripts/bootstrap-agent.sh`](../../../scripts/bootstrap-agent.sh), which reads the committed manifest, clones `source.repo` at `source.ref`, validates paths against `catalog.json`, copies skills to `~/.cursor/skills/`, and copies rules to `.cursor/rules/` in the workspace. The script must be **idempotent**.

```json
{
  "build": {
    "install": "curl -fsSL https://raw.githubusercontent.com/KroniK907/agent-config/v1.0.0/scripts/bootstrap-agent.sh | bash"
  }
}
```

**Legacy (deprecated):** copying [bootstrap/install-skills.sh](bootstrap/install-skills.sh) to `.cursor/install-wayfinder-skills.sh` with inline `WAYFINDER_SKILLS_TAG` env vars. Use only when bootstrap-agent.sh is unavailable at your pinned tag.

**Local smoke (manifest-driven):**

```bash
# From app repo root with .cursor/agent-manifest.json committed
export AGENT_CONFIG_WORKSPACE="$PWD"
bash /path/to/agent-config/scripts/bootstrap-agent.sh
test -f ~/.cursor/skills/wayfinder/orchestrators/implement-task/SKILL.md
test -f .cursor/rules/unslop.mdc   # when rules/unslop.mdc is in manifest
```

**Local smoke (legacy global install):**

```powershell
$env:WAYFINDER_SKILLS_TAG = "v1.0.0"
.\skills\wayfinder\utilities\bootstrap\install-skills.ps1
```

Legacy scripts copy a fixed wayfinder skill set globally. Prefer **agent-config-wizard** for per-project desktop setup and **bootstrap-agent.sh** for cloud agents.

---

## 3. Configure GH_TOKEN

Cloud agents need **`gh`** and GitHub API access for issue edits, comments, and label operations.

| Source | When |
|--------|------|
| Cursor dashboard **Secrets** | Default - set `GH_TOKEN` (PAT or fine-grained token with repo + issues scope) |
| `environment.json` **`env.GH_TOKEN`** | Optional override when dashboard secret is not set |

Do **not** commit tokens. Verify in a Cloud Agent shell:

```bash
gh auth status
```

---

## 4. Duplicate Cursor automation (one per repo)

1. Open Cursor **Automations** for the app repo.
2. Create **one** repo-scoped automation - trigger: **issue comment** containing **`Approved - AFK implement`** (exact phrase).
3. Paste prompt from [bootstrap/automation-prompt.md](bootstrap/automation-prompt.md).
4. **Disable PR creation** in automation settings.
5. Save and note the automation name for your runbook.

The prompt references **`implement-task`** as the sole orchestration entry. Task bodies are identical for HITL and AFK - contract lives in implement-task + this automation. Skills add label **`wf:approved`** for reviewer visibility; the **comment phrase** starts the automation until issue-label triggers are supported ([afk-pickup-comment.md](../orchestrators/implement-task/references/afk-pickup-comment.md)).

---

## 5. HITL smoke before unattended pickup

Run at least **one** implementation task manually before enabling AFK on production frontier work:

1. Chart / define-bundle / create-tasks through to a **`wf:approved`** HITL task with **Status:** `ready`.
2. In chat: `/implement-task` on that task (or invoke implement-task skill with issue `#N`).
3. Confirm: bundle branch checkout, Method build, code-review, push, resolution comment, **Status:** `awaiting-reconcile`.
4. Reconcile with **`Approved - reconcile and close`**.

Only after HITL smoke passes:

- [ ] Add **`wf:afk`** label to AFK-mode tasks at create-tasks approval time
- [ ] Enable the comment automation from step 4 (phrase **`Approved - AFK implement`**)
- [ ] Confirm serial queue: only one open issue should hold **`wf:afk-running`** at a time
- [ ] Confirm first AFK pickup posts **`wf:approved`** + pickup comment (create-tasks or implement-task handoff)

---

## 6. Ongoing operations

| Event | Action |
|-------|--------|
| Skills pack update | Cut new semver tag in skills repo ([RELEASE.md](RELEASE.md)); bump **`source.ref`** in app `.cursor/agent-manifest.json`; update tag in `.cursor/environment.json` **install** curl URL; rebuild Cloud Agent environment |
| New wayfinder label | Add to [labels-manifest.json](bootstrap/labels-manifest.json) in skills repo; re-run bootstrap script in app repos |
| Bundle complete | Human opens **one PR** from `afk/bundle-{N}-{slug}` - agents do not |
| Task shipped | Human Reconcile **`Approved - reconcile and close`** per task resolution comment |

---

## File index

| Path | Purpose |
|------|---------|
| [AFK-BOOTSTRAP.md](AFK-BOOTSTRAP.md) | This checklist |
| [bootstrap/labels-manifest.json](bootstrap/labels-manifest.json) | Canonical `wf:*` labels |
| [bootstrap/bootstrap-labels.ps1](bootstrap/bootstrap-labels.ps1) | Label bootstrap (Windows) |
| [bootstrap/bootstrap-labels.sh](bootstrap/bootstrap-labels.sh) | Label bootstrap (Unix) |
| [bootstrap/environment.json.example](bootstrap/environment.json.example) | Legacy Cloud Agent env template (use `.cursor/examples/environment.json.example` instead) |
| [bootstrap/install-skills.ps1](bootstrap/install-skills.ps1) | Legacy local / Windows skills install |
| [bootstrap/install-skills.sh](bootstrap/install-skills.sh) | Legacy Cloud Agent skills install |
| [scripts/bootstrap-agent.sh](../../../scripts/bootstrap-agent.sh) | Manifest-driven Cloud Agent bootstrap (primary) |
| [bootstrap/automation-prompt.md](bootstrap/automation-prompt.md) | Cursor automation prompt template |
| [RELEASE.md](RELEASE.md) | Skills repo semver release process |
| [implement-task/SKILL.md](../../orchestrators/implement-task/SKILL.md) | AFK/HITL orchestration contract |
| [create-tasks/REFERENCE.md](../../actions/create-tasks/REFERENCE.md) | Implementation task body template |

Linked from [REFERENCE.md - Ecosystem integration](REFERENCE.md#ecosystem-integration).

# AFK bootstrap checklist

Cross-repo setup for **wayfinder AFK v1** unattended implementation pickup. Complete every step in an **app implementation repo** before adding `wf:afk` tasks or enabling label-trigger automation.

**Binding contract:** tracker lives in each app repo; skills come from [`KroniK907/skills`](https://github.com/KroniK907/skills) pinned to a **semver tag**; one automation per repo; agents **never open PRs** - bundle branch + resolution comment only.

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

## 2. Pin skills pack (environment.json)

1. Copy [bootstrap/environment.json.example](bootstrap/environment.json.example) to **`.cursor/environment.json`** in the app repo root.
2. Copy [bootstrap/install-skills.sh](bootstrap/install-skills.sh) to **`.cursor/install-wayfinder-skills.sh`** in the app repo (same script; stable path for the install command).
3. Set **`WAYFINDER_SKILLS_TAG`** in `environment.json` to the exact semver tag you released (e.g. `v0.1.0`).
4. Commit both files under `.cursor/`.

The **`install`** command runs on Cloud Agent Build creation. It must be **idempotent** - see [bootstrap/install-skills.sh](bootstrap/install-skills.sh).

**Local smoke (optional):**

```powershell
$env:WAYFINDER_SKILLS_TAG = "v0.1.0"
.\wayfinder\utilities\bootstrap\install-skills.ps1
```

Skills land in `~/.cursor/skills/` (wayfinder hub + actions + repo-root utilities synced by the script).

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
2. Create **one** repo-scoped automation - trigger: **`wf:approved` label added**.
3. Paste prompt from [bootstrap/automation-prompt.md](bootstrap/automation-prompt.md).
4. **Disable PR creation** in automation settings.
5. Save and note the automation name for your runbook.

The prompt references **`implement-task`** as the sole orchestration entry. Task bodies are identical for HITL and AFK - contract lives in implement-task + this automation.

---

## 5. HITL smoke before unattended pickup

Run at least **one** implementation task manually before enabling AFK on production frontier work:

1. Chart / define-bundle / create-tasks through to a **`wf:approved`** HITL task with **Status:** `ready`.
2. In chat: `/implement-task` on that task (or invoke implement-task skill with issue `#N`).
3. Confirm: bundle branch checkout, Method build, code-review, push, resolution comment, **Status:** `awaiting-reconcile`.
4. Reconcile with **`Approved - reconcile and close`**.

Only after HITL smoke passes:

- [ ] Add **`wf:afk`** label to AFK-mode tasks at create-tasks approval time
- [ ] Enable the label automation from step 4
- [ ] Confirm serial queue: only one open issue should hold **`wf:afk-running`** at a time

---

## 6. Ongoing operations

| Event | Action |
|-------|--------|
| Skills pack update | Cut new semver tag in skills repo ([RELEASE.md](RELEASE.md)); bump `WAYFINDER_SKILLS_TAG` in app `.cursor/environment.json`; rebuild Cloud Agent environment |
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
| [bootstrap/environment.json.example](bootstrap/environment.json.example) | Cloud Agent env template |
| [bootstrap/install-skills.ps1](bootstrap/install-skills.ps1) | Local / Windows skills install |
| [bootstrap/install-skills.sh](bootstrap/install-skills.sh) | Cloud Agent skills install |
| [bootstrap/automation-prompt.md](bootstrap/automation-prompt.md) | Cursor automation prompt template |
| [RELEASE.md](RELEASE.md) | Skills repo semver release process |
| [implement-task/SKILL.md](../../orchestrators/implement-task/SKILL.md) | AFK/HITL orchestration contract |
| [create-tasks/REFERENCE.md](../../actions/create-tasks/REFERENCE.md) | Implementation task body template |

Linked from [REFERENCE.md - Ecosystem integration](REFERENCE.md#ecosystem-integration).

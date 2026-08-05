# Action skill patterns

Shared scaffold for **`wayfinder/actions/<name>/`** skills — Method playbooks dispatched by [`implement-task`](../implement-task/SKILL.md) (task #29). Map-frontier skills (e.g. `research`, `grill-me`) live as **`wayfinder/<skill>/`** peers, not under `actions/`.

Each action skill is a **`wayfinder/**/<name>/SKILL.md`** entry in the Method pool ([WF-ECO-GM-030](https://github.com/KroniK907/skills/issues/11)). `create-tasks` sets **## Method** to the skill `name` from frontmatter.

## Mandatory sections

Every action skill **REFERENCE.md** (or equivalent detail doc referenced from a minimal SKILL.md) must include these five sections, in order:

### 1. Output artifact

What tangible deliverable this Method produces for the implementation task — file paths, issue comment shape, branch state, or other verifiable artifact. Be specific enough that `implement-task` resolution **Done when** can cite it.

### 2. Prerequisites

What must be true before the action runs: loaded task + bundle context, bundle branch checked out, upstream dependencies merged, env/tools, or human inputs. List hard gates that should fail closed when unmet.

### 3. Workflow

Step-by-step playbook the delegated agent follows after `implement-task` startup gates pass. Numbered steps; no orchestration duplicated from `implement-task` (git push, resolution comment, queue handoff stay in the orchestrator).

### 4. Done when mapping

How each bullet in the task issue **Done when** maps to workflow steps or output checks. Table or bullet mapping preferred — resolution comment **Done when** section copies this mapping at end-of-run.

### 5. Division of labor

Explicit split between **`implement-task`** (orchestration: gates, branch, push, resolution comment, unblock/handoff, never close task) and **this action skill** (build work only). Call out AFK vs HITL differences only when they affect the action itself, not queue semantics.

## SKILL.md conventions

- YAML frontmatter with unique `name` matching **## Method** on tasks
- Short description for agent discovery
- Body: when to use, pointer to REFERENCE, link to this file
- Keep SKILL.md minimal; detail lives in REFERENCE and bundled references

## Adding a new action skill

1. Create `wayfinder/actions/<name>/SKILL.md` + REFERENCE implementing all five sections above
2. Register in wayfinder REFERENCE routing / ecosystem tables when the Method replaces an inline or root skill default
3. Ensure `create-tasks` can set **Method:** `<name>` for the ticket types that use it

**Deferred in task #28:** `actions/prototype/` — first concrete action skill lands in task #29.

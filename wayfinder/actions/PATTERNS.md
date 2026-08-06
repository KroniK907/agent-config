# Action skill patterns

Shared scaffold for **`wayfinder/actions/<name>/`** skills - focused playbooks for map-scoped work. Invoked **directly** (Route / user attach), via ticket **## Method** (`create-tasks`, `one-off`, `implement-task`), or **delegated** by [orchestrators](../orchestrators/implement-task/SKILL.md) after startup gates pass.

**Ideation** skills (`grill-me`, `feature-discovery`, …) live under **`wayfinder/ideation/`**. **Orchestrators** (`wayfinder`, `one-off`, `implement-task`) live at the hub root or under **`wayfinder/orchestrators/`**.

Each action skill is a **`wayfinder/**/<name>/SKILL.md`** entry in the Method pool. Default pool: skills at **`wayfinder/**/<name>/SKILL.md`** in the pinned pack; repo-root one-offs are valid only when task **## Method** explicitly names them. `create-tasks` sets **## Method** to the skill `name` from frontmatter.

## Mandatory sections

Every action skill **REFERENCE.md** (or equivalent detail doc referenced from a minimal SKILL.md) must include these five sections, in order:

### 1. Output artifact

What tangible deliverable this skill produces - file paths, issue comment shape, branch state, draft bundle issue, or other verifiable artifact. Be specific enough that resolution **Done when** (task, bundle, or one-off ticket) can cite it.

### 2. Prerequisites

What must be true before the action runs: loaded ticket + map/bundle context, branch checked out, upstream dependencies merged, env/tools, or human inputs. List hard gates that should fail closed when unmet.

### 3. Workflow

Step-by-step playbook the agent follows once entry gates pass. Numbered steps; no orchestration duplicated from orchestrators (git push, resolution comment, queue handoff, Reconcile approval phrases stay in the orchestrator or wayfinder hub).

### 4. Done when mapping

How each bullet in the ticket **Done when** (or skill-specific completion checklist) maps to workflow steps or output checks. Table or bullet mapping preferred - resolution comment **Done when** section copies this mapping at end-of-run when applicable.

### 5. Division of labor

Explicit split between the **entry orchestrator** (when any - e.g. [implement-task](../orchestrators/implement-task/SKILL.md), [one-off](../orchestrators/one-off/SKILL.md), or direct HITL invoke) and **this action skill** (domain work only). For skills invoked directly with no orchestrator tail, state that explicitly. Call out AFK vs HITL differences only when they affect the action itself.

## SKILL.md conventions

- YAML frontmatter with unique `name` matching **## Method** on tasks (when used as Method)
- Short description for agent discovery
- Body: when to use, pointer to REFERENCE, link to this file
- Keep SKILL.md minimal; detail lives in REFERENCE and bundled references

## Adding a new action skill

1. Create `wayfinder/actions/<name>/SKILL.md` + REFERENCE implementing all five sections above
2. Register in wayfinder REFERENCE routing / ecosystem tables when the skill becomes a Route or Method default
3. Ensure `create-tasks` can set **Method:** `<name>` when the ticket type uses it

**Examples:** [write-code](write-code/SKILL.md) and [prototype](prototype/SKILL.md) (bundle build Methods); [research](research/SKILL.md), [define-bundle](define-bundle/SKILL.md), [create-tasks](create-tasks/SKILL.md), [design-modules](design-modules/SKILL.md), [code-review](code-review/SKILL.md) (direct or orchestrated entry).

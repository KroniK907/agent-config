# Design it twice

When interface shape is open, explore alternatives with parallel sub-agents. Based on "Design It Twice" (Ousterhout) - your first idea is unlikely to be the best.

Uses vocabulary from [REFERENCE.md](REFERENCE.md) - **module**, **interface**, **seam**, **adapter**, **leverage**.

Adapted from [mattpocock/skills - engineering/codebase-design/DESIGN-IT-TWICE](https://github.com/mattpocock/skills/blob/main/skills/engineering/codebase-design/DESIGN-IT-TWICE.md).

## Process

### 1. Frame the problem space

Before spawning sub-agents, write a user-facing explanation:

- Constraints any new interface must satisfy
- Dependencies and category ([DEEPENING.md](DEEPENING.md))
- Rough illustrative sketch - not a proposal

Show this to the user, then proceed to step 2. The user reads while sub-agents work in parallel.

### 2. Spawn sub-agents

Spawn 3+ sub-agents in parallel via Task tool. Each must produce a **radically different** interface.

Prompt each sub-agent with a separate technical brief (bundle **Decisions**, file paths when known, dependency category, what sits behind the seam). Give each agent a different design constraint:

- Agent 1: "Minimize the interface - aim for 1-3 entry points max. Maximise leverage per entry point."
- Agent 2: "Maximise flexibility - support many use cases and extension."
- Agent 3: "Optimise for the most common caller - make the default case trivial."
- Agent 4 (if applicable): "Design around ports and adapters for cross-seam dependencies."

Each sub-agent outputs:

1. Interface (types, methods, params - plus invariants, ordering, error modes)
2. Usage example
3. What the implementation hides behind the seam
4. Dependency strategy and adapters ([DEEPENING.md](DEEPENING.md))
5. Trade-offs - where leverage is high, where it is thin

### 3. Present and compare

Present designs sequentially. Compare in prose by **depth** (leverage at the interface), **locality** (where change concentrates), and **seam placement**.

Give a recommendation: strongest design and why. Propose a hybrid when elements from different designs combine well. Be opinionated - the user wants a strong read, not just a menu.

### 4. Capture in artifact

Fold the recommendation into the [module-design artifact](REFERENCE.md#module-design-artifact-template) on the bundle or ticket comment.

## Anti-patterns

- Do not let sub-agents produce similar designs - enforce radical difference
- Do not skip comparison - the value is in contrast
- Do not implement - this step is interface shape only unless the human explicitly requests spike code
- Do not evaluate based on implementation effort alone

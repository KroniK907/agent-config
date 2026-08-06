# Design-module reference

HITL module-design skill. Implements bundle and planning entry paths. Vocabulary adapted from [mattpocock/skills - engineering/codebase-design](https://github.com/mattpocock/skills/tree/main/skills/engineering/codebase-design).

---

## Glossary

Use these terms exactly - consistent language is the point.

**Module** - anything with an interface and an implementation. Scale-agnostic: function, class, package, or tier-spanning slice.

**Interface** - everything a caller must know to use the module: signature, invariants, ordering, error modes, configuration, performance characteristics.

**Implementation** - code inside the module. Distinct from **Adapter** (concrete thing at a seam).

**Depth** - leverage at the interface: behaviour exercisable per unit of interface learned. **Deep** = lots of behaviour behind a small interface; **shallow** = interface nearly as complex as implementation.

**Seam** _(Michael Feathers)_ - where behaviour can change without editing in that place; where the module's interface lives.

**Adapter** - concrete thing satisfying an interface at a seam.

**Leverage** - capability callers get from depth.

**Locality** - maintainers get change concentrated in one place.

See full definitions and principles in workflow step 2.

---

## Prerequisites

- `gh` authenticated on the target repo
- **Bundle entry:** `wf:bundle`, **Status:** `approved`; parent map + decision log links
- **Planning entry:** open ticket with **Question** and **Map** link, or human-declared module question on a map
- HITL only - never add **`wf:afk`**

**Gate (bundle):** If bundle **Status** is `draft`, hand off to [define-bundle](../define-bundle/SKILL.md) for **`bundle approved`**.

**Gate (planning):** If the question is fact-gathering not interface shape, hand off to [research](../research/SKILL.md).

---

## Workflow

### 1. Load context

**Bundle entry:**

```text
gh issue view <bundle-num> --json body,title,url,labels
gh issue view <map-num> --json body,title,url
```

From bundle: **Decisions**, **Constraints**, scope summary, boundaries, open questions.

**Planning entry:**

```text
gh issue view <ticket-num> --json body,title,url,labels
gh issue view <map-num> --json body,title,url
```

From ticket: **Question**, **Done when**, **Map** link.

### 2. Frame the module

From loaded context, write a user-facing problem-space summary:

- Module name (working)
- Problem the module solves; callers
- Constraints from **Decisions** / **Constraints** / ticket **Question**
- Dependencies and category (see [DEEPENING.md](DEEPENING.md))
- Rough illustrative sketch - not a proposal, just grounding

Apply deep-module principles:

- **Deletion test** - if deleting the module makes complexity vanish, it was pass-through; if complexity reappears across callers, it earns its keep
- **Interface is the test surface** - callers and tests cross the same seam
- **One adapter = hypothetical seam; two adapters = real seam**

Optional: short [grill-me](../grill-me/SKILL.md) pass when bundle **Open questions** are interface-shaped.

### 3. Explore interfaces (when shape is open)

When more than one viable interface exists, follow [DESIGN-IT-TWICE.md](DESIGN-IT-TWICE.md): parallel sub-agents, present sequentially, compare on depth/locality/seam placement, give a recommendation.

When the bundle or ticket already specifies the interface tightly, skip sub-agents and document the chosen shape with rationale.

### 4. Post artifact comment

Post on the **bundle issue** (bundle entry) or **ticket issue** (planning entry) using the template below.

End with handoff:

- **Bundle:** suggest [create-tasks](../create-tasks/SKILL.md); note whether task split hints in the artifact should drive the split
- **Planning:** human review; Reconcile when ticket **Done when** satisfied

Do not edit the map, decision log, or bundle **Status**. Do not add **`wf:approved`** or implementation labels.

---

## Module-design artifact template

```markdown
## Module design: {module name}

**Entry:** bundle | planning
**Map:** [{FeatureName}:Map](map-url)

### Problem and callers

{1-2 paragraphs}

### Recommended interface

{Signature, invariants, error modes - the small surface}

### Seam and dependency category

- **Seam placement:** …
- **Category:** in-process | local-substitutable | remote-but-owned | true-external (see DEEPENING.md)
- **Adapters:** production vs test (when applicable)

### Alternatives considered

{Brief - only when design-it-twice ran}

### Recommendation

{Which shape and why - depth, locality, leverage}

### Task split hints

{Optional bullets for create-tasks - vertical slices, suggested Method per slice}

### Open gaps

{Anything still foggy - link grill or research if needed}
```

---

## Bundle vs planning deltas

| Topic | Bundle entry | Planning entry |
|-------|--------------|------------------|
| Comment target | Bundle issue | Ticket (or map comment if no ticket) |
| Input authority | Bundle **Decisions** verbatim | Ticket **Question** |
| Next step | create-tasks | Reconcile ticket when done |
| Repo edits | Only on explicit human request | Same |

---

## Interaction rules

1. **Comment-only default** - artifact is the deliverable; no product code unless asked
2. **HITL only** - no AFK pickup path in v1
3. **No tracker writes** - map, coverage, and log unchanged by this skill
4. **Replace design-an-interface** - all former `design-an-interface` routing uses **design-module**
5. **Portable docs** - no live issue numbers in committed skill markdown

---

## Route trigger (for wayfinder)

Suggest **design-module** when:

- User finished **`bundle approved`** and asks to shape the module before splitting tasks
- [define-bundle](../define-bundle/SKILL.md) hand off - optional step before create-tasks
- Map **To Do** `wf:prototype` ticket explores interface variants (planning entry)
- User mentions "design module", "deep module", "design it twice", or module interface shaping

Do **not** suggest for:

- Codebase-wide refactor RFCs → **improve-codebase-architecture**
- Throwaway bundle-branch demos → **prototype** action via implement-task
- Binding GM rows → **grill-me** + Reconcile

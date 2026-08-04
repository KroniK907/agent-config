# Reference: session state, map-discovery artifact

## Session state line (every assistant reply)

Emit one compact line after **Recap**:

`Session state: Phase DISCOVERY · Zone [name] · Round [current]/5 · Next: [next zone or post map-discovery comment]`

Examples:

- `Session state: Phase DISCOVERY · Zone Surfaces & experience · Round 1/5 · Next: Behavior & correctness`
- `Session state: Phase DISCOVERY · Zone Change, risk & evidence · Round 5/5 · Next: post map-discovery comment · Status ready for materialize`

**First reply:** Recap states seed and map issue #; Session state starts at Surfaces & experience 1/5.

## Map-discovery artifact

**Canonical store:** A comment on the **`{FeatureName}:Map`** issue. Use the exact heading `## Map discovery` so wayfinder **Materialize** can find it.

Also emit the block in chat on completion so the user can Materialize from the current session without re-fetching GitHub.

Layout:

```markdown
## Map discovery

**Target outcome:** [one line from map]
**Status:** in progress | ready for materialize | materialized

### Zone matrix

| Zone | Known | Unknown / needs research |
|------|-------|---------------------------|
| Surfaces & experience | … | … |
| Behavior & correctness | … | … |
| Boundaries & integration | … | … |
| Persistence & data | … | … |
| Change, risk & evidence | … | … |

### Ticket candidates

| Title | Type | Mode | Question |
|-------|------|------|----------|
| … | research / prototype / grilling / task | HITL / AFK | One-line ticket Question |

### Fog (not yet specified)

- …

### Out of scope suggestions

- …

### Notes

- …
```

### Ticket candidate types

| Type | When |
|------|------|
| `research` | Facts, prior art, constraints to gather |
| `prototype` | Stub, outline, or interface exploration |
| `grilling` | Single decision needing depth-first Q&A (`grill-me`) |
| `task` | Checklist, doc update, wiring, human errand |

### Research ticket shape (materialize)

When **Ticket candidates** Type is `research`, wayfinder **Materialize** creates a child issue with the [research ticket template](../research/REFERENCE.md#research-ticket-template):

| Field | Materialize action |
|-------|-------------------|
| **Question** column | Becomes `## Question` — one sharp investigation, fact-gathering not binding decision |
| **Done when** | Add section with 1–3 verifiable acceptance bullets derived from the question and zone context |
| **Map** | Auto-link parent map issue |
| **Source hints** | Optional — add when discovery notes name docs, repos, or code paths |
| **Perspectives** | Optional — add when discovery notes name stakeholders or alternate framings |

**Mode:** Default `HITL` for research in v1 ([research](../research/SKILL.md) skill). Labels: `wayfinder:todo` + `wayfinder:research` + `wayfinder:hitl`.

**Not research:** If the candidate Question is a binding decision or scope/strategy shape, use Type `grilling` instead (→ grill-me or strategic-ideation).

### Persist as map comment

1. Build artifact markdown with updated sections.
2. `gh issue comment <map-num> --body-file …`
3. On completion: **Status:** `ready for materialize`.

**Partial updates:** Post an updated comment with the full artifact so far and **Status:** `in progress`. Materialize uses the **latest** comment whose body contains `## Map discovery` and **Status:** `ready for materialize` (or user confirms an `in progress` draft).

**Local fallback:** Same sections in `wayfinder/plans/{FeatureName}.Map-Discovery.md` when GitHub is unavailable.

## Discovery vs strategic-ideation vs grill-me

| | feature-discovery | strategic-ideation | grill-me |
|--|-------------------|-------------------|----------|
| Layer | Edge-finding; inventory | Idea/strategy; scope shape | Implementation depth |
| Pace | One zone per reply (default) | Ideation/tension bands | One question per reply |
| Output | Map-discovery comment on map issue | Scope handoff (chat) | `{MAP-SLUG}-GM-xx` |

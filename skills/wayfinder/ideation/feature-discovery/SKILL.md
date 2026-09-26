---
name: feature-discovery
description: feature discovery, wayfinder Chart handoff, zone triage, map-discovery artifact, ready for materialize, breadth-first discovery, known vs unknown, five coverage zones, wayfinder Materialize, feature start, before grilling
agent-config-sync: true
---

# Feature discovery (breadth-first zone triage)

Walk all **five coverage zones** once at **triage depth** - collect what is **known**, **unknown**, or **needs research**. Do **not** depth-first grill; do **not** resolve unknowns. On completion, post a **[map-discovery artifact](REFERENCE.md#map-discovery-artifact)** as a **comment on the map issue** for wayfinder **Materialize**.

Typically runs **after** wayfinder **Chart** creates the map skeleton and decision log.

## Not this skill

| Skill | When instead |
|-------|----------------|
| [constrain-fog](../constrain-fog/SKILL.md) | Groom existing map **Not yet specified** fog - reuses zone matrix at triage depth per item |
| [strategic-ideation](../strategic-ideation/SKILL.md) | Scope/strategy expand → tension → prune; idea-level tradeoffs |
| [grill-me](../grill-me/SKILL.md) | Depth-first Q&A on one branch; resolves a single ticket **Question** |
| [wayfinder](../../SKILL.md) | Create map, materialize To Do tickets, reconcile, routing |

## Prerequisites

Wayfinder **Chart** must have created `{FeatureName}:Map` (with target outcome in the body).

**Local fallback:** If GitHub is unavailable, write the artifact to `wayfinder/utilities/plans/{FeatureName}.Map-Discovery.md` - same layout per [REFERENCE.md](REFERENCE.md#map-discovery-artifact).

## Coverage zones

Load [grill-me ZONES.md](../grill-me/ZONES.md) once per session and use each zone's **Skip when** / **In scope if** lines. Do not re-read it per zone.

# Wayfinder-Ecosystem:Map

**Phase:** deciding  
**Map slug:** `WF-ECO`  
**Decision log:** #{{LOG}} — prefix `WF-ECO-GM-`  
**Tracker:** GitHub issues (`KroniK907/skills`)

## Target outcome

A cohesive **wayfinder ecosystem**: opinionated `wayfinder` skill, updated sibling skills (`grill-me`, future `research`, PRD pipeline), and documented conventions so large features can be planned from ideation through implementation tickets — with human-closed Completed tickets and scoped decision logs.

## Notes

- **Repo:** `KroniK907/skills`
- **Model:** Composer 2.5 / Cursor Cloud Agents
- **Sibling skills:** `grill-me`, `write-a-prd`, `prd-to-issues`, `feature-ideation`, `design-an-interface`
- **Human closes all** To Do issues — agents post resolution only
- **Labels:** Apply `wayfinder:*` labels when available (see To Do ticket #{{T4}})

## Subfeatures

<!-- None yet -->

## To Do

| Ticket | Type | Mode | Issue |
|--------|------|------|-------|
| Define decision-log handoff to write-a-prd | grilling | HITL | #{{T1}} |
| Update grill-me to append scoped GM-xx to map decision log | task | HITL | #{{T2}} |
| Specify research ticket workflow (skill or Task pattern) | grilling | HITL | #{{T3}} |
| GitHub tracker setup for skills repo maps | task | AFK | #{{T4}} |
| Subfeature map worked example in REFERENCE or EXAMPLES | prototype | HITL | #{{T5}} |
| Reconcile / post-implementation map update skill or mode | grilling | HITL | #{{T6}} |
| Cloud AFK automation contract for wayfinder:afk tickets | grilling | HITL | #{{T7}} |

## Completed

- Initial wayfinder SKILL.md and REFERENCE.md (PR #1) — ideation interview, To Do/Completed, subfeatures, scoped decision log
- Bootstrap map migrated from local plan files to GitHub issues (WF-ECO-GM-007)

## Not yet specified

- **Canvas dashboard** — read-only view of map frontier
- **feature-ideation entry** — always before wayfinder, or only for raw seeds?

## Out of scope

- Implementing cloud automation itself (assumed future)
- Matt Pocock tracker compatibility
- Linear/Jira trackers

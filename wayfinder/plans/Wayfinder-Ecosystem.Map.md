# Wayfinder-Ecosystem:Map

**Phase:** deciding  
**Map slug:** `WF-ECO`  
**Decision log:** `wayfinder/plans/Wayfinder-Ecosystem.Decision-Log.md` — prefix `WF-ECO-GM-`  
**Bootstrap note:** This map plans the wayfinder skill ecosystem using wayfinder itself. Tracker = local markdown until GitHub issues are created on `KroniK907/skills`.

## Target outcome

A cohesive **wayfinder ecosystem**: opinionated `wayfinder` skill, updated sibling skills (`grill-me`, future `research`, PRD pipeline), and documented conventions so large features can be planned from ideation through implementation tickets — with human-closed Completed tickets and scoped decision logs.

## Notes

- **Repo:** `KroniK907/skills`
- **Model:** Composer 2.5 / Cursor Cloud Agents
- **Sibling skills:** `grill-me`, `write-a-prd`, `prd-to-issues`, `feature-ideation`, `design-an-interface`
- **Human closes all** `wayfinder:todo` issues — agents post resolution only
- **Matt Pocock:** informed by original wayfinder; this map intentionally diverges (To Do/Completed, ideation interview, decision-log prefix, subfeatures)

## Subfeatures

<!-- None yet — spawn when a child domain needs its own map (e.g. Cloud-AFK-Automation:Map) -->

## To Do

| Ticket | Type | Mode | Assignee | Blocked by |
|--------|------|------|----------|------------|
| [Define decision-log handoff to write-a-prd](#todo-decision-log-prd) | grilling | HITL | unclaimed | — |
| [Update grill-me to append scoped GM-xx to map decision log](#todo-grill-me-integration) | task | HITL | unclaimed | Define decision-log handoff to write-a-prd |
| [Specify research ticket workflow (skill or Task pattern)](#todo-research-workflow) | grilling | HITL | unclaimed | — |
| [GitHub tracker setup for skills repo maps](#todo-github-tracker) | task | AFK | unclaimed | — |
| [Subfeature map worked example in REFERENCE or EXAMPLES](#todo-subfeature-example) | prototype | HITL | unclaimed | — |
| [Reconcile / post-implementation map update skill or mode](#todo-reconcile) | grilling | HITL | unclaimed | Define decision-log handoff to write-a-prd |
| [Cloud AFK automation contract for wayfinder:afk tickets](#todo-afk-contract) | grilling | HITL | unclaimed | GitHub tracker setup for skills repo maps |

### Ticket stubs

<a id="todo-decision-log-prd"></a>
#### Define decision-log handoff to write-a-prd
**Question:** When map To Do is empty, how does `write-a-prd` ingest `{MAP-SLUG}-GM-xx` — keep prefixes, renumber to `GM-001`, or hybrid? What is the minimum Completed + log state required?

<a id="todo-grill-me-integration"></a>
#### Update grill-me to append scoped GM-xx to map decision log
**Question:** What changes does `grill-me` need (session startup, session complete) to write `{MAP-SLUG}-GM-xx` when invoked from a `wayfinder:grilling` ticket?

<a id="todo-research-workflow"></a>
#### Specify research ticket workflow (skill or Task pattern)
**Question:** New `research` skill vs `Task` subagent convention — resolution comment shape, branch naming, map updates?

<a id="todo-github-tracker"></a>
#### GitHub tracker setup for skills repo maps
**Question:** Labels, issue templates, and `gh` commands for `wayfinder:map`, `:todo`, `:decision-log` on `KroniK907/skills`.

<a id="todo-subfeature-example"></a>
#### Subfeature map worked example in REFERENCE or EXAMPLES
**Question:** Document one parent/child map pair (e.g. hypothetical greenfield app) showing Subfeatures + integration ticket.

<a id="todo-reconcile"></a>
#### Reconcile / post-implementation map update skill or mode
**Question:** After `agent-queue` slices ship, what updates Completed / decision log / new To Do — separate skill or wayfinder mode?

<a id="todo-afk-contract"></a>
#### Cloud AFK automation contract for wayfinder:afk tickets
**Question:** What must an AFK ticket body contain so Cursor cloud automation can claim, implement, and comment without closing?

## Completed

- [Initial wayfinder SKILL.md and REFERENCE.md](../SKILL.md) — Charted opinionated skill: ideation interview, To Do/Completed, subfeatures, scoped decision log, ecosystem table; human-close rule.

## Not yet specified

- **Canvas dashboard** — read-only view of map frontier; defer until GitHub tracker works.
- **Local vs GitHub source of truth** when both exist — sync direction unclear.
- **feature-ideation entry** — always before wayfinder, or only for raw seeds?
- **Decision log file vs issue** — single file per map vs dedicated GitHub issue for log only.

## Out of scope

- Implementing cloud automation itself (assumed future).
- Matt Pocock `/setup-matt-pocock-skills` tracker compatibility.
- Linear/Jira trackers (GitHub-first for now).

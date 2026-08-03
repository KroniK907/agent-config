# Wayfinder-Ecosystem:Decision-Log

**Map:** [Wayfinder-Ecosystem:Map](Wayfinder-Ecosystem.Map.md)  
**Prefix:** `WF-ECO-GM-`

Append-only log. `grill-me` and ticket resolutions add rows here. `write-a-prd` consolidates when map To Do is empty.

---

**WF-ECO-GM-001** — Wayfinder maps use `{FeatureName}:Map` titles, **To Do** and **Completed** sections (not Matt's "Decisions so far" index-only pattern). Agents post ticket resolutions; **humans close** tickets and move rows to Completed.

**WF-ECO-GM-002** — Charting starts with a **breadth-first ideation interview** across the five grill-me coverage zones. Answers may be concrete, `unknown`, or `needs research`; ideation collects gaps without resolving them.

**WF-ECO-GM-003** — Decision IDs use map-scoped prefixes `{MAP-SLUG}-GM-NNN` (this map: `WF-ECO-GM-`). Subfeature maps extend the slug (e.g. `CMD-PAL-SEARCH-GM-001`). Consolidation into PRD may renumber — TBD in open To Do ticket.

**WF-ECO-GM-004** — To Do tickets carry labels for **type** (research, prototype, grilling, task) and **mode** (HITL, AFK). AFK tickets are intended for cloud automation pickup; human still closes the issue.

**WF-ECO-GM-005** — Large features may spawn **subfeature maps** linked under parent **Subfeatures**; parent owns cross-cutting integration grilling tickets.

**WF-ECO-GM-006** — Local map files use `{FeatureName}.Map.md` and `{FeatureName}.Decision-Log.md` (dot separator) because `:` is invalid in paths on Windows. GitHub issue titles keep the `{FeatureName}:Map` display name.

*(Bootstrap entries from initial chart — refine via frontier tickets.)*

---
name: wayfinder
description: Overview of the wayfinder map lifecycle - map, grill, build, release - and which skill runs each step.
disable-model-invocation: true
agent-config-sync: true
---

# Wayfinder

A feature lives in one GitHub **map** issue ([MAP.md](MAP.md)): the plan, its decisions and its build order. Four skills move it along:

| Step | Skill | Ends with |
|------|-------|-----------|
| Shape | `map` | Research done and a map issue filed (Phase `shaping`) |
| Decide | `grill` | Every open question settled as a GM comment, **Slices** written (Phase `shaped`) |
| Build | `build` | A stack of reviewed PRs on `dev` and an HTML report (Phase `built`) |
| Ship | `release` | `dev` merged to `main`, tagged, back-merged; finished maps closed |

Helpers along the way: `research` (findings comment on the map), `design-modules` (interface shapes before a build), `babysit` and `pr` (each PR), `subagent` (model tiers), `hatch` (serve mocks and reports).

Small changes skip the map: build them directly, open a PR and babysit it.

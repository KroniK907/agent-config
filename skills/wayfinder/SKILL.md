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
| Build (Sol) | `sol-build` | The same for one or more maps in one stack: each builder runs its own sol-loop, and Sol reviews the seams |
| Hand off | `release-steps` | The build's human steps recorded on the release issue (Phase `built`) |
| Ship | `release` | `dev` merged to `main`, tagged, back-merged; human steps done; finished maps closed (Phase `released`) |

Helpers along the way: `research` (findings comment on the map), `design-modules` (interface shapes before a build), `babysit` and `pr` (each PR), `subagent` (model tiers), `hatch` (serve mocks and reports).

The phases are defined in [MAP.md](MAP.md). Human steps that prod needs (infra, secrets, data, DNS, settings) live on the repo's one open release issue from the moment work lands on `dev` until the release that ships it.

Small changes skip the map: build them directly and babysit them. `babysit` opens the PR once Sol approves.

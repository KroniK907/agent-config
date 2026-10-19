---
name: map
description: Create a feature map issue ({FeatureName}:Map), or post a research or cross-map note on one. Use when the user asks for a map issue, wants to shape or explore a new feature, or names a map with a note to add.
agent-config-sync: true
---

# Map

The map issue format lives in the wayfinder folder's `MAP.md`. Resolve this skill's folder with `readlink -f`; `MAP.md` is in its parent. Read it before writing anything.

Two forms:
- `/map <Name>: <paragraph>` creates a map.
- `/map #N <note>` posts a note on an existing map: research findings (follow the `research` skill) or a cross-map note that names the decisions it affects on both maps.

## Create a map

### 1. Research

Read the code the feature touches, related open issues and maps, and any sources the user gave: docs, files, drives, other repos. When the user asks how others solve it, or the feature brings in an outside API or library, research that too.

**Done when:** you can say what exists today, what the feature changes, and which issues and maps it touches.

### 2. Draft the body

Fill in the template from `MAP.md`: Target outcome, What's wrong today, Shape, and Open questions grouped by area. Each question is a decision the build needs. Leave **Slices** and **Decisions** empty. Pick a short uppercase slug.

**Done when:** every section except Slices and Decisions has content, and every open question can be settled with a single decision.

### 3. Explore, when the shape is not settled

The user may say grilling is too detailed, or describe a long-term vision without a first step. In that case, post a comment with up to about eight numbered shaping questions about edges, phases and naming, each with your lean. Fold the answers into the draft. Repeat until the user is happy with the shape.

**Done when:** the user accepts the shape, or the questions were not needed.

### 4. File it

Create the issue with `gh issue create --label wf:map --body-file <file>`. Post research worth keeping as a comment and link it from the body. Set **Phase** to `shaping` and **Next** to `/grill #N`. Comment on each related map and issue with one line linking back.

**Done when:** the map URL is in your reply, along with the related issues you linked.

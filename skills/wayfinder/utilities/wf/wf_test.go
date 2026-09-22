package main

// Run: go test ./wf.go ./wf_test.go   (from this directory; no go.mod by design)

import (
	"strings"
	"testing"
)

const testMap = `# Demo:Map

**Phase:** deciding
**Map slug:** ` + "`DEMO`" + `

## Target outcome

Ship it.

## Notes

- note

## Subfeatures

## To Do

| Ticket | Type | Mode | Assignee | Blocked by |
|--------|------|------|----------|------------|
| [Grill: alpha](https://github.com/o/r/issues/12) | grilling | HITL | unclaimed | - |
| [Research: beta](https://github.com/o/r/issues/13) | research | HITL | unclaimed | - |

## Implementing

| Ticket | Bundle | Mode | Status | Blocked by |
|--------|--------|------|--------|------------|
| [Task: gamma](https://github.com/o/r/issues/20) | #19 | AFK | ready | - |

## Completed

- [Map discovery](https://github.com/o/r/issues/2) - materialized 2 tickets

## Not yet specified

- fog one
- fog two

## Out of scope

## Decision coverage

| GM ID | Status | Linked issue |
|-------|--------|--------------|
| DEMO-GM-001 | open | - |
| DEMO-GM-002 | scoped | #19 |
`

func TestLogRowsAmendAndGlobal(t *testing.T) {
	is := &issue{Body: "# Demo:Decision-Log\n\n---\n\n**DEMO-GM-001** - first. `[global]`\n(from x)\n\n**DEMO-GM-002** - second\n"}
	is.Comments = append(is.Comments, struct {
		Body string `json:"body"`
	}{Body: "**DEMO-GM-002** - second, amended\n\n**DEMO-GM-003** - third"})

	order, byID := logRows(is)
	if strings.Join(order, ",") != "DEMO-GM-001,DEMO-GM-002,DEMO-GM-003" {
		t.Fatalf("order = %v", order)
	}
	if !strings.Contains(byID["DEMO-GM-002"].Text, "amended") {
		t.Errorf("comment should supersede body row: %q", byID["DEMO-GM-002"].Text)
	}
	if !isGlobal(byID["DEMO-GM-001"]) || isGlobal(byID["DEMO-GM-003"]) {
		t.Error("global tag detection wrong")
	}
	if !strings.Contains(byID["DEMO-GM-001"].Text, "(from x)") {
		t.Error("continuation line lost")
	}
	if id, _ := nextID(order, byID); id != "DEMO-GM-004" {
		t.Errorf("nextID = %s", id)
	}
}

func TestApplyMapOpsCompleteCoverageFog(t *testing.T) {
	out, _, err := applyMapOps(testMap, mapOps{
		complete: "12",
		gist:     "settled alpha",
		coverage: []string{"DEMO-GM-002=implemented", "DEMO-GM-003=open"},
		appends:  []string{"Out of scope=no mobile"},
		removes:  []string{"Not yet specified=fog one"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "issues/12) | grilling") {
		t.Error("To Do row not removed")
	}
	if !strings.Contains(out, "- [Grill: alpha](https://github.com/o/r/issues/12) - settled alpha") {
		t.Error("Completed gist missing")
	}
	if !strings.Contains(out, "| DEMO-GM-002 | implemented | #19 |") {
		t.Error("coverage update should keep existing link")
	}
	if !strings.HasSuffix(strings.TrimSpace(out), "| DEMO-GM-003 | open | - |") {
		t.Error("new coverage row should land at end of table")
	}
	if strings.Contains(out, "fog one") || !strings.Contains(out, "fog two") {
		t.Error("fog removal wrong")
	}
	if !strings.Contains(out, "## Out of scope\n\n- no mobile\n") {
		t.Errorf("append to empty section wrong:\n%s", out)
	}
	if err := validateMap(out); err != nil {
		t.Errorf("edited map invalid: %v", err)
	}
}

func TestApplyMapOpsImplementingAndErrors(t *testing.T) {
	out, _, err := applyMapOps(testMap, mapOps{complete: "20", gist: "shipped"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "- [Task: gamma](https://github.com/o/r/issues/20) - shipped") {
		t.Error("Implementing row should move to Completed")
	}
	if _, _, err := applyMapOps(testMap, mapOps{complete: "99", gist: "x"}); err == nil {
		t.Error("missing ticket should error")
	}
	if _, _, err := applyMapOps(testMap, mapOps{removes: []string{"Not yet specified=fog"}}); err == nil {
		t.Error("ambiguous remove should error")
	}
}

func TestValidateMap(t *testing.T) {
	if err := validateMap(testMap); err != nil {
		t.Fatalf("valid map rejected: %v", err)
	}
	collapsed := strings.ReplaceAll(testMap, "\n", " ")
	if validateMap(collapsed) == nil {
		t.Error("collapsed body accepted")
	}
	if validateMap(testMap+"\n## Extra\n") == nil {
		t.Error("coverage not last accepted")
	}
	if validateMap(strings.Replace(testMap, "Ship it.", "Ship it \u2014 now.", 1)) == nil {
		t.Error("em dash accepted")
	}
}

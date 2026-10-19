package main

// Run: go test ./wf.go ./wf_test.go   (from this directory; no go.mod by design)

import (
	"strings"
	"testing"
)

const testMap = `# Demo:Map

**Phase:** shaping. Next: /grill #N
**Map slug:** ` + "`DEMO`" + `

## Target outcome

Ship it.

## Shape

` + "```" + `
request ──► queue ├─ worker
` + "```" + `

## Slices

1. **Core.** ([DEMO-GM-001](https://github.com/o/r/issues/9#issuecomment-1))

## Decisions

| ID | Decision |
|---|---|
| DEMO-GM-001 | Core first |
`

func TestSectionBounds(t *testing.T) {
	lines := strings.Split(testMap, "\n")
	s, e, ok := sectionBounds(lines, "Slices")
	if !ok {
		t.Fatal("Slices not found")
	}
	got := strings.TrimSpace(strings.Join(lines[s:e], "\n"))
	if !strings.HasPrefix(got, "## Slices") || strings.Contains(got, "## Decisions") {
		t.Fatalf("bad bounds: %q", got)
	}
	if _, _, ok := sectionBounds(lines, "Missing"); ok {
		t.Fatal("found a missing section")
	}
}

func TestValidateBody(t *testing.T) {
	if err := validateBody(testMap); err != nil {
		t.Fatalf("valid map rejected (box drawing must pass): %v", err)
	}
	noDecisions := strings.Replace(testMap, "## Decisions", "## Log", 1)
	if err := validateBody(noDecisions); err == nil {
		t.Fatal("map without ## Decisions accepted")
	}
	if err := validateBody("plain comment, not a map"); err != nil {
		t.Fatalf("non-map rejected: %v", err)
	}
	if err := validateBody(testMap + "\u00c3©"); err == nil {
		t.Fatal("mojibake accepted")
	}
}

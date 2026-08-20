package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// minimalCatalog returns a small catalog with one rule, a skill group, and a leaf skill.
func minimalCatalog() *catalogFile {
	c := &catalogFile{}
	c.Catalog.Version = "test-1.0"
	c.Rules = map[string]catalogEntry{
		"unslop": {Path: "rules/unslop.mdc", Label: "Unslop"},
	}
	c.Skills = map[string]catalogEntry{
		"wayfinder": {
			Path:  "skills/wayfinder",
			Label: "Wayfinder",
		},
		"wayfinder-code-review": {
			Path:  "skills/wayfinder/actions/code-review",
			Label: "Code review",
		},
		"commit": {
			Path:  "skills/commit",
			Label: "Commit",
		},
	}
	return c
}

func writeTeamFixture(t *testing.T, cat *catalogFile) string {
	t.Helper()
	if cat == nil {
		cat = minimalCatalog()
	}
	team := t.TempDir()
	if err := os.MkdirAll(filepath.Join(team, ".cursor", "examples"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(team, ".cursor", "examples", "README.md"), []byte("# examples\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	raw, err := json.MarshalIndent(cat, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(team, "catalog.json"), raw, 0o644); err != nil {
		t.Fatal(err)
	}
	for _, e := range cat.Rules {
		p := filepath.Join(team, filepath.FromSlash(e.Path))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		body := "---\nagent-config-sync: true\n---\n\n# rule\n"
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	for _, e := range cat.Skills {
		p := filepath.Join(team, filepath.FromSlash(e.Path), "SKILL.md")
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		body := "---\nagent-config-sync: true\nname: test\n---\n\n# skill\n"
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return team
}

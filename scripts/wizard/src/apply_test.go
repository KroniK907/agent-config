package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func enabledItems(paths ...string) []treeItem {
	set := map[string]bool{}
	for _, p := range paths {
		set[p] = true
	}
	var items []treeItem
	for p := range set {
		kind := "skill"
		if strings.HasPrefix(p, "rules/") {
			kind = "rule"
		}
		items = append(items, treeItem{
			Path:    p,
			Kind:    kind,
			Enabled: true,
		})
	}
	return items
}

func TestApplyEnabled_copiesEnabledEntries(t *testing.T) {
	cat := minimalCatalog()
	team := writeTeamFixture(t, cat)
	project := t.TempDir()

	m := newManifest(team, project, cat.Catalog.Version)
	m.LastCatalogPaths = catalogPaths(cat)
	items := enabledItems("rules/unslop.mdc", "skills/commit")

	res := applyEnabled(team, project, items, m)

	if len(res.Errors) > 0 {
		t.Fatalf("apply errors: %v", res.Errors)
	}
	destRule := filepath.Join(project, ".cursor", "rules", "unslop.mdc")
	if _, err := os.Stat(destRule); err != nil {
		t.Fatalf("rule not copied: %v", err)
	}
	destSkill := filepath.Join(project, ".cursor", "skills", "commit", "SKILL.md")
	if _, err := os.Stat(destSkill); err != nil {
		t.Fatalf("skill not copied: %v", err)
	}
	sync, ok := readAgentConfigSync(destRule, "rule")
	if !ok || sync != "true" {
		t.Errorf("stamped sync on rule = %q, ok=%v", sync, ok)
	}
	if !contains(m.Rules, "rules/unslop.mdc") {
		t.Errorf("manifest Rules = %v", m.Rules)
	}
	if !contains(m.Skills, "skills/commit") {
		t.Errorf("manifest Skills = %v", m.Skills)
	}
}

func TestApplyEnabled_removesDeselectedManaged(t *testing.T) {
	cat := minimalCatalog()
	team := writeTeamFixture(t, cat)
	project := t.TempDir()

	// Pre-copy a managed rule, then apply with it disabled.
	dest := filepath.Join(project, ".cursor", "rules", "unslop.mdc")
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		t.Fatal(err)
	}
	body := "---\nagent-config-sync: true\n---\n\n# managed\n"
	if err := os.WriteFile(dest, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	m := &manifest{
		Rules:       []string{"rules/unslop.mdc"},
		LastApplied: map[string]string{"rules/unslop.mdc": "2026-01-01T00:00:00Z"},
		LastCatalogPaths: catalogPaths(cat),
	}
	items := buildTree(cat, m, project)
	for i := range items {
		if items[i].Path == "rules/unslop.mdc" {
			items[i].Enabled = false
		}
	}

	res := applyEnabled(team, project, items, m)
	if len(res.Removed) != 1 || res.Removed[0] != "rules/unslop.mdc" {
		t.Errorf("Removed = %v, want [rules/unslop.mdc]", res.Removed)
	}
	if _, err := os.Stat(dest); !os.IsNotExist(err) {
		t.Error("deselected managed rule should be removed")
	}
}

func TestApplyEnabled_skipsDetachedOnDeselect(t *testing.T) {
	cat := minimalCatalog()
	team := writeTeamFixture(t, cat)
	project := t.TempDir()

	dest := filepath.Join(project, ".cursor", "rules", "unslop.mdc")
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		t.Fatal(err)
	}
	body := "---\nagent-config-sync: false\n---\n\n# detached\n"
	if err := os.WriteFile(dest, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	m := &manifest{
		Rules:            []string{"rules/unslop.mdc"},
		LastApplied:      map[string]string{"rules/unslop.mdc": "2026-01-01T00:00:00Z"},
		LastCatalogPaths: catalogPaths(cat),
	}
	items := buildTree(cat, m, project)
	for i := range items {
		if items[i].Path == "rules/unslop.mdc" {
			items[i].Enabled = false
		}
	}

	res := applyEnabled(team, project, items, m)
	if len(res.Removed) != 0 {
		t.Errorf("detached file should not be removed, Removed = %v", res.Removed)
	}
	if len(res.Skipped) == 0 || !strings.Contains(res.Skipped[0], "detached") {
		t.Errorf("Skipped = %v, want detached skip message", res.Skipped)
	}
	if _, err := os.Stat(dest); err != nil {
		t.Fatalf("detached file should remain: %v", err)
	}
}

func TestApplyEnabled_skipsProjectOverride(t *testing.T) {
	cat := minimalCatalog()
	team := writeTeamFixture(t, cat)
	project := t.TempDir()

	dest := filepath.Join(project, ".cursor", "skills", "commit", "SKILL.md")
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		t.Fatal(err)
	}
	body := "---\nagent-config-sync: false\n---\n\n# local override\n"
	if err := os.WriteFile(dest, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	m := newManifest(team, project, cat.Catalog.Version)
	m.LastCatalogPaths = catalogPaths(cat)
	items := buildTree(cat, m, project)
	for i := range items {
		if items[i].Path == "skills/commit" {
			items[i].Enabled = true
		}
	}

	res := applyEnabled(team, project, items, m)
	foundSkip := false
	for _, s := range res.Skipped {
		if strings.Contains(s, "skills/commit") && strings.Contains(s, "project override") {
			foundSkip = true
		}
	}
	if !foundSkip {
		t.Errorf("expected project override skip, Skipped = %v", res.Skipped)
	}
	data, err := os.ReadFile(dest)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "local override") {
		t.Error("project file content should be unchanged")
	}
}

func TestPreviouslyManaged_includesLastApplied(t *testing.T) {
	m := &manifest{
		Skills:      []string{"skills/a"},
		LastApplied: map[string]string{"rules/b.mdc": "t"},
	}
	pm := previouslyManaged(m)
	if !pm["skills/a"] || !pm["rules/b.mdc"] {
		t.Errorf("previouslyManaged = %v", pm)
	}
}

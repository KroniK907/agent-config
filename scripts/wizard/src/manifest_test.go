package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadManifest_missingReturnsNil(t *testing.T) {
	dir := t.TempDir()
	m, err := loadManifest(dir)
	if err != nil {
		t.Fatalf("loadManifest: %v", err)
	}
	if m != nil {
		t.Fatalf("expected nil manifest, got %+v", m)
	}
}

func TestSaveAndLoadManifest_roundtrip(t *testing.T) {
	dir := t.TempDir()
	orig := &manifest{
		Source: sourceMeta{
			Repo:     sourceRepo,
			Ref:      "1.0.0",
			TeamPath: "/team",
		},
		ProjectPath:      dir,
		Skills:           []string{"skills/commit"},
		Rules:            []string{"rules/unslop.mdc"},
		LastApplied:      map[string]string{"skills/commit": "2026-01-01T00:00:00Z"},
		LastCatalogPaths: []string{"skills/commit", "rules/unslop.mdc"},
	}
	if err := saveManifest(dir, orig); err != nil {
		t.Fatalf("saveManifest: %v", err)
	}

	got, err := loadManifest(dir)
	if err != nil {
		t.Fatalf("loadManifest: %v", err)
	}
	if got.ProjectPath != orig.ProjectPath {
		t.Errorf("ProjectPath = %q, want %q", got.ProjectPath, orig.ProjectPath)
	}
	if len(got.Skills) != 1 || got.Skills[0] != "skills/commit" {
		t.Errorf("Skills = %v, want [skills/commit]", got.Skills)
	}
	if len(got.Rules) != 1 || got.Rules[0] != "rules/unslop.mdc" {
		t.Errorf("Rules = %v, want [rules/unslop.mdc]", got.Rules)
	}
	if got.LastApplied["skills/commit"] != "2026-01-01T00:00:00Z" {
		t.Errorf("LastApplied = %v", got.LastApplied)
	}

	path := manifestPath(dir)
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("manifest file missing: %v", err)
	}
}

func TestNewManifest_defaults(t *testing.T) {
	m := newManifest("/team", "/project", "2.0.0")
	if m.Source.Repo != sourceRepo {
		t.Errorf("Repo = %q, want %q", m.Source.Repo, sourceRepo)
	}
	if m.Source.Ref != "2.0.0" {
		t.Errorf("Ref = %q, want 2.0.0", m.Source.Ref)
	}
	if m.ProjectPath != "/project" {
		t.Errorf("ProjectPath = %q", m.ProjectPath)
	}
	if m.Skills != nil || m.Rules != nil {
		t.Errorf("expected nil Skills/Rules on new manifest")
	}
	if m.LastApplied == nil {
		t.Fatal("LastApplied should be initialized")
	}
}

func TestEnabledSet(t *testing.T) {
	m := &manifest{
		Rules:  []string{"rules/a.mdc"},
		Skills: []string{"skills/b", "skills/c"},
	}
	set := enabledSet(m)
	for _, p := range []string{"rules/a.mdc", "skills/b", "skills/c"} {
		if !set[p] {
			t.Errorf("expected %q enabled", p)
		}
	}
	if set["skills/missing"] {
		t.Error("unexpected path in enabled set")
	}
	if len(enabledSet(nil)) != 0 {
		t.Error("nil manifest should yield empty set")
	}
}

func TestManifestPath(t *testing.T) {
	got := manifestPath("/proj")
	want := filepath.Join("/proj", ".cursor", manifestName)
	if got != want {
		t.Errorf("manifestPath = %q, want %q", got, want)
	}
}

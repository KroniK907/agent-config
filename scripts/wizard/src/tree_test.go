package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuildTree_envDetailsThenRulesBeforeSkills(t *testing.T) {
	cat := minimalCatalog()
	project := t.TempDir()
	items := buildTree(cat, nil, project)

	if len(items) == 0 {
		t.Fatal("expected tree items")
	}
	if items[0].Kind != "env" {
		t.Errorf("first item kind = %q, want env", items[0].Kind)
	}
	var sawRule, sawSkill bool
	for _, it := range items[1:] {
		if it.Kind == "rule" {
			sawRule = true
		}
		if it.Kind == "skill" && !it.IsGroup {
			if !sawRule {
				t.Error("expected rules before leaf skills")
			}
			sawSkill = true
		}
	}
	if !sawRule || !sawSkill {
		t.Fatalf("expected rules and skills in tree, got rule=%v skill=%v", sawRule, sawSkill)
	}
}

func TestBuildTree_groupHasDescendants(t *testing.T) {
	cat := minimalCatalog()
	project := t.TempDir()
	items := buildTree(cat, nil, project)

	var group *treeItem
	for i := range items {
		if items[i].IsGroup && items[i].Path == "skills/wayfinder" {
			group = &items[i]
			break
		}
	}
	if group == nil {
		t.Fatal("wayfinder group not found")
	}
	if len(group.DescendantPaths) < 2 {
		t.Fatalf("DescendantPaths = %v, want hub + child paths", group.DescendantPaths)
	}
	hasHub := false
	hasChild := false
	for _, p := range group.DescendantPaths {
		if p == "skills/wayfinder" {
			hasHub = true
		}
		if p == "skills/wayfinder/actions/code-review" {
			hasChild = true
		}
	}
	if !hasHub || !hasChild {
		t.Errorf("DescendantPaths = %v, want hub and nested child", group.DescendantPaths)
	}
}

func TestBuildTree_newEntryDefaultsOff(t *testing.T) {
	cat := minimalCatalog()
	project := t.TempDir()
	prev := &manifest{
		LastCatalogPaths: []string{"rules/unslop.mdc", "skills/wayfinder", "skills/wayfinder/actions/code-review"},
		Skills:           []string{"skills/commit"},
	}
	items := buildTree(cat, prev, project)

	var commit *treeItem
	for i := range items {
		if items[i].Path == "skills/commit" {
			commit = &items[i]
			break
		}
	}
	if commit == nil {
		t.Fatal("skills/commit leaf not found")
	}
	if !commit.IsNew {
		t.Error("expected IsNew for path absent from LastCatalogPaths")
	}
	if commit.Enabled {
		t.Error("new entries should default to disabled")
	}
}

func TestBuildTree_restoresPreviouslyEnabled(t *testing.T) {
	cat := minimalCatalog()
	project := t.TempDir()
	prev := &manifest{
		LastCatalogPaths: catalogPaths(cat),
		Skills:           []string{"skills/commit"},
	}
	items := buildTree(cat, prev, project)

	for _, it := range items {
		if it.Path == "skills/commit" {
			if !it.Enabled {
				t.Error("previously enabled skill should stay enabled")
			}
			if it.IsNew {
				t.Error("known catalog path should not be NEW")
			}
			return
		}
	}
	t.Fatal("skills/commit not found")
}

func TestProjectOverride_detachedSyncFalse(t *testing.T) {
	project := t.TempDir()
	dest := filepath.Join(project, ".cursor", "rules", "unslop.mdc")
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		t.Fatal(err)
	}
	body := "---\nagent-config-sync: false\n---\n\n# detached\n"
	if err := os.WriteFile(dest, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	if !projectOverride(project, "rules/unslop.mdc", "rule", nil) {
		t.Error("expected project override when sync is false")
	}
}

func TestProjectOverride_managedCopyNotOverride(t *testing.T) {
	project := t.TempDir()
	dest := filepath.Join(project, ".cursor", "skills", "commit", "SKILL.md")
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		t.Fatal(err)
	}
	body := "---\nagent-config-sync: true\n---\n\n# managed\n"
	if err := os.WriteFile(dest, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	prev := &manifest{Skills: []string{"skills/commit"}}
	if projectOverride(project, "skills/commit", "skill", prev) {
		t.Error("managed copy with sync true should not be override")
	}
}

func TestSetGroupEnabled_cascadesToChildren(t *testing.T) {
	cat := minimalCatalog()
	project := t.TempDir()
	items := buildTree(cat, nil, project)

	var groupIdx int
	for i, it := range items {
		if it.IsGroup && it.Path == "skills/wayfinder" {
			groupIdx = i
			break
		}
	}

	m := &model{items: items, projectRoot: project}
	m.setGroupEnabled(&m.items[groupIdx], true)

	leaf := m.leafByPath("skills/wayfinder/actions/code-review")
	if leaf == nil {
		t.Fatal("nested leaf not found")
	}
	if !leaf.Enabled {
		t.Error("nested child should be enabled after group select")
	}
}

func TestGroupCheck_partialSelection(t *testing.T) {
	cat := &catalogFile{}
	cat.Catalog.Version = "test-1.0"
	cat.Skills = map[string]catalogEntry{
		"wf-cr": {Path: "skills/wayfinder/actions/code-review", Label: "Code review"},
		"wf-wc": {Path: "skills/wayfinder/actions/write-code", Label: "Write code"},
	}
	project := t.TempDir()
	items := buildTree(cat, nil, project)

	var groupIdx int
	found := false
	for i, it := range items {
		if it.IsGroup && it.Path == "skills/wayfinder/actions" {
			groupIdx = i
			found = true
			break
		}
	}
	if !found {
		t.Fatal("actions group not found")
	}

	m := &model{items: items, projectRoot: project}
	if leaf := m.leafByPath("skills/wayfinder/actions/code-review"); leaf != nil {
		leaf.Enabled = true
	}
	check := m.groupCheck(groupIdx)
	if check != "~" {
		t.Errorf("groupCheck = %q, want ~ for partial selection", check)
	}
}

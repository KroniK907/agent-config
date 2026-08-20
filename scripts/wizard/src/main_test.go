package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

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

func TestParseFrontmatter_readsSyncField(t *testing.T) {
	data := []byte("---\nname: foo\nagent-config-sync: true\n---\n\nbody")
	fm := parseFrontmatter(data)
	if fm["agent-config-sync"] != "true" {
		t.Errorf("agent-config-sync = %q, want true", fm["agent-config-sync"])
	}
	if fm["name"] != "foo" {
		t.Errorf("name = %q", fm["name"])
	}
}

func TestParseFrontmatter_noFrontmatter(t *testing.T) {
	if parseFrontmatter([]byte("no frontmatter")) != nil {
		t.Error("expected nil for missing frontmatter")
	}
}

func TestUpsertFrontmatterField_updatesExisting(t *testing.T) {
	orig := []byte("---\nagent-config-sync: true\nname: x\n---\n\ncontent")
	got := upsertFrontmatterField(orig, "agent-config-sync", "false")
	if !bytes.Contains(got, []byte("agent-config-sync: false")) {
		t.Fatalf("expected updated sync field, got:\n%s", got)
	}
	if bytes.Contains(got, []byte("agent-config-sync: true")) {
		t.Error("old sync value should be replaced")
	}
}

func TestUpsertFrontmatterField_insertsWhenMissing(t *testing.T) {
	orig := []byte("---\nname: x\n---\n\ncontent")
	got := upsertFrontmatterField(orig, "agent-config-sync", "true")
	if !bytes.Contains(got, []byte("agent-config-sync: true")) {
		t.Fatalf("expected inserted field, got:\n%s", got)
	}
}

func TestUpsertFrontmatterField_addsFrontmatterWhenAbsent(t *testing.T) {
	got := upsertFrontmatterField([]byte("plain body"), "agent-config-sync", "true")
	if !bytes.HasPrefix(got, []byte("---\n")) {
		t.Fatalf("expected new frontmatter block, got:\n%s", got)
	}
	if !bytes.Contains(got, []byte("agent-config-sync: true")) {
		t.Error("expected sync field in new frontmatter")
	}
}

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

	dest := filepath.Join(project, ".cursor", "rules", "unslop.mdc")
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		t.Fatal(err)
	}
	body := "---\nagent-config-sync: true\n---\n\n# managed\n"
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

func TestRenderEnvironmentRule_includesScopeAndRuntime(t *testing.T) {
	content := string(renderEnvironmentRule(envProbe{
		OSLine:     "Windows (Microsoft Windows [Version 10.0.26200])",
		ShellLine:  "PowerShell 7.5.0 (default for terminal commands)",
		ScratchDir: "C:/Users/me/.cursor/scratch",
		ToolLines: []string{
			"`go` installed (go version go1.22.0 windows/amd64)",
			"`git` installed (git version 2.43.0.windows.1)",
		},
		AbsentLines: []string{
			"Do not assume `node`, `npm`, or `npx`.",
			"Do not assume `python` or `pip`.",
		},
	}))

	for _, want := range []string{
		"alwaysApply: true",
		"generated-by: agent-config-wizard",
		"**Scope:**",
		"context bloat",
		"**OS:** Windows",
		"**Shell:** PowerShell",
		"Agent scratch directory",
		"`go` installed",
		"Absent runtimes",
		"Do not assume `node`",
	} {
		if !strings.Contains(content, want) {
			t.Errorf("render missing %q\n%s", want, content)
		}
	}
}

func TestRefreshEnvironmentRule_writesGeneratedFile(t *testing.T) {
	project := t.TempDir()
	if err := refreshEnvironmentRule(project); err != nil {
		t.Fatalf("refreshEnvironmentRule: %v", err)
	}
	path := environmentRulePath(project)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read rule: %v", err)
	}
	if !strings.Contains(string(data), "generated-by: agent-config-wizard") {
		t.Fatalf("expected generated marker in %q", string(data))
	}
	if !isGeneratedEnvironmentRule(path) {
		t.Fatal("isGeneratedEnvironmentRule should be true")
	}
}

func TestRemoveEnvironmentRule_onlyGenerated(t *testing.T) {
	project := t.TempDir()
	path := environmentRulePath(project)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	custom := "---\nagent-config-sync: false\n---\n\n# custom\n"
	if err := os.WriteFile(path, []byte(custom), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := removeEnvironmentRule(project); err != nil {
		t.Fatalf("removeEnvironmentRule: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal("custom rule should remain when not generated")
	}

	if err := refreshEnvironmentRule(project); err != nil {
		t.Fatal(err)
	}
	if err := removeEnvironmentRule(project); err != nil {
		t.Fatalf("remove generated: %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("generated rule should be removed")
	}
}

func TestApplyEnabled_refreshesEnvRuleWhenEnabled(t *testing.T) {
	cat := minimalCatalog()
	team := writeTeamFixture(t, cat)
	project := t.TempDir()

	m := newManifest(team, project, cat.Catalog.Version)
	m.EnvDetails = true
	m.LastCatalogPaths = catalogPaths(cat)

	res := applyEnabled(team, project, nil, m)
	if len(res.Errors) > 0 {
		t.Fatalf("apply errors: %v", res.Errors)
	}
	if !contains(res.Copied, envRuleRelPath) {
		t.Fatalf("expected env rule in copied, got %v", res.Copied)
	}
	if !isGeneratedEnvironmentRule(environmentRulePath(project)) {
		t.Fatal("environment rule not written")
	}
}

func TestApplyEnabled_removesEnvRuleWhenDisabled(t *testing.T) {
	cat := minimalCatalog()
	team := writeTeamFixture(t, cat)
	project := t.TempDir()

	m := newManifest(team, project, cat.Catalog.Version)
	m.EnvDetails = true
	if err := refreshEnvironmentRule(project); err != nil {
		t.Fatal(err)
	}
	m.EnvDetails = false

	res := applyEnabled(team, project, nil, m)
	if len(res.Errors) > 0 {
		t.Fatalf("apply errors: %v", res.Errors)
	}
	if !contains(res.Removed, envRuleRelPath) {
		t.Fatalf("expected env rule removed, got removed=%v", res.Removed)
	}
}

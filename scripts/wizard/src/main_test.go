package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
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

func TestEnvRuleProjectOverride_generatedNotOverride(t *testing.T) {
	project := t.TempDir()
	if err := refreshEnvironmentRule(project); err != nil {
		t.Fatal(err)
	}
	if envRuleProjectOverride(project) {
		t.Fatal("generated env rule should not be project override")
	}
	if projectOverride(project, envRuleRelPath, "env", nil) {
		t.Fatal("projectOverride should be false for generated env rule")
	}
}

func TestEnvRuleProjectOverride_customIsOverride(t *testing.T) {
	project := t.TempDir()
	path := environmentRulePath(project)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	custom := "---\nagent-config-sync: false\n---\n\n# custom\n"
	if err := os.WriteFile(path, []byte(custom), 0o644); err != nil {
		t.Fatal(err)
	}
	if !envRuleProjectOverride(project) {
		t.Fatal("custom env rule should be project override")
	}
}

func TestApplyEnabled_envRuleNotMarkedOverrideAfterApply(t *testing.T) {
	cat := minimalCatalog()
	team := writeTeamFixture(t, cat)
	project := t.TempDir()
	m := newManifest(team, project, cat.Catalog.Version)
	m.EnvDetails = true
	items := buildTree(cat, m, project)
	for i := range items {
		if items[i].Kind == "env" {
			items[i].Enabled = true
		}
	}
	res := applyEnabled(team, project, items, m)
	if len(res.Errors) > 0 {
		t.Fatalf("apply errors: %v", res.Errors)
	}
	for _, it := range items {
		if it.Kind != "env" {
			continue
		}
		if projectOverride(project, it.Path, it.Kind, m) {
			t.Fatal("env item should not show project override after apply")
		}
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

func TestApplyStatusMessage(t *testing.T) {
	ok := applyResult{Copied: []string{"a"}, Removed: []string{"b"}, Skipped: []string{"c"}}
	if msg := applyStatusMessage(ok, nil); msg != "applied: 1 copied, 1 removed, 1 skipped" {
		t.Errorf("success message = %q", msg)
	}

	withErr := applyResult{
		Copied: []string{"a"},
		Errors: []string{"rules/x: copy failed"},
	}
	if msg := applyStatusMessage(withErr, nil); !strings.Contains(msg, "1 error") {
		t.Errorf("error message = %q", msg)
	}
	if msg := applyStatusMessage(withErr, os.ErrPermission); !strings.Contains(msg, "NOT saved") {
		t.Errorf("save failed message = %q", msg)
	}
}

func TestRenderErrors_showsVisibleLines(t *testing.T) {
	m := model{errors: []string{"first", "second", "third"}}
	out := m.renderErrors()
	for _, want := range []string{"ERRORS:", "! first", "! second", "! third"} {
		if !strings.Contains(out, want) {
			t.Errorf("renderErrors missing %q:\n%s", want, out)
		}
	}
}

func TestRenderErrors_truncatesLongList(t *testing.T) {
	var errs []string
	for i := 0; i < 8; i++ {
		errs = append(errs, "err")
	}
	m := model{errors: errs}
	out := m.renderErrors()
	if !strings.Contains(out, "... and 3 more") {
		t.Errorf("expected truncation notice, got:\n%s", out)
	}
}

func TestFooterLineCount_growsWithErrors(t *testing.T) {
	m := model{}
	if m.footerLineCount() != minFooterLines {
		t.Fatalf("empty footer = %d, want %d", m.footerLineCount(), minFooterLines)
	}
	m.errors = []string{"one", "two"}
	if m.footerLineCount() <= minFooterLines {
		t.Fatalf("footer with errors should grow, got %d", m.footerLineCount())
	}
}

func TestRefreshStateDump_listsApplyErrorsFirst(t *testing.T) {
	m := model{manifest: newManifest("/team", "/proj", "1")}
	res := applyResult{Errors: []string{"rules/x: boom"}}
	m.refreshStateDump(&res)
	if !strings.HasPrefix(m.stateDump, "Apply errors:") {
		t.Fatalf("expected errors header first, got:\n%s", m.stateDump)
	}
	if !strings.Contains(m.stateDump, "rules/x: boom") {
		t.Fatalf("expected error line in dump, got:\n%s", m.stateDump)
	}
}

func TestCursorMarker_spinnerWhileApplying(t *testing.T) {
	if got := cursorMarker(true, 0); got != "| " {
		t.Fatalf("frame 0 = %q, want %q", got, "| ")
	}
	if got := cursorMarker(true, 1); got != "/ " {
		t.Fatalf("frame 1 = %q, want %q", got, "/ ")
	}
	if got := cursorMarker(true, 4); got != "| " {
		t.Fatalf("frame 4 = %q, want %q", got, "| ")
	}
	if got := cursorMarker(false, 0); got != "> " {
		t.Fatalf("idle = %q, want %q", got, "> ")
	}
}

func TestFormatVersionLabel_firstLineAndSingleLine(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "wsl utf16 decoded",
			in:   "WSL version: 2.7.12.0\r\nKernel version: 6.18.33.2-2\r\n",
			want: "WSL version: 2.7.12.0",
		},
		{
			name: "gcloud multiline",
			in:   "Google Cloud SDK 564.0.0\nbeta 2026.04.03\nbq 2.1.31\n",
			want: "Google Cloud SDK 564.0.0",
		},
		{
			name: "carriage return before closing paren scenario",
			in:   "Google Cloud SDK 564.0.0\r",
			want: "Google Cloud SDK 564.0.0",
		},
		{
			name: "collapses internal whitespace",
			in:   "git version  2.43.0.windows.1",
			want: "git version 2.43.0.windows.1",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := formatVersionLabel(tc.in); got != tc.want {
				t.Fatalf("formatVersionLabel() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestNormalizeCommandOutput_decodesWSLUTF16(t *testing.T) {
	raw := []byte{
		'W', 0, 'S', 0, 'L', 0, ' ', 0, 'v', 0, 'e', 0, 'r', 0, 's', 0, 'i', 0, 'o', 0, 'n', 0, ':', 0, ' ', 0,
		'2', 0, '.', 0, '7', 0, '.', 0, '1', 0, '2', 0, '.', 0, '0', 0, '\r', 0, '\n', 0,
		'K', 0, 'e', 0, 'r', 0, 'n', 0, 'e', 0, 'l', 0,
	}
	got := formatVersionLabel(normalizeCommandOutput(raw))
	want := "WSL version: 2.7.12.0"
	if got != want {
		t.Fatalf("formatVersionLabel(normalizeCommandOutput()) = %q, want %q", got, want)
	}
}

func TestNormalizeCommandOutput_keepsUTF8(t *testing.T) {
	raw := []byte("git version 2.43.0.windows.1\n")
	got := normalizeCommandOutput(raw)
	if got != "git version 2.43.0.windows.1" {
		t.Fatalf("normalizeCommandOutput() = %q", got)
	}
}

func TestWrapText_breaksLongErrorLines(t *testing.T) {
	long := strings.Repeat("segment ", 20)
	wrapped := wrapText("  ! "+long, 40)
	for _, line := range strings.Split(wrapped, "\n") {
		if len(line) > 40 {
			t.Fatalf("line exceeds width %d: %q", len(line), line)
		}
	}
	if !strings.Contains(wrapped, "\n") {
		t.Fatal("expected wrapped output to span multiple lines")
	}
}

func TestToolProbeSpecs_uniqueNamesAndCoverage(t *testing.T) {
	specs := toolProbeSpecs()
	if len(specs) < 20 {
		t.Fatalf("expected expanded probe list, got %d specs", len(specs))
	}
	seen := map[string]bool{}
	for _, spec := range specs {
		if spec.name == "" {
			t.Fatal("tool spec with empty name")
		}
		if seen[spec.name] {
			t.Fatalf("duplicate probe name %q", spec.name)
		}
		seen[spec.name] = true
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
		"Generated by agent-config-wizard",
		"Manual edits are overwritten",
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

func TestApplyEnabled_envDetailsTreeItem_skipsTeamCopy(t *testing.T) {
	cat := minimalCatalog()
	team := writeTeamFixture(t, cat)
	project := t.TempDir()

	m := newManifest(team, project, cat.Catalog.Version)
	m.LastCatalogPaths = catalogPaths(cat)
	items := buildTree(cat, m, project)
	for i := range items {
		if items[i].Kind == "env" {
			items[i].Enabled = true
		}
	}
	m.EnvDetails = true

	res := applyEnabled(team, project, items, m)
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

func TestParseReleaseTagNames_extractsTags(t *testing.T) {
	data := []byte(`[{"tag_name":"v1.0.1"},{"tag_name":"v0.1.0"}]`)
	tags, err := parseReleaseTagNames(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) != 2 || tags[0] != "v1.0.1" {
		t.Fatalf("tags = %v", tags)
	}
}

func TestSemverTagForCatalogVersion_matchesVPrefix(t *testing.T) {
	tags := []string{"v1.0.1", "v0.1.0"}
	tag, ok := semverTagForCatalogVersion("1.0.1", tags)
	if !ok || tag != "v1.0.1" {
		t.Fatalf("tag = %q ok = %v", tag, ok)
	}
}

func TestDefaultCloudRef_prefersExisting(t *testing.T) {
	tags := []string{"v1.0.1", "v0.1.0"}
	got := defaultCloudRef("1.0.1", tags, "v0.1.0")
	if got != "v0.1.0" {
		t.Fatalf("ref = %q, want v0.1.0", got)
	}
}

func TestEnvironmentInstallURL_pinsRef(t *testing.T) {
	url := environmentInstallURL("v1.0.1")
	if !strings.Contains(url, "v1.0.1") || !strings.Contains(url, "bootstrap-agent.sh") {
		t.Fatalf("url = %q", url)
	}
}

func TestDiffCloudPaths_addedAndRemoved(t *testing.T) {
	diff := diffCloudPaths(
		[]string{"skills/a"},
		[]string{"rules/old.mdc"},
		[]string{"skills/b"},
		[]string{"rules/new.mdc"},
	)
	if len(diff.Added) != 2 || len(diff.Removed) != 2 {
		t.Fatalf("diff = %+v", diff)
	}
}

func TestValidateCloudPaths_rejectsUnknown(t *testing.T) {
	cat := minimalCatalog()
	errs := validateCloudPaths(cat, []string{"skills/missing"}, nil)
	if len(errs) != 1 {
		t.Fatalf("errs = %v", errs)
	}
}

func TestValidateCloudPaths_acceptsKnown(t *testing.T) {
	cat := minimalCatalog()
	errs := validateCloudPaths(cat, []string{"skills/commit"}, []string{"rules/unslop.mdc"})
	if len(errs) != 0 {
		t.Fatalf("errs = %v", errs)
	}
}

func TestWriteCloudConfig_writesBothFiles(t *testing.T) {
	project := t.TempDir()
	cm := newCloudManifest("v1.0.1")
	cm.Skills = []string{"skills/commit"}
	cm.Rules = []string{"rules/unslop.mdc"}
	if err := writeCloudConfig(project, cm); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(cloudManifestPath(project)); err != nil {
		t.Fatalf("manifest missing: %v", err)
	}
	if _, err := os.Stat(environmentJSONPath(project)); err != nil {
		t.Fatalf("environment.json missing: %v", err)
	}
	if _, err := os.Stat(cloudDockerfilePath(project)); err != nil {
		t.Fatalf("Dockerfile missing: %v", err)
	}
	raw, err := os.ReadFile(environmentJSONPath(project))
	if err != nil {
		t.Fatal(err)
	}
	var env environmentFile
	if err := json.Unmarshal(raw, &env); err != nil {
		t.Fatal(err)
	}
	if env.Build.Dockerfile != "Dockerfile" {
		t.Fatalf("dockerfile = %q", env.Build.Dockerfile)
	}
	if !strings.Contains(env.Install, "v1.0.1") || !strings.Contains(env.Install, "bootstrap-agent.sh") {
		t.Fatalf("install = %q", env.Install)
	}
	gi, err := os.ReadFile(filepath.Join(project, ".gitignore"))
	if err != nil {
		t.Fatalf(".gitignore missing: %v", err)
	}
	body := string(gi)
	for _, want := range []string{
		gitignoreBegin,
		".cursor/agent-config.local.json",
		"!.cursor/agent-manifest.json",
		"!.cursor/environment.json",
		"!.cursor/Dockerfile",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf(".gitignore missing %q:\n%s", want, body)
		}
	}
}

func TestMergeGitignoreBlock_createsAndUpdates(t *testing.T) {
	created := mergeGitignoreBlock("")
	if !strings.Contains(created, gitignoreBegin) || !strings.Contains(created, "!.cursor/agent-manifest.json") {
		t.Fatalf("create block:\n%s", created)
	}
	withExtra := "# custom\n*.log\n"
	merged := mergeGitignoreBlock(withExtra)
	if !strings.Contains(merged, "*.log") || !strings.Contains(merged, gitignoreBegin) {
		t.Fatalf("merge append:\n%s", merged)
	}
	updated := mergeGitignoreBlock(merged)
	if updated != merged {
		t.Fatal("expected idempotent merge")
	}
	replaced := mergeGitignoreBlock(strings.Replace(merged, "!.cursor/environment.json", "# old", 1))
	if strings.Contains(replaced, "# old") {
		t.Fatalf("expected block replaced, got:\n%s", replaced)
	}
	if !strings.Contains(replaced, "!.cursor/environment.json") {
		t.Fatal("expected fresh environment negation")
	}
}

func TestEnsureProjectGitignore_writesFile(t *testing.T) {
	project := t.TempDir()
	if err := ensureProjectGitignore(project); err != nil {
		t.Fatal(err)
	}
	if err := ensureProjectGitignore(project); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(filepath.Join(project, ".gitignore"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), gitignoreBegin) {
		t.Fatalf("gitignore = %s", raw)
	}
}

func TestSyncCloudFromDesktop_copiesEnabledPaths(t *testing.T) {
	cat := minimalCatalog()
	desktop := &manifest{
		Skills: []string{"skills/commit"},
		Rules:  []string{"rules/unslop.mdc"},
	}
	items := buildCloudTree(cat, newCloudManifest("v1.0.1"))
	syncCloudFromDesktop(items, desktop)
	var enabled []string
	for _, it := range items {
		if !it.IsGroup && it.Enabled {
			enabled = append(enabled, it.Path)
		}
	}
	sort.Strings(enabled)
	want := []string{"rules/unslop.mdc", "skills/commit"}
	if strings.Join(enabled, ",") != strings.Join(want, ",") {
		t.Fatalf("enabled = %v want %v", enabled, want)
	}
}

func TestCatalogAtRef_usesLocalWhenRefMatches(t *testing.T) {
	cat := minimalCatalog()
	cat.Catalog.Version = "1.0.1"
	got, err := catalogAtRef("/team", "v1.0.1", cat)
	if err != nil {
		t.Fatal(err)
	}
	if got != cat {
		t.Fatal("expected local catalog pointer")
	}
}

func TestLoadCloudManifest_missingReturnsNil(t *testing.T) {
	dir := t.TempDir()
	cm, err := loadCloudManifest(dir)
	if err != nil {
		t.Fatal(err)
	}
	if cm != nil {
		t.Fatalf("expected nil, got %+v", cm)
	}
}

package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

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

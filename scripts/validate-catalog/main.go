// validate-catalog checks catalog.json against the repo tree.
// Exit 0 when every catalog path exists and every cataloged skill/rule is listed.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type catalogFile struct {
	Catalog struct {
		Version string `json:"version"`
	} `json:"catalog"`
	Skills  map[string]entry `json:"skills"`
	Rules   map[string]entry `json:"rules"`
	Scripts map[string]scriptEntry `json:"scripts"`
}

type entry struct {
	Path  string `json:"path"`
	Label string `json:"label"`
}

type scriptEntry struct {
	Path  string `json:"path"`
	Label string `json:"label"`
	Role  string `json:"role"`
}

func main() {
	root, err := os.Getwd()
	if err != nil {
		fail("cwd: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(root, "catalog.json"))
	if err != nil {
		fail("read catalog.json: %v", err)
	}

	var cat catalogFile
	if err := json.Unmarshal(data, &cat); err != nil {
		fail("parse catalog.json: %v", err)
	}

	if cat.Catalog.Version == "" {
		fail("catalog.version is required")
	}

	var errors []string

	checkEntry := func(kind, key string, path string) {
		if path == "" {
			errors = append(errors, fmt.Sprintf("%s %q: missing path", kind, key))
			return
		}
		full := filepath.Join(root, filepath.FromSlash(path))
		if _, err := os.Stat(full); err != nil {
			errors = append(errors, fmt.Sprintf("%s %q: path missing on disk: %s", kind, key, path))
		}
	}

	for key, e := range cat.Skills {
		checkEntry("skill", key, e.Path)
	}
	for key, e := range cat.Rules {
		checkEntry("rule", key, e.Path)
	}
	for key, e := range cat.Scripts {
		checkEntry("script", key, e.Path)
	}

	catalogSkillPaths := map[string]bool{}
	for _, e := range cat.Skills {
		catalogSkillPaths[filepath.ToSlash(e.Path)] = true
	}

	err = filepath.WalkDir(filepath.Join(root, "skills"), func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Name() != "SKILL.md" {
			return nil
		}
		dir := filepath.Dir(path)
		rel, err := filepath.Rel(root, dir)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if !catalogSkillPaths[rel] {
			errors = append(errors, fmt.Sprintf("disk skill not in catalog: %s", rel))
		}
		return nil
	})
	if err != nil {
		fail("walk skills: %v", err)
	}

	catalogRulePaths := map[string]bool{}
	for _, e := range cat.Rules {
		catalogRulePaths[filepath.ToSlash(e.Path)] = true
	}

	ruleFiles, err := filepath.Glob(filepath.Join(root, "rules", "*.mdc"))
	if err != nil {
		fail("glob rules: %v", err)
	}
	for _, rf := range ruleFiles {
		rel, err := filepath.Rel(root, rf)
		if err != nil {
			fail("rel rule: %v", err)
		}
		rel = filepath.ToSlash(rel)
		if !catalogRulePaths[rel] {
			errors = append(errors, fmt.Sprintf("disk rule not in catalog: %s", rel))
		}
	}

	if len(errors) > 0 {
		fmt.Fprintf(os.Stderr, "catalog validation failed (%d issues):\n", len(errors))
		for _, e := range errors {
			fmt.Fprintf(os.Stderr, "  - %s\n", e)
		}
		os.Exit(1)
	}

	fmt.Printf("OK: catalog.json valid (version %s, %d skills, %d rules, %d scripts)\n",
		cat.Catalog.Version, len(cat.Skills), len(cat.Rules), len(cat.Scripts))
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "FAIL: "+format+"\n", args...)
	os.Exit(1)
}

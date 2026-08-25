// agent-config-wizard - desktop TUI for opt-in agent-config apply.
//
// Team root: walk up from binary or source dir until catalog.json is found.
// Project root: cwd, or -project flag.
//
// Run from a project repo:
//
//	C:\path\to\agent-config\scripts\wizard\bin\agent-config-wizard.exe
//
// Dev (go run resolves modules from -C directory):
//
//	go run -C C:\path\to\agent-config\scripts\wizard\src . -project C:\path\to\my-app
//
// Manifest: .cursor/agent-config.local.json (gitignored in real projects).
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
	"unicode/utf16"

	"golang.org/x/text/unicode/norm"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/viewport"
)

const (
	manifestName        = "agent-config.local.json"
	cloudManifestName   = "agent-manifest.json"
	environmentJSONName = "environment.json"
	sourceRepo          = "KroniK907/agent-config"
	envRuleRelPath      = ".cursor/rules/environment.mdc"
	envRuleGeneratedBy  = "agent-config-wizard"
	githubAPIReleases   = "https://api.github.com/repos/" + sourceRepo + "/releases"
	githubRawCatalog    = "https://raw.githubusercontent.com/" + sourceRepo + "/%s/catalog.json"
	gitignoreBegin      = "# >>> agent-config-wizard >>>"
	gitignoreEnd        = "# <<< agent-config-wizard <<<"
)

// --- manifest (AGENT-CFG-GM-006) ---

type manifest struct {
	Source           sourceMeta        `json:"source"`
	ProjectPath      string            `json:"projectPath"`
	Skills           []string          `json:"skills"`
	Rules            []string          `json:"rules"`
	EnvDetails       bool              `json:"envDetails"`
	LastApplied      map[string]string `json:"lastApplied,omitempty"`
	LastCatalogPaths []string          `json:"lastCatalogPaths,omitempty"`
}

type sourceMeta struct {
	Repo     string `json:"repo"`
	Ref      string `json:"ref"`
	TeamPath string `json:"teamPath"`
}

func manifestPath(projectRoot string) string {
	return filepath.Join(projectRoot, ".cursor", manifestName)
}

func loadManifest(projectRoot string) (*manifest, error) {
	data, err := os.ReadFile(manifestPath(projectRoot))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var m manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

func saveManifest(projectRoot string, m *manifest) error {
	dir := filepath.Join(projectRoot, ".cursor")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(manifestPath(projectRoot), data, 0o644)
}

func enabledSet(m *manifest) map[string]bool {
	out := map[string]bool{}
	if m == nil {
		return out
	}
	for _, p := range m.Rules {
		out[p] = true
	}
	for _, p := range m.Skills {
		out[p] = true
	}
	return out
}

func newManifest(teamRoot, projectRoot, ref string) *manifest {
	return &manifest{
		Source: sourceMeta{
			Repo:     sourceRepo,
			Ref:      ref,
			TeamPath: teamRoot,
		},
		ProjectPath: projectRoot,
		Skills:      nil,
		Rules:       nil,
		LastApplied: map[string]string{},
	}
}

// --- cloud manifest (AGENT-CFG-GM-006) ---

type cloudManifest struct {
	Source cloudSourceMeta `json:"source"`
	Skills []string        `json:"skills"`
	Rules  []string        `json:"rules"`
}

type cloudSourceMeta struct {
	Repo string `json:"repo"`
	Ref  string `json:"ref"`
}

type environmentFile struct {
	Build   environmentBuild `json:"build"`
	Install string           `json:"install"`
}

type environmentBuild struct {
	Dockerfile string `json:"dockerfile"`
	Context    string `json:"context,omitempty"`
}

type cloudPathDiff struct {
	Added   []string `json:"added"`
	Removed []string `json:"removed"`
}

var (
	fetchReleaseTagsFn  = fetchReleaseTagsHTTP
	fetchCatalogAtRefFn = fetchCatalogAtRefHTTP
)

func cloudManifestPath(projectRoot string) string {
	return filepath.Join(projectRoot, ".cursor", cloudManifestName)
}

func environmentJSONPath(projectRoot string) string {
	return filepath.Join(projectRoot, ".cursor", environmentJSONName)
}

func cloudDockerfilePath(projectRoot string) string {
	return filepath.Join(projectRoot, ".cursor", "Dockerfile")
}

const cloudDockerfileBody = `# Minimal base for agent-config cloud bootstrap.
# Extend for project toolchain (Go version, Node, etc.).
FROM ubuntu:24.04

RUN apt-get update \
    && apt-get install -y --no-install-recommends git curl jq ca-certificates sudo \
    && rm -rf /var/lib/apt/lists/*
`

func loadCloudManifest(projectRoot string) (*cloudManifest, error) {
	data, err := os.ReadFile(cloudManifestPath(projectRoot))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var cm cloudManifest
	if err := json.Unmarshal(data, &cm); err != nil {
		return nil, err
	}
	return &cm, nil
}

func saveCloudManifest(projectRoot string, cm *cloudManifest) error {
	dir := filepath.Join(projectRoot, ".cursor")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cm, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(cloudManifestPath(projectRoot), data, 0o644)
}

func saveEnvironmentJSON(projectRoot, ref string) error {
	dir := filepath.Join(projectRoot, ".cursor")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	env := environmentFile{
		Build: environmentBuild{
			Dockerfile: "Dockerfile",
			Context:    ".",
		},
		Install: environmentInstallCommand(ref),
	}
	data, err := json.MarshalIndent(env, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(environmentJSONPath(projectRoot), data, 0o644)
}

func saveCloudDockerfile(projectRoot string) error {
	dir := filepath.Join(projectRoot, ".cursor")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	path := cloudDockerfilePath(projectRoot)
	if _, err := os.Stat(path); err == nil {
		return nil // keep project customizations
	}
	return os.WriteFile(path, []byte(cloudDockerfileBody), 0o644)
}

func environmentInstallCommand(ref string) string {
	return fmt.Sprintf(
		"curl -fsSL https://raw.githubusercontent.com/%s/%s/scripts/bootstrap-agent.sh | bash",
		sourceRepo, ref,
	)
}

func environmentInstallURL(ref string) string {
	return environmentInstallCommand(ref)
}

func writeCloudConfig(projectRoot string, cm *cloudManifest) error {
	if err := saveCloudManifest(projectRoot, cm); err != nil {
		return fmt.Errorf("agent-manifest.json: %w", err)
	}
	if err := saveEnvironmentJSON(projectRoot, cm.Source.Ref); err != nil {
		return fmt.Errorf("environment.json: %w", err)
	}
	if err := saveCloudDockerfile(projectRoot); err != nil {
		return fmt.Errorf("Dockerfile: %w", err)
	}
	if err := ensureProjectGitignore(projectRoot); err != nil {
		return fmt.Errorf(".gitignore: %w", err)
	}
	return nil
}

func agentConfigGitignoreLines() []string {
	return []string{
		gitignoreBegin,
		"# Desktop local state stays untracked; cloud bootstrap JSON is committed.",
		"!.cursor/",
		".cursor/agent-config.local.json",
		".cursor/agent-config/",
		"!.cursor/agent-manifest.json",
		"!.cursor/environment.json",
		"!.cursor/Dockerfile",
		gitignoreEnd,
	}
}

func mergeGitignoreBlock(existing string) string {
	block := strings.Join(agentConfigGitignoreLines(), "\n") + "\n"
	begin := strings.Index(existing, gitignoreBegin)
	end := strings.Index(existing, gitignoreEnd)
	if begin >= 0 && end > begin {
		end += len(gitignoreEnd)
		for end < len(existing) && (existing[end] == '\n' || existing[end] == '\r') {
			end++
		}
		return existing[:begin] + block + existing[end:]
	}
	if existing == "" {
		return block
	}
	if !strings.HasSuffix(existing, "\n") {
		existing += "\n"
	}
	return existing + "\n" + block
}

func ensureProjectGitignore(projectRoot string) error {
	path := filepath.Join(projectRoot, ".gitignore")
	existing := ""
	if data, err := os.ReadFile(path); err == nil {
		existing = string(data)
	} else if !os.IsNotExist(err) {
		return err
	}
	merged := mergeGitignoreBlock(existing)
	if merged == existing {
		return nil
	}
	return os.WriteFile(path, []byte(merged), 0o644)
}

func enabledCloudSet(cm *cloudManifest) map[string]bool {
	out := map[string]bool{}
	if cm == nil {
		return out
	}
	for _, p := range cm.Rules {
		out[p] = true
	}
	for _, p := range cm.Skills {
		out[p] = true
	}
	return out
}

func newCloudManifest(ref string) *cloudManifest {
	return &cloudManifest{
		Source: cloudSourceMeta{Repo: sourceRepo, Ref: ref},
		Skills: nil,
		Rules:  nil,
	}
}

func buildCloudTree(cat *catalogFile, cm *cloudManifest) []treeItem {
	prevEnabled := enabledCloudSet(cm)

	var items []treeItem
	type ruleRow struct {
		key string
		e   catalogEntry
	}
	var rules []ruleRow
	for k, e := range cat.Rules {
		rules = append(rules, ruleRow{k, e})
	}
	sort.Slice(rules, func(i, j int) bool { return rules[i].e.Path < rules[j].e.Path })
	for _, r := range rules {
		items = append(items, cloudLeafItem(r.key, r.e, "rule", 0, prevEnabled))
	}
	items = append(items, buildCloudSkillItems(cat, prevEnabled)...)
	return items
}

func buildCloudSkillItems(cat *catalogFile, prevEnabled map[string]bool) []treeItem {
	root := &skillNode{children: map[string]*skillNode{}}
	for key, e := range cat.Skills {
		insertSkillNode(root, key, e)
	}
	var items []treeItem
	for _, name := range sortedSkillKeys(root.children) {
		emitCloudSkillNode(root.children[name], 0, prevEnabled, &items)
	}
	return items
}

func emitCloudSkillNode(n *skillNode, depth int, prevEnabled map[string]bool, items *[]treeItem) {
	if len(n.children) > 0 {
		label := n.segment + "/"
		if n.entry != nil {
			label = n.entry.Label
		}
		*items = append(*items, treeItem{
			Key:             n.key,
			Path:            n.fullPath,
			Label:           label,
			Kind:            "skill",
			Depth:           depth,
			IsGroup:         true,
			DescendantPaths: collectSkillCatalogPaths(n),
		})
		for _, name := range sortedSkillKeys(n.children) {
			emitCloudSkillNode(n.children[name], depth+1, prevEnabled, items)
		}
		return
	}
	if n.entry != nil {
		*items = append(*items, cloudLeafItem(n.key, *n.entry, "skill", depth, prevEnabled))
	}
}

func cloudLeafItem(key string, e catalogEntry, kind string, depth int, prevEnabled map[string]bool) treeItem {
	return treeItem{
		Key:     key,
		Path:    e.Path,
		Label:   e.Label,
		Kind:    kind,
		Depth:   depth,
		Enabled: prevEnabled[e.Path],
	}
}

func syncCloudDraftFromItems(cm *cloudManifest, items []treeItem) {
	var rules, skills []string
	for _, it := range items {
		if it.IsGroup || !it.Enabled {
			continue
		}
		if it.Kind == "rule" {
			rules = append(rules, it.Path)
		} else if it.Kind == "skill" {
			skills = append(skills, it.Path)
		}
	}
	sort.Strings(rules)
	sort.Strings(skills)
	cm.Rules = rules
	cm.Skills = skills
}

func syncCloudFromDesktop(items []treeItem, desktop *manifest) {
	if desktop == nil {
		return
	}
	desktopPaths := enabledSet(desktop)
	for i := range items {
		if items[i].IsGroup {
			continue
		}
		items[i].Enabled = desktopPaths[items[i].Path]
	}
}

func diffCloudPaths(oldSkills, oldRules, newSkills, newRules []string) cloudPathDiff {
	oldSet := map[string]bool{}
	for _, p := range append(append([]string{}, oldSkills...), oldRules...) {
		oldSet[p] = true
	}
	newSet := map[string]bool{}
	for _, p := range append(append([]string{}, newSkills...), newRules...) {
		newSet[p] = true
	}
	var diff cloudPathDiff
	for p := range newSet {
		if !oldSet[p] {
			diff.Added = append(diff.Added, p)
		}
	}
	for p := range oldSet {
		if !newSet[p] {
			diff.Removed = append(diff.Removed, p)
		}
	}
	sort.Strings(diff.Added)
	sort.Strings(diff.Removed)
	return diff
}

func validateCloudPaths(cat *catalogFile, skills, rules []string) []string {
	if cat == nil {
		return []string{"catalog unavailable"}
	}
	catalogPaths := map[string]bool{}
	for _, e := range cat.Rules {
		catalogPaths[e.Path] = true
	}
	for _, e := range cat.Skills {
		catalogPaths[e.Path] = true
	}
	var errs []string
	for _, p := range append(append([]string{}, skills...), rules...) {
		if !catalogPaths[p] {
			errs = append(errs, fmt.Sprintf("path not in catalog at ref: %s", p))
		}
	}
	return errs
}

func parseReleaseTagNames(data []byte) ([]string, error) {
	var releases []struct {
		TagName string `json:"tag_name"`
	}
	if err := json.Unmarshal(data, &releases); err != nil {
		return nil, err
	}
	var tags []string
	for _, r := range releases {
		if r.TagName != "" {
			tags = append(tags, r.TagName)
		}
	}
	return tags, nil
}

func semverTagForCatalogVersion(version string, tags []string) (string, bool) {
	if version == "" {
		return "", false
	}
	candidates := []string{"v" + version, version}
	for _, c := range candidates {
		for _, t := range tags {
			if t == c {
				return t, true
			}
		}
	}
	return "", false
}

func defaultCloudRef(catalogVersion string, tags []string, existing string) string {
	if existing != "" && contains(tags, existing) {
		return existing
	}
	if tag, ok := semverTagForCatalogVersion(catalogVersion, tags); ok {
		return tag
	}
	if len(tags) > 0 {
		return tags[0]
	}
	if catalogVersion != "" {
		return "v" + catalogVersion
	}
	return ""
}

func fetchReleaseTagsHTTP(repo string) ([]string, error) {
	url := githubAPIReleases
	if repo != "" && repo != sourceRepo {
		url = "https://api.github.com/repos/" + repo + "/releases"
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "agent-config-wizard")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github releases API: %s", resp.Status)
	}
	return parseReleaseTagNames(body)
}

func fetchCatalogAtRefHTTP(repo, ref string) (*catalogFile, error) {
	if repo == "" {
		repo = sourceRepo
	}
	url := fmt.Sprintf(githubRawCatalog, ref)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "agent-config-wizard")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("catalog fetch at %s: %s", ref, resp.Status)
	}
	var cat catalogFile
	if err := json.Unmarshal(body, &cat); err != nil {
		return nil, err
	}
	return &cat, nil
}

func catalogAtRef(teamRoot, ref string, localCat *catalogFile) (*catalogFile, error) {
	if localCat != nil {
		localRef := localCat.Catalog.Version
		if ref == localRef || ref == "v"+localRef {
			return localCat, nil
		}
	}
	return fetchCatalogAtRefFn(sourceRepo, ref)
}

func formatCloudDiff(diff cloudPathDiff) string {
	if len(diff.Added) == 0 && len(diff.Removed) == 0 {
		return "sync diff: no path changes"
	}
	var b strings.Builder
	b.WriteString("sync diff:")
	for _, p := range diff.Added {
		b.WriteString("\n  + ")
		b.WriteString(p)
	}
	for _, p := range diff.Removed {
		b.WriteString("\n  - ")
		b.WriteString(p)
	}
	return b.String()
}

// --- catalog ---

type catalogFile struct {
	Catalog struct {
		Version string `json:"version"`
	} `json:"catalog"`
	Skills map[string]catalogEntry `json:"skills"`
	Rules  map[string]catalogEntry `json:"rules"`
}

type catalogEntry struct {
	Path  string `json:"path"`
	Label string `json:"label"`
}

type treeItem struct {
	Key             string
	Path            string // catalog path for leaves; group prefix for groups
	Label           string
	Kind            string // rule | skill
	Depth           int
	Enabled         bool
	IsGroup         bool
	DescendantPaths []string // catalog paths under a group (includes hub path when present)
	ProjectOverride bool
	IsNew           bool
}

type skillNode struct {
	segment  string
	fullPath string
	entry    *catalogEntry
	key      string
	children map[string]*skillNode
}

func loadCatalog(teamRoot string) (*catalogFile, error) {
	data, err := os.ReadFile(filepath.Join(teamRoot, "catalog.json"))
	if err != nil {
		return nil, err
	}
	var c catalogFile
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func envDetailsItem(enabled bool, projectRoot string) treeItem {
	return treeItem{
		Key:             "env-details",
		Path:            envRuleRelPath,
		Label:           "Environment details rule",
		Kind:            "env",
		Enabled:         enabled,
		ProjectOverride: envRuleProjectOverride(projectRoot),
	}
}

func buildTree(cat *catalogFile, prev *manifest, projectRoot string) []treeItem {
	prevEnabled := enabledSet(prev)

	var items []treeItem
	envEnabled := false
	if prev != nil {
		envEnabled = prev.EnvDetails
	}
	items = append(items, envDetailsItem(envEnabled, projectRoot))

	type ruleRow struct {
		key string
		e   catalogEntry
	}
	var rules []ruleRow
	for k, e := range cat.Rules {
		rules = append(rules, ruleRow{k, e})
	}
	sort.Slice(rules, func(i, j int) bool { return rules[i].e.Path < rules[j].e.Path })
	for _, r := range rules {
		items = append(items, makeLeafItem(r.key, r.e, "rule", 0, prev, prevEnabled, projectRoot))
	}

	items = append(items, buildSkillItems(cat, prev, prevEnabled, projectRoot)...)

	return items
}

func buildSkillItems(cat *catalogFile, prev *manifest, prevEnabled map[string]bool, projectRoot string) []treeItem {
	root := &skillNode{children: map[string]*skillNode{}}
	for key, e := range cat.Skills {
		insertSkillNode(root, key, e)
	}

	var items []treeItem
	for _, name := range sortedSkillKeys(root.children) {
		emitSkillNode(root.children[name], 0, prev, prevEnabled, projectRoot, &items)
	}
	return items
}

func insertSkillNode(root *skillNode, key string, e catalogEntry) {
	rel := strings.TrimPrefix(e.Path, "skills/")
	parts := strings.Split(rel, "/")
	cur := root
	for i, part := range parts {
		if cur.children == nil {
			cur.children = map[string]*skillNode{}
		}
		if cur.children[part] == nil {
			cur.children[part] = &skillNode{segment: part}
		}
		cur = cur.children[part]
		cur.fullPath = "skills/" + strings.Join(parts[:i+1], "/")
	}
	cur.entry = &e
	cur.key = key
}

func sortedSkillKeys(m map[string]*skillNode) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func emitSkillNode(n *skillNode, depth int, prev *manifest, prevEnabled map[string]bool, projectRoot string, items *[]treeItem) {
	if len(n.children) > 0 {
		label := n.segment + "/"
		if n.entry != nil {
			label = n.entry.Label
		}
		*items = append(*items, treeItem{
			Key:             n.key,
			Path:            n.fullPath,
			Label:           label,
			Kind:            "skill",
			Depth:           depth,
			IsGroup:         true,
			DescendantPaths: collectSkillCatalogPaths(n),
		})
		for _, name := range sortedSkillKeys(n.children) {
			emitSkillNode(n.children[name], depth+1, prev, prevEnabled, projectRoot, items)
		}
		return
	}
	if n.entry != nil {
		*items = append(*items, makeLeafItem(n.key, *n.entry, "skill", depth, prev, prevEnabled, projectRoot))
	}
}

func collectSkillCatalogPaths(n *skillNode) []string {
	var paths []string
	var walk func(*skillNode)
	walk = func(node *skillNode) {
		if node.entry != nil {
			paths = append(paths, node.fullPath)
		}
		for _, name := range sortedSkillKeys(node.children) {
			walk(node.children[name])
		}
	}
	walk(n)
	sort.Strings(paths)
	return paths
}

func makeLeafItem(key string, e catalogEntry, kind string, depth int, prev *manifest, prevEnabled map[string]bool, projectRoot string) treeItem {
	override := projectOverride(projectRoot, e.Path, kind, prev)
	wasEnabled := prevEnabled[e.Path]
	isNew := prev != nil && !contains(prev.LastCatalogPaths, e.Path)

	enabled := wasEnabled
	if isNew {
		enabled = false // AGENT-CFG-GM-007: new entries default opt-out
	}

	return treeItem{
		Key:             key,
		Path:            e.Path,
		Label:           e.Label,
		Kind:            kind,
		Depth:           depth,
		Enabled:         enabled,
		ProjectOverride: override,
		IsNew:           isNew,
	}
}

func projectDest(projectRoot, catalogPath, kind string) string {
	switch kind {
	case "rule":
		base := strings.TrimPrefix(catalogPath, "rules/")
		return filepath.Join(projectRoot, ".cursor", "rules", base)
	case "skill":
		rest := strings.TrimPrefix(catalogPath, "skills/")
		return filepath.Join(projectRoot, ".cursor", "skills", rest)
	default:
		return filepath.Join(projectRoot, catalogPath)
	}
}

func catalogPaths(cat *catalogFile) []string {
	var paths []string
	for _, e := range cat.Rules {
		paths = append(paths, e.Path)
	}
	for _, e := range cat.Skills {
		paths = append(paths, e.Path)
	}
	sort.Strings(paths)
	return paths
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

func projectOverride(projectRoot, catalogPath, kind string, prev *manifest) bool {
	if kind == "env" || catalogPath == envRuleRelPath {
		return envRuleProjectOverride(projectRoot)
	}
	dest := projectDest(projectRoot, catalogPath, kind)
	if _, err := os.Stat(dest); err != nil {
		return false
	}
	if syncVal, ok := readAgentConfigSync(dest, kind); ok {
		return syncVal == "false"
	}
	if prev != nil {
		if contains(prev.Skills, catalogPath) || contains(prev.Rules, catalogPath) {
			return false
		}
		if prev.LastApplied != nil {
			if _, managed := prev.LastApplied[catalogPath]; managed {
				return false
			}
		}
	}
	return true
}

func envRuleProjectOverride(projectRoot string) bool {
	path := environmentRulePath(projectRoot)
	if _, err := os.Stat(path); err != nil {
		return false
	}
	if isGeneratedEnvironmentRule(path) {
		return false
	}
	return true
}

func previouslyManaged(m *manifest) map[string]bool {
	out := map[string]bool{}
	if m == nil {
		return out
	}
	for _, p := range m.Skills {
		out[p] = true
	}
	for _, p := range m.Rules {
		out[p] = true
	}
	for p := range m.LastApplied {
		out[p] = true
	}
	return out
}

func readAgentConfigSync(dest, kind string) (string, bool) {
	path := dest
	if kind == "skill" {
		path = filepath.Join(dest, "SKILL.md")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", false
	}
	fm := parseFrontmatter(data)
	v, ok := fm["agent-config-sync"]
	return v, ok
}

func parseFrontmatter(data []byte) map[string]string {
	s := string(data)
	if !strings.HasPrefix(s, "---") {
		return nil
	}
	parts := strings.SplitN(s, "---", 3)
	if len(parts) < 3 {
		return nil
	}
	out := map[string]string{}
	for _, line := range strings.Split(parts[1], "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		out[strings.TrimSpace(k)] = strings.TrimSpace(v)
	}
	return out
}

func stampAgentConfigSync(dest, kind string, enabled bool) error {
	path := dest
	if kind == "skill" {
		path = filepath.Join(dest, "SKILL.md")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	val := "false"
	if enabled {
		val = "true"
	}
	return os.WriteFile(path, upsertFrontmatterField(data, "agent-config-sync", val), 0o644)
}

func upsertFrontmatterField(data []byte, key, value string) []byte {
	s := string(data)
	if !strings.HasPrefix(s, "---") {
		return []byte(fmt.Sprintf("---\n%s: %s\n---\n\n%s", key, value, s))
	}
	parts := strings.SplitN(s, "---", 3)
	if len(parts) < 3 {
		return data
	}
	var lines []string
	found := false
	for _, line := range strings.Split(parts[1], "\n") {
		k, _, ok := strings.Cut(strings.TrimSpace(line), ":")
		if ok && strings.TrimSpace(k) == key {
			lines = append(lines, fmt.Sprintf("%s: %s", key, value))
			found = true
		} else {
			lines = append(lines, line)
		}
	}
	if !found {
		lines = append([]string{fmt.Sprintf("%s: %s", key, value)}, lines...)
	}
	return []byte("---" + strings.Join(lines, "\n") + "---" + parts[2])
}

func removeManaged(projectRoot, catalogPath, kind string) error {
	dest := projectDest(projectRoot, catalogPath, kind)
	if kind == "rule" {
		return os.Remove(dest)
	}
	return os.RemoveAll(dest)
}

// --- apply ---

type applyResult struct {
	Copied  []string
	Removed []string
	Skipped []string
	Errors  []string
}

func applyEnabled(teamRoot, projectRoot string, items []treeItem, m *manifest) applyResult {
	var res applyResult
	now := time.Now().UTC().Format(time.RFC3339)
	if m.LastApplied == nil {
		m.LastApplied = map[string]string{}
	}

	for _, it := range items {
		if it.Kind == "env" || it.Path == envRuleRelPath {
			m.EnvDetails = it.Enabled
			break
		}
	}

	prevManaged := previouslyManaged(m)
	newEnabled := map[string]bool{}
	for _, it := range items {
		if it.IsGroup || isEnvTreeItem(it) || !it.Enabled {
			continue
		}
		newEnabled[it.Path] = true
	}
	if m.EnvDetails {
		newEnabled[envRuleRelPath] = true
	}

	for path := range prevManaged {
		if path == envRuleRelPath {
			continue
		}
		if newEnabled[path] {
			continue
		}
		kind := "skill"
		if strings.HasPrefix(path, "rules/") {
			kind = "rule"
		}
		dest := projectDest(projectRoot, path, kind)
		if syncVal, ok := readAgentConfigSync(dest, kind); ok && syncVal == "false" {
			res.Skipped = append(res.Skipped, fmt.Sprintf("%s (detached, left in place)", path))
			continue
		}
		if err := removeManaged(projectRoot, path, kind); err != nil && !os.IsNotExist(err) {
			res.Errors = append(res.Errors, fmt.Sprintf("%s: remove: %v", path, err))
		} else if err == nil {
			res.Removed = append(res.Removed, path)
		}
		delete(m.LastApplied, path)
	}

	m.Rules = nil
	m.Skills = nil

	for _, it := range items {
		if it.IsGroup || isEnvTreeItem(it) {
			continue
		}
		if !it.Enabled {
			continue
		}
		if it.Kind == "rule" {
			m.Rules = append(m.Rules, it.Path)
		} else {
			m.Skills = append(m.Skills, it.Path)
		}

		if it.ProjectOverride {
			res.Skipped = append(res.Skipped, fmt.Sprintf("%s (project override)", it.Path))
			continue
		}

		src := filepath.Join(teamRoot, filepath.FromSlash(it.Path))
		dest := projectDest(projectRoot, it.Path, it.Kind)

		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			res.Errors = append(res.Errors, fmt.Sprintf("%s: mkdir: %v", it.Path, err))
			continue
		}

		var err error
		if it.Kind == "rule" {
			err = copyFile(src, dest)
		} else {
			err = copyTree(src, dest)
		}
		if err != nil {
			res.Errors = append(res.Errors, fmt.Sprintf("%s: %v", it.Path, err))
			continue
		}
		if err := stampAgentConfigSync(dest, it.Kind, true); err != nil {
			res.Errors = append(res.Errors, fmt.Sprintf("%s: stamp sync: %v", it.Path, err))
		}
		m.LastApplied[it.Path] = now
		res.Copied = append(res.Copied, it.Path)
	}

	sort.Strings(m.Rules)
	sort.Strings(m.Skills)

	if err := refreshAgentConfigFramework(teamRoot, projectRoot); err != nil {
		res.Errors = append(res.Errors, "agent-config framework: "+err.Error())
	} else {
		res.Copied = append(res.Copied, ".cursor/agent-config/ (framework copy)")
	}

	if m.EnvDetails {
		if err := refreshEnvironmentRule(projectRoot); err != nil {
			res.Errors = append(res.Errors, envRuleRelPath+": "+err.Error())
		} else {
			res.Copied = append(res.Copied, envRuleRelPath)
		}
	} else {
		path := environmentRulePath(projectRoot)
		if isGeneratedEnvironmentRule(path) {
			if err := removeEnvironmentRule(projectRoot); err != nil {
				res.Errors = append(res.Errors, envRuleRelPath+": remove: "+err.Error())
			} else {
				res.Removed = append(res.Removed, envRuleRelPath)
			}
		}
	}

	return res
}

func refreshAgentConfigFramework(teamRoot, projectRoot string) error {
	src := filepath.Join(teamRoot, ".cursor", "examples")
	dest := filepath.Join(projectRoot, ".cursor", "agent-config")
	if _, err := os.Stat(src); err != nil {
		return err
	}
	return copyTree(src, dest)
}

func copyFile(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func copyTree(src, dest string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dest, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		return copyFile(path, target)
	})
}

// --- environment details rule ---

type envProbe struct {
	OSLine      string
	ShellLine   string
	ScratchDir  string
	ToolLines   []string
	AbsentLines []string
}

func probeEnvironment() envProbe {
	p := envProbe{
		ScratchDir: defaultScratchDir(),
	}
	p.OSLine = probeOSLine()
	p.ShellLine = probeShellLine()
	p.ToolLines, p.AbsentLines = probeTools()
	return p
}

func defaultScratchDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".cursor", "scratch")
}

func probeOSLine() string {
	switch runtime.GOOS {
	case "windows":
		if ver := windowsVersionString(); ver != "" {
			return fmt.Sprintf("Windows (%s)", ver)
		}
		return "Windows 10/11"
	case "darwin":
		if ver := commandOutput("sw_vers", "-productVersion"); ver != "" {
			return fmt.Sprintf("macOS %s", encodeText(ver))
		}
		return "macOS"
	case "linux":
		if ver := readLinuxOSRelease(); ver != "" {
			return ver
		}
		return "Linux"
	default:
		return fmt.Sprintf("%s (%s)", runtime.GOOS, runtime.GOARCH)
	}
}

func windowsVersionString() string {
	out := commandOutput("cmd", "/c", "ver")
	if out == "" {
		return ""
	}
	return encodeText(out)
}

func readLinuxOSRelease() string {
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return ""
	}
	var name, version string
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			return strings.Trim(strings.TrimPrefix(line, "PRETTY_NAME="), `"`)
		}
		if strings.HasPrefix(line, "NAME=") {
			name = strings.Trim(strings.TrimPrefix(line, "NAME="), `"`)
		}
		if strings.HasPrefix(line, "VERSION_ID=") {
			version = strings.Trim(strings.TrimPrefix(line, "VERSION_ID="), `"`)
		}
	}
	if name != "" && version != "" {
		return name + " " + version
	}
	return name
}

func probeShellLine() string {
	if runtime.GOOS == "windows" {
		if os.Getenv("PSModulePath") != "" {
			ver := commandOutput("powershell", "-NoProfile", "-Command", "$PSVersionTable.PSVersion.ToString()")
			if ver != "" {
				return fmt.Sprintf("PowerShell %s (default for terminal commands)", encodeText(ver))
			}
			return "PowerShell (default for terminal commands)"
		}
		if comspec := os.Getenv("ComSpec"); comspec != "" {
			return fmt.Sprintf("%s (default shell)", filepath.Base(comspec))
		}
		return "cmd.exe (default shell)"
	}
	if shell := os.Getenv("SHELL"); shell != "" {
		return fmt.Sprintf("%s (default shell)", shell)
	}
	return "sh (default shell)"
}

type toolSpec struct {
	name        string
	versionArgs []string
	absentNote  string
	osFilter    string // "", "windows", or "!windows"
}

// toolProbeSpecs lists optional dev tools that are not OS defaults but often appear on PATH.
// absentNote belongs on the primary name in each ecosystem so partial installs do not warn twice.
func toolProbeSpecs() []toolSpec {
	return []toolSpec{
		// Runtimes and languages
		{name: "go", versionArgs: []string{"version"}, absentNote: "Prefer Go for small one-off scripts and utilities."},
		{name: "node", versionArgs: []string{"--version"}, absentNote: "Do not assume `node`, `npm`, or `npx`."},
		{name: "npm", versionArgs: []string{"--version"}},
		{name: "npx", versionArgs: []string{"--version"}},
		{name: "pnpm", versionArgs: []string{"--version"}, absentNote: "Do not assume alternate JS package managers (`pnpm`, `yarn`, `bun`)."},
		{name: "yarn", versionArgs: []string{"--version"}},
		{name: "bun", versionArgs: []string{"--version"}},
		{name: "python", versionArgs: []string{"--version"}, absentNote: "Do not assume `python`, `python3`, or `pip`."},
		{name: "python3", versionArgs: []string{"--version"}},
		{name: "pip", versionArgs: []string{"--version"}},
		{name: "pip3", versionArgs: []string{"--version"}},
		{name: "uv", versionArgs: []string{"--version"}, absentNote: "Do not assume Python package tooling (`uv`)."},
		{name: "rustc", versionArgs: []string{"--version"}},
		{name: "cargo", versionArgs: []string{"--version"}, absentNote: "Do not assume Rust (`rustc`, `cargo`)."},
		{name: "dotnet", versionArgs: []string{"--version"}, absentNote: "Do not assume `.NET` (`dotnet`)."},
		{name: "java", versionArgs: []string{"-version"}},
		// Version control and GitHub
		{name: "git", versionArgs: []string{"--version"}},
		{name: "gh", versionArgs: []string{"--version"}},
		// Containers and orchestration
		{name: "docker", versionArgs: []string{"--version"}, absentNote: "Do not assume `docker` or compose subcommands."},
		{name: "kubectl", versionArgs: []string{"version", "--client"}},
		{name: "helm", versionArgs: []string{"version", "--short"}},
		// Cloud and IaC CLIs
		{name: "aws", versionArgs: []string{"--version"}, absentNote: "Do not assume cloud CLIs (`aws`, `az`, `gcloud`)."},
		{name: "az", versionArgs: []string{"--version"}},
		{name: "gcloud", versionArgs: []string{"--version"}},
		{name: "terraform", versionArgs: []string{"version"}, absentNote: "Do not assume IaC CLIs (`terraform`)."},
		// Native build tooling
		{name: "make", versionArgs: []string{"--version"}, absentNote: "Do not assume native build tools (`make`, `gcc`, `clang`)."},
		{name: "gcc", versionArgs: []string{"--version"}},
		{name: "clang", versionArgs: []string{"--version"}},
		// Shell utilities often installed separately
		{name: "jq", versionArgs: []string{"--version"}},
		{name: "sqlite3", versionArgs: []string{"--version"}},
		{name: "pwsh", versionArgs: []string{"--version"}},
		// Windows-only optional installs (Git Bash, WSL, package managers)
		{name: "bash", versionArgs: []string{"--version"}, osFilter: "windows"},
		{name: "wsl", versionArgs: []string{"--version"}, osFilter: "windows"},
		{name: "winget", versionArgs: []string{"--version"}, osFilter: "windows"},
		{name: "choco", versionArgs: []string{"--version"}, osFilter: "windows"},
		{name: "scoop", versionArgs: []string{"--version"}, osFilter: "windows"},
	}
}

func toolSpecApplies(spec toolSpec) bool {
	switch spec.osFilter {
	case "windows":
		return runtime.GOOS == "windows"
	case "!windows":
		return runtime.GOOS != "windows"
	default:
		return true
	}
}

func probeTools() (present []string, absent []string) {
	seenAbsent := map[string]bool{}
	for _, spec := range toolProbeSpecs() {
		if !toolSpecApplies(spec) {
			continue
		}
		path, err := exec.LookPath(spec.name)
		if err != nil {
			if spec.absentNote != "" && !seenAbsent[spec.absentNote] {
				absent = append(absent, spec.absentNote)
				seenAbsent[spec.absentNote] = true
			}
			continue
		}
		line := fmt.Sprintf("`%s` on PATH (%s)", spec.name, path)
		if len(spec.versionArgs) > 0 {
			if ver := commandOutput(spec.name, spec.versionArgs...); ver != "" {
				if label := formatVersionLabel(ver); label != "" {
					line = fmt.Sprintf("`%s` installed (%s)", spec.name, label)
				}
			}
		}
		present = append(present, line)
	}
	return present, absent
}

func commandOutput(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return normalizeCommandOutput(out)
}

// formatVersionLabel picks the first non-empty output line and encodes it as a single-line label.
func formatVersionLabel(s string) string {
	line := firstNonEmptyLine(s)
	if line == "" {
		return ""
	}
	return encodeText(line)
}

func firstNonEmptyLine(s string) string {
	s = strings.ToValidUTF8(s, "")
	for _, line := range strings.FieldsFunc(s, func(r rune) bool {
		return r == '\n' || r == '\r'
	}) {
		line = strings.TrimSpace(line)
		if line != "" {
			return line
		}
	}
	return ""
}

// encodeText normalizes arbitrary command text to plain UTF-8 on one line.
func encodeText(s string) string {
	s = strings.ToValidUTF8(s, "")
	s = norm.NFC.String(s)
	s = strings.Map(func(r rune) rune {
		switch r {
		case '\t', '\n', '\r':
			return ' '
		default:
			if r < 32 {
				return -1
			}
			return r
		}
	}, s)
	return strings.Join(strings.Fields(s), " ")
}

// normalizeCommandOutput converts raw command bytes to normalized UTF-8 text.
// Some Windows tools (notably wsl.exe) write UTF-16LE to stdout without a BOM.
func normalizeCommandOutput(raw []byte) string {
	if len(raw) == 0 {
		return ""
	}
	text := string(raw)
	if decoded, ok := decodeUTF16LEOutput(raw); ok {
		text = decoded
	}
	text = strings.ToValidUTF8(text, "")
	text = norm.NFC.String(text)
	text = strings.Map(func(r rune) rune {
		if r == '\n' || r == '\r' || r == '\t' || r >= 32 {
			return r
		}
		return -1
	}, text)
	return strings.TrimSpace(text)
}

func decodeUTF16LEOutput(raw []byte) (string, bool) {
	if len(raw) >= 2 && raw[0] == 0xFF && raw[1] == 0xFE {
		raw = raw[2:]
	}
	if len(raw) < 4 || len(raw)%2 != 0 {
		return "", false
	}
	// UTF-16LE ASCII text puts a zero byte after each character.
	if raw[1] != 0 || raw[3] != 0 {
		return "", false
	}
	u16 := make([]uint16, len(raw)/2)
	for i := range u16 {
		u16[i] = uint16(raw[2*i]) | uint16(raw[2*i+1])<<8
	}
	return string(utf16.Decode(u16)), true
}

func renderEnvironmentRule(p envProbe) []byte {
	var b strings.Builder
	b.WriteString("---\n")
	b.WriteString("description: Local machine environment - regenerated by agent-config-wizard on apply.\n")
	b.WriteString("generated-by: " + envRuleGeneratedBy + "\n")
	b.WriteString("alwaysApply: true\n")
	b.WriteString("---\n\n")
	b.WriteString("# Environment\n\n")
	b.WriteString("Generated by agent-config-wizard from the local machine. Manual edits are overwritten on apply.\n\n")
	b.WriteString("## Runtime\n\n")
	b.WriteString(fmt.Sprintf("- **OS:** %s\n", p.OSLine))
	b.WriteString(fmt.Sprintf("- **Shell:** %s\n", p.ShellLine))
	if p.ScratchDir != "" {
		b.WriteString(fmt.Sprintf("- **Agent scratch directory:** `%s`\n", filepath.ToSlash(p.ScratchDir)))
	}
	b.WriteString("\n## Tools on PATH\n\n")
	if len(p.ToolLines) == 0 {
		b.WriteString("- No probed tools found on PATH.\n")
	} else {
		for _, line := range p.ToolLines {
			b.WriteString("- " + line + "\n")
		}
	}
	if len(p.AbsentLines) > 0 {
		b.WriteString("\n## Absent runtimes\n\n")
		for _, line := range p.AbsentLines {
			b.WriteString("- " + line + "\n")
		}
	}
	b.WriteString("\n## Paths\n\n")
	b.WriteString("- Use Windows path syntax and PowerShell idioms when OS is Windows.\n")
	b.WriteString("- Put ephemeral agent artifacts in the scratch directory, not the project root.\n")
	return []byte(b.String())
}

func environmentRulePath(projectRoot string) string {
	return filepath.Join(projectRoot, filepath.FromSlash(envRuleRelPath))
}

func isEnvTreeItem(it treeItem) bool {
	return it.Kind == "env" || it.Path == envRuleRelPath
}

func refreshEnvironmentRule(projectRoot string) error {
	path, err := filepath.Abs(environmentRulePath(projectRoot))
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", dir, err)
	}
	content := renderEnvironmentRule(probeEnvironment())
	if err := os.WriteFile(path, content, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

func removeEnvironmentRule(projectRoot string) error {
	path := environmentRulePath(projectRoot)
	if !isGeneratedEnvironmentRule(path) {
		return nil
	}
	err := os.Remove(path)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func isGeneratedEnvironmentRule(path string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	fm := parseFrontmatter(data)
	return fm["generated-by"] == envRuleGeneratedBy
}

func teamRootFromScript() (string, error) {
	start, err := scriptDir()
	if err != nil {
		return "", err
	}
	dir := start
	for {
		if _, err := os.Stat(filepath.Join(dir, "catalog.json")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("catalog.json not found above %s", start)
		}
		dir = parent
	}
}

func scriptDir() (string, error) {
	if exe, err := os.Executable(); err == nil {
		return filepath.Dir(exe), nil
	}
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("resolve script dir failed")
	}
	return filepath.Dir(file), nil
}

// --- TUI ---

type focusPane int

const (
	paneList focusPane = iota
	paneState
)

type uiMode int

const (
	modeDesktop uiMode = iota
	modeCloud
)

const headerLines = 6
const minFooterLines = 3 // separator, status, help
const maxVisibleErrors = 5

var applySpinnerChars = []rune{'|', '/', '-', '\\'}

type spinnerTickMsg struct{}
type applyDoneMsg struct {
	res     applyResult
	saveErr error
}
type releasesFetchedMsg struct {
	tags []string
	err  error
}
type cloudWriteDoneMsg struct {
	err   error
	diff  cloudPathDiff
	wrote bool
}

type model struct {
	teamRoot      string
	projectRoot   string
	catalog       *catalogFile
	manifest      *manifest
	items         []treeItem
	cloudItems    []treeItem
	cloudDraft    *cloudManifest
	savedCloud    *cloudManifest
	releaseTags   []string
	refIndex      int
	cloudDiff     cloudPathDiff
	cursor        int
	status        string
	errors        []string
	lastApply     applyResult
	stateDump     string
	firstRun      bool
	quitting      bool
	err           error
	focus         focusPane
	mode          uiMode
	width, height int
	listVP        viewport.Model
	stateVP       viewport.Model
	applying      bool
	writingCloud  bool
	fetchingTags  bool
	spinnerFrame  int
}

func initialModel(teamRoot, projectRoot string, cat *catalogFile, m *manifest) model {
	firstRun := m == nil
	if m == nil {
		m = newManifest(teamRoot, projectRoot, cat.Catalog.Version)
	}
	items := buildTree(cat, m, projectRoot)

	savedCloud, _ := loadCloudManifest(projectRoot)
	cloudDraft := newCloudManifest("")
	if savedCloud != nil {
		copy := *savedCloud
		cloudDraft = &copy
	}
	cloudItems := buildCloudTree(cat, cloudDraft)

	mod := model{
		teamRoot:    teamRoot,
		projectRoot: projectRoot,
		catalog:     cat,
		manifest:    m,
		items:       items,
		cloudItems:  cloudItems,
		cloudDraft:  cloudDraft,
		savedCloud:  savedCloud,
		firstRun:    firstRun,
		focus:       paneList,
		mode:        modeDesktop,
		status:      "tab state · space toggle · a apply · c cloud bootstrap · q quit",
		listVP:      viewport.New(80, 20),
		stateVP:     viewport.New(80, 20),
	}
	mod.refreshStateDump(nil)
	mod.syncListContent()
	mod.syncStateContent()
	return mod
}

func (m model) Init() tea.Cmd {
	return tea.WindowSize()
}

func tickSpinnerCmd() tea.Cmd {
	return tea.Tick(120*time.Millisecond, func(time.Time) tea.Msg {
		return spinnerTickMsg{}
	})
}

func (m model) runApplyCmd() tea.Cmd {
	teamRoot := m.teamRoot
	projectRoot := m.projectRoot
	items := m.items
	manifest := m.manifest
	return func() tea.Msg {
		res := applyEnabled(teamRoot, projectRoot, items, manifest)
		saveErr := saveManifest(projectRoot, manifest)
		return applyDoneMsg{res: res, saveErr: saveErr}
	}
}

func (m model) fetchReleasesCmd() tea.Cmd {
	return func() tea.Msg {
		tags, err := fetchReleaseTagsFn(sourceRepo)
		return releasesFetchedMsg{tags: tags, err: err}
	}
}

func (m model) runCloudWriteCmd() tea.Cmd {
	projectRoot := m.projectRoot
	teamRoot := m.teamRoot
	catalog := m.catalog
	items := append([]treeItem(nil), m.cloudItems...)
	draft := *m.cloudDraft
	ref := draft.Source.Ref
	savedSkills := []string{}
	savedRules := []string{}
	if m.savedCloud != nil {
		savedSkills = append(savedSkills, m.savedCloud.Skills...)
		savedRules = append(savedRules, m.savedCloud.Rules...)
	}
	return func() tea.Msg {
		syncCloudDraftFromItems(&draft, items)
		diff := diffCloudPaths(savedSkills, savedRules, draft.Skills, draft.Rules)
		cat, err := catalogAtRef(teamRoot, ref, catalog)
		if err != nil {
			return cloudWriteDoneMsg{err: fmt.Errorf("catalog at %s: %w", ref, err)}
		}
		if pathErrs := validateCloudPaths(cat, draft.Skills, draft.Rules); len(pathErrs) > 0 {
			return cloudWriteDoneMsg{err: fmt.Errorf("%s", strings.Join(pathErrs, "; "))}
		}
		if err := writeCloudConfig(projectRoot, &draft); err != nil {
			return cloudWriteDoneMsg{err: err}
		}
		return cloudWriteDoneMsg{wrote: true, diff: diff}
	}
}

func (m *model) activeItems() *[]treeItem {
	if m.mode == modeCloud {
		return &m.cloudItems
	}
	return &m.items
}

func (m *model) enterCloudMode() tea.Cmd {
	m.mode = modeCloud
	m.cursor = 0
	m.focus = paneList
	m.cloudDiff = cloudPathDiff{}
	saved, _ := loadCloudManifest(m.projectRoot)
	m.savedCloud = saved
	if saved != nil {
		copy := *saved
		m.cloudDraft = &copy
	} else {
		m.cloudDraft = newCloudManifest("")
	}
	m.cloudItems = buildCloudTree(m.catalog, m.cloudDraft)
	m.fetchingTags = true
	m.status = "cloud bootstrap: fetching GitHub releases..."
	m.syncListContent()
	return m.fetchReleasesCmd()
}

func (m *model) applyReleaseTags(tags []string, fetchErr error) {
	m.fetchingTags = false
	m.releaseTags = tags
	if fetchErr != nil {
		m.errors = append(m.errors, "release fetch: "+fetchErr.Error())
		m.status = "cloud bootstrap: release fetch failed (ref picker may be limited)"
	}
	ref := defaultCloudRef(m.catalog.Catalog.Version, tags, m.cloudDraft.Source.Ref)
	if ref == "" {
		ref = "v" + m.catalog.Catalog.Version
	}
	m.cloudDraft.Source.Ref = ref
	m.refIndex = 0
	for i, t := range tags {
		if t == ref {
			m.refIndex = i
			break
		}
	}
	if _, ok := semverTagForCatalogVersion(m.catalog.Catalog.Version, tags); !ok && len(tags) > 0 && fetchErr == nil {
		m.status = fmt.Sprintf("cloud bootstrap: local catalog %s has no matching release tag; using %s", m.catalog.Catalog.Version, ref)
	} else {
		m.status = fmt.Sprintf("cloud bootstrap: ref %s · space toggle · [ ] ref · l sync desktop · w write · c back", ref)
	}
	m.cloudItems = buildCloudTree(m.catalog, m.cloudDraft)
	m.refreshCloudStateDump(nil)
	m.syncListContent()
	m.syncStateContent()
}

func (m *model) cycleCloudRef(delta int) {
	if len(m.releaseTags) == 0 {
		m.status = "cloud bootstrap: no release tags loaded"
		return
	}
	m.refIndex += delta
	if m.refIndex < 0 {
		m.refIndex = len(m.releaseTags) - 1
	}
	if m.refIndex >= len(m.releaseTags) {
		m.refIndex = 0
	}
	m.cloudDraft.Source.Ref = m.releaseTags[m.refIndex]
	m.status = fmt.Sprintf("cloud ref -> %s", m.cloudDraft.Source.Ref)
	m.refreshCloudStateDump(nil)
	m.syncStateContent()
}

func (m *model) syncCloudFromLocal() {
	oldSkills := append([]string(nil), m.cloudDraft.Skills...)
	oldRules := append([]string(nil), m.cloudDraft.Rules...)
	m.syncManifestFromItems()
	syncCloudFromDesktop(m.cloudItems, m.manifest)
	syncCloudDraftFromItems(m.cloudDraft, m.cloudItems)
	m.cloudDiff = diffCloudPaths(oldSkills, oldRules, m.cloudDraft.Skills, m.cloudDraft.Rules)
	m.status = formatCloudDiff(m.cloudDiff)
	m.refreshCloudStateDump(nil)
	m.syncListContent()
	m.syncStateContent()
	m.focus = paneState
	m.stateVP.GotoTop()
}

func cursorMarker(applying bool, frame int) string {
	if applying {
		return string(applySpinnerChars[frame%len(applySpinnerChars)]) + " "
	}
	return "> "
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinnerTickMsg:
		if !m.applying && !m.writingCloud {
			return m, nil
		}
		m.spinnerFrame++
		m.syncListContent()
		if m.applying || m.writingCloud {
			return m, tickSpinnerCmd()
		}
		return m, nil

	case releasesFetchedMsg:
		m.applyReleaseTags(msg.tags, msg.err)
		return m, nil

	case cloudWriteDoneMsg:
		m.writingCloud = false
		if msg.err != nil {
			m.errors = append(m.errors, "cloud write: "+msg.err.Error())
			m.status = "cloud write failed"
			m.focus = paneState
			m.refreshCloudStateDump(nil)
			m.syncStateContent()
			return m, nil
		}
		m.savedCloud = m.cloudDraft
		written := *m.cloudDraft
		m.savedCloud = &written
		m.cloudDraft = &written
		m.status = fmt.Sprintf("cloud config written (.cursor/%s, .cursor/%s)", cloudManifestName, environmentJSONName)
		if len(msg.diff.Added) > 0 || len(msg.diff.Removed) > 0 {
			m.status += " · " + formatCloudDiff(msg.diff)
		}
		m.refreshCloudStateDump(nil)
		m.syncStateContent()
		return m, nil

	case applyDoneMsg:
		m.applying = false
		m.lastApply = msg.res
		m.manifest.LastCatalogPaths = catalogPaths(m.catalog)
		m.errors = append([]string(nil), m.lastApply.Errors...)
		if msg.saveErr != nil {
			m.errors = append(m.errors, "manifest save: "+msg.saveErr.Error())
		}
		m.status = applyStatusMessage(m.lastApply, msg.saveErr)
		m.refreshOverrideFlags()
		m.refreshStateDump(&m.lastApply)
		m.syncListContent()
		m.syncStateContent()
		m.adjustViewportHeights()
		if len(m.errors) > 0 {
			m.focus = paneState
			m.stateVP.GotoTop()
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.adjustViewportHeights()
		m.syncListContent()
		m.syncStateContent()
		m.ensureCursorVisible()
		return m, nil

	case tea.KeyMsg:
		if m.applying || m.writingCloud || m.fetchingTags {
			return m, nil
		}
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "c":
			if m.mode == modeCloud {
				m.mode = modeDesktop
				m.cursor = 0
				m.status = "desktop mode · a apply · c cloud bootstrap · q quit"
				m.refreshStateDump(nil)
				m.syncListContent()
				m.syncStateContent()
				return m, nil
			}
			return m, m.enterCloudMode()
		case "tab":
			if m.focus == paneList {
				m.focus = paneState
				m.status = "state pane · tab back to list"
			} else {
				m.focus = paneList
				m.status = "list pane · tab for state"
			}
		case "up", "k":
			if m.focus == paneList {
				if m.cursor > 0 {
					m.cursor--
				}
				m.syncListContent()
				m.ensureCursorVisible()
			} else {
				m.stateVP.LineUp(1)
			}
		case "down", "j":
			if m.focus == paneList {
				items := m.activeItems()
				if m.cursor < len(*items)-1 {
					m.cursor++
				}
				m.syncListContent()
				m.ensureCursorVisible()
			} else {
				m.stateVP.LineDown(1)
			}
		case "pgup":
			if m.focus == paneList {
				m.listVP.ViewUp()
				m.cursor -= m.listVP.Height
				if m.cursor < 0 {
					m.cursor = 0
				}
				m.syncListContent()
				m.ensureCursorVisible()
			} else {
				m.stateVP.ViewUp()
			}
		case "pgdown":
			if m.focus == paneList {
				m.listVP.ViewDown()
				items := m.activeItems()
				m.cursor += m.listVP.Height
				if m.cursor >= len(*items) {
					m.cursor = len(*items) - 1
				}
				m.syncListContent()
				m.ensureCursorVisible()
			} else {
				m.stateVP.ViewDown()
			}
		case "home":
			if m.focus == paneList {
				m.cursor = 0
				m.listVP.GotoTop()
				m.syncListContent()
			} else {
				m.stateVP.GotoTop()
			}
		case "end":
			if m.focus == paneList {
				items := m.activeItems()
				m.cursor = len(*items) - 1
				m.syncListContent()
				m.ensureCursorVisible()
			} else {
				m.stateVP.GotoBottom()
			}
		case " ":
			if m.focus != paneList {
				break
			}
			m.toggleAtCursor()
			if m.mode == modeCloud {
				syncCloudDraftFromItems(m.cloudDraft, m.cloudItems)
				m.refreshCloudStateDump(nil)
			} else {
				m.syncManifestFromItems()
				m.refreshStateDump(nil)
			}
			m.syncListContent()
			m.syncStateContent()
		case "a":
			if m.mode != modeDesktop {
				m.status = "desktop apply only in desktop mode (press c to exit cloud)"
				break
			}
			m.syncManifestFromItems()
			m.applying = true
			m.spinnerFrame = 0
			m.status = "applying..."
			m.syncListContent()
			return m, tea.Batch(m.runApplyCmd(), tickSpinnerCmd())
		case "l":
			if m.mode != modeCloud {
				break
			}
			m.syncCloudFromLocal()
		case "w":
			if m.mode != modeCloud {
				break
			}
			syncCloudDraftFromItems(m.cloudDraft, m.cloudItems)
			m.writingCloud = true
			m.spinnerFrame = 0
			m.status = "writing cloud config..."
			m.syncListContent()
			return m, tea.Batch(m.runCloudWriteCmd(), tickSpinnerCmd())
		case "[", "left":
			if m.mode == modeCloud {
				m.cycleCloudRef(-1)
			}
		case "]", "right":
			if m.mode == modeCloud {
				m.cycleCloudRef(1)
			}
		case "s":
			if m.mode == modeCloud {
				m.refreshCloudStateDump(nil)
			} else {
				m.refreshStateDump(&m.lastApply)
			}
			m.syncStateContent()
			m.focus = paneState
			m.status = "state refreshed · tab back to list"
		}
	}
	return m, nil
}

func applyStatusMessage(res applyResult, saveErr error) string {
	if saveErr != nil {
		return fmt.Sprintf("apply finished with errors and manifest was NOT saved (%d copied, %d removed, %d skipped)",
			len(res.Copied), len(res.Removed), len(res.Skipped))
	}
	if len(res.Errors) > 0 {
		return fmt.Sprintf("apply finished with %d error(s) (%d copied, %d removed, %d skipped)",
			len(res.Errors), len(res.Copied), len(res.Removed), len(res.Skipped))
	}
	return fmt.Sprintf("applied: %d copied, %d removed, %d skipped",
		len(res.Copied), len(res.Removed), len(res.Skipped))
}

func (m *model) footerLineCount() int {
	n := minFooterLines
	if len(m.errors) == 0 {
		return n
	}
	n++ // ERRORS header
	shown := len(m.errors)
	if shown > maxVisibleErrors {
		shown = maxVisibleErrors
		n++ // "... and N more"
	}
	return n + shown
}

func (m *model) adjustViewportHeights() {
	bodyH := m.height - headerLines - m.footerLineCount()
	if bodyH < 1 {
		bodyH = 1
	}
	m.listVP.Width = m.width
	m.listVP.Height = bodyH
	m.stateVP.Width = m.width
	m.stateVP.Height = bodyH
}

func (m model) wrapWidth() int {
	w := m.width - 2
	if w < 24 {
		return 24
	}
	return w
}

func wrapText(text string, width int) string {
	if width <= 0 || text == "" {
		return text
	}
	var out strings.Builder
	for i, line := range strings.Split(text, "\n") {
		if i > 0 {
			out.WriteByte('\n')
		}
		out.WriteString(wrapLine(line, width))
	}
	return out.String()
}

func wrapLine(line string, width int) string {
	if len(line) <= width {
		return line
	}
	indent := leadingWhitespace(line)
	content := strings.TrimLeft(line, " \t")
	prefix := indent
	var parts []string
	for len(content) > 0 {
		room := width - len(prefix)
		if room < 1 {
			room = 1
		}
		if len(content) <= room {
			parts = append(parts, prefix+content)
			break
		}
		cut := room
		if idx := strings.LastIndex(content[:cut], " "); idx > 0 {
			cut = idx
		}
		parts = append(parts, prefix+strings.TrimRight(content[:cut], " "))
		content = strings.TrimLeft(content[cut:], " ")
		prefix = indent + "  "
	}
	return strings.Join(parts, "\n")
}

func leadingWhitespace(s string) string {
	i := 0
	for i < len(s) && (s[i] == ' ' || s[i] == '\t') {
		i++
	}
	return s[:i]
}

func (m model) renderErrors() string {
	if len(m.errors) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("ERRORS:\n")
	width := m.wrapWidth()
	for i, e := range m.errors {
		if i >= maxVisibleErrors {
			b.WriteString(fmt.Sprintf("  ... and %d more (tab -> state JSON for full dump)\n", len(m.errors)-maxVisibleErrors))
			break
		}
		b.WriteString(wrapText("  ! "+e, width))
		b.WriteByte('\n')
	}
	return b.String()
}

func (m *model) syncListContent() {
	m.listVP.SetContent(m.renderList())
}

func (m *model) syncStateContent() {
	header := "STATE"
	if m.mode == modeCloud {
		header = "CLOUD STATE"
	}
	m.stateVP.SetContent(header + "\n" + m.stateDump)
}

func (m *model) renderList() string {
	items := m.activeItems()
	var b strings.Builder
	if m.mode == modeCloud {
		b.WriteString(fmt.Sprintf("source.ref: %s  ([ ] cycle ref", m.cloudDraft.Source.Ref))
		if len(m.releaseTags) > 0 {
			b.WriteString(fmt.Sprintf(" · %d releases", len(m.releaseTags)))
		}
		b.WriteString(")\n")
	}
	var section string
	for i, it := range *items {
		if it.Kind != section {
			section = it.Kind
			switch section {
			case "env":
				b.WriteString("\nENV DETAILS\n")
			default:
				b.WriteString("\n" + strings.ToUpper(section) + "S\n")
			}
		}
		prefix := strings.Repeat("  ", it.Depth)
		box := "[ ]"
		if it.IsGroup {
			box = "[" + m.groupCheck(i) + "]"
		} else if it.Enabled {
			box = "[x]"
		}
		line := fmt.Sprintf("%s%s %s", prefix, box, it.Label)
		if it.IsGroup {
			line += "  (group)"
		}
		if m.mode == modeDesktop {
			if it.IsNew {
				line += "  NEW"
			}
			if it.ProjectOverride {
				line += "  (project override)"
			}
		}
		if i == m.cursor {
			line = cursorMarker(m.applying || m.writingCloud, m.spinnerFrame) + line
		} else {
			line = "  " + line
		}
		b.WriteString(line + "\n")
	}
	return b.String()
}

// cursorLine returns the content line index (0-based) of the list cursor.
func (m *model) cursorLine() int {
	items := m.activeItems()
	line := 0
	if m.mode == modeCloud {
		line = 1 // ref header line
	}
	var section string
	for i, it := range *items {
		if it.Kind != section {
			section = it.Kind
			if i > 0 || m.mode == modeCloud {
				line++
			}
			line++
		}
		if i == m.cursor {
			return line
		}
		line++
	}
	return line
}

func (m *model) ensureCursorVisible() {
	items := m.activeItems()
	if len(*items) == 0 {
		return
	}
	cl := m.cursorLine()
	vp := &m.listVP
	if cl < vp.YOffset {
		vp.SetYOffset(cl)
	}
	if cl >= vp.YOffset+vp.Height {
		vp.SetYOffset(cl - vp.Height + 1)
	}
}

func (m *model) refreshOverrideFlags() {
	for i := range m.items {
		if m.items[i].IsGroup {
			continue
		}
		m.items[i].ProjectOverride = projectOverride(m.projectRoot, m.items[i].Path, m.items[i].Kind, m.manifest)
	}
}

func (m *model) syncManifestFromItems() {
	m.manifest.ProjectPath = m.projectRoot
	m.manifest.Source.TeamPath = m.teamRoot
	m.manifest.Source.Ref = m.catalog.Catalog.Version
	var rules, skills []string
	m.manifest.EnvDetails = false
	for _, it := range m.items {
		if it.IsGroup {
			continue
		}
		if it.Kind == "env" {
			m.manifest.EnvDetails = it.Enabled
			continue
		}
		if !it.Enabled {
			continue
		}
		if it.Kind == "rule" {
			rules = append(rules, it.Path)
		} else {
			skills = append(skills, it.Path)
		}
	}
	m.manifest.Rules = rules
	m.manifest.Skills = skills
}

func (m *model) leafByPath(path string) *treeItem {
	items := m.activeItems()
	for i := range *items {
		it := &(*items)[i]
		if it.IsGroup {
			continue
		}
		if it.Path == path {
			return it
		}
	}
	return nil
}

func (m *model) groupCheck(idx int) string {
	items := m.activeItems()
	g := (*items)[idx]
	enabled, total := 0, 0
	for _, p := range g.DescendantPaths {
		leaf := m.leafByPath(p)
		if leaf == nil {
			continue
		}
		if m.mode == modeDesktop && leaf.ProjectOverride {
			continue
		}
		total++
		if leaf.Enabled {
			enabled++
		}
	}
	if total == 0 || enabled == 0 {
		return " "
	}
	if enabled == total {
		return "x"
	}
	return "~"
}

func (m *model) setGroupEnabled(group *treeItem, enable bool) {
	items := m.activeItems()
	pathSet := map[string]bool{}
	for _, p := range group.DescendantPaths {
		pathSet[p] = true
	}
	for i := range *items {
		it := &(*items)[i]
		if it.IsGroup || !pathSet[it.Path] {
			continue
		}
		if m.mode == modeDesktop && it.ProjectOverride && enable {
			continue
		}
		it.Enabled = enable
	}
}

func (m *model) toggleAtCursor() {
	items := m.activeItems()
	it := &(*items)[m.cursor]
	if it.IsGroup {
		enable := m.groupCheck(m.cursor) != "x"
		m.setGroupEnabled(it, enable)
		m.status = fmt.Sprintf("group %s -> all %s", it.Label, map[bool]string{true: "on", false: "off"}[enable])
		return
	}
	if m.mode == modeDesktop && it.ProjectOverride {
		m.status = fmt.Sprintf("locked: project override at %s", projectDest(m.projectRoot, it.Path, it.Kind))
		return
	}
	it.Enabled = !it.Enabled
	m.status = fmt.Sprintf("toggled %s -> enabled=%v", it.Path, it.Enabled)
}

func (m *model) refreshCloudStateDump(writeErr error) {
	snapshot := struct {
		CloudDraft  *cloudManifest `json:"cloudDraft"`
		SavedCloud  *cloudManifest `json:"savedCloud,omitempty"`
		ReleaseTags []string       `json:"releaseTags,omitempty"`
		SyncDiff    cloudPathDiff  `json:"syncDiff,omitempty"`
		Environment string         `json:"environmentInstall,omitempty"`
		Timestamp   string         `json:"timestamp"`
	}{
		CloudDraft:  m.cloudDraft,
		SavedCloud:  m.savedCloud,
		ReleaseTags: m.releaseTags,
		SyncDiff:    m.cloudDiff,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}
	if m.cloudDraft != nil && m.cloudDraft.Source.Ref != "" {
		snapshot.Environment = environmentInstallURL(m.cloudDraft.Source.Ref)
	}
	data, _ := json.MarshalIndent(snapshot, "", "  ")
	dump := string(data)
	if writeErr != nil {
		dump = "Cloud write error: " + writeErr.Error() + "\n\n" + dump
	}
	if diffLine := formatCloudDiff(m.cloudDiff); diffLine != "sync diff: no path changes" {
		dump = diffLine + "\n\n" + dump
	}
	m.stateDump = dump
}

func (m *model) refreshStateDump(apply *applyResult) {
	snapshot := struct {
		Manifest  *manifest    `json:"manifest"`
		Apply     *applyResult `json:"apply,omitempty"`
		FirstRun  bool         `json:"firstRun"`
		Timestamp string       `json:"timestamp"`
	}{
		Manifest:  m.manifest,
		Apply:     apply,
		FirstRun:  m.firstRun,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	data, _ := json.MarshalIndent(snapshot, "", "  ")
	dump := string(data)
	if apply != nil && len(apply.Errors) > 0 {
		var b strings.Builder
		width := m.wrapWidth()
		b.WriteString("Apply errors:\n")
		for _, e := range apply.Errors {
			b.WriteString(wrapText("  - "+e, width))
			b.WriteByte('\n')
		}
		b.WriteByte('\n')
		b.WriteString(dump)
		dump = b.String()
	}
	m.stateDump = dump
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("error: %v\n", m.err)
	}
	if m.quitting {
		return "quit\n"
	}

	var b strings.Builder
	b.WriteString("agent-config-wizard\n")
	b.WriteString(fmt.Sprintf("Team: %s\n", m.teamRoot))
	b.WriteString(fmt.Sprintf("Project: %s\n", m.projectRoot))
	if m.mode == modeCloud {
		b.WriteString("Mode: cloud bootstrap (optional · press c to return to desktop)\n")
	} else if m.firstRun {
		b.WriteString("Mode: first-run wizard (c for optional cloud bootstrap)\n")
	} else {
		b.WriteString("Mode: desktop re-run (c for cloud bootstrap)\n")
	}

	if m.focus == paneList {
		if m.mode == modeCloud {
			b.WriteString("\n--- cloud opt-in (tab for state) ---\n")
		} else {
			b.WriteString("\n--- desktop opt-in tree (tab for state) ---\n")
		}
		b.WriteString(m.listVP.View())
	} else {
		b.WriteString("\n--- state JSON (tab for list) ---\n")
		b.WriteString(m.stateVP.View())
	}

	b.WriteString("\n")
	b.WriteString(strings.Repeat("-", min(m.width, 72)) + "\n")
	b.WriteString(m.status + "\n")
	if errBlock := m.renderErrors(); errBlock != "" {
		b.WriteString(errBlock)
	}
	if m.mode == modeCloud {
		b.WriteString("cloud: space toggle · [ ] ref · l sync desktop · w write · c desktop · q quit\n")
	} else {
		b.WriteString("desktop: space toggle · a apply · c cloud · tab state · q quit\n")
	}
	return b.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	projectFlag := flag.String("project", "", "project repo root (default: cwd)")
	flag.Parse()

	teamRoot, err := teamRootFromScript()
	if err != nil {
		fmt.Fprintf(os.Stderr, "team root: %v\n", err)
		os.Exit(1)
	}

	projectRoot := *projectFlag
	if projectRoot == "" {
		projectRoot, err = os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "cwd: %v\n", err)
			os.Exit(1)
		}
	} else {
		projectRoot, err = filepath.Abs(projectRoot)
		if err != nil {
			fmt.Fprintf(os.Stderr, "project path: %v\n", err)
			os.Exit(1)
		}
	}

	cat, err := loadCatalog(teamRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "catalog: %v\n", err)
		os.Exit(1)
	}

	prev, err := loadManifest(projectRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "manifest: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(initialModel(teamRoot, projectRoot, cat, prev), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "tui: %v\n", err)
		os.Exit(1)
	}
}

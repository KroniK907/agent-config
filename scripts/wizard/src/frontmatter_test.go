package main

import (
	"bytes"
	"testing"
)

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

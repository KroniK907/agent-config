#!/usr/bin/env bash
# Cloud Agent Build install hook - manifest-driven copy-only bootstrap.
# Reads .cursor/agent-manifest.json from the workspace (CWD or AGENT_CONFIG_WORKSPACE).
#
# AGENT-CFG-GM-002: skills -> ~/.cursor/skills/<path-after-skills/>; rules -> .cursor/rules/
# AGENT-CFG-GM-005: validate manifest paths against catalog.json at source.ref; no walk fallback.
#
# Requires: git, jq
# Optional env: GH_TOKEN (private clone), AGENT_CONFIG_WORKSPACE (default: pwd)
set -euo pipefail

WORKSPACE="${AGENT_CONFIG_WORKSPACE:-$(pwd)}"
MANIFEST="${WORKSPACE}/.cursor/agent-manifest.json"
SKILLS_HOME="${HOME}/.cursor/skills"
RULES_DIR="${WORKSPACE}/.cursor/rules"
CACHE_ROOT="${HOME}/.cache/agent-config"

die() {
  echo "bootstrap-agent: $*" >&2
  exit 1
}

need_cmd() {
  command -v "$1" >/dev/null 2>&1 || die "missing required command: $1"
}

need_cmd git
need_cmd jq

[[ -f "$MANIFEST" ]] || die "manifest not found: $MANIFEST"

REPO="$(jq -r '.source.repo // empty' "$MANIFEST")"
REF="$(jq -r '.source.ref // empty' "$MANIFEST")"
[[ -n "$REPO" && "$REPO" != "null" ]] || die "manifest source.repo is required"
[[ -n "$REF" && "$REF" != "null" ]] || die "manifest source.ref is required"

CACHE="${CACHE_ROOT}/${REPO//\//-}"
mkdir -p "$SKILLS_HOME" "$RULES_DIR" "$CACHE_ROOT"

clone_url="https://github.com/${REPO}.git"
if [[ -n "${GH_TOKEN:-}" ]]; then
  clone_url="https://x-access-token:${GH_TOKEN}@github.com/${REPO}.git"
fi

if [[ ! -d "${CACHE}/.git" ]]; then
  echo "bootstrap-agent: cloning ${REPO}@${REF} -> ${CACHE}"
  if ! git clone --depth 1 --branch "$REF" "$clone_url" "$CACHE" 2>/dev/null; then
    rm -rf "$CACHE"
    git clone --depth 1 "$clone_url" "$CACHE"
    git -C "$CACHE" checkout "$REF"
  fi
else
  echo "bootstrap-agent: refreshing cache ${CACHE} @ ${REF}"
  git -C "$CACHE" fetch origin tag "$REF" --force 2>/dev/null \
    || git -C "$CACHE" fetch origin "$REF" --force 2>/dev/null \
    || git -C "$CACHE" fetch origin --tags --force
  git -C "$CACHE" checkout "$REF"
fi

CATALOG="${CACHE}/catalog.json"
[[ -f "$CATALOG" ]] || die "catalog.json missing at ${REF} in ${REPO}"

catalog_has_path() {
  local path="$1"
  jq -e --arg p "$path" '
    (.skills[]?.path, .rules[]?.path) | select(. == $p)
  ' "$CATALOG" >/dev/null
}

copy_tree() {
  local src="$1"
  local dest="$2"
  [[ -d "$src" ]] || die "source missing: $src"
  rm -rf "$dest"
  mkdir -p "$(dirname "$dest")"
  cp -R "$src" "$dest"
}

copy_file() {
  local src="$1"
  local dest="$2"
  [[ -f "$src" ]] || die "source missing: $src"
  mkdir -p "$(dirname "$dest")"
  cp "$src" "$dest"
}

skill_count="$(jq '.skills | length' "$MANIFEST")"
rule_count="$(jq '.rules | length' "$MANIFEST")"
echo "bootstrap-agent: ${REPO}@${REF} -> ${skill_count} skill(s), ${rule_count} rule(s)"

if [[ "$skill_count" -eq 0 && "$rule_count" -eq 0 ]]; then
  echo "bootstrap-agent: manifest lists no skills or rules; nothing to copy"
  exit 0
fi

while IFS= read -r path; do
  [[ -n "$path" ]] || continue
  catalog_has_path "$path" || die "path not in catalog at ${REF}: ${path}"
  case "$path" in
    skills/*)
      rel="${path#skills/}"
      src="${CACHE}/${path}"
      dest="${SKILLS_HOME}/${rel}"
      copy_tree "$src" "$dest"
      echo "  skill: ${path} -> ${dest}"
      ;;
    *)
      die "invalid skill path (expected skills/ prefix): ${path}"
      ;;
  esac
done < <(jq -r '.skills[]?' "$MANIFEST")

while IFS= read -r path; do
  [[ -n "$path" ]] || continue
  catalog_has_path "$path" || die "path not in catalog at ${REF}: ${path}"
  case "$path" in
    rules/*.mdc)
      base="${path#rules/}"
      src="${CACHE}/${path}"
      dest="${RULES_DIR}/${base}"
      copy_file "$src" "$dest"
      echo "  rule: ${path} -> ${dest}"
      ;;
    *)
      die "invalid rule path (expected rules/*.mdc): ${path}"
      ;;
  esac
done < <(jq -r '.rules[]?' "$MANIFEST")

echo "bootstrap-agent: complete"

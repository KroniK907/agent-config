#!/usr/bin/env bash
# Idempotent install of KroniK907/skills wayfinder pack into ~/.cursor/skills/
# Called from .cursor/environment.json install on Cloud Agent VMs.
#
# Env (optional):
#   WAYFINDER_SKILLS_REPO  default KroniK907/skills
#   WAYFINDER_SKILLS_TAG   default v0.1.0
#   GH_TOKEN               used for private clone; public repo works without it
set -euo pipefail

REPO="${WAYFINDER_SKILLS_REPO:-KroniK907/skills}"
TAG="${WAYFINDER_SKILLS_TAG:-v0.1.0}"
DEST="${HOME}/.cursor/skills"
CACHE="${HOME}/.cache/wayfinder-skills/${REPO//\//-}"

mkdir -p "$DEST" "$CACHE"

clone_url="https://github.com/${REPO}.git"
if [[ -n "${GH_TOKEN:-}" ]]; then
  clone_url="https://x-access-token:${GH_TOKEN}@github.com/${REPO}.git"
fi

if [[ ! -d "$CACHE/.git" ]]; then
  git clone --depth 1 --branch "$TAG" "$clone_url" "$CACHE"
else
  git -C "$CACHE" fetch origin tag "$TAG" --force 2>/dev/null || git -C "$CACHE" fetch origin --tags --force
  git -C "$CACHE" checkout "$TAG"
fi

# Sync wayfinder ecosystem + repo-root utilities referenced by wayfinder REFERENCE
sync_dir() {
  local src="$1"
  local name="$2"
  if [[ -d "$src" ]]; then
    rm -rf "$DEST/$name"
    cp -R "$src" "$DEST/$name"
    echo "  synced $name"
  fi
}

echo "Installing skills from ${REPO}@${TAG} -> ${DEST}"

sync_dir "$CACHE/wayfinder" "wayfinder"

for util in tdd commit writing-for-agents write-a-prd prd-to-plan prd-to-issues request-refactor-plan triage-issue improve-codebase-architecture ubiquitous-language write-a-skill; do
  sync_dir "$CACHE/$util" "$util"
done

echo "Skills install complete."

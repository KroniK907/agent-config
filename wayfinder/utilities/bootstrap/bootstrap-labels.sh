#!/usr/bin/env bash
# Bootstrap wf:* GitHub labels from labels-manifest.json.
# Idempotent: gh label create --force updates existing labels.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MANIFEST="${1:-$SCRIPT_DIR/labels-manifest.json}"
REPO="${REPO:-$(gh repo view --json nameWithOwner -q .nameWithOwner)}"

echo "Bootstrapping wayfinder labels on $REPO from $MANIFEST ..."

while IFS= read -r line; do
  name=$(echo "$line" | jq -r '.name')
  color=$(echo "$line" | jq -r '.color')
  desc=$(echo "$line" | jq -r '.description // ""')
  gh label create "$name" --repo "$REPO" --color "$color" --description "$desc" --force
  echo "  OK $name"
done < <(jq -c '.[]' "$MANIFEST")

echo "Done."

#!/usr/bin/env bash
# Validate a draft map issue body before gh issue edit --body-file.
# Usage: bash validate-map-body.sh path/to/map-body.md
set -euo pipefail

path="${1:?usage: validate-map-body.sh <map-body.md>}"

if [[ ! -f "$path" ]]; then
  echo "FAIL: file not found: $path" >&2
  exit 1
fi

lines=$(wc -l < "$path" | tr -d ' ')
if (( lines < 40 )); then
  echo "FAIL: body has $lines lines (expected >= 40; collapsed maps are often ~10 lines)" >&2
  exit 1
fi

for section in "## To Do" "## Completed" "## Decision coverage"; do
  if ! grep -qx "$section" "$path"; then
    echo "FAIL: missing section header on its own line: $section" >&2
    exit 1
  fi
done

if grep -qP '[\x{2013}\x{2014}\x{00b7}]' "$path" 2>/dev/null || LC_ALL=C grep -E $'[\xe2\x80\x93\xe2\x80\x94\xc2\xb7]' "$path"; then
  echo "FAIL: unicode em dash, en dash, or middle dot found - use ASCII hyphen only" >&2
  exit 1
fi

if grep -qE 'Ã|Γ|╬|├|┤|┬|┐' "$path"; then
  echo "FAIL: mojibake sequences detected (encoding corruption)" >&2
  exit 1
fi

echo "OK: map body valid ($lines lines, required sections present, ASCII punctuation)"

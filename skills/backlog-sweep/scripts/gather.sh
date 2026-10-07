#!/usr/bin/env bash
# Gather evidence for a backlog sweep into one directory:
#   issues/<N>.md   every open issue with labels, body, and comments
#   index.tsv       number, updated date, labels, title (one line per open issue)
#   merged-prs.tsv  recently merged PRs: number, merged date, base, closing issue refs, title
#   labels.txt      the repo's label names
#
# Usage: gather.sh <owner/repo> <out-dir> [merged-pr-limit]
set -euo pipefail

repo="${1:?usage: gather.sh <owner/repo> <out-dir> [merged-pr-limit]}"
out="${2:?usage: gather.sh <owner/repo> <out-dir> [merged-pr-limit]}"
limit="${3:-100}"

mkdir -p "$out/issues"

gh issue list -R "$repo" --state open --limit 500 \
  --json number,title,labels,updatedAt \
  --jq '.[] | [.number, .updatedAt[:10], ([.labels[].name] | join(",")), .title] | @tsv' \
  > "$out/index.tsv"

cut -f1 "$out/index.tsv" | while read -r n; do
  gh issue view "$n" -R "$repo" --json number,title,body,labels,comments --jq '
    "# \(.number) \(.title)\nLABELS: \([.labels[].name] | join(","))\n\n\(.body)\n\n--- COMMENTS ---\n" +
    ([.comments[] | "[\(.author.login) \(.createdAt[:10])]\n\(.body)\n"] | join("\n"))
  ' > "$out/issues/$n.md"
done

owner="${repo%%/*}"
name="${repo#*/}"
gh api graphql -F owner="$owner" -F name="$name" -F limit="$limit" -f query='
  query($owner: String!, $name: String!, $limit: Int!) {
    repository(owner: $owner, name: $name) {
      pullRequests(states: MERGED, first: $limit, orderBy: {field: UPDATED_AT, direction: DESC}) {
        nodes {
          number title mergedAt baseRefName
          closingIssuesReferences(first: 10) { nodes { number } }
        }
      }
    }
  }' --jq '.data.repository.pullRequests.nodes | sort_by(.mergedAt) | reverse | .[] |
    [.number, .mergedAt[:10], .baseRefName, ([.closingIssuesReferences.nodes[].number | tostring] | join(",")), .title] | @tsv' \
  > "$out/merged-prs.tsv"

gh label list -R "$repo" --limit 200 --json name --jq '.[].name' > "$out/labels.txt"

echo "open issues: $(wc -l < "$out/index.tsv")  merged PRs: $(wc -l < "$out/merged-prs.tsv")  -> $out"

#!/usr/bin/env bash
# Wait until reviewbot posts a review of a PR's head commit, then print it.
#
# Usage: wait-for-review.sh <pr-number> [timeout-seconds] [poll-seconds]
#   Run from inside the repo checkout. Timeout defaults to 540s so a foreground
#   call fits under a 10-minute tool limit; call again after a timeout.
#
# Output on success: the review summary, then one block per inline comment
#   (comment id, path:line, body). The head SHA is on the first line.
# Exit codes: 0 review found, 2 timed out, 3 PR is a draft (reviewbot skips drafts),
#   4 PR is not open.
set -euo pipefail

pr="${1:?usage: wait-for-review.sh <pr-number> [timeout-seconds] [poll-seconds]}"
timeout="${2:-540}"
poll="${3:-30}"

repo="$(gh repo view --json nameWithOwner -q .nameWithOwner)"
deadline=$(( $(date +%s) + timeout ))

while :; do
  read -r state draft sha < <(gh pr view "$pr" --json state,isDraft,headRefOid \
    -q '"\(.state) \(.isDraft) \(.headRefOid)"')
  [ "$state" = OPEN ] || { echo "PR #$pr is $state"; exit 4; }
  [ "$draft" = false ] || { echo "PR #$pr is a draft; run gh pr ready $pr"; exit 3; }

  review_id="$(gh api --paginate "repos/$repo/pulls/$pr/reviews" \
    --jq ".[] | select(.body | contains(\"reviewbot:sha=$sha\")) | .id" | tail -n1)"

  if [ -n "$review_id" ]; then
    echo "head $sha  review $review_id"
    gh api "repos/$repo/pulls/$pr/reviews/$review_id" --jq .body
    echo
    echo "--- INLINE COMMENTS ---"
    gh api --paginate "repos/$repo/pulls/$pr/reviews/$review_id/comments" \
      --jq '.[] | "[comment \(.id)] \(.path):\(.line // .original_line)\n\(.body)\n"'
    exit 0
  fi

  if [ "$(date +%s)" -ge "$deadline" ]; then
    echo "no reviewbot review of $sha after ${timeout}s"
    exit 2
  fi
  sleep "$poll"
done

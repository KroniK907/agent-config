#!/usr/bin/env bash
# Bootstrap Wayfinder-Ecosystem map on GitHub. Requires gh auth with issue write.
set -euo pipefail
REPO="KroniK907/skills"

create_ticket() {
  local title="$1"
  local body_file="$2"
  local labels="$3"
  gh issue create -R "$REPO" --title "$title" --label "$labels" --body-file "$body_file"
}

add_blocked_by() {
  local issue_num="$1"
  local blocker_num="$2"
  local blocker_id
  blocker_id=$(gh api "repos/$REPO/issues/$blocker_num" --jq .id)
  gh api -X POST "repos/$REPO/issues/$issue_num/dependencies/blocked_by" -F "issue_id=$blocker_id" >/dev/null
}

add_sub_issue() {
  local parent_num="$1"
  local child_num="$2"
  local parent_id child_id
  parent_id=$(gh api "repos/$REPO/issues/$parent_num" --jq .node_id)
  child_id=$(gh api "repos/$REPO/issues/$child_num" --jq .node_id)
  gh api graphql -f query="
    mutation(\$parent: ID!, \$child: ID!) {
      addSubIssue(input: { issueId: \$parent, subIssueId: \$child }) { issue { number } }
    }" -f parent="$parent_id" -f child="$child_id" >/dev/null
}

DIR="$(cd "$(dirname "$0")/issue-bodies" && pwd)"

echo "Creating To Do tickets..."
T1=$(create_ticket "Define decision-log handoff to write-a-prd" "$DIR/todo-decision-log-prd.md" "documentation")
T1N=$(echo "$T1" | rg -o '[0-9]+$')
T2=$(create_ticket "Update grill-me to append scoped GM-xx to map decision log" "$DIR/todo-grill-me-integration.md" "documentation")
T2N=$(echo "$T2" | rg -o '[0-9]+$')
T3=$(create_ticket "Specify research ticket workflow (skill or Task pattern)" "$DIR/todo-research-workflow.md" "documentation")
T3N=$(echo "$T3" | rg -o '[0-9]+$')
T4=$(create_ticket "GitHub tracker setup for skills repo maps" "$DIR/todo-github-tracker.md" "documentation")
T4N=$(echo "$T4" | rg -o '[0-9]+$')
T5=$(create_ticket "Subfeature map worked example in REFERENCE or EXAMPLES" "$DIR/todo-subfeature-example.md" "documentation")
T5N=$(echo "$T5" | rg -o '[0-9]+$')
T6=$(create_ticket "Reconcile / post-implementation map update skill or mode" "$DIR/todo-reconcile.md" "documentation")
T6N=$(echo "$T6" | rg -o '[0-9]+$')
T7=$(create_ticket "Cloud AFK automation contract for wayfinder:afk tickets" "$DIR/todo-afk-contract.md" "documentation")
T7N=$(echo "$T7" | rg -o '[0-9]+$')

echo "Creating decision log..."
LOG=$(create_ticket "Wayfinder-Ecosystem:Decision-Log" "$DIR/decision-log.md" "documentation")
LOGN=$(echo "$LOG" | rg -o '[0-9]+$')

echo "Creating map..."
# Inject issue numbers into map body
sed -e "s/{{LOG}}/$LOGN/g" \
    -e "s/{{T1}}/$T1N/g" -e "s/{{T2}}/$T2N/g" -e "s/{{T3}}/$T3N/g" \
    -e "s/{{T4}}/$T4N/g" -e "s/{{T5}}/$T5N/g" -e "s/{{T6}}/$T6N/g" -e "s/{{T7}}/$T7N/g" \
    "$DIR/map.md" > /tmp/wayfinder-map-body.md
MAP=$(gh issue create -R "$REPO" --title "Wayfinder-Ecosystem:Map" --label "documentation" --body-file /tmp/wayfinder-map-body.md)
MAPN=$(echo "$MAP" | rg -o '[0-9]+$')

echo "Patching cross-links..."
# Update decision log and tickets with map number via comments in body - already have MAP placeholder in some files
# Re-create is not possible; patch decision log body reference in a follow-up if edit works for user

echo "Wiring blocked-by..."
add_blocked_by "$T2N" "$T1N"
add_blocked_by "$T6N" "$T1N"
add_blocked_by "$T7N" "$T4N"

echo "Linking sub-issues to map..."
for n in "$T1N" "$T2N" "$T3N" "$T4N" "$T5N" "$T6N" "$T7N" "$LOGN"; do
  add_sub_issue "$MAPN" "$n"
done

echo "Done."
echo "Map: $MAP"
echo "Decision log: $LOG"
echo "Tickets: #$T1N #$T2N #$T3N #$T4N #$T5N #$T6N #$T7N"

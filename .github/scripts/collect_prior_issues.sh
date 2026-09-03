#!/usr/bin/env bash
# Collect earlier digest / doc-drift issues so the model does not re-report known findings.
# Writes /tmp/prior-issues.md and a skip=true|false line to $GITHUB_OUTPUT.
# Needs: GH_TOKEN, GITHUB_REPOSITORY, GITHUB_OUTPUT.
set -euo pipefail

# Two title generations exist: "Doc drift check ..." (early issues) and
# "Claude Code weekly digest ..." (current). Search both so changelog cropping
# and re-report suppression see every prior issue regardless of title era.
queries=('Doc drift check in:title' 'weekly digest in:title')

open_count=0
for q in "${queries[@]}"; do
  n=$(gh issue list --repo "$GITHUB_REPOSITORY" --state open \
    --search "$q" --json number --jq 'length')
  open_count=$((open_count + n))
done
echo "open drift issues: $open_count"
if [ "$open_count" -ge 10 ]; then
  echo "skip=true" >> "$GITHUB_OUTPUT"
  echo "10+ unaddressed drift issues; skipping this run instead of piling on."
  exit 0
fi
echo "skip=false" >> "$GITHUB_OUTPUT"

for q in "${queries[@]}"; do
  gh issue list --repo "$GITHUB_REPOSITORY" --state all \
    --search "$q" --json number,state,title,body,comments
done | jq -rs 'add | unique_by(.number) | .[] | "### #\(.number) [\(.state)] \(.title)\nFindings:\n\(.body | [scan("(?m)^## .*$")] | join("\n"))\nLast comment:\n\(if (.comments|length) > 0 then .comments[-1].body[0:1500] else "(none)" end)\n"' \
  > /tmp/prior-issues.md || true
wc -c /tmp/prior-issues.md

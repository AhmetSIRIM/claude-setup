#!/usr/bin/env bash
# Open one assigned issue when a new setup version is released, so every machine
# pulls the new version. TAG comes from the release event (or the latest release on
# manual dispatch). Older open pull-nudges are superseded. DRY_RUN=1 prints instead.
# Needs: GH_TOKEN, GITHUB_REPOSITORY, TAG.
set -euo pipefail

owner="${GITHUB_REPOSITORY%%/*}"
title="Pull setup $TAG on every machine"

if gh issue list --repo "$GITHUB_REPOSITORY" --state all \
     --search "\"$title\" in:title" --json title --jq '.[].title' | grep -qx "$title"; then
  echo "already announced: $TAG"; exit 0
fi

body="Setup $TAG is released. Per machine (Mac and Windows):

- \`cd\` into the clone and \`git pull\` (links and junctions pick the changes up on their own);
- open a new session: a silent setup-doctor means the machine is current;
- skim the release notes for anything that needs a hand step (a new tool, a settings key).

Close this issue once every machine is on $TAG."
if [ "${DRY_RUN:-0}" = "1" ]; then
  echo "DRY RUN: would open '$title' assigned to $owner"; exit 0
fi
gh issue create --repo "$GITHUB_REPOSITORY" --title "$title" --body "$body" --assignee "$owner"

gh issue list --repo "$GITHUB_REPOSITORY" --state open \
  --search 'Pull setup in:title' --json number,title --jq '.[] | "\(.number)\t\(.title)"' \
  | while IFS=$'\t' read -r num t; do
      [ "$t" = "$title" ] && continue
      gh issue close "$num" --repo "$GITHUB_REPOSITORY" \
        --comment "Superseded by the newer release nudge." || true
    done

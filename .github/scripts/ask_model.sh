#!/usr/bin/env bash
# Send the digest prompt and classify the answer as drift / no drift.
# The prompt goes over stdin (argv has a single-argument kernel limit). opencode is
# occasionally silent forever on a request; one bounded retry covers that.
# Needs: OPENCODE_MODEL, GITHUB_OUTPUT. Reads /tmp/prompt.txt and /tmp/digest-meta.env.
set -euo pipefail
ok=false
for attempt in 1 2; do
  if timeout 300 opencode run -m "$OPENCODE_MODEL" --format default \
       < /tmp/prompt.txt > /tmp/answer.txt 2> /tmp/opencode.err && [ -s /tmp/answer.txt ]; then
    ok=true; break
  fi
  echo "attempt $attempt produced no answer; stderr:"; cat /tmp/opencode.err
done
$ok || { echo "opencode failed twice"; exit 1; }
echo "--- answer ---"; cat /tmp/answer.txt
if grep -q '^NOTHING_TO_REPORT' /tmp/answer.txt; then
  echo "drift=false" >> "$GITHUB_OUTPUT"
else
  echo "drift=true" >> "$GITHUB_OUTPUT"
fi
cat /tmp/digest-meta.env >> "$GITHUB_OUTPUT"

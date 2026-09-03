#!/usr/bin/env bash
# A 60s micro-request separates "provider/account is down" from "the big request
# is the problem", and fails fast instead of burning two 300s retries.
# Needs: OPENCODE_MODEL.
set -euo pipefail
echo "Reply with exactly the word PONG and nothing else." > /tmp/probe.txt
timeout 60 opencode run -m "$OPENCODE_MODEL" --format default \
  < /tmp/probe.txt > /tmp/probe-answer.txt 2> /tmp/probe.err \
  || { echo "provider probe FAILED (provider or account problem, not the digest prompt):"; cat /tmp/probe.err; exit 1; }
echo "probe answer: $(head -c 200 /tmp/probe-answer.txt)"

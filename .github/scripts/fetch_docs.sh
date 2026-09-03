#!/usr/bin/env bash
# Fetch the Claude Code docs and changelog the digest compares against.
set -euo pipefail
curl -fsSL https://code.claude.com/docs/llms.txt -o /tmp/llms.txt
curl -fsSL https://code.claude.com/docs/en/memory.md -o /tmp/memory.md || \
curl -fsSL https://code.claude.com/docs/en/memory -o /tmp/memory.md
curl -fsSL https://code.claude.com/docs/en/hooks.md -o /tmp/hooks.md || \
curl -fsSL https://code.claude.com/docs/en/hooks -o /tmp/hooks.md
curl -fsSL https://raw.githubusercontent.com/anthropics/claude-code/main/CHANGELOG.md -o /tmp/changelog.md
wc -c /tmp/llms.txt /tmp/memory.md /tmp/hooks.md /tmp/changelog.md

#!/usr/bin/env bash
# Fetch the Claude Code docs and changelog the digest compares against.
set -euo pipefail
curl -fsSL https://code.claude.com/docs/llms.txt -o /tmp/llms.txt
curl -fsSL https://code.claude.com/docs/en/memory.md -o /tmp/memory.md || \
curl -fsSL https://code.claude.com/docs/en/memory -o /tmp/memory.md
curl -fsSL https://code.claude.com/docs/en/hooks.md -o /tmp/hooks.md || \
curl -fsSL https://code.claude.com/docs/en/hooks -o /tmp/hooks.md
curl -fsSL https://code.claude.com/docs/en/agent-teams.md -o /tmp/agent-teams.md || \
curl -fsSL https://code.claude.com/docs/en/agent-teams -o /tmp/agent-teams.md
curl -fsSL https://code.claude.com/docs/en/permission-modes.md -o /tmp/permission-modes.md || \
curl -fsSL https://code.claude.com/docs/en/permission-modes -o /tmp/permission-modes.md
curl -fsSL https://code.claude.com/docs/en/auto-mode-config.md -o /tmp/auto-mode-config.md || \
curl -fsSL https://code.claude.com/docs/en/auto-mode-config -o /tmp/auto-mode-config.md
curl -fsSL https://code.claude.com/docs/en/costs.md -o /tmp/costs.md || \
curl -fsSL https://code.claude.com/docs/en/costs -o /tmp/costs.md
curl -fsSL https://code.claude.com/docs/en/cross-session-messaging.md -o /tmp/cross-session-messaging.md || \
curl -fsSL https://code.claude.com/docs/en/cross-session-messaging -o /tmp/cross-session-messaging.md
curl -fsSL https://code.claude.com/docs/en/remote-control.md -o /tmp/remote-control.md || \
curl -fsSL https://code.claude.com/docs/en/remote-control -o /tmp/remote-control.md
curl -fsSL https://raw.githubusercontent.com/anthropics/claude-code/main/CHANGELOG.md -o /tmp/changelog.md
wc -c /tmp/llms.txt /tmp/memory.md /tmp/hooks.md /tmp/agent-teams.md \
  /tmp/permission-modes.md /tmp/auto-mode-config.md /tmp/costs.md \
  /tmp/cross-session-messaging.md /tmp/remote-control.md /tmp/changelog.md

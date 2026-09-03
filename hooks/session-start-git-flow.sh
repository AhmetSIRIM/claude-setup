#!/bin/bash
# SessionStart hook: name the git mode so the user can switch it in the first message.
# stdout is added to Claude's context (SessionStart is one of the events where it is).
echo "Git mode: ask (default). Say 'plan' or 'free' to switch for this session; see ~/.claude/rules/git.md for what each mode covers."

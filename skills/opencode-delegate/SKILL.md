---
name: opencode-delegate
description: Use when a task is about to run on the opencode CLI (`opencode run`), when the user wants a second opinion from another model family, or when an `opencode run` call hangs or errors, or when choosing between a Claude subagent and opencode for a subtask.
---

# opencode-delegate

Reference for running subtasks on another model through the `opencode` CLI. Claude may
propose delegation and a fitting model, stays the orchestrator either way: it writes the
task, reads the result, verifies it, then reports.

## Running
```bash
cd <target-dir>
opencode run -m <provider/model> "<task>" --format default
```
- Start long calls in the background; if nothing comes back in two minutes, kill the
  process and diagnose.
- Never pass `--auto` (it approves permissions on the model's behalf); delegated work
  stays read-only unless the user explicitly asks for file writes.
- Auth doubts: `opencode auth list`. Logs: `ls -t ~/.local/share/opencode/log | head -1`;
  sessions: `opencode session list`.
- Several models at once: batch by four; concurrent runs share one SQLite state and
  serialize when overloaded.

## Lane choice
| Work | Lane |
|---|---|
| Needs this session's tools, permissions, or a structured return | Claude subagent |
| Noisy discovery whose output would flood the conversation | Claude subagent on a cheap model |
| Second opinion on design or architecture | opencode, strong tier; a different family sees different things |
| Low-risk mechanical bulk: scans, summaries, format conversion | opencode, free or light tier |
| House-style code, files next to secrets, anything needing this conversation's context | stays in the session |

## Model choice
The live list and prices are in `~/.cache/opencode/models.json`; do not trust this table
for numbers.

| Tier | Models | Fit |
|---|---|---|
| Free | `opencode/*-free` (list: `opencode models opencode`) | first try for mechanical work: scans, summaries, format conversion |
| Go plan, light | `opencode-go/mimo-v2.5`, `opencode-go/gpt-5.6-luna`, `opencode-go/minimax-m2.7` | mechanical work when free fails; medium tasks |
| Go plan, strong | `opencode-go/kimi-k3`, `opencode-go/qwen3.8-max`, `opencode-go/grok-4.5` | second opinion on design or architecture; a different family sees different things, so pick one the user has not consulted on this question yet |

The Go plan is prepaid monthly credit that expires unused; spending it is not a cost to
avoid. `opencode/claude-*` models sit outside the plan and return `UnknownError`.

## Known failures
| Symptom | Cause | Fix |
|---|---|---|
| Hangs; log shows `Cannot connect to API` retry every 30 s | the network blocks `opencode.ai` at the TLS layer; opencode never gives up | kill the process; verify with `curl -sS --max-time 6 https://models.opencode.ai/api.json`, switch network |
| Output carries stray CJK characters | free-tier models | clean them when the content is otherwise right |
| `only available hosted in China and requires explicit opt in` | `opencode-go/deepseek-v4-flash` is China-hosted | pick another light model |
| Session opens, no output within the two-minute window, log shows no connect errors | `opencode-go/glm-5.1` / `glm-5.2` | pick another model |
| `UnknownError` | `opencode/claude-*` not in the Go plan | pick a Go-plan model |

Re-test a model with one small request before recommending it; this list goes stale.

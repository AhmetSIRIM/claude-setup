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
opencode run --dir <target-dir> -m <provider/model> "<task>" --format default
```
- Start long calls in the background; if nothing comes back in two minutes, kill the
  process and diagnose.
- Never pass `--auto` (it approves permissions on the model's behalf); delegated work
  stays read-only unless the user explicitly asks for file writes.
- `--variant <effort>` sets provider-specific reasoning effort where the model supports
  it; `--thinking` shows the reasoning blocks.
- Auth doubts: `opencode auth list`. Log:
  `~/.local/share/opencode/log/opencode.log`; sessions: `opencode session list`, run
  from inside the target project (it lists that project's sessions; it has no `--dir`,
  and directories outside a git checkout share one global list).
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
`opencode models` prints what this machine can actually reach: the providers it holds a
credential for, narrowed again by the account's own console settings, so two machines on
the same plan can see different lists. The catalog in `~/.cache/opencode/models.json`
holds every provider models.dev knows, reachable or not; it answers "what does it cost",
never "can I call it". An id the catalog has and `opencode models` does not print fails
at run time with `Unexpected server error` and an `err_...` ref, not with a clear "no
such model", so a stale id reads as an outage.

Read the answer from the CLI at the moment of use, never from a name written here:
```bash
opencode models                       # reachable now
opencode models <provider> --verbose  # adds cost and context metadata per model
opencode models --refresh             # re-pull the catalog from models.dev
```

| Tier | How to recognize it | Fit |
|---|---|---|
| Free | the id ends in `-free` | first try for mechanical work: scans, summaries, format conversion |
| Light | `cost.input` below $1 per 1M | mechanical work when free fails; medium tasks |
| Strong | `cost.input` at $1 per 1M or above | second opinion on design or architecture |

Two providers carry opencode's own models and their ids overlap, so the prefix decides
what is billed. The same console API key unlocks both, but a credential is stored per
provider: `opencode auth login` registers it under one provider at a time, so a Zen
entry alone leaves `opencode-go/*` unreachable. The `OPENCODE_API_KEY` environment
variable is the one route both providers read from a single value.
- `opencode/*` is OpenCode Zen, charged per token against an account balance. It is the
  provider that carries the `-free` ids and the `claude-*` ones.
- `opencode-go/*` is the OpenCode Go subscription, a flat monthly fee with
  dollar-denominated limits per rolling window. Open models only, no `claude-*`.
  Prefer it over Zen for paid work: the fee is already spent, the tokens are not.

For a second opinion pick a family the user has not consulted on this question yet; the
family is the part of the id before the version (`kimi`, `qwen`, `glm`, `grok`,
`minimax`, `gpt`, `gemini`, `deepseek`). The catalog's own `family` field is finer
than that (it splits one vendor by generation), so it does not serve this purpose.

## Known failures
| Symptom | Cause | Fix |
|---|---|---|
| Hangs; log shows `Cannot connect to API` retry every 30 s | the network blocks `opencode.ai` at the TLS layer; opencode never gives up | kill the process; verify with `curl -sS --max-time 6 https://models.opencode.ai/api.json`, switch network |
| Session opens, no output within the two-minute window, log shows no connect errors | the model stalls on the provider side | kill the process; pick another family and say which one stalled |
| `UnknownError` on stderr, body `Unexpected server error` with an `err_...` ref, exit 1 | the id is not in this account's live list, usually one that was retired | `opencode models <provider>`, use an id it prints |
| `Insufficient balance`, with a billing link | a Zen (`opencode/*`) paid model on an empty balance | move to `opencode-go/*` while the plan is active, else a `-free` id, and say the balance is empty |
| Go calls start failing after heavy use | a Go usage window is exhausted and the console's balance fallback is off | wait for the window to reset, or ask the user to enable the fallback |
| `Provider not found: <id>` | no credential for that provider here | `opencode auth list`; only configured providers are usable |
| `only available hosted in China and requires explicit opt in` | the account's China-hosting setting is off | ask the user; it is a console toggle, not a reason to drop the model |
| Output carries stray CJK characters | free-tier models | clean them when the content is otherwise right |

The rows are keyed to symptoms, not to an id list. Re-test a model with one small
request before recommending it; what an account exposes changes under the skill.

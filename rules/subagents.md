# Subagents

## Constraints travel in the dispatch prompt
A Claude Code subagent loads the CLAUDE.md hierarchy (these rules included) and the
session-start git status on its own; those need no restating. Everything else that
binds the task is restated as explicit lines in the dispatch prompt or plan task: a
standing decision from the conversation, a frozen wire contract, a naming decision, a
constraint recorded only in memory. Never assume a subagent inherits the conversation
or reads memory on its own. A plan whose tasks will run in subagents states the shared
constraints once in the plan itself, and each dispatch copies the ones that apply.

Two kinds of delegate still get the full treatment, every constraint restated: the
built-in Explore and Plan agents (they skip CLAUDE.md and rules for speed), and
anything running outside Claude Code (an opencode delegate loads no rules at all).

Scenario this prevents: a subagent renames a frozen constant, drops required semantics,
or hardcodes an LTR layout because the constraint lived only in memory or in the
conversation, and the violation surfaces late, in review or in production.

Source: code.claude.com/docs/en/sub-agents, subagent context behavior. Harness
mechanics; re-verify against the docs before leaning on it in a new harness version.

## The return is a contract, and a claim is not proof
A dispatch names the shape of the answer it expects: verdict or summary first, then
the files touched as absolute paths, then open questions; raw tool output never
travels back. The dispatch also carries its acceptance criteria, and the orchestrator
verifies the result with its own reads and checks, in proportion to the dispatch's
write scope: a read-only research dispatch earns a spot-check of the claims it cites,
a write dispatch earns opening the files it names. A subagent saying "completed" is a
signal, not evidence; the same discipline external-systems.md applies to services
("A success response is not proof") applies to delegates.

Scenario this prevents: a subagent reports done over a file that was never written,
or pastes its whole transcript back and floods the orchestrator's context.

Lane choice (which work goes to a Claude subagent, which to opencode, which stays in
the session) lives in the opencode-delegate skill.

## An agent file is earned, not planned
A custom agent definition is written after the same delegation has been dispatched by
hand three times. Speculative agents die unused; the ones that survive absorb noise a
real task keeps producing.

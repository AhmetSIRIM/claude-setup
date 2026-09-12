# Session hygiene

## Session files have a home and a lifetime
A file that is meaningless without its session (scratch script, draft, probe output,
intermediate dump, backup copy) is session-scoped. Session-scoped files are written to
the session's scratch directory (the scratchpad or job tmp dir the harness provides;
`/tmp` as fallback), never into the repository or the home directory.

When a session-scoped file must live inside the repository (a tool resolves it by
relative path, the build must see it), it is an exception: name it and its removal
point out loud when creating it, keep it on a running list, and delete everything on
that list at session end or when the user asks for a sweep. A session never ends with
an undeclared session-scoped file in the repository or in HOME.

Promotion is explicit: a scratch file that turns out to be worth keeping is moved to a
permanent location and named for its role, in the open, not left where it was born.

Why: an orphaned scratch file outlives the only context that explained it. Weeks later
it reads as possibly load-bearing, so nobody deletes it, and the repository and HOME
silt up with files no one can account for.

## A session carries one piece of work
A session is scoped to one piece of work: a feature, an investigation, a study topic.
Claude keeps it that way from its side:

- When the conversation turns to an unrelated subject, Claude says so and proposes a
  new session (or `/clear`) before continuing; the user decides.
- When the piece of work is done, Claude says so and notes that the session can be
  closed; `/resume` brings it back with its summary.

Each proposal names this rule and offers to read the Claude Code doc behind it on
request, so the user weighs a sourced suggestion rather than a passing remark. The
form is one sentence, once; the user's answer settles it for that session.
Example: "The subject changed; per the session rule, a new session fits better. Want
me to check the docs on why?"

Why: context that outlives its task turns into noise. Claude Code's own guidance is
that too much context makes Claude less effective, that skills may stop triggering and
conventions get lost, and that compaction keeps decisions but drops detail. A finished
piece of work is cheaper to rebuild from its commits than to carry, and a session left
open across days re-reads its whole history at full price once the prompt cache has
expired.

Source: Claude Code docs, "Sub-agents" (context noise), "Explore the context window"
(compaction), "Manage costs effectively" and "How Claude Code uses prompt caching"
(cache lifetime), "Manage sessions" (`/resume`); Anthropic, "Claude Code best
practices" (multiple sessions per piece of work).

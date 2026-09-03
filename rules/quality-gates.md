# Quality gates

## A red gate is fixed, never bypassed
A failing quality gate (pre-commit hook, linter, type check, test suite, CI check) is
fixed at its cause. Any mechanism that silences a gate instead of fixing it needs the
user's explicit yes in that same turn, and the ask names the exact gate, the silencing
mechanism, and the reason. Silencing mechanisms include: `git commit --no-verify` or
otherwise skipping a hook, suppression annotations (`@Suppress`, `@ts-ignore`,
`// nolint` and equivalents), adding a finding to a lint baseline, disabling or skipping
a test, and lowering a gate's severity. Time pressure from the user is not that yes; the
ask still happens.

This holds in every git mode: `plan` and `free` delegate the commit ask, never the
bypass ask. A project may sanction a specific mechanism (for example, a lint baseline
during legacy adoption) in its own CLAUDE.md; additions to it are still named in the
turn they happen.

Why: a bypassed gate ships the exact defect the gate exists to catch, and each silent
suppression normalizes the next one, until the pipeline is green only because it no
longer looks.

Source: harperreed/dotfiles and obra/dotfiles, .claude/CLAUDE.md (forbidden flags such
as --no-verify; never skip, evade, or disable a pre-commit hook). Borrowed convention,
extended to suppression annotations and lint baselines.

# Communication

## Memory writes ask first
Nothing is written to auto memory (create, update, or delete) without the user's
approval. A proposal states three things in one short block: the intent (what would be
saved), the reason (what correction or fact triggered it), and the expected outcome (how
future sessions behave differently). An explicit "save this to memory" from the user is
already approval.

Why: unapproved writes accumulate facts that go stale, and the user wants to see the
reasoning behind each memory so that their model of the work and Claude's stay in sync.

## No em dashes
Nothing written contains an em dash: responses, code comments, docs, commit messages,
issue and PR text. Use commas, parentheses, semicolons, or separate sentences instead.

## Assumptions are visible
No silent assumptions. When information is missing: ask, or request permission to
assume. Once permission is granted, every assumption is marked in place with
`// TODO (Assumption): <the assumption, self-contained>`. This covers code, logic,
requirements, and implementation details.

Why: review depends on knowing exactly which parts followed instructions as-is and
which parts were guessed.

## Plan files are scaffolds
A plan holds intent, decisions, and acceptance criteria; no code listings. Deviate
consciously when execution surfaces better options; never silently.

Why: code inside a plan is a coupled copy of the code to come; the first deviation in
execution turns it stale and starts drift between plan sections. Intent and acceptance
criteria survive change; listings do not.

## Volatile values are read at use time, never recalled
A value that can drift (a model id, a partner key, a config constant, a branch-dependent or environment-dependent setting, a number kept on a single source-of-truth page) is read from its source of truth at the moment it is used: the repo's config, the current branch, the SSOT page. Memory and prior conversations may record where such a value lives, never what it currently is. When a value differs per branch or environment, the branch or environment checked is named alongside the value.

Why: a recalled value looks as authoritative as a fresh one, so acting on it fails silently; the config changed, the checked-out branch differs, or the source page was updated after the value was memorized. The docs rule keeps stale facts out of written artifacts; this rule keeps them out of decisions at read time.

## A memory is written timeless
Auto memory is for learnings, corrections, and project context the code cannot carry;
this rule keeps that mission and narrows only the form.

- An entry records where truth lives and facts that do not drift: a pointer to the
  source, a decision with its why, an absolute date.
- It never stores a value that can change under it: a price, a version, a count, an
  environment-dependent setting.
- It never references a session ("as discussed today"); the reader has neither.
- Exception: in-progress state with no other home may be written. It carries an
  absolute date and is treated as expired by default; verify before acting on it.
- A write proposal (see "Memory writes ask first") says how the entry meets this rule.

Scenario this prevents: a remembered value reads as authoritative long after the
config, price, or plan changed under it, and the session acts on the stale copy.

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

## Decisions are asked as choices
A decision that needs the user's answer is asked with `AskUserQuestion`: the
recommended option first, each option naming its trade-off in one sentence. A
plain-text question is for an open request that does not fit into options.

Why: the options show the axis the agent is reasoning on, so the user sees it and can
accept, shift, or reframe it; the two models of the work stay in sync. With
`askUserQuestionTimeout` set, an unanswered question lets the agent go on with its own
judgment; a plain-text question holds the session until someone answers.

Source: Claude Code docs, "Tools reference" (AskUserQuestion tool behavior; Question
auto-continue timeout). The rule is the owner's decision on top of that mechanism.

## Plan files are scaffolds
A plan holds intent, decisions, and acceptance criteria; no code listings. Deviate
consciously when execution surfaces better options; never silently.

Why: code inside a plan is a coupled copy of the code to come; the first deviation in
execution turns it stale and starts drift between plan sections. Intent and acceptance
criteria survive change; listings do not.

## No time forecasts
Claude never forecasts duration or dates: no hour or day estimates, no deadlines, no
week or sprint plans, no "this takes ten minutes". Size is stated in what can be
counted and checked: the files and steps a change touches, the unknowns it depends
on, the risk if a step fails. A plan orders its work by dependency, not by calendar.
A measured duration (a build that took four minutes) or a sourced ratio (a documented
slowdown) is a fact, not a forecast, and may be reported.

Effort is still sized. The ban covers the unit, not the judgment: a four-line edit is
called a four-line edit, and a small change is never declined as "not worth the
complexity" without naming what the complexity is.

Why: models learned duration from human estimates made before agents, so they
over-predict short tasks several times over and give the same guess whatever the
size; a forecast reads as authoritative and anchors the plan. A ban on estimates
alone pushes the other way: without effort-sizing, the agent inflates small tasks
into "complex" ones and declines them.

Source: Ofengenden and Andriushchenko, "Your Agents Are Not Time Aware" (LessWrong);
anthropics/claude-code issue #20270 (the effort-sizing side effect). Borrowed
evidence; the rule is the owner's decision.

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

# Git

## Git mode
Every session runs in one of three modes. `ask` is the default. The user switches mode for
the current session by naming it; a mode never carries over to the next session, never
comes from a skill, and never comes from an approval given in a past conversation.

| Mode   | Commit                                        | Push          | PR create / merge / release |
|--------|-----------------------------------------------|---------------|-------------------|
| `ask`  | ask for an explicit yes, in a separate turn   | ask           | ask               |
| `plan` | tasks of the agreed plan commit without asking; anything outside the plan asks | ask | ask     |
| `free` | commit without asking                         | push without asking | ask         |

"Ask" means: show the final diff, message, files and target, then wait for a yes in the
next turn. "Commit it" said while approving the approach is not that yes; the finished
work is what gets approved. PR create, merge, and release ask in every mode.

Before any push, PR, or release, in every mode, the controlling agent verifies first
(tests, diff review, acceptance) and only then asks; an approval given before
verification is uninformed.

Why: each write is a review checkpoint, and the user decides per session how much of it
to delegate. `free` exists for flows such as a push-triggered CI run; PR and release
stay gated because they are the outward-facing, hard-to-undo steps.

## Cadence, when a plan exists
- One atomic commit per plan task, after the user confirms understanding of that task.
  One commit carries one self-contained change: a fix discovered mid-task, or a
  refactoring next to a feature, gets its own commit rather than riding along.
- Push at plan boundaries (end of a phase or milestone), not after every commit.
- Why: `git log` reads as one concept per entry, so revert, bisect and cherry-pick work at
  concept granularity; a reviewer understands a change faster when it is not mixed with a
  refactoring. Ad-hoc work without a plan is not bound by the cadence.

Source for the one-change-per-commit framing: Google Engineering Practices, "Small CLs",
sections What is Small and Separate Out Refactorings. Borrowed convention.

## History stands alone
Everything recorded into history (commit message, PR title and body, tag, branch name)
describes the change by its effect on the product or codebase and is understandable with
zero conversation context: no plan codenames, no task numbers, no references to sessions,
drafts, or internal documents. Branch names follow the same rule because they embed in
merge commits.

Per-repository conventions (subject prefix, ticket-id format, whether commit bodies carry
a teaching note, conventional-commits or free-form) live in that project's CLAUDE.md, not
here.

Why: history is permanent and read by people who never saw the plan.

## Every repository is versioned; a release is a changelog entry published on the host
- Every repository ships versions unless stated otherwise or a special case applies.
- Semantic Versioning is the default scheme.
- A versioned repository keeps `CHANGELOG.md` in the Keep a Changelog format; that file
  is the single source of release notes.
- A release is the host's published release object (a GitHub Release, or what the host
  offers in its place). It is tagged with the version, `vX.Y.Z` by default, and titled
  the same; its body is that version's `CHANGELOG.md` section, copied as is. The
  changelog section is written and committed first; release notes are never drafted on
  the host and copied back. A tag alone is not a release. On a host with no release
  object, ask the user what stands in for it before tagging.

Why: versions and a changelog carry a project through long-term use and maintenance;
without them, what changed and what an upgrade needs is dug out of history each time.

Source: Keep a Changelog (https://keepachangelog.com/en/1.1.0/).

## A worktree is proposed, never entered unannounced
Work happens in the checkout the user opened. A worktree is entered only when the user
asks for one, or when Claude proposes it with the reason (parallel edits to the same
repository, a long change the user wants kept off the working copy) and the user says
yes. Where a repository's own settings keep automatic isolation on, the first message
of the session says that a worktree was entered and where it lives.

Why: a worktree opened without notice moves the work to a path the user's IDE does not
show; nested inside the project, JetBrains reads it as a multi-root project and its
Git integration breaks. The isolation is worth that cost only when the user knows it is
there.

Source: JetBrains, "Use Git worktrees" (avoid nesting a worktree inside the project
directory); the proposal rule is the user's decision.

## Write discipline
- Every git write names its repository: `git -C <absolute path> ...`. The shell working
  directory resets between tool calls; with worktrees, two checkouts of the same
  repository sit side by side and a bare `git commit` or `git commit --amend` can land in
  the wrong one.
- Before merging, read the branch content from the remote (`git ls-tree origin/<branch>`
  or the host's contents API), not from the local tree. A green pipeline proves that what
  was pushed passes; it does not prove that what was pushed is what was written.
- After a deploy, prove on the host that the new code is running: an embedded string,
  build id, or schema version. A passing health check is also passed by the old binary.
- Never run `git checkout <file>` on uncommitted work. It discards the file's changes and
  leaves no reflog entry, so unlike most git mistakes it cannot be recovered.

Scenario these prevent: an amend meant for a worktree runs in the main checkout and
rewrites the wrong commit; a force-push then replaces the pull request's real commit; the
pipeline stays green because the replacement is documentation-only; the merge ships no
code; the release tag deploys the old binary; nothing surfaces until someone inspects the
host by hand. Recovery costs a force-push to a shared branch, a deleted tag and release,
and a reopened issue.

## Remote branch delete or rename consults the host first
Before deleting or renaming a remote branch, query the hosting service's pull request
API (Azure DevOps, GitHub, or whichever host owns the repository) for PRs where the
branch is the source and for PRs where it is the target, with status filters covering
both open and closed PRs; default listings often show only active ones. Git data alone
is not evidence: a branch merged into every other branch can still back an active PR, a
branch can be the target of someone else's open PR with no local trace, deleting a
source branch orphans its PR instead of closing it, and renaming a source branch
destroys the PR outright.

Before deciding, also compute which commits only that branch holds (`git rev-list`
against all other refs), so the decision is made with the full loss in view. After the
operation, re-read the affected PRs' status on the host instead of assuming the outcome.

Why: PR state lives on the host, not in git. A deleted branch is usually recoverable; a
dead PR's review history and discussion are not. This keeps routine branch cleanup from
silently orphaning or killing team-visible reviews.

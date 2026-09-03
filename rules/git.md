# Git

## Git mode
Every session runs in one of three modes. `ask` is the default. The user switches mode for
the current session by naming it; a mode never carries over to the next session, never
comes from a skill, and never comes from an approval given in a past conversation.

| Mode   | Commit                                        | Push          | PR create / merge |
|--------|-----------------------------------------------|---------------|-------------------|
| `ask`  | ask for an explicit yes, in a separate turn   | ask           | ask               |
| `plan` | tasks of the agreed plan commit without asking; anything outside the plan asks | ask | ask     |
| `free` | commit without asking                         | push without asking | ask         |

"Ask" means: show the final diff, message, files and target, then wait for a yes in the
next turn. "Commit it" said while approving the approach is not that yes; the finished
work is what gets approved. PR create and merge ask in every mode.

Before any push or PR, in every mode, the controlling agent verifies first (tests, diff
review, acceptance) and only then asks; an approval given before verification is
uninformed.

Why: each write is a review checkpoint, and the user decides per session how much of it
to delegate. `free` exists for flows such as a push-triggered CI run; PR stays gated
because it is the outward-facing, hard-to-undo step.

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

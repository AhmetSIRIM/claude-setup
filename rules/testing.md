# Testing

## A red test is fixed or surfaced, never removed
A failing test is never deleted, disabled (@Disabled, @Ignore, commented out, skipped),
or weakened (assertion loosened, expected value rewritten to match current output) in
order to make the suite green. Three exits exist: fix the code under test; fix a test
that is genuinely wrong, stating in the change why the old expectation was wrong; or
stop and surface the failure to the user. Quarantining a flaky test is a user decision,
made after the failure is surfaced, and it carries a tracked reason in the change.

Reducing test coverage is treated as worse than a failing test. A red test is visible
evidence of a defect; a deleted or loosened test converts that caught regression into a
shipped one, and the loss hides in review because the diff reads as cleanup.

Source: Jesse Vincent, "That time it tried to delete all my tests"
(https://blog.fsck.com/2026/04/30/that-time-it-tried-to-delete-all-my-tests/);
obra/dotfiles .claude/CLAUDE.md (reducing test coverage is worse than failing tests).
Borrowed convention.

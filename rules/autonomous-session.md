# Autonomous session

## Bypass permissions means nobody is at the keyboard
A session whose permission mode is `bypassPermissions` is an autonomous run: the user
accepted every action in advance and may not be watching. In that session no rule
blocks on the user's answer; a question nobody can answer becomes a decision and a line
in the final report.

This section overrides every ask in the other rules. Wherever a rule says to wait for
the user (a git mode default, a memory write proposal, a plan approval, a "confirm
first"), this section says what happens instead, and the other rules do not repeat the
exception. A new rule that adds an ask does not need to mention this file; the override
already covers it.

- Git mode starts as `free`: commits and pushes go without asking. The user still
  switches it by naming `ask` or `plan`.
- PR create, merge, and release are neither done nor asked. The work is verified and
  left ready, and the final report says which of them wait for a yes.
- A memory write is not proposed in chat; the proposal goes into the final report.
- Missing information does not stop the work: the assumption is marked in place with
  `// TODO (Assumption): ...` and listed in the report.
- A destructive step the task did not ask for is skipped and reported, never taken
  because nobody objected.
- A red quality gate is never silenced and a red test is never removed, disabled, or
  weakened. Those two asks stay closed here as they do in every mode: the run stops
  that path, or goes on around it with the failure left in place, and the report
  names the gate or the test.

The mode is learned from the prompt hook's line ("Autonomous session: bypass
permissions is on ...") or from the user saying so. The hook prints the line when the
session's permission mode is first seen and again when it changes, so a Shift+Tab
switch into bypass is announced on the next prompt. `--allow-dangerously-skip-permissions`
alone does not make a session autonomous: it only puts bypass in the Shift+Tab cycle,
and the session starts in its usual mode.

Why: bypass is the user's declaration of full trust for one session. A commit ask after
it is a checkpoint nobody will answer, and a blocked autonomous run wastes the session.
The gates that stay closed (PR, merge, release) are the outward-facing steps whose undo
needs another person; the docs also note that bypass has no protection against prompt
injection, so the step that publishes stays with a human.

Source: Claude Code docs, "Choose a permission mode" (the bypassPermissions section and
what puts the mode in the cycle) and "Hooks" (`permission_mode` in the UserPromptSubmit
input; the SessionStart input does not carry it, and no event fires on a mode change).
The behavior list is the owner's decision.

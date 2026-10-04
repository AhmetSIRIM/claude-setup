# Autonomous session

## The owner's word starts an autonomous run
An autonomous run starts when the owner says the session should go on without them
(they are leaving, "run it overnight", "carry on on your own") and ends when they are
back or say so. Nobody is at the keyboard: the owner accepted the run's actions in
advance and may not be watching. In that run no rule blocks on the owner's answer; a
question nobody can answer becomes a decision and a line in the final report. The
owner is the person the other rules call the user.

When the owner declares a run, the session says before they leave which steps will wait
for them (PR create, merge, release) and asks once whether this run may spawn a
teammate or start a session outside an approved team plan, or with no plan at all; an
unanswered question keeps that gate closed.

This section overrides every ask in the other rules. Wherever a rule says to wait for
the owner (a git mode default, a memory write proposal, an implementation plan approval,
a "confirm first"), this section says what happens instead, and the other rules do not
repeat the exception. A new rule that adds an ask does not need to mention this file;
the override already covers it.

- Git mode starts as `free`: commits and pushes go without asking. The owner still
  switches it by naming `ask` or `plan`. The owner's approval of an autonomous team
  plan declares the run for every teammate in that plan, so a teammate's git mode
  starts as `free` too.
- PR create, merge, and release are neither done nor asked. The work is verified and
  left ready, and the final report says which of them wait for a yes.
- Unless the owner opened it when declaring the run, spawning a teammate or starting a
  session outside an approved team plan, or with no plan at all, is neither done nor
  asked. The final report names it and the work that waits for it.
- A memory write is not proposed in chat; the proposal goes into the final report.
- Missing information does not stop the work: the assumption is marked in place with
  `// TODO (Assumption): ...` and listed in the report.
- A destructive step the task did not ask for is skipped and reported, never taken
  because nobody objected.
- A red quality gate is never silenced and a red test is never removed, disabled, or
  weakened. Those two asks stay closed here as they do in every mode: the run stops
  that path, or goes on around it with the failure left in place, and the report
  names the gate or the test.

The run happens in auto mode. Its classifier blocks actions that go beyond the request
(force push, discarding uncommitted work, merging an unapproved pull request, sending
secrets out), and the run goes on after a block; a block is reported, never worked
around. Repeated blocks pause auto mode until someone approves, so the report names the
work left waiting. The prompt hook's line names the permission mode; when the owner
declares a run and that mode is not auto, the session says so before they leave, since
every permission prompt would wait for them. When no hook line is left in context, as
after a compaction, the session asks the owner which mode the status bar shows. Auto
mode still asks in a few cases, such as the first read outside the working
directories, so the session names the directories the run will need and the owner adds
them with `/add-dir` before leaving.

Why: the owner's word is their declaration of trust for one run. A commit ask after it
is a checkpoint nobody will answer, and a blocked run wastes the time it was given.
The gates that stay closed (PR, merge, release, a teammate or session outside the
approved team plan) are the outward-facing or trust-widening steps whose undo needs the
owner. Auto mode is the starting permission mode, so the mode cannot tell an attended
session from an unattended one; only the owner's word can. Its classifier also turns
several of these rules' prohibitions into blocks the harness enforces, a floor under the
prompt rules.

Source: Claude Code docs, "Choose a permission mode" (auto mode as the built-in starting
mode, what the classifier blocks, the repeated-block fallback); "Orchestrate teams of
Claude Code sessions", Messages between agents, and "Cross-session messaging", How a
session treats an incoming message (a message from another agent is never consent).
The behavior list is the owner's decision.

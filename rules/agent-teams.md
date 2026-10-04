# Agent teams

Applies when work runs as a Claude Code agent team: a lead session the owner talks to,
and teammates it spawns. Delegation inside one repository stays with subagents; a team
fits work across several repositories. The first section applies to every session,
because the setting is on everywhere.

## Agent teams are on in user settings
Agent teams are switched on in user settings, so a team can start from any folder.
While they are on, a subagent that Claude names launches as a teammate; that cost is
accepted, and outside the teammates the owner approved an Agent call carries no
`name`, so the delegation stays a subagent. A team that touches several repositories
works best with its lead opened in the folder that holds them, because a lead opened
inside one repository gives its teammates only that repository's instructions.

Scenario this prevents: a team the owner asks for cannot start because the folder it
was asked from has no settings file switching teams on; a named delegation launches a
teammate nobody approved.

Source: Claude Code docs, "Orchestrate teams of Claude Code sessions" (a named subagent
launches as a teammate while teams are on) and "How Claude remembers your project" (a
subdirectory's CLAUDE.md loads on demand); a teammate of a lead in a parent folder
loads a sub-repository's CLAUDE.md, confirmed in practice.

## The owner approves the teammates; each is named by scope and job
The lead names the teammates it wants and the owner approves them before any spawn;
that approval is the team plan the other rules refer to. A teammate is named
`<scope>.<job>`: the repository or file set it writes, then its job
(`whaletales-ios.writer`). Fields are joined by `.` and use lowercase letters, digits,
and `-` only; the name is fixed at spawn and never names a human title or a person.
Each review is a new teammate, so a reviewer's job carries what sets it apart
(`whaletales-server.reviewer-t3`).

Scenario this prevents: the lead launches teammates the owner never saw (Claude Code
does not ask before a launch); a panel that shows only names hides which repository
each teammate writes.

## One writer per repository
Each teammate writes only the repository or files the owner approved it for
(subagents.md, "One writer per scope").

Scenario this prevents: two teammates overwrite each other's files.

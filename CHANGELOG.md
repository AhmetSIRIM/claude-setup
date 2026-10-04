# Changelog

## [Unreleased]

### Added
- `rules/subagents.md`, "A dispatch says where it stops": every dispatch names its
  stop condition, and a return ends with one of four statuses.
- `rules/subagents.md`, "One writer per scope": parallel work is split into separate
  scopes, never two writers on one.
- `rules/subagents.md`, "The reviewer starts clean": a reviewer starts fresh from the
  diff and the acceptance criteria, never from the author's rationale.
- `rules/subagents.md`, "The model is named by role": every dispatch names its model
  (`haiku`, `sonnet`, `opus`, or `fable`, each with its role), and the verifying role
  never runs on a weaker tier than the work it judges.
- `rules/agent-teams.md`: rules for a Claude Code agent team; agent teams are on in
  user settings, a lead for several repositories works best opened in the folder that
  holds them, a delegation outside the approved teammates carries no `name`, the owner
  approves the teammates before any spawn, a teammate is named `<scope>.<job>`, and
  each teammate writes one repository.

### Changed
- `rules/autonomous-session.md`: an autonomous run starts on the owner's word, not on a
  permission mode, and runs in auto mode, whose classifier blocks actions beyond the
  request; the session names the steps that will wait for the owner before they leave
  and the directories the run will need, asks for the mode when no hook line is left in
  context, and the setup no longer uses `bypassPermissions`. Spawning a teammate or
  starting a session outside an approved team plan, or with no plan, is a closed gate
  in an autonomous run, reported instead of done, unless the owner opens it when
  declaring the run; the session asks once, and an unanswered question keeps it
  closed. The owner's approval of an autonomous team plan declares the run, and git
  mode `free`, for every teammate in it.
- `cmd/announce-git-mode`: names the permission mode and git mode ask in every mode and
  points at the autonomous rule, where git mode is `free`; no mode makes a session
  autonomous.
- README, "A session that runs on its own": replaces the section on starting a session
  without permission checks.
- `rules/subagents.md`: a delegate outside Claude Code gets a prompt written for a
  reader with no context, and a return from another model family is weighed finding
  by finding.
- `rules/subagents.md`: drops the pointer to the opencode-delegate skill; the rules
  stand on their own and name no skill.
- `settings.template.json`: switches agent teams on in user settings with
  `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS`.

### Removed
- README, "Agent teams considered and held": work across several repositories needs
  more than one session's subagents, and the experimental label and the cost are
  accepted, so agent teams are now in use under `rules/agent-teams.md`.

### Hand steps
- Add `"CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS": "1"` to the `env` block of
  `~/.claude/settings.json` on every machine, then restart open sessions.
- Run `go install ./cmd/announce-git-mode` on every machine, so the hook stops treating
  bypass as an autonomous run.

## [0.4.0] - 2026-10-04

### Added
- `skills/driving-apps-with-maestro`: rules for driving a mobile app with Maestro on a
  simulator; reading the screen tree, selector semantics (including a `$` in an `id`),
  waiting, commands that pass without doing anything, iOS traps, running flows, driver
  processes left behind even after passing runs, the MCP and CLI split, and a simulator
  window the user can watch by default.
- `rules/communication.md`, "Decisions are asked as choices": a decision that needs
  the user's answer goes through `AskUserQuestion` with the recommended option first,
  so the axis the agent reasons on is visible and an unanswered question can
  auto-continue.
- `rules/communication.md`, "No time forecasts": no duration, deadline, or calendar
  estimates; size is stated in countable scope, and effort is still sized so small
  changes are not declined as complex.

### Fixed
- `cmd/announce-git-mode`: the per-session state lives under
  `claude-setup.announce-git-mode` in the temp dir. A file named plain
  `announce-git-mode` there made every state write fail with "not a directory", and the
  git mode line repeated on every prompt.

## [0.3.0] - 2026-09-24

### Added
- `cmd/budget-check`: reads the Go plan usage windows before the weekly digest, the
  repository's first Go tool (`go.mod` at the root).
- `cmd/announce-git-mode`: the git mode hook as a Go binary on UserPromptSubmit, the
  event whose input carries `permission_mode` (SessionStart does not). It announces
  the git mode when a session's permission mode is first seen and again when it
  changes, with the autonomous-session warning when the mode is `bypassPermissions`.
- `skills/learn`: `/learn <topic>` sets a tutoring contract for the session; one source
  at a time, the learner chooses how to work, Claude checks and gives graduated hints.
- `rules/git.md`, "A worktree is proposed, never entered unannounced".
- `rules/session-hygiene.md`, "A session carries one piece of work".
- `tools/`: ccstatusline pinned with a lockfile, its config, and a Dependabot entry.
- `README.md`, "New machine": the Terminal.app Option-as-Meta switch and agent
  view as the opening screen, two preferences that live outside the repo.
- `README.md`, "A session without permission checks": `bypassPermissions` enters the
  mode cycle only at launch, and the `--allow-` flag puts it there without activating it.
- `rules/autonomous-session.md`: a session in `bypassPermissions` mode is an autonomous
  run; git mode starts as `free`, no rule blocks on the user, and PR, merge and release
  wait in the final report. The file overrides every ask in the other rules, so they
  do not repeat the exception.

### Changed
- `doc-drift-check.yml`: a rate-limited Go plan window marks the review job skipped,
  with the reset time in the job summary, instead of failing the run.
- `settings.template.json`: `worktree.bgIsolation` is `none`; background sessions work
  in the checkout unless a worktree is asked for.
- `rules/git.md`: a release is the host's published release object, cut from a
  changelog section written first; a tag alone is not one, and a host without a
  release object gets asked. Cutting a release asks in every git mode, like PR create
  and merge.
- `settings.template.json`: `showThinkingSummaries` is on; thinking summaries show in
  the transcript view.
- `skills/opencode-delegate`: the model table says how to recognize a tier (the `-free`
  suffix, `cost.input`) instead of listing ids that go stale under the skill.
  `opencode models` is the reachability source and depends on the account's console
  settings; `models.json` is the price source and lists models no credential reaches.
  The Zen and Go providers are separated, and the failure rows are keyed to symptoms.
- `doc-drift-check.yml` and `README.md`: `OPENCODE_GO_KEY` is documented as an API key
  from the OpenCode console, which serves both the Zen and the Go provider, instead of
  a field of a credential file that exists only on one machine.
- `settings.template.json`: `model` is `opus`, which runs with the 1M window natively
  on the Anthropic API.
- `settings.template.json`: the status line runs ccstatusline from `tools/`.
- `hooks/session-start-doctor.sh`: warns when the installed ccstatusline version
  differs from the pin in `tools/package.json`.

### Removed
- `.github/scripts/probe_provider.sh`, superseded by `cmd/budget-check`.
- The model dropdown of `doc-drift-check.yml`; the model is set in the workflow's
  `env`.
- `hooks/statusline.sh` and its license, replaced by ccstatusline.
- `hooks/session-start-git-flow.sh`, replaced by `cmd/announce-git-mode`.
- `effortLevel` from `settings.template.json`; Opus 5.5 and later ignore it.
- Windows support: the junction, Git Bash and winget steps in the README, the release
  notice's Windows mention, and the doctor's search for a Windows python launcher. The
  setup targets macOS; the removed steps stay in history for a future port.

### Hand steps
- Add `"worktree": {"bgIsolation": "none"}` to `~/.claude/settings.json` on every
  machine (the template carries it for new installs).
- Set `"model": "opus"` in `~/.claude/settings.json` on every machine.
- Install `node`, run `npm ci --prefix tools`, then copy `statusLine` from the
  template and drop `effortLevel` in `~/.claude/settings.json`, on every machine.
- `opencode auth login -p opencode-go -m api` on every machine, with the OpenCode
  console API key; a Zen credential alone does not reach `opencode-go/*`.
- Terminal.app: "Use Option as Meta: Left Option only" on the default profile, and
  agent view as the opening screen, on every machine (README, "New machine").
- `brew install go`, `go install ./cmd/announce-git-mode`, then in
  `~/.claude/settings.json` drop the `session-start-git-flow.sh` entry and add a
  `UserPromptSubmit` hook running `~/go/bin/announce-git-mode` (as in the template),
  on every machine.

## [0.2.0] - 2026-09-06

### Added
- `rules/scripting.md`: shell stays inside the Google Shell Style Guide's boundary, Go
  beyond it.
- `rules/git.md`, "Every repository is versioned; a release is a changelog entry":
  `CHANGELOG.md` is the single source of release notes.

### Hand steps
- Install `go` and `shellcheck` on every machine.

## [0.1.0] - 2026-09-03

### Added
- Initial tracked setup: rules, skills, hooks, status line, and the weekly doc-drift
  digest.

[Unreleased]: https://github.com/AhmetSIRIM/claude-setup/compare/v0.4.0...HEAD
[0.4.0]: https://github.com/AhmetSIRIM/claude-setup/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/AhmetSIRIM/claude-setup/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/AhmetSIRIM/claude-setup/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/AhmetSIRIM/claude-setup/releases/tag/v0.1.0

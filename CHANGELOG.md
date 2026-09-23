# Changelog

## [Unreleased]

### Added
- `cmd/budget-check`: reads the Go plan usage windows before the weekly digest, the
  repository's first Go tool (`go.mod` at the root).
- `skills/learn`: `/learn <topic>` sets a tutoring contract for the session; one source
  at a time, the learner chooses how to work, Claude checks and gives graduated hints.
- `rules/git.md`, "A worktree is proposed, never entered unannounced".
- `rules/session-hygiene.md`, "A session carries one piece of work".
- `tools/`: ccstatusline pinned with a lockfile, its config, and a Dependabot entry.
- `README.md`, "New machine": the Terminal.app Option-as-Meta switch and agent
  view as the opening screen, two preferences that live outside the repo.
- `README.md`, "A session without permission checks": `bypassPermissions` enters the
  mode cycle only at launch, and the `--allow-` flag puts it there without activating it.

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

[Unreleased]: https://github.com/AhmetSIRIM/claude-setup/compare/v0.2.0...HEAD
[0.2.0]: https://github.com/AhmetSIRIM/claude-setup/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/AhmetSIRIM/claude-setup/releases/tag/v0.1.0

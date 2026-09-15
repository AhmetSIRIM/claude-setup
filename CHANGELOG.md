# Changelog

## [Unreleased]

### Added
- `cmd/budget-check`: reads the Go plan usage windows before the weekly digest, the
  repository's first Go tool (`go.mod` at the root).
- `skills/learn`: `/learn <topic>` sets a tutoring contract for the session; one source
  at a time, the learner chooses how to work, Claude checks and gives graduated hints.
- `rules/git.md`, "A worktree is proposed, never entered unannounced".
- `rules/session-hygiene.md`, "A session carries one piece of work".

### Changed
- `doc-drift-check.yml`: a rate-limited Go plan window marks the review job skipped,
  with the reset time in the job summary, instead of failing the run.
- `settings.template.json`: `worktree.bgIsolation` is `none`; background sessions work
  in the checkout unless a worktree is asked for.
- `rules/git.md`: a release is the host's published release object, cut from a
  changelog section written first; a tag alone is not one, and a host without a
  release object gets asked. Cutting a release asks in every git mode, like PR create
  and merge.

### Removed
- `.github/scripts/probe_provider.sh`, superseded by `cmd/budget-check`.
- The model dropdown of `doc-drift-check.yml`; the model is set in the workflow's
  `env`.

### Hand steps
- Add `"worktree": {"bgIsolation": "none"}` to `~/.claude/settings.json` on every
  machine (the template carries it for new installs).

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

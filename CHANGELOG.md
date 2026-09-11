# Changelog

## [Unreleased]

### Removed
- The model dropdown of `doc-drift-check.yml`; the model is set in the workflow's
  `env`.

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

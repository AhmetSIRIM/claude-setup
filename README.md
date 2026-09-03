# claude-setup

Personal Claude Code setup, tracked so it survives a machine change and stays
identical across machines.

> This setup carries fixes for problems I hit in my daily Claude Code work, and
> guards against failure modes the wider community keeps reporting:
>
> - working rules that load every session,
> - a doctor hook that catches silent breakage,
> - templates that ship nothing personal.
>
> The setup also maintains itself: a
> [weekly digest](.github/workflows/doc-drift-check.yml) watches the Claude Code docs
> and changelog for drift.

## What lives where

| Path in this repo | Wiring |
|---|---|
| `CLAUDE.template.md` | copy to `~/.claude/CLAUDE.md`, fill; stays local |
| `rules/*.md` | symlinked as `~/.claude/rules/`; `kotlin.md` carries `paths:` scoping |
| `skills/*/SKILL.md` | symlinked as `~/.claude/skills/` |
| `hooks/*.sh` | symlinked as `~/.claude/hooks/`; wired via `hooks` and `statusLine` in settings |
| `settings.template.json` | copy to `~/.claude/settings.json`, fill `env` |
| `.github/workflows/doc-drift-check.yml` | weekly digest + breakage issue, assigned to the owner |

## Why this layout

Each rule file carries its own "Why" next to the rule; this section covers only the
choices between files.

- **Rules in `~/.claude/rules/`, not one big `CLAUDE.md`.** Three reasons:
  - The [memory docs](https://code.claude.com/docs/en/memory) recommend about 200 lines
    per instruction file; past that, Claude follows instructions less reliably. Topic
    files stay small by nature.
  - A rule file can do something a section of a big file cannot: with a
    [`paths:` header](https://code.claude.com/docs/en/memory#path-specific-rules) it
    loads only when Claude touches matching files. `kotlin.md` loads for Kotlin work
    and costs nothing anywhere else.
  - `@import` was considered and dropped: an imported file is pasted into the main file
    at launch and costs the same context as one big file. It only looks tidier.
- **Agent teams considered and held.** The feature is experimental, and its cost grows
  linearly by design: the orchestrator and every teammate each carry a full context.
  One orchestrator with subagents and opencode delegates covers the same need today;
  the hold ends when the docs drop the experimental label.
- **Where each kind of content goes.**
  - A rule goes into this repo only when it applies in every project. If even one
    project would not want it, it goes into that project's own `CLAUDE.md` instead;
    issue writing, notes format and TDD-by-default all moved out this way.
  - Personal preferences inside a team repo go into a gitignored `CLAUDE.local.md`.
  - Rules are never stored in auto memory: only memory's index is guaranteed to load,
    so a rule kept there may never be read. Memory holds facts, not rules.
- **Two kinds of duplication; only one is a problem.**
  - Coupled copies: the same rule kept in several places with a promise to update them
    together. The promise always breaks, the copies drift apart, and the stale copy
    still reads as authoritative. So one rule lives in exactly one file.
  - Independent copies: two projects picking the same convention separately. Changing
    one does not break the other, so nothing needs merging.
- **A rule describes the scenario it prevents, not the event that produced it.**
  Written that way, the rule stays true after everyone forgets the day it was born.
  Written as "we once did X wrong", it ages into a story.
- **Borrowed rules carry a `Source:` line.**
  - A rule taken from Kotlin conventions, Google engineering practices, or another
    source is marked; a rule distilled from the user's own incidents is not.
  - The line exists for the moment a rule reads unfamiliar: it says where the rule
    came from, so the reader reconnects and moves on instead of doubting it.
- **Personal content stays local.**
  - CLAUDE.md carries identity and preferences, so the repo tracks only
    `CLAUDE.template.md`; the filled copy lives at `~/.claude/CLAUDE.md` and nothing
    personal ships with the repo.
  - Accepted cost: that one file is not version-tracked.
- **No secrets in the repo.**
  - `settings.template.json` shows the env pattern with one self-describing example
    key; real keys and values live only in the local `settings.json`.
  - Accepted cost: a new machine fills the env block by hand. Nothing in the repo can
    do it, and repository secrets could not either; they are readable only inside a
    workflow.

## New machine

1. Clone the repo:
   ```bash
   git clone git@github.com:AhmetSIRIM/claude-setup.git ~/Projects/oss/claude-setup
   cd ~/Projects/oss/claude-setup
   ```
2. Link the always-loaded pieces:
   ```bash
   # macOS / Linux
   for i in rules skills hooks; do ln -sfn "$PWD/$i" ~/.claude/$i; done
   ```
   ```bat
   :: Windows (cmd): Git Bash's ln -s silently copies instead of linking.
   :: Junctions link for real and need neither admin rights nor Developer Mode.
   mklink /J "%USERPROFILE%\.claude\rules"  "<repo>\rules"
   mklink /J "%USERPROFILE%\.claude\skills" "<repo>\skills"
   mklink /J "%USERPROFILE%\.claude\hooks"  "<repo>\hooks"
   ```
3. Copy the personal-instructions template and fill it in; it stays local:
   ```bash
   cp CLAUDE.template.md ~/.claude/CLAUDE.md
   ```
4. Copy the settings template and fill the `env` block from your secret store:
   ```bash
   cp settings.template.json ~/.claude/settings.json
   ```
   On Windows, hook and status-line commands also need an explicit interpreter,
   because bash is not on PATH there:
   `"C:/Program Files/Git/bin/bash.exe" "C:/Users/<user>/.claude/hooks/<hook>.sh"`.
5. Tools the setup leans on (`jq` feeds the status line, `gitleaks` guards pushes):
   ```bash
   # macOS
   brew install jq gitleaks
   ```
   ```bat
   :: Windows
   winget install jqlang.jq gitleaks.gitleaks
   ```
   Then create `.git/hooks/pre-push` (chmod +x) so outgoing commits are scanned
   before they reach the remote; fall back to a full scan while `origin/main` does
   not exist yet:
   ```sh
   #!/bin/sh
   if git rev-parse --verify --quiet origin/main >/dev/null; then
     exec gitleaks git --log-opts="origin/main..HEAD" .
   fi
   exec gitleaks git .
   ```
6. Open a new session: a silent `setup-doctor` means the install is complete; anything
   missing shows up as a warning.

## Weekly digest

`doc-drift-check.yml` runs every Monday and opens at most one issue.

- It fetches the Claude Code changelog and docs, then asks one model for two things:
  a short newsletter of what changed since the last digest, and a drift check over
  the files in this repo.
- A drift finding is either a mechanism the setup uses that the docs no longer
  support, or a claim a file makes that the docs now contradict.
- The issue is assigned to the repository owner, because assignment is what makes
  GitHub send an e-mail. Nothing new and nothing broken means no issue.
- Prior issues travel back into the prompt, so closing an issue as accepted keeps the
  same finding from coming back.
- The model runs through the opencode CLI. `OPENCODE_GO_KEY` must exist as a
  repository secret: the `key` field of the `opencode-go` entry in the credential
  store `opencode auth login` writes (`~/.local/share/opencode/auth.json` on
  macOS/Linux). The default model and per-run choices live in
  the workflow's dispatch dropdown.
- The CLI version is pinned in `.github/package.json` and the workflow actions to
  commit SHAs; Dependabot updates both weekly.

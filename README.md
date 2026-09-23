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
| `hooks/*.sh` | symlinked as `~/.claude/hooks/`; wired via `hooks` in settings |
| `tools/` | `npm ci --prefix tools`; `statusLine` in settings runs ccstatusline from here with `tools/ccstatusline.json` |
| `settings.template.json` | copy to `~/.claude/settings.json`, fill `env` |
| `.github/workflows/doc-drift-check.yml` | weekly digest + breakage issue, assigned to the owner |
| `cmd/*/`, `go.mod` | Go tools; a workflow runs one with `go run ./cmd/<name>`, a hook calls the binary `go install ./cmd/<name>` puts in `~/go/bin` |

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
- **The status line is ccstatusline, pinned in `tools/`.**
  - The [status line docs](https://code.claude.com/docs/en/statusline) point to
    ccstatusline as a community project, and it is actively maintained and takes
    outside contributions. A hand-kept script would have to track the stdin schema
    alone.
  - The version sits in `tools/package.json` with a lockfile instead of a global
    install, so every machine runs the same version, the lockfile checks the package
    hash, and Dependabot turns each new release into a pull request to review.
  - Accepted cost: after a pull that changes `tools/`, each machine runs
    `npm ci --prefix tools`; the doctor warns while the installed version differs
    from the pin.
  - Accepted risk: the weekly drift check does not read `tools/`. The stdin schema is
    ccstatusline's to track, and each of its releases arrives as a Dependabot pull
    request to review.
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
   for i in rules skills hooks; do ln -sfn "$PWD/$i" ~/.claude/$i; done
   ```
3. Copy the personal-instructions template and fill it in; it stays local:
   ```bash
   cp CLAUDE.template.md ~/.claude/CLAUDE.md
   ```
4. Copy the settings template and fill the `env` block from your secret store:
   ```bash
   cp settings.template.json ~/.claude/settings.json
   ```
   The status line command names the clone at `$HOME/Projects/oss/claude-setup`; a
   clone anywhere else means editing that path in the copied settings.
5. Tools the setup leans on (`node` runs the status line, `gitleaks` guards pushes,
   `go` builds the hook that names the git mode), then the pinned status line and the
   hook binary:
   ```bash
   brew install node gitleaks go
   ```
   ```bash
   npm ci --prefix tools
   go install ./cmd/announce-git-mode
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
7. Two preferences live outside this repo and outside `settings.json`; set them by hand:
   - Terminal.app, so that Option+Backspace deletes a word in the Claude Code prompt:
     Settings > Profiles > the default profile > Keyboard > "Use Option as Meta: Left
     Option only". The right Option key keeps producing accented characters.
   - Agent view as the screen Claude Code opens on. The documented path is to start
     with `claude agents` instead of `claude` (see the
     [agent view doc](https://code.claude.com/docs/en/agent-view)). The undocumented
     alternative is `"defaultToAgentsView": true` in `~/.claude.json`, Claude Code's
     state file in the home directory, beside the `~/.claude` folder; an undocumented
     key can change or disappear in any release, because that file is Claude Code's
     own and not part of the settings contract.

## A session without permission checks

`bypassPermissions` joins the Shift+Tab mode cycle only when the session starts with
it enabled; a running session cannot add it later (see the
[permission modes doc](https://code.claude.com/docs/en/permission-modes)). Two flags:

```bash
claude --allow-dangerously-skip-permissions   # selectable in the cycle, not selected
claude --dangerously-skip-permissions         # active from the first prompt
```

The first is the one to reach for: the session starts in the usual mode, and bypass
sits one Shift+Tab past `plan` for the moment it is needed. Having bypass in the cycle
has one side effect of its own: in that session `plan` no longer blocks an edit or a
command Claude attempts while planning, so a planning session that must stay
read-only starts without either flag. Neither flag belongs in a shell alias, and
`permissions.defaultMode: "bypassPermissions"` does not belong in `settings.json`: with
bypass in every cycle an extra Shift+Tab lands on it by accident, and in that mode
writes to protected paths such as `.git` and `.claude` run without a prompt.

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
- A `budget` job runs first: `cmd/budget-check` reads the Go plan's usage windows
  (rolling, weekly, monthly). A rate-limited window marks the review job skipped,
  with the reset time in the job summary; anything else wrong with the plan or the
  key fails the run with the API's own message.
- The model runs through the opencode CLI and is set in the workflow's `env`.
  `OPENCODE_GO_KEY` must exist as a repository secret: an API key from the OpenCode
  console. One key serves both the Zen and the Go provider; the workflow registers it
  under `opencode-go` only.
- The CLI version is pinned in `.github/package.json` and the workflow actions to
  commit SHAs; Dependabot updates both weekly and watches `go.mod` and `tools/`.

# Scripting

## Shell inside the Google Shell Style Guide's boundary, Go beyond it
Shell scripts follow the Google Shell Style Guide, and shell is used only where that
guide's "When to use Shell" section allows it: a small utility or a simple wrapper that
runs a few commands in sequence. Bash is the shell. A new shell script passes ShellCheck
before it is committed. An existing script is not cleaned up wholesale: a change fixes
the findings on the lines it edits and leaves the rest. Fixing means changing the code;
a `# shellcheck disable` directive is not a fix.

Anything past the guide's boundary is written in Go. The boundary is crossed when a
script would exceed 100 lines as `wc -l` counts them, or would use non-straightforward
control flow. These cross it too, whatever the line count:
- parsing JSON or other structured data (a Claude Code hook reading its JSON stdin
  included);
- arrays or string manipulation as the main work (an array that only holds a command's
  arguments does not count);
- an error path beyond `set -euo pipefail` (retry, partial success).
The guide's own judgment on control flow still applies outside this list.

An existing script that already crosses the boundary is not rewritten for its own sake.
A bug fix stays in place; a change that adds capability moves the script to Go.

Anything else is a deviation: a shell other than Bash, or a language other than Go past
the boundary. A deviation is named with its reason before it is written, so the user
decides.

Go scripts live under `cmd/<name>/` with the module root at the repository root. The
standard library comes first; a dependency needs a stated reason, something the standard
library cannot do. A hook calls a binary installed with `go install ./cmd/<name>` from
the checkout, never `go run`; a workflow may `go run` the same path. Compiled binaries
are never committed.

Why: left to itself, a model reaches for shell first and keeps extending it past the
point where a person can still read it; the guide's boundary stops that. On the far
side, Go is the language the user has chosen as a long-term investment.

Source: Google Shell Style Guide (https://google.github.io/styleguide/shellguide.html),
"Background" (Which Shell to Use; When to use Shell) and "Features and Bugs"
(ShellCheck). The 100-line figure and the control-flow criterion are the guide's; the
trigger list, the Go choice, the layout, and the distribution path are the user's
decisions on top.

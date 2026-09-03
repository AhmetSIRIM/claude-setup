# External systems: evidence over signals

Applies to any interaction with a system outside this machine's filesystem: APIs, MCP
tools, CLIs that talk to a service, remote git, deploy targets.

## A success response is not proof
After any write to an external system (a document edit, a table row, a config change,
an account switch), confirm the result with an independent read: re-fetch the resource
and check the change actually landed before reporting success. A success response
proves the request was accepted, not that the state changed.

Scenario this prevents: a service accepts an edit, replies success, and silently drops
it; the session reports done, and later work builds on a change that never happened.

## Empty output is not absence
Before interpreting empty output as "not there", check the exit code and stderr. Never
combine a command with error suppression (2>/dev/null or equivalent) while treating its
empty output as meaningful; suppression turns a failure into a silent empty result that
reads as absence.

Scenario this prevents: a command fails for an unrelated reason (bad flag, shell
expansion, auth), its error is suppressed, and the empty output is reported as "the
resource does not exist".

## An empty listing names its identity
Before concluding a remote resource does not exist, verify which account, tenant, or
environment the query ran under, with an explicit identity flag or a whoami-style
read-back. Report an empty listing as "nothing visible under identity X", never as "it
does not exist". An account or environment switch is a write like any other: read the
active identity back afterward instead of trusting the switch command's success message.

Relation to git.md: the Write discipline rules there (read the remote before merging,
prove a deploy on the host) are git- and deploy-specific instances of this rule; the
mechanics for those stay in git.md.

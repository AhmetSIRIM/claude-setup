---
name: driving-apps-with-maestro
description: Use when driving a mobile app with Maestro on a simulator or emulator, or when you need to read what is on screen, tap an element, enter text, open a deeplink, make a flow repeatable, or see that a change works in the running app. Also triggers on "Maestro", "write a flow", "try it on the simulator", "what is on the screen", "flow yaz", "simülatörde dene", "ekranda ne var", "uygulamayı sür", and on a Maestro command failing unexpectedly.
---

# Driving Apps with Maestro

## Overview

Maestro drives an app as a black box through the accessibility layer. It needs no app code, test
target, or instrumentation; the installed `.app` or `.apk` is enough. Underneath, iOS runs
XCTest and Android runs UIAutomator.

**First principle:** **read** what is on screen; do not guess. The tree Maestro sees is not the
same as what you infer from a screenshot. Pull the tree before you write a selector.

**Second principle:** a passing step does not mean the work happened. Some Maestro commands
return success without doing anything (see "Silent failures").

The general rules hold on both platforms; the trap tables and the driver and process sections
are verified on iOS only.

## Install and version

Install from Maestro's official Homebrew tap. The plain name is a trap: `brew install maestro`
**installs the wrong tool**, a cask from runmaestro.ai that shares the name. Always give the tap:

```bash
brew install mobile-dev-inc/tap/maestro
```

The tap's formula follows each Maestro release and carries its checksum; `brew upgrade` moves to
the next one. `brew info mobile-dev-inc/tap/maestro` reports "no available formula" until the
tap is added; that is a missing tap, not a missing formula.

Maestro's own install script (`https://get.maestro.mobile.dev`) is the other path, for machines
without Homebrew. It writes a PATH line into `~/.zshrc` and `~/.bash_profile`, and it refuses to
run when a Homebrew `maestro` is already on PATH. Use one path, not both.

The CLI sends anonymous analytics by default; `MAESTRO_CLI_NO_ANALYTICS=1` turns them off, and
`MAESTRO_CLI_ANALYSIS_NOTIFICATION_DISABLED=true` hides the "Analyze with AI" banner.

Do not take a new version without checking it. On Apple Silicon with a current iOS runtime,
certain versions hang with no output (upstream issue #3137; in the report one version hangs
while an older one works). When a version works, hold it with `brew pin maestro` and upgrade on
purpose. Before an upgrade, check whether the issues described below (#2448, #3611) still apply
to the new version.

When searching GitHub issues, do not confuse `mobile-dev-inc/Maestro` with `RunMaestro/Maestro`;
the second is an unrelated project.

## Reading the screen

```bash
maestro hierarchy --no-reinstall-driver --compact
```

With `--no-reinstall-driver`, Maestro does not remove the XCUITest driver when the command ends,
and the next command reuses it if a driver already answers on the port it connects to.
`maestro hierarchy` connects to the fixed default port (22087), so repeated calls share one
driver; without the flag, every call installs the driver and removes it again, which makes each
call many times slower. `maestro test` is different: it picks a free port on every run (unless
the hidden `--driver-host-port` option fixes one), so each run starts a new driver and the flag
only leaves the old one running; see "Orphaned drivers".

`--compact` prints CSV (`element_num,depth,attributes,parent_num`) and drops empty and `false`
attributes. It is the right choice for an agent's context; the default JSON is much longer.

Fields to read: `accessibilityText` (`accessibilityLabel` on iOS), `resource-id`
(`accessibilityIdentifier` on iOS), `text`, `bounds`. On iOS, SwiftUI `Text` and button labels
can arrive only as `accessibilityText`, with `text` empty; a `text:` selector still matches
them.

In the compact output, `element_num` is not unique: the same number can appear on more than
one row, sometimes with a different depth and parent. Identify an element by its attributes
and `bounds`, not by its number. A repeated attribute (`enabled=true; enabled=true`) is an
output duplicate.

**When a step fails, do not pull a live tree.** Maestro already wrote the tree of that moment as
an artifact: `<test-output>/<flow>/screen-hierarchy/`. By the time you pull a live tree, the
screen may have changed; the artifact is the right evidence. The same output holds
`commands.json`, `logs/`, and a screenshot of the failing step. A `takeScreenshot` file is also
written inside this output (`takeScreenshot/`), not to the working directory.

## Selector semantics

**Matching is a full-string regex, not partial, and ignores case.** `text: "Step 1 of 3"` also
matches `step 1 OF 3`, while `text: "Step 1"` does not match it. For a partial match, say so:
`text: ".*Step 1.*"`. Line breaks in the text become spaces before matching.

Regex special characters (`(`, `?`, `+`, `.`) in a `text` selector: written as the full text
they still match, because the `text` match also runs a literal equality check. Inside a
partial pattern, escape them (`.*\(3\).*`).

`id` has no literal check; it is a regex only. A special character in an identifier therefore
changes what the selector means. The common trap is `$`, the end-of-string anchor:
`id: "paywall.plan.$rc_monthly"` can never match, because it asks for more characters after
the end, and the step fails with "element not found" even when an element carries exactly that
identifier. Escape such characters with a backslash (`\$`, `\(`); simpler still, keep them out
of accessibility identifiers. It is not a variable: flow
interpolation needs braces (`${NAME}`), and a bare `$name` is left as written.

`accessibilityText` is **not** a selector; it is one of the fields the `text` selector scans
(`text`, `hintText`, `accessibilityText`, and on Android only, `error`).

Order of preference, from most robust to most fragile:

| Rank | Selector | When |
|---|---|---|
| 1 | `id` | Most stable. `accessibilityIdentifier` on iOS. The right choice for icons, controls without text, and localized interfaces |
| 2 | `text` | Checks what the user sees, reads well. Depends on language and copy changes |
| 3 | Relational and state | `below`, `above`, `leftOf`, `rightOf`, `childOf`, `containsChild`, `containsDescendants`, `traits`; for state, `enabled`, `checked`, `focused`, `selected` |
| 4 | `index` | Visual order (top to bottom, then left to right). Negative index works: `index: -1` is the last |
| 5 | `point` | Last resort. A screen coordinate depends on resolution and layout |

The legitimate use of `point` is an offset **inside an element**, together with a selector
(`tapOn: {id: x, point: "50%,85%"}`); that does not carry the fragility of a screen coordinate.

Where text selectors break in a localized interface: short words match more than once (`OK`,
`All`), text gets truncated, the locale changes. If you keep tapping an element that has no
label, the fix is not a coordinate; add an `accessibilityIdentifier` to the app.

When several nodes match, the **deepest** one wins (on a tie, the tappable one comes first). If
a text matches both a cell and the label inside it, `tapOn` presses the inner label; that
explains most "it matched but tapped the wrong place" cases.

## Waiting

- `assertVisible` and the other assertions poll, pulling the tree again, until their default
  timeout. That timeout is clearly longer than the short timeout of `optional: true` steps and
  of `when` / `while` conditions. The values live in `Orchestra.kt` as `lookupTimeoutMs` and
  `optionalLookupTimeoutMs`. Most places need no extra wait around them.
- As a result, `when: visible` waits only briefly. Guarding a late element with `when` can skip
  the block because the element is not drawn yet.
- For waits longer than the default, use `extendedWaitUntil: {visible: <element>, timeout: <ms>}`.
  A rough starting rule (from experience, not from a source): app launch 15-30 s, an API call
  10-15 s.
- **Wait for the content to appear, not for the spinner to go away.** `notVisible` is a race:
  it is also true before the spinner is drawn.
- **No fixed `sleep`.** Maestro leaves this command out on purpose; every command retries
  against the latest tree until its own timeout.
- The first command after `launchApp` hits the splash screen. Wait for the first real screen
  element with `extendedWaitUntil`.

## Silent failures

These return success but did not do the work. Check the outcome of each one separately:

| Command | Silent behavior | Countermeasure |
|---|---|---|
| `back` (iOS) | A no-op in the iOS driver; the step reports success and nothing happens | `tapOn` the back button in the navigation bar, then `assertVisible` an element of the previous screen |
| `tapOn` (tab bar and similar) | On iOS the tap reports success but no navigation happens. One suspected cause is a press on a child view that does not take taps; upstream issue #2448 has not settled the cause. The tree can also keep the views of tabs that are not on screen, so an `assertVisible` of a common element can pass on the wrong tab | After every navigation tap, `assertVisible` an element only the target screen draws (a title or a row, never the tab's own label, which is on screen from every tab), or check a screenshot |
| `hideKeyboard` (iOS) | There is no native API; it tries two small swipes near the center of the screen. It does not fail if the keyboard stays open, and the swipes can scroll content | `tapOn` an empty area, then check that the keyboard is gone |
| `waitForAnimationToEnd` | Does **not fail** when the timeout runs out; it counts as success and moves on | Wait for a concrete element instead. With an endless animation (shimmer, an auto-rotating banner) this command always runs into its timeout |
| `takeScreenshot` | Called with an undefined variable, it raises no error and writes `undefined.png` | Give a fixed file name |
| `scrollUntilVisible` | Visibility is the element's area clipped to the screen bounds; a tab bar or keyboard covering it is ignored. A covered element counts as "fully visible" and scrolling stops. Also, any `visibilityPercentage` below 100 acts as 0 in practice (integer division) | Leave `visibilityPercentage` at its default (100); add `centerElement: true` when needed; then `tapOn` the element and check the result. Under a floating tab bar the `tapOn` lands on the tab instead of the row: `tapOn` always taps the center of the element it selects, and a relational selector (`above`, `below`) changes only which element is selected, not where the tap lands. Swipe the row clear of the bar first, or tap a point inside the row that the bar does not cover (`tapOn: {id: <row id>, point: "50%,20%"}`) |

## iOS-specific behavior

- **Use a simulator.** The CLI has a path for real devices, but basic operations such as
  `clearState`, `clearKeychain`, `launchApp`, `openLink`, and screenshots are not supported
  there. MCP accepts simulators only.
- When `openLink` opens a deeplink, iOS can show a confirmation dialog ("Open in ...?"). The
  choice **persists** in the simulator and **`clearState` does not reset it**. Pass the dialog
  with a conditional block.
- On iOS, `clearState` **reinstalls** the app, which is slow. It does not clear the keychain; use
  `clearKeychain` for that, but it wipes the keychain of the **whole simulator**, not only the
  app's.
- `launchApp` without parameters restarts the app (`stopApp` defaults to `true`). To bring a
  backgrounded app to the front, use `stopApp: false`.
- `maestro hierarchy` takes an XCTest snapshot, which can stall while the main thread is busy
  (upstream issue #1967: the hierarchy request fails with "main thread busy for 30.0s"), for
  example on a screen with a live map. Run it under a time limit
  (`perl -e 'alarm 120; exec @ARGV' maestro hierarchy`) and use screenshots on such screens.
- On a map, a `swipe` with a long `duration` can act as a press and not pan; a short one (about
  300 ms) pans. The threshold is observed, not documented.
- To test a path with no location fix, reboot the simulator after `xcrun simctl location <udid>
  clear`; until then the location service can still hand out the last cached fix.
- In mixed UIKit and SwiftUI hosting, if an element never shows up in the tree, try
  `platform.ios.snapshotKeyHonorModalViews: false` in `config.yaml`. This setting works **only in
  the CLI**; the MCP session does not read it.

## Conditional flow

`when` works on `runFlow` and `runScript`; `repeat` takes `while`. Single commands such as
`tapOn` take no condition.

```yaml
- runFlow:
    when:
      visible: "Close"
    commands:
      - tapOn: "Close"
```

Conditions: `visible`, `notVisible`, `platform`, `true` (JS). Several in one block are ANDed.

`optional: true` has different semantics: `when` checks the condition and **skips**, while
`optional` tries and **swallows the failure**, marking the step as a warning. Use it only where
skipping the step is truly acceptable, and check the expected state in the next step.

## Running flows

A flow file is a config header, a `---` line, then a list of commands. `commands:` is not a
top-level key; it only nests inside commands such as `runFlow`, `repeat`, and `retry`.

```yaml
appId: com.example.app
---
- launchApp
- assertVisible:
    id: home.title
```

**The default is to run the directory with one command** (`maestro test flows/`): the driver is
installed once and the flows run in the same session. In a directory run, flow order is **not
alphabetical and not defined**; when order matters, write `executionOrder.flowsOrder` in
`config.yaml`.

Keep subflows shared through `runFlow` (such as login) outside the directory you run (for
example `shared/` next to `flows/`); otherwise the directory run also runs them as standalone tests.

**If the driver dies in batch mode.** There are reports of the driver going unresponsive after a
few flows on the iOS 26 runtime, with `Failed to connect to 127.0.0.1:<port>` (upstream issues
#3254, #3318). Both are closed without a confirmed cause, so treat this as a possible failure,
not a known one. If you hit this error, run each flow in its own process, so each one starts a fresh driver.
Leave `--no-reinstall-driver` out here: each `maestro test` picks a new port, so the flag
reuses nothing and only leaves one driver behind per flow. Save the snippet as a script file;
pasted into an interactive shell, `exit` closes the shell:

```bash
#!/bin/bash
rc=0
for f in flows/*.yaml flows/*.yml; do
  [[ -e "$f" ]] || continue
  maestro test "$f" || { echo "FAILED: $f"; rc=1; }
done
exit "$rc"
```

`rc` is required: `|| echo` alone swallows every failure, the loop exits 0, and CI shows green.

**The simulator runs where the user can see it, by default.** The user may want to watch the
run even without saying so. `xcrun simctl boot` starts the simulator headless: Maestro drives
it, but no window opens. So after booting, open the window. From Xcode 27 on there is no
Simulator.app; Device Hub takes its place (`open -b com.apple.dt.Devices`). On older Xcode,
`open -a Simulator`. Maestro Viewer over MCP also works. Stay headless only when the user asks
for it.

Do not quit Device Hub while a run is in progress. Quitting it shuts the simulator down and
kills the Maestro driver; the next command fails with a connection error, and the user loses
the view they were watching.

**Orphaned drivers.** On Xcode 26 and later, a driver session that fails or is killed can leave
its `xcodebuild test-without-building` process and a `simctl diagnose` behind, and the diagnose
can run until its own `--timeout` (upstream issue #3633). Never end a driver while a Maestro command
is running; that command fails with `DeviceUnreachableException`. In a wait loop, use `pgrep -x <name>` or wait on a PID;
`pgrep -f <pattern>` matches the loop's own shell, and the loop never ends.

Leftovers also pile up after runs that pass. With `--no-reinstall-driver`, the CLI does not stop
the driver's `xcodebuild` process when it exits (`LocalXCTestInstaller`: `uninstall()` returns
early when the driver is not reinstalled), and a `simctl diagnose` was seen next to most of them.
An agent that calls `maestro test` with the flag many times collects one pair per call, all on
the same simulator. Count them after every few calls. While no Maestro command is running, end
all of them for your simulator's UDID; the next command starts a driver of its own. The
`simctl diagnose` command line carries no UDID, so find it through its parent `xcodebuild`.
The bracket in each pattern keeps it from matching the shell that runs it:

```bash
for p in $(pgrep -f "[t]est-without-building.*-destination id=<udid>"); do
  pkill -P "$p" -f "[s]imctl diagnose"
  kill "$p"
done
```

At the end of a run, end whatever is left the same way and shut the simulator down.

**Name the simulator for its job.** `xcodebuild -destination 'name=<device name>'` picks the first
simulator with that name, so a test lane can land on the device a manual run is using. Give each
job its own named simulator, or address it by UDID.

When the work is done, shut the simulator down (`xcrun simctl shutdown <udid>`); a booted
simulator holds several GB of memory.

## MCP or CLI

Use both, in different phases. Add the MCP server with `claude mcp add maestro -- maestro mcp`;
its tools arrive as `mcp__maestro__<tool>`.

- **Exploring and writing flows: MCP.** `mcp__maestro__inspect_screen` returns the tree as
  compact JSON, filters the noise, and carries the selector rules in its tool description.
  `mcp__maestro__open_maestro_viewer` opens a web view of the live device and the running
  commands. Inside the CLI, the Viewer ships only with the MCP server; `maestro studio` points to
  a separate desktop app.
- **Verification and repeatability: CLI.** The MCP `run` tool carries none of the output and
  report flags, reinstalls the driver on every session, and an MCP call cannot be committed; a
  flow file can.

**Single simulator rule:** in MCP, the iOS driver port is fixed at 22087 and every Simulator
session shares it. If a second simulator is booted, `inspect_screen` returns the wrong device's
tree and taps go to the wrong device, **and the result still looks valid** (upstream issue
#3611). Keep a single booted simulator while an agent works. In the CLI, `maestro hierarchy` uses the
same fixed port, while `maestro test` picks a free port on every run.

## Commands that send data out

When the screen shows personal or confidential data, do not use these without a deliberate
decision:

- `maestro record` renders the recording on Maestro's servers **by default**, which means it
  uploads it. To keep it local, use `--local`, or `startRecording` / `stopRecording` inside a
  flow.
- `assertWithAI`, `assertNoDefectsWithAI`, `extractTextWithAI`, and `maestro test --analyze`
  send the screen to Maestro Cloud and need a Maestro account (login, `--api-key`, or
  `MAESTRO_CLOUD_API_KEY`). There is no option to supply your own OpenAI or Anthropic key.

## Common mistakes

The costliest breaks of the rules above, for a quick check:

| Mistake | Instead |
|---|---|
| Assuming what is on screen | Pull the tree. The knowledge that an element "should be there" may come from another context |
| Treating a passing step as work done | See "Silent failures"; after navigation, `assertVisible` the target screen |
| Jumping to `point` for an unlabeled element | First `id`, then a relational selector, then `index`. `point` is the last resort |
| Guarding a late element with `when` | `when` waits briefly; use `extendedWaitUntil` first |
| Booting headless and never opening the window | Open Device Hub (or Simulator.app) after the boot; the user may be watching |
| Using MCP with two simulators booted | Keep one booted simulator; otherwise the wrong device is driven and nobody notices |
| `brew install maestro` | Installs a different product. Use `brew install mobile-dev-inc/tap/maestro` |

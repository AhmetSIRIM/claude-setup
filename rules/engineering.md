# Engineering defaults

## Simplest complete solution
Write the simplest implementation that fully meets the current requirement. Before
adding a library or hand-writing a capability, check what the project's existing
dependencies already provide; read their docs and types rather than assuming.

## Check for the mainstream path first
Before hand-writing what a tool may already provide, look for its preset, documented
default, or sanctioned path, and say what you found. A hand-written equivalent has to be
maintained here and drifts from what it copied, so it should be a deliberate choice, not
the first reflex. Deviating is sometimes right; when you do, name the deviation and its
reason before writing it, so the user decides with the trade-off in view.

Why: after an implementation with AI it is sometimes unclear whether a solution is a
workaround or a real fix, project-specific or mainstream. That ambiguity is unhealthy.
A model will happily hand-write what a preset already provides; this rule makes the
choice visible instead of silent.

## Names carry the role
Name an abstraction by its role, not its shape. `Clock`, `Helper`, `Util`, and a bare
`Manager` force the reader to open the file; a role-shaped name carries half the
contract on the file name.

Gateway family (one seam to one external concern, no orchestration):

| Name     | Connotation                                   | Typical shape              |
|----------|-----------------------------------------------|----------------------------|
| Gateway  | neutral, Fowler                               | anything                   |
| Port     | Hexagonal seam where an adapter plugs in      | anything                   |
| Provider | single-source, on-demand read                 | `getX()`, `findX()`        |
| Observer | push-style stream of changes                  | `observeX(): Flow<X>`      |
| Listener | push-style single-event subscription          | `onEvent(...)` registration|
| Notifier | domain pushes outward to an external system   | `notify(event)`            |

Repository vs Gateway: a Repository orchestrates two or more sources for one domain
concept (cache plus remote); a Gateway talks to exactly one. Downgrading a Repository
to "Provider" because it sounds lighter hides the orchestration.

Orchestrator family:

| The class...                                              | Name                |
|-----------------------------------------------------------|---------------------|
| owns state across calls (queue, mutex, cache, lifecycle)  | `<Concern>Manager`  |
| runs one stateless command and returns a result           | Use Case            |
| coordinates several repositories under one boundary       | Application Service |

`Manager` is reserved for a stateful orchestrator that owns a long-lived concern, and the
prefix names that concern (`SessionManager`, `BatchManager`). Suffix-only `Manager`,
`DataManager`, `AppManager` are not allowed; the name must constrain what belongs inside.

Applies to new abstractions; existing classes are not retro-renamed.

Source: Gateway and Repository follow Fowler, Patterns of Enterprise Application
Architecture; Port comes from Cockburn's hexagonal architecture; Observer and Listener
are the classic GoF/platform meanings. The tables and the Manager constraint are
decisions made on top of those sources, not quotes.

## One concept, one word
Pick one term per concept and keep it: `startIndex` everywhere, not `offset` in one place
and `beginIndex` in another; one of `element`, `item`, `entry` for collection values, not
a mix. Related operations come in regular pairs (`first` / `firstOrNull`,
`single` / `singleOrNull`). Overloads behave identically regardless of the argument type.
Parameters run from essential to optional, in the same order across related functions.

Source: Kotlin API Guidelines, "Consistency". Borrowed convention.

## Retire, do not remove
Nothing public is deleted outright. It goes through a deprecation cycle: `@Deprecated`
with a `message` that says what changed and why and a `replaceWith` that migrates the
call, at level `WARNING` first, then `ERROR`, then `HIDDEN`, across releases. The same
applies inside an app: an old path is marked and migrated over time, not cut in one
change. The exception is a project that has declared itself pre-1.0 in its own CLAUDE.md.

Source: Kotlin API Guidelines, "Backward compatibility", section Evolve APIs
pragmatically. Borrowed convention, extended by the user to application code.

## Descriptive names, no noise
A name says what the thing is; length grows with scope. A longer name that says more
beats a shorter one that says less, but noise words (data, info, helper) and type
echoes say nothing.

Source: Kotlin Coding Conventions, "Choose good names"; Robert C. Martin, Clean Code
ch. 2. Borrowed convention.

## Declarative over imperative
Prefer expression bodies, collection operations over manual loops, immutability, and
sealed hierarchies with exhaustive `when`. When editing existing code, follow the
file's established style and raise a refactor as a suggestion rather than applying it.
Deviate for performance or clarity, and say why. Language idiom wins: in a culture that
is deliberately imperative (Go), the idiom is the established style.

## Test doubles: the species is in the name
Fake, Stub, Saboteur, or Mock appears in the class name, so the double's contract reads
from the name alone:

- **Fake**: a lightweight working implementation (an in-memory repository).
- **Stub**: canned answers, no logic.
- **Saboteur**: injected failures, for error-path tests.
- **Mock**: canned returns plus interaction assertions.

A Fake's doc comment lists Matches and Divergences against the real implementation;
a fake's classic trap is silent drift, and that list keeps the divergence auditable.

Choose fake vs mock by behavior depth: the test needs the dependency to behave
(state, sequencing) means fake; canned returns or call-verification only means mock
or stub. When in doubt, prefer the fake.

Source: vocabulary from Meszaros, xUnit Test Patterns (Fake Object; Saboteur as a Test
Stub variation); the fake-first preference matches Software Engineering at Google
ch. 13, "Prefer Realism Over Isolation", and the classical school in Fowler's "Mocks
Aren't Stubs".

## Rename and delete need a multi-front search
Before renaming or deleting a symbol, run a separate search for each way it can be
referenced: direct code references, string literals, reflection and DI wiring,
configuration files, build scripts, resource files, and tests. One grep for the
identifier is never accepted as proof that nothing else points at it; each front is
checked and the result of each check is reported. On Android this binds without
exception: manifests, XML resources, ProGuard/R8 keep rules, and serialization reach
classes by string name, so the compiler stays green after the break. The same holds
wherever a name lives in a string: Spring property keys and bean names, JSON field
names, Go struct tags, storyboard and plist entries.

Scenario this prevents: a rename compiles cleanly and ships, but a manifest entry, a
keep rule, or a serialized field name still carries the old string; nothing fails until
the release build or the first deserialization in production.

Source: harperreed/dotfiles, .claude/CLAUDE.md, "Rename safety" section ("One grep
always misses something"): https://github.com/harperreed/dotfiles

## No silent fallbacks
An error is handled meaningfully or propagated; it never disappears. Concretely banned:
- a catch block that is empty or log-only and lets the flow continue as if the operation succeeded;
- a `runCatching` or try/catch that returns a default value masking the failure (`getOrDefault` and friends), or an ignored error return in languages with explicit errors;
- a fallback path added without the user deciding that a fallback is the desired behavior.

When a fallback is genuinely wanted, the code states what failure it absorbs, and the degraded behavior is observable: logged at warning level or above, or surfaced to the caller.

Scenario this prevents: a feature works in the happy path and fails silently in production; the swallowed exception leaves no signal, so the defect is found weeks later from user reports instead of minutes later from a stack trace.

Source: Kirill Markin, "Claude Code Rules: CLAUDE.md Global Instructions"; echoed in harperreed/dotfiles. Borrowed convention, extended by the user with the observability requirement.

## CI logic lives in script files
A workflow step's `run` block stays a one-liner that calls a script in the
repository's scripts directory; logic beyond a few lines moves into that script,
and the workflow file reads as orchestration only.

Why: an inline block can be neither syntax-checked nor run locally; a script file
gets both, plus room for its own comments.

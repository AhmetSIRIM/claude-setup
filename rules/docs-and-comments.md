# Docs and comments

Applies to code comments, KDoc and equivalents, comments in build and CI files, and
tracked prose (README, docs/). Test files welcome generous KDoc; the durability rules
below still apply to them.

## Timeless
No volatile facts: no counts ("13 fields", "67 tests"), no versions, dates, or default
values, no claim tied to the current shape of neighboring code ("the only caller"). No
references to plans, specs, tasks, sprints, drafts, or sessions. No debate-responsive
phrasing ("we decided to change this because..."). Docs read as present-tense snapshots:
no past-relative words (dropped, previously, originally); future-relative is fine
(deferred, planned).

Self-check before finalizing a line: "if someone changes the neighboring code tomorrow,
does this sentence become false?" If yes, rewrite toward role and rationale.

Why: a stale but authoritative sentence is worse than no sentence; readers trust it.

## Standalone
A block is understandable without opening other files. State the contract in place. When
the contract is genuinely owned elsewhere, give one pointer, never a copy; copies drift.

A committed artifact never points at a file that is not itself committed (git-ignored
notes, `.superpowers/`, `.claude/`, private files). Such a file exists only on the disk
where it was written; anyone who clones the repository finds the reference pointing at
nothing.

## What earns a comment
Only a non-obvious invariant earns a comment. Three things that do not:
- general language or framework knowledge (what `copy()` does on a data class);
- behavior the runtime already signals (a warning the code logs on misuse);
- narration of an obvious mechanism, in tests as well as in production code.
Prefer removal over enrichment. When unsure whether a comment is non-obvious, cut it.

Source for the why-over-what framing: Google Engineering Practices, "What to look for in
a code review", section Comments. The three exclusions above are the user's own.

## KDoc shape
- Open with a summary fragment: a noun or verb phrase ("Returns the foo"), not a full
  sentence ("This method returns..."). Indexes show only that fragment.
- Avoid `@param` and `@return`. Embed the parameter in the prose as `[name]`; it links.
  Use the tags only when an explanation is too long for the flow of the text.
- A self-explanatory member (`getFoo`, `foo`) needs no KDoc when there is truly nothing
  to say beyond "returns the foo". The exception does not cover a term the reader may
  not know: `canonicalName` still gets a line explaining what canonical means here.

Source: Kotlin Coding Conventions, "Documentation comments" (param tags); Google Kotlin
Style Guide, "Documentation" (summary fragment, self-explanatory exception). These three
are borrowed conventions, not lessons from this user's own projects; check the source
before removing one.

## TODO format
Continuation lines sit one space deeper than the opener; no blank comment lines inside
the block. Why: IntelliJ extends the TODO highlight onto continuation lines only when
they are indented one deeper; a same-level line or a blank `//` drops the highlight for
the rest of the block.

## Plain English
Written artifacts in English stay plain and short: common vocabulary, no stacked
subordinate constructions. The reader may be a non-native speaker; reader cost matters
more than register.

## Docs move with the code
The commit that changes the code changes its docs. Doc updates are not deferred to a
later cleanup; that later never arrives and the doc is stale in the meantime.

Source: Google Documentation Best Practices, "Update docs with code" and "Delete dead
documentation". Borrowed convention.

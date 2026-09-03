# Session hygiene

## Session files have a home and a lifetime
A file that is meaningless without its session (scratch script, draft, probe output,
intermediate dump, backup copy) is session-scoped. Session-scoped files are written to
the session's scratch directory (the scratchpad or job tmp dir the harness provides;
`/tmp` as fallback), never into the repository or the home directory.

When a session-scoped file must live inside the repository (a tool resolves it by
relative path, the build must see it), it is an exception: name it and its removal
point out loud when creating it, keep it on a running list, and delete everything on
that list at session end or when the user asks for a sweep. A session never ends with
an undeclared session-scoped file in the repository or in HOME.

Promotion is explicit: a scratch file that turns out to be worth keeping is moved to a
permanent location and named for its role, in the open, not left where it was born.

Why: an orphaned scratch file outlives the only context that explained it. Weeks later
it reads as possibly load-bearing, so nobody deletes it, and the repository and HOME
silt up with files no one can account for.

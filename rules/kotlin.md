---
paths:
  - "**/*.kt"
  - "**/*.kts"
---

# Kotlin

## One top-level declaration per file
Each file has a single focus. Cohesion exceptions: a sealed hierarchy with its
variants, extension functions grouped per receiver, and private helpers used only by
that file's main declaration.

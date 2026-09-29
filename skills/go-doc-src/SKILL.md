---
name: go-doc-src
description: Retrieve the source code of a Go symbol (function, type, method, const, var) that isn't in the current context — e.g. in another package, or in a dependency outside the working tree. Use during Go code review, debugging, or when using an unfamiliar API and you need to see the actual implementation rather than guess from its name or signature. Triggers on requests like "how does X work internally", "show me the source of Y", "what does this stdlib/vendored function do".
license: Apache-2.0
compatibility: Requires the Go toolchain (the `go` command) on PATH.
metadata:
  author: Olivier Mengué
  version: "0.2.0"
---

# Go doc source lookup

When you need the source of a Go symbol that is not already visible in context (a
different package in the same module, a dependency, or the standard library), use
`go doc -src` instead of trying to locate and open the file manually. It resolves
build tags, GOPATH/module cache locations, and package aliases correctly, and prints
just the focused declaration (plus its doc comment) rather than a whole file.

## Usage

```
go doc -src <pkg>.<Symbol>          # exported symbol
go doc -src -u <pkg>.<symbol>       # unexported symbol (requires -u)
go doc -src <pkg>.<Type>.<Method>   # method on a type
```

`<pkg>` can be an import path (`net/http`), a relative import path from the current
module, or just the package name if it's unambiguous in context. Run the command from
within the relevant module (or a directory where the module is resolvable) so `go doc`
can find it — for stdlib and already-downloaded module-cache packages this works from
anywhere.

Without `-u`, `go doc` only resolves exported identifiers — looking up an unexported
one fails with `doc: no symbol <name> in package <pkg>`. Add `-u` to also match
unexported identifiers (it still matches exported ones too).

## Case matching — use `-c` for an exact name

By default, lower-case letters in the argument match either case, while upper-case
letters match exactly. So `net/http.client` (even with `-u`) silently returns the
exported `http.Client`, and a lower-case query matching several symbols prints all
of them.

When you know the exact identifier, add `-c` to respect case:

```
go doc -src -u -c <pkg>.<symbol>
```

With `-c`, a symbol that doesn't exist with that exact case fails with
`doc: no symbol <name> in package <pkg>` instead of returning a different one.

## When to reach for this vs. reading the file directly

- Prefer `go doc -src` when you only need one symbol's implementation — it's faster
  and avoids pulling an entire (possibly large or vendored) file into context.
- Fall back to reading the file directly when you need surrounding context (adjacent
  helpers, package-level vars/consts it relies on, file structure) that `-src` won't
  include.

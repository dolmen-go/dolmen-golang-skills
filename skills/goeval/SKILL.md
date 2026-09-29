---
name: goeval
description: How and when to run Go as a script with goeval (github.com/dolmen-go/goeval). Use this skill whenever a task calls for a throwaway program — a bulk or structured file edit, parsing or validating YAML/JSON/Markdown, probing a library or an API, computing something a shell pipeline would do badly — and before writing any scratch Go program. Scripting here is done in Go, never in Python, so reach for this rather than for a python3 heredoc.
---

# Running Go as a script with goeval

`goeval` compiles and runs a fragment of Go — the body of `main` — without a
module, a file or any boilerplate. It makes Go a practical scripting language,
which is what this user wants instead of Python.

## Invocation

The program can be an argument, but prefer standard input with a *quoted*
here-document: the shell then leaves the Go source completely alone, so `$`,
backticks and quotes need no escaping.

```sh
goeval -i os,log,strings - <<'EOF'
data, err := os.ReadFile("CHANGELOG.md")
if err != nil {
	log.Fatal(err)
}
os.Stdout.WriteString(strings.ToUpper(string(data)))
EOF
```

Arguments after `-` reach the program as `os.Args[1:]`, and `os.Exit` status
propagates to the shell, so a goeval program works in a pipeline or as a test:

```sh
goeval -i fmt,os - a b <<'EOF'
fmt.Println("args:", os.Args[1:])
os.Exit(3)
EOF
```

## Imports

Pack the imports into a single `-i`, comma-separated, rather than repeating the
flag: `-i os,log,strings`.

Standard library imports are resolved automatically, so `-i` is optional for
them; listing them anyway documents the program and avoids a surprise when a
name is ambiguous. A third-party package needs `path@version`, which switches
goeval to module mode and downloads it (so it needs the network the first time),
and it packs with the others:

```sh
goeval -i os,log,gopkg.in/yaml.v3@latest - <<'EOF'
…
EOF
```

## When to use it

Reach for goeval when the work needs real data structures, a parser or a
library, but not a repository:

- **Structured edits to a file** — inserting a line at the end of every bullet,
  converting inline links to reference-style, redistributing definitions per
  section. Anything where `sed` would need lookahead or state.
- **Parsing and validating** — check a GitHub workflow with `gopkg.in/yaml.v3`
  and `KnownFields(true)`, walk JSON, inspect Go source with `go/ast`.
- **Probing an API** — what does this stdlib function actually return for that
  input, what does a dependency do with a malformed value. Faster and more
  honest than reasoning about it, and it uses the toolchain that is installed.
- **Computing a number** for a report, from files or command output.

Do not use it when:

- **The program will be kept or rerun by others.** Write a real package and run
  it with `go run ./path/to/tool` (a directory under `.claude/skills/*/scripts/`
  or a `cmd/` is fine — the go tool ignores directories starting with `.`).
  A committed tool deserves a doc comment, flags and a name.
- **A single substitution would do.** `sed -i` or the Edit tool is clearer than
  ten lines of Go.

## Writing the program

**Fail loudly.** A bulk edit that half-applies is worse than one which refuses:
check every assumption and `log.Fatal` on anything unexpected — an anchor that
does not match, a label already bound to another URL, a count that differs from
what was expected. Then the program is safe to re-run.

```sh
if !strings.Contains(s, anchor) {
	log.Fatalf("anchor not found: %q", anchor)
}
```

**Backticks.** Go raw strings are delimited by backticks, so a literal cannot
hold the backticks of a Markdown code span. Use interpreted strings, or key the
replacements on backtick-free substrings.

**Verify the result, don't trust the edit.** After rewriting a file, prove that
only the intended dimension changed — for a rewrap, that the sequence of words
is identical:

```sh
diff <(tr -s ' \n' '\n\n' <before) <(tr -s ' \n' '\n\n' <after)
```

## Inspecting what runs

- `-E` prints the assembled source instead of running it, which is the quickest
  way to see how the fragment was wrapped and which imports were added.
- `-x` prints the commands executed, `-o <file>` builds a binary instead of
  running it.
- `-play`, `-share` run or publish the code on the Go playground; do not send a
  program which carries the user's data there.

## Scratch modules

When a fragment grows into a small module in a temporary directory, do not name
the module after a standard library package (`cmp`, `slices`, `os`…): the import
becomes ambiguous and the build fails with a confusing message. Name it after
the task.

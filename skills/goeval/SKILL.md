---
name: goeval
description: How and when to run Go as a script with goeval (github.com/dolmen-go/goeval). Use this skill whenever a task calls for a throwaway program — a bulk or structured file edit, parsing or validating YAML/JSON/Markdown, probing a library or an API, computing something a shell pipeline would do badly — and before writing any scratch Go program. Prefer it over a python3 heredoc for throwaway scripts.
---

# Go as a script with goeval

`goeval` compiles and runs the body of `main`: no module, file or boilerplate.
Pass it on stdin via a *quoted* here-document, so `$`, backticks and quotes need
no escaping. Arguments after `-` are `os.Args[1:]`; the exit status propagates.

```sh
goeval -i os,log,strings,gopkg.in/yaml.v3@latest - arg1 arg2 <<'EOF'
data, err := os.ReadFile(os.Args[1])
if err != nil {
	log.Fatal(err)
}
…
EOF
```

- One comma-separated `-i`, not repeated flags. Stdlib imports are
  auto-resolved but listing them avoids ambiguity. Third-party packages need
  `path@version` (module mode, network on first use).
- `-E` prints the assembled source instead of running it, `-x` prints the
  commands, `-o <file>` builds a binary. `-play`/`-share` send code to the Go
  playground: never with the user's data.

## When

Use it for work that needs data structures, a parser or a library, but not a
repository: stateful file edits where `sed` falls short, parsing/validating
(`yaml.v3` with `KnownFields(true)`, JSON, `go/ast`), checking what an API
really returns for an input, computing a number for a report.

Not for: a single substitution (use `sed` or Edit), or a program that will be
kept or rerun by others — write a real package run with `go run ./path`
(`cmd/` or `.claude/skills/*/scripts/`), with a doc comment and flags.

## How

- **Fail loudly**: `log.Fatal` on every unmet assumption (anchor not found,
  unexpected count, conflicting value), so a bulk edit never half-applies and
  is safe to re-run.
- Raw strings cannot contain backticks: use interpreted strings for Markdown
  code spans.
- **Verify the result**: prove only the intended dimension changed, e.g. for a
  rewrap `diff <(tr -s ' \n' '\n\n' <before) <(tr -s ' \n' '\n\n' <after)`.

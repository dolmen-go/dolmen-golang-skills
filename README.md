# dolmen-golang-skills

[Olivier Mengué](https://github.com/dolmen)'s collection of [Agent
Skills](https://agentskills.io) for Go development — reusable instructions
that teach Claude (or any compatible agent) Go-specific techniques it
wouldn't otherwise reach for.

## Skills

### [go-doc-src](./skills/go-doc-src/SKILL.md)

Fetch the source of a Go symbol that isn't already in context — a function,
type, method, const, or var in another package, a dependency, or the standard
library — using `go doc -src` instead of guessing from the name or hunting
down the file by hand.

It also documents a caveat verified against real `go doc` behavior: symbol
lookups are case-insensitive, so a lookup can silently return a different
symbol than the one you asked for (e.g. `net/http.client` resolves to the
exported `http.Client`). The skill tells the agent to check the returned
declaration's exact name before trusting it.

Neither the official Claude Code marketplace nor the community Go-skill
collections found at the time of writing (e.g. `samber/cc-skills-golang`)
document this technique. The closest existing tool, the `gopls-lsp` plugin
(Go language-server integration), solves a different problem: interactive
navigation/refactoring in code you already have open, not on-demand source
retrieval for a symbol you don't have in context yet. Benchmarked against
`strings.ToUpper`, `go doc -src` was ~35-45x faster and needed one exact tool
call versus two-plus imprecise ones via `gopls`.

### [goeval](./skills/goeval/SKILL.md)

Use Go, not Python, as the scripting language for throwaway programs — bulk or
structured file edits, parsing and validating YAML/JSON/Markdown, probing a
library or an API — by running the body of `main` directly with
[goeval](https://github.com/dolmen-go/goeval), without a module, a file or any
boilerplate.

## Install (Claude Code)

```
/plugin marketplace add dolmen-go/dolmen-golang-skills
/plugin install <skill-name>@dolmen-golang-skills
```

For example:

```
/plugin install go-doc-src@dolmen-golang-skills
```

Or, for local development/testing:

```
claude --plugin-dir /path/to/dolmen-golang-skills
```

## Install (other agents)

Each skill lives at `skills/<name>/SKILL.md` and follows the plain [Agent
Skills open standard](https://agentskills.io/specification) — no Claude
Code-specific frontmatter or syntax — so any harness that supports the
standard can load one directly by pointing at its directory, or by copying
it into its own skills directory.

## License

Apache-2.0, see [LICENSE](./LICENSE).

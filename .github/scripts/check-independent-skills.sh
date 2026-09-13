#!/usr/bin/env bash
# Enforces this repo's marketplace convention: every skill under skills/ is
# declared as its own plugin entry in .claude-plugin/marketplace.json, scoped
# to exactly that skill via "skills", with "strict": false (no plugin.json,
# marketplace entry is the sole authority — see README/architecture notes).
set -euo pipefail

MARKETPLACE=".claude-plugin/marketplace.json"
fail=0

# 1. Every plugin entry must have "strict": false.
non_strict=$(jq -r '.plugins[] | select(.strict != false) | .name' "$MARKETPLACE")
if [[ -n "$non_strict" ]]; then
  echo "::error::Plugin entries missing \"strict\": false: $non_strict"
  fail=1
fi

# 2. Every plugin entry must scope to exactly one skill (independent plugin).
not_single=$(jq -r '.plugins[] | select((.skills // []) | length != 1) | .name' "$MARKETPLACE")
if [[ -n "$not_single" ]]; then
  echo "::error::Plugin entries not scoped to exactly one skill via \"skills\": $not_single"
  fail=1
fi

# 3. The skills declared in marketplace.json must exactly match skills/*/SKILL.md.
declared=$(jq -r '.plugins[].skills[]?' "$MARKETPLACE" | sed -E 's#^\./skills/##' | sort -u)
actual=$(for d in skills/*/; do
  name="$(basename "$d")"
  [[ -f "${d}SKILL.md" ]] && echo "$name"
done | sort -u)

missing=$(comm -23 <(echo "$actual") <(echo "$declared"))
extra=$(comm -13 <(echo "$actual") <(echo "$declared"))

if [[ -n "$missing" ]]; then
  echo "::error::Skill directories with no plugin entry: $missing"
  fail=1
fi
if [[ -n "$extra" ]]; then
  echo "::error::Plugin entries reference a nonexistent skill directory: $extra"
  fail=1
fi

if [[ "$fail" -eq 0 ]]; then
  echo "All skills are declared as independent, strict:false plugin entries."
fi

exit "$fail"

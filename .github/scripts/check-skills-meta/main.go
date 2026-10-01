// Command check-skills-meta enforces this repo's marketplace convention: every
// skill under skills/ is declared as its own plugin entry in
// .claude-plugin/marketplace.json, scoped to exactly that skill via "skills",
// with "strict": false (no plugin.json, marketplace entry is the sole
// authority — see README/architecture notes).
//
// Usage, from the repository root:
//
//	go -C .github/scripts/check-skills-meta run . "$PWD"
//
// The optional argument is the repository root (default: current directory).
//
// Violations are reported as GitHub Actions error annotations, and the exit
// status is 1 if any is found.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const marketplacePath = ".claude-plugin/marketplace.json"

type marketplace struct {
	Plugins []struct {
		Name   string   `json:"name"`
		Strict *bool    `json:"strict"`
		Skills []string `json:"skills"`
	} `json:"plugins"`
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("check-skills-meta: ")

	switch len(os.Args) {
	case 1:
	case 2:
		if err := os.Chdir(os.Args[1]); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("usage: check-skills-meta [<repo-root>]")
	}

	data, err := os.ReadFile(marketplacePath)
	if err != nil {
		log.Fatal(err)
	}
	var m marketplace
	if err := json.Unmarshal(data, &m); err != nil {
		log.Fatalf("%s: %v", marketplacePath, err)
	}

	fail := false
	report := func(msg string, names []string) {
		if len(names) > 0 {
			fmt.Printf("::error::%s: %s\n", msg, strings.Join(names, ", "))
			fail = true
		}
	}

	var nonStrict, notSingle, declared []string
	for _, p := range m.Plugins {
		// 1. Every plugin entry must have "strict": false.
		if p.Strict == nil || *p.Strict {
			nonStrict = append(nonStrict, p.Name)
		}
		// 2. Every plugin entry must scope to exactly one skill (independent plugin).
		if len(p.Skills) != 1 {
			notSingle = append(notSingle, p.Name)
		}
		for _, s := range p.Skills {
			declared = append(declared, strings.TrimPrefix(s, "./skills/"))
		}
	}
	report(`Plugin entries missing "strict": false`, nonStrict)
	report(`Plugin entries not scoped to exactly one skill via "skills"`, notSingle)

	// 3. The skills declared in marketplace.json must exactly match skills/*/SKILL.md.
	skillFiles, err := filepath.Glob("skills/*/SKILL.md")
	if err != nil {
		log.Fatal(err)
	}
	var actual []string
	for _, f := range skillFiles {
		actual = append(actual, filepath.Base(filepath.Dir(f)))
	}
	slices.Sort(declared)
	declared = slices.Compact(declared)

	var missing, extra []string
	for _, name := range actual {
		if _, found := slices.BinarySearch(declared, name); !found {
			missing = append(missing, name)
		}
	}
	for _, name := range declared {
		if !slices.Contains(actual, name) {
			extra = append(extra, name)
		}
	}
	report("Skill directories with no plugin entry", missing)
	report("Plugin entries reference a nonexistent skill directory", extra)

	if fail {
		os.Exit(1)
	}
	fmt.Println("All skills are declared as independent, strict:false plugin entries.")
}

// Command lang guards the zones Cyrillic has already left.
//
// The move to English goes in steps, and every step has to have a checkable
// mark of completion. Without one, "the zone is translated" is the author's
// memory: it holds until the first edit made out of habit, and it is
// discovered on the consumer's side.
//
// Zones are switched on ONE AT A TIME and at the END of a step, not at the
// start: a gate raised red lives with its check disabled and guards nothing.
// The list of zones lies here in full, the switched-off ones included —
// otherwise "how much is done" and "how much is left" would have to be counted
// from somebody's memory.
//
// The command stands on its own rather than inside the site build: the build
// sees src/ and docs/, while the zones of the move also cover dist/, tools/,
// assets/ and the files at the root.
//
//	go run ./cmd/lang        check the zones that are switched on
//	go run ./cmd/lang -v     show the switched-off ones as well
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
)

type zone struct {
	name  string
	paths []string // files and directories, relative to the repository root
	on    bool
	step  string // the step of the move on which the zone is switched on
	why   string
	// except lists files inside the zone where Cyrillic is DATA rather than
	// text, together with the reason. A reason is mandatory: an exception
	// without one is indistinguishable from a forgotten file, and a list of
	// those turns the gate back into a wish.
	except []exception
	// skipSuffix takes a whole CLASS of files out of the zone, and it exists
	// for exactly one thing: the Russian half of a bilingual corpus. Ninety-four
	// files cannot be listed one by one as exceptions, and a list that long
	// would say nothing anyway — what matters is the rule, and the rule is that
	// a page named `*.ru.md` is Russian BY ITS NAME. Everything else under the
	// zone is English and is checked.
	skipSuffix string
}

type exception struct {
	path string
	why  string
}

var zones = []zone{
	{
		name:  "package manifest",
		paths: []string{"package.json"},
		on:    true,
		step:  "1",
		why:   "description is visible in the npm registry",
	},
	{
		name:  "marks",
		paths: []string{"assets"},
		on:    true,
		step:  "1",
		why:   "the files ship inside the package together with the comments in them",
	},
	// Step 2 went file by file rather than in one piece: 3 392 Cyrillic lines
	// is 15–22 sessions, and all that time the zone "src as a whole" would have
	// stood switched off. The work of twenty sessions comes back out of habit
	// exactly as the work of one does, so the zone grew: a translated file
	// entered the list in the same commit that translated it.
	//
	// The twentieth file — tokens.css — closed the list, and it is collapsed
	// into "src". The collapse is the mark of completion for the step: while
	// the zone held names, a missing name was visible; now nothing is missing.
	{
		name:  "kit sources",
		paths: []string{"src", "dist/instrument.js", "dist/instrument.css"},
		on:    true,
		step:  "2",
		why:   "src/ ships inside the package through the files field",
	},
	{
		name:  "tools",
		paths: []string{"tools"},
		on:    true,
		step:  "3",
		why:   "the gates' output lives in the same lines as their code",
		except: []exception{
			{
				path: "tools/cmd/mutate/main.go",
				why: "the mutation table holds two kinds of Russian, and neither is text. " +
					"An anchor points into a `*.ru.md` page, and a payload is the very " +
					"string a gate has to reject: feeding cmd/lang or the comment rule " +
					"English would test nothing.",
			},
			{
				path: "tools/cmd/registry/main.go",
				why: "the frontmatter says the kind of an API entry in Russian, and the " +
					"ownership map is built by matching it. Both spellings are listed in " +
					"the pattern because both kinds of page are permanent.",
			},
			{
				path: "tools/cmd/docscheck/contract.go",
				why: "the Russian words heading the contract table on a page, and the " +
					"Russian for yes and for if, are the PATTERN a row is parsed by rather " +
					"than a phrase. Both spellings stand here because the corpus is " +
					"bilingual by design rather than in transit: a `*.ru.md` page is read " +
					"by the same parser as its English base.",
			},
			{
				path: "tools/cmd/docscheck/tokens.go",
				why: "the Russian words for black and white are how an alpha cell is " +
					"written in docs/foundations/tokens.ru.md. They are the pattern a cell " +
					"is matched against rather than a phrase.",
			},
		},
	},
	{
		name: "site",
		// The build output is not on the list, and not because it is
		// gitignored: `site/dist` holds the RENDERED pages of both languages,
		// so a zone that watched it would be red on every Russian page for as
		// long as the site is bilingual — which is to say for good — while
		// saying nothing about the site's own sources. What the built page IS
		// checked for is the other direction: `site/internal/check` forbids
		// Cyrillic on a page that is not Russian. `public` holds CNAME and
		// robots.txt.
		paths: []string{"site/cmd", "site/internal", "site/public", "site/go.mod"},
		on:    true,
		step:  "3",
		why:   "the output of the build lives in the same lines as the code",
		// Six files, and in every one of them the Cyrillic is a KEY or a
		// PATTERN rather than a phrase: it is matched against a Russian page,
		// or it is one half of a translation that already exists. None of it
		// is translated — all of it leaves with the last Russian page.
		except: []exception{
			{
				path: "site/internal/i18n/i18n.go",
				why: "the Russian half of every entry IS the translation rather than text " +
					"awaiting one. The base language flipped on step 5 and the map was not " +
					"touched by it: only Base and the order of All changed.",
			},
			{
				path: "site/internal/content/sections.go",
				why: "the aliases are the headings printed on a page, and a page is " +
					"matched to a section by them. Both spellings stand there because both " +
					"kinds of page are permanent.",
			},
			{
				path: "site/internal/content/markdown.go",
				why: "the transliteration table turns the heading of a Russian page into an " +
					"anchor, and the kind vocabulary beside it lists both spellings. Both " +
					"are read against a page rather than shown to anybody.",
			},
			{
				path: "site/internal/content/content.go",
				why: "apiKinds is the vocabulary a page writes in its frontmatter, and the " +
					"kind generated for a token is a key of that same dictionary. What is " +
					"printed comes from i18n, not from here.",
			},
			{
				path: "site/internal/render/assets/docs.js",
				why: "the two Russian letters folded onto one another are what makes a " +
					"heading findable when the reader types the other of the pair, and the " +
					"one visible phrase beside them is the translation itself.",
			},
			{
				path: "site/internal/check/assets.go",
				why: "the chronicle phrases are the PATTERN the comment gate matches, and " +
					"feeding it English would test nothing while Russian comments can still " +
					"appear. The English half stands beside them.",
			},
		},
	},
	{
		name:       "documentation",
		paths:      []string{"docs/start", "docs/foundations", "docs/components", "docs/agent", "docs/layout", "docs/blocks", "docs/about"},
		on:         true,
		step:       "5",
		skipSuffix: ".ru.md",
		why:        "the base language is flipped: `page.md` is English, and Russian lives in `page.ru.md`",
	},
	{
		name:  "root",
		paths: []string{"README.md", "CONTRIBUTING.md", "ROADMAP.md", "docs/README.md"},
		step:  "6",
		why:   "CHANGELOG.md is not translated: a translated chronicle is a rewritten one",
	},
}

func main() {
	root := flag.String("root", "..", "repository root")
	verbose := flag.Bool("v", false, "show the switched-off zones")
	flag.Parse()

	var problems []string
	onCount, exCount := 0, 0
	for _, z := range zones {
		if !z.on {
			continue
		}
		onCount++
		exCount += len(z.except)
		for _, p := range z.paths {
			problems = append(problems, z.scan(filepath.Join(*root, filepath.FromSlash(p)), *root)...)
		}
	}
	sort.Strings(problems)

	if len(problems) > 0 {
		fmt.Printf("── Cyrillic in a zone it has left (%d) ──\n", len(problems))
		for _, p := range problems {
			fmt.Println("  ·", p)
		}
		fmt.Println()
		fmt.Println("The zone is declared translated, and the mark of completion for the step")
		fmt.Println("rests on this check. Put the line back into English, or drop the zone")
		fmt.Println("deliberately.")
		os.Exit(1)
	}

	fmt.Printf("· zones under guard: %d of %d, exceptions: %d — no Cyrillic in them\n",
		onCount, len(zones), exCount)
	if *verbose {
		for _, z := range zones {
			if z.on {
				for _, e := range z.except {
					fmt.Printf("    exception in «%s»: %s\n      %s\n", z.name, e.path, e.why)
				}
				continue
			}
			fmt.Printf("    waiting for step %s: %s (%s)\n", z.step, z.name, strings.Join(z.paths, " · "))
		}
	}
}

// scan walks a file or a directory of a zone.
func (z zone) scan(path, root string) []string {
	var out []string
	filepath.WalkDir(path, func(p string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return nil
		}
		rel := p
		if r, err := filepath.Rel(root, p); err == nil {
			rel = filepath.ToSlash(r)
		}
		if z.excepted(rel) {
			return nil
		}
		if z.skipSuffix != "" && strings.HasSuffix(rel, z.skipSuffix) {
			return nil
		}
		text := strings.ReplaceAll(string(b), "\r\n", "\n")
		for i, line := range strings.Split(text, "\n") {
			if w := cyrillicIn(line); w != "" {
				out = append(out, fmt.Sprintf("%s:%d: %s — zone «%s»", rel, i+1, w, z.name))
			}
		}
		return nil
	})
	return out
}

func (z zone) excepted(rel string) bool {
	for _, e := range z.except {
		if e.path == rel {
			return true
		}
	}
	return false
}

// cyrillicIn returns the first Cyrillic word of a line — the message has to
// show WHAT WAS FOUND rather than only where: from the word alone it is
// immediately clear whether the line was forgotten or meant.
func cyrillicIn(line string) string {
	start := -1
	for i, r := range line {
		if unicode.Is(unicode.Cyrillic, r) {
			if start < 0 {
				start = i
			}
			continue
		}
		if start >= 0 {
			return "«" + line[start:i] + "»"
		}
	}
	if start >= 0 {
		return "«" + line[start:] + "»"
	}
	return ""
}

// Command mutate checks the CHECKS: it breaks one invariant at a time and
// demands that the matching gate go red.
//
// WHY IT EXISTS. A gate can be green and empty, and reading it will not reveal
// that. The "button weight ladder" passed that way forever: it compared the
// lightness of two colours as if they were opaque, while the kit's recessed
// surfaces are films, and a film contributed its own lightness, that is zero.
// The step against any backdrop came out enormous. The check existed, reported
// green and verified nothing — right up until the default button in the dark
// themes stopped differing from the background.
//
// Hence the rule: a check nobody has ever seen red is not a check but
// decoration. Every mutation below is the question "and what if I break it?"
// asked of a machine rather than of memory.
//
// HOW IT WORKS. The tree is copied into a temporary directory, exactly one
// thing is spoilt in the copy, the gate is run against the copy and MUST
// return a non-zero code. A caught mutation is a "caught" line; a missed one
// is a hole in the gate, and the command fails.
//
// The control pass comes first: on an untouched copy every gate has to be
// green. Without it the harness would measure its own noise rather than the
// mutations.
//
//	go run ./cmd/mutate        the whole list
//	go run ./cmd/mutate -v     with the output of the failing gate
//	go run ./cmd/mutate -only contrast   only the mutations of one gate
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// A gate: its name, its package, and how to tell it where the tree is.
//
// mod is the module the gate is built from; empty means tools. perRun means
// rebuild for EVERY mutation rather than once at the start. Exactly one gate
// needs that: the site build carries its styles and templates inside itself
// through go:embed, and a mutation in them is visible only to a rebuilt
// binary. Building them all that way would mean paying a build for each of
// forty-odd mutations for the sake of one.
type gate struct {
	pkg    string
	mod    string
	perRun bool
	args   func(tree string) []string
}

var gates = map[string]gate{
	"contrast":   {"./cmd/contrast", "", false, func(t string) []string { return []string{"-tokens", t + "/src/tokens.css"} }},
	"targets":    {"./cmd/targets", "", false, func(t string) []string { return []string{"-tokens", t + "/src/tokens.css"} }},
	"proportion": {"./cmd/proportion", "", false, func(t string) []string { return []string{"-tokens", t + "/src/tokens.css"} }},
	"dist":       {"./cmd/dist", "", false, func(t string) []string { return []string{"-src", t + "/src", "-out", t + "/dist"} }},
	"docscheck": {"./cmd/docscheck", "", false, func(t string) []string {
		return []string{"-kit", t + "/src", "-docs", t + "/docs", "-stage", t + "/stage.css"}
	}},
	"registry": {"./cmd/registry", "", false, func(t string) []string {
		return []string{"-src", t + "/src", "-docs", t + "/docs", "-registry", t + "/components.json"}
	}},
	"lang": {"./cmd/lang", "", false, func(t string) []string { return []string{"-root", t} }},
	// The site build is the largest gate by volume of measurement and until
	// recently the only one with no mutation standing against it at all. It is
	// built from the COPIED tree and run inside it: it looks for internal, cmd
	// and tools relative to its own directory, and it keeps its styles inside
	// the binary — there is nothing to mutate in it from outside.
	"site": {"./cmd/site", "site", true, func(t string) []string {
		return []string{"-docs", "../docs", "-kit", "../src", "-assets", "../assets", "-out", t + "/site-out"}
	}},
}

// A mutation is one breakage. `from` has to exist in the file: a mutation that
// replaced nothing would produce a false "caught" on an untouched tree.
type mutation struct {
	name string
	gate string
	file string
	from string
	to   string
	why  string
}

// RUSSIAN INSIDE THE TABLE IS DATA, NOT TEXT — and it is here on purpose.
//
// Two kinds, and neither translates with the rest of the file.
//
// An ANCHOR (`from`) points at documentation that is still Russian: step 4 of
// the migration has not happened, and an anchor translated ahead of the page
// it points at would stop matching. When those pages turn English the anchors
// follow them, in the same commit.
//
// A PAYLOAD (`to`) is the very string a gate has to reject. cmd/lang exists to
// find Cyrillic in a translated zone, and check.Comments looks for Russian
// chronicle phrases — feeding them English would test nothing.
//
// This is why tools/cmd/mutate/main.go is named in the exception list of
// cmd/lang: the file cannot be free of Cyrillic while it is doing its job.
var mutations = []mutation{
	// ── contrast ───────────────────────────────────────────────────────────
	{"text: muted raised until unreadable", "contrast", "src/tokens.css",
		"--text-muted:     light-dark(var(--n-8),  var(--n-6));",
		"--text-muted:     light-dark(var(--n-5),  var(--n-6));",
		"a label stops taking 4.5:1 on a panel"},
	{"stack: the panel levelled with the page", "contrast", "src/tokens.css",
		"--surface-page:     light-dark(var(--n-1), var(--n-13));",
		"--surface-page:     light-dark(var(--n-0), var(--n-13));",
		"depth is carried by the order of lightness, and the step is gone"},
	{"button levelled with its backdrop", "contrast", "src/tokens.css",
		"--surface-recessed:        light-dark(oklch(0 0 0 / 0.060), oklch(1 0 0 / 0.050));",
		"--surface-recessed:        light-dark(oklch(0 0 0 / 0.002), oklch(1 0 0 / 0.002));",
		"precisely the defect this harness was written for"},
	{"ladder inverted: soft louder than default", "contrast", "src/tokens.css",
		"--surface-recessed-hover:  light-dark(oklch(0 0 0 / 0.025), oklch(1 0 0 / 0.015));",
		"--surface-recessed-hover:  light-dark(oklch(0 0 0 / 0.025), oklch(1 0 0 / 0.090));",
		"the third step sounds louder than the second, and an absolute difference cannot see it"},
	{"ladder: hover levelled with default", "contrast", "src/tokens.css",
		"--surface-recessed-hover:  light-dark(oklch(0 0 0 / 0.025), oklch(1 0 0 / 0.015));",
		"--surface-recessed-hover:  light-dark(oklch(0 0 0 / 0.058), oklch(1 0 0 / 0.048));",
		"two weights of a button produce one look"},
	{"text painted with a token that has no pair", "contrast", "src/base.css",
		"body {",
		".mut-probe { color: var(--n-7); }\nbody {",
		"a text colour absent from every pair is a threshold nobody measured"},

	// ── targets ────────────────────────────────────────────────────────────
	{"control height pushed below the minimum", "targets", "src/tokens.css",
		"--control-h-sm: 26px;",
		"--control-h-sm: 14px;",
		"a tap target under 24 and without the gap that compensates for the size"},

	// ── proportion ─────────────────────────────────────────────────────────
	{"type steps merged", "proportion", "src/tokens.css",
		"--text-xs:  0.78125rem;",
		"--text-xs:  0.86rem;",
		"the scale declares a size that cannot be seen"},
	{"radius turned odd", "proportion", "src/tokens.css",
		"--radius-md: 8px;",
		"--radius-md: 7px;",
		"the arc lands off the device grid at density 1.5"},
	{"geometry shrinks as the scale grows", "proportion", "src/tokens.css",
		"  --label-col:     120px;",
		"  --label-col:     8px;",
		"the scale only goes up"},
	{"row vertical wider than horizontal", "proportion", "src/tokens.css",
		"--row-pad-y:   var(--space-3);",
		"--row-pad-y:   var(--space-6);",
		"part of the vertical air was already delivered by the leading"},

	// ── dist: the "Forbidden" section ──────────────────────────────────────
	{"!important in a component", "dist", "src/actions.css",
		".inst-btn {", ".mut-probe { color: red !important; }\n.inst-btn {",
		"the only one allowed is [hidden] in base.css"},
	{"weight 700", "dist", "src/actions.css",
		".inst-btn {", ".mut-probe { font-weight: 700; }\n.inst-btn {",
		"the kit has two weights, and 700 shouts louder than the data"},
	{"raw colour in a component", "dist", "src/actions.css",
		".inst-btn {", ".mut-probe { color: #333; }\n.inst-btn {",
		"a number instead of semantics is a hardcoded light theme"},
	{"spacing utility", "dist", "src/layout.css",
		".inst-stack {", ".mt-3 { margin-top: 12px; }\n.inst-stack {",
		"rhythm is set by the flow primitives, not by the markup of every screen"},
	{"colour inside a data URI", "dist", "src/agent.css",
		"stroke='%23000'", "stroke='%23ff0000'",
		"the shape is drawn by a mask and coloured by a token"},
	{"image as a background rather than a mask", "dist", "src/data.css",
		".inst-tag {", ".mut-probe { background-image: url(\"x.svg\"); }\n.inst-tag {",
		"a background image colours itself and does not follow the theme"},
	{"gap bypassing the role tier", "dist", "src/layout.css",
		".inst-cluster {", ".mut-probe { gap: var(--space-8); }\n.inst-cluster {",
		"a large step of the scale taken directly"},
	{"stray closing brace", "dist", "src/motion.css",
		".inst-caret { animation: none; opacity: 1; }",
		".inst-caret { animation: none; opacity: 1; } }",
		"the rest of the file leaves the layer and starts winning against the application"},

	// ── docscheck ──────────────────────────────────────────────────────────
	{"a class in an example that does not exist", "docscheck", "docs/components/actions/button.md",
		"<button class=\"inst-btn\"", "<button class=\"inst-btn inst-btn--mut\"",
		"a reader copies the markup and it silently does not work"},
	{"a number in the token table lied", "docscheck", "docs/foundations/tokens.md",
		"| `--surface-recessed` | чёрный 6% | белый 5% |",
		"| `--surface-recessed` | чёрный 9% | белый 5% |",
		"the class is there, the token is there, only the number lies"},
	{"a number in the density table lied", "docscheck", "docs/foundations/density.md",
		"| `--row-pad-y` | `--space-2` | `--space-3` | `--space-4` |",
		"| `--row-pad-y` | `--space-2` | `--space-3` | `--space-6` |",
		"the mode table drifted away from the density block"},
	{"a number in the scale table lied", "docscheck", "docs/foundations/scale.md",
		"| `--text-sm` | 14 | 15 | 16 | 17 | 18 |",
		"| `--text-sm` | 14 | 15 | 19 | 17 | 18 |",
		"the type size in the table is not the one in the code"},
	{"source leads nowhere", "docscheck", "docs/components/actions/button.md",
		"source: src/actions.css", "source: src/nowhere.css",
		"a person goes to look at how it is built and lands in nothing"},

	// ── docscheck: the language of the tables ──────────────────────────────
	//
	// The five mutations below check one thing: that the gate matches tables by
	// STRUCTURE rather than by Russian words in the header and in the cells.
	// Each translates one word and lies with a number at the same time. If the
	// gate went blind on the translation the lie passes, and "documentation and
	// kit agree completely" becomes a report about nothing.
	{"mode table header translated", "docscheck", "docs/foundations/density.md",
		"| Токен | `compact` | по умолчанию | `comfortable` |\n|---|---|---|---|\n| `--control-h-xs` | 18px | 20px | 22px |",
		"| Token | `compact` | по умолчанию | `comfortable` |\n|---|---|---|---|\n| `--control-h-xs` | 18px | 20px | 24px |",
		"one word in the header used to silence the check over the whole table"},
	{"default column translated", "docscheck", "docs/foundations/density.md",
		"| Токен | `compact` | по умолчанию | `comfortable` |\n|---|---|---|---|\n| `--size-check` | 13px | 15px | 17px |",
		"| Токен | `compact` | default | `comfortable` |\n|---|---|---|---|\n| `--size-check` | 13px | 19px | 17px |",
		"the column was recognised by its caption, and the caption was the only Russian one of three"},
	{"mode code in the header translated", "docscheck", "docs/foundations/density.md",
		"| Токен | `compact` | по умолчанию | `comfortable` |\n|---|---|---|---|\n| `--pad-panel` | `--space-4` | `--space-5` | `--space-6` |",
		"| Токен | `плотная` | по умолчанию | `comfortable` |\n|---|---|---|---|\n| `--pad-panel` | `--space-2` | `--space-5` | `--space-6` |",
		"a mode code is part of the kit, and a translated code has to be heard"},
	{"alpha cell translated, value lied", "docscheck", "docs/foundations/tokens.md",
		"| `--surface-recessed` | чёрный 6% | белый 5% |",
		"| `--surface-recessed` | black 9% | белый 5% |",
		"translating a cell used to carry away the check of its number"},
	{"alpha cell unreadable", "docscheck", "docs/foundations/tokens.md",
		"| `--surface-recessed-hover` | чёрный 2.5% | белый 1.5% |",
		"| `--surface-recessed-hover` | тёмный 2.5% | белый 1.5% |",
		"an unparsed cell in a machine column is silence, not permission"},
	{"token table lied in translation", "docscheck", "docs/foundations/tokens.en.md",
		"",
		"---\ntitle: Tokens\n---\n\n| Token | Light | Dark |\n|---|---|---|\n| `--surface-recessed` | black 9% | white 5% |\n",
		"translating a page used to carry its tables out from under the check entirely"},
	{"mode table lied in translation", "docscheck", "docs/foundations/density.en.md",
		"",
		"---\ntitle: Density\n---\n\n| Token | `compact` | default | `comfortable` |\n|---|---|---|---|\n| `--row-pad-y` | `--space-2` | `--space-3` | `--space-6` |\n",
		"the check was nailed to a Russian file name rather than to a page"},

	// ── registry ───────────────────────────────────────────────────────────
	{"component outside the registry", "registry", "components.json",
		"\"toolbar\": {", "\"toolbar-renamed\": {",
		"no cross-cutting axis watches that component any more"},
	{"an axis lost its value", "registry", "components.json",
		"\"flow\": \"inline\",\n      \"print\": \"hide\",\n      \"page\": \"flow\",\n      \"forced\": \"none\"",
		"\"flow\": \"inline\",\n      \"page\": \"flow\",\n      \"forced\": \"none\"",
		"an omission and a deliberate \"nothing\" have to be distinguishable"},
	{"CSS assigns a class to print, the registry is silent", "registry", "src/print.css",
		"  .inst-segmented { display: none; }",
		"  .inst-segmented, .inst-timeline { display: none; }",
		"the list quietly grows with things nobody decided on"},
	{"infinite animation without a slowdown", "registry", "src/status.css",
		".inst-dot {", ".inst-mut-spin { animation: inst-pulse 1s linear infinite; }\n.inst-dot {",
		"an infinite animation squeezed to 0.01ms is a stop, not a slowdown"},
	{"a group item lost its state", "registry", "docs/components/actions/segmented.md",
		"<button type=\"button\" role=\"radio\" aria-checked=\"false\" tabindex=\"-1\">Сетка</button>",
		"<button type=\"button\" role=\"radio\" tabindex=\"-1\">Сетка</button>",
		"state lives in an attribute rather than in a class"},
	{"two roving tabindexes in one group", "registry", "docs/components/actions/segmented.md",
		"<button type=\"button\" role=\"radio\" aria-checked=\"false\" tabindex=\"-1\">Сетка</button>",
		"<button type=\"button\" role=\"radio\" aria-checked=\"false\" tabindex=\"0\">Сетка</button>",
		"Tab would walk every item and the group would stop being one control"},
	{"a meter without aria-valuenow", "registry", "docs/components/charts/meter.md",
		"aria-valuenow=\"43\"", "data-mut=\"43\"",
		"the role promises a number to assistive technology"},
	{"a scripted role not declared in the registry", "registry", "components.json",
		"\"role\": \"tablist\",", "\"role\": \"tablist-renamed\",",
		"the script promises a keyboard on the kit's behalf and the registry knows nothing about it"},
	// The "Related" section of a shape 2 page is rendered from the registry and
	// from nothing else. A name pointing at nothing drops a link from the page
	// silently; a one-way link is a statement about two components that only one
	// of them was told about — and that is exactly how the four hand-written
	// lists of "Actions" had already drifted.
	{"a neighbour that is not in the registry", "registry", "components.json",
		"\"related\": [\n      \"segmented\",\n      \"button\",",
		"\"related\": [\n      \"segmented-renamed\",\n      \"button\",",
		"the link would vanish from the page without a word"},
	{"a link that goes only one way", "registry", "components.json",
		"\"related\": [\n      \"segmented\",\n      \"button\",\n      \"tag\",",
		"\"related\": [\n      \"segmented\",\n      \"button\",",
		"a neighbour that does not know it is one is a list nobody keeps in step"},
	// ── lang ───────────────────────────────────────────────────────────────
	//
	// A zone is declared translated — so a Russian string that comes back into
	// it out of habit has to fail here rather than in a consumer's aria-label.
	{"a Russian string came back into the public defaults", "lang", "dist/instrument.js",
		"copied: 'Copied',", "copied: 'Скопировано',",
		"a screen reader speaks this phrase, and a consumer should replace it by choice, not out of necessity"},
	{"a Russian string in the tag template", "lang", "dist/instrument.js",
		"tagRemoved: (label) => `Tag ${label} removed`,",
		"tagRemoved: (label) => `Метка ${label} снята`,",
		"a template literal is just as spoken a phrase as a plain one"},
	{"a Russian description in the package manifest", "lang", "package.json",
		"\"description\": \"CSS kit for interfaces",
		"\"description\": \"CSS-кит для интерфейсов",
		"description is visible in the npm registry to any outsider"},
	{"a Russian comment rode along into a mark", "lang", "assets/logo.svg",
		"<!-- The mark: a vertical run",
		"<!-- Знак: вертикальный ряд",
		"the file ships inside the package together with the comments in it"},
	{"Russian came back into a translated file", "lang", "src/motion.css",
		"  /* Transitions collapse rather than switch off",
		"  /* Переходы схлопываются, а не выключаются",
		"step 2 runs for twenty sessions, and all that time a file rests on memory alone"},
	// The tools zone carries named exceptions, and an exception list is the one
	// thing that can quietly grow until the zone guards nothing. The mutation
	// lands in a file that is NOT excepted: it proves the zone is alive, and it
	// fails the moment somebody widens an exception to cover the whole
	// directory.
	{"Russian came back into the tools", "lang", "tools/cmd/targets/main.go",
		"// min is the WCAG 2.2 AA 2.5.8 minimum (Target Size, Minimum).",
		"// Норма WCAG 2.2 AA 2.5.8 (Target Size, Minimum).",
		"a zone with exceptions has to be checked outside them, or the exceptions become the zone"},

	{"two roving tabindexes on a translated page", "registry", "docs/components/actions/segmented.en.md",
		"",
		"---\ntitle: Segmented\nsource: src/actions.css\n---\n\n```html preview\n<div class=\"inst-segmented\" role=\"radiogroup\" aria-label=\"View\">\n  <button type=\"button\" role=\"radio\" aria-checked=\"true\" tabindex=\"0\">List</button>\n  <button type=\"button\" role=\"radio\" aria-checked=\"false\" tabindex=\"0\">Grid</button>\n</div>\n```\n",
		"a translation is written from scratch and errs more often than the original, and only the original was checked"},

	// ── the site build ─────────────────────────────────────────────────────
	//
	// The largest gate by volume of measurement: the page contract, links, the
	// sprite, the site's tokens against the kit, and the comment rule across
	// the whole repository. Not a single mutation stood against it, which means
	// exactly as much was known about its green colour as about any unchecked
	// one.
	{"a site token the kit does not have", "site", "site/internal/render/assets/docs.css",
		".site-logo:hover { background: var(--surface-hover); }",
		".site-logo:hover { background: var(--surface-nonesuch); }",
		"the site scaffolding takes from the kit something the kit never promised"},
	{"raw colour in the site scaffolding", "site", "site/internal/render/assets/docs.css",
		".site-logo:hover { background: var(--surface-hover); }",
		".site-logo:hover { background: #f3f3f3; }",
		"a number instead of semantics is a hardcoded light theme on the kit's own site"},
	{"a chronicle comment in the kit", "site", "src/actions.css",
		".inst-btn {",
		"/* Раньше здесь стоял другой радиус. */\n.inst-btn {",
		"the history of edits instead of why it is the way it is now"},
	{"a chronicle comment in the script", "site", "src/kit.js",
		"function rovingOwned(group) {",
		"// Раньше здесь стоял флаг на самом пункте.\nfunction rovingOwned(group) {",
		"the gate's zone is the whole repository, not just CSS"},
	{"an early */ inside an explanation", "site", "src/status.css",
		".inst-dot {",
		"/* Полоса 1.20*/1.34 задаёт вид. */\n.inst-dot {",
		"the comment closes early and eats the rule that follows"},
	{"a chronicle comment in the tools", "site", "tools/cmd/contrast/main.go",
		"func main() {",
		"// Раньше здесь стоял отдельный порог для крупного кегля.\nfunc main() {",
		"tools/** entered the gate's zone; without a mutation that rested on nothing"},
	// Two page-contract dictionaries live side by side while the pages migrate
	// one at a time, and the risk of that arrangement is specific: a page that
	// declares the new shape while keeping the old sections satisfies neither
	// dictionary and nobody notices, because the gate that used to judge it now
	// judges something else.
	//
	// The mutation declares shape 2 on a page written under shape 1 and demands
	// that the gate say so. It fires on both halves of the new rule at once —
	// the missing markup contract and the sections shape 2 no longer keeps —
	// and that is honest while no page has migrated yet. The half that
	// separates them lands together with the first migrated page.
	// The page has to be one still in shape 1, and it is picked from a section
	// that has no index yet: such a section cannot migrate at all until the
	// index exists, so the anchor will not move out from under the harness on
	// the next migrated page.
	//
	// When the last page moves to shape 2 this mutation goes with it: it tests
	// the transition, and the transition will be over.
	{"a page declares shape 2 without moving to it", "site", "docs/components/display/tag.md",
		"layout: component\nsource: src/data.css",
		"layout: component\nshape: 2\nsource: src/data.css",
		"a half-migrated page passes neither dictionary, and the gate has to name which one it fails"},
	// The two halves of the new rule, each on its own, now that a migrated page
	// exists to break. The mutation above fires on both at once and cannot tell
	// them apart; these two can, and a rule whose halves are never seen
	// separately is a rule nobody has checked.
	{"a shape 2 page lost its markup contract", "site", "docs/components/actions/chip.md",
		"## Контракт", "### Контракт",
		"the contract is the section shape 2 exists for; without it the page promises a role and explains nothing"},
	{"a shape 2 page kept a section of the old shape", "site", "docs/components/actions/chip.md",
		"## Состояния", "## Использование",
		"while both shapes are allowed on one page, nobody can tell which one it follows"},
}

func main() {
	verbose := flag.Bool("v", false, "show the output of the failing gate")
	only := flag.String("only", "", "only the mutations of one gate")
	keep := flag.Bool("keep", false, "do not delete the tree of a missed mutation")
	root := flag.String("root", "..", "repository root")
	flag.Parse()

	repo, err := filepath.Abs(*root)
	must(err)

	// The gates are built ONCE: thirty `go run` calls in a row spend minutes
	// recompiling the same thing.
	binDir, err := os.MkdirTemp("", "instrument-gates-")
	must(err)
	defer os.RemoveAll(binDir)

	fmt.Println("building the gates…")
	bin := map[string]string{}
	build := func(name, root string) string {
		g := gates[name]
		mod := g.mod
		if mod == "" {
			mod = "tools"
		}
		out := filepath.Join(binDir, name+exeSuffix())
		cmd := exec.Command("go", "build", "-o", out, g.pkg)
		cmd.Dir = filepath.Join(root, mod)
		if b, err := cmd.CombinedOutput(); err != nil {
			fmt.Fprintf(os.Stderr, "cannot build %s: %v\n%s", name, err, b)
			os.Exit(1)
		}
		return out
	}
	for name, g := range gates {
		if g.perRun {
			continue
		}
		bin[name] = build(name, repo)
	}

	run := func(name, tree string) (bool, string) {
		g := gates[name]
		exe := bin[name]
		if g.perRun {
			exe = build(name, tree)
		}
		cmd := exec.Command(exe, g.args(tree)...)
		// The working directory is the gate's module IN THE COPY, if a module
		// is named. The site gate looks for src, docs and tools relative to its
		// own directory: run from the repository it would honestly check the
		// repository and never notice the mutation, and the harness would score
		// that as "the invariant is not guarded".
		if g.mod != "" {
			cmd.Dir = filepath.Join(tree, g.mod)
		} else {
			cmd.Dir = filepath.Join(repo, "tools")
		}
		b, err := cmd.CombinedOutput()
		return err == nil, string(b)
	}

	// ── control: an untouched tree has to be green ─────────────────────────
	base, err := os.MkdirTemp("", "instrument-base-")
	must(err)
	defer os.RemoveAll(base)
	must(seed(repo, base))

	fmt.Println("control pass on an untouched copy…")
	for name := range gates {
		if *only != "" && *only != name {
			continue
		}
		ok, out := run(name, base)
		if !ok {
			fmt.Fprintf(os.Stderr, "\nCONTROL FAILED: %s is red on an untouched tree.\n"+
				"The harness would measure its own noise rather than the mutations.\n%s", name, out)
			os.Exit(1)
		}
	}

	// ── mutations ──────────────────────────────────────────────────────────
	fmt.Println()
	width := 0
	for _, m := range mutations {
		if n := len([]rune(m.name)); n > width {
			width = n
		}
	}

	missed, ran := 0, 0
	for _, m := range mutations {
		if *only != "" && *only != m.gate {
			continue
		}
		ran++
		tree, err := os.MkdirTemp("", "instrument-mut-")
		must(err)
		must(seed(repo, tree))

		applied, err := apply(filepath.Join(tree, m.file), m.from, m.to)
		if err != nil || !applied {
			os.RemoveAll(tree)
			fmt.Printf("  %-*s  ✗ MUTATION DID NOT APPLY (%s)\n", width, m.name, m.file)
			fmt.Printf("      the text looked for is not in the file — the harness has drifted from the kit\n")
			missed++
			continue
		}

		ok, out := run(m.gate, tree)

		pad := strings.Repeat(" ", width-len([]rune(m.name)))
		if ok {
			missed++
			fmt.Printf("  %s%s  ✗ MISSED    by gate %s\n", m.name, pad, m.gate)
			fmt.Printf("      %s\n", m.why)
			// The output of the GREEN gate is shown just like the output of a
			// red one: the question "why did it not catch this" is answered by
			// the gate's own words, and silence in this place sends you off to
			// investigate blind.
			if *verbose {
				for _, l := range strings.Split(strings.TrimSpace(out), "\n") {
					fmt.Printf("      %s\n", l)
				}
			}
			// The tree is kept under -keep. A missed mutation is the question
			// "what did the gate see", and only the tree itself can answer it.
			if *keep {
				fmt.Printf("      tree: %s\n", tree)
				continue
			}
		} else {
			fmt.Printf("  %s%s  · caught    %s\n", m.name, pad, m.gate)
			if *verbose {
				for _, l := range strings.Split(strings.TrimSpace(out), "\n") {
					fmt.Printf("      %s\n", l)
				}
			}
		}
		os.RemoveAll(tree)
	}

	fmt.Println()
	if missed > 0 {
		fmt.Printf("── holes in the gates: %d of %d ──\n", missed, ran)
		fmt.Println("A missed mutation means the invariant is declared but not guarded.")
		os.Exit(1)
	}
	fmt.Printf("· %d mutations, all caught: every invariant on the list is guarded\n", ran)
}

// seed lays out in a temporary directory the minimal tree every gate needs:
// the sources, the documentation, the registry, the version and the manifest.
func seed(repo, tree string) error {
	for _, d := range []string{"src", "docs"} {
		if err := copyTree(filepath.Join(repo, d), filepath.Join(tree, d)); err != nil {
			return err
		}
	}
	// The site build is built and run FROM THE COPY: it carries its styles and
	// templates inside the binary through go:embed, and a mutation in them
	// shows up only for a rebuilt one. tools/ is copied along with it — the
	// comment gate looks for it next to the kit, and without it the zone the
	// hole was opened for would be outside again.
	//
	// site/dist is skipped: nine megabytes of built site are of no use to
	// anyone, and copying them for each of forty-odd mutations means paying for
	// that every time.
	for _, d := range []string{"site", "tools"} {
		src := filepath.Join(repo, d)
		if err := copyTreeExcept(src, filepath.Join(tree, d), filepath.Join(src, "dist")); err != nil {
			return err
		}
	}
	// Besides src/ and docs/ the pages point with their source field at two
	// more files — the built CSS and the pixel reader. Without them docscheck
	// honestly fails on an UNTOUCHED copy, and the control pass catches it: a
	// snapshot missing something would measure its own incompleteness.
	for _, f := range []string{
		"components.json", "VERSION", "package.json",
		"tools/audit.js", "dist/instrument.min.css",
		// Zones of migration step 1: the built module and the marks. The lang
		// gate looks at the OUTPUT rather than at the source — a Russian string
		// that never reached dist/ is none of a consumer's business.
		"dist/instrument.js",
		"assets/logo.svg", "assets/mark.svg", "assets/sprite.svg",
	} {
		if err := copyFile(filepath.Join(repo, f), filepath.Join(tree, filepath.FromSlash(f))); err != nil {
			return err
		}
	}
	// The example stage lives in site/, and docscheck reads it for the
	// properties of the demo scene. It is copied under its own name so that the
	// whole site directory does not have to come along.
	return copyFile(filepath.Join(repo, "site", "internal", "render", "assets", "docs.css"),
		filepath.Join(tree, "stage.css"))
}

func apply(path, from, to string) (bool, error) {
	// An empty `from` is a mutation BY A NEW FILE: the gate has to see a page
	// that was not in the tree. An already existing file here means the same as
	// replacement text that was not found — the harness has drifted from the
	// kit — so this is a refusal rather than a silent overwrite.
	if from == "" {
		if _, err := os.Stat(path); err == nil {
			return false, nil
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return false, err
		}
		return true, os.WriteFile(path, []byte(to), 0o644)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}
	s := strings.ReplaceAll(string(b), "\r\n", "\n")
	if !strings.Contains(s, from) {
		return false, nil
	}
	return true, os.WriteFile(path, []byte(strings.Replace(s, from, to, 1)), 0o644)
}

func copyTree(src, dst string) error {
	return copyTreeExcept(src, dst, "")
}

func copyTreeExcept(src, dst, skip string) error {
	return filepath.WalkDir(src, func(p string, d os.DirEntry, err error) error {
		if skip != "" && p == skip {
			return filepath.SkipDir
		}
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, p)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		return copyFile(p, target)
	})
}

func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

func exeSuffix() string {
	if os.Getenv("OS") == "Windows_NT" || filepath.Separator == '\\' {
		return ".exe"
	}
	return ""
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

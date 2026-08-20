// The registry command cross-checks the component registry against the kit.
//
// WHY IT EXISTS. Cross-cutting axes — print, high contrast, flow — are recorded
// in CSS by enumerating component names: eight lists across three files. Each
// list must be edited by hand when adding a component, a forgotten entry causes
// no browser error, and none of the five previous checks could see it.
// A comment in layout.css describes how one such list already diverged, and
// the hole was found in the application, not in the kit.
//
// WHAT LIVES IN THE REGISTRY. Only what cannot be derived from the source. Which
// classes exist, which file they are in, which tokens they read, whether there
// is a page — all of that is obtained by traversal and is absent from the
// registry: a duplicated fact silently diverges from the original. What remains
// is the author's decision: whether the component is printed, whether it breaks
// between pages, what forced-colors does to it, whether it is measured by its
// content.
//
// WHAT IS NOT AND WILL NOT BE IN THE REGISTRY. prefers-reduced-motion slowing
// is DERIVED: only things with infinite animation need to be slowed. The check
// below compares the motion.css list with that derivation and does not ask the
// registry at all. An axis that can be derived is not added — otherwise the
// registry becomes an eighth list that must be remembered.
//
// THE "CLASS -> COMPONENT" MAPPING is not stored either: it is already declared
// in page frontmatter, and docscheck guards its completeness in both directions.
// The registry names the component, while the documentation names its classes.
//
//	go run ./cmd/registry
//	go run ./cmd/registry -v   with a list of what was parsed
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// ── axis dictionaries ───────────────────────────────────────────────────────
//
// Closed, like the kit's attribute dictionaries: a value outside the dictionary
// causes no error anywhere, it simply does nothing. And "none" here is a VALUE,
// not the absence of a key: omission and an intentional "nothing needed" must
// remain distinguishable, otherwise the gate becomes a wish.
var vocab = map[string][]string{
	"flow":   {"inline", "block"},
	"print":  {"keep", "hide"},
	"page":   {"flow", "whole"},
	"forced": {"none", "border", "fill", "selection", "glyph"},
}

var axisOrder = []string{"flow", "print", "page", "forced"}

// The axis value for which CSS should contain NOTHING. The reverse check
// relies on it: a class included in the list at the base value is an error
// just as much as a class missing from the list at the declared policy.
var base = map[string]string{"flow": "block", "print": "keep", "page": "flow", "forced": "none"}

type entry struct {
	Axes  map[string]string            `json:"axes"`
	Aria  *aria                        `json:"aria,omitempty"`
	Parts map[string]map[string]string `json:"parts,omitempty"`

	// Neighbours of the component, by registry name. The reference renders the
	// "Related" section from this and nothing else, so a name that resolves to
	// nothing produces a page with a link missing and no word about it.
	Related []string `json:"related,omitempty"`
}

// Markup contract: a role is a promise, and an unfulfilled promise is worse
// than an undeclared one. This contains what the MARKUP AUTHOR must write;
// what instrument.js does on their behalf is stated by the roving field.
type aria struct {
	On       string   `json:"on"` // class carrying the role, if not the root
	Role     string   `json:"role"`
	Roving   string   `json:"roving"` // author · js
	Requires []string `json:"requires"`
	Item     *struct {
		Role     string   `json:"role"`
		State    string   `json:"state"`
		Requires []string `json:"requires"`
	} `json:"item"`
	Exceptions map[string]struct {
		When string `json:"when"`
		Why  string `json:"why"`
	} `json:"exceptions"`
}

// ── CSS parsing ─────────────────────────────────────────────────────────────

type decl struct {
	sel, at, prop, value string
}

var commentRe = regexp.MustCompile(`(?s)/\*.*?\*/`)
var classRe = regexp.MustCompile(`\.(inst-[A-Za-z0-9_-]+)`)
var wsRe = regexp.MustCompile(`\s+`)

// walk parses a file into a flat list of declarations with the full selector
// and chain of at-rules. CSS nesting is handled with a stack: the kit uses it
// everywhere, and a rule inside `&:hover` must know its owner.
func walk(css string) []decl {
	css = commentRe.ReplaceAllString(css, "")
	var out []decl
	var stack []string
	var prelude strings.Builder

	flush := func() {
		d := strings.TrimSpace(prelude.String())
		prelude.Reset()
		if d == "" || len(stack) == 0 {
			return
		}
		i := strings.Index(d, ":")
		if i < 0 {
			return
		}
		prop := strings.TrimSpace(d[:i])
		if prop == "" || strings.ContainsAny(prop, "{}()") {
			return
		}
		var sel, at []string
		for _, s := range stack {
			if strings.HasPrefix(s, "@") {
				at = append(at, s)
			} else {
				sel = append(sel, s)
			}
		}
		out = append(out, decl{
			sel:   strings.Join(sel, " "),
			at:    strings.Join(at, " "),
			prop:  prop,
			value: strings.TrimSpace(d[i+1:]),
		})
	}

	for _, ch := range css {
		switch ch {
		case '{':
			stack = append(stack, wsRe.ReplaceAllString(strings.TrimSpace(prelude.String()), " "))
			prelude.Reset()
		case '}':
			flush()
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
			prelude.Reset()
		case ';':
			flush()
		default:
			prelude.WriteRune(ch)
		}
	}
	return out

}

func classesIn(sel string) []string {
	seen := map[string]bool{}
	var out []string
	for _, m := range classRe.FindAllStringSubmatch(sel, -1) {
		if !seen[m[1]] {
			seen[m[1]] = true
			out = append(out, m[1])
		}
	}
	return out
}

type set map[string]bool

func (s set) addAll(cs []string) {
	for _, c := range cs {
		s[c] = true
	}
}
func (s set) sorted() []string {
	out := make([]string, 0, len(s))
	for k := range s {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// ── what is declared in CSS ────────────────────────────────────────────────

// declared collects the actual lists from the cross-cutting files. It reads
// the real CSS, not a copy of the values: a copy silently diverges.
func declared(src map[string]string) map[string]set {
	out := map[string]set{
		"print:hide":       {},
		"page:whole":       {},
		"forced:border":    {},
		"forced:fill":      {},
		"forced:selection": {},
		"forced:glyph":     {},
		"flow:inline":      {},
	}

	for _, d := range walk(src["print.css"]) {
		switch {
		case d.prop == "display" && d.value == "none":
			out["print:hide"].addAll(classesIn(d.sel))
		case d.prop == "break-inside" && strings.Contains(d.value, "avoid"):
			out["page:whole"].addAll(classesIn(d.sel))
		}
	}

	for name, css := range src {
		for _, d := range walk(css) {
			if !strings.Contains(d.at, "forced-colors") {
				continue
			}
			cs := classesIn(d.sel)
			switch {
			case d.prop == "border" && strings.Contains(d.value, "CanvasText"):
				out["forced:border"].addAll(cs)
			case d.prop == "background" && strings.Contains(d.value, "CanvasText"):
				out["forced:glyph"].addAll(cs)
			case d.prop == "background" && d.value == "Highlight":
				if strings.Contains(d.sel, "aria-selected") {
					out["forced:selection"].addAll(cs)
				} else {
					out["forced:fill"].addAll(cs)
				}
			}
			_ = name
		}
	}

	// Flow is declared by the component itself: `align-self: var(--flow-self, …)`
	// sits next to its `display: inline-flex`. Print is excluded — it has its own
	// flow rules, unrelated to the axis.
	for name, css := range src {
		if name == "print.css" {
			continue
		}
		for _, d := range walk(css) {
			if d.prop == "align-self" && strings.Contains(d.value, "--flow-self") {
				out["flow:inline"].addAll(classesIn(d.sel))
			}
		}
	}

	// Priority: a class has one forced-colors policy, the most specific one.
	// Glyph is a mask, selection is a list row, fill is the semantic carrier,
	// border is simply a surface that needs an outline.
	for _, c := range out["forced:glyph"].sorted() {
		delete(out["forced:selection"], c)
		delete(out["forced:fill"], c)
		delete(out["forced:border"], c)
	}
	for _, c := range out["forced:selection"].sorted() {
		delete(out["forced:fill"], c)
		delete(out["forced:border"], c)
	}
	for _, c := range out["forced:fill"].sorted() {
		delete(out["forced:border"], c)
	}
	return out

}

// ── class -> component, from documentation frontmatter ─────────────────────

// Both spellings of the kind are listed by name, exactly as the alpha cell is
// handled in docscheck: the frontmatter is Russian until step 4 of the move,
// and a page translated ahead of this line would drop out of the ownership map
// without a word. The Russian alternative goes when the pages do.
var apiClassRe = regexp.MustCompile(`name:\s*"(inst-[a-z0-9-]+)"\s*,\s*kind:\s*"(?:класс|class)"`)

func ownership(docs string) (map[string]string, map[string][]string, error) {
	owner := map[string]string{}
	classes := map[string][]string{}
	var pages []string
	for _, dir := range []string{"components", "agent", "layout"} {
		err := filepath.WalkDir(filepath.Join(docs, dir), func(p string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return err
			}
			if strings.HasSuffix(p, ".md") && !strings.HasSuffix(p, ".en.md") {
				pages = append(pages, p)
			}
			return nil
		})
		if err != nil {
			return nil, nil, err
		}
	}
	// A page without classes in the api is not considered a component and does
	// not enter the registry. There is exactly one — the chart palette: it
	// documents --chart-* tokens, which have no cross-cutting axes and cannot
	// have any. Requiring an entry for it would mean creating four decisions
	// about things that do not exist.
	sort.Strings(pages)
	for _, p := range pages {
		b, err := os.ReadFile(p)
		if err != nil {
			return nil, nil, err
		}
		text := strings.ReplaceAll(string(b), "\r\n", "\n")
		i := strings.Index(text, "\n---")
		if !strings.HasPrefix(text, "---\n") || i < 0 {
			continue
		}
		fm := text[4 : i+1]
		name := strings.TrimSuffix(filepath.Base(p), ".md")
		for _, m := range apiClassRe.FindAllStringSubmatch(fm, -1) {
			classes[name] = append(classes[name], m[1])
			// The first page to name a class is its owner. A page may name
			// another page's class in its own api — the search shows the field,
			// the excerpt shows the property — but the owner alone declares
			// its policy, and only the owner does.
			if _, ok := owner[m[1]]; !ok {
				owner[m[1]] = name
			}
		}
	}
	return owner, classes, nil
}

// ── checks ─────────────────────────────────────────────────────────────────

func main() {
	var (
		src      = flag.String("src", "../src", "kit directory")
		docs     = flag.String("docs", "../docs", "documentation directory")
		registry = flag.String("registry", "../components.json", "component registry")
		verbose  = flag.Bool("v", false, "show parsed data")
	)
	flag.Parse()

	css := map[string]string{}
	ents, err := os.ReadDir(*src)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot read kit:", err)
		os.Exit(1)
	}
	for _, e := range ents {
		if strings.HasSuffix(e.Name(), ".css") {
			b, err := os.ReadFile(filepath.Join(*src, e.Name()))
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			css[e.Name()] = string(b)
		}
	}

	raw, err := os.ReadFile(*registry)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot read registry:", err)
		os.Exit(1)
	}
	var reg map[string]entry
	if err := json.Unmarshal(raw, &reg); err != nil {
		fmt.Fprintln(os.Stderr, "registry cannot be parsed:", err)
		os.Exit(1)
	}

	kitJSBytes, err := os.ReadFile(filepath.Join(*src, "kit.js"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot read instrument.js:", err)
		os.Exit(1)
	}
	kitJS := groupsBlockRe.FindString(string(kitJSBytes))

	owner, classes, err := ownership(*docs)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot read documentation:", err)
		os.Exit(1)
	}

	decl := declared(css)
	var problems []string
	checks := 0
	add := func(f string, a ...any) { problems = append(problems, "  · "+fmt.Sprintf(f, a...)) }

	// ── 1. completeness: component exists in both docs and registry ─────────
	for name := range classes {
		checks++
		if _, ok := reg[name]; !ok {
			add("component %q is documented but missing from the registry.\n      "+
				"Add an entry: without it, no cross-cutting axis watches it", name)
		}
	}
	for name := range reg {
		checks++
		if _, ok := classes[name]; !ok {
			add("registry knows %q, but there is no page with its classes", name)
		}
	}

	// ── 2. axis requirements: omission differs from intentional "nothing" ──
	for name, e := range reg {
		for _, ax := range axisOrder {
			checks++
			v, ok := e.Axes[ax]
			if !ok {
				add("%s: axis %q is not declared.\n      "+
					"If nothing is needed, write %q explicitly: omission and a decision must differ",
					name, ax, base[ax])
				continue
			}
			if !contains(vocab[ax], v) {
				add("%s: axis %q = %q is outside the dictionary (%s)", name, ax, v, strings.Join(vocab[ax], " · "))
			}
		}
		for part, ax := range e.Parts {
			checks++
			if owner[part] != name {
				if owner[part] == "" {
					add("%s: part %q is not named on any page", name, part)
				} else {
					add("%s: part %q belongs to component %q — it declares the policy for it",
						name, part, owner[part])
				}
			}
			for k, v := range ax {
				if !contains(vocab[k], v) {
					add("%s / %s: axis %q = %q is outside the dictionary", name, part, k, v)
				}
			}
		}
	}

	// ── 3. consistency: registry <-> CSS lists, in both directions ─────────
	//
	// Forward: declare a policy — the class must be in the list. Backward: class
	// is in the list — the policy must be declared. Without the second half,
	// the list quietly grows with things nobody decided on.
	want := map[string]string{} // "axis:value|class" -> component
	for name, e := range reg {
		root := ""
		if cs := classes[name]; len(cs) > 0 {
			root = cs[0]
		}
		for ax, v := range e.Axes {
			if v == base[ax] || root == "" {
				continue
			}
			want[ax+":"+v+"|"+root] = name
		}
		for part, ax := range e.Parts {
			for k, v := range ax {
				if v == base[k] {
					continue
				}
				want[k+":"+v+"|"+part] = name
			}
		}
	}
	for key, name := range want {
		i := strings.Index(key, "|")
		list, class := key[:i], key[i+1:]
		checks++
		if s, ok := decl[list]; !ok || !s[class] {
			add("%s declares %s for .%s, but it is not in that CSS list", name, strings.Replace(list, ":", "=", 1), class)
		}
	}
	for list, s := range decl {
		for _, class := range s.sorted() {
			checks++
			if _, ok := want[list+"|"+class]; !ok {
				who := owner[class]
				if who == "" {
					add("CSS assigns .%s to %s, but the class is not named on any page",
						class, strings.Replace(list, ":", "=", 1))
				} else {
					add("CSS assigns .%s to %s, but %s does not declare it",
						class, strings.Replace(list, ":", "=", 1), who)
				}
			}
		}
	}

	// ── 4. motion is DERIVED, not declared ─────────────────────────────────
	//
	// Under prefers-reduced-motion, slow exactly what has infinite animation:
	// an infinite animation compressed to 0.01ms is a stop, and a stopped
	// indicator says that work has stalled. The motion.css list must match this
	// derivation, and the registry is not involved here — an axis that can be
	// derived is not added.
	spins, slowed := set{}, set{}
	for name, text := range css {
		if name == "motion.css" || name == "print.css" {
			continue
		}
		for _, d := range walk(text) {
			if strings.HasPrefix(d.prop, "animation") && strings.Contains(d.value, "infinite") {
				spins.addAll(classesIn(d.sel))
			}
		}
	}
	for _, d := range walk(css["motion.css"]) {
		slowed.addAll(classesIn(d.sel))
	}
	for _, c := range spins.sorted() {
		checks++
		if !slowed[c] {
			add(".%s spins infinitely, but motion.css does not slow it.\n      "+
				"An infinite animation compressed to 0.01ms is a stop, and a stopped "+
				"indicator says that work has stalled", c)
		}
	}
	for _, c := range slowed.sorted() {
		checks++
		if !spins[c] {
			add("motion.css slows .%s, which has no infinite animation", c)
		}
	}

	// ── 5. kit.js role dictionary against registry ─────────────────────────
	//
	// kit.js has a closed list of roles with a keyboard contract, and above it
	// says that it matches the markup contract table. This match has so far
	// relied on memory — and has already diverged: the script handles keyboard
	// input for tablist, while there was no row for tabs in the table.
	//
	// A role silently handled by the script promises keyboard behavior on behalf
	// of the kit; a role declared in the registry without script support
	// promises it for nothing. Both directions are checked.
	jsRoles := set{}
	for _, m := range groupsRe.FindAllStringSubmatch(kitJS, -1) {
		jsRoles[m[1]] = true
	}
	regRoles := set{}
	for name, e := range reg {
		if e.Aria == nil || e.Aria.Item == nil {
			continue
		}
		regRoles[e.Aria.Role] = true
		checks++
		if e.Aria.Roving != "author" && e.Aria.Roving != "js" {
			add("%s: roving = %q is outside the dictionary (author · js)", name, e.Aria.Roving)
		}
	}
	for _, r := range jsRoles.sorted() {
		checks++
		if !regRoles[r] {
			add("instrument.js handles keyboard input for role=%q, but nobody declares this role in the registry.\n      "+
				"The role promises keyboard behavior on behalf of the kit — the promise must be recorded", r)
		}
	}
	for _, r := range regRoles.sorted() {
		checks++
		if !jsRoles[r] {
			add("registry declares role=%q group, but instrument.js does not handle it — the keyboard promise is empty", r)
		}
	}

	// ── 6. neighbours: names resolve, and the link goes both ways ───────────
	//
	// The reference renders "Related" from this and nothing else, so a name
	// pointing at nothing yields a page with a link silently missing.
	//
	// Reciprocity is required, and it is the whole reason the list moved into
	// the registry. Four pages of "Actions" each carried a hand-written list of
	// neighbours, and they had already drifted: the chip named the segmented
	// control, the segmented control did not name the chip. A one-sided
	// statement about two components is a statement that one of them has not
	// been told about.
	regNames := make([]string, 0, len(reg))
	for name := range reg {
		regNames = append(regNames, name)
	}
	sort.Strings(regNames)
	for _, name := range regNames {
		for _, other := range reg[name].Related {
			checks++
			o, ok := reg[other]
			if !ok {
				add("%s: related names %q, and the registry has no such component — the link would vanish from the page without a word", name, other)
				continue
			}
			back := false
			for _, r := range o.Related {
				if r == name {
					back = true
					break
				}
			}
			if !back {
				add("%s and %s: the link goes one way — %s names %s, %s does not name %s", name, other, name, other, other, name)
			}
		}
	}

	// ── 7. markup contract against live examples ────────────────────────────
	markup, markupChecks := checkMarkup(*docs, reg, classes)
	problems = append(problems, markup...)
	checks += markupChecks

	if *verbose {
		fmt.Printf("components in registry: %d\n", len(reg))
		fmt.Printf("classes named by documentation: %d\n", len(owner))
		for _, k := range sortedKeys(decl) {
			fmt.Printf("  %-18s %d: %s\n", k, len(decl[k]), strings.Join(decl[k].sorted(), " "))
		}
		fmt.Printf("  %-18s %d\n", "motion", len(spins))
	}

	fmt.Println()
	if len(problems) > 0 {
		sort.Strings(problems)
		fmt.Printf("── registry diverged from kit (%d) ──\n", len(problems))
		for _, p := range problems {
			fmt.Println(p)
		}
		fmt.Println()
		fmt.Printf("· %d checks, %d failed\n", checks, len(problems))
		os.Exit(1)
	}
	fmt.Printf("· %d checks: %d components, %d cross-cutting lists — registry and kit match\n",
		checks, len(reg), len(decl))

}

func contains(xs []string, x string) bool {
	for _, v := range xs {
		if v == x {
			return true
		}
	}
	return false
}

func sortedKeys(m map[string]set) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// Keys of the GROUPS dictionary in instrument.js — roles with a keyboard
// contract. The module itself is read, not a copied list: the copy is exactly
// what diverged.
var groupsRe = regexp.MustCompile(`(?m)^\s{2}([a-z]+):\s*\{`)

// The GROUPS block itself, not the whole module: a word like `menu:` occurs in
// the file outside the dictionary too, and without extracting the block a
// random line would enter the role list.
var groupsBlockRe = regexp.MustCompile(`(?s)const GROUPS = \{.*?\n\};`)

// ── Markup contract against live examples ─────────────────────────────────
//
// A role is a promise: `role="listbox"` tells assistive technology that arrow
// keys work. Until the promise is fulfilled, the component is not "unfinished";
// it is LYING. Nothing has checked the contract table until now: the words
// "aria" and "role" appeared in none of the six checks.
//
// Examples, not CSS, are checked, and this is the only correct place. An example
// is what the reader copies into their own code; an example missing a required
// attribute spreads the error. It is also the component's executable test:
// the kit has no other one, and adding a second would create two diverging
// descriptions.
//
// The parser is custom, about fifty lines, because tools/ intentionally has no
// dependencies. Example markup is handwritten and deliberately correct in
// form: a full parser here would solve a problem that does not exist.

var fenceRe = regexp.MustCompile("(?s)`html preview[^\n]*\n(.*?)`")

// tag — opening tag with its attributes and subtree boundaries.
type tag struct {
	name  string
	attrs map[string]string
	start int // index in tag list
	end   int // index of first tag AFTER the subtree
}

var voidTags = map[string]bool{
	"area": true, "base": true, "br": true, "col": true, "embed": true,
	"hr": true, "img": true, "input": true, "link": true, "meta": true,
	"source": true, "track": true, "wbr": true,
}

var tagRe = regexp.MustCompile(`<(/?)([a-zA-Z][\w-]*)((?:\s+[^<>"']+(?:"[^"]*")?)*)\s*(/?)>`)
var attrRe = regexp.MustCompile(`([a-zA-Z_:][-\w:.]*)(?:\s*=\s*"([^"]*)")?`)

// scan parses a fragment into a flat list of opening tags, each with its
// subtree boundary. This is enough for both "is the attribute on the element"
// and "which elements are inside it."
func scan(html string) []tag {
	var out []tag
	var stack []int
	for _, m := range tagRe.FindAllStringSubmatch(html, -1) {
		closing, name, rawAttrs, selfClose := m[1] == "/", strings.ToLower(m[2]), m[3], m[4] == "/"
		if closing {
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				out[top].end = len(out)
				if out[top].name == name {
					break
				}
			}
			continue
		}
		attrs := map[string]string{}
		for _, a := range attrRe.FindAllStringSubmatch(rawAttrs, -1) {
			if a[1] != "" {
				attrs[strings.ToLower(a[1])] = a[2]
			}
		}
		out = append(out, tag{name: name, attrs: attrs, start: len(out), end: -1})
		if !selfClose && !voidTags[name] {
			stack = append(stack, len(out)-1)
		} else {
			out[len(out)-1].end = len(out)
		}
	}
	for i := range out {
		if out[i].end < 0 {
			out[i].end = len(out)
		}
	}
	return out
}

func hasClass(t tag, c string) bool {
	for _, f := range strings.Fields(t.attrs["class"]) {
		if f == c {
			return true
		}
	}
	return false
}

// checkMarkup cross-checks page examples against the registry's ARIA contracts.
func checkMarkup(docs string, reg map[string]entry, classes map[string][]string) ([]string, int) {
	var problems []string
	checks := 0

	for name, e := range reg {
		if e.Aria == nil {
			continue
		}
		host := ""
		if cs := classes[name]; len(cs) > 0 {
			host = cs[0]
		}
		if e.Aria.On != "" {
			host = e.Aria.On
		}

		for _, page := range pagesOf(docs) {
			b, err := os.ReadFile(page)
			if err != nil {
				continue
			}
			text := strings.ReplaceAll(string(b), "\r\n", "\n")
			for _, f := range fenceRe.FindAllStringSubmatch(text, -1) {
				tags := scan(f[1])
				for i, t := range tags {
					if !hasClass(t, host) || t.attrs["role"] != e.Aria.Role {
						continue
					}
					// filepath.Rel, not TrimPrefix: the path comes from the flag
					// and may be written with any separator, while the message
					// must read the same way.
					short := page
					if rel, err := filepath.Rel(docs, page); err == nil {
						short = filepath.ToSlash(rel)
					}
					where := fmt.Sprintf("%s example with .%s", short, host)
					checks++
					for _, req := range e.Aria.Requires {
						if _, ok := t.attrs[req]; ok {
							continue
						}
						if ex, has := e.Aria.Exceptions[req]; has && matchesWhen(t, ex.When) {
							continue
						}
						problems = append(problems, fmt.Sprintf(
							"  · %s: no %s with role=%q.\n      The role promises this to assistive technology — without the attribute the promise is false",
							where, req, e.Aria.Role))
					}
					if e.Aria.Item == nil {
						continue
					}
					var items []tag
					for _, c := range tags[i+1 : t.end] {
						if c.attrs["role"] == e.Aria.Item.Role {
							items = append(items, c)
						}
					}
					if len(items) == 0 {
						problems = append(problems, fmt.Sprintf(
							"  · %s: role=%q without a single role=%q inside", where, e.Aria.Role, e.Aria.Item.Role))
						continue
					}
					checks++
					if st := e.Aria.Item.State; st != "" {
						for _, it := range items {
							if _, ok := it.attrs[st]; !ok {
								problems = append(problems, fmt.Sprintf(
									"  · %s: role=%q has no %s — state lives in an attribute, not a class",
									where, e.Aria.Item.Role, st))
								break
							}
						}
					}
					for _, req := range e.Aria.Item.Requires {
						for _, it := range items {
							if _, ok := it.attrs[req]; !ok {
								problems = append(problems, fmt.Sprintf(
									"  · %s: role=%q has no %s", where, e.Aria.Item.Role, req))
								break
							}
						}
					}
					// A roving tabindex is mandatory where the AUTHOR writes it.
					// Where instrument.js assigns it, nothing is required — but
					// that exemption is based on the script actually assigning it.
					// The premise is testable, and tools/behavior.js tests it in
					// the popover section: it opens the popover and reads the group,
					// WITHOUT setting focus. Checking after focus is meaningless —
					// the first Tab inside fixes the items itself, and the order is
					// visible where it is not. The exemption is used by the single
					// component with roving: "js", meaning 100% of its population.
					if e.Aria.Roving == "author" {
						checks++
						zero := 0
						for _, it := range items {
							if it.attrs["tabindex"] == "0" {
								zero++
							}
						}
						if zero != 1 {
							problems = append(problems, fmt.Sprintf(
								"  · %s: roving tabindex — %d items with tabindex=\"0\" out of %d, exactly one is required.\n"+
									"      Without it, Tab moves through every item, and the group stops being one control",
								where, zero, len(items)))
						}
					}
				}
			}
		}
	}
	sort.Strings(problems)
	return problems, checks

}

// matchesWhen checks an exception condition of the form [data-state="indeterminate"].
func matchesWhen(t tag, when string) bool {
	m := regexp.MustCompile(`\[([\w-]+)="([^"]*)"\]`).FindStringSubmatch(when)
	if m == nil {
		return false
	}
	return t.attrs[m[1]] == m[2]
}

var pageCache []string

func pagesOf(docs string) []string {
	if pageCache != nil {
		return pageCache
	}
	for _, dir := range []string{"components", "agent", "layout", "blocks"} {
		filepath.WalkDir(filepath.Join(docs, dir), func(p string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			// Translation is the same kind of page, and its markup is just as
			// live: the reader copies it from the reference and gets either a
			// working or broken component. Excluding `.en.md` would release the
			// ARIA contract exactly on pages that are written from scratch and
			// therefore make mistakes more often than the original.
			if strings.HasSuffix(p, ".md") {
				pageCache = append(pageCache, p)
			}
			return nil
		})
	}
	sort.Strings(pageCache)
	return pageCache
}

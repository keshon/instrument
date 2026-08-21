// The docscheck command checks the documentation against the kit. Both sides at once.
//
// The kit documentation can diverge from the kit in two directions, and both
// cost the same:
//
//	forward   an example contains a class that does not exist in the kit — the reader
//	          copies the markup, it silently does not work, and the kit gets blamed;
//	backward  the class exists in the kit, but has no page — the component exists
//	          and cannot be found, which for the consumer is equivalent to its absence.
//
// The check reads the actual kit sources and the actual pages, so it cannot
// diverge from them — the same principle as contrast.
//
//	go run ./cmd/docscheck        briefly
//	go run ./cmd/docscheck -v     with a list of undocumented classes
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var (
	commentRe = regexp.MustCompile(`(?s)/\*.*?\*/`)
	urlRe     = regexp.MustCompile(`url\([^)]*\)`)
	dqRe      = regexp.MustCompile(`"(?:[^"\\]|\\.)*"`)
	sqRe      = regexp.MustCompile(`'(?:[^'\\]|\\.)*'`)

	classRe = regexp.MustCompile(`\.(-?[A-Za-z_][\w-]*)`)
	declRe  = regexp.MustCompile(`(--[a-z][\w-]*)\s*:`)
	varUse  = regexp.MustCompile(`var\(\s*(--[a-z][\w-]*)`)
	dataSel = regexp.MustCompile(`\[data-([a-z-]+)="([^"]+)"\]`)

	// Strip everything from the selector that is NOT a bare element name:
	// classes, IDs, attributes, pseudo-classes together with their parentheses.
	selStrip = regexp.MustCompile(`::?[a-z-]+(\([^)]*\))?|\.[\w-]+|#[\w-]+|\[[^\]]*\]|[>+~,*]`)
	selBare  = regexp.MustCompile(`[a-z][a-z0-9]*`)

	// Elements that occur in the prose of the reference and inside examples
	// at the same time. The list is closed: the rule for <html> or <svg> does
	// not belong here, and expanding it randomly means catching false positives.
	htmlTags = map[string]bool{
		"table": true, "thead": true, "tbody": true, "tfoot": true, "tr": true,
		"th": true, "td": true, "caption": true,
		"p": true, "ul": true, "ol": true, "li": true, "dl": true, "dt": true, "dd": true,
		"blockquote": true, "pre": true, "code": true, "kbd": true, "samp": true,
		"a": true, "img": true, "figure": true, "figcaption": true, "hr": true,
		"h1": true, "h2": true, "h3": true, "h4": true, "h5": true, "h6": true,
		"button": true, "input": true, "select": true, "label": true, "fieldset": true,
		"legend": true, "details": true, "summary": true, "dialog": true,
	}

	// A class name, and NOT the tail of a longer hyphenated word. \b would
	// match inside data-inst-live — the hyphen counts as a boundary — and the
	// gate then reported .inst-live as a class the kit does not have, on a page
	// that was correctly naming an ATTRIBUTE the kit sets. RE2 has no lookbehind,
	// so the character before is captured and thrown away by the caller.
	instRe  = regexp.MustCompile(`(^|[^\w-])(inst-[a-z0-9-]+)`)
	linkRe  = regexp.MustCompile(`\]\((\.[^)#]+\.md)(#[^)]*)?\)`)
	tokRe   = regexp.MustCompile(`--[a-z][\w-]*(?:/[a-z0-9-]+)*`)
	dataAtt = regexp.MustCompile(`data-([a-z-]+)="([^"]*)"`)

	// Attributes that kit.js READS. They have no selector because they have
	// no styling: data-copy and data-value are data sources for behavior, not
	// state. They are found in the module itself rather than listed by hand:
	// the list would diverge from the kit at the first edit.
	jsDataset = regexp.MustCompile(`dataset\.([a-zA-Z]+)`)
)

// camelToDash converts a name from dataset back to an attribute name:
// copiedLabel -> copied-label.
func camelToDash(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			b.WriteByte('-')
			r += 'a' - 'A'
		}
		b.WriteRune(r)
	}
	return b.String()
}

// strip removes what is not a selector: comments, strings, url().
//
// Naive file-wide searching produces garbage: .w3 and .org come from
// xmlns='[http://www.w3.org/2000/svg](http://www.w3.org/2000/svg)' inside data-URI, while .tokens and .base
// come from file names in comments.
func strip(s string) string {
	s = commentRe.ReplaceAllString(s, " ")
	s = urlRe.ReplaceAllString(s, " ")
	s = dqRe.ReplaceAllString(s, " ")
	return sqRe.ReplaceAllString(s, " ")
}

func main() {
	var (
		srcDir  = flag.String("kit", "../src", "kit directory")
		docsDir = flag.String("docs", "../docs", "documentation directory")
		stage   = flag.String("stage", "../site/internal/render/assets/docs.css", "example scene styles")
		verbose = flag.Bool("v", false, "show list of classes without a page")
	)
	flag.Parse()

	// ── What exists in the kit ────────────────────────────────────────────
	kit := map[string]bool{}
	var allCSS strings.Builder
	entries, err := os.ReadDir(*srcDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot read kit:", err)
		os.Exit(1)
	}
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".css") {
			continue
		}
		b, err := os.ReadFile(filepath.Join(*srcDir, e.Name()))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		allCSS.Write(b)
		allCSS.WriteByte('\n')
		for _, line := range strings.Split(strip(string(b)), "\n") {
			i := strings.Index(line, "{")
			if i < 0 {
				continue
			}
			// "At-rules" are not selectors, and a dot inside them is not a class.
			// "@layer kit.components {" produced a nonexistent .components,
			// and it looked like a class without a page, that is, a real
			// finding. A false positive in the check costs more than a missed one:
			// misses get investigated, false positives get trusted.
			if strings.HasPrefix(strings.TrimSpace(line), "@") {
				continue
			}
			for _, m := range classRe.FindAllStringSubmatch(line[:i], -1) {
				kit[m[1]] = true
			}
		}
	}
	css := allCSS.String()

	// Tokens. A fabricated token is just as silent as a fabricated class: var()
	// on a nonexistent name silently returns empty, and the page lies without a sound.
	tokens := map[string]bool{}
	for _, m := range declRe.FindAllStringSubmatch(css, -1) {
		tokens[m[1]] = true
	}
	// Variables that the kit only READS, while markup sets them: tree node depth,
	// ring fraction, row number. They intentionally have no declarations.
	for _, m := range varUse.FindAllStringSubmatch(css, -1) {
		tokens[m[1]] = true
	}
	// The example stage has its own properties — --c on the color sample. They are
	// not kit tokens and never will be, but they are not fabricated either.
	if b, err := os.ReadFile(*stage); err == nil {
		for _, m := range declRe.FindAllStringSubmatch(string(b), -1) {
			tokens[m[1]] = true
		}
		for _, m := range varUse.FindAllStringSubmatch(string(b), -1) {
			tokens[m[1]] = true
		}
	}

	// Modifier suffixes: inst-btn--danger produces danger. Needed so that
	// "variant --danger" in prose is not treated as a nonexistent token.
	modifiers := map[string]bool{}
	for c := range kit {
		if i := strings.Index(c, "--"); i > 0 {
			modifiers[c[i+2:]] = true
		}
	}

	// data-attribute dictionaries — from selectors, not from memory.
	vocab := map[string]map[string]bool{}
	for _, m := range dataSel.FindAllStringSubmatch(css, -1) {
		if vocab[m[1]] == nil {
			vocab[m[1]] = map[string]bool{}
		}
		vocab[m[1]][m[2]] = true
	}
	// Behavior attributes: their values are free-form, only the name is checked.
	free := map[string]bool{}
	if b, err := os.ReadFile(filepath.Join(*srcDir, "kit.js")); err == nil {
		for _, m := range jsDataset.FindAllStringSubmatch(string(b), -1) {
			free[camelToDash(m[1])] = true
		}
	}

	// Base state values have no styling, so they are not present in selectors.
	// They are declared by the constitution — without them the check would
	// complain about valid markup.
	for _, v := range []string{"queued", "todo", "approved"} {
		if vocab["state"] == nil {
			vocab["state"] = map[string]bool{}
		}
		vocab["state"][v] = true
	}

	// ── What exists in the documentation ─────────────────────────────────
	var pages []string
	filepath.WalkDir(*docsDir, func(p string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if strings.Contains(filepath.ToSlash(p), "/internal/") || !strings.HasSuffix(p, ".md") {
			return nil
		}
		pages = append(pages, p)
		return nil
	})
	sort.Strings(pages)

	// Paths are printed relative to the repository root, not the working directory:
	// "docs/button.md" is equally clear in the output and in the editor, regardless
	// of where the check was launched from.
	root := filepath.Dir(*docsDir)
	rootRel := func(p string) string {
		if r, err := filepath.Rel(root, p); err == nil {
			return filepath.ToSlash(r)
		}
		return filepath.ToSlash(p)
	}

	documented := map[string]bool{}
	var problems []string
	pending := map[string]int{}

	// ── The reference has no right to RESTYLE the kit ─────────────────────
	//
	// The live specification promises that examples on its pages are working
	// markup, not pictures. The markup did work; the rendering lied.
	// `docs.css` held the rule `.site-body .inst-page-title { font-size:
	// --text-2xl }`, and the screen heading INSIDE the live example rendered at
	// 27px instead of the kit's 21px. Anyone copying what they saw got something else.
	//
	// Class coverage cannot see this hole by design: the class is documented,
	// the page exists, the checkmark is there. So this check is separate and
	// looks at PROPERTIES: the site may give a kit class spacing (that is the
	// layout of its own article), but it cannot touch what the component drawing
	// consists of — font size, line height, color, fill, border, radius.
	//
	// Example scene rules (`.demo-root`, `.demo-stage`) are excluded: they are
	// the frame around the component, not the component itself.
	if b, err := os.ReadFile(*stage); err == nil {
		lines := strings.Split(strip(string(b)), "\n")
		banned := []string{"font-size", "line-height", "font-weight", "color:",
			"background", "border", "border-radius", "letter-spacing"}
		rel := rootRel(*stage)
		for i, line := range lines {
			br := strings.Index(line, "{")
			if br < 0 {
				continue
			}
			sel := line[:br]
			if !strings.Contains(sel, ".inst-") {
				continue
			}
			if strings.Contains(sel, ".demo-root") || strings.Contains(sel, ".demo-stage") {
				continue
			}
			body := line[br:]
			// A single-line rule carries its body here; a block rule carries it
			// below, up to "}".
			if !strings.Contains(body, "}") {
				for j := i + 1; j < len(lines) && !strings.Contains(lines[j], "}"); j++ {
					body += lines[j]
				}
			}
			for _, prop := range banned {
				if strings.Contains(body, prop) {
					problems = append(problems, fmt.Sprintf(
						"%s:%d  reference restyles kit: %q sets %s",
						rel, i+1, strings.TrimSpace(sel), strings.TrimSuffix(prop, ":")))
					break
				}
			}
		}

		// ── Second door: BARE ELEMENT SELECTOR ─────────────────────────────
		//
		// The check above looks for rules NAMING a kit class. `.site-body td`
		// gets past it because `td` names no kit class, — and still restyles
		// every table inside every live example, because site styles live outside
		// layers and override the kit.
		//
		// Measurement: cell padding in the "Dashboard" example remained 8px
		// at all three densities, although --row-pad-y honestly changed 6 → 4 → 12.
		// Density did not work on the table, while the page promised otherwise.
		//
		// The door is closed not by attentiveness, but by syntax: a prose rule
		// touching the drawing must exclude the example subtree through
		// `:not(.demo-stage *)`. Then it physically cannot reach the kit.
		for i, line := range lines {
			br := strings.Index(line, "{")
			if br < 0 {
				continue
			}
			sel := strings.TrimSpace(line[:br])
			if !strings.HasPrefix(sel, ".site-body") || strings.Contains(sel, ".demo-stage *") {
				continue
			}
			// Bare element names: everything left after stripping classes,
			// pseudo-classes, attributes, and combinators.
			bare := selBare.FindAllString(selStrip.ReplaceAllString(sel, " "), -1)
			hit := ""
			for _, w := range bare {
				if htmlTags[w] {
					hit = w
					break
				}
			}
			if hit == "" {
				continue
			}
			body := line[br:]
			if !strings.Contains(body, "}") {
				for j := i + 1; j < len(lines) && !strings.Contains(lines[j], "}"); j++ {
					body += lines[j]
				}
			}
			for _, prop := range banned {
				if strings.Contains(body, prop) {
					problems = append(problems, fmt.Sprintf(
						"%s:%d  prose rule reaches inside live example: %q sets %s on <%s>. Add :not(.demo-stage *)",
						rel, i+1, sel, strings.TrimSuffix(prop, ":"), hit))
					break
				}
			}
		}
	}

	skipAttr := map[string]bool{"theme": true, "density": true, "dir": true}

	promises := 0

	for _, p := range pages {
		b, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		rel := rootRel(p)

		// What the page PROMISES about markup, against what it SHOWS.
		kept, n := promisesKept(func(ln int) string { return fmt.Sprintf("%s:%d", rel, ln) }, string(b))
		problems = append(problems, kept...)
		promises += n

		for i, line := range strings.Split(string(b), "\n") {
			at := fmt.Sprintf("%s:%d", rel, i+1)

			for _, m := range instRe.FindAllStringSubmatch(line, -1) {
				name := m[2]
				documented[name] = true
				if !kit[name] {
					problems = append(problems, at+"  class not in kit: ."+name)
				}
			}

			// Links between pages. Not an "error" and not a "metric", but a third
			// kind: a link to a page not written yet is a normal trace of unfinished
			// work, but a forgotten one is a hole nobody notices until the reader
			// falls into it.
			for _, m := range linkRe.FindAllStringSubmatch(line, -1) {
				target := filepath.Join(filepath.Dir(p), filepath.FromSlash(m[1]))
				if _, err := os.Stat(target); err != nil {
					pending[rootRel(target)]++
				}
			}

			// Tokens. Abbreviations like --text-xs/sm/md are expanded: this is
			// the accepted notation on pages, not a typo.
			for _, raw := range tokRe.FindAllString(line, -1) {
				if j := strings.Index(line, raw); j > 0 {
					prev := line[j-1]
					if prev == '-' || prev == '_' || (prev >= 'a' && prev <= 'z') ||
						(prev >= 'A' && prev <= 'Z') || (prev >= '0' && prev <= '9') {
						continue // part of a class name: inst-btn--primary
					}
				}
				parts := strings.Split(raw, "/")
				head := parts[0]
				// Two valid forms that are not tokens:
				// a class modifier in prose ("variant --danger") and
				// a truncated family abbreviation ("--space-*").
				if strings.HasSuffix(head, "-") || modifiers[strings.TrimPrefix(head, "--")] {
					continue
				}
				stem := head
				if k := strings.LastIndex(head, "-"); k > 1 {
					stem = head[:k]
				}
				names := []string{head}
				for _, t := range parts[1:] {
					names = append(names, stem+"-"+t)
				}
				for _, n := range names {
					if !tokens[n] {
						problems = append(problems, at+"  token not in kit: "+n)
					}
				}
			}

			for _, m := range dataAtt.FindAllStringSubmatch(line, -1) {
				attr, val := m[1], m[2]
				if skipAttr[attr] {
					continue
				}
				if free[attr] {
					continue
				}
				known, ok := vocab[attr]
				if !ok {
					problems = append(problems, at+"  unknown attribute: data-"+attr)
					continue
				}
				if val != "" && !known[val] {
					var list []string
					for v := range known {
						list = append(list, v)
					}
					sort.Strings(list)
					problems = append(problems,
						fmt.Sprintf("%s  data-%s=%q — not in dictionary (%s)", at, attr, val, strings.Join(list, " ")))
				}
			}
		}
	}

	// ── Report ────────────────────────────────────────────────────────────
	var undocumented []string
	for c := range kit {
		if !documented[c] {
			undocumented = append(undocumented, c)
		}
	}
	sort.Strings(undocumented)

	covered := len(kit) - len(undocumented)
	pct := 0
	if len(kit) > 0 {
		pct = covered * 100 / len(kit)
	}
	fmt.Printf("pages: %d  ·  classes in kit: %d  ·  covered: %d (%d%%)\n\n",
		len(pages), len(kit), covered, pct)
	if promises > 0 {
		fmt.Printf("markup promises held to the examples: %d\n\n", promises)
	}

	// Token tables: compare VALUES, not existence. Kept as a separate list,
	// because this is a different kind of error: class is present, token is present,
	// only the number lies.
	staleTable := false
	if miss := checkSourceFields(*docsDir); len(miss) > 0 {
		fmt.Printf("── source leads nowhere (%d) ──\n", len(miss))
		for _, s := range miss {
			fmt.Println("  " + s)
		}
		fmt.Println()
		staleTable = true
	}

	if stale := checkModeTables(*srcDir, *docsDir); len(stale) > 0 {
		fmt.Printf("── mode tables diverged from code (%d) ──\n", len(stale))
		for _, s := range stale {
			fmt.Println("  " + s)
		}
		fmt.Println()
		staleTable = true
	}

	if stale := checkTokenTables(*srcDir, *docsDir); len(stale) > 0 {
		fmt.Printf("── token table diverged from code (%d) ──\n", len(stale))
		for _, s := range stale {
			fmt.Println("  " + s)
		}
		fmt.Println()
		staleTable = true
	}

	if len(problems) > 0 {
		seen := map[string]bool{}
		var uniq []string
		for _, p := range problems {
			if !seen[p] {
				seen[p] = true
				uniq = append(uniq, p)
			}
		}
		fmt.Printf("── documentation references nonexistent items (%d) ──\n", len(uniq))
		for _, p := range uniq {
			fmt.Println("  " + p)
		}
		fmt.Println()
	}

	if len(pending) > 0 {
		var keys []string
		for k := range pending {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		fmt.Printf("── links to pages not yet written (%d) ──\n", len(keys))
		for _, k := range keys {
			fmt.Printf("  %s  ← %d\n", k, pending[k])
		}
		fmt.Println()
	}

	if len(undocumented) > 0 {
		if *verbose {
			fmt.Printf("── without a page (%d) ──\n", len(undocumented))
			for i, c := range undocumented {
				undocumented[i] = "." + c
			}
			fmt.Println("  " + strings.Join(undocumented, " ") + "\n")
		} else {
			fmt.Printf("── without a page: %d classes (list with -v) ──\n\n", len(undocumented))
		}
	}

	if len(problems) == 0 && len(undocumented) == 0 && len(pending) == 0 && !staleTable {
		fmt.Println("· documentation and kit match completely")
	}

	if len(problems) > 0 || staleTable {
		os.Exit(1)
	}

}

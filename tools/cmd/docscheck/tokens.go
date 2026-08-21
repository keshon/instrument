// Verify the TOKEN TABLES in the reference against declarations in tokens.css.
//
// The main part of docscheck verifies EXISTENCE: whether a class exists, whether
// a token exists, whether a page exists. It cannot catch the most common kind of
// mismatch — when documentation names a value that the token no longer has.
//
// The example this check was written for: the surface table claimed that
// --surface-recessed was "black 5%", while the code already had 0.060. The
// class was present, the token was present, docscheck was silent — only the
// number was wrong.
//
// Only machine-readable cells are checked. A row like
//
//	| `--surface-sunken` | `--n-2` | `--n-14` |
//
// requires the code to contain light-dark(var(--n-2), var(--n-14)); a row like
//
//	| `--surface-hover` | black 3.5% | white 4.5% |
//
// requires light-dark(oklch(0 0 0 / 0.035), oklch(1 0 0 / 0.045)).
//
// A row whose token is declared through light-dark is checked AS A WHOLE: both
// of its cells must be machine-readable. Skipping an unrecognized cell would
// mean accepting "black 9%" as prose and staying silent: an unrecognized word
// carries away the number next to it.
package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"instrument/tools/internal/css"
)

var (
	// a table row whose first cell names a token: | `--name` | … |
	//
	// WHICH cells hold the two themes is decided by the HEAD of the table, not
	// by their position. While the row had to be exactly three cells wide and
	// the values had to be the second and the third, two shapes fell out of the
	// check without a word:
	//
	//	| Token | Light | Dark | Threshold |   ← a fourth column of prose
	//	| Token | Hue | Light | Dark |         ← the values shifted right
	//
	// Everything in those two shapes was skipped, and --focus-ring stood in the
	// reference as `--a-4` / `--a-3` long after the code had moved it to
	// `--a-6` / `--a-2`. The check was green because it never looked.
	rowRe     = regexp.MustCompile("^\\|\\s*`(--[a-z][\\w-]*)`\\s*\\|")
	headSep   = regexp.MustCompile(`^\|(?:\s*:?-{3,}:?\s*\|)+\s*$`)
	lightHead = regexp.MustCompile(`(?i)^(светлая|light)$`)
	darkHead  = regexp.MustCompile(`(?i)^(тёмная|темная|dark)$`)
	// A hue column belongs to neither theme: both ends of a light-dark carry
	// the same hue, and the reader is promised exactly that one number. The
	// chart table stood at 280° after --chart-1 had moved to 292 — the row was
	// read, the two lightnesses matched, and the number the palette is chosen
	// by was the one nobody compared.
	hueHead = regexp.MustCompile(`(?i)^(тон|hue)$`)
	// A LIGHTNESS column belongs to the ramp tables, where a token is a plain
	// oklch rather than a light-dark. Those tables are the most factual thing
	// on the page — a step, its lightness, the role that took it — and they
	// were the ones nothing looked at: the accent ramp said 250° and 0.560
	// after the default had moved to petrol at 215° and 0.545.
	lumHead = regexp.MustCompile(`(?i)^(светлота|lightness)$`)
	// the first number of an oklch: its lightness. The chroma of the neutral is
	// a calc(), so the three-number form does not match there.
	oklchLum = regexp.MustCompile(`oklch\(\s*([0-9.]+)`)
	// a plain declaration: --name: oklch(…)
	plainRe = regexp.MustCompile(`(--[a-z][\w-]*):\s*(oklch\([^;]*)`)
	// a cell naming one component of an oklch: "L 0.520" or "280°". The chart
	// table is written that way, and both are checkable — the lightness and the
	// hue stand in the declaration next to each other.
	lumRe = regexp.MustCompile(`^L\s+([0-9.]+)$`)
	hueRe = regexp.MustCompile(`^([0-9.]+)°$`)
	// the head of a light-dark declaration; the tail is parsed by counting
	// parentheses, because oklch(...) sits inside with parentheses of its own
	ldHead = regexp.MustCompile(`(--[a-z][\w-]*):\s*light-dark\(`)
	// a cell that references another token
	refRe = regexp.MustCompile("^`(--[a-z][\\w-]*)`$")
	// a blend cell: `--ok-4` 16% — a role colour diluted with transparency
	refPctRe = regexp.MustCompile("^`(--[a-z][\\w-]*)`\\s+([0-9.]+)%$")
	// a cell of the form "black 3.5%" or "white 8%", in either language.
	//
	// There are exactly two words per language, and both are listed here by
	// name rather than reduced to "any word before a percentage". A third
	// language has to fail parsing: an unparsed cell is an error now, and that
	// is the only thing separating a translation from a typo.
	//
	// The Russian spellings are what docs/foundations/tokens.md still says.
	// They leave together with the page on step 4 of the move.
	alphaRe = regexp.MustCompile(`^(чёрный|белый|black|white)\s+([0-9.]+)%$`)
	// the three numbers of an oklch(), in order: lightness, chroma, hue
	oklchNums = regexp.MustCompile(`oklch\(\s*([0-9.]+)\s+([0-9.]+)\s+([0-9.]+)`)
	// numeric alpha, so that 0.05 and 0.050 count as one value
	alphaNum = regexp.MustCompile(`oklch\(([\d ]+)\s*/\s*([\d.]+)\)`)
)

// splitArgs splits the contents of light-dark(...) at a TOP-LEVEL comma.
// Returns nil if there are not two arguments.
func splitArgs(s string) []string {
	var out []string
	depth, start := 0, 0
	for i, r := range s {
		switch r {
		case '(':
			depth++
		case ')':
			depth--
		case ',':
			if depth == 0 {
				out = append(out, strings.TrimSpace(s[start:i]))
				start = i + 1
			}
		}
	}
	out = append(out, strings.TrimSpace(s[start:]))
	if len(out) != 2 {
		return nil
	}
	return out
}

// lightDarkPairs collects the first declaration of each token: this is :root,
// so it is the base cell. Theme and scale overrides are not described by the
// table.
func lightDarkPairs(css string) map[string][2]string {
	out := map[string][2]string{}
	for _, loc := range ldHead.FindAllStringSubmatchIndex(css, -1) {
		name := css[loc[2]:loc[3]]
		if _, seen := out[name]; seen {
			continue
		}
		depth, end := 1, -1
		for i := loc[1]; i < len(css); i++ {
			if css[i] == '(' {
				depth++
			} else if css[i] == ')' {
				depth--
				if depth == 0 {
					end = i
					break
				}
			}
		}
		if end < 0 {
			continue
		}
		if a := splitArgs(css[loc[1]:end]); a != nil {
			out[name] = [2]string{a[0], a[1]}
		}
	}
	return out
}

// expect translates a table cell into what must appear in the code.
// The second value is false if the cell is not machine-readable.
func expect(cell string) (string, bool) {
	cell = strings.TrimSpace(cell)
	if m := refRe.FindStringSubmatch(cell); m != nil {
		return "var(" + m[1] + ")", true
	}
	if m := refPctRe.FindStringSubmatch(cell); m != nil {
		pct, err := strconv.ParseFloat(m[2], 64)
		if err != nil {
			return "", false
		}
		return fmt.Sprintf("color-mix(in oklab, var(%s) %g%%, transparent)", m[1], pct), true
	}
	if m := alphaRe.FindStringSubmatch(cell); m != nil {
		pct, err := strconv.ParseFloat(m[2], 64)
		if err != nil {
			return "", false
		}
		lum := "0 0 0"
		if m[1] == "белый" || m[1] == "white" {
			lum = "1 0 0"
		}
		return fmt.Sprintf("oklch(%s / %g)", lum, pct/100), true
	}
	return "", false
}

// cells splits a markdown row into cells, dropping the empty edges the leading
// and trailing pipes produce.
func cells(line string) []string {
	line = strings.TrimSpace(line)
	line = strings.TrimPrefix(line, "|")
	line = strings.TrimSuffix(line, "|")
	return strings.Split(line, "|")
}

// component compares a cell that names ONE component of an oklch — the
// lightness or the hue — against the declaration. The chart table is written
// that way: a colour whose saturation nobody quotes still has a lightness and a
// hue, and both of those are what the row promises the reader.
func component(cell, got string) (ok, machine bool) {
	nums := oklchNums.FindStringSubmatch(got)
	if nums == nil {
		return false, false
	}
	if m := lumRe.FindStringSubmatch(cell); m != nil {
		return sameNum(m[1], nums[1]), true
	}
	if m := hueRe.FindStringSubmatch(cell); m != nil {
		return sameNum(m[1], nums[3]), true
	}
	return false, false
}

// checkLum holds a lightness cell of a ramp table to the declaration. A step
// the reference names but the code does not declare is an error too: the table
// is a list of the ramp, and a row about a step that is gone is worse than a
// missing row.
func checkLum(out *[]string, name string, line int, token, cell, got string) {
	if got == "" {
		*out = append(*out, fmt.Sprintf(
			"%s:%d: %s — the table names a step the code does not declare",
			name, line, token))
		return
	}
	nums := oklchLum.FindStringSubmatch(got)
	if nums == nil {
		return
	}
	if !sameNum(cell, nums[1]) {
		*out = append(*out, fmt.Sprintf(
			"%s:%d: %s, lightness — %q in the table, %q in the code",
			name, line, token, cell, nums[1]))
	}
}

func sameNum(a, b string) bool {
	x, err1 := strconv.ParseFloat(a, 64)
	y, err2 := strconv.ParseFloat(b, 64)
	return err1 == nil && err2 == nil && x == y
}

func sameValue(want, got string) bool {
	if strings.Join(strings.Fields(want), " ") == strings.Join(strings.Fields(got), " ") {
		return true
	}
	w := alphaNum.FindStringSubmatch(want)
	g := alphaNum.FindStringSubmatch(got)
	if w == nil || g == nil {
		return false
	}
	if strings.Join(strings.Fields(w[1]), " ") != strings.Join(strings.Fields(g[1]), " ") {
		return false
	}
	a, err1 := strconv.ParseFloat(w[2], 64)
	b, err2 := strconv.ParseFloat(g[2], 64)
	return err1 == nil && err2 == nil && a == b
}

// langVariants — the page and all its translations: tokens.md and tokens.en.md.
//
// Checking only the base name would mean releasing the tables on exactly those
// pages that are rewritten and therefore make mistakes more often than the
// original. The cost is known: the same filter in cmd/registry released half
// of the ARIA contract. It guards against the mutation "the token table in a
// translation is wrong".
func langVariants(dir, stem string) []string {
	base, _ := filepath.Glob(filepath.Join(dir, stem+".md"))
	tr, _ := filepath.Glob(filepath.Join(dir, stem+".*.md"))
	return append(base, tr...)
}

// checkTokenTables returns table/code mismatches.
func checkTokenTables(srcDir, docsDir string) []string {
	css, err := os.ReadFile(filepath.Join(srcDir, "tokens.css"))
	if err != nil {
		return nil
	}
	decl := lightDarkPairs(string(css))
	// The FIRST declaration only: --a-* is redeclared by every accent, and the
	// base ramp is what the reference tabulates.
	plain := map[string]string{}
	for _, m := range plainRe.FindAllStringSubmatch(string(css), -1) {
		if _, seen := plain[m[1]]; !seen {
			plain[m[1]] = m[2]
		}
	}

	var out []string
	// BOTH pages, not just the reference. The tables of the ramps live on the
	// colour page, and while only tokens.md was read they were the most factual
	// thing in the documentation with nothing checking it: the accent ramp
	// stood at 250° and 0.560 after the default had moved to petrol.
	var pages []string
	for _, stem := range []string{"tokens", "colors"} {
		pages = append(pages, langVariants(filepath.Join(docsDir, "foundations"), stem)...)
	}
	for _, page := range pages {
		out = append(out, tokenTableOf(page, decl, plain)...)
	}
	return out

}

func tokenTableOf(page string, decl map[string][2]string, plain map[string]string) []string {
	md, err := os.ReadFile(page)
	if err != nil {
		return nil
	}
	name := filepath.Base(page)

	var out []string
	lines := strings.Split(string(md), "\n")
	lightCol, darkCol, hueCol, lumCol := -1, -1, -1, -1
	for i, line := range lines {
		line = strings.TrimRight(line, "\r")

		// A separator row means the line above it was the head: read from it
		// which columns hold the two themes, and keep that until the next
		// table. A table with no such head leaves the pair at -1, and its rows
		// are skipped — a heading nobody recognised is not a licence to guess
		// at positions.
		if headSep.MatchString(line) && i > 0 {
			lightCol, darkCol, hueCol, lumCol = -1, -1, -1, -1
			for k, h := range cells(strings.TrimRight(lines[i-1], "\r")) {
				switch {
				case lightHead.MatchString(strings.TrimSpace(h)):
					lightCol = k
				case darkHead.MatchString(strings.TrimSpace(h)):
					darkCol = k
				case hueHead.MatchString(strings.TrimSpace(h)):
					hueCol = k
				case lumHead.MatchString(strings.TrimSpace(h)):
					lumCol = k
				}
			}
			continue
		}

		m := rowRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		if lumCol >= 0 {
			row := cells(line)
			if lumCol < len(row) {
				checkLum(&out, name, i+1, m[1], strings.TrimSpace(row[lumCol]), plain[m[1]])
			}
		}

		got, ok := decl[m[1]]
		if !ok {
			continue // the token is not declared through light-dark — not our case
		}
		if lightCol < 0 || darkCol < 0 {
			continue
		}
		row := cells(line)
		if lightCol >= len(row) || darkCol >= len(row) {
			continue
		}
		if hueCol >= 0 && hueCol < len(row) {
			cell := strings.TrimSpace(row[hueCol])
			if ok, machine := component(cell, got[0]); machine && !ok {
				out = append(out, fmt.Sprintf(
					"%s:%d: %s, hue — %q in the table, %q in the code",
					name, i+1, m[1], cell, got[0]))
			}
		}
		for k, cell := range []string{row[lightCol], row[darkCol]} {
			end := "light"
			if k == 1 {
				end = "dark"
			}
			if ok, machine := component(strings.TrimSpace(cell), got[k]); machine {
				if !ok {
					out = append(out, fmt.Sprintf(
						"%s:%d: %s, %s — %q in the table, %q in the code",
						name, i+1, m[1], end, strings.TrimSpace(cell), got[k]))
				}
				continue
			}
			want, machine := expect(cell)
			if !machine {
				out = append(out, fmt.Sprintf(
					"%s:%d: %s, %s — cannot parse the cell %q.\n      "+
						"The token is declared through light-dark, so the column is machine "+
						"readable: an unparsed cell stays quiet where it has to object",
					name, i+1, m[1], end, strings.TrimSpace(cell)))
				continue
			}
			if !sameValue(want, got[k]) {
				out = append(out, fmt.Sprintf(
					"%s:%d: %s, %s — %q in the table, %q in the code",
					name, i+1, m[1], end, strings.TrimSpace(cell), got[k]))
			}
		}
	}
	return out
}

// ── DENSITY AND SCALE TABLES ──────────────────────────────────────────────
//
// The check above compares two-column color tables: token, light theme, dark
// theme. density.md and scale.md have a different shape — "token × mode", and
// values live in [data-density] and [data-scale] blocks. There was no check for
// this shape, and it cost exactly what this check was written for: a change to
// --row-pad-y passed, while the table kept naming the previous step. The class
// was present, the token was present, only the number was wrong.
//
// Only reference cells (`--space-4`) are parsed. Numbers and prose are skipped:
// forcing the whole table to be machine-readable would hollow it out.

// A table in markup is recognized by the SEPARATOR ROW: it must be second, with
// the header one row above. The word in the first header cell means nothing:
// modeRowRe identifies the token column by backticks around the name. Identifying
// the header by the word "Token" would mean disabling the entire table check with
// one translated word — and continuing to print "fully consistent".
var sepRe = regexp.MustCompile(`^\|(?:\s*:?-{3,}:?\s*\|)+\s*$`)

// A backtick inside a Go raw literal does not work, so concatenate: it reads
// better than escaping every table delimiter.
const bq = "`"

var modeRowRe = regexp.MustCompile(`^\|\s*` + bq + `(--[a-z][\w-]*)` + bq + `\s*\|(.+)\|\s*$`)
var refCellRe = regexp.MustCompile(`^` + bq + `(--[a-z][\w-]*)` + bq + `$`)

// Own declaration regex: declRe in main.go has one capture group — it needs
// only the TOKEN NAME, while this one also needs the value.
var modeDeclRe = regexp.MustCompile(`(--[a-z][\w-]*)\s*:\s*([^;{}]+);`)

// blockDecls extracts declarations from the FIRST block whose selector matches.
func blockDecls(css, selector string) map[string]string {
	i := strings.Index(css, selector)
	if i < 0 {
		return nil
	}
	j := strings.Index(css[i:], "{")
	if j < 0 {
		return nil
	}
	depth, start := 0, i+j+1
	for k := i + j; k < len(css); k++ {
		switch css[k] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				out := map[string]string{}
				for _, m := range modeDeclRe.FindAllStringSubmatch(css[start:k], -1) {
					out[m[1]] = strings.TrimSpace(m[2])
				}
				return out
			}
		}
	}
	return nil
}

// Mode code is declared not here, but in the kit: [data-density="compact"],
// [data-scale="17"]. A list hardcoded into this command would duplicate the kit
// from memory and silently diverge from it. Worse: when matching a label against
// the LABEL dictionary, "dense" instead of "compact" would have to be treated
// as a non-mode column — and released entirely.
var modeSelRe = regexp.MustCompile(`(?m)^\[data-(?:density|scale)="([^"]+)"\] \{$`)

func modeSelectors(sheet string) map[string]string {
	out := map[string]string{}
	for _, m := range modeSelRe.FindAllStringSubmatch(sheet, -1) {
		out[m[1]] = m[0]
	}
	return out
}

// Parsed "token × mode" table.
type modeRow struct {
	line  int
	token string
	cells []string
}

type modeTable struct {
	line int      // header row
	head []string // column labels, without the token column
	rows []modeRow
}

// headCells splits the header into cells and removes backticks: `compact` and
// compact are the same label, and distinguishing them would mean requiring the
// table author to care about markup rather than meaning.
func headCells(line string) []string {
	line = strings.TrimSpace(line)
	line = strings.TrimSuffix(strings.TrimPrefix(line, "|"), "|")
	var out []string
	for _, c := range strings.Split(line, "|") {
		out = append(out, strings.Trim(strings.TrimSpace(c), bq))
	}
	return out
}

func modeTablesOf(md string) []modeTable {
	lines := strings.Split(strings.ReplaceAll(md, "\r\n", "\n"), "\n")
	var out []modeTable
	for i := 1; i < len(lines); i++ {
		if !sepRe.MatchString(lines[i]) {
			continue
		}
		head := headCells(lines[i-1])
		if len(head) < 2 {
			continue
		}
		t := modeTable{line: i, head: head[1:]}
		for j := i + 1; j < len(lines) && strings.HasPrefix(lines[j], "|"); j++ {
			m := modeRowRe.FindStringSubmatch(lines[j])
			if m == nil {
				continue
			}
			var cells []string
			for _, c := range strings.Split(m[2], "|") {
				cells = append(cells, strings.TrimSpace(c))
			}
			t.rows = append(t.rows, modeRow{j + 1, m[1], cells})
		}
		// A table without a single token row is not about modes: density.md and
		// scale.md have five such tables, from "Value · Purpose" to the 24px
		// threshold.
		if len(t.rows) > 0 {
			out = append(out, t)
		}
	}
	return out
}

// machineCols marks columns where at least one cell is machine-readable — a
// token reference or a number in px. A column of pure prose ("When", "Why")
// does not pretend to be a mode, and there is nothing to ask of it.
func machineCols(t modeTable) []bool {
	out := make([]bool, len(t.head))
	for _, r := range t.rows {
		for k, cell := range r.cells {
			if k >= len(out) || out[k] {
				continue
			}
			if refCellRe.MatchString(cell) {
				out[k] = true
				continue
			}
			if _, err := strconv.ParseFloat(strings.TrimSuffix(cell, "px"), 64); err == nil {
				out[k] = true
			}
		}
	}
	return out
}

// column — a table column resolved to a selector in the kit.
type column struct {
	head string
	sel  string // "" — base column: value lives in the role tier
	use  bool
}

// resolveCols maps columns to kit modes.
//
// One rule: a machine-readable column must be either a mode code from the kit,
// or the SINGLE unnamed one — the base. Two unnamed machine-readable columns
// mean one of them was not recognized, and that is an error, not a reason to
// stay silent. Silence here is exactly the useless report this check was
// written to prevent: without this rule, an unknown label produces no values
// at all, and the column disappears.
func resolveCols(name string, t modeTable, modeSel map[string]string) ([]column, []string) {
	machine := machineCols(t)
	cols := make([]column, len(t.head))
	var unnamed []int
	count := 0
	for k, h := range t.head {
		cols[k] = column{head: h}
		if !machine[k] {
			continue
		}
		count++
		if sel, ok := modeSel[h]; ok {
			cols[k].sel, cols[k].use = sel, true
			continue
		}
		unnamed = append(unnamed, k)
	}
	// One machine-readable column is not "token × mode", but a list: there is
	// nothing to compare it with, and no reason to require a mode code from it.
	if count < 2 {
		return nil, nil
	}
	if len(unnamed) > 1 {
		var names []string
		for _, k := range unnamed {
			names = append(names, "«"+t.head[k]+"»")
		}
		return nil, []string{fmt.Sprintf(
			"%s:%d: columns %s are not named with a mode code.\n      "+
				"The kit declares the mode code ([data-density], [data-scale]); "+
				"only one unnamed column — the base — is allowed",
			name, t.line, strings.Join(names, ", "))}
	}
	if len(unnamed) == 1 {
		cols[unnamed[0]].use = true
	}
	return cols, nil
}

func checkModeTables(srcDir, docsDir string) []string {
	raw, err := os.ReadFile(filepath.Join(srcDir, "tokens.css"))
	if err != nil {
		return nil
	}
	sheet := commentRe.ReplaceAllString(strings.ReplaceAll(string(raw), "\r\n", "\n"), "")
	modeSel := modeSelectors(sheet)

	// The role tier: the base a mode overrides its own values on top of.
	base := map[string]string{}
	for _, sel := range []string{":root {", `:where(:root) {`} {
		for k, v := range blockDecls(sheet, sel) {
			base[k] = v
		}
	}

	cache := map[string]map[string]string{}
	values := func(sel string) map[string]string {
		if v, ok := cache[sel]; ok {
			return v
		}
		out := map[string]string{}
		for k, v := range base {
			out[k] = v
		}
		if sel != "" {
			for k, v := range blockDecls(sheet, sel) {
				out[k] = v
			}
		}
		cache[sel] = out
		return out
	}

	var problems []string
	var pages []string
	for _, stem := range []string{"density", "scale"} {
		pages = append(pages, langVariants(filepath.Join(docsDir, "foundations"), stem)...)
	}
	for _, page := range pages {
		md, err := os.ReadFile(page)
		if err != nil {
			continue
		}
		name := filepath.Base(page)
		for _, t := range modeTablesOf(string(md)) {
			cols, bad := resolveCols(name, t, modeSel)
			if bad != nil {
				problems = append(problems, bad...)
				continue
			}
			for _, r := range t.rows {
				for k, cell := range r.cells {
					if k >= len(cols) || !cols[k].use {
						continue
					}
					got, ok := values(cols[k].sel)[r.token]
					if !ok {
						continue // the token is not declared in this mode
					}

					// A reference cell: the name is what gets compared.
					if ref := refCellRe.FindStringSubmatch(cell); ref != nil {
						if got != "var("+ref[1]+")" {
							problems = append(problems, fmt.Sprintf(
								"%s:%d: %s in mode %q — %q in the table, %q in the code",
								name, r.line, r.token, cols[k].head, ref[1], got))
						}
						continue
					}

					// A numeric cell: PIXELS are what gets compared. The scale tables
					// are set in px while the tokens are declared in rem — without
					// resolving them a hundred numbers would stay unchecked, exactly
					// as the density step once did.
					want, err := strconv.ParseFloat(strings.TrimSuffix(cell, "px"), 64)
					if err != nil {
						continue // prose
					}
					px, err := css.ResolvePx(values(cols[k].sel), r.token)
					if err != nil {
						continue // not a length
					}
					if math.Abs(px-want) > 0.01 {
						problems = append(problems, fmt.Sprintf(
							"%s:%d: %s in mode %q — %g in the table, %g in the code (%s)",
							name, r.line, r.token, cols[k].head, want, px, got))
					}
				}
			}
		}
	}
	return problems
}

// ── source field in the page header ───────────────────────────────────────
//
// Each page names the file its component is made from, and the site prints it
// as a link. Nobody checked it: splitting components.css into eight files left
// twenty-five pages pointing into empty space, and none of the seven gates said
// a word.
//
// This error costs exactly what the reference exists for: a person goes to see
// how the component is built and ends up nowhere.
var sourceRe = regexp.MustCompile(`(?m)^source:\s*(\S+)\s*$`)

func checkSourceFields(docsDir string) []string {
	root := filepath.Dir(docsDir)
	var problems []string
	filepath.WalkDir(docsDir, func(p string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(p, ".md") {
			return nil
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return nil
		}
		text := strings.ReplaceAll(string(b), "\r\n", "\n")
		i := strings.Index(text, "\n---")
		if !strings.HasPrefix(text, "---\n") || i < 0 {
			return nil
		}
		m := sourceRe.FindStringSubmatch(text[:i])
		if m == nil {
			return nil
		}
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(m[1]))); err != nil {
			rel, _ := filepath.Rel(docsDir, p)
			problems = append(problems, fmt.Sprintf(
				"%s: source points to %s — that file does not exist",
				filepath.ToSlash(rel), m[1]))
		}
		return nil
	})
	sort.Strings(problems)
	return problems
}

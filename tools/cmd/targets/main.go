// Command targets checks tap targets against WCAG 2.2 AA (2.5.8).
//
// The minimum is 24×24 CSS pixels for EVERY interactive target. The criterion
// has four exceptions and exactly one of them applies to a kit: an "inline"
// target inside a line of text. Everything else — a link in a sentence does not
// count, a styled range slider does — has to hold the minimum.
//
// It reads the REAL src/tokens.css and resolves var() and calc() the same way
// contrast does, so it cannot drift away from the kit. It checks ALL THREE
// densities: it is in compact that the whole size scale drops below the
// minimum, and that cannot be seen at the default.
//
//	go run ./cmd/targets
//	go run ./cmd/targets -v   with the list of passes
package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"

	"instrument/tools/internal/css"
)

// min is the WCAG 2.2 AA 2.5.8 minimum (Target Size, Minimum).
const min = 24.0

// target is an interactive target of the kit.
//
// w and h are token expressions, exactly the ones standing in the CSS. Writing
// numbers here would set up a second source of truth: the check would drift
// away from the kit silently, and that is precisely what contrast was written
// against.
type target struct {
	label string
	where string // where in the kit
	w, h  string

	// gap is the distance to the neighbouring target.
	//
	// The criterion has a DISTANCE exception, and without it the check lies. A
	// target below the minimum counts if a circle of diameter 24 placed at its
	// centre does not cross a neighbour's circle: for targets of size S with a
	// gap G that holds when S + G >= 24. Otherwise the compact density would be
	// a violation by construction, even though a 22px menu item with a 2px gap
	// satisfies the criterion.
	gap string

	// hit is the token a component grows its HIT AREA to with an invisible
	// pseudo-element. When it is set, the area is checked rather than the
	// shape: a target may stay small to the eye and still hold the minimum.
	hit string

	// alone marks a target that has no neighbours at all, so the distance
	// exception does not apply to it: it has to hold the minimum on its own.
	alone bool
}

var targets = []target{
	// Buttons. The icon one is square, so its width equals its height. In a
	// row, buttons are separated by the row gap.
	{label: "button sm", where: "actions.css .inst-btn--sm", w: "--control-h-sm", h: "--control-h-sm", gap: "--gap-inline"},
	{label: "button md", where: "actions.css .inst-btn", w: "--control-h-md", h: "--control-h-md", gap: "--gap-inline"},
	{label: "button lg", where: "actions.css .inst-btn--lg", w: "--control-h-lg", h: "--control-h-lg", gap: "--gap-inline"},
	{label: "icon button sm", where: "actions.css .inst-btn--sm.inst-btn--icon", w: "--control-h-sm", h: "--control-h-sm", gap: "--gap-inline"},

	// Checkboxes and switches. They have no neighbour: in a table's selection
	// column a checkbox stands alone, without a label, so it holds the minimum
	// itself.
	{label: "checkbox", where: "forms.css .inst-checkbox > input", w: "--size-check", h: "--size-check", hit: "--tap-min", alone: true},
	{label: "radio", where: "forms.css .inst-radio > input", w: "--size-check", h: "--size-check", hit: "--tap-min", alone: true},
	{label: "switch", where: "forms.css .inst-switch > input", w: "--size-switch-w", h: "--size-switch-h", hit: "--tap-min", alone: true},

	// Slider. The target is the whole control: you drag the thumb but you aim
	// at the track.
	{label: "slider", where: "forms.css .inst-slider", w: "--control-h-md", h: "--tap-min", alone: true},

	// Fields: the target is the whole control, and there are no horizontal
	// neighbours.
	{label: "field sm", where: "forms.css .inst-input--sm", w: "--control-h-sm", h: "--control-h-sm", gap: "--gap-row"},
	{label: "field md", where: "forms.css .inst-input", w: "--control-h-md", h: "--control-h-md", gap: "--gap-row"},

	// Navigation and lists: neighbours are separated by the component's own
	// gap.
	{label: "nav item", where: "layout.css .inst-nav-item", w: "--control-h-md", h: "--control-h-md", gap: "--space-1"},
	{label: "pager item", where: "layout.css .inst-pager-item", w: "--control-h-sm", h: "--control-h-sm", gap: "--space-2"},
	{label: "menu item", where: "overlay.css .inst-menu-item", w: "--control-h-sm", h: "--control-h-sm", gap: "--space-1"},
	{label: "tree node", where: "text.css .inst-tree-item", w: "--control-h-sm", h: "--control-h-sm", gap: "--space-1"},
	{label: "calendar day", where: "data.css .inst-calendar-day", w: "--control-h-sm", h: "--control-h-sm", gap: "--space-1"},
	{label: "segmented item", where: "actions.css .inst-segmented > button", w: "--control-h-md", h: "calc(var(--control-h-md) - var(--space-2))", gap: "--space-1"},

	// Insert. Its width comes from its text, so the height is what gets counted
	// for the width as well: even the shortest label has side padding, and it
	// cannot end up narrower than it is tall.
	{label: "insert", where: "rows.css .inst-insert", w: "--control-h-sm", h: "--control-h-sm", gap: "--space-2"},

	// Filter chip. Neighbours are separated by the strip's own gap, as with the
	// pager: the width comes from the text, so the height is counted for it.
	{label: "filter chip", where: "actions.css .inst-chip", w: "--control-h-sm", h: "--control-h-sm", gap: "--space-2"},

	// The tag's remove cross. The smallest target of the kit: there is no label
	// beside it and the neighbouring tag stands flush against it, so the
	// distance exception does not apply.
	{label: "tag remove", where: "data.css .inst-tag-remove", w: "--size-chevron", h: "--size-chevron", hit: "--tap-min", alone: true},
}

var densities = []struct{ id, label string }{
	{"", "default"},
	{"compact", "compact"},
	{"comfortable", "comfortable"},
}

// Scale is the second size axis. For this check it matters more than density:
// density presses targets down towards the floor of the criterion, while scale
// pulls them up, and the dangerous combination is a large scale with a compact
// layout, where the cell returns the heights almost to the default ones while
// leaving the type large.
//
// There is no step downwards: the compact mode already sits on the floor, and a
// scale below one would break 2.5.8 by construction.
var scales = []struct{ id, label string }{
	{"", "14px"},
	{"15", "15px"},
	{"16", "16px"},
	{"17", "17px"},
	{"18", "18px"},
}

// combo is one cell of the "scale × density" grid. The axes are flattened into
// a single list so that the body of the check stays one level deep.
type combo struct {
	scale, dens string
	label       string
}

func combos() []combo {
	var out []combo
	for _, sc := range scales {
		for _, d := range densities {
			out = append(out, combo{scale: sc.id, dens: d.id, label: sc.label + " · " + d.label})
		}
	}
	return out
}

func main() {
	tokens := flag.String("tokens", "../src/tokens.css", "path to tokens.css")
	verbose := flag.Bool("v", false, "show the passes")
	flag.Parse()

	src, err := css.Load(*tokens)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot read the tokens:", err)
		os.Exit(1)
	}

	// The role tier is declared through :where(:root), the densities through an
	// attribute. The --space-* scale is declared in :root and the roles in
	// :where(:root). Both are needed: the gaps come from the first, the sizes
	// from the second.
	base := src.Decls(regexp.MustCompile(`(?m)^:root \{`))
	for k, v := range src.Decls(regexp.MustCompile(`:where\(:root\)\s*\{`)) {
		base[k] = v
	}
	if len(base) == 0 {
		fmt.Fprintln(os.Stderr, "the role tier was not found — :where(:root)")
		os.Exit(1)
	}

	width := 0
	for _, t := range targets {
		if n := utf8.RuneCountInString(t.label); n > width {
			width = n
		}
	}

	failed, total := 0, 0
	for _, d := range combos() {
		vals := map[string]string{}
		for k, v := range base {
			vals[k] = v
		}
		// The one-dimensional part of the scale goes down first, the
		// two-dimensional cell on top of it.
		if d.scale != "" {
			for k, v := range src.Decls(regexp.MustCompile(`\[data-scale="` + d.scale + `"\]\s*\{`)) {
				vals[k] = v
			}
		}
		if d.dens != "" {
			re := regexp.MustCompile(`\[data-density="` + d.dens + `"\]\s*\{`)
			if d.scale != "" {
				re = regexp.MustCompile(`\[data-scale="` + d.scale + `"\]\[data-density="` + d.dens + `"\]`)
			}
			for k, v := range src.Decls(re) {
				vals[k] = v
			}
		}

		fmt.Printf("\nSCALE · DENSITY %s\n", d.label)
		fmt.Println(strings.Repeat("─", width+34))

		for _, t := range targets {
			total++
			w, errW := resolve(vals, t.w)
			h, errH := resolve(vals, t.h)
			pad := strings.Repeat(" ", width-utf8.RuneCountInString(t.label))
			if errW != nil || errH != nil {
				failed++
				fmt.Printf("  ✗ %s%s  ERROR: %v %v\n", t.label, pad, errW, errH)
				continue
			}
			gap := 0.0
			if t.gap != "" && !t.alone {
				if g, err := resolve(vals, t.gap); err == nil {
					gap = g
				}
			}
			if t.hit != "" {
				if v, err := resolve(vals, t.hit); err == nil {
					if v > w {
						w = v
					}
					if v > h {
						h = v
					}
				}
			}
			ok := w+gap >= min && h+gap >= min
			mark := "·"
			if !ok {
				mark, failed = "✗", failed+1
			}
			note := ""
			if gap > 0 {
				note = fmt.Sprintf("  +%.0f gap", gap)
			}
			fmt.Printf("  %s %s%s  %4.0f×%-4.0f  (need %.0f×%.0f)%s\n",
				mark, t.label, pad, w, h, min, min, note)
		}
	}

	fmt.Println()
	if failed > 0 {
		fmt.Printf("✗ failures: %d of %d\n", failed, total)
		fmt.Println()
		fmt.Println("A target is grown by its hit area, not by its shape:")
		fmt.Println("  .inst-tag-remove { position: relative }")
		fmt.Println("  .inst-tag-remove::before { content: \"\"; position: absolute; inset: -7px }")
		os.Exit(1)
	}
	if *verbose {
		fmt.Printf("· all %d targets hold %.0f×%.0f across three densities\n", total, min, min)
	} else {
		fmt.Printf("· all %d checks passed in %d scale and density combinations\n", total, len(combos()))
	}
}

var (
	varRe  = regexp.MustCompile(`var\(\s*(--[a-z][\w-]*)\s*\)`)
	calcRe = regexp.MustCompile(`calc\(([^()]*)\)`)
	pxRe   = regexp.MustCompile(`([\d.]+)px`)
)

// resolve expands var() and evaluates calc(). The values here are geometric,
// that is always in px: colours and percentages never arrive, and that is why
// the parsing fits into thirty lines rather than the colour resolver contrast
// needs.
func resolve(vals map[string]string, expr string) (float64, error) {
	if strings.HasPrefix(expr, "--") {
		v, ok := vals[expr]
		if !ok {
			return 0, fmt.Errorf("no token %s", expr)
		}
		expr = v
	}
	for i := 0; i < 12 && strings.Contains(expr, "var("); i++ {
		var missing string
		expr = varRe.ReplaceAllStringFunc(expr, func(m string) string {
			name := varRe.FindStringSubmatch(m)[1]
			v, ok := vals[name]
			if !ok {
				missing = name
				return m
			}
			return strings.TrimSpace(v)
		})
		if missing != "" {
			return 0, fmt.Errorf("no token %s", missing)
		}
	}
	for strings.Contains(expr, "calc(") {
		before := expr
		expr = calcRe.ReplaceAllStringFunc(expr, func(m string) string {
			inner := calcRe.FindStringSubmatch(m)[1]
			v, ok := arith(inner)
			if !ok {
				return m
			}
			return fmt.Sprintf("%gpx", v)
		})
		if expr == before {
			return 0, fmt.Errorf("cannot parse the calc in %q", expr)
		}
	}
	if v, ok := arith(expr); ok {
		return v, nil
	}
	return 0, fmt.Errorf("cannot parse %q", expr)
}

// arith evaluates an expression of numbers in px. Multiplication and division
// come before addition and subtraction — as in arithmetic and as in CSS.
func arith(expr string) (float64, bool) {
	expr = pxRe.ReplaceAllString(strings.TrimSpace(expr), "$1")
	f := strings.Fields(expr)
	if len(f) == 0 || len(f)%2 == 0 {
		return 0, false
	}
	num := func(s string) (float64, bool) {
		var v float64
		_, err := fmt.Sscanf(s, "%g", &v)
		return v, err == nil
	}
	first, ok := num(f[0])
	if !ok {
		return 0, false
	}
	vals := []float64{first}
	var ops []string
	for i := 1; i < len(f); i += 2 {
		rhs, ok := num(f[i+1])
		if !ok {
			return 0, false
		}
		switch f[i] {
		case "*":
			vals[len(vals)-1] *= rhs
		case "/":
			if rhs == 0 {
				return 0, false
			}
			vals[len(vals)-1] /= rhs
		case "+", "-":
			ops = append(ops, f[i])
			vals = append(vals, rhs)
		default:
			return 0, false
		}
	}
	out := vals[0]
	for i, op := range ops {
		if op == "+" {
			out += vals[i+1]
		} else {
			out -= vals[i+1]
		}
	}
	return out, true
}

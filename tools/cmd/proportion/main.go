// proportion — checks ROLE PROPORTIONS against the type scale.
//
// Why a separate command. contrast guards color, targets — tap targets,
// docscheck — reference consistency, dist — delivery. None of them saw the
// relationship between sizes, and this surfaced expensively: when the base
// type size grew from 13 to 14, the glyph and gutter tier stayed in place.
// All four checks remained green, while on screen the icon fell 9.8% behind
// the uppercase letter, and the label column overflowed and broke form
// alignment.
//
// These checks measure RELATIONSHIPS, not values: a specific number can be
// changed deliberately, but a ratio that leaves its band is almost always
// an oversight.
//
// The bands come from the kit's current state, not from guesswork: each is
// labeled with what happens when it falls outside the band.
package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"

	"instrument/tools/internal/css"
)

// Type size distinguishability threshold. Below this, adjacent steps read as
// one, and the scale loses a step while continuing to declare it. 1.12 is the
// lower bound of the band that contains the kit's working steps.
const stepMin = 1.12

type rule struct {
	label    string
	a, b     string  // what is divided by what
	min, max float64 // allowed ratio band
	why      string  // what happens when it leaves the band
	perDens  bool    // check at each density
}

var rules = []rule{
	// ── Type size ladder ────────────────────────────────────────────────────
	// Each step must differ from its neighbor enough for the difference to be
	// visible. The upper end of the scale spreads wider — that is where the
	// difference between headings lives.
	{label: "type: 2xs → xs", a: "--text-xs", b: "--text-2xs", min: stepMin, max: 1.30,
		why: "steps merge: the scale declares a size that cannot be seen"},
	{label: "type: xs → sm", a: "--text-sm", b: "--text-xs", min: stepMin, max: 1.30,
		why: "metadata stops differing from the base by size"},
	{label: "type: sm → md", a: "--text-md", b: "--text-sm", min: stepMin, max: 1.30,
		why: "panel name merges with the data below it"},
	{label: "type: md → lg", a: "--text-lg", b: "--text-md", min: stepMin, max: 1.30,
		why: "block heading merges with prose"},
	{label: "type: lg → xl", a: "--text-xl", b: "--text-lg", min: stepMin, max: 1.40,
		why: "section heading merges with block heading"},
	{label: "type: xl → 2xl", a: "--text-2xl", b: "--text-xl", min: stepMin, max: 1.40,
		why: "hero number merges with heading"},

	// ── Glyph against text ─────────────────────────────────────────────────
	// The icon and spinner sit next to a label and are measured NOT by the box,
	// but by how their ink looks against the uppercase letter. Uppercase is
	// 0.700 of the type size — measured with the CSS `cap` unit, which asks the
	// font rather than the raster; the 0.71 that stood here was a guess that
	// happened to be close. Icon ink is 9.02 from a box of 16, or 0.564.
	// Hence the band: icon box to base type size.
	{label: "icon to base type size", a: "--size-icon", b: "--text-sm", min: 1.20, max: 1.34, perDens: true,
		why: "icon next to the label reads tiny or overwhelms it instead"},
	// The same measure, but against ITS OWN type size. The status row is set in
	// --text-2xs, so its icon must be measured against the uppercase of eleven,
	// not fourteen: ordinary 18 against 11 gives 1.64 and outweighs the readout.
	//
	// This band used to be 1.20…1.42 — wider than its neighbour — and the
	// widening had a reason written out at length: the model "uppercase = 0.71
	// of the type size" could not resolve one device pixel, so the fine part of
	// the invariant was delegated to the pixel gate where uppercase was REAL.
	// Every word of that was wrong. The pixel gate measured uppercase with a
	// canvas, which reports the RASTERISED ascent rounded to a whole pixel and
	// was off by 0.8px; the real cap height of this family is a clean 0.700 of
	// the type size and resolves nothing in pixels at all. Three rungs of the
	// --size-icon-sm ladder were raised to satisfy that canvas, and the band was
	// widened to 1.42 to admit them. The band was fitted to the correction, and
	// the correction was fitted to a broken ruler.
	//
	// The ladder went back to 14-14-16-16-18, and with it the band is DERIVED
	// rather than chosen. The five ratios are 1.2174…1.2857. Move any one rung
	// by a single pixel and the highest such value below the body is 1.2143
	// (17/14, scale 18) and the lowest above it is 1.3043 (15/11.5, scale 15).
	// The body therefore fits strictly between two defects, and a band exists
	// that admits every honest rung and refuses every one-pixel error:
	//
	//     1.2143  <  [ 1.2174 … 1.2857 ]  <  1.3043
	//              ^                    ^
	//              1.216                1.300
	//
	// The lower margin is 0.003 and that is the POINT rather than a fragility:
	// it is smaller than the smallest step this quantity can take, so nothing
	// broken can hide inside it. Its neighbour above cannot be tightened this
	// way — the --size-icon ladder spans 1.25…1.3333 while a one-pixel error
	// reaches 1.2778 and 1.3125, so its body and its defects OVERLAP and no
	// global band separates them. The two bands are different numbers because
	// the two ladders are different objects, and that is now a derivation
	// rather than a preference.
	{label: "small icon to its type size", a: "--size-icon-sm", b: "--text-2xs", min: 1.216, max: 1.300, perDens: true,
		why: "icon in the status row overwhelms the readout or disappears beside it"},
	{label: "spinner to base type size", a: "--size-spinner", b: "--text-sm", min: 0.92, max: 1.20, perDens: true,
		why: "busy indicator takes the label's place and must be its height"},
	{label: "chevron to base type size", a: "--size-chevron", b: "--text-sm", min: 0.62, max: 0.86, perDens: true,
		why: "disclosure arrow competes with the label or gets lost beside it"},

	// An icon WITHOUT a label is measured not against uppercase, but against
	// the tap target: there are no words beside it, so there is nothing to
	// compare against. Its own rule, rather than the icon-with-label band —
	// otherwise --size-icon-lg (24 against base 14, or 1.71) would fail the
	// gate, doing exactly what the rule was introduced to prevent.
	//
	// Lower bound: a glyph smaller than half the target gets lost in an empty
	// button. Upper bound: a glyph larger than three quarters hits the edges,
	// and the target stops reading as a button.
	{label: "rail icon to its button", a: "--size-icon-lg", b: "--control-h-lg", min: 0.5, max: 0.76, perDens: true,
		why: "glyph gets lost in an empty button or hits its edges"},

	// ── Rhythm ──────────────────────────────────────────────────────────────
	// A DIFFERENT KIND OF RULE from the others here. The other bands are
	// perceptual: "the ratio must land in a range, otherwise the difference is
	// not visible." Here the ratio is EXACT, and the band is narrowed to the
	// division error.
	//
	// --gap-section is derived as --gap-row raised by four scale steps. The
	// scale 2·4·6·8·12·16·24·32·48·64 is arranged so that +4 steps gives
	// exactly ×4 at every point: 4→16, 6→24, 8→32, 12→48, 16→64. That is what
	// makes the gap between sections a multiple of the spacing within them in
	// all fifteen cells without separate tuning.
	//
	// The gate guards the derivation, not the number: editing one cell catches
	// it immediately. Measured on assembled screens, which is where the four
	// itself came from: Anthropic uses 8·12·16 inside a group, 54–61 between
	// groups.
	{label: "section to row", a: "--gap-section", b: "--gap-row", min: 3.99, max: 4.01, perDens: true,
		why: "gap between sections is no longer a multiple of the spacing within them — the screen reads as one continuous flow again"},

	// ── Control against text ───────────────────────────────────────────────
	// The band is intentionally wide: control height is retuned by density,
	// while type size is not, so the ratio must vary. Bounds are taken from the
	// edges of the ladder itself (26/14 at dense and 36/14 at loose) with one
	// step of margin. The rule catches not drift within the ladder, but falling
	// outside it: a control whose declared height misses the row.
	{label: "control height to base", a: "--control-h-md", b: "--text-sm", min: 1.8, max: 2.8, perDens: true,
		why: "label hits the button ceiling or sinks into it"},
	{label: "badge to its type size", a: "--control-h-xs", b: "--text-2xs", min: 1.5, max: 2.2, perDens: true,
		why: "badge squeezes the label or inflates around it"},

	// ── Control shape ──────────────────────────────────────────────────────
	// Padding is derived from height through --control-ratio-*, and these three
	// rules guard the derivation itself: if someone declares padding as a value
	// again, the shape will move with density, as it did before the derivation.
	{label: "sm button shape", a: "--control-pad-sm", b: "--control-h-sm", min: 0.30, max: 0.32, perDens: true,
		why: "padding is no longer derived from height — button SHAPE changes with density"},
	{label: "md button shape", a: "--control-pad-md", b: "--control-h-md", min: 0.365, max: 0.385, perDens: true,
		why: "padding is no longer derived from height — button SHAPE changes with density"},
	{label: "lg button shape", a: "--control-pad-lg", b: "--control-h-lg", min: 0.41, max: 0.43, perDens: true,
		why: "padding is no longer derived from height — button SHAPE changes with density"},

	// Radius is the second half of the same shape. One radius for all three
	// sizes reads as almost a pill on a small control: --radius-md, or 8 at
	// height 26, is 0.31, while in dense mode 8 at 22 is already 0.36. Lower
	// bound: below 0.16 the control stops being rounded and competes with the
	// card around it. Upper bound: above 0.30 it reads as a pill, and a pill in
	// the kit means "toggle".
	{label: "sm control rounding", a: "--radius-control-sm", b: "--control-h-sm", min: 0.16, max: 0.30, perDens: true,
		why: "small control reads as a pill or loses its rounding entirely"},
	{label: "md control rounding", a: "--radius-control-md", b: "--control-h-md", min: 0.16, max: 0.30, perDens: true,
		why: "control reads as a pill or loses its rounding entirely"},
	{label: "lg control rounding", a: "--radius-control-lg", b: "--control-h-lg", min: 0.16, max: 0.30, perDens: true,
		why: "large control reads as a pill or loses its rounding entirely"},

	// ── Mark shape ─────────────────────────────────────────────────────────
	// The same guard as above, pointing the other way. A control is big enough
	// that its declared radius is the one that gets painted; a mark is not.
	// Below 8px the browser scales every radius on the box until two of them
	// fit the edge, so a radius at or above half the box does not read as "very
	// rounded" — it reads as the shape at half, and the shapes at half are
	// SPOKEN FOR. Half a dot is a circle, which marks a reversible consequence.
	// Half a bar is a pill, which is a badge and a tag.
	//
	// That is what the upper bound buys: not taste, but keeping --radius-mark
	// out of the two shapes the kit assigns meaning to. The lower one keeps it
	// from vanishing — at 1px on a 6px dot the corner is a rasterising artefact
	// rather than a rounding.
	{label: "mark rounding to dot", a: "--radius-mark", b: "--size-dot", min: 0.20, max: 0.40,
		why: "the state dot clamps into a circle, and the circle marks a reversible consequence"},
	{label: "mark rounding to meter", a: "--radius-mark", b: "--size-meter", min: 0.20, max: 0.40,
		why: "the meter and the lane track clamp into a pill, and a pill in the kit is a badge and a tag"},

	// The tick is the one mark the tier cannot save from clamping: it is 3px
	// wide and any corner worth drawing is more than half of that. The band is
	// therefore looser and guards a different failure — a corner that outgrows
	// its own tick, at which point the painted shape is set by the tick's WIDTH,
	// and on a grouped strip the width is set by the data.
	{label: "mark rounding to tick", a: "--radius-mark", b: "--size-tick", min: 0.30, max: 0.70,
		why: "the corner outgrows the tick, and the strip's shape starts following the traffic through it"},

	// ── Label column ────────────────────────────────────────────────────────
	// It holds TEXT, so it must grow with type size. This is exactly what
	// overflowed when the base grew: "Token limit" required 99.9px at a width
	// of 92. Same caveat: the column grows with density, type size does not.
	// The band is wide and catches only gross deviation.
	//
	// COLUMN OVERFLOW IS NOT CAUGHT BY THIS RULE, and that is the honest limit
	// of a token-level gate: whether "Token limit" fits into its pixels depends
	// on the font, not the tokens. tools/audit.js measures that on the rendered
	// output.
	{label: "label column to base", a: "--label-col", b: "--text-sm", min: 6.0, max: 9.0, perDens: true,
		why: "label column fell outside its ladder"},

	// ── Vertical rhythm ────────────────────────────────────────────────────
	{label: "block vertical to row", a: "--pad-block-y", b: "--row-pad-y", min: 0.9, max: 1.7, perDens: true,
		why: "panel header becomes taller than the row it labels"},

	// "Tighter vertically than horizontally: some vertical air comes from
	// line-height" — a layout rule that had lived only in text until now. It
	// held in fourteen of fifteen cells: in loose density at the base type size
	// it was 12 against 12, exactly what the rule forbids. The upper bound is
	// strictly less than one — that is the entire rule; the lower bound rejects
	// the opposite skew, where the row gets flattened against its own sides.
	{label: "row vertical to horizontal", a: "--row-pad-y", b: "--pad-cell-x", min: 0.55, max: 0.95, perDens: true,
		why: "vertical is no longer tighter than horizontal — the row expands upward even though line-height already supplied the air"},
}

// Not everything is checked by ratio. Type floor is absolute, radii are a
// property of the number itself. Such rules live here.
type absCheck struct {
	label string
	fn    func(v func(string) float64) (bool, string)
}

var absolutes = []absCheck{
	{"type floor not breached", func(v func(string) float64) (bool, string) {
		// 11px — the declared kit floor: there is no size below it and none
		// should be introduced.
		//
		// The check exists because the floor is declared in rem, so the floor
		// is NOT SELF-CONTAINED, but depends on a root of 16. A consumer can
		// change the root — this is the only text-density control — and at root
		// 14 the smallest step silently drops to 9.63px.
		//
		// Clamping it through max(11px, …) is impossible, and this needs to be
		// recorded so it is not proposed again: at root 14 the clamped smallest
		// step gives 11, while the unclamped neighbor gives 10.94, inverting
		// the step order. Clamping both means collapsing them into one number
		// and losing a step. Therefore the floor is guarded by a check, not an
		// expression.
		got := v("--text-2xs")
		if got < 11-0.01 {
			return false, fmt.Sprintf("--text-2xs = %.2fpx at root %g — below the declared floor of 11px", got, css.RootPx)
		}
		return true, fmt.Sprintf("--text-2xs = %.2fpx", got)
	}},
	{"radii are even", func(v func(string) float64) (bool, string) {
		for _, n := range []string{"--radius-2xs", "--radius-xs", "--radius-sm", "--radius-md", "--radius-lg"} {
			r := v(n)
			if int(r)%2 != 0 || r != float64(int(r)) {
				return false, fmt.Sprintf("%s = %g: odd radius produces half a device pixel at density 1.5, and the arc meets the edge off-grid", n, r)
			}
		}
		return true, fmt.Sprintf("%g · %g · %g · %g · %g",
			v("--radius-2xs"), v("--radius-xs"), v("--radius-sm"), v("--radius-md"), v("--radius-lg"))
	}},
	{"radius ladder grows", func(v func(string) float64) (bool, string) {
		xs, sm, md, lg := v("--radius-xs"), v("--radius-sm"), v("--radius-md"), v("--radius-lg")
		if !(xs < sm && sm < md && md < lg) {
			return false, fmt.Sprintf("%g · %g · %g · %g — nested radius must be smaller than the outer one", xs, sm, md, lg)
		}
		return true, fmt.Sprintf("%g < %g < %g < %g", xs, sm, md, lg)
	}},
	// A segment inside the track subtracts the field from the track radius
	// directly in CSS, so concentricity maintains itself. What remains to guard
	// is the subtraction result: it must be positive and EVEN, otherwise the
	// inner corner rasterizes off the device grid at density 1.5.
	{"segment radius subtraction leaves an even remainder", func(v func(string) float64) (bool, string) {
		md, gap := v("--radius-control-md"), v("--space-1")
		in := md - gap
		if in <= 0 || int(in)%2 != 0 {
			return false, fmt.Sprintf("%g − %g = %g: inner radius must be positive and even", md, gap, in)
		}
		return true, fmt.Sprintf("%g − %g = %g", md, gap, in)
	}},
}

var densities = []struct{ id, label string }{
	{"", "normal"},
	{"compact", "compact"},
	{"comfortable", "comfortable"},
}

// Scale is the SECOND dimensional axis, and it cannot be checked separately
// from density: both move control heights, and their combination is described
// by its own token cell. Hence 3 × 3 checks, not 3 + 3.
//
// There is intentionally no step below: compact mode is already at the WCAG
// 2.5.8 floor, and a scale below one would breach the criterion by construction.
var scales = []struct{ id, label string }{
	{"", "14px"},
	{"15", "15px"},
	{"16", "16px"},
	{"17", "17px"},
	{"18", "18px"},
}

// combo is one cell of the "scale × density" grid. The axes are flattened
// into one list deliberately: this keeps the check body flat, and adding a
// third axis will not turn it into a ladder of nested loops.
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
	verbose := flag.Bool("v", false, "show passed checks")
	root := flag.Float64("root", css.RootPx, "document root in px: type sizes are declared in rem and calculated from it")
	flag.Parse()
	css.RootPx = *root

	src, err := css.Load(*tokens)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot read tokens:", err)
		os.Exit(1)
	}

	base := src.Decls(regexp.MustCompile(`(?m)^:root \{`))
	for k, v := range src.Decls(regexp.MustCompile(`:where\(:root\)\s*\{`)) {
		base[k] = v
	}
	// Control padding is derived below the density blocks by a common rule.
	for k, v := range src.Decls(regexp.MustCompile(`:where\(:root\), \[data-density\], \[data-scale\]\s*\{`)) {
		base[k] = v
	}
	if len(base) == 0 {
		fmt.Fprintln(os.Stderr, "role tier not found — :where(:root)")
		os.Exit(1)
	}

	width := 0
	for _, r := range rules {
		if n := utf8.RuneCountInString(r.label); n > width {
			width = n
		}
	}
	for _, r := range absolutes {
		if n := utf8.RuneCountInString(r.label); n > width {
			width = n
		}
	}

	var problems []string
	total := 0

	valOf := func(vals map[string]string) func(string) float64 {
		return func(n string) float64 {
			v, err := css.ResolvePx(vals, n)
			if err != nil {
				problems = append(problems, fmt.Sprintf("  %s: %v", n, err))
				return 0
			}
			return v
		}
	}

	// scaleVals — base plus the one-dimensional scale part: type sizes, radius
	// ladder, glyphs, and text widths.
	scaleVals := func(id string) map[string]string {
		vals := map[string]string{}
		for k, v := range base {
			vals[k] = v
		}
		if id != "" {
			for k, v := range src.Decls(regexp.MustCompile(`\[data-scale="` + id + `"\]\s*\{`)) {
				vals[k] = v
			}
		}
		return vals
	}

	// Radii do not depend on density, but they do depend on SCALE: the ladder
	// is redeclared as a whole. Therefore absolute checks run once per scale,
	// not once for the entire kit.
	for _, sc := range scales {
		vals := scaleVals(sc.id)
		for _, rc := range absolutes {
			total++
			ok, note := rc.fn(valOf(vals))
			if !ok {
				problems = append(problems, fmt.Sprintf("  · %-*s  %s  (scale %s)", width, rc.label, note, sc.label))
			} else if *verbose {
				fmt.Printf("  · %-*s  %s  (scale %s)\n", width, rc.label, note, sc.label)
			}
		}
	}

	var grid []cellVals

	for _, d := range combos() {
		vals := scaleVals(d.scale)
		// Density comes from its own block when the scale is normal, and from a
		// two-dimensional cell otherwise: the cell fully describes the pair.
		if d.dens != "" {
			re := regexp.MustCompile(`\[data-density="` + d.dens + `"\]\s*\{`)
			if d.scale != "" {
				re = regexp.MustCompile(`\[data-scale="` + d.scale + `"\]\[data-density="` + d.dens + `"\]`)
			}
			for k, v := range src.Decls(re) {
				vals[k] = v
			}
		}
		grid = append(grid, cellVals{scale: d.scale, dens: d.dens, label: d.label, vals: vals})
		if *verbose {
			fmt.Printf("\n── %s ──\n", d.label)
		}
		for _, r := range rules {
			// perDens means "depends on DENSITY", not "on any cell".
			// The type-size ladder is not touched by density, so in its cells
			// the rule is skipped — but SCALE redeclares it as a whole, and it
			// cannot be skipped there. Therefore the condition looks at d.dens,
			// not at d.perDens: with the latter, scale ladders disappear from
			// the check entirely, and the gate turns green undeservedly.
			if d.dens != "" && !r.perDens {
				continue
			}
			total++
			a, err1 := css.ResolvePx(vals, r.a)
			b, err2 := css.ResolvePx(vals, r.b)
			if err1 != nil || err2 != nil {
				problems = append(problems, fmt.Sprintf("  · %-*s  cannot resolve: %v %v", width, r.label, err1, err2))
				continue
			}
			if b == 0 {
				problems = append(problems, fmt.Sprintf("  · %-*s  divisor is zero", width, r.label))
				continue
			}
			got := a / b
			if got < r.min || got > r.max {
				problems = append(problems, fmt.Sprintf(
					"  · %-*s  %.3f  (band %.2f–%.2f, %s)\n      %s",
					width, r.label, got, r.min, r.max, d.label, r.why))
			} else if *verbose {
				fmt.Printf("  · %-*s  %.3f  (band %.2f–%.2f)\n", width, r.label, got, r.min, r.max)
			}
		}
	}

	mono, monoChecked := checkMonotonic(grid, base)
	problems = append(problems, mono...)
	total += monoChecked

	fmt.Println()
	if len(problems) > 0 {
		sort.Strings(problems)
		fmt.Printf("── proportions left the band (%d) ──\n", len(problems))
		for _, p := range problems {
			fmt.Println(p)
		}
		fmt.Println()
		fmt.Printf("· %d checks at root %gpx, %d failed\n", total, css.RootPx, len(problems))
		os.Exit(1)
	}
	fmt.Printf("· all %d checks passed in %d scale and density combinations at root %gpx\n",
		total, len(combos()), css.RootPx)

}

// ── Cross-cell check ────────────────────────────────────────────────────────
//
// The rules above measure ratios WITHIN one cell and are therefore blind to
// the shape of the table itself. That is exactly how the --row-pad-y inversion
// passed through them: at the base type size, loose density gave 12px, while
// at scale 15 it gave 8px, meaning the geometry DECREASED as scale increased.
// Both cells individually sat inside the "block vertical to row" band (1.00
// and 1.50 with a 0.9–1.7 tolerance), and the gate stayed silent.
//
// Here the opposite is checked: not the ratio at a point, but the behavior of
// the value ALONG THE AXIS. Scale only goes upward, so the size it moves must
// also only go upward.
//
// The token list is not introduced: it is derived from the tokens themselves.
// Everything that resolves to pixels across all five scales is checked —
// colors, fractions, and keywords filter themselves out by failing to resolve.
// A manual list here would be the eighth registry that has to be remembered.
//
// Plateaus are intentionally allowed: the ladder is set in whole pixels, and
// adjacent steps can legitimately coincide after rounding. Only DECREASE is
// forbidden.

type cellVals struct {
	scale, dens string
	label       string
	vals        map[string]string
}

// tokenNames — names in stable order so the report does not jump between
// runs.
func tokenNames(base map[string]string) []string {
	out := make([]string, 0, len(base))
	for k := range base {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func checkMonotonic(grid []cellVals, base map[string]string) (problems []string, checked int) {
	byKey := map[string]map[string]string{}
	for _, c := range grid {
		byKey[c.scale+"/"+c.dens] = c.vals
	}

	for _, name := range tokenNames(base) {
		for _, d := range densities {
			// The token's track along scale at a fixed density.
			type step struct {
				label string
				px    float64
			}
			var track []step
			ok := true
			for _, sc := range scales {
				vals := byKey[sc.id+"/"+d.id]
				if vals == nil {
					ok = false
					break
				}
				px, err := css.ResolvePx(vals, name)
				if err != nil {
					// Not a length — not our concern. Color, fraction, keyword.
					ok = false
					break
				}
				track = append(track, step{sc.label, px})
			}
			if !ok || len(track) < 2 {
				continue
			}
			checked++
			for i := 1; i < len(track); i++ {
				// Tolerance for fractional rem: 0.01px, not exact comparison.
				if track[i].px < track[i-1].px-0.01 {
					var line []string
					for _, s := range track {
						line = append(line, fmt.Sprintf("%s→%g", s.label, s.px))
					}
					problems = append(problems, fmt.Sprintf(
						"  · %s  decreases as scale increases, density \"%s\"\n      %s\n"+
							"      scale only goes up: the value it moves has no right to decrease",
						name, d.label, strings.Join(line, "  ")))
					break
				}
			}
		}
	}
	return problems, checked

}

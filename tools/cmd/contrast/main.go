// The contrast command checks the kit's token contrast against WCAG thresholds.
//
// It reads the ACTUAL src/tokens.css and resolves semantics the same way the
// browser does: var() recursively, light-dark() by scheme,
// color-mix(… transparent) into alpha, oklch() into sRGB. Therefore the check
// cannot diverge from the kit — it does not duplicate values, it computes them.
//
// Acceptance rule:
//
//	text  < 18px   — 4.5:1  (WCAG 1.4.3)
//	text ≥ 18px    — 3.0:1
//	non-decorative border — 3.0:1  (WCAG 1.4.11): checkbox, field, switch
//	                         track, state indicator. A decorative frame with
//	                         a nearby surface change is not included and is
//	                         intentionally quieter.
//
//	go run ./cmd/contrast
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"

	"instrument/tools/internal/css"
)

const (
	text  = 4.5 // < 18px
	large = 3.0 // ≥ 18px, as well as non-decorative borders and indicators

	// The DISTINGUISHABILITY threshold, not the accessibility threshold, and
	// the only one in this file.
	//
	// It is needed where two elements sit flush together and must read as
	// different, but accessibility says nothing about this pair. There is
	// exactly one case: the focus ring around a SOLID fill. By the letter of
	// WCAG, the color adjacent to the ring is the surface in the gap around
	// it, and that is checked by the three pairs below; the ring's fill is not
	// touching it. Visually, though, a ring matching the fill is not focus, but
	// "the button became thicker."
	//
	// 3.0 cannot go here: the farthest accent step gives 2.63 against the fill,
	// so the threshold would be unreachable for any accent ring and would
	// require a neutral one. 1.5 is a confident lightness step; a ring matching
	// the fill (that was: exactly 1.00) does not reach it.
	distinct = 1.5

	// The SURFACE STACK STEP threshold — and it has its own measure.
	//
	// The kit's first surface law ("depth is conveyed by lightness order") has
	// never been checked by anything: every case above is a FOREGROUND on a
	// BACKGROUND, while surface against surface was never asked. The cost:
	// panel and card both sat on --surface-raised, the difference was exactly
	// 1.00, and on the "Height and surfaces" page the card inside the panel was
	// not readable.
	//
	// Contrast ratio cannot measure this. The WCAG formula has a 0.05 term —
	// a highlight model — and at the lightness floor it compresses all ratios
	// toward one: the same 0.035 step gives 1.10 at the white end of the ramp
	// and 1.04 at the black end. A threshold based on the light end would
	// declare dark themes broken, while a threshold based on the dark end would
	// stop catching anything in light themes.
	//
	// Therefore the step is measured as OKLCH LIGHTNESS DIFFERENCE — the same
	// axis used to build the ramp, and perceptually uniform by construction.
	// One threshold works across all five themes.
	//
	// 0.022 is just below the kit's tightest step (light end, 0.024). The margin
	// is intentionally small: a threshold that passes by twice catches nothing.
	step = 0.022
)

// A check case: description, foreground token, background stack, threshold.
type kase struct {
	label string
	fg    string
	bg    []string
	min   float64

	// alt is the SECOND stack against which the first is compared, and this
	// comparison is DIRECTIONAL: fg must be QUIETER than alt on background bg,
	// not merely different from it.
	//
	// It is needed where two states are alternatives, not layers. The soft
	// button weight REPLACES the default fill (`--btn-bg` is reassigned), rather
	// than sitting on top of it, so the question is: does the soft button differ
	// from the default on the same background — and in the correct direction?
	//
	// Direction matters, and it cost a regression. The absolute difference
	// allows a REVERSED ladder: a soft weight that becomes louder than the
	// default differs from it by exactly as much as required, and the gate stays
	// silent. On screen, the third step is then louder than the second, and only
	// the eye can notice it.
	//
	// Loudness is measured by presence — the absolute step against the
	// background. The sign of the film is not part of it: in light themes the
	// control is recessed, in dark themes it is raised, and it must be quieter
	// in both.
	alt []string
}

// The MEASURE is derived from the threshold, rather than stored as a separate
// field.
//
// A contrast ratio by definition cannot be less than one: it is the quotient
// of greater luminance by lesser. Therefore a threshold below one cannot be a
// ratio — and exactly one case remains for which it exists: the surface stack
// step, measured as an OKLCH lightness difference.
//
// A fifth field in the struct would say the same thing, but would require
// adding it to each of the eighty-eight existing cases — and would create the
// possibility of the flag becoming inconsistent with the threshold.
func (c kase) isStep() bool { return c.min < 1 }

var cases = []kase{
	// Text on three surfaces. Each step is checked everywhere it can appear —
	// this exact rule was not being enforced before.
	{label: "text: primary on panel", fg: "--text-primary", bg: []string{"--surface-raised"}, min: text},
	{label: "text: primary on page", fg: "--text-primary", bg: []string{"--surface-page"}, min: text},
	{label: "text: primary in inset", fg: "--text-primary", bg: []string{"--surface-sunken"}, min: text},
	{label: "text: secondary on panel", fg: "--text-secondary", bg: []string{"--surface-raised"}, min: text},
	{label: "text: secondary in inset", fg: "--text-secondary", bg: []string{"--surface-sunken"}, min: text},
	{label: "text: muted on panel", fg: "--text-muted", bg: []string{"--surface-raised"}, min: text},
	{label: "text: muted on page", fg: "--text-muted", bg: []string{"--surface-page"}, min: text},
	{label: "text: muted in inset (log)", fg: "--text-muted", bg: []string{"--surface-sunken"}, min: text},
	// Recessed surfaces are translucent and therefore declared as a STACK:
	// the background comes first, the fill second. Flatten combines them.
	// This pair cannot be declared as one layer — the gate would read the alpha
	// as an opaque color and measure contrast against black.
	{label: "text: primary in field on panel", fg: "--text-primary", bg: []string{"--surface-raised", "--surface-field"}, min: text},
	{label: "text: primary in field on page", fg: "--text-primary", bg: []string{"--surface-page", "--surface-field"}, min: text},
	{label: "text: button label", fg: "--text-primary", bg: []string{"--surface-raised", "--surface-recessed"}, min: text},
	{label: "text: button label on hover", fg: "--text-primary", bg: []string{"--surface-raised", "--surface-recessed-hover"}, min: text},
	{label: "text: tag on panel", fg: "--text-secondary", bg: []string{"--surface-raised", "--surface-recessed"}, min: text},
	{label: "text: selected chip", fg: "--accent-text", bg: []string{"--surface-raised", "--surface-recessed"}, min: text},

	// faint — a decoration threshold, not a reading threshold. It is forbidden
	// for readable text.
	{label: "decor: faint on panel", fg: "--text-faint", bg: []string{"--surface-raised"}, min: large},
	{label: "decor: faint in inset", fg: "--text-faint", bg: []string{"--surface-sunken"}, min: large},

	// Badges: 11px, so the full text threshold applies.
	{label: "badge: accent on own background", fg: "--accent-text", bg: []string{"--surface-raised", "--accent-bg"}, min: text},
	{label: "badge: ok on own background", fg: "--ok-text", bg: []string{"--surface-raised", "--ok-bg"}, min: text},
	{label: "badge: warn on own background", fg: "--warn-text", bg: []string{"--surface-raised", "--warn-bg"}, min: text},
	{label: "badge: err on own background", fg: "--err-text", bg: []string{"--surface-raised", "--err-bg"}, min: text},
	{label: "badge: neutral in inset", fg: "--text-secondary", bg: []string{"--surface-sunken"}, min: text},

	// Status text also lives outside badges — metric delta, field error, footnote.
	{label: "status: ok-text on panel", fg: "--ok-text", bg: []string{"--surface-raised"}, min: text},
	{label: "status: warn-text on panel", fg: "--warn-text", bg: []string{"--surface-raised"}, min: text},
	{label: "status: err-text on panel", fg: "--err-text", bg: []string{"--surface-raised"}, min: text},
	{label: "status: ok-text in inset", fg: "--ok-text", bg: []string{"--surface-sunken"}, min: text},
	{label: "status: warn-text in inset", fg: "--warn-text", bg: []string{"--surface-sunken"}, min: text},
	{label: "status: err-text in inset", fg: "--err-text", bg: []string{"--surface-sunken"}, min: text},

	// Solid button. Hover must INCREASE label contrast, not reduce it.
	{label: "button: label on accent-solid", fg: "--accent-on", bg: []string{"--accent-solid"}, min: text},
	{label: "button: label on accent-hover", fg: "--accent-on", bg: []string{"--accent-hover"}, min: text},
	{label: "link: accent-text on page", fg: "--accent-text", bg: []string{"--surface-page"}, min: text},

	// State indicators — non-text, but non-decorative: 3:1. They are measured
	// with the same token they are painted with: the dot, measure fill, and
	// history stroke use --*-mark, not the text step.
	{label: "mark: ok on panel", fg: "--ok-mark", bg: []string{"--surface-raised"}, min: large},
	{label: "mark: ok in inset", fg: "--ok-mark", bg: []string{"--surface-sunken"}, min: large},
	{label: "mark: ok on track", fg: "--ok-mark", bg: []string{"--surface-raised", "--track"}, min: large},
	{label: "mark: warn on panel", fg: "--warn-mark", bg: []string{"--surface-raised"}, min: large},
	{label: "mark: warn in inset", fg: "--warn-mark", bg: []string{"--surface-sunken"}, min: large},
	{label: "mark: err on panel", fg: "--err-mark", bg: []string{"--surface-raised"}, min: large},
	{label: "mark: err in inset", fg: "--err-mark", bg: []string{"--surface-sunken"}, min: large},
	{label: "mark: err on track", fg: "--err-mark", bg: []string{"--surface-raised", "--track"}, min: large},
	{label: "dot: running on panel", fg: "--accent-mark", bg: []string{"--surface-raised"}, min: large},
	{label: "dot: running in inset", fg: "--accent-mark", bg: []string{"--surface-sunken"}, min: large},
	{label: "caret on panel", fg: "--accent-mark", bg: []string{"--surface-raised"}, min: large},
	{label: "slider on track", fg: "--accent-mark", bg: []string{"--surface-raised", "--track"}, min: large},

	// Non-decorative borders: the border IS the control.
	{label: "control border on panel", fg: "--border-control", bg: []string{"--surface-raised"}, min: large},
	{label: "control border on page", fg: "--border-control", bg: []string{"--surface-page"}, min: large},
	{label: "control border in inset", fg: "--border-control", bg: []string{"--surface-raised", "--surface-field"}, min: large},

	// Measure fill against its own track, and the track on every surface where
	// the measure can sit.
	{label: "measure: fill on track (panel)", fg: "--accent-mark", bg: []string{"--surface-raised", "--track"}, min: large},
	{label: "measure: fill on track (inset)", fg: "--accent-mark", bg: []string{"--surface-sunken", "--track"}, min: large},
	{label: "measure: ok on track", fg: "--ok-text", bg: []string{"--surface-raised", "--track"}, min: large},
	{label: "measure: warn on track", fg: "--warn-text", bg: []string{"--surface-raised", "--track"}, min: large},
	{label: "measure: err on track", fg: "--err-text", bg: []string{"--surface-raised", "--track"}, min: large},

	// Categorical palette: every series must separate from the surface.
	{label: "chart: series 1 on panel", fg: "--chart-1", bg: []string{"--surface-raised"}, min: large},
	{label: "chart: series 1 on page", fg: "--chart-1", bg: []string{"--surface-page"}, min: large},
	{label: "chart: series 2 on panel", fg: "--chart-2", bg: []string{"--surface-raised"}, min: large},
	{label: "chart: series 2 on page", fg: "--chart-2", bg: []string{"--surface-page"}, min: large},
	{label: "chart: series 3 on panel", fg: "--chart-3", bg: []string{"--surface-raised"}, min: large},
	{label: "chart: series 3 on page", fg: "--chart-3", bg: []string{"--surface-page"}, min: large},
	{label: "chart: series 4 on panel", fg: "--chart-4", bg: []string{"--surface-raised"}, min: large},
	{label: "chart: series 4 on page", fg: "--chart-4", bg: []string{"--surface-page"}, min: large},
	{label: "chart: series 5 on panel", fg: "--chart-5", bg: []string{"--surface-raised"}, min: large},
	{label: "chart: series 5 on page", fg: "--chart-5", bg: []string{"--surface-page"}, min: large},
	{label: "chart: series 6 on panel", fg: "--chart-6", bg: []string{"--surface-raised"}, min: large},
	{label: "chart: series 6 on page", fg: "--chart-6", bg: []string{"--surface-page"}, min: large},

	// Overlays: everything that sits on --surface-overlay.
	{label: "popover: text", fg: "--text-primary", bg: []string{"--surface-overlay"}, min: text},
	{label: "menu: keyboard shortcut", fg: "--text-muted", bg: []string{"--surface-overlay"}, min: text},
	{label: "menu: dangerous item", fg: "--err-text", bg: []string{"--surface-overlay"}, min: text},
	{label: "menu: marked item", fg: "--accent-text", bg: []string{"--surface-overlay"}, min: text},
	{label: "tooltip: text", fg: "--text-primary", bg: []string{"--surface-overlay"}, min: text},

	// Banner: text over a tinted fill.
	{label: "banner ok: heading", fg: "--text-primary", bg: []string{"--surface-page", "--ok-bg"}, min: text},
	{label: "banner warn: heading", fg: "--text-primary", bg: []string{"--surface-page", "--warn-bg"}, min: text},
	{label: "banner error: heading", fg: "--text-primary", bg: []string{"--surface-page", "--err-bg"}, min: text},
	{label: "banner warn: explanation", fg: "--text-secondary", bg: []string{"--surface-page", "--warn-bg"}, min: text},
	{label: "banner warn: icon", fg: "--warn-text", bg: []string{"--surface-page", "--warn-bg"}, min: large},

	// Forms.
	{label: "choice card: heading", fg: "--text-primary", bg: []string{"--surface-raised", "--accent-bg"}, min: text},
	{label: "choice card: description", fg: "--text-secondary", bg: []string{"--surface-raised", "--accent-bg"}, min: text},
	{label: "multi-select: selected item", fg: "--accent-text", bg: []string{"--surface-raised", "--surface-field", "--surface-selected"}, min: text},
	{label: "field prefix", fg: "--text-muted", bg: []string{"--surface-sunken"}, min: text},
	{label: "readonly: text in inset", fg: "--text-primary", bg: []string{"--surface-sunken"}, min: text},
	{label: "file zone dashed border", fg: "--border-control", bg: []string{"--surface-raised", "--surface-field"}, min: large},
	{label: "required marker", fg: "--err-text", bg: []string{"--surface-raised"}, min: text},

	// Layout and navigation.
	{label: "text on sidebar", fg: "--text-secondary", bg: []string{"--surface-sunken"}, min: text},
	{label: "navigation: current item", fg: "--accent-text", bg: []string{"--surface-sunken", "--surface-selected"}, min: text},
	{label: "navigation: edge marker", fg: "--accent-solid", bg: []string{"--surface-sunken", "--surface-selected"}, min: large},
	{label: "tab: underline", fg: "--accent-solid", bg: []string{"--surface-page"}, min: large},
	{label: "breadcrumbs: separator", fg: "--text-faint", bg: []string{"--surface-page"}, min: large},
	{label: "pagination: current page", fg: "--accent-text", bg: []string{"--surface-page", "--surface-selected"}, min: text},
	{label: "steps: completed bar", fg: "--accent-mark", bg: []string{"--surface-page", "--track"}, min: large},

	// Inverse plate: tooltip and everything that explains the interface without
	// being the interface itself.
	// The pair is its own because no text step sits on it: it has its own
	// foreground, and nothing else checks it.
	{label: "annotation: text on inverse", fg: "--text-on-inverse", bg: []string{"--surface-inverse"}, min: text},
	// ...and the plate itself must separate from what it floats over, otherwise
	// the meaning of "this is not content" is lost in the first dark theme.
	{label: "annotation: plate on page", fg: "--surface-inverse", bg: []string{"--surface-page"}, min: large},
	{label: "annotation: plate on panel", fg: "--surface-inverse", bg: []string{"--surface-raised"}, min: large},

	// Focus ring — against what is underneath it.
	{label: "focus: ring on page", fg: "--focus-ring", bg: []string{"--surface-page"}, min: large},
	{label: "focus: ring on panel", fg: "--focus-ring", bg: []string{"--surface-raised"}, min: large},
	{label: "focus: ring in inset", fg: "--focus-ring", bg: []string{"--surface-sunken"}, min: large},
	// ...and against what it OUTLINES. This pair was not here, so the gate
	// stayed green for years while in light themes --focus-ring and
	// --accent-solid were the same color: contrast 1.00 on the main
	// screen button. The threshold is distinguishability; see the constant.
	{label: "focus: ring around primary fill", fg: "--focus-ring", bg: []string{"--accent-solid"}, min: distinct},

	// ── SURFACE STACK STEPS ─────────────────────────────────────────────────
	//
	// These pairs were not here, and the gap followed directly from the design
	// of the check: every case is a FOREGROUND on a BACKGROUND, meaning text or
	// a marker on a surface. Surface against SURFACE was never checked.
	//
	// The cost: panel and card both sat on --surface-raised, the difference was
	// exactly 1.00, and on the "Height and surfaces" page the card inside the
	// panel was not readable. The gate was green — all its pairs honestly
	// passed, because none of them asked "are two neighboring tiers
	// distinguishable?"
	//
	// The kit's first surface law — "depth is conveyed by lightness order" —
	// was not expressed by any check. Now it is.
	//
	// The threshold is distinguishability (see the constant), not accessibility:
	// WCAG says nothing about adjacent surfaces, and should not. 1.5 is not
	// suitable here — it is the threshold for a ring over a FILL, where nothing
	// else is adjacent. For a stack, the border also helps neighboring layers,
	// so a smaller value is sufficient, but it must NOT be ONE: one means there
	// is no step at all.
	{label: "stack: page over recess", fg: "--surface-page", bg: []string{"--surface-sunken"}, min: step},
	{label: "stack: panel over page", fg: "--surface-raised", bg: []string{"--surface-page"}, min: step},
	{label: "stack: panel over recess", fg: "--surface-raised", bg: []string{"--surface-sunken"}, min: step},
	{label: "stack: card in panel", fg: "--surface-sunken", bg: []string{"--surface-raised"}, min: step},

	// The step between BUTTON WEIGHTS, not between surface stack layers. The
	// weight ladder sits on the same recesses, and its steps must separate just
	// as well: when they converge, soft and default become two names for one
	// appearance.
	//
	// This has to be checked here because weights are composed ON TOP of the
	// background and are not part of the surface stack. The pair also catches
	// the reverse case: at the absolute step, soft went deeper than default,
	// and the ladder flipped.
	{label: "ladder: soft against default", fg: "--surface-recessed-hover", bg: []string{"--surface-raised"},
		alt: []string{"--surface-raised", "--surface-recessed"}, min: step},

	// BUTTON AGAINST WHAT IT SITS ON. This pair was not here, and its absence
	// was costly.
	//
	// Above, surface stack layers are checked against each other, and button
	// weights are checked against each other. Neither asks the main question:
	// does the button differ from its background? For the default, this is its
	// ONLY distinguishing feature — there is no border, no shadow, both removed
	// in the failure.
	//
	// The failure is caused by a black film in both themes, because the recess
	// is conceptually darker. But the film subtracts lightness, and at the
	// bottom of the ramp there is nothing left to subtract: the dark theme page
	// is rgb(12,12,12), and even SOLID black gives 1.073 against it versus 1.142
	// for the light ones. The button stopped existing: difference 1.014.
	//
	// Both backgrounds, because the button sits on both: on the panel — in the
	// form and in the card, on the page — in the screen header and shell toolbar.
	{label: "button against panel", fg: "--surface-recessed", bg: []string{"--surface-raised"}, min: step},
	{label: "button against page", fg: "--surface-recessed", bg: []string{"--surface-page"}, min: step},
	// There is NO inset field here, and this is not an omission. The field
	// surface matches the raised one: the control border identifies the field,
	// not depth, and the "control border" pairs above handle it. Requiring a
	// step from the field would mean requiring depth where the kit intentionally
	// does not provide it.

}

var themes = []*css.Theme{
	{ID: "light-neutral", Label: "light neutral", Scheme: "light"},
	{ID: "light", Label: "light warm", Scheme: "light"},
	{ID: "light-cool", Label: "light cool", Scheme: "light"},
	{ID: "dark-soft", Label: "dark gray", Scheme: "dark"},
	{ID: "dark", Label: "dark black", Scheme: "dark"},
}

// Accent is the second axis of the check, and it must be complete here.
//
// While there was one set, each accent pair was checked exactly in the form
// it was drawn. With a four-position control, that stops working: three
// quarters of the accent pairs go into production untested, and their
// thresholds are DIFFERENT — the lightness ceiling of a fill under a white
// label is hue-dependent (0.545 on petrol versus 0.580 on clay), and a set
// assembled from the numbers of a neighboring hue will honestly fail.
//
// Empty ID is the base: the set declared in :root without an attribute. It must
// come first because detailed output shows exactly that one.
var accents = []struct{ ID, Label string }{
	{"", "petrol"},
	{"graphite", "graphite"},
	{"indigo", "indigo"},
	{"clay", "clay"},
}

func main() {
	tokens := flag.String("tokens", "../src/tokens.css", "path to tokens.css")
	// Full output is 528 lines per set, meaning more than two thousand per run.
	// The base accent is printed in detail; for the others, only failures and
	// the result are printed: a list nobody reads is not a check, it is noise.
	verbose := flag.Bool("v", false, "print all pairs for each accent")
	flag.Parse()

	src, err := css.Load(*tokens)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot read tokens:", err)
		os.Exit(1)
	}

	// Base is the first :root block. [data-density] blocks and media queries do
	// NOT belong here: otherwise control heights would overwrite base
	// declarations.
	base := src.Decls(regexp.MustCompile(`(?m)^:root \{`))
	for _, t := range themes {
		t.Base = base
		t.Vars = src.Decls(regexp.MustCompile(`\[data-theme="` + t.ID + `"\]\s*\{`))
	}

	width := 0
	for _, c := range cases {
		if n := utf8.RuneCountInString(c.label); n > width {
			width = n
		}
	}

	failed, total := 0, 0
	for _, a := range accents {
		vars := map[string]string{}
		if a.ID != "" {
			vars = src.Decls(regexp.MustCompile(`\[data-accent="` + a.ID + `"\]\s*\{`))
			if len(vars) == 0 {
				fmt.Fprintf(os.Stderr, "accent %s is declared in the check but not in tokens\n", a.ID)
				os.Exit(1)
			}
		}
		detail := *verbose || a.ID == ""

		fmt.Printf("\n═══ ACCENT %s ═══\n", a.Label)
		for _, t := range themes {
			t.Accent, t.AccentID, t.AccentLabel = vars, a.ID, a.Label

			bad := 0
			var lines []string
			for _, c := range cases {
				total++
				pad := strings.Repeat(" ", width-utf8.RuneCountInString(c.label))

				fg, err := t.Token(c.fg)
				if err == nil {
					var bg css.RGBA
					bg, err = t.Flatten(c.bg)
					if err == nil {
						// The measure depends on the case: surface stack steps are
						// measured by lightness, everything else by contrast ratio.
						r, unit := css.Ratio(fg, bg), ""
						if c.isStep() {
							r, unit = css.Step(fg, bg), " ΔL"
							if len(c.alt) > 0 {
								var other css.RGBA
								if other, err = t.Flatten(c.alt); err == nil {
									// The presence of each alternative is the
									// step against the shared background. The
									// difference is taken WITH ITS SIGN:
									// negative means a reversed ladder and
									// must fail.
									quiet := css.Step(fg, bg)
									loud := css.StepOf(other, bg)
									r = loud - quiet
								}
							}
						}
						if err != nil {
							failed, bad = failed+1, bad+1
							lines = append(lines, fmt.Sprintf("  ✗ %s%s  ERROR: %v", c.label, pad, err))
							continue
						}
						mark := "·"
						if r < c.min {
							mark, failed, bad = "✗", failed+1, bad+1
						}
						if detail || r < c.min {
							lines = append(lines, fmt.Sprintf("  %s %s%s  %6.3f%s  (need %.3f)", mark, c.label, pad, r, unit, c.min))
						}
						continue
					}
				}
				failed, bad = failed+1, bad+1
				lines = append(lines, fmt.Sprintf("  ✗ %s%s  ERROR: %v", c.label, pad, err))
			}

			if detail {
				fmt.Printf("\nTHEME %s — %s\n", t.ID, t.Label)
				fmt.Println(strings.Repeat("─", width+26))
			} else if bad > 0 {
				fmt.Printf("\nTHEME %s — %s\n", t.ID, t.Label)
			}
			for _, l := range lines {
				fmt.Println(l)
			}
			if !detail && bad == 0 {
				fmt.Printf("  · %s — %d pairs passed\n", t.Label, len(cases))
			}
		}
	}

	// Table coverage: a text color that is absent from it is a threshold that
	// nobody measured. Count it once, not once per theme: the question is not
	// the value, but whether the row exists.
	gaps, inkCount, err := checkInkCoverage(filepath.Dir(*tokens), cases)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot read kit:", err)
		os.Exit(1)
	}
	total += inkCount
	failed += len(gaps)
	if len(gaps) > 0 {
		fmt.Println("\nTABLE COVERAGE")
		for _, g := range gaps {
			fmt.Println(g)
		}
	}

	fmt.Println()
	if failed > 0 {
		fmt.Printf("✗ failures: %d of %d\n", failed, total)
		os.Exit(1)
	}
	fmt.Printf("· all %d checks passed: %d themes × %d accents, coverage %d text colors\n",
		total, len(themes), len(accents), inkCount)

}

// ── Coverage: every text color must be checked by something ────────────────
//
// The pair table above is handwritten, and that is correct: a pair carries not
// only two tokens, but also a COMPOSITION STACK ("label in field on panel" —
// two backgrounds, because the inset is translucent and combines with the
// panel beneath it). That fact about markup nesting is not derivable from
// tokens.css, and generating pairs would strip the table of the human name
// every row has.
//
// But a GAP in it is derivable, and costs almost nothing to detect. A component
// that paints text with a new token nobody considered will reveal itself
// nowhere: the color will apply, nobody will measure its threshold, and this
// will be discovered by a person who cannot read the label.
//
// Therefore this checks not the result, but COVERAGE: every token the kit uses
// to paint text must appear as a foreground in at least one pair. What lies
// beneath it is still decided by the pair author.
var (
	colorDecl = regexp.MustCompile(`(?:^|[;{])\s*color\s*:\s*([^;}]+)`)
	varUse    = regexp.MustCompile(`var\(\s*(--[a-z][\w-]*)`)
)

// inkTokens — tokens the kit actually uses to paint text, with the location of
// the first occurrence. tokens.css is excluded: color is declared there, not
// applied; print.css is excluded too, because paper has its own set.
func inkTokens(dir string) (map[string]string, error) {
	out := map[string]string{}
	ents, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, e := range ents {
		name := e.Name()
		if !strings.HasSuffix(name, ".css") || name == "tokens.css" || name == "print.css" {
			continue
		}
		b, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return nil, err
		}
		// Comments are BLANKED rather than cut out: they take three quarters of
		// the file, and after cutting the line number drifts by hundreds of
		// lines — that is, the message points at somebody else's rule.
		text := string(css.Blank([]byte(strings.ReplaceAll(string(b), "\r\n", "\n"))))
		for _, m := range colorDecl.FindAllStringSubmatchIndex(text, -1) {
			value := text[m[2]:m[3]]
			line := strings.Count(text[:m[2]], "\n") + 1
			for _, v := range varUse.FindAllStringSubmatch(value, -1) {
				if _, seen := out[v[1]]; !seen {
					out[v[1]] = fmt.Sprintf("%s:%d", name, line)
				}
			}
		}
	}
	return out, nil
}

// checkInkCoverage returns tokens used to paint text but not used as a
// foreground in any pair.
//
// Component variables are skipped: --btn-fg and --tone-ink are not colors, but
// SUBSTITUTIONS — semantics sits behind them and is measured under its own
// name. Checking them would mean requiring a pair for every intermediary name.
func checkInkCoverage(dir string, cases []kase) ([]string, int, error) {
	ink, err := inkTokens(dir)
	if err != nil {
		return nil, 0, err
	}
	covered := map[string]bool{}
	for _, c := range cases {
		covered[c.fg] = true
	}
	var bad []string
	for tok, where := range ink {
		if indirect.MatchString(tok) || covered[tok] {
			continue
		}
		bad = append(bad, fmt.Sprintf(
			"  ✗ %s  paints text (%s), but is not a foreground in any pair.\n"+
				"      Add a row to the table above: an unchecked text color is a threshold nobody measured",
			tok, where))
	}
	sort.Strings(bad)
	return bad, len(ink), nil
}

// Intermediary names: component variable and tone. Semantics sits behind them,
// which is measured separately under its own name.
var indirect = regexp.MustCompile(`^--(btn|tone|level|change|chart)-`)

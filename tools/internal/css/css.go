// Package css parses custom-property declarations out of the kit's real CSS
// and resolves them the way a browser does.
//
// The point is that the checks must NOT DUPLICATE the kit's values but compute
// them: a duplicated value drifts away from the original silently, and at the
// moment when nobody is looking at it.
package css

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var commentRe = regexp.MustCompile(`(?s)/\*.*?\*/`)

// Source is a parsed token file.
type Source struct{ text string }

func Load(path string) (*Source, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return &Source{text: commentRe.ReplaceAllString(string(b), "")}, nil
}

// bodyAt returns the body of the block starting at position idx, with
// balanced braces.
func (s *Source) bodyAt(idx int) string {
	depth, start := 0, -1
	for i := idx; i < len(s.text); i++ {
		switch s.text[i] {
		case '{':
			if depth == 0 {
				start = i + 1
			}
			depth++
		case '}':
			depth--
			if depth == 0 {
				return s.text[start:i]
			}
		}
	}
	return ""
}

var declRe = regexp.MustCompile(`(--[\w-]+)\s*:\s*([^;]+);`)
var spaceRe = regexp.MustCompile(`\s+`)

// Decls extracts custom properties from the FIRST block matching re.
//
// The first one specifically: [data-density] blocks and media queries must not
// get in here, or the control heights would overwrite the base declarations.
func (s *Source) Decls(re *regexp.Regexp) map[string]string {
	out := map[string]string{}
	loc := re.FindStringIndex(s.text)
	if loc == nil {
		return out
	}
	for _, m := range declRe.FindAllStringSubmatch(s.bodyAt(loc[0]), -1) {
		out[m[1]] = spaceRe.ReplaceAllString(strings.TrimSpace(m[2]), " ")
	}
	return out
}

// Theme is a theme as two independent handles rather than a separate set of
// tokens.
type Theme struct {
	ID, Label, Scheme string
	Vars              map[string]string

	// The accent set is a SEPARATE axis rather than a kind of theme, because it
	// combines with any of them: 4 accents × 5 themes. In the kit itself this is
	// [data-accent="…"], declared next to the theme blocks.
	//
	// It cannot collide with Vars by construction: a theme overrides the
	// SEMANTICS (--accent-text and the like), an accent overrides the RAMP
	// (--a-*), and the semantics read the ramp through var(). So the order
	// between these two maps changes nothing, and they have nothing to argue
	// about.
	Accent      map[string]string
	AccentID    string
	AccentLabel string

	Base map[string]string
}

func (t *Theme) lookup(name string) (string, error) {
	if v, ok := t.Vars[name]; ok {
		return v, nil
	}
	if v, ok := t.Accent[name]; ok {
		return v, nil
	}
	if v, ok := t.Base[name]; ok {
		return v, nil
	}
	return "", fmt.Errorf("token %s is not declared", name)
}

var varRe = regexp.MustCompile(`var\((--[\w-]+)\)`)

// expand substitutes var() textually — the way a browser does it BEFORE
// parsing the function. Without that, oklch(0.994 0.002 var(--hue-neutral))
// cannot be parsed.
func (t *Theme) expand(value string, depth int) (string, error) {
	if depth > 30 {
		return "", fmt.Errorf("var() substitution is nested too deep")
	}
	if !strings.Contains(value, "var(") {
		return value, nil
	}
	var err error
	out := varRe.ReplaceAllStringFunc(value, func(m string) string {
		name := varRe.FindStringSubmatch(m)[1]
		v, e := t.lookup(name)
		if e != nil && err == nil {
			err = e
		}
		return v
	})
	if err != nil {
		return "", err
	}
	return t.expand(out, depth+1)
}

// RGBA is a colour in sRGB with alpha, components 0…1.
type RGBA struct{ R, G, B, A float64 }

func splitArgs(s string) []string {
	var out []string
	depth := 0
	var cur strings.Builder
	for _, ch := range s {
		switch ch {
		case '(':
			depth++
		case ')':
			depth--
		}
		if ch == ',' && depth == 0 {
			out = append(out, strings.TrimSpace(cur.String()))
			cur.Reset()
			continue
		}
		cur.WriteRune(ch)
	}
	if strings.TrimSpace(cur.String()) != "" {
		out = append(out, strings.TrimSpace(cur.String()))
	}
	return out
}

func inner(s string) string {
	return s[strings.Index(s, "(")+1 : strings.LastIndex(s, ")")]
}

// OklchToSrgb converts to sRGB. The matrices and the gamma are a browser's.
func OklchToSrgb(L, C, H, alpha float64) RGBA {
	h := H * math.Pi / 180
	a := C * math.Cos(h)
	b := C * math.Sin(h)
	l := math.Pow(L+0.3963377774*a+0.2158037573*b, 3)
	m := math.Pow(L-0.1055613458*a-0.0638541728*b, 3)
	s := math.Pow(L-0.0894841775*a-1.2914855480*b, 3)

	enc := func(v float64) float64 {
		if v <= 0.0031308 {
			v = 12.92 * v
		} else {
			v = 1.055*math.Pow(math.Max(v, 0), 1.0/2.4) - 0.055
		}
		return math.Min(1, math.Max(0, v))
	}
	return RGBA{
		R: enc(4.0767416621*l - 3.3077115913*m + 0.2309699292*s),
		G: enc(-1.2684380046*l + 2.6097574011*m - 0.3413193965*s),
		B: enc(-0.0041960863*l - 0.7034186147*m + 1.7076147010*s),
		A: alpha,
	}
}

// Lightness returns the lightness of a colour in OKLCH — the same L axis the
// kit's ramp is set on.
//
// It is needed separately from Ratio because the step between surfaces and the
// readability of text are DIFFERENT measures, and one must not stand in for
// the other.
//
// The WCAG formula carries a 0.05 term — a model of glare and scattered light
// on a screen. It is right for text, but at the very bottom of the lightness
// range it squeezes ALL ratios towards one: a step of 0.035 in lightness gives
// 1.10 at the white end of the ramp and 1.04 at the black one, although the eye
// sees them as the same. A threshold set by the light end then declares the
// dark themes broken, and one set by the dark end stops catching anything in
// the light ones.
//
// OKLCH is perceptually uniform by construction: a difference in lightness
// means the same thing at both ends, and one threshold works for all five
// themes.
func Lightness(c RGBA) float64 {
	dec := func(v float64) float64 {
		if v <= 0.04045 {
			return v / 12.92
		}
		return math.Pow((v+0.055)/1.055, 2.4)
	}
	r, g, b := dec(c.R), dec(c.G), dec(c.B)
	l := math.Cbrt(0.4122214708*r + 0.5363325363*g + 0.0514459929*b)
	m := math.Cbrt(0.2119034982*r + 0.6806995451*g + 0.1073969566*b)
	s := math.Cbrt(0.0883024619*r + 0.2817188376*g + 0.6299787005*b)
	return 0.2104542553*l + 0.7936177850*m - 0.0040720468*s
}

// Step is the lightness difference between what WAS LAID DOWN and what it was
// laid on.
//
// Compositing is mandatory, and its absence cost dearly. The kit's recessed
// surfaces are films: --surface-recessed is black with alpha rather than a
// colour. Without compositing, such a film contributed its own lightness, that
// is zero, and the step against any backdrop came out enormous — the check
// passed always and verified nothing. That is how the "button weight ladder"
// stayed quiet while the default in the dark themes differed from the
// background by 1.4%.
//
// For an opaque foreground the compositing returns it unchanged, so the pairs
// from the surface stack are indifferent to this.
func Step(a, b RGBA) float64 { return math.Abs(Lightness(Composite(a, b)) - Lightness(b)) }

// calcRe matches a simple arithmetic expression whole.
//
// The kit uses calc() for exactly one thing: multiplying the chroma of a ramp
// step by the tint handle. A full CSS evaluator is not needed here and would be
// a lie about the capabilities — numbers and four operations are supported, and
// everything else honestly fails with an error.
var calcRe = regexp.MustCompile(`calc\(([^()]*)\)`)

// evalCalc collapses calc() into a number. Multiplication and division come
// before addition and subtraction — as in arithmetic and as in CSS.
func evalCalc(expr string) (float64, bool) {
	f := strings.Fields(expr)
	if len(f) == 0 || len(f)%2 == 0 {
		return 0, false
	}
	num := func(s string) (float64, bool) { v, err := strconv.ParseFloat(s, 64); return v, err == nil }

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

var (
	oklchRe = regexp.MustCompile(`^oklch\(\s*([\d.]+)\s+([\d.]+)\s+([\d.]+)\s*(?:/\s*([\d.]+)\s*)?\)$`)
	mixRe   = regexp.MustCompile(`^(.*?)\s+([\d.]+)%$`)
)

// Resolve resolves a value into an sRGB colour for a theme.
func (t *Theme) Resolve(value string) (RGBA, error) {
	v, err := t.expand(value, 0)
	if err != nil {
		return RGBA{}, err
	}
	v = strings.TrimSpace(v)

	// calc() is collapsed BEFORE the colour is parsed: by the time it paints, a
	// browser does the same, and oklch() receives a number already.
	for calcRe.MatchString(v) {
		bad := false
		v = calcRe.ReplaceAllStringFunc(v, func(m string) string {
			n, ok := evalCalc(calcRe.FindStringSubmatch(m)[1])
			if !ok {
				bad = true
				return m
			}
			return strconv.FormatFloat(n, 'f', -1, 64)
		})
		if bad {
			return RGBA{}, fmt.Errorf("cannot evaluate the calc: %s", v)
		}
	}

	if v == "transparent" {
		return RGBA{}, nil
	}

	if strings.HasPrefix(v, "light-dark(") {
		args := splitArgs(inner(v))
		if len(args) != 2 {
			return RGBA{}, fmt.Errorf("light-dark expects two arguments: %s", v)
		}
		if t.Scheme == "dark" {
			return t.Resolve(args[1])
		}
		return t.Resolve(args[0])
	}

	if strings.HasPrefix(v, "color-mix(") {
		args := splitArgs(inner(v))
		if len(args) != 3 || !strings.EqualFold(args[0], "in oklab") {
			return RGBA{}, fmt.Errorf(`only "in oklab" is supported: %s`, v)
		}
		m := mixRe.FindStringSubmatch(args[1])
		if m == nil {
			return RGBA{}, fmt.Errorf("cannot parse the share: %s", args[1])
		}
		pct, _ := strconv.ParseFloat(m[2], 64)
		pct /= 100
		base, err := t.Resolve(m[1])
		if err != nil {
			return RGBA{}, err
		}
		other, err := t.Resolve(args[2])
		if err != nil {
			return RGBA{}, err
		}
		a := base.A*pct + other.A*(1-pct)
		if a == 0 {
			return RGBA{}, nil
		}
		return RGBA{
			R: (base.R*base.A*pct + other.R*other.A*(1-pct)) / a,
			G: (base.G*base.A*pct + other.G*other.A*(1-pct)) / a,
			B: (base.B*base.A*pct + other.B*other.A*(1-pct)) / a,
			A: a,
		}, nil
	}

	if m := oklchRe.FindStringSubmatch(v); m != nil {
		l, _ := strconv.ParseFloat(m[1], 64)
		c, _ := strconv.ParseFloat(m[2], 64)
		h, _ := strconv.ParseFloat(m[3], 64)
		alpha := 1.0
		if m[4] != "" {
			alpha, _ = strconv.ParseFloat(m[4], 64)
		}
		return OklchToSrgb(l, c, h, alpha), nil
	}

	return RGBA{}, fmt.Errorf("cannot parse the colour: %s", v)
}

// Token resolves a token by name.
func (t *Theme) Token(name string) (RGBA, error) {
	v, err := t.lookup(name)
	if err != nil {
		return RGBA{}, err
	}
	return t.Resolve(v)
}

// Composite lays fg over bg.
func Composite(fg, bg RGBA) RGBA {
	return RGBA{
		R: fg.R*fg.A + bg.R*(1-fg.A),
		G: fg.G*fg.A + bg.G*(1-fg.A),
		B: fg.B*fg.A + bg.B*(1-fg.A),
		A: fg.A + bg.A*(1-fg.A),
	}
}

// Flatten collapses a stack of tokens into an opaque colour. The first is the
// base.
func (t *Theme) Flatten(stack []string) (RGBA, error) {
	var out RGBA
	for i, name := range stack {
		c, err := t.Token(name)
		if err != nil {
			return RGBA{}, err
		}
		if i == 0 {
			out = c
			continue
		}
		out = Composite(c, out)
	}
	return out, nil
}

func luminance(c RGBA) float64 {
	f := func(v float64) float64 {
		if v <= 0.04045 {
			return v / 12.92
		}
		return math.Pow((v+0.055)/1.055, 2.4)
	}
	return 0.2126*f(c.R) + 0.7152*f(c.G) + 0.0722*f(c.B)
}

// Ratio is the contrast of a foreground against a background, per WCAG.
func Ratio(fg, bg RGBA) float64 {
	solid := Composite(fg, bg)
	hi, lo := luminance(solid), luminance(bg)
	if lo > hi {
		hi, lo = lo, hi
	}
	return (hi + 0.05) / (lo + 0.05)
}

// ── Geometry ────────────────────────────────────────────────────────────────

var (
	varUseRe = regexp.MustCompile(`var\(\s*(--[a-z][\w-]*)\s*\)`)
	pxNumRe  = regexp.MustCompile(`([\d.]+)px`)
	remNumRe = regexp.MustCompile(`([\d.]+)rem`)
)

// RootPx is the root that rem is measured against.
//
// The kit keeps exactly one tier in rem — the type sizes — and it does so in
// order to respect an enlarged font size in the browser. The check needs a
// concrete root, or there is nothing to compare a type size with: the geometry
// is declared in px. 16 is every browser's default, and the comments in
// tokens.css are tied to it as well ("0.875rem × 16 = 14px").
//
// A VARIABLE rather than a constant, and that matters. A consumer is entitled
// to change the root: it is the only lever with which an application makes
// itself denser in TEXT without touching the geometry — control heights and tap
// targets stay in px and do not shrink along with the letters. But the lever is
// not safe: at a root of 14 the smallest step falls to 9.63px, that is below
// the floor the kit declares, and the icon stops landing inside its band
// against the cap height.
//
// While the root was a constant there was nothing to check somebody else's
// choice with. Now there is: `proportion -root 14` computes the whole ladder
// the way a browser will see it.
var RootPx = 16.0

// ResolvePx expands var() and evaluates calc(), returning a number of pixels.
//
// It lives here rather than in a command because there are two consumers now:
// cmd/targets measures tap targets and cmd/proportion measures proportions. A
// copy of the resolver in each of them would be exactly the "future
// divergence" the design principles warn about: one of them would fix the calc
// parsing and the other would not.
//
// The values here are always geometric, that is in px: colours and percentages
// never arrive, and that is why the parsing fits into thirty lines rather than
// the colour resolver above.
func ResolvePx(vals map[string]string, expr string) (float64, error) {
	if strings.HasPrefix(expr, "--") {
		v, ok := vals[expr]
		if !ok {
			return 0, fmt.Errorf("no token %s", expr)
		}
		expr = v
	}
	for i := 0; i < 12 && strings.Contains(expr, "var("); i++ {
		var missing string
		expr = varUseRe.ReplaceAllStringFunc(expr, func(m string) string {
			name := varUseRe.FindStringSubmatch(m)[1]
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
			v, ok := evalCalc(stripPx(calcRe.FindStringSubmatch(m)[1]))
			if !ok {
				return m
			}
			return fmt.Sprintf("%gpx", v)
		})
		if expr == before {
			return 0, fmt.Errorf("cannot parse the calc in %q", expr)
		}
	}
	if v, ok := evalCalc(stripPx(expr)); ok {
		return v, nil
	}
	return 0, fmt.Errorf("cannot parse %q", expr)
}

// stripPx reduces a length to a dimensionless number of pixels: px loses its
// suffix and rem is multiplied by the root. That is why a mixed expression such
// as calc(1rem + 2px) is computed correctly rather than by luck.
func stripPx(s string) string {
	s = remNumRe.ReplaceAllStringFunc(strings.TrimSpace(s), func(m string) string {
		v, err := strconv.ParseFloat(remNumRe.FindStringSubmatch(m)[1], 64)
		if err != nil {
			return m
		}
		return strconv.FormatFloat(v*RootPx, 'g', -1, 64)
	})
	return pxNumRe.ReplaceAllString(s, "$1")
}

// Blank replaces comment bodies with spaces, PRESERVING the line breaks.
//
// It is needed wherever a check reports a LINE NUMBER. Cutting comments out
// entirely will not do: in the kit they take three quarters of a file, and
// after cutting the number drifts by hundreds of lines — the message points at
// somebody else's rule and costs more than no message at all.
//
// Comments in CSS do not nest, so a single cursor pass is enough.
func Blank(src []byte) []byte {
	out := make([]byte, len(src))
	copy(out, src)
	for i := 0; i+1 < len(out); {
		if out[i] != '/' || out[i+1] != '*' {
			i++
			continue
		}
		j := i
		for ; j+1 < len(out); j++ {
			if out[j] == '*' && out[j+1] == '/' {
				j += 2
				break
			}
		}
		if j+1 >= len(out) {
			j = len(out)
		}
		for k := i; k < j; k++ {
			if out[k] != '\n' {
				out[k] = ' '
			}
		}
		i = j
	}
	return out
}

// StepOf is the lightness difference between two ALREADY opaque colours.
//
// It differs from Step in compositing nothing: both sides are assembled in
// advance. It is needed for alternatives — two states of one element on one
// backdrop, where the question is not "how much does the film change the
// background" but "how much do these two looks differ from each other".
func StepOf(a, b RGBA) float64 { return math.Abs(Lightness(a) - Lightness(b)) }

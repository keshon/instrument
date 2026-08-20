// The dist command builds the kit into a single file.
//
// The kit is connected with one line and requires no build BY THE CONSUMER —
// that remains true. But "no build for the consumer" and "no built file in
// the repository" are different things, and the latter used to be costly:
//
//	request per file instead of one (@import — a waterfall: the browser first
//	fetches kit.css and only then learns about the rest);
//	half a megabyte of raw source, three quarters of which are comments
//	in Russian.
//
// src/ remains the source. dist/ is output, and it lives in the repository
//
//	only so the consumer can take one file.
//
//	go run ./cmd/dist            build
//	go run ./cmd/dist -check     fail if the built output diverged from source
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"instrument/tools/internal/css"
)

var (
	layerStmtRe = regexp.MustCompile(`@layer\s+([^;{]+);`)
	importRe    = regexp.MustCompile(`@import\s+url\(["']?\./([^"')]+)["']?\)\s*layer\(([^)]+)\)\s*;`)
)

func main() {
	var (
		src     = flag.String("src", "../src", "kit directory")
		out     = flag.String("out", "../dist", "output directory")
		version = flag.String("version", "", "version; defaults to VERSION")
		check   = flag.Bool("check", false, "do not write, only compare")
	)
	flag.Parse()

	ver := *version
	if ver == "" {
		b, err := os.ReadFile(filepath.Join(filepath.Dir(*src), "VERSION"))
		if err != nil {
			fmt.Fprintln(os.Stderr, "cannot read VERSION:", err)
			os.Exit(1)
		}
		ver = strings.TrimSpace(string(b))
	}

	// The version is declared TWICE: in VERSION and in package.json. There is
	// nothing to make them diverge — except forgetfulness, and that happens
	// exactly at release time.
	//
	// The cost of divergence is not cosmetic: the dist/ header prints VERSION,
	// while npm publishes package.json, and a package containing a different
	// version is sent to the registry. After that it cannot be traced: the CDN
	// file carries one number, while the comment inside it names another.
	if err := checkPkgVersion(filepath.Dir(*src), ver); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	entry, err := os.ReadFile(filepath.Join(*src, "kit.css"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot read kit.css:", err)
		os.Exit(1)
	}

	// Layer order comes from kit.css itself rather than being rewritten here:
	// a second list would diverge from the first, and diverge silently.
	stmt := layerStmtRe.FindSubmatch(entry)
	if stmt == nil {
		fmt.Fprintln(os.Stderr, "kit.css has no @layer declaration")
		os.Exit(1)
	}

	var css bytes.Buffer
	fmt.Fprintf(&css, "/*! instrument %s — https://github.com/keshon/instrument\n"+
		"    Built from src/. Edit the SOURCE, not this file. */\n", ver)
	fmt.Fprintf(&css, "@layer %s;\n\n", strings.Join(strings.Fields(string(stmt[1])), " "))

	imports := importRe.FindAllSubmatch(entry, -1)
	if len(imports) == 0 {
		fmt.Fprintln(os.Stderr, "kit.css has no @import entries")
		os.Exit(1)
	}
	for _, m := range imports {
		name, layer := string(m[1]), string(m[2])
		b, err := os.ReadFile(filepath.Join(*src, name))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		// Line endings are normalized on input. Otherwise the build depends on
		// the system it was built on: git delivers source with CRLF on Windows,
		// and -check would not match itself.
		b = bytes.ReplaceAll(b, []byte("\r\n"), []byte("\n"))
		if err := checkBraces(name, b); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := checkRoleTier(name, b); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := checkBans(name, b); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		// Importing into a layer nests ALL file contents in that layer — that
		// is exactly what the block reproduces. Otherwise the layers would
		// diverge, and the kit's layer order carries meaning: motion and print
		// must override components.
		fmt.Fprintf(&css, "/* ── %s → %s ─────────────────────────────── */\n", name, layer)
		fmt.Fprintf(&css, "@layer %s {\n%s\n}\n\n", layer, b)
	}

	js, err := os.ReadFile(filepath.Join(*src, "kit.js"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot read kit.js:", err)
		os.Exit(1)
	}
	js = bytes.ReplaceAll(js, []byte("\r\n"), []byte("\n"))

	files := map[string][]byte{
		"instrument.css":     css.Bytes(),
		"instrument.min.css": []byte(minify(css.String())),
		"instrument.js":      js,
	}

	if *check {
		bad := 0
		for name, want := range files {
			got, err := os.ReadFile(filepath.Join(*out, name))
			if err != nil {
				fmt.Fprintf(os.Stderr, "  missing %s: %v\n", name, err)
				bad++
				continue
			}
			if !bytes.Equal(bytes.ReplaceAll(got, []byte("\r\n"), []byte("\n")), want) {
				fmt.Fprintf(os.Stderr, "  %s diverged from source\n", name)
				bad++
			}
		}
		if bad > 0 {
			fmt.Fprintln(os.Stderr, "\nBuilt output is behind src/. Rebuild:")
			fmt.Fprintln(os.Stderr, "  go -C tools run ./cmd/dist")
			os.Exit(1)
		}
		fmt.Println("· dist matches src")
		return
	}

	if err := os.MkdirAll(*out, 0o755); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for name, b := range files {
		if err := os.WriteFile(filepath.Join(*out, name), b, 0o644); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	fmt.Printf("instrument %s\n", ver)
	for _, name := range []string{"instrument.css", "instrument.min.css", "instrument.js"} {
		fmt.Printf("  %-18s %6.1f KB\n", name, float64(len(files[name]))/1024)
	}
	fmt.Printf("\nbefore: %d files, %.1f KB, the same number of requests\n",
		len(imports)+1, float64(len(css.Bytes()))/1024)

}

/*
minify — intentionally cautious.

It removes comments and extra whitespace and DOES NOT TOUCH anything else.
The temptation to squeeze more out is strong, but every next optimization is
dangerous on this particular CSS:

  - whitespace before ":" cannot be removed. ".a :hover" and ".a:hover" are
    different selectors, and distinguishing them without parsing the selector
    is impossible;
  - string contents are untouchable. The kit draws shapes with data-URI masks,
    and spaces, ";", and "}" live inside them — collapsing them means erasing
    half the glyphs;
  - "url(...)" without quotes is the same.

Comments make up three quarters of the file, so this caution has ample margin.
*/
func minify(s string) string {
	var b strings.Builder
	b.Grow(len(s))

	var quote byte     // current quote, 0 — outside a string
	inURL := false     // inside url(...) without quotes
	pendingWS := false // saw whitespace, have not decided whether to write it
	keepFirst := true  // first /*! comment is preserved — it is the header

	writeByte := func(c byte) { b.WriteByte(c) }

	for i := 0; i < len(s); i++ {
		c := s[i]

		if quote != 0 {
			writeByte(c)
			if c == '\\' && i+1 < len(s) {
				i++
				writeByte(s[i])
			} else if c == quote {
				quote = 0
			}
			continue
		}

		if inURL {
			writeByte(c)
			if c == ')' {
				inURL = false
			}
			continue
		}

		switch {
		case c == '"' || c == '\'':
			flushWS(&b, &pendingWS)
			quote = c
			writeByte(c)
			continue

		case c == '/' && i+1 < len(s) && s[i+1] == '*':
			end := strings.Index(s[i+2:], "*/")
			if end < 0 {
				return b.String()
			}
			if keepFirst && i+2 < len(s) && s[i+2] == '!' {
				b.WriteString(s[i : i+2+end+2])
				b.WriteByte('\n')
				keepFirst = false
			} else {
				// A comment is equivalent to whitespace: "a/**/b" is two
				// tokens.
				pendingWS = true
			}
			i += 2 + end + 1
			continue

		case c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f':
			pendingWS = true
			continue
		}

		// url( without quotes
		if c == '(' && strings.HasSuffix(strings.ToLower(b.String()), "url") {
			flushWS(&b, &pendingWS)
			writeByte(c)
			// Skip whitespace immediately after "url("; if a quote follows,
			// the normal string branch will pick it up.
			j := i + 1
			for j < len(s) && (s[j] == ' ' || s[j] == '\n' || s[j] == '\t') {
				j++
			}
			if j < len(s) && s[j] != '"' && s[j] != '\'' {
				inURL = true
			}
			i = j - 1
			continue
		}

		// Whitespace is unnecessary next to braces and declaration separators.
		if pendingWS {
			prev := lastByte(&b)
			if !isDrop(prev) && !isDrop(c) {
				writeByte(' ')
			}
			pendingWS = false
		}
		writeByte(c)
	}
	return strings.TrimSpace(b.String()) + "\n"

}

func flushWS(b *strings.Builder, pending *bool) {
	if *pending {
		if p := lastByte(b); p != 0 && !isDrop(p) {
			b.WriteByte(' ')
		}
		*pending = false
	}
}

func lastByte(b *strings.Builder) byte {
	s := b.String()
	if len(s) == 0 {
		return 0
	}
	return s[len(s)-1]
}

// checkBraces checks the balance of curly braces in a kit file.
//
// Importing into a layer nests the file inside an `@layer X { ... }` block,
// so an extra `}` closes the layer itself rather than the rule: the rest of
// the file — and everything subsequently imported into that same layer — ends
// up OUTSIDE the layers. An unlayered rule wins over every layer, so overrides
// from the application, prefers-reduced-motion, @media print, and
// forced-colors stop working for such components.
//
// In src/, the bug is invisible through @import: an extra top-level brace is
// a syntax error, the browser discards it, and the import assigns the layer.
// Comparing dist/ with src/ does not catch it either: the brace is identical
// in both. Therefore the check runs against source before the build.
//
// Strings, comments, and unquoted url(...) are skipped: inside data-URI masks,
// braces are part of the drawing, not structure.
func checkBraces(name string, css []byte) error {
	s := string(css)
	line := 1
	depth := 0
	openLine := 0 // line of the last unclosed "{"

	var quote byte
	inURL := false

	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '\n' {
			line++
		}

		switch {
		case quote != 0:
			if c == '\\' && i+1 < len(s) {
				i++
			} else if c == quote {
				quote = 0
			}
		case inURL:
			if c == ')' {
				inURL = false
			}
		case c == '"' || c == '\'':
			quote = c
		case c == '/' && i+1 < len(s) && s[i+1] == '*':
			end := strings.Index(s[i+2:], "*/")
			if end < 0 {
				return fmt.Errorf("%s:%d: comment is not closed", name, line)
			}
			line += strings.Count(s[i:i+2+end+2], "\n")
			i += 2 + end + 1
		case c == '(' && i >= 3 && strings.EqualFold(s[i-3:i], "url"):
			j := i + 1
			for j < len(s) && (s[j] == ' ' || s[j] == '\t' || s[j] == '\n') {
				j++
			}
			if j < len(s) && s[j] != '"' && s[j] != '\'' {
				inURL = true
			}
		case c == '{':
			if depth == 0 {
				openLine = line
			}
			depth++
		case c == '}':
			depth--
			if depth < 0 {
				return fmt.Errorf(
					"%s:%d: stray \"}\".\n"+
						"The file is nested into an @layer block, so it closes the layer rather\n"+
						"than a rule: the rest of the file ends up outside the layers and starts\n"+
						"winning against application styles, prefers-reduced-motion and @media print.",
					name, line)
			}
		}
	}

	if depth > 0 {
		return fmt.Errorf("%s: unclosed blocks: %d, the first opened on line %d",
			name, depth, openLine)
	}
	if quote != 0 {
		return fmt.Errorf("%s: string is not closed", name)
	}
	return nil
}

// checkPkgVersion checks the version in package.json against VERSION.
//
// Missing package.json is not an error: the kit works as a file without a
// registry, and a repository consumed by link does not have to contain
// package.json. But DIVERGENCE is an error, and a silent one: it is visible
// only after publication.
func checkPkgVersion(root, ver string) error {
	b, err := os.ReadFile(filepath.Join(root, "package.json"))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("cannot read package.json: %w", err)
	}

	var pkg struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(b, &pkg); err != nil {
		return fmt.Errorf("package.json cannot be parsed: %w", err)
	}
	if pkg.Version != ver {
		return fmt.Errorf(
			"versions diverged: VERSION = %s, package.json = %s.\n"+
				"The dist/ header prints the first, while the registry receives the second —\n"+
				"the CDN will contain a file that names itself differently from the package around it.",
			ver, pkg.Version)
	}
	return nil

}

// isDrop — characters next to which whitespace means nothing.
//
// ":" and "," are intentionally NOT included: whitespace changes the selector
// for the first, while it is safe for the second, but the gain is negligible,
// and keeping one rule is simpler.
func isDrop(c byte) bool {
	switch c {
	case '{', '}', ';', 0:
		return true
	}
	return false
}

// rawGap catches direct access to the spacing scale where a role already
// exists.
//
// Kit rule: "components take this only through the role tiers." It relied on
// an honor system and broke once: .inst-stack--loose carried
// gap: var(--space-7), because there was no role for the space between
// sections, while 24px were needed. The role was added (--gap-section), the
// workaround removed — but nothing prevented adding it again.
//
// Only LARGE steps are caught, and only in the open. Small ones (1–6) serve
// micro-spacing, where there is no role and none should be created: spacing
// between an icon and its label is not the same as spacing between sections.
// Padding is not caught at all: it has its own roles, and viewport arithmetic
// such as 100vw - var(--space-8) has nothing to do with the tier.
var rawGap = regexp.MustCompile(`(?:^|[;{[:space:]])(?:row-|column-)?gap:[[:space:]]*var\(--space-(7|8|9|10)\)`)

// checkRoleTier forbids spacing that bypasses the role tier. tokens.css is
// excluded: roles are DECLARED there through scale steps; that is its job.
func checkRoleTier(name string, css []byte) error {
	if name == "tokens.css" {
		return nil
	}
	for i, line := range strings.Split(string(css), "\n") {
		if m := rawGap.FindStringSubmatch(line); m != nil {
			return fmt.Errorf("%s:%d: spacing bypasses role tier — var(--space-%s). "+
				"Between sections use --gap-section, inside use --gap-row or --gap-inline",
				name, i+1, m[1])
		}
	}
	return nil
}

// The "Forbidden" section of the design principles, translated into regular
// expressions.
//
// All five rules are GREEN TODAY — these are locks, not work. The point of the
// gate is not to find something now, but to stop the prohibition from relying
// on memory: each of the five can be violated by one innocuous-looking line
// that produces no browser error.
var (
	banImportant = regexp.MustCompile(`!\s*important`)
	banBold      = regexp.MustCompile(`font-weight:\s*(700|800|900|bold)\b`)
	banColor     = regexp.MustCompile(`#[0-9a-fA-F]{3,8}\b|\brgba?\(|\bhsla?\(|\boklch\(|\bcolor\(`)
	banUtility   = regexp.MustCompile(`\.[mp][trblxy]?-[0-9]`)

	maskDecl = regexp.MustCompile(`\bmask(-image)?\s*:`)
	maskInk  = regexp.MustCompile(`^#(000|000000)$`)
	uriHex   = regexp.MustCompile(`%23[0-9a-fA-F]{3,8}`)
	urlDecl  = regexp.MustCompile(`([\w-]+)\s*:\s*[^;]*url\(`)
)

// Token layer. Here a raw color is WORK, not a violation: tokens.css declares
// ramps and semantics, print.css restores the light theme on paper by
// reassigning the same semantic names. Both are tiers 1–2, and both are
// supposed to name color by number. Everyone else is not.
var colorLayer = map[string]bool{"tokens.css": true, "print.css": true}

// checkBans guards prohibitions that have so far relied on attentiveness.
//
// The check is line-based and therefore does not see a declaration split by a
// line break. That does not matter for four of the five rules — property and
// value are on the same line — and masks were checked separately: all 27 kit
// declarations are single-line.
func checkBans(name string, raw []byte) error {
	css := css.Blank(raw)
	for i, line := range strings.Split(string(css), "\n") {
		at := fmt.Sprintf("%s:%d", name, i+1)

		// !important. Exactly one exception, named explicitly: [hidden] in
		// base.css is correctness, not presentation.
		if banImportant.MatchString(line) {
			if !(name == "base.css" && strings.Contains(line, "[hidden]")) {
				return fmt.Errorf("%s: !important. The only allowed one is [hidden] in base.css. "+
					"Layer order does the same thing without breaking the promise that application styles win", at)
			}
		}

		// Weight 700. base.css closes the only door through which the platform
		// brings it (strong and b are reset to --weight-medium); this closes the
		// door through which a change can bring it.
		if m := banBold.FindStringSubmatch(line); m != nil {
			return fmt.Errorf("%s: weight %s. The kit has two weights — --weight-normal and --weight-medium (600). "+
				"Segoe UI has no real 500, and 700 in an instrumental interface shouts louder than data", at, m[1])
		}

		// Spacing utilities. The scale is intentionally sparse at the top, and
		// a utility set would bring back "a little more" as the first class.
		if m := banUtility.FindString(line); m != "" {
			return fmt.Errorf("%s: spacing utility %q. Rhythm is set by flow primitives "+
				"(.inst-stack, .inst-cluster, .inst-grid) with a gap named by intent", at, m)
		}

		// A raw color outside the token layer is a hardcoded light theme.
		if !colorLayer[name] {
			for _, c := range banColor.FindAllString(line, -1) {
				// A mask is an exception, and it is narrow: the color in it is
				// not color but the alpha channel that cuts the shape. Fill
				// comes from a token, so black is the only meaningful value
				// here, and only it is allowed.
				if maskDecl.MatchString(line) && maskInk.MatchString(c) {
					continue
				}
				return fmt.Errorf("%s: raw color %q. The component uses semantics "+
					"(--text-primary, --surface-raised), not a number: a number is a hardcoded light theme", at, c)
			}
		}

		// Color inside a data-URI. The shape is drawn with a mask, color comes
		// from a token — so the image itself contains no color, only shape ink.
		for _, h := range uriHex.FindAllString(line, -1) {
			if h != "%23000" {
				return fmt.Errorf("%s: color %s inside data-URI. The shape is drawn with a mask and colored by a token; "+
					"only %%23000 is allowed inside the image — the shape's own ink", at, h)
			}
		}

		// url() only in a mask. An image placed as a background colors itself
		// and therefore cannot follow the theme.
		if m := urlDecl.FindStringSubmatch(line); m != nil && !strings.HasPrefix(m[1], "mask") {
			return fmt.Errorf("%s: url() in property %q. An image is allowed only as a mask (mask, mask-image): "+
				"a background colors itself and follows neither the theme nor the tone", at, m[1])
		}
	}
	return nil

}

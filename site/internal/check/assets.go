package check

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"instrument/site/internal/content"
)

// The site's own variables: those that have no place in the kit and never
// will. --c and --v are the data channel of a colour swatch and of an API row.
// --code-* is syntax highlighting: the kit does not highlight code, that is
// the reference's business with its own content.
var ownVars = map[string]bool{
	"--c": true, "--v": true,
	"--code-tag": true, "--code-attr": true, "--code-val": true,
	// --code-line is the listing's line box: the copy button in the corner is
	// positioned against a LINE rather than against the padding edge, and the
	// number has to be stated once for the two rules that need it.
	"--code-line": true,
}

var (
	cssVarUse = regexp.MustCompile(`var\((--[a-z0-9-]+)`)

	bannedProps = regexp.MustCompile(`(?m)^\s*(text-transform)\s*:`)

	hexColor = regexp.MustCompile(`#[0-9a-fA-F]{3,8}\b`)

	// TWO HALVES, ONE RULE. The Russian half is the pattern the gate was
	// written with; the English half exists because `src` and `tools` are now
	// English, and against them the Russian half matches nothing at all. A gate
	// whose phrase list is in a language the corpus no longer speaks is green
	// for the same reason an empty gate is.
	//
	// The English phrases are narrower than a literal translation, and that is
	// deliberate: «раньше» carries the sense of a chronicle by itself, while a
	// bare "previously" is also how a sentence about a browser starts. Each of
	// the six was run over every comment in `src`, `tools` and `site` before
	// being added — zero matches on the corpus as it stands. A gate that
	// complains about what a human did not write is the first one people stop
	// reading.
	pastTense = regexp.MustCompile(`(?i)(^|[.!?;] )раньше|ни разу не|выяснилось|` +
		`здесь сто(ял|яла|яли|яло)|было объявлено|` +
		`(^|[.!?;] )(previously|formerly|originally)\b|` +
		`\bhere (stood|used to stand|used to be)\b|there used to be|` +
		`\bit turned out that\b|\bnot once did\b|\bwas declared here\b`)
)

// Base checks that a token declared in a CELL has a declaration in the base
// as well.
//
// The "From" column of the reference stands on this. It prints the first
// declaration and calls it the base, counting the rest as deviations: "base
// +14: density · scale". If a token is declared only in
// `[data-density="compact"]`, that is what comes first — and the page calls a
// base what is not one, without saying so once. The number will be right and
// will mean something other than what is written beside it, which is exactly
// the illness the column was started for.
//
// Today the invariant holds across all one hundred and eighty-two tokens of
// `tokens.css`, and the check is needed not to establish it but to keep it
// from drifting in silence.
func Base(tokens map[string]content.Token) []string {
	var problems []string
	var names []string
	for n := range tokens {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		t := tokens[n]
		// tokens.css only. The other files have no cells: there a component's
		// variable is declared on the component itself — `--btn-bg` on
		// `.inst-btn` — and that is its home rather than a deviation from a
		// base. Demanding `:root` of it would mean demanding a global name of
		// something deliberately local: the first version of this check did
		// exactly that and produced twenty-three false ones.
		if t.File != "src/tokens.css" || len(t.Cells) == 0 {
			continue
		}
		base := false
		for _, c := range t.Cells {
			if !strings.Contains(c, "@media") &&
				(strings.Contains(c, ":root") || strings.HasPrefix(c, ".inst-theme")) {
				base = true
				break
			}
		}
		if !base {
			problems = append(problems, fmt.Sprintf(
				"%s  the token %s is declared only in a cell (%s): the reference will call a base what is not one",
				t.File, n, t.Cells[0]))
		}
	}
	return problems
}

func Assets(files map[string]string, tokens map[string]content.Token) []string {
	var problems []string

	for name, css := range files {
		code := stripComments(css)

		seen := map[string]bool{}
		for _, m := range cssVarUse.FindAllStringSubmatch(code, -1) {
			v := m[1]
			if seen[v] || ownVars[v] {
				continue
			}
			seen[v] = true
			if _, ok := tokens[v]; !ok {
				problems = append(problems, fmt.Sprintf(
					"%s  there is no token %s in the kit: the declaration is dropped in silence", name, v))
			}
		}

		for _, m := range bannedProps.FindAllStringSubmatch(code, -1) {
			problems = append(problems, fmt.Sprintf(
				"%s  the property %s does not occur in the kit at all — the site has nothing to justify starting it with", name, m[1]))
		}

		for _, m := range hexColor.FindAllString(code, -1) {
			problems = append(problems, fmt.Sprintf(
				"%s  the colour %s bypasses the semantics: one theme is hard-coded", name, m))
		}

	}

	sort.Strings(problems)
	return problems
}

// StrayCommentEnd looks for a `*/` that closed a comment ahead of time.
//
// Comments in CSS do not nest. A line like `--pad-*/--gap-*` inside an
// explanation closes it on the spot: the rest of the text becomes rubbish,
// and the parser, recovering, eats the WHOLE of the NEXT rule. The error is
// silent — the file loads, the styles partly work, and it can be found only
// by the behaviour that is missing.
//
// That is how the kit lost the `neutral` tone (a note with neither icon nor
// fill) and the tap area of the checkbox, the radio and the switch — that is,
// the WCAG 2.2 AA requirement the rule was written for.
func StrayCommentEnd(files map[string]string) []string {
	var problems []string
	for name, src := range files {
		if !strings.HasSuffix(name, ".css") {
			continue
		}
		var code strings.Builder
		rest := src
		for {
			i := strings.Index(rest, "/*")
			if i < 0 {
				code.WriteString(rest)
				break
			}
			code.WriteString(rest[:i])
			j := strings.Index(rest[i+2:], "*/")
			if j < 0 {
				break
			}
			rest = rest[i+2+j+2:]
		}
		if k := strings.Index(code.String(), "*/"); k >= 0 {
			line := strings.Count(code.String()[:k], "\n") + 1
			problems = append(problems, fmt.Sprintf(
				"%s:%d  a comment closed ahead of time: a `*/` inside an explanation eats the next rule",
				name, line))
		}
	}
	sort.Strings(problems)
	return problems
}

func Comments(files map[string]string) []string {
	var problems []string
	for name, src := range files {
		text := strings.Join(strings.Fields(comments(src)), " ")
		seen := map[string]bool{}
		for _, m := range pastTense.FindAllString(text, -1) {
			key := strings.ToLower(m)
			if seen[key] {
				continue
			}
			seen[key] = true
			problems = append(problems, fmt.Sprintf(
				"%s  the comment tells a story (\"%s\"): what is wanted is why it is so now", name, m))
		}
	}
	sort.Strings(problems)
	return problems
}

// comments takes out of a source what a human wrote to a HUMAN.
//
// String literals are skipped, and that is not pedantry. Parsing by the first
// `//` on a line counts as a comment everything standing behind the quotes —
// and the mutation stand keeps the forbidden phrases as DATA: for it they are
// the text of a mutation rather than an explanation beside a function. Without
// that distinction the stand brings down the site build with its own content,
// and deservedly so: to the parser it is indistinguishable from an offender.
//
// The parsing is line by line: a multi-line raw Go literal is beyond it. The
// cost of that error is one-sided — a missed comment, never an invented one —
// and that is the right side: a gate that complains about what a human did not
// write is the first one people stop reading.
func comments(src string) string {
	var b strings.Builder
	for _, line := range strings.Split(src, "\n") {
		if i := lineComment(line); i >= 0 {
			b.WriteString(line[i+2:])
			b.WriteByte('\n')
		}
	}
	for {
		i := blockOpen(src)
		if i < 0 {
			return b.String()
		}
		src = src[i+2:]
		j := strings.Index(src, "*/")
		if j < 0 {
			b.WriteString(src)
			return b.String()
		}
		b.WriteString(src[:j])
		b.WriteByte('\n')
		src = src[j+2:]
	}
}

// lineComment is the position of the `//` that starts a comment, or -1.
func lineComment(line string) int {
	return outsideQuotes(line, "//", false)
}

// blockOpen is the position of the `/*` that starts a comment, or -1.
func blockOpen(src string) int {
	return outsideQuotes(src, "/*", true)
}

// outsideQuotes looks for the first occurrence of want outside quotes. A quote
// left unclosed by the end of the line does not count as one: it is an
// apostrophe in prose rather than a literal.
func outsideQuotes(src, want string, multiline bool) int {
	var quote byte
	for i := 0; i < len(src); i++ {
		c := src[i]
		if c == '\n' {
			if !multiline {
				break
			}
			quote = 0
			continue
		}
		if quote != 0 {
			if c == '\\' {
				i++
				continue
			}
			if c == quote {
				quote = 0
			}
			continue
		}
		if c == '"' || c == '\'' || c == '`' {
			// A lone quote with no pair before the line ends is an apostrophe.
			if !closes(src[i+1:], c) {
				continue
			}
			quote = c
			continue
		}
		if strings.HasPrefix(src[i:], want) {
			return i
		}
	}
	return -1
}

// closes answers whether the quote occurs again before the line ends.
func closes(rest string, q byte) bool {
	for i := 0; i < len(rest); i++ {
		if rest[i] == '\n' {
			return false
		}
		if rest[i] == '\\' {
			i++
			continue
		}
		if rest[i] == q {
			return true
		}
	}
	return false
}

func stripComments(css string) string {
	var b strings.Builder
	b.Grow(len(css))
	for {
		i := strings.Index(css, "/*")
		if i < 0 {
			break
		}
		j := strings.Index(css[i:], "*/")
		if j < 0 {
			css = css[:i]
			break
		}
		b.WriteString(css[:i])
		css = css[i+j+2:]
	}
	b.WriteString(css)

	return dataURI.ReplaceAllString(b.String(), `url("")`)
}

var dataURI = regexp.MustCompile(`url\("data:[^"]*"\)`)

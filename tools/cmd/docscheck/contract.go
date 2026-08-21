package main

import (
	"fmt"
	"regexp"
	"strings"
)

// A page PROMISES things about markup, and the promise is prose. This check
// asks whether the page keeps it in its own examples.
//
// WHY IT EXISTS. The contract table is the densest part of a component page —
// 294 rows over eighty pages, and after the rules section went away it is the
// only place a requirement is written down at all. Nothing read it. A row
// saying an attribute is REQUIRED sat next to live examples that did not carry
// it, and both went to press.
//
// A measurement over the sixty-one unconditional ARIA promises found exactly
// one broken, and it was broken in the direction nobody looks for: code.md
// demanded the application add `aria-live="polite"` for the copy
// confirmation, while kit.js announces it itself through a shared
// [data-inst-live] region. Not a missing attribute — an INVENTED OBLIGATION,
// work charged to the reader that the kit had already done. A gate that only
// looked for missing attributes would never have found it; this one did,
// because it asks a blunter question: is the thing you demand anywhere in the
// examples you show?
//
// WHAT IT DOES NOT CHECK, and why that is not laziness. Two thirds of the
// contract rows name an element, a class, or nothing in code at all, and ten
// of the ARIA rows are CONDITIONAL — "yes, if the banner appeared in response
// to an action", "yes, if the panel is a landmark". A static example does not
// meet such a condition and must not be failed for it. The condition is
// recognised by the word "если" / "if" in either cell; a promise that carries
// no condition is the only kind held to the examples.
//
// A page in Russian and the same page in English will exist side by side
// through step 4 of the language move, so both spellings of the header, of
// "yes" and of "if" are matched. The Russian here is a PATTERN, not text —
// the same standing as the Russian in cmd/registry and cmd/docscheck/tokens.go.
var (
	contractHead = regexp.MustCompile(`(?m)^## (?:Контракт|Contract)\s*$`)
	nextH2       = regexp.MustCompile(`(?m)^## `)
	promiseTable = regexp.MustCompile(
		`(?m)^\| ?(?:Что|What) ?\| ?(?:Обязательно|Required) ?\| ?(?:Почему|Why) ?\|\s*$\r?\n\|[-\s|:]+\|\r?\n((?:\|.*\r?\n)+)`)
	demoBlock    = regexp.MustCompile("(?s)```html preview.*?\n(.*?)```")
	inlineCodeRe = regexp.MustCompile("`([^`]+)`")
	ariaName     = regexp.MustCompile(`^(role|aria-[a-z]+)`)
	// Not \b: Go word boundaries are ASCII-only, and «да» ends in a letter
	// the engine does not consider one — the boundary never matches and the whole
	// check silently held nothing. The same trap ate a Cyrillic phrase count
	// earlier in the audit; \p{L} is the fix in both languages.
	saysYes   = regexp.MustCompile(`(?i)^(да|yes)([^\p{L}]|$)`)
	condition = regexp.MustCompile(`(?i)(^|[^\p{L}])(если|if)([^\p{L}]|$)`)
)

// promisesKept returns one problem per unconditional promise the page's own
// examples do not keep, and the number of promises it managed to check.
func promisesKept(at func(int) string, body string) (problems []string, checked int) {
	h := contractHead.FindStringIndex(body)
	if h == nil {
		return nil, 0
	}
	rest := body[h[1]:]
	if n := nextH2.FindStringIndex(rest); n != nil {
		rest = rest[:n[0]]
	}

	var demos strings.Builder
	for _, d := range demoBlock.FindAllStringSubmatch(body, -1) {
		demos.WriteString(d[1])
	}
	if demos.Len() == 0 {
		// A page with no live example promises nothing it could break here.
		// Kept deliberately rather than reported: index pages and the prose of
		// the foundations legitimately have no frame.
		return nil, 0
	}
	shown := demos.String()

	for _, tbl := range promiseTable.FindAllStringSubmatch(rest, -1) {
		for _, line := range strings.Split(strings.TrimSpace(tbl[1]), "\n") {
			cells := strings.Split(strings.Trim(strings.TrimSpace(line), "|"), "|")
			if len(cells) != 3 {
				continue
			}
			what := strings.TrimSpace(cells[0])
			must := strings.TrimSpace(cells[1])
			if !saysYes.MatchString(must) {
				continue
			}
			if condition.MatchString(what) || condition.MatchString(must) {
				continue
			}
			for _, tok := range inlineCodeRe.FindAllStringSubmatch(what, -1) {
				name := ariaName.FindString(tok[1])
				if name == "" {
					continue
				}
				// The VALUE counts when the promise names one. Searching for `role`
				// alone passed a page promising role="listbox" whose examples carried
				// only role="option" — the substring was there and meant something
				// else. Caught by the stand: the mutation that strips the listbox role
				// off the chip strip went through untouched.
				want := name
				if i := strings.IndexByte(tok[1], 61); i >= 0 {
					want = strings.TrimSpace(tok[1])
				}
				checked++
				if !strings.Contains(shown, want) {
					// The line number is the contract heading rather than the row: a
					// table row moves, and a report pointing at a line that has since
					// become something else is worse than one pointing at the section.
					problems = append(problems, fmt.Sprintf(
						"%s  promised without a condition and shown in no example: %s",
						at(countLines(body[:h[0]])+1), what))
				}
				break // one attribute per row is enough to hold it to account
			}
		}
	}
	return problems, checked
}

func countLines(s string) int { return strings.Count(s, "\n") }

package check

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"instrument/site/internal/content"
	"instrument/site/internal/i18n"
)

var (
	hrefRe = regexp.MustCompile(`href="(/[^"#?]*)"`)
	// The same link WITH its fragment. Two expressions rather than one with an
	// optional group: the first is asked "does the page exist", the second
	// "does the heading exist", and a link with no fragment has nothing to
	// answer to the second.
	fragRe   = regexp.MustCompile(`href="(/[^"#?]*)#([^"]+)"`)
	idRe2    = regexp.MustCompile(`\sid="([^"]+)"`)
	leftMdRe = regexp.MustCompile(`href="([^"]*\.md(?:#[^"]*)?)"`)

	useRe = regexp.MustCompile(`<use\s[^>]*href="#([^"]*)"`)

	nodeRe = regexp.MustCompile(`(?m)^<([a-z][\w-]*)([^>]*)>`)

	hiddenAttrRe = regexp.MustCompile(`\s(popover|hidden)(\s|=|>|$)`)

	cyrillicRe = regexp.MustCompile(`[\p{Cyrillic}]+`)
)

var hidden = map[string]bool{
	"dialog": true, "template": true, "script": true, "style": true,
}

// A LINK INTO A HEADING IS CHECKED TOO, and it is the half that was missing.
// A `#fragment` that names nothing lands the reader at the top of the page
// instead of at the paragraph they were sent to — no error, no empty page, and
// nothing that looks wrong to whoever wrote the link.
//
// The move to English broke exactly this: a fragment naming a Russian heading
// on a page whose headings had been English for four commits. The rule now
// stands where every anchor in the corpus is visible at once.
func Verify(pages []*content.Page, sprite string) []string {
	routes := map[string]bool{}
	anchors := map[string]map[string]bool{}
	for _, p := range pages {
		routes[p.Route] = true
		ids := map[string]bool{}
		for _, m := range idRe2.FindAllStringSubmatch(p.HTML, -1) {
			ids[m[1]] = true
		}
		anchors[p.Route] = ids
	}

	var problems []string
	for _, p := range pages {
		for _, m := range leftMdRe.FindAllStringSubmatch(p.HTML, -1) {
			problems = append(problems,
				fmt.Sprintf("%s  the link stayed a file link: %s", p.Route, m[1]))
		}
		for _, m := range hrefRe.FindAllStringSubmatch(p.HTML, -1) {
			t := m[1]
			if strings.HasPrefix(t, "/kit/") || strings.HasPrefix(t, "/assets/") {
				continue
			}
			if !strings.HasSuffix(t, "/") {
				t += "/"
			}
			if !routes[t] {
				problems = append(problems,
					fmt.Sprintf("%s  a link to nowhere: %s", p.Route, m[1]))
			}
		}
		for _, m := range fragRe.FindAllStringSubmatch(p.HTML, -1) {
			t := m[1]
			if !strings.HasSuffix(t, "/") {
				t += "/"
			}
			// A link to nowhere has already been reported by the loop above;
			// naming it twice would say the same thing about one line.
			if !routes[t] {
				continue
			}
			if !anchors[t][m[2]] {
				problems = append(problems,
					fmt.Sprintf("%s  a link into a heading that is not there: %s#%s",
						p.Route, m[1], m[2]))
			}
		}
		if strings.Contains(p.HTML, "```html preview") {
			problems = append(problems, p.Route+"  a preview fence did not unfold")
		}

		// A page of the other language must carry no Cyrillic AT ALL, and
		// the check earns its place: the body of a page is assembled from
		// three sources, and two of them handed out Russian without a word.
		// The Related row resolved a neighbour through an index that had no
		// language in the key, and the kind of a generated token row was a
		// literal. Both came out as Russian inside an English page, and
		// neither was visible to any other gate.
		if p.Lang != i18n.RU {
			for _, m := range cyrillicRe.FindAllString(p.HTML, -1) {
				problems = append(problems,
					fmt.Sprintf("%s  Russian in a page of another language: %s", p.Route, m))
			}
		}

		problems = append(problems, dupIDs(p)...)

		for _, d := range p.Demos {
			problems = append(problems, verifyDemo(p, d, sprite)...)
		}
	}
	sort.Strings(problems)
	return problems
}

func dupIDs(p *content.Page) []string {
	seen := map[string]bool{}
	dup := map[string]bool{}
	for _, d := range p.Demos {
		for _, m := range idRe.FindAllStringSubmatch(d.Markup, -1) {
			if seen[m[1]] {
				dup[m[1]] = true
			}
			seen[m[1]] = true
		}
	}
	var out []string
	for id := range dup {
		out = append(out, fmt.Sprintf(
			"%s  id=%q is declared in two examples: on one page they now stand side by side", p.Route, id))
	}
	sort.Strings(out)
	return out
}

var idRe = regexp.MustCompile(`\bid="([^"]+)"`)

func verifyDemo(p *content.Page, d content.Demo, sprite string) []string {
	var out []string

	for _, m := range useRe.FindAllStringSubmatch(d.Markup, -1) {
		id := m[1]
		if id == "" {
			out = append(out, fmt.Sprintf("%s  an empty reference to a symbol: <use href=\"#\">", p.Route))
			continue
		}
		if !strings.Contains(sprite, `id="`+id+`"`) {
			out = append(out, fmt.Sprintf("%s  no such symbol in the sprite: #%s", p.Route, id))
		}
	}

	nodes := nodeRe.FindAllStringSubmatch(d.Markup, -1)
	if len(nodes) == 0 {
		return out
	}
	visible := 0
	for _, n := range nodes {
		tag, attrs := n[1], n[2]
		if hidden[tag] {
			continue
		}
		if hiddenAttrRe.MatchString(attrs + ">") {
			continue
		}
		visible++
	}
	if visible == 0 {
		out = append(out, fmt.Sprintf(
			"%s  an empty frame: the whole content of the example is hidden by default (%s)",
			p.Route, d.ID))
	}
	return out
}

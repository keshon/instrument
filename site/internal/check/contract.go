package check

import (
	"fmt"
	"sort"
	"strings"

	"instrument/site/internal/content"
	"instrument/site/internal/i18n"
)

// What a page OWES.
//
// There was a second entry here while the pages migrated: a page declared its
// shape and was judged by it, because judging a page written under the old
// contract by the new one would fail it for something nobody asked of it. The
// last page took that entry with it, and the collapse is the mark of
// completion for the step — the same way the `src` zone closed in cmd/lang.
//
// Shape 2 asks for little because the prose left the component page altogether
// (docs/internal/DOCS-SHAPE.md): the comparison of neighbours moved to the
// section index, accessibility merged into the markup contract, and the list of
// neighbours is derived from the registry rather than written.
var required = []string{"contract", "api"}

var gone = map[string]string{
	"usage": "the markup moved into the demo frame and the code panel",
	"a11y":  "merged into the markup contract",
	// A page cannot migrate before its section has an index, and that is the
	// point rather than a snag: "take another one instead" is a statement about
	// NEIGHBOURS, and repeated on each of four pages it drifts from them by one
	// line per edit. The gate makes the destination exist first.
	"when": "moved to the section index, where the comparison of neighbours stands in one place",
}

func underContract(dir string) bool {
	for _, p := range []string{"components/", "foundations", "layout", "agent"} {
		if strings.HasPrefix(dir, p) {
			return true
		}
	}
	return false
}

func Contract(pages []*content.Page) (errs, warns []string) {
	for _, p := range pages {
		if p.Lang != i18n.RU || p.Splash || p.Slug == "index" {
			continue
		}
		strict := p.Layout == "component" || p.Layout == "foundation"

		if !strict && !underContract(p.Dir) {
			continue
		}

		var out []string
		for _, s := range p.Sections {
			if !s.Known {
				out = append(out, fmt.Sprintf("%s  the section %q is not in the vocabulary (make it an H3 inside a canonical section)",
					p.Route, s.Title))
			}
		}

		prev, prevID := 0, ""
		for _, s := range p.Sections {
			if !s.Known {
				continue
			}
			if s.Order < prev {
				out = append(out, fmt.Sprintf("%s  the section %q stands after %q and belongs before it",
					p.Route, s.ID, prevID))
			}
			prev, prevID = s.Order, s.ID
		}

		seen := map[string]bool{}
		for _, s := range p.Sections {
			if s.Known && seen[s.ID] {
				out = append(out, fmt.Sprintf("%s  the section %q occurs twice", p.Route, s.ID))
			}
			seen[s.ID] = true
		}

		for _, id := range required {
			if !seen[id] {
				out = append(out, fmt.Sprintf("%s  the required section %q is missing", p.Route, id))
			}
		}
		for _, s := range p.Sections {
			if why, gone := gone[s.ID]; gone {
				out = append(out, fmt.Sprintf(
					"%s  the section %q no longer lives on a page: %s", p.Route, s.ID, why))
			}
		}
		if !p.Hero {
			out = append(out, p.Route+"  no lead example: a live example has to stand in the header, before the first section")
		}

		if p.JS != "" || p.JSOpt != "" {
			if !p.HasJS {
				out = append(out, p.Route+"  js is declared, and there is not one example on JS")
			}
			if !seen["js"] {
				out = append(out, p.Route+"  js is declared, and there is no 'JS' section")
			}
		}

		if strict {
			errs = append(errs, out...)
		} else {
			warns = append(warns, out...)
		}
	}
	sort.Strings(errs)
	sort.Strings(warns)
	return errs, warns
}

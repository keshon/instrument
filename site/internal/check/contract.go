package check

import (
	"fmt"
	"sort"
	"strings"

	"instrument/site/internal/content"
	"instrument/site/internal/i18n"
)

// What a page OWES, per shape of the contract.
//
// Two dictionaries live side by side on purpose, and the reason is the same one
// the migration zones in cmd/lang are built on: switching the single dictionary
// would turn all fifty pages red in one commit, and a gate raised red lives
// with its check disabled. A page declares its shape and is judged by it; the
// last page to switch takes the old entry with it, and that collapse is the
// mark of completion for the step.
//
// Shape 2 asks for less because the prose leaves the component page altogether
// (docs/internal/DOCS-SHAPE.md): «Правила» and «Когда использовать» move to
// foundations and to the section index, «Доступность» merges into the markup
// contract, and «Связанное» is derived from the registry rather than written.
var requiredByShape = map[int][]string{
	1: {"usage", "a11y", "api", "related"},
	2: {"contract", "api"},
}

// Sections that shape 2 does not merely stop demanding but forbids: while they
// are still allowed, a half-migrated page passes both dictionaries and nobody
// can tell which one it follows.
//
// The list holds only what has somewhere to go TODAY. Three more sections are
// meant to leave and cannot yet, and saying so here is cheaper than discovering
// it while migrating a page:
//
//	related  — meant to be derived from the registry, but a registry entry
//	           carries one field, `axes`. There are no relations in it to
//	           derive from; that line of DOCS-SHAPE.md was written on an
//	           assumption rather than on a reading.
//	rules    — belongs in foundations, and which foundation page takes which
//	           rule is not decided per component.
//
// Forbidding a section whose destination does not exist would force the prose
// to be deleted rather than moved, and deleted reasoning does not come back.
var goneInShape2 = map[string]string{
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
				out = append(out, fmt.Sprintf("%s  раздел %q не из словаря (сделайте его H3 внутри канонического раздела)",
					p.Route, s.Title))
			}
		}

		prev, prevID := 0, ""
		for _, s := range p.Sections {
			if !s.Known {
				continue
			}
			if s.Order < prev {
				out = append(out, fmt.Sprintf("%s  раздел %q стоит после %q, а должен перед ним",
					p.Route, s.ID, prevID))
			}
			prev, prevID = s.Order, s.ID
		}

		seen := map[string]bool{}
		for _, s := range p.Sections {
			if s.Known && seen[s.ID] {
				out = append(out, fmt.Sprintf("%s  раздел %q встречается дважды", p.Route, s.ID))
			}
			seen[s.ID] = true
		}

		for _, id := range requiredByShape[p.Shape] {
			if !seen[id] {
				out = append(out, fmt.Sprintf("%s  нет обязательного раздела %q (форма %d)",
					p.Route, id, p.Shape))
			}
		}
		if p.Shape >= 2 {
			for _, s := range p.Sections {
				if why, gone := goneInShape2[s.ID]; gone {
					out = append(out, fmt.Sprintf(
						"%s  раздел %q не живёт в форме 2: %s", p.Route, s.ID, why))
				}
			}
		}
		if !p.Hero {
			out = append(out, p.Route+"  нет главного примера: живой пример обязан стоять в шапке, до первого раздела")
		}

		if p.JS != "" || p.JSOpt != "" {
			if !p.HasJS {
				out = append(out, p.Route+"  объявлен js, но нет ни одного примера на JS")
			}
			if !seen["js"] {
				out = append(out, p.Route+"  объявлен js, но нет раздела «JS»")
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

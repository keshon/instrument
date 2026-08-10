package check

import (
	"fmt"
	"sort"
	"strings"

	"instrument/site/internal/content"
	"instrument/site/internal/i18n"
)

var required = []string{"usage", "a11y", "api", "related"}

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

		for _, id := range required {
			if !seen[id] {
				out = append(out, fmt.Sprintf("%s  нет обязательного раздела %q", p.Route, id))
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

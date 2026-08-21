package nav

import (
	"sort"

	"instrument/site/internal/content"
	"instrument/site/internal/i18n"
)

type Section struct {
	Label string
	Items []*content.Page
}

var sections = []struct {
	dir   string
	order []string
}{
	{"start", []string{"install"}},
	{"foundations", []string{
		"colors", "typography", "spacing", "elevation", "motion",
		"density", "scale", "icons", "behavior", "utilities", "tokens",
	}},
	{"layout", []string{
		"index",
		"shell", "container", "flow", "split", "page-header", "section",
	}},
	{"components/actions", []string{"index", "button", "button-group", "segmented", "chip"}},
	{"components/inputs", []string{
		"index",
		"input", "select", "cascader", "toggles", "slider", "num-field", "search",
		"choice-card", "file", "inserts", "form",
	}},
	{"components/display", []string{
		"index",
		"panel", "card", "table", "kv", "metric", "badge", "tag",
		"avatar", "timeline", "calendar", "code",
	}},
	{"components/charts", []string{
		"index",
		"meter", "ring", "sparkline", "legend", "palette",
	}},
	{"components/navigation", []string{
		"index",
		"nav", "tabs", "breadcrumbs", "pagination", "steps", "toolbar",
	}},
	{"components/overlays", []string{
		"index",
		"popover", "menu", "tooltip", "dialog", "sheet",
	}},
	{"components/feedback", []string{
		"index",
		"banner", "toast", "note", "empty", "skeleton", "spinner", "states", "accordion",
	}},
	{"agent", []string{
		"index",
		"run", "task", "step", "approval", "failure", "diff", "output",
		"log", "lane", "history", "budget", "tree",
	}},

	{"blocks", []string{"dashboard", "inspector", "settings-screen"}},

	{"about", nil},
}

func Build(lang i18n.Lang, pages []*content.Page) []Section {
	byDir := map[string][]*content.Page{}
	for _, p := range pages {
		if p.Slug == "index" && p.Dir == "" {
			continue // the front page lives in the header, not in the list
		}
		byDir[p.Dir] = append(byDir[p.Dir], p)
	}

	var out []Section
	for _, s := range sections {
		items := byDir[s.dir]
		if len(items) == 0 {
			continue
		}
		// The index of a section always comes first, whether the order names
		// it or not: it is the page ABOVE the rest, and its place matches
		// that. Putting it into every list by hand would start a rule that
		// gets forgotten on the first new section.
		rank := map[string]int{"index": -1}
		for i, slug := range s.order {
			if slug != "index" {
				rank[slug] = i
			}
		}
		sort.Slice(items, func(i, j int) bool {
			ri, oki := rank[items[i].Slug]
			rj, okj := rank[items[j].Slug]
			switch {
			case oki && okj:
				return ri < rj
			case oki:
				return true
			case okj:
				return false
			default:
				return items[i].Title < items[j].Title
			}
		})
		out = append(out, Section{Label: i18n.Section(lang, s.dir), Items: items})
		delete(byDir, s.dir)
	}

	var rest []string
	for dir := range byDir {
		if dir != "" {
			rest = append(rest, dir)
		}
	}
	sort.Strings(rest)
	for _, dir := range rest {
		items := byDir[dir]
		sort.Slice(items, func(i, j int) bool { return items[i].Title < items[j].Title })
		out = append(out, Section{Label: i18n.Section(lang, dir), Items: items})
	}
	return out
}

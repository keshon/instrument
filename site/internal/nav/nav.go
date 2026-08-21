package nav

import (
	"sort"
	"strings"

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
		"shell", "rail", "statusbar",
		"container", "split", "flow", "page-header", "section",
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
		"avatar", "change", "timeline", "calendar", "code",
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
		"run", "task", "step", "approval", "failure", "diff",
		"log", "lane", "history", "budget", "tree",
	}},

	{"blocks", []string{"console"}},

	{"about", nil},
}

// Check reports an order that names a page which does not exist.
//
// All four kinds of drift the list collected were silent. `agent/output` was a
// page that left; `blocks` named three assemblies that were never written; and
// `rail`, `statusbar` and `change` existed while going unnamed, so they fell
// through to the alphabetical tail — the side column put the status bar
// between the page header and the section, and nobody could say why.
//
// A name with no page is an ERROR and a page with no name is not: the second
// is how a page is added, and demanding both at once would mean a new page
// could not be written without editing this file first. What is guarded is the
// direction that cannot be right.
func Check(pages []*content.Page) []string {
	have := map[string]bool{}
	for _, p := range pages {
		have[strings.TrimPrefix(p.Dir+"/"+p.Slug, "/")] = true
	}
	var out []string
	for _, s := range sections {
		for _, slug := range s.order {
			if key := s.dir + "/" + slug; !have[key] {
				out = append(out, "the order of the side column names "+key+
					", and there is no such page")
			}
		}
	}
	sort.Strings(out)
	return out
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

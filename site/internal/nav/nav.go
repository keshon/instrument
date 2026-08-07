package nav

import (
	"sort"

	"instrument/site/internal/content"
)

type Section struct {
	Label string
	Items []*content.Page
}

// Порядок разделов — по пути читателя, а не по алфавиту: сначала из чего
// собран кит, потом чем собирают экран, потом компоненты по работе, и
// агентный слой отдельным разделом, потому что он и есть причина брать
// именно этот кит.
var sections = []struct {
	dir, label string
	order      []string // явный порядок внутри; остальное по заголовку
}{
	{"foundations", "Основания", []string{
		"colors", "typography", "spacing", "elevation", "motion",
		"density", "icons", "utilities", "tokens",
	}},
	{"layout", "Раскладка", []string{
		"shell", "container", "flow", "split", "page-header", "section",
	}},
	{"components/actions", "Действия", nil},
	{"components/inputs", "Ввод", []string{
		"input", "select", "toggles", "slider", "num-field", "search",
		"choice-card", "file", "form",
	}},
	{"components/display", "Отображение данных", []string{
		"panel", "card", "table", "kv", "metric", "badge", "tag",
		"avatar", "timeline", "calendar", "code",
	}},
	{"components/charts", "Графики", []string{
		"meter", "ring", "sparkline", "legend", "palette",
	}},
	{"components/navigation", "Навигация", []string{
		"nav", "tabs", "breadcrumbs", "pagination", "steps", "toolbar",
	}},
	{"components/overlays", "Оверлеи", []string{
		"popover", "menu", "tooltip", "dialog", "sheet",
	}},
	{"components/feedback", "Обратная связь", []string{
		"banner", "note", "empty", "skeleton", "spinner", "states", "accordion",
	}},
	{"agent", "Агентный слой", []string{
		"run", "task", "step", "approval", "failure", "diff", "output",
		"log", "lane", "budget", "tree",
	}},
	{"about", "О проекте", nil},
}

func Build(pages []*content.Page) []Section {
	byDir := map[string][]*content.Page{}
	for _, p := range pages {
		if p.Slug == "index" && p.Dir == "" {
			continue // главная живёт в шапке, а не в списке
		}
		byDir[p.Dir] = append(byDir[p.Dir], p)
	}

	var out []Section
	for _, s := range sections {
		items := byDir[s.dir]
		if len(items) == 0 {
			continue
		}
		rank := map[string]int{}
		for i, slug := range s.order {
			rank[slug] = i
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
		out = append(out, Section{Label: s.label, Items: items})
		delete(byDir, s.dir)
	}

	// Каталог, не попавший в список выше, всё равно показывается: молча
	// потерянная страница хуже некрасивого раздела.
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
		out = append(out, Section{Label: dir, Items: items})
	}
	return out
}

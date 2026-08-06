package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// verify проверяет собранный сайт до того, как его увидит человек.
//
// docs-check сверяет ИСХОДНИКИ: что классы существуют, ссылки ведут в файлы,
// токены не выдуманы. Но между исходником и страницей стоит генератор, и он
// ломает своё — причём так, что страница остаётся внешне целой:
//
//   · ссылка .md не переписалась в маршрут — перелинковка ведёт в 404;
//   · ограда ```html preview не развернулась — код без демонстрации;
//   · маршрут в ссылке есть, а страницы по нему нет — опечатка в пути.
//
// Ни одно из трёх не видно при беглом просмотре, поэтому проверка идёт по
// файлам, а не по впечатлению, и роняет сборку.

var (
	hrefRe   = regexp.MustCompile(`href="(/[^"#?]*)"`)
	leftMdRe = regexp.MustCompile(`href="([^"]*\.md(?:#[^"]*)?)"`)
)

func verify(pages []*Page) []string {
	routes := map[string]bool{}
	for _, p := range pages {
		routes[p.Route] = true
	}

	var problems []string
	for _, p := range pages {
		for _, m := range leftMdRe.FindAllStringSubmatch(p.HTML, -1) {
			problems = append(problems,
				fmt.Sprintf("%s  ссылка осталась файловой: %s", p.Route, m[1]))
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
					fmt.Sprintf("%s  ссылка в никуда: %s", p.Route, m[1]))
			}
		}
		// Ограда html preview обязана была развернуться в сцену.
		if strings.Contains(p.HTML, "```html preview") {
			problems = append(problems, p.Route+"  ограда preview не развернулась")
		}
	}
	sort.Strings(problems)
	return problems
}

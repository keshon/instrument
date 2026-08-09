package check

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"instrument/site/internal/content"
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

	// Ссылка на символ спрайта. Пустая (`#`) и выдуманная одинаково молчаливы:
	// <svg> рисует ничего, кнопка выходит пустым квадратом — и сборку это
	// проходит.
	useRe = regexp.MustCompile(`<use\s[^>]*href="#([^"]*)"`)

	// Верхнеуровневые узлы разметки примера.
	nodeRe = regexp.MustCompile(`(?m)^<([a-z][\w-]*)([^>]*)>`)

	// Атрибут целиком, а не подстрока. Наивный поиск " popover" находил его
	// внутри popovertarget у КНОПКИ, вызывающей поповер, и объявлял пустыми
	// четыре живых примера — ложное срабатывание, то есть худший вид ошибки
	// в проверке: пропуск ищут, ложному верят.
	hiddenAttrRe = regexp.MustCompile(`\s(popover|hidden)(\s|=|>|$)`)
)

// hidden — элементы, которые сами по себе не рисуют ничего.
//
// Пример, у которого ВСЕ верхнеуровневые узлы такие, — пустой кадр. Ровно так
// выглядели примеры модалки и шторки: закрытый <dialog> это display: none, а
// кнопки-триггера в разметке не было. Читатель видел белый прямоугольник и
// решал, что сломан кит.
var hidden = map[string]bool{
	"dialog": true, "template": true, "script": true, "style": true,
}

// Verify проверяет собранные страницы и разметку примеров.
//
// sprite — содержимое assets/sprite.svg: ссылки на символы проверяются по
// нему, а не по списку в голове.
func Verify(pages []*content.Page, sprite string) []string {
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

		problems = append(problems, dupIDs(p)...)

		for _, d := range p.Demos {
			problems = append(problems, verifyDemo(p, d, sprite)...)
		}
	}
	sort.Strings(problems)
	return problems
}

// dupIDs ловит повторяющийся id в пределах страницы.
//
// Проверка появилась вместе с отказом от iframe. Пока каждый пример жил в
// своём документе, одинаковые `id` в двух примерах были безобидны — документы
// разные. В общем потоке страницы это уже поломка, и поломка тихая:
// `aria-labelledby` находит ПЕРВЫЙ узел с таким именем, поэтому второй пример
// начинает объявляться подписью первого. Ровно так и было у согласования: два
// блока с `id="ap1"`, и второй представлялся чужим заголовком.
//
// Ищутся только id, объявленные в разметке примеров: id самой страницы —
// `main`, `sidebar`, разделы контракта — генерируются, и повтор среди них уже
// невозможен по построению.
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
			"%s  id=%q объявлен в двух примерах: на одной странице они теперь рядом", p.Route, id))
	}
	sort.Strings(out)
	return out
}

var idRe = regexp.MustCompile(`\bid="([^"]+)"`)

// verifyDemo проверяет разметку одного стола примера.
//
// Обе проверки существуют по факту найденной ошибки, и обе невидимы при
// беглом просмотре: страница остаётся целой, а кадр внутри неё — пустым или с
// пустыми квадратами вместо кнопок.
func verifyDemo(p *content.Page, d content.Demo, sprite string) []string {
	var out []string

	for _, m := range useRe.FindAllStringSubmatch(d.Markup, -1) {
		id := m[1]
		if id == "" {
			out = append(out, fmt.Sprintf("%s  пустая ссылка на символ: <use href=\"#\">", p.Route))
			continue
		}
		if !strings.Contains(sprite, `id="`+id+`"`) {
			out = append(out, fmt.Sprintf("%s  нет символа в спрайте: #%s", p.Route, id))
		}
	}

	// Верхнеуровневые узлы примера. Разметка отформатирована с нулевым
	// отступом у корней, поэтому «строка, начинающаяся с тега» и есть корень.
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
			"%s  пустой кадр: всё содержимое примера скрыто по умолчанию (%s)",
			p.Route, d.ID))
	}
	return out
}

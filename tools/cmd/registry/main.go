// Команда registry сверяет реестр компонентов с китом.
//
// ЗАЧЕМ ОН ЕСТЬ. Сквозные оси — печать, высокий контраст, поток — записаны в
// CSS перечислением имён компонентов: восемь списков в трёх файлах. Каждый
// список надо править руками при добавлении компонента, забытая запись не даёт
// ошибки в браузере, и ни одна из пяти прежних проверок этого не видела.
// Комментарий в layout.css описывает, как один такой список уже разошёлся, а
// дыру нашли в приложении, не в ките.
//
// ЧТО ЛЕЖИТ В РЕЕСТРЕ. Только то, чего нельзя вывести из исходников. Какие
// классы существуют, в каком они файле, какие токены читают, есть ли страница —
// всё это добывается обходом и в реестре отсутствует: продублированный факт
// расходится с оригиналом молча. Остаётся решение автора: печатается ли
// компонент, рвётся ли между листами, что с ним делает forced-colors, меряется
// ли он содержимым.
//
// ЧЕГО В РЕЕСТРЕ НЕТ И НЕ БУДЕТ. Замедление при prefers-reduced-motion
// ВЫВОДИТСЯ: замедлять надо ровно то, у чего бесконечная анимация. Проверка
// ниже сверяет список motion.css с этим выводом и не спрашивает реестр вовсе.
// Ось, которую можно вывести, не заводится — иначе реестр становится восьмым
// списком, который надо помнить.
//
// ОТОБРАЖЕНИЕ «КЛАСС → КОМПОНЕНТ» тоже не хранится: оно уже объявлено в
// frontmatter страниц, и docscheck сторожит его полноту в обе стороны. Реестр
// называет компонент, а классы за него называет документация.
//
//	go run ./cmd/registry
//	go run ./cmd/registry -v   со списком разобранного
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// ── словари осей ───────────────────────────────────────────────────────────
//
// Закрытые, как и атрибутные словари кита: значение вне словаря не даёт ошибки
// нигде, оно просто ничего не делает. И «none» здесь ЗНАЧЕНИЕ, а не отсутствие
// ключа: пропуск и осознанное «ничего не нужно» обязаны различаться, иначе гейт
// превращается в пожелание.
var vocab = map[string][]string{
	"flow":   {"inline", "block"},
	"print":  {"keep", "hide"},
	"page":   {"flow", "whole"},
	"forced": {"none", "border", "fill", "selection", "glyph"},
}

var axisOrder = []string{"flow", "print", "page", "forced"}

// Значение оси, при котором в CSS не должно быть НИЧЕГО. Обратная сверка
// опирается на него: класс, попавший в список при базовом значении, — ошибка
// ровно так же, как класс, не попавший туда при объявленной политике.
var base = map[string]string{"flow": "block", "print": "keep", "page": "flow", "forced": "none"}

type entry struct {
	Axes  map[string]string            `json:"axes"`
	Aria  *aria                        `json:"aria,omitempty"`
	Parts map[string]map[string]string `json:"parts,omitempty"`
}

// Разметочный контракт: роль — это обещание, и невыполненное обещание хуже
// необъявленного. Здесь лежит то, что обязан написать АВТОР разметки; что
// делает за него instrument.js, сказано полем roving.
type aria struct {
	On       string   `json:"on"` // класс-носитель роли, если не корень
	Role     string   `json:"role"`
	Roving   string   `json:"roving"` // author · js
	Requires []string `json:"requires"`
	Item     *struct {
		Role     string   `json:"role"`
		State    string   `json:"state"`
		Requires []string `json:"requires"`
	} `json:"item"`
	Exceptions map[string]struct {
		When string `json:"when"`
		Why  string `json:"why"`
	} `json:"exceptions"`
}

// ── разбор CSS ─────────────────────────────────────────────────────────────

type decl struct {
	sel, at, prop, value string
}

var commentRe = regexp.MustCompile(`(?s)/\*.*?\*/`)
var classRe = regexp.MustCompile(`\.(inst-[A-Za-z0-9_-]+)`)
var wsRe = regexp.MustCompile(`\s+`)

// walk разбирает файл в плоский список объявлений с полным селектором и
// цепочкой at-правил. Вложенность CSS учитывается стеком: у кита она повсюду,
// и правило внутри `&:hover` обязано знать своего хозяина.
func walk(css string) []decl {
	css = commentRe.ReplaceAllString(css, "")
	var out []decl
	var stack []string
	var prelude strings.Builder

	flush := func() {
		d := strings.TrimSpace(prelude.String())
		prelude.Reset()
		if d == "" || len(stack) == 0 {
			return
		}
		i := strings.Index(d, ":")
		if i < 0 {
			return
		}
		prop := strings.TrimSpace(d[:i])
		if prop == "" || strings.ContainsAny(prop, "{}()") {
			return
		}
		var sel, at []string
		for _, s := range stack {
			if strings.HasPrefix(s, "@") {
				at = append(at, s)
			} else {
				sel = append(sel, s)
			}
		}
		out = append(out, decl{
			sel:   strings.Join(sel, " "),
			at:    strings.Join(at, " "),
			prop:  prop,
			value: strings.TrimSpace(d[i+1:]),
		})
	}

	for _, ch := range css {
		switch ch {
		case '{':
			stack = append(stack, wsRe.ReplaceAllString(strings.TrimSpace(prelude.String()), " "))
			prelude.Reset()
		case '}':
			flush()
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
			prelude.Reset()
		case ';':
			flush()
		default:
			prelude.WriteRune(ch)
		}
	}
	return out
}

func classesIn(sel string) []string {
	seen := map[string]bool{}
	var out []string
	for _, m := range classRe.FindAllStringSubmatch(sel, -1) {
		if !seen[m[1]] {
			seen[m[1]] = true
			out = append(out, m[1])
		}
	}
	return out
}

type set map[string]bool

func (s set) addAll(cs []string) {
	for _, c := range cs {
		s[c] = true
	}
}
func (s set) sorted() []string {
	out := make([]string, 0, len(s))
	for k := range s {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// ── что объявлено в CSS ────────────────────────────────────────────────────

// declared собирает фактические списки из сквозных файлов. Читается настоящий
// CSS, а не копия значений: копия расходится молча.
func declared(src map[string]string) map[string]set {
	out := map[string]set{
		"print:hide":       {},
		"page:whole":       {},
		"forced:border":    {},
		"forced:fill":      {},
		"forced:selection": {},
		"forced:glyph":     {},
		"flow:inline":      {},
	}

	for _, d := range walk(src["print.css"]) {
		switch {
		case d.prop == "display" && d.value == "none":
			out["print:hide"].addAll(classesIn(d.sel))
		case d.prop == "break-inside" && strings.Contains(d.value, "avoid"):
			out["page:whole"].addAll(classesIn(d.sel))
		}
	}

	for name, css := range src {
		for _, d := range walk(css) {
			if !strings.Contains(d.at, "forced-colors") {
				continue
			}
			cs := classesIn(d.sel)
			switch {
			case d.prop == "border" && strings.Contains(d.value, "CanvasText"):
				out["forced:border"].addAll(cs)
			case d.prop == "background" && strings.Contains(d.value, "CanvasText"):
				out["forced:glyph"].addAll(cs)
			case d.prop == "background" && d.value == "Highlight":
				if strings.Contains(d.sel, "aria-selected") {
					out["forced:selection"].addAll(cs)
				} else {
					out["forced:fill"].addAll(cs)
				}
			}
			_ = name
		}
	}

	// Поток объявляется компонентом у себя: `align-self: var(--flow-self, …)`
	// стоит рядом с его `display: inline-flex`. Печать исключена — там свои
	// правила потока, к оси отношения не имеющие.
	for name, css := range src {
		if name == "print.css" {
			continue
		}
		for _, d := range walk(css) {
			if d.prop == "align-self" && strings.Contains(d.value, "--flow-self") {
				out["flow:inline"].addAll(classesIn(d.sel))
			}
		}
	}

	// Приоритет: у класса одна политика forced-colors, самая специфичная.
	// Глиф — это маска, выделение — строка списка, заливка — носитель смысла,
	// рамка — просто поверхность, которой нужен контур.
	for _, c := range out["forced:glyph"].sorted() {
		delete(out["forced:selection"], c)
		delete(out["forced:fill"], c)
		delete(out["forced:border"], c)
	}
	for _, c := range out["forced:selection"].sorted() {
		delete(out["forced:fill"], c)
		delete(out["forced:border"], c)
	}
	for _, c := range out["forced:fill"].sorted() {
		delete(out["forced:border"], c)
	}
	return out
}

// ── класс → компонент, из frontmatter документации ─────────────────────────

var apiClassRe = regexp.MustCompile(`name:\s*"(inst-[a-z0-9-]+)"\s*,\s*kind:\s*"класс"`)

func ownership(docs string) (map[string]string, map[string][]string, error) {
	owner := map[string]string{}
	classes := map[string][]string{}
	var pages []string
	for _, dir := range []string{"components", "agent", "layout"} {
		err := filepath.WalkDir(filepath.Join(docs, dir), func(p string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return err
			}
			if strings.HasSuffix(p, ".md") && !strings.HasSuffix(p, ".en.md") {
				pages = append(pages, p)
			}
			return nil
		})
		if err != nil {
			return nil, nil, err
		}
	}
	// Страница без классов в api компонентом не считается и в реестр не
	// попадает. Такая ровно одна — палитра графика: она документирует токены
	// --chart-*, у которых сквозных осей нет и быть не может. Требовать от
	// неё запись значило бы завести четыре решения о том, чего не существует.
	sort.Strings(pages)
	for _, p := range pages {
		b, err := os.ReadFile(p)
		if err != nil {
			return nil, nil, err
		}
		text := strings.ReplaceAll(string(b), "\r\n", "\n")
		i := strings.Index(text, "\n---")
		if !strings.HasPrefix(text, "---\n") || i < 0 {
			continue
		}
		fm := text[4 : i+1]
		name := strings.TrimSuffix(filepath.Base(p), ".md")
		for _, m := range apiClassRe.FindAllStringSubmatch(fm, -1) {
			classes[name] = append(classes[name], m[1])
			// Первая страница, назвавшая класс, и есть его владелец. Чужой
			// класс в своём api страница называть вправе — поиск показывает
			// поле, врез показывает свойство, — но политику за него объявляет
			// владелец, и только он.
			if _, ok := owner[m[1]]; !ok {
				owner[m[1]] = name
			}
		}
	}
	return owner, classes, nil
}

// ── проверки ───────────────────────────────────────────────────────────────

func main() {
	var (
		src      = flag.String("src", "../src", "каталог кита")
		docs     = flag.String("docs", "../docs", "каталог документации")
		registry = flag.String("registry", "../components.json", "реестр компонентов")
		verbose  = flag.Bool("v", false, "показать разобранное")
	)
	flag.Parse()

	css := map[string]string{}
	ents, err := os.ReadDir(*src)
	if err != nil {
		fmt.Fprintln(os.Stderr, "не прочитать кит:", err)
		os.Exit(1)
	}
	for _, e := range ents {
		if strings.HasSuffix(e.Name(), ".css") {
			b, err := os.ReadFile(filepath.Join(*src, e.Name()))
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			css[e.Name()] = string(b)
		}
	}

	raw, err := os.ReadFile(*registry)
	if err != nil {
		fmt.Fprintln(os.Stderr, "не прочитать реестр:", err)
		os.Exit(1)
	}
	var reg map[string]entry
	if err := json.Unmarshal(raw, &reg); err != nil {
		fmt.Fprintln(os.Stderr, "реестр не разбирается:", err)
		os.Exit(1)
	}

	kitJSBytes, err := os.ReadFile(filepath.Join(*src, "kit.js"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "не прочитать instrument.js:", err)
		os.Exit(1)
	}
	kitJS := groupsBlockRe.FindString(string(kitJSBytes))

	owner, classes, err := ownership(*docs)
	if err != nil {
		fmt.Fprintln(os.Stderr, "не прочитать документацию:", err)
		os.Exit(1)
	}

	decl := declared(css)
	var problems []string
	checks := 0
	add := func(f string, a ...any) { problems = append(problems, "  · "+fmt.Sprintf(f, a...)) }

	// ── 1. полнота: компонент есть и в документации, и в реестре ───────────
	for name := range classes {
		checks++
		if _, ok := reg[name]; !ok {
			add("компонент %q задокументирован, но в реестре его нет.\n      "+
				"Заведите запись: без неё ни одна сквозная ось за ним не следит", name)
		}
	}
	for name := range reg {
		checks++
		if _, ok := classes[name]; !ok {
			add("реестр знает %q, но такой страницы с классами нет", name)
		}
	}

	// ── 2. обязательность осей: пропуск отличается от осознанного «ничего» ──
	for name, e := range reg {
		for _, ax := range axisOrder {
			checks++
			v, ok := e.Axes[ax]
			if !ok {
				add("%s: не объявлена ось %q.\n      "+
					"Если ничего не нужно — напишите %q явно: пропуск и решение обязаны различаться",
					name, ax, base[ax])
				continue
			}
			if !contains(vocab[ax], v) {
				add("%s: ось %q = %q вне словаря (%s)", name, ax, v, strings.Join(vocab[ax], " · "))
			}
		}
		for part, ax := range e.Parts {
			checks++
			if owner[part] != name {
				if owner[part] == "" {
					add("%s: часть %q не названа ни на одной странице", name, part)
				} else {
					add("%s: часть %q принадлежит компоненту %q — политику за неё объявляет он",
						name, part, owner[part])
				}
			}
			for k, v := range ax {
				if !contains(vocab[k], v) {
					add("%s / %s: ось %q = %q вне словаря", name, part, k, v)
				}
			}
		}
	}

	// ── 3. согласованность: реестр ⇄ списки CSS, в обе стороны ─────────────
	//
	// Вперёд: объявил политику — класс обязан стоять в списке. Назад: класс
	// стоит в списке — политика обязана быть объявлена. Без второй половины
	// список тихо обрастает тем, о чём никто не решал.
	want := map[string]string{} // "ось:значение|класс" -> компонент
	for name, e := range reg {
		root := ""
		if cs := classes[name]; len(cs) > 0 {
			root = cs[0]
		}
		for ax, v := range e.Axes {
			if v == base[ax] || root == "" {
				continue
			}
			want[ax+":"+v+"|"+root] = name
		}
		for part, ax := range e.Parts {
			for k, v := range ax {
				if v == base[k] {
					continue
				}
				want[k+":"+v+"|"+part] = name
			}
		}
	}
	for key, name := range want {
		i := strings.Index(key, "|")
		list, class := key[:i], key[i+1:]
		checks++
		if s, ok := decl[list]; !ok || !s[class] {
			add("%s объявляет %s для .%s, но в CSS его там нет", name, strings.Replace(list, ":", "=", 1), class)
		}
	}
	for list, s := range decl {
		for _, class := range s.sorted() {
			checks++
			if _, ok := want[list+"|"+class]; !ok {
				who := owner[class]
				if who == "" {
					add("CSS относит .%s к %s, но класс не назван ни на одной странице",
						class, strings.Replace(list, ":", "=", 1))
				} else {
					add("CSS относит .%s к %s, но %s этого не объявляет",
						class, strings.Replace(list, ":", "=", 1), who)
				}
			}
		}
	}

	// ── 4. движение ВЫВОДИТСЯ, а не объявляется ───────────────────────────
	//
	// Замедлять при prefers-reduced-motion надо ровно то, у чего бесконечная
	// анимация: сжатая до 0.01ms бесконечная анимация — это остановка, а
	// остановленный индикатор сообщает, что работа встала. Список motion.css
	// обязан совпадать с этим выводом, и реестр здесь не участвует — ось,
	// которую можно вывести, не заводится.
	spins, slowed := set{}, set{}
	for name, text := range css {
		if name == "motion.css" || name == "print.css" {
			continue
		}
		for _, d := range walk(text) {
			if strings.HasPrefix(d.prop, "animation") && strings.Contains(d.value, "infinite") {
				spins.addAll(classesIn(d.sel))
			}
		}
	}
	for _, d := range walk(css["motion.css"]) {
		slowed.addAll(classesIn(d.sel))
	}
	for _, c := range spins.sorted() {
		checks++
		if !slowed[c] {
			add(".%s крутится бесконечно, но motion.css его не замедляет.\n      "+
				"Сжатая до 0.01ms бесконечная анимация — это остановка, а остановленный "+
				"индикатор сообщает, что работа встала", c)
		}
	}
	for _, c := range slowed.sorted() {
		checks++
		if !spins[c] {
			add("motion.css замедляет .%s, у которого нет бесконечной анимации", c)
		}
	}

	// ── 5. словарь ролей kit.js против реестра ────────────────────────────
	//
	// В kit.js стоит закрытый список ролей с контрактом клавиатуры, и над ним
	// написано, что он совпадает с таблицей разметочного контракта. Совпадение
	// это до сих пор держалось на памяти — и уже разошлось: скрипт выполняет
	// клавиатуру для tablist, а строки про вкладки в таблице не было.
	//
	// Роль, которую скрипт обслуживает молча, обещает клавиатуру от имени кита;
	// роль, объявленная в реестре без поддержки скрипта, обещает её впустую.
	// Сверяются оба направления.
	jsRoles := set{}
	for _, m := range groupsRe.FindAllStringSubmatch(kitJS, -1) {
		jsRoles[m[1]] = true
	}
	regRoles := set{}
	for name, e := range reg {
		if e.Aria == nil || e.Aria.Item == nil {
			continue
		}
		regRoles[e.Aria.Role] = true
		checks++
		if e.Aria.Roving != "author" && e.Aria.Roving != "js" {
			add("%s: roving = %q вне словаря (author · js)", name, e.Aria.Roving)
		}
	}
	for _, r := range jsRoles.sorted() {
		checks++
		if !regRoles[r] {
			add("instrument.js выполняет клавиатуру для role=%q, но в реестре эту роль не объявляет никто.\n      "+
				"Роль обещает клавиатуру от имени кита — обещание должно быть записано", r)
		}
	}
	for _, r := range regRoles.sorted() {
		checks++
		if !jsRoles[r] {
			add("реестр объявляет группу role=%q, но instrument.js её не обслуживает — обещание клавиатуры пустое", r)
		}
	}

	// ── 6. разметочный контракт по живым примерам ─────────────────────────
	markup, markupChecks := checkMarkup(*docs, reg, classes)
	problems = append(problems, markup...)
	checks += markupChecks

	if *verbose {
		fmt.Printf("компонентов в реестре: %d\n", len(reg))
		fmt.Printf("классов, названных документацией: %d\n", len(owner))
		for _, k := range sortedKeys(decl) {
			fmt.Printf("  %-18s %d: %s\n", k, len(decl[k]), strings.Join(decl[k].sorted(), " "))
		}
		fmt.Printf("  %-18s %d\n", "движение", len(spins))
	}

	fmt.Println()
	if len(problems) > 0 {
		sort.Strings(problems)
		fmt.Printf("── реестр разошёлся с китом (%d) ──\n", len(problems))
		for _, p := range problems {
			fmt.Println(p)
		}
		fmt.Println()
		fmt.Printf("· %d сверок, провалено %d\n", checks, len(problems))
		os.Exit(1)
	}
	fmt.Printf("· %d сверок: %d компонентов, %d сквозных списков — реестр и кит сходятся\n",
		checks, len(reg), len(decl))
}

func contains(xs []string, x string) bool {
	for _, v := range xs {
		if v == x {
			return true
		}
	}
	return false
}

func sortedKeys(m map[string]set) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// Ключи словаря GROUPS в instrument.js — роли, у которых есть контракт
// клавиатуры. Читается сам модуль, а не копия списка: копия и есть то, что
// разошлось.
var groupsRe = regexp.MustCompile(`(?m)^\s{2}([a-z]+):\s*\{`)

// Сам блок GROUPS, а не весь модуль: слово вида `menu:` встречается в файле и
// вне словаря, и без вырезки блока в список ролей попадала бы случайная строка.
var groupsBlockRe = regexp.MustCompile(`(?s)const GROUPS = \{.*?\n\};`)

// ── Разметочный контракт по живым примерам ────────────────────────────────
//
// Роль — это обещание: `role="listbox"` сообщает вспомогательной технологии,
// что стрелки работают. Пока обещание не выполнено, компонент не «не доделан»,
// он ОБМАНЫВАЕТ. До сих пор за таблицей контракта не следило ничего: слов
// «aria» и «role» не было ни в одной из шести проверок.
//
// Проверяются ПРИМЕРЫ, а не CSS, и это единственное правильное место. Пример —
// то, что читатель копирует к себе; пример, забывший обязательный атрибут,
// расходится тиражом. Он же и есть исполняемый тест компонента: другого у кита
// нет, и заводить второй значило бы получить два расходящихся описания.
//
// Разбор — свой, на полсотни строк, потому что зависимостей у tools/ нет
// намеренно. Разметка примеров написана руками и заведомо правильна по форме:
// полноценный парсер здесь решал бы задачу, которой нет.

var fenceRe = regexp.MustCompile("(?s)```html preview[^\n]*\n(.*?)```")

// tag — открывающий тег с его атрибутами и границами поддерева.
type tag struct {
	name  string
	attrs map[string]string
	start int // индекс в списке тегов
	end   int // индекс первого тега ЗА поддеревом
}

var voidTags = map[string]bool{
	"area": true, "base": true, "br": true, "col": true, "embed": true,
	"hr": true, "img": true, "input": true, "link": true, "meta": true,
	"source": true, "track": true, "wbr": true,
}

var tagRe = regexp.MustCompile(`<(/?)([a-zA-Z][\w-]*)((?:\s+[^<>"']+(?:"[^"]*")?)*)\s*(/?)>`)
var attrRe = regexp.MustCompile(`([a-zA-Z_:][-\w:.]*)(?:\s*=\s*"([^"]*)")?`)

// scan разбирает фрагмент в плоский список открывающих тегов, у каждого из
// которых известна граница поддерева. Этого достаточно и для «есть ли атрибут
// на элементе», и для «какие элементы лежат внутри».
func scan(html string) []tag {
	var out []tag
	var stack []int
	for _, m := range tagRe.FindAllStringSubmatch(html, -1) {
		closing, name, rawAttrs, selfClose := m[1] == "/", strings.ToLower(m[2]), m[3], m[4] == "/"
		if closing {
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				out[top].end = len(out)
				if out[top].name == name {
					break
				}
			}
			continue
		}
		attrs := map[string]string{}
		for _, a := range attrRe.FindAllStringSubmatch(rawAttrs, -1) {
			if a[1] != "" {
				attrs[strings.ToLower(a[1])] = a[2]
			}
		}
		out = append(out, tag{name: name, attrs: attrs, start: len(out), end: -1})
		if !selfClose && !voidTags[name] {
			stack = append(stack, len(out)-1)
		} else {
			out[len(out)-1].end = len(out)
		}
	}
	for i := range out {
		if out[i].end < 0 {
			out[i].end = len(out)
		}
	}
	return out
}

func hasClass(t tag, c string) bool {
	for _, f := range strings.Fields(t.attrs["class"]) {
		if f == c {
			return true
		}
	}
	return false
}

// checkMarkup сверяет примеры страниц с ARIA-контрактами реестра.
func checkMarkup(docs string, reg map[string]entry, classes map[string][]string) ([]string, int) {
	var problems []string
	checks := 0

	for name, e := range reg {
		if e.Aria == nil {
			continue
		}
		host := ""
		if cs := classes[name]; len(cs) > 0 {
			host = cs[0]
		}
		if e.Aria.On != "" {
			host = e.Aria.On
		}

		for _, page := range pagesOf(docs) {
			b, err := os.ReadFile(page)
			if err != nil {
				continue
			}
			text := strings.ReplaceAll(string(b), "\r\n", "\n")
			for _, f := range fenceRe.FindAllStringSubmatch(text, -1) {
				tags := scan(f[1])
				for i, t := range tags {
					if !hasClass(t, host) || t.attrs["role"] != e.Aria.Role {
						continue
					}
					// filepath.Rel, а не TrimPrefix: путь приходит флагом и
					// может быть записан с любым разделителем, а сообщение
					// обязано читаться одинаково.
					short := page
					if rel, err := filepath.Rel(docs, page); err == nil {
						short = filepath.ToSlash(rel)
					}
					where := fmt.Sprintf("%s пример с .%s", short, host)
					checks++
					for _, req := range e.Aria.Requires {
						if _, ok := t.attrs[req]; ok {
							continue
						}
						if ex, has := e.Aria.Exceptions[req]; has && matchesWhen(t, ex.When) {
							continue
						}
						problems = append(problems, fmt.Sprintf(
							"  · %s: нет %s при role=%q.\n      Роль обещает это вспомогательной технологии — без атрибута обещание ложно",
							where, req, e.Aria.Role))
					}
					if e.Aria.Item == nil {
						continue
					}
					var items []tag
					for _, c := range tags[i+1 : t.end] {
						if c.attrs["role"] == e.Aria.Item.Role {
							items = append(items, c)
						}
					}
					if len(items) == 0 {
						problems = append(problems, fmt.Sprintf(
							"  · %s: role=%q без единого role=%q внутри", where, e.Aria.Role, e.Aria.Item.Role))
						continue
					}
					checks++
					if st := e.Aria.Item.State; st != "" {
						for _, it := range items {
							if _, ok := it.attrs[st]; !ok {
								problems = append(problems, fmt.Sprintf(
									"  · %s: у role=%q нет %s — состояние живёт в атрибуте, а не в классе",
									where, e.Aria.Item.Role, st))
								break
							}
						}
					}
					for _, req := range e.Aria.Item.Requires {
						for _, it := range items {
							if _, ok := it.attrs[req]; !ok {
								problems = append(problems, fmt.Sprintf(
									"  · %s: у role=%q нет %s", where, e.Aria.Item.Role, req))
								break
							}
						}
					}
					// Бегущий tabindex обязателен там, где его пишет АВТОР.
					// Где его расставляет instrument.js, требовать нечего:
					// у меню пункты лежат в закрытом поповере и в обход не
					// попадают, пока его не открыли.
					if e.Aria.Roving == "author" {
						checks++
						zero := 0
						for _, it := range items {
							if it.attrs["tabindex"] == "0" {
								zero++
							}
						}
						if zero != 1 {
							problems = append(problems, fmt.Sprintf(
								"  · %s: бегущий tabindex — %d пунктов с tabindex=\"0\" из %d, нужен ровно один.\n"+
									"      Без него Tab пройдёт по каждому пункту, и группа перестанет быть одним контролом",
								where, zero, len(items)))
						}
					}
				}
			}
		}
	}
	sort.Strings(problems)
	return problems, checks
}

// matchesWhen проверяет условие исключения вида [data-state="indeterminate"].
func matchesWhen(t tag, when string) bool {
	m := regexp.MustCompile(`\[([\w-]+)="([^"]*)"\]`).FindStringSubmatch(when)
	if m == nil {
		return false
	}
	return t.attrs[m[1]] == m[2]
}

var pageCache []string

func pagesOf(docs string) []string {
	if pageCache != nil {
		return pageCache
	}
	for _, dir := range []string{"components", "agent", "layout", "blocks"} {
		filepath.WalkDir(filepath.Join(docs, dir), func(p string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			// Перевод — такая же страница, и разметка в нём такая же живая:
			// читатель копирует её из справочника и получает работающий
			// компонент или неработающий. Исключить `.en.md` значит отпустить
			// ARIA-контракт ровно на тех страницах, которые пишутся заново и
			// потому ошибаются чаще оригинала.
			if strings.HasSuffix(p, ".md") {
				pageCache = append(pageCache, p)
			}
			return nil
		})
	}
	sort.Strings(pageCache)
	return pageCache
}

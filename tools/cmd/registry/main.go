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
	Parts map[string]map[string]string `json:"parts,omitempty"`
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

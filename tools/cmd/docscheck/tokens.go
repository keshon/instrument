// Сверка ТАБЛИЦ ТОКЕНОВ в справочнике с объявлениями в tokens.css.
//
// Основная часть docscheck проверяет СУЩЕСТВОВАНИЕ: есть ли класс, есть ли
// токен, есть ли страница. Она не может поймать самый частый вид расхождения —
// когда документация называет значение, которого у токена больше нет.
//
// Пример, ради которого проверка написана: таблица поверхностей утверждала, что
// --surface-recessed это «чёрный 5%», а в коде уже стояло 0.060. Класс на месте,
// токен на месте, docscheck молчал — врало только число.
//
// Проверяются лишь машиночитаемые ячейки. Строка вида
//
//	| `--surface-sunken` | `--n-2` | `--n-14` |
//
// требует, чтобы в коде было light-dark(var(--n-2), var(--n-14)); строка вида
//
//	| `--surface-hover` | чёрный 3.5% | белый 4.5% |
//
// требует light-dark(oklch(0 0 0 / 0.035), oklch(1 0 0 / 0.045)).
//
// Строка, чей токен объявлен через light-dark, проверяется ЦЕЛИКОМ: обе её
// ячейки обязаны читаться машинно. Пропускать непонятую ячейку значило бы
// принять «black 9%» за прозу и промолчать: непонятое слово уносит с собой и
// число, которое стоит рядом.
package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"instrument/tools/internal/css"
)

var (
	// строка таблицы: | `--имя` | ячейка | ячейка |
	rowRe = regexp.MustCompile("^\\|\\s*`(--[a-z][\\w-]*)`\\s*\\|([^|]*)\\|([^|]*)\\|\\s*$")
	// начало объявления через light-dark; хвост разбирается счётом скобок,
	// потому что внутри лежит oklch(...) со своими
	ldHead = regexp.MustCompile(`(--[a-z][\w-]*):\s*light-dark\(`)
	// ячейка-ссылка на другой токен
	refRe = regexp.MustCompile("^`(--[a-z][\\w-]*)`$")
	// ячейка-подмес: `--ok-4` 16% — цвет роли, разведённый прозрачным
	refPctRe = regexp.MustCompile("^`(--[a-z][\\w-]*)`\\s+([0-9.]+)%$")
	// ячейка вида «чёрный 3.5%» или «белый 8%».
	//
	// Слов ровно два на язык, и оба перечислены здесь поимённо, а не сведены
	// к «любое слово перед процентом». Третий язык обязан упасть на разборе:
	// непонятая ячейка теперь ошибка, и это единственное, что отличает
	// перевод от опечатки.
	alphaRe = regexp.MustCompile(`^(чёрный|белый|black|white)\s+([0-9.]+)%$`)
	// численная альфа, чтобы 0.05 и 0.050 считались одним значением
	alphaNum = regexp.MustCompile(`oklch\(([\d ]+)\s*/\s*([\d.]+)\)`)
)

// splitArgs делит содержимое light-dark(...) по запятой ВЕРХНЕГО уровня.
// Возвращает nil, если аргументов не два.
func splitArgs(s string) []string {
	var out []string
	depth, start := 0, 0
	for i, r := range s {
		switch r {
		case '(':
			depth++
		case ')':
			depth--
		case ',':
			if depth == 0 {
				out = append(out, strings.TrimSpace(s[start:i]))
				start = i + 1
			}
		}
	}
	out = append(out, strings.TrimSpace(s[start:]))
	if len(out) != 2 {
		return nil
	}
	return out
}

// lightDarkPairs собирает первое объявление каждого токена: это :root, то есть
// базовая ячейка. Переопределения тем и масштабов таблица не описывает.
func lightDarkPairs(css string) map[string][2]string {
	out := map[string][2]string{}
	for _, loc := range ldHead.FindAllStringSubmatchIndex(css, -1) {
		name := css[loc[2]:loc[3]]
		if _, seen := out[name]; seen {
			continue
		}
		depth, end := 1, -1
		for i := loc[1]; i < len(css); i++ {
			if css[i] == '(' {
				depth++
			} else if css[i] == ')' {
				depth--
				if depth == 0 {
					end = i
					break
				}
			}
		}
		if end < 0 {
			continue
		}
		if a := splitArgs(css[loc[1]:end]); a != nil {
			out[name] = [2]string{a[0], a[1]}
		}
	}
	return out
}

// expect переводит ячейку таблицы в то, что обязано стоять в коде.
// Второе значение — false, если ячейка не машиночитаема.
func expect(cell string) (string, bool) {
	cell = strings.TrimSpace(cell)
	if m := refRe.FindStringSubmatch(cell); m != nil {
		return "var(" + m[1] + ")", true
	}
	if m := refPctRe.FindStringSubmatch(cell); m != nil {
		pct, err := strconv.ParseFloat(m[2], 64)
		if err != nil {
			return "", false
		}
		return fmt.Sprintf("color-mix(in oklab, var(%s) %g%%, transparent)", m[1], pct), true
	}
	if m := alphaRe.FindStringSubmatch(cell); m != nil {
		pct, err := strconv.ParseFloat(m[2], 64)
		if err != nil {
			return "", false
		}
		lum := "0 0 0"
		if m[1] == "белый" || m[1] == "white" {
			lum = "1 0 0"
		}
		return fmt.Sprintf("oklch(%s / %g)", lum, pct/100), true
	}
	return "", false
}

func sameValue(want, got string) bool {
	if strings.Join(strings.Fields(want), " ") == strings.Join(strings.Fields(got), " ") {
		return true
	}
	w := alphaNum.FindStringSubmatch(want)
	g := alphaNum.FindStringSubmatch(got)
	if w == nil || g == nil {
		return false
	}
	if strings.Join(strings.Fields(w[1]), " ") != strings.Join(strings.Fields(g[1]), " ") {
		return false
	}
	a, err1 := strconv.ParseFloat(w[2], 64)
	b, err2 := strconv.ParseFloat(g[2], 64)
	return err1 == nil && err2 == nil && a == b
}

// langVariants — страница и все её переводы: tokens.md и tokens.en.md.
//
// Сверять только базовое имя значило бы отпустить таблицы ровно на тех
// страницах, которые пишутся заново и потому ошибаются чаще оригинала. Цена
// известна: тот же фильтр в cmd/registry отпускал половину ARIA-контракта.
// Сторожит мутация «таблица токенов в переводе соврала».
func langVariants(dir, stem string) []string {
	base, _ := filepath.Glob(filepath.Join(dir, stem+".md"))
	tr, _ := filepath.Glob(filepath.Join(dir, stem+".*.md"))
	return append(base, tr...)
}

// checkTokenTables возвращает список расхождений таблицы с кодом.
func checkTokenTables(srcDir, docsDir string) []string {
	css, err := os.ReadFile(filepath.Join(srcDir, "tokens.css"))
	if err != nil {
		return nil
	}
	decl := lightDarkPairs(string(css))

	var out []string
	for _, page := range langVariants(filepath.Join(docsDir, "foundations"), "tokens") {
		out = append(out, tokenTableOf(page, decl)...)
	}
	return out
}

func tokenTableOf(page string, decl map[string][2]string) []string {
	md, err := os.ReadFile(page)
	if err != nil {
		return nil
	}
	name := filepath.Base(page)

	var out []string
	for i, line := range strings.Split(string(md), "\n") {
		m := rowRe.FindStringSubmatch(strings.TrimRight(line, "\r"))
		if m == nil {
			continue
		}
		got, ok := decl[m[1]]
		if !ok {
			continue // токен объявлен не через light-dark — не наш случай
		}
		for k, cell := range []string{m[2], m[3]} {
			end := "светлая"
			if k == 1 {
				end = "тёмная"
			}
			want, machine := expect(cell)
			if !machine {
				out = append(out, fmt.Sprintf(
					"%s:%d: %s, %s — ячейку «%s» не разобрать.\n      "+
						"Токен объявлен через light-dark, значит колонка машинная: "+
						"непонятая ячейка молчит там, где обязана возразить",
					name, i+1, m[1], end, strings.TrimSpace(cell)))
				continue
			}
			if !sameValue(want, got[k]) {
				out = append(out, fmt.Sprintf(
					"%s:%d: %s, %s — в таблице «%s», в коде «%s»",
					name, i+1, m[1], end, strings.TrimSpace(cell), got[k]))
			}
		}
	}
	return out
}

// ── Таблицы ПЛОТНОСТИ И МАСШТАБА ──────────────────────────────────────────
//
// Проверка выше сверяет двухколоночные таблицы цвета: токен, светлая тема,
// тёмная. У density.md и scale.md форма другая — «токен × режим», и значения
// лежат не в :root, а в блоках [data-density] и [data-scale]. Под эту форму
// сверки не было, и это стоило ровно того, ради чего проверка вообще написана:
// правка --row-pad-y прошла, а таблица осталась называть прежнюю ступень.
// Класс на месте, токен на месте, врало только число.
//
// Разбираются лишь ячейки-ссылки (`--space-4`). Числа и проза пропускаются:
// заставлять таблицу быть машиночитаемой целиком значило бы её выхолостить.

// Таблица в разметке узнаётся по СТРОКЕ-РАЗДЕЛИТЕЛЮ: она обязана стоять
// второй, а шапка — строкой выше. Слово в первой ячейке шапки при этом не
// значит ничего: колонку токена опознаёт modeRowRe по обратным кавычкам вокруг
// имени. Узнавать шапку по слову «Токен» значило бы гасить сверку всей таблицы
// одним переведённым словом — и продолжать печатать «сходятся полностью».
var sepRe = regexp.MustCompile(`^\|(?:\s*:?-{3,}:?\s*\|)+\s*$`)

// Обратная кавычка внутри строки Go не уживается с сырым литералом, поэтому
// склейка: она читается лучше, чем экранирование каждого разделителя таблицы.
const bq = "`"

var modeRowRe = regexp.MustCompile(`^\|\s*` + bq + `(--[a-z][\w-]*)` + bq + `\s*\|(.+)\|\s*$`)
var refCellRe = regexp.MustCompile(`^` + bq + `(--[a-z][\w-]*)` + bq + `$`)

// Своя регулярка объявления: у declRe в main.go одна группа захвата — ему
// нужно только ИМЯ токена, а здесь нужно ещё и значение.
var modeDeclRe = regexp.MustCompile(`(--[a-z][\w-]*)\s*:\s*([^;{}]+);`)

// blockDecls вынимает объявления из ПЕРВОГО блока, чей селектор совпал.
func blockDecls(css, selector string) map[string]string {
	i := strings.Index(css, selector)
	if i < 0 {
		return nil
	}
	j := strings.Index(css[i:], "{")
	if j < 0 {
		return nil
	}
	depth, start := 0, i+j+1
	for k := i + j; k < len(css); k++ {
		switch css[k] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				out := map[string]string{}
				for _, m := range modeDeclRe.FindAllStringSubmatch(css[start:k], -1) {
					out[m[1]] = strings.TrimSpace(m[2])
				}
				return out
			}
		}
	}
	return nil
}

// Код режима объявлен не здесь, а в ките: [data-density="compact"],
// [data-scale="17"]. Список, зашитый в эту команду, повторял бы кит по памяти
// и расходился бы с ним молча. Хуже: сверяя подпись со словарём ПОДПИСЕЙ,
// «плотная» вместо «compact» пришлось бы считать колонкой не про режим — и
// отпустить её целиком.
var modeSelRe = regexp.MustCompile(`(?m)^\[data-(?:density|scale)="([^"]+)"\] \{$`)

func modeSelectors(sheet string) map[string]string {
	out := map[string]string{}
	for _, m := range modeSelRe.FindAllStringSubmatch(sheet, -1) {
		out[m[1]] = m[0]
	}
	return out
}

// Разобранная таблица «токен × режим».
type modeRow struct {
	line  int
	token string
	cells []string
}

type modeTable struct {
	line int      // строка шапки
	head []string // подписи колонок, без колонки токена
	rows []modeRow
}

// headCells режет шапку на ячейки и снимает обратные кавычки: `compact` и
// compact — одна и та же подпись, и различать их значило бы требовать от
// автора таблицы вёрстки, а не смысла.
func headCells(line string) []string {
	line = strings.TrimSpace(line)
	line = strings.TrimSuffix(strings.TrimPrefix(line, "|"), "|")
	var out []string
	for _, c := range strings.Split(line, "|") {
		out = append(out, strings.Trim(strings.TrimSpace(c), bq))
	}
	return out
}

func modeTablesOf(md string) []modeTable {
	lines := strings.Split(strings.ReplaceAll(md, "\r\n", "\n"), "\n")
	var out []modeTable
	for i := 1; i < len(lines); i++ {
		if !sepRe.MatchString(lines[i]) {
			continue
		}
		head := headCells(lines[i-1])
		if len(head) < 2 {
			continue
		}
		t := modeTable{line: i, head: head[1:]}
		for j := i + 1; j < len(lines) && strings.HasPrefix(lines[j], "|"); j++ {
			m := modeRowRe.FindStringSubmatch(lines[j])
			if m == nil {
				continue
			}
			var cells []string
			for _, c := range strings.Split(m[2], "|") {
				cells = append(cells, strings.TrimSpace(c))
			}
			t.rows = append(t.rows, modeRow{j + 1, m[1], cells})
		}
		// Таблица без единой строки-токена — не про режимы: в density.md и
		// scale.md таких пять, от «Значение · Для чего» до порога 24px.
		if len(t.rows) > 0 {
			out = append(out, t)
		}
	}
	return out
}

// machineCols отмечает колонки, у которых хотя бы одна ячейка читается
// машинно — ссылкой на токен или числом в px. Колонка сплошной прозы
// («Когда», «Зачем») режимом не притворяется, и спрашивать с неё нечего.
func machineCols(t modeTable) []bool {
	out := make([]bool, len(t.head))
	for _, r := range t.rows {
		for k, cell := range r.cells {
			if k >= len(out) || out[k] {
				continue
			}
			if refCellRe.MatchString(cell) {
				out[k] = true
				continue
			}
			if _, err := strconv.ParseFloat(strings.TrimSuffix(cell, "px"), 64); err == nil {
				out[k] = true
			}
		}
	}
	return out
}

// column — колонка таблицы, доведённая до селектора в ките.
type column struct {
	head string
	sel  string // "" — базовая колонка: значение стоит в ярусе ролей
	use  bool
}

// resolveCols раскладывает колонки по режимам кита.
//
// Правило одно: машинная колонка обязана быть либо кодом режима из кита, либо
// ЕДИНСТВЕННОЙ безымянной — базовой. Две безымянные машинные колонки значат,
// что одну из них не узнали, и это ошибка, а не повод промолчать. Молчание
// здесь и есть тот отчёт ни о чём, ради которого проверка написана: без этого
// правила неизвестная подпись не даёт значений вовсе, и колонка выпадает.
func resolveCols(name string, t modeTable, modeSel map[string]string) ([]column, []string) {
	machine := machineCols(t)
	cols := make([]column, len(t.head))
	var unnamed []int
	count := 0
	for k, h := range t.head {
		cols[k] = column{head: h}
		if !machine[k] {
			continue
		}
		count++
		if sel, ok := modeSel[h]; ok {
			cols[k].sel, cols[k].use = sel, true
			continue
		}
		unnamed = append(unnamed, k)
	}
	// Одна машинная колонка — это не «токен × режим», а список: сверять её не
	// с чем, и требовать от неё кода режима не за что.
	if count < 2 {
		return nil, nil
	}
	if len(unnamed) > 1 {
		var names []string
		for _, k := range unnamed {
			names = append(names, "«"+t.head[k]+"»")
		}
		return nil, []string{fmt.Sprintf(
			"%s:%d: колонки %s не названы кодом режима.\n      "+
				"Код режима объявляет кит ([data-density], [data-scale]); "+
				"безымянной — базовой — может быть только одна",
			name, t.line, strings.Join(names, ", "))}
	}
	if len(unnamed) == 1 {
		cols[unnamed[0]].use = true
	}
	return cols, nil
}

func checkModeTables(srcDir, docsDir string) []string {
	raw, err := os.ReadFile(filepath.Join(srcDir, "tokens.css"))
	if err != nil {
		return nil
	}
	sheet := commentRe.ReplaceAllString(strings.ReplaceAll(string(raw), "\r\n", "\n"), "")
	modeSel := modeSelectors(sheet)

	// Ярус ролей: база, поверх которой режим переопределяет своё.
	base := map[string]string{}
	for _, sel := range []string{":root {", `:where(:root) {`} {
		for k, v := range blockDecls(sheet, sel) {
			base[k] = v
		}
	}

	cache := map[string]map[string]string{}
	values := func(sel string) map[string]string {
		if v, ok := cache[sel]; ok {
			return v
		}
		out := map[string]string{}
		for k, v := range base {
			out[k] = v
		}
		if sel != "" {
			for k, v := range blockDecls(sheet, sel) {
				out[k] = v
			}
		}
		cache[sel] = out
		return out
	}

	var problems []string
	var pages []string
	for _, stem := range []string{"density", "scale"} {
		pages = append(pages, langVariants(filepath.Join(docsDir, "foundations"), stem)...)
	}
	for _, page := range pages {
		md, err := os.ReadFile(page)
		if err != nil {
			continue
		}
		name := filepath.Base(page)
		for _, t := range modeTablesOf(string(md)) {
			cols, bad := resolveCols(name, t, modeSel)
			if bad != nil {
				problems = append(problems, bad...)
				continue
			}
			for _, r := range t.rows {
				for k, cell := range r.cells {
					if k >= len(cols) || !cols[k].use {
						continue
					}
					got, ok := values(cols[k].sel)[r.token]
					if !ok {
						continue // токен в этом режиме не объявлен
					}

					// Ячейка-ссылка: сверяем имя.
					if ref := refCellRe.FindStringSubmatch(cell); ref != nil {
						if got != "var("+ref[1]+")" {
							problems = append(problems, fmt.Sprintf(
								"%s:%d: %s в режиме «%s» — в таблице «%s», в коде «%s»",
								name, r.line, r.token, cols[k].head, ref[1], got))
						}
						continue
					}

					// Ячейка-число: сверяем ПИКСЕЛИ. Таблицы масштаба набраны
					// кеглями в px, а токены объявлены в rem — без разрешения
					// сотня чисел осталась бы непроверенной, как осталась бы
					// ступень плотности.
					want, err := strconv.ParseFloat(strings.TrimSuffix(cell, "px"), 64)
					if err != nil {
						continue // проза
					}
					px, err := css.ResolvePx(values(cols[k].sel), r.token)
					if err != nil {
						continue // не длина
					}
					if math.Abs(px-want) > 0.01 {
						problems = append(problems, fmt.Sprintf(
							"%s:%d: %s в режиме «%s» — в таблице %g, в коде %g (%s)",
							name, r.line, r.token, cols[k].head, want, px, got))
					}
				}
			}
		}
	}
	return problems
}

// ── Поле source в шапке страницы ──────────────────────────────────────────
//
// Каждая страница называет файл, из которого её компонент сделан, и сайт
// печатает это ссылкой. Проверял её никто: разбиение components.css на восемь
// файлов оставило двадцать пять страниц указывающими в пустоту, и ни один из
// семи гейтов не сказал ни слова.
//
// Стоит эта ошибка ровно того, ради чего справочник существует: человек идёт
// смотреть, как компонент устроен, и попадает в никуда.
var sourceRe = regexp.MustCompile(`(?m)^source:\s*(\S+)\s*$`)

func checkSourceFields(docsDir string) []string {
	root := filepath.Dir(docsDir)
	var problems []string
	filepath.WalkDir(docsDir, func(p string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(p, ".md") {
			return nil
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return nil
		}
		text := strings.ReplaceAll(string(b), "\r\n", "\n")
		i := strings.Index(text, "\n---")
		if !strings.HasPrefix(text, "---\n") || i < 0 {
			return nil
		}
		m := sourceRe.FindStringSubmatch(text[:i])
		if m == nil {
			return nil
		}
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(m[1]))); err != nil {
			rel, _ := filepath.Rel(docsDir, p)
			problems = append(problems, fmt.Sprintf(
				"%s: source указывает на %s — такого файла нет",
				filepath.ToSlash(rel), m[1]))
		}
		return nil
	})
	sort.Strings(problems)
	return problems
}

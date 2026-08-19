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
// требует light-dark(oklch(0 0 0 / 0.035), oklch(1 0 0 / 0.045)). Ячейки, где
// стоит проза («--a-4 14%, в обеих темах»), пропускаются молча: заставлять
// таблицу быть машиночитаемой целиком значило бы выхолостить её.
package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
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
	// ячейка вида «чёрный 3.5%» или «белый 8%»
	alphaRe = regexp.MustCompile(`^(чёрный|белый)\s+([0-9.]+)%$`)
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
	if m := alphaRe.FindStringSubmatch(cell); m != nil {
		pct, err := strconv.ParseFloat(m[2], 64)
		if err != nil {
			return "", false
		}
		lum := "0 0 0"
		if m[1] == "белый" {
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

// checkTokenTables возвращает список расхождений таблицы с кодом.
func checkTokenTables(srcDir, docsDir string) []string {
	css, err := os.ReadFile(filepath.Join(srcDir, "tokens.css"))
	if err != nil {
		return nil
	}
	decl := lightDarkPairs(string(css))

	md, err := os.ReadFile(filepath.Join(docsDir, "foundations", "tokens.md"))
	if err != nil {
		return nil
	}

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
			want, machine := expect(cell)
			if !machine {
				continue
			}
			if !sameValue(want, got[k]) {
				end := "светлая"
				if k == 1 {
					end = "тёмная"
				}
				out = append(out, fmt.Sprintf(
					"tokens.md:%d: %s, %s — в таблице «%s», в коде «%s»",
					i+1, m[1], end, strings.TrimSpace(cell), got[k]))
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

// Шапка таблицы режимов: первая колонка — токен, дальше режимы.
var modeHeadRe = regexp.MustCompile(`^\|\s*Токен\s*\|(.+)\|\s*$`)

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

// Имя колонки → селектор блока, в котором лежит её значение. Базовая колонка
// разрешается пустым селектором: её значения стоят в ярусе ролей.
var modeSelector = map[string]string{
	"compact":      `[data-density="compact"] {`,
	"comfortable":  `[data-density="comfortable"] {`,
	"по умолчанию": "",
	"14":           "",
	"15":           `[data-scale="15"] {`,
	"16":           `[data-scale="16"] {`,
	"17":           `[data-scale="17"] {`,
	"18":           `[data-scale="18"] {`,
}

func checkModeTables(srcDir, docsDir string) []string {
	raw, err := os.ReadFile(filepath.Join(srcDir, "tokens.css"))
	if err != nil {
		return nil
	}
	sheet := commentRe.ReplaceAllString(string(raw), "")

	// Ярус ролей: база, поверх которой режим переопределяет своё.
	base := map[string]string{}
	for _, sel := range []string{":root {", `:where(:root) {`} {
		for k, v := range blockDecls(sheet, sel) {
			base[k] = v
		}
	}

	cache := map[string]map[string]string{}
	values := func(mode string) map[string]string {
		if v, ok := cache[mode]; ok {
			return v
		}
		sel, known := modeSelector[mode]
		if !known {
			cache[mode] = nil
			return nil
		}
		out := map[string]string{}
		for k, v := range base {
			out[k] = v
		}
		for k, v := range blockDecls(sheet, sel) {
			out[k] = v
		}
		cache[mode] = out
		return out
	}

	var problems []string
	for _, name := range []string{"density.md", "scale.md"} {
		md, err := os.ReadFile(filepath.Join(docsDir, "foundations", name))
		if err != nil {
			continue
		}
		var header []string
		for i, line := range strings.Split(string(md), "\n") {
			line = strings.TrimRight(line, "\r")
			if m := modeHeadRe.FindStringSubmatch(line); m != nil {
				header = nil
				for _, c := range strings.Split(m[1], "|") {
					header = append(header, strings.Trim(strings.TrimSpace(c), "`"))
				}
				continue
			}
			m := modeRowRe.FindStringSubmatch(line)
			if m == nil || header == nil {
				continue
			}
			cells := strings.Split(m[2], "|")
			for k, cell := range cells {
				if k >= len(header) {
					break
				}
				text := strings.TrimSpace(cell)
				vals := values(header[k])
				if vals == nil {
					continue // колонка не про режим: «когда», «зачем»
				}
				got, ok := vals[m[1]]
				if !ok {
					continue // токен в этом режиме не объявлен
				}

				// Ячейка-ссылка: сверяем имя.
				if ref := refCellRe.FindStringSubmatch(text); ref != nil {
					if got != "var("+ref[1]+")" {
						problems = append(problems, fmt.Sprintf(
							"%s:%d: %s в режиме «%s» — в таблице «%s», в коде «%s»",
							name, i+1, m[1], header[k], ref[1], got))
					}
					continue
				}

				// Ячейка-число: сверяем ПИКСЕЛИ. Таблицы масштаба набраны
				// кеглями в px, а токены объявлены в rem — без разрешения
				// сотня чисел оставалась бы непроверенной ровно так же, как
				// ступень плотности до сегодня.
				want, err := strconv.ParseFloat(strings.TrimSuffix(text, "px"), 64)
				if err != nil {
					continue // проза
				}
				px, err := css.ResolvePx(vals, m[1])
				if err != nil {
					continue // не длина
				}
				if math.Abs(px-want) > 0.01 {
					problems = append(problems, fmt.Sprintf(
						"%s:%d: %s в режиме «%s» — в таблице %g, в коде %g (%s)",
						name, i+1, m[1], header[k], want, px, got))
				}
			}
		}
	}
	return problems
}

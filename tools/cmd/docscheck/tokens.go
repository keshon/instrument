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
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
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

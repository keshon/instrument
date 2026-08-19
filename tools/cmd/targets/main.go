// Команда targets проверяет цели нажатия против WCAG 2.2 AA (2.5.8).
//
// Норма — 24×24 CSS-пикселя на КАЖДУЮ интерактивную цель. Исключений в
// критерии четыре, и к киту относится ровно одно: «инлайновая» цель внутри
// строки текста. Всё остальное — ссылка в предложении не считается, а
// стилизованный range-бегунок считается — обязано держать норму.
//
// Читает НАСТОЯЩИЙ src/tokens.css и разрешает var() и calc() так же, как
// contrast, поэтому разойтись с китом не может. Проверяет ВСЕ ТРИ плотности:
// именно в compact вся размерная шкала уходит под норму, и увидеть это на
// умолчании нельзя.
//
//	go run ./cmd/targets
//	go run ./cmd/targets -v   со списком пройденных
package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"

	"instrument/tools/internal/css"
)

// min — норма WCAG 2.2 AA 2.5.8 (Target Size, Minimum).
const min = 24.0

// target — интерактивная цель кита.
//
// w и h — выражения из токенов, ровно те, что стоят в CSS. Писать сюда числа
// значило бы завести второй источник истины: проверка разошлась бы с китом
// молча, а это ровно то, против чего написан contrast.
type target struct {
	label string
	where string // где в ките
	w, h  string

	// gap — расстояние до соседней цели.
	//
	// У критерия есть исключение по РАССТОЯНИЮ, и без него проверка врёт.
	// Цель меньше нормы засчитывается, если круг диаметром 24, помещённый в
	// её центр, не пересекает круг соседа: для целей размера S с зазором G
	// это выполняется при S + G ≥ 24. Иначе плотная плотность оказалась бы
	// нарушением по построению, хотя пункт меню в 22px с зазором в 2px
	// критерию удовлетворяет.
	gap string

	// hit — токен, до которого компонент растит ОБЛАСТЬ нажатия невидимым
	// псевдоэлементом. Если он задан, проверяется область, а не форма: цель
	// может остаться маленькой на вид и при этом держать норму.
	hit string

	// alone — цель, у которой соседей нет вовсе, поэтому исключение по
	// расстоянию к ней неприменимо: она обязана держать норму сама.
	alone bool
}

var targets = []target{
	// Кнопки. Иконочная квадратная, поэтому ширина равна высоте. В ряду
	// кнопки разделены зазором ряда.
	{label: "кнопка sm", where: "actions.css .inst-btn--sm", w: "--control-h-sm", h: "--control-h-sm", gap: "--gap-inline"},
	{label: "кнопка md", where: "actions.css .inst-btn", w: "--control-h-md", h: "--control-h-md", gap: "--gap-inline"},
	{label: "кнопка lg", where: "actions.css .inst-btn--lg", w: "--control-h-lg", h: "--control-h-lg", gap: "--gap-inline"},
	{label: "иконочная кнопка sm", where: "actions.css .inst-btn--sm.inst-btn--icon", w: "--control-h-sm", h: "--control-h-sm", gap: "--gap-inline"},

	// Флажки и переключатели. Соседа у них нет: в колонке выбора таблицы
	// флажок стоит один, без подписи, — значит норму держит сам.
	{label: "флажок", where: "forms.css .inst-checkbox > input", w: "--size-check", h: "--size-check", hit: "--tap-min", alone: true},
	{label: "радио", where: "forms.css .inst-radio > input", w: "--size-check", h: "--size-check", hit: "--tap-min", alone: true},
	{label: "свитч", where: "forms.css .inst-switch > input", w: "--size-switch-w", h: "--size-switch-h", hit: "--tap-min", alone: true},

	// Слайдер. Цель — весь контрол: тянут за бегунок, но попадают в дорожку.
	{label: "слайдер", where: "forms.css .inst-slider", w: "--control-h-md", h: "--tap-min", alone: true},

	// Поля: цель — весь контрол, соседей по горизонтали нет.
	{label: "поле sm", where: "forms.css .inst-input--sm", w: "--control-h-sm", h: "--control-h-sm", gap: "--gap-row"},
	{label: "поле md", where: "forms.css .inst-input", w: "--control-h-md", h: "--control-h-md", gap: "--gap-row"},

	// Навигация и списки: соседи разделены собственным зазором компонента.
	{label: "пункт навигации", where: "layout.css .inst-nav-item", w: "--control-h-md", h: "--control-h-md", gap: "--space-1"},
	{label: "пункт пагинации", where: "layout.css .inst-pager-item", w: "--control-h-sm", h: "--control-h-sm", gap: "--space-2"},
	{label: "пункт меню", where: "overlay.css .inst-menu-item", w: "--control-h-sm", h: "--control-h-sm", gap: "--space-1"},
	{label: "узел дерева", where: "text.css .inst-tree-item", w: "--control-h-sm", h: "--control-h-sm", gap: "--space-1"},
	{label: "день календаря", where: "data.css .inst-calendar-day", w: "--control-h-sm", h: "--control-h-sm", gap: "--space-1"},
	{label: "сегмент контрола", where: "actions.css .inst-segmented > button", w: "--control-h-md", h: "calc(var(--control-h-md) - var(--space-2))", gap: "--space-1"},

	// Подстановка. Ширина у неё по тексту, поэтому по ширине считается
	// высота: даже у самой короткой подписи есть боковые отступы, и уже
	// высоты она быть не может.
	{label: "подстановка", where: "rows.css .inst-insert", w: "--control-h-sm", h: "--control-h-sm", gap: "--space-2"},

	// Чип фильтра. Соседи разделены собственным зазором полосы, как у
	// пагинации: ширина по тексту, поэтому по ширине считается высота.
	{label: "чип фильтра", where: "actions.css .inst-chip", w: "--control-h-sm", h: "--control-h-sm", gap: "--space-2"},

	// Крестик тега. Самая маленькая цель кита: подписи рядом нет, а соседний
	// тег стоит вплотную — исключение по расстоянию неприменимо.
	{label: "крестик тега", where: "data.css .inst-tag-remove", w: "--size-chevron", h: "--size-chevron", hit: "--tap-min", alone: true},
}

var densities = []struct{ id, label string }{
	{"", "обычная"},
	{"compact", "плотная"},
	{"comfortable", "свободная"},
}

// Масштаб — вторая размерная ось. Для этой проверки она важнее плотности:
// плотность жмёт цели вниз, к полу критерия, а масштаб тянет вверх, и
// опасное сочетание — крупный масштаб с плотной раскладкой, где ячейка
// возвращает высоты почти к обычным, оставляя кегль крупным.
//
// Ступени вниз нет: плотный режим уже стоит на полу, и масштаб меньше
// единицы пробил бы 2.5.8 по построению.
var scales = []struct{ id, label string }{
	{"", "14px"},
	{"15", "15px"},
	{"16", "16px"},
	{"17", "17px"},
	{"18", "18px"},
}

// combo — ячейка сетки «масштаб × плотность». Оси расплющены в один список,
// чтобы тело проверки осталось одноуровневым.
type combo struct {
	scale, dens string
	label       string
}

func combos() []combo {
	var out []combo
	for _, sc := range scales {
		for _, d := range densities {
			out = append(out, combo{scale: sc.id, dens: d.id, label: sc.label + " · " + d.label})
		}
	}
	return out
}

func main() {
	tokens := flag.String("tokens", "../src/tokens.css", "путь к tokens.css")
	verbose := flag.Bool("v", false, "показать пройденные")
	flag.Parse()

	src, err := css.Load(*tokens)
	if err != nil {
		fmt.Fprintln(os.Stderr, "не прочитать токены:", err)
		os.Exit(1)
	}

	// Ярус ролей объявлен через :where(:root), плотности — атрибутом.
	// Шкала --space-* объявлена в :root, роли — в :where(:root). Нужны обе:
	// зазоры берутся из первой, размеры из второй.
	base := src.Decls(regexp.MustCompile(`(?m)^:root \{`))
	for k, v := range src.Decls(regexp.MustCompile(`:where\(:root\)\s*\{`)) {
		base[k] = v
	}
	if len(base) == 0 {
		fmt.Fprintln(os.Stderr, "не найден ярус ролей — :where(:root)")
		os.Exit(1)
	}

	width := 0
	for _, t := range targets {
		if n := utf8.RuneCountInString(t.label); n > width {
			width = n
		}
	}

	failed, total := 0, 0
	for _, d := range combos() {
		vals := map[string]string{}
		for k, v := range base {
			vals[k] = v
		}
		// Одномерная часть масштаба ложится первой, двумерная ячейка — поверх.
		if d.scale != "" {
			for k, v := range src.Decls(regexp.MustCompile(`\[data-scale="` + d.scale + `"\]\s*\{`)) {
				vals[k] = v
			}
		}
		if d.dens != "" {
			re := regexp.MustCompile(`\[data-density="` + d.dens + `"\]\s*\{`)
			if d.scale != "" {
				re = regexp.MustCompile(`\[data-scale="` + d.scale + `"\]\[data-density="` + d.dens + `"\]`)
			}
			for k, v := range src.Decls(re) {
				vals[k] = v
			}
		}

		fmt.Printf("\nМАСШТАБ · ПЛОТНОСТЬ %s\n", d.label)
		fmt.Println(strings.Repeat("─", width+34))

		for _, t := range targets {
			total++
			w, errW := resolve(vals, t.w)
			h, errH := resolve(vals, t.h)
			pad := strings.Repeat(" ", width-utf8.RuneCountInString(t.label))
			if errW != nil || errH != nil {
				failed++
				fmt.Printf("  ✗ %s%s  ОШИБКА: %v %v\n", t.label, pad, errW, errH)
				continue
			}
			gap := 0.0
			if t.gap != "" && !t.alone {
				if g, err := resolve(vals, t.gap); err == nil {
					gap = g
				}
			}
			if t.hit != "" {
				if v, err := resolve(vals, t.hit); err == nil {
					if v > w {
						w = v
					}
					if v > h {
						h = v
					}
				}
			}
			ok := w+gap >= min && h+gap >= min
			mark := "·"
			if !ok {
				mark, failed = "✗", failed+1
			}
			note := ""
			if gap > 0 {
				note = fmt.Sprintf("  +%.0f зазор", gap)
			}
			fmt.Printf("  %s %s%s  %4.0f×%-4.0f  (нужно %.0f×%.0f)%s\n",
				mark, t.label, pad, w, h, min, min, note)
		}
	}

	fmt.Println()
	if failed > 0 {
		fmt.Printf("✗ провалов: %d из %d\n", failed, total)
		fmt.Println()
		fmt.Println("Цель растится областью нажатия, а не формой:")
		fmt.Println("  .inst-tag-remove { position: relative }")
		fmt.Println("  .inst-tag-remove::before { content: \"\"; position: absolute; inset: -7px }")
		os.Exit(1)
	}
	if *verbose {
		fmt.Printf("· все %d целей держат %.0f×%.0f в трёх плотностях\n", total, min, min)
	} else {
		fmt.Printf("· все %d проверок пройдены в %d сочетаниях масштаба и плотности\n", total, len(combos()))
	}
}

var (
	varRe  = regexp.MustCompile(`var\(\s*(--[a-z][\w-]*)\s*\)`)
	calcRe = regexp.MustCompile(`calc\(([^()]*)\)`)
	pxRe   = regexp.MustCompile(`([\d.]+)px`)
)

// resolve разворачивает var() и считает calc(). Значения здесь геометрические,
// то есть всегда в px: цвета и проценты сюда не приходят, и поэтому разбор
// умещается в тридцать строк вместо резолвера цветов из contrast.
func resolve(vals map[string]string, expr string) (float64, error) {
	if strings.HasPrefix(expr, "--") {
		v, ok := vals[expr]
		if !ok {
			return 0, fmt.Errorf("нет токена %s", expr)
		}
		expr = v
	}
	for i := 0; i < 12 && strings.Contains(expr, "var("); i++ {
		var missing string
		expr = varRe.ReplaceAllStringFunc(expr, func(m string) string {
			name := varRe.FindStringSubmatch(m)[1]
			v, ok := vals[name]
			if !ok {
				missing = name
				return m
			}
			return strings.TrimSpace(v)
		})
		if missing != "" {
			return 0, fmt.Errorf("нет токена %s", missing)
		}
	}
	for strings.Contains(expr, "calc(") {
		before := expr
		expr = calcRe.ReplaceAllStringFunc(expr, func(m string) string {
			inner := calcRe.FindStringSubmatch(m)[1]
			v, ok := arith(inner)
			if !ok {
				return m
			}
			return fmt.Sprintf("%gpx", v)
		})
		if expr == before {
			return 0, fmt.Errorf("не разобрать calc в %q", expr)
		}
	}
	if v, ok := arith(expr); ok {
		return v, nil
	}
	return 0, fmt.Errorf("не разобрать %q", expr)
}

// arith считает выражение из чисел в px. Умножение и деление раньше сложения
// и вычитания — как в арифметике и как в CSS.
func arith(expr string) (float64, bool) {
	expr = pxRe.ReplaceAllString(strings.TrimSpace(expr), "$1")
	f := strings.Fields(expr)
	if len(f) == 0 || len(f)%2 == 0 {
		return 0, false
	}
	num := func(s string) (float64, bool) {
		var v float64
		_, err := fmt.Sscanf(s, "%g", &v)
		return v, err == nil
	}
	first, ok := num(f[0])
	if !ok {
		return 0, false
	}
	vals := []float64{first}
	var ops []string
	for i := 1; i < len(f); i += 2 {
		rhs, ok := num(f[i+1])
		if !ok {
			return 0, false
		}
		switch f[i] {
		case "*":
			vals[len(vals)-1] *= rhs
		case "/":
			if rhs == 0 {
				return 0, false
			}
			vals[len(vals)-1] /= rhs
		case "+", "-":
			ops = append(ops, f[i])
			vals = append(vals, rhs)
		default:
			return 0, false
		}
	}
	out := vals[0]
	for i, op := range ops {
		if op == "+" {
			out += vals[i+1]
		} else {
			out -= vals[i+1]
		}
	}
	return out, true
}

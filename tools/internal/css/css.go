// Package css разбирает объявления кастомных свойств из настоящего CSS кита
// и резолвит их так же, как это делает браузер.
//
// Смысл в том, чтобы проверки НЕ ДУБЛИРОВАЛИ значения кита, а вычисляли их:
// продублированное значение расходится с оригиналом молча и в тот момент,
// когда его никто не смотрит.
package css

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var commentRe = regexp.MustCompile(`(?s)/\*.*?\*/`)

// Source — разобранный файл токенов.
type Source struct{ text string }

func Load(path string) (*Source, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return &Source{text: commentRe.ReplaceAllString(string(b), "")}, nil
}

// bodyAt возвращает тело блока, начинающегося на позиции idx, со
// сбалансированными скобками.
func (s *Source) bodyAt(idx int) string {
	depth, start := 0, -1
	for i := idx; i < len(s.text); i++ {
		switch s.text[i] {
		case '{':
			if depth == 0 {
				start = i + 1
			}
			depth++
		case '}':
			depth--
			if depth == 0 {
				return s.text[start:i]
			}
		}
	}
	return ""
}

var declRe = regexp.MustCompile(`(--[\w-]+)\s*:\s*([^;]+);`)
var spaceRe = regexp.MustCompile(`\s+`)

// Decls вынимает кастомные свойства из ПЕРВОГО блока, совпавшего с re.
//
// Именно первого: блоки [data-density] и медиазапросы попадать сюда не
// должны, иначе высоты контролов затрут базовые объявления.
func (s *Source) Decls(re *regexp.Regexp) map[string]string {
	out := map[string]string{}
	loc := re.FindStringIndex(s.text)
	if loc == nil {
		return out
	}
	for _, m := range declRe.FindAllStringSubmatch(s.bodyAt(loc[0]), -1) {
		out[m[1]] = spaceRe.ReplaceAllString(strings.TrimSpace(m[2]), " ")
	}
	return out
}

// Theme — тема как две независимые ручки, а не отдельный набор токенов.
type Theme struct {
	ID, Label, Scheme string
	Vars              map[string]string
	Base              map[string]string
}

func (t *Theme) lookup(name string) (string, error) {
	if v, ok := t.Vars[name]; ok {
		return v, nil
	}
	if v, ok := t.Base[name]; ok {
		return v, nil
	}
	return "", fmt.Errorf("токен %s не объявлен", name)
}

var varRe = regexp.MustCompile(`var\((--[\w-]+)\)`)

// expand подставляет var() текстом — так же, как браузер делает ДО разбора
// функции. Без этого oklch(0.994 0.002 var(--hue-neutral)) не разобрать.
func (t *Theme) expand(value string, depth int) (string, error) {
	if depth > 30 {
		return "", fmt.Errorf("слишком глубокая подстановка var()")
	}
	if !strings.Contains(value, "var(") {
		return value, nil
	}
	var err error
	out := varRe.ReplaceAllStringFunc(value, func(m string) string {
		name := varRe.FindStringSubmatch(m)[1]
		v, e := t.lookup(name)
		if e != nil && err == nil {
			err = e
		}
		return v
	})
	if err != nil {
		return "", err
	}
	return t.expand(out, depth+1)
}

// RGBA — цвет в sRGB с альфой, компоненты 0…1.
type RGBA struct{ R, G, B, A float64 }

func splitArgs(s string) []string {
	var out []string
	depth := 0
	var cur strings.Builder
	for _, ch := range s {
		switch ch {
		case '(':
			depth++
		case ')':
			depth--
		}
		if ch == ',' && depth == 0 {
			out = append(out, strings.TrimSpace(cur.String()))
			cur.Reset()
			continue
		}
		cur.WriteRune(ch)
	}
	if strings.TrimSpace(cur.String()) != "" {
		out = append(out, strings.TrimSpace(cur.String()))
	}
	return out
}

func inner(s string) string {
	return s[strings.Index(s, "(")+1 : strings.LastIndex(s, ")")]
}

// OklchToSrgb — перевод в sRGB. Матрицы и гамма те же, что у браузера.
func OklchToSrgb(L, C, H, alpha float64) RGBA {
	h := H * math.Pi / 180
	a := C * math.Cos(h)
	b := C * math.Sin(h)
	l := math.Pow(L+0.3963377774*a+0.2158037573*b, 3)
	m := math.Pow(L-0.1055613458*a-0.0638541728*b, 3)
	s := math.Pow(L-0.0894841775*a-1.2914855480*b, 3)

	enc := func(v float64) float64 {
		if v <= 0.0031308 {
			v = 12.92 * v
		} else {
			v = 1.055*math.Pow(math.Max(v, 0), 1.0/2.4) - 0.055
		}
		return math.Min(1, math.Max(0, v))
	}
	return RGBA{
		R: enc(4.0767416621*l - 3.3077115913*m + 0.2309699292*s),
		G: enc(-1.2684380046*l + 2.6097574011*m - 0.3413193965*s),
		B: enc(-0.0041960863*l - 0.7034186147*m + 1.7076147010*s),
		A: alpha,
	}
}

// calcRe находит простое арифметическое выражение целиком.
//
// Кит использует calc() ровно для одного: домножить цветность шага рампы на
// ручку уклона. Полноценный вычислитель CSS здесь не нужен и был бы враньём
// о возможностях — поддержаны числа и четыре действия, всё остальное честно
// падает с ошибкой.
var calcRe = regexp.MustCompile(`calc\(([^()]*)\)`)

// evalCalc сворачивает calc() в число. Умножение и деление раньше сложения
// и вычитания — как в арифметике и как в CSS.
func evalCalc(expr string) (float64, bool) {
	f := strings.Fields(expr)
	if len(f) == 0 || len(f)%2 == 0 {
		return 0, false
	}
	num := func(s string) (float64, bool) { v, err := strconv.ParseFloat(s, 64); return v, err == nil }

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

var (
	oklchRe = regexp.MustCompile(`^oklch\(\s*([\d.]+)\s+([\d.]+)\s+([\d.]+)\s*(?:/\s*([\d.]+)\s*)?\)$`)
	mixRe   = regexp.MustCompile(`^(.*?)\s+([\d.]+)%$`)
)

// Resolve разрешает значение в цвет sRGB для темы.
func (t *Theme) Resolve(value string) (RGBA, error) {
	v, err := t.expand(value, 0)
	if err != nil {
		return RGBA{}, err
	}
	v = strings.TrimSpace(v)

	// calc() сворачивается ДО разбора цвета: браузер к моменту отрисовки
	// делает то же самое, и oklch() получает уже число.
	for calcRe.MatchString(v) {
		bad := false
		v = calcRe.ReplaceAllStringFunc(v, func(m string) string {
			n, ok := evalCalc(calcRe.FindStringSubmatch(m)[1])
			if !ok {
				bad = true
				return m
			}
			return strconv.FormatFloat(n, 'f', -1, 64)
		})
		if bad {
			return RGBA{}, fmt.Errorf("не вычислить calc: %s", v)
		}
	}

	if v == "transparent" {
		return RGBA{}, nil
	}

	if strings.HasPrefix(v, "light-dark(") {
		args := splitArgs(inner(v))
		if len(args) != 2 {
			return RGBA{}, fmt.Errorf("light-dark ожидает два аргумента: %s", v)
		}
		if t.Scheme == "dark" {
			return t.Resolve(args[1])
		}
		return t.Resolve(args[0])
	}

	if strings.HasPrefix(v, "color-mix(") {
		args := splitArgs(inner(v))
		if len(args) != 3 || !strings.EqualFold(args[0], "in oklab") {
			return RGBA{}, fmt.Errorf(`поддержан только "in oklab": %s`, v)
		}
		m := mixRe.FindStringSubmatch(args[1])
		if m == nil {
			return RGBA{}, fmt.Errorf("не разобрать долю: %s", args[1])
		}
		pct, _ := strconv.ParseFloat(m[2], 64)
		pct /= 100
		base, err := t.Resolve(m[1])
		if err != nil {
			return RGBA{}, err
		}
		other, err := t.Resolve(args[2])
		if err != nil {
			return RGBA{}, err
		}
		a := base.A*pct + other.A*(1-pct)
		if a == 0 {
			return RGBA{}, nil
		}
		return RGBA{
			R: (base.R*base.A*pct + other.R*other.A*(1-pct)) / a,
			G: (base.G*base.A*pct + other.G*other.A*(1-pct)) / a,
			B: (base.B*base.A*pct + other.B*other.A*(1-pct)) / a,
			A: a,
		}, nil
	}

	if m := oklchRe.FindStringSubmatch(v); m != nil {
		l, _ := strconv.ParseFloat(m[1], 64)
		c, _ := strconv.ParseFloat(m[2], 64)
		h, _ := strconv.ParseFloat(m[3], 64)
		alpha := 1.0
		if m[4] != "" {
			alpha, _ = strconv.ParseFloat(m[4], 64)
		}
		return OklchToSrgb(l, c, h, alpha), nil
	}

	return RGBA{}, fmt.Errorf("не разобрать цвет: %s", v)
}

// Token разрешает токен по имени.
func (t *Theme) Token(name string) (RGBA, error) {
	v, err := t.lookup(name)
	if err != nil {
		return RGBA{}, err
	}
	return t.Resolve(v)
}

// Composite накладывает fg на bg.
func Composite(fg, bg RGBA) RGBA {
	return RGBA{
		R: fg.R*fg.A + bg.R*(1-fg.A),
		G: fg.G*fg.A + bg.G*(1-fg.A),
		B: fg.B*fg.A + bg.B*(1-fg.A),
		A: fg.A + bg.A*(1-fg.A),
	}
}

// Flatten схлопывает стопку токенов в непрозрачный цвет. Первый — база.
func (t *Theme) Flatten(stack []string) (RGBA, error) {
	var out RGBA
	for i, name := range stack {
		c, err := t.Token(name)
		if err != nil {
			return RGBA{}, err
		}
		if i == 0 {
			out = c
			continue
		}
		out = Composite(c, out)
	}
	return out, nil
}

func luminance(c RGBA) float64 {
	f := func(v float64) float64 {
		if v <= 0.04045 {
			return v / 12.92
		}
		return math.Pow((v+0.055)/1.055, 2.4)
	}
	return 0.2126*f(c.R) + 0.7152*f(c.G) + 0.0722*f(c.B)
}

// Ratio — контраст переднего плана против фона по WCAG.
func Ratio(fg, bg RGBA) float64 {
	solid := Composite(fg, bg)
	hi, lo := luminance(solid), luminance(bg)
	if lo > hi {
		hi, lo = lo, hi
	}
	return (hi + 0.05) / (lo + 0.05)
}

// ── Геометрия ───────────────────────────────────────────────────────────────

var (
	varUseRe = regexp.MustCompile(`var\(\s*(--[a-z][\w-]*)\s*\)`)
	pxNumRe  = regexp.MustCompile(`([\d.]+)px`)
	remNumRe = regexp.MustCompile(`([\d.]+)rem`)
)

// RootPx — корень, относительно которого считается rem.
//
// В rem кит держит ровно один ярус — кегли, — и делает это затем, чтобы
// уважать увеличенный размер шрифта в браузере. Проверке нужен конкретный
// корень, иначе кегль не с чем сравнить: геометрия объявлена в px. 16 —
// умолчание всех браузеров, и к нему же привязаны комментарии в tokens.css
// («0.875rem × 16 = 14px»).
const RootPx = 16.0

// ResolvePx разворачивает var() и считает calc(), возвращая число пикселей.
//
// Живёт здесь, а не в команде, потому что таких потребителей стало двое:
// cmd/targets меряет цели нажатия, cmd/proportion — пропорции. Копия резолвера
// в каждой из них была бы ровно тем «будущим расхождением», о котором говорит
// конституция: одна поправила бы разбор calc, вторая нет.
//
// Значения здесь всегда геометрические, то есть в px: цвета и проценты сюда
// не приходят, и поэтому разбор умещается в тридцать строк вместо резолвера
// цветов выше.
func ResolvePx(vals map[string]string, expr string) (float64, error) {
	if strings.HasPrefix(expr, "--") {
		v, ok := vals[expr]
		if !ok {
			return 0, fmt.Errorf("нет токена %s", expr)
		}
		expr = v
	}
	for i := 0; i < 12 && strings.Contains(expr, "var("); i++ {
		var missing string
		expr = varUseRe.ReplaceAllStringFunc(expr, func(m string) string {
			name := varUseRe.FindStringSubmatch(m)[1]
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
			v, ok := evalCalc(stripPx(calcRe.FindStringSubmatch(m)[1]))
			if !ok {
				return m
			}
			return fmt.Sprintf("%gpx", v)
		})
		if expr == before {
			return 0, fmt.Errorf("не разобрать calc в %q", expr)
		}
	}
	if v, ok := evalCalc(stripPx(expr)); ok {
		return v, nil
	}
	return 0, fmt.Errorf("не разобрать %q", expr)
}

// stripPx приводит длину к безразмерному числу пикселей: px теряет суффикс,
// rem домножается на корень. Поэтому смешанное выражение вроде
// calc(1rem + 2px) считается верно, а не «как повезёт».
func stripPx(s string) string {
	s = remNumRe.ReplaceAllStringFunc(strings.TrimSpace(s), func(m string) string {
		v, err := strconv.ParseFloat(remNumRe.FindStringSubmatch(m)[1], 64)
		if err != nil {
			return m
		}
		return strconv.FormatFloat(v*RootPx, 'g', -1, 64)
	})
	return pxNumRe.ReplaceAllString(s, "$1")
}

// Команда contrast проверяет контраст токенов кита против порогов WCAG.
//
// Читает НАСТОЯЩИЙ src/tokens.css и резолвит семантику так же, как браузер:
// var() рекурсивно, light-dark() по режиму, color-mix(… transparent) в альфу,
// oklch() в sRGB. Поэтому проверка не может разойтись с китом — она не
// дублирует значения, она их вычисляет.
//
// Правило приёмки:
//
//	текст  < 18px   — 4.5:1  (WCAG 1.4.3)
//	текст ≥ 18px    — 3.0:1
//	несущая граница — 3.0:1  (WCAG 1.4.11): чекбокс, поле, дорожка свитча,
//	                  индикатор состояния. Декоративная рамка, у которой
//	                  рядом есть перепад поверхностей, сюда не входит и
//	                  намеренно тише.
//
//	go run ./cmd/contrast
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

const (
	text  = 4.5 // < 18px
	large = 3.0 // ≥ 18px, а также несущие границы и индикаторы

	// Порог РАЗЛИЧИМОСТИ, а не доступности, и единственный в этом файле.
	//
	// Нужен там, где два элемента лежат вплотную и обязаны читаться разными,
	// но доступность про эту пару ничего не говорит. Случай ровно один:
	// кольцо фокуса вокруг СПЛОШНОЙ заливки. По букве WCAG соседний к кольцу
	// цвет — это поверхность в зазоре отступа, и она проверена тремя парами
	// ниже; заливка кольца не касается. По рисунку же кольцо цвета заливки —
	// это не фокус, а «кнопка стала толще».
	//
	// 3.0 сюда поставить нельзя: самый дальний шаг акцента даёт против заливки
	// 2.63, то есть порог был бы недостижим для любого акцентного кольца и
	// требовал бы нейтрального. 1.5 — это уверенная ступень светлоты; кольцо,
	// совпавшее с заливкой (так было: ровно 1.00), её не берёт.
	distinct = 1.5

	// Порог СТУПЕНИ стопки поверхностей — и мера у него своя.
	//
	// Первый закон поверхностей кита («глубину передаёт порядок светлоты») до
	// сих пор не проверялся ничем: каждый случай выше — это ПЕРЕДНИЙ ПЛАН на
	// ФОНЕ, а поверхность против поверхности не спрашивал никто. Цена: панель и
	// карточка обе стояли на --surface-raised, перепад ровно 1.00, и на
	// странице «Высота и поверхности» карточка внутри панели не читалась.
	//
	// Отношением контраста это не измерить. В формуле WCAG есть слагаемое 0.05
	// — модель бликов, — и у дна светлоты оно сжимает все отношения к единице:
	// один и тот же шаг в 0.035 даёт 1.10 у белого конца рампы и 1.04 у
	// чёрного. Порог по светлому концу объявил бы сломанными тёмные темы,
	// порог по тёмному перестал бы ловить что-либо в светлых.
	//
	// Поэтому ступень меряется РАЗНИЦЕЙ СВЕТЛОТЫ В OKLCH — той же осью, которой
	// набрана рампа, и перцептивно равномерной по построению. Один порог
	// работает во всех шести темах.
	//
	// 0.022 — чуть ниже самой тесной ступени кита (светлый конец, 0.024).
	// Запас намеренно мал: порог, проходящий вдвое, не ловит ничего.
	step = 0.022
)

// Случай проверки: описание, токен переднего плана, стопка фона, порог.
type kase struct {
	label string
	fg    string
	bg    []string
	min   float64
}

// МЕРА выводится из порога, а не хранится отдельным полем.
//
// Отношение контраста по определению не бывает меньше единицы: это частное
// большей светимости на меньшую. Значит, порог ниже единицы не может быть
// отношением — и остаётся ровно один случай, для которого он заведён: ступень
// стопки поверхностей, меряемая разницей светлоты в OKLCH.
//
// Пятое поле в структуре сказало бы то же самое, но потребовало бы дописать
// его в каждый из восьмидесяти восьми существующих случаев — и завести
// возможность рассогласовать флаг с порогом.
func (c kase) isStep() bool { return c.min < 1 }

var cases = []kase{
	// Текст на трёх поверхностях. Каждая ступень проверяется везде, где может
	// оказаться, — именно это правило и не выполнялось раньше.
	{"текст: primary на панели", "--text-primary", []string{"--surface-raised"}, text},
	{"текст: primary на странице", "--text-primary", []string{"--surface-page"}, text},
	{"текст: primary во врезе", "--text-primary", []string{"--surface-sunken"}, text},
	{"текст: secondary на панели", "--text-secondary", []string{"--surface-raised"}, text},
	{"текст: secondary во врезе", "--text-secondary", []string{"--surface-sunken"}, text},
	{"текст: muted на панели", "--text-muted", []string{"--surface-raised"}, text},
	{"текст: muted на странице", "--text-muted", []string{"--surface-page"}, text},
	{"текст: muted во врезе (лог)", "--text-muted", []string{"--surface-sunken"}, text},
	// Утопленные поверхности полупрозрачны и потому объявляются СТОПКОЙ:
	// первым идёт подложка, вторым — сама заливка. Flatten их складывает.
	// Одним слоем такую пару объявить нельзя — гейт прочитает альфу как
	// непрозрачный цвет и померяет контраст с чёрным.
	{"текст: primary в поле на панели", "--text-primary", []string{"--surface-raised", "--surface-field"}, text},
	{"текст: primary в поле на странице", "--text-primary", []string{"--surface-page", "--surface-field"}, text},
	{"текст: подпись кнопки", "--text-primary", []string{"--surface-raised", "--surface-recessed"}, text},
	{"текст: подпись кнопки под курсором", "--text-primary", []string{"--surface-raised", "--surface-recessed-hover"}, text},
	{"текст: тег на панели", "--text-secondary", []string{"--surface-raised", "--surface-recessed"}, text},
	{"текст: чип выбранный", "--accent-text", []string{"--surface-raised", "--surface-recessed"}, text},

	// faint — порог декорации, а не чтения. Читаемому тексту он запрещён.
	{"декор: faint на панели", "--text-faint", []string{"--surface-raised"}, large},
	{"декор: faint во врезе", "--text-faint", []string{"--surface-sunken"}, large},

	// Бейджи: 11px, значит полный текстовый порог.
	{"бейдж: accent на своём фоне", "--accent-text", []string{"--surface-raised", "--accent-bg"}, text},
	{"бейдж: ok на своём фоне", "--ok-text", []string{"--surface-raised", "--ok-bg"}, text},
	{"бейдж: warn на своём фоне", "--warn-text", []string{"--surface-raised", "--warn-bg"}, text},
	{"бейдж: err на своём фоне", "--err-text", []string{"--surface-raised", "--err-bg"}, text},
	{"бейдж: нейтральный во врезе", "--text-secondary", []string{"--surface-sunken"}, text},

	// Статусный текст живёт и вне бейджа — дельта метрики, ошибка поля, сноска.
	{"статус: ok-text на панели", "--ok-text", []string{"--surface-raised"}, text},
	{"статус: warn-text на панели", "--warn-text", []string{"--surface-raised"}, text},
	{"статус: err-text на панели", "--err-text", []string{"--surface-raised"}, text},
	{"статус: ok-text во врезе", "--ok-text", []string{"--surface-sunken"}, text},
	{"статус: warn-text во врезе", "--warn-text", []string{"--surface-sunken"}, text},
	{"статус: err-text во врезе", "--err-text", []string{"--surface-sunken"}, text},

	// Сплошная кнопка. Ховер обязан УВЕЛИЧИВАТЬ контраст подписи, а не ронять.
	{"кнопка: подпись на accent-solid", "--accent-on", []string{"--accent-solid"}, text},
	{"кнопка: подпись на accent-hover", "--accent-on", []string{"--accent-hover"}, text},
	{"ссылка: accent-text на странице", "--accent-text", []string{"--surface-page"}, text},

	// Индикаторы состояния — нетекстовые, но несущие: 3:1. Меряются тем же
	// токеном, которым красятся: точка, заливка меры и штрих истории берут
	// --*-mark, а не текстовую ступень.
	{"метка: ok на панели", "--ok-mark", []string{"--surface-raised"}, large},
	{"метка: ok во врезе", "--ok-mark", []string{"--surface-sunken"}, large},
	{"метка: ok на дорожке", "--ok-mark", []string{"--surface-raised", "--track"}, large},
	{"метка: warn на панели", "--warn-mark", []string{"--surface-raised"}, large},
	{"метка: warn во врезе", "--warn-mark", []string{"--surface-sunken"}, large},
	{"метка: err на панели", "--err-mark", []string{"--surface-raised"}, large},
	{"метка: err во врезе", "--err-mark", []string{"--surface-sunken"}, large},
	{"метка: err на дорожке", "--err-mark", []string{"--surface-raised", "--track"}, large},
	{"точка: running на панели", "--accent-mark", []string{"--surface-raised"}, large},
	{"точка: running во врезе", "--accent-mark", []string{"--surface-sunken"}, large},
	{"каретка на панели", "--accent-mark", []string{"--surface-raised"}, large},
	{"бегунок на дорожке", "--accent-mark", []string{"--surface-raised", "--track"}, large},

	// Несущие границы: граница И ЕСТЬ контрол.
	{"граница контрола на панели", "--border-control", []string{"--surface-raised"}, large},
	{"граница контрола на странице", "--border-control", []string{"--surface-page"}, large},
	{"граница контрола во врезе", "--border-control", []string{"--surface-raised", "--surface-field"}, large},

	// Заполнение меры относительно собственной дорожки, и дорожки — на всех
	// поверхностях, где мера может стоять.
	{"мера: заливка на дорожке (панель)", "--accent-mark", []string{"--surface-raised", "--track"}, large},
	{"мера: заливка на дорожке (врез)", "--accent-mark", []string{"--surface-sunken", "--track"}, large},
	{"мера: ok на дорожке", "--ok-text", []string{"--surface-raised", "--track"}, large},
	{"мера: warn на дорожке", "--warn-text", []string{"--surface-raised", "--track"}, large},
	{"мера: err на дорожке", "--err-text", []string{"--surface-raised", "--track"}, large},

	// Категориальная палитра: каждый ряд обязан отделяться от поверхности.
	{"график: ряд 1 на панели", "--chart-1", []string{"--surface-raised"}, large},
	{"график: ряд 1 на странице", "--chart-1", []string{"--surface-page"}, large},
	{"график: ряд 2 на панели", "--chart-2", []string{"--surface-raised"}, large},
	{"график: ряд 2 на странице", "--chart-2", []string{"--surface-page"}, large},
	{"график: ряд 3 на панели", "--chart-3", []string{"--surface-raised"}, large},
	{"график: ряд 3 на странице", "--chart-3", []string{"--surface-page"}, large},
	{"график: ряд 4 на панели", "--chart-4", []string{"--surface-raised"}, large},
	{"график: ряд 4 на странице", "--chart-4", []string{"--surface-page"}, large},
	{"график: ряд 5 на панели", "--chart-5", []string{"--surface-raised"}, large},
	{"график: ряд 5 на странице", "--chart-5", []string{"--surface-page"}, large},
	{"график: ряд 6 на панели", "--chart-6", []string{"--surface-raised"}, large},
	{"график: ряд 6 на странице", "--chart-6", []string{"--surface-page"}, large},

	// Оверлеи: всё, что лежит на --surface-overlay.
	{"поповер: текст", "--text-primary", []string{"--surface-overlay"}, text},
	{"меню: горячая клавиша", "--text-muted", []string{"--surface-overlay"}, text},
	{"меню: опасный пункт", "--err-text", []string{"--surface-overlay"}, text},
	{"меню: отмеченный пункт", "--accent-text", []string{"--surface-overlay"}, text},
	{"тултип: текст", "--text-primary", []string{"--surface-overlay"}, text},

	// Баннер: текст поверх тонированной заливки.
	{"баннер ok: заголовок", "--text-primary", []string{"--surface-page", "--ok-bg"}, text},
	{"баннер warn: заголовок", "--text-primary", []string{"--surface-page", "--warn-bg"}, text},
	{"баннер error: заголовок", "--text-primary", []string{"--surface-page", "--err-bg"}, text},
	{"баннер warn: пояснение", "--text-secondary", []string{"--surface-page", "--warn-bg"}, text},
	{"баннер warn: значок", "--warn-text", []string{"--surface-page", "--warn-bg"}, large},

	// Формы.
	{"карточка выбора: заголовок", "--text-primary", []string{"--surface-raised", "--accent-bg"}, text},
	{"карточка выбора: описание", "--text-secondary", []string{"--surface-raised", "--accent-bg"}, text},
	{"множественный выбор: выбранный пункт", "--accent-text", []string{"--surface-raised", "--surface-field", "--surface-selected"}, text},
	{"приставка поля", "--text-muted", []string{"--surface-sunken"}, text},
	{"readonly: текст на врезе", "--text-primary", []string{"--surface-sunken"}, text},
	{"пунктир зоны файла", "--border-control", []string{"--surface-raised", "--surface-field"}, large},
	{"обязательность", "--err-text", []string{"--surface-raised"}, text},

	// Раскладка и навигация.
	{"текст на боковой колонке", "--text-secondary", []string{"--surface-sunken"}, text},
	{"навигация: текущий пункт", "--accent-text", []string{"--surface-sunken", "--surface-selected"}, text},
	{"навигация: метка у края", "--accent-solid", []string{"--surface-sunken", "--surface-selected"}, large},
	{"вкладка: подчёркивание", "--accent-solid", []string{"--surface-page"}, large},
	{"крошки: разделитель", "--text-faint", []string{"--surface-page"}, large},
	{"пагинация: текущая страница", "--accent-text", []string{"--surface-page", "--surface-selected"}, text},
	{"шаги: полоса пройденного", "--accent-mark", []string{"--surface-page", "--track"}, large},

	// Инверсная плашка: тултип и всё, что поясняет интерфейс, не будучи им.
	// Пара своя, потому что ни одна ступень текста на ней не лежит: там свой
	// передний план, и проверять его больше некому.
	{"аннотация: текст на инверсии", "--text-on-inverse", []string{"--surface-inverse"}, text},
	// …и сама плашка обязана отделяться от того, над чем висит, иначе смысл
	// «это не содержимое» теряется на первой же тёмной теме.
	{"аннотация: плашка на странице", "--surface-inverse", []string{"--surface-page"}, large},
	{"аннотация: плашка на панели", "--surface-inverse", []string{"--surface-raised"}, large},

	// Кольцо фокуса — против того, что под ним.
	{"фокус: кольцо на странице", "--focus-ring", []string{"--surface-page"}, large},
	{"фокус: кольцо на панели", "--focus-ring", []string{"--surface-raised"}, large},
	{"фокус: кольцо во врезе", "--focus-ring", []string{"--surface-sunken"}, large},
	// …и против того, что оно ОБВОДИТ. Этой пары здесь не было, и потому
	// гейт годами оставался зелёным, пока в светлых темах --focus-ring и
	// --accent-solid были одним и тем же цветом: контраст 1.00 на главной
	// кнопке экрана. Порог — различимость, см. константу.
	{"фокус: кольцо вокруг заливки primary", "--focus-ring", []string{"--accent-solid"}, distinct},

	// ── СТУПЕНИ СТОПКИ ПОВЕРХНОСТЕЙ ────────────────────────────────────────
	//
	// Этих пар здесь не было, и дыра прямо следовала из устройства проверки:
	// каждый случай — это ПЕРЕДНИЙ ПЛАН на ФОНЕ, то есть текст или метка на
	// поверхности. Поверхность против ПОВЕРХНОСТИ не проверял никто.
	//
	// Цена: панель и карточка обе стояли на --surface-raised, перепад ровно
	// 1.00, и на странице «Высота и поверхности» карточка внутри панели не
	// читалась. Гейт был зелёным — все его пары честно проходили, потому что
	// ни одна из них не спрашивала «а различимы ли два соседних яруса».
	//
	// Первый закон поверхностей кита — «глубину передаёт порядок светлоты» —
	// до сих пор не был выражен ни одной проверкой. Теперь выражен.
	//
	// Порог — различимость (см. константу), а не доступность: WCAG про
	// соседние поверхности ничего не говорит, и не должен. 1.5 сюда не годится
	// — это порог для кольца поверх ЗАЛИВКИ, где рядом нет ничего другого.
	// У стопки соседям помогает ещё и рамка, поэтому довольно меньшего, но
	// НЕ ЕДИНИЦЫ: единица означает, что ступени нет вовсе.
	{label: "стопка: страница над углублением", fg: "--surface-page", bg: []string{"--surface-sunken"}, min: step},
	{label: "стопка: панель над страницей", fg: "--surface-raised", bg: []string{"--surface-page"}, min: step},
	{label: "стопка: панель над углублением", fg: "--surface-raised", bg: []string{"--surface-sunken"}, min: step},
	{label: "стопка: карточка в панели", fg: "--surface-sunken", bg: []string{"--surface-raised"}, min: step},

	// Ступень между ВЕСАМИ кнопки, а не между поверхностями стопки. Лестница
	// весов стоит на тех же утоплениях, и её ступени обязаны расходиться так же:
	// сойдясь, мягкая и умолчание дают два имени для одного вида.
	//
	// Проверять приходится здесь, потому что веса собираются КОМПОЗИТОМ поверх
	// подложки и в стопку поверхностей не входят. Пара ловит и обратный случай:
	// на абсолютной ступени мягкая уходила глубже умолчания, и лестница
	// переворачивалась.
	{label: "лестница: мягкая против умолчания", fg: "--surface-recessed-hover", bg: []string{"--surface-raised", "--surface-recessed"}, min: step},
	// Поля вреза здесь НЕТ, и это не пропуск. Поверхность поля совпадает с
	// приподнятой: опознаёт поле несущая рамка, а не глубина, и отвечают за
	// это пары «граница контрола» выше. Требовать от поля ступени значило бы
	// требовать глубины там, где кит её намеренно не даёт.
}

var themes = []*css.Theme{
	{ID: "light-neutral", Label: "светлая нейтральная", Scheme: "light"},
	{ID: "light", Label: "светлая тёплая", Scheme: "light"},
	{ID: "light-cool", Label: "светлая холодная", Scheme: "light"},
	{ID: "dark-light", Label: "тёмная светло-серая", Scheme: "dark"},
	{ID: "dark-soft", Label: "тёмная серая", Scheme: "dark"},
	{ID: "dark", Label: "тёмная чёрная", Scheme: "dark"},
}

// Акцент — вторая ось проверки, и она обязана быть здесь целиком.
//
// Пока набор был один, каждая акцентная пара проверялась ровно в том виде,
// в каком её и рисовали. С ручкой на четыре положения это перестаёт работать:
// три четверти акцентных пар уезжают в поставку непроверенными, а пороги у
// них РАЗНЫЕ — потолок светлоты заливки под белой подписью хью-зависимый
// (0.545 на петроли против 0.580 на глине), и набор, собранный по числам
// соседнего хью, честно провалится.
//
// Пустой ID — база: набор, объявленный в :root без атрибута. Он должен идти
// первым, потому что при подробном выводе показывается именно он.
var accents = []struct{ ID, Label string }{
	{"", "петроль"},
	{"graphite", "графит"},
	{"indigo", "индиго"},
	{"clay", "глина"},
}

func main() {
	tokens := flag.String("tokens", "../src/tokens.css", "путь к tokens.css")
	// Полный вывод — это 528 строк на набор, то есть больше двух тысяч на
	// прогон. Подробно печатается базовый акцент, у остальных — только
	// провалы и итог: список, который никто не читает, не проверка, а шум.
	verbose := flag.Bool("v", false, "печатать все пары для каждого акцента")
	flag.Parse()

	src, err := css.Load(*tokens)
	if err != nil {
		fmt.Fprintln(os.Stderr, "не прочитать токены:", err)
		os.Exit(1)
	}

	// База — первый блок :root. Блоки [data-density] и медиазапросы сюда НЕ
	// попадают: иначе высоты контролов затёрли бы базовые объявления.
	base := src.Decls(regexp.MustCompile(`(?m)^:root \{`))
	for _, t := range themes {
		t.Base = base
		t.Vars = src.Decls(regexp.MustCompile(`\[data-theme="` + t.ID + `"\]\s*\{`))
	}

	width := 0
	for _, c := range cases {
		if n := utf8.RuneCountInString(c.label); n > width {
			width = n
		}
	}

	failed, total := 0, 0
	for _, a := range accents {
		vars := map[string]string{}
		if a.ID != "" {
			vars = src.Decls(regexp.MustCompile(`\[data-accent="` + a.ID + `"\]\s*\{`))
			if len(vars) == 0 {
				fmt.Fprintf(os.Stderr, "акцент %s объявлен в проверке, но не в токенах\n", a.ID)
				os.Exit(1)
			}
		}
		detail := *verbose || a.ID == ""

		fmt.Printf("\n═══ АКЦЕНТ %s ═══\n", a.Label)
		for _, t := range themes {
			t.Accent, t.AccentID, t.AccentLabel = vars, a.ID, a.Label

			bad := 0
			var lines []string
			for _, c := range cases {
				total++
				pad := strings.Repeat(" ", width-utf8.RuneCountInString(c.label))

				fg, err := t.Token(c.fg)
				if err == nil {
					var bg css.RGBA
					bg, err = t.Flatten(c.bg)
					if err == nil {
						// Мера зависит от случая: ступень поверхностей меряется
						// светлотой, всё остальное — отношением контраста.
						r, unit := css.Ratio(fg, bg), ""
						if c.isStep() {
							r, unit = css.Step(fg, bg), " ΔL"
						}
						mark := "·"
						if r < c.min {
							mark, failed, bad = "✗", failed+1, bad+1
						}
						if detail || r < c.min {
							lines = append(lines, fmt.Sprintf("  %s %s%s  %6.3f%s  (нужно %.3f)", mark, c.label, pad, r, unit, c.min))
						}
						continue
					}
				}
				failed, bad = failed+1, bad+1
				lines = append(lines, fmt.Sprintf("  ✗ %s%s  ОШИБКА: %v", c.label, pad, err))
			}

			if detail {
				fmt.Printf("\nТЕМА %s — %s\n", t.ID, t.Label)
				fmt.Println(strings.Repeat("─", width+26))
			} else if bad > 0 {
				fmt.Printf("\nТЕМА %s — %s\n", t.ID, t.Label)
			}
			for _, l := range lines {
				fmt.Println(l)
			}
			if !detail && bad == 0 {
				fmt.Printf("  · %s — %d пар пройдено\n", t.Label, len(cases))
			}
		}
	}

	fmt.Println()
	if failed > 0 {
		fmt.Printf("✗ провалов: %d из %d\n", failed, total)
		os.Exit(1)
	}
	fmt.Printf("· все %d проверок пройдены: %d тем × %d акцентов\n", total, len(themes), len(accents))
}

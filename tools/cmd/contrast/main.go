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
)

// Случай проверки: описание, токен переднего плана, стопка фона, порог.
type kase struct {
	label string
	fg    string
	bg    []string
	min   float64
}

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
	{"текст: primary в поле ввода", "--text-primary", []string{"--surface-field"}, text},

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

	// Индикаторы состояния — нетекстовые, но несущие: 3:1.
	{"точка: ok на панели", "--ok-text", []string{"--surface-raised"}, large},
	{"точка: ok во врезе", "--ok-text", []string{"--surface-sunken"}, large},
	{"точка: warn на панели", "--warn-text", []string{"--surface-raised"}, large},
	{"точка: warn во врезе", "--warn-text", []string{"--surface-sunken"}, large},
	{"точка: err на панели", "--err-text", []string{"--surface-raised"}, large},
	{"точка: err во врезе", "--err-text", []string{"--surface-sunken"}, large},
	{"точка: running на панели", "--accent-mark", []string{"--surface-raised"}, large},
	{"точка: running во врезе", "--accent-mark", []string{"--surface-sunken"}, large},
	{"каретка на панели", "--accent-mark", []string{"--surface-raised"}, large},
	{"бегунок на дорожке", "--accent-mark", []string{"--surface-raised", "--track"}, large},

	// Несущие границы: граница И ЕСТЬ контрол.
	{"граница контрола на панели", "--border-control", []string{"--surface-raised"}, large},
	{"граница контрола на странице", "--border-control", []string{"--surface-page"}, large},
	{"граница контрола во врезе", "--border-control", []string{"--surface-field"}, large},

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
	{"множественный выбор: выбранный пункт", "--accent-text", []string{"--surface-field", "--surface-selected"}, text},
	{"приставка поля", "--text-muted", []string{"--surface-sunken"}, text},
	{"readonly: текст на врезе", "--text-primary", []string{"--surface-sunken"}, text},
	{"пунктир зоны файла", "--border-control", []string{"--surface-field"}, large},
	{"обязательность", "--err-text", []string{"--surface-raised"}, text},

	// Раскладка и навигация.
	{"текст на боковой колонке", "--text-secondary", []string{"--surface-sunken"}, text},
	{"навигация: текущий пункт", "--accent-text", []string{"--surface-sunken", "--surface-selected"}, text},
	{"навигация: метка у края", "--accent-solid", []string{"--surface-sunken", "--surface-selected"}, large},
	{"вкладка: подчёркивание", "--accent-solid", []string{"--surface-page"}, large},
	{"крошки: разделитель", "--text-faint", []string{"--surface-page"}, large},
	{"пагинация: текущая страница", "--accent-text", []string{"--surface-page", "--surface-selected"}, text},
	{"шаги: полоса пройденного", "--accent-mark", []string{"--surface-page", "--track"}, large},

	// Кольцо фокуса — против того, что под ним.
	{"фокус: кольцо на странице", "--focus-ring", []string{"--surface-page"}, large},
	{"фокус: кольцо на панели", "--focus-ring", []string{"--surface-raised"}, large},
	{"фокус: кольцо во врезе", "--focus-ring", []string{"--surface-sunken"}, large},
}

var themes = []*css.Theme{
	{ID: "light", Label: "светлая тёплая", Scheme: "light"},
	{ID: "light-cool", Label: "светлая холодная", Scheme: "light"},
	{ID: "dark", Label: "тёмная глубокая", Scheme: "dark"},
	{ID: "dark-soft", Label: "тёмная мягкая", Scheme: "dark"},
}

func main() {
	tokens := flag.String("tokens", "../src/tokens.css", "путь к tokens.css")
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
	for _, t := range themes {
		fmt.Printf("\nТЕМА %s — %s\n", t.ID, t.Label)
		fmt.Println(strings.Repeat("─", width+26))
		for _, c := range cases {
			total++
			pad := strings.Repeat(" ", width-utf8.RuneCountInString(c.label))

			fg, err := t.Token(c.fg)
			if err == nil {
				var bg css.RGBA
				bg, err = t.Flatten(c.bg)
				if err == nil {
					r := css.Ratio(fg, bg)
					mark := "·"
					if r < c.min {
						mark, failed = "✗", failed+1
					}
					fmt.Printf("  %s %s%s  %6.2f  (нужно %.1f)\n", mark, c.label, pad, r, c.min)
					continue
				}
			}
			failed++
			fmt.Printf("  ✗ %s%s  ОШИБКА: %v\n", c.label, pad, err)
		}
	}

	fmt.Println()
	if failed > 0 {
		fmt.Printf("✗ провалов: %d из %d\n", failed, total)
		os.Exit(1)
	}
	fmt.Printf("· все %d проверок пройдены в %d темах\n", total, len(themes))
}

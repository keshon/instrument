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
	"path/filepath"
	"regexp"
	"sort"
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

	// alt — ВТОРАЯ стопка, с которой сравнивается первая.
	//
	// Нужна там, где два состояния — альтернативы, а не слои. Наведение
	// кнопки ЗАМЕНЯЕТ заливку умолчания (`--btn-bg` переприсваивается), а не
	// ложится поверх неё, и вопрос стоит так: отличается ли наведённая кнопка
	// от ненаведённой на одной и той же подложке. Без alt пара складывала
	// плёнку наведения ПОВЕРХ плёнки умолчания и мерила композицию, которая
	// на экране не рисуется никогда.
	alt []string
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
	{label: "текст: primary на панели", fg: "--text-primary", bg: []string{"--surface-raised"}, min: text},
	{label: "текст: primary на странице", fg: "--text-primary", bg: []string{"--surface-page"}, min: text},
	{label: "текст: primary во врезе", fg: "--text-primary", bg: []string{"--surface-sunken"}, min: text},
	{label: "текст: secondary на панели", fg: "--text-secondary", bg: []string{"--surface-raised"}, min: text},
	{label: "текст: secondary во врезе", fg: "--text-secondary", bg: []string{"--surface-sunken"}, min: text},
	{label: "текст: muted на панели", fg: "--text-muted", bg: []string{"--surface-raised"}, min: text},
	{label: "текст: muted на странице", fg: "--text-muted", bg: []string{"--surface-page"}, min: text},
	{label: "текст: muted во врезе (лог)", fg: "--text-muted", bg: []string{"--surface-sunken"}, min: text},
	// Утопленные поверхности полупрозрачны и потому объявляются СТОПКОЙ:
	// первым идёт подложка, вторым — сама заливка. Flatten их складывает.
	// Одним слоем такую пару объявить нельзя — гейт прочитает альфу как
	// непрозрачный цвет и померяет контраст с чёрным.
	{label: "текст: primary в поле на панели", fg: "--text-primary", bg: []string{"--surface-raised", "--surface-field"}, min: text},
	{label: "текст: primary в поле на странице", fg: "--text-primary", bg: []string{"--surface-page", "--surface-field"}, min: text},
	{label: "текст: подпись кнопки", fg: "--text-primary", bg: []string{"--surface-raised", "--surface-recessed"}, min: text},
	{label: "текст: подпись кнопки под курсором", fg: "--text-primary", bg: []string{"--surface-raised", "--surface-recessed-hover"}, min: text},
	{label: "текст: тег на панели", fg: "--text-secondary", bg: []string{"--surface-raised", "--surface-recessed"}, min: text},
	{label: "текст: чип выбранный", fg: "--accent-text", bg: []string{"--surface-raised", "--surface-recessed"}, min: text},

	// faint — порог декорации, а не чтения. Читаемому тексту он запрещён.
	{label: "декор: faint на панели", fg: "--text-faint", bg: []string{"--surface-raised"}, min: large},
	{label: "декор: faint во врезе", fg: "--text-faint", bg: []string{"--surface-sunken"}, min: large},

	// Бейджи: 11px, значит полный текстовый порог.
	{label: "бейдж: accent на своём фоне", fg: "--accent-text", bg: []string{"--surface-raised", "--accent-bg"}, min: text},
	{label: "бейдж: ok на своём фоне", fg: "--ok-text", bg: []string{"--surface-raised", "--ok-bg"}, min: text},
	{label: "бейдж: warn на своём фоне", fg: "--warn-text", bg: []string{"--surface-raised", "--warn-bg"}, min: text},
	{label: "бейдж: err на своём фоне", fg: "--err-text", bg: []string{"--surface-raised", "--err-bg"}, min: text},
	{label: "бейдж: нейтральный во врезе", fg: "--text-secondary", bg: []string{"--surface-sunken"}, min: text},

	// Статусный текст живёт и вне бейджа — дельта метрики, ошибка поля, сноска.
	{label: "статус: ok-text на панели", fg: "--ok-text", bg: []string{"--surface-raised"}, min: text},
	{label: "статус: warn-text на панели", fg: "--warn-text", bg: []string{"--surface-raised"}, min: text},
	{label: "статус: err-text на панели", fg: "--err-text", bg: []string{"--surface-raised"}, min: text},
	{label: "статус: ok-text во врезе", fg: "--ok-text", bg: []string{"--surface-sunken"}, min: text},
	{label: "статус: warn-text во врезе", fg: "--warn-text", bg: []string{"--surface-sunken"}, min: text},
	{label: "статус: err-text во врезе", fg: "--err-text", bg: []string{"--surface-sunken"}, min: text},

	// Сплошная кнопка. Ховер обязан УВЕЛИЧИВАТЬ контраст подписи, а не ронять.
	{label: "кнопка: подпись на accent-solid", fg: "--accent-on", bg: []string{"--accent-solid"}, min: text},
	{label: "кнопка: подпись на accent-hover", fg: "--accent-on", bg: []string{"--accent-hover"}, min: text},
	{label: "ссылка: accent-text на странице", fg: "--accent-text", bg: []string{"--surface-page"}, min: text},

	// Индикаторы состояния — нетекстовые, но несущие: 3:1. Меряются тем же
	// токеном, которым красятся: точка, заливка меры и штрих истории берут
	// --*-mark, а не текстовую ступень.
	{label: "метка: ok на панели", fg: "--ok-mark", bg: []string{"--surface-raised"}, min: large},
	{label: "метка: ok во врезе", fg: "--ok-mark", bg: []string{"--surface-sunken"}, min: large},
	{label: "метка: ok на дорожке", fg: "--ok-mark", bg: []string{"--surface-raised", "--track"}, min: large},
	{label: "метка: warn на панели", fg: "--warn-mark", bg: []string{"--surface-raised"}, min: large},
	{label: "метка: warn во врезе", fg: "--warn-mark", bg: []string{"--surface-sunken"}, min: large},
	{label: "метка: err на панели", fg: "--err-mark", bg: []string{"--surface-raised"}, min: large},
	{label: "метка: err во врезе", fg: "--err-mark", bg: []string{"--surface-sunken"}, min: large},
	{label: "метка: err на дорожке", fg: "--err-mark", bg: []string{"--surface-raised", "--track"}, min: large},
	{label: "точка: running на панели", fg: "--accent-mark", bg: []string{"--surface-raised"}, min: large},
	{label: "точка: running во врезе", fg: "--accent-mark", bg: []string{"--surface-sunken"}, min: large},
	{label: "каретка на панели", fg: "--accent-mark", bg: []string{"--surface-raised"}, min: large},
	{label: "бегунок на дорожке", fg: "--accent-mark", bg: []string{"--surface-raised", "--track"}, min: large},

	// Несущие границы: граница И ЕСТЬ контрол.
	{label: "граница контрола на панели", fg: "--border-control", bg: []string{"--surface-raised"}, min: large},
	{label: "граница контрола на странице", fg: "--border-control", bg: []string{"--surface-page"}, min: large},
	{label: "граница контрола во врезе", fg: "--border-control", bg: []string{"--surface-raised", "--surface-field"}, min: large},

	// Заполнение меры относительно собственной дорожки, и дорожки — на всех
	// поверхностях, где мера может стоять.
	{label: "мера: заливка на дорожке (панель)", fg: "--accent-mark", bg: []string{"--surface-raised", "--track"}, min: large},
	{label: "мера: заливка на дорожке (врез)", fg: "--accent-mark", bg: []string{"--surface-sunken", "--track"}, min: large},
	{label: "мера: ok на дорожке", fg: "--ok-text", bg: []string{"--surface-raised", "--track"}, min: large},
	{label: "мера: warn на дорожке", fg: "--warn-text", bg: []string{"--surface-raised", "--track"}, min: large},
	{label: "мера: err на дорожке", fg: "--err-text", bg: []string{"--surface-raised", "--track"}, min: large},

	// Категориальная палитра: каждый ряд обязан отделяться от поверхности.
	{label: "график: ряд 1 на панели", fg: "--chart-1", bg: []string{"--surface-raised"}, min: large},
	{label: "график: ряд 1 на странице", fg: "--chart-1", bg: []string{"--surface-page"}, min: large},
	{label: "график: ряд 2 на панели", fg: "--chart-2", bg: []string{"--surface-raised"}, min: large},
	{label: "график: ряд 2 на странице", fg: "--chart-2", bg: []string{"--surface-page"}, min: large},
	{label: "график: ряд 3 на панели", fg: "--chart-3", bg: []string{"--surface-raised"}, min: large},
	{label: "график: ряд 3 на странице", fg: "--chart-3", bg: []string{"--surface-page"}, min: large},
	{label: "график: ряд 4 на панели", fg: "--chart-4", bg: []string{"--surface-raised"}, min: large},
	{label: "график: ряд 4 на странице", fg: "--chart-4", bg: []string{"--surface-page"}, min: large},
	{label: "график: ряд 5 на панели", fg: "--chart-5", bg: []string{"--surface-raised"}, min: large},
	{label: "график: ряд 5 на странице", fg: "--chart-5", bg: []string{"--surface-page"}, min: large},
	{label: "график: ряд 6 на панели", fg: "--chart-6", bg: []string{"--surface-raised"}, min: large},
	{label: "график: ряд 6 на странице", fg: "--chart-6", bg: []string{"--surface-page"}, min: large},

	// Оверлеи: всё, что лежит на --surface-overlay.
	{label: "поповер: текст", fg: "--text-primary", bg: []string{"--surface-overlay"}, min: text},
	{label: "меню: горячая клавиша", fg: "--text-muted", bg: []string{"--surface-overlay"}, min: text},
	{label: "меню: опасный пункт", fg: "--err-text", bg: []string{"--surface-overlay"}, min: text},
	{label: "меню: отмеченный пункт", fg: "--accent-text", bg: []string{"--surface-overlay"}, min: text},
	{label: "тултип: текст", fg: "--text-primary", bg: []string{"--surface-overlay"}, min: text},

	// Баннер: текст поверх тонированной заливки.
	{label: "баннер ok: заголовок", fg: "--text-primary", bg: []string{"--surface-page", "--ok-bg"}, min: text},
	{label: "баннер warn: заголовок", fg: "--text-primary", bg: []string{"--surface-page", "--warn-bg"}, min: text},
	{label: "баннер error: заголовок", fg: "--text-primary", bg: []string{"--surface-page", "--err-bg"}, min: text},
	{label: "баннер warn: пояснение", fg: "--text-secondary", bg: []string{"--surface-page", "--warn-bg"}, min: text},
	{label: "баннер warn: значок", fg: "--warn-text", bg: []string{"--surface-page", "--warn-bg"}, min: large},

	// Формы.
	{label: "карточка выбора: заголовок", fg: "--text-primary", bg: []string{"--surface-raised", "--accent-bg"}, min: text},
	{label: "карточка выбора: описание", fg: "--text-secondary", bg: []string{"--surface-raised", "--accent-bg"}, min: text},
	{label: "множественный выбор: выбранный пункт", fg: "--accent-text", bg: []string{"--surface-raised", "--surface-field", "--surface-selected"}, min: text},
	{label: "приставка поля", fg: "--text-muted", bg: []string{"--surface-sunken"}, min: text},
	{label: "readonly: текст на врезе", fg: "--text-primary", bg: []string{"--surface-sunken"}, min: text},
	{label: "пунктир зоны файла", fg: "--border-control", bg: []string{"--surface-raised", "--surface-field"}, min: large},
	{label: "обязательность", fg: "--err-text", bg: []string{"--surface-raised"}, min: text},

	// Раскладка и навигация.
	{label: "текст на боковой колонке", fg: "--text-secondary", bg: []string{"--surface-sunken"}, min: text},
	{label: "навигация: текущий пункт", fg: "--accent-text", bg: []string{"--surface-sunken", "--surface-selected"}, min: text},
	{label: "навигация: метка у края", fg: "--accent-solid", bg: []string{"--surface-sunken", "--surface-selected"}, min: large},
	{label: "вкладка: подчёркивание", fg: "--accent-solid", bg: []string{"--surface-page"}, min: large},
	{label: "крошки: разделитель", fg: "--text-faint", bg: []string{"--surface-page"}, min: large},
	{label: "пагинация: текущая страница", fg: "--accent-text", bg: []string{"--surface-page", "--surface-selected"}, min: text},
	{label: "шаги: полоса пройденного", fg: "--accent-mark", bg: []string{"--surface-page", "--track"}, min: large},

	// Инверсная плашка: тултип и всё, что поясняет интерфейс, не будучи им.
	// Пара своя, потому что ни одна ступень текста на ней не лежит: там свой
	// передний план, и проверять его больше некому.
	{label: "аннотация: текст на инверсии", fg: "--text-on-inverse", bg: []string{"--surface-inverse"}, min: text},
	// …и сама плашка обязана отделяться от того, над чем висит, иначе смысл
	// «это не содержимое» теряется на первой же тёмной теме.
	{label: "аннотация: плашка на странице", fg: "--surface-inverse", bg: []string{"--surface-page"}, min: large},
	{label: "аннотация: плашка на панели", fg: "--surface-inverse", bg: []string{"--surface-raised"}, min: large},

	// Кольцо фокуса — против того, что под ним.
	{label: "фокус: кольцо на странице", fg: "--focus-ring", bg: []string{"--surface-page"}, min: large},
	{label: "фокус: кольцо на панели", fg: "--focus-ring", bg: []string{"--surface-raised"}, min: large},
	{label: "фокус: кольцо во врезе", fg: "--focus-ring", bg: []string{"--surface-sunken"}, min: large},
	// …и против того, что оно ОБВОДИТ. Этой пары здесь не было, и потому
	// гейт годами оставался зелёным, пока в светлых темах --focus-ring и
	// --accent-solid были одним и тем же цветом: контраст 1.00 на главной
	// кнопке экрана. Порог — различимость, см. константу.
	{label: "фокус: кольцо вокруг заливки primary", fg: "--focus-ring", bg: []string{"--accent-solid"}, min: distinct},

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
	{label: "лестница: мягкая против умолчания", fg: "--surface-recessed-hover", bg: []string{"--surface-raised"},
		alt: []string{"--surface-raised", "--surface-recessed"}, min: step},

	// КНОПКА ПРОТИВ ТОГО, НА ЧЁМ ЛЕЖИТ. Этой пары не было, и её отсутствие
	// стоило дорого.
	//
	// Выше проверено, что поверхности стопки расходятся между собой, и что веса
	// кнопки расходятся между собой. Ни то, ни другое не спрашивает главного:
	// отличается ли кнопка от своей подложки. А у умолчания это ЕДИНСТВЕННЫЙ
	// признак — рамки нет, тени нет, обе снял провал.
	//
	// Провал делается чёрной плёнкой в обеих темах, потому что углубление по
	// смыслу темнее. Но плёнка вычитает светлоту, а вычитать её на дне рампы
	// уже не из чего: страница тёмной темы — rgb(12,12,12), и даже СПЛОШНОЙ
	// чёрный даёт против неё 1.073 против 1.142 у светлых. Кнопка переставала
	// существовать: перепад 1.014.
	//
	// Обе подложки, потому что кнопка стоит и там и там: на панели — в форме и
	// в карточке, на странице — в шапке экрана и в тулбаре оболочки.
	{label: "кнопка против панели", fg: "--surface-recessed", bg: []string{"--surface-raised"}, min: step},
	{label: "кнопка против страницы", fg: "--surface-recessed", bg: []string{"--surface-page"}, min: step},
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
							if len(c.alt) > 0 {
								var other css.RGBA
								if other, err = t.Flatten(c.alt); err == nil {
									r = css.StepOf(css.Composite(fg, bg), other)
								}
							}
						}
						if err != nil {
							failed, bad = failed+1, bad+1
							lines = append(lines, fmt.Sprintf("  ✗ %s%s  ОШИБКА: %v", c.label, pad, err))
							continue
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

	// Охват таблицы: цвет текста, которого в ней нет, — это порог, который
	// никто не мерил. Считается один раз, а не на каждую тему: вопрос не в
	// значении, а в наличии строки.
	gaps, inkCount, err := checkInkCoverage(filepath.Dir(*tokens), cases)
	if err != nil {
		fmt.Fprintln(os.Stderr, "не прочитать кит:", err)
		os.Exit(1)
	}
	total += inkCount
	failed += len(gaps)
	if len(gaps) > 0 {
		fmt.Println("\nОХВАТ ТАБЛИЦЫ")
		for _, g := range gaps {
			fmt.Println(g)
		}
	}

	fmt.Println()
	if failed > 0 {
		fmt.Printf("✗ провалов: %d из %d\n", failed, total)
		os.Exit(1)
	}
	fmt.Printf("· все %d проверок пройдены: %d тем × %d акцентов, охват %d цветов текста\n",
		total, len(themes), len(accents), inkCount)
}

// ── Покрытие: каждый цвет текста обязан быть кем-то проверен ───────────────
//
// Таблица пар выше — рукописная, и это правильно: пара несёт не только два
// токена, но и СТОПКУ композиции («подпись в поле на панели» — два фона,
// потому что врез полупрозрачен и складывается с панелью под ним). Такой факт
// о вложенности разметки из tokens.css не выводится, и генерация пар обезличила
// бы таблицу, у каждой строки которой есть человеческое имя.
//
// А вот ДЫРА в ней выводится, и стоит копейки. Компонент, покрасивший текст
// новым токеном, о котором не подумали, ниоткуда себя не проявит: цвет
// применится, порог никто не померяет, и узнается это от человека, который не
// смог прочитать подпись.
//
// Поэтому здесь сверяется не результат, а ОХВАТ: всякий токен, которым кит
// красит текст, обязан стоять передним планом хотя бы в одной паре. Что именно
// под ним лежит, по-прежнему решает автор пары.
var (
	colorDecl = regexp.MustCompile(`(?:^|[;{])\s*color\s*:\s*([^;}]+)`)
	varUse    = regexp.MustCompile(`var\(\s*(--[a-z][\w-]*)`)
)

// inkTokens — токены, которыми кит реально красит текст, с местом первой
// встречи. tokens.css исключён: там цвет объявляется, а не применяется;
// print.css — тоже, у бумаги свой набор.
func inkTokens(dir string) (map[string]string, error) {
	out := map[string]string{}
	ents, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, e := range ents {
		name := e.Name()
		if !strings.HasSuffix(name, ".css") || name == "tokens.css" || name == "print.css" {
			continue
		}
		b, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return nil, err
		}
		// Комментарии ГАСЯТСЯ, а не вырезаются: они занимают три четверти
		// файла, и после вырезания номер строки уезжает на сотни строк — то
		// есть сообщение показывает на чужое правило.
		text := string(css.Blank([]byte(strings.ReplaceAll(string(b), "\r\n", "\n"))))
		for _, m := range colorDecl.FindAllStringSubmatchIndex(text, -1) {
			value := text[m[2]:m[3]]
			line := strings.Count(text[:m[2]], "\n") + 1
			for _, v := range varUse.FindAllStringSubmatch(value, -1) {
				if _, seen := out[v[1]]; !seen {
					out[v[1]] = fmt.Sprintf("%s:%d", name, line)
				}
			}
		}
	}
	return out, nil
}

// checkInkCoverage возвращает токены, которыми красят текст, но которые не
// стоят передним планом ни в одной паре.
//
// Компонентные переменные пропускаются: --btn-fg и --tone-ink не цвета, а
// ПОДСТАНОВКИ — за ними стоит семантика, и меряется она под своим именем.
// Проверять их значило бы требовать пару на каждое имя-посредник.
func checkInkCoverage(dir string, cases []kase) ([]string, int, error) {
	ink, err := inkTokens(dir)
	if err != nil {
		return nil, 0, err
	}
	covered := map[string]bool{}
	for _, c := range cases {
		covered[c.fg] = true
	}
	var bad []string
	for tok, where := range ink {
		if indirect.MatchString(tok) || covered[tok] {
			continue
		}
		bad = append(bad, fmt.Sprintf(
			"  ✗ %s  красит текст (%s), но не стоит передним планом ни в одной паре.\n"+
				"      Заведите строку в таблице выше: непроверенный цвет текста — это порог, который никто не мерил",
			tok, where))
	}
	sort.Strings(bad)
	return bad, len(ink), nil
}

// Имена-посредники: компонентная переменная и тон. За ними стоит семантика,
// которая меряется отдельно и под своим именем.
var indirect = regexp.MustCompile(`^--(btn|tone|level|change|chart)-`)

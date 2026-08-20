// Команда mutate проверяет ПРОВЕРКИ: ломает по одному инварианту и требует,
// чтобы соответствующий гейт упал.
//
// ЗАЧЕМ ОНА ЕСТЬ. Гейт бывает зелёным и пустым, и чтением это не выявляется.
// Так «лестница весов кнопки» проходила всегда: она сравнивала светлоту двух
// цветов как непрозрачных, а утопления кита — плёнки, и у плёнки бралась её
// собственная светлота, то есть ноль. Перепад против любой подложки выходил
// огромным. Проверка существовала, отчитывалась зелёным и не проверяла ничего
// — пока умолчание кнопки в тёмных темах не перестало отличаться от фона.
//
// Отсюда правило: проверка, которую никогда не видели красной, не проверка, а
// украшение. Каждая мутация ниже — это ответ на вопрос «а если сломать?»,
// заданный машине, а не памяти.
//
// КАК ЭТО РАБОТАЕТ. Дерево копируется во временный каталог, в копии портится
// ровно одна вещь, гейт запускается против копии и ОБЯЗАН вернуть ненулевой
// код. Пойманная мутация — строка «поймана»; непойманная — дыра в гейте, и
// команда падает.
//
// Контрольный прогон идёт первым: на нетронутой копии все гейты обязаны быть
// зелёными. Без него стенд мерил бы не мутации, а собственный шум.
//
//	go run ./cmd/mutate        весь список
//	go run ./cmd/mutate -v     с выводом упавшего гейта
//	go run ./cmd/mutate -only contrast   только мутации одного гейта
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Гейт: имя, пакет и то, как ему объяснить, где лежит дерево.
type gate struct {
	pkg  string
	args func(tree string) []string
}

var gates = map[string]gate{
	"contrast":   {"./cmd/contrast", func(t string) []string { return []string{"-tokens", t + "/src/tokens.css"} }},
	"targets":    {"./cmd/targets", func(t string) []string { return []string{"-tokens", t + "/src/tokens.css"} }},
	"proportion": {"./cmd/proportion", func(t string) []string { return []string{"-tokens", t + "/src/tokens.css"} }},
	"dist":       {"./cmd/dist", func(t string) []string { return []string{"-src", t + "/src", "-out", t + "/dist"} }},
	"docscheck": {"./cmd/docscheck", func(t string) []string {
		return []string{"-kit", t + "/src", "-docs", t + "/docs", "-stage", t + "/stage.css"}
	}},
	"registry": {"./cmd/registry", func(t string) []string {
		return []string{"-src", t + "/src", "-docs", t + "/docs", "-registry", t + "/components.json"}
	}},
}

// mutation — одна поломка. from обязан существовать в файле: мутация, которая
// ничего не заменила, дала бы ложное «поймана» на нетронутом дереве.
type mutation struct {
	name string
	gate string
	file string
	from string
	to   string
	why  string
}

var mutations = []mutation{
	// ── contrast ───────────────────────────────────────────────────────────
	{"текст: muted подняли до нечитаемого", "contrast", "src/tokens.css",
		"--text-muted:     light-dark(var(--n-8),  var(--n-6));",
		"--text-muted:     light-dark(var(--n-5),  var(--n-6));",
		"подпись перестаёт брать 4.5:1 на панели"},
	{"стопка: панель сравнялась со страницей", "contrast", "src/tokens.css",
		"--surface-page:     light-dark(var(--n-1), var(--n-13));",
		"--surface-page:     light-dark(var(--n-0), var(--n-13));",
		"глубина передаётся порядком светлоты, а ступени не стало"},
	{"кнопка сравнялась с подложкой", "contrast", "src/tokens.css",
		"--surface-recessed:        light-dark(oklch(0 0 0 / 0.060), oklch(1 0 0 / 0.050));",
		"--surface-recessed:        light-dark(oklch(0 0 0 / 0.002), oklch(1 0 0 / 0.002));",
		"ровно тот дефект, из-за которого стенд и написан"},
	{"лестница перевёрнута: мягкая громче умолчания", "contrast", "src/tokens.css",
		"--surface-recessed-hover:  light-dark(oklch(0 0 0 / 0.025), oklch(1 0 0 / 0.015));",
		"--surface-recessed-hover:  light-dark(oklch(0 0 0 / 0.025), oklch(1 0 0 / 0.090));",
		"третья ступень звучит громче второй, а модуль разницы этого не видит"},
	{"лестница: наведение сравнялось с умолчанием", "contrast", "src/tokens.css",
		"--surface-recessed-hover:  light-dark(oklch(0 0 0 / 0.025), oklch(1 0 0 / 0.015));",
		"--surface-recessed-hover:  light-dark(oklch(0 0 0 / 0.058), oklch(1 0 0 / 0.048));",
		"два веса кнопки дают один вид"},
	{"текст красят токеном без пары", "contrast", "src/base.css",
		"body {",
		".mut-probe { color: var(--n-7); }\nbody {",
		"цвет текста, которого нет ни в одной паре, — порог, который никто не мерил"},

	// ── targets ────────────────────────────────────────────────────────────
	{"высота контрола ушла под норму", "targets", "src/tokens.css",
		"--control-h-sm: 26px;",
		"--control-h-sm: 14px;",
		"цель нажатия меньше 24 и без зазора, компенсирующего размер"},

	// ── proportion ─────────────────────────────────────────────────────────
	{"ступени кегля слились", "proportion", "src/tokens.css",
		"--text-xs:  0.78125rem;",
		"--text-xs:  0.86rem;",
		"шкала объявляет размер, которого не видно"},
	{"радиус стал нечётным", "proportion", "src/tokens.css",
		"--radius-md: 8px;",
		"--radius-md: 7px;",
		"дуга садится мимо устройственной сетки на плотности 1.5"},
	{"геометрия убывает при росте масштаба", "proportion", "src/tokens.css",
		"  --label-col:     120px;",
		"  --label-col:     8px;",
		"масштаб идёт только вверх"},
	{"вертикаль строки шире горизонтали", "proportion", "src/tokens.css",
		"--row-pad-y:   var(--space-3);",
		"--row-pad-y:   var(--space-6);",
		"часть вертикального воздуха уже принёс интерлиньяж"},

	// ── dist: раздел «Запрещено» ───────────────────────────────────────────
	{"!important в компоненте", "dist", "src/actions.css",
		".inst-btn {", ".mut-probe { color: red !important; }\n.inst-btn {",
		"единственный разрешённый — [hidden] в base.css"},
	{"вес 700", "dist", "src/actions.css",
		".inst-btn {", ".mut-probe { font-weight: 700; }\n.inst-btn {",
		"в ките два веса, и 700 кричит громче данных"},
	{"сырой цвет в компоненте", "dist", "src/actions.css",
		".inst-btn {", ".mut-probe { color: #333; }\n.inst-btn {",
		"число вместо семантики — захардкоженная светлая тема"},
	{"утилита отступа", "dist", "src/layout.css",
		".inst-stack {", ".mt-3 { margin-top: 12px; }\n.inst-stack {",
		"ритм задают примитивы потока, а не разметка каждого экрана"},
	{"цвет внутри data-URI", "dist", "src/agent.css",
		"stroke='%23000'", "stroke='%23ff0000'",
		"форма рисуется маской, цвет приходит токеном"},
	{"картинка фоном, а не маской", "dist", "src/data.css",
		".inst-tag {", ".mut-probe { background-image: url(\"x.svg\"); }\n.inst-tag {",
		"фоновая картинка красится собой и не следует теме"},
	{"просвет мимо яруса ролей", "dist", "src/layout.css",
		".inst-cluster {", ".mut-probe { gap: var(--space-8); }\n.inst-cluster {",
		"крупная ступень шкалы взята напрямую"},
	{"лишняя фигурная скобка", "dist", "src/motion.css",
		".inst-caret { animation: none; opacity: 1; }",
		".inst-caret { animation: none; opacity: 1; } }",
		"остаток файла уходит из слоя и начинает выигрывать у приложения"},

	// ── docscheck ──────────────────────────────────────────────────────────
	{"в примере класс, которого нет", "docscheck", "docs/components/actions/button.md",
		"<button class=\"inst-btn\"", "<button class=\"inst-btn inst-btn--mut\"",
		"читатель копирует разметку, она молча не работает"},
	{"число в таблице токенов соврало", "docscheck", "docs/foundations/tokens.md",
		"| `--surface-recessed` | чёрный 6% | белый 5% |",
		"| `--surface-recessed` | чёрный 9% | белый 5% |",
		"класс на месте, токен на месте, врёт только число"},
	{"число в таблице плотности соврало", "docscheck", "docs/foundations/density.md",
		"| `--row-pad-y` | `--space-2` | `--space-3` | `--space-4` |",
		"| `--row-pad-y` | `--space-2` | `--space-3` | `--space-6` |",
		"таблица режимов разошлась с блоком плотности"},
	{"число в таблице масштаба соврало", "docscheck", "docs/foundations/scale.md",
		"| `--text-sm` | 14 | 15 | 16 | 17 | 18 |",
		"| `--text-sm` | 14 | 15 | 19 | 17 | 18 |",
		"кегль в таблице не тот, что в коде"},
	{"source ведёт в никуда", "docscheck", "docs/components/actions/button.md",
		"source: src/actions.css", "source: src/nowhere.css",
		"человек идёт смотреть устройство и попадает в пустоту"},

	// ── docscheck: язык таблиц ─────────────────────────────────────────────
	//
	// Пять мутаций ниже проверяют одно: гейт сверяет таблицы по СТРУКТУРЕ, а
	// не по русским словам в шапке и в ячейках. Каждая переводит одно слово и
	// одновременно врёт числом. Если гейт ослеп на переводе, ложь проходит, и
	// «документация и кит сходятся полностью» становится отчётом ни о чём.
	{"шапка таблицы режимов переведена", "docscheck", "docs/foundations/density.md",
		"| Токен | `compact` | по умолчанию | `comfortable` |\n|---|---|---|---|\n| `--control-h-xs` | 18px | 20px | 22px |",
		"| Token | `compact` | по умолчанию | `comfortable` |\n|---|---|---|---|\n| `--control-h-xs` | 18px | 20px | 24px |",
		"одно слово в шапке гасило сверку всей таблицы"},
	{"колонка умолчания переведена", "docscheck", "docs/foundations/density.md",
		"| Токен | `compact` | по умолчанию | `comfortable` |\n|---|---|---|---|\n| `--size-check` | 13px | 15px | 17px |",
		"| Токен | `compact` | default | `comfortable` |\n|---|---|---|---|\n| `--size-check` | 13px | 19px | 17px |",
		"колонка узнавалась по подписи, а подпись — единственная русская из трёх"},
	{"код режима в шапке переведён", "docscheck", "docs/foundations/density.md",
		"| Токен | `compact` | по умолчанию | `comfortable` |\n|---|---|---|---|\n| `--pad-panel` | `--space-4` | `--space-5` | `--space-6` |",
		"| Токен | `плотная` | по умолчанию | `comfortable` |\n|---|---|---|---|\n| `--pad-panel` | `--space-2` | `--space-5` | `--space-6` |",
		"код режима — часть кита, и переведённый код обязан быть слышен"},
	{"ячейка альфы переведена, значение соврало", "docscheck", "docs/foundations/tokens.md",
		"| `--surface-recessed` | чёрный 6% | белый 5% |",
		"| `--surface-recessed` | black 9% | белый 5% |",
		"перевод ячейки уносил с собой проверку её числа"},
	{"ячейка альфы не читается", "docscheck", "docs/foundations/tokens.md",
		"| `--surface-recessed-hover` | чёрный 2.5% | белый 1.5% |",
		"| `--surface-recessed-hover` | тёмный 2.5% | белый 1.5% |",
		"непонятая ячейка в машинной колонке — это молчание, а не разрешение"},
	{"таблица токенов в переводе соврала", "docscheck", "docs/foundations/tokens.en.md",
		"",
		"---\ntitle: Tokens\n---\n\n| Token | Light | Dark |\n|---|---|---|\n| `--surface-recessed` | black 9% | white 5% |\n",
		"перевод страницы уносил её таблицы из-под сверки целиком"},
	{"таблица режимов в переводе соврала", "docscheck", "docs/foundations/density.en.md",
		"",
		"---\ntitle: Density\n---\n\n| Token | `compact` | default | `comfortable` |\n|---|---|---|---|\n| `--row-pad-y` | `--space-2` | `--space-3` | `--space-6` |\n",
		"проверка была прибита к русскому имени файла, а не к странице"},

	// ── registry ───────────────────────────────────────────────────────────
	{"компонент вне реестра", "registry", "components.json",
		"\"toolbar\": {", "\"toolbar-renamed\": {",
		"за компонентом перестаёт следить любая сквозная ось"},
	{"у оси убрали значение", "registry", "components.json",
		"\"flow\": \"inline\",\n      \"print\": \"hide\",\n      \"page\": \"flow\",\n      \"forced\": \"none\"",
		"\"flow\": \"inline\",\n      \"page\": \"flow\",\n      \"forced\": \"none\"",
		"пропуск и осознанное «ничего» обязаны различаться"},
	{"CSS относит класс к печати, реестр молчит", "registry", "src/print.css",
		"  .inst-segmented { display: none; }",
		"  .inst-segmented, .inst-timeline { display: none; }",
		"список тихо обрастает тем, о чём никто не решал"},
	{"бесконечная анимация без замедления", "registry", "src/status.css",
		".inst-dot {", ".inst-mut-spin { animation: inst-pulse 1s linear infinite; }\n.inst-dot {",
		"сжатая до 0.01ms бесконечная анимация — это остановка"},
	{"у пункта группы убрали состояние", "registry", "docs/components/actions/segmented.md",
		"<button type=\"button\" role=\"radio\" aria-checked=\"false\" tabindex=\"-1\">Сетка</button>",
		"<button type=\"button\" role=\"radio\" tabindex=\"-1\">Сетка</button>",
		"состояние живёт в атрибуте, а не в классе"},
	{"два бегущих tabindex в одной группе", "registry", "docs/components/actions/segmented.md",
		"<button type=\"button\" role=\"radio\" aria-checked=\"false\" tabindex=\"-1\">Сетка</button>",
		"<button type=\"button\" role=\"radio\" aria-checked=\"false\" tabindex=\"0\">Сетка</button>",
		"Tab пройдёт по каждому пункту, группа перестанет быть одним контролом"},
	{"мера без aria-valuenow", "registry", "docs/components/charts/meter.md",
		"aria-valuenow=\"43\"", "data-mut=\"43\"",
		"роль обещает число вспомогательной технологии"},
	{"роль скрипта не объявлена в реестре", "registry", "components.json",
		"\"role\": \"tablist\",", "\"role\": \"tablist-renamed\",",
		"скрипт обещает клавиатуру от имени кита, а реестр об этом не знает"},
	{"два бегущих tabindex на переведённой странице", "registry", "docs/components/actions/segmented.en.md",
		"",
		"---\ntitle: Segmented\nsource: src/actions.css\n---\n\n```html preview\n<div class=\"inst-segmented\" role=\"radiogroup\" aria-label=\"View\">\n  <button type=\"button\" role=\"radio\" aria-checked=\"true\" tabindex=\"0\">List</button>\n  <button type=\"button\" role=\"radio\" aria-checked=\"false\" tabindex=\"0\">Grid</button>\n</div>\n```\n",
		"перевод пишется заново и ошибается чаще оригинала, а проверялся только оригинал"},
}

func main() {
	verbose := flag.Bool("v", false, "показать вывод упавшего гейта")
	only := flag.String("only", "", "только мутации одного гейта")
	root := flag.String("root", "..", "корень репозитория")
	flag.Parse()

	repo, err := filepath.Abs(*root)
	must(err)

	// Гейты собираются ОДИН раз: тридцать `go run` подряд тратят минуты на
	// повторную компиляцию одного и того же.
	binDir, err := os.MkdirTemp("", "instrument-gates-")
	must(err)
	defer os.RemoveAll(binDir)

	fmt.Println("сборка гейтов…")
	bin := map[string]string{}
	for name, g := range gates {
		out := filepath.Join(binDir, name+exeSuffix())
		cmd := exec.Command("go", "build", "-o", out, g.pkg)
		cmd.Dir = filepath.Join(repo, "tools")
		if b, err := cmd.CombinedOutput(); err != nil {
			fmt.Fprintf(os.Stderr, "не собрать %s: %v\n%s", name, err, b)
			os.Exit(1)
		}
		bin[name] = out
	}

	run := func(name, tree string) (bool, string) {
		cmd := exec.Command(bin[name], gates[name].args(tree)...)
		cmd.Dir = filepath.Join(repo, "tools")
		b, err := cmd.CombinedOutput()
		return err == nil, string(b)
	}

	// ── контроль: нетронутое дерево обязано быть зелёным ───────────────────
	base, err := os.MkdirTemp("", "instrument-base-")
	must(err)
	defer os.RemoveAll(base)
	must(seed(repo, base))

	fmt.Println("контроль на нетронутой копии…")
	for name := range gates {
		if *only != "" && *only != name {
			continue
		}
		ok, out := run(name, base)
		if !ok {
			fmt.Fprintf(os.Stderr, "\nКОНТРОЛЬ ПРОВАЛЕН: %s красный на нетронутом дереве.\n"+
				"Стенд мерил бы собственный шум, а не мутации.\n%s", name, out)
			os.Exit(1)
		}
	}

	// ── мутации ────────────────────────────────────────────────────────────
	fmt.Println()
	width := 0
	for _, m := range mutations {
		if n := len([]rune(m.name)); n > width {
			width = n
		}
	}

	missed, ran := 0, 0
	for _, m := range mutations {
		if *only != "" && *only != m.gate {
			continue
		}
		ran++
		tree, err := os.MkdirTemp("", "instrument-mut-")
		must(err)
		must(seed(repo, tree))

		applied, err := apply(filepath.Join(tree, m.file), m.from, m.to)
		if err != nil || !applied {
			os.RemoveAll(tree)
			fmt.Printf("  %-*s  ✗ МУТАЦИЯ НЕ ПРИМЕНИЛАСЬ (%s)\n", width, m.name, m.file)
			fmt.Printf("      искомого текста в файле нет — стенд разошёлся с китом\n")
			missed++
			continue
		}

		ok, out := run(m.gate, tree)
		os.RemoveAll(tree)

		pad := strings.Repeat(" ", width-len([]rune(m.name)))
		if ok {
			missed++
			fmt.Printf("  %s%s  ✗ ПРОПУЩЕНА  гейтом %s\n", m.name, pad, m.gate)
			fmt.Printf("      %s\n", m.why)
		} else {
			fmt.Printf("  %s%s  · поймана   %s\n", m.name, pad, m.gate)
			if *verbose {
				for _, l := range strings.Split(strings.TrimSpace(out), "\n") {
					fmt.Printf("      %s\n", l)
				}
			}
		}
	}

	fmt.Println()
	if missed > 0 {
		fmt.Printf("── дыр в гейтах: %d из %d ──\n", missed, ran)
		fmt.Println("Пропущенная мутация означает, что инвариант объявлен, но не охраняется.")
		os.Exit(1)
	}
	fmt.Printf("· %d мутаций, все пойманы: каждый инвариант из списка охраняется\n", ran)
}

// seed раскладывает во временном каталоге минимальное дерево, которого хватает
// всем гейтам: исходники, документация, реестр, версия и манифест.
func seed(repo, tree string) error {
	for _, d := range []string{"src", "docs"} {
		if err := copyTree(filepath.Join(repo, d), filepath.Join(tree, d)); err != nil {
			return err
		}
	}
	// Кроме src/ и docs/ страницы ссылаются полем source ещё на два файла —
	// собранный CSS и читалку по пикселям. Без них docscheck честно падает на
	// НЕТРОНУТОЙ копии, и контроль это ловит: слепок, в котором чего-то нет,
	// мерил бы собственную неполноту.
	for _, f := range []string{
		"components.json", "VERSION", "package.json",
		"tools/audit.js", "dist/instrument.min.css",
	} {
		if err := copyFile(filepath.Join(repo, f), filepath.Join(tree, filepath.FromSlash(f))); err != nil {
			return err
		}
	}
	// Стол примера живёт в site/, а docscheck читает его ради собственных
	// свойств демо-сцены. Копируется под своим именем, чтобы не тащить весь
	// каталог сайта.
	return copyFile(filepath.Join(repo, "site", "internal", "render", "assets", "docs.css"),
		filepath.Join(tree, "stage.css"))
}

func apply(path, from, to string) (bool, error) {
	// Пустой from — мутация НОВЫМ файлом: гейт обязан увидеть страницу,
	// которой в дереве не было. Уже существующий файл здесь значит то же,
	// что ненайденный текст замены, — стенд разошёлся с китом, — поэтому
	// это отказ, а не тихая перезапись.
	if from == "" {
		if _, err := os.Stat(path); err == nil {
			return false, nil
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return false, err
		}
		return true, os.WriteFile(path, []byte(to), 0o644)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}
	s := strings.ReplaceAll(string(b), "\r\n", "\n")
	if !strings.Contains(s, from) {
		return false, nil
	}
	return true, os.WriteFile(path, []byte(strings.Replace(s, from, to, 1)), 0o644)
}

func copyTree(src, dst string) error {
	return filepath.WalkDir(src, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, p)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		return copyFile(p, target)
	})
}

func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

func exeSuffix() string {
	if os.Getenv("OS") == "Windows_NT" || filepath.Separator == '\\' {
		return ".exe"
	}
	return ""
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

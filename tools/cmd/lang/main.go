// Команда lang сторожит зоны, из которых кириллица уже ушла.
//
// Переезд на английский идёт шагами, и у каждого шага обязан быть проверяемый
// признак завершённости. Без него «зона переведена» — это память автора: она
// держится до первой правки, сделанной по привычке, и обнаруживается на
// стороне потребителя.
//
// Зоны включаются ПО ОДНОЙ и в КОНЦЕ шага, а не в начале: гейт, заведённый
// красным, живёт с отключённой проверкой и не сторожит ничего. Список зон
// лежит здесь целиком, вместе с невключёнными, — иначе «сколько сделано» и
// «сколько осталось» пришлось бы считать по чужой памяти.
//
// Команда стоит отдельно, а не внутри сборки сайта: сборка видит src/ и docs/,
// а зоны переезда — ещё и dist/, tools/, assets/ и корневые файлы.
//
//	go run ./cmd/lang        проверить включённые зоны
//	go run ./cmd/lang -v     показать заодно и невключённые
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
)

// Как читать файл зоны.
//
// whole — кириллицы нет вообще. strings — кириллицы нет в СТРОКОВЫХ ЛИТЕРАЛАХ,
// комментарии не в счёт.
//
// Второй режим нужен ровно потому, что dist/instrument.js — копия src/kit.js
// со всеми его комментариями, а src/ переводится шагом позже. Ждать шага 2,
// чтобы начать сторожить фразы, которые кит ГОВОРИТ человеку, значило бы
// оставить единственную зону, видимую постороннему, без охраны на пятнадцать
// сеансов. Разделение честное: строка уезжает в aria-label, комментарий не
// уезжает никуда.
type mode int

const (
	whole mode = iota
	literals
)

type zone struct {
	name  string
	paths []string // файлы и каталоги, относительно корня репозитория
	mode  mode
	on    bool
	step  string // шаг переезда, на котором зона включается
	why   string
}

var zones = []zone{
	{
		name:  "публичные строки кита",
		paths: []string{"dist/instrument.js"},
		mode:  literals,
		on:    true,
		step:  "1",
		why:   "уезжают в aria-label и в объявления скринридеру",
	},
	{
		name:  "манифест пакета",
		paths: []string{"package.json"},
		mode:  whole,
		on:    true,
		step:  "1",
		why:   "description виден в реестре npm",
	},
	{
		name:  "знаки",
		paths: []string{"assets"},
		mode:  whole,
		on:    true,
		step:  "1",
		why:   "файлы уезжают в пакет вместе с комментариями внутри",
	},
	// Шаг 2 идёт файлами, а не одним куском: 3 392 кириллические строки —
	// это 15–22 сеанса, и всё это время зона «src целиком» стояла бы
	// выключенной. Работа двадцати сеансов без охраны возвращается по
	// привычке ровно так же, как работа одного.
	//
	// Поэтому зона растёт: переведённый файл попадает в список в том же
	// коммите, которым переведён. Когда в списке окажутся все двадцать, он
	// схлопнется в "src" — и это будет признак завершённости шага, а не
	// объявление о нём.
	{
		name: "исходники кита: переведённые файлы",
		paths: []string{
			"src/kit.css",
			"src/forced.css",
			"src/motion.css",
			"src/rows.css",
			"src/tone.css",
			"src/print.css",
			"src/feedback.css",
			"src/text.css",
			"src/status.css",
			"src/base.css",
			"src/surfaces.css",
			"src/table.css",
		},
		mode: whole,
		on:   true,
		step: "2",
		why:  "переведённый файл не возвращается к русскому по привычке",
	},
	{
		name:  "исходники кита целиком",
		paths: []string{"src", "dist/instrument.js", "dist/instrument.css"},
		mode:  whole,
		step:  "2",
		why:   "src/ уезжает в пакет полем files",
	},
	{
		name:  "инструменты",
		paths: []string{"tools"},
		mode:  whole,
		step:  "3",
		why:   "вывод гейтов живёт в тех же строках, что и код",
	},
	{
		name:  "сайт",
		paths: []string{"site"},
		mode:  whole,
		step:  "3",
		why:   "",
	},
	{
		name:  "документация",
		paths: []string{"docs/start", "docs/foundations", "docs/components", "docs/agent", "docs/layout", "docs/blocks", "docs/about"},
		mode:  whole,
		step:  "4",
		why:   "русский переезжает в *.ru.md разворотом базового языка",
	},
	{
		name:  "корень",
		paths: []string{"README.md", "CONTRIBUTING.md", "ROADMAP.md"},
		mode:  whole,
		step:  "6",
		why:   "CHANGELOG.md не переводится: переведённая летопись — переписанная",
	},
}

func main() {
	root := flag.String("root", "..", "корень репозитория")
	verbose := flag.Bool("v", false, "показать невключённые зоны")
	flag.Parse()

	var problems []string
	onCount := 0
	for _, z := range zones {
		if !z.on {
			continue
		}
		onCount++
		for _, p := range z.paths {
			problems = append(problems, z.scan(filepath.Join(*root, filepath.FromSlash(p)), *root)...)
		}
	}
	sort.Strings(problems)

	if len(problems) > 0 {
		fmt.Printf("── кириллица в зоне, из которой она ушла (%d) ──\n", len(problems))
		for _, p := range problems {
			fmt.Println("  ·", p)
		}
		fmt.Println()
		fmt.Println("Зона объявлена переведённой, и признак завершённости шага держится на")
		fmt.Println("этой проверке. Верните строку на английский или снимите зону осознанно.")
		os.Exit(1)
	}

	fmt.Printf("· зон под охраной: %d из %d — кириллицы в них нет\n", onCount, len(zones))
	if *verbose {
		for _, z := range zones {
			if z.on {
				continue
			}
			fmt.Printf("    ждёт шага %s: %s (%s)\n", z.step, z.name, strings.Join(z.paths, " · "))
		}
	}
}

// scan обходит файл или каталог зоны.
func (z zone) scan(path, root string) []string {
	var out []string
	filepath.WalkDir(path, func(p string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return nil
		}
		rel := p
		if r, err := filepath.Rel(root, p); err == nil {
			rel = filepath.ToSlash(r)
		}
		text := strings.ReplaceAll(string(b), "\r\n", "\n")
		if z.mode == literals {
			text = onlyJSLiterals(text)
		}
		for i, line := range strings.Split(text, "\n") {
			if w := cyrillicIn(line); w != "" {
				out = append(out, fmt.Sprintf("%s:%d: %s — зона «%s»", rel, i+1, w, z.name))
			}
		}
		return nil
	})
	return out
}

// cyrillicIn возвращает первое кириллическое слово строки — сообщение обязано
// показывать НАЙДЕННОЕ, а не только координату: по слову видно сразу, забытая
// это строка или намеренная.
func cyrillicIn(line string) string {
	start := -1
	for i, r := range line {
		if unicode.Is(unicode.Cyrillic, r) {
			if start < 0 {
				start = i
			}
			continue
		}
		if start >= 0 {
			return "«" + line[start:i] + "»"
		}
	}
	if start >= 0 {
		return "«" + line[start:] + "»"
	}
	return ""
}

// onlyJSLiterals гасит всё, кроме строковых литералов, сохраняя номера строк.
//
// Разбор посимвольный, а не регулярным выражением: у комментария и литерала
// общие символы, и регулярка, не помнящая, где она находится, принимает
// апостроф внутри комментария за начало строки и глотает половину файла.
// Подстановка ${...} внутри шаблона не разбирается — код внутри неё
// кириллицы не содержит по тем же правилам, что и код снаружи. Регулярное
// выражение литералом не считается: человек его не слышит, а шагом 2 зона
// всё равно переходит в режим whole, где не считается уже ничего.
func onlyJSLiterals(src string) string {
	var out strings.Builder
	out.Grow(len(src))

	keep := func(r rune) {
		if r == '\n' {
			out.WriteRune('\n')
			return
		}
		out.WriteRune(r)
	}
	drop := func(r rune) {
		if r == '\n' {
			out.WriteRune('\n')
			return
		}
		out.WriteRune(' ')
	}

	rs := []rune(src)
	for i := 0; i < len(rs); i++ {
		r := rs[i]
		switch {
		case r == '/' && i+1 < len(rs) && rs[i+1] == '/':
			for ; i < len(rs) && rs[i] != '\n'; i++ {
				drop(rs[i])
			}
			i--
		case r == '/' && i+1 < len(rs) && rs[i+1] == '*':
			drop(rs[i])
			i++
			for ; i < len(rs); i++ {
				drop(rs[i])
				if rs[i] == '/' && rs[i-1] == '*' {
					break
				}
			}
		case r == '\'' || r == '"' || r == '`':
			quote := r
			drop(rs[i])
			i++
			for ; i < len(rs); i++ {
				if rs[i] == '\\' && i+1 < len(rs) {
					keep(rs[i])
					i++
					keep(rs[i])
					continue
				}
				if rs[i] == quote {
					drop(rs[i])
					break
				}
				// Незакрытая кавычка до конца строки — это не литерал, а
				// апостроф в тексте: делить регулярное выражение и деление
				// эта команда не умеет и не обязана.
				if rs[i] == '\n' && quote != '`' {
					drop(rs[i])
					break
				}
				keep(rs[i])
			}
		default:
			drop(r)
		}
	}
	return out.String()
}

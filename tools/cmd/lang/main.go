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

type zone struct {
	name  string
	paths []string // файлы и каталоги, относительно корня репозитория
	on    bool
	step  string // шаг переезда, на котором зона включается
	why   string
}

var zones = []zone{
	{
		name:  "манифест пакета",
		paths: []string{"package.json"},
		on:    true,
		step:  "1",
		why:   "description виден в реестре npm",
	},
	{
		name:  "знаки",
		paths: []string{"assets"},
		on:    true,
		step:  "1",
		why:   "файлы уезжают в пакет вместе с комментариями внутри",
	},
	// Шаг 2 шёл файлами, а не одним куском: 3 392 кириллические строки — это
	// 15–22 сеанса, и всё это время зона «src целиком» стояла бы выключенной.
	// Работа двадцати сеансов без охраны возвращается по привычке ровно так
	// же, как работа одного, поэтому зона росла: переведённый файл попадал в
	// список в том же коммите, которым переведён.
	//
	// Двадцатый файл — tokens.css — закрыл список, и он схлопнут в "src".
	// Схлопывание и есть признак завершённости шага: пока в зоне стояли
	// имена, недостающее имя было видно; теперь недостающего нет.
	{
		name:  "исходники кита",
		paths: []string{"src", "dist/instrument.js", "dist/instrument.css"},
		on:    true,
		step:  "2",
		why:   "src/ уезжает в пакет полем files",
	},
	{
		name:  "инструменты",
		paths: []string{"tools"},
		step:  "3",
		why:   "вывод гейтов живёт в тех же строках, что и код",
	},
	{
		name:  "сайт",
		paths: []string{"site"},
		step:  "3",
		why:   "",
	},
	{
		name:  "документация",
		paths: []string{"docs/start", "docs/foundations", "docs/components", "docs/agent", "docs/layout", "docs/blocks", "docs/about"},
		step:  "4",
		why:   "русский переезжает в *.ru.md разворотом базового языка",
	},
	{
		name:  "корень",
		paths: []string{"README.md", "CONTRIBUTING.md", "ROADMAP.md"},
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

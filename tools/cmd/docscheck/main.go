// Команда docscheck сверяет документацию с китом. Обе стороны разом.
//
// Документация о ките расходится с китом в двух направлениях, и оба
// одинаково дороги:
//
//	вперёд   в примере стоит класс, которого в ките нет — читатель копирует
//	         разметку, она молча не работает, и виноват оказывается кит;
//	назад    класс в ките есть, а страницы у него нет — компонент существует
//	         и не находится, что для потребителя равно его отсутствию.
//
// Проверка читает настоящие исходники кита и настоящие страницы, поэтому
// разойтись с ними не может — тот же принцип, что у contrast.
//
//	go run ./cmd/docscheck        коротко
//	go run ./cmd/docscheck -v     со списком неописанных классов
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var (
	commentRe = regexp.MustCompile(`(?s)/\*.*?\*/`)
	urlRe     = regexp.MustCompile(`url\([^)]*\)`)
	dqRe      = regexp.MustCompile(`"(?:[^"\\]|\\.)*"`)
	sqRe      = regexp.MustCompile(`'(?:[^'\\]|\\.)*'`)

	classRe = regexp.MustCompile(`\.(-?[A-Za-z_][\w-]*)`)
	declRe  = regexp.MustCompile(`(--[a-z][\w-]*)\s*:`)
	varUse  = regexp.MustCompile(`var\(\s*(--[a-z][\w-]*)`)
	dataSel = regexp.MustCompile(`\[data-([a-z-]+)="([^"]+)"\]`)

	instRe  = regexp.MustCompile(`\binst-[a-z0-9-]+`)
	linkRe  = regexp.MustCompile(`\]\((\.[^)#]+\.md)(#[^)]*)?\)`)
	tokRe   = regexp.MustCompile(`--[a-z][\w-]*(?:/[a-z0-9-]+)*`)
	dataAtt = regexp.MustCompile(`data-([a-z-]+)="([^"]*)"`)
)

// strip убирает то, что не является селектором: комментарии, строки, url().
//
// Наивный поиск по файлу даёт мусор: .w3 и .org приезжают из
// xmlns='http://www.w3.org/2000/svg' внутри data-URI, а .tokens и .base — из
// имён файлов в комментариях.
func strip(s string) string {
	s = commentRe.ReplaceAllString(s, " ")
	s = urlRe.ReplaceAllString(s, " ")
	s = dqRe.ReplaceAllString(s, " ")
	return sqRe.ReplaceAllString(s, " ")
}

func main() {
	var (
		srcDir  = flag.String("kit", "../src", "каталог кита")
		docsDir = flag.String("docs", "../docs", "каталог документации")
		stage   = flag.String("stage", "../site/internal/render/assets/demo.css", "стили стола примера")
		verbose = flag.Bool("v", false, "показать список классов без страницы")
	)
	flag.Parse()

	// ── Что есть в ките ──────────────────────────────────────────────────
	kit := map[string]bool{}
	var allCSS strings.Builder
	entries, err := os.ReadDir(*srcDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "не прочитать кит:", err)
		os.Exit(1)
	}
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".css") {
			continue
		}
		b, err := os.ReadFile(filepath.Join(*srcDir, e.Name()))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		allCSS.Write(b)
		allCSS.WriteByte('\n')
		for _, line := range strings.Split(strip(string(b)), "\n") {
			i := strings.Index(line, "{")
			if i < 0 {
				continue
			}
			for _, m := range classRe.FindAllStringSubmatch(line[:i], -1) {
				kit[m[1]] = true
			}
		}
	}
	css := allCSS.String()

	// Токены. Выдуманный токен так же нем, как выдуманный класс: var() от
	// несуществующего имени молча отдаёт пустоту, и страница врёт беззвучно.
	tokens := map[string]bool{}
	for _, m := range declRe.FindAllStringSubmatch(css, -1) {
		tokens[m[1]] = true
	}
	// Переменные, которые кит только ЧИТАЕТ, а задаёт разметка: глубина узла
	// дерева, доля кольца, номер ряда. Объявления у них нет по замыслу.
	for _, m := range varUse.FindAllStringSubmatch(css, -1) {
		tokens[m[1]] = true
	}
	// Стол примера имеет собственные свойства — --c у образца цвета. Они не
	// токены кита и никогда ими не станут, но и выдумкой не являются.
	if b, err := os.ReadFile(*stage); err == nil {
		for _, m := range declRe.FindAllStringSubmatch(string(b), -1) {
			tokens[m[1]] = true
		}
		for _, m := range varUse.FindAllStringSubmatch(string(b), -1) {
			tokens[m[1]] = true
		}
	}

	// Хвосты модификаторов: из inst-btn--danger получается danger. Нужны,
	// чтобы не принимать «вариант --danger» в прозе за несуществующий токен.
	modifiers := map[string]bool{}
	for c := range kit {
		if i := strings.Index(c, "--"); i > 0 {
			modifiers[c[i+2:]] = true
		}
	}

	// Словари data-атрибутов — из селекторов, а не из головы.
	vocab := map[string]map[string]bool{}
	for _, m := range dataSel.FindAllStringSubmatch(css, -1) {
		if vocab[m[1]] == nil {
			vocab[m[1]] = map[string]bool{}
		}
		vocab[m[1]][m[2]] = true
	}
	// Базовые значения оформления не имеют, поэтому в селекторах их нет. Они
	// объявлены в конституции — без них проверка ругалась бы на правильную
	// разметку.
	for _, v := range []string{"queued", "todo", "approved"} {
		if vocab["state"] == nil {
			vocab["state"] = map[string]bool{}
		}
		vocab["state"][v] = true
	}

	// ── Что есть в документации ──────────────────────────────────────────
	var pages []string
	filepath.WalkDir(*docsDir, func(p string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if strings.Contains(filepath.ToSlash(p), "/internal/") || !strings.HasSuffix(p, ".md") {
			return nil
		}
		pages = append(pages, p)
		return nil
	})
	sort.Strings(pages)

	// Пути печатаются от корня репозитория, а не от рабочего каталога:
	// «docs/button.md» одинаково понятен и в выводе, и в редакторе, откуда
	// бы проверку ни запустили.
	root := filepath.Dir(*docsDir)
	rootRel := func(p string) string {
		if r, err := filepath.Rel(root, p); err == nil {
			return filepath.ToSlash(r)
		}
		return filepath.ToSlash(p)
	}

	documented := map[string]bool{}
	var problems []string
	pending := map[string]int{}

	skipAttr := map[string]bool{"theme": true, "density": true, "dir": true}

	for _, p := range pages {
		b, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		rel := rootRel(p)
		for i, line := range strings.Split(string(b), "\n") {
			at := fmt.Sprintf("%s:%d", rel, i+1)

			for _, m := range instRe.FindAllString(line, -1) {
				documented[m] = true
				if !kit[m] {
					problems = append(problems, at+"  класса нет в ките: ."+m)
				}
			}

			// Ссылки между страницами. Не «ошибка» и не «метрика», а третий
			// вид: ссылка на ещё не написанную страницу — нормальный след
			// недоделанной работы, но забытая это дыра, которую никто не
			// заметит, пока читатель в неё не провалится.
			for _, m := range linkRe.FindAllStringSubmatch(line, -1) {
				target := filepath.Join(filepath.Dir(p), filepath.FromSlash(m[1]))
				if _, err := os.Stat(target); err != nil {
					pending[rootRel(target)]++
				}
			}

			// Токены. Сокращения вида --text-xs/sm/md раскрываются: это
			// принятая на страницах форма записи, а не опечатка.
			for _, raw := range tokRe.FindAllString(line, -1) {
				if j := strings.Index(line, raw); j > 0 {
					prev := line[j-1]
					if prev == '-' || prev == '_' || (prev >= 'a' && prev <= 'z') ||
						(prev >= 'A' && prev <= 'Z') || (prev >= '0' && prev <= '9') {
						continue // часть имени класса: inst-btn--primary
					}
				}
				parts := strings.Split(raw, "/")
				head := parts[0]
				// Две законные формы, которые токенами не являются:
				// модификатор класса в прозе («вариант --danger») и
				// сокращение семейства с оборванным хвостом («--space-*»).
				if strings.HasSuffix(head, "-") || modifiers[strings.TrimPrefix(head, "--")] {
					continue
				}
				stem := head
				if k := strings.LastIndex(head, "-"); k > 1 {
					stem = head[:k]
				}
				names := []string{head}
				for _, t := range parts[1:] {
					names = append(names, stem+"-"+t)
				}
				for _, n := range names {
					if !tokens[n] {
						problems = append(problems, at+"  токена нет в ките: "+n)
					}
				}
			}

			for _, m := range dataAtt.FindAllStringSubmatch(line, -1) {
				attr, val := m[1], m[2]
				if skipAttr[attr] {
					continue
				}
				known, ok := vocab[attr]
				if !ok {
					problems = append(problems, at+"  неизвестный атрибут: data-"+attr)
					continue
				}
				if val != "" && !known[val] {
					var list []string
					for v := range known {
						list = append(list, v)
					}
					sort.Strings(list)
					problems = append(problems,
						fmt.Sprintf("%s  data-%s=%q — нет в словаре (%s)", at, attr, val, strings.Join(list, " ")))
				}
			}
		}
	}

	// ── Отчёт ────────────────────────────────────────────────────────────
	var undocumented []string
	for c := range kit {
		if !documented[c] {
			undocumented = append(undocumented, c)
		}
	}
	sort.Strings(undocumented)

	covered := len(kit) - len(undocumented)
	pct := 0
	if len(kit) > 0 {
		pct = covered * 100 / len(kit)
	}
	fmt.Printf("страниц: %d  ·  классов в ките: %d  ·  покрыто: %d (%d%%)\n\n",
		len(pages), len(kit), covered, pct)

	if len(problems) > 0 {
		seen := map[string]bool{}
		var uniq []string
		for _, p := range problems {
			if !seen[p] {
				seen[p] = true
				uniq = append(uniq, p)
			}
		}
		fmt.Printf("── документация ссылается на несуществующее (%d) ──\n", len(uniq))
		for _, p := range uniq {
			fmt.Println("  " + p)
		}
		fmt.Println()
	}

	if len(pending) > 0 {
		var keys []string
		for k := range pending {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		fmt.Printf("── ссылки на ещё не написанные страницы (%d) ──\n", len(keys))
		for _, k := range keys {
			fmt.Printf("  %s  ← %d\n", k, pending[k])
		}
		fmt.Println()
	}

	if len(undocumented) > 0 {
		if *verbose {
			fmt.Printf("── без страницы (%d) ──\n", len(undocumented))
			for i, c := range undocumented {
				undocumented[i] = "." + c
			}
			fmt.Println("  " + strings.Join(undocumented, " ") + "\n")
		} else {
			fmt.Printf("── без страницы: %d классов (список — с ключом -v) ──\n\n", len(undocumented))
		}
	}

	if len(problems) == 0 && len(undocumented) == 0 && len(pending) == 0 {
		fmt.Println("· документация и кит сходятся полностью")
	}

	if len(problems) > 0 {
		os.Exit(1)
	}
}

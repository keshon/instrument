// Генератор сайта документации.
//
// Сайт собран САМИМ КИТОМ: раскладка — .inst-shell, навигация — .inst-nav,
// проза — .inst-prose. Это не поза. Кит для дашбордов, показанный внутри
// чужой дизайн-системы, ничего не доказывает; кит, из которого собран его
// собственный сайт, доказывает всё сразу и ломается на глазах, если врёт.
//
// Отсюда же отсутствие зависимостей тяжелее markdown-парсера: проект весь
// про «без сборки, на платформе, честно», и 200 МБ node_modules ради показа
// CSS-кита противоречили бы его характеру.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Page struct {
	Route   string // /components/actions/button/
	Rel     string // components/actions/button.md
	Dir     string // components/actions
	Slug    string // button
	Title   string
	Group   string
	Source  string
	NeedsJS string
	Status  string
	Splash  bool
	HTML    string
	TOC     []Heading
	Text    string // очищенный текст для поиска
}

type Heading struct {
	Level int
	ID    string
	Text  string
}

func main() {
	var (
		docs = flag.String("docs", "../docs", "каталог с исходниками страниц")
		out  = flag.String("out", "../site", "каталог вывода")
		src  = flag.String("kit", "../src", "каталог кита")
		asts = flag.String("assets", "../assets", "каталог ресурсов кита")
		serve = flag.String("serve", "", "после сборки поднять сервер на этом адресе, например :4321")
	)
	flag.Parse()

	pages, err := collect(*docs)
	if err != nil {
		log.Fatalf("сбор страниц: %v", err)
	}
	if len(pages) == 0 {
		log.Fatalf("в %s не найдено ни одной страницы", *docs)
	}

	nav := buildNav(pages)

	if problems := verify(pages); len(problems) > 0 {
		for _, x := range problems {
			fmt.Fprintln(os.Stderr, "  "+x)
		}
		log.Fatalf("сборка остановлена: %d проблем", len(problems))
	}

	if err := os.RemoveAll(*out); err != nil {
		log.Fatalf("очистка вывода: %v", err)
	}
	if err := writeSite(*out, pages, nav); err != nil {
		log.Fatalf("сборка: %v", err)
	}
	if err := copyKit(*src, *asts, *out); err != nil {
		log.Fatalf("копирование кита: %v", err)
	}
	if err := writeSearchIndex(*out, pages); err != nil {
		log.Fatalf("индекс поиска: %v", err)
	}

	demos := 0
	for _, p := range pages {
		demos += strings.Count(p.HTML, `class="demo-stage"`)
	}
	fmt.Printf("страниц: %d  ·  живых примеров: %d  ·  разделов навигации: %d\n",
		len(pages), demos, len(nav))

	if *serve != "" {
		fmt.Printf("сервер: http://localhost%s\n", *serve)
		log.Fatal(http.ListenAndServe(*serve, http.FileServer(http.Dir(*out))))
	}
}

// collect обходит каталог документации и разбирает каждую страницу.
// internal/ пропускается намеренно: это процессные документы, а не
// документация, и они прямо об этом говорят у себя в шапке.
func collect(root string) ([]*Page, error) {
	var pages []*Page
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == "internal" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".md") {
			return nil
		}
		rel, _ := filepath.Rel(root, path)
		rel = filepath.ToSlash(rel)
		if rel == "README.md" {
			return nil
		}
		p, err := parsePage(path, rel)
		if err != nil {
			return fmt.Errorf("%s: %w", rel, err)
		}
		pages = append(pages, p)
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(pages, func(i, j int) bool { return pages[i].Rel < pages[j].Rel })
	return pages, nil
}

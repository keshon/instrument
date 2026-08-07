// Команда site собирает документацию instrument.
//
// Сайт собран САМИМ КИТОМ: раскладка — inst-shell, навигация — inst-nav,
// проза — inst-prose. Это не поза: кит для дашбордов, показанный внутри
// чужой дизайн-системы, ничего не доказывает, а кит, из которого собран его
// собственный сайт, ломается на глазах, если врёт.
//
//	go run ./cmd/site                 собрать в ../dist
//	go run ./cmd/site -serve :4321    собрать и поднять сервер
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"instrument/site/internal/check"
	"instrument/site/internal/content"
	"instrument/site/internal/nav"
	"instrument/site/internal/render"
)

func main() {
	var (
		docs   = flag.String("docs", "../docs", "каталог с исходниками страниц")
		out    = flag.String("out", "dist", "каталог вывода")
		kit    = flag.String("kit", "../src", "каталог кита")
		assets = flag.String("assets", "../assets", "каталог ресурсов кита")
		serve  = flag.String("serve", "", "поднять сервер после сборки, например :4321")
	)
	flag.Parse()

	pages, err := content.Collect(*docs)
	if err != nil {
		log.Fatalf("сбор страниц: %v", err)
	}
	if len(pages) == 0 {
		log.Fatalf("в %s не найдено ни одной страницы", *docs)
	}

	// Иконка выводится из слага, поэтому новая страница получает ссылку на
	// символ, которого может не быть. Пустой <use> ничего не рисует и ничего
	// не сообщает — ловим здесь.
	sprite, err := os.ReadFile(filepath.Join(*assets, "sprite.svg"))
	if err != nil {
		log.Fatalf("не прочитать спрайт: %v", err)
	}
	var missing []string
	for _, p := range pages {
		if p.Icon != "" && !strings.Contains(string(sprite), `id="`+p.Icon+`"`) {
			missing = append(missing, fmt.Sprintf("%s  нет символа %s в спрайте", p.Route, p.Icon))
		}
	}

	// Проверка идёт ДО записи: страница с битой ссылкой выглядит целой, и
	// увидеть её глазами нельзя.
	if problems := append(check.Verify(pages), missing...); len(problems) > 0 {
		for _, p := range problems {
			fmt.Fprintln(os.Stderr, "  "+p)
		}
		log.Fatalf("сборка остановлена: %d проблем", len(problems))
	}

	sections := nav.Build(pages)
	if err := render.Site(pages, sections, render.Options{
		Out: *out, Kit: *kit, Assets: *assets,
	}); err != nil {
		log.Fatalf("сборка: %v", err)
	}

	demos := 0
	for _, p := range pages {
		demos += len(p.Demos)
	}
	fmt.Printf("страниц: %d  ·  живых примеров: %d  ·  разделов навигации: %d\n",
		len(pages), demos, len(sections))

	if *serve != "" {
		fmt.Printf("сервер: http://localhost%s\n", *serve)
		log.Fatal(http.ListenAndServe(*serve, http.FileServer(http.Dir(*out))))
	}
}

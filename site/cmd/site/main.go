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
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"instrument/site/internal/check"
	"instrument/site/internal/content"
	"instrument/site/internal/i18n"
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

		// Замечания контракта — это план работ на сорок четыре страницы, и
		// печатать его при каждой сборке значит утопить в нём настоящие
		// ошибки. Печатается счётчик; список — по требованию.
		verbose = flag.Bool("contract", false, "печатать замечания контракта по непереносенным страницам")
	)
	flag.Parse()

	byLang, err := content.Collect(*docs)
	if err != nil {
		log.Fatalf("сбор страниц: %v", err)
	}
	pages := content.Flat(byLang)
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

	// Значения токенов приходят из самого кита, а не из страниц. Написанное
	// руками разошлось бы с tokens.css при первой же правке.
	tokens, err := content.TokenValues(*kit)
	if err != nil {
		log.Fatalf("не прочитать токены кита: %v", err)
	}
	content.ResolveTokens(pages, tokens)

	// Разметка рендерится ПОСЛЕ того, как справочник достроен.
	if err := content.Render(pages); err != nil {
		log.Fatalf("разметка: %v", err)
	}

	// Проверка идёт ДО записи: страница с битой ссылкой выглядит целой, и
	// увидеть её глазами нельзя.
	// Контракт страницы. Ошибки — только у страниц, которые его объявили;
	// остальным идёт список того, что помешает переносу, и он не роняет
	// сборку: волна без своего критерия готовности не засчитывается, но и не
	// блокирует соседние.
	contractErrs, contractWarns := check.Contract(pages)

	// CSS самого сайта — против кита. Проверка появилась после того, как в
	// docs.css обнаружилась ссылка на несуществующий --tracking-wide рядом с
	// text-transform, которого в ките нет ни разу: страница выглядела целой,
	// а из задуманного работала половина.
	styles, err := render.Stylesheets()
	if err != nil {
		log.Fatalf("не прочитать стили сайта: %v", err)
	}
	assetErrs := check.Assets(styles, tokens)

	// Комментарии — и кита, и сайта. Правило одно: объясняем, почему так
	// сейчас, а не что было раньше.
	sources := map[string]string{}
	for n, s := range styles {
		sources["site/"+n] = s
	}
	for _, dir := range []string{*kit, "internal"} {
		if err := collectSources(sources, dir); err != nil {
			log.Fatalf("не прочитать исходники: %v", err)
		}
	}
	assetErrs = append(assetErrs, check.Comments(sources)...)

	problems := check.Verify(pages, string(sprite))
	problems = append(problems, missing...)
	problems = append(problems, contractErrs...)
	problems = append(problems, assetErrs...)
	if len(problems) > 0 {
		for _, p := range problems {
			fmt.Fprintln(os.Stderr, "  "+p)
		}
		log.Fatalf("сборка остановлена: %d проблем", len(problems))
	}

	sections := map[i18n.Lang][]nav.Section{}
	for lang, ps := range byLang {
		sections[lang] = nav.Build(lang, ps)
	}
	if err := render.Site(byLang, sections, render.Options{
		Out: *out, Kit: *kit, Assets: *assets,
	}); err != nil {
		log.Fatalf("сборка: %v", err)
	}

	demos := map[string]bool{}
	for _, p := range pages {
		for _, d := range p.Demos {
			demos[d.ID] = true
		}
	}
	translated := 0
	for _, p := range byLang[i18n.EN] {
		if p.Translated {
			translated++
		}
	}
	fmt.Printf("страниц: %d на язык  ·  живых примеров: %d  ·  разделов навигации: %d\n",
		len(byLang[i18n.RU]), len(demos), len(sections[i18n.RU]))
	onContract := 0
	for _, p := range byLang[i18n.RU] {
		if p.Layout == "component" {
			onContract++
		}
	}
	fmt.Printf("языков: %d  ·  переведено на английский: %d из %d\n",
		len(i18n.All), translated, len(byLang[i18n.EN]))
	fmt.Printf("под контрактом страницы: %d  ·  ждут переноса: %d замечаний\n",
		onContract, len(contractWarns))
	if *verbose {
		for _, w := range contractWarns {
			fmt.Fprintln(os.Stderr, "  ~ "+w)
		}
	}

	if *serve != "" {
		fmt.Printf("сервер: http://localhost%s\n", *serve)
		log.Fatal(http.ListenAndServe(*serve, http.FileServer(http.Dir(*out))))
	}
}

// collectSources собирает .css, .js и .go для проверки комментариев.
func collectSources(out map[string]string, dir string) error {
	return filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		switch filepath.Ext(p) {
		case ".css", ".js", ".go":
		default:
			return nil
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		out[filepath.ToSlash(p)] = string(b)
		return nil
	})
}

// go run ./cmd/site                 собрать в ../dist
// go run ./cmd/site -serve :4321    собрать и поднять сервер
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
		public = flag.String("public", "public", "каталог, копируемый в корень вывода как есть")
		serve  = flag.String("serve", "", "поднять сервер после сборки, например :4321")

		registry = flag.String("registry", "../components.json", "реестр компонентов: связи для раздела «Связанное»")

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

	tokens, err := content.TokenValues(*kit)
	if err != nil {
		log.Fatalf("не прочитать токены кита: %v", err)
	}
	content.ResolveTokens(pages, tokens)

	content.SetSprite(string(sprite))
	rel, err := content.LoadRelations(*registry)
	if err != nil {
		log.Fatalf("не прочитать реестр: %v", err)
	}
	content.SetRelations(rel)
	if err := content.Render(pages); err != nil {
		log.Fatalf("разметка: %v", err)
	}

	contractErrs, contractWarns := check.Contract(pages)

	styles, err := render.Stylesheets()
	if err != nil {
		log.Fatalf("не прочитать стили сайта: %v", err)
	}
	assetErrs := check.Assets(styles, tokens)

	sources := map[string]string{}
	for n, s := range styles {
		sources["site/"+n] = s
	}
	// Правило комментария одно на весь репозиторий, значит и зона у гейта
	// одна. Пока сюда не входили `tools` и `cmd`, правило держалось на
	// внимательности ровно там, где живут сами проверки: 4 384 строки Go и
	// две команды на JS оставались снаружи, и «Раньше все три размера брали
	// --radius-md» проходило молча.
	//
	// Каталог инструментов ищется рядом с китом, а не задаётся флагом: он
	// лежит на одном уровне с `src`, и второй флаг разошёлся бы с первым.
	tools := filepath.Join(filepath.Dir(*kit), "tools")
	for _, dir := range []string{*kit, "internal", "cmd", tools} {
		// Отсутствующий каталог — не ошибка: сайт собирается из своего
		// модуля, и снаружи его дерева может не оказаться ничего.
		if _, err := os.Stat(dir); err != nil {
			continue
		}
		if err := collectSources(sources, dir); err != nil {
			log.Fatalf("не прочитать исходники: %v", err)
		}
	}
	assetErrs = append(assetErrs, check.Comments(sources)...)
	assetErrs = append(assetErrs, check.StrayCommentEnd(sources)...)

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
		Out: *out, Kit: *kit, Assets: *assets, Public: *public,
	}); err != nil {
		log.Fatalf("сборка: %v", err)
	}

	demos := map[string]bool{}
	for _, p := range pages {
		for _, d := range p.Demos {
			demos[d.ID] = true
		}
	}
	fmt.Printf("страниц: %d  ·  живых примеров: %d  ·  разделов навигации: %d\n",
		len(byLang[i18n.RU]), len(demos), len(sections[i18n.RU]))
	onContract := 0
	for _, p := range byLang[i18n.RU] {
		if p.Layout == "component" {
			onContract++
		}
	}
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

func collectSources(out map[string]string, dir string) error {
	return filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		// `.mjs` в списке потому, что бегунок пиксельной проверки написан
		// именно так, и без этого расширения он был бы единственным файлом
		// репозитория, где комментарий не сторожит никто.
		switch filepath.Ext(p) {
		case ".css", ".js", ".mjs", ".go":
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

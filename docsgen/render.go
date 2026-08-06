package main

import (
	"embed"
	"encoding/json"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
)

//go:embed templates/*.html assets/*
var embedded embed.FS

type pageData struct {
	Page    *Page
	Nav     []NavSection
	Sprite  template.HTML
	Body    template.HTML
	Version string
}

func writeSite(out string, pages []*Page, nav []NavSection) error {
	tpl, err := template.New("page.html").Funcs(template.FuncMap{
		"eq3": func(a, b string) bool { return a == b },
	}).ParseFS(embedded, "templates/*.html")
	if err != nil {
		return err
	}

	sprite, err := os.ReadFile(filepath.Join("..", "assets", "sprite.svg"))
	if err != nil {
		return err
	}

	for _, p := range pages {
		dir := filepath.Join(out, filepath.FromSlash(strings.Trim(p.Route, "/")))
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
		f, err := os.Create(filepath.Join(dir, "index.html"))
		if err != nil {
			return err
		}
		err = tpl.Execute(f, pageData{
			Page:   p,
			Nav:    nav,
			Sprite: template.HTML(sprite),
			Body:   template.HTML(p.HTML),
		})
		f.Close()
		if err != nil {
			return err
		}
	}

	// Собственные стили и скрипт сайта.
	for _, name := range []string{"docs.css", "docs.js"} {
		b, err := embedded.ReadFile("assets/" + name)
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(out, name), b, 0o644); err != nil {
			return err
		}
	}
	return nil
}

// copyKit кладёт кит рядом с сайтом КАК ЕСТЬ, без обработки.
// Сайт обязан показывать ровно тот CSS, который получит потребитель, —
// пропущенный через сборщик это уже другой файл.
func copyKit(src, assets, out string) error {
	for _, pair := range [][2]string{{src, "kit"}, {assets, "assets"}} {
		dst := filepath.Join(out, pair[1])
		if err := os.MkdirAll(dst, 0o755); err != nil {
			return err
		}
		entries, err := os.ReadDir(pair[0])
		if err != nil {
			return err
		}
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			b, err := os.ReadFile(filepath.Join(pair[0], e.Name()))
			if err != nil {
				return err
			}
			if err := os.WriteFile(filepath.Join(dst, e.Name()), b, 0o644); err != nil {
				return err
			}
		}
	}
	return nil
}

type searchDoc struct {
	R string `json:"r"` // маршрут
	T string `json:"t"` // заголовок
	G string `json:"g"` // группа
	B string `json:"b"` // текст
}

// writeSearchIndex пишет один JSON на весь сайт.
//
// Готовый поисковик (Pagefind и подобные) тянет бинарник и WASM на каждую
// страницу. На 73 страницах это дороже самого индекса: весь текст кита
// помещается в несколько сотен килобайт, а поиск по нему — в полсотни строк.
func writeSearchIndex(out string, pages []*Page) error {
	docs := make([]searchDoc, 0, len(pages))
	for _, p := range pages {
		body := p.Text
		if len(body) > 1200 {
			body = body[:1200]
		}
		docs = append(docs, searchDoc{R: p.Route, T: p.Title, G: p.Group, B: body})
	}
	f, err := os.Create(filepath.Join(out, "search.json"))
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(docs)
}

var _ io.Writer = (*os.File)(nil)

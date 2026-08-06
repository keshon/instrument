// Package render пишет собранный сайт на диск.
package render

import (
	"embed"
	"encoding/json"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"instrument/site/internal/content"
	"instrument/site/internal/nav"
)

//go:embed templates/*.html assets/*
var files embed.FS

type Options struct {
	Out    string // куда писать
	Kit    string // каталог кита (src)
	Assets string // каталог ресурсов кита
}

type pageData struct {
	Page   *content.Page
	Nav    []nav.Section
	Sprite template.HTML
	Body   template.HTML
}

type demoData struct {
	Sprite template.HTML
	Markup template.HTML
	Title  string
}

func Site(pages []*content.Page, sections []nav.Section, o Options) error {
	tpl, err := template.New("").Funcs(template.FuncMap{
		"same": func(a, b string) bool { return a == b },
	}).ParseFS(files, "templates/*.html")
	if err != nil {
		return err
	}

	sprite, err := os.ReadFile(filepath.Join(o.Assets, "sprite.svg"))
	if err != nil {
		return err
	}

	if err := os.RemoveAll(o.Out); err != nil {
		return err
	}

	for _, p := range pages {
		dir := filepath.Join(o.Out, filepath.FromSlash(strings.Trim(p.Route, "/")))
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
		if err := write(filepath.Join(dir, "index.html"), tpl, "page.html", pageData{
			Page: p, Nav: sections,
			Sprite: template.HTML(sprite),
			Body:   template.HTML(p.HTML),
		}); err != nil {
			return err
		}

		// Каждый пример — самостоятельный документ. Он подключает кит той
		// же строкой, что и любое приложение, и ничего не знает ни о
		// справочнике, ни о его теме.
		for _, d := range p.Demos {
			f := filepath.Join(o.Out, "demo", d.ID+".html")
			if err := os.MkdirAll(filepath.Dir(f), 0o755); err != nil {
				return err
			}
			if err := write(f, tpl, "demo.html", demoData{
				Sprite: template.HTML(sprite),
				Markup: template.HTML(d.Markup),
				Title:  p.Title,
			}); err != nil {
				return err
			}
		}
	}

	for _, name := range []string{"docs.css", "docs.js", "demo.css", "demo.js"} {
		b, err := files.ReadFile("assets/" + name)
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(o.Out, name), b, 0o644); err != nil {
			return err
		}
	}

	if err := copyDir(o.Kit, filepath.Join(o.Out, "kit")); err != nil {
		return err
	}
	if err := copyDir(o.Assets, filepath.Join(o.Out, "assets")); err != nil {
		return err
	}
	return searchIndex(o.Out, pages)
}

func write(path string, tpl *template.Template, name string, data any) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return tpl.ExecuteTemplate(f, name, data)
}

// copyDir кладёт кит рядом с сайтом КАК ЕСТЬ, без обработки: сайт обязан
// показывать ровно тот CSS, который получит потребитель, — пропущенный
// через сборщик это уже другой файл.
func copyDir(src, dst string) error {
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		b, err := os.ReadFile(filepath.Join(src, e.Name()))
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(dst, e.Name()), b, 0o644); err != nil {
			return err
		}
	}
	return nil
}

type doc struct {
	R string `json:"r"`
	T string `json:"t"`
	G string `json:"g"`
	B string `json:"b"`
}

// searchIndex пишет один JSON на весь сайт.
//
// Готовый поисковик тянет бинарник и WASM на каждую страницу. На семидесяти
// страницах это дороже самого индекса: весь текст помещается в полсотни
// килобайт, а поиск по нему — в полсотни строк.
func searchIndex(out string, pages []*content.Page) error {
	docs := make([]doc, 0, len(pages))
	for _, p := range pages {
		b := p.Text
		if len(b) > 1200 {
			b = b[:1200]
		}
		docs = append(docs, doc{R: p.Route, T: p.Title, G: p.Group, B: b})
	}
	f, err := os.Create(filepath.Join(out, "search.json"))
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(docs)
}

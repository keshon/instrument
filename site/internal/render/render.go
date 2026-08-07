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
	"instrument/site/internal/i18n"
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

	// Lang и Langs — для шапки: атрибут документа и переключатель.
	Lang  i18n.Lang
	Langs []langLink
}

// langLink — та же страница на другом языке. Адрес считается здесь, потому
// что переключатель обязан вести НА ЭТУ ЖЕ страницу, а не на главную: увести
// читателя на главную при смене языка — потерять место, где он был.
type langLink struct {
	Lang    i18n.Lang
	Label   string
	Route   string
	Current bool
}

type demoData struct {
	Sprite  template.HTML
	Markup  template.HTML
	Title   string
	Context bool
}

func Site(byLang map[i18n.Lang][]*content.Page, sections map[i18n.Lang][]nav.Section, o Options) error {
	tpl, err := template.New("").Funcs(template.FuncMap{
		"same": func(a, b string) bool { return a == b },
		"t":    i18n.T,
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

	// Столы примеров пишутся ОДИН раз на все языки: разметка примера языка не
	// имеет, а второй комплект из 144 документов пришлось бы держать в
	// синхроне руками.
	seenDemo := map[string]bool{}

	for _, lang := range i18n.All {
		for _, p := range byLang[lang] {
			dir := filepath.Join(o.Out, filepath.FromSlash(strings.Trim(p.Route, "/")))
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return err
			}
			if err := write(filepath.Join(dir, "index.html"), tpl, "page.html", pageData{
				Page: p, Nav: sections[lang],
				Sprite: template.HTML(sprite),
				Body:   template.HTML(p.HTML),
				Lang:   lang,
				Langs:  langLinks(p),
			}); err != nil {
				return err
			}

			// Каждый пример — самостоятельный документ. Он подключает кит той
			// же строкой, что и любое приложение, и ничего не знает ни о
			// справочнике, ни о его теме.
			for _, d := range p.Demos {
				if seenDemo[d.ID] {
					continue
				}
				seenDemo[d.ID] = true
				f := filepath.Join(o.Out, "demo", d.ID+".html")
				if err := os.MkdirAll(filepath.Dir(f), 0o755); err != nil {
					return err
				}
				if err := write(f, tpl, "demo.html", demoData{
					Sprite:  template.HTML(sprite),
					Markup:  template.HTML(d.Markup),
					Title:   p.Title,
					Context: d.Context,
				}); err != nil {
					return err
				}
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
	for _, lang := range i18n.All {
		if err := searchIndex(o.Out, lang, byLang[lang]); err != nil {
			return err
		}
	}
	return nil
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
	S string `json:"s"` // слаг: латинское имя страницы
	G string `json:"g"`
	O string `json:"o"` // имена, которые страница ОПИСЫВАЕТ
	N string `json:"n"` // имена, которые на ней просто встречаются
	B string `json:"b"`
}

// searchIndex пишет один JSON на весь сайт.
//
// Готовый поисковик тянет бинарник и WASM на каждую страницу. На семидесяти
// страницах это дороже самого индекса: весь текст помещается в несколько
// сотен килобайт, а поиск по нему — в полсотни строк.
//
// Три вещи, которых здесь раньше не было, и без каждой поиск не находил
// компоненты:
//
//	слаг    латинское имя страницы уже существует и никем не использовалось.
//	        Документация по-русски, API по-английски, и слаг — единственный
//	        мост между ними: «dialog» находит «Модалку» без словаря синонимов;
//	имена   классы и токены живут в блоках кода, а прозаический индекс их
//	        вырезает — запрос «inst-btn» давал ноль результатов;
//	тело    целиком. Обрезка в 1200 БАЙТ оставляла от страницы два абзаца и
//	        вдобавок рвала кириллицу посередине руны.
//
// langLinks — адреса этой же страницы на всех языках.
func langLinks(p *content.Page) []langLink {
	var out []langLink
	for _, l := range i18n.All {
		route := l.Prefix() + "/"
		if p.Slug != "index" {
			route = l.Prefix() + "/" + strings.Trim(p.Dir+"/"+p.Slug, "/") + "/"
		} else if p.Dir != "" {
			route = l.Prefix() + "/" + p.Dir + "/"
		}
		out = append(out, langLink{Lang: l, Label: l.Label(), Route: route, Current: l == p.Lang})
	}
	return out
}

// Индекс СВОЙ у каждого языка: искать по русским телам из английской версии
// значило бы отдавать читателю страницы, которых он не прочтёт.
func searchIndex(out string, lang i18n.Lang, pages []*content.Page) error {
	docs := make([]doc, 0, len(pages))
	for _, p := range pages {
		docs = append(docs, doc{
			R: p.Route, T: p.Title, S: p.Slug, G: p.Group,
			O: p.Own, N: p.Names, B: p.Text,
		})
	}
	name := "search.json"
	if lang.Prefix() != "" {
		name = strings.TrimPrefix(lang.Prefix(), "/") + "-search.json"
	}
	f, err := os.Create(filepath.Join(out, name))
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(docs)
}

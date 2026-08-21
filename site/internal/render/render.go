// Package render writes the built site to disk.
package render

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"instrument/site/internal/content"
	"instrument/site/internal/i18n"
	"instrument/site/internal/nav"
)

//go:embed templates/*.html assets/*
var files embed.FS

func Stylesheets() (map[string]string, error) {
	out := map[string]string{}
	for _, name := range []string{"docs.css"} {
		b, err := files.ReadFile("assets/" + name)
		if err != nil {
			return nil, err
		}
		out[name] = string(b)
	}
	return out, nil
}

type Options struct {
	Out    string
	Kit    string
	Assets string
	// Public is a directory whose contents are laid into the ROOT of the
	// output as they are.
	//
	// What lives here is neither a page nor a resource of the kit, yet has to
	// end up beside them: the CNAME for the domain, robots.txt. It could be
	// appended in CI — and then a locally built site would differ from the
	// published one by precisely the files that decide where the site opens
	// at all.
	Public string
}

type pageData struct {
	Page   *content.Page
	Nav    []nav.Section
	Sprite template.HTML
	Body   template.HTML

	Lang  i18n.Lang
	Langs []langLink
}

type langLink struct {
	Lang    i18n.Lang
	Label   string
	Route   string
	Current bool
}

func Site(byLang map[i18n.Lang][]*content.Page, sections map[i18n.Lang][]nav.Section, o Options) error {
	tpl, err := template.New("").Funcs(template.FuncMap{
		"same":  func(a, b string) bool { return a == b },
		"t":     i18n.T,
		"index": searchIndexPath,
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

		}
	}

	for _, name := range []string{"docs.css", "docs.js"} {
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
	// The directory may be absent: the site builds without it.
	if o.Public != "" {
		if _, err := os.Stat(o.Public); err == nil {
			if err := copyDir(o.Public, o.Out); err != nil {
				return err
			}
		} else if !os.IsNotExist(err) {
			return err
		}
	}
	for _, lang := range i18n.All {
		if err := searchIndex(o.Out, lang, byLang[lang]); err != nil {
			return err
		}
	}
	return verifyRefs(o.Out, string(sprite))
}

var indexRefRe = regexp.MustCompile(`data-index="([^"]*)"`)
var useRefRe = regexp.MustCompile(`<use href="#([^"]*)"`)
var symbolIDRe = regexp.MustCompile(`<symbol id="([^"]+)"`)

// verifyIndexRefs holds every page to the file it asks the browser for.
//
// The search index is the one asset a page names by an address computed at
// build time, and a wrong address fails NOWHERE the build can see it: the page
// renders, the gates pass, and the box simply returns nothing when a reader
// types into it. That is what the base-language flip did — the English page
// asked for a file nobody writes, and the Russian one searched English bodies.
//
// A GLYPH IS CHECKED THE SAME WAY, and it was the second half of one hole. A
// `<use>` whose href names no symbol draws NOTHING — no error, no box, no
// broken-image mark — and an empty `href="#"` draws nothing just as quietly.
// Nine items of the side column stood that way, their labels shifted against
// the rest of the list, and the only thing that ever reported it was somebody
// looking at the page.
func verifyRefs(out, sprite string) error {
	symbols := map[string]bool{}
	for _, m := range symbolIDRe.FindAllStringSubmatch(sprite, -1) {
		symbols[m[1]] = true
	}
	return filepath.WalkDir(out, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || d.Name() != "index.html" {
			return err
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(out, p)
		rel = filepath.ToSlash(rel)
		text := string(b)
		for _, m := range indexRefRe.FindAllStringSubmatch(text, -1) {
			ref := filepath.Join(out, filepath.FromSlash(strings.TrimPrefix(m[1], "/")))
			if _, err := os.Stat(ref); err != nil {
				return fmt.Errorf("%s asks for %s, which the build never wrote", rel, m[1])
			}
		}
		for _, m := range useRefRe.FindAllStringSubmatch(text, -1) {
			if m[1] == "" {
				return fmt.Errorf("%s draws a glyph with an empty href", rel)
			}
			if !symbols[m[1]] {
				return fmt.Errorf("%s draws #%s, which the sprite has no symbol for", rel, m[1])
			}
		}
		return nil
	})
}

func write(path string, tpl *template.Template, name string, data any) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return tpl.ExecuteTemplate(f, name, data)
}

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
	S string `json:"s"`
	G string `json:"g"`
	O string `json:"o"`
	N string `json:"n"`
	B string `json:"b"`
}

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

// searchIndexPath names the index of a language, and it is ONE function
// because the page has to fetch exactly the file the build wrote. The rule
// used to be spelled twice — here in Go and again in docs.js, where the base
// language was the literal "ru". The day English took the bare name, the
// English page asked for a file nobody writes and the Russian page searched
// English bodies; both failed in the browser alone, where no gate was looking.
// Now the address is printed into the markup and the script reads it.
func searchIndexPath(lang i18n.Lang) string {
	if lang.Prefix() == "" {
		return "/search.json"
	}
	return lang.Prefix() + "-search.json"
}

func searchIndex(out string, lang i18n.Lang, pages []*content.Page) error {
	docs := make([]doc, 0, len(pages))
	for _, p := range pages {
		docs = append(docs, doc{
			R: p.Route, T: p.Title, S: p.Slug, G: p.Group,
			O: p.Own, N: p.Names, B: p.Text,
		})
	}
	f, err := os.Create(filepath.Join(out, strings.TrimPrefix(searchIndexPath(lang), "/")))
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(docs)
}

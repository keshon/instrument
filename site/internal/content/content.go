// Package content читает страницы документации: frontmatter, markdown,
// оглавление и живые примеры.
package content

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"

	"instrument/site/internal/i18n"
)

type Page struct {
	Route      string
	Rel        string
	Dir        string
	Slug       string
	Title      string
	Group      string
	Source     string
	JS         string
	JSOpt      string
	Splash     bool
	Lang       i18n.Lang
	Translated bool
	APIFrom    string
	Icon       string
	Layout     string
	Shape      int
	HTML       string
	TOC        []Heading
	Sections   []Section
	Hero       bool
	HasJS      bool
	Demos      []Demo
	Text       string
	Names      string
	Own        string
	API        []APIRow
	body       []byte
}

type APIRow struct {
	Name  string `yaml:"name"`
	Kind  string `yaml:"kind"`
	Value string `yaml:"value"`
	Doc   string `yaml:"doc"`
}

var apiKinds = []string{"класс", "модификатор", "атрибут", "событие", "переменная", "токен"}

var declRe = regexp.MustCompile(`(?m)(?:^|[\s{;])(--[a-z][\w-]*)\s*:\s*([^;{}]+)[;}]`)

type Token struct{ Value, File string }

func TokenValues(kitDir string) (map[string]Token, error) {
	out := map[string]Token{}
	entries, err := os.ReadDir(kitDir)
	if err != nil {
		return nil, err
	}

	names := []string{"tokens.css"}
	for _, e := range entries {
		n := e.Name()
		if e.IsDir() || !strings.HasSuffix(n, ".css") {
			continue
		}
		if n == "tokens.css" || n == "print.css" || n == "motion.css" || n == "kit.css" {
			continue
		}
		names = append(names, n)
	}
	for _, name := range names {
		b, err := os.ReadFile(filepath.Join(kitDir, name))
		if err != nil {
			return nil, err
		}
		for _, m := range declRe.FindAllStringSubmatch(string(b), -1) {
			if _, ok := out[m[1]]; !ok {
				out[m[1]] = Token{
					Value: strings.Join(strings.Fields(m[2]), " "),
					File:  "src/" + name,
				}
			}
		}
	}
	return out, nil
}

func ResolveTokens(pages []*Page, tokens map[string]Token) {
	for _, p := range pages {
		// Страница-справочник перечисляет всё, что объявлено в её файле.
		if p.APIFrom == "kit" {
			var names []string
			for name, t := range tokens {
				if t.File == p.Source {
					names = append(names, name)
				}
			}
			sort.Strings(names)
			for _, name := range names {
				p.API = append(p.API, APIRow{
					Name: name, Kind: "токен", Value: tokens[name].Value,
				})
			}
		}

		var own []string
		for i, r := range p.API {
			isToken := strings.HasPrefix(r.Name, "--")
			t, known := tokens[r.Name]
			if isToken && known && r.Value == "" {
				p.API[i].Value = t.Value
			}
			if isToken && known && t.File != p.Source {
				continue // используется, но объявлено не здесь
			}
			own = append(own, r.Name)
		}
		p.Own = strings.Join(own, " ")
	}
}

// sprite — содержимое спрайта кита. Нужен галерее иконок: она печатает
// список символов из него, а не из таблицы, набранной руками.
var sprite string

func SetSprite(s string) { sprite = s }

// Связи между компонентами: имя → имена соседей. Приходят из components.json,
// того же файла, по которому cmd/registry сверяет сквозные оси.
//
// Реестр, а не фронтматтер страницы, по одной причине: связь — утверждение о
// ДВОИХ, и лежать ей полагается там, где обе стороны видны сразу. Списки,
// набранные руками на каждой странице, расходятся по одной строке за правку, и
// заметить это можно только прочитав обе.
var relations map[string][]string

func SetRelations(m map[string][]string) { relations = m }

// pageBySlug — страницы по слагу, для разрешения имён из реестра в адреса.
// Заполняется один раз перед разметкой: соседи бывают из другого каталога, и
// искать их обходом на каждую ссылку значило бы обходить дерево двести раз.
var pageBySlug = map[string]*Page{}

func indexPages(pages []*Page) {
	for _, p := range pages {
		if p.Lang == i18n.RU {
			pageBySlug[p.Slug] = p
		}
	}
}

// Related возвращает соседей страницы: имя компонента совпадает со слагом
// страницы у всех шестидесяти восьми записей реестра, поэтому отображение
// хранить негде и незачем.
func (p *Page) Related() []string { return relations[p.Slug] }

// LoadRelations читает связи из реестра компонентов.
//
// Отсутствие файла — ошибка, а не пустой список: раздел «Связанное» на странице
// формы 2 берётся только отсюда, и молча отдать ноль связей значит собрать
// страницу без соседей и не сказать об этом.
func LoadRelations(path string) (map[string][]string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var reg map[string]struct {
		Related []string `json:"related"`
	}
	if err := json.Unmarshal(b, &reg); err != nil {
		return nil, err
	}
	out := map[string][]string{}
	for name, e := range reg {
		if len(e.Related) > 0 {
			out[name] = e.Related
		}
	}
	return out, nil
}

func Render(pages []*Page) error {
	indexPages(pages)
	for _, p := range pages {
		if err := renderMarkdown(p, p.body); err != nil {
			return fmt.Errorf("%s: %w", p.Rel, err)
		}
	}
	return nil
}

func apiRank(kind string) (int, bool) {
	for i, k := range apiKinds {
		if k == kind {
			return i, true
		}
	}
	return 0, false
}

type Heading struct {
	Level int
	ID    string
	Text  string
}

type Demo struct {
	ID      string
	Markup  string
	Context bool
}

type frontmatter struct {
	Title    string   `yaml:"title"`
	Group    string   `yaml:"group"`
	Source   string   `yaml:"source"`
	JS       string   `yaml:"js"`
	JSOpt    string   `yaml:"js-optional"`
	Template string   `yaml:"template"`
	Layout   string   `yaml:"layout"`
	API      []APIRow `yaml:"api"`
	APIFrom  string   `yaml:"api-from"`

	// Shape is the version of the page contract this page follows. Empty means
	// 1: the dictionary the reference was written on. 2 is the shape from
	// docs/internal/DOCS-SHAPE.md.
	//
	// It is a field of its own rather than a second value of `layout` because
	// the two answer different questions: layout says how the page is RENDERED,
	// shape says which sections it OWES. During the move one shape is rendered
	// by both layouts, and one layout serves both shapes.
	Shape int `yaml:"shape"`

	TitleEN string `yaml:"title-en"`
	GroupEN string `yaml:"group-en"`
}

// shapeOf normalises the declared page shape.
//
// The default is 1 rather than "the newest": a page that says nothing keeps
// the contract it was written under. Otherwise switching the dictionary would
// turn all fifty pages red in one commit, and a gate raised red lives with its
// check disabled — the same rule the migration zones in cmd/lang are built on.
func shapeOf(v int) int {
	if v == 0 {
		return 1
	}
	return v
}

var fmRe = regexp.MustCompile(`(?s)\A---\r?\n(.*?)\r?\n---\r?\n`)

func Collect(root string) (map[i18n.Lang][]*Page, error) {
	var sources []string
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == "internal" {
				return filepath.SkipDir
			}
			return nil
		}
		rel, _ := filepath.Rel(root, p)
		rel = filepath.ToSlash(rel)
		if !strings.HasSuffix(rel, ".md") || rel == "README.md" {
			return nil
		}
		base := strings.TrimSuffix(rel, ".md")
		for _, l := range i18n.Known {
			if l != i18n.RU && strings.HasSuffix(base, l.Suffix()) {
				return nil
			}
		}
		sources = append(sources, rel)
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(sources)

	out := map[i18n.Lang][]*Page{}
	for _, lang := range i18n.All {
		for _, rel := range sources {
			file := filepath.Join(root, filepath.FromSlash(rel))
			translated := false
			if lang != i18n.RU {
				alt := filepath.Join(root, filepath.FromSlash(
					strings.TrimSuffix(rel, ".md")+lang.Suffix()+".md"))
				if _, err := os.Stat(alt); err == nil {
					file, translated = alt, true
				}
			} else {
				translated = true
			}
			page, err := parse(file, rel, lang, translated)
			if err != nil {
				return nil, fmt.Errorf("%s [%s]: %w", rel, lang, err)
			}
			out[lang] = append(out[lang], page)
		}
	}
	return out, nil
}

func Flat(byLang map[i18n.Lang][]*Page) []*Page {
	var all []*Page
	for _, l := range i18n.All {
		all = append(all, byLang[l]...)
	}
	return all
}

func parse(fsPath, rel string, lang i18n.Lang, translated bool) (*Page, error) {
	raw, err := os.ReadFile(fsPath)
	if err != nil {
		return nil, err
	}

	var meta frontmatter
	body := raw
	if m := fmRe.FindSubmatch(raw); m != nil {
		if err := yaml.Unmarshal(m[1], &meta); err != nil {
			return nil, fmt.Errorf("frontmatter: %w", err)
		}
		body = raw[len(m[0]):]
	}

	slug := strings.TrimSuffix(path.Base(rel), ".md")
	dir := path.Dir(rel)
	if dir == "." {
		dir = ""
	}
	route := lang.Prefix() + "/"
	if slug != "index" {
		route = lang.Prefix() + "/" + path.Join(dir, slug) + "/"
	} else if dir != "" {
		route = lang.Prefix() + "/" + dir + "/"
	}

	p := &Page{
		Route: route, Rel: rel, Dir: dir, Slug: slug,
		Title: meta.Title, Group: meta.Group, Source: meta.Source,
		JS:         meta.JS,
		JSOpt:      meta.JSOpt,
		Splash:     meta.Template == "splash",
		Layout:     meta.Layout,
		Shape:      shapeOf(meta.Shape),
		API:        meta.API,
		APIFrom:    meta.APIFrom,
		Lang:       lang,
		Translated: translated,
	}

	if !translated && lang == i18n.EN {
		if meta.TitleEN != "" {
			p.Title = meta.TitleEN
		}
		if meta.GroupEN != "" {
			p.Group = meta.GroupEN
		}
	}

	for i, row := range p.API {
		if row.Name == "" {
			return nil, fmt.Errorf("api[%d]: пустое имя", i)
		}
		if _, ok := apiRank(row.Kind); !ok {
			return nil, fmt.Errorf("api[%d] (%s): вид %q не из словаря (%s)",
				i, row.Name, row.Kind, strings.Join(apiKinds, " · "))
		}
	}

	sort.SliceStable(p.API, func(i, j int) bool {
		ri, _ := apiRank(p.API[i].Kind)
		rj, _ := apiRank(p.API[j].Kind)
		return ri < rj
	})
	if p.Title == "" {
		p.Title = slug
	}
	if slug != "index" {
		p.Icon = "i-p-" + slug
	}

	p.body = body
	p.Text = plainText(body)
	p.Names = apiNames(raw)

	own := make([]string, 0, len(p.API))
	for _, r := range p.API {
		own = append(own, r.Name)
	}
	p.Own = strings.Join(own, " ")
	return p, nil
}

var stripRe = regexp.MustCompile("(?s)```.*?```|`[^`]*`|[|*_#>\\[\\]()]")

func plainText(body []byte) string {
	s := stripRe.ReplaceAllString(string(body), " ")
	return strings.Join(strings.Fields(s), " ")
}

var namesRe = regexp.MustCompile(`\binst-[a-z0-9-]+|--[a-z][a-z0-9-]*|\baria-[a-z-]+|\bdata-[a-z-]+`)

func apiNames(body []byte) string {
	seen := map[string]bool{}
	var out []string
	for _, m := range namesRe.FindAllString(string(body), -1) {
		if seen[m] {
			continue
		}
		seen[m] = true
		out = append(out, m)
	}
	sort.Strings(out)
	return strings.Join(out, " ")
}

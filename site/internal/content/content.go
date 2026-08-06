// Package content читает страницы документации: frontmatter, markdown,
// оглавление и живые примеры.
package content

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Page — одна страница документации.
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

	HTML  string // готовое тело страницы
	TOC   []Heading
	Demos []Demo // живые примеры, каждый уедет в свой документ
	Text  string // очищенный текст для поиска
}

type Heading struct {
	Level int
	ID    string
	Text  string
}

// Demo — живой пример. Он рисуется в ОТДЕЛЬНОМ документе внутри iframe,
// а не в потоке страницы, и это не техническая деталь:
//
//   · пример должен читаться как чужой экран, а не как часть справочника.
//     Пока он был бесшовен со страницей, компонент выглядел «сиротой»:
//     непонятно, где кончается документация и начинается демонстрация;
//   · у примера своя тема и своя плотность. На общей странице переключение
//     темы вики перекрашивало и примеры — то есть демонстрировало не кит,
//     а страницу про кит;
//   · изоляция стилей достаётся бесплатно и навсегда.
type Demo struct {
	ID     string // demo/components-actions-button-1
	Markup string
}

type frontmatter struct {
	Title    string `yaml:"title"`
	Group    string `yaml:"group"`
	Source   string `yaml:"source"`
	Status   string `yaml:"status"`
	NeedsJS  string `yaml:"needs-js"`
	Template string `yaml:"template"`
}

var fmRe = regexp.MustCompile(`(?s)\A---\r?\n(.*?)\r?\n---\r?\n`)

// Collect обходит каталог документации.
//
// internal/ пропускается намеренно: это процессные документы, а не
// документация, и они прямо говорят об этом у себя в шапке.
func Collect(root string) ([]*Page, error) {
	var pages []*Page
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
		if !strings.HasSuffix(p, ".md") {
			return nil
		}
		rel, _ := filepath.Rel(root, p)
		rel = filepath.ToSlash(rel)
		if rel == "README.md" {
			return nil
		}
		page, err := parse(p, rel)
		if err != nil {
			return fmt.Errorf("%s: %w", rel, err)
		}
		pages = append(pages, page)
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(pages, func(i, j int) bool { return pages[i].Rel < pages[j].Rel })
	return pages, nil
}

func parse(fsPath, rel string) (*Page, error) {
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
	route := "/"
	if slug != "index" {
		route = "/" + path.Join(dir, slug) + "/"
	} else if dir != "" {
		route = "/" + dir + "/"
	}

	p := &Page{
		Route: route, Rel: rel, Dir: dir, Slug: slug,
		Title: meta.Title, Group: meta.Group, Source: meta.Source,
		NeedsJS: meta.NeedsJS, Status: meta.Status,
		Splash: meta.Template == "splash",
	}
	if p.Title == "" {
		p.Title = slug
	}

	if err := renderMarkdown(p, body); err != nil {
		return nil, err
	}
	p.Text = plainText(body)
	return p, nil
}

var stripRe = regexp.MustCompile("(?s)```.*?```|`[^`]*`|[|*_#>\\[\\]()]")

func plainText(body []byte) string {
	s := stripRe.ReplaceAllString(string(body), " ")
	return strings.Join(strings.Fields(s), " ")
}

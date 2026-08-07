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

	// APIFrom == "kit" — справочник перечисляется генератором из файла,
	// указанного в source. Заводится для страницы «Токены»: рукописный
	// список из двухсот имён разойдётся с китом на первой же неделе, а
	// перечислить их машина умеет и так.
	APIFrom string

	// Icon — идентификатор символа в спрайте, выведенный из слага.
	//
	// Не поле frontmatter намеренно: связь «страница ↔ иконка» взаимно
	// однозначна, и выводить её из имени дешевле, чем поддерживать в
	// семидесяти одной шапке. Отсутствие символа ловит сборка — страница
	// без иконки иначе потерялась бы молча.
	Icon string

	HTML  string // готовое тело страницы
	TOC   []Heading
	Demos []Demo // живые примеры, каждый уедет в свой документ
	Text  string // очищенная проза для поиска
	Names string // все идентификаторы, встреченные на странице
	Own   string // идентификаторы, которые страница ОПИСЫВАЕТ (из api:)

	// API — справочник страницы ДАННЫМИ, а не прозой.
	//
	// Раньше он был четырьмя-шестью markdown-таблицами с разными колонками:
	// «Класс | Работа», «Атрибут | Значения | Где», «Переменная | По
	// умолчанию», а токены — вообще списком через точку, без значений. Глаз
	// учился читать каждую страницу заново, и это ровно то, на что жалуются
	// в первую очередь.
	//
	// Данными это становится единообразным ПО ПОСТРОЕНИЮ, а не по
	// дисциплине автора: колонки задаёт шаблон, порядок задаёт вид, а
	// выдуманный вид роняет сборку.
	API []APIRow

	// body хранится до второго прохода: разметка рендерится ПОСЛЕ того, как
	// справочник достроен. Пока рендер шёл прямо в parse, таблица собиралась
	// из ещё пустых данных — значения токенов и перечисление по api-from
	// приходят позже и в HTML не попадали.
	body []byte
}

// APIRow — одна строка справочника.
type APIRow struct {
	Name  string `yaml:"name"`
	Kind  string `yaml:"kind"`
	Value string `yaml:"value"`
	Doc   string `yaml:"doc"`
}

// Виды закрыты и упорядочены. Порядок — это порядок чтения: сначала чем
// пользуются, потом чем настраивают. Сортировка по нему делает все страницы
// одинаковыми на вид без единой договорённости между авторами.
var apiKinds = []string{"класс", "модификатор", "атрибут", "переменная", "токен"}

// tokenValues читает объявления токенов из самого кита.
//
// Значение токена НЕ пишется в странице руками. Написанное руками — это
// второй источник истины, и он разойдётся: `--control-h-md: 32px` в
// документации переживёт правку в tokens.css и будет врать беззвучно. Тот же
// довод, по которому contrast читает настоящий tokens.css, а не свою копию.
//
// Значение берётся СЫРЫМ, без разрешения var(): `--pad-panel: var(--space-5)`
// полезнее, чем `12px`, — оно показывает, откуда отступ берётся, а это и есть
// то, что нужно знать при настройке.
var declRe = regexp.MustCompile(`(?m)^\s*(--[a-z][\w-]*)\s*:\s*([^;]+);`)

// Token — объявление токена в ките: значение и файл, где оно стоит.
type Token struct{ Value, File string }

func TokenValues(kitDir string) (map[string]Token, error) {
	out := map[string]Token{}
	entries, err := os.ReadDir(kitDir)
	if err != nil {
		return nil, err
	}
	// tokens.css читается ПЕРВЫМ, а слои переопределений не читаются вовсе.
	//
	// Порядок здесь не косметика. os.ReadDir отдаёт файлы по алфавиту, и
	// print.css оказывался раньше tokens.css: печатное переопределение
	// `--surface-sunken: var(--n-2)` становилось «значением по умолчанию»
	// и попадало в справочник вместо настоящего light-dark(). Первое
	// объявление считается базовым — значит, первым обязан идти базовый файл.
	names := []string{"tokens.css"}
	for _, e := range entries {
		n := e.Name()
		if e.IsDir() || !strings.HasSuffix(n, ".css") {
			continue
		}
		// print и motion — слои переопределений, значений по умолчанию в них нет.
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
			// Первое объявление — базовое. Переопределения тем и плотностей
			// идут ниже и значением по умолчанию не являются.
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

// ResolveTokens проставляет значения токенов и уточняет владение.
//
// Значение — из кита, не из страницы: написанное руками разошлось бы с
// tokens.css при первой правке.
//
// Владение — тонкий момент, и без него поиск по имени токена не работает.
// Класс `inst-btn` описывает одна страница, а токен `--space-4` перечислен в
// справочнике у полутора десятков: он там ИСПОЛЬЗУЕТСЯ. Если считать
// владельцем каждого, то на запрос «--space-4» первым выходит лог или код —
// то есть кто угодно, кроме страницы про пространство.
//
// Правило: страница владеет токеном, только если документирует ТОТ ФАЙЛ, где
// он объявлен. Для остальных это упоминание, и вес соответствующий.
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

// Render — второй проход: разметка страниц.
//
// Отдельно от чтения, потому что справочник достраивается между ними:
// значения токенов приходят из кита, а страница «Токены» перечисляет их
// целиком. Рендер обязан видеть готовые данные, иначе печатает пустоту.
func Render(pages []*Page) error {
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

// Demo — живой пример. Он рисуется в ОТДЕЛЬНОМ документе внутри iframe,
// а не в потоке страницы, и это не техническая деталь:
//
//	· пример должен читаться как чужой экран, а не как часть справочника.
//	  Пока он был бесшовен со страницей, компонент выглядел «сиротой»:
//	  непонятно, где кончается документация и начинается демонстрация;
//	· у примера своя тема и своя плотность. На общей странице переключение
//	  темы вики перекрашивало и примеры — то есть демонстрировало не кит,
//	  а страницу про кит;
//	· изоляция стилей достаётся бесплатно и навсегда.
type Demo struct {
	ID     string // demo/components-actions-button-1
	Markup string

	// Context отмечает пример второго рода — собранный экран, а не
	// вычлёненный компонент.
	//
	// Оба нужны, и по разным причинам. Минимальный отвечает «что
	// скопировать»: в нём нет ничего, кроме показываемого компонента.
	// Контекстный отвечает «как это живёт»: полоса действий внутри панели,
	// поле внутри формы, панель внутри оболочки. Без второго компонент
	// читается сиротой — видно класс, но не видно работы.
	Context bool
}

type frontmatter struct {
	Title    string   `yaml:"title"`
	Group    string   `yaml:"group"`
	Source   string   `yaml:"source"`
	Status   string   `yaml:"status"`
	NeedsJS  string   `yaml:"needs-js"`
	Template string   `yaml:"template"`
	API      []APIRow `yaml:"api"`
	APIFrom  string   `yaml:"api-from"`
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
		Splash:  meta.Template == "splash",
		API:     meta.API,
		APIFrom: meta.APIFrom,
	}

	// Вид проверяется здесь, а не глазами при вычитке. Опечатка «класс »
	// или выдуманный «свойство» иначе просто встали бы в колонку и выглядели
	// как решение автора — тот же случай, что с data-state="quued".
	for i, row := range p.API {
		if row.Name == "" {
			return nil, fmt.Errorf("api[%d]: пустое имя", i)
		}
		if _, ok := apiRank(row.Kind); !ok {
			return nil, fmt.Errorf("api[%d] (%s): вид %q не из словаря (%s)",
				i, row.Name, row.Kind, strings.Join(apiKinds, " · "))
		}
	}
	// Порядок задаёт вид, а не автор: справочник обязан читаться одинаково
	// на всех страницах. Внутри вида сохраняется порядок написания — там он
	// осмысленный (базовый класс раньше модификаторов размера).
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
	// Имена берутся из ВСЕГО файла, включая frontmatter: справочник теперь
	// живёт там, и его строки — самые точные имена, какие у страницы есть.
	p.Names = apiNames(raw)

	// Упоминание и описание — разные вещи, и без этого различия поиск по
	// имени класса не работает. `inst-btn` встречается на двадцати шести
	// страницах, `inst-badge` на семнадцати; упоминание получает такой же
	// вес, как описание, и первой выходит конституция, где перечислено всё.
	//
	// Владение объявляется, а не угадывается: слаг здесь не помогает —
	// страница `button` описывает `inst-btn`, а `toggles` сразу чекбокс,
	// радио и свитч. Единственный честный источник — справочник страницы.
	// Владение уточняется позже, в ResolveTokens: чтобы отличить «описывает
	// токен» от «использует его», нужен весь кит, а здесь есть одна страница.
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

// namesRe — идентификаторы API страницы: классы, токены, ARIA- и
// data-атрибуты.
//
// Они живут в блоках кода и в инлайновом коде, а plainText оба ВЫРЕЗАЕТ —
// иначе прозаический индекс забился бы разметкой примеров. Из-за этого
// поиск не находил ни одного компонента по имени класса: запрос «inst-btn»
// давал ноль, «table» не возвращал страницу таблицы. Имена собираются
// отдельным полем, и вес у них выше прозы: тот, кто ищет inst-badge, ищет
// именно его.
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

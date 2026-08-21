package i18n

import gostrings "strings"

// Lang is the code of a language. An empty string is impossible: the language
// is always explicit.
type Lang string

const (
	RU Lang = "ru"
	EN Lang = "en"
)

// Base is the language that owns the bare name: `page.md` rather than
// `page.ru.md`, and `/layout/split/` rather than `/ru/layout/split/`.
//
// IT IS ONE VARIABLE AND NOT A CONSTANT ON PURPOSE. The whole of the flip that
// made English the base is this line plus the renaming of the files: every
// other place in the site asks Base rather than naming a language, so nothing
// else had to know that the answer had changed. What still names RU outright
// names it because the RUSSIAN language is what it means — the label of the
// switch, the transliteration of a heading, the Cyrillic check — and those
// would be wrong to reroute through Base.
var Base = EN

// All is what gets built, and the base goes FIRST: reports count the base
// language, and a reader who lands on the site with no prefix lands here.
var All = []Lang{EN, RU}

var Known = []Lang{RU, EN}

func (l Lang) Prefix() string {
	if l == Base {
		return ""
	}
	return "/" + string(l)
}

func (l Lang) Suffix() string {
	if l == Base {
		return ""
	}
	return "." + string(l)
}

// Short is the two letters of the switch in the header. It is the CODE
// upper-cased rather than a dictionary entry: a language names itself the same
// way in every language, and a second spelling would be a second thing to keep
// in agreement for no gain.
func (l Lang) Short() string {
	return gostrings.ToUpper(string(l))
}

func (l Lang) Label() string {
	switch l {
	case EN:
		return "English"
	default:
		return "Русский"
	}
}

// THE RUSSIAN HALF OF EVERY ENTRY IS THE TRANSLATION ITSELF, not text waiting
// for one, and that is why this file is named in the exception list of the
// language gate. The keys are Russian too where a page writes them — the
// section titles and the kinds — and both halves stand here so that a page in
// either language finds its own.
//
// The base language flips at step five: RU stops being the empty prefix, EN
// takes its place, and the map is untouched by that — only Prefix, Suffix and
// the order of All change.
var strings = map[string]map[Lang]string{
	// A format string, and the only one in the dictionary: the count is
	// computed from the sprite, and a phrase built by gluing pieces around it
	// would put the word order of one language into both.
	"icons.note": {
		RU: `<p class="icon-note">Плюс %d символов страниц справочника ` +
			`(<code>i-p-*</code>) — они принадлежат сайту, а не киту.</p>`,
		EN: `<p class="icon-note">Plus %d symbols of the reference's own pages ` +
			`(<code>i-p-*</code>) — they belong to the site, not to the kit.</p>`,
	},
	"skip":          {RU: "Перейти к содержимому", EN: "Skip to content"},
	"sections":      {RU: "Разделы", EN: "Sections"},
	"nav.aria":      {RU: "Разделы документации", EN: "Documentation sections"},
	"search":        {RU: "Поиск  /", EN: "Search  /"},
	"search.aria":   {RU: "Поиск по документации", EN: "Search the documentation"},
	"search.result": {RU: "Результаты", EN: "Results"},
	"search.empty":  {RU: "Ничего не найдено", EN: "Nothing found"},
	"view":          {RU: "Вид", EN: "View"},
	"density":       {RU: "Плотность", EN: "Density"},
	"density.sm":    {RU: "Плотно", EN: "Compact"},
	"density.md":    {RU: "Обычно", EN: "Default"},
	"density.lg":    {RU: "Свободно", EN: "Comfortable"},
	// A scale item is labelled by type size, and the height of a control moves
	// into the explanation: "large" says nothing until it is pressed, "16px"
	// says it at once.
	"scale":            {RU: "Масштаб", EN: "Scale"},
	"scale.aria":       {RU: "Базовый кегль интерфейса", EN: "Interface base size"},
	"scale.base":       {RU: "умолчание", EN: "default"},
	"scale.ctl":        {RU: "контрол", EN: "control"},
	"theme":            {RU: "Тема", EN: "Theme"},
	"theme.aria":       {RU: "Тема оформления", EN: "Colour theme"},
	"theme.system":     {RU: "По системе", EN: "System"},
	"theme.ln":         {RU: "Светлая нейтральная", EN: "Light neutral"},
	"theme.l":          {RU: "Светлая тёплая", EN: "Light warm"},
	"theme.lc":         {RU: "Светлая холодная", EN: "Light cool"},
	"theme.ds":         {RU: "Тёмная серая", EN: "Dark grey"},
	"theme.d":          {RU: "Тёмная чёрная", EN: "Dark black"},
	"accent":           {RU: "Акцент", EN: "Accent"},
	"accent.aria":      {RU: "Акцентный тон", EN: "Accent hue"},
	"accent.pe":        {RU: "Петроль", EN: "Petrol"},
	"accent.gr":        {RU: "Графит", EN: "Graphite"},
	"accent.in":        {RU: "Индиго", EN: "Indigo"},
	"accent.cl":        {RU: "Глина", EN: "Clay"},
	"lang":             {RU: "Язык", EN: "Language"},
	"toc":              {RU: "На этой странице", EN: "On this page"},
	"source":           {RU: "Исходник", EN: "Source"},
	"github":           {RU: "Исходный код на GitHub", EN: "Source code on GitHub"},
	"license":          {RU: "Открытый код под MIT", EN: "Open source under MIT"},
	"why":              {RU: "Почему кит устроен именно так", EN: "Why the kit works this way"},
	"constitution":     {RU: "конституция", EN: "the constitution"},
	"js":               {RU: "JS", EN: "JS"},
	"js.optional":      {RU: "JS по желанию", EN: "JS optional"},
	"demo":             {RU: "Пример", EN: "Example"},
	"demo.context":     {RU: "В контексте", EN: "In context"},
	"demo.markup":      {RU: "Разметка", EN: "Markup"},
	"demo.code":        {RU: "Код", EN: "Code"},
	"demo.hero":        {RU: "Компонент", EN: "The component"},
	"contents":         {RU: "Разделы страницы", EN: "Page sections"},
	"copy":             {RU: "копировать", EN: "copy"},
	"copy.done":        {RU: "скопировано", EN: "copied"},
	"copy.fail":        {RU: "не вышло", EN: "failed"},
	"api.name":         {RU: "Имя", EN: "Name"},
	"api.value":        {RU: "Значение", EN: "Value"},
	"api.doc":          {RU: "Что делает", EN: "What it does"},
	"api.classes":      {RU: "Классы", EN: "Classes"},
	"api.attrs":        {RU: "Атрибуты", EN: "Attributes"},
	"api.tokens":       {RU: "Токены", EN: "Tokens"},
	"api.where":        {RU: "Откуда", EN: "Where from"},
	"api.one":          {RU: "одно на все ячейки", EN: "one across every cell"},
	"api.set":          {RU: "задаёт приложение", EN: "the application sets it"},
	"api.base":         {RU: "база", EN: "base"},
	"api.axis.scale":   {RU: "масштаб", EN: "scale"},
	"api.axis.density": {RU: "плотность", EN: "density"},
	"api.axis.theme":   {RU: "тема", EN: "theme"},
	"api.axis.accent":  {RU: "акцент", EN: "accent"},
	"api.axis.screen":  {RU: "экран", EN: "screen"},
	// NEITHER HALF NAMES A LANGUAGE, and that is the point. The banner stands
	// on a page whose own half is missing, so what shows through is whatever
	// Base happens to be — and Base changed hands once already. Naming Russian
	// here was true while Russian was the base; the day English took the bare
	// name, the Russian half of this entry was empty and the English one said
	// the opposite of what a reader would see.
	"untranslated": {
		RU: "Эта страница ещё не переведена — ниже текст на базовом языке " +
			"справочника. Справочник API от этого не меняется: имена классов, " +
			"токенов и атрибутов одинаковы в обоих языках.",
		EN: "This page is not translated yet — the text below is in the " +
			"reference's base language. The API reference is unaffected: the " +
			"names of classes, tokens and attributes are the same in both.",
	},
}

// An EMPTY half is a missing translation that no panic reports: T hands the
// empty string back, the template prints nothing, and the page comes out with
// a hole where a sentence was. That is exactly how the "not translated yet"
// banner spent the whole of the bilingual period blank on one side — the side
// nobody could reach until the base language flipped and made it the only
// reachable one.
//
// Checked at init rather than in a gate: the dictionaries live here, the loop
// costs microseconds, and a build that would print a hole must not start.
func init() {
	for name, dict := range map[string]map[string]map[Lang]string{
		"strings": strings, "kinds": kinds,
		"sectionTitles": sectionTitles, "sections": sections,
	} {
		for key, m := range dict {
			for _, l := range Known {
				v, ok := m[l]
				if !ok {
					panic("i18n: " + name + "[" + key + "] has no " + string(l))
				}
				if v == "" {
					panic("i18n: " + name + "[" + key + "] is empty in " + string(l))
				}
			}
		}
	}
}

func T(l Lang, key string) string {
	m, ok := strings[key]
	if !ok {
		panic("i18n: no such key " + key)
	}
	s, ok := m[l]
	if !ok {
		panic("i18n: no translation of " + key + " into " + string(l))
	}
	return s
}

var kinds = map[string]map[Lang]string{
	"класс":       {RU: "класс", EN: "class"},
	"модификатор": {RU: "модификатор", EN: "modifier"},
	"атрибут":     {RU: "атрибут", EN: "attribute"},
	"событие":     {RU: "событие", EN: "event"},
	"переменная":  {RU: "переменная", EN: "variable"},
	"токен":       {RU: "токен", EN: "token"},
}

func Kind(l Lang, kind string) string {
	if m, ok := kinds[kind]; ok {
		if s, ok := m[l]; ok {
			return s
		}
	}
	return kind
}

var sectionTitles = map[string]map[Lang]string{
	"install":       {RU: "Установка", EN: "Installation"},
	"usage":         {RU: "Использование", EN: "Usage"},
	"contract":      {RU: "Контракт", EN: "Contract"},
	"when":          {RU: "Когда использовать", EN: "When to use"},
	"anatomy":       {RU: "Устройство", EN: "Anatomy"},
	"scale":         {RU: "Шкала", EN: "Scale"},
	"variants":      {RU: "Варианты", EN: "Variants"},
	"sizes":         {RU: "Размеры", EN: "Sizes"},
	"states":        {RU: "Состояния", EN: "States"},
	"icons":         {RU: "С иконкой", EN: "With icon"},
	"behavior":      {RU: "Поведение", EN: "Behavior"},
	"js":            {RU: "JS", EN: "JS"},
	"composition":   {RU: "Композиции", EN: "Composition"},
	"patterns":      {RU: "Сценарии", EN: "Patterns"},
	"rules":         {RU: "Правила", EN: "Rules"},
	"a11y":          {RU: "Доступность", EN: "Accessibility"},
	"customization": {RU: "Настройка", EN: "Customization"},
	"api":           {RU: "API", EN: "API"},
	"related":       {RU: "Связанное", EN: "Related"},
}

func SectionTitle(l Lang, id string) string {
	m, ok := sectionTitles[id]
	if !ok {
		panic("i18n: no title for section " + id)
	}
	s, ok := m[l]
	if !ok {
		panic("i18n: no translation of section " + id + " into " + string(l))
	}
	return s
}

var sections = map[string]map[Lang]string{
	"start":                 {RU: "Начало", EN: "Start"},
	"foundations":           {RU: "Основания", EN: "Foundations"},
	"layout":                {RU: "Раскладка", EN: "Layout"},
	"components/actions":    {RU: "Действия", EN: "Actions"},
	"components/inputs":     {RU: "Ввод", EN: "Inputs"},
	"components/display":    {RU: "Отображение данных", EN: "Data display"},
	"components/charts":     {RU: "Графики", EN: "Charts"},
	"components/navigation": {RU: "Навигация", EN: "Navigation"},
	"components/overlays":   {RU: "Оверлеи", EN: "Overlays"},
	"components/feedback":   {RU: "Обратная связь", EN: "Feedback"},
	"agent":                 {RU: "Агентный слой", EN: "Agent layer"},
	"blocks":                {RU: "Сборки", EN: "Blocks"},
	"about":                 {RU: "О проекте", EN: "About"},
}

func Section(l Lang, dir string) string {
	if m, ok := sections[dir]; ok {
		if s, ok := m[l]; ok {
			return s
		}
	}
	return dir
}

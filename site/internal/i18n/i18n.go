// Package i18n держит язык сайта в одном месте.
//
// Русский остаётся ИСХОДНИКОМ и живёт в корне: /components/actions/button/.
// Английский — под префиксом: /en/components/actions/button/. Обратный
// порядок сломал бы каждый существующий адрес и каждый якорь, а выигрыша не
// дал бы никакого.
//
// Перевод раскатывается ПОСТРАНИЧНО. Страница без своего button.en.md не
// исчезает и не отдаёт 404 — она показывает русский текст с пометкой. Иначе
// перевод семидесяти пяти страниц был бы «всё или ничего», то есть не
// начался бы никогда.
package i18n

// Lang — код языка. Пустая строка невозможна: язык всегда явный.
type Lang string

const (
	RU Lang = "ru"
	EN Lang = "en"
)

// All — порядок в переключателе и порядок сборки.
var All = []Lang{RU, EN}

// Prefix — префикс маршрута. У исходного языка его нет.
func (l Lang) Prefix() string {
	if l == RU {
		return ""
	}
	return "/" + string(l)
}

// Suffix — хвост имени файла: button.md для исходного, button.en.md для
// остальных. Соседний файл, а не параллельное дерево: дерево пришлось бы
// держать в синхроне руками, и оно бы разъехалось.
func (l Lang) Suffix() string {
	if l == RU {
		return ""
	}
	return "." + string(l)
}

// Label — имя языка в переключателе, на самом языке.
func (l Lang) Label() string {
	switch l {
	case EN:
		return "English"
	default:
		return "Русский"
	}
}

// strings — вся видимая пользователю обвязка сайта.
//
// Ключи английские, потому что это идентификаторы, а не текст. Отсутствие
// перевода ловится сборкой: T() падает на неизвестном ключе, а не отдаёт
// пустоту, — молчаливая пустота в интерфейсе выглядит как решение дизайнера.
var strings = map[string]map[Lang]string{
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
	"theme":         {RU: "Тема", EN: "Theme"},
	"theme.aria":    {RU: "Тема оформления", EN: "Colour theme"},
	"theme.system":  {RU: "По системе", EN: "System"},
	"theme.ln":      {RU: "Светлая нейтральная", EN: "Light neutral"},
	"theme.l":       {RU: "Светлая тёплая", EN: "Light warm"},
	"theme.lc":      {RU: "Светлая холодная", EN: "Light cool"},
	"theme.dl":      {RU: "Тёмная светло-серая", EN: "Dark light-grey"},
	"theme.ds":      {RU: "Тёмная серая", EN: "Dark grey"},
	"theme.d":       {RU: "Тёмная чёрная", EN: "Dark black"},
	"lang":          {RU: "Язык", EN: "Language"},
	"toc":           {RU: "На этой странице", EN: "On this page"},
	"source":        {RU: "Исходник", EN: "Source"},
	"why":           {RU: "Почему кит устроен именно так", EN: "Why the kit works this way"},
	"constitution":  {RU: "конституция", EN: "the constitution"},
	"needs-js":      {RU: "нужен JS", EN: "needs JS"},
	"demo":          {RU: "Пример", EN: "Example"},
	"demo.context":  {RU: "В контексте", EN: "In context"},
	"demo.markup":   {RU: "Разметка", EN: "Markup"},
	"demo.theme":    {RU: "Тема примера", EN: "Example theme"},
	"demo.same":     {RU: "как у сайта", EN: "same as site"},
	"demo.hero":     {RU: "Компонент", EN: "The component"},
	"rule.do":       {RU: "Так", EN: "Do"},
	"rule.dont":     {RU: "Не так", EN: "Don't"},
	"contents":      {RU: "Разделы страницы", EN: "Page sections"},
	"copy":          {RU: "копировать", EN: "copy"},
	"copy.done":     {RU: "скопировано", EN: "copied"},
	"copy.fail":     {RU: "не вышло", EN: "failed"},
	"api.name":      {RU: "Имя", EN: "Name"},
	"api.value":     {RU: "Значение", EN: "Value"},
	"api.doc":       {RU: "Что делает", EN: "What it does"},
	"untranslated": {
		RU: "",
		EN: "This page is not translated yet — the text below is in Russian. " +
			"The API reference is in English either way: class, token and " +
			"attribute names are the same in both languages.",
	},
}

// T возвращает строку. Неизвестный ключ — ошибка программиста, и она обязана
// быть громкой: пустая строка в интерфейсе неотличима от замысла.
func T(l Lang, key string) string {
	m, ok := strings[key]
	if !ok {
		panic("i18n: нет ключа " + key)
	}
	s, ok := m[l]
	if !ok {
		panic("i18n: нет перевода " + key + " на " + string(l))
	}
	return s
}

// Kind — имя вида в справочнике API. Словарь закрыт и совпадает с тем, что
// проверяет content.apiRank: выдуманный вид роняет сборку.
var kinds = map[string]map[Lang]string{
	"класс":       {RU: "класс", EN: "class"},
	"модификатор": {RU: "модификатор", EN: "modifier"},
	"атрибут":     {RU: "атрибут", EN: "attribute"},
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

// Заголовок канонического раздела страницы. Словарь закрыт и совпадает с тем,
// что перечисляет content.sectionDefs: неизвестный идентификатор — ошибка
// программиста, и она обязана быть громкой.
//
// Заголовок печатается ОТСЮДА, а не со страницы. Из-за этого «Справочник» на
// одной странице и «API» на другой перестают быть разными разделами на вид:
// узнан один и тот же раздел — напечатано одно и то же слово.
var sectionTitles = map[string]map[Lang]string{
	"install":       {RU: "Установка", EN: "Installation"},
	"usage":         {RU: "Использование", EN: "Usage"},
	"when":          {RU: "Когда использовать", EN: "When to use"},
	"anatomy":       {RU: "Устройство", EN: "Anatomy"},
	"scale":         {RU: "Шкала", EN: "Scale"},
	"variants":      {RU: "Варианты", EN: "Variants"},
	"sizes":         {RU: "Размеры", EN: "Sizes"},
	"states":        {RU: "Состояния", EN: "States"},
	"icons":         {RU: "С иконкой", EN: "With icon"},
	"behavior":      {RU: "Поведение", EN: "Behavior"},
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
		panic("i18n: нет заголовка раздела " + id)
	}
	s, ok := m[l]
	if !ok {
		panic("i18n: нет перевода раздела " + id + " на " + string(l))
	}
	return s
}

// Раздел навигации. Метки живут здесь, а не в nav.go, чтобы у языка был один
// источник и не пришлось искать строки по трём файлам.
var sections = map[string]map[Lang]string{
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

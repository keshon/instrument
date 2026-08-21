package i18n

// Lang — код языка. Пустая строка невозможна: язык всегда явный.
type Lang string

const (
	RU Lang = "ru"
	EN Lang = "en"
)

var All = []Lang{RU}

var Known = []Lang{RU, EN}

func (l Lang) Prefix() string {
	if l == RU {
		return ""
	}
	return "/" + string(l)
}

func (l Lang) Suffix() string {
	if l == RU {
		return ""
	}
	return "." + string(l)
}

func (l Lang) Label() string {
	switch l {
	case EN:
		return "English"
	default:
		return "Русский"
	}
}

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
	// Пункт масштаба подписан кеглем, а высота контрола уходит в пояснение:
	// «крупно» ни о чём не сообщает, пока не нажато, «16px» — сообщает сразу.
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
	"untranslated": {
		RU: "",
		EN: "This page is not translated yet — the text below is in Russian. " +
			"The API reference is in English either way: class, token and " +
			"attribute names are the same in both languages.",
	},
}

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
		panic("i18n: нет заголовка раздела " + id)
	}
	s, ok := m[l]
	if !ok {
		panic("i18n: нет перевода раздела " + id + " на " + string(l))
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

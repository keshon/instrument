package content

import "strings"

type SectionDef struct {
	ID      string
	Order   int
	Aliases []string
}

var sectionDefs = []SectionDef{
	{"install", 10, []string{"Установка", "Installation"}},
	{"usage", 20, []string{"Использование", "Usage", "Разметка", "Разметка обязательна"}},
	{"contract", 25, []string{"Контракт", "Contract", "Контракт разметки", "Markup contract"}},
	{"when", 30, []string{"Когда использовать", "When to use"}},
	{"anatomy", 40, []string{"Устройство", "Anatomy", "Части", "Пункты"}},
	{"scale", 45, []string{"Шкала", "Scale", "Значения", "Values"}},
	{"variants", 50, []string{"Варианты", "Variants", "Модификаторы", "Тон", "Тона", "Темы"}},
	{"sizes", 60, []string{"Размеры", "Sizes"}},
	{"states", 70, []string{"Состояния", "States"}},
	{"icons", 80, []string{"С иконкой", "With icon", "Иконки"}},
	{"behavior", 90, []string{"Поведение", "Behavior", "Behaviour"}},
	{"js", 95, []string{"JS", "JavaScript", "Скрипт"}},
	{"composition", 100, []string{"Композиции", "Composition", "Композиция"}},
	{"patterns", 110, []string{"Сценарии", "Patterns", "Собранный экран"}},
	{"rules", 120, []string{"Правила", "Rules"}},
	{"a11y", 130, []string{"Доступность", "Accessibility"}},
	{"customization", 140, []string{"Настройка", "Customization", "Свой вариант"}},
	// The markup contract absorbs what «Использование» and «Доступность» used to
	// keep apart: both said the same thing — what has to be written for the
	// component to work rather than to lie — and across fifty pages that came to
	// 1 480 lines of one rule told twice.
	{"api", 150, []string{"API", "Справочник", "API Reference"}},
	{"related", 160, []string{"Связанное", "Related"}},
}

var byAlias = func() map[string]*SectionDef {
	m := map[string]*SectionDef{}
	for i := range sectionDefs {
		d := &sectionDefs[i]
		for _, a := range d.Aliases {
			m[normTitle(a)] = d
		}
	}
	return m
}()

var byID = func() map[string]*SectionDef {
	m := map[string]*SectionDef{}
	for i := range sectionDefs {
		m[sectionDefs[i].ID] = &sectionDefs[i]
	}
	return m
}()

func normTitle(s string) string {
	return strings.ToLower(strings.Join(strings.Fields(s), " "))
}

func LookupSection(title string) (*SectionDef, bool) {
	d, ok := byAlias[normTitle(title)]
	return d, ok
}

func SectionByID(id string) (*SectionDef, bool) {
	d, ok := byID[id]
	return d, ok
}

type Section struct {
	ID    string
	Title string
	Order int
	Known bool
}

package content

import (
	"fmt"
	"regexp"
	"strings"

	"instrument/site/internal/i18n"
)

var (
	directiveOpen  = regexp.MustCompile(`^:::([a-z][a-z-]*)\s*(.*)$`)
	directiveClose = regexp.MustCompile(`^:::\s*$`)
	fenceRe        = regexp.MustCompile("^(```|~~~)")
)

var ruleKinds = map[string]bool{"do": true, "dont": true}

// Тон сноски: «так» — это подтверждение, «не так» — ошибка. Оба значка уже
// есть у кита, и заводить для документации свои значило бы рисовать второй
// набор для той же пары смыслов.
var ruleTones = map[string]string{"do": "ok", "dont": "error"}

var noteTones = map[string]string{
	"note":   "info",
	"warn":   "warn",
	"danger": "error",
	"ok":     "ok",
}

func expandDirectives(body []byte, lang i18n.Lang) ([]byte, error) {
	lines := strings.Split(string(body), "\n")
	var out []string

	inFence := false
	rowOpen := false

	type openBlock struct{ kind, closes string }
	var open []openBlock

	closeRow := func() {
		if rowOpen {
			out = append(out, "", "</div>", "")
			rowOpen = false
		}
	}

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		if fenceRe.MatchString(strings.TrimSpace(line)) {
			inFence = !inFence
			out = append(out, line)
			continue
		}
		if inFence {
			out = append(out, line)
			continue
		}

		if m := directiveOpen.FindStringSubmatch(strings.TrimSpace(line)); m != nil {
			kind, head := m[1], strings.TrimSpace(m[2])
			var closes string
			switch {
			case ruleKinds[kind]:
				// Пара правил — это БАННЕР, а не сноска с дописанной шапкой.
				//
				// У баннера кита ровно та структура, которую сайт до этого
				// изображал вручную: значок от тона, заголовок, текст, рамка
				// и тонированная заливка. Сноска же рассчитана на строку в
				// потоке — заголовка у неё нет, и его приходилось подставлять
				// своим классом.
				//
				// Цена прежней сборки читалась на экране: значок говорил «так»
				// формой, плашка повторяла это словом, заголовок — фразой. Три
				// пересказа одного на блок в три строки. Теперь слово живёт в
				// заголовке баннера, а форму по-прежнему несёт значок: закон
				// кита «цвет не единственный носитель» выполняется, а повтора
				// нет.
				//
				// Ряд — СВОЯ сетка, а не китовая. `.inst-grid` адаптивная:
				// auto-fit подбирает число колонок по ширине, и на широком
				// мониторе даёт три. Пары идут в разметке подряд, поэтому ряд
				// поехал бы так: do · dont · do / dont · do · dont — сравнение,
				// ради которого правила и стоят рядом, сломалось бы именно
				// там, где места больше всего.
				//
				// Примитива «ровно две равные колонки» в ките нет, и заводить
				// его незачем: это раскладка СТАТЬИ справочника, а не
				// компонент. Сайту она принадлежит по тому же правилу, по
				// которому ему принадлежит отступ вокруг примера.
				if !rowOpen {
					out = append(out, "", `<div class="rule-row">`, "")
					rowOpen = true
				}
				label := i18n.T(lang, "rule."+kind)
				out = append(out,
					fmt.Sprintf(`<div class="inst-banner rule" data-tone="%s">`, ruleTones[kind]),
					`<div class="inst-banner-body">`,
					"",
					fmt.Sprintf(`<p class="inst-banner-title">%s</p>`, escape(ruleTitle(label, head))),
					`<div class="inst-banner-text">`,
					"")
				closes = "</div></div></div>"
			case noteTones[kind] != "":
				closeRow()
				out = append(out, "",
					fmt.Sprintf(`<div class="inst-note" data-tone="%s"><div class="note-body">`, noteTones[kind]),
					"")
				closes = "</div></div>"
			default:
				return nil, fmt.Errorf("строка %d: неизвестный блок :::%s (есть do · dont · note · warn · danger · ok)", i+1, kind)
			}
			open = append(open, openBlock{kind, closes})
			continue
		}

		if directiveClose.MatchString(strings.TrimSpace(line)) {
			if len(open) == 0 {
				return nil, fmt.Errorf("строка %d: закрытие ::: без открытия", i+1)
			}
			b := open[len(open)-1]
			open = open[:len(open)-1]
			out = append(out, "", b.closes, "")
			if !ruleKinds[b.kind] {
				closeRow()
			}
			continue
		}

		if rowOpen && len(open) == 0 && strings.TrimSpace(line) != "" {
			closeRow()
		}
		out = append(out, line)
	}

	if len(open) > 0 {
		return nil, fmt.Errorf("блок :::%s не закрыт", open[len(open)-1].kind)
	}
	closeRow()
	return []byte(strings.Join(out, "\n")), nil
}

// Слово и фраза в ОДНОМ заголовке, а не в двух строках.
//
// «Так» без продолжения — не заголовок, а метка; «Слайдер для приблизительного»
// без слова теряет то, что кит требует от статуса: носитель помимо цвета.
// Тире соединяет их в одну фразу, которую читают целиком.
func ruleTitle(label, head string) string {
	if head == "" {
		return label
	}
	return label + " — " + head
}

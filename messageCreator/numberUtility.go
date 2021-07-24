package messageCreator

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Number int

func (num Number) AddComma() String {
	p := message.NewPrinter(language.English)
	return String(p.Sprintf("%d", num))
}

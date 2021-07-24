package messageCreator

import (
	"fmt"
)

type String string

// todo: write test for these functions

func (str String) Bold() String {
	return String(fmt.Sprintf("<b>%s</b>", str))
}

func (str String) Strike() String {
	return String(fmt.Sprintf("<s>%s</s>", str))
}

func (str String) UnderLine() String {
	return String(fmt.Sprintf("<u>%s</u>", str))
}

func (str String) ToLink(url string) String {
	return String(fmt.Sprintf(`<a href="%s">%s</a>`, url, str))
}

func (str String) AddNewLine() String {
	return String(fmt.Sprintf("%+v \n", str))
}

func (str String) Append(str2 string) String {
	return String(fmt.Sprintf("%+v %+v", str, str2))
}

func (str String) ToString() string {
	return string(str)
}

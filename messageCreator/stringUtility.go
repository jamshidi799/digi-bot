package messageCreator

import (
	"fmt"
	"regexp"
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

func FixNumber(str string) int {
	re := regexp.MustCompile(`[^۰۱۲۳۴۵۶۷۸۹]`)
	return convertToEnglishNumber(re.ReplaceAllString(str, ""))
}

func convertToEnglishNumber(number string) int {
	enNum := 0
	dict := map[rune]int{'۱': 1, '۲': 2, '۳': 3, '۴': 4, '۵': 5, '۶': 6, '۷': 7, '۸': 8, '۹': 9, '۰': 0}
	for _, char := range number {
		enNum = enNum*10 + dict[char]
	}

	return enNum
}

func CleaningString(str string) string {
	re := regexp.MustCompile(`\n`)
	str = re.ReplaceAllString(str, " ")
	re = regexp.MustCompile(` {2,}`)
	return re.ReplaceAllString(str, " ")
}

func CleaningStringWithDelimiter(str string, delimiter string) string {
	re := regexp.MustCompile(`\n`)
	str = re.ReplaceAllString(str, delimiter)
	re = regexp.MustCompile(` {2,}`)
	return re.ReplaceAllString(str, delimiter)
}

package util

import "strings"

type StringUtil struct{}

var String StringUtil = StringUtil{}

func (_ StringUtil) IsEmpty(str string) bool {
	return strings.TrimSpace(str) == ""
}

func GetExcelColumnName(columnNumber int) string {
	dividend := columnNumber
	var columnName string
	var modulo int

	for dividend > 0 {
		modulo = (dividend - 1) % 26
		columnName = toCharStr(modulo) + columnName
		dividend = int((dividend - modulo) / 26)
	}

	return columnName
}

func toCharStr(i int) string {
	return string(rune('A' + i))
}

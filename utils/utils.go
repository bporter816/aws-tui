package utils

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	DefaultTimeFormat = "2006-01-02 15:04:05"
)

var (
	titleCaser = cases.Title(language.English)
	upperCaser = cases.Upper(language.English)
)

func TitleCase(str string) string {
	return titleCaser.String(str)
}

func UpperCase(str string) string {
	return upperCaser.String(str)
}

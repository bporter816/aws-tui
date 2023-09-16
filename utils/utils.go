package utils

import (
	"fmt"
	"math"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	DefaultTimeFormat = "2006-01-02 15:04:05"
)

var (
	titleCaser = cases.Title(language.English)
	upperCaser = cases.Upper(language.English)

	abbreviations = map[string]string{
		"Milliseconds": "ms",
		"Seconds":      "sec",
		"Minutes":      "min",
		"Bytes":        "B",
		"Kilobytes":    "KB",
		"Megabytes":    "MB",
		"Gigabytes":    "GB",
	}
)

func TitleCase(str string) string {
	return titleCaser.String(str)
}

func UpperCase(str string) string {
	return upperCaser.String(str)
}

func SimplifyFloat(value float64) string {
	if value == math.Trunc(value) {
		return fmt.Sprintf("%v", int(value))
	} else {
		return fmt.Sprintf("%f", value)
	}
}

func AbbreviateUnit(unit string) string {
	if u, ok := abbreviations[unit]; ok {
		return u
	} else {
		return unit
	}
}

func BoolToString(b bool, y string, n string) string {
	if b {
		return y
	} else {
		return n
	}
}

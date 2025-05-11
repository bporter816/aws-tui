package utils

import (
	"fmt"
)

// TODO add tests
func FormatServiceQuotasValue(value *float64, unit *string) string {
	if value == nil {
		return ""
	}

	v := SimplifyFloat(*value)

	if unit != nil && *unit != "None" {
		return fmt.Sprintf("%v %v", v, AbbreviateUnit(*unit))
	} else {
		return v
	}
}

package utils

import "strings"

func ReplaceDotWithDash(s string) string {
	return strings.Replace(s, ".", "-", -1)
}

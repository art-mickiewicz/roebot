package util

import (
	"strings"
)

func TrimLine(s string, upTo int) string {
	var s1 string
	s1 = strings.ReplaceAll(s, "\n", " ")
	s1 = strings.ReplaceAll(s1, "\r", " ")
	s1 = strings.ReplaceAll(s1, "  ", " ")
	r := []rune(s1)

	if len(r) > upTo {
		return string(r[:upTo]) + "..."
	} else {
		return s1
	}
}

func PadLine(s string, upTo int, with string) string {
	ret := []rune(s)
	for len(ret) < upTo {
		ret = append(ret, []rune(with)...)
	}
	return string(ret)
}

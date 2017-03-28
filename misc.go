package main

import (
	"strings"
	"time"
	"unicode"
)

func splitStringByCase(splitString string) (result string) {
	l := 0
	for s := splitString; s != ""; s = s[l:] {
		l = strings.IndexFunc(s[1:], unicode.IsUpper) + 1
		if l <= 0 {
			l = len(s)
		}
		if result == "" {
			result = s[:l]
		} else {
			result += " " + s[:l]
		}
	}
	return
}

func inTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}

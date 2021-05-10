package main

import (
	"regexp"
	"strings"
)

func Filter(word string) string {
	// reject opening and closing HTML tags
	if regexp.MustCompile(`^<.*|\/>$`).MatchString(word) {
		return ""
	}

	// reject URLs
	if regexp.MustCompile(`^https?:\/\/`).MatchString(word) {
		return ""
	}

	// reject HTML attributes
	if regexp.MustCompile(`^.*=`).MatchString(word) {
		return ""
	}

	// reject markdown horizontal rules
	if regexp.MustCompile(`^[\-]{2,3}$`).MatchString(word) {
		return ""
	}

	// reject paths
	if regexp.MustCompile(`^(\/.*?)+$`).MatchString(word) {
		return ""
	}

	word = strings.ToLower(word)

	// strip out illegal characters
	word = regexp.MustCompile(`[~“”‘’–…\/\{\}\[\]\$\+\%\\\*#()"_?!,.]`).ReplaceAllString(word, "")

	// strip out back ticks (nice one, go)
	word = regexp.MustCompile("`").ReplaceAllString(word, "")

	// reject times
	if regexp.MustCompile(`^[0-9]{1,2}[:]?[0-9]{0,2}(am|pm|AM|PM)?$`).MatchString(word) {
		return ""
	}

	// wrapped single quotes, e.g. "He said 'hello',"
	word = regexp.MustCompile(`^'(.*?)'$`).ReplaceAllString(word, "$1")

	// leading single quotes, e.g. 'Hello
	word = regexp.MustCompile(`^'(.*?)$`).ReplaceAllString(word, "$1")

	// check for dangling symbols
	word = regexp.MustCompile(`^[\|à§\/\-\+%–&<>]$`).ReplaceAllString(word, "")

	// check for dangling numbers
	word = regexp.MustCompile(`^\d+$`).ReplaceAllString(word, "")

	return word
}

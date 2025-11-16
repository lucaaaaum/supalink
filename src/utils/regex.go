package utils

import "strings"

const RegexConstants = ".*+?[]()|{}"

func AddStopSuffixToPattern(pattern *string) {
	if !strings.HasSuffix(*pattern, "$") {
		*pattern += "$"
	}
}

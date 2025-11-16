package utils

import "strings"

func AddStopSuffixToPattern(pattern *string) {
	if !strings.HasSuffix(*pattern, "$") {
		*pattern += "$"
	}
}

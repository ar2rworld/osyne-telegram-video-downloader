package cleaner

import "strings"

func CleanURL(s string) string {
	withoutNewLine := strings.ReplaceAll(s, "\n", " ")
	stringParts := strings.Split(withoutNewLine, " ")
	var out string
	for _, part := range stringParts {
		if strings.Contains(part, "https") && strings.Contains(part, ".com") {
			out = part
			break
		}
	}
	return out
}

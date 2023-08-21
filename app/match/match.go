package match

import (
	"regexp"
)

func Match(s string) string {
	pattern := `(https?:\/\/)?(www|vm)?\.?(youtube\.com\/shorts|twitter\.com|tiktok\.com|instagram\.com)[^\s]+`
	re := regexp.MustCompile(pattern)
	matches := re.FindAllString(s, -1)
	if len(matches) > 0 {
		return matches[0]
	}
	return ""
}

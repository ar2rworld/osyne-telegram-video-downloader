package match

import (
	"regexp"
)

func Match(s string) string {
	pattern := `(https?:\/\/)?(www|vm)?\.?(youtube|youtu(\.com|\.be)|twitter\.com|tiktok\.com|instagram\.com)[^\s]+`
	re := regexp.MustCompile(pattern)
	matches := re.FindAllString(s, -1)
	if len(matches) > 0 {
		return matches[0]
	}
	return ""
}

// Match only youtube video, but don't match shorts
func Youtube(s string) string {
	pattern := `(https?:\/\/)?(www)?\.?(youtube|youtu)(\.com|\.be)?\/([^s][^h][^o][^r][^t][^s])[^\s]+`
	re := regexp.MustCompile(pattern)
	matches := re.FindAllString(s, -1)
	if len(matches) > 0 {
		return matches[0]
	}
	return ""
}

package match

import (
	"regexp"
)

func Match(s string) string {
	pattern := `(https?:\/\/)?(www|vm)?\.?(youtube\.com\/(?:watch\?v=|embed\/|v\/|shorts\/)|youtu\.be\/|twitter\.com\/|tiktok\.com\/|instagram\.com\/)[^\s]+`
	re := regexp.MustCompile(pattern)
	matches := re.FindAllString(s, -1)
	if len(matches) > 0 {
		return matches[0]
	}
	return ""
}

func DownloadSectionsArgument(s string) string {
	re := regexp.MustCompile(`\s+`)
	s = re.ReplaceAllString(s, " ")
	pattern := `(-s\s+\*\d+:\d+-\d+:\d+)|(-s\s+)`
	re = regexp.MustCompile(pattern)
	matches := re.FindAllString(s, -1)
	if len(matches) > 0 {
		return matches[0]
	}
	return ""
}

func DownloadAudioArgument(s string) string {
	// Pattern to match -x argument
	pattern := `(?:^|\s)(-x)(?:\s|$)`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(s)

	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}

// Match only youtube video, but don't match shorts
func Youtube(s string) string {
	pattern := `(https?:\/\/)?(www\.)?(youtube\.com\/(?:watch\?v=|embed\/|v\/)|youtu\.be\/)([a-zA-Z0-9_-]{11})`
	re := regexp.MustCompile(pattern)
	matches := re.FindAllString(s, -1)
	if len(matches) > 0 {
		return matches[0]
	}
	return ""
}

// Match only youtube shorts
func YoutubeShorts(s string) string {
	pattern := `(?i)(https?:\/\/)?(www\.)?(youtube|youtu)(\.com|\.be)?\/shorts\/[^\s]+`
	re := regexp.MustCompile(pattern)
	matches := re.FindString(s)
	return matches
}

// Match only instagram content
func Instagram(s string) string {
	pattern := `(https?:\/\/)?(www)?\.?(instagram\.com)[^\s]+`
	re := regexp.MustCompile(pattern)
	matches := re.FindAllString(s, -1)
	if len(matches) > 0 {
		return matches[0]
	}
	return ""
}

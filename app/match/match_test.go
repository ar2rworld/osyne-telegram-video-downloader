package match

import "testing"

func TestMatch(t *testing.T) {
	MatchingCases := []struct {
		Name string
		Text string
		URL  string
	}{
		{"Empty", "Some text", ""},
		{"Only isntagram", "https://www.instagram.com/reel/Cv-hdiDt9Ix/", "https://www.instagram.com/reel/Cv-hdiDt9Ix/"},
		{"Only tiktok", "https://vm.tiktok.com/ZM2KGqk1v/", "https://vm.tiktok.com/ZM2KGqk1v/"},
		{"Only twitter", "https://twitter.com/webflite/status/1692079842689159520?s=20", "https://twitter.com/webflite/status/1692079842689159520?s=20"},
		{"Only youtubeshorts", "https://youtube.com/shorts/G90KEDm_G28?feature=share", "https://youtube.com/shorts/G90KEDm_G28?feature=share"},
		{"Without https", "Some text youtube.com/shorts/G90KEDm_G28?feature=share", "youtube.com/shorts/G90KEDm_G28?feature=share"},
	}

	for _, matchingCase := range MatchingCases {
		t.Run(matchingCase.Name, func(tt *testing.T) {
			got := Match(matchingCase.Text)
			want := matchingCase.URL
			if got != want {
				tt.Errorf("Couldn't match: \"%s\" -> \"%s\" , got:\"%s\"", matchingCase.Text, want, got)
			}
		})
	}
}

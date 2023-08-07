package tests

import (
	"testing"
	
	"ar2rworld/golang-telegram-video-downloader/internal/cleaner"
)

func TestClearUrl(t *testing.T) {
	type urlCase struct {
		gotMessage string
		wantUrl string
	}
	cases := []urlCase{
		{
			gotMessage: "Omg look at this empty input!",
			wantUrl: "",
		},
		{
			gotMessage: "Omg look at this! https://youtube.com/shorts/runningonwaterbird",
			wantUrl: "https://youtube.com/shorts/runningonwaterbird",
		},
		{
			gotMessage: "Omg!\nhttps://youtube.com/shorts/just-an-owl",
			wantUrl: "https://youtube.com/shorts/just-an-owl",
		},
		{
			gotMessage: "wowwwwwo!\nhttps://instagram.com/bird-does-fishing-with-bread i cannot imaging to...",
			wantUrl: "https://instagram.com/bird-does-fishing-with-bread",
		},
		{
			gotMessage: "u need to see this! https://tiktok.com/doka2-funny-moments\namazing!!! \n gg wp",
			wantUrl: "https://tiktok.com/doka2-funny-moments",
		},
		{
			gotMessage: "Omg look at this! https://tiktok.com/biggest-rc-helicopter i m very impressed by this helicopter",
			wantUrl: "https://tiktok.com/biggest-rc-helicopter",
		},
		{
			gotMessage: "https://twitter.com/some-politics democracy is not perfect but it is the best we can do",
			wantUrl: "https://twitter.com/some-politics",
		},
		{
			gotMessage: "https://twitter.com/how-to-change-political-regime-in-your-country\nAn interesting point u got",
			wantUrl: "https://twitter.com/how-to-change-political-regime-in-your-country",
		},
	}
	for _, c := range(cases) {
		t.Run(c.gotMessage, func(t *testing.T) {
			cleaned := cleaner.CleanUrl(c.gotMessage)
			if cleaned != c.wantUrl {
				t.Errorf("Failed to clean url:%s\n!=\n%s", c.wantUrl, cleaned)
			}
		})
	}
}

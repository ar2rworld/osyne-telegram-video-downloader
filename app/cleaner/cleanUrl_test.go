package cleaner

import "testing"

func TestClearUrl(t *testing.T) {
	type urlCase struct {
		gotMessage string
		wantURL    string
	}
	cases := []urlCase{
		{
			gotMessage: "Omg look at this empty input!",
			wantURL:    "",
		},
		{
			gotMessage: "Omg look at this! https://youtube.com/shorts/runningonwaterbird",
			wantURL:    "https://youtube.com/shorts/runningonwaterbird",
		},
		{
			gotMessage: "wowwwwwo!\nhttps://instagram.com/bird-does-fishing-with-bread i cannot imaging to...",
			wantURL:    "https://instagram.com/bird-does-fishing-with-bread",
		},
		{
			gotMessage: "u need to see this! https://tiktok.com/doka2-funny-moments\namazing!!! \n gg wp",
			wantURL:    "https://tiktok.com/doka2-funny-moments",
		},
		{
			gotMessage: "Omg look at this! https://tiktok.com/biggest-rc-helicopter i m very impressed by this helicopter",
			wantURL:    "https://tiktok.com/biggest-rc-helicopter",
		},
		{
			gotMessage: "https://twitter.com/some-politics democracy is not perfect but it is the best we can do",
			wantURL:    "https://twitter.com/some-politics",
		},
		{
			gotMessage: "https://twitter.com/how-to-change-political-regime-in-your-country\nAn interesting point u got",
			wantURL:    "https://twitter.com/how-to-change-political-regime-in-your-country",
		},
		{
			gotMessage: "https://vm.tiktok.com/somevideoid/\n@a или мальчики из америки, что такое \"show reading\"?",
			wantURL:    "https://vm.tiktok.com/somevideoid/",
		},
		{
			gotMessage: "Dua Lipa's emoji — https://www.youtube.com/shorts/zrws7lzoQJQ",
			wantURL:    "https://www.youtube.com/shorts/zrws7lzoQJQ",
		},
		{
			gotMessage: "https://google.com/search?q=what\nshould be the message with two links\nhttps://www.tiktok.com/someInterestingId",
			wantURL:    "https://www.tiktok.com/someInterestingId",
		},
		{
			gotMessage: "https://www.tiktok.com/someInterestingId\nshould be the message with two links but reversed\nhttps://google.com/search?q=what",
			wantURL:    "https://www.tiktok.com/someInterestingId",
		},
	}
	for _, c := range cases {
		t.Run(c.gotMessage, func(t *testing.T) {
			cleaned := CleanURL(c.gotMessage)
			if cleaned != c.wantURL {
				t.Errorf("Failed to clean url:%s\n!=\n%s", c.wantURL, cleaned)
			}
		})
	}
}

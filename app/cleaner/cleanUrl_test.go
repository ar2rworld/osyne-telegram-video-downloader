package cleaner

import "testing"

func TestClearUrl(t *testing.T) {
	type urlCase struct {
		gotMessage string
		wantUrl    string
	}
	cases := []urlCase{
		{
			gotMessage: "Omg look at this empty input!",
			wantUrl:    "",
		},
		{
			gotMessage: "Omg look at this! https://youtube.com/shorts/runningonwaterbird",
			wantUrl:    "https://youtube.com/shorts/runningonwaterbird",
		},
		{
			gotMessage: "Omg!\nhttps://youtube.com/shorts/just-an-owl",
			wantUrl:    "https://youtube.com/shorts/just-an-owl",
		},
		{
			gotMessage: "wowwwwwo!\nhttps://instagram.com/bird-does-fishing-with-bread i cannot imaging to...",
			wantUrl:    "https://instagram.com/bird-does-fishing-with-bread",
		},
		{
			gotMessage: "u need to see this! https://tiktok.com/doka2-funny-moments\namazing!!! \n gg wp",
			wantUrl:    "https://tiktok.com/doka2-funny-moments",
		},
		{
			gotMessage: "Omg look at this! https://tiktok.com/biggest-rc-helicopter i m very impressed by this helicopter",
			wantUrl:    "https://tiktok.com/biggest-rc-helicopter",
		},
		{
			gotMessage: "https://twitter.com/some-politics democracy is not perfect but it is the best we can do",
			wantUrl:    "https://twitter.com/some-politics",
		},
		{
			gotMessage: "https://twitter.com/how-to-change-political-regime-in-your-country\nAn interesting point u got",
			wantUrl:    "https://twitter.com/how-to-change-political-regime-in-your-country",
		},
		{
			gotMessage: "https://vm.tiktok.com/somevideoid/\n@a или мальчики из америки, что такое \"show reading\"?",
			wantUrl:    "https://vm.tiktok.com/somevideoid/",
		},
		{
			gotMessage: "Dua Lipa's emoji — https://www.youtube.com/shorts/zrws7lzoQJQ",
			wantUrl:    "https://www.youtube.com/shorts/zrws7lzoQJQ",
		},
		{
			gotMessage: "https://google.com/search?q=what\nshould be the message with two links\nhttps://www.tiktok.com/someInterestingId",
			wantUrl:    "https://www.tiktok.com/someInterestingId",
		},
		{
			gotMessage: "https://www.tiktok.com/someInterestingId\nshould be the message with two links but reversed\nhttps://google.com/search?q=what",
			wantUrl:    "https://www.tiktok.com/someInterestingId",
		},
	}
	for _, c := range cases {
		t.Run(c.gotMessage, func(t *testing.T) {
			cleaned := CleanURL(c.gotMessage)
			if cleaned != c.wantUrl {
				t.Errorf("Failed to clean url:%s\n!=\n%s", c.wantUrl, cleaned)
			}
		})
	}
}

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
			gotMessage: "https://instagram.com/bird-does-fishing-with-bread i cannot imaging to...",
			wantURL:    "https://instagram.com/bird-does-fishing-with-bread",
		},
		{
			gotMessage: "u need to see this! https://tiktok.com/doka2-funny-moments amazing!!! gg wp",
			wantURL:    "https://tiktok.com/doka2-funny-moments",
		},
		{
			gotMessage: "Omg look at this!\nhttps://tiktok.com/biggest-rc-helicopter\ni m very impressed by this helicopter",
			wantURL:    "https://tiktok.com/biggest-rc-helicopter",
		},
		{
			gotMessage: "https://twitter.com/how-to-change-political-regime-in-your-country\nAn interesting point u got",
			wantURL:    "https://twitter.com/how-to-change-political-regime-in-your-country",
		},
		{
			gotMessage: "Dua Lipa's emoji — https://www.youtube.com/shorts/zrws7lzoQJQ",
			wantURL:    "https://www.youtube.com/shorts/zrws7lzoQJQ",
		},
		{
			gotMessage: "https://google.com/search?q=what THE WRONG TEST, CLEANER SHOULD PICK UP THE TIKTOK LINK OR NOT? https://www.tiktok.com/someInterestingId",
			wantURL:    "https://google.com/search?q=what",
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

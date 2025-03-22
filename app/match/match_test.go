package match

import (
	"testing"
)

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
		{"Only x.com", "https://x.com/KeyDatch/status/1903332041124368783", "https://x.com/KeyDatch/status/1903332041124368783"},
		{"Some other blahblahx.com", "blahblahx.com", ""},
		{"Only youtubeshorts", "https://youtube.com/shorts/G90KEDm_G28?feature=share", "https://youtube.com/shorts/G90KEDm_G28?feature=share"},
		{"Without https", "Some text youtube.com/shorts/G90KEDm_G28?feature=share", "youtube.com/shorts/G90KEDm_G28?feature=share"},
		{"Short Youtube link", "Some text https://youtu.be/rfwnQzS9KkA", "https://youtu.be/rfwnQzS9KkA"},
		{"Youtube video", "Some text https://www.youtube.com/watch?v=rfwnQzS9KkA", "https://www.youtube.com/watch?v=rfwnQzS9KkA"},
		{"Single youtube word", "https://gist.ly/youtube-summarizer/effective-communication-strategies-lessons-from-ilya-yashin", ""},
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

func TestMatchYoutube(t *testing.T) {
	MatchingCases := []struct {
		Name string
		Text string
		URL  string
	}{
		{"Empty", "Some text", ""},
		{"Don't match tiktok", "https://vm.tiktok.com/ZM2KGqk1v/", ""},
		{"Only youtubeshorts", "https://youtube.com/shorts/G90KEDm_G28?feature=share", ""},
		{"Without https", "Some text youtube.com/shorts/G90KEDm_G28?feature=share", ""},
		{"Youtube video", "Some text https://www.youtube.com/watch?v=rfwnQzS9KkA", "https://www.youtube.com/watch?v=rfwnQzS9KkA"},
		{"Short link", "Some text https://youtu.be/rfwnQzS9KkA", "https://youtu.be/rfwnQzS9KkA"},
		{"Single youtube word", "youtube", ""},
		{"Single youtu.be word", "youtu.be", ""},
		{"", "https://youtu.be/GG5u0sUe7sQ", "https://youtu.be/GG5u0sUe7sQ"},
		{"Link with single word youtube", "https://gist.ly/youtube-summarizer/effective-communication-strategies-lessons-from-ilya-yashin", ""},
	}

	for _, matchingCase := range MatchingCases {
		t.Run(matchingCase.Name, func(tt *testing.T) {
			got := Youtube(matchingCase.Text)
			want := matchingCase.URL
			if got != want {
				tt.Errorf("Couldn't match: \"%s\" -> \"%s\" , got:\"%s\"", matchingCase.Text, want, got)
			}
		})
	}
}

func TestMatchYoutubeShorts(t *testing.T) {
	MatchingCases := []struct {
		Name string
		Text string
		URL  string
	}{
		{"Empty", "Some text", ""},
		{"Don't match tiktok", "https://vm.tiktok.com/ZM2KGqk1v/", ""},
		{"Only youtubeshorts", "https://youtube.com/shorts/G90KEDm_G28?feature=share", "https://youtube.com/shorts/G90KEDm_G28?feature=share"},
		{"Without https", "Some text youtube.com/shorts/G90KEDm_G28?feature=share", "youtube.com/shorts/G90KEDm_G28?feature=share"},
		{"Youtube video", "Some text https://www.youtube.com/watch?v=rfwnQzS9KkA", ""},
		{"Short link", "Some text https://youtu.be/rfwnQzS9KkA", ""},
		{"http youtubeshorts no params", "https://youtube.com/shorts/id", "https://youtube.com/shorts/id"},
	}

	for _, matchingCase := range MatchingCases {
		t.Run(matchingCase.Name, func(tt *testing.T) {
			got := YoutubeShorts(matchingCase.Text)
			want := matchingCase.URL
			if got != want {
				tt.Errorf("Couldn't match: \"%s\" -> \"%s\" , got:\"%s\"", matchingCase.Text, want, got)
			}
		})
	}
}

func TestDownloadSectionsArgument(t *testing.T) { //nolint: all
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "s and url",
			args: args{s: "-s   *0:0-0:10   https://youtube.com/watch?v=dQw4w9WgXcQ&asjfse=1"},
			want: "-s *0:0-0:10",
		},
		{
			name: "s with spaces and url",
			args: args{s: `-s   *0:0-0:10   https://youtube.com/watch?v=dQw4w9WgXcQ&asjfse=1`},
			want: "-s *0:0-0:10",
		},
		{
			name: "many digits",
			args: args{s: `-s    *000:0000-0000:0001000  https://youtube.com/watch?v=dQw4w9WgXcQ&asjfse=1`},
			want: "-s *000:0000-0000:0001000",
		},
		{
			name: "onle s with single space",
			args: args{s: `-s `},
			want: "-s ",
		},
		{
			name: "some digits",
			args: args{s: `-s  *100:050-100:100 https://youtube.com/watch?v=dQw4w9WgXcQ&asjfse=1`},
			want: "-s *100:050-100:100",
		},
		{
			name: "missing space",
			args: args{s: `-s*100:050-100:100 https://youtube.com/watch?v=dQw4w9WgXcQ&asjfse=1`},
			want: "",
		},
		{
			name: "some text and s",
			args: args{s: `asfheuihfsi -re -erf eje -s *0:0-0:10 https://youtube.com/watch?v=dQw4w9WgXcQ&asjfse=1`},
			want: "-s *0:0-0:10",
		},
		{
			name: "url and s order",
			args: args{s: `asfheuihfsi -re -erf eje https://youtube.com/watch?v=dQw4w9WgXcQ&asjfse=1 -s *0:0-0:10 asfieuheiufa`},
			want: "-s *0:0-0:10",
		},
		{
			name: "url and s with some bad input",
			args: args{s: `asfheuihfsi -re -erf eje https://youtube.com/watch?v=dQw4w9WgXcQ&comment=rickroll -s *0:0-0:10 -rm -rf / fjsoe asfieuheiufa`},
			want: "-s *0:0-0:10",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DownloadSectionsArgument(tt.args.s); got != tt.want {
				t.Errorf(`DownloadSectionsArgument() = "%v", want "%v"`, got, tt.want)
			}
		})
	}
}

func TestDownloadAudioArgument(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			"x",
			"-x",
			"-x",
		},
		{
			"x with s",
			"-x -s *0:0-0:10",
			"-x",
		},
		{
			"x with s url",
			"-x -s *0:0-0:10 https://someurl.com",
			"-x",
		},
		{
			"s with x",
			"-s *0:0-0:10 -x",
			"-x",
		},
		{
			"s with x and url",
			"-s *0:0-0:10 -x https://someurl.com text",
			"-x",
		},
		{
			"s with x and url and rm -rf",
			"-s *0:0-0:10 -x https://someurl.com text rm -rf /",
			"-x",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DownloadAudioArgument(tt.input); got != tt.want {
				t.Errorf(`DownloadSectionsArgument() = "%v", want "%v"`, got, tt.want)
			}
		})
	}
}

func TestMatchInstagram(t *testing.T) {
	MatchingCases := []struct {
		Name string
		Text string
		URL  string
	}{
		{"Empty", "Some text", ""},
		{"Don't match tiktok", "https://vm.tiktok.com/ZM2KGqk1v/", ""},
		{"Don't match youtube", "https://youtube.com/shorts/G90KEDm_G28?feature=share", ""},
		{"Instagram reel", "https://www.instagram.com/reel/C6zLk3Op7b7/", "https://www.instagram.com/reel/C6zLk3Op7b7/"},
		{"Instagram reel no https", "www.instagram.com/reel/C6zLk3Op7b7/", "www.instagram.com/reel/C6zLk3Op7b7/"},
		{"Instagram reel no https no www", "instagram.com/reel/C6zLk3Op7b7/", "instagram.com/reel/C6zLk3Op7b7/"},
	}

	for _, matchingCase := range MatchingCases {
		t.Run(matchingCase.Name, func(tt *testing.T) {
			got := Instagram(matchingCase.Text)
			want := matchingCase.URL
			if got != want {
				tt.Errorf("Couldn't match: \"%s\" -> \"%s\" , got:\"%s\"", matchingCase.Text, want, got)
			}
		})
	}
}

func TestMatchFacebookReels(t *testing.T) {
	MatchingCases := []struct {
		Name string
		Text string
		URL  string
	}{
		{"Empty", "Some text", ""},
		{"Don't match tiktok", "https://vm.tiktok.com/ZM2KGqk1v/", ""},
		{"Don't match youtube", "https://youtube.com/shorts/G90KEDm_G28?feature=share", ""},
		{"Facebook reels direct link", "https://www.facebook.com/reel/643236801729548", "https://www.facebook.com/reel/643236801729548"},
		{"Facebook reels direct link with some text", "Hello guys, my dad just sent me this https://www.facebook.com/reel/643236801729548 amazing right?", "https://www.facebook.com/reel/643236801729548"},
		{"Facebook reels share option", "https://www.facebook.com/share/r/1DBVN8ZL9n/", "https://www.facebook.com/share/r/1DBVN8ZL9n/"},
		{"Facebook reels no https", "www.facebook.com/share/r/1DBVN8ZL9n/", "www.facebook.com/share/r/1DBVN8ZL9n/"},
		{"Facebook reels no https no www", "facebook.com/share/r/1DBVN8ZL9n/", "facebook.com/share/r/1DBVN8ZL9n/"},
	}

	for _, matchingCase := range MatchingCases {
		t.Run(matchingCase.Name, func(tt *testing.T) {
			got := FacebookReels(matchingCase.Text)
			want := matchingCase.URL
			if got != want {
				tt.Errorf("Couldn't match: \"%s\" -> \"%s\", got:\"%s\"", matchingCase.Text, want, got)
			}
		})
	}
}

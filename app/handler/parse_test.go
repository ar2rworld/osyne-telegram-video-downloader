package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parse(t *testing.T) { //nolint: all
	type want struct {
		s string
		x bool
	}
	tests := []struct {
		name    string
		args    string
		want    want
		wantErr bool
	}{
		{
			name:    "empty",
			args:    "just a link",
			want:    want{s: "", x: false},
			wantErr: false,
		},
		{
			name:    "test1",
			args:    "-s *0:0-0:61",
			want:    want{s: "*0:0-0:61", x: false},
			wantErr: false,
		},
		{
			name: "test2",
			args: "text",
			want: want{s: "", x: false},
		},
		{
			name:    "test3 only -s",
			args:    "-s",
			want:    want{s: "", x: false},
			wantErr: true,
		},
		{
			name: "test3 only -s with single space",
			args: "-s ",
			want: want{s: "", x: false},
		},
		{
			name: "test4",
			args: `-s *1111:00-1111:30`,
			want: want{s: "*1111:00-1111:30", x: false},
		},
		{
			name: "test5 extract audio",
			args: `-s *1111:00-1111:30 -x`,
			want: want{s: "*1111:00-1111:30", x: true},
		},
		{
			name: "test6 extract audio",
			args: `-x -s *1111:00-1111:30`,
			want: want{s: "*1111:00-1111:30", x: true},
		},
		{
			name: "test7 extract audio",
			args: `-x`,
			want: want{s: "*1111:00-1111:30", x: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parse(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got.Sections != nil {
				assert.Equal(t, tt.want.s, *got.Sections)
			}
			if got.ExtractAudio != nil {
				assert.Equal(t, tt.want.x, *got.ExtractAudio)
			}
		})
	}
}

func TestParseStartTime(t *testing.T) { //nolint:funlen
	testcases := []struct {
		name   string
		arg    string
		result string
	}{
		{
			"1 second",
			"https://youtu.be/tAleyWfgypM?t=1",
			"1",
		},
		{
			"1163 seconds",
			"https://youtu.be/tAleyWfgypM?t=1163",
			"1163",
		},
		{
			"youtube 90 with s",
			"https://www.youtube.com/watch?v=tAleyWfgypM&t=90s",
			"90",
		},
		{
			"youtube 90 without s",
			"https://www.youtube.com/watch?v=tAleyWfgypM&t=90",
			"90",
		},
		{
			"3089",
			"https://youtu.be/T_JKIkSf93Y?t=3089",
			"3089",
		},
		{
			"no current time",
			"https://youtu.be/T_JKIkSf93Y",
			"",
		},
		{
			"no current time youtube",
			"https://www.youtube.com/watch?v=T_JKIkSf93Y",
			"",
		},
		{
			"empty current time argument",
			"https://youtu.be/T_JKIkSf93Y?t=",
			"",
		},
		{
			"invalid current time argument",
			"https://youtu.be/T_JKIkSf93Y?t=asdf",
			"",
		},
		{
			"empty current time argument with other argument",
			"https://youtu.be/T_JKIkSf93Y?t=&a=1&b=asdf",
			"",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			r := parseCurrentTime(tc.arg)
			assert.Equal(t, tc.result, r)
		})
	}
}

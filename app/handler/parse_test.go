package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parse(t *testing.T) { //nolint: all
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "test1",
			args:    args{s: "-s *0:0-0:61"},
			want:    "*0:0-0:61",
			wantErr: false,
		},
		{
			name:    "test2",
			args:    args{s: "text"},
			want:    "*0:0-0:30",
			wantErr: false,
		},
		{
			name:    "test3 only -s",
			args:    args{s: "-s"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "test3 only -s with single space",
			args:    args{s: "-s "},
			want:    "*0:0-0:30",
			wantErr: false,
		},
		{
			name:    "test4",
			args:    args{s: `-s *1111:00-1111:30`},
			want:    "*1111:00-1111:30",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parse(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parse() = %v, want %v", got, tt.want)
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

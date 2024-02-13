package handler

import "testing"

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
			args:    args{s: "-s *0:0-0:61 https://youtube.com/watch?v=123456"},
			want:    "*0:0-0:61",
			wantErr: false,
		},
		{
			name:    "test2",
			args:    args{s: "text https://youtube.com/watch?v=123456"},
			want:    "*0:0-0:30",
			wantErr: false,
		},
		{
			name:    "test3",
			args:    args{s: "-s"},
			want:    "",
			wantErr: true,
		},
		{
			name: "test4 injection",
			args: args{s: `
		";
rm -rf / -s *1111:00-1111:30 https://www.youtube.com/watch?v=oENx7uPX-hc`},
			want:    "*0:0-0:30",
			wantErr: false,
		},
		// {
		// 	name:    "test5 injection",
		// 	args:    args{s: `-rm -rf / -s *1111:00-1111:30 https://www.youtube.com/watch?v=oENx7uPX-hc`},
		// 	want:    "*0:0-0:30",
		// 	wantErr: true,
		// },
		{
			name: "test6 injection",
			args: args{s: `);
			exec.Command("rm", "-rf", "/")
			-s *1111:00-1111:30 https://www.youtube.com/watch?v=oENx7uPX-hc`},
			want:    "*0:0-0:30",
			wantErr: false,
		},
		// {
		// 	name:    "test7 injection",
		// 	args:    args{s: `-s bash https://www.youtube.com/watch?v=oENx7uPX-hc`},
		// 	want:    "*0:0-0:30",
		// 	wantErr: false,
		// },
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

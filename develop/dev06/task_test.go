package main

import (
	"fmt"
	"testing"
)

func Test_cut(t *testing.T) {
	type args struct {
		file      string
		fields    string
		sep       bool
		delimeter string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "base01",
			args: args{
				file:      "resources/file.txt",
				fields:    "1,3",
				delimeter: " ",
			},
			want: "[[a c] [1 3] [2 4] [asdsvsvbefkbefb]]",
		},
		{
			name: "base02",
			args: args{
				file:      "resources/file.txt",
				fields:    "1,3",
				delimeter: " ",
				sep:       true,
			},
			want: "[[a c] [1 3] [2 4]]",
		},
		{
			name: "base03",
			args: args{
				file:      "resources/file.txt",
				fields:    "1-3",
				delimeter: " ",
				sep:       true,
			},
			want: "[[a b c] [1 2 3] [2 3 4]]",
		},
		{
			name: "base04",
			args: args{
				file:      "resources/file.txt",
				fields:    "2-",
				delimeter: " ",
				sep:       true,
			},
			want: "[[b c] [2 3] [3 4 5]]",
		},
		{
			name: "base05",
			args: args{
				file:      "resources/file.txt",
				fields:    "2-",
				delimeter: "d",
				sep:       true,
			},
			want: "[[svsvbefkbefb]]",
		},
		{
			name: "base06",
			args: args{
				file:      "resources/file.txt",
				fields:    "5",
				delimeter: "d",
				sep:       true,
			},
			want: "[[]]",
		},
	}
	for _, tt := range tests {
		a := cutArgs{}
		a.d = tt.args.delimeter
		a.s = tt.args.sep
		a.file = tt.args.file
		a.f = tt.args.fields
		t.Run(tt.name, func(t *testing.T) {
			r := cut(&a)
			if fmt.Sprintf("%+v", r) != tt.want {
				t.Errorf("Wanted '%+v', have got other value '%+v'", tt.want, r)
			}
		})
	}
}

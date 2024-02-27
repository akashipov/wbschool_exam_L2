package main

import (
	"testing"
)

func TestSort(t *testing.T) {
	type args struct {
		filename  string
		columnN   int
		reverse   bool
		isNumberC bool
		isUnique  bool
	}
	a := ReadFlags()

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "base01",
			args: args{
				filename: "./resources/test01.txt",
				columnN:  1,
				reverse:  false,
			},
			want:    "a b c\nb r c\nl o p\nt d g",
			wantErr: false,
		},
		{
			name: "base02",
			args: args{
				filename: "./resources/test01.txt",
				columnN:  2,
				reverse:  false,
			},
			want:    "a b c\nt d g\nl o p\nb r c",
			wantErr: false,
		},
		{
			name: "base03",
			args: args{
				filename: "./resources/test01.txt",
				columnN:  2,
				reverse:  true,
			},
			want:    "b r c\nl o p\nt d g\na b c",
			wantErr: false,
		},
		{
			name: "base04",
			args: args{
				filename:  "./resources/test02.txt",
				columnN:   3,
				reverse:   false,
				isNumberC: false,
			},
			want:    "l o 1\nt d 17\nb r 2\na b 21",
			wantErr: false,
		},
		{
			name: "base05",
			args: args{
				filename:  "./resources/test02.txt",
				columnN:   3,
				reverse:   false,
				isNumberC: true,
			},
			want:    "l o 1\nb r 2\nt d 17\na b 21",
			wantErr: false,
		},
		{
			name: "base06_test_unique",
			args: args{
				filename:  "./resources/test03.txt",
				columnN:   3,
				reverse:   false,
				isNumberC: true,
				isUnique:  true,
			},
			want:    "l o 1\nb r 2\nt d 17\nt a 17\na b 21",
			wantErr: false,
		},
		{
			name: "base07_test_unique",
			args: args{
				filename:  "./resources/test03.txt",
				columnN:   3,
				reverse:   false,
				isNumberC: true,
				isUnique:  false,
			},
			want:    "l o 1\nb r 2\nt d 17\nt d 17\nt a 17\na b 21",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a.filename = tt.args.filename
			a.k = tt.args.columnN
			a.r = tt.args.reverse
			a.n = tt.args.isNumberC
			a.u = tt.args.isUnique
			got, err := Sort(a)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sort() error = '%v', wantErr '%v'", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Sort() = '%v', want '%v'", got, tt.want)
			}
		})
	}
}

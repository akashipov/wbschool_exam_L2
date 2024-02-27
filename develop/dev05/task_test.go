package main

import (
	"strings"
	"testing"
)

func Test_grep(t *testing.T) {
	type args struct {
		f     string
		m     string
		r     bool
		n     bool
		a     int
		b     int
		c     int
		i     bool
		fixed bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "base01",
			args: args{
				f: "kra.txt",
				m: ".*Hello.*",
				r: false,
			},
			want: "",
		},
		{
			name: "base02",
			args: args{
				f: "kra.txt",
				m: ".*Hello.*",
				r: true,
			},
			want: "Matched in file: 'resources/abc/kra.txt'\nHello world\n",
		},
		{
			name: "base03",
			args: args{
				f: "*/kra.txt",
				m: ".*Hello.*",
				r: false,
			},
			want: "",
		},
		{
			name: "base04",
			args: args{
				f: "*/*/kra.txt",
				m: ".*Hello.*",
				r: false,
			},
			want: "Matched in file: 'resources/abc/kra.txt'\nHello world\n",
		},
		{
			name: "base05",
			args: args{
				f: "*/*/kra.txt",
				m: ".*Hello.*",
				r: false,
				n: true,
			},
			want: "Matched in file: 'resources/abc/kra.txt'\nNumber of line: 10\nHello world\n",
		},
		{
			name: "base06",
			args: args{
				f: "*/*/kra.txt",
				m: ".*Hello.*",
				r: false,
				c: 1,
			},
			want: "Matched in file: 'resources/abc/kra.txt'\n4\nHello world\n5\n",
		},
		{
			name: "base07",
			args: args{
				f: "*/*/kra.txt",
				m: ".*Hello.*",
				r: false,
				a: 1,
				b: 2,
			},
			want: "Matched in file: 'resources/abc/kra.txt'\n3\n4\nHello world\n5\n",
		},
		{
			name: "base08",
			args: args{
				f: "kra.txt",
				m: ".*Hello.*",
				r: true,
				i: true,
			},
			want: "Matched in file: 'resources/abc/kra.txt'\nhello world\nMatched in file: 'resources/abc/kra.txt'\nheLlo world\nMatched in file: 'resources/abc/kra.txt'\nHello world\n",
		},
		{
			name: "base09",
			args: args{
				f: "kra.txt",
				m: ".*hello world.*",
				r: true,
				i: false,
				b: 2,
			},
			want: "Matched in file: 'resources/abc/kra.txt'\nhello world\n",
		},
		{
			name: "base10",
			args: args{
				f: "kra.txt",
				m: ".*Hello world.*",
				r: true,
				i: false,
				a: 10,
			},
			want: "Matched in file: 'resources/abc/kra.txt'\nHello world\n5\n6\n7\n",
		},
		{
			name: "base11_fixed",
			args: args{
				f:     "kra.txt",
				m:     ".*Hello world.*",
				r:     true,
				i:     false,
				fixed: true,
			},
			want: "",
		},
		{
			name: "base12_fixed",
			args: args{
				f:     "kra.txt",
				m:     "Hello world",
				r:     true,
				i:     false,
				fixed: true,
			},
			want: "Matched in file: 'resources/abc/kra.txt'\nHello world\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := grepArgs{}
			a.sublineToFind = tt.args.m
			a.isRecursive = tt.args.r
			a.filename = tt.args.f
			a.n = tt.args.n
			a.context = tt.args.c
			a.after = tt.args.a
			a.before = tt.args.b
			a.ignoreCase = tt.args.i
			a.fixed = tt.args.fixed
			if a.ignoreCase {
				a.sublineToFind = strings.ToLower(a.sublineToFind)
			}
			r := grep(&a)
			if tt.want != r {
				t.Errorf("Wanted '%s', but got '%s'\n", tt.want, r)
			}
		})
	}
}

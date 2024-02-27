package main

import (
	"reflect"
	"testing"
)

func TestGetAnagramMap(t *testing.T) {
	type args struct {
		words []string
	}
	tests := []struct {
		name string
		args args
		want map[string][]string
	}{
		{
			name: "common_test01",
			args: args{
				words: []string{"фис", "сиф", "вба", "абв", "хчы", "Абв", "с"},
			},
			want: map[string][]string{"вба": {"абв", "вба"}, "фис": {"сиф", "фис"}, "хчы": []string{"хчы"}},
		},
		{
			name: "test02_one_letter_words",
			args: args{
				words: []string{"a", "р", "с"},
			},
			want: map[string][]string{},
		},
		{
			name: "test03_to_lower",
			args: args{
				words: []string{"фис", "Фис", "сИф"},
			},
			want: map[string][]string{"фис": {"сиф", "фис"}},
		},
		{
			name: "test04_order",
			args: args{
				words: []string{"бав", "абв"},
			},
			want: map[string][]string{"бав": {"абв", "бав"}},
		},
		{
			name: "test05_from_example",
			args: args{
				words: []string{"листок", "слиток", "столик"},
			},
			want: map[string][]string{"листок": {"листок", "слиток", "столик"}},
		},
		{
			name: "test05_doubled_letter_only_one",
			args: args{
				words: []string{"аа", "аа", "аи"},
			},
			want: map[string][]string{"аи": {"аи"}},
		},
		{
			name: "test05_doubled_letter_not_one",
			args: args{
				words: []string{"ааи", "аиа", "аи"},
			},
			want: map[string][]string{"аи": {"аи"}, "ааи": {"ааи", "аиа"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getAnagramMap(tt.args.words); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAnagramMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

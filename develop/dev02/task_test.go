package main

import (
	"strings"
	"testing"
)

func TestExpander_ExpandStr(t *testing.T) {
	type fields struct {
		state        State
		Idx          int
		Str          []rune
		RBuilder     *strings.Builder
		LastToken    *rune
		RepeatNumber []rune
	}
	b := strings.Builder{}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr string
	}{
		{
			"ok01",
			fields{
				state:        Basic,
				Idx:          0,
				Str:          []rune("a3t6"),
				RBuilder:     &b,
				LastToken:    nil,
				RepeatNumber: nil,
			},
			"aaatttttt",
			"",
		},
		{
			"ok02",
			fields{
				state:        Basic,
				Idx:          0,
				Str:          []rune("a3t6\\\\"),
				RBuilder:     &b,
				LastToken:    nil,
				RepeatNumber: nil,
			},
			"aaatttttt\\",
			"",
		},
		{
			"ok03",
			fields{
				state:        Basic,
				Idx:          0,
				Str:          []rune("a3t6\\74"),
				RBuilder:     &b,
				LastToken:    nil,
				RepeatNumber: nil,
			},
			"aaatttttt7777",
			"",
		},
		{
			"ok04",
			fields{
				state:        Basic,
				Idx:          0,
				Str:          []rune("a10"),
				RBuilder:     &b,
				LastToken:    nil,
				RepeatNumber: nil,
			},
			"aaaaaaaaaa",
			"",
		},
		{
			"ok05",
			fields{
				state:        Basic,
				Idx:          0,
				Str:          []rune(""),
				RBuilder:     &b,
				LastToken:    nil,
				RepeatNumber: nil,
			},
			"",
			"",
		},
		{
			"ok06",
			fields{
				state:        Basic,
				Idx:          0,
				Str:          []rune("abcd"),
				RBuilder:     &b,
				LastToken:    nil,
				RepeatNumber: nil,
			},
			"abcd",
			"",
		},
		{
			"ok07",
			fields{
				state:        Basic,
				Idx:          0,
				Str:          []rune("qwe\\5"),
				RBuilder:     &b,
				LastToken:    nil,
				RepeatNumber: nil,
			},
			"qwe5",
			"",
		},
		{
			"ok08",
			fields{
				state:        Basic,
				Idx:          0,
				Str:          []rune("qwe\\4\\5"),
				RBuilder:     &b,
				LastToken:    nil,
				RepeatNumber: nil,
			},
			"qwe45",
			"",
		},
		{
			"wrong01",
			fields{
				state:        Basic,
				Idx:          0,
				Str:          []rune("qwe\\"),
				RBuilder:     &b,
				LastToken:    nil,
				RepeatNumber: nil,
			},
			"",
			"escape letter at the end of line",
		},
		{
			"wrong02",
			fields{
				state:        Basic,
				Idx:          0,
				Str:          []rune("qwe\\a"),
				RBuilder:     &b,
				LastToken:    nil,
				RepeatNumber: nil,
			},
			"",
			"string value is wrong. After '\\' comes wrong letter",
		},
		{
			"wrong03",
			fields{
				state:        Basic,
				Idx:          0,
				Str:          []rune("45"),
				RBuilder:     &b,
				LastToken:    nil,
				RepeatNumber: nil,
			},
			"",
			"adding '45' times 'nil' token, passed string is wrong",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Expander{
				state:        tt.fields.state,
				Idx:          tt.fields.Idx,
				Str:          tt.fields.Str,
				RBuilder:     tt.fields.RBuilder,
				LastToken:    tt.fields.LastToken,
				RepeatNumber: tt.fields.RepeatNumber,
			}
			got, err := e.ExpandStr()
			if err != nil && tt.wantErr != err.Error() {
				t.Errorf("Error value is = '%v', wantedErr '%v'", err.Error(), tt.wantErr)
			} else {
				if got != tt.want {
					t.Errorf("Expander.ExpandStr() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

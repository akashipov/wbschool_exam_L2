package pattern

import (
	"fmt"
	"strings"
	"unicode"
)

/*
	Реализовать паттерн «состояние».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/State_pattern
*/

type State int

const (
	Basic State = iota
	Token
	Numeric
)

type Scanner struct {
	S        State
	Str      []rune
	Idx      int
	Tokens   []string
	Numerics []string
}

func (s *Scanner) TokenScan(b *strings.Builder) {
	var v rune
	if s.Idx < len(s.Str) {
		v = s.Str[s.Idx]
	} else {
		v = ' '
	}
	if unicode.IsDigit(v) || unicode.IsLetter(v) {
		s.Idx++
		b.WriteRune(v)
	} else {
		s.Tokens = append(s.Tokens, b.String())
		b.Reset()
		s.S = Basic
	}
}

func (s *Scanner) NumericScan(b *strings.Builder) {
	var v rune
	if s.Idx < len(s.Str) {
		v = s.Str[s.Idx]
	} else {
		v = ' '
	}
	if unicode.IsDigit(v) {
		b.WriteRune(v)
		s.Idx++
	} else if unicode.IsLetter(v) {
		s.S = Token
	} else {
		s.Numerics = append(s.Numerics, b.String())
		b.Reset()
		s.S = Basic
	}
}

func (s *Scanner) BasicScan(b *strings.Builder) {
	var v rune
	if s.Idx < len(s.Str) {
		v = s.Str[s.Idx]
	} else {
		v = ' '
	}
	if unicode.IsLetter(v) {
		s.S = Token
	} else if unicode.IsDigit(v) {
		s.S = Numeric
	} else {
		s.Idx++
	}
}

func (s *Scanner) Scan() {
	b := strings.Builder{}
	for s.Idx < len(s.Str) {
		switch s.S {
		case Basic:
			s.BasicScan(&b)
		case Token:
			s.TokenScan(&b)
		case Numeric:
			s.NumericScan(&b)
		}
	}
	switch s.S {
	case Token:
		s.TokenScan(&b)
	case Numeric:
		s.NumericScan(&b)
	}
}

func (s *Scanner) Print() {
	fmt.Println(s.Tokens)
	fmt.Println(s.Numerics)
}

// func main() {
// 	s := Scanner{
// 		S:        Basic,
// 		Tokens:   make([]string, 0),
// 		Numerics: make([]string, 0),
// 		Idx:      0,
// 	}
// 	s.Str = []rune("aspifjhrigh3i 6 askjdon askdl 123")
// 	s.Scan()
// 	s.Print()
// }

// +
// декомпозирует код в зависимости от количества состояний
// концентрирует код с определенным состоянием в одном месте, очень удобно и комфортно для понимания
// -
// может неоправдано усложнить код для понимания, если применения паттерна не уместно.
//
// пример:
// может быть использовано например при написания статического анализатора, который разбивает на токены

package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

/*
=== Задача на распаковку ===

Создать Go функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы / руны, например:
	- "a4bc2d5e" => "aaaabccddddde"
	- "abcd" => "abcd"
	- "45" => "" (некорректная строка)
	- "" => ""
Дополнительное задание: поддержка escape - последовательностей
	- qwe\4\5 => qwe45 (*)
	- qwe\45 => qwe44444 (*)
	- qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка функция должна возвращать ошибку. Написать unit-тесты.

Функция должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

// State - list of all states for Expander
type State int

const (
	// Basic triggers with BasicLogic
	Basic State = iota
	// Escape triggers with EscapeLogic
	Escape
	// Numeric triggers with NumericLogic
	Numeric
)

// Expander - description above, decodes string to wider form
type Expander struct {
	state        State
	Idx          int
	Str          []rune
	RBuilder     *strings.Builder
	LastToken    *rune
	RepeatNumber []rune
}

// Add adds last letter with repeating n times in result
func (e *Expander) Add(n int) error {
	if e.LastToken == nil {
		return fmt.Errorf("adding '%d' times 'nil' token, passed string is wrong", n)
	}
	e.RBuilder.WriteString(strings.Repeat(string(*e.LastToken), n))
	e.LastToken = nil
	return nil
}

// EscapeLogic works with state Escape, there is all about '\' sign
func (e *Expander) EscapeLogic() error {
	if e.Idx >= len(e.Str) {
		e.Idx++
		return errors.New("escape letter at the end of line")
	}
	v := e.Str[e.Idx]
	if unicode.IsDigit(v) || v == '\\' {
		e.LastToken = &v
		e.state = Basic
	} else {
		return errors.New("string value is wrong. After '\\' comes wrong letter")
	}
	e.Idx++
	return nil
}

// NumericLogic works with state Numeric, there is all about how much need to repeat last letter
func (e *Expander) NumericLogic() error {
	if e.Idx >= len(e.Str) {
		n, err := strconv.Atoi(string(e.RepeatNumber))
		e.Idx++
		if err != nil {
			return err
		}
		if err = e.Add(n); err != nil {
			return err
		}
		e.RepeatNumber = nil
		e.state = Basic
		return nil
	}
	v := e.Str[e.Idx]
	if unicode.IsDigit(v) {
		if e.RepeatNumber == nil {
			e.RepeatNumber = []rune{v}
		} else {
			e.RepeatNumber = append(e.RepeatNumber, v)
		}
		e.Idx++
		return nil
	}
	n, err := strconv.Atoi(string(e.RepeatNumber))
	if err != nil {
		return err
	}

	if err = e.Add(n); err != nil {
		return err
	}
	e.RepeatNumber = nil
	e.state = Basic
	return nil
}

// BaseLogic works with basical situations, sets last token and sents to other state Numeric or Escape
func (e *Expander) BaseLogic() error {
	if e.Idx >= len(e.Str) {
		e.Idx++
		if e.LastToken != nil {

			if err := e.Add(1); err != nil {
				return err
			}
		}
		return nil
	}
	v := e.Str[e.Idx]
	if unicode.IsDigit(v) {
		e.state = Numeric
		return nil
	}
	if e.LastToken != nil {
		if err := e.Add(1); err != nil {
			return err
		}
	}
	if v != '\\' {
		e.LastToken = &v
	} else {
		e.state = Escape
	}
	e.Idx++
	return nil
}

// ExpandStr runs decoding logic for string Expander.Str
func (e *Expander) ExpandStr() (string, error) {
	defer func() {
		e.RBuilder.Reset()
		e.Idx = 0
		e.RepeatNumber = nil
		e.LastToken = nil
	}()
	var err error
	for e.Idx <= len(e.Str) {
		switch e.state {
		case Basic:
			err = e.BaseLogic()
		case Escape:
			err = e.EscapeLogic()
		case Numeric:
			err = e.NumericLogic()
		}
		if err != nil {
			return "", err
		}
	}
	return e.RBuilder.String(), nil
}

// func main() {
// 	// Прошу прощения что без конструктора, я подумал здесь упор на паттерн проектирования
// 	// или на само решение
// 	b := strings.Builder{}
// 	exp := Expander{
// 		state:     Basic,
// 		Idx:       0,
// 		Str:       []rune("a3t7"),
// 		RBuilder:  &b,
// 		LastToken: nil,
// 	}
// 	fmt.Println(exp.ExpandStr())
// 	exp.Str = []rune("a5t\\73")
// 	fmt.Println(exp.ExpandStr())
// 	exp.Str = []rune("a5t10")
// 	fmt.Println(exp.ExpandStr())
// 	exp.Str = []rune("a5t10\\\\")
// 	fmt.Println(exp.ExpandStr())
// 	exp.Str = []rune("")
// 	fmt.Println(exp.ExpandStr())
// 	exp.Str = []rune("a5t10\\a")
// 	fmt.Println(exp.ExpandStr())
// 	exp.Str = []rune("45")
// 	fmt.Println(exp.ExpandStr())
// 	exp.Str = []rune("a5t10\\")
// 	fmt.Println(exp.ExpandStr())
// }

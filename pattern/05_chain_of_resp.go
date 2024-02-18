package pattern

import (
	"fmt"
)

/*
	Реализовать паттерн «цепочка вызовов».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Chain-of-responsibility_pattern
*/

type Request struct {
	Mode Flag
	Data []byte
}

type Flag int8

const (
	R Flag = iota
	W
	S
	SW
)

type Action interface {
	Process(r Request)
	SetNext(a Action)
}

type Logger struct {
	N Action
}

func (l *Logger) Process(r Request) {
	if r.Mode&W != 0 {
		fmt.Println("Printing to log file:", r.Data)
	}
	if l.N != nil {
		l.N.Process(r)
	}
}

func (l *Logger) SetNext(a Action) {
	l.N = a
}

type Gzip struct {
	N Action
}

func (l *Gzip) Process(r Request) {
	fmt.Println("Encoding data... %s -> 2045811285852", string(r.Data))
	r.Data = []byte("2045811285852")
	if l.N != nil {
		l.N.Process(r)
	}
}

func (l *Gzip) SetNext(a Action) {
	l.N = a
}

type Saver struct {
	N Action
}

func (l *Saver) Process(r Request) {
	if r.Mode&S != 0 {
		fmt.Println("Saving data... %s", string(r.Data))
	}
	if l.N != nil {
		l.N.Process(r)
	}
}

func (l *Saver) SetNext(a Action) {
	l.N = a
}

// func main() {
// 	l := Logger{}
// 	g := Gzip{}
// 	s := Saver{}
// 	g.SetNext(&s)
// 	l.SetNext(&g)
// 	l.Process(Request{Data: []byte("blabla"), Mode: R})
// 	fmt.Println()
// 	l.Process(Request{Data: []byte("blabla"), Mode: SW})
// 	fmt.Println()
// 	l.Process(Request{Data: []byte("blabla"), Mode: S})
// 	fmt.Println()
// 	l.Process(Request{Data: []byte("blabla"), Mode: W})
// }

// используется например в http серверах как мидлвейр так называется, добавляются логи, проверка на авторизацию и так далее
// +
// обработчики не влияют друг на друга - их легко писать, отлаживать и переносить между проектами
// каждый обработчик занимается конкретным делом
//
// -
// усложение программы за счет дополнительных классов
// клиенту нужно разобраться в сложности стратегий, чтобы выбрать подходящую для конкретной задачи

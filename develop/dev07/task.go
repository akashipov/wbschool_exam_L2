package main

import (
	"fmt"
	"time"
)

/*
=== Or channel ===

Реализовать функцию, которая будет объединять один или более done каналов в single канал если один из его составляющих каналов закроется.
Одним из вариантов было бы очевидно написать выражение при помощи select, которое бы реализовывало эту связь,
однако иногда неизестно общее число done каналов, с которыми вы работаете в рантайме.
В этом случае удобнее использовать вызов единственной функции, которая, приняв на вход один или более or каналов, реализовывала весь функционал.

Определение функции:
var or func(channels ...<- chan interface{}) <- chan interface{}

Пример использования функции:
sig := func(after time.Duration) <- chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
}()
return c
}

start := time.Now()
<-or (
	sig(2*time.Hour),
	sig(5*time.Minute),
	sig(1*time.Second),
	sig(1*time.Hour),
	sig(1*time.Minute),
)

fmt.Printf(“fone after %v”, time.Since(start))
*/

// Or - function from description to resolve the problem
func Or(channels ...<-chan interface{}) <-chan struct{ channelNumber int } {
	done := make(chan struct{ channelNumber int })
	for idx, channel := range channels {
		go func(c <-chan interface{}) {
			<-c
			done <- struct{ channelNumber int }{channelNumber: idx}
		}(channel)
	}
	return done
}

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}
	start := time.Now()
	done := Or(
		sig(2*time.Hour),
		sig(500*time.Millisecond),
		sig(5*time.Minute),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	r := <-done
	fmt.Println(r.channelNumber)

	fmt.Printf("Done after %v", time.Since(start))
}

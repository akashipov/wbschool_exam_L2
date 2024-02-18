package pattern

import "fmt"

/*
	Реализовать паттерн «стратегия».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Strategy_pattern
*/

type DeleteAlg interface {
	Delete(c *cache)
}

type cache struct {
	e           DeleteAlg
	capacity    int
	maxCapacity int
	elements    map[string]string
}

func (c *cache) SetMode(a DeleteAlg) {
	c.e = a
}

func (c *cache) Add(key string, value string) {
	if c.capacity > c.maxCapacity {
		c.Delete()
	}
	c.capacity++
	c.elements[key] = value
}

func (c *cache) Delete() {
	c.e.Delete(c)
	c.capacity--
}

type Fifo struct {
}

func (f Fifo) Delete(c *cache) {
	fmt.Println("Fifo algorithm")
}

type Lru struct {
}

func (f Lru) Delete(c *cache) {
	fmt.Println("Lru algorithm")
}

// func main() {
// 	l := Lru{}
// 	f := Fifo{}
// 	c := cache{
// 		e:           l,
// 		maxCapacity: 1,
// 		capacity:    0,
// 		elements:    make(map[string]string),
// 	}
// 	c.Add("5", "6")
// 	c.Add("5", "6")
// 	c.Add("10", "6")
// 	c.SetMode(f)
// 	c.Add("6", "7")
// }

// паттерн позволяет очень легко поменять поведение программы в зависимости от предпочтений, не меняя всего остального
// (только схожие алгоритмы могут быть использованы, то есть которые выполняют примерно одну и ту же задачу только разным образом)
// -
// усложнение программы, возникают дополнительные классы
// пользователь должен понимать в чем разница алгоритмов

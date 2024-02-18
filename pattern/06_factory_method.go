package pattern

import (
	"errors"
)

/*
	Реализовать паттерн «фабричный метод».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Factory_method_pattern
*/

type DeliverTransport interface {
	setName(n string)
	GetName() string
	setSpeed(s uint)
	GetSpeed() uint
}

type transport struct {
	name  string
	speed uint
}

func (t *transport) setName(n string) {
	t.name = n
}

func (t *transport) GetName() string {
	return t.name
}

func (t *transport) setSpeed(s uint) {
	t.speed = s
}

func (t *transport) GetSpeed() uint {
	return t.speed
}

type car struct {
	transport
}

func newCar() DeliverTransport {
	return &car{transport: transport{
		name:  "Car",
		speed: 10,
	}}
}

type byfoot struct {
	transport
}

func newByFoot() DeliverTransport {
	return &byfoot{transport: transport{
		name:  "ByFoot",
		speed: 1,
	}}
}

func getDeliveryTransport(n string) (DeliverTransport, error) {
	switch n {
	case "car":
		return newCar(), nil
	case "byfoot":
		return newByFoot(), nil
	default:
		return nil, errors.New("Wrong type of deliver transport")
	}
}

// func main() {
// 	a, _ := getDeliveryTransport("byfoot")
// 	fmt.Println(a)
// }

// Один удобный конструктор по входящим данным
// универсальный интерфейс возвращается, то есть некая обстракция, не нужно понимать тонкости реализации всего остального
// легкое добавление новых транспортов по аналогии в будущем
// -
// работа с интерфейсом имеет свои минусы
// создании фабрики даже для одного объекта на начальном этапе
//

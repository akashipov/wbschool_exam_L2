package pattern

/*
	Реализовать паттерн «посетитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Visitor_pattern
*/

/*
Добавление функционала более сложному объекту, которые нужно выполнить независимо на его части
например: тестирование
*/
import (
	"fmt"
)

// CarPart — семейство типов, которым хотим добавить
// функциональность детали автомобиля.
type CarPart interface {
	Accept(CarPartVisitor)
}

// CarPartVisitor — интерфейс visitor,
// в его коде и содержится новая функциональность.
type CarPartVisitor interface {
	testWheel(wheel *Wheel)
	testEngine(engine *Engine)
}

// Wheel — реализация деталей.
type Wheel struct {
	Name string
}

// Accept — единственный метод, который нужно добавить типам семейства,
// ссылка на метод visitor.
func (w *Wheel) Accept(visitor CarPartVisitor) {
	visitor.testWheel(w)
}

type Engine struct{}

func (e *Engine) Accept(visitor CarPartVisitor) {
	visitor.testEngine(e)
}

type Car struct {
	parts []CarPart
}

// NewCar — конструктор автомобиля.
func NewCar() *Car {
	this := new(Car)
	this.parts = []CarPart{
		&Wheel{"front left"},
		&Wheel{"front right"},
		&Wheel{"rear right"},
		&Wheel{"rear left"},
		&Engine{}}
	return this
}

func (c *Car) Accept(visitor CarPartVisitor) {
	for _, part := range c.parts {
		part.Accept(visitor)
	}
}

// TestVisitor — конкретная реализация visitor,
// которая может проверять колёса и двигатель.
type TestVisitor struct {
}

func (v *TestVisitor) testWheel(wheel *Wheel) {
	fmt.Printf("Testing the %v wheel\n", wheel.Name)
}

func (v *TestVisitor) testEngine(engine *Engine) {
	fmt.Println("Testing engine")
}

// func main() {
//     // клиентский код
//     car := NewCar()
//     visitor := new(TestVisitor)
//     car.Accept(visitor)
// }

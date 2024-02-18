package pattern

import "fmt"

/*
	Реализовать паттерн «фасад».
Объяснить применимость паттерна, его плюсы и минусы,а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Facade_pattern
*/

type Site struct {
}

func (s *Site) GetOrder() {
	fmt.Println("Creating order...")
	fmt.Println("Order has been created")
}

func (s *Site) CloseOrder() {
	fmt.Println("Closing order...")
	fmt.Println("Order has been closed")
}

type DeliveryMan struct {
	Name     string
	Lastname string
	IsBusy   bool
}

func (d *DeliveryMan) Deliver() {
	fmt.Println("Delivering order...")
	fmt.Println("Order has been delivered")
}

func (d *DeliveryMan) GetPizza() {
	fmt.Println("Order has been got")
}

func (d *DeliveryMan) GetPayment() {
	fmt.Println("Order was paid")
}

type Kitchen struct {
}

func (p *Kitchen) CookPizza() {
	fmt.Println("Pizza has been cooked")
}

type Pizza struct {
	site       Site
	deliverman DeliveryMan
	kitchen    Kitchen
}

func (p *Pizza) StartOrder() {
	p.site.GetOrder()
	p.kitchen.CookPizza()
	p.deliverman.GetPizza()
	p.deliverman.Deliver()
}

func (p *Pizza) EndOrder() {
	p.deliverman.GetPayment()
	p.site.CloseOrder()
}

// делается для того чтобы декомпозировать сложный процесс на более простые по логике
// клиенту по сути нужно лишь только создать заказ и оплатить
// происходит некоторое упращение логики, сложная логика(излишняя) скрывается внутри других объектов
// из минусов например: приходится создавать дополнительный объект

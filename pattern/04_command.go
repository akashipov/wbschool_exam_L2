package pattern

import "fmt"

/*
	Реализовать паттерн «комманда».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Command_pattern
*/

type Command interface {
	execute()
}

type Invoker struct {
	Commands map[string]Command
}

func (i *Invoker) Do(c string) {
	i.Commands[c].execute()
}

type Receiver interface {
	Action(value string)
}

// выглядит как некий адаптер
type Printer struct {
	R     Receiver
	Value string
}

func (c *Printer) execute() {
	c.R.Action(c.Value)
}

type RemoteMachine struct {
	Host string
	Port string
}

func (p *RemoteMachine) Action(value string) {
	fmt.Println("remote printing...:", p.Host, p.Port, value)
}

func main() {
	p := RemoteMachine{"123.1.1.123", "8080"}
	pri := Printer{&p, "John Smith"}
	inv := Invoker{}
	inv.Commands = make(map[string]Command)
	inv.Commands["print name"] = &pri
	inv.Do("print name")
}

/*
Например я видел на старой работе такой паттерн использовался для вызова нашего скрипта
на аирфлоу(аирфлоу оператор) в докер контейнере или в кубернетисе
то есть формировался некий класс с полями, которые потом уходили во флаги скрипта
(грубо говоря исполнялся bash скрипт по итогу на удаленной машине)
*/

/*
+
Мы даем удобный интерфейс для работы с командами, инаксулируем какие то вещи от конечного потребителя
-
мы все таки накручиваем и усложняем код в каком то смысле, это не всегда может быть оправдано
*/

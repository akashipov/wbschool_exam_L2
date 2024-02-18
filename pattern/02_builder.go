package pattern

/*
	Реализовать паттерн «строитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Builder_pattern
*/

type Person struct {
	name    string
	surname string
}

func (p *Person) SetName(name string) *Person {
	p.name = name
	return p
}

func (p *Person) SetSurname(surname string) *Person {
	p.surname = surname
	return p
}

func NewPerson() *Person {
	return &Person{}
}

// func main() {
// 	p := NewPerson().SetName("Artyom").SetSurname("Kashipov")
// 	fmt.Println(p)
// }

// мы можем использовать его чтобы декомпозировать логику формирования каждого параметра отдельно
// (так как у каждого параметра может быть очень сложная логика), все вместе смотрелось бы ужасно
// а также мы можем делать частичную инициализацию, то есть есть некая свобода в этом плане и гибкость

// например у нас есть какой то объект который может работать и по сети и локально в оффлайн
// мы можем за счет параметров запуска запускать лишь в одном каком то режиме и не запускать лишнюю возможно очень сложную логику
// в данном случае если нет интернета могут возникать некоторые лишние ожидания, которые займут какое то время

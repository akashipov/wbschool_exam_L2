Что выведет программа? Объяснить вывод программы.

```go
package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}
```

Ответ:
```
ну это похожий кейс как в одной из задач выше 03 по моему. Выведется ошибка то есть error
Так как у нас на выходе будет считаться что мы возвращаем уже что то вроде
var a *customError = nil
то есть произойдет каст на выходе в нулевой значение определенного типа. а раз это уже определенный тип то это при касте в интерфейс выдаст не 0 значение
```

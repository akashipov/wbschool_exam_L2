Что выведет программа? Объяснить вывод программы. Рассказать про внутреннее устройство слайсов и что происходит при передачи их в качестве аргументов функции.

```go
package main

import (
	"fmt"
)

func main() {
	var s = []string{"1", "2", "3"}
	modifySlice(s)
	fmt.Println(s)
}

func modifySlice(i []string) {
	i[0] = "3"
	i = append(i, "4")
	i[1] = "5"
	i = append(i, "6")
}
```

Ответ:
```
будет вывод
[3, 2, 3]

значение i[0] = "3" поменяется и там снаружи так как мы ссылаемся пока еще на тот же массив в этой первой строке

когда мы добавим элемент (append) произойдет переполнение за капасити и создастся новый массив(копия с большим капасити) и мы уже тут потеряем связь с реальностью снаружи
дальнейшие любые изменения не повлияют на вывод

```

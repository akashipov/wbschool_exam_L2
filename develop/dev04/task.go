package main

import (
	"fmt"
	"slices"
	"strings"
)

/*
=== Поиск анаграмм по словарю ===

Напишите функцию поиска всех множеств анаграмм по словарю.
Например:
'пятак', 'пятка' и 'тяпка' - принадлежат одному множеству,
'листок', 'слиток' и 'столик' - другому.

Входные данные для функции: ссылка на массив - каждый элемент которого - слово на русском языке в кодировке utf8.
Выходные данные: Ссылка на мапу множеств анаграмм.
Ключ - первое встретившееся в словаре слово из множества
Значение - ссылка на массив, каждый элемент которого, слово из множества. Массив должен быть отсортирован по возрастанию.
Множества из одного элемента не должны попасть в результат.
Все слова должны быть приведены к нижнему регистру.
В результате каждое слово должно встречаться только один раз.

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func getKeyFromWord(word string) (string, int) {
	m := make(map[rune]rune)
	r := make([]rune, 0)
	for _, v := range word {
		_, ok := m[v]
		if !ok {
			m[v] = v
			r = append(r, v)
		}
	}
	slices.Sort(r)
	a := []rune(word)
	slices.Sort(a)
	return string(a), len(r)
}

func getSetFromArray(a []string) []string {
	m := make(map[string]string)
	r := make([]string, 0)
	for _, v := range a {
		_, ok := m[v]
		if !ok {
			m[v] = v
			r = append(r, v)
		}
	}
	slices.Sort(r)
	return r
}

func getAnagramMap(words []string) map[string][]string {
	m := make(map[string][]string)
	firstWordMapping := make(map[string]string)
	for _, word := range words {
		word = strings.ToLower(word)
		wordKey, countLetters := getKeyFromWord(word)
		if countLetters == 1 {
			continue
		}
		v, ok := firstWordMapping[wordKey]
		if !ok {
			firstWordMapping[wordKey] = word
			m[word] = []string{word}
		} else {
			m[v] = append(m[v], word)
		}
	}
	for key := range m {
		m[key] = getSetFromArray(m[key])
	}
	return m
}

func main() {
	words := []string{"фис", "сиф", "вба", "абв", "хчы", "Абв", "с"}
	fmt.Printf("%+v\n", getAnagramMap(words))
}

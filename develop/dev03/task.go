package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
)

/*
=== Утилита sort ===

Отсортировать строки (man sort)
Основное

Поддержать ключи

-k — указание колонки для сортировки
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

Дополнительное

Поддержать ключи

-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учётом суффиксов

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

// Arguments - possibly argument to save it in program
type Arguments struct {
	filename string
	k        int
	n        bool
	r        bool
	u        bool
	M        bool
	b        bool
	c        bool
	h        bool
}

// ReadFlags - read all arguments for 'sort' from the description
func ReadFlags() Arguments {
	k := flag.Int("k", 1, "Number of column for sorting")
	f := flag.String("f", "", "file's name for scan")
	r := flag.Bool("r", false, "Reverse result")
	n := flag.Bool("n", false, "Compare column 'k' like a number")
	h := flag.Bool("h", false, "Compare column 'k' like a number with suffix")
	c := flag.Bool("c", false, "Check is sorted lines or not?")
	b := flag.Bool("b", false, "Ignore extra space at the end of line")
	M := flag.Bool("M", false, "Work with column like month's name")
	u := flag.Bool("u", false, "Distinct all repeated strings")
	flag.Parse()
	return Arguments{
		k: *k, filename: *f,
		r: *r, n: *n, u: *u,
		M: *M, b: *b, c: *c, h: *h,
	}
}

// AppendToResult - add to result the list of lines with repeating for example
func AppendToResult(m map[string][]string, key string, unique bool) []string {
	r := make([]string, 0, len(m[key]))
	o := make(map[string]string)
	for _, elem := range m[key] {
		_, ok := o[elem]
		if !ok || !unique {
			r = append(r, elem)
			o[elem] = elem
		}
	}
	delete(m, key)
	return r
}

// ReadLines - read content from the file to work with
func ReadLines(isNumberColumn bool, columnNumber int, filename string) ([]string, map[string][]string, []string, []int, error) {
	results := make([]string, 0)
	m := make(map[string][]string)
	toSortListStrs := make([]string, 0)
	toSortListInts := make([]int, 0)
	f, err := os.OpenFile(filename, os.O_RDONLY, 0000)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	err = nil
	idx := 0
	for err != io.EOF {
		var lineBytes []byte
		lineBytes, _, err = r.ReadLine()
		if err != nil && err != io.EOF {
			return nil, nil, nil, nil, err
		}
		if err == io.EOF {
			break
		}
		lineStr := string(lineBytes)
		l := strings.Split(lineStr, " ")
		if columnNumber > len(l) {
			results = append(results, lineStr)
			idx++
			continue
		}
		_, ok := m[l[columnNumber-1]]
		if ok {
			m[l[columnNumber-1]] = append(m[l[columnNumber-1]], lineStr)
		} else {
			m[l[columnNumber-1]] = []string{lineStr}
		}
		if isNumberColumn {
			n, err := strconv.Atoi(l[columnNumber-1])
			if err == nil {
				toSortListInts = append(toSortListInts, n)
				idx++
				continue
			}
		}
		toSortListStrs = append(toSortListStrs, l[columnNumber-1])
		idx++
	}
	return results, m, toSortListStrs, toSortListInts, nil
}

// Sort - sorting algorithm for columns
func Sort(args Arguments) (string, error) {
	results, m, toSortListStrs, toSortListInts, err := ReadLines(args.n, args.k, args.filename)
	if err != nil {
		return "", err
	}
	slices.Sort(toSortListInts)
	slices.Sort(toSortListStrs)
	for _, v := range toSortListStrs {
		results = append(results, AppendToResult(m, v, args.u)...)
	}
	for _, v := range toSortListInts {
		results = append(results, AppendToResult(m, strconv.Itoa(v), args.u)...)
	}
	if args.r {
		slices.Reverse(results)
	}
	return strings.Join(results, "\n"), nil
}

func main() {
	args := ReadFlags()
	r, err := Sort(args)
	if err != nil {
		fmt.Println("Error has got:", err.Error())
	}
	fmt.Printf("Result is: '%s'\n", r)
}

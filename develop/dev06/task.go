package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

/*
=== Утилита cut ===

Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type cutArgs struct {
	file string
	f    string
	d    string
	s    bool
}

const (
	mincolumnidx = 1
	maxcolumnidx = 2147483647
)

func getCutArgs() *cutArgs {
	file := flag.String("file", "", "Path to file")
	f := flag.String("f", "1-", "Columns to select")
	d := flag.String("d", " ", "Delimeter for lines")
	s := flag.Bool("s", false, "Select lines only with separator")
	flag.Parse()
	if file == nil {
		log.Fatalln("There is no file path for input")
	}
	return &cutArgs{file: *file, f: *f, d: *d, s: *s}
}

type bound struct {
	start int
	end   int
}

type boundes struct {
	intervals []bound
}

func (b *boundes) getLength() int {
	l := 0
	for _, i := range b.intervals {
		l += i.end - i.start + 1
	}
	return l
}

func (b *boundes) update(newB bound) {
	first := -1
	last := -1
	idx := 0
	if len(b.intervals) == 0 {
		b.intervals = append(b.intervals, newB)
		return
	}
	for idx < len(b.intervals) {
		v := b.intervals[idx]
		if first == -1 && newB.start >= v.start && newB.start <= v.end {
			first = idx
		}
		if last == -1 && newB.end >= v.start && newB.end <= v.end {
			last = idx
		}
		if v.start > newB.end {
			break
		}
		idx++
	}
	if first != -1 && last == -1 {
		b.intervals[first].end = newB.end
	} else if first == -1 && last != -1 {
		b.intervals[last].start = newB.start
	} else if first == -1 && last == -1 {
		if newB.end < b.intervals[0].start {
			b.intervals = append([]bound{newB}, b.intervals...)
		} else if newB.start > b.intervals[len(b.intervals)-1].end {
			b.intervals = append(b.intervals, newB)
		}
	} else if first != -1 && last != -1 && first != last {
		b.intervals[last].start = b.intervals[first].start
		a := make([]bound, 0, first+len(b.intervals)-last)
		a = append(a, b.intervals[:first]...)
		a = append(a, b.intervals[last:]...)
		b.intervals = a
	}
}

func parseColumnsNumber(fields string) boundes {
	intervals := strings.Split(fields, ",")
	o := boundes{make([]bound, 0)}
	for _, interval := range intervals {
		ranges := strings.Split(interval, "-")
		var l, r int
		var err error
		if len(ranges) == 2 {
			if ranges[0] == "" {
				l = mincolumnidx
			} else {
				l, err = strconv.Atoi(ranges[0])
				if err != nil {
					log.Fatalln(err.Error())
				}
			}
			if ranges[1] == "" {
				r = maxcolumnidx
			} else {
				r, err = strconv.Atoi(ranges[1])
				if err != nil {
					log.Fatalln(err.Error())
				}
			}

		} else if len(ranges) == 1 {
			if ranges[0] == "" {
				continue
			}
			r, err = strconv.Atoi(ranges[0])
			if err != nil {
				log.Fatalln(err.Error())
			}
			l = r
		}
		if l < mincolumnidx || r > maxcolumnidx {
			log.Fatalf("Bounderies must be into the interval [%d;%d]\n", mincolumnidx, maxcolumnidx)
		}
		o.update(bound{l, r})
	}
	return o
}

func cut(args *cutArgs) [][]string {
	b := parseColumnsNumber(args.f)
	f, err := os.OpenFile(args.file, os.O_RDONLY, 0)
	if err != nil {
		log.Fatalln(err.Error())
	}
	text, err := io.ReadAll(f)
	if err != nil {
		log.Fatalln(err.Error())
	}
	lines := strings.Split(string(text), "\n")
	result := make([][]string, 0, len(lines))
	lengthPerLine := b.getLength()
	for _, line := range lines {
		columns := strings.Split(line, args.d)
		if len(columns) == 1 && args.s {
			continue
		}
		r := make([]string, 0, lengthPerLine)
		for _, v := range b.intervals {
			e := v.end
			if v.start > len(columns) {
				break
			}
			if len(columns) < v.end {
				e = len(columns)
			}
			r = append(r, columns[v.start-1:e]...)
			if e != v.end {
				break
			}
		}
		result = append(result, r)
	}
	return result
}

func main() {
	args := getCutArgs()
	fmt.Println(cut(args))
}

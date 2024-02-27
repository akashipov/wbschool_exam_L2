package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

/*
=== Утилита grep ===

Реализовать утилиту фильтрации (man grep)

Поддержать флаги:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", печатать номер строки

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type grepArgs struct {
	filename      string
	sublineToFind string
	isRecursive   bool
	isCounted     bool
	before        int
	after         int
	context       int
	n             bool
	ignoreCase    bool
	fixed         bool
	invert        bool
}

func parseArgs() *grepArgs {
	f := flag.String("f", "", "Path to file")
	r := flag.Bool("r", false, "Match all files recursively")
	c := flag.Bool("c", false, "Count matches or not")
	n := flag.Bool("n", false, "Print num of line or not")
	i := flag.Bool("i", false, "Ignore case or not")
	v := flag.Bool("v", false, "Invert result or not")
	F := flag.Bool("F", false, "It is absolutly the same line like in 'M' flag or not")
	a := flag.Int("A", 0, "\"after\" печатать +N строк после совпадения")
	b := flag.Int("B", 0, "\"before\" печатать +N строк до совпадения")
	C := flag.Int("C", 0, "\"context\" (A+B) печатать ±N строк вокруг совпадения")
	m := flag.String("m", "", "Строка, которую мы хотим найти")
	flag.Parse()
	if *m == "" {
		log.Fatalln("Need to set 'm' flag to detect which string is needed to find")
	}
	if f == nil {
		log.Fatalln("Need to path argument 'f'(filepath)")
	}
	if *C != 0 && (*b != 0 || *a != 0) {
		log.Fatal("Flag of context(C) cannot be set with After or Before(A,B)")
	}
	if *i {
		*m = strings.ToLower(*m)
	}
	return &grepArgs{
		filename: *f, isRecursive: *r, isCounted: *c,
		after: *a, before: *b, context: *C, sublineToFind: *m,
		n: *n, ignoreCase: *i, fixed: *F, invert: *v,
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func findLine(text string, path string, args *grepArgs, b *strings.Builder) error {
	lines := strings.Split(text, "\n")
	c := 0
	var gErr error
	excludeStatus := make([]bool, len(lines))
	for idx, line := range lines {
		workingLine := line
		if args.ignoreCase {
			workingLine = strings.ToLower(line)
		}
		var m bool
		var err error
		if args.fixed {
			m, err = line == args.sublineToFind, nil
		} else {
			m, err = regexp.Match(args.sublineToFind, []byte(workingLine))
		}
		if err != nil {
			gErr = errors.Join(gErr, err)
		}
		if m {
			c++
			b.WriteString(fmt.Sprintf("Matched in file: '%s'\n", path))
			if args.n {
				b.WriteString(fmt.Sprintf("Number of line: %d\n", idx+1))
			}
			for i := max(idx-max(args.context, args.before), 0); i <= min(idx+max(args.context, args.after), len(lines)-1); i++ {
				if args.invert {
					excludeStatus[i] = true
				} else {
					b.WriteString(lines[i] + "\n")
				}
			}
		}
	}
	if args.invert {
		for idx := range lines {
			if !excludeStatus[idx] {
				b.WriteString(lines[idx] + "\n")
			}
		}
	}
	if args.isCounted {
		b.WriteString(fmt.Sprintf("Count of found lines in file '%s': '%d'\n", path, c))
	}
	return gErr
}

func grep(args *grepArgs) string {
	// все это можно распараллелить конечно же
	builder := strings.Builder{}
	if err := filepath.Walk("./", func(path string, info fs.FileInfo, err error) error {
		var matchedName string
		if args.isRecursive {
			matchedName = info.Name()
		} else {
			matchedName = path
		}
		b, err := filepath.Match(args.filename, matchedName)
		if err != nil {
			return err
		}
		if !b {
			// builder.WriteString(fmt.Sprintf("File isn't satisfied the pattern '%s'\n", path))
			return nil
		}
		if info.IsDir() {
			builder.WriteString(fmt.Sprintf("It is directory '%s'\n", path))
			return nil
		}
		if file, err := os.OpenFile(path, os.O_RDONLY, 0000); err == nil {
			defer file.Close()
			b, err := io.ReadAll(file)
			if err != nil {
				return err
			}
			return findLine(string(b), path, args, &builder)
		}
		return errors.New("file cannot be opened")
	}); err != nil {
		log.Fatalln(err.Error())
	}
	return builder.String()
}

func main() {
	args := parseArgs()
	r := grep(args)
	fmt.Print(r)
}

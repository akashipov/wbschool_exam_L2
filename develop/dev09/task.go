package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/chigopher/pathlib"
	"github.com/go-resty/resty/v2"
)

/*
=== Утилита wget ===

Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/
// https://www.google.com

type wgetArgs struct {
	recursive bool
	directory *pathlib.Path
	maxDepth  int
}

// Wget - simple analog of 'wget' from the linux command
type Wget struct {
	args           wgetArgs
	siteRegex      *regexp.Regexp
	processedSites sync.Map
}

// NewWget - constructor for 'wget'
func NewWget() (*Wget, error) {
	reg, err := regexp.Compile("(\"|')(http|ftp)(s)?://[^\"']+(\"|')")
	if err != nil {
		return nil, err
	}
	return &Wget{siteRegex: reg, processedSites: sync.Map{}}, nil
}

func (w *Wget) execute(
	site string, dir *pathlib.Path, depth int,
	waitGP *sync.WaitGroup, filename string,
) error {
	defer waitGP.Done()
	if _, ok := w.processedSites.Load(site); ok {
		fmt.Printf("site %s has been already processed\n", site)
		return nil
	}
	if depth >= w.args.maxDepth {
		return nil
	}
	if dir == nil {
		dir = w.args.directory
	}
	fmt.Printf("execute %s. depth: %d, dir: %s\n", site, depth, dir.String())
	c := resty.New()
	response, err := c.R().Get(site)
	if err != nil {
		return fmt.Errorf("problem with request to site: %w", err)
	}
	if w.args.recursive {
		dir, err = w.createDir(site, dir, depth)
		if err != nil {
			return err
		}
		results := w.siteRegex.FindAll(response.Body(), -1)
		if len(results) == 0 {
			fmt.Printf("Nothing to find recursively. Depth is %d\n", depth)
		}
		waitG := sync.WaitGroup{}
		for idx, v := range results {
			waitG.Add(1)
			go w.execute(
				string(v[1:len(v)-1]), dir, depth+1, &waitG,
				fmt.Sprintf("file_%d.html", idx))
		}
		waitG.Wait()
	}
	pathToFile := filepath.Join(dir.String(), filename)
	err = w.saveHTML(pathToFile, response.Body())
	if err != nil {
		return err
	}
	w.processedSites.Store(site, site)
	return nil
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func (w *Wget) genNameOfFile(depth int) string {
	b := strings.Builder{}
	for i := 0; i < 64; i++ {
		b.WriteRune(letterRunes[rand.Intn(len(letterRunes))])
	}
	return fmt.Sprintf("%s_%d.html", b.String(), depth)
}

func (w *Wget) parseCMDLine() {
	r := flag.Bool("r", false, "Recursive downloading")
	l := flag.Int("l", 2, "Recursive downloading max depth")
	d := flag.String("d", "./", "Directory to save files")
	flag.Parse()
	args := wgetArgs{recursive: *r, directory: pathlib.NewPath(*d), maxDepth: *l}
	w.args = args
}

func (w *Wget) saveHTML(filename string, data []byte) error {
	// fmt.Println(string(data))
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if err != nil {
		return fmt.Errorf("error with saving HTML: %w", err)
	}
	return nil
}

func (w *Wget) createDir(site string, path *pathlib.Path, depth int) (*pathlib.Path, error) {
	var folderName string
	if depth == 0 {
		folderName = strings.Replace(site, "/", "", -1)
	} else {
		folderName = fmt.Sprintf("depth_%d", depth)
	}
	newDir := filepath.Join(path.String(), folderName)
	if _, err := os.Stat(newDir); os.IsNotExist(err) {
		err = os.Mkdir(newDir, 0777)
		if err != nil && os.IsExist(err) {
			return nil, fmt.Errorf("error with creating dir: %w", err)
		}
	}
	return pathlib.NewPath(newDir), nil
}

func (w *Wget) readSiteName() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Please write your site here:")
	site, err := reader.ReadString('\n')
	site = strings.TrimRight(site, "\n")
	if err != nil {
		log.Fatalf(err.Error())
	}
	return site
}

func main() {
	f, err := os.OpenFile("test_4003833119/main.html", os.O_RDONLY, 0)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	textBytes, err := io.ReadAll(f)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(textBytes))
}

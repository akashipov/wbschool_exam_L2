package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/chigopher/pathlib"
	proc "github.com/mitchellh/go-ps"
)

/*
=== Взаимодействие с ОС ===

Необходимо реализовать собственный шелл

встроенные команды: cd/pwd/echo/kill/ps
поддержать fork/exec команды
конвеер на пайпах

Реализовать утилиту netcat (nc) клиент
принимать данные из stdin и отправлять в соединение (tcp/udp)
Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type shell struct {
	PreviousPath *pathlib.Path
	CurrentDir   *pathlib.Path
	stdin        io.Reader
	stdout       io.Writer
	stderr       io.Writer
}

func (s *shell) cd(args []string) {
	// Я мог тут работать с PWD env переменной, но не стал переделывать старый вариант
	// наверное тут хотели чтобы я работал с подгрузкой переменных через os и изменением их
	if !s.CurrentDir.IsAbsolute() {
		r := make(chan string, 1)
		s.pwd(r)
		s.CurrentDir = pathlib.NewPath(<-r)
	}
	if len(args) == 0 {
		fmt.Println("Wrong count of arguments for 'cd' command")
	}
	home := pathlib.NewPath(os.Getenv("HOME"))
	if args[0] == "" {
		s.PreviousPath = s.CurrentDir
		s.CurrentDir = home
	} else if args[0] == "-" {
		if s.PreviousPath != nil {
			s.CurrentDir, s.PreviousPath = s.PreviousPath, s.CurrentDir
		}
	} else if args[0] == ".." {
		s.PreviousPath = pathlib.NewPath(s.CurrentDir.String())
		s.CurrentDir = s.CurrentDir.Parent()
	} else if args[0] == "." {
		s.PreviousPath = s.CurrentDir
	} else if newP := pathlib.NewPath(args[0]); newP.IsAbsolute() {
		s.PreviousPath = s.CurrentDir
		s.CurrentDir = newP
	} else if args[0] == "~/" {
		s.PreviousPath = s.CurrentDir
		s.CurrentDir = home
	} else {
		newP := pathlib.NewPath(filepath.Join(s.CurrentDir.String(), args[0]))
		isDir, err := newP.IsDir()
		if !isDir || err != nil {
			s.errorPrint("it is not directory")
			return
		}
		dirExists, err := newP.DirExists()
		if !dirExists || err != nil {
			s.errorPrint("Directory doesn't exist")
			return
		}
		s.PreviousPath = s.CurrentDir
		s.CurrentDir = newP
	}
}

func (s *shell) pwd(p chan string) {
	f, err := filepath.Abs(s.CurrentDir.String())
	if err != nil {
		s.errorPrint("it is problem extract absolute path for current directory")
		return
	}
	p <- f
}

func (s *shell) echo(args []string, p chan string) {
	p <- args[0]
}

func (s *shell) kill(args []string) {
	if len(args) == 0 {
		fmt.Fprint(os.Stdout, "need to pass id number to kill\n")
		return
	}
	n, err := strconv.Atoi(args[0])
	if err != nil {
		s.errorPrint(err.Error())
		return
	}
	p, err := os.FindProcess(n)
	if err != nil {
		s.errorPrint(err.Error())
		return
	}
	err = p.Kill()
	if err != nil {
		s.errorPrint(err.Error())
		return
	}
	fmt.Fprint(os.Stdout, "process is successfuly killed\n")
}

func (s *shell) ps(args []string, p chan string) {
	// Я на mac os может тут хотели чтобы я читал файл /proc для линукс, но у меня как будто нет аналога пока
	// или мне было сложновато найти аналог (sysctl может, но он какой то непонятный)
	processList, err := proc.Processes()
	if err != nil {
		p <- "ps.Processes() Failed, are you using windows?\n"
		return
	}
	b := strings.Builder{}
	for x := range processList {
		process := processList[x]
		b.WriteString(fmt.Sprintf("%d\t%s\n", process.Pid(), process.Executable()))
	}
	p <- b.String()
}

func (s *shell) getAbsPathToFile(p string) string {
	path := pathlib.NewPath(p)
	p = path.String()
	if !path.IsAbsolute() {
		r := make(chan string, 1)
		s.pwd(r)
		p = filepath.Join(<-r, p)
	}
	return p
}

func (s *shell) createTCP(host string, port int) {
	addr := net.TCPAddr{
		IP:   net.ParseIP(host),
		Port: port,
	}
	srv, err := net.ListenTCP("tcp", &addr)
	if err != nil {
		s.errorPrint(err.Error())
		return
	}
	for {
		conn, err := srv.AcceptTCP()
		if err != nil {
			s.errorPrint(err.Error())
			return
		}
		defer conn.Close()
		for {
			p := make([]byte, 2048)
			_, err = conn.Read(p)
			if err == io.EOF {
				break
			}
			if err != nil {
				s.errorPrint(err.Error())
				return
			}
			fmt.Fprintf(s.stdout, "Have got msg: '%s'\n", string(p))
		}
	}
}

func (s *shell) createUDP(host string, port int) {
	addr := net.UDPAddr{
		IP:   net.ParseIP(host),
		Port: port,
	}
	srv, err := net.ListenUDP("udp", &addr)
	if err != nil {
		s.errorPrint(err.Error())
		return
	}
	for {
		p := make([]byte, 2048)
		_, d, err := srv.ReadFromUDP(p)
		if err != nil {
			s.errorPrint(err.Error())
			return
		}
		fmt.Fprintf(s.stdout, "Have got msg from '%s': '%s'\n", d, string(p))
	}
}

func (s *shell) dialTCP(host string, port int) {
	raddr := &net.TCPAddr{
		IP:   net.ParseIP(host),
		Port: port,
	}
	conn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		s.errorPrint(err.Error())
		return
	}
	line, err := s.readFromStdin()
	if err != nil {
		s.errorPrint(err.Error())
		return
	}
	_, err = conn.Write([]byte(line))
	if err != nil {
		s.errorPrint(err.Error())
		return
	}
}

func (s *shell) dialUDP(host string, port int) {
	raddr := &net.UDPAddr{
		IP:   net.ParseIP(host),
		Port: port,
	}
	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		s.errorPrint(err.Error())
		return
	}
	line, err := s.readFromStdin()
	if err != nil {
		s.errorPrint(err.Error())
		return
	}
	_, err = conn.Write([]byte(line))
	if err != nil {
		s.errorPrint(err.Error())
		return
	}
}

func (s *shell) netcat(args []string) {
	var toCreate, isUDP bool
	var host string = "127.0.0.1"
	var port int
	for idx, a := range args {
		switch a {
		case "-l":
			toCreate = true
		case "-u":
			isUDP = true
		case "-h":
			if toCreate {
				s.errorPrint("it is not allowed to pass host to the listening server")
				return
			}
			// need to check by pattern the value of it
			if idx+1 >= len(args) {
				s.errorPrint("need to pass host value")
				return
			}
			host = args[idx+1]
		case "-p":
			var err error
			if idx+1 >= len(args) {
				s.errorPrint("need to pass port value")
				return
			}
			port, err = strconv.Atoi(args[idx+1])
			if err != nil {
				s.errorPrint("problem with parsing of port value, please check it")
				return
			}
		}
	}
	if toCreate {
		if !isUDP {
			s.createTCP(host, port)
		} else {
			s.createUDP(host, port)
		}
	} else {
		if !isUDP {
			s.dialTCP(host, port)
		} else {
			s.dialUDP(host, port)
		}
	}

}

func (s *shell) errorPrint(msg string) string {
	msg = strings.ToLower(msg)
	msg = strings.TrimFunc(msg, func(r rune) bool {
		if r == '\n' || r == ' ' {
			return true
		}
		return false
	})
	fmt.Fprint(s.stderr, msg+"\n")
	return msg
}

func (s *shell) exec(args []string) {
	if len(args) == 0 {
		fmt.Println("Wrong number of arguments for 'exec' command")
		return
	}
	const (
		inputSign     = "<"
		outputSign    = ">"
		terminalInOut = "/dev/tty"
	)
	if strings.HasPrefix(args[0], inputSign) {
		p := s.getAbsPathToFile(args[0][len(inputSign):])
		if p == terminalInOut {
			s.stdin = os.Stdin
			return
		}
		f, err := os.OpenFile(p, os.O_RDONLY|os.O_CREATE, 0777)
		if err != nil {
			s.errorPrint(fmt.Sprintf("problem with opening file for reading: '%s'", err.Error()))
			return
		}
		s.stdin = f
	} else if strings.HasPrefix(args[0], outputSign) {
		p := s.getAbsPathToFile(args[0][len(outputSign):])
		if p == terminalInOut {
			s.stdout = os.Stdout
			return
		}
		f, err := os.OpenFile(p, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
		if err != nil {
			s.errorPrint(
				fmt.Sprintf("problem with opening file for writing: '%s'", err.Error()),
			)
			return
		}
		s.stdout = f
	} else {
		s.errorPrint(fmt.Sprintf("wrong argument for command 'exec': '%s'", args[0]))
	}
}

func (s *shell) readFromStdin() (string, error) {
	reader := bufio.NewReader(s.stdin)
	var line string
	var err error
	if s.stdin == os.Stdin {
		line, err = reader.ReadString('\n')
	} else {
		var b []byte
		b, err = io.ReadAll(reader)
		if err != nil {
			return "", err
		}
		_, err = s.stdin.(*os.File).Seek(0, io.SeekStart)
		line = string(b)
	}
	if err != nil {
		return "", err
	}
	return line, nil
}

func (s *shell) cat(args []string, p chan string) {
	for _, v := range args {
		if v == "" {
			continue
		} else {
			p <- v + "\n"
			return
		}
	}
	line, err := s.readFromStdin()
	if err != nil {
		s.errorPrint(err.Error())
		return
	}
	p <- line
}

func (s *shell) readCMDLine(f *os.File) ([]string, error) {
	reader := bufio.NewReader(f)
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	line = strings.TrimSuffix(line, "\n")
	line = strings.TrimSpace(line)
	pipes := strings.Split(line, " | ")
	return pipes, nil
}

func (s *shell) cmd(f *os.File, w *sync.WaitGroup) {
	globalForForkW := sync.WaitGroup{}
loop:
	for {
		fmt.Print(s.CurrentDir.String() + "$")
		pipes, err := s.readCMDLine(f)
		if err == io.EOF {
			globalForForkW.Wait()
			break
		}
		if err != nil {
			log.Fatalf("something wrong with reading of input: '%s'", err.Error())
		}
		LastResult := make(chan string, 1)
		for idx, pipe := range pipes {
			arguments := strings.Split(pipe, " ")
			if len(arguments) == 1 && arguments[0] == "" {
				continue
			}
			args := []string{strings.Join(arguments[1:], " ")}
			if idx != 0 {
				args = append(args, <-LastResult)
			}
			if arguments[len(arguments)-1] == "&" {
				s.fork(arguments[:len(arguments)-1], LastResult, &globalForForkW)
			} else {
				switch arguments[0] {
				case "cd":
					s.cd(args)
					LastResult <- ""
				case "pwd":
					s.pwd(LastResult)
				case "echo":
					s.echo(args, LastResult)
				case "kill":
					s.kill(args)
					LastResult <- ""
				case "ps":
					s.ps(args, LastResult)
				case "cat":
					s.cat(args, LastResult)
				case "\\q":
					fmt.Fprint(s.stdout, "quiting from shell\n")
					globalForForkW.Wait()
					break loop
				case "nc":
					s.netcat(arguments[1:])
					LastResult <- ""
				case "exec":
					s.exec(args)
					LastResult <- ""
				}
			}
			if idx == len(pipes)-1 {
				fmt.Fprint(s.stdout, <-LastResult)
			}
		}
	}
	w.Done()
}

func (s *shell) fork(arguments []string, p chan string, globalW *sync.WaitGroup) {
	// executing the process as a child process
	if len(arguments) == 0 {
		s.errorPrint("wrong number of arguments for fork function\n")
		return
	}
	cmd := exec.Command(arguments[0], arguments[1:]...)
	cmd.Stdin = s.stdin

	// creating a std pipeline
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		s.errorPrint(err.Error())
		return
	}

	err = cmd.Start()
	if err != nil {
		s.errorPrint(err.Error())
		return
	}
	fmt.Fprintf(s.stdout, "Child process with ID has been started: '%d'\n", cmd.Process.Pid)
	globalW.Add(1)
	go func(globalW *sync.WaitGroup) {
		out := make([]byte, 1024)
		w := sync.WaitGroup{}
		w.Add(1)
		go func() {
			b := strings.Builder{}
			for {
				n, err := stdout.Read(out)

				if err == io.EOF {
					break
				}
				if err != nil {
					s.errorPrint(err.Error())
					break
				}
				_, err = b.Write(out[:n])
				if err != nil {
					s.errorPrint(err.Error())
					break
				}
			}
			p <- b.String()
			w.Done()
		}()

		w.Wait()
		fmt.Fprintf(s.stdout, "Process with id '%d' is finished\n", cmd.Process.Pid)
		cmd.Wait()
		globalW.Done()
	}(globalW)
}

func newShell(startPath string) *shell {
	p := pathlib.NewPath(startPath)
	return &shell{
		CurrentDir:   p,
		PreviousPath: p,
		stdin:        os.Stdin,
		stdout:       os.Stdout,
		stderr:       os.Stderr,
	}
}

func main() {
	shell := newShell("./")
	w := sync.WaitGroup{}
	w.Add(1)
	go shell.cmd(os.Stdin, &w)
	w.Wait()
}

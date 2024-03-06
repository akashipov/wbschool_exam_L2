package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

/*
=== Утилита telnet ===

Реализовать примитивный telnet клиент:
Примеры вызовов:
go-telnet --timeout=10s host port
go-telnet mysite.ru 8080
go-telnet --timeout=3s 1.1.1.1 123

Программа должна подключаться к указанному хосту (ip или доменное имя) и порту по протоколу TCP.
После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s).

При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера, программа должна также завершаться.
При подключении к несуществующему сервер, программа должна завершаться через timeout.
*/

type TelnetArgs struct {
	host    string
	port    int64
	timeout int
}

type Telnet struct {
	args *TelnetArgs
}

func (t *Telnet) ParseArgs() {
	time := flag.Int("timeout", 10, "Timeout for connection to host")
	flag.Parse()
	if len(flag.Args()) != 2 {
		log.Fatal("Wrong number of arguments - need to pass host and port")
	}
	port, err := strconv.ParseInt(flag.Args()[1], 10, 64)
	if err != nil {
		log.Fatalf("Value '%s' cannot be casted to int32", flag.Args()[1])
	}
	t.args = &TelnetArgs{timeout: *time, host: flag.Args()[0], port: port}
}

func (t *Telnet) Read(done chan struct{}, info chan string) {
	reader := bufio.NewReader(os.Stdin)
	defer func() {
		close(info)
		os.Stdin.Close()
		fmt.Println("Reading is closed")
	}()
	buf := make([]byte, 1024)
	builder := strings.Builder{}
	fmt.Println("Please write the message for server:")
	for {
		n, err := reader.Read(buf)
		if err == io.EOF {
			if builder.String() != "" {
				info <- builder.String()
			}
			return
		}
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		builder.Write(buf[:n])
	}
}

func (t *Telnet) Connect(info chan string) *net.TCPConn {
	raddr := &net.TCPAddr{
		IP:   net.ParseIP(t.args.host),
		Port: int(t.args.port),
	}
	ticker := time.NewTicker(time.Second * time.Duration(t.args.timeout))
	var err error
	var conn *net.TCPConn
loop:
	for {
		select {
		case <-ticker.C:
			fmt.Println("Timeout to connect:", err.Error())
			return nil
		default:
			conn, err = net.DialTCP("tcp", nil, raddr)
			if err == nil {
				break loop
			}
		}
	}
	return conn
}

func (t *Telnet) Send(info chan string, done chan struct{}) {
	defer close(done)
	conn := t.Connect(info)
	if conn == nil {
		return
	}
	fmt.Println("TCP is set")
	defer func() {
		if conn != nil {
			conn.Close()
		}
		fmt.Println("TCP connection is closed")
	}()

	for i := range info {
		_, err := conn.Write([]byte(i))
		if err != nil {
			conn.Close()
			conn = t.Connect(info)
			if conn == nil {
				fmt.Println("Problem with writing to TCP server:", err.Error())
				break
			}
		}
	}
	fmt.Println("Got signal to close TCP connection")
}

func (t *Telnet) run() {
	w := sync.WaitGroup{}
	w.Add(2)
	done := make(chan struct{})
	info := make(chan string)
	go func() {
		t.Send(info, done)
		w.Done()
		fmt.Println("Sending is stopped")
	}()
	go t.Read(done, info)
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		select {
		case sig := <-sigint:
			fmt.Printf("\nSignal has got: '%v'\n", sig)
		case <-done:
		}
		w.Done()
		fmt.Println("Signal is stopped")
	}()
	w.Wait()
	fmt.Println("Telnet is Done")
}

func NewTelnet() *Telnet {
	t := Telnet{}
	t.ParseArgs()
	return &t
}

func main() {
	t := NewTelnet()
	t.run()
}

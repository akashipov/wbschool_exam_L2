package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/chigopher/pathlib"
)

func TestShell_errorPrint(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "base01",
			args: args{
				msg: "hello \n",
			},
			want: "hello",
		},
		{
			name: "base02",
			args: args{
				msg: "\nhello ",
			},
			want: "hello",
		},
		{
			name: "base03",
			args: args{
				msg: "Hello",
			},
			want: "hello",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := pathlib.NewPath("./")
			s := &Shell{
				PreviousPath: p,
				CurrentDir:   p,
				stdin:        os.Stdin,
				stdout:       os.Stdout,
				stderr:       os.Stderr,
			}
			r := s.errorPrint(tt.args.msg)
			if r != tt.want {
				t.Errorf("Wanted '%s', but has got '%s'\n", tt.want, r)
			}
		})
	}
}

func TestShell_kill(t *testing.T) {
	type fields struct {
		PreviousPath *pathlib.Path
		CurrentDir   *pathlib.Path
		stdin        io.Reader
		stdout       io.Writer
		stderr       io.Writer
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "base01",
			fields: fields{
				PreviousPath: nil,
				CurrentDir:   nil,
				stdin:        os.Stdin,
				stdout:       os.Stdout,
				stderr:       os.Stderr,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Shell{
				PreviousPath: tt.fields.PreviousPath,
				CurrentDir:   tt.fields.CurrentDir,
				stdin:        tt.fields.stdin,
				stdout:       tt.fields.stdout,
				stderr:       tt.fields.stderr,
			}
			cmd := exec.Command("sleep", "10")
			cmd.Start()
			defer func() {
				fmt.Println("Killing by defer function at the end")
				cmd.Process.Kill()
			}()
			pid := cmd.Process.Pid
			fmt.Println(pid)
			err := cmd.Process.Signal(syscall.Signal(0))
			if err != nil {
				t.Errorf("Something is wrong: %s\n", err.Error())
			}
			args := []string{fmt.Sprintf("%d", pid)}
			s.kill(args)
			p, err := os.FindProcess(pid)
			if err != nil {
				t.Error("Something is wrong:", err.Error())
			}
			state, err := p.Wait()
			if err != nil {
				t.Errorf("Problem with waiting")
			}
			fmt.Println(state)
		})
	}
}

func TestShell_exec(t *testing.T) {
	type args struct {
		filename string
		prefix   string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "output01",
			args: args{filename: "output.txt", prefix: ">"},
		},
		{
			name: "output02",
			args: args{filename: "/dev/tty", prefix: ">"},
		},
		{
			name: "input01",
			args: args{filename: "input.txt", prefix: "<"},
		},
		{
			name: "input02",
			args: args{filename: "/dev/tty", prefix: "<"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewShell("./")
			args := []string{tt.args.prefix + tt.args.filename}
			s.exec(args)
			defer func() {
				if strings.HasPrefix(tt.args.filename, "output") {
					fmt.Println("Removing files after creating")
					os.Remove(tt.args.filename)
				}
			}()
			switch tt.args.prefix {
			case ">":
				if tt.args.filename == "/dev/tty" {
					if s.stdout != os.Stdout {
						t.Error("Stdout wasn't changed to default")
					}
					return
				}
				wStr := "Hello World"
				fmt.Fprint(s.stdout, wStr)
				f, err := os.OpenFile(tt.args.filename, os.O_RDONLY, 0)
				if err != nil {
					t.Error(err.Error())
					return
				}
				b, err := io.ReadAll(f)
				if err != nil {
					t.Error(err.Error())
					return
				}
				if string(b) != wStr {
					t.Error("Stdout hasn't been changed\n")
					return
				}
			case "<":
				if tt.args.filename == "/dev/tty" {
					if s.stdin != os.Stdin {
						t.Error("Stdin wasn't changed to default")
					}
					return
				}
				line, err := s.readFromStdin()
				if err != nil {
					t.Error(err.Error())
					return
				}
				exp := "blabla\nbla"
				if line != exp {
					t.Errorf("We have got: '%s', but should be '%s'", line, exp)
					return
				}
			}
		})
	}
}

func TestShell_Fork(t *testing.T) {
	type args struct {
		arguments []string
	}
	output := make(chan string, 1)
	tests := []struct {
		name     string
		args     args
		want     string
		duration time.Duration
	}{
		{
			name:     "basetest01",
			args:     args{arguments: []string{"echo", "a"}},
			want:     "a\n",
			duration: 0,
		},
		{
			name:     "sleep02",
			args:     args{arguments: []string{"sleep", "1"}},
			want:     "",
			duration: time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewShell("./")
			s.Fork(tt.args.arguments, output, &sync.WaitGroup{})
			Time := time.Now()
			r := <-output
			since := time.Since(Time)
			if r != tt.want {
				t.Errorf("Wanted '%s', got '%s'", tt.want, r)
			}
			if since <= tt.duration {
				t.Errorf("Wanted duration '%v', got '%v'", tt.duration, since)
			}
		})
	}
}

func TestShell_cd(t *testing.T) {
	type args struct {
		args []string
	}
	currentDir := pathlib.NewPath(os.Getenv("PWD"))
	homeDir := pathlib.NewPath(os.Getenv("HOME"))
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "root01",
			args: args{[]string{""}},
			want: homeDir.String(),
		},
		{
			name: "back02",
			args: args{[]string{".."}},
			want: currentDir.Parent().String(),
		},
		{
			name: "back021",
			args: args{[]string{"./.."}},
			want: currentDir.Parent().String(),
		},
		{
			name: "back022",
			args: args{[]string{"../.."}},
			want: currentDir.Parent().Parent().String(),
		},
		{
			name: "home030",
			args: args{[]string{"/"}},
			want: "/",
		},
		{
			name: "home031",
			args: args{[]string{"/home"}},
			want: "/home",
		},
		{
			name: "previous04",
			args: args{[]string{"-"}},
			want: currentDir.String(),
		},
		{
			name: "next05_file",
			args: args{[]string{"./input.txt"}},
			want: currentDir.String(),
		},
		{
			name: "next06_dir",
			args: args{[]string{"./test"}},
			want: pathlib.NewPath(filepath.Join(os.Getenv("PWD"), "test")).String(),
		},
		{
			name: "next07_not_exist",
			args: args{[]string{"./abc"}},
			want: currentDir.String(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewShell("./")
			s.cd(tt.args.args)
			r := make(chan string, 1)
			s.pwd(r)
			result := <-r
			if tt.want != result {
				t.Errorf("Wanted '%s', got '%s'", tt.want, result)
			}
		})
	}
}

func TestShell_cmd(t *testing.T) {
	type args struct {
		commands string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "base01",
			args: args{"echo a\n"},
			want: []string{"./$a./$"},
		},
		{
			name: "quit02",
			args: args{"echo a\n\\q\n"},
			want: []string{"./$a./$"},
		},
		{
			name: "quit03",
			args: args{"\\q\n"},
			want: []string{"./$quiting from shell\n"},
		},
		{
			name: "pipe04",
			args: args{"echo a | echo b\n"},
			want: []string{"./$b./$"},
		},
		{
			name: "fork05",
			args: args{"echo a &\n"},
			want: []string{"Child process with ID has been started", "is finished"},
		},
		{
			name: "sleep06",
			args: args{"sleep 1 &\n"},
			want: []string{"Child process with ID has been started", "is finished"},
		},
		{
			name: "cd07",
			args: args{"cd ~/\n"},
			want: []string{fmt.Sprintf("Users/%s", os.Getenv("USER"))},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpfileOUT, err := os.CreateTemp("", "inputCMDout.txt")
			if err != nil {
				return
			}
			oldstdout := os.Stdout
			os.Stdout = tmpfileOUT
			defer func() {
				tmpfileOUT.Close()
				os.Remove(tmpfileOUT.Name())
				os.Stdout = oldstdout
			}()

			tmpfileIN, err := os.CreateTemp("", "inputCMDin.txt")
			if err != nil {
				return
			}
			defer func() {
				tmpfileIN.Close()
				os.Remove(tmpfileIN.Name())
			}()
			if _, err = tmpfileIN.Write([]byte(tt.args.commands)); err != nil {
				log.Fatal(err)
			}
			if _, err := tmpfileIN.Seek(0, 0); err != nil {
				log.Fatal(err)
			}
			oldstdin := os.Stdin
			os.Stdin = tmpfileIN
			defer func() {
				os.Stdin = oldstdin
			}()
			s := NewShell("./")
			var w sync.WaitGroup
			w.Add(1)
			s.cmd(tmpfileIN, &w)
			w.Wait()
			if _, err := tmpfileOUT.Seek(0, 0); err != nil {
				log.Fatal(err)
			}
			builder := strings.Builder{}
			b := make([]byte, 1024)
			for {
				n, err := tmpfileOUT.Read(b)
				if err == io.EOF {
					break
				}
				if err != nil {
					t.Error(err.Error())
					return
				}
				builder.Write(b[:n])
			}
			result := builder.String()
			for _, wanted := range tt.want {
				if !strings.Contains(result, wanted) {
					t.Errorf("Wanted '%s' inside the string '%s'", wanted, result)
				}
			}
		})
	}
}

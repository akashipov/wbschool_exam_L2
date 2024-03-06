package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/chigopher/pathlib"
	"github.com/go-chi/chi"
)

func Root(w http.ResponseWriter, request *http.Request) {
	w.Write([]byte(
		fmt.Sprintf("<html>Hello World! \"http://%s/a\"</html>", request.Host)),
	)
}

func A(w http.ResponseWriter, request *http.Request) {
	w.Write([]byte(
		fmt.Sprintf("<html>Hello It is A! Go to the \"http://%s/deeper\"</html>", request.Host),
	))
}

func FileTest(t *testing.T, dir string, want string) {
	f, err := os.OpenFile(dir, os.O_RDONLY, 0)
	if err != nil {
		t.Error(err.Error())
		return
	}
	textBytes, err := io.ReadAll(f)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if !strings.Contains(string(textBytes), want) {
		t.Errorf("'%s' must contain '%s'", string(textBytes), want)
		return
	}
}

func GetRoute() *chi.Mux {
	r := chi.NewRouter()
	r.Get(
		"/", Root,
	)
	r.Get(
		"/a", A,
	)
	return r
}
func Test_wget_execute(t *testing.T) {
	// Здесь можно было очень много всяких тестов написать но у меня просто не было времени
	type args struct {
		site        string
		dir         *pathlib.Path
		depth       int
		waitGP      *sync.WaitGroup
		isRecursive bool
	}

	srv := httptest.NewServer(GetRoute())
	defer srv.Close()
	w := sync.WaitGroup{}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantTexts map[string]string
	}{
		{
			name: "base01",
			args: args{
				site: srv.URL, depth: 0, waitGP: &w, dir: nil,
				isRecursive: true,
			},
			wantErr:   false,
			wantTexts: map[string]string{"main.html": "<html>Hello World!"},
		},
		{
			name: "base02",
			args: args{
				site: srv.URL, depth: 0, waitGP: &w, dir: nil,
				isRecursive: false,
			},
			wantErr:   false,
			wantTexts: map[string]string{"main.html": "<html>Hello World!"},
		},
		{
			name: "base03",
			args: args{
				site: srv.URL, depth: 0, waitGP: &w, dir: nil,
				isRecursive: true,
			},
			wantErr: false,
			wantTexts: map[string]string{
				"main.html":           "<html>Hello World! \"http",
				"depth_1/file_0.html": "<html>Hello It is A! Go to the \"http:",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("./", "test_")
			os.Chmod(tmpDir, 0777)
			if err != nil {
				t.Error(err.Error())
				return
			}
			defer func() {
				err := os.RemoveAll(tmpDir)
				if err != nil {
					t.Error(err.Error())
					return
				}
			}()
			_, err = os.Stat(tmpDir)
			if os.IsNotExist(err) {
				t.Error("File doesn't exist:", tmpDir)
				return
			}
			w, err := NewWget()
			w.args = wgetArgs{recursive: tt.args.isRecursive, directory: pathlib.NewPath(tmpDir), maxDepth: 2}
			if err != nil {
				t.Error(err.Error())
				return
			}
			tt.args.waitGP.Add(1)
			err = w.execute(
				tt.args.site, tt.args.dir,
				tt.args.depth, tt.args.waitGP, "main.html",
			)
			tt.args.waitGP.Wait()
			if (err != nil) != tt.wantErr {
				t.Errorf("wget.execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for k, v := range tt.wantTexts {
				var dir string
				if tt.args.isRecursive {
					dir = filepath.Join(tmpDir, strings.Replace(srv.URL, "/", "", -1), k)
				} else {
					dir = filepath.Join(tmpDir, k)
				}
				FileTest(t, dir, v)
			}
		})
	}
}

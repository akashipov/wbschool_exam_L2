package main

import (
	"io"
	"net"
	"os"
	"testing"
	"time"
)

func TestTelnet_run(t *testing.T) {
	// very simple test, I had no time to work with it more. Here need to work more
	type fields struct {
		args *TelnetArgs
	}
	addr := net.TCPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: 1234,
	}
	srv, err := net.ListenTCP("tcp", &addr)
	if err != nil {
		t.Error(err.Error())
		return
	}
	resp := ""
	go func() {
		conn, err := srv.AcceptTCP()
		if err != nil {
			t.Error(err.Error())
			return
		}
		defer conn.Close()
		for {
			p := make([]byte, 2048)
			n, err := conn.Read(p)
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Error(err.Error())
				return
			}
			resp += string(p[:n])
		}
	}()
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "base01",
			fields: fields{
				args: &TelnetArgs{
					timeout: 1,
					host:    "0.0.0.0",
					port:    1234,
				},
			},
			want: "abcdsfjsvg",
		},
		{
			name: "base02",
			fields: fields{
				args: &TelnetArgs{
					timeout: 1,
					host:    "0.0.0.0",
					port:    1235,
				},
			},
			want: "",
		},
		{
			name: "base03",
			fields: fields{
				args: &TelnetArgs{
					timeout: 1,
					host:    "0.0.0.128",
					port:    1234,
				},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp = ""
			tr := &Telnet{
				args: tt.fields.args,
			}
			f, err := os.CreateTemp("./", "test_")
			if err != nil {
				t.Error(err.Error())
				return
			}
			defer os.Remove(f.Name())
			f.Write([]byte("abcdsfjsvg"))
			f.Seek(0, 0)
			os.Stdin = f
			go tr.run()
			time.Sleep(time.Second)
			if resp != tt.want {
				t.Errorf("Want '%s', got '%s'", tt.want, resp)
			}
		})
	}
}

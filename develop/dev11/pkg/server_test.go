package pkg

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
)

func TestServerCalendar_GetRouter(t *testing.T) {
	type args struct {
		url         string
		requestData []byte
		method      string
	}
	log, _ := GetLogger()
	r := (&ServerCalendar{}).GetRouter(log)
	srv := httptest.NewServer(r)
	defer srv.Close()
	client := resty.New()
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "base01",
			args: args{
				url:         srv.URL + "/create_event",
				method:      http.MethodPost,
				requestData: []byte("{\"msg\":\"Hello\",\"date\":\"1995-01-01\"}"),
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r *resty.Response
			var err error
			if tt.args.method == http.MethodPost {
				r, err = client.R().SetBody(tt.args.requestData).Post(tt.args.url)
			} else {
				r, err = client.R().SetBody(tt.args.requestData).Get(tt.args.url)
			}
			if err != nil {
				t.Errorf("Something wrong with request to the server -> %s\n", err.Error())
				return
			}
			if string(r.Body()) != tt.want {
				t.Errorf("Wanted - '%s', got '%s'", tt.want, string(r.Body()))
			}
		})
	}
}

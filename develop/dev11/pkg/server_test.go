package pkg

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
)

func TestServerCalendar_GetRouter(t *testing.T) {
	type args struct {
		url          string
		requestData  []byte
		method       string
		clearStorage bool
	}
	log, _ := GetLogger()
	r := (&ServerCalendar{}).GetRouter(log)
	Storage = GetStorageOfEvents()
	srv := httptest.NewServer(r)
	defer srv.Close()
	client := resty.New()
	tests := []struct {
		name           string
		args           args
		want           string
		wantStatusCode int
	}{
		{
			name: "create01",
			args: args{
				url:         srv.URL + "/create_event",
				method:      http.MethodPost,
				requestData: []byte("{\"msg\":\"Hello\",\"date\":\"1995-01-01\"}"),
			},
			want:           "",
			wantStatusCode: 200,
		},
		{
			name: "create02",
			args: args{
				url:         srv.URL + "/create_event",
				method:      http.MethodPost,
				requestData: []byte("\"msg\":\"Hello\",\"date\":\"1995-01-01\"}"),
			},
			want:           "{\"error\":\"Problem with unmarshaling of request body: 'invalid character ':' after top-level value'\\n\"}",
			wantStatusCode: 400,
		},
		{
			name: "create03",
			args: args{
				url:         srv.URL + "/create_event",
				method:      http.MethodGet,
				requestData: []byte("\"msg\":\"Hello\",\"date\":\"1995-01-01\"}"),
			},
			want:           "{\"error\":\"Wrong method is used: allowed only 'POST'\\n\"}",
			wantStatusCode: 400,
		},
		{
			name: "create04",
			args: args{
				url:         srv.URL + "/create_event",
				method:      http.MethodPost,
				requestData: []byte("{\"msg\":\"Hello\",\"date\":\"1995-01-01\"}"),
			},
			want:           "{\"error\":\"message already there\"}",
			wantStatusCode: 400,
		},
		{
			name: "delete01",
			args: args{
				url:         srv.URL + "/delete_event",
				method:      http.MethodPost,
				requestData: []byte("{\"msg\":\"Hello\",\"date\":\"1995-01-01\"}"),
			},
			want:           "",
			wantStatusCode: 200,
		},
		{
			name: "delete02",
			args: args{
				url:         srv.URL + "/delete_event",
				method:      http.MethodGet,
				requestData: []byte("{\"msg\":\"Hello\",\"date\":\"1995-01-01\"}"),
			},
			want:           "{\"error\":\"Wrong method is used: allowed only 'POST'\\n\"}",
			wantStatusCode: 400,
		},
		{
			name: "delete03_again",
			args: args{
				url:         srv.URL + "/delete_event",
				method:      http.MethodPost,
				requestData: []byte("{\"msg\":\"Hello\",\"date\":\"1995-01-01\"}"),
			},
			want:           "",
			wantStatusCode: 200,
		},
		{
			name: "create01_after_deleting",
			args: args{
				url:         srv.URL + "/create_event",
				method:      http.MethodPost,
				requestData: []byte("{\"msg\":\"Hello0\",\"date\":\"1995-01-01\"}"),
			},
			want:           "",
			wantStatusCode: 200,
		},
		{
			name: "create02_after_deleting",
			args: args{
				url:         srv.URL + "/create_event",
				method:      http.MethodPost,
				requestData: []byte("{\"msg\":\"Hello1\",\"date\":\"1994-12-28\"}"),
			},
			want:           "",
			wantStatusCode: 200,
		},
		{
			name: "create03_after_deleting",
			args: args{
				url:         srv.URL + "/create_event",
				method:      http.MethodPost,
				requestData: []byte("{\"msg\":\"Hello2\",\"date\":\"1994-12-18\"}"),
			},
			want:           "",
			wantStatusCode: 200,
		},
		{
			name: "update01",
			args: args{
				url:         srv.URL + "/update_event",
				method:      http.MethodPost,
				requestData: []byte("{\"oldMsg\":\"Hello0\",\"msg\":\"Updated\",\"date\":\"1995-01-01\"}"),
			},
			want:           "",
			wantStatusCode: 200,
		},
		{
			name: "update02_wrong",
			args: args{
				url:         srv.URL + "/update_event",
				method:      http.MethodPost,
				requestData: []byte("{\"oldMsg\":\"Hello213\",\"msg\":\"New\",\"date\":\"1995-01-01\"}"),
			},
			want:           "{\"error\":\"There is no message for this day\"}",
			wantStatusCode: 400,
		},
		{
			name: "get_daily01",
			args: args{
				url:    srv.URL + "/events_for_day?date=1995-01-01",
				method: http.MethodGet,
			},
			want:           "{\"result\":[\"Updated\"]}",
			wantStatusCode: 200,
		},
		{
			name: "get_week01",
			args: args{
				url:    srv.URL + "/events_for_week?date=1995-01-01",
				method: http.MethodGet,
			},
			want:           "{\"result\":[\"Updated\",\"Hello1\"]}",
			wantStatusCode: 200,
		},
		{
			name: "get_month01",
			args: args{
				url:    srv.URL + "/events_for_month?date=1995-01-01",
				method: http.MethodGet,
			},
			want:           "{\"result\":[\"Updated\",\"Hello1\",\"Hello2\"]}",
			wantStatusCode: 200,
		},
		{
			name: "get_month02",
			args: args{
				url:         srv.URL + "/events_for_month?date=1995-01-01",
				method:      http.MethodPost,
				requestData: []byte(""),
			},
			want:           "{\"error\":\"Wrong method is used: allowed only 'GET'\\n\"}",
			wantStatusCode: 400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r *resty.Response
			var err error
			if tt.args.clearStorage {
				Storage = GetStorageOfEvents()
			}
			if tt.args.method == http.MethodPost {
				r, err = client.R().SetBody(tt.args.requestData).Post(tt.args.url)
			} else {
				r, err = client.R().Get(tt.args.url)
			}
			if err != nil {
				t.Errorf("Something wrong with request to the server -> %s\n", err.Error())
				fmt.Println("Body", r.Body())
				return
			}
			if string(r.Body()) != tt.want {
				t.Errorf("Wanted body - '%s', got '%s'", tt.want, string(r.Body()))
			}
			if r.StatusCode() != tt.wantStatusCode {
				t.Errorf("Wanted status code - '%d', got '%d'", tt.wantStatusCode, r.StatusCode())
			}
		})
	}
}

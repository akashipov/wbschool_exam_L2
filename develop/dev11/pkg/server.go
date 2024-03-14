package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

// ServerCalendar - our api, with host and port fields
type ServerCalendar struct {
	Host string
	Port int
}

// AllowedMethod - checking that method of reqeust is correct(satisfied our 'method' param)
func AllowedMethod(method string, logic func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case method:
			logic(w, r)
		default:
			err := Error{Value: fmt.Sprintf("Wrong method is used: allowed only '%s'\n", method)}
			b, e := json.Marshal(err)
			if e != nil {
				fmt.Println(e)
			}
			_, e = w.Write(b)
			if e != nil {
				fmt.Println(e)
			}
		}
	}
}

// GetRouter - get router with all main handlers
func (s *ServerCalendar) GetRouter(log *zap.SugaredLogger) *http.ServeMux {
	customMux := http.NewServeMux()
	customMux.HandleFunc("/create_event", WithLogging(http.HandlerFunc(AllowedMethod("POST", CreateEvent)), log))
	customMux.HandleFunc("/update_event", WithLogging(http.HandlerFunc(AllowedMethod("POST", UpdateEvent)), log))
	customMux.HandleFunc("/delete_event", WithLogging(http.HandlerFunc(AllowedMethod("POST", DeleteEvent)), log))
	customMux.HandleFunc("/events_for_day", WithLogging(http.HandlerFunc(AllowedMethod("GET", GetEventsForDaily)), log))
	customMux.HandleFunc("/events_for_month", WithLogging(http.HandlerFunc(AllowedMethod("GET", GetEventsForMonth)), log))
	customMux.HandleFunc("/events_for_week", WithLogging(http.HandlerFunc(AllowedMethod("GET", GetEventsForWeek)), log))
	return customMux
}

// RunServer - general structor of our api, list of all handlers
// and starting to listen requests
func (s *ServerCalendar) RunServer(log *zap.SugaredLogger) {
	fullHost := fmt.Sprintf("%s:%d", s.Host, s.Port)
	err := http.ListenAndServe(fullHost, s.GetRouter(log))
	if err != nil {
		fmt.Println(err.Error())
	}
}

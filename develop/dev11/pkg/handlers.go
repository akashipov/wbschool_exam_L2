package pkg

import (
	"fmt"
	"net/http"
)

// GetEventsForDaily - get all events in the calendar for day (1 day)
func GetEventsForDaily(w http.ResponseWriter, request *http.Request) {
	ReturnEventsForDates(w, request, 1)
}

// GetEventsForWeek - get all events in the calendar for week (7 days)
func GetEventsForWeek(w http.ResponseWriter, request *http.Request) {
	ReturnEventsForDates(w, request, 7)
}

// GetEventsForMonth - get all events in the calendar for month (30 days)
func GetEventsForMonth(w http.ResponseWriter, request *http.Request) {
	// it is like example... it is not correct, need to change it to more properly version
	ReturnEventsForDates(w, request, 30)
}

// CreateEvent - create event, add to storage
func CreateEvent(w http.ResponseWriter, request *http.Request) {
	defer func() {
		fmt.Println("create finished")
		fmt.Printf("%+v\n", Storage)
	}()
	e, err := UnmarshalEvent(w, request)
	if err != nil {
		return
	}
	err = Storage.AddEvent(*e)
	if err != nil {
		ReportError(
			w,
			http.StatusBadRequest,
			err.Error(),
		)
	}
}

// UpdateEvent - update element from storage
func UpdateEvent(w http.ResponseWriter, request *http.Request) {
	defer func() {
		fmt.Println("update finished")
		fmt.Printf("%+v\n", Storage)
	}()
	e, err := UnmarshalEvent(w, request)
	if err != nil {
		return
	}
	err = Storage.UpdateEvent(*e)
	if err != nil {
		ReportError(
			w,
			http.StatusBadRequest,
			err.Error(),
		)
	}
}

// DeleteEvent - delete element from storage
func DeleteEvent(w http.ResponseWriter, request *http.Request) {
	defer func() {
		fmt.Println("delete finished")
		fmt.Printf("%+v\n", Storage)
	}()
	e, err := UnmarshalEvent(w, request)
	if err != nil {
		return
	}
	Storage.DeleteEvent(*e)
}

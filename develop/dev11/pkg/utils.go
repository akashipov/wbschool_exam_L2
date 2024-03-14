// I know that it is bad practice to use utils name for file or package
// also it is bad decision to save everything in one folder. There is left only one day to do this task
package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var format string = "2006-01-02"

// ValidateDate - function with logic of validation of passed date(string)
// format is taken from the global variable
func ValidateDate(date string) (time.Time, error) {
	return time.Parse(format, date)
}

// UnmarshalEvent - parse method for data struct Event(checking on validation of date included)
func UnmarshalEvent(w http.ResponseWriter, r *http.Request) (*Event, error) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		ReportError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("Problem with reading of request body: '%s'\n", err),
		)
		return nil, err
	}
	e := Event{}
	err = json.Unmarshal(b, &e)
	if err != nil {
		ReportError(
			w,
			http.StatusBadRequest,
			fmt.Sprintf("Problem with unmarshaling of request body: '%s'\n", err),
		)
		return nil, err
	}
	_, err = ValidateDate(e.Date)
	if err != nil {
		ReportError(
			w,
			http.StatusBadRequest,
			fmt.Sprintf("Problem with date format: '%s'\n", e.Date),
		)
		return nil, err
	}
	return &e, nil
}

// ReturnEventsForDates - getting all list of events in the calendar
// from interval [passed day - Day * (length - 1), passed day]
func ReturnEventsForDates(w http.ResponseWriter, r *http.Request, length int) {
	date := r.URL.Query().Get("date")
	t, err := ValidateDate(date)
	if err != nil {
		ReportError(
			w,
			http.StatusBadRequest,
			err.Error(),
		)
	}
	events := Events{}
	for i := 0; i < length; i++ {
		t := t.Add(-time.Hour * time.Duration(24*i))
		for k := range Storage.DateMap[t.Format(format)] {
			events.Results = append(events.Results, k)
		}
	}
	b, err := json.Marshal(events)
	if err != nil {
		ReportError(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
	}
	w.Write(b)
}

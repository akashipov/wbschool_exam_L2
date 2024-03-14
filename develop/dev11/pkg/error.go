package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Error - json working struct with errors
type Error struct {
	Value string `json:"error"`
}

// ReportError - trying to marshal error like json, assign status of error to the header
func ReportError(w http.ResponseWriter, status int, errMsg string) {
	w.WriteHeader(status)
	b, err := json.Marshal(Error{Value: errMsg})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	_, err = w.Write([]byte(b))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

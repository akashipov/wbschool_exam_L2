package pkg

import (
	"errors"
	"sync"
)

// Storage - Global variable for our storage(as simple example) better to use Postgres for example
var Storage StorageOfEvents

// StorageOfEvents - basical storage to work with events of calendar
type StorageOfEvents struct {
	DateMap map[string]map[string]string
	m       sync.Mutex
}

// Events - is used to get results of all our events for the period of time
type Events struct {
	Results []string `json:"result"`
}

// Event - general json structor for our events in the our calendar application
type Event struct {
	Date       string `json:"date"`
	Message    string `json:"msg"`
	OldMessage string `json:"oldMsg,omitempty"`
}

// GetStorageOfEvents - default constructor for 'StorageOfEvents'
func GetStorageOfEvents() StorageOfEvents {
	s := StorageOfEvents{
		make(map[string]map[string]string),
		sync.Mutex{},
	}
	return s
}

// AddEvent - add event to the calendar if it doesn't exist, otherwise error
func (strg *StorageOfEvents) AddEvent(e Event) error {
	strg.m.Lock()
	defer func() {
		strg.m.Unlock()
	}()
	v, ok := strg.DateMap[e.Date]
	if ok {
		_, ok := v[e.Message]
		if ok {
			return errors.New("message already there")
		}
		v[e.Message] = e.Message
		return nil
	}
	strg.DateMap[e.Date] = map[string]string{e.Message: e.Message}
	return nil
}

// UpdateEvent - update if it has the event
func (strg *StorageOfEvents) UpdateEvent(e Event) error {
	strg.m.Lock()
	defer func() {
		strg.m.Unlock()
	}()
	errMsg := "There is no message for this day"
	v, ok := strg.DateMap[e.Date]
	if ok {
		_, ok := v[e.OldMessage]
		if !ok {
			return errors.New(errMsg)
		}
		delete(v, e.OldMessage)
		v[e.Message] = e.Message
		return nil
	}
	return errors.New(errMsg)
}

// DeleteEvent - delete if it exists
func (strg *StorageOfEvents) DeleteEvent(e Event) {
	strg.m.Lock()
	defer func() {
		strg.m.Unlock()
	}()
	v, ok := strg.DateMap[e.Date]
	if ok {
		_, ok := v[e.Message]
		if ok {
			delete(v, e.Message)
			if len(v) == 0 {
				delete(strg.DateMap, e.Date)
			}
		}
	}
}

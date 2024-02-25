package models

import "time"

// our events sturcture and models package will have methods related to deatabase stuff
type Event struct {
	ID          int
	Name        string    `binding:"required"`
	Description string    `binding:"required"`
	Location    string    `binding:"required"`
	DateTime    time.Time `binding:"required"`
	UserID      int
}

var events = []Event{}

func (e *Event) Save() {
	// adding to the database

	events = append(events, *e)
}

// call all available events
func GetAllEvents() []Event {

	return events
}

package main

import "errors"

type events struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// For this demo, we're storing the event list in memory
// In a real application, this list will most likely be fetched
// from a database or from static files
var eventList = []events{
	events{ID: 1, Title: "events 1", Content: "events 1 body"},
	events{ID: 2, Title: "events 2", Content: "events 2 body"},
}

// Return a list of all the events
func getAllEvent() []events {
	return eventList
}

// Fetch an event based on the ID supplied
func getEventByID(id int) (*events, error) {
	for _, a := range eventList {
		if a.ID == id {
			return &a, nil
		}
	}
	return nil, errors.New("Event not found")
}

// Create a new event with the title and content provided
func createNewEvent(title, content string) (*events, error) {
	// Set the ID of a new event to one more than the number of events
	a := events{ID: len(eventList) + 1, Title: title, Content: content}

	// Add the event to the list of events
	eventList = append(eventList, a)

	return &a, nil
}

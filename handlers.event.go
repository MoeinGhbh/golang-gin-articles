package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func showIndexPage(c *gin.Context) {
	events := getAllEvent()

	// Call the render function with the name of the template to render
	render(c, gin.H{
		"title":   "Home Page",
		"payload": events}, "index.html")
}

func showeventCreationPage(c *gin.Context) {
	// Call the render function with the name of the template to render
	render(c, gin.H{
		"title": "Create New event"}, "create-event.html")
}

func getevent(c *gin.Context) {
	// Check if the event ID is valid
	if eventID, err := strconv.Atoi(c.Param("event_id")); err == nil {
		// Check if the event exists
		if event, err := getEventByID(eventID); err == nil {
			// Call the render function with the title, event and the name of the
			// template
			render(c, gin.H{
				"title":   event.Title,
				"payload": event}, "event.html")

		} else {
			// If the event is not found, abort with an error
			c.AbortWithError(http.StatusNotFound, err)
		}

	} else {
		// If an invalid event ID is specified in the URL, abort with an error
		c.AbortWithStatus(http.StatusNotFound)
	}
}

func createevent(c *gin.Context) {
	// Obtain the POSTed title and content values
	title := c.PostForm("title")
	content := c.PostForm("content")

	if a, err := createNewEvent(title, content); err == nil {
		// If the event is created successfully, show success message
		render(c, gin.H{
			"title":   "Submission Successful",
			"payload": a}, "submission-successful.html")
	} else {
		// if there was an error while creating the event, abort with an error
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

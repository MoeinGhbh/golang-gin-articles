// handlers.event_test.go

package main

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

// Test that a GET request to the home page returns the home page with
// the HTTP code 200 for an unauthenticated user
func TestShowIndexPageUnauthenticated(t *testing.T) {
	r := getRouter(true)

	r.GET("/", showIndexPage)

	// Create a request to send to the above route
	req, _ := http.NewRequest("GET", "/", nil)

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		// Test that the http status code is 200
		statusOK := w.Code == http.StatusOK

		// Test that the page title is "Home Page"
		// You can carry out a lot more detailed tests using libraries that can
		// parse and process HTML pages
		p, err := ioutil.ReadAll(w.Body)
		pageOK := err == nil && strings.Index(string(p), "<title>Home Page</title>") > 0

		return statusOK && pageOK
	})
}

// Test that a GET request to the home page returns the home page with
// the HTTP code 200 for an authenticated user
func TestShowIndexPageAuthenticated(t *testing.T) {
	// Create a response recorder
	w := httptest.NewRecorder()

	// Get a new router
	r := getRouter(true)

	// Set the token cookie to simulate an authenticated user
	http.SetCookie(w, &http.Cookie{Name: "token", Value: "123"})

	// Define the route similar to its definition in the routes file
	r.GET("/", showIndexPage)

	// Create a request to send to the above route
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header = http.Header{"Cookie": w.HeaderMap["Set-Cookie"]}

	// Create the service and process the above request.
	r.ServeHTTP(w, req)

	// Test that the http status code is 200
	if w.Code != http.StatusOK {
		t.Fail()
	}

	// Test that the page title is "Home Page"
	// You can carry out a lot more detailed tests using libraries that can
	// parse and process HTML pages
	p, err := ioutil.ReadAll(w.Body)
	if err != nil || strings.Index(string(p), "<title>Home Page</title>") < 0 {
		t.Fail()
	}

}

// Test that a GET request to an event page returns the event page with
// the HTTP code 200 for an unauthenticated user
func TesteventUnauthenticated(t *testing.T) {
	r := getRouter(true)

	// Define the route similar to its definition in the routes file
	r.GET("/event/view/:event_id", getevent)

	// Create a request to send to the above route
	req, _ := http.NewRequest("GET", "/event/view/1", nil)

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		// Test that the http status code is 200
		statusOK := w.Code == http.StatusOK

		// Test that the page title is "event 1"
		// You can carry out a lot more detailed tests using libraries that can
		// parse and process HTML pages
		p, err := ioutil.ReadAll(w.Body)
		pageOK := err == nil && strings.Index(string(p), "<title>event 1</title>") > 0

		return statusOK && pageOK
	})
}

// Test that a GET request to an event page returns the event page with
// the HTTP code 200 for an authenticated user
func TesteventAuthenticated(t *testing.T) {
	// Create a response recorder
	w := httptest.NewRecorder()

	// Get a new router
	r := getRouter(true)

	// Set the token cookie to simulate an authenticated user
	http.SetCookie(w, &http.Cookie{Name: "token", Value: "123"})

	// Define the route similar to its definition in the routes file
	r.GET("/event/view/:event_id", getevent)

	// Create a request to send to the above route
	req, _ := http.NewRequest("GET", "/event/view/1", nil)
	req.Header = http.Header{"Cookie": w.HeaderMap["Set-Cookie"]}

	// Create the service and process the above request.
	r.ServeHTTP(w, req)

	// Test that the http status code is 200
	if w.Code != http.StatusOK {
		t.Fail()
	}

	// Test that the page title is "event 1"
	// You can carry out a lot more detailed tests using libraries that can
	// parse and process HTML pages
	p, err := ioutil.ReadAll(w.Body)
	if err != nil || strings.Index(string(p), "<title>event 1</title>") < 0 {
		t.Fail()
	}

}

// Test that a GET request to the home page returns the list of events
// in JSON format when the Accept header is set to application/json
func TesteventListJSON(t *testing.T) {
	r := getRouter(true)

	// Define the route similar to its definition in the routes file
	r.GET("/", showIndexPage)

	// Create a request to send to the above route
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept", "application/json")

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		// Test that the http status code is 200
		statusOK := w.Code == http.StatusOK

		// Test that the response is JSON which can be converted to
		// an array of event structs
		p, err := ioutil.ReadAll(w.Body)
		if err != nil {
			return false
		}
		var events []event
		err = json.Unmarshal(p, &events)

		return err == nil && len(events) >= 2 && statusOK
	})
}

// Test that a GET request to an event page returns the event in XML
// format when the Accept header is set to application/xml
func TesteventXML(t *testing.T) {
	r := getRouter(true)

	// Define the route similar to its definition in the routes file
	r.GET("/event/view/:event_id", getevent)

	// Create a request to send to the above route
	req, _ := http.NewRequest("GET", "/event/view/1", nil)
	req.Header.Add("Accept", "application/xml")

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		// Test that the http status code is 200
		statusOK := w.Code == http.StatusOK

		// Test that the response is JSON which can be converted to
		// an array of event structs
		p, err := ioutil.ReadAll(w.Body)
		if err != nil {
			return false
		}
		var a event
		err = xml.Unmarshal(p, &a)

		return err == nil && a.ID == 1 && len(a.Title) >= 0 && statusOK
	})
}

// Test that a GET request to the event creation page returns the
// event creation page with the HTTP code 200 for an authenticated user
func TesteventCreationPageAuthenticated(t *testing.T) {
	// Create a response recorder
	w := httptest.NewRecorder()

	// Get a new router
	r := getRouter(true)

	// Set the token cookie to simulate an authenticated user
	http.SetCookie(w, &http.Cookie{Name: "token", Value: "123"})

	// Define the route similar to its definition in the routes file
	r.GET("/event/create", ensureLoggedIn(), showeventCreationPage)

	// Create a request to send to the above route
	req, _ := http.NewRequest("GET", "/event/create", nil)
	req.Header = http.Header{"Cookie": w.HeaderMap["Set-Cookie"]}

	// Create the service and process the above request.
	r.ServeHTTP(w, req)

	// Test that the http status code is 200
	if w.Code != http.StatusOK {
		t.Fail()
	}

	// Test that the page title is "Create New event"
	// You can carry out a lot more detailed tests using libraries that can
	// parse and process HTML pages
	p, err := ioutil.ReadAll(w.Body)
	if err != nil || strings.Index(string(p), "<title>Create New event</title>") < 0 {
		t.Fail()
	}

}

// Test that a GET request to the event creation page returns
// an HTTP 401 error for an unauthorized user
func TesteventCreationPageUnauthenticated(t *testing.T) {
	r := getRouter(true)

	// Define the route similar to its definition in the routes file
	r.GET("/event/create", ensureLoggedIn(), showeventCreationPage)

	// Create a request to send to the above route
	req, _ := http.NewRequest("GET", "/event/create", nil)

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		// Test that the http status code is 401
		return w.Code == http.StatusUnauthorized
	})
}

// Test that a POST request to create an event returns
// an HTTP 200 code along with a success message for an authenticated user
func TesteventCreationAuthenticated(t *testing.T) {
	// Create a response recorder
	w := httptest.NewRecorder()

	// Get a new router
	r := getRouter(true)

	// Set the token cookie to simulate an authenticated user
	http.SetCookie(w, &http.Cookie{Name: "token", Value: "123"})

	// Define the route similar to its definition in the routes file
	r.POST("/event/create", ensureLoggedIn(), createevent)

	// Create a request to send to the above route
	eventPayload := geteventPOSTPayload()
	req, _ := http.NewRequest("POST", "/event/create", strings.NewReader(eventPayload))
	req.Header = http.Header{"Cookie": w.HeaderMap["Set-Cookie"]}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(eventPayload)))

	// Create the service and process the above request.
	r.ServeHTTP(w, req)

	// Test that the http status code is 200
	if w.Code != http.StatusOK {
		t.Fail()
	}

	// Test that the page title is "Submission Successful"
	// You can carry out a lot more detailed tests using libraries that can
	// parse and process HTML pages
	p, err := ioutil.ReadAll(w.Body)
	if err != nil || strings.Index(string(p), "<title>Submission Successful</title>") < 0 {
		t.Fail()
	}

}

// Test that a POST request to create an event returns
// an HTTP 401 error for an unauthorized user
func TesteventCreationUnauthenticated(t *testing.T) {
	r := getRouter(true)

	// Define the route similar to its definition in the routes file
	r.POST("/event/create", ensureLoggedIn(), createevent)

	// Create a request to send to the above route
	eventPayload := geteventPOSTPayload()
	req, _ := http.NewRequest("POST", "/event/create", strings.NewReader(eventPayload))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(eventPayload)))

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		// Test that the http status code is 401
		return w.Code == http.StatusUnauthorized
	})
}

func geteventPOSTPayload() string {
	params := url.Values{}
	params.Add("title", "Test event Title")
	params.Add("content", "Test event Content")

	return params.Encode()
}

package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"os"
)

const (

	// WhitePixelB64 is a PNG that is 1x1 pixels that are white.
	WhitePixelB64 = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABAQMAAAAl21bKAAAAA1BMVEUAAACnej3aAAAAAXRSTlMAQObYZgAAAApJREFUCNdjYAAAAAIAAeIhvDMAAAAASUVORK5CYII="
)

// HandleEverything is an HTTP handler that will handle everything.
type HandleEverything struct {
	Body           []byte
	HandleError    func(err error)
	Headers        http.Header
	HandleRequests func(req *http.Request)
	Status         int
}

func (h HandleEverything) ServeHTTP(writer http.ResponseWriter, req *http.Request) {

	// Write the response code.
	writer.WriteHeader(h.Status)

	// Write the response headers.
	for key, values := range h.Headers {
		for _, value := range values {
			writer.Header().Set(key, value)
		}
	}

	// Write the body of the response.
	var err error
	if _, err = writer.Write(h.Body); err != nil {
		go h.HandleError(err)
		return
	}

	// Process the request asynchronously.
	go h.HandleRequests(req)
}

func main() {

	// Create a logger.
	logger := log.New(os.Stdout, "", 0)

	// Create an error handler.
	errorHandler := func(err error) {
		logger.Printf("An error happened asynchronously.\nError: %s\n", err.Error())
	}

	// Create a request handler.
	requestHandler := func(req *http.Request) {
		logger.Printf("URL requested: %s", req.URL.String())
		logger.Printf("  RemoteAddr: %s", req.RemoteAddr)
		logger.Printf("  ForwardedFor: %s", req.Header.Get("X-Forwarded-For"))
	}

	// Create the response headers.
	headers := map[string][]string{
		"Content-Type": {"image/png"},
	}

	// Create the HTTP handler.
	handler := HandleEverything{
		HandleError:    errorHandler,
		Headers:        headers,
		HandleRequests: requestHandler,
		Status:         200,
	}

	// Create the response body as a byte slice.
	var err error
	if handler.Body, err = base64.StdEncoding.DecodeString(WhitePixelB64); err != nil {
		logger.Fatalf("Failed to decode response body.\nError: %s", err.Error())
	}

	// Start the service.
	if err = http.ListenAndServe(":8080", handler); err != nil {
		logger.Fatalf("Failed to serve.\nError: %s", err.Error())
	}
}

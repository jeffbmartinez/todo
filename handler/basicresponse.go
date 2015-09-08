package handler

import (
	"net/http"
)

// BasicResponse creates a handler which responds with a standard response
// code and message string.
func BasicResponse(code int) func(response http.ResponseWriter, request *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		response.WriteHeader(code)
		response.Write([]byte(http.StatusText(code)))
	}
}

// WriteBasicResponse responds the the request with a BasicResponse
func WriteBasicResponse(code int, response http.ResponseWriter) {
	BasicResponse(code)(response, nil)
}

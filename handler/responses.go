package handler

import (
	"encoding/json"
	"net/http"

	"github.com/jeffbmartinez/log"
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

/*
WriteJSONResponse writes a json response (with correct http header).
*/
func WriteJSONResponse(response http.ResponseWriter, message interface{}, statusCode int) {
	responseString, err := json.Marshal(message)
	if err != nil {
		log.Errorf("Couldn't marshal json: %v", err)

		WriteBasicResponse(http.StatusInternalServerError, response)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(statusCode)
	response.Write([]byte(responseString))
}

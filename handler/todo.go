package handler

import (
	"net/http"
)

// Todo handles requests to the todo endpoint.
func Todo(response http.ResponseWriter, request *http.Request) {
	handler := BasicResponse(http.StatusMethodNotAllowed)

	switch request.Method {
	case "GET":
		handler = getTodo
	case "POST":
		handler = postTodo
	}

	handler(response, request)
}

func getTodo(response http.ResponseWriter, request *http.Request) {
	BasicResponse(http.StatusOK)
}

func postTodo(response http.ResponseWriter, request *http.Request) {
	BasicResponse(http.StatusOK)
}

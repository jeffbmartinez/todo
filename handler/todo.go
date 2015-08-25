package handler

import (
	"net/http"

	"github.com/jeffbmartinez/log"

	"github.com/jeffbmartinez/todo/storage"
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
	taskset, err := storage.GetTaskset()
	if err != nil {
		log.Errorf("Couldn't get taskset (%v)", err)
		BasicResponse(http.StatusInternalServerError)
		return
	}

	WriteJSONResponse(response, taskset, http.StatusOK)
}

func postTodo(response http.ResponseWriter, request *http.Request) {
	taskset, err := storage.GetTaskset()
	if err != nil {
		log.Errorf("Couldn't get taskset (%v)", err)
		BasicResponse(http.StatusInternalServerError)
		return
	}

	if err := storage.SaveTaskset(taskset); err != nil {
		log.Errorf("Couldn't save taskset (%v)", err)
		BasicResponse(http.StatusInternalServerError)
		return
	}

	BasicResponse(http.StatusOK)
}

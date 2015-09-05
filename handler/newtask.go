package handler

import (
	"net/http"

	"github.com/jeffbmartinez/log"

	"github.com/jeffbmartinez/todo/storage"
)

// NewTask handles requests to the /tasks/new endpoint.
func NewTask(response http.ResponseWriter, request *http.Request) {
	handler := BasicResponse(http.StatusMethodNotAllowed)

	switch request.Method {
	case "POST":
		handler = postNewTask
	}

	handler(response, request)
}

func postNewTask(response http.ResponseWriter, request *http.Request) {
	taskset, err := storage.GetTaskset()
	if err != nil {
		log.Errorf("Couldn't get taskset (%v)", err)
		WriteBasicResponse(http.StatusInternalServerError, response, request)
		return
	}

	if err := storage.SaveTaskset(taskset); err != nil {
		log.Errorf("Couldn't save taskset (%v)", err)
		WriteBasicResponse(http.StatusInternalServerError, response, request)
		return
	}

	WriteBasicResponse(http.StatusNotImplemented, response, request)
}

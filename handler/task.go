package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jeffbmartinez/log"

	"github.com/jeffbmartinez/todo/storage"
)

// Task handles requests to the /task/{id} endpoint.
func Task(response http.ResponseWriter, request *http.Request) {
	handler := BasicResponse(http.StatusMethodNotAllowed)

	switch request.Method {
	case "GET":
		handler = getTask
	case "PUT":
		handler = putTask
	case "DELETE":
		handler = deleteTask
	}

	handler(response, request)
}

func getTask(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	taskID := vars["id"]

	taskset, err := storage.GetTaskset()
	if err != nil {
		log.Errorf("Couldn't get taskset (%v)", err)
		WriteBasicResponse(http.StatusInternalServerError, response, request)
		return
	}

	task, ok := taskset.Get(taskID)
	if !ok {
		WriteBasicResponse(http.StatusNotFound, response, request)
		return
	}

	WriteJSONResponse(response, task, http.StatusOK)
}

func putTask(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	taskID := vars["id"]

	taskset, err := storage.GetTaskset()
	if err != nil {
		log.Errorf("Couldn't get taskset (%v)", err)
		WriteBasicResponse(http.StatusInternalServerError, response, request)
		return
	}

	task, ok := taskset.Get(taskID)
	if !ok {
		WriteBasicResponse(http.StatusNotFound, response, request)
		return
	}

	// TODO: Modify the task

	err = taskset.Add(task)
	if err != nil {
		log.Errorf("Couldn't add modified task to taskset (%v)", err)
		WriteBasicResponse(http.StatusInternalServerError, response, request)
		return
	}

	err = storage.SaveTaskset(taskset)
	if err != nil {
		log.Errorf("Couldn't save taskset (%v)", err)
		WriteBasicResponse(http.StatusInternalServerError, response, request)
		return
	}

	WriteBasicResponse(http.StatusNotImplemented, response, request)
}

func deleteTask(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	taskID := vars["id"]

	taskset, err := storage.GetTaskset()
	if err != nil {
		log.Errorf("Couldn't get taskset (%v)", err)
		WriteBasicResponse(http.StatusInternalServerError, response, request)
		return
	}

	task, ok := taskset.Get(taskID)
	if !ok {
		WriteBasicResponse(http.StatusNotFound, response, request)
		return
	}

	taskset.Delete(task)

	err = storage.SaveTaskset(taskset)
	if err != nil {
		log.Errorf("Couldn't save taskset (%v)", err)
		WriteBasicResponse(http.StatusInternalServerError, response, request)
		return
	}

	WriteBasicResponse(http.StatusOK, response, request)
}

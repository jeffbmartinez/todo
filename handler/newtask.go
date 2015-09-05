package handler

import (
	"encoding/json"
	"net/http"

	"github.com/jeffbmartinez/log"

	"github.com/jeffbmartinez/todo/storage"
	"github.com/jeffbmartinez/todo/task"
)

/*
TaskParams is the json struct that gets passed in the request to create
a new Task object.
*/
type TaskParams struct {
	Name     string `json:"name"`
	ParentID string `json:"parentID"`
}

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
	if request.Body == nil {
		WriteBasicResponse(http.StatusBadRequest, response, request)
		return
	}

	defer request.Body.Close()
	decoder := json.NewDecoder(request.Body)

	var taskParams TaskParams
	err := decoder.Decode(&taskParams)
	if err != nil {
		WriteBasicResponse(http.StatusBadRequest, response, request)
		return
	}

	newTask, err := task.NewTask(taskParams.Name, taskParams.ParentID)
	if err != nil {
		log.Errorf("Couldn't create new task (%v)", err)
		WriteBasicResponse(http.StatusInternalServerError, response, request)
		return
	}

	taskset, err := storage.GetTaskset()
	if err != nil {
		log.Errorf("Couldn't get taskset (%v)", err)
		WriteBasicResponse(http.StatusInternalServerError, response, request)
		return
	}

	taskset.Add(newTask)

	if err := storage.SaveTaskset(taskset); err != nil {
		log.Errorf("Couldn't save taskset (%v)", err)
		WriteBasicResponse(http.StatusInternalServerError, response, request)
		return
	}

	WriteBasicResponse(http.StatusOK, response, request)
}

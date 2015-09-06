package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jeffbmartinez/log"

	"github.com/jeffbmartinez/todo/storage"
)

/*
UpdateTaskParams is the json struct that gets passed in the request to update
a Task object.
*/
type UpdateTaskParams struct {
	Name       string   `json:"name"`
	Complete   bool     `json:"complete"`
	SubtaskIDs []string `json:"subtaskIDs"`
	ParentIDs  []string `json:"parentIDs"`
	DueDate    int64    `json:"dueDate"`
	Categories []string `json:"categories"`
}

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

	tasklist, err := storage.GetTasklist()
	if err != nil {
		log.Errorf("Couldn't get tasklist (%v)", err)
		WriteBasicResponse(http.StatusInternalServerError, response, request)
		return
	}

	task, ok := tasklist.Registry[taskID]
	if !ok {
		WriteBasicResponse(http.StatusNotFound, response, request)
		return
	}

	WriteJSONResponse(response, task, http.StatusOK)
}

func putTask(response http.ResponseWriter, request *http.Request) {
	if request.Body == nil {
		WriteBasicResponse(http.StatusBadRequest, response, request)
		return
	}

	defer request.Body.Close()
	decoder := json.NewDecoder(request.Body)

	var params UpdateTaskParams
	err := decoder.Decode(&params)
	if err != nil {
		WriteBasicResponse(http.StatusBadRequest, response, request)
		return
	}

	vars := mux.Vars(request)
	taskID := vars["id"]

	tasklist, err := storage.GetTasklist()
	if err != nil {
		log.Errorf("Couldn't get tasklist (%v)", err)
		WriteBasicResponse(http.StatusInternalServerError, response, request)
		return
	}

	task, ok := tasklist.Registry[taskID]
	if !ok {
		WriteBasicResponse(http.StatusNotFound, response, request)
		return
	}

	// TODO: Modify the task
	if params.Name != "" {
		task.Name = params.Name
	}

	task.SetComplete(params.Complete)
	if params.Categories != nil {
		task.Categories = params.Categories
	}

	if params.DueDate != 0 {
		task.DueDate = params.DueDate
	}

	for _, subtaskID := range params.SubtaskIDs {
		subtask, ok := tasklist.Registry[subtaskID]
		if !ok {
			log.Warnf("Could not find subtask (%v)", subtaskID)
			WriteBasicResponse(http.StatusBadRequest, response, request)
			return
		}

		task.AddSubtask(subtask)
	}

	for _, parentID := range params.ParentIDs {
		parent, ok := tasklist.Registry[parentID]
		if !ok {
			log.Warnf("Could not find parent task (%v)", parentID)
			WriteBasicResponse(http.StatusBadRequest, response, request)
			return
		}

		task.AddParent(parent)
	}

	err = storage.SaveTasklist(tasklist)

	if err != nil {
		log.Errorf("Couldn't save tasklist (%v)", err)
		WriteBasicResponse(http.StatusInternalServerError, response, request)
		return
	}

	WriteBasicResponse(http.StatusOK, response, request)
}

func deleteTask(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	taskID := vars["id"]

	tasklist, err := storage.GetTasklist()
	if err != nil {
		log.Errorf("Couldn't get tasklist (%v)", err)
		WriteBasicResponse(http.StatusInternalServerError, response, request)
		return
	}

	task, ok := tasklist.Registry[taskID]
	if !ok {
		WriteBasicResponse(http.StatusNotFound, response, request)
		return
	}

	tasklist.Delete(task)

	err = storage.SaveTasklist(tasklist)
	if err != nil {
		log.Errorf("Couldn't save tasklist (%v)", err)
		WriteBasicResponse(http.StatusInternalServerError, response, request)
		return
	}

	WriteBasicResponse(http.StatusOK, response, request)
}

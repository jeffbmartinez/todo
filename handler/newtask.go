package handler

import (
	"encoding/json"
	"net/http"

	"github.com/jeffbmartinez/log"

	"github.com/jeffbmartinez/todo-persistence/storage"
	"github.com/jeffbmartinez/todo-persistence/task"
)

/*
NewTaskParams is the json struct that gets passed in the request to create
a new Task object.
*/
type NewTaskParams struct {
	Name       string   `json:"name"`
	ParentIDs  []string `json:"parentIDs"`
	DueDate    int64    `json:"dueDate"`
	Categories []string `json:"categories"`
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
		log.Warn("Empty body in new task POST request")
		WriteBasicResponse(http.StatusBadRequest, response, request)
		return
	}

	defer request.Body.Close()
	decoder := json.NewDecoder(request.Body)

	var params NewTaskParams
	err := decoder.Decode(&params)
	if err != nil {
		log.Warn("Couldn't decode params")
		WriteBasicResponse(http.StatusBadRequest, response, request)
		return
	}

	tasklist, err := storage.GetTasklist()
	if err != nil {
		log.Errorf("Could not get tasklist")
		WriteBasicResponse(http.StatusInternalServerError, response, request)
		return
	}

	var parentTasks []*task.Task
	for _, parentID := range params.ParentIDs {
		parentTask, ok := tasklist.Registry[parentID]
		if !ok {
			log.Warnf("Couldn't find parent task (%v)", parentID)
			WriteBasicResponse(http.StatusBadRequest, response, request)
			return
		}

		parentTasks = append(parentTasks, parentTask)
	}

	newTask := tasklist.AddTask(params.Name, parentTasks)
	newTask.DueDate = params.DueDate
	newTask.Categories = params.Categories

	if err := storage.SaveTasklist(tasklist); err != nil {
		log.Errorf("Couldn't save tasklist (%v)", err)
		WriteBasicResponse(http.StatusInternalServerError, response, request)
		return
	}

	WriteJSONResponse(response, newTask, http.StatusOK)
}

package handler

import (
	"net/http"

	"github.com/jeffbmartinez/log"

	"github.com/jeffbmartinez/todo/storage"
)

// Tasks handles requests to the /tasks endpoint.
func Tasks(response http.ResponseWriter, request *http.Request) {
	handler := BasicResponse(http.StatusMethodNotAllowed)

	switch request.Method {
	case "GET":
		handler = getTasks
	}

	handler(response, request)
}

func getTasks(response http.ResponseWriter, request *http.Request) {
	tasklist, err := storage.GetTasklist()
	if err != nil {
		log.Errorf("Couldn't get tasklist (%v)", err)
		WriteBasicResponse(http.StatusInternalServerError, response, request)
		return
	}

	WriteJSONResponse(response, tasklist.RootTasks, http.StatusOK)
}

package task

import (
	"github.com/jeffbmartinez/log"
)

func findTaskInSlice(tasks []*Task, taskID string) int {
	for i, task := range tasks {
		if task.ID == taskID {
			return i
		}
	}

	return -1
}

func deleteFromSliceByIndex(tasks []*Task, index int) []*Task {
	// https://github.com/golang/go/wiki/SliceTricks
	copy(tasks[index:], tasks[index+1:])
	tasks[len(tasks)-1] = nil
	return tasks[:len(tasks)-1]
}

func deleteFromSliceByID(tasks []*Task, taskID string) []*Task {
	taskIndex := findTaskInSlice(tasks, taskID)
	if taskIndex == -1 {
		log.Errorf("Couldn't find task ('%v') in parent subtask list, should never happen", taskID)
		return tasks
	}

	return deleteFromSliceByIndex(tasks, taskIndex)
}

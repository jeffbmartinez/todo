package storage

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"

	"github.com/jeffbmartinez/todo/task"
)

const storageFilename = ".todo.storage"

var lock sync.Mutex

/*
GetTasklist retrieves a tasklist in a threadsafe manner. It's not efficient
but gets the job done for now.
*/
func GetTasklist() (task.Tasklist, error) {
	lock.Lock()
	defer lock.Unlock()

	serializableTasks, err := getTasklist(storageFilename)
	if err != nil {
		return task.NewTasklist(), err
	}

	tasklist := task.NewTasklist()

	for _, serializableTask := range serializableTasks {
		newTask := &task.Task{
			ID:           serializableTask.ID,
			Name:         serializableTask.Name,
			Complete:     serializableTask.Complete,
			CreatedDate:  serializableTask.CreatedDate,
			ModifiedDate: serializableTask.ModifiedDate,
			DueDate:      serializableTask.DueDate,
			Categories:   serializableTask.Categories,
			Parents:      []*task.Task{},
			Subtasks:     []*task.Task{},
		}

		tasklist.Registry[serializableTask.ID] = newTask
	}

	for _, serializableTask := range serializableTasks {
		task := tasklist.Registry[serializableTask.ID]

		for _, parentID := range serializableTask.ParentIDs {
			parentTask := tasklist.Registry[parentID]
			task.Parents = append(task.Parents, parentTask)
		}

		for _, subtaskID := range serializableTask.SubtaskIDs {
			subtask := tasklist.Registry[subtaskID]
			task.Subtasks = append(task.Subtasks, subtask)
		}
	}

	for _, task := range tasklist.Registry {
		if task.IsRootTask() {
			tasklist.RootTasks = append(tasklist.RootTasks, task)
		}
	}

	return tasklist, err
}

/*
SaveTasklist saves a tasklist in a threadsafe manner.
*/
func SaveTasklist(tasklist task.Tasklist) error {
	lock.Lock()
	defer lock.Unlock()

	var serializableTasks []Task
	for _, task := range tasklist.Registry {
		var parentIDs []string
		for _, parent := range task.Parents {
			parentIDs = append(parentIDs, parent.ID)
		}

		var subtaskIDs []string
		for _, subtask := range task.Subtasks {
			subtaskIDs = append(subtaskIDs, subtask.ID)
		}

		newTask := Task{
			ID:           task.ID,
			Name:         task.Name,
			Complete:     task.Complete,
			CreatedDate:  task.CreatedDate,
			ModifiedDate: task.ModifiedDate,
			DueDate:      task.DueDate,
			Categories:   task.Categories,
			ParentIDs:    parentIDs,
			SubtaskIDs:   subtaskIDs,
		}

		serializableTasks = append(serializableTasks, newTask)
	}

	return saveTasklist(serializableTasks, storageFilename)
}

func getTasklist(filename string) ([]Task, error) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var serializableTasks []Task
	err = json.Unmarshal(contents, &serializableTasks)
	if err != nil {
		return nil, err
	}

	return serializableTasks, nil
}

func saveTasklist(tasks []Task, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	err = encoder.Encode(tasks)

	return err
}

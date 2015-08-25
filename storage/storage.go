package storage

import (
	"sync"

	"github.com/jeffbmartinez/todo/task"
)

const storageFilename = ".todo.storage"

var taskset task.Taskset

var lock sync.Mutex

/*
GetTaskset retrieves a taskset in a threadsafe manner.
*/
func GetTaskset() (task.Taskset, error) {
	lock.Lock()
	defer lock.Unlock()

	err := taskset.Restore(storageFilename)
	return taskset, err
}

/*
SaveTaskset saves a taskset in a threadsafe manner.
*/
func SaveTaskset(ts task.Taskset) error {
	lock.Lock()
	defer lock.Unlock()

	return ts.Store(storageFilename)
}

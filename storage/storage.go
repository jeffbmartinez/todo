package storage

import (
	"sync"

	"github.com/jeffbmartinez/todo/task"
)

const storageFilename = ".todo.storage"

var lock sync.Mutex
var tasklist task.Tasklist

func init() {
	tasklist = task.NewTasklist()
}

/*
GetTasklist retrieves a tasklist in a threadsafe manner.
*/
func GetTasklist() (task.Tasklist, error) {
	lock.Lock()
	defer lock.Unlock()

	err := tasklist.Restore(storageFilename)
	return tasklist, err
}

/*
SaveTasklist saves a tasklist in a threadsafe manner.
*/
func SaveTasklist(ts task.Tasklist) error {
	lock.Lock()
	defer lock.Unlock()

	return ts.Store(storageFilename)
}

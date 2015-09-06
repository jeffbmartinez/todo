package storage

import (
	"time"
)

/*
Task is an easier to serialize version of task.Task.
*/
type Task struct {
	ID       string
	Name     string
	Complete bool

	CreatedDate  time.Time
	ModifiedDate time.Time
	DueDate      time.Time

	Categories []string

	ParentIDs  []string
	SubtaskIDs []string
}

package storage

import (
	"time"
)

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

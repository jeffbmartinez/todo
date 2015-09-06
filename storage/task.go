package storage

/*
Task is an easier to serialize version of task.Task.
*/
type Task struct {
	ID       string
	Name     string
	Complete bool

	CreatedDate  int64
	ModifiedDate int64
	DueDate      int64

	Categories []string

	ParentIDs  []string
	SubtaskIDs []string
}

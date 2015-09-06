package task

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

/*
Tasklist provides an organized set of tasks.
*/
type Tasklist struct {
	/*
		Registry allows a quick lookup of any task by ID.

		Tasks have references to their parents and subtasks, but looking up a
		task by ID without this registry would require a potentially expensive
		traversal of the tasks.
	*/
	Registry map[string]*Task

	/*
		RootTasks contains the root tasks in the list. These tasks have no
		parents.
	*/
	RootTasks []*Task `json:"-"`
}

/*
NewTasklist returns an empty set of tasks.
*/
func NewTasklist() Tasklist {
	return Tasklist{
		Registry:  make(map[string]*Task),
		RootTasks: make([]*Task, 0),
	}
}

/*
AddTask creates and adds a new task to the task list.
*/
func (ts *Tasklist) AddTask(name string, parents []*Task) *Task {
	task := NewTask(name, parents)

	if task.IsRootTask() {
		ts.RootTasks = append(ts.RootTasks, task)
	}

	ts.Registry[task.ID] = task

	return task
}

/*
Delete task from tasklist
*/
func (ts *Tasklist) Delete(task *Task) {
	delete(ts.Registry, task.ID)

	if task.IsRootTask() {
		ts.RootTasks = deleteFromSliceByID(ts.RootTasks, task.ID)
	}

	task.Delete()
}

/*
Store serializes the contents of the tasklist to the specified file.
*/
func (ts Tasklist) Store(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	err = encoder.Encode(ts)

	return err
}

/*
Restore restores a serialized tasklist to the tasklist object from
the file specified.
*/
func (ts *Tasklist) Restore(filename string) error {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(contents, &ts)
	if err != nil {
		return err
	}

	ts.RootTasks = []*Task{}

	for _, task := range ts.Registry {
		if task.IsRootTask() {
			ts.RootTasks = append(ts.RootTasks, task)
		}
	}

	return nil
}

package task

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

/*
Taskset contains a map of task IDs to the tasks themselves.
It provides a quick lookup to find a task, given a task ID.
*/
type Taskset struct {
	Tasks map[string]*Task
}

/*
NewTaskset returns an empty set of tasks.
*/
func NewTaskset() Taskset {
	return Taskset{
		Tasks: make(map[string]*Task),
	}
}

/*
Get returns a task by it's ID, and an 'ok' parameter. If ok
is false, no task with that ID was found.
*/
func (ts *Taskset) Get(ID string) (*Task, bool) {
	task, ok := ts.Tasks[ID]
	return task, ok
}

/*
Put adds a new task or replaces an existing task in the set of tasks.
*/
func (ts *Taskset) Put(task *Task) {
	ts.Tasks[task.ID] = task
}

/*
Delete task from taskset
*/
func (ts *Taskset) Delete(task *Task) {
	delete(ts.Tasks, task.ID)
}

/*
Store serializes the contents of the taskset to the specified file.
*/
func (ts Taskset) Store(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return nil
	}

	encoder := json.NewEncoder(file)
	err = encoder.Encode(ts)

	return err
}

/*
Restore restores a serialized taskset to the taskset object from
the file specified.
*/
func (ts *Taskset) Restore(filename string) error {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(contents, &ts)
	if err != nil {
		return err
	}

	return nil
}

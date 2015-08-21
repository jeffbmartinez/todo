package task

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Taskset struct {
	Tasks []*Task
}

func (ts *Taskset) AddTask(task *Task) {
	ts.Tasks = append(ts.Tasks, task)
}

func (ts Taskset) GetSerializable() serializableTaskset {
	newTaskset := serializableTaskset{}

	for _, task := range ts.Tasks {
		newTask := task.getSerializable()
		newTaskset.Tasks = append(newTaskset.Tasks, &newTask)
	}

	return newTaskset
}

func (ts Taskset) Store(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return nil
	}

	encoder := json.NewEncoder(file)
	err = encoder.Encode(ts.GetSerializable())

	return err
}

func (ts *Taskset) Restore(filename string) error {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	var taskset serializableTaskset
	err = json.Unmarshal(contents, &taskset)
	if err != nil {
		return err
	}

	ts.Tasks = make([]*Task, 0)
	for _, task := range taskset.Tasks {
		ts.Tasks = append(ts.Tasks, task.getTask(nil))
	}

	return nil
}

type serializableTaskset struct {
	Tasks []*serializableTask
}

func (ts serializableTaskset) getTaskset() Taskset {
	taskset := Taskset{}

	for _, task := range ts.Tasks {
		taskset.Tasks = append(taskset.Tasks, task.getTask(nil))
	}

	return taskset
}

package task

import (
	"github.com/twinj/uuid"
)

var allTasks Taskset

func init() {
	allTasks = NewTaskset()
}

/*
Task represents a task or a subtask. Tasks can be stand-alone todo items
or they can be broken down into subtasks, which can then have their own
subtasks, and so on. A Task with a ParentID of "" (empty string) is called
a "root" task, which means it has no parent.
*/
type Task struct {
	ID       string
	Name     string
	Complete bool

	ParentID string
	Subtasks map[string]bool
}

func newTask(name string, parentID string) (*Task, error) {
	task := &Task{
		ID:       uuid.NewV4().String(),
		Name:     name,
		Complete: false,
		ParentID: parentID,
		Subtasks: make(map[string]bool),
	}

	err := allTasks.Add(task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

/*
NewTask creates a new subtask with the assigned parent, in an
incomplete state, with no subtasks of it's own. Passing in "" (empty
string) as the argument for parent means creating a root task.
*/
func NewTask(name string, parentID string) (*Task, error) {
	if parentID == "" {
		return newTask(name, parentID)
	}

	parent, ok := allTasks.Get(parentID)
	if !ok {
		return nil, NewNotFoundError(parentID)
	}

	subtask, err := newTask(name, parentID)
	if err != nil {
		return nil, UnableToCreateTaskError{}
	}

	parent.addSubtask(subtask)

	return subtask, nil
}

/*
IsRootTask returns true if the task is a root task. In other words, it is the
main task.
*/
func (t Task) IsRootTask() bool {
	return t.ParentID == ""
}

/*
ParentTask retrieves the Task object associated with the ParentID, or
nil if the task has no parent.
*/
func (t Task) ParentTask() (*Task, error) {
	if t.IsRootTask() {
		return nil, nil
	}

	task, ok := allTasks.Get(t.ParentID)
	if !ok {
		return nil, NewNotFoundError(t.ParentID)
	}

	return task, nil
}

/*
MarkAsComplete marks a task and all of it's subtasks as complete or finished.
It the task itself is a subtask of a parent task and the task was
previously the only remaining task to be completed, the parent will
marked as complete as well. This works all the way up the chain.
If this method is called on a task already marked as complete, nothing
happens.
*/
func (t *Task) MarkAsComplete() error {
	if t.Complete {
		return nil
	}

	t.Complete = true

	if t.IsRootTask() {
		return nil
	}

	for subtaskID := range t.Subtasks {
		subtask, ok := allTasks.Get(subtaskID)
		if !ok {
			return NewNotFoundError(subtaskID)
		}

		if err := subtask.MarkAsComplete(); err != nil {
			return err
		}
	}

	parentTask, err := t.ParentTask()
	if err != nil {
		return err
	}

	subtasksAreComplete, err := parentTask.allSubtasksAreComplete()
	if err != nil {
		return err
	}

	if subtasksAreComplete {
		if err := parentTask.MarkAsComplete(); err != nil {
			return err
		}
	}

	return nil
}

/*
MarkAsIncomplete marks a task, as well as all parents up the chain, as
incomplete. If the task is not marked as complete, nothing happens.
*/
func (t *Task) MarkAsIncomplete() error {
	if !t.Complete {
		return nil
	}

	t.Complete = false

	parent, err := t.ParentTask()
	if err != nil {
		return err
	}

	if err := parent.MarkAsIncomplete(); err != nil {
		return err
	}

	return nil
}

/*
addSubtask adds a subtask to a task. If the subtask is incomplete,
the task will be marked as incomplete as well.
*/
func (t *Task) addSubtask(subtask *Task) error {
	if t.Complete && !subtask.Complete {
		if err := t.MarkAsIncomplete(); err != nil {
			return err
		}
	}

	t.Subtasks[subtask.ID] = true

	return nil
}

/*
allSubtasksAreComplete returns true if all subtasks of the task have been
marked as complete. Also returns true if there are no subtasks.
*/
func (t Task) allSubtasksAreComplete() (bool, error) {
	for subtaskID := range t.Subtasks {
		subtask, ok := allTasks.Get(subtaskID)
		if !ok {
			return false, NewNotFoundError(subtaskID)
		}

		if !subtask.Complete {
			return false, nil
		}
	}

	return true, nil
}

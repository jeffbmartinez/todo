package task

import (
	"time"

	"github.com/twinj/uuid"
)

/*
Registry allows a quick lookup of any task by ID.

Tasks have references to their parents and subtasks, but looking up a
task by ID without this registry would require a potentially expensive
traversal of the tasks.
*/
var Registry map[string]*Task

func init() {
	Registry = make(map[string]*Task)
}

/*
Task represents a task or a subtask. Tasks can be stand-alone todo items
or they can be broken down into subtasks, which can then have their own
subtasks, and so on. Subtasks have parents which they can be grouped into.
Tasks can have multiple parents to cover the possibility of a single task
accomplishing two parent tasks, for example "clean room" can be a subtask
for "clean house" as well as "prepare for parents' visit".

A task with no parents is called a "root" task.
*/
type Task struct {
	ID       string
	Name     string
	Complete bool

	CreatedDate  time.Time
	ModifiedDate time.Time
	DueDate      time.Time

	Categories []string

	Parents  []*Task
	Subtasks []*Task
}

/*
NewTask creates a new task with the assigned parents. The new task is in an
incomplete state, with no subtasks of it's own. Passing in nil or an empty
array as the argument for parent means it has no parent(s) and so it
is a root task.
*/
func NewTask(name string, parents []*Task) *Task {
	if parents == nil {
		parents = make([]*Task, 0)
	}

	now := time.Now()

	newTask := &Task{
		ID:           uuid.NewV4().String(),
		Name:         name,
		Complete:     false,
		CreatedDate:  now,
		ModifiedDate: now,
		DueDate:      time.Time{},
		Categories:   make([]string, 0),
		Parents:      parents,
		Subtasks:     make([]*Task, 0),
	}

	for _, parent := range parents {
		parent.AddSubtask(newTask)
	}

	Registry[newTask.ID] = newTask

	return newTask
}

/*
IsRootTask returns true if the task is a root task. A root task is one
that has no parents.
*/
func (t Task) IsRootTask() bool {
	return len(t.Parents) == 0
}

/*
MarkAsComplete marks a task and all of it's subtasks as complete.
It the task itself is a subtask of a parent task and the task was
previously the only remaining task to be completed, the parent will be
marked as complete as well. This works all the way up the chain.
If this method is called on a task already marked as complete, nothing
happens.
*/
func (t *Task) MarkAsComplete() {
	if t.Complete {
		return
	}

	t.Complete = true

	for _, subtask := range t.Subtasks {
		subtask.MarkAsComplete()
	}

	for _, parent := range t.Parents {
		if parent.allSubtasksAreComplete() {
			parent.MarkAsComplete()
		}
	}
}

/*
MarkAsIncomplete marks a task, as well as all parents up it's chain, as
incomplete. If the task is already incomplete, nothing happens.
*/
func (t *Task) MarkAsIncomplete() {
	if !t.Complete {
		return
	}

	t.Complete = false

	for _, parent := range t.Parents {
		parent.MarkAsIncomplete()
	}
}

/*
AddSubtask adds a subtask to a task. If the subtask is incomplete,
the task will be marked as incomplete as well.
*/
func (t *Task) AddSubtask(subtask *Task) {
	if t.Complete && !subtask.Complete {
		t.MarkAsIncomplete()
	}

	t.Subtasks = append(t.Subtasks, subtask)
}

/*
allSubtasksAreComplete returns true if all subtasks of the task have been
marked as complete. Also returns true if there are no subtasks.
*/
func (t Task) allSubtasksAreComplete() bool {
	for _, subtask := range t.Subtasks {
		if !subtask.Complete {
			return false
		}
	}

	return true
}

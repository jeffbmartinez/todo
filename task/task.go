package task

import (
	"github.com/twinj/uuid"
	"time"
)

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
	ID       string `json:"id"`
	Name     string `json:"name"`
	Complete bool   `json:complete`

	CreatedDate  int64 `json:"createdDate"`
	ModifiedDate int64 `json:"modifiedDate"`
	DueDate      int64 `json:"dueDate"`

	Categories []string `json:"categories"`

	Parents  []*Task `json:"-"`
	Subtasks []*Task `json:"subtasks"`
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

	now := time.Now().Unix()

	newTask := &Task{
		ID:           uuid.NewV4().String(),
		Name:         name,
		Complete:     false,
		CreatedDate:  now,
		ModifiedDate: now,
		DueDate:      0,
		Categories:   make([]string, 0),
		Parents:      parents,
		Subtasks:     make([]*Task, 0),
	}

	for _, parent := range parents {
		parent.AddSubtask(newTask)
	}

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
SetComplete is a convenience method/shortcut for marking a task as complete
or incomplete. It simply delegates the work to either MarkAsComplete()
or MarkAsIncomplete() as appropriate.
*/
func (t *Task) SetComplete(complete bool) {
	if complete {
		t.MarkAsComplete()
	} else {
		t.MarkAsIncomplete()
	}
}

/*
AddSubtask adds a subtask to a task. If the subtask is incomplete,
the task will be marked as incomplete as well. If the provided task
is already a listed subtask, nothing happens.
*/
func (t *Task) AddSubtask(subtask *Task) {
	if findTaskInSlice(t.Subtasks, subtask.ID) != -1 {
		return
	}

	if t.Complete && !subtask.Complete {
		t.MarkAsIncomplete()
	}

	t.Subtasks = append(t.Subtasks, subtask)
}

/*
AddParent adds a parent to a task. If the task is incomplete, it marks
the new parent as incomplete as well. If the parent is already included
in the list of parents, nothing happens.
*/
func (t *Task) AddParent(parent *Task) {
	if findTaskInSlice(t.Parents, parent.ID) != -1 {
		return
	}

	if !t.Complete && parent.Complete {
		parent.MarkAsIncomplete()
	}

	t.Parents = append(t.Parents, parent)
	parent.AddSubtask(t)
}

/*
Delete removes the current Node and all subtasks from the Task structure.
Similar to MarkAsComplete() any parent tasks are then marked as complete
if they have no remaining incomplete subtasks.
*/
func (t *Task) Delete() {
	t.MarkAsComplete()

	for _, subtask := range t.Subtasks {
		subtask.Delete()
	}

	for _, parent := range t.Parents {
		parent.Subtasks = deleteFromSliceByID(parent.Subtasks, t.ID)
	}
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

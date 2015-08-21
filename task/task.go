package task

// Task represents a task or a subtask. Tasks can be stand-alone todo items
// or they can be broken down into subtasks, which can then have their own
// subtasks, and so on.
type Task struct {
	Name     string
	Complete bool

	Parent   *Task
	Subtasks []*Task
}

func newTask(name string, parent *Task) *Task {
	return &Task{
		Name:     name,
		Complete: false,
		Parent:   parent,
		Subtasks: make([]*Task, 0),
	}
}

// NewRootTask creates a new root task in an incomplete state, with no
// subtasks.
func NewRootTask(name string) *Task {
	return newTask(name, nil)
}

// NewSubtask creates a new subtask with the assigned parent, in an
// incomplete state, with no subtasks of it's own. Passing in *nil*
// as the argument for parent is equivalent to using NewRootTask.
// Using NewRootTask for this purpose is recommended as the intention
// is more clear.
func NewSubtask(name string, parent *Task) *Task {
	subtask := newTask(name, parent)

	parent.addSubtask(subtask)

	return subtask
}

// IsRootTask returns true if the task is a root task. In other words, it is the
// main task.
func (t Task) IsRootTask() bool {
	return t.Parent == nil
}

// MarkAsComplete marks a task and all of it's subtasks as complete or finished.
// It the task itself is a subtask of a parent task, and the task was
// previously the only remaining task to be completed, the parent will
// marked as complete as well. This works all the way up the chain.
// If this method is called on a task already marked as complete, nothing
// happens.
func (t *Task) MarkAsComplete() {
	if t.Complete {
		return
	}

	t.Complete = true

	for _, subtask := range t.Subtasks {
		subtask.MarkAsComplete()
	}

	if t.IsRootTask() {
		return
	}

	if t.Parent.allSubtasksAreComplete() {
		t.Parent.MarkAsComplete()
	}
}

// MarkAsIncomplete marks a task, as well as all parents up the chain, as
// incomplete. If the task is not marked as complete, nothing happens.
func (t *Task) MarkAsIncomplete() {
	if !t.Complete {
		return
	}

	t.Complete = false
	t.Parent.MarkAsIncomplete()
}

// addSubtask adds a subtask to a task. If the subtask is incomplete,
// the task will be marked as incomplete as well.
func (t *Task) addSubtask(subtask *Task) {
	if t.Complete && !subtask.Complete {
		t.MarkAsIncomplete()
	}

	t.Subtasks = append(t.Subtasks, subtask)
}

// allSubtasksAreComplete returns true if all subtasks of the task have been
// marked as complete. Also returns true if there are no subtasks.
func (t Task) allSubtasksAreComplete() bool {
	for _, subtask := range t.Subtasks {
		if !subtask.Complete {
			return false
		}
	}

	return true
}

// getSerializable returns a serializable version of a Task. It is almost
// identical, except the Parent is omitted to avoid cycles in serialization
func (t Task) getSerializable() serializableTask {
	newTask := serializableTask{
		Name:     t.Name,
		Complete: t.Complete,
		Subtasks: make([]*serializableTask, 0),
	}

	for _, subtask := range t.Subtasks {
		newSubtask := subtask.getSerializable()
		newTask.Subtasks = append(newTask.Subtasks, &newSubtask)
	}

	return newTask
}

// serializableTask is similar to a Task, but with the circular dependency
// of Parent removed, to make it easier to serialize
type serializableTask struct {
	Name     string
	Complete bool
	Subtasks []*serializableTask
}

func (t serializableTask) getTask(parent *Task) *Task {
	task := newTask(t.Name, parent)
	task.Complete = t.Complete

	for _, subtask := range t.Subtasks {
		task.Subtasks = append(task.Subtasks, subtask.getTask(task))
	}

	return task
}

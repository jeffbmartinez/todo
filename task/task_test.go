package task

import (
	"testing"
)

func TestIsRootTask(t *testing.T) {
	task := NewTask("root", nil)
	if !task.IsRootTask() {
		t.FailNow()
	}
}

func TestMarkAsComplete1(t *testing.T) {
	root := NewTask("root", nil)

	if root.Complete {
		t.FailNow()
	}

	child1 := NewTask("child1", []*Task{root})

	if root.Complete {
		t.FailNow()
	}

	child1.MarkAsComplete()

	if !root.Complete {
		t.FailNow()
	}
}

func TestMarkAsComplete2(t *testing.T) {
	root := NewTask("root", nil)
	child := NewTask("child", []*Task{root})
	grandchild := NewTask("grandchild", []*Task{child})

	if root.Complete || child.Complete || grandchild.Complete {
		t.FailNow()
	}

	grandchild.MarkAsComplete()

	if !(root.Complete && child.Complete && grandchild.Complete) {
		t.FailNow()
	}
}

func TestMarkAsComplete3(t *testing.T) {
	root := NewTask("root", nil)
	child := NewTask("child", []*Task{root})
	grandchild1 := NewTask("grandchild1", []*Task{child})
	grandchild2 := NewTask("grandchild2", []*Task{child})

	if root.Complete {
		t.Fatal("Root task should not be complete by default")
	}

	grandchild1.MarkAsComplete()

	if root.Complete {
		t.Fatal("Root task should not be complete, it has an incomplete subtask")
	}

	grandchild2.MarkAsComplete()

	if !root.Complete {
		t.Fatal("Root task should be complete, all subtasks are complete")
	}
}

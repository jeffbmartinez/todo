package task

import (
	"testing"
)

func TestIsRootTask(t *testing.T) {
	task, err := NewRootTask("root")
	if err != nil || !task.IsRootTask() {
		t.FailNow()
	}
}

func TestMarkAsComplete1(t *testing.T) {
	root, err := NewRootTask("root")

	if err != nil || root.Complete {
		t.FailNow()
	}

	child1, err := NewSubtask("child1", root.ID)

	if err != nil || root.Complete {
		t.FailNow()
	}

	err = child1.MarkAsComplete()

	if err != nil || !root.Complete {
		t.FailNow()
	}
}

func TestMarkAsComplete2(t *testing.T) {
	root, err := NewRootTask("root")
	if err != nil {
		t.FailNow()
	}

	child, err := NewSubtask("child", root.ID)
	if err != nil {
		t.FailNow()
	}

	grandchild, err := NewSubtask("grandchild", child.ID)
	if err != nil {
		t.FailNow()
	}

	if root.Complete || child.Complete || grandchild.Complete {
		t.FailNow()
	}

	grandchild.MarkAsComplete()

	if !(root.Complete && child.Complete && grandchild.Complete) {
		t.FailNow()
	}
}

func TestMarkAsComplete3(t *testing.T) {
	root, err := NewRootTask("root")
	if err != nil {
		t.FailNow()
	}

	child, err := NewSubtask("child", root.ID)
	if err != nil {
		t.FailNow()
	}

	grandchild1, err := NewSubtask("grandchild1", child.ID)
	if err != nil {
		t.FailNow()
	}

	grandchild2, err := NewSubtask("grandchild2", child.ID)
	if err != nil {
		t.FailNow()
	}

	if root.Complete {
		t.FailNow()
	}

	err = grandchild1.MarkAsComplete()

	if err != nil || root.Complete {
		t.FailNow()
	}

	err = grandchild2.MarkAsComplete()

	if err != nil || !root.Complete {
		t.FailNow()
	}
}

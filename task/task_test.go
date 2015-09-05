package task

import (
	"testing"
)

func TestIsRootTask(t *testing.T) {
	task, err := NewTask("root", "")
	if err != nil || !task.IsRootTask() {
		t.FailNow()
	}
}

func TestMarkAsComplete1(t *testing.T) {
	root, err := NewTask("root", "")

	if err != nil || root.Complete {
		t.FailNow()
	}

	child1, err := NewTask("child1", root.ID)

	if err != nil || root.Complete {
		t.FailNow()
	}

	err = child1.MarkAsComplete()

	if err != nil || !root.Complete {
		t.FailNow()
	}
}

func TestMarkAsComplete2(t *testing.T) {
	root, err := NewTask("root", "")
	if err != nil {
		t.FailNow()
	}

	child, err := NewTask("child", root.ID)
	if err != nil {
		t.FailNow()
	}

	grandchild, err := NewTask("grandchild", child.ID)
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
	root, err := NewTask("root", "")
	if err != nil {
		t.FailNow()
	}

	child, err := NewTask("child", root.ID)
	if err != nil {
		t.FailNow()
	}

	grandchild1, err := NewTask("grandchild1", child.ID)
	if err != nil {
		t.FailNow()
	}

	grandchild2, err := NewTask("grandchild2", child.ID)
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

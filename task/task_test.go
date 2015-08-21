package task

import (
	"testing"
)

func TestIsRootTask(t *testing.T) {
	if !NewRootTask("root").IsRootTask() {
		t.FailNow()
	}
}

func TestMarkAsComplete1(t *testing.T) {
	root := NewRootTask("root")

	if root.Complete {
		t.FailNow()
	}

	child1 := NewSubtask("child1", root)

	if root.Complete {
		t.FailNow()
	}

	child1.MarkAsComplete()

	if !root.Complete {
		t.FailNow()
	}
}

func TestMarkAsComplete2(t *testing.T) {
	root := NewRootTask("root")
	child := NewSubtask("child", root)
	grandchild := NewSubtask("grandchild", child)

	if root.Complete || child.Complete || grandchild.Complete {
		t.FailNow()
	}

	grandchild.MarkAsComplete()

	if !(root.Complete && child.Complete && grandchild.Complete) {
		t.FailNow()
	}
}

func TestMarkAsComplete3(t *testing.T) {
	root := NewRootTask("root")
	child := NewSubtask("child", root)
	grandchild1 := NewSubtask("grandchild1", child)
	grandchild2 := NewSubtask("grandchild2", child)

	if root.Complete {
		t.FailNow()
	}

	grandchild1.MarkAsComplete()

	if root.Complete {
		t.FailNow()
	}

	grandchild2.MarkAsComplete()

	if !root.Complete {
		t.FailNow()
	}
}

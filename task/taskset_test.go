package task

import (
	"fmt"
	"os"
	"testing"
)

func TestStoreTaskset1(t *testing.T) {
	task1 := NewRootTask("one")
	task2 := NewRootTask("two")
	task3 := NewRootTask("three")
	NewSubtask("three-one", task3)

	task2.Complete = true

	ts := Taskset{}
	ts.AddTask(task1)
	ts.AddTask(task2)
	ts.AddTask(task3)

	const tempStorageFile = "_tempstoragefile_.json"

	if err := ts.Store(tempStorageFile); err != nil {
		fmt.Printf("Couldn't store taskset (%v)\n", err)
		t.FailNow()
	}

	if err := ts.Restore(tempStorageFile); err != nil {
		fmt.Printf("Couln't restore file (%v)\n", err)
		t.FailNow()
	}

	if len(ts.Tasks) != 3 || ts.Tasks[0].Name != "one" {
		fmt.Printf("Restored tasks not as expected (1)\n")
		t.FailNow()
	}

	if !ts.Tasks[1].Complete || ts.Tasks[2].Subtasks[0].Name != "three-one" {
		fmt.Printf("Restored tasks not as expected (2)\n")
		t.FailNow()
	}

	if err := os.Remove(tempStorageFile); err != nil {
		fmt.Printf("Couldn't delete storage file (%v)\n", err)
		t.FailNow()
	}
}

func TestRestoreTaskset1(t *testing.T) {
	ts := Taskset{}
	if err := ts.Restore("testdata/restore1.json"); err != nil {
		t.FailNow()
	}

	if len(ts.Tasks) != 3 {
		t.FailNow()
	}

	task1 := ts.Tasks[1]
	if task1.Name != "one" || !task1.Complete || len(task1.Subtasks) != 0 {
		t.FailNow()
	}

	task2 := ts.Tasks[2]
	if task2.Name != "two" || task2.Complete {
		t.FailNow()
	}

	task210 := task2.Subtasks[1].Subtasks[0]
	if task210.Name != "two-one-zero" || task210.Complete || len(task210.Subtasks) != 0 {
		t.FailNow()
	}
}

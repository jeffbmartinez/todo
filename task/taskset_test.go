package task

import (
	"fmt"
	"os"
	"testing"
)

func TestStoreTaskset1(t *testing.T) {
	task1, err := NewRootTask("one")
	if err != nil {
		t.FailNow()
	}

	task2, err := NewRootTask("two")
	if err != nil {
		t.FailNow()
	}

	task3, err := NewRootTask("three")
	if err != nil {
		t.FailNow()
	}

	NewSubtask("three-one", task3.ID)

	task2.Complete = true

	ts := NewTaskset()
	ts.Add(task1)
	ts.Add(task2)
	ts.Add(task3)

	const tempStorageFile = "_tempstoragefile_.json"

	if err := ts.Store(tempStorageFile); err != nil {
		fmt.Printf("Couldn't store taskset (%v)\n", err)
		t.FailNow()
	}

	if err := ts.Restore(tempStorageFile); err != nil {
		fmt.Printf("Couln't restore file (%v)\n", err)
		t.FailNow()
	}

	if len(ts.Tasks) != 3 {
		fmt.Printf("Restored tasks not as expected (1)\n")
		t.FailNow()
	}

	task, ok := ts.Get(task1.ID)
	if !ok || task.Name != "one" {
		t.FailNow()
	}

	if err := os.Remove(tempStorageFile); err != nil {
		fmt.Printf("Couldn't delete storage file (%v)\n", err)
		t.FailNow()
	}
}

func TestRestoreTaskset1(t *testing.T) {
	ts := NewTaskset()
	if err := ts.Restore("testdata/restore1.json"); err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	if len(ts.Tasks) != 6 {
		fmt.Println("here zero", len(ts.Tasks))
		t.FailNow()
	}

	task1, ok := ts.Get("598d22ee-dc92-40fc-b4c8-15f94bfd4d4a")
	if !ok || task1.Name != "one" || !task1.Complete || len(task1.Subtasks) != 0 {
		fmt.Println("here one")
		t.FailNow()
	}

	task2, ok := ts.Get("c30f4d91-6db3-4967-93fe-de76ac40a8cd")
	if !ok || task2.Name != "two" || task2.Complete || len(task2.Subtasks) != 2 {
		fmt.Println("here two")
		t.FailNow()
	}
}

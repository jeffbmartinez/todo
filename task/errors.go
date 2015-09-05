package task

import (
	"fmt"
)

/*
UnableToCreateError is an error object used to signify
that there was a problem creating a new task.
*/
type UnableToCreateError struct {
	reason string
}

func NewUnableToCreateError(reason error) UnableToCreateError {
	return UnableToCreateError{
		reason: reason,
	}
}

func (t UnableToCreateError) Error() string {
	return fmt.Sprintf("Unable to create new task (%v)", t.reason)
}

////////////////////////////////

/*
NotFoundError is used to signify that a task which was
expected to exist could not be found or retrieved.
*/
type NotFoundError struct {
	errorMessage string
}

/*
NewNotFoundError uses the supplied taskID argument to create
an error message.
*/
func NewNotFoundError(taskID string) NotFoundError {
	return NotFoundError{
		errorMessage: fmt.Sprintf("Could not find task with id '%v'", taskID),
	}
}

func (t NotFoundError) Error() string {
	return t.errorMessage
}

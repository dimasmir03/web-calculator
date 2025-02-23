package workerpool

import (
	"fmt"
)

type TaskFunc func(interface{}) error

type Task struct {
	Data interface{}
	Err  error
	fn   TaskFunc
}

func NewTask(fn func(interface{}) error, data interface{}) *Task {
	return &Task{
		Data: data,
		fn:   fn,
	}
}

func process(workerID int, task *Task) {
	fmt.Printf("Worker %d processes task %v\n\n", workerID, task.Data)
	task.Err = task.fn(task.Data)
}

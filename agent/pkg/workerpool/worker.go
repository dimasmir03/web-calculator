package workerpool

import (
	"context"
)

// Worker контролирует всю работу
type Worker struct {
	ID       int
	taskChan <-chan *Task
}

// NewWorker возвращает новый экземпляр worker-а
func NewWorker(ID int, taskChan <-chan *Task) *Worker {
	return &Worker{
		ID:       ID,
		taskChan: taskChan,
	}
}

// Start запуск worker
func (wr *Worker) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case task := <-wr.taskChan:
			if task != nil {
				task.Err = task.fn(task.Data)
			}
		}
	}
}

// StartBackground запускает worker-а в фоне
// func (wr *Worker) StartBackground() {
// 	fmt.Printf("Starting worker %d\n", wr.ID)
// 	for {
// 		select {
// 		case task := <-wr.taskChan:
// 			fmt.Println("обрабатываю")
// 			process(wr.ID, task)
// 		case <-wr.quit:
// 			return
// 		}
// 	}
// }

// Stop Остановка quits для воркера
// func (wr *Worker) Stop() {
// 	fmt.Printf("Closing worker %d\n", wr.ID)
// 	go func() {
// 		wr.quit <- true
// 	}()
// }

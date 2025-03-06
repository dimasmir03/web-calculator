package workerpool

import (
	"context"
	"sync"
)

// Pool воркера
type Pool struct {
	Workers  []*Worker
	taskChan chan *Task

	concurrency int
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
}

// NewPool инициализирует новый пул с заданными задачами и
func NewPool(ctx context.Context, concurrency int) *Pool {
	ctx, cancel := context.WithCancel(ctx)
	return &Pool{
		concurrency: concurrency,
		taskChan:    make(chan *Task, 1000),
		ctx:         ctx,
		cancel:      cancel,
	}
}

// AddTask добавляет таски в pool
func (p *Pool) AddTask(fn TaskFunc, data interface{}) {
	p.taskChan <- &Task{
		fn:   fn,
		Data: data,
	}
}

// Run запускает всю работу в Pool и блокирует ее до тех пор,
// пока она не будет закончена.
func (p *Pool) Run() {
	for i := 1; i <= p.concurrency; i++ {
		worker := NewWorker(i, p.taskChan)
		p.Workers = append(p.Workers, worker)
		p.wg.Add(1)
		go func(w *Worker) {
			defer p.wg.Done()
			w.Start(p.ctx)
		}(worker)
		// worker.Start()
	}

}

// RunBackground запускает pool в фоне
// func (p *Pool) RunBackground() {
// 	fmt.Print("Start workers\n")
// 	// go func() {
// 	// 	for {
// 	// 		time.Sleep(10 * time.Second)
// 	// 	}
// 	// }()
//
// 	for i := 1; i <= p.concurrency; i++ {
// 		worker := NewWorker(p.taskchan, i)
// 		p.Workers = append(p.Workers, worker)
// 		go worker.StartBackground()
// 	}
//
// 	fmt.Print("⌛ Waiting for tasks to come in ...\n")
//
// 	p.runBackground = make(chan bool)
// 	fmt.Println("stop3")
// 	<-p.runBackground
// 	fmt.Println("stop4")
// }

// Stop останавливает запущенных в фоне worker-ов
// func (p *Pool) Stop() {
// 	for i := range p.Workers {
// 		p.Workers[i].Stop()
// 	}
// 	go func() {
// 		p.runBackground <- true
// 	}()
// }

func (p *Pool) Stop() {
	p.cancel()
	p.wg.Wait()
	close(p.taskChan)
}

package worker

import "log/slog"

// WorkerPool manages the workers, distributes tasks, and collects results.
type WorkerPoolService struct {
	TaskQueue     chan Tasker
	activeWorkers int
	maxCount      int
}

func New(maxCount int) *WorkerPoolService {
	in := make(chan Tasker)

	return &WorkerPoolService{
		TaskQueue: in,
		// result:      &sync.Map{},
		maxCount: maxCount,
	}
}

func (wp *WorkerPoolService) Start() {

	for i := wp.activeWorkers; i < wp.maxCount; i++ {
		worker := Worker{id: i, taskQueue: wp.TaskQueue}
		worker.Start()
		wp.activeWorkers++
	}
	slog.Info("worker", slog.Any("activeWorkers", wp.activeWorkers))
}

// Worker is a goroutine that processes tasks and sends the results through a channel.
type Worker struct {
	id        int
	taskQueue chan Tasker
}

// This should be
func (w *Worker) Start() {
	go func() {
		for t := range w.taskQueue {
			// slog.Info("worker", slog.Any("id", w.id))

			t.Task()

		}
	}()
}

type Tasker interface {
	Task()
}

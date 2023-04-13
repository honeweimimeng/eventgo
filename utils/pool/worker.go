package pool

type Worker struct {
	pool Pool
	Next func() Task
}

func (w *Worker) doWork(task Task) {
	if task == nil {
		return
	}
	go func() {
		for {
			func() {
				defer func() {
					if r := recover(); r != nil {
						w.pool.ExHandler(task, r.(error))
						return
					}
				}()
				if task != nil {
					task.Run()
				}
			}()
			task = w.Next()
		}
	}()
}

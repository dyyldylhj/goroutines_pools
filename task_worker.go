package goroutines_pool

type TaskWorker struct {
	Worker
}

func GetTaskWorker() *TaskWorker {
	return &TaskWorker{
		GetWorker(),
	}
}

func GetFixSizeTaskWorker(size int) *TaskWorker {
	return &TaskWorker{
		GetFixSizeWorker(size),
	}
}

func (w *TaskWorker) Do(f func()) {
	w.add()
	w.GetPool().Submit(func() {
		defer w.done()
		f()
	})
}

package goroutines_pool

type workerFunc func(interface{})

type TaskFuncWorker struct {
	Worker
	f workerFunc
}

func GetTaskFuncWorker(f workerFunc) *TaskFuncWorker {
	return &TaskFuncWorker{
		f:      f,
		Worker: GetWorker(),
	}
}

func GetFixSizeTaskFuncWorker(f workerFunc, size int) *TaskFuncWorker {
	return &TaskFuncWorker{
		f:      f,
		Worker: GetFixSizeWorker(size),
	}
}

func (w *TaskFuncWorker) Invoke(args interface{}) {
	w.add()
	w.GetPool().Submit(func() {
		defer w.done()
		w.f(args)
	})
}

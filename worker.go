package goroutines_pool

import (
	"github.com/panjf2000/ants"
	"sync"
	"time"
)

type Worker interface {
	done()
	add()
	Wait()
	Cancel()
	WaitWithTimeOut(timeout time.Duration) (isTimeOut bool)
	GetPool() *ants.Pool
}

//不限大小的协程池worker
type WaitGroupWorker struct {
	pool       *ants.Pool
	cancelChan chan struct{} //取消的chan
	doneChan   chan struct{} //完成的chan
	waitGroup  sync.WaitGroup
}

func GetWorker() *WaitGroupWorker {
	return &WaitGroupWorker{
		pool:       GetPool(),
		cancelChan: make(chan struct{}, 1),
		doneChan:   make(chan struct{}, 1),
		waitGroup:  sync.WaitGroup{},
	}
}

func (w *WaitGroupWorker) Cancel() {
	if len(w.cancelChan) < 1 {
		w.cancelChan <- struct{}{}
	}
}

func (w *WaitGroupWorker) Wait() {
	go func() {
		defer func() {
			recover()
			w.doneChan <- struct{}{}

		}()
		w.waitGroup.Wait()

	}()

	select {
	case <-w.doneChan:
	case <-w.cancelChan:
	}
}

func (w *WaitGroupWorker) WaitWithTimeOut(timeout time.Duration) (isTimeOut bool) {
	isTimeOut = false

	time.AfterFunc(timeout, func() {
		isTimeOut = true
		w.Cancel()
	})

	w.Wait()
	return
}

func (w *WaitGroupWorker) done() {
	w.waitGroup.Done()
}

func (w *WaitGroupWorker) GetPool() *ants.Pool {
	return w.pool
}

func (w *WaitGroupWorker) add() {
	w.waitGroup.Add(1)
}

//固定大小的协程池worker
type FixSizeWorker struct {
	pool *ants.Pool

	cancelChan chan struct{} //取消的chan
	doneChan   chan struct{} //完成的chan

	sizeWaitGroup SizedWaitGroup
}

func GetFixSizeWorker(size int) *FixSizeWorker {
	if size < 0 {
		panic("GetFixSizeWorker, must be greater than 0")
	}
	return &FixSizeWorker{
		pool:          GetPool(),
		cancelChan:    make(chan struct{}, 1),
		doneChan:      make(chan struct{}, 1),
		sizeWaitGroup: New(size),
	}
}

func (w *FixSizeWorker) Cancel() {
	if len(w.cancelChan) < 1 {
		w.cancelChan <- struct{}{}
	}
}

func (w *FixSizeWorker) Wait() {
	go func() {
		defer func() {
			recover()
			w.doneChan <- struct{}{}

		}()
		w.sizeWaitGroup.Wait()

	}()

	select {
	case <-w.doneChan:
	case <-w.cancelChan:
	}
}

func (w *FixSizeWorker) WaitWithTimeOut(timeout time.Duration) (isTimeOut bool) {
	isTimeOut = false

	timer := time.AfterFunc(timeout, func() {
		isTimeOut = true
		w.Cancel()
	})
	defer func() {
		timer.Stop()
	}()

	w.Wait()
	return
}

func (w *FixSizeWorker) done() {
	w.sizeWaitGroup.Done()
}

func (w *FixSizeWorker) add() {
	w.sizeWaitGroup.Add()
}

func (w *FixSizeWorker) GetPool() *ants.Pool {
	return w.pool
}

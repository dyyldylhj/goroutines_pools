package goroutines_pool

import (
	"sync"
)

//使用与多协程互斥干某件事，且有优先级
//优先级必须连续，且唯一
type PriorityGroupoperation struct {
	cond             *sync.Cond
	maxPriorityLevel int
	curPriorityLevel int
	isOver           bool
}

func GetPriorityGroupoperation(maxPriorityLevel int) *PriorityGroupoperation {
	return &PriorityGroupoperation{
		cond:             sync.NewCond(&sync.Mutex{}),
		maxPriorityLevel: maxPriorityLevel,
		curPriorityLevel: maxPriorityLevel,
		isOver:           false,
	}
}

func (g *PriorityGroupoperation) SetMaxPriorityLevel(maxPriorityLevel int) {
	g.maxPriorityLevel = maxPriorityLevel
	g.curPriorityLevel = maxPriorityLevel
}

func (g *PriorityGroupoperation) Broadcast() {
	g.cond.Broadcast()
}

func (g *PriorityGroupoperation) Start(priorityLevel int, f func() bool) {
	if g.isOver {
		return
	}
	g.cond.L.Lock()
	defer func() {
		g.cond.L.Unlock()
	}()
Reentry:
	if g.isOver {
		return
	}

	if g.maxPriorityLevel < 1 || g.curPriorityLevel > priorityLevel {
		g.cond.Wait()
		//fmt.Printf("wake up priorityLevel:%d\n", priorityLevel)
		//time.Sleep(1 * time.Second)
		goto Reentry
	}

	runRt := f()
	if runRt {
		g.isOver = true
	} else {
		g.curPriorityLevel--
	}
	g.cond.Broadcast()
	return
}

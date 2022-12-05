package goroutines_pool

import (
	"fmt"
	"github.com/panjf2000/ants"
	"math"
)

var antsPool *ants.Pool

//获取协程池
func GetPool() *ants.Pool {
	if antsPool == nil {
		var err error
		antsPool, err = ants.NewPool(math.MaxInt32)
		if err != nil {
			panic(fmt.Sprintf("创建ants pool异常, error:%v", err))
		}
	}
	return antsPool
}

//释放协程池
func Release() bool {
	if antsPool == nil {
		return false
	}

	antsPool.Release()
	antsPool = nil
	return true
}

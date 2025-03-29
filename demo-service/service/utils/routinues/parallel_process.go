package routinues

import (
	"fmt"
	"runtime/debug"
	"sync"

	logger "github.com/achilsh/go-dev-tools-collection/base-lib/log"
)

type RoutineGroupWrap struct {
	wg sync.WaitGroup
}

func NewRoutineGroupWrap() *RoutineGroupWrap {
	return new(RoutineGroupWrap)
}

func (g *RoutineGroupWrap) Run(fn func()) {
	g.wg.Add(1)
	GoSafe(func() {
		defer g.wg.Done()
		fn()
	})
}
func (g *RoutineGroupWrap) AsyncRun(async bool, fn func()) {
	g.wg.Add(1)
	AsyncRun(async, func() {
		defer g.wg.Done()
		fn()
	})
}

func (g *RoutineGroupWrap) Wait() {
	g.wg.Wait()
}

func Recover(cleanup ...func()) {
	for _, cl := range cleanup {
		cl()
	}
	if p := recover(); p != nil {
		logger.Errorf(fmt.Sprintf("%s\n%s", fmt.Sprint(p), string(debug.Stack())))
	}
}

func RunSafe(fn func()) {
	defer Recover()
	fn()
}

// GoSafe 对外使用，可以捕获异常并打印异常的栈信息
func GoSafe(fn func()) {
	go RunSafe(fn)
}

// AsyncRun 异步运行任务
func AsyncRun(async bool, fn func()) {
	if async == true {
		go RunSafe(fn)
	} else {
		RunSafe(fn)
	}
}

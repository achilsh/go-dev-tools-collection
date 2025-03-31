package routinues

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

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

func WrapperFnWithTimeout(milliSecond int, fn func()) {
	if milliSecond <= 0 {
		milliSecond = 1000
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(milliSecond*int(time.Millisecond)))
	defer cancel() // 确保资源释放
	resultChan := make(chan string, 1)
	//
	go func(ctxIn context.Context, ch chan string) {
		defer close(resultChan)
		for {
			select {
			case <-ctxIn.Done():
				logger.Infof("biz logic timemout.")
				resultChan <- "timeout"
				return
			default:
				fn()
				logger.Debugf("biz logic fn is over.")
				resultChan <- "succ"
				return
			}
		}
	}(ctx, resultChan)

	// 等待结果或超时
	select {
	case result := <-resultChan:
		logger.Infof("Result: %v", result)
	case <-ctx.Done():
		logger.Errorf("Main context timeout!")
	}
}
func (g *RoutineGroupWrap) AsyncRun(async bool, fn func()) {
	g.wg.Add(1)
	AsyncRun(async, func() {
		defer g.wg.Done()
		fn()
	})
}

func (g *RoutineGroupWrap) AsyncTimeoutRun(async bool, milliSecond int, fn func()) {
	g.wg.Add(1)
	AsyncRun(async, func() {
		defer g.wg.Done()
		if async {
			WrapperFnWithTimeout(milliSecond, fn)
		} else {
			fn()
		}
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

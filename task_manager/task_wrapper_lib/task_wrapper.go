package taskwrapperlib

import (
	"fmt"
	"sync"
	"time"
)

type AsyncsTaskWrapper struct {
	// taskid
	id string

	// 等待异步任务结果最大时间
	maxWaitTime time.Duration

	// 任务状态完成通知通道
	statusNotifyCh chan struct{}

	// 任务运行的结果
	taskResult any

	// 返回结果
	err error

	// 任务是否关闭
	isClose bool

	//
	mutex sync.Mutex
}

func (atw *AsyncsTaskWrapper) SetResult(ret any, err error) {
	atw.taskResult, atw.err = ret, err
}
func NewAsyncTAskWrapper(opts ...AsyncTaskOption) *AsyncsTaskWrapper {
	task := &AsyncsTaskWrapper{
		id:             fmt.Sprintf("%v", time.Now().UTC().UnixMicro()),
		maxWaitTime:    5 * time.Second,
		statusNotifyCh: make(chan struct{}, 1),
		err:            nil,
	}
	for _, opt := range opts {
		opt(task)
	}

	return task
}

type AsyncTaskOption func(a *AsyncsTaskWrapper)

func WithID(id string) AsyncTaskOption {
	return func(a *AsyncsTaskWrapper) {
		a.id = id
	}
}

func WithWaitMaxTime(tm time.Duration) AsyncTaskOption {
	return func(a *AsyncsTaskWrapper) {
		a.maxWaitTime = tm
	}
}

func WithTaskResult(r any) AsyncTaskOption {
	return func(a *AsyncsTaskWrapper) {
		a.taskResult = r
	}
}

func WithErr(e error) AsyncTaskOption {
	return func(a *AsyncsTaskWrapper) {
		a.err = e
	}
}

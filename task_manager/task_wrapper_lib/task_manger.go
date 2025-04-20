package taskwrapperlib

import (
	"errors"
	"fmt"
	"sync"
	"time"

	logger "github.com/achilsh/go-dev-tools-collection/base-lib/log"
	singletonimpl "github.com/achilsh/go-dev-tools-collection/singleton_wrapper/singleton_impl"
)

/////////////////////////

type TaskKeyer interface {
	Key() string
	//确定这种task是不能通过key来确定唯一性，有task 列表，而且不同key的是task 是串行处理
	IsAppend() bool
}

// 通过调用 GetAsyncTaskMngInstance() 获取 *AsyncTaskMng 任务管理器的 单例
var GetAsyncTaskMngInstance = singletonimpl.SingletonObjCreateFuncFactory(func(item *AsyncTaskMng) {
	item.taskMap = make(map[TaskKeyer][]*AsyncsTaskWrapper)
})

type AsyncTaskMng struct {
	locker  sync.RWMutex
	taskMap map[TaskKeyer][]*AsyncsTaskWrapper
}

// 等待入参 task 任务执行结果
// 1.添加任务，
// 2.设置定时器等待结果通知；超时则删除任务；否则就接收结果（结果值和错误）
func (atm *AsyncTaskMng) SyncWait(key TaskKeyer, task *AsyncsTaskWrapper) (any, error) {
	if err := atm.addTask(key, task); err != nil { // 添加时必须要保持原子性，因为删除的时依赖于添加的顺序。
		logger.Errorf("add task fail on wait.")
		return nil, err
	}

	tmr := time.NewTicker(task.maxWaitTime)
	defer tmr.Stop()
	//
	for {
		select {
		case <-tmr.C:
			if err := atm.delTask(key, task); err != nil {
				return nil, fmt.Errorf("del task fail as task timeout, key: %v", key.Key())
			}
			return nil, fmt.Errorf("task run timeout, key: %v", key.Key())

		case <-task.statusNotifyCh: //在发送结果通知时，会删除 移除 任务。
			return task.taskResult, nil
		}
	}
}

// 其中 参数 task 是通知的结果，用临时task mock结果，把任务结果 复制给 任务集中找到的任务。
func (atm *AsyncTaskMng) NotifyDone(key TaskKeyer, taskTemp *AsyncsTaskWrapper) error {
	resultTask, err := atm.getTask(key)
	if err != nil {
		return err
	}
	if resultTask.isClose {
		return nil
	}

	resultTask.mutex.Lock()
	defer resultTask.mutex.Unlock()

	if !resultTask.isClose {
		resultTask.taskResult = taskTemp.taskResult
		resultTask.err = taskTemp.err
		close(resultTask.statusNotifyCh)
		resultTask.isClose = true
	}
	if err := atm.delTask(key, resultTask); err != nil {
		logger.Errorf("del task fail err: %v, key: %v", err, key.Key())
	}
	return nil
}

// 获取 key 对应的任务，或者最早的任务（取出任务来通知结果）
func (atm *AsyncTaskMng) getTask(key TaskKeyer) (*AsyncsTaskWrapper, error) {
	atm.locker.RLock()
	defer atm.locker.RUnlock()

	if key.IsAppend() {
		tasks, ok := atm.taskMap[key]
		if !ok {
			return nil, fmt.Errorf("no task for key: %v", key.Key())
		}
		for _, task := range tasks {
			//取最早的那个任务
			return task, nil
		}
		return nil, fmt.Errorf("not find task for key: %v", key.Key())
	}
	//只有一个任务
	tasks, ok := atm.taskMap[key]
	if !ok || len(tasks) <= 0 {
		return nil, fmt.Errorf("not find task for key: %v", key.Key())
	}
	return tasks[0], nil
}

// 把任务从 任务管理器中移除
func (atm *AsyncTaskMng) delTask(key TaskKeyer, task *AsyncsTaskWrapper) error {
	atm.locker.Lock()
	defer atm.locker.Unlock()
	//
	taskList, ok := atm.taskMap[key]
	if !ok {
		return nil
	}

	//
	index := 0
	if key.IsAppend() {
		for _, item := range taskList {
			if task.id == item.id {
				logger.Infof("task should to delete, task_Id: %v", item.id)
				continue
			}
			taskList[index] = item
			index++
		}
		atm.taskMap[key] = atm.taskMap[key][:index]
		//
		return nil
	}

	delete(atm.taskMap, key)
	return nil
}
func (atm *AsyncTaskMng) addTask(key TaskKeyer, task *AsyncsTaskWrapper) error {
	atm.locker.Lock()
	defer atm.locker.Unlock()

	_, ok := atm.taskMap[key]
	if key.IsAppend() {
		if !ok {
			atm.taskMap[key] = make([]*AsyncsTaskWrapper, 0)
		}
		atm.taskMap[key] = append(atm.taskMap[key], task)
		//
	} else {
		if ok {
			logger.Info("exist some record, key: %v", key.Key())
			return errors.New("exist some record.")
		}
		atm.taskMap[key] = []*AsyncsTaskWrapper{task}
	}
	return nil
}

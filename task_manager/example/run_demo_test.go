package example

import (
	"fmt"
	"log"
	"testing"
	"time"

	taskwrapperlib "github.com/achilsh/go-dev-tools-collection/task_manager/task_wrapper_lib"
	"github.com/stretchr/testify/assert"
)

type MyKeyer interface {
	Key() string
}

type myKey struct {
	id   int64
	name string
	addr *int
}

func (m myKey) Key() string {
	return fmt.Sprintf("%v:%v", m.id, m.name)
}

func TestMyKeyInMap(t *testing.T) {
	var xyz = map[MyKeyer]int{
		myKey{
			id:   10,
			name: "abc",
		}: 11,
		myKey{
			id:   101,
			name: "abc",
		}: 1000,
	}

	for k, v := range xyz {
		fmt.Printf("key: %v, value: %v\n", k, v)
	}

	r := xyz[myKey{
		id:   10,
		name: "abc",
	}]
	fmt.Printf("r value: %v\n", r)
}
func TestTaskWrapper(t *testing.T) {
	log.Printf("this is demo.")
}

type BizKey struct {
	xField int
	yField string
}

func NewBizKey(x int, y string) BizKey {
	return BizKey{
		xField: x,
		yField: y,
	}
}

func (k BizKey) Key() string {
	return fmt.Sprintf("%v:%v", k.xField, k.yField)
}

// 确定这种task是不能通过key来确定唯一性，有task 列表，而且不同key的是task 是串行处理
func (k BizKey) IsAppend() bool {
	return false
}

type TasksProcessor struct {
	dataCh  chan any
	isClose chan struct{}
}

func NewTaskProcessor() *TasksProcessor {
	return &TasksProcessor{
		dataCh:  make(chan any, 100),
		isClose: make(chan struct{}, 1),
	}
}

func (Tp *TasksProcessor) Close() {
	Tp.isClose <- struct{}{}
}
func (Tp *TasksProcessor) Send(data any) {
	select {
	case Tp.dataCh <- data:
		log.Printf("send data to process ok.")
	}
}

func (Tp *TasksProcessor) logic(data any) {
	keyData, ok := data.(BizKey)
	if !ok {
		log.Printf("is not BizKey type data process.")
		return
	}
	//
	log.Printf("receive BizKey to process, begin to process and send result to caller, data: %+v", keyData)

	tmpTaskResult := taskwrapperlib.NewAsyncTAskWrapper(taskwrapperlib.WithTaskResult("------------"))
	if err := taskwrapperlib.GetAsyncTaskMngInstance().NotifyDone(keyData, tmpTaskResult); err != nil {
		log.Printf("send task result notify fail, err: %v", err)
		return
	} else {
		log.Printf("send task result notify succ, key: %v", keyData)
	}
	//

}
func (Tp *TasksProcessor) Loop() {

	for {
		select {
		case <-Tp.isClose:
			log.Printf("receive close signal notify.")
			close(Tp.dataCh)
			return
		case data, ok := <-Tp.dataCh:
			if !ok {
				log.Printf("close data chan")
				return
			}
			Tp.logic(data)
		}
	}
}
func TestTaskProducerAndProcess(t *testing.T) {
	//1. 业务要做一件事
	var tpHandle = NewTaskProcessor()
	go func() {
		tpHandle.Loop()
	}()
	var testCases = []struct {
		id  string
		tm  time.Duration
		key BizKey
	}{
		{
			id: func() string {
				return fmt.Sprintf("%v", time.Now().UTC().UnixMicro())
			}(),
			tm: 1 * time.Second,
			key: func() BizKey {
				fy := fmt.Sprintf("demo_test_yField:%v", time.Now().UTC().UnixMilli())
				return NewBizKey(int(time.Now().UTC().UnixMilli()), fy)
			}(),
		},
	}

	for _, tcase := range testCases {
		t.Run("test tasks", func(tt *testing.T) {

			newTask := taskwrapperlib.NewAsyncTAskWrapper(taskwrapperlib.WithID(tcase.id), taskwrapperlib.WithWaitMaxTime(tcase.tm))
			go func() {
				tpHandle.Send(tcase.key)
			}()

			retsult, err := taskwrapperlib.GetAsyncTaskMngInstance().SyncWait(tcase.key, newTask)
			t.Logf("recevie task process result: %v, err: %v", retsult, err)
			assert.Equal(tt, err, nil)
		})
	}
	t.Logf("task case run over.")
	time.Sleep(100 * time.Microsecond)
	tpHandle.Close()
	time.Sleep(500 * time.Millisecond)
}

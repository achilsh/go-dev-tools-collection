package inoutputtasks

import (
	"context"
	"sync"
)

// 封装 map-reduce logic.

type applyInputJob struct {
	input any
	i     int
}
type applyResult struct {
	result any
	err    error
	i      int
}

type MapFunc func(any) (any, error)
type ReduceFunc func([]any) (any, error)
type TasksWorks struct {
	TaskWorkNums int
	//  函数入参是对输入列表每个元素。 出参是输出列表的每个元素
	MapCall    MapFunc
	ReduceCall ReduceFunc
}

func sendApplyInputJobs(inputJobs chan applyInputJob, inputValues []any) {
	for i, input := range inputValues {
		inputJobs <- applyInputJob{
			input: input,
			i:     i,
		}
	}
	close(inputJobs)
}

func getApplyResults(
	ctx context.Context,
	resultsChan chan applyResult,
	inputValues []any,
) ([]any, error) { //nolint:lll

	results := make([]any, len(inputValues))
	for range results {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()

		case r, ok := <-resultsChan:
			if !ok {
				return results, nil
			}

			if r.err != nil {
				return nil, r.err
			}

			results[r.i] = r.result
		}
	}

	return results, nil
}

//	并发模型是： 从输入 chan中 多个协程消费数据，每个 协程调用逻辑处理函数，处理完后把结果写入到 输出 chan 中。
//
// .......................... ======== process =====
//
// ---- input chan[]----->   ======== process =====   >----- output chan [] -------
//
// ........................	======== process =====
// 输入是有序的，输出结果按输入顺序存放，但是内部处理逻辑是并行处理
func (tw *TasksWorks) Run(ctx context.Context, inputValus []any) ([]any, error) {
	if tw.TaskWorkNums <= 0 {
		tw.TaskWorkNums = 10
	}

	inputJobsChan := make(chan applyInputJob, len(inputValus))
	resultChan := make(chan applyResult, len(inputValus))

	var wg sync.WaitGroup
	wg.Add(tw.TaskWorkNums)

	for w := 0; w < tw.TaskWorkNums; w++ {

		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return

				case input, ok := <-inputJobsChan:
					if !ok {
						return
					}

					if tw.MapCall != nil {
						res, err := tw.MapCall(input.input)

						resultChan <- applyResult{
							result: res,
							err:    err,
							i:      input.i,
						}
					}
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		// 等所有的任务都处理完，包括往输出chan 写完数据。最后关闭输出chan.
		close(resultChan)
	}()

	//把所有输入数据 发送到 输入 chan； 发送完关闭 输入 chan.
	sendApplyInputJobs(inputJobsChan, inputValus)

	// 从 输出 chan 接收数据；
	return getApplyResults(ctx, resultChan, inputValus)
}

func (tw *TasksWorks) Reduce(ctx context.Context, inputValus []any) (any, error) {
	if tw.ReduceCall != nil {
		return tw.ReduceCall(inputValus)
	}
	return nil, nil
}

package inoutputtasks

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTaskwroks(t *testing.T) {
	// 入参是：输入列表中的每个元素。
	mFunc := func(a any) (any, error) {
		input, ok := (a).(map[string]any)
		if !ok {
			return nil, fmt.Errorf("input is not map[string]error type")
		}
		return input, nil
	}

	rFunc := func(a []any) (any, error) {
		var ret string
		for _, v := range a {
			input, ok := (v).(map[string]any)
			if !ok {
				continue
			}

			for k1, v1 := range input {
				if len(ret) > 0 {
					ret += "\n"
				}
				ret += fmt.Sprintf("%v: %v", k1, v1)

			}
		}
		return ret, nil

	}

	var tskw *TasksWorks = &TasksWorks{
		TaskWorkNums: 2,
		MapCall:      mFunc,
		ReduceCall:   rFunc,
	}

	ret, err := tskw.Run(context.Background(), []any{
		map[string]any{
			"i":    0,
			"a1":   123,
			"a12":  345,
			"a123": 678,
		},
		map[string]any{
			"i":    1,
			"b1":   "b1",
			"b12":  "bb2",
			"b123": "bbb3",
		},
		map[string]any{
			"i":    2,
			"c1":   "c1",
			"c12":  "cc2",
			"c123": "cc3",
		},
	})

	assert.Nil(t, err)
	for _, v := range ret {
		vv, ok := v.(map[string]any)
		if !ok {
			continue
		}
		for a, b := range vv {
			t.Logf("k: %v, v: %v", a, b)
		}
	}

	rRet, err := tskw.Reduce(context.Background(), ret)
	if err != nil {
		t.Errorf("err for reduce: %v", err)
		return
	}
	sRet, ok := (rRet).(string)
	if !ok {
		return
	}
	t.Logf("reduce ret: %v", sRet)
}

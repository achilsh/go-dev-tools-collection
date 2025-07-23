package agentlibs

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrAgentNoReturn          = errors.New("no actions or finish was returned by the agent")
	ErrNotFinished            = errors.New("agent not finished before max iterations")
	ErrExecutorInputNotString = errors.New("input to executor not string")
)

const _intermediateStepsOutputKey = "intermediateSteps"

// 定义执行器
func NewExecutor(agent Agent) *Executor {
	ret := &Executor{
		AgentField: agent,
	}
	return ret
}

type Executor struct {
	AgentField              Agent
	ReturnIntermediateSteps bool
}

// 调用 action，调用某个具体的任务； 该任务是第三方或者业务提供的能力
func (e *Executor) doAction(ctx context.Context, steps []AgentStep,
	nameToTool map[string]Tool, action AgentAction) ([]AgentStep, error) {

	// 根据 action 获取任务信息
	tool, ok := nameToTool[strings.ToUpper(action.Tool)]
	if !ok {
		return append(steps, AgentStep{
			Action:      action,
			Observation: fmt.Sprintf("%s is not a valid tool, try another one", action.Tool),
		}), nil
	}

	// 调用 任务具体实现
	observation, err := tool.Call(ctx, strings.TrimSuffix(action.ToolInput, "\nObservation:"))
	if err != nil {
		return nil, err
	}

	return append(steps, AgentStep{
		Action:      action,
		Observation: observation,
	}), nil
}

// 每个任务执行
func (e *Executor) doIteration( // nolint
	ctx context.Context,
	steps []AgentStep,
	nameToTool map[string]Tool,
	inputs map[string]string,

) ([]AgentStep, map[string]any, error) {

	actions, finish, err := e.AgentField.Plan(ctx, steps, inputs)
	if err != nil {
		return steps, nil, err
	}

	if len(actions) == 0 && finish == nil {
		return steps, nil, ErrAgentNoReturn
	}
	if finish != nil {
		return steps, e.getReturn(finish, steps), nil
	}

	// 获取要做的任务列表
	for _, action := range actions {
		steps, err = e.doAction(ctx, steps, nameToTool, action)
		if err != nil {
			return steps, nil, err
		}
	}
	return steps, nil, nil
}

func inputsToString(inputValues map[string]any) (map[string]string, error) {
	inputs := make(map[string]string, len(inputValues))
	for key, value := range inputValues {
		valueStr, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("%w: %s", ErrExecutorInputNotString, key)
		}

		inputs[key] = valueStr
	}

	return inputs, nil
}

func (e *Executor) getReturn(finish *AgentFinish, steps []AgentStep) map[string]any {
	if e.ReturnIntermediateSteps {
		finish.ReturnValues[_intermediateStepsOutputKey] = steps
	}

	return finish.ReturnValues
}
func (e *Executor) Call(ctx context.Context, inputValues map[string]any) (map[string]any, error) {
	nameToTool := getNameToTool(e.AgentField.GetTools())

	inputs, err := inputsToString(inputValues)
	if err != nil {
		return nil, err
	}

	steps := make([]AgentStep, 0)
	var finish map[string]any

	// 对所有的任务，根据 大模型选择任务
	steps, finish, err = e.doIteration(ctx, steps, nameToTool, inputs)
	if err != nil || finish != nil {
		return finish, err
	}

	return e.getReturn(
		&AgentFinish{ReturnValues: make(map[string]any)},
		steps,
	), ErrNotFinished
}

func getNameToTool(t []Tool) map[string]Tool {
	if len(t) == 0 {
		return nil
	}

	nameToTool := make(map[string]Tool, len(t))
	for _, tool := range t {
		nameToTool[strings.ToUpper(tool.Name())] = tool
	}

	return nameToTool
}

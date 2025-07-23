package main

import (
	"context"
	"fmt"

	agentlibs "github.com/achilsh/go-dev-tools-collection/agent_demo/agent_libs"
)

type HelloWorldTool struct {
}

func (c *HelloWorldTool) Name() string {
	return "hello-world-tool"
}
func (c *HelloWorldTool) Description() string {
	return "demo hello world tool."
}
func (c *HelloWorldTool) Call(ctx context.Context, input string) (string, error) {
	return fmt.Sprintf("call demo, input: %v", input), nil
}

var _ agentlibs.Tool = &HelloWorldTool{}
var _ agentlibs.Agent = (*FunctionAgent)(nil)

// 下面 agent 定义
type FunctionAgent struct {
	Tools []agentlibs.Tool
}

// 计划获取 action 和结束状态
// Plan decides what action to take or returns the final result of the input.
func (fa *FunctionAgent) Plan(
	ctx context.Context,
	previousSteps []agentlibs.AgentStep,
	inputs map[string]string,
) ([]agentlibs.AgentAction, *agentlibs.AgentFinish, error) {

	// 可以根据输入的工具 call_tool和输入数据，调用 llm 选择哪个 tool 决定采取哪个 tool.
	return nil, nil, nil
}
func (fa *FunctionAgent) GetTools() []agentlibs.Tool {
	return nil
}

func NewFunctionAgent(tools []agentlibs.Tool) *FunctionAgent {
	ret := &FunctionAgent{
		Tools: tools,
	}
	return ret
}

package main

import (
	"context"
	"encoding/json"
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

// Tool is a tool that can be used by the model.
type Tool struct {
	// Type is the type of the tool.
	Type string `json:"type"`
	// Function is the function to call.
	Function *FunctionDefinition `json:"function,omitempty"`
}

// FunctionDefinition is a definition of a function that can be called by the model.
type FunctionDefinition struct {
	// Name is the name of the function.
	Name string `json:"name"`
	// Description is a description of the function.
	Description string `json:"description"`
	// Parameters is a list of parameters for the function.
	Parameters any `json:"parameters,omitempty"`
	// Strict is a flag to indicate if the function should be called strictly. Only used for openai llm structured output.
	Strict bool `json:"strict,omitempty"`
}

func (fa *FunctionAgent) tools() []Tool {
	res := make([]Tool, 0)
	for _, tool := range fa.Tools {
		res = append(res, Tool{
			Type: "function",
			Function: &FunctionDefinition{
				Name:        tool.Name(),
				Description: tool.Description(),
				Parameters: map[string]any{
					"properties": map[string]any{
						"__arg1": map[string]string{"title": "__arg1", "type": "string"},
					},
					"required": []string{"__arg1"},
					"type":     "object",
				},
			},
		})
	}
	return res
}

func callLLM(tools []Tool) {

}

// FunctionCall is the name and arguments of a function call.
type FunctionCall struct {
	// The name of the function to call.
	Name string `json:"name"`
	// The arguments to pass to the function, as a JSON string.
	Arguments string `json:"arguments"`
}

// ToolCall is a call to a tool (as requested by the model) that should be executed.
type ToolCall struct {
	// ID is the unique identifier of the tool call.
	ID string `json:"id"`
	// Type is the type of the tool call. Typically, this would be "function".
	Type string `json:"type"`
	// FunctionCall is the function call to be executed.
	FunctionCall *FunctionCall `json:"function,omitempty"`
}

// ContentChoice is one of the response choices returned by GenerateContent
// calls.
type ContentChoice struct {
	// Content is the textual content of a response
	Content string

	// StopReason is the reason the model stopped generating output.
	StopReason string

	// GenerationInfo is arbitrary information the model adds to the response.
	GenerationInfo map[string]any

	// FuncCall is non-nil when the model asks to invoke a function/tool.
	// If a model invokes more than one function/tool, this field will only
	// contain the first one.
	FuncCall *FunctionCall

	// ToolCalls is a list of tool calls the model asks to invoke.
	ToolCalls []ToolCall

	// This field is only used with the deepseek-reasoner model and represents the reasoning contents of the assistant message before the final answer.
	ReasoningContent string
}

// ContentResponse is the response returned by a GenerateContent call.
// It can potentially return multiple content choices.
type ContentResponse struct {
	Choices []*ContentChoice
}

func processLLMResult(contentResp *ContentResponse) ([]agentlibs.AgentAction, *agentlibs.AgentFinish, error) {
	choice := contentResp.Choices[0]

	// Check for new-style tool calls first
	if len(choice.ToolCalls) > 0 {
		toolCall := choice.ToolCalls[0]
		functionName := toolCall.FunctionCall.Name
		toolInputStr := toolCall.FunctionCall.Arguments
		toolInputMap := make(map[string]any, 0)
		err := json.Unmarshal([]byte(toolInputStr), &toolInputMap)
		if err != nil {
			// If it's not valid JSON, it might be a raw expression for the calculator
			// Try to use it directly as tool input
			return []agentlibs.AgentAction{
				{
					Tool:      functionName,
					ToolInput: toolInputStr,
					Log:       fmt.Sprintf("Invoking: %s with %s\n", functionName, toolInputStr),
					ToolID:    toolCall.ID,
				},
			}, nil, nil
		}

		toolInput := toolInputStr
		if arg1, ok := toolInputMap["__arg1"]; ok {
			toolInputCheck, ok := arg1.(string)
			if ok {
				toolInput = toolInputCheck
			}
		}

		contentMsg := "\n"
		if choice.Content != "" {
			contentMsg = fmt.Sprintf("responded: %s\n", choice.Content)
		}

		return []agentlibs.AgentAction{
			{
				Tool:      functionName,
				ToolInput: toolInput,
				Log:       fmt.Sprintf("Invoking: %s with %s \n %s \n", functionName, toolInputStr, contentMsg),
				ToolID:    toolCall.ID,
			},
		}, nil, nil
	}

	// Check for legacy function call
	if choice.FuncCall != nil {
		functionCall := choice.FuncCall
		functionName := functionCall.Name
		toolInputStr := functionCall.Arguments
		toolInputMap := make(map[string]any, 0)
		err := json.Unmarshal([]byte(toolInputStr), &toolInputMap)
		if err != nil {
			// If it's not valid JSON, it might be a raw expression for the calculator
			// Try to use it directly as tool input
			return []agentlibs.AgentAction{
				{
					Tool:      functionName,
					ToolInput: toolInputStr,
					Log:       fmt.Sprintf("Invoking: %s with %s\n", functionName, toolInputStr),
					ToolID:    "", // Legacy function calls don't have tool IDs
				},
			}, nil, nil
		}

		toolInput := toolInputStr
		if arg1, ok := toolInputMap["__arg1"]; ok {
			toolInputCheck, ok := arg1.(string)
			if ok {
				toolInput = toolInputCheck
			}
		}

		contentMsg := "\n"
		if choice.Content != "" {
			contentMsg = fmt.Sprintf("responded: %s\n", choice.Content)
		}

		return []agentlibs.AgentAction{
			{
				Tool:      functionName,
				ToolInput: toolInput,
				Log:       fmt.Sprintf("Invoking: %s with %s \n %s \n", functionName, toolInputStr, contentMsg),
				ToolID:    "", // Legacy function calls don't have tool IDs
			},
		}, nil, nil
	}

	// No function/tool call - this is a finish
	return nil, &agentlibs.AgentFinish{
		ReturnValues: map[string]any{
			"output": choice.Content,
		},
		Log: choice.Content,
	}, nil
}

// 计划获取 action 和结束状态
// Plan decides what action to take or returns the final result of the input.
func (fa *FunctionAgent) Plan(
	ctx context.Context,
	previousSteps []agentlibs.AgentStep,
	inputs map[string]string,
) ([]agentlibs.AgentAction, *agentlibs.AgentFinish, error) {

	tools := fa.tools()
	// Tools is a list of tools to use. Each tool can be a specific tool or a function.
	// Tools []Tool `json:"tools,omitempty"`
	//  需要把 tools 设置 llm的请求参数Tools字段。

	callLLM(tools)

	// 可以根据输入的工具 call_tool和输入数据，调用 llm 选择哪个 tool 决定采取哪个 tool.
	//此处是调用大模型的 call_tools, 比如openai大模型的 function_tools.

	// 然后解析 llm 返回的数据
	processLLMResult(nil)
	return nil, nil, nil
}
func (fa *FunctionAgent) GetTools() []agentlibs.Tool {
	return fa.Tools
}

func NewFunctionAgent(tools []agentlibs.Tool) *FunctionAgent {
	ret := &FunctionAgent{
		Tools: tools,
	}
	return ret
}

package agentlibs

import "context"

type Agent interface {
	// 根据输入和之前的step；决定下一步做什么。返回动作或结束状态
	Plan(ctx context.Context, previousSteps []AgentStep, inputs map[string]string) ([]AgentAction, *AgentFinish, error)
	GetTools() []Tool
}

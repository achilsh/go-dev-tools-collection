// AgentAction is the agent's action to take.
package agentlibs

// agent动作定义
type AgentAction struct {
	Tool      string
	ToolInput string
	Log       string
	ToolID    string
}

// AgentStep is a step of the agent: agent 步骤
type AgentStep struct {
	Action      AgentAction
	Observation string //结果
}

// AgentFinish is the agent's return value: agent返回结果
type AgentFinish struct {
	ReturnValues map[string]any
	Log          string
}

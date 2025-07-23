package main

import (
	"context"

	agentlibs "github.com/achilsh/go-dev-tools-collection/agent_demo/agent_libs"
)

func main() {
	tools := []agentlibs.Tool{&HelloWorldTool{}}
	agent := NewFunctionAgent(tools)
	executor := agentlibs.NewExecutor(agent)
	executor.Call(context.Background(), map[string]any{})
}

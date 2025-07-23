package agentlibs

import "context"

// tool 是 agent与不同应用交互的工具
type Tool interface {
	Name() string
	Description() string
	Call(ctx context.Context, input string) (string, error)
}

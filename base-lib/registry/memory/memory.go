package memory

import "go-micro.dev/v4/registry"

var (
	Registry = registry.NewMemoryRegistry()
)

// NewRegistry memory注册中心
func NewRegistry() registry.Registry {
	return Registry
}

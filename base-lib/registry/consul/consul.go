package consul

import (
	"github.com/go-micro/plugins/v4/registry/consul"
	"go-micro.dev/v4/registry"
)

const (
	DefaultAddr = "127.0.0.1:8500"
)

// NewRegistry consul注册中心
//
// addr: consul地址
func NewRegistry(addr string) registry.Registry {
	if addr == "" {
		addr = DefaultAddr
	}
	return consul.NewRegistry(registry.Addrs(addr))
}

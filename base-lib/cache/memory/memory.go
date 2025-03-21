package memory

import "go-micro.dev/v4/cache"

var (
	Cache = cache.NewCache()
)

// NewCache memory cache
func NewCache() cache.Cache {
	return Cache
}

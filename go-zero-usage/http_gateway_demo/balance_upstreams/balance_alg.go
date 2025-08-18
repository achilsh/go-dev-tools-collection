package balanceupstreams

import "sync"

type SmoothWeightedNode struct {
	weight  int // 配置的权重
	current int // 动态权重
	addr    string
}

type SmoothWeightedBalancer struct {
	nodes []*SmoothWeightedNode
	mu    sync.Mutex
}

func NewSmoothWeightedBalancer(weights map[string]int) *SmoothWeightedBalancer {
	var nodes []*SmoothWeightedNode
	for addr, w := range weights {
		nodes = append(nodes, &SmoothWeightedNode{
			weight: w,
			addr:   addr,
		})
	}
	return &SmoothWeightedBalancer{nodes: nodes}
}

func (b *SmoothWeightedBalancer) Pick() string {
	b.mu.Lock()
	defer b.mu.Unlock()

	var best *SmoothWeightedNode
	total := 0

	for _, n := range b.nodes {
		n.current += n.weight
		total += n.weight
		if best == nil || n.current > best.current {
			best = n
		}
	}

	if best == nil {
		return ""
	}

	best.current -= total
	return best.addr
}

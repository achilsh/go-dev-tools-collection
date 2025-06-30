package workflowusage

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// chain 的错误信息，
type BaseErrorFlowCtx struct {
	mu              sync.Mutex
	dependentErrMap map[string]map[string]error
}

func NewBaseErrorFlowCtx() *BaseErrorFlowCtx {
	return &BaseErrorFlowCtx{
		dependentErrMap: make(map[string]map[string]error),
	}
}
func (b *BaseErrorFlowCtx) SetDependentErrors(node string, errs map[string]error) {
	if node == "" || len(errs) <= 0 {
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()
	if b.dependentErrMap == nil {
		b.dependentErrMap = make(map[string]map[string]error)
	}
	b.dependentErrMap[node] = errs
}

func (b *BaseErrorFlowCtx) GetDependentErrors(node string) map[string]error {
	if node == "" {
		return nil
	}

	b.mu.Lock()
	defer b.mu.Unlock()
	return b.dependentErrMap[node]
}

// A
// ├── B (依赖 A)
// │   └── C (依赖 B)
// 假设节点 A 执行失败了，B 和 C 都依赖于 A。
// 在 B 执行时，会调用：ctx.GetDependentErrors("B")，你就可以拿到 A 的错误。
// 在 C 执行时，会调用：ctx.GetDependentErrors("C")，你就可以拿到 B 的错误（如果 B 失败了），进而间接感知 A 是否成功。

type WorkflowCtxBasic interface {
	//设置node 下游节点的错误信息
	SetDependentErrors(node string, depErrs map[string]error)
	//获取 node 依赖下游节点的运行错误信息（不是点node自身的错误，是他依赖的），返回下游节点对应的错误信息
	GetDependentErrors(node string) map[string]error
}

type WorkflowNode struct {
	Name        string
	Dependences []string
	Action      func(ctx WorkflowCtxBasic) error
}

func NewWorkflowManager() *WorkflowManager {
	return &WorkflowManager{
		Nodes:  make(map[string]*WorkflowNode),
		Orders: []string{},
	}
}

type WorkflowManager struct {
	Nodes  map[string]*WorkflowNode
	Orders []string
}

func (wfm *WorkflowManager) AddNode(name string, deps []string, proc func(WorkflowCtxBasic) error) {
	wfm.Nodes[name] = &WorkflowNode{
		Name:        name,
		Dependences: deps,
		Action:      proc,
	}
	wfm.Orders = append(wfm.Orders, name)
}

type nodeSignal struct {
	done   chan struct{}
	failed atomic.Bool // 如果失败就标记
	err    error       // 存储 error（可为 nil）
}

// 只要依赖有一个失败，上层节点就不会执行逻辑
func (wfm *WorkflowManager) Run(ctx WorkflowCtxBasic) error {
	var wg sync.WaitGroup
	errs := make(chan error, len(wfm.Nodes))
	signals := make(map[string]*nodeSignal)

	for name := range wfm.Nodes {
		signals[name] = &nodeSignal{
			done: make(chan struct{}),
		}
	}

	runNode := func(node *WorkflowNode) {
		defer wg.Done()
		//wait all dependence done.
		dependentErrors := make(map[string]error) //key: is dependence name. value: error
		for _, dep := range node.Dependences {
			signal := signals[dep]
			<-signal.done // dependence is done.
			if signal.failed.Load() {
				dependentErrors[dep] = signal.err
			}
		}
		if len(dependentErrors) > 0 {
			ctx.SetDependentErrors(node.Name, dependentErrors)
		}

		err := node.Action(ctx)
		if err != nil {
			signals[node.Name].failed.Store(true)
			signals[node.Name].err = err
			errs <- fmt.Errorf("node: %s failed: %w", node.Name, err)
		}
		close(signals[node.Name].done)
	}

	for _, name := range wfm.Orders {
		node := wfm.Nodes[name]
		wg.Add(1)
		go runNode(node)
	}
	wg.Wait()

	close(errs)
	if len(errs) > 0 {
		ret := <-errs
		for range errs {
		}
		return ret
	}
	return nil
}

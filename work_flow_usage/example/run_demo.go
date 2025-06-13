package main

import (
	"fmt"
	"sync"

	workflowusage "github.com/achilsh/go-dev-tools-collection/work_flow_usage"
)

var _ workflowusage.WorkflowCtxBasic = (*DemoWorkflowCtx)(nil)

type DemoWorkflowCtx struct {
	mu              sync.Mutex
	dependentErrMap map[string]map[string]error
	AInput          string
	AOutPut         string
	BInput          string
	BOutput         string
	CInput          string
	COutput         string
}

func (dwfc *DemoWorkflowCtx) SetDependentErrors(node string, errs map[string]error) {
	if node == "" || len(errs) <= 0 {
		return
	}

	dwfc.mu.Lock()
	defer dwfc.mu.Unlock()
	if dwfc.dependentErrMap == nil {
		dwfc.dependentErrMap = make(map[string]map[string]error)
	}
	dwfc.dependentErrMap[node] = errs
}

func (dwfc *DemoWorkflowCtx) GetDependentErrors(node string) map[string]error {
	if node == "" {
		return nil
	}

	dwfc.mu.Lock()
	defer dwfc.mu.Unlock()
	return dwfc.dependentErrMap[node]
}

func demoCallWorkflow(input string) {
	wf := workflowusage.NewWorkflowManager()

	wf.AddNode("A", nil, func(wcb workflowusage.WorkflowCtxBasic) error {
		wc, _ := wcb.(*DemoWorkflowCtx)
		wc.AInput = input
		wc.AOutPut = "A output data: " + input
		wc.BInput = wc.AOutPut
		return nil
	})

	wf.AddNode("B", []string{"A"}, func(wcb workflowusage.WorkflowCtxBasic) error {
		wc, _ := wcb.(*DemoWorkflowCtx)

		depNodesErr := wcb.GetDependentErrors("B")
		if len(depNodesErr) > 0 {
			for dep, err := range depNodesErr {
				fmt.Printf("dep: %v, err: %v\n", dep, err)
			}
			return nil
		}

		wc.BOutput = "B output + " + wc.BInput
		wc.CInput = wc.BOutput

		return fmt.Errorf("mock b error: %v", wc.BOutput)
	})

	wf.AddNode("C", []string{"B"}, func(wcb workflowusage.WorkflowCtxBasic) error {
		wc, _ := wcb.(*DemoWorkflowCtx)

		depNodesErr := wcb.GetDependentErrors("C")
		if len(depNodesErr) > 0 {
			for dep, err := range depNodesErr {
				fmt.Printf("dep: %v, err: %v\n", dep, err)
			}
			return nil
		}

		wc.COutput = "C output + " + wc.CInput
		return nil
	})

	wfCtx := &DemoWorkflowCtx{
		dependentErrMap: make(map[string]map[string]error),
		AInput:          "this is a input.",
	}
	wf.Run(wfCtx)
	//
	fmt.Printf("last output: %v\n", wfCtx.COutput)
}
func main() {
	demoCallWorkflow("------")

}

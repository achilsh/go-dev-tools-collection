package main

import (
	"fmt"

	workflowusage "github.com/achilsh/go-dev-tools-collection/work_flow_usage"
)

var _ workflowusage.WorkflowCtxBasic = (*DemoWorkFlowCtx)(nil)

type DemoWorkFlowCtx struct {
	*workflowusage.BaseErrorFlowCtx
	//
	AInput  string
	AOutPut string
	BInput  string
	BOutput string
	CInput  string
	COutput string
}

func demoCallWorkflow(input string) {
	wf := workflowusage.NewWorkflowManager()

	wf.AddNode("A", nil, func(wcb workflowusage.WorkflowCtxBasic) error {
		wc, _ := wcb.(*DemoWorkFlowCtx)
		wc.AInput = input
		wc.AOutPut = "A output data: " + input
		wc.BInput = wc.AOutPut
		return nil
	})

	wf.AddNode("B", []string{"A"}, func(wcb workflowusage.WorkflowCtxBasic) error {
		wc, _ := wcb.(*DemoWorkFlowCtx)

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
		wc, _ := wcb.(*DemoWorkFlowCtx)

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

	wfCtx := &DemoWorkFlowCtx{
		BaseErrorFlowCtx: workflowusage.NewBaseErrorFlowCtx(),

		AInput: "this is a input.",
	}
	wf.Run(wfCtx)
	//
	fmt.Printf("last output: %v\n", wfCtx.COutput)
}
func main() {
	demoCallWorkflow("------")

}

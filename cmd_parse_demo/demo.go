package main

import (
	"context"

	"github.com/spf13/cobra"

	"cmd_parse_demo/cmd"
)

func main() {
	cobra.CheckErr(cmd.NewCLI().ExecuteContext(context.Background()))
}
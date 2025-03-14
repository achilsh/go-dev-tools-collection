package cmd

import (
	"log"
	"os"
	"runtime"

	"github.com/containerd/console"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func versionHandler(cmd *cobra.Command, _ []string) { 
	log.Println("get version.")
}

func preRunEOnCreateCmd(cmd *cobra.Command, _ []string) error {
	log.Printf("call preRunEOnCreateCmd")
	return nil
}	
func RunHandler(cmd *cobra.Command, args []string) error {
	log.Printf("run cmd args: %v", args)
	k, _ := cmd.Flags().GetString("keepalive")
	v, _ := cmd.Flags().GetBool("verbose")
	log.Printf("k, v: %v, %v", k,v)
	return nil
}
func RunServer(_ *cobra.Command, _ []string) error { 
	log.Printf("run server.....")
	return nil
}


// args 是命令行后面的参数
func CreateHandler(cmd *cobra.Command, args []string) error { 
	log.Printf("create handle args: %v", args)
	
	// 这是获取参数选项。
	f, _ := cmd.Flags().GetString("file")
	log.Printf("get input parameters options: %v", f)
	return nil
}
func NewCLI()  *cobra.Command {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	cobra.EnableCommandSorting = false
	//
	if runtime.GOOS == "windows" && term.IsTerminal(int(os.Stdout.Fd())) {
		console.ConsoleFromFile(os.Stdin) //nolint:errcheck
	}

	rootCmd := &cobra.Command {
		Use: "cmd_parse_demo", //命令字 没有携带参数
		Short: "is xyz root command.",
		SilenceUsage:  true,
		SilenceErrors: true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		Run: func(cmd *cobra.Command, args []string) {
			if version, _ := cmd.Flags().GetBool("version"); version {
				versionHandler(cmd, args)
				return
			}
			cmd.Print(cmd.UsageString())
		},
	}
	//设置参数选项名字， 后续通过  cmd.Flags().GetBool("version")  获取.
	rootCmd.Flags().BoolP("version", "v", false, "Show version information")// 定义 父命令 

	///
	createCmd := &cobra.Command{
		Use:     "create demo", // 命令字 create 和 一个 参数 demo
		Short:   "Create a demo from a file",
		Args:    cobra.ExactArgs(1), // 限制 create 后面的参数个数只有一个。
		PreRunE: preRunEOnCreateCmd,
		RunE:    CreateHandler,

	}

	//设置参数选项名字， 后续通过 cmd.Flags().GetString("file") 获取
	createCmd.Flags().StringP("file", "f", "", "Name of the Modelfile (default \"Modelfile\"")
	createCmd.Flags().StringP("quantize", "q", "", "Quantize model to this level (e.g. q4_0)")
	//

	runCmd := &cobra.Command {
		Use: "run dmeo", //命令字run 加 一个参数 demo
		Short:"run a demo.",
		// 指定 run 后面参数至少有一个
		Args:    cobra.MinimumNArgs(1),
		RunE: RunHandler,
	}

	runCmd.Flags().String("keepalive", "", "Duration to keep a model loaded (e.g. 5m)")
	runCmd.Flags().Bool("verbose", false, "Show timings for response")
	runCmd.Flags().Bool("insecure", false, "Use an insecure registry")
	runCmd.Flags().Bool("nowordwrap", false, "Don't wrap words to the next line automatically")
	runCmd.Flags().String("format", "", "Response format (e.g. json)")


	serveCmd := &cobra.Command{
		Use:     "serve", // 这个命令字
		Aliases: []string{"start"},
		Short:   "Start server",
		Args:    cobra.ExactArgs(0), //没有参数
		RunE:    RunServer,
	}
	// 自定义 命令 help 的处理. 
	serveCmd.SetHelpFunc(func(cmd *cobra.Command, _ []string){
		log.Printf("run serve cmd help.....")
		log.Printf("%v", cmd.UsageTemplate())

	})

	///
	rootCmd.AddCommand(
		createCmd,
		runCmd,
		serveCmd,
	)
	return rootCmd
}
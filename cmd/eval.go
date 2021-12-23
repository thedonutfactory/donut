package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/thedonutfactory/donutbox/repl"
)

// replCmd represents the repl command
var replCmd = &cobra.Command{
	Use:   "eval",
	Short: "Run donut language inline",
	Long:  `Run the donut programming environment in interactive mode.`,
	Run: func(cmd *cobra.Command, args []string) {
		engine := "eval"
		repl.Start(os.Stdin, os.Stdout, &engine)
	},
}

func init() {
	rootCmd.AddCommand(replCmd)
}

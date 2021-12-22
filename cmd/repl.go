/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/thedonutfactory/donutbox/repl"
)

// replCmd represents the repl command
var replCmd = &cobra.Command{
	Use:   "repl",
	Short: "Run donut language inline",
	Long:  `Run the donut programming environment in interactive mode.`,
	Run: func(cmd *cobra.Command, args []string) {
		repl.Start(os.Stdin, os.Stdout)
	},
}

func init() {
	rootCmd.AddCommand(replCmd)
}

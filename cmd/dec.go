/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thedonutfactory/donutbox/global"
	"github.com/thedonutfactory/go-tfhe/io"
)

// decCmd represents the dec command
var decCmd = &cobra.Command{
	Use:   "dec",
	Short: "Decrypt an encrypted ciphertext",
	Long:  `Decrypt an encrypted ciphertext from file, using the default private key.`,
	Run: func(cmd *cobra.Command, args []string) {
		encFile, err := cmd.Flags().GetString("file")
		if err != nil {
			fmt.Printf("Error reading input flag `file`\n %s\n", err)
			return
		}

		ctxt, err := io.ReadCiphertext(encFile)
		if err != nil {
			fmt.Printf("Error reading Ciphertext file\n %s\n", err)
			return
		}
		result := global.PriKey.Decrypt(*ctxt)
		fmt.Println(result)
	},
}

func init() {
	rootCmd.AddCommand(decCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// decCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// decCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	decCmd.Flags().StringP("file", "f", "in.txn", "Name of the encrypted input file")
	decCmd.MarkFlagRequired("file")
}

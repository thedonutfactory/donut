package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// transactionCmd represents the transaction command
var (
	inputs []int32

	transactionCmd = &cobra.Command{
		Aliases: []string{"txn"},
		Use:     "transaction",
		Short:   "Create a transaction call",
		Long: `Create a new transaction to the runtime environment, passing function
	name and arguments`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("transaction called")
			fn, err := cmd.Flags().GetString("function")
			if err != nil {
				fmt.Println(err)
				return
			}
			file, err := cmd.Flags().GetString("file")
			if err != nil {
				fmt.Println(err)
				return
			}

			txn := NewDonutTransaction(inputs, fn)
			fmt.Println(txn)
			txn.write(file)
		},
	}
)

func init() {
	rootCmd.AddCommand(transactionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// transactionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// transactionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	transactionCmd.Flags().Int32SliceVarP(&inputs, "input", "i", []int32{}, "Input Variables")

	transactionCmd.Flags().StringP("function", "f", "main", "Entry function call to transaction")
	// transactionCmd.MarkFlagRequired("function")

	transactionCmd.Flags().StringP("file", "o", "txn.out", "File name for transaction")
	// transactionCmd.MarkFlagRequired("file")
}

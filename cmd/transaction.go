package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thedonutfactory/donut/code"
	"github.com/thedonutfactory/donut/global"
	"github.com/thedonutfactory/donut/object"
)

// transactionCmd represents the transaction command
var (
	inputs []int32

	transactionCmd = &cobra.Command{
		Aliases: []string{"txn"},
		Use:     "transaction foo.cipher",
		Short:   "Create a transaction call",
		Long: `Create a new transaction for the given bytecode, passing function
	name and arguments`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			bc := NewDonutByteCode()
			err := bc.read(args[0])
			if err != nil {
				fmt.Printf("Error reading bytecode:\n %s\n", err)
				return
			}

			fn, err := cmd.Flags().GetInt("func")
			if err != nil {
				fmt.Println(err)
				return
			}
			file, err := cmd.Flags().GetString("file")
			if err != nil {
				fmt.Println(err)
				return
			}

			txn := createTransactionCall(fn, inputs, bc)
			txn.write(file)
			fmt.Println("transaction created")
		},
	}
)

func createTransactionCall(funcIndex int, inputs []int32, bc *DonutBytecode) *DonutTransaction {
	txn := NewDonutTransaction()
	lenConst := len(bc.Bytecode.Constants)
	for _, val := range inputs {
		if global.PriKey == nil {
			fmt.Errorf("key is null")
			return nil
		}
		ctxt := global.PriKey.Encrypt(int8(val))
		txn.Bytecode.Constants = append(txn.Bytecode.Constants, &object.Ciphertext{Value: ctxt})
	}
	txn.Bytecode.Instructions = append(txn.Bytecode.Instructions, code.Make(code.OpGetGlobal, funcIndex)...)
	for i := 0; i < len(inputs); i++ {
		txn.Bytecode.Instructions = append(txn.Bytecode.Instructions, code.Make(code.OpConstant, lenConst+i)...)
	}
	txn.Bytecode.Instructions = append(txn.Bytecode.Instructions, code.Make(code.OpCall, len(inputs))...)
	txn.Bytecode.Instructions = append(txn.Bytecode.Instructions, code.Make(code.OpPop)...)
	return txn
}

func createTransactionCall2(funcIndex int, inputs []int32, bc *DonutBytecode) *DonutTransaction {
	txn := NewDonutTransaction()
	lenConst := len(bc.Bytecode.Constants)
	for _, val := range inputs {
		txn.Bytecode.Constants = append(txn.Bytecode.Constants, &object.Integer{Value: int64(val)})
	}
	txn.Bytecode.Instructions = append(txn.Bytecode.Instructions, code.Make(code.OpGetGlobal, funcIndex)...)
	for i := 0; i < len(inputs); i++ {
		txn.Bytecode.Instructions = append(txn.Bytecode.Instructions, code.Make(code.OpConstant, lenConst+i)...)
	}
	txn.Bytecode.Instructions = append(txn.Bytecode.Instructions, code.Make(code.OpCall, len(inputs))...)
	txn.Bytecode.Instructions = append(txn.Bytecode.Instructions, code.Make(code.OpPop)...)
	return txn
}

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

	transactionCmd.Flags().IntP("func", "n", 0, "Function index in the constant pool")
	// transactionCmd.MarkFlagRequired("function")

	transactionCmd.Flags().StringP("file", "o", "in.txn", "File name for transaction")
	// transactionCmd.MarkFlagRequired("file")
}

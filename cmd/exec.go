package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thedonutfactory/donutbox/compiler"
	"github.com/thedonutfactory/donutbox/global"
	"github.com/thedonutfactory/donutbox/object"
	"github.com/thedonutfactory/donutbox/vm"
	"github.com/thedonutfactory/go-tfhe/io"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec *.txn *.cipher",
	Short: "Execute a transaction on a compiled .cipher file",
	Long: `Execute a transaction function call on compiled .cipher file 
containing compiled homomorphic bytecode.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		txn := &DonutTransaction{Version: 1}
		err := txn.read(args[0])
		if err != nil {
			fmt.Printf("Error reading transaction:\n %s\n", err)
			return
		}
		bc := NewDonutByteCode()
		err = bc.read(args[1])
		if err != nil {
			fmt.Printf("Error reading bytecode:\n %s\n", err)
			return
		}

		globals := make([]object.Object, vm.GlobalsSize)
		symbolTable := compiler.NewSymbolTable()
		for i, v := range object.Builtins {
			symbolTable.DefineBuiltin(i, v.Name)
		}

		// add function call to bytecode
		appendFunctionCall(txn, bc)
		machine := vm.NewWithGlobalsState(bc.Bytecode, globals)
		err = machine.Run()
		if err != nil {
			fmt.Printf("Error executing bytecode:\n %s\n", err)
			return
		}

		lastPopped := machine.LastPoppedStackElement()

		outputFile, err := cmd.Flags().GetString("file")
		if err != nil {
			fmt.Printf("Error reading input flag `file`\n %s\n", err)
			return
		}
		if lastPopped.Type() == object.CiphertextObj {
			//left.(*object.Ciphertext).Value
			ctxt := lastPopped.(*object.Ciphertext).Value
			io.WriteCiphertext(&ctxt, outputFile)

			// test
			result := global.PriKey.Decrypt(ctxt)
			fmt.Println(result)
			// test
		}
		fmt.Println(lastPopped.Inspect())
	},
}

func appendFunctionCall(txn *DonutTransaction, bc *DonutBytecode) {
	bc.Bytecode.Constants = append(bc.Bytecode.Constants, txn.Bytecode.Constants...)
	bc.Bytecode.Instructions = append(bc.Bytecode.Instructions, txn.Bytecode.Instructions...)
}

func init() {
	rootCmd.AddCommand(execCmd)
	execCmd.Flags().StringP("file", "o", "out.txn", "File name for transaction output")

}

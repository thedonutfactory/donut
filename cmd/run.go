/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thedonutfactory/donutbox/code"
	"github.com/thedonutfactory/donutbox/compiler"
	"github.com/thedonutfactory/donutbox/object"
	"github.com/thedonutfactory/donutbox/vm"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run [*.cipher]",
	Short: "Execute a compiled .cipher file",
	Long:  `Execute a compiled .cipher file containing compiled homomorphic bytecode.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		bc := NewDonutByteCode()
		err := bc.read(args[0])
		if err != nil {
			fmt.Printf("Error reading bytecode:\n %s\n", err)
			return
		}

		bytecode := bc.Bytecode

		//constants := []object.Object{}
		globals := make([]object.Object, vm.GlobalsSize)
		symbolTable := compiler.NewSymbolTable()
		for i, v := range object.Builtins {
			symbolTable.DefineBuiltin(i, v.Name)
		}

		// add call to main
		appendCall(bytecode)

		machine := vm.NewWithGlobalsState(bytecode, globals)
		err = machine.Run()
		if err != nil {
			fmt.Printf("Error executing bytecode:\n %s\n", err)
			return
		}

		lastPopped := machine.LastPoppedStackElement()
		fmt.Println(lastPopped.Inspect())
	},
}

func appendCall(bytecode *compiler.Bytecode) {
	// add call to main
	bytecode.Instructions = append(bytecode.Instructions, code.Make(code.OpGetGlobal, 1)...)
	bytecode.Instructions = append(bytecode.Instructions, code.Make(code.OpCall, 0)...)
	bytecode.Instructions = append(bytecode.Instructions, code.Make(code.OpPop)...)
}

func init() {
	rootCmd.AddCommand(runCmd)
	//compileCmd.Flags().StringP("file", "ff", "", "donut box run -f foo.cipher")
	//compileCmd.MarkFlagRequired("file")

}

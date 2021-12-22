/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/thedonutfactory/donutbox/compiler"
	"github.com/thedonutfactory/donutbox/lexer"
	"github.com/thedonutfactory/donutbox/object"
	"github.com/thedonutfactory/donutbox/parser"
)

// compileCmd represents the compile command
var compileCmd = &cobra.Command{
	Use:   "compile [*.donut]",
	Short: "Compile .donut programs into intermediate bytecode",
	Long: `Compile .donut programs into intermediate bytecode to be executed
by fully homomorphic runtime environments. Resultant bytecode files (.cipher)
can be executed by running the "run" command`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		f := args[0]
		b, err := ioutil.ReadFile(args[0]) // just pass the file name
		if err != nil {
			fmt.Print(err)
		}
		contents := string(b)

		// runtime
		constants := []object.Object{}
		//globals := make([]object.Object, vm.GlobalSize)
		symbolTable := compiler.NewSymbolTable()
		for i, v := range object.Builtins {
			symbolTable.DefineBuiltin(i, v.Name)
		}
		l := lexer.New(contents)
		p := parser.New(l)
		program := p.ParseProgram()

		for _, msg := range p.Errors() {
			fmt.Print(msg)
		}

		comp := compiler.NewWithState(symbolTable, constants)
		err = comp.Compile(program)
		if err != nil {
			fmt.Printf("Compilation failed:\n %s\n", err)
			return
		}

		debug, err := cmd.Flags().GetBool("debug")
		if err != nil {
			fmt.Print(err)
		}
		if debug {
			fmt.Println("\n-- Constants --")
			for _, c := range comp.Bytecode().Constants {
				fmt.Println(c.Inspect())
			}
			fmt.Println("\n-- Instructions --")
			fmt.Println(comp.Bytecode().Instructions.String())
		}

		code := comp.Bytecode()
		bc := &DonutBytecode{
			Version:  1,
			Bytecode: code,
		}
		bc.write(f[0:len(f)-6] + ".cipher")
	},
}

func init() {
	rootCmd.AddCommand(compileCmd)

	compileCmd.Flags().BoolP("debug", "d", false, "donutbox compile -d -f foo.cipher")
}
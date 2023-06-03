package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/thedonutfactory/donut/compiler"
	"github.com/thedonutfactory/donut/lexer"
	"github.com/thedonutfactory/donut/object"
	"github.com/thedonutfactory/donut/parser"

	"github.com/thedonutfactory/donut/global"
)

// compileCmd represents the compile command
var compileCmd = &cobra.Command{
	Use:   "compile *.donut",
	Short: "Compile a .donut program into intermediate bytecode",
	Long: `Compile a .donut program into intermediate bytecode to be executed
by fully homomorphic runtime environments. Resultant bytecode files (.cipher)
can be executed with the "run" command`,
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
			fmt.Println(msg)
		}

		comp := compiler.NewWithState(symbolTable, constants)
		err = comp.Compile(program)
		if err != nil {
			fmt.Printf("Compilation failed:\n %s\n", err)
			return
		}

		// Encrypt the Constants Pool
		for i, c := range comp.Bytecode().Constants {
			// fmt.Println(c.Inspect())
			if c.Type() == object.IntegerObj {
				comp.Bytecode().Constants[i], err = encryptObject(c)
				if err != nil {
					fmt.Printf("Error encrypting constant: %s\n", err)
					return
				}
				// fmt.Println(comp.Bytecode().Constants[i].Inspect())
			}
		}

		debug, err := cmd.Flags().GetBool("instr")
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

		bc := NewDonutByteCode()
		bc.Bytecode = comp.Bytecode()
		bc.write(f[0:len(f)-6] + ".cipher")
	},
}

func encryptObject(obj object.Object) (*object.Ciphertext, error) {
	switch {
	default:
		return nil, fmt.Errorf("unknown ciphertext conversion for type : %s", obj.Type())
	//case leftType == object.IntegerObj && rightType == object.IntegerObj:
	//	return vm.executeBinaryIntegerOperation(op, left, right)
	//case leftType == object.StringObj && rightType == object.StringObj:
	//	return vm.executeBinaryStringOperation(op, left, right)
	case obj.Type() == object.IntegerObj:
		// return vm.executeBinaryCiphertextOperation(op, left, right)
		val := obj.(*object.Integer).Value
		return &object.Ciphertext{Value: global.PubKey.CipherBits(int(val), global.NB_BITS)}, nil
	}
}

func init() {
	rootCmd.AddCommand(compileCmd)
	compileCmd.Flags().BoolP("instr", "i", false, "view bytecode instructions")
}

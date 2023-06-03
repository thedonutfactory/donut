package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/thedonutfactory/donut/ast"
	"github.com/thedonutfactory/donut/compiler"
	"github.com/thedonutfactory/donut/evaluator"
	"github.com/thedonutfactory/donut/lexer"
	"github.com/thedonutfactory/donut/object"
	"github.com/thedonutfactory/donut/parser"
	"github.com/thedonutfactory/donut/vm"
)

// Oops ...need I explain?
const Oops = `            __,__
ooooooo                                      oo
o888   888o   ooooooo  ooooooooo    oooooooo8 8888
888     888 888     888 888    888 888ooooooo 8888
888o   o888 888     888 888    888         888 88
  88ooo88     88ooo88   888ooo88   88oooooo88  oo
                       o888
`

// Start - starts REPL, passes stdin to lexer line by line
func Start(in io.Reader, out io.Writer, engine *string) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	constants := []object.Object{}
	globals := make([]object.Object, vm.GlobalsSize)
	symbolTable := compiler.NewSymbolTable()

	for i, v := range object.Builtins {
		symbolTable.DefineBuiltin(i, v.Name)
	}

	for {
		fmt.Printf(">> ")

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		if *engine == "eval" {
			evaluate(program, env, out)
		} else if *engine == "vm" {
			if err := compileAndExecute(symbolTable, constants, program, globals, out); err != nil {
				fmt.Fprintf(out, "Woops! Compilation failed:\n %s\n", err)
				continue
			}
		}
	}
}

func evaluate(program *ast.RootNode, env *object.Environment, out io.Writer) {
	evaluated := evaluator.Eval(program, env)
	if evaluated != nil {
		io.WriteString(out, evaluated.Inspect())
		io.WriteString(out, "\n")
	}
}

func compileAndExecute(
	symbolTable *compiler.SymbolTable,
	constants []object.Object,
	program *ast.RootNode,
	globals []object.Object,
	out io.Writer,
) error {
	comp := compiler.NewWithState(symbolTable, constants)
	err := comp.Compile(program)
	if err != nil {
		return err
	}

	code := comp.Bytecode()
	constants = code.Constants

	machine := vm.NewWithGlobalsState(code, globals)
	err = machine.Run()
	if err != nil {
		fmt.Fprintf(out, "Woops! Executing bytecode failed:\n %s\n", err)
	}

	lastPopped := machine.LastPoppedStackElement()
	io.WriteString(out, lastPopped.Inspect())
	io.WriteString(out, "\n")

	return nil
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, Oops)
	io.WriteString(out, "Oops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

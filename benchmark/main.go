package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/thedonutfactory/donutbox/compiler"
	"github.com/thedonutfactory/donutbox/evaluator"
	"github.com/thedonutfactory/donutbox/lexer"
	"github.com/thedonutfactory/donutbox/object"
	"github.com/thedonutfactory/donutbox/parser"
	"github.com/thedonutfactory/donutbox/vm"
)

var engine = flag.String("engine", "vm", "use 'vm' or 'eval'")
var input = `
let fibonacci = func(x) {
	if (x == 0) {
		0
	} else {
		if (x == 1) {
			return 1;
		} else {
			fibonacci(x - 1) + fibonacci(x - 2);
		}
	}
};
fibonacci(35);
`

func main() {
	flag.Parse()

	var duration time.Duration
	var result object.Object

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if *engine == "vm" {
		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			fmt.Printf("Compiler error: %s", err)
			return
		}

		machine := vm.New(comp.Bytecode())
		start := time.Now()

		err = machine.Run()
		if err != nil {
			fmt.Printf("VM error: %s", err)
			return
		}

		duration = time.Since(start)
		result = machine.LastPoppedStackElement()
	} else {
		env := object.NewEnvironment()
		start := time.Now()
		result = evaluator.Eval(program, env)
		duration = time.Since(start)
	}

	fmt.Printf(
		"engine=%s, result=%s, duration=%s\n",
		*engine,
		result.Inspect(),
		duration,
	)
}

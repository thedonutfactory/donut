# DonutBox 🍩
DonutBox is an easy to use runtime for developing and executing fully homomorphic programs with a new programming language called `donut`.

## Buidl

`go build -o donutbox`

## Eval

Run `donutbox eval` in your terminal to enter into interactive mode and play around with the donut programming language.

## Quickstart

Create a new donut file, `foo.donut` and add the following code:

```js
// a simple function that adds two numbers
let addTwoNumbers = func(x, y) {
    return x + y
}

// how about three
let addThreeNumbers = func(x, y, z) {
    return x + y + z
}
```

Compile the donut source file into intermediate bytecode (`foo.cipher`):

`donutbox compile foo.donut`

Create a transaction to execute against the 0th function, with inputs 234 and 100, outputting to a file called `in.txn`

`donutbox txn foo.cipher -n 0 -i 234,100 -o in.txn`

Execute the transaction against compiled bytecode:

`donutbox exec in.txn foo.cipher`

## References
[Writing an interpreter in Go](https://interpreterbook.com) and [Writing a compiler in Go](https://compilerbook.com) books have been instrumental in providing a bedrock with which to build the donut runtime.

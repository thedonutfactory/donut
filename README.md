# DonutBox 🍩
DonutBox is an easy to use runtime for developing and executing fully homomorphic programs with a new programming language called `donut`.

## Fully Homomorphic Encryption
Donutbox is a cryptosystem that supports arbitrary computation on ciphertexts known as fully homomorphic encryption (FHE). Under the hood it uses [go-tfhe](https://github.com/thedonutfactory/go-tfhe), which is the golang implementation of the TFHE scheme. Donutbox enables the construction of programs for any desirable functionality, which can be run on encrypted inputs to produce an encryption of the result. Since Donutbox never decrypts its inputs, it can be run by an untrusted party without revealing its inputs and internal state. Donutbox has practical implications in the outsourcing of private computations, for instance, in the context of cloud computing.

*Funfact: the `T` in TFHE stands for `Torus`, which [geometrically resembles a donut 🍩](https://mathworld.wolfram.com/Torus.html). Also, I am 🇨🇦Canadian, and we love a good ol' toroidal confectionary with coffee. Ergo, the name of the language and runtime)*

## Buidl

`go build -o donutbox`

## Eval

Run `donutbox eval` in your terminal to enter into interactive mode and play around with the donut programming language.

## Quickstart

Let's put together the workflow: develop an FHE program with Donut, compile it into intermediate encrypted bytecode, build an encrypted transaction, and execute it against the bytecode.

1. Create a new donut file, `foo.donut` and add the following code:

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

2. Compile the donut source file into intermediate bytecode (`foo.cipher`):

`donutbox compile foo.donut`

3. Create a transaction to execute against the 0th function (`addTwoNumbers`), with inputs 234 and 100, outputting to a file called `in.txn`

`donutbox txn foo.cipher -n 0 -i 234,100 -o in.txn`

4. Execute the transaction against compiled bytecode:

`donutbox exec in.txn foo.cipher`

## References
[Writing an interpreter in Go](https://interpreterbook.com) and [Writing a compiler in Go](https://compilerbook.com) books have been instrumental in providing a bedrock with which to build the donut runtime.

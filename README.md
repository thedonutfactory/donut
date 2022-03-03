<img src="img/gopher.png" alt="FHE Gopher" width="200"/>

# Donut Programming Language 🍩
Donut is an easy to use runtime for developing and executing fully homomorphic programs with a new programming language by the same name.

## Fully Homomorphic Encryption and Computation
Donut is a cryptosystem that supports arbitrary computation on ciphertexts known as fully homomorphic encryption (FHE). Under the hood it uses [go-tfhe](https://github.com/thedonutfactory/go-tfhe), which is the golang implementation of the TFHE scheme. Donut enables the construction of programs for any desirable functionality, which can be run on encrypted inputs to produce an encryption of the result. Since Donut never decrypts its inputs, it can be run by an untrusted party without revealing its inputs and internal state. Donut has practical implications in the outsourcing of private computations, for instance, in the context of cloud computing.

*(Funfact: the `T` in TFHE stands for `Torus`, which [geometrically resembles a donut 🍩](https://mathworld.wolfram.com/Torus.html). Also, I am 🇨🇦Canadian, and we love a good ol' toroidal confectionary with coffee 🍩☕. Ergo, the name of the language and runtime)*

## Buidl

`go build -o donut`

## Eval

Run `donut eval` in your terminal to enter into interactive mode and play around with the donut programming language.

## Quickstart

Let's put together the workflow: develop an FHE program with Donut, compile it into intermediate encrypted bytecode, build an encrypted transaction, and execute it against the bytecode.

0. Generate the public and private keys (only have to perform this step once):

`donut -keys`

The provider should only have access to the `public.key` file.

1. Create a new donut file, `foo.donut` and add the following code:

```js
// a simple function that adds two numbers
// when x is less than 10 and subtracts two
// when x is greater than or equal to 10
let calculate = func(x, y) {
    if (x < 10) {
        return x + y
    } else {
        return x - y
    }
}
```

2. Compile the donut source file into intermediate bytecode (`foo.cipher`):

`donut compile foo.donut`

3. Create a transaction to execute against the 0th function (`addTwoNumbers`), with inputs 3 and 2, outputting to a file called `in.txn`

`donut txn foo.cipher -n 0 -i 3,2 -o in.txn`

4. Execute the transaction against compiled bytecode, saving the resulting ciphertext to `out.txn`:

`donut exec in.txn foo.cipher -o out.txn`

5. Decrypt and view the resulting ciphertext:

`donut dec -f out.txn`

## TODO

- [x] Homomorphic branching
- [x] Serializable bytecode and transactions
- [ ] Strongly typed ( u8, u16, i32, etc. )
- [ ] Generalized serialization format for bytecode (via protocol buffers)
- [ ] Addition of OO features for composable contracts ( classes, inheritance, etc )

## References

[CGGI19]: I. Chillotti, N. Gama, M. Georgieva, and M. Izabachène. TFHE: Fast Fully Homomorphic Encryption over the Torus. In Journal of Cryptology, volume 33, pages 34–91 (2020). [PDF](https://eprint.iacr.org/2018/421.pdf)

[CGGI16]: I. Chillotti, N. Gama, M. Georgieva, and M. Izabachène. Faster fully homomorphic encryption: Bootstrapping in less than 0.1 seconds. In Asiacrypt 2016 (Best Paper), pages 3-33. [PDF](https://eprint.iacr.org/2016/870.pdf)

[A Security Site]: Buchanan, W. et al. Security and So Many Things. https://asecuritysite.com/

[Writing an interpreter in Go](https://interpreterbook.com) and [Writing a compiler in Go](https://compilerbook.com) books have been instrumental in providing a bedrock with which to build the donut runtime.

package object

import (
	"fmt"

	"github.com/thedonutfactory/go-tfhe/gates"
)

// Ciphertext type holds the value of the integer as an int64
type Ciphertext struct {
	Value gates.Int
}

// Type returns our Integer's ObjectType
func (i *Ciphertext) Type() ObjectType { return CiphertextObj }

// Inspect returns a string representation of the Integer's Value
func (i *Ciphertext) Inspect() string { return fmt.Sprintf("%v", i.Value) }

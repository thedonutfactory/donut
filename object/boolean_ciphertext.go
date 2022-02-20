package object

import (
	"fmt"

	"github.com/thedonutfactory/go-tfhe/core"
)

// Boolean type holds the value of the boolean as a bool
type BooleanCiphertext struct {
	Value *core.LweSample
}

// Type returns our Boolean's ObjectType (BooleanCiphertext)
func (b *BooleanCiphertext) Type() ObjectType { return BooleanCiphertextObj }

// Inspect returns a string representation of the Boolean's Value
func (b *BooleanCiphertext) Inspect() string { return fmt.Sprintf("%v", b.Value) }

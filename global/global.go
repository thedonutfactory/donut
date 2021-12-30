package global

import (
	"fmt"
	"os"

	"github.com/thedonutfactory/go-tfhe/gates"
	"github.com/thedonutfactory/go-tfhe/io"
)

const NB_BITS = 8

// generate params
var (
	// generate the keys
	Ctx            = gates.Default128bitGateBootstrappingParameters()
	PubKey, PriKey = Keys(Ctx)
	Ops            = &gates.CipheredOperations{Pk: PubKey}
)

/*
func init() {
	Ctx = gates.Default128bitGateBootstrappingParameters()
	PubKey, PriKey = Keys(Ctx)
	Ops = &gates.CipheredOperations{Pk: PubKey}
}
*/

func Keys2(ctx *gates.GateBootstrappingParameterSet) (*gates.PublicKey, *gates.PrivateKey) {
	return ctx.GenerateKeys()

}

func Keys(params *gates.GateBootstrappingParameterSet) (*gates.PublicKey, *gates.PrivateKey) {
	var pubKey *gates.PublicKey
	var privKey *gates.PrivateKey
	if _, err := os.Stat("private.key"); err == nil {
		fmt.Println("------ Reading keys from file ------")
		privKey, _ = io.ReadPrivKey("private.key")
		pubKey, _ = io.ReadPubKey("public.key")

	} else {
		fmt.Errorf("no keys found in current directory, please generate keys first")
	}
	return pubKey, privKey
}

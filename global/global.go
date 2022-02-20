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
	Ctx    *gates.GateBootstrappingParameterSet
	PubKey *gates.PublicKey
	PriKey *gates.PrivateKey
	Ops    *gates.CipheredOperations
)

func init() {
	Ctx = gates.TestGateBootstrappingParameters()
	PubKey, PriKey = Keys(Ctx)
	Ops = &gates.CipheredOperations{Pk: PubKey}
}

func InitDebug() {
	Ctx = gates.TestGateBootstrappingParameters()
	PubKey, PriKey = Ctx.GenerateKeys()
	Ops = &gates.CipheredOperations{Pk: PubKey}
}

func Keys(params *gates.GateBootstrappingParameterSet) (*gates.PublicKey, *gates.PrivateKey) {
	var pubKey *gates.PublicKey
	var privKey *gates.PrivateKey
	if _, err := os.Stat("private.key"); err == nil {
		fmt.Println("------ Reading keys from file ------")
		privKey, _ = io.ReadPrivKey("private.key")
		pubKey, _ = io.ReadPubKey("public.key")
	} else {
		//fmt.Println("error: no keys found in current directory, please generate keys first")
		fmt.Println("------ Key Generation ------")
		// generate the keys
		pubKey, privKey = params.GenerateKeys()
		io.WritePrivKey(privKey, "private.key")
		io.WritePubKey(pubKey, "public.key")
	}
	return pubKey, privKey
}

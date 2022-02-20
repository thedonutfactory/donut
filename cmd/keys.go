package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thedonutfactory/go-tfhe/gates"
	"github.com/thedonutfactory/go-tfhe/io"
)

// keysCmd represents the keys command
var keysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Generate public and private keys",
	Long: `Generate public and private keys to execute fully homomorphic programs. Public keys
can be used by third parties to execute programs securely.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("------ Key Generation ------")
		// generate the keys
		var pubKey *gates.PublicKey
		var privKey *gates.PrivateKey
		debug, err := cmd.PersistentFlags().GetBool("debug")
		if err != nil {
			fmt.Println(err)
			return
		}
		if debug {
			fmt.Println("*** Warning: generating insecure debug keys, do not use in production ***")
			pubKey, privKey = gates.TestGateBootstrappingParameters().GenerateKeys()
		} else {
			pubKey, privKey = gates.Default128bitGateBootstrappingParameters().GenerateKeys()
		}
		pubKeyFile, _ := cmd.Flags().GetString("public")
		prvKeyFile, _ := cmd.Flags().GetString("private")
		io.WritePrivKey(privKey, prvKeyFile)
		io.WritePubKey(pubKey, pubKeyFile)
		fmt.Printf("Generated keys: %s, %s\n", pubKeyFile, prvKeyFile)
	},
}

func init() {
	rootCmd.AddCommand(keysCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	keysCmd.PersistentFlags().StringP("public", "p", "public.key", "Public key filename")
	keysCmd.PersistentFlags().StringP("private", "r", "private.key", "Private key filename")

	keysCmd.PersistentFlags().BoolP("debug", "d", false, "create insecure debug keys")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// keysCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

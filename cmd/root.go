package cmd

import (
	"encoding/gob"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thedonutfactory/donutbox/object"
)

var (
	cfgFile string

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "donutbox",
		Short: "🍩 Donutbox Fully Homomorphic Runtime Environment 🍩",
		Long: `🍩 Donutbox 🍩 is a CLI that empowers the development of fully homomorphic
programs. It allows users to write programs in a new language called Donut,
compile the programs into encrypted intermediate bytecode, and build 
transactions to execute fully encrypted programs.`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) { },
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "$HOME/.donutbox/config.yaml", "config file")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {

	// register gob types
	gob.Register(&object.HashKey{})
	gob.Register(&object.Integer{})
	gob.Register(&object.Boolean{})
	gob.Register(&object.Null{})
	gob.Register(&object.ReturnValue{})
	gob.Register(&object.Error{})
	gob.Register(&object.Function{})
	gob.Register(&object.String{})
	gob.Register(&object.Builtin{})
	gob.Register(&object.Array{})
	gob.Register(&object.HashPair{})
	gob.Register(&object.Hash{})
	gob.Register(&object.CompiledFunction{})
	gob.Register(&object.Closure{})

	viper.AddConfigPath("$HOME/.donutbox")
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home + "/.donutbox")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")

	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

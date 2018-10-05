package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var webhook string

func init() {
	// Read Config File
	cobra.OnInitialize(initConfig)

	// Persistent Global Flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.disco)")
	rootCmd.PersistentFlags().StringVar(&webhook, "webhook", "", "webhook url to notify (required)")

	// Bind persistent flags to Viper
	viper.BindPFlag("webhook", rootCmd.PersistentFlags().Lookup("webhook"))

	// Bind ENV Vars to Viper
	viper.BindEnv("webhook") // WEBHOOK='...'
}

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigName(".disco")
	}

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config:", viper.ConfigFileUsed())
	}
}

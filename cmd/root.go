package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "container-juggler",
	Short: "container-juggler manages different environment-scenarios",
	Long: `container-juggler manages different environment-scenarios for docker-compose. 
it generates the 'docker-compose.yml'-file based on specified scenario. 
	
for missing-services in the specied scenario compared to 'all'-scenario, it adds 'extra-hosts'-entries to all other services. so your services running in docker will try to connect the missing services on your host-machine`,
}

// Execute runs cobra-command-tool
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./container-juggler.yml)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("container-juggler")
		viper.AddConfigPath(".")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println("could not read config-file, aborting ...", err)
		os.Exit(1)
	}
}

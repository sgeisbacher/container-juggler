package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "lists all available environments",
	Run: func(cmd *cobra.Command, args []string) {
		scenarios := viper.GetStringMapStringSlice("scenarios")
		for key := range scenarios {
			fmt.Println(key)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

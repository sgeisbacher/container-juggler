package cmd

import (
	"fmt"
	"log"

	"github.com/sgeisbacher/compose-env-manager/generation"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate [scenario]",
	Short: "generates docker-compose.yml for specified scenario",
	Long: `TODO: longer description

[scenario] defaults to all`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("running generate command ...")
		generator := generation.CreateGenerator()
		scenario := ""
		if len(args) > 0 {
			scenario = args[0]
		}
		if err := generator.Generate(scenario); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(generateCmd)
}

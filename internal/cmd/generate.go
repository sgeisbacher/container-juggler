package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/sgeisbacher/container-juggler/internal/generation"
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
		composeFile, err := os.Create("docker-compose.yml")
		if err != nil {
			log.Fatal("could not create docker-compose.yml")
		}
		if err := generator.Generate(scenario, composeFile); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}

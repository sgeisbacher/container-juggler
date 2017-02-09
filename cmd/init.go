package cmd

import (
	"log"

	"github.com/sgeisbacher/compose-env-manager/volumeadmin"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initializes all volumes (based on 'volume-init'-config)",
	Long:  `'init' downloads the zips and extracts them based on 'volume-init'-configuration in compose-env-manager.yml`,
	Run: func(cmd *cobra.Command, args []string) {
		loader := volumeadmin.New()
		if err := loader.Load(false); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}

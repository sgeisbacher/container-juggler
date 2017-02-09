package cmd

import (
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "runs 'docker-compose up'",
	Run: func(cmd *cobra.Command, args []string) {
		path, err := exec.LookPath("docker-compose")
		if err != nil {
			log.Fatal(err)
		}
		dockerComposeCmd := exec.Command(path, "up")
		dockerComposeCmd.Stdout = os.Stdout
		dockerComposeCmd.Stderr = os.Stderr
		dockerComposeCmd.Run()
		dockerComposeCmd.Wait()
	},
}

func init() {
	RootCmd.AddCommand(runCmd)
}

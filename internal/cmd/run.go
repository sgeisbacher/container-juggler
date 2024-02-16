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
		err = dockerComposeCmd.Run()
		if err != nil {
			log.Fatal(err)
		}
		err = dockerComposeCmd.Wait()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

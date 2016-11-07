// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

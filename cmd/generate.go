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
	"strings"

	"github.com/spf13/cobra"
)

var toogle bool

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		/*		dict := viper.GetStringMap("scenarios")
				for key, _ := range dict {
					fmt.Println(key, "- ...")
					for _, lelem := range viper.GetStringSlice("scenarios." + key) {
						fmt.Println("   ", lelem)
					}
				}
				fmt.Println("")
				fmt.Println("")
				fmt.Println("")
				dict2 := viper.GetStringMap("*")
				fmt.Println(dict2) */
		fmt.Println("gen-cmd-args:", strings.Join(args, " "))
		if toogle {
			fmt.Println("gen-cmd toogled:")
		}
	},
}

func init() {
	RootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	generateCmd.Flags().BoolVarP(&toogle, "toggle", "t", false, "Help message for toggle")

}

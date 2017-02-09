// Copyright © 2016 Stefan Geisbacher
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

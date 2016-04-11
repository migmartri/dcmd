// Copyright © 2016 Miguel Martinez <migmartri@gmail.com>
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
	"github.com/spf13/cobra"
	"os"
)

var cfgFile string

const composeFileName = "docker-compose.yml"

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "dcmd",
	Short: "Command line tool for development environments based on Docker Compose",
	Long: `Wraps common operations executed while you are developing in a Docker Compose based environment.
i.e, "dcmd exec rails c" instead of "docker exec -it <container name> rails c`,
}

func checkDockerComposeFile() {
	if _, err := os.Stat(composeFileName); err != nil {
		fmt.Printf("\"%v\" file not found in current directory\n", composeFileName)
		os.Exit(0)
	}
}

// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	checkDockerComposeFile()
}

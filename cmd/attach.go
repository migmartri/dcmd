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
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// attachCmd represents the attach command
var attachCmd = &cobra.Command{
	Use:   "attach",
	Short: "Attach to a running container",
	Run:   attach,
}

func attach(cmd *cobra.Command, args []string) {
	container := chooseContainer()
	dockerAttach(container)
}

func dockerAttach(containerName string) {
	fmt.Printf("Attached to \"%s\"\n", containerName)

	cmd := exec.Command(`docker`, `attach`, containerName)

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		log.Fatal(`Error attaching to container`)
	}
}

func init() {
	RootCmd.AddCommand(attachCmd)
}

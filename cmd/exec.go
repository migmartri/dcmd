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
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec <command to run>",
	Short: "Run any command inside a running container",
	Long: `Execute commands inside one of the containers running in your docker-compose defined cluster.
i/e "dcmd exec rails c" will translate in "docker exec -it <your container> rails c
`,
	Run: invokeCommand,
}

func invokeCommand(cmd *cobra.Command, args []string) {
	/** Check number of args and show usage information if needed */
	if len(args) == 0 {
		cmd.Usage()
		os.Exit(0)
	}
	container := chooseContainer()
	dockerExec(container, args)
}

func dockerExec(containerName string, command []string) {
	res := fmt.Sprintf("docker exec -it %v %s", containerName, strings.Join(command, " "))
	fmt.Printf("Executing: \"%s\"\n", res)

	/* exec uses first argument program name and then a list of options */
	expandedCommand := strings.Split(res, " ")
	cmd := exec.Command(expandedCommand[0], expandedCommand[1:]...)

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		fmt.Println(`Error executing the command`)
		os.Exit(0)
	}
}

/*
  TODO: Mark default option based on your previously choice, supporting clicking enter
*/
func chooseContainer() string {
	names := getContainerNames()
	var choice int
	if len(names) == 0 {
		fmt.Println("You do not have any container available, please run docker-compose up first.")
		os.Exit(0)
	}

	fmt.Println(`Select the container you want to run the command in:`)
	for i, elem := range names {
		fmt.Printf("  %d. %s\n", i+1, elem)
	}

	if _, err := fmt.Scanf("%d", &choice); err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(0)
	}

	// Choice out of bounds
	if len(names) < choice {
		fmt.Println(`Invalid choice.`)
		os.Exit(0)
	}

	/* Clear the terminal */
	cmd := exec.Command(`clear`)
	cmd.Stdout = os.Stdout
	cmd.Run()

	return names[choice-1]
}

func getContainerNames() []string {
	out, err := exec.Command("docker-compose", "ps").Output()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// Get first lines of the output on docker-compose
	r := regexp.MustCompile(`(?m)^[a-z\d_]+`)
	return r.FindAllString(string(out), -1)
}

func init() {
	RootCmd.AddCommand(execCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// execCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// execCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

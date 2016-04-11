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
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
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

/*
  TODO: Mark default option based on your previously choice, supporting clicking enter
*/
func chooseContainer() string {
	names := getContainerNames()
	var containerName string

	var choice int
	if len(names) == 0 {
		fmt.Println("You do not have any container available, please run docker-compose up first.")
		os.Exit(0)
	}

	fmt.Println(`Select the container you want to run the command in:`)
	for i, elem := range names {
		fmt.Printf("  %d. %s", i+1, elem)
		/* Set the value to the one that is in the configuration */
		if nameInConfig := viper.GetString("container"); elem == nameInConfig {
			fmt.Printf(" (default)")
			containerName = elem
		}
		fmt.Println()
	}

	if _, err := fmt.Scanf("%d", &choice); err != nil {
		if err.Error() == "unexpected newline" {
			/* We pressed enter but we have an option in the config file */
			if containerName != "" {
				return containerName
			}
		}
		fmt.Println("Invalid choice")
		os.Exit(0)
	}

	// Choice out of bounds
	if len(names) < choice {
		fmt.Println(`Invalid choice.`)
		os.Exit(0)
	}

	/* Clear the terminal */
	// cmd := exec.Command(`clear`)
	// cmd.Stdout = os.Stdout
	// cmd.Run()

	fmt.Printf("%s", containerName)
	containerName = names[choice-1]
	writeInConfigFile(fmt.Sprintf("container=\"%s\"", containerName))
	return containerName
}

func writeInConfigFile(keyValue string) {
	configFile := fmt.Sprintf("%s.toml", configFile)
	_, err := os.Stat(configFile)

	if err != nil {
		if os.IsNotExist(err) {
			newFile, err := os.Create(configFile)
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
			newFile.Close()
		}
	}
	err = ioutil.WriteFile(configFile, []byte(keyValue), 0666)
	if err != nil {
		log.Fatal(err)
	}
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

const configFile = ".dcmd"

func initConfig() {
	viper.SetConfigName(configFile) // name of config file (without extension)
	viper.AddConfigPath(".")        // adding home directory as first search path
	// If a config file is found, read it in.
	viper.ReadInConfig()
}

func init() {
	cobra.OnInitialize(checkDockerComposeFile)
	cobra.OnInitialize(initConfig)

	// RootCmd.PersistentFlags().String("container", "", "docker container name/id (optional)")
	// viper.BindPFlag("container", RootCmd.PersistentFlags().Lookup("container"))
}

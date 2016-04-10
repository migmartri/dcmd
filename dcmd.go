package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Commands struct {
	Commands []Command
}

type Command struct {
	Name    string
	Command string
	Service string
}

func check(e error) {
	if e != nil {
		log.Fatalf("error: %v", e)
		panic(e)
	}
}

func (commands Commands) availableCommandNames() []string {
	var res []string
	for _, elem := range commands.Commands {
		res = append(res, elem.Name)
	}
	return res
}

/* Parses program name and args */
func (c Command) dockerCommand() string {
	res := fmt.Sprintf("docker exec -it %v %v", c.Service, c.Command)
	return res
}

func (c Command) invoke() {
	fmt.Printf("Executing \"%v\" in container %v: %v\n", c.Command, c.Service, c.dockerCommand())
	expandedCommand := strings.Split(c.dockerCommand(), " ")
	cmd := exec.Command(expandedCommand[0], expandedCommand[1:]...)

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	check(err)
}

func (commands Commands) invoke(name string) {
	var toInvoke Command
	for _, elem := range commands.Commands {
		if elem.Name == name {
			toInvoke = elem
		}
	}

	if toInvoke == (Command{}) {
		fmt.Printf("Command: \"%v\" not found. Available options: %v \n", name, strings.Join(commands.availableCommandNames(), ", "))
		os.Exit(0)
	}

	toInvoke.invoke()
}

func checkFileExists(fileName string) {
	if _, err := os.Stat(fileName); err != nil {
		fmt.Printf("\"%v\" not found in current directory\n", fileName)
		os.Exit(0)
	}
}

func checkArgs() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: \"dcmd <command>\"")
		os.Exit(0)
	}
}

const composeFileName = "docker-compose.yml"
const dcmdFileName = "dcmd.yml"

func loadDcmdYamlFile() Commands {
	data, err := ioutil.ReadFile(dcmdFileName)
	check(err)

	var commands = Commands{}
	err = yaml.Unmarshal([]byte(data), &commands)
	check(err)
	return commands
}

func main() {
	checkFileExists(composeFileName)
	checkFileExists(dcmdFileName)
	checkArgs()

	/* load commandName introduced via CLI args */
	commandName := os.Args[1]
	/* Load Yaml file */
	commands := loadDcmdYamlFile()
	/* Invoke dcmd command */
	commands.invoke(commandName)
}

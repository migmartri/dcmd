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

func (commands Commands) availableCommands() []string {
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
		fmt.Printf("Command: \"%v\" not found. Available options: %v \n", name, strings.Join(commands.availableCommands(), ", "))
		os.Exit(0)
	}

	toInvoke.invoke()
}

func main() {
	data, err := ioutil.ReadFile("dcmd.yaml")
	check(err)

	args := os.Args[1:]
	commandName := args[0]

	var commands = Commands{}
	err = yaml.Unmarshal([]byte(data), &commands)
	check(err)

	commands.invoke(commandName)
}

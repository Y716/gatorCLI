package main

import (
	"fmt"
	"os"

	"github.com/Y716/gatorcli/gatorcli/internal/config"
)

func main() {
	cfgFile, err := config.Read()
	if err != nil {
		fmt.Printf("Error found: %v\n", err)
	}

	s := state{
		config: &cfgFile,
	}

	cmds := commands{
		all_commands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Arguments not sufficent: Need at least 2 arguments")
		os.Exit(1)
	}
	cmd_name := args[1]
	args_slice := args[2:]

	cmd := command{
		name: cmd_name,
		args: args_slice,
	}

	err = cmds.run(&s, cmd)
	if err != nil {
		fmt.Printf("Error found: %v", err)
		os.Exit(1)
	}
}

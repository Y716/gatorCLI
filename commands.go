package main

import (
	"fmt"

	"github.com/Y716/gatorcli/gatorcli/internal/config"
)

type state struct {
	config *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	all_commands map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("Expect an argument <username>")
	}

	username := cmd.args[0]
	err := s.config.SetUser(username)
	if err != nil {
		return err
	}

	fmt.Printf("The username has been set to: %s\n", username)
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	err := c.all_commands[cmd.name](s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.all_commands[name] = f
}

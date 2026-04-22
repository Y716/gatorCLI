package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Y716/gatorcli/gatorcli/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("expect an argument <username>\n")
	}

	username := cmd.args[0]
	err := s.config.SetUser(username)
	if err != nil {
		return err
	}

	ctx := context.Background()

	user, err := s.db.GetUser(ctx, username)
	if err != nil {
		return fmt.Errorf("username is not registered: %s", username)
	}

	fmt.Printf("Login success. Welcome: %s\n", user.Name)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("expect an argument <username>\n")
	}

	username := cmd.args[0]
	ctx := context.Background()

	user, err := s.db.GetUser(ctx, username)
	if user.Name != "" {
		fmt.Printf("username %s already exists\n", user.Name)
		os.Exit(1)
		return err
	}

	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}

	newUser, err := s.db.CreateUser(ctx, params)
	if err != nil {
		return err
	}

	fmt.Println("A new user has been created")
	fmt.Printf("ID: %d\n", newUser.ID)
	fmt.Printf("Name: %s\n", newUser.Name)
	fmt.Printf("CreatedAt: %s\n", newUser.CreatedAt)
	fmt.Printf("UpdatedAt: %s\n", newUser.UpdatedAt)

	err = s.config.SetUser(username)
	if err != nil {
		return err
	}

	return nil
}

func handlerReset(s *state, cmd command) error {
	ctx := context.Background()

	err := s.db.DeleteUsers(ctx)
	if err != nil {
		os.Exit(1)
		return fmt.Errorf("error found while reseting users table: %v", err)
	}

	fmt.Println("Reset table successful")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	ctx := context.Background()

	users, err := s.db.GetUsers(ctx)
	if err != nil {
		return err
	}

	for _, user := range users {
		fmt.Printf("* %s", user)
		if s.config.CurrentUserName == user {
			fmt.Print(" (current)")
		} else {
			fmt.Print("\n")
		}
	}

	return nil
}

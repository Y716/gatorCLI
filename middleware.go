package main

import (
	"context"

	"github.com/Y716/gatorcli/gatorcli/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, c command) error {
		ctx := context.Background()
		user_db, err := s.db.GetUser(ctx, s.config.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, c, user_db)
	}
}

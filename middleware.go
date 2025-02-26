package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Ernestlph/Blog_Aggregator/internal/database"
)

// middleware.go
func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		// Get current user from config
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("must be logged in to use %s command", cmd.Name)
			}
			return fmt.Errorf("couldn't get user: %w", err)
		}

		// Call original handler with user parameter
		return handler(s, cmd, user)
	}
}

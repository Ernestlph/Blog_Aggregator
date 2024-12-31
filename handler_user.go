package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Ernestlph/Blog_Aggregator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	name := cmd.Args[0]

	// Check if username exists, if it doesn't exist return an error with code 1
	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user %v does not exist", name)
		}
		return fmt.Errorf("couldn't get user: %w", err)
	}

	// Set current user
	if err := s.cfg.SetUser(name); err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Printf("Switched to user: %s\n", name)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	username := cmd.Args[0]

	user, err := s.db.CreateUser(context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      username},
	)
	// Exit with code 1 if there's already a user with the same name
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return fmt.Errorf("user %v already exists", username)
		}
	}

	if err != nil {
		return fmt.Errorf("couldn't create user: %w", err)

	}

	fmt.Printf("User %v created successfully!, now you can login\n", username)
	s.cfg.SetUser(user.Name)
	return nil
}

func handlerResetDatabase(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	if err := s.db.ResetDatabase(context.Background()); err != nil {
		return fmt.Errorf("couldn't reset database: %w", err)
	}
	fmt.Println("Database reset successfully!")
	return nil
}

func handlerGetUsers(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	users, err := s.db.GetUsers(context.Background(), s.cfg.CurrentUserName)

	if err != nil {
		return fmt.Errorf("couldn't get users: %w", err)
	}

	fmt.Println("Users:")
	for _, user := range users {
		fmt.Printf("  - %s\n", user.Name_2)
	}
	return nil
}

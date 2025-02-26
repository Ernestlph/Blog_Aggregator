package main

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/Ernestlph/Blog_Aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command, user database.User) error { // Modified handler signature
	// 1. Input validation
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}
	name := cmd.Args[0]
	url := cmd.Args[1]

	// 2. Get current user (No longer needed - user is passed by middleware)

	// 3. Create feed with proper error handling
	feed, err := s.db.CreateFeed(context.Background(),
		database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(), // Use UTC for consistency
			UpdatedAt: time.Now().UTC(),
			Name:      name,
			Url:       url,
			UserID:    user.ID, // Use user from middleware
		})

	if err != nil {
		return fmt.Errorf("failed to create feed: %w", err)
	}
	_ = feed

	// 4. Get and verify the created feed
	showfeed, err := s.db.GetFeed(context.Background(), name)
	if err != nil {
		return fmt.Errorf("couldn't get feed: %w", err)
	}

	// 5. Print success message
	fmt.Printf("Feed %q added successfully!\n", showfeed.Name)

	// 6. Print feed details using reflection
	v := reflect.ValueOf(showfeed)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldName := t.Field(i).Name
		fieldValue := v.Field(i).Interface()

		if fieldValue == "" {
			fmt.Printf("Warning: Field %s is empty\n", fieldName)
		} else {
			fmt.Printf("%s: %v\n", fieldName, fieldValue)
		}
	}

	// 7. Create the feedfollow relationship with the created feed
	err = handlerFollowFeed(s, command{Name: "followfeed", Args: []string{showfeed.Url}}, user) // Pass user to handlerFollowFeed
	if err != nil {
		return err // Return error from handlerFollowFeed if it occurs
	}

	return nil
}

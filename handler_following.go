package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Ernestlph/Blog_Aggregator/internal/database"
)

func handlerFollowing(s *state, cmd command, user database.User) error {
	// Validate input
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	// Get following
	following, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("you are not following anyone")
		}
		return err
	}
	// Print out all the names of the feeds the user is following
	fmt.Printf("You are following:\n")
	for _, feed := range following {
		fmt.Printf("  - %s\n", feed.FeedName)
	}

	return nil
}

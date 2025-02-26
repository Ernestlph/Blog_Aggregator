package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Ernestlph/Blog_Aggregator/internal/database"
	"github.com/lib/pq"
)

func handlerFollowFeed(s *state, cmd command, user database.User) error { // Modified handler signature to accept user
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <feed_url>", cmd.Name)
	}

	feedURL := cmd.Args[0]
	ctx := context.Background()

	// Verify feed exists
	feed, err := s.db.GetFeedByURL(ctx, feedURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("feed %v does not exist", feedURL)
		}
		return fmt.Errorf("couldn't get feed: %w", err)
	}

	// Create follow relationship
	_, err = s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		FeedID: feed.ID,
		UserID: user.ID, // Use the user passed by the middleware
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return fmt.Errorf("user %s is already following feed %s", user.Name, feed.Name)
		}
		return fmt.Errorf("failed to create feed follow: %w", err)
	}

	fmt.Printf("User %s is now following feed %s\n", user.Name, feed.Name)
	return nil
}

package main

import (
	"context"
	"fmt"

	"github.com/Ernestlph/Blog_Aggregator/internal/database"
)

func handlerUnfollowFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <feed_url>", cmd.Name)
	}

	feedURL := cmd.Args[0]
	ctx := context.Background()

	err := s.db.DeleteFeedFollowByUserAndFeedURL(ctx, database.DeleteFeedFollowByUserAndFeedURLParams{
		UserID: user.ID, // Assuming 'user' is the logged-in user object
		Url:    feedURL,
	})
	if err != nil {
		return fmt.Errorf("failed to unfollow feed: %w", err)
	}

	fmt.Printf("User %s is no longer following feed with URL: %s\n", user.Name, feedURL)
	return nil
}

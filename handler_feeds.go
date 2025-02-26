package main

import (
	"context"
	"fmt"
)

func handlerFeeds(s *state, cmd command) error {
	// Validate input
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	// get all feeds
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	// Print out details of all feeds
	for _, feed := range feeds {

		// get User name from feed.UserID
		user, err := s.db.GetUserfromID(context.Background(), feed.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("%s: %s Added by: %s\n", feed.Name, feed.Url, user.Name)

	}

	return nil
}

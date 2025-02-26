package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Ernestlph/Blog_Aggregator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2 // Default limit

	if len(cmd.Args) > 1 {
		return fmt.Errorf("usage: %s <optional limit>", cmd.Name)
	}
	if len(cmd.Args) == 1 {
		providedLimit, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("invalid limit: %v", cmd.Args[0])
		}
		limit = providedLimit
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit), // Convert limit to int32 as expected by sqlc
	})
	if err != nil {
		return fmt.Errorf("couldn't get posts: %w", err)
	}

	fmt.Println("Posts for you:")
	if len(posts) == 0 {
		fmt.Println("  No posts found from feeds you follow.")
	} else {
		for _, post := range posts {
			fmt.Printf("  - Title: %s\n", post.Title)
			fmt.Printf("    URL: %s\n", post.Url)
			fmt.Printf("    Published At: %s\n", post.PublishedAt.Format(time.RFC3339)) // Format time for display
			if post.Description.Valid {                                                 // Check if description is valid (not NULL)
				description := post.Description.String // Access the string value
				// Truncate description if it's too long for display
				if len(description) > 200 {
					description = description[:200] + "..."
				}
				fmt.Printf("    Description: %s\n", description)
			}
			fmt.Println() // Add an empty line between posts
		}
	}

	return nil
}

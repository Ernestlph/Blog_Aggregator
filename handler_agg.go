package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/Ernestlph/Blog_Aggregator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func parseTime(dateStr string) (time.Time, error) {
	// Attempt to parse with multiple formats
	formats := []string{
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		"Mon, 2 Jan 2006 15:04:05 MST", // Common RSS format
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.DateTime,
		time.DateOnly,
		time.TimeOnly,
	}

	for _, format := range formats {
		parsedTime, err := time.Parse(format, dateStr)
		if err == nil {
			return parsedTime.UTC(), nil // Convert to UTC for database storage
		}
	}
	return time.Time{}, fmt.Errorf("could not parse time: %s", dateStr)
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <time_between_reqs>", cmd.Name)
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("could not parse duration: %w", err)
	}

	fmt.Printf("Collecting feeds every %s...\n", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)

	scrapeFeeds(s) // Run immediately when agg command starts

	for range ticker.C { // Run every time the ticker ticks
		scrapeFeeds(s)
	}
	return nil // This part of the code will likely not be reached in a ticker loop
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	// Create a new request with NewRequestWithContext
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("couldn't create request: %w", err)
	}
	// Set a user agent header
	req.Header.Set("User-Agent", "gator")

	// Create a http client and make request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("couldn't make request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Decode response body

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("couldn't read response body: %w", err)
	}

	var feed RSSFeed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, fmt.Errorf("couldn't unmarshal response body: %w", err)
	}
	// Unescape HTML entities for the channel
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	// Unescape HTML entities for each item
	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}

	// For debugging: Print the entire struct
	fmt.Printf("Feed Title: %s\nDescription: %s\nLink: %s\n",
		feed.Channel.Title, feed.Channel.Description, feed.Channel.Link)
	for _, item := range feed.Channel.Item {
		fmt.Printf("Item Title: %s\nDescription: %s\nLink: %s\nPubDate: %s\n",
			item.Title, item.Description, item.Link, item.PubDate)
	}

	return &feed, err
}

func scrapeFeeds(s *state) {
	ctx := context.Background()

	fmt.Println("Fetching next feed...")
	feedRow, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		fmt.Printf("Error getting next feed to fetch: %v\n", err)
		return
	}

	if feedRow.ID.String() == "" {
		fmt.Println("No feeds to fetch at the moment.")
		return
	}

	fmt.Printf("Fetching feed: %s from %s\n", feedRow.Name, feedRow.Url)
	err = s.db.MarkFeedFetched(ctx, feedRow.ID)
	if err != nil {
		fmt.Printf("Error marking feed as fetched: %v\n", err)
		return
	}

	rssFeed, err := fetchFeed(ctx, feedRow.Url)
	if err != nil {
		fmt.Printf("Error fetching feed content: %v\n", err)
		return
	}

	fmt.Printf("Feed fetched, saving posts...\n") // Indicate saving posts

	for _, item := range rssFeed.Channel.Item {
		publishedAt := time.Now().UTC() // Default to now if parsing fails
		if item.PubDate != "" {
			parsedTime, perr := parseTime(item.PubDate)
			if perr != nil {
				fmt.Printf("Error parsing time '%s': %v. Using current time.\n", item.PubDate, perr)
			} else {
				publishedAt = parsedTime
			}
		}

		_, err := s.db.CreatePost(ctx, database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: true}, // Convert item.Description to sql.NullString
			PublishedAt: publishedAt,
			FeedID:      feedRow.ID,
		})
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				continue // Ignore duplicate URL errors
			}
			fmt.Printf("Error creating post: %v\n", err) // Log other errors
		}
	}
	fmt.Printf("Feed %s posts saved!\n----------------------\n", feedRow.Name) // Indicate posts saved
}

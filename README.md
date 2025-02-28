# Blog_Aggregator
Blog Aggregator following Boot.dev

# Gator - Blog Aggregator CLI

Gator is a command-line interface (CLI) application that aggregates content from RSS feeds and allows you to browse the latest posts from your followed sources directly in your terminal.

## Prerequisites

Before running Gator, ensure you have the following installed:

*   **Go:** Gator is written in Go, so you need the Go toolchain installed. You can download and install Go from the [official Go website](https://go.dev/dl/). Follow the installation instructions for your operating system.

*   **PostgreSQL:** Gator uses PostgreSQL as its database. You'll need to have PostgreSQL installed and running. You can download and install PostgreSQL from the [official PostgreSQL website](https://www.postgresql.org/download/). Make sure to set up a PostgreSQL database and user for Gator.

## Installation

To install the `gator` CLI, use the `go install` command:

```bash
go install github.com/Ernestlph/Blog_Aggregator@latest
```

## Configuration

Gator requires a configuration file to connect to your PostgreSQL database and manage user settings.

1.  **Create a `.gatorconfig.json` file** in the home directory

2.  **Add the following configuration** to `.gatorconfig.json`, modifying the **placeholder values** below to match your PostgreSQL setup and desired settings:
```json
{"db_url":"postgres://[default_username]:postgres@localhost:5432/gator?sslmode=disable","current_user_name":"[default_username]"}
```
   

    **Important:** **Replace the placeholder values for username**


## Available Commands

Here are the commands you can use with the Gator CLI:

*   **`create_user <username> <password>`**: Creates a new user account.
    ```bash
    gator create_user john.doe mySecretPassword
    ```

*   **`login <username> <password>`**: Logs in an existing user.
    ```bash
    gator login john.doe mySecretPassword
    ```

*   **`add_feed <name> <feed_url>`**: Adds a new RSS feed to be tracked.
    ```bash
    gator add_feed "TechCrunch" [https://techcrunch.com/feed/](https://techcrunch.com/feed/)
    ```

*   **`follow_feed <feed_url>`**: Starts following a specific feed.
    ```bash
    gator follow_feed [https://techcrunch.com/feed/](https://techcrunch.com/feed/)
    ```

*   **`following`**: Lists the feeds you are currently following.
    ```bash
    gator following
    ```

*   **`browse [limit]`**: Browses the latest posts from the feeds you follow.  Optionally, you can specify a limit for the number of posts to display.
    ```bash
    gator browse
    gator browse 10 # Browse the latest 10 posts
    ```

*   **`agg <time_between_requests>`**: Continuously aggregates feeds and saves new posts to the database.  `<time_between_requests>` is a duration string like `10s`, `1m`, `1h`.
    ```bash
    gator agg 1m # Aggregate feeds every 1 minute
    ```

*   **`help`**: Displays a list of available commands and their descriptions.
    ```bash
    gator help
    ```
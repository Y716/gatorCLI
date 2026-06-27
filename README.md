# 🐊 gatorCLI
**THIS IS A COURSE PROJECT FROM BOOT.DEV. MORE OF IT IN THEIR WEBSITE**

A command-line RSS feed aggregator built in Go. gatorCLI lets you register users, subscribe to RSS feeds, automatically scrape new posts on a configurable interval, and browse aggregated content — all from your terminal.

## Features

- **Multi-user support** — register and switch between users, with the current session persisted in a local config file
- **Feed management** — add RSS feeds, list all available feeds, and see who added each one
- **Follow system** — follow/unfollow feeds independently of other users; each user has their own subscription list
- **Automatic aggregation** — a long-running `agg` command scrapes feeds on a configurable interval, parsing RSS/XML and storing posts to the database
- **Post browsing** — browse posts from your followed feeds, sorted by most recent, with a configurable limit
- **Auth middleware** — commands that require a logged-in user are protected by a middleware pattern
- **Database migrations** — incremental schema evolution using Goose
- **Type-safe SQL** — query generation with sqlc

## Tech Stack

| Layer       | Technology           |
|------------|---------------------|
| Language    | Go                  |
| Database    | PostgreSQL          |
| SQL Tooling | sqlc, Goose         |
| Config      | JSON (~/.gatorconfig.json) |
| IDs         | UUID (google/uuid)  |
| RSS Parsing | encoding/xml (stdlib) |

## Project Structure

```
.
├── main.go                  # Entry point, command registration
├── commands.go              # State, command, and command registry types
├── handlerUser.go           # login, register, reset, users commands
├── handlerFeed.go           # agg, addfeed, feeds, follow, unfollow, following, browse
├── middleware.go             # Logged-in user middleware
├── rss_feed.go              # RSS feed fetcher and XML parser
├── internal/
│   ├── config/              # Config file read/write (~/.gatorconfig.json)
│   └── database/            # sqlc-generated database layer
├── sql/
│   ├── schema/              # Goose migration files (users, feeds, feed_follows, posts)
│   └── queries/             # sqlc query definitions
└── sqlc.yaml                # sqlc configuration
```

## Commands

| Command | Auth Required | Usage | Description |
|---------|:---:|-------|-------------|
| `register` | — | `gatorCLI register <username>` | Create a new user and log in |
| `login` | — | `gatorCLI login <username>` | Switch to an existing user |
| `users` | — | `gatorCLI users` | List all users (marks current) |
| `reset` | — | `gatorCLI reset` | Delete all users (dev tool) |
| `addfeed` | ✓ | `gatorCLI addfeed <name> <url>` | Add a feed and auto-follow it |
| `feeds` | — | `gatorCLI feeds` | List all feeds with their creators |
| `follow` | ✓ | `gatorCLI follow <url>` | Follow an existing feed |
| `unfollow` | ✓ | `gatorCLI unfollow <url>` | Unfollow a feed |
| `following` | ✓ | `gatorCLI following` | List feeds you follow |
| `agg` | — | `gatorCLI agg <duration>` | Start scraping feeds on an interval (e.g. `1m`, `30s`) |
| `browse` | ✓ | `gatorCLI browse [limit]` | Browse posts from followed feeds (default: 2) |

## Database Schema

The app uses five tables managed through Goose migrations:

- **users** — id, name, timestamps
- **feeds** — id, name, url (unique), user_id (creator), last_fetched_at, timestamps
- **feed_follows** — many-to-many join between users and feeds (unique per user-feed pair)
- **posts** — id, title, description, url (unique), published_at, feed_id, timestamps

## Getting Started

### Prerequisites

- Go 1.22+
- PostgreSQL
- [Goose](https://github.com/pressly/goose) (for migrations)
- [sqlc](https://sqlc.dev/) (for query generation)

### Setup

1. **Clone the repo**
   ```bash
   git clone https://github.com/Y716/gatorCLI.git
   cd gatorCLI
   ```

2. **Create a config file** at `~/.gatorconfig.json`
   ```json
   {
     "db_url": "postgres://user:password@localhost:5432/gator?sslmode=disable"
   }
   ```

3. **Create the database and run migrations**
   ```bash
   createdb gator
   goose -dir sql/schema postgres "your-db-url" up
   ```

4. **Generate database code** (only if you modify queries)
   ```bash
   sqlc generate
   ```

5. **Build and run**
   ```bash
   go build -o gatorCLI && ./gatorCLI <command> [args]
   ```

## Usage Example

```bash
# Register and start using
./gatorCLI register yasin

# Add some feeds
./gatorCLI addfeed "Go Blog" https://go.dev/blog/feed.atom
./gatorCLI addfeed "Hacker News" https://hnrss.org/frontpage

# Start the aggregator (scrapes every 2 minutes)
./gatorCLI agg 2m

# In another terminal, browse your posts
./gatorCLI browse 5
```

## License

This project is unlicensed — feel free to use it however you like.

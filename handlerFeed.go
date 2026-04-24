package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Y716/gatorcli/gatorcli/internal/database"
	"github.com/google/uuid"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("expect an argument <duration>\n")
	}
	time_between_reqs, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Collecting feeds every %v\n", time_between_reqs)

	ticker := time.NewTicker(time_between_reqs)

	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			return err
		}
	}
	return nil
}

func scrapeFeeds(s *state) error {
	ctx := context.Background()

	feedData, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}

	params := database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt:     time.Now(),
		ID:            feedData.ID,
	}
	s.db.MarkFeedFetched(ctx, params)

	feedPosts, err := fetchFeed(ctx, feedData.Url)
	if err != nil {
		return err
	}
	fmt.Printf("Scraping from '%s':\n", feedData.Name)
	for i, post := range feedPosts.Channel.Item {

		publishedAt := sql.NullTime{}
		format := []string{time.RFC1123Z, time.RFC1123, time.RFC3339}
		for _, v := range format {
			if parsed, err := time.Parse(v, post.PubDate); err == nil {
				publishedAt = sql.NullTime{Time: parsed, Valid: true}
				break
			}
		}

		postParams := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       post.Title,
			Description: post.Description,
			Url:         post.Link,
			PublishedAt: publishedAt,
			FeedID:      feedData.ID,
		}

		postDB, err := s.db.CreatePost(ctx, postParams)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}

		fmt.Printf("%d. %s\n", i+1, postDB.Title)

	}

	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("expect an argument <name> <url>\n")
	}
	feedName := cmd.args[0]
	url := cmd.args[1]

	ctx := context.Background()
	user_id := user.ID

	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       url,
		UserID:    user_id,
	}

	newFeed, err := s.db.CreateFeed(ctx, params)
	if err != nil {
		return err
	}

	feed_db, err := s.db.GetFeedByUrl(ctx, url)
	if err != nil {
		return err
	}
	feed_id := feed_db.ID

	ff_params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user_id,
		FeedID:    feed_id,
	}

	_, err = s.db.CreateFeedFollow(ctx, ff_params)
	if err != nil {
		return err
	}

	fmt.Println("A new feed has been created")
	fmt.Printf("ID: %d\n", newFeed.ID)
	fmt.Printf("Feed Name: %s\n", newFeed.Name)
	fmt.Printf("Feed URL: %s\n", newFeed.Url)
	fmt.Printf("User ID: %s\n", newFeed.UserID)
	fmt.Printf("CreatedAt: %s\n", newFeed.CreatedAt)
	fmt.Printf("UpdatedAt: %s\n", newFeed.UpdatedAt)

	return nil
}

func handlerListFeeds(s *state, cmd command) error {
	ctx := context.Background()

	feedList, err := s.db.GetFeeds(ctx)
	if err != nil {
		return err
	}

	fmt.Println("List of all feeds:")
	for i, feed := range feedList {
		feedName := feed.Name
		feedURL := feed.Url
		feedUserID := feed.UserID

		user_post_feed, err := s.db.GetUserByID(ctx, feedUserID)
		if err != nil {
			return err
		}
		fmt.Printf("%d. '%s': '%s'. Posted by: '%s'\n", i+1, feedName, feedURL, user_post_feed)
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("expect an argument <url>\n")
	}
	url := cmd.args[0]
	ctx := context.Background()

	user_id := user.ID

	feed_db, err := s.db.GetFeedByUrl(ctx, url)
	if err != nil {
		return err
	}
	feed_id := feed_db.ID

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user_id,
		FeedID:    feed_id,
	}

	_, err = s.db.CreateFeedFollow(ctx, params)
	if err != nil {
		return err
	}

	fmt.Printf("%s has followed %s\n", s.config.CurrentUserName, feed_db.Name)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	ctx := context.Background()

	ff_db, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return err
	}

	fmt.Printf("List of all feeds followed by '%s':\n", s.config.CurrentUserName)
	for i, onefeed := range ff_db {
		feedName := onefeed.FeedName
		fmt.Printf("%d. %s\n", i+1, feedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("expect an argument <url>\n")
	}
	url := cmd.args[0]
	ctx := context.Background()
	feed_db, err := s.db.GetFeedByUrl(ctx, url)
	if err != nil {
		return err
	}

	params := database.DeleteFollowingParams{
		UserID: user.ID,
		FeedID: feed_db.ID,
	}

	s.db.DeleteFollowing(ctx, params)
	fmt.Printf("%s has unfollowed %s\n", s.config.CurrentUserName, feed_db.Name)
	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := "2"
	if len(cmd.args) == 1 {
		limit = cmd.args[0]
	}
	ctx := context.Background()
	intlimit, _ := strconv.Atoi(limit)
	params := database.GetPostForUserParams{
		UserID: user.ID,
		Limit:  int32(intlimit),
	}
	posts, err := s.db.GetPostForUser(ctx, params)
	if err != nil {
		return err
	}

	for _, post := range posts {
		fmt.Printf("Title: %s\n", post.Title)
		fmt.Printf("%s\n\n", post.Description)
	}

	return nil
}

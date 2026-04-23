package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Y716/gatorcli/gatorcli/internal/database"
	"github.com/google/uuid"
)

func handlerAgg(s *state, cmd command) error {
	ctx := context.Background()

	feed, err := fetchFeed(ctx, "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Println(feed)

	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("expect an argument <name> <url>\n")
	}
	feedName := cmd.args[0]
	url := cmd.args[1]

	ctx := context.Background()
	db, err := s.db.GetUser(ctx, s.config.CurrentUserName)
	if err != nil {
		return err
	}
	user_id := db.ID

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

func handlerFollow(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("expect an argument <url>\n")
	}
	url := cmd.args[0]
	ctx := context.Background()

	user_db, err := s.db.GetUser(ctx, s.config.CurrentUserName)
	if err != nil {
		return err
	}
	user_id := user_db.ID

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

func handlerFollowing(s *state, cmd command) error {
	ctx := context.Background()

	ff_db, err := s.db.GetFeedFollowsForUser(ctx, s.config.CurrentUserName)
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

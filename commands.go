package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/abdullah0iq/gator/internal/database"
	"github.com/google/uuid"
)

type commands struct {
	commandsMap map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commandsMap[name] = f
}
func (c *commands) run(s *state, cmd command) error {
	f, ok := c.commandsMap[cmd.name]
	if !ok {
		log.Fatal("there is no such command")
	}
	if err := f(s, cmd); err != nil {
		return err
	}
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("the login handler expects a single argument, the username.")
	}
	user, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		log.Fatal("user Does not exists")
	}
	s.config.SetUser(user.Name)
	fmt.Printf("logged in with username: %s", user.Name)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("the register handler expects a single argument, the username.")
	}
	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err == nil {
		log.Fatal("user already exists")
	}
	dbUser, err := s.db.CreateUser(context.Background(), database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: cmd.args[0]})
	if err != nil {

		log.Fatalf("Couldnt create user : %v", err)
	}
	s.config.SetUser(dbUser.Name)
	fmt.Println("User was created")
	fmt.Println(dbUser)
	return nil
}

func handlerReset(s *state, cmd command) error {
	if err := s.db.ResetTable(context.Background()); err != nil {
		log.Fatalf("Failed to reset the table: %v", err)
	}
	log.Println("the table was successfully reset")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.ListUsers(context.Background())
	if err != nil {
		log.Fatal("Failed to fetch all the users from table")
	}
	for _, user := range users {
		if user == s.config.CurrentUserName {
			fmt.Printf("* %s (current)\n", user)
		} else {
			fmt.Printf("* %s \n", user)

		}
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) < 1 || len(cmd.args) > 2 {
		return fmt.Errorf("usage: %v <time_between_reqs>", cmd.name)
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}

	log.Printf("Collecting feeds every %s...", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Println("Couldn't get next feeds to fetch", err)
		return
	}
	log.Println("Found a feed to fetch!")
	scrapeFeed(s.db, feed)
}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	 err := db.MarkFeedFetched(context.Background(), feed.Url)
	if err != nil {
		log.Printf("Couldn't mark feed %s fetched: %v", feed.Name, err)
		return
	}

	feedData, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("Couldn't collect feed %s: %v", feed.Name, err)
		return
	}
	for _, item := range feedData.Channel.Items {
		fmt.Printf("Found post: %s\n", item.Title)
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Items))
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("the addfeed command needs two argument: nameOfTheFeed feedURL")
	}

	feed, err := s.db.InsertFeed(context.Background(), database.InsertFeedParams{Name: cmd.args[0], Url: cmd.args[1], UserID: user.ID})
	if err != nil {
		return fmt.Errorf("couldn't not add the feed : %v", err)
	}
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{ID: uuid.New(), UserID: user.ID, FeedID: feed.Url})

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("the feeds command doesnt take any arguments")
	}
	feeds, err := s.db.GetAllFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Printf("- %v\n- %v\n", feed.Name, feed.Name_2)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("you need a to type the blog url after the follow command")
	} else if len(cmd.args) > 1 {
		return fmt.Errorf("follow command takes one argument (blog url) but found %v", len(cmd.args))
	}
	feed, err := s.db.GetFeed(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{ID: uuid.New(), UserID: user.ID, FeedID: feed.Url})
	if err != nil {
		return err
	}

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("following command expect no argument")
	}

	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("couldn't fetch your following feed")
	}
	for _, feed := range feeds {
		fmt.Println(feed.FeedName)
	}
	return nil
}

func handlerUnFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("unfollow command need one argument: feed url")
	}
	if err := s.db.UnFollowFeed(context.Background(), database.UnFollowFeedParams{UserID: user.ID, FeedID: cmd.args[0]}); err != nil {
		return fmt.Errorf("could not unfollow the feed: %v", err)
	}
	return nil

}

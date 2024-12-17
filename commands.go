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
	rssFeed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(rssFeed)

	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) <= 1 {
		return fmt.Errorf("the addfeed command needs to argument: nameOfTheFeed feedURL")
	}
	user, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("could not get the user: %v", err)
	}
	feed, err := s.db.InsertFeed(context.Background(), database.InsertFeedParams{Name: cmd.args[0], Url: cmd.args[1], UserID: user.ID})
	if err != nil {
		return fmt.Errorf("couldn't not add the feed : %v", err)
	}
	fmt.Println(feed)
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
	
	for _,feed := range feeds {
		fmt.Printf("- %v\n- %v\n", feed.Name,feed.Name_2)
	}
	return nil
}

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/abdullah0iq/gator/internal/config"
	"github.com/abdullah0iq/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	// initiate config
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	cfg.DBURL = "postgres://Abdullah:@localhost:5432/gator?sslmode=disable"

	//initiate state
	s := state{config: &cfg}

	//initiate commands
	cmds := commands{commandsMap: map[string]func(*state, command) error{}}

	//initiate register command
	cmds.register("login", handlerLogin)
	//initiate register command
	cmds.register("register", handlerRegister)
	//initiate reset command
	cmds.register("reset", handlerReset)
	//initiate users command
	cmds.register("users", handlerUsers)
	// initiate agg command
	cmds.register("agg", handlerAgg)
	// initiate addfeed command
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	//initiate feeds command
	cmds.register("feeds", handlerFeeds)
	//initiate follow command
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	//initiate following command
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow",middlewareLoggedIn(handlerUnFollow))

	//initiating db
	db, err := sql.Open("postgres", s.config.DBURL)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	} else {
		log.Println("Successfully connected to the database!")
	}
	dbQueries := database.New(db)
	s.db = dbQueries

	//getting arguments
	args := os.Args
	if len(args) < 2 {
		log.Fatal("No argument")
	}
	commandName := args[1]
	args = args[2:]
	cmd := command{name: commandName, args: args}
	err = cmds.run(&s, cmd)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	os.Exit(0)
}

func middlewareLoggedIn(handler func(*state, command, database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
		if err != nil {
			return fmt.Errorf("you must be logged in to perform this action")
		}
		return handler(s, cmd, user) // Pass the logged-in user to the handler
	}
}

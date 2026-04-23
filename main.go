package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Y716/gatorcli/gatorcli/internal/config"
	"github.com/Y716/gatorcli/gatorcli/internal/database"
	_ "github.com/lib/pq"
)

func main() {

	cfgFile, err := config.Read()
	if err != nil {
		fmt.Printf("error found accessing config file: %v\n", err)
	}

	s := state{
		config: &cfgFile,
	}

	db, err := sql.Open("postgres", s.config.DbURL)
	if err != nil {
		fmt.Printf("Error found accessing database: %v\n", err)
	}

	dbQueries := database.New(db)

	s.db = dbQueries

	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerAddFeed)
	cmds.register("feeds", handlerListFeeds)
	cmds.register("follow", handlerFollow)
	cmds.register("following", handlerFollowing)
	args := os.Args

	if len(args) < 2 {
		fmt.Println("arguments not sufficent: Need at least 2 arguments")
		os.Exit(1)
	}
	cmd_name := args[1]
	args_slice := args[2:]

	cmd := command{
		name: cmd_name,
		args: args_slice,
	}

	err = cmds.run(&s, cmd)
	if err != nil {
		fmt.Printf("error found running command: %v\n", err)
		os.Exit(1)
	}
}

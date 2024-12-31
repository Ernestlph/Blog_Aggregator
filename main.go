package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Ernestlph/Blog_Aggregator/internal/config"
	"github.com/Ernestlph/Blog_Aggregator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func parseArgs() (cmdName string, cmdArgs []string, err error) {
	if len(os.Args) == 1 {
		fmt.Println("Error: not enough arguments were provided")
		os.Exit(1)
		return
	}
	if (len(os.Args) == 2) && ((os.Args[1] == "login") || (os.Args[1] == "register")) {
		fmt.Println("Error: a username is required")
		os.Exit(1)
		return
	}
	if len(os.Args) < 2 && os.Args[1] != "reset" {
		fmt.Println("Error: unknown command")
		os.Exit(1)
		return
	}
	cmdName = os.Args[1]
	cmdArgs = os.Args[2:]
	return cmdName, cmdArgs, nil

}

func main() {
	// Parse command line args first
	cmdName, cmdArgs, err := parseArgs()
	if err != nil {
		log.Fatal(err)
	}

	// Reads config which also returns a config variable
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	// Open database connection
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatal(err)
	}

	// Create queries object
	dbQueries := database.New(db)

	// Create a state object, which contains the config from Read
	programState := &state{
		cfg: &cfg,
		db:  dbQueries,
	}

	// Create commands struct and initializes empty map
	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	// Register commands
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerResetDatabase)
	cmds.register("users", handlerGetUsers)

	// Runs command if command not found return error with code 1
	err = cmds.run(programState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		if err.Error() == "unknown command" {
			os.Exit(1)
		}
		log.Fatal(err)
	}
}

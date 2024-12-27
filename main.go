package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Ernestlph/Blog_Aggregator/internal/config"
)

type state struct {
	cfg *config.Config
}

func parseArgs() (cmdName string, cmdArgs []string, err error) {
	if len(os.Args) == 1 {
		fmt.Println("Error: not enough arguments were provided")
		os.Exit(1)
		return
	}
	if (len(os.Args) == 2) && (os.Args[1] == "login") {
		fmt.Println("Error: a username is required")
		os.Exit(1)
		return
	}
	if len(os.Args) < 3 {
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

	// Create a state object, which contains the config from Read
	programState := &state{
		cfg: &cfg,
	}

	// Create commands struct and initializes empty map
	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	// Register commands
	cmds.register("login", handlerLogin)

	// Runs command
	err = cmds.run(programState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Ernestlph/Blog_Aggregator/internal/config"
)

func main() {
	// Step 1: Read the config file
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}
	fmt.Printf("Read config: %+v\n", cfg)

	// Store config file in a new instance of state struct
	newstate := config.State{Config: &cfg}

	// Create a new instance of the commands struct with an initialized map of handler functions
	commands := config.Commands{
		// Ensure any map or lists are initiated
		Handlers: make(map[string]func(*config.State, config.Command) error),
	}

	// Register a handler function for the login command
	commands.Register("login", handlerLogin)

	// Use os.Args to get the command-line arguments passed in by the user
	args := os.Args[1:] // Skip the first argument as it is the program name

	if len(args) < 1 {
		log.Println("No command provided!")
		return
	}

	// Construct the Command struct
	cmd := config.Command{
		Name: args[0],
		Args: args[1:], // Remaining arguments are for the command
	}

	// Run the command using the commands struct
	if err := commands.Run(&state, cmd); err != nil {
		log.Fatalf("Error running command: %v", err)
	}

	// Read the config file again after update
	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("Error reading updated config: %v", err)
	}

	fmt.Printf("Updated Config: %+v\n", cfg)
}

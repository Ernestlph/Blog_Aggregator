package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

// Config represents the structure of the JSON configuration file.
type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

// Helper function to get the config file path
func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, configFileName), nil
}

// Helper function to write the config back to the file
func write(cfg Config) error {
	configPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// Read loads the JSON file into a Config struct
func Read() (Config, error) {
	fullPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(fullPath)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	cfg := Config{}
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// SetUser updates the current_user_name and writes the updated config to the JSON file
func (config *Config) SetUser(userName string) error {
	config.CurrentUserName = userName

	if err := write(*config); err != nil {
		return err
	}

	return nil
}

// Create a state struct that holds a pointer to a config
type State struct {
	Config *Config
}

// Create a command struct that contains a name and a slice of string arguments for example login command , the name would be "login" and the handler will expect the arguments slice to contain one string, the username
type Command struct {
	Name string
	Args []string
}

// Create a login handler function: func handlerLogin(s *state, cmd command) error. This will be the function signature of all command handlers.
// If the command's arg's slice is empty, return an error; the login handler expects a single argument, the username.
// Use the state's access to the config struct to set the user to the given username. Remember to return any errors.
// Print a message to the terminal that the user has been set.

func handlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return os.ErrInvalid
	}
	userName := cmd.Args[0] // assumes the first arg is the username
	err := s.Config.SetUser(userName)
	if err != nil {
		return err
	}
	fmt.Println("User set to", userName)
	return nil
}

//Create a commands struct. This will hold all the commands the CLI can handle. Add a map[string]func(*state, command) error field to it. This will be a map of command names to their handler functions.

type Commands struct {
	Handlers map[string]func(*State, Command) error
}

func (c *Commands) Register(name string, handler func(*State, Command) error) {
	if c.Handlers == nil {
		c.Handlers = make(map[string]func(*State, Command) error)
	}
	c.Handlers[name] = handler
}

func (c *Commands) Run(s *State, cmd Command) error {
	if handler, ok := c.Handlers[cmd.Name]; ok {
		return handler(s, cmd)
	}
	return os.ErrInvalid
}

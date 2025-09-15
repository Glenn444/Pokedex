package cli


import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/Glenn444/pokedexcli/internal/api"
	"github.com/Glenn444/pokedexcli/internal/storage"
)

type CLI struct {
	config    *Config
	apiClient *api.Client
	storage   *storage.FileStorage
	commands  map[string]Command
}

type Command struct {
	Name        string
	Description string
	Execute     func([]string) error
}

func NewCLI(cfg *Config, client *api.Client, storage *storage.FileStorage) *CLI {
	cli := &CLI{
		config:    cfg,
		apiClient: client,
		storage:   storage,
		commands:  make(map[string]Command),
	}
	
	cli.registerCommands()
	return cli
}

func (c *CLI) registerCommands() {
	c.commands["exit"] = Command{
		Name:        "exit",
		Description: "Exit the Pokedex",
		Execute:     c.exitCommand,
	}
	
	c.commands["help"] = Command{
		Name:        "help",
		Description: "Displays a help message",
		Execute:     c.helpCommand,
	}
	
	c.commands["map"] = Command{
		Name:        "map",
		Description: "Displays names of 20 location areas",
		Execute:     c.mapCommand,
	}
	
	c.commands["mapb"] = Command{
		Name:        "mapb",
		Description: "Displays previous 20 locations",
		Execute:     c.mapbCommand,
	}
	
	c.commands["explore"] = Command{
		Name:        "explore",
		Description: "Explores locations by name",
		Execute:     c.exploreCommand,
	}
	
	c.commands["catch"] = Command{
		Name:        "catch",
		Description: "Catch Pokemon and add to Pokedex",
		Execute:     c.catchCommand,
	}
	
	c.commands["inspect"] = Command{
		Name:        "inspect",
		Description: "Inspect caught Pokemon details",
		Execute:     c.inspectCommand,
	}
	
	c.commands["pokedex"] = Command{
		Name:        "pokedex",
		Description: "Show all caught Pokemon",
		Execute:     c.pokedexCommand,
	}
}

func (c *CLI) Execute(tokens []string) error {
	cmdName := tokens[0]
	
	if cmd, exists := c.commands[cmdName]; exists {
		return cmd.Execute(tokens)
	}
	
	fmt.Println("Unknown command. Type 'help' for available commands.")
	return nil
}

func (c *CLI) exitCommand(args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	c.storage.Save()
	os.Exit(0)
	return nil
}

func (c *CLI) helpCommand(args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	
	for _, cmd := range c.commands {
		fmt.Printf("%s: %s\n", cmd.Name, cmd.Description)
	}
	return nil
}

func (c *CLI) mapCommand(args []string) error {
	locations, err := c.apiClient.GetLocations(c.config.Next)
	if err != nil {
		return err
	}
	
	c.config.Next = locations.Next
	c.config.Previous = locations.Previous
	
	for _, result := range locations.Results {
		fmt.Println(result.Name)
	}
	return nil
}

func (c *CLI) mapbCommand(args []string) error {
	if c.config.Previous == "" {
		fmt.Println("You're on the first page")
		return nil
	}
	
	locations, err := c.apiClient.GetLocations(c.config.Previous)
	if err != nil {
		return err
	}
	
	c.config.Next = locations.Next
	c.config.Previous = locations.Previous
	
	for _, result := range locations.Results {
		fmt.Println(result.Name)
	}
	return nil
}

func (c *CLI) exploreCommand(args []string) error {
	if len(args) < 2 {
		fmt.Println("Please provide a location name")
		return nil
	}
	
	locationName := args[1]
	locationArea, err := c.apiClient.GetLocationArea(locationName)
	if err != nil {
		return err
	}
	
	fmt.Printf("Exploring %s...\n", locationName)
	fmt.Println("Found Pokemon:")
	
	for _, encounter := range locationArea.PokemonEncounters {
		fmt.Println("-", encounter.Pokemon.Name)
	}
	return nil
}

func (c *CLI) catchCommand(args []string) error {
	if len(args) < 2 {
		fmt.Println("Please provide a Pokemon name")
		return nil
	}
	
	pokemonName := args[1]
	poke, err := c.apiClient.GetPokemon(pokemonName)
	if err != nil {
		return err
	}
	
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)
	
	// Calculate catch chance based on base experience
	catchChance := 1.0 / (float64(poke.BaseExperience)/50 + 1)
	rand.Seed(time.Now().UnixNano())
	
	if rand.Float64() <= catchChance {
		c.storage.AddPokemon(pokemonName, *poke)
		fmt.Printf("%s was caught!\n", pokemonName)
		fmt.Println("You may now inspect it with the inspect command.")
	} else {
		fmt.Printf("%s escaped!\n", pokemonName)
	}
	
	return nil
}

func (c *CLI) inspectCommand(args []string) error {
	if len(args) < 2 {
		fmt.Println("Please provide a Pokemon name")
		return nil
	}
	
	pokemonName := args[1]
	poke, exists := c.storage.GetPokemon(pokemonName)
	if !exists {
		fmt.Println("You have not caught that Pokemon")
		return nil
	}
	
	fmt.Printf("Name: %s\n", poke.Name)
	fmt.Printf("Height: %d\n", poke.Height)
	fmt.Printf("Weight: %d\n", poke.Weight)
	fmt.Println("Stats:")
	for _, stat := range poke.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, pokemonType := range poke.Types {
		fmt.Printf("  - %s\n", pokemonType.Type.Name)
	}
	
	return nil
}

func (c *CLI) pokedexCommand(args []string) error {
	caughtPokemon := c.storage.GetAllPokemon()
	
	if len(caughtPokemon) == 0 {
		fmt.Println("Go catch some Pokemon using the catch command")
		return nil
	}
	
	fmt.Println("Your Pokedex:")
	for name := range caughtPokemon {
		fmt.Printf("  - %s\n", name)
	}
	return nil
}
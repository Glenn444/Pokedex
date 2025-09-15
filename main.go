package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Glenn444/pokedexcli/internal/api"
	"github.com/Glenn444/pokedexcli/internal/cli"
	"github.com/Glenn444/pokedexcli/internal/storage"
	"github.com/Glenn444/pokedexcli/internal/pokecache"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	
	// Initialize dependencies
	cache := pokecache.NewCache(15 * time.Second)
	apiClient := api.NewClient(cache)
	storage := storage.NewFileStorage("data/pokedex.json")
	cfg := cli.NewConfig()
	
	// Load existing caught Pokemon
	storage.Load()
	
	// Initialize CLI with dependencies
	cliApp := cli.NewCLI(cfg, apiClient, storage)
	
	for {
		fmt.Printf("Pokedex > ")
		if !scanner.Scan() {
			break
		}
		
		tokens := CleanInput(scanner.Text())
		if len(tokens) == 0 {
			continue
		}
		
		if err := cliApp.Execute(tokens); err != nil {
			fmt.Println("Error:", err)
		}
	}
}

func CleanInput(text string) []string {
	text = strings.ToLower(strings.TrimSpace(text))
	return strings.Fields(text)
}
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
  "net/http"
  "io"
  "encoding/json"
  "time"
  "math/rand"

  "internal/pokecache"
)

const interval = 15 * time.Second
var cache = pokecache.NewCache(interval)
type Pokemon struct {
	BaseExperience int `json:"base_experience"`
	Height    int `json:"height"`
	Name          string `json:"name"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}


var caughtPokemon = make(map[string]Pokemon)

type cliCommand struct{
  name string
  description string
  callback func(*config,string) error
}

type config struct{
  Next string
  Previous string
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
  cfg := &config{}
  supportedCmds := map[string]cliCommand{
    "exit":{
      name:"exit",
      description: "Exit the Pokedex",
      callback: commandExit,
    },
  }

  supportedCmds["help"] = cliCommand{
    name:"help",
    description: "Displays a help message",
    callback: func(cfg *config,locationArea string) error {return helpFunc(cfg,locationArea,supportedCmds)},
  }

  supportedCmds["map"] = cliCommand{
    name:"map",
    description: "Displays names of 20 location areas in the Pokemon world",
    callback: mapLocationsFunc,
  }

	supportedCmds["mapb"] = cliCommand{
    name:"mapb",
    description:"Displays Previous 20 locations",
    callback: mapbFunc,
  }

  supportedCmds["explore"] = cliCommand{
    name:"explore",
    description:"Explores locations and by name or ID",
    callback: exploreFunc,
  }

  supportedCmds["catch"] = cliCommand{
    name:"catch",
    description:"Catch Pokemon adds to the users Pokedex",
    callback: catchFunc,
  }

  supportedCmds["inspect"] = cliCommand{
    name:"inspect",
    description:"Inspects name of Pokemon and prints the name,height,weight,stats and types of the Pokemon",
    callback: inspectFunc,
  }

  supportedCmds["pokedex"] = cliCommand{
    name:"pokedex",
    description: "Prints all the caught Pokemon",
    callback: pokedexFunc,
  }

	for {
    fmt.Printf("Pokedex > ")
		if !scanner.Scan(){
      break
    }
			tokens := cleaninput(scanner.Text())
			if len(tokens) == 0 {
        continue
			}
      cmdName := tokens[0]
      if len(tokens) == 1{
        tokens = append(tokens, " ")
      }

      //fmt.Println("CMD: ",tokens)

      if cmd, ok := supportedCmds[cmdName]; ok{

        err := cmd.callback(cfg,tokens[1])
        if err != nil{
        tokens[1] = " "
          fmt.Println("Error:",err)
        }
      }else{
        fmt.Println("Unknown command")
      }
  }
}

func cleaninput(text string) []string {
	//implement the logic
	//1. trim space
	//2. to lowercase
	//3. split to whitespace
	text = strings.ToLower(strings.TrimSpace(text))
	return strings.Fields(text) //fields splits on any whitespace
}

func commandExit(_ *config, _ string) error{
  fmt.Printf("Closing the Pokedex... Goodbye!\n")
  os.Exit(0)
  return nil
}

func helpFunc(_ *config,_ string, supportedCmds map[string]cliCommand)error{
    fmt.Println("Welcome to the Pokedex!")
  fmt.Println("Usage:")
  for _, cmd := range supportedCmds{
    fmt.Printf("%s: %s\n", cmd.name,cmd.description)
  }
  return nil
}

func mapLocationsFunc(cfg *config,_ string) error{
  data := []byte{}
  url := cfg.Next
  if url == ""{
    url = "https://pokeapi.co/api/v2/location-area/"
  }
  if dt, ok := cache.Get(url);ok{
    data = dt
  }else{
  res, err := http.Get(url)
  if err != nil{
    return err
  }

  defer res.Body.Close()

  data,err = io.ReadAll(res.Body)
  if err != nil{
    return err
  }
  cache.Add(url,data)
}
  type PokeLocationResponse struct{
    Count int `json:"count"`
    Next string `json:"next"`
    Previous string `json:"previous"`
    Results []struct{
      Name string `json:"name"`
      URL string `json:"url"`
    }`json:"results"`
  }


  var locations PokeLocationResponse 
  
  err := json.Unmarshal(data,&locations) 
  url = locations.Next
  if err != nil{
    return err
  }

  cfg.Next = locations.Next
  cfg.Previous = locations.Previous
  
  for _,n := range locations.Results{
    fmt.Println(n.Name)
  }
   return nil
}

func mapbFunc(cfg *config,_ string)error{
  if cfg.Previous == ""{
    fmt.Println("youre on the first Page")
    return nil
  }
  data := []byte{}
if dt,ok := cache.Get(cfg.Previous);ok{
    data = dt
  }else{
  res, err := http.Get(cfg.Previous)
  if err != nil{
    return err
  }
  defer res.Body.Close()

  data, err = io.ReadAll(res.Body)
  if err != nil{
    return err
  }
  cache.Add(cfg.Previous,data)
}

  type pLocations struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}
  var previousLocations pLocations

  err := json.Unmarshal(data,&previousLocations)
  if err != nil{
    return err
  }
    cfg.Next = previousLocations.Next
  cfg.Previous = previousLocations.Previous
  for _,loc := range previousLocations.Results{
    fmt.Println(loc.Name)
  }
  return nil
}

func exploreFunc(_ *config, locationArea string)error{
  if locationArea == " "{
    return nil
  }
  locArea := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s",locationArea)
  res, err := http.Get(locArea)
  if err != nil{
  err = fmt.Errorf("Error Get: %v",err)

    return err
  }
  defer res.Body.Close()

  if res.StatusCode != http.StatusOK{
    bodyBytes,_ := io.ReadAll(res.Body)
    return fmt.Errorf(
      "Unexpected status code %d; body: %s",
      res.StatusCode,string(bodyBytes),
    )
  }
  data,err := io.ReadAll(res.Body)
  if err != nil{
  return fmt.Errorf("ReadAll: %v",err)
  }

  type pokemonLocationArea struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}
  var pokeLocationArea pokemonLocationArea

  err = json.Unmarshal(data,&pokeLocationArea)
  if err != nil{
    err = fmt.Errorf("Unmarshal: %v",err)
    return err
  }
  fmt.Printf("Exploring %s...\n", locationArea)
  fmt.Println("Found Pokemon:")
  for _,area := range pokeLocationArea.PokemonEncounters{
    fmt.Println("-",area.Pokemon.Name)
  }
  return nil
}

func catchFunc(cfg *config,pokemonName string)error{
  if pokemonName == " "{
    return nil
  }
  pokemonNameUrl := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s",pokemonName)
  res, err  := http.Get(pokemonNameUrl)
  if err != nil{
    return fmt.Errorf("Get: %v",err)
  }
  defer res.Body.Close()

  data,err := io.ReadAll(res.Body)
  if err != nil{
    return fmt.Errorf("ReadAll: ",err)
  }


  var pokeName Pokemon 

  err = json.Unmarshal(data, &pokeName)
  if err != nil{
    return fmt.Errorf("Unmarshal: ",err)
  }

  rand.Seed(time.Now().UnixNano()) //we seed because PRNG sarts with a default seed 1
  //sequece of random numbers is the same each time the program runs if not set
  catchChance := 1.0 / (float64(pokeName.BaseExperience)/50 + 1)
  r := rand.Float64()

  fmt.Printf("Throwing a Pokeball at %s...\n",pokemonName)
  if r <= catchChance{
    caughtPokemon[pokemonName] = pokeName 

    fmt.Printf("%s was caught!\n",pokemonName)
    fmt.Println("You may now inspect it with the inspect command.")
  }else{
  fmt.Printf("%s escaped!\n",pokemonName)
  }
  return nil
}

func inspectFunc(_ *config, pokemonname string)error{
  if pokemonname == " "{
    return nil
  }

  if poke,ok := caughtPokemon[pokemonname]; ok{
    fmt.Printf("Name: %s\n",poke.Name)
    fmt.Printf("Height: %d\n",poke.Height)
    fmt.Printf("Weight: %d\n",poke.Weight)
    fmt.Println("Stats:")
    for _,s := range poke.Stats{
      fmt.Printf("  -%s: %d\n",s.Stat.Name,s.BaseStat)
    }
    fmt.Println("Types:")
    for _,t := range poke.Types{
      fmt.Printf("  - %s\n",t.Type.Name)
    }
  }else{
    fmt.Println("you have not caught that pokemon")
  }
  return nil
}

func pokedexFunc(_ *config,_ string)error{

  if len(caughtPokemon) == 0{
    fmt.Println("Go catch some Pokemon using the catch command")
    return nil
  }
  fmt.Println("Your Pokedex:")
  for _,pokemon := range caughtPokemon{
    fmt.Printf("  - %s\n",pokemon.Name)
  }
  return nil
}

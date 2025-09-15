package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Glenn444/pokedexcli/internal/pokemon"
	"github.com/Glenn444/pokedexcli/internal/pokecache"
)

type Client struct {
	baseURL string
	cache   *pokecache.Cache
	client  *http.Client
}

func NewClient(cache *pokecache.Cache) *Client {
	return &Client{
		baseURL: "https://pokeapi.co/api/v2",
		cache:   cache,
		client:  &http.Client{},
	}
}

func (c *Client) GetLocations(url string) (*pokemon.LocationResponse, error) {
	if url == "" {
		url = c.baseURL + "/location-area/"
	}
	
	data, err := c.fetchWithCache(url)
	if err != nil {
		return nil, err
	}
	
	var locations pokemon.LocationResponse
	if err := json.Unmarshal(data, &locations); err != nil {
		return nil, fmt.Errorf("failed to unmarshal locations: %w", err)
	}
	
	return &locations, nil
}

func (c *Client) GetLocationArea(name string) (*pokemon.LocationArea, error) {
	url := fmt.Sprintf("%s/location-area/%s", c.baseURL, name)
	
	data, err := c.fetchWithCache(url)
	if err != nil {
		return nil, err
	}
	
	var locationArea pokemon.LocationArea
	if err := json.Unmarshal(data, &locationArea); err != nil {
		return nil, fmt.Errorf("failed to unmarshal location area: %w", err)
	}
	
	return &locationArea, nil
}

func (c *Client) GetPokemon(name string) (*pokemon.Pokemon, error) {
	url := fmt.Sprintf("%s/pokemon/%s", c.baseURL, name)
	
	data, err := c.fetchWithCache(url)
	if err != nil {
		return nil, err
	}
	
	var poke pokemon.Pokemon
	if err := json.Unmarshal(data, &poke); err != nil {
		return nil, fmt.Errorf("failed to unmarshal pokemon: %w", err)
	}
	
	return &poke, nil
}

func (c *Client) fetchWithCache(url string) ([]byte, error) {
	if data, ok := c.cache.Get(url); ok {
		return data, nil
	}
	
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", url, err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d for %s", resp.StatusCode, url)
	}
	
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	
	c.cache.Add(url, data)
	return data, nil
}
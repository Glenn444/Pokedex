package storage

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/Glenn444/pokedexcli/internal/pokemon"
)

type FileStorage struct {
	filename       string
	caughtPokemon map[string]pokemon.Pokemon
}

func NewFileStorage(filename string) *FileStorage {
	return &FileStorage{
		filename:      filename,
		caughtPokemon: make(map[string]pokemon.Pokemon),
	}
}

func (fs *FileStorage) AddPokemon(name string, poke pokemon.Pokemon) {
	fs.caughtPokemon[name] = poke
}

func (fs *FileStorage) GetPokemon(name string) (pokemon.Pokemon, bool) {
	poke, exists := fs.caughtPokemon[name]
	return poke, exists
}

func (fs *FileStorage) GetAllPokemon() map[string]pokemon.Pokemon {
	return fs.caughtPokemon
}

func (fs *FileStorage) Save() error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(fs.filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	data, err := json.MarshalIndent(fs.caughtPokemon, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(fs.filename, data, 0644)
}

func (fs *FileStorage) Load() error {
	data, err := os.ReadFile(fs.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // File doesn't exist yet, that's okay
		}
		return err
	}
	
	return json.Unmarshal(data, &fs.caughtPokemon)
}
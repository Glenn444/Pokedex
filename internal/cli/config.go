package cli

type Config struct {
	Next     string
	Previous string
}

func NewConfig() *Config {
	return &Config{}
}
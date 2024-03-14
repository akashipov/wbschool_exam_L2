package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

// Config - struct to load config from the file
type Config struct {
	Host string
	Port int
}

// LoadConfig - func to load config from the file
func LoadConfig(path string) Config {
	f, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		log.Fatalln(fmt.Errorf("problem with opening file of config: %w", err))
	}
	data, err := io.ReadAll(f)
	if err != nil {
		log.Fatalln(fmt.Errorf("problem with reading of file with config: %w", err))
	}
	c := Config{}
	err = json.Unmarshal(data, &c)
	if err != nil {
		log.Fatalln(fmt.Errorf("problem with unmarshaling of file with config: %w", err))
	}
	return c
}

package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	GorutineCount int
	Hostname      string
	Port          int
}

func NewConfig() (*Config, error) {
	gorutineCountStr := os.Getenv("GORUTINE_COUNT")
	if len(gorutineCountStr) == 0 {
		return nil, fmt.Errorf("environment variable GORUTINE_COUNT not specified")
	}

	gorutineCount, err := strconv.Atoi(gorutineCountStr)
	if err != nil {
		return nil, fmt.Errorf("couldn't convert GORUTINE_COUNT in the number: %v", err)
	}

	if gorutineCount < 0 {
		return nil, fmt.Errorf("GORUTINE_COUNT must be > 0")
	}

	flag.Parse()

	hostname := flag.String("h", "localhost", "orchestrator adress")
	port := flag.Int("p", 8081, "orchestrator port")

	if *port < 0 {
		return nil, fmt.Errorf("port must be > 0")
	}

	agcfg := Config{
		GorutineCount: gorutineCount,
		Hostname:      *hostname,
		Port:          *port,
	}

	return &agcfg, nil
}

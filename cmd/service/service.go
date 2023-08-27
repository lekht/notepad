package main

import (
	"flag"
	"log"

	"github.com/lekht/notepad/config"
	"github.com/lekht/notepad/internal/app"
)

var cfg config.Config

func init() {
	path := flag.String("config", "", "path to config file")
	flag.Parse()

	err := cfg.LoadEnv(*path)
	if err != nil {
		log.Fatalf("failed to load environment variables: %s", err)
	}
}

func main() {
	app.Run(&cfg)
}

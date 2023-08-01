package main

import (
	"fmt"
	"os"
)

var ConfigFile = "config.yml"

func main() {
	cfg := ReadConfig(ConfigFile)
	p := NewProducer(cfg)
	defer p.Close()

	server := InitServer(p)

	err := server.Run(":8080")
	if err != nil {
		fmt.Printf("Failed to send event: %s\n", err)
		os.Exit(1)
	}
}

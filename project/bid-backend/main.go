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

	dto := BidDto{
		ProductId: "bb",
		Value:     "1.234",
		BidderId:  "aa",
	}

	err := p.Send(dto)
	if err != nil {
		fmt.Printf("Failed to send event: %s\n", err)
		os.Exit(1)
	}
}

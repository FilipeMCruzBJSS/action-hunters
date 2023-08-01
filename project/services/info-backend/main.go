package main

var ConfigFile = "config.yml"

func main() {
	cfg := ReadConfig(ConfigFile)
	p := NewConsumer(cfg)
	defer p.Close()

	p.Subscribe()
}

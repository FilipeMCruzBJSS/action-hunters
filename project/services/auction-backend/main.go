package main

var ConfigFile = "config.yml"

func main() {
	cfg := ReadConfig(ConfigFile)

	s := NewStore()

	c := NewConsumer(cfg, s)
	defer c.Close()

	c.Subscribe()
}

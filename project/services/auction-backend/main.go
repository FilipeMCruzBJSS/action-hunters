package main

var ConfigFile = "config.yml"

func main() {
	cfg := ReadConfig(ConfigFile)

	updateChannel := make(chan AuctionDto)

	s := NewStore(updateChannel)

	p := NewProducer(cfg, updateChannel)
	defer p.Close()

	c := NewConsumer(cfg, s)
	defer c.Close()

	c.Subscribe()

	close(updateChannel)
}

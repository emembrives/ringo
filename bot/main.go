package main

import (
	"flag"
	"io/ioutil"
	"log"

  "gopkg.in/yaml.v2"
	"github.com/emembrives/ringo/bot/telegram"
)

var (
  configPath = flag.String("config", "config.yaml",
		"Path to the YAML configuration file")
)

type config struct {
	TelegramBot telegram.Config
}

func main() {
	flag.Parse()

	configData, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.Fatalf("Unable to read configuration file: %v", err)
	}

	c := &config{}
	yaml.Unmarshal(configData, c)

	bot := telegram.NewBot(c.TelegramBot)
	go bot.Run()
}

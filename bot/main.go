package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/emembrives/ringo/bot/telegram"
	"gopkg.in/yaml.v2"
)

var (
	configPath = flag.String("config", "config.yaml",
		"Path to the YAML configuration file")
)

type config struct {
	TelegramBot telegram.Config `yaml:"telegram-bot"`
	Ringer      RingerConfig    `yaml:"ringer"`
}

func main() {
	flag.Parse()

	configData, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.Fatalf("Unable to read configuration file: %v", err)
	}

	c := &config{}
	err = yaml.Unmarshal(configData, c)
	if err != nil {
		log.Fatalf("Unable to parse configuration file: %v", err)
	}

	bot := telegram.NewWebhookBot(c.TelegramBot, &Processor{
		Config: c.Ringer,
	})
	bot.Run()
}

package main

import (
	"log"
	"time"

	"github.com/emembrives/ringo/bot/telegram"

	zmq "github.com/pebbe/zmq4"
)

const (
	zmqURL = "tcp://0.0.0.0:9101"
)

// RingerConfig holds the configuration for the ringer comm server.
type RingerConfig struct {
	ChatIDs []int64 `yaml:"chat-ids"`
}

// Processor is the brain of the telegram bot.
type Processor struct {
	Config      RingerConfig
	botIncoming <-chan telegram.Update
	botOutgoing chan<- telegram.SendMessage
}

// SetIncomingChannel sets the incoming message channel.
func (p *Processor) SetIncomingChannel(message <-chan telegram.Update) {
	p.botIncoming = message
}

// SetOutgoingChannel sets the outgoing message channel
func (p *Processor) SetOutgoingChannel(message chan<- telegram.SendMessage) {
	p.botOutgoing = message
}

// Run runs the bot logic.
func (p *Processor) Run() {
	go p.receiveMessages()
	for {
		_ = <-p.botIncoming
		// Maybe do something with that
	}
}

func (p *Processor) receiveMessages() {
	responder, err := zmq.NewSocket(zmq.REP)
	defer responder.Close()
	if err != nil {
		log.Fatalf("Unable to open ZMQ socket: %v", err)
	}

	err = responder.Bind(zmqURL)
	if err != nil {
		log.Fatalf("Unable to bind to port: %v", err)
	}

	for {
		msg, err := responder.RecvBytes(0)
		if err != nil {
			log.Printf("Error while receiving ZMQ request: %v", err)
			responder.SendBytes([]byte("err"), 0)
			continue
		}

		msgStr := string(msg)
		log.Println(msgStr)

		if msgStr == "t" {
			// Return the current time.
			dateStr := time.Now().Format(time.RFC3339)
			responder.SendBytes([]byte(dateStr), 0)
		} else if msgStr == "r" {
			for _, chatID := range p.Config.ChatIDs {
				m := telegram.SendMessage{
					ChatID: chatID,
					Text:   "Etienne a besoin d'aide.",
				}
				p.botOutgoing <- m
			}
			responder.SendBytes([]byte("ok"), 0)
		}
	}
}

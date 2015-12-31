package telegram

import (
	"log"
)

const (
	webhookPath = "/webhook/%s/"
	telegramAPI = "https://api.telegram.org/bot%s/%s"
)

func dieOnError(err error, str string) {
	if err != nil {
		log.Fatalf(str, err)
	}
}

// Bot is an implementation of a Telegram bot.
type Bot struct {
	server *WebhookServer
}

// NewBot creates a new telegram bot.
func NewBot(config Config) *Bot {
	tb := &Bot{}
	tb.server = NewWebhookServer(config)
	return tb
}

// Run runs the bot. This method is long-running.
func (tb *Bot) Run() {
	tb.server.Run()
}

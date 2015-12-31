package telegram

// BotLogic implements the logic of the bot
type BotLogic interface {
	SetIncomingChannel(message <-chan Update)
	SetOutgoingChannel(message chan<- SendMessage)

	Run()
}

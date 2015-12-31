package telegram

// Config holds the configuration for a Telegram bot.
type Config struct {
	Token    string `yaml:"token"`
	BasePath string `yaml:"base-url"`
}

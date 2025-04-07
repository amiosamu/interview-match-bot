package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// handleStartCommand processes the /start command
func (b *Bot) handleStartCommand(message *tgbotapi.Message) {
	welcomeText := "Welcome to Interview Match Bot! Please select your field of interest:"
	b.sendMessage(message.Chat.ID, welcomeText, CreateCategoriesKeyboard())
}

// handleHelpCommand processes the /help command
func (b *Bot) handleHelpCommand(message *tgbotapi.Message) {
	helpText := `*Interview Match Bot Help*
	
This bot helps you find partners for interview practice.

*Commands:*
/start - Start the bot and select your category
/help - Show this help message

*How to use:*
1. Select your field of interest (e.g., Backend, Frontend)
2. Select your experience level (e.g., Junior, Middle)
3. The bot will match you with others looking for the same criteria

If you have any issues, please try restarting the bot with /start.`

	msg := tgbotapi.NewMessage(message.Chat.ID, helpText)
	msg.ParseMode = "Markdown"
	b.api.Send(msg)
}

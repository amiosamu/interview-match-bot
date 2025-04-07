package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// handleMessage processes incoming messages
func (b *Bot) handleMessage(message *tgbotapi.Message) {
	// Save or update user information
	b.saveUserInfo(message.From)
	
	// Handle commands
	if message.IsCommand() {
		switch message.Command() {
		case "start":
			b.handleStartCommand(message)
		case "help":
			b.handleHelpCommand(message)
		default:
			b.sendMessage(message.Chat.ID, "Unknown command. Type /start to begin or /help for assistance.", nil)
		}
		return
	}
	
	// For non-command messages, just prompt the user to use the commands
	b.sendMessage(message.Chat.ID, "Please use the buttons or type /start to begin.", nil)
}

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

// handleCallbackQuery processes button presses
func (b *Bot) handleCallbackQuery(query *tgbotapi.CallbackQuery) {
	// Always answer the callback query
	callback := tgbotapi.NewCallback(query.ID, "")
	b.api.Request(callback)
	
	user := b.saveUserInfo(query.From)
	
	// Parse the callback data
	data := query.Data
	
	if len(data) > 9 && data[:9] == "category:" {
		category := data[9:]
		user.Field = category
		b.userStore.SaveUser(user)
		
		// Now ask for the level
		msg := tgbotapi.NewMessage(query.Message.Chat.ID, "Great! Now select your experience level:")
		msg.ReplyMarkup = CreateLevelsKeyboard()
		b.api.Send(msg)
		
	} else if len(data) > 6 && data[:6] == "level:" {
		level := data[6:]
		user.Level = level
		b.userStore.SaveUser(user)
		
		// Confirm the selection
		confirmMessage := "Perfect! I'll notify you when I find someone matching your criteria.\n\n" +
			"You selected: " + user.Field + " - " + user.Level
		
		b.sendMessage(query.Message.Chat.ID, confirmMessage, nil)
		
		// Find matches
		b.notifyMatches(user)
	}
}
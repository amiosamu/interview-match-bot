package bot

import (
	"database/sql"
	"log"
	"strings"

	"github.com/amiosamu/interview-match-bot/internal/models"
	"github.com/amiosamu/interview-match-bot/internal/service"
	"github.com/amiosamu/interview-match-bot/internal/store"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Bot represents the interview bot application
type Bot struct {
	api         *tgbotapi.BotAPI
	db          *sql.DB
	userStore   *store.UserStore
	quizService *service.QuizService
	// These services would be added when implementing other features
	// analyticsService  *service.AnalyticsService
	// moderationService *service.ModerationService
}

// NewBot creates a new Bot instance
func NewBot(token string, db *sql.DB) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		api:         api,
		db:          db,
		userStore:   store.NewUserStore(),
		quizService: service.NewQuizService(db),
	}, nil
}

// Start begins listening for updates
func (b *Bot) Start() {
	log.Printf("Authorized on account %s", b.api.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			b.handleMessage(update.Message)
		} else if update.CallbackQuery != nil {
			b.handleCallbackQuery(update.CallbackQuery)
		}
	}
}

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
		case "prepare":
			b.handlePrepareCommand(message)
		default:
			b.sendMessage(message.Chat.ID, "Unknown command. Type /start to begin or /help for assistance.", nil)
		}
		return
	}

	// For non-command messages, just prompt the user to use the commands
	b.sendMessage(message.Chat.ID, "Please use the buttons or type /start to begin.", nil)
}

// handleCallbackQuery processes button presses
func (b *Bot) handleCallbackQuery(query *tgbotapi.CallbackQuery) {
	// Always answer the callback query to stop the loading indicator
	callback := tgbotapi.NewCallback(query.ID, "")
	b.api.Request(callback)

	user := b.saveUserInfo(query.From)

	// Parse the callback data
	data := query.Data

	if strings.HasPrefix(data, "category:") {
		category := strings.TrimPrefix(data, "category:")
		user.Field = category
		b.userStore.SaveUser(user)

		// Now ask for the level
		msg := tgbotapi.NewMessage(query.Message.Chat.ID, "Great! Now select your experience level:")
		msg.ReplyMarkup = CreateLevelsKeyboard()
		b.api.Send(msg)

	} else if strings.HasPrefix(data, "level:") {
		level := strings.TrimPrefix(data, "level:")
		user.Level = level
		b.userStore.SaveUser(user)

		// Confirm the selection
		confirmMessage := "Perfect! I'll notify you when I find someone matching your criteria.\n\n" +
			"You selected: " + user.Field + " - " + user.Level

		b.sendMessage(query.Message.Chat.ID, confirmMessage, nil)

		// Find matches
		b.notifyMatches(user)
	} else if strings.HasPrefix(data, "quiz:") {
		// Handle quiz-related callbacks
		b.handleQuizCallback(query)
	} else if strings.HasPrefix(data, "main:") {
		// Handle main menu callbacks
		if data == "main:menu" {
			welcomeText := "Welcome to Interview Match Bot! What would you like to do?"
			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Find Interview Partner", "main:find_partner"),
					tgbotapi.NewInlineKeyboardButtonData("Practice Questions", "quiz:new"),
				),
			)
			b.sendMessage(query.Message.Chat.ID, welcomeText, keyboard)
		} else if data == "main:find_partner" {
			b.sendMessage(query.Message.Chat.ID, "Please select your field of interest:", CreateCategoriesKeyboard())
		}
	}
}

// saveUserInfo stores or updates user information
func (b *Bot) saveUserInfo(tgUser *tgbotapi.User) *models.User {
	user, exists := b.userStore.GetUser(tgUser.ID)

	if !exists {
		user = &models.User{
			ID:        tgUser.ID,
			Username:  tgUser.UserName,
			FirstName: tgUser.FirstName,
			LastName:  tgUser.LastName,
		}
		b.userStore.SaveUser(user)
	}

	return user
}

// sendMessage sends a message to a chat
func (b *Bot) sendMessage(chatID int64, text string, markup interface{}) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = markup
	b.api.Send(msg)
}

// notifyMatches notifies users about matches
func (b *Bot) notifyMatches(user *models.User) {
	matches := b.userStore.FindMatches(user.ID, user.Field, user.Level)

	if len(matches) == 0 {
		return
	}

	// Notify the current user about all matches
	for _, match := range matches {
		messageText := "I found a match! User " + match.DisplayName() + " is also looking for " +
			user.Field + " " + user.Level + " positions."
		b.sendMessage(user.ID, messageText, nil)

		// Also notify the matched user
		matchMessageText := "I found a match! User " + user.DisplayName() + " is also looking for " +
			user.Field + " " + user.Level + " positions."
		b.sendMessage(match.ID, matchMessageText, nil)
	}
}

// Debug enables debug mode for the bot
func (b *Bot) Debug(enable bool) {
	b.api.Debug = enable
}

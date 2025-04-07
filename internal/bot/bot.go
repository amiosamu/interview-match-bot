package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/amiosamu/interview-match-bot/internal/models"
	"github.com/amiosamu/interview-match-bot/internal/store"
)

// Bot represents the interview bot application
type Bot struct {
	api       *tgbotapi.BotAPI
	userStore *store.UserStore
}

// NewBot creates a new Bot instance
func NewBot(token string) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		api:       api,
		userStore: store.NewUserStore(),
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
package bot

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/amiosamu/interview-match-bot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// handlePrepareCommand initiates the interview preparation quiz flow
func (b *Bot) handlePrepareCommand(message *tgbotapi.Message) {
	user := b.saveUserInfo(message.From)

	// Check if the user already has an active quiz session
	session, err := b.quizService.GetActiveQuizSession(user.ID)
	if err != nil {
		log.Printf("Error retrieving active quiz session: %v", err)
		b.sendMessage(message.Chat.ID, "Sorry, I encountered an error. Please try again later.", nil)
		return
	}

	if session != nil {
		// User has an active quiz, ask if they want to continue or start a new one
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Continue Current Quiz", "quiz:continue"),
				tgbotapi.NewInlineKeyboardButtonData("Start New Quiz", "quiz:new"),
			),
		)

		b.sendMessage(message.Chat.ID, fmt.Sprintf("You have an unfinished %s quiz. Would you like to continue or start a new one?", session.Language), keyboard)
		return
	}

	// No active quiz, show language selection
	b.sendQuizLanguageSelection(message.Chat.ID)
}

// sendQuizLanguageSelection shows available programming languages for quizzes
func (b *Bot) sendQuizLanguageSelection(chatID int64) {
	// Get available languages or use a predefined list if the database query fails
	languages, err := b.quizService.GetQuizLanguages()
	if err != nil || len(languages) == 0 {
		// Fallback to hardcoded languages
		languages = []string{"golang", "python", "javascript", "java", "csharp", "ruby"}
	}

	// Create a keyboard with available languages
	var rows [][]tgbotapi.InlineKeyboardButton

	// Group languages into rows of 2 buttons each
	for i := 0; i < len(languages); i += 2 {
		var row []tgbotapi.InlineKeyboardButton

		// Add the first language in this row
		buttonText := formatLanguageName(languages[i])
		row = append(row, tgbotapi.NewInlineKeyboardButtonData(buttonText, "quiz:lang:"+languages[i]))

		// Add the second language if there is one
		if i+1 < len(languages) {
			buttonText = formatLanguageName(languages[i+1])
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(buttonText, "quiz:lang:"+languages[i+1]))
		}

		rows = append(rows, row)
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	b.sendMessage(chatID, "Choose a programming language for your interview preparation quiz:", keyboard)
}

// formatLanguageName converts language code to a nice display name
func formatLanguageName(lang string) string {
	switch strings.ToLower(lang) {
	case "golang":
		return "Go"
	case "javascript":
		return "JavaScript"
	case "csharp":
		return "C#"
	case "cpp":
		return "C++"
	default:
		// Capitalize the first letter
		if len(lang) > 0 {
			return strings.ToUpper(lang[:1]) + lang[1:]
		}
		return lang
	}
}

// handleQuizCallback processes quiz-related button actions
func (b *Bot) handleQuizCallback(query *tgbotapi.CallbackQuery) {
	// Always acknowledge the callback query to stop the loading indicator
	callback := tgbotapi.NewCallback(query.ID, "")
	b.api.Request(callback)

	data := query.Data
	user := b.saveUserInfo(query.From)

	// Handle different quiz callbacks
	if data == "quiz:continue" {
		b.continueQuiz(query.Message.Chat.ID, user.ID)
		return
	}

	if data == "quiz:new" {
		// End any active session and show language selection
		session, err := b.quizService.GetActiveQuizSession(user.ID)
		if err == nil && session != nil {
			b.quizService.CompleteQuizSession(session.ID)
		}

		b.sendQuizLanguageSelection(query.Message.Chat.ID)
		return
	}

	if strings.HasPrefix(data, "quiz:lang:") {
		// Extract the language
		language := strings.TrimPrefix(data, "quiz:lang:")
		b.startNewQuiz(query.Message.Chat.ID, user.ID, language)
		return
	}

	if strings.HasPrefix(data, "quiz:answer:") {
		// Extract session ID and answer index
		parts := strings.Split(data, ":")
		if len(parts) != 4 {
			b.sendMessage(query.Message.Chat.ID, "Invalid option. Please try again.", nil)
			return
		}

		sessionID, err := strconv.Atoi(parts[2])
		if err != nil {
			b.sendMessage(query.Message.Chat.ID, "Invalid session. Please try again.", nil)
			return
		}

		answerIndex, err := strconv.Atoi(parts[3])
		if err != nil {
			b.sendMessage(query.Message.Chat.ID, "Invalid answer. Please try again.", nil)
			return
		}

		b.processQuizAnswer(query.Message.Chat.ID, user.ID, sessionID, answerIndex)
		return
	}
}

// startNewQuiz begins a new quiz session for a user
func (b *Bot) startNewQuiz(chatID int64, userID int64, language string) {
	// Get random questions for this language (10 questions per quiz)
	questions, err := b.quizService.GetQuestionsByLanguage(language, 10)
	if err != nil {
		log.Printf("Error fetching questions: %v", err)
		b.sendMessage(chatID, "Sorry, I couldn't start a quiz for this language. Please try again later.", nil)
		return
	}

	if len(questions) == 0 {
		b.sendMessage(chatID, fmt.Sprintf("Sorry, no questions are available for %s yet. Please try another language.", formatLanguageName(language)), nil)
		return
	}

	// Extract question IDs
	var questionIDs []int
	for _, q := range questions {
		questionIDs = append(questionIDs, q.ID)
	}

	// Create a new quiz session
	session, err := b.quizService.CreateQuizSession(userID, language, questionIDs)
	if err != nil {
		log.Printf("Error creating quiz session: %v", err)
		b.sendMessage(chatID, "Sorry, I couldn't start the quiz. Please try again later.", nil)
		return
	}

	// Send introduction message
	b.sendMessage(chatID, fmt.Sprintf("Starting a new %s quiz with %d questions. Let's begin!", formatLanguageName(language), len(questionIDs)), nil)

	// Wait a moment before sending the first question
	time.Sleep(1 * time.Second)

	// Send the first question
	b.sendQuizQuestion(chatID, userID, session)
}

// continueQuiz resumes an existing quiz session
func (b *Bot) continueQuiz(chatID int64, userID int64) {
	session, err := b.quizService.GetActiveQuizSession(userID)
	if err != nil {
		log.Printf("Error retrieving active session: %v", err)
		b.sendMessage(chatID, "Sorry, I couldn't retrieve your quiz. Please try starting a new one.", nil)
		return
	}

	if session == nil {
		b.sendMessage(chatID, "You don't have an active quiz. Please start a new one.", nil)
		b.sendQuizLanguageSelection(chatID)
		return
	}

	// Send the current question
	b.sendQuizQuestion(chatID, userID, session)
}

// sendQuizQuestion sends the current question to the user
func (b *Bot) sendQuizQuestion(chatID int64, userID int64, session *models.QuizSession) {
	// Check if the quiz is complete
	if session.IsComplete() {
		b.completeQuiz(chatID, userID, session)
		return
	}

	// Get the current question
	questionID := session.QuestionIDs[session.CurrentQuestionIndex]
	question, err := b.quizService.GetQuestionByID(questionID)
	if err != nil {
		log.Printf("Error fetching question %d: %v", questionID, err)
		b.sendMessage(chatID, "Sorry, I couldn't retrieve the question. Please try again later.", nil)
		return
	}

	// Create the question message
	questionNumber := session.CurrentQuestionIndex + 1
	totalQuestions := len(session.QuestionIDs)
	messageText := fmt.Sprintf("*Question %d of %d*\n\n%s", questionNumber, totalQuestions, question.QuestionText)

	// Create answer buttons
	var rows [][]tgbotapi.InlineKeyboardButton
	for i, option := range question.AnswerOptions {
		callbackData := fmt.Sprintf("quiz:answer:%d:%d", session.ID, i)
		rows = append(rows, []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(option, callbackData),
		})
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	// Send the message with the keyboard
	msg := tgbotapi.NewMessage(chatID, messageText)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	b.api.Send(msg)
}

// processQuizAnswer handles a user's answer to a quiz question
func (b *Bot) processQuizAnswer(chatID int64, userID int64, sessionID int, answerIndex int) {
	// Get the active session
	session, err := b.quizService.GetActiveQuizSession(userID)
	if err != nil {
		log.Printf("Error retrieving active session: %v", err)
		b.sendMessage(chatID, "Sorry, I couldn't process your answer. Please try again later.", nil)
		return
	}

	if session == nil || session.ID != sessionID {
		b.sendMessage(chatID, "This quiz is no longer active. Please start a new one.", nil)
		return
	}

	// Get the current question
	questionID := session.QuestionIDs[session.CurrentQuestionIndex]
	question, err := b.quizService.GetQuestionByID(questionID)
	if err != nil {
		log.Printf("Error fetching question %d: %v", questionID, err)
		b.sendMessage(chatID, "Sorry, I couldn't retrieve the question. Please try again later.", nil)
		return
	}

	// Check if the answer index is valid
	if answerIndex < 0 || answerIndex >= len(question.AnswerOptions) {
		b.sendMessage(chatID, "Invalid answer selection. Please try again.", nil)
		return
	}

	// Get the selected answer
	selectedAnswer := question.AnswerOptions[answerIndex]

	// Check if the answer is correct
	isCorrect := selectedAnswer == question.CorrectAnswer

	// Record the answer
	err = b.quizService.RecordAnswer(userID, session.ID, questionID, selectedAnswer, isCorrect)
	if err != nil {
		log.Printf("Error recording answer: %v", err)
		// Continue anyway - this isn't critical
	}

	// Prepare feedback message
	var feedbackMessage string
	if isCorrect {
		feedbackMessage = "✅ *Correct!*\n\n"
		// Update the session's correct answers count locally as well
		session.CorrectAnswers++
	} else {
		feedbackMessage = fmt.Sprintf("❌ *Incorrect*\n\nThe correct answer is: *%s*\n\n", question.CorrectAnswer)
	}

	// Add explanation
	feedbackMessage += "*Explanation:*\n" + question.Explanation

	// Send feedback
	msg := tgbotapi.NewMessage(chatID, feedbackMessage)
	msg.ParseMode = "Markdown"
	b.api.Send(msg)

	// Advance to the next question
	err = b.quizService.AdvanceQuizSession(session.ID)
	if err != nil {
		log.Printf("Error advancing session: %v", err)
		b.sendMessage(chatID, "Sorry, I couldn't advance to the next question. Please try again later.", nil)
		return
	}

	// Update session locally
	session.CurrentQuestionIndex++

	// Wait before showing the next question
	time.Sleep(3 * time.Second)

	// Send the next question or complete the quiz
	if session.CurrentQuestionIndex >= len(session.QuestionIDs) {
		b.completeQuiz(chatID, userID, session)
	} else {
		b.sendQuizQuestion(chatID, userID, session)
	}
}

// completeQuiz finishes a quiz session and shows results
func (b *Bot) completeQuiz(chatID int64, userID int64, session *models.QuizSession) {
	// Mark the session as complete in the database
	if session.CompletedAt == nil {
		err := b.quizService.CompleteQuizSession(session.ID)
		if err != nil {
			log.Printf("Error completing quiz session: %v", err)
			// Continue anyway - the user should still see their results
		}

		// Set completed time locally
		now := time.Now()
		session.CompletedAt = &now
	}

	// Calculate score percentage
	score := session.GetScore()

	// Prepare results message
	var resultsMessage string

	if score >= 90 {
		resultsMessage = "🏆 *Excellent job!* "
	} else if score >= 70 {
		resultsMessage = "👍 *Good work!* "
	} else if score >= 50 {
		resultsMessage = "👌 *Not bad!* "
	} else {
		resultsMessage = "📚 *Keep practicing!* "
	}

	resultsMessage += fmt.Sprintf("You completed the %s quiz.\n\n", formatLanguageName(session.Language))
	resultsMessage += fmt.Sprintf("*Your score: %.1f%%* (%d correct out of %d questions)\n\n",
		score, session.CorrectAnswers, session.CurrentQuestionIndex)

	// Add recommendations based on score
	if score < 60 {
		resultsMessage += "It looks like you might benefit from studying this topic more. Would you like to try another quiz or return to the main menu?"
	} else if score < 80 {
		resultsMessage += "You have a good understanding, but there's room for improvement. Would you like to try another quiz or return to the main menu?"
	} else {
		resultsMessage += "Great job! You have a strong understanding of this topic. Would you like to try another quiz or return to the main menu?"
	}

	// Create keyboard with options
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Take Another Quiz", "quiz:new"),
			tgbotapi.NewInlineKeyboardButtonData("Main Menu", "main:menu"),
		),
	)

	// Send results
	msg := tgbotapi.NewMessage(chatID, resultsMessage)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	b.api.Send(msg)
}

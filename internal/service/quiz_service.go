package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/amiosamu/interview-match-bot/internal/models"
)

// QuizService handles quiz-related operations
type QuizService struct {
	db *sql.DB
}

// NewQuizService creates a new QuizService
func NewQuizService(db *sql.DB) *QuizService {
	return &QuizService{db: db}
}

// GetQuestionsByLanguage retrieves random questions for a specific language
func (s *QuizService) GetQuestionsByLanguage(language string, limit int) ([]*models.QuizQuestion, error) {
	// Query randomly selected questions for the specified language
	rows, err := s.db.Query(`
		SELECT id, language, category, difficulty, question_text, answer_options, correct_answer, explanation, created_at 
		FROM quiz_questions 
		WHERE language = $1 
		ORDER BY RANDOM() 
		LIMIT $2
	`, language, limit)
	
	if err != nil {
		return nil, fmt.Errorf("error querying questions: %w", err)
	}
	defer rows.Close()
	
	var questions []*models.QuizQuestion
	for rows.Next() {
		q := &models.QuizQuestion{}
		var answerOptionsJSON string
		
		err := rows.Scan(
			&q.ID, 
			&q.Language, 
			&q.Category, 
			&q.Difficulty, 
			&q.QuestionText, 
			&answerOptionsJSON, 
			&q.CorrectAnswer, 
			&q.Explanation,
			&q.CreatedAt,
		)
		
		if err != nil {
			return nil, fmt.Errorf("error scanning question row: %w", err)
		}
		
		// Parse the JSON array of answer options
		err = json.Unmarshal([]byte(answerOptionsJSON), &q.AnswerOptions)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling answer options: %w", err)
		}
		
		questions = append(questions, q)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating question rows: %w", err)
	}
	
	return questions, nil
}

// GetQuestionByID retrieves a specific question by ID
func (s *QuizService) GetQuestionByID(questionID int) (*models.QuizQuestion, error) {
	q := &models.QuizQuestion{}
	var answerOptionsJSON string
	
	err := s.db.QueryRow(`
		SELECT id, language, category, difficulty, question_text, answer_options, correct_answer, explanation, created_at 
		FROM quiz_questions 
		WHERE id = $1
	`, questionID).Scan(
		&q.ID, 
		&q.Language, 
		&q.Category, 
		&q.Difficulty, 
		&q.QuestionText, 
		&answerOptionsJSON, 
		&q.CorrectAnswer, 
		&q.Explanation,
		&q.CreatedAt,
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("question not found: %w", err)
		}
		return nil, fmt.Errorf("error querying question: %w", err)
	}
	
	// Parse the JSON array of answer options
	err = json.Unmarshal([]byte(answerOptionsJSON), &q.AnswerOptions)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling answer options: %w", err)
	}
	
	return q, nil
}

// CreateQuizSession starts a new quiz session for a user
func (s *QuizService) CreateQuizSession(userID int64, language string, questionIDs []int) (*models.QuizSession, error) {
	// Convert question IDs to JSON
	questionIDsJSON, err := json.Marshal(questionIDs)
	if err != nil {
		return nil, fmt.Errorf("error marshaling question IDs: %w", err)
	}
	
	// Create the session in the database
	var sessionID int
	err = s.db.QueryRow(`
		INSERT INTO user_quiz_sessions (user_id, language, question_ids)
		VALUES ($1, $2, $3)
		RETURNING id, started_at
	`, userID, language, questionIDsJSON).Scan(&sessionID, &time.Time{})
	
	if err != nil {
		return nil, fmt.Errorf("error creating quiz session: %w", err)
	}
	
	// Return the new session
	return &models.QuizSession{
		ID:                  sessionID,
		UserID:              userID,
		Language:            language,
		CurrentQuestionIndex: 0,
		QuestionIDs:         questionIDs,
		CorrectAnswers:      0,
		StartedAt:           time.Now(),
	}, nil
}

// GetActiveQuizSession retrieves the active quiz session for a user
func (s *QuizService) GetActiveQuizSession(userID int64) (*models.QuizSession, error) {
	var session models.QuizSession
	var questionIDsJSON string
	var completedAt sql.NullTime
	
	err := s.db.QueryRow(`
		SELECT id, user_id, language, current_question_index, question_ids, correct_answers, started_at, completed_at
		FROM user_quiz_sessions
		WHERE user_id = $1 AND completed_at IS NULL
		ORDER BY started_at DESC
		LIMIT 1
	`, userID).Scan(
		&session.ID,
		&session.UserID,
		&session.Language,
		&session.CurrentQuestionIndex,
		&questionIDsJSON,
		&session.CorrectAnswers,
		&session.StartedAt,
		&completedAt,
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No active session
		}
		return nil, fmt.Errorf("error querying active session: %w", err)
	}
	
	// Parse the JSON array of question IDs
	err = json.Unmarshal([]byte(questionIDsJSON), &session.QuestionIDs)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling question IDs: %w", err)
	}
	
	if completedAt.Valid {
		session.CompletedAt = &completedAt.Time
	}
	
	return &session, nil
}

// RecordAnswer records a user's answer to a question
func (s *QuizService) RecordAnswer(userID int64, sessionID int, questionID int, answerGiven string, isCorrect bool) error {
	_, err := s.db.Exec(`
		INSERT INTO user_quiz_answers (user_id, session_id, question_id, answer_given, is_correct)
		VALUES ($1, $2, $3, $4, $5)
	`, userID, sessionID, questionID, answerGiven, isCorrect)
	
	if err != nil {
		return fmt.Errorf("error recording answer: %w", err)
	}
	
	// If the answer is correct, increment the correct_answers count
	if isCorrect {
		_, err = s.db.Exec(`
			UPDATE user_quiz_sessions
			SET correct_answers = correct_answers + 1
			WHERE id = $1
		`, sessionID)
		
		if err != nil {
			return fmt.Errorf("error updating correct answers: %w", err)
		}
	}
	
	return nil
}

// AdvanceQuizSession moves to the next question in a session
func (s *QuizService) AdvanceQuizSession(sessionID int) error {
	_, err := s.db.Exec(`
		UPDATE user_quiz_sessions
		SET current_question_index = current_question_index + 1
		WHERE id = $1
	`, sessionID)
	
	if err != nil {
		return fmt.Errorf("error advancing quiz session: %w", err)
	}
	
	return nil
}

// CompleteQuizSession marks a quiz session as completed
func (s *QuizService) CompleteQuizSession(sessionID int) error {
	_, err := s.db.Exec(`
		UPDATE user_quiz_sessions
		SET completed_at = NOW()
		WHERE id = $1
	`, sessionID)
	
	if err != nil {
		return fmt.Errorf("error completing quiz session: %w", err)
	}
	
	return nil
}

// GetQuizLanguages returns a list of available quiz languages
func (s *QuizService) GetQuizLanguages() ([]string, error) {
	rows, err := s.db.Query(`
		SELECT DISTINCT language
		FROM quiz_questions
		ORDER BY language
	`)
	
	if err != nil {
		return nil, fmt.Errorf("error querying languages: %w", err)
	}
	defer rows.Close()
	
	var languages []string
	for rows.Next() {
		var language string
		if err := rows.Scan(&language); err != nil {
			return nil, fmt.Errorf("error scanning language: %w", err)
		}
		languages = append(languages, language)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating language rows: %w", err)
	}
	
	return languages, nil
}

// GetUserQuizStats gets statistics about a user's quiz performance
func (s *QuizService) GetUserQuizStats(userID int64) (map[string]map[string]int, error) {
	rows, err := s.db.Query(`
		SELECT 
			q.language,
			COUNT(DISTINCT s.id) AS completed_quizzes,
			COUNT(a.id) AS total_questions,
			SUM(CASE WHEN a.is_correct THEN 1 ELSE 0 END) AS correct_answers
		FROM user_quiz_sessions s
		JOIN user_quiz_answers a ON s.id = a.session_id
		JOIN quiz_questions q ON a.question_id = q.id
		WHERE s.user_id = $1 AND s.completed_at IS NOT NULL
		GROUP BY q.language
	`, userID)
	
	if err != nil {
		return nil, fmt.Errorf("error querying user stats: %w", err)
	}
	defer rows.Close()
	
	stats := make(map[string]map[string]int)
	for rows.Next() {
		var language string
		var completedQuizzes, totalQuestions, correctAnswers int
		
		err := rows.Scan(&language, &completedQuizzes, &totalQuestions, &correctAnswers)
		if err != nil {
			return nil, fmt.Errorf("error scanning stats row: %w", err)
		}
		
		stats[language] = map[string]int{
			"completed_quizzes": completedQuizzes,
			"total_questions":   totalQuestions,
			"correct_answers":   correctAnswers,
		}
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating stats rows: %w", err)
	}
	
	return stats, nil
}
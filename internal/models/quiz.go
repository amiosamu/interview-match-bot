package models

import (
	"encoding/json"
	"time"
)

// QuizQuestion represents a single quiz question
type QuizQuestion struct {
	ID            int       `json:"id"`
	Language      string    `json:"language"`
	Category      string    `json:"category"`
	Difficulty    string    `json:"difficulty"`
	QuestionText  string    `json:"question_text"`
	AnswerOptions []string  `json:"answer_options"`
	CorrectAnswer string    `json:"correct_answer"`
	Explanation   string    `json:"explanation"`
	CreatedAt     time.Time `json:"created_at"`
}

// QuizSession represents an active quiz session for a user
type QuizSession struct {
	ID                  int        `json:"id"`
	UserID              int64      `json:"user_id"`
	Language            string     `json:"language"`
	CurrentQuestionIndex int       `json:"current_question_index"`
	QuestionIDs         []int      `json:"question_ids"`
	CorrectAnswers      int        `json:"correct_answers"`
	StartedAt           time.Time  `json:"started_at"`
	CompletedAt         *time.Time `json:"completed_at,omitempty"`
}

// QuizAnswer represents a user's answer to a quiz question
type QuizAnswer struct {
	ID         int       `json:"id"`
	UserID     int64     `json:"user_id"`
	SessionID  int       `json:"session_id"`
	QuestionID int       `json:"question_id"`
	AnswerGiven string    `json:"answer_given"`
	IsCorrect  bool      `json:"is_correct"`
	AnsweredAt time.Time `json:"answered_at"`
}

// IsComplete returns true if the quiz session is complete
func (s *QuizSession) IsComplete() bool {
	return s.CurrentQuestionIndex >= len(s.QuestionIDs) || s.CompletedAt != nil
}

// GetScore returns the score as a percentage
func (s *QuizSession) GetScore() float64 {
	if s.CurrentQuestionIndex == 0 {
		return 0
	}
	return float64(s.CorrectAnswers) / float64(s.CurrentQuestionIndex) * 100
}

// ToJSON converts the QuestionIDs to a JSON string
func (s *QuizSession) ToJSON() (string, error) {
	b, err := json.Marshal(s.QuestionIDs)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// FromJSON populates QuestionIDs from a JSON string
func (s *QuizSession) FromJSON(jsonStr string) error {
	return json.Unmarshal([]byte(jsonStr), &s.QuestionIDs)
}
package store

import (
	"sync"

	"github.com/amiosamu/interview-match-bot/internal/models"
)

// UserStore handles storing and retrieving user information
type UserStore struct {
	users map[int64]*models.User
	mutex sync.RWMutex
}

// NewUserStore creates a new UserStore instance
func NewUserStore() *UserStore {
	return &UserStore{
		users: make(map[int64]*models.User),
	}
}

// SaveUser stores or updates a user
func (s *UserStore) SaveUser(user *models.User) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.users[user.ID] = user
}

// GetUser retrieves a user by ID
func (s *UserStore) GetUser(userID int64) (*models.User, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	user, exists := s.users[userID]
	return user, exists
}

// FindMatches returns users that match the given criteria
func (s *UserStore) FindMatches(userID int64, field, level string) []*models.User {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var matches []*models.User
	for id, user := range s.users {
		if id != userID && user.Field == field && user.Level == level {
			matches = append(matches, user)
		}
	}
	return matches
}

// SetUserField updates the field for a specific user
func (s *UserStore) SetUserField(userID int64, field string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	user, exists := s.users[userID]
	if !exists {
		return
	}
	user.Field = field
}

// SetUserLevel updates the level for a specific user
func (s *UserStore) SetUserLevel(userID int64, level string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	user, exists := s.users[userID]
	if !exists {
		return
	}
	user.Level = level
}

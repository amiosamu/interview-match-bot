package models

// User represents a user of the interview bot
type User struct {
	ID        int64  // Telegram user ID
	Username  string // Telegram username (may be empty)
	FirstName string // Telegram first name
	LastName  string // Telegram last name (may be empty)
	Field     string // Selected field of interest (e.g., "backend", "frontend")
	Level     string // Selected experience level (e.g., "intern", "junior", "middle", "senior")
}

// DisplayName returns the best available name for the user
// If username is available, it returns the username
// Otherwise, it returns the full name or just the first name
func (u *User) DisplayName() string {
	if u.Username != "" {
		return "@" + u.Username
	}
	
	if u.LastName != "" {
		return u.FirstName + " " + u.LastName
	}
	
	return u.FirstName
}

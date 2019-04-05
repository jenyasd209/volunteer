package user

import "time"

// Session struct for save user session
type Session struct {
	ID        int
	UUID      string
	Email     string
	UserID    int
	CreatrdAt time.Time
}

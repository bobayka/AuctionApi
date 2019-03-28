package session

import "time"

type Session struct {
	SessionId  string    `json:"session_id"`
	UserID     int64     `json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
	ValidUntil time.Time `json:"valid_until"`
}

func (s *Session) CheckTokenTime() bool {
	if s.ValidUntil.After(time.Now()) {
		return true
	}
	return false
}

package domains

import "time"

type Session struct {
	SessionID  string    `json:"session_id"`
	UserID     int64     `json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
	ValidUntil time.Time `json:"valid_until"`
}

func (s *Session) CheckTokenTime() bool {
	return s.ValidUntil.After(time.Now())
}

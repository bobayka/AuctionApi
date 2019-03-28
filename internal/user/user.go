package user

import "time"

type User struct {
	ID        int        `json:"id"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Birthday  *time.Time `json:"birthday"`
	Email     string     `json:"email"`
	Password  string     `json:"-"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"-"`
}

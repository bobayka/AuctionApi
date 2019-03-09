package user

import "time"

type User struct {
	ID        int
	FirstName string
	LastName  string
	Birthday  time.Time
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

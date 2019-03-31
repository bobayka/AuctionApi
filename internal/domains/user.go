package domains

import (
	"gitlab.com/bobayka/courseproject/pkg/customTime"
	"time"
)

type User struct {
	ID        int                    `json:"id"`
	FirstName string                 `json:"first_name"`
	LastName  string                 `json:"last_name"`
	Birthday  *customTime.CustomTime `json:"birthday"`
	Email     string                 `json:"email"`
	Password  string                 `json:"-"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"-"`
}

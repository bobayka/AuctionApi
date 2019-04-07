package domains

import (
	"gitlab.com/bobayka/courseproject/pkg/customTime"
	"time"
)

type User struct {
	ID        int64                  `json:"id"`
	FirstName string                 `json:"first_name"`
	LastName  string                 `json:"last_name"`
	Birthday  *customtime.CustomTime `json:"birthday"`
	Email     string                 `json:"email"`
	Password  string                 `json:"-"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"-"`
}

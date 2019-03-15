package request

import "time"

type BasicInfo struct {
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Birthday  *time.Time `json:"birthday,omitempty"`
}

type GeneralInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type RegUser struct {
	BasicInfo
	GeneralInfo
}

func (r *RegUser) GetEmail() string {
	return r.Email
}

type AuthUser struct {
	GeneralInfo
	Token string `json:"authorization"`
}

func (a *AuthUser) GetEmail() string {
	return a.Email
}
func (a *AuthUser) GetToken() string {
	return a.Token
}

type UpdateUser struct {
	BasicInfo
	Token string `json:"authorization"`
}

func (u *UpdateUser) GetToken() string {
	return u.Token
}

type TokenGetter interface {
	GetToken() string
}

type EmailGetter interface {
	GetEmail() string
}

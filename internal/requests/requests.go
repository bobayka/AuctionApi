package request

import (
	"gitlab.com/bobayka/courseproject/pkg/customTime"
)

type BasicInfo struct {
	FirstName string                 `json:"first_name"`
	LastName  string                 `json:"last_name"`
	Birthday  *customTime.CustomTime `json:"birthday,omitempty"`
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
}

func (a *AuthUser) GetEmail() string {
	return a.Email
}

type UpdateUser struct {
	BasicInfo
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
}

func (u *UpdateUser) GetToken() string {
	return u.AccessToken
}

type TokenGetter interface {
	GetToken() string
}

type EmailGetter interface {
	GetEmail() string
}

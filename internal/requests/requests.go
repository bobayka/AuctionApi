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

func (u RegUser) GetEmail() string {
	return u.Email
}

func (u RegUser) GetToken() string {
	return ""
}

type AuthUser struct {
	GeneralInfo
	Token string `json:"authorization"`
}

func (u AuthUser) GetEmail() string {
	return u.Email
}

func (u AuthUser) GetToken() string {
	return u.Token
}

type UpdateUser struct {
	BasicInfo
	Token string `json:"authorization"`
}

func (u UpdateUser) GetEmail() string {
	return ""
}

func (u UpdateUser) GetToken() string {
	return u.Token
}

type UserEmail interface {
	GetEmail() string
}

type UserToken interface {
	GetToken() string
}
type User interface {
	UserEmail
	UserToken
}

package request

import (
	"gitlab.com/bobayka/courseproject/pkg/customTime"
	"time"
)

type BasicInfo struct {
	FirstName string                 `json:"first_name"`
	LastName  string                 `json:"last_name"`
	Birthday  *customtime.CustomTime `json:"birthday,omitempty"`
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

func (u *UpdateUser) GetTokenType() string {
	return u.TokenType
}

type TokenTypeGetter interface {
	GetTokenType() string
}

type EmailGetter interface {
	GetEmail() string
}

type LotToCreateUpdate struct {
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	MinPrice    float64   `json:"min_price"`
	PriceStep   *float64  `json:"price_step"`
	EndAt       time.Time `json:"end_at"`
	Status      string    `json:"status"`
}

type WebLotToCreateUpdate struct {
	Lot      LotToCreateUpdate
	BuyPrice string `json:"buy_price"`
	Status   string `json:"status"`
}

type Token struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
}

type Price struct {
	Price float64 `json:"price"`
}

package request

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/cmd/Protobuf"
	"gitlab.com/bobayka/courseproject/internal/MyGRPCLib"
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

type LotCreateUpdate struct {
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	MinPrice    float64   `json:"min_price"`
	PriceStep   *float64  `json:"price_step"`
	EndAt       time.Time `json:"end_at"`
	Status      *string   `json:"status"`
}

type WebLotToCreateUpdate struct {
	Lot      LotCreateUpdate
	BuyPrice string  `json:"buy_price"`
	Status   *string `json:"status"`
}

type Token struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
}

type Price struct {
	Price float64 `json:"price"`
}

func ConvLotCreateUpdateToGRPC(lot *LotCreateUpdate, userID *int64, lotID *int64) (*lotspb.LotCreateUpdate, error) {
	description := MyGRPCLib.ConvStringPointerToString(lot.Description)
	priceStep := MyGRPCLib.ConvFloat64PointerToFloat64(lot.PriceStep)
	endAt, err := ptypes.TimestampProto(lot.EndAt)
	if err != nil {
		return nil, errors.Wrap(err, "can't convert time to timestamp")
	}
	status := MyGRPCLib.ConvStringPointerToString(lot.Status)
	uID := MyGRPCLib.ConvInt64PointerToInt64(userID)
	lID := MyGRPCLib.ConvInt64PointerToInt64(lotID)
	return &lotspb.LotCreateUpdate{Title: lot.Title, Description: description,
		MinPrice: lot.MinPrice, PriceStep: priceStep, EndAt: endAt, Status: status, UserID: uID, LotID: lID}, nil
}

func ConvGRPCToLotCreateUpdate(lot *lotspb.LotCreateUpdate) (*LotCreateUpdate, *int64, *int64, error) {
	description := MyGRPCLib.ConvStringToStringPointer(lot.Description)
	priceStep := MyGRPCLib.ConvFloat64ToFloat64Pointer(lot.PriceStep)
	endAt, err := ptypes.Timestamp(lot.EndAt)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "can't convert timestamp to time")
	}
	status := MyGRPCLib.ConvStringToStringPointer(lot.Status)
	userID := MyGRPCLib.ConvInt64ToInt64Pointer(lot.UserID)
	lotID := MyGRPCLib.ConvInt64ToInt64Pointer(lot.LotID)
	return &LotCreateUpdate{Title: lot.Title, Description: description,
		MinPrice: lot.MinPrice, PriceStep: priceStep, EndAt: endAt, Status: status}, userID, lotID, nil
}

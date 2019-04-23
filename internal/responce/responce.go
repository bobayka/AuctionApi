package responce

import (
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/cmd/Protobuf"
	"gitlab.com/bobayka/courseproject/internal/MyGRPCLib"
	"gitlab.com/bobayka/courseproject/internal/domains"
	"net/http"
)

type RespLot struct {
	domains.LotGeneral
	Creator ShortUser  `json:"creator"`
	Buyer   *ShortUser `json:"buyer,omitempty"`
}

type ShortUser struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func RespondJSON(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	fmt.Fprintln(w, msg)
}

func ConvShortUserToGRPC(su *ShortUser) *lotspb.ShortUser {
	if su == nil {
		return nil
	}
	return &lotspb.ShortUser{ID: su.ID, FirstName: su.FirstName, LastName: su.LastName}
}

func ConvGRPCToShortUser(su *lotspb.ShortUser) *ShortUser {
	if su == nil {
		return nil
	}
	return &ShortUser{ID: su.ID, FirstName: su.FirstName, LastName: su.LastName}
}

func ConvertRespLotToGRPC(resp *RespLot) (*lotspb.Lot, error) {
	description := MyGRPCLib.ConvStringPointerToString(resp.Description)
	buyPrice := MyGRPCLib.ConvFloat64PointerToFloat64(resp.BuyPrice)
	endAt, err := ptypes.TimestampProto(resp.EndAt)
	if err != nil {
		return nil, errors.Wrap(err, "can't convert time to timestamp")
	}
	createdAt, err := ptypes.TimestampProto(resp.CreatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "can't convert time to timestamp")
	}
	updatedAt, err := ptypes.TimestampProto(resp.UpdatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "can't convert time to timestamp")
	}
	creator := ConvShortUserToGRPC(&resp.Creator)

	buyer := ConvShortUserToGRPC(resp.Buyer)
	return &lotspb.Lot{ID: resp.ID, Title: resp.Title, Description: description,
		BuyPrice: buyPrice, MinPrice: resp.MinPrice, PriceStep: resp.PriceStep,
		Status: resp.Status, EndAt: endAt, CreatedAt: createdAt, UpdatedAt: updatedAt,
		Creator: creator, Buyer: buyer}, nil
}

func ConvertGRPCToRespLot(resp *lotspb.Lot) (*RespLot, error) {
	price := MyGRPCLib.ConvFloat64ToFloat64Pointer(resp.BuyPrice)
	endAt, err := ptypes.Timestamp(resp.EndAt)
	if err != nil {
		return nil, errors.Wrap(err, "can't convert timestamp to time")
	}
	createdAt, err := ptypes.Timestamp(resp.CreatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "can't convert timestamp to time")
	}
	updatedAt, err := ptypes.Timestamp(resp.UpdatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "can't convert timestamp to time")
	}
	buyer := ConvGRPCToShortUser(resp.Buyer)

	return &RespLot{LotGeneral: domains.LotGeneral{ID: resp.ID, Title: resp.Title,
		Description: &resp.Description, BuyPrice: price, MinPrice: resp.MinPrice,
		PriceStep: resp.PriceStep, Status: resp.Status, EndAt: endAt, CreatedAt: createdAt,
		UpdatedAt: updatedAt},
		Creator: ShortUser{ID: resp.Creator.ID, FirstName: resp.Creator.FirstName, LastName: resp.Creator.LastName},
		Buyer:   buyer,
	}, nil
}

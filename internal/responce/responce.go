package responce

import "gitlab.com/bobayka/courseproject/internal/domains"

type LotToResponce struct {
	*domains.Lot
	Creator *ShortUSer `json:"creator"`
	Buyer   *ShortUSer `json:"buyer"`
}

type ShortUSer struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

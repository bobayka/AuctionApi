package responce

import (
	"fmt"
	"gitlab.com/bobayka/courseproject/internal/domains"
	"net/http"
)

type RespLot struct {
	domains.LotGeneral
	Creator ShortUSer  `json:"creator"`
	Buyer   *ShortUSer `json:"buyer,omitempty"`
}

type ShortUSer struct {
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

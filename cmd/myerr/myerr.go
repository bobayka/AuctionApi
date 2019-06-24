package myerr

import (
	"encoding/json"
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/internal/responce"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
	"strings"
)

// nolint: gochecknoglobals
var (
	ErrConflict     = errors.New("conflict")
	ErrUnauthorized = errors.New("unauthorized")
	ErrBadRequest   = errors.New("bad request")
	ErrNotFound     = errors.New("not found")
)

func Split(str string) []string {
	return strings.Split(str, "$")
}

type ErrString struct {
	Err string `json:"error"`
}

func ConvGRPCStatusToMyError(status *status.Status) error {
	switch status.Code() {
	case http.StatusNotFound:
		return errors.Wrap(ErrNotFound, status.Message())
	case http.StatusBadRequest:
		return errors.Wrap(ErrBadRequest, status.Message())
	case http.StatusUnauthorized:
		return errors.Wrap(ErrUnauthorized, status.Message())
	case http.StatusConflict:
		return errors.Wrap(ErrConflict, status.Message())
	}
	return errors.New(status.Message())
}
func GiveErrToClient(str string) string {
	lastInd := strings.LastIndex(str, "$")
	initialInd := strings.Index(str, "$")
	if lastInd == -1 {
		log.Fatalf(`cant find $ in: \n --%s--`, str)
		return ""
	}
	if lastInd == initialInd {
		log.Fatalf(`cant find second $ in: \n --%s--`, str)
	}
	return str[initialInd+1 : lastInd]
}

func ErrMarshal(v string) ([]byte, error) {
	return json.Marshal(ErrString{Err: v})
}

func JSONErrRespond(w http.ResponseWriter, msg string, code int) {
	resp, _ := ErrMarshal(msg)
	responce.RespondJSON(w, string(resp), code)
}

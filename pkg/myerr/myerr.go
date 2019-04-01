package myerr

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

var (
	UnprocessableEntity = errors.New("unprocessable entity")
	Conflict            = errors.New("conflict")
	Created             = errors.New("created")
	Accepted            = errors.New("accepted")
	Unauthorized        = errors.New("unauthorized")
	BadRequest          = errors.New("bad request")
	Success             = errors.New("success")
	NotFound            = errors.New("not found")
)

func Split(str string) []string {
	return strings.Split(str, "$")
}

type ErrString struct {
	Err string `json:"error"`
}

func GetClientErr(str string) string {
	return str[strings.Index(str, "$")+1 : strings.LastIndex(str, "$")]
}

func Error(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	fmt.Fprintln(w, error)
}

func ErrMarshal(v string) ([]byte, error) {
	return json.Marshal(ErrString{Err: v})
}

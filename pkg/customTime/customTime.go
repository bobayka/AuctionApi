package customtime

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
	"time"
)

const CTLayout = "2006-01-02"

type CustomTime time.Time

func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse(CTLayout, s)
	if err != nil {
		return errors.Wrap(err, "cant parse time")
	}
	*ct = CustomTime(t)
	return
}

func (ct *CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", time.Time(*ct).Format(CTLayout))), nil
}

package customTime

import (
	"fmt"
	"strings"
	"time"
)

const CTLayout = "2006-01-02"

type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return
	}
	ct.Time, err = time.Parse(CTLayout, s)
	return
}

func (ct *CustomTime) MarshalJSON() ([]byte, error) {
	if ct.Time.UnixNano() == NilTime {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", ct.Time.Format(CTLayout))), nil
}

var NilTime = (time.Time{}).UnixNano()

func (ct *CustomTime) IsSet() bool {
	return ct.UnixNano() != NilTime
}

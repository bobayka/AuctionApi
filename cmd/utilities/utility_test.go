package utility

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"gitlab.com/bobayka/courseproject/cmd/myerr"
	"net/http/httptest"
	"testing"
)

func TestReadReqData(t *testing.T) {
	r := require.New(t)

	type testCase struct {
		Name     string
		body     []byte
		Expected error
	}

	testCases := []testCase{
		{
			Name:     "Check word bearer",
			body:     []byte("fdffd"),
			Expected: myerr.ErrBadRequest,
		},
		{
			Name:     "Check word bearer",
			body:     []byte(`{"id": 72}`),
			Expected: nil,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://testing", bytes.NewReader(tc.body))
			dd := struct {
				id int64
			}{}
			err := ReadReqData(req, &dd)
			fmt.Println(dd.id)
			r.Equal(tc.Expected, errors.Cause(err), "Wrong body: %s", tc.body)

		})
	}
}

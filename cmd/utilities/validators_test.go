package utility

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"gitlab.com/bobayka/courseproject/cmd/myerr"
	"testing"
)

func TestCheckEmail(t *testing.T) {
	r := require.New(t)
	type testCase struct {
		Name     string
		Email    string
		Expected error
	}

	testCases := []testCase{
		{
			Name:     "Check right email",
			Email:    "andrew@mail.ru",
			Expected: nil,
		},

		{
			Name:     "Check wrong email #1",
			Email:    "andrewmail.ru",
			Expected: myerr.ErrBadRequest,
		},

		{
			Name:     "Check wrong email #2",
			Email:    "andr---1112ed&@w@mailru",
			Expected: myerr.ErrBadRequest,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			r.Equal(tc.Expected, errors.Cause(CheckEmail(tc.Email)), "Wrong result: %s", tc.Email)
		})

	}

}

func TestCheckBearer(t *testing.T) {
	r := require.New(t)
	type testCase struct {
		Name        string
		Bearer      string
		ExpectedErr error
		ExpectedTok string
	}

	testCases := []testCase{
		{
			Name:        "Check right email",
			Bearer:      "Bearer ",
			ExpectedErr: nil,
			ExpectedTok: "",
		},

		{
			Name:        "Check wrong email #1",
			Bearer:      "Bearer 123333",
			ExpectedErr: nil,
			ExpectedTok: "",
		},

		{
			Name:        "Check wrong email #2",
			Bearer:      "Beare ffffkfkfkf",
			ExpectedErr: myerr.ErrBadRequest,
			ExpectedTok: "",
		},
		{
			Name:        "Check wrong email #2",
			Bearer:      "Bearer ffffkfkfkf",
			ExpectedErr: nil,
			ExpectedTok: "ffffkfkfkf",
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			str, err := CheckBearer(tc.Bearer)
			r.Equal(tc.ExpectedErr, errors.Cause(err), "Wrong result: %s", tc.Bearer)
			r.Equal(tc.ExpectedTok, str, "Wrong result: %s", tc.Bearer)
		})

	}
}

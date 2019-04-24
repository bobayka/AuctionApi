package utility

import (
	"github.com/stretchr/testify/require"
	"gitlab.com/bobayka/courseproject/internal/postgres"
	"gitlab.com/bobayka/courseproject/internal/postgres/storage"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckTokenMiddleware(t *testing.T) {

	r := require.New(t)
	db, err := postgres.PGInit("localhost", 5432, "bobayka", "12345", "TinkoffProj")
	r.NoError(err)
	sessionStorage, err := storage.NewSessionsStorage(db)
	r.NoError(err)
	type head struct {
		key   string
		value string
	}
	type testCase struct {
		Name     string
		headers  head
		Expected int
	}

	testCases := []testCase{
		{
			Name:     "Check word bearer",
			Expected: http.StatusBadRequest,
			headers:  head{key: "Authorization", value: "Beare tqqcwfsjedlikwaxqwid"},
		},
		{
			Name:     "Check authorization",
			Expected: http.StatusUnauthorized,
			headers:  head{key: "Authorization", value: "Bearer tqqcwfsjedlikwaxqwi"},
		},
		{
			Name:     "Check access",
			Expected: http.StatusOK,
			headers:  head{key: "Authorization", value: "Bearer tqqcwfsjedlikwaxqwid"},
		},
		{
			Name:     "Check key",
			Expected: http.StatusBadRequest,
			headers:  head{key: "Authorizatio", value: "Bearer tqqcwfsjedlikwaxqwid"},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://testing", nil)
			req.Header.Add(tc.headers.key, tc.headers.value)
			handlerToTest := CheckTokenMiddleware(sessionStorage)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

			// call the handler using a mock response recorder (we'll not use that anyway)
			recorder := httptest.NewRecorder()
			handlerToTest.ServeHTTP(recorder, req)

			respBody := recorder.Body.String()
			r.Equal(tc.Expected, recorder.Code, "Wrong http code, response: %s", respBody)

		})
	}
}

package lothandlers

//import (
//	"github.com/go-chi/chi"
//	"github.com/stretchr/testify/require"
//	lotspb "gitlab.com/bobayka/courseproject/cmd/Protobuf"
//	utility "gitlab.com/bobayka/courseproject/cmd/utilities"
//	"gitlab.com/bobayka/courseproject/internal/postgres"
//	"gitlab.com/bobayka/courseproject/internal/postgres/storage"
//	"google.golang.org/grpc"
//	"io"
//	"io/ioutil"
//	"log"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//)
//
//const path = "/v1/auction/lots/"
//const method = "PUT"
//func TestLotServiceHandler_UpdatePriceHandler(t *testing.T) {
//	r := require.New(t)
//	db, err := postgres.PGInit("localhost", 5432, "bobayka", "12345", "TinkoffProj")
//	r.NoError(err)
//	sessionStorage, err := storage.NewSessionsStorage(db)
//	r.NoError(err)
//
//	router := chi.NewRouter()
//
//	router.Use(utility.CheckTokenMiddleware(sessionStorage))
//
//	conn, err := grpc.Dial("localhost:5001", grpc.WithInsecure())
//	if err != nil {
//		log.Fatalf("can't connect to server: %v", err)
//	}
//	defer conn.Close()
//	client := lotspb.NewLotsServiceClient(conn)
//	lot := NewLotHandler(client)
//	lotRouter := lot.Routes()
//	router.Mount("/v1/auction/lots", lotRouter)
//	ts := httptest.NewServer(router)
//	defer ts.Close()
//	type testCase struct {
//		Name     string
//		path     string
//		headers  map[string]string
//		Expected int
//	}
//	testCases := []testCase{
//		{
//			Name:     "Check user token",
//			path:     "0",
//			Expected: http.StatusOK,
//			headers:  map[string]string{"Authorization": "Bearer tqqcwfsjedlikwaxqwid",},
//		},
//
//		{
//			Name:     "Check word bearer",
//			path:     "0",
//			Expected: http.StatusBadRequest,
//			headers:  map[string]string{"Authorization": "Beare tqqcwfsjedlikwaxqwid",},
//		},
//		{
//			Name:     "Check authorization",
//			path:     "2222222",
//			Expected: http.StatusUnauthorized,
//			headers:  map[string]string{"Authorization": "Bearer tqqcwfsjedlikwaxqw",},
//		},
//		{
//			Name:     "Check access",
//			path:     "2222222",
//			Expected: http.StatusNotFound,
//			headers:  map[string]string{"Authorization": "Bearer tqqcwfsjedlikwaxqwid",},
//		},
//	}
//	for _, tc := range testCases {
//		tc := tc
//		t.Run(tc.Name, func(t *testing.T) {
//			resp, body := testRequest(t, ts, method, path+tc.path, nil, tc.headers)
//			r.Equal(tc.Expected, resp.StatusCode,"Wrong http code, response: %s", body)
//		})
//
//	}
//	log.Fatal(http.ListenAndServe(":5000", router))
//}
//func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader, headers map[string]string) (*http.Response, string) {
//	req, err := http.NewRequest(method, ts.URL+path, body)
//	if err != nil {
//		t.Fatal(err)
//		return nil, ""
//	}
//	for key, val := range headers {
//		req.Header.Add(key, val)
//	}
//	resp, err := http.DefaultClient.Do(req)
//	if err != nil {
//		t.Fatal(err)
//		return nil, ""
//	}
//
//	respBody, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		t.Fatal(err)
//		return nil, ""
//	}
//	defer resp.Body.Close()
//	return resp, string(respBody)
//}

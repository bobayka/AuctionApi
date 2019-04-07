package utility

import (
	"context"
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/cmd/myerr"
	"gitlab.com/bobayka/courseproject/internal/postgres"
	"gitlab.com/bobayka/courseproject/internal/services"
	"net/http"
)

func CheckTokenMiddleware(store *postgres.UsersStorage) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) error {
			if r.Header.Get("token_type") != "bearer" {
				return errors.Wrap(myerr.ErrBadRequest, "$invalid token type$")
			}
			s, err := services.CheckValidToken(r.Header.Get("access_token"), store)
			if err != nil {
				return errors.Wrap(err, "cant check valid token")
			}
			ctx := context.WithValue(r.Context(), UserIDKey, s.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
			return nil
		}
		return MakeHandler(fn)
	}
}

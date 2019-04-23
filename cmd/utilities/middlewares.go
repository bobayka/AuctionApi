package utility

import (
	"context"
	"github.com/pkg/errors"
	"gitlab.com/bobayka/courseproject/internal/postgres/storage"
	"gitlab.com/bobayka/courseproject/internal/services"
	"net/http"
)

func CheckTokenMiddleware(store *storage.SessionStorage) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) error {
			token := r.Header.Get("Authorization")
			token, err := CheckBearer(token)
			if err != nil {
				return errors.Wrap(err, "doesn't valid token")
			}
			s, err := services.CheckValidToken(token, store)
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

func CheckCookieMiddleware(store *storage.SessionStorage) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) error {
			token, err := r.Cookie("Authorization")
			if err != nil {
				switch err {
				case http.ErrNoCookie:
					http.Redirect(w, r, "http://localhost:5000/signin", http.StatusFound)
					return nil

				default:
					return errors.Wrap(err, "Cant read cookie")
				}
			}
			s, err := services.CheckValidToken(token.Value, store)
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

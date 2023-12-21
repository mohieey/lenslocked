package middlewares

import (
	"net/http"

	"github.com/mohieey/lenslocked/appctx"
	"github.com/mohieey/lenslocked/controllers"
	"github.com/mohieey/lenslocked/services"
)

type UserMiddleware struct {
	SessionService *services.SessionService
}

func (umw *UserMiddleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie(controllers.CookieSession)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		user, err := umw.SessionService.User(tokenCookie.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()
		ctx = appctx.WithUser(ctx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (umw *UserMiddleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := appctx.User(r.Context())
		if user == nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

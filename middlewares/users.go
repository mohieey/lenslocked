package middlewares

import (
	"net/http"

	"github.com/mohieey/lenslocked/appctx"
	"github.com/mohieey/lenslocked/controllers"
	"github.com/mohieey/lenslocked/models"
)

type UserMiddleware struct {
	SessionService *models.SessionService
}

func (usm *UserMiddleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie(controllers.CookieSession)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		user, err := usm.SessionService.User(tokenCookie.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()
		ctx = appctx.WithUser(ctx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

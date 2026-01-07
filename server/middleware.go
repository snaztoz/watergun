package server

import (
	"context"
	"crypto"
	"net/http"
	"strings"

	"github.com/snaztoz/watergun/serverctx"
	"github.com/snaztoz/watergun/socket"
)

func bearerTokenParser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")

		ctx := r.Context()

		if strings.HasPrefix(authorization, "Bearer ") {
			key := strings.TrimPrefix(authorization, "Bearer ")
			ctx = context.WithValue(ctx, serverctx.AccessTokenKey, key)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func socketRouteAuth(key crypto.PublicKey) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accessToken, ok := r.Context().Value(serverctx.AccessTokenKey).(string)
			if !ok {
				http.Error(w, http.StatusText(403), 403)
				return
			}

			userID, err := socket.Authenticate(accessToken, key)
			if err != nil {
				http.Error(w, http.StatusText(403), 403)
				return
			}

			ctx := context.WithValue(r.Context(), serverctx.UserIDKey, userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func adminRoutesAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key, ok := r.Context().Value(serverctx.AccessTokenKey).(string)
		if !ok {
			http.Error(w, http.StatusText(403), 403)
			return
		}

		if key != adminAPIKey() {
			http.Error(w, http.StatusText(403), 403)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func jsonContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

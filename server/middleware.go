package server

import (
	"context"
	"crypto"
	"net/http"
	"strings"

	"github.com/snaztoz/watergun/serverctx"
	"github.com/snaztoz/watergun/socket"
)

const (
	allowQueryParamToken = true
)

func accessTokenParser(allowQueryParamToken bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hasBearerAuthorizationHeader := func(r *http.Request) bool {
			header := r.Header.Get("Authorization")
			return header != "" && strings.HasPrefix(header, "Bearer ")
		}

		hasQueryParamToken := func(r *http.Request) bool {
			return r.URL.Query().Has("token")
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			switch {
			case hasBearerAuthorizationHeader(r):
				header := r.Header.Get("Authorization")
				key := strings.TrimPrefix(header, "Bearer ")
				ctx = context.WithValue(ctx, serverctx.AccessTokenKey, key)

			case allowQueryParamToken && hasQueryParamToken(r):
				token := r.URL.Query().Get("token")
				ctx = context.WithValue(ctx, serverctx.AccessTokenKey, token)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
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

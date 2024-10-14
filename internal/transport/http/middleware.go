package http

import (
	"context"
	"fmt"

	"net/http"
	"strings"

	"github.com/yervsil/auth_service/internal/token"
	"github.com/yervsil/auth_service/internal/utils"
)

type contextKey string // Определение типа ключа для контекста

const UserDataKey = contextKey("userData")

func(h *Handler) JWTMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc( 
			func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			h.log.Warn("empty auth header")

			utils.SendJson(w, "empty auth header", http.StatusUnauthorized)
			return
		}

		headersSlice := strings.Split(authHeader, " ")
		if len(headersSlice) < 2 {
			h.log.Warn("invalid auth header")

			utils.SendJson(w, "invalid auth header", http.StatusUnauthorized)
			return
		}

		authToken := headersSlice[1]
		
		data, err := token.ParseToken(authToken)
		if err != nil {
			h.log.Warn(fmt.Sprintf("error token parsing: %s", err.Error()))

			utils.SendJson(w, err.Error(), http.StatusUnauthorized)
			return
		}
		
		ctx := context.WithValue(r.Context(), UserDataKey, data)

		r = r.WithContext(ctx)
		
		next.ServeHTTP(w, r)
	})
}
}
package middlewares

import (
	"context"
	"go-identity/pkg/models"
	"go-identity/pkg/responses"
	"net/http"
	"strings"
)

func ContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(model models.Model, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			responses.NewUnauthorized(w, "Unauthorized.")
			return
		}

		user, err := model.GetUserByToken(r.Context(), parts[1])
		if err != nil {
			responses.NewInternalServerErrorResponse(w, err.Error())
			return
		}

		nContext := context.WithValue(r.Context(), "user", user)
		r = r.WithContext(nContext)

		next.ServeHTTP(w, r)
	})
}

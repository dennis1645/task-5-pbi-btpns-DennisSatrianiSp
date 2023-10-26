package middleware

import (
	"net/http"

	"github.com/dennis1645/go-api/config"
	"github.com/dennis1645/go-api/helper"
	"github.com/golang-jwt/jwt/v4"
)

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err == http.ErrNoCookie {
			response := map[string]string{"Peringatan ": "Tidak Boleh ! Login dulu"}
			helper.ResponseJSON(w, http.StatusUnauthorized, response)
			return
		}
		// Mengambil token value
		tokenString := c.Value

		claims := &config.JWTClaim{}

		// parsing token
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return config.JWT_KEY, nil
		})
		if err != nil {
			v, _ := err.(*jwt.ValidationError)
			switch v.Errors {
			case jwt.ValidationErrorSignatureInvalid:
				// Token Invalid
				response := map[string]string{"Peringatan ": "Token Invalid"}
				helper.ResponseJSON(w, http.StatusUnauthorized, response)
				return
			case jwt.ValidationErrorExpired:
				// Token Expired
				response := map[string]string{"Peringatan ": "Token Expired"}
				helper.ResponseJSON(w, http.StatusUnauthorized, response)
				return
			default:
				response := map[string]string{"Peringatan ": "Tidak Boleh Masuk"}
				helper.ResponseJSON(w, http.StatusUnauthorized, response)
				return
			}
		}
		if !token.Valid {
			response := map[string]string{"Peringatan ": "Tidak Boleh ! Login dulu"}
			helper.ResponseJSON(w, http.StatusUnauthorized, response)
			return
		}

		next.ServeHTTP(w, r)

	})
}

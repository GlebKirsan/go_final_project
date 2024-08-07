package middleware

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"

	"github.com/GlebKirsan/go-final-project/internal/config"
	"github.com/GlebKirsan/go-final-project/internal/service"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg := config.Get()
		pass := cfg.Pass
		if len(pass) > 0 {
			var jwtToken string

			cookie, err := r.Cookie("token")
			if err != nil {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
			jwtToken = cookie.Value

			token, err := jwt.Parse(jwtToken, func(t *jwt.Token) (interface{}, error) {
				return []byte(cfg.Secret), nil
			})

			if !token.Valid {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}

			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			res, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return
			}

			tokenPassHash := res["hash"]
			if tokenPassHash != service.GetMD5Hash(pass) {
				http.Error(w, "wrong password hash", http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

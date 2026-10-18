package middleware

import (
	"net/http"

	"github.com/icekidtech/breakthrough-palace-mission/backend/auth"
	"github.com/icekidtech/breakthrough-palace-mission/backend/config"
)

func AuthMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get token from cookie
			cookie, err := r.Cookie("auth_token")
			if err != nil {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}

			// Verify JWT
			claims, err := auth.VerifyJWT(cookie.Value, cfg.JWTSecret)
			if err != nil {
				http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
				return
			}

			// Add claims to request context
			r.Header.Set("X-Member-ID", claims.MemberID)
			r.Header.Set("X-Family-ID", claims.FamilyID)
			r.Header.Set("X-Role", claims.Role)

			next.ServeHTTP(w, r)
		})
	}
}

func CORSMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Whitelist frontend URL
			allowedOrigins := map[string]bool{
				cfg.FrontendURL: true,
			}

			if allowedOrigins[origin] {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func ExtractFamilyID(r *http.Request) string {
	return r.Header.Get("X-Family-ID")
}

func ExtractMemberID(r *http.Request) string {
	return r.Header.Get("X-Member-ID")
}

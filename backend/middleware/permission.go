package middleware

import (
	"net/http"

	"github.com/icekidtech/breakthrough-palace-mission/backend/rbac"
)

// RequirePermission middleware checks if user has required permission
func RequirePermission(permission rbac.Permission) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract role from request (set by auth middleware)
			role := r.Header.Get("X-Role")
			if role == "" {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}

			// Check permission
			if !rbac.HasPermission(rbac.Role(role), permission) {
				http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

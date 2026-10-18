package middleware

import (
	"net/http"

	"github.com/icekidtech/breakthrough-palace-mission/backend/rbac"
	"gorm.io/gorm"
)

// RequirePermission middleware checks if user has required permission
func RequirePermission(db *gorm.DB, permission rbac.Permission) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract role from request (set by auth middleware)
			role := r.Header.Get("X-Role")
			if role == "" {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}

			memberID := r.Header.Get("X-Member-ID")
			if memberID == "" {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}

			// Super admin has all permissions
			if rbac.IsSuperAdmin(role) {
				next.ServeHTTP(w, r)
				return
			}

			// Check if user's role has this permission (from DB)
			var count int64
			db.Table("role_permissions").
				Joins("JOIN roles ON role_permissions.role_id = roles.id").
				Joins("JOIN permissions ON role_permissions.permission_id = permissions.id").
				Where("roles.name = ? AND permissions.name = ?", role, string(permission)).
				Count(&count)

			if count == 0 {
				http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

package migrations

import (
	"github.com/google/uuid"
	"github.com/icekidtech/breakthrough-palace-mission/backend/models"
	"gorm.io/gorm"
)

// SeedRolesAndPermissions creates default roles and permissions
func SeedRolesAndPermissions(db *gorm.DB) error {
	// Check if roles already exist
	var count int64
	db.Model(&models.Role{}).Count(&count)
	if count > 0 {
		return nil // Already seeded
	}

	// Define permissions
	permissions := []models.Permission{
		// Family permissions
		{Name: "family:view", Description: "View family details", Category: "family"},
		{Name: "family:edit", Description: "Edit family details", Category: "family"},
		{Name: "family:view_stats", Description: "View family statistics", Category: "family"},

		// Member permissions
		{Name: "members:view", Description: "View family members", Category: "members"},
		{Name: "members:add", Description: "Add new members", Category: "members"},
		{Name: "members:edit", Description: "Edit member details", Category: "members"},
		{Name: "members:remove", Description: "Remove members", Category: "members"},
		{Name: "members:approve", Description: "Approve new members", Category: "members"},

		// Meeting permissions
		{Name: "meetings:view", Description: "View meetings", Category: "meetings"},
		{Name: "meetings:schedule", Description: "Schedule meetings", Category: "meetings"},
		{Name: "meetings:edit", Description: "Edit meetings", Category: "meetings"},
		{Name: "meetings:delete", Description: "Delete meetings", Category: "meetings"},

		// Minutes permissions
		{Name: "minutes:view", Description: "View minutes", Category: "minutes"},
		{Name: "minutes:publish", Description: "Publish minutes", Category: "minutes"},
		{Name: "minutes:delete", Description: "Delete minutes", Category: "minutes"},
		{Name: "minutes:acknowledge", Description: "Acknowledge minutes", Category: "minutes"},

		// Attendance permissions
		{Name: "attendance:view", Description: "View attendance", Category: "attendance"},
		{Name: "attendance:record", Description: "Record attendance", Category: "attendance"},
		{Name: "attendance:edit", Description: "Edit attendance", Category: "attendance"},

		// Contribution permissions
		{Name: "contributions:view", Description: "View contributions", Category: "contributions"},
		{Name: "contributions:record", Description: "Record contributions", Category: "contributions"},
		{Name: "contributions:confirm", Description: "Confirm contributions", Category: "contributions"},
		{Name: "contributions:edit", Description: "Edit contributions", Category: "contributions"},

		// Biodata permissions
		{Name: "biodata:view", Description: "View biodata", Category: "biodata"},
		{Name: "biodata:edit", Description: "Edit biodata", Category: "biodata"},

		// Admin permissions
		{Name: "admin:manage_families", Description: "Manage families", Category: "admin"},
		{Name: "admin:manage_users", Description: "Manage users", Category: "admin"},
		{Name: "admin:manage_roles", Description: "Manage roles", Category: "admin"},
		{Name: "admin:manage_departments", Description: "Manage departments", Category: "admin"},
		{Name: "admin:view_audit_log", Description: "View audit log", Category: "admin"},
		{Name: "admin:system_settings", Description: "System settings", Category: "admin"},
	}

	// Create permissions with generated UUIDs
	for i := range permissions {
		permissions[i].ID = uuid.New()
	}
	if err := db.CreateInBatches(permissions, 100).Error; err != nil {
		return err
	}

	// Build permission map for easier access
	permMap := make(map[string]uuid.UUID)
	for _, p := range permissions {
		permMap[p.Name] = p.ID
	}

	// Define roles with their permissions
	roles := []models.Role{
		{
			Name:        "super_admin",
			Description: "Super Administrator with full system access",
			Scope:       "platform",
		},
		{
			Name:        "family_leader",
			Description: "Family Leader",
			Scope:       "family",
		},
		{
			Name:        "family_assistant",
			Description: "Family Assistant",
			Scope:       "family",
		},
		{
			Name:        "secretary",
			Description: "Family Secretary",
			Scope:       "family",
		},
		{
			Name:        "evang_coordinator",
			Description: "Evangelism Coordinator",
			Scope:       "family",
		},
		{
			Name:        "member",
			Description: "Family Member",
			Scope:       "family",
		},
		{
			Name:        "church_executive",
			Description: "Church Executive",
			Scope:       "platform",
		},
		{
			Name:        "finance_oversight",
			Description: "Finance Oversight",
			Scope:       "platform",
		},
		{
			Name:        "minister",
			Description: "Minister",
			Scope:       "department",
		},
	}

	// Create roles with generated UUIDs
	for i := range roles {
		roles[i].ID = uuid.New()
	}
	if err := db.Create(roles).Error; err != nil {
		return err
	}

	// Build role map
	roleMap := make(map[string]uuid.UUID)
	for _, r := range roles {
		roleMap[r.Name] = r.ID
	}

	// Assign permissions to roles
	rolePermissions := map[string][]string{
		"super_admin": {
			// All permissions
			"family:view", "family:edit", "family:view_stats",
			"members:view", "members:add", "members:edit", "members:remove", "members:approve",
			"meetings:view", "meetings:schedule", "meetings:edit", "meetings:delete",
			"minutes:view", "minutes:publish", "minutes:delete", "minutes:acknowledge",
			"attendance:view", "attendance:record", "attendance:edit",
			"contributions:view", "contributions:record", "contributions:confirm", "contributions:edit",
			"biodata:view", "biodata:edit",
			"admin:manage_families", "admin:manage_users", "admin:manage_roles", "admin:manage_departments",
			"admin:view_audit_log", "admin:system_settings",
		},
		"family_leader": {
			"family:view", "family:edit", "family:view_stats",
			"members:view", "members:edit", "members:approve",
			"meetings:view",
			"minutes:view", "minutes:acknowledge",
			"attendance:view", "attendance:record",
			"contributions:view", "contributions:confirm",
			"biodata:view",
		},
		"family_assistant": {
			"family:view", "family:view_stats",
			"members:view", "members:edit", "members:approve",
			"meetings:view",
			"minutes:view",
			"attendance:view", "attendance:record",
			"contributions:view",
			"biodata:view",
		},
		"secretary": {
			"family:view", "family:edit", "family:view_stats",
			"members:view", "members:add", "members:edit", "members:approve",
			"meetings:view", "meetings:schedule", "meetings:edit", "meetings:delete",
			"minutes:view", "minutes:publish",
			"attendance:view", "attendance:record", "attendance:edit",
			"contributions:view", "contributions:record", "contributions:edit",
			"biodata:view",
		},
		"evang_coordinator": {
			"family:view", "family:view_stats",
			"members:view",
			"meetings:view",
			"minutes:view",
			"attendance:view", "attendance:record",
			"contributions:view",
			"biodata:view",
		},
		"member": {
			"family:view",
			"members:view",
			"meetings:view",
			"minutes:view",
			"attendance:view", "attendance:record",
			"contributions:view", "contributions:record",
			"biodata:view", "biodata:edit",
		},
		"church_executive": {
			"family:view", "family:view_stats",
			"members:view",
			"meetings:view",
			"minutes:view",
			"attendance:view",
			"contributions:view",
			"biodata:view",
			"admin:view_audit_log", "admin:manage_roles",
		},
		"finance_oversight": {
			"family:view", "family:view_stats",
			"contributions:view", "contributions:confirm",
			"biodata:view",
			"admin:view_audit_log",
		},
		"minister": {
			"family:view",
			"members:view",
			"meetings:view",
			"minutes:view",
			"attendance:view",
			"contributions:view",
			"biodata:view",
		},
	}

	// Create role-permission associations
	for roleName, permNames := range rolePermissions {
		roleID := roleMap[roleName]
		var role models.Role
		if err := db.First(&role, "id = ?", roleID).Error; err != nil {
			return err
		}

		for _, permName := range permNames {
			permID := permMap[permName]
			var perm models.Permission
			if err := db.First(&perm, "id = ?", permID).Error; err != nil {
				return err
			}

			if err := db.Model(&role).Association("Permissions").Append(&perm); err != nil {
				return err
			}
		}
	}

	return nil
}

package rbac

// Only keep super admin as a built-in role
// All other roles and permissions are managed in the database

type Role string

const (
	RoleSuperAdmin Role = "super_admin"
)

type Permission string

// These are just constants for reference in code
// Actual permissions are stored in the database
const (
	// Family permissions
	PermViewFamily Permission = "family:view"
	PermEditFamily Permission = "family:edit"
	PermViewStats  Permission = "family:view_stats"

	// Member permissions
	PermViewMembers   Permission = "members:view"
	PermAddMembers    Permission = "members:add"
	PermEditMembers   Permission = "members:edit"
	PermRemoveMembers Permission = "members:remove"
	PermApproveMember Permission = "members:approve"

	// Meeting permissions
	PermViewMeetings    Permission = "meetings:view"
	PermScheduleMeeting Permission = "meetings:schedule"
	PermEditMeetings    Permission = "meetings:edit"
	PermDeleteMeetings  Permission = "meetings:delete"

	// Minutes permissions
	PermViewMinutes    Permission = "minutes:view"
	PermPublishMinutes Permission = "minutes:publish"
	PermDeleteMinutes  Permission = "minutes:delete"
	PermAcknowMinutes  Permission = "minutes:acknowledge"

	// Attendance permissions
	PermViewAttendance   Permission = "attendance:view"
	PermRecordAttendance Permission = "attendance:record"
	PermEditAttendance   Permission = "attendance:edit"

	// Contribution permissions
	PermViewContributions   Permission = "contributions:view"
	PermRecordContribution  Permission = "contributions:record"
	PermConfirmContribution Permission = "contributions:confirm"
	PermEditContribution    Permission = "contributions:edit"

	// Admin permissions
	PermManageFamilies    Permission = "admin:manage_families"
	PermManageUsers       Permission = "admin:manage_users"
	PermManageRoles       Permission = "admin:manage_roles"
	PermManageDepartments Permission = "admin:manage_departments"
	PermViewAuditLog      Permission = "admin:view_audit_log"
	PermSystemSettings    Permission = "admin:system_settings"
)

// IsSuperAdmin is a quick check for super admin role
func IsSuperAdmin(role string) bool {
	return role == string(RoleSuperAdmin)
}

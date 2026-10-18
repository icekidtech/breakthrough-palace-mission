package rbac

// Role defines family member roles
type Role string

const (
	RoleCoordinator Role = "coordinator"
	RoleSecretary   Role = "secretary"
	RoleWelfareLead Role = "welfare_lead"
	RoleMember      Role = "member"
)

// Permission defines what actions can be performed
type Permission string

const (
	// Family permissions
	PermViewFamily   Permission = "family:view"
	PermEditFamily   Permission = "family:edit"
	PermViewStats    Permission = "family:view_stats"

	// Member permissions
	PermViewMembers   Permission = "members:view"
	PermAddMembers    Permission = "members:add"
	PermEditMembers   Permission = "members:edit"
	PermRemoveMembers Permission = "members:remove"

	// Meeting permissions
	PermViewMeetings   Permission = "meetings:view"
	PermScheduleMeeting Permission = "meetings:schedule"
	PermEditMeetings   Permission = "meetings:edit"
	PermDeleteMeetings Permission = "meetings:delete"

	// Attendance permissions
	PermViewAttendance   Permission = "attendance:view"
	PermRecordAttendance Permission = "attendance:record"
	PermEditAttendance   Permission = "attendance:edit"

	// Contribution permissions
	PermViewContributions   Permission = "contributions:view"
	PermRecordContribution  Permission = "contributions:record"
	PermConfirmContribution Permission = "contributions:confirm"
	PermEditContribution    Permission = "contributions:edit"

	// Minutes permissions
	PermViewMinutes   Permission = "minutes:view"
	PermPublishMinutes Permission = "minutes:publish"
	PermDeleteMinutes Permission = "minutes:delete"
)

// RolePermissions maps each role to its permissions
var RolePermissions = map[Role][]Permission{
	RoleCoordinator: {
		// Family
		PermViewFamily,
		PermEditFamily,
		PermViewStats,

		// Members
		PermViewMembers,
		PermAddMembers,
		PermEditMembers,
		PermRemoveMembers,

		// Meetings
		PermViewMeetings,
		PermScheduleMeeting,
		PermEditMeetings,
		PermDeleteMeetings,

		// Attendance
		PermViewAttendance,
		PermRecordAttendance,
		PermEditAttendance,

		// Contributions
		PermViewContributions,
		PermRecordContribution,
		PermConfirmContribution,
		PermEditContribution,

		// Minutes
		PermViewMinutes,
		PermPublishMinutes,
		PermDeleteMinutes,
	},

	RoleSecretary: {
		// Family
		PermViewFamily,
		PermViewStats,

		// Members
		PermViewMembers,
		PermAddMembers,
		PermEditMembers,

		// Meetings
		PermViewMeetings,
		PermScheduleMeeting,
		PermEditMeetings,

		// Attendance
		PermViewAttendance,
		PermRecordAttendance,

		// Contributions
		PermViewContributions,
		PermRecordContribution,

		// Minutes
		PermViewMinutes,
		PermPublishMinutes,
	},

	RoleWelfareLead: {
		// Family
		PermViewFamily,
		PermViewStats,

		// Members
		PermViewMembers,

		// Meetings
		PermViewMeetings,

		// Attendance
		PermViewAttendance,
		PermRecordAttendance,

		// Contributions
		PermViewContributions,
		PermRecordContribution,
		PermConfirmContribution,

		// Minutes
		PermViewMinutes,
	},

	RoleMember: {
		// Family
		PermViewFamily,

		// Members
		PermViewMembers,

		// Meetings
		PermViewMeetings,

		// Attendance
		PermViewAttendance,

		// Contributions
		PermViewContributions,
		PermRecordContribution,

		// Minutes
		PermViewMinutes,
	},
}

// HasPermission checks if a role has a specific permission
func HasPermission(role Role, permission Permission) bool {
	permissions, exists := RolePermissions[role]
	if !exists {
		return false
	}

	for _, p := range permissions {
		if p == permission {
			return true
		}
	}
	return false
}

package shared

var ROLES = map[string]string{
	"ADMIN":    "admin",
	"USER":     "user",
	"SELLER":   "seller",
	"APPROVER": "approver",
}

var ADMIN_ASSIGNED_ROLES = []string{
	ROLES["SELLER"],
	ROLES["APPROVER"],
	ROLES["USER"],
}

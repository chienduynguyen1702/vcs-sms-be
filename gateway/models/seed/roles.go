package seed

import "github.com/chienduynguyen1702/vcs-sms-be/models"

var Roles = []models.Role{
	{
		Name:        "Admin",
		Description: "Admin role has full access to the system",
	},
	{
		Name:        "User",
		Description: "User role has limited access to the system",
	},
}

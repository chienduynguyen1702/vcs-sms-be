package seed

import "github.com/chienduynguyen1702/vcs-sms-be/models"

var Users = []models.User{
	{
		Email:          "admin@vcs.vn",
		Username:       "admin@vcs.vn",
		OrganizationID: 1,
		Phone:          "0123456789",
		RoleID:         1,
	},
	{
		Email:          "chien@vcs.vn",
		Username:       "chien@vcs.vn",
		OrganizationID: 1,
		Phone:          "0123456789",
		RoleID:         2,
	},
}

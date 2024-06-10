package seed

import "github.com/chienduynguyen1702/vcs-sms-be/models"

var Users = []models.User{
	{
		Email:          "admin@vcs.vn",
		Username:       "admin@vcs.vn",
		OrganizationID: 1,
		Phone:          "0123456789",
		RoleID:         1,
		Password:       "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjAyNDc2NzMsInVzZXJfaWQiOiIxIn0.SwdcD9J72jpfGlmj1aT2KbTlbSsdXZIt6ZyfncH_14Y",
	},
	{
		Email:          "chien@vcs.vn",
		Username:       "chien@vcs.vn",
		OrganizationID: 1,
		Phone:          "0123456789",
		RoleID:         2,
		Password:       "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjAyNDc2NzMsInVzZXJfaWQiOiIxIn0.SwdcD9J72jpfGlmj1aT2KbTlbSsdXZIt6ZyfncH_14Y",
	},
}

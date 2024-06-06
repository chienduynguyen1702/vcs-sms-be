package seed

import "github.com/chienduynguyen1702/vcs-sms-be/models"

var Servers = []models.Server{
	{
		Name:           "Server 1",
		IP:             "8.8.8.8",
		OrganizationID: 1,
	},
	{
		Name:           "Server 2",
		IP:             "153.252.34.54",
		OrganizationID: 1,
	},
	{
		Name:           "Server 3",
		IP:             "1.1.1.1",
		OrganizationID: 1,
	},
}

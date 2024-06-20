package services

import (
	"github.com/chienduynguyen1702/vcs-sms-be/dtos"
	"github.com/chienduynguyen1702/vcs-sms-be/repositories"
)

type MailService struct {
	userRepo      *repositories.UserRepository
	serverService *ServerService
}

func NewMailService(userRepo *repositories.UserRepository, serverRepo *ServerService) *MailService {
	return &MailService{userRepo, serverRepo}
}

func (ms *MailService) GetMailInfoToSend(from, to string) (dtos.MailBody, error) {
	count, onlineServer, offlineServer, apus, err := ms.serverService.GetMailData(from, to)
	if err != nil {
		return dtos.MailBody{}, err
	}

	adminUser, err := ms.userRepo.GetAdminUser()
	if err != nil {
		return dtos.MailBody{}, err
	}
	var adminMails []string
	for _, user := range adminUser {
		adminMails = append(adminMails, user.Email)
	}
	var mailBody dtos.MailBody
	mailBody.AdminMails = adminMails
	mailBody.From = from
	mailBody.To = to
	mailBody.TotalServer = count
	mailBody.TotalServerOnline = onlineServer
	mailBody.TotalServerOffline = offlineServer
	mailBody.AvgUptime = float64(apus)

	return mailBody, nil
}

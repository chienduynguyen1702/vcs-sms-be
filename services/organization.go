package services

import (
	"github.com/chienduynguyen1702/vcs-sms-be/models"
	"github.com/chienduynguyen1702/vcs-sms-be/repositories"
)

type IOrganizationService interface {
	CreateOrganization(organization *models.Organization) error
	GetOrganizationByEmail(email string) (*models.Organization, error)
}

type OrganizationService struct {
	organizationRepo *repositories.OrganizationRepository
}

func NewOrganizationService(organizationRepo *repositories.OrganizationRepository) *OrganizationService {
	return &OrganizationService{organizationRepo: organizationRepo}
}

func (us *OrganizationService) CreateOrganization(organization *models.Organization) error {
	return us.organizationRepo.CreateOrganization(organization)
}

func (us *OrganizationService) GetOrganizationByName(email string) *models.Organization {
	return us.organizationRepo.GetOrganizationByName(email)
}

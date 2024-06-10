package services

import (
	"github.com/chienduynguyen1702/vcs-sms-be/dtos"
	"github.com/chienduynguyen1702/vcs-sms-be/models"
	"github.com/chienduynguyen1702/vcs-sms-be/repositories"
	"github.com/chienduynguyen1702/vcs-sms-be/utilities"
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

func (us *OrganizationService) CreateOrganization(organization *models.Organization) (*models.Organization, error) {
	return us.organizationRepo.CreateOrganization(organization)
}

func (us *OrganizationService) GetOrganizationByName(email string) *models.Organization {
	return us.organizationRepo.GetOrganizationByName(email)
}

func (us *OrganizationService) GetOrganizationByID(ID string) (dtos.OrganizationResponse, error) {
	org, err := us.organizationRepo.GetOrganizationByID(ID)
	if err != nil {
		return dtos.OrganizationResponse{}, err
	}
	return dtos.MakeOrganizationResponse(*org), nil
}

func (us *OrganizationService) UpdateOrganization(organization *dtos.UpdateOrganizationRequest, orgID string) (dtos.OrganizationResponse, error) {
	org, err := us.organizationRepo.GetOrganizationByID(orgID)
	if err != nil {
		return dtos.OrganizationResponse{}, err
	}
	parsedEstablishmentDate := utilities.ParseStringToDate(organization.EstablishmentDate)

	// Update organization
	org.Name = organization.Name
	org.AliasName = organization.AliasName
	org.EstablishmentDate = parsedEstablishmentDate
	org.Description = organization.Description
	org.Address = organization.Address
	org.WebsiteURL = organization.WebsiteURL

	err = us.organizationRepo.UpdateOrganization(org)
	if err != nil {
		return dtos.OrganizationResponse{}, err
	}

	return dtos.MakeOrganizationResponse(*org), nil
}

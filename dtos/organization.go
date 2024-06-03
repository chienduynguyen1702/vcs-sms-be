package dtos

import (
	"time"

	"github.com/chienduynguyen1702/vcs-sms-be/models"
)

type OrganizationResponse struct {
	ID                uint      `json:"id"`
	Name              string    `json:"name"`
	AliasName         string    `json:"alias_name"`
	EstablishmentDate time.Time `json:"establishment_date"`
	Description       string    `json:"description"`
	Address           string    `json:"address"`
	WebsiteURL        string    `json:"website_url"`
}

func MakeOrganizationResponse(organization models.Organization) OrganizationResponse {
	return OrganizationResponse{
		ID:                organization.ID,
		Name:              organization.Name,
		AliasName:         organization.AliasName,
		EstablishmentDate: organization.EstablishmentDate,
		Description:       organization.Description,
		Address:           organization.Address,
		WebsiteURL:        organization.WebsiteURL,
	}
}

type UpdateOrganizationRequest struct {
	Name              string `json:"name"`
	AliasName         string `json:"alias_name"`
	EstablishmentDate string `json:"establishment_date"`
	Description       string `json:"description"`
	Address           string `json:"address"`
	WebsiteURL        string `json:"website_url"`
}

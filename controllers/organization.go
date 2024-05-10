package controllers

import "github.com/chienduynguyen1702/vcs-sms-be/services"

type OrganizationController struct {
	organizationService *services.OrganizationService
}

func NewOrganizationController(organizationService *services.OrganizationService) *OrganizationController {
	return &OrganizationController{organizationService}
}

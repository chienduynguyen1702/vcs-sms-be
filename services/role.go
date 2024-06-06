package services

import (
	"github.com/chienduynguyen1702/vcs-sms-be/dtos"
	"github.com/chienduynguyen1702/vcs-sms-be/models"
	"github.com/chienduynguyen1702/vcs-sms-be/repositories"
)

type IRoleService interface {
	GetRole() ([]models.Role, error)
}

type RoleService struct {
	roleRepo *repositories.RoleRepository
}

func NewRoleService(roleRepo *repositories.RoleRepository) *RoleService {
	return &RoleService{roleRepo}
}

func (rs *RoleService) GetRoles() ([]dtos.RoleResponse, error) {

	roles, err := rs.roleRepo.GetRoles()
	if err != nil {
		return nil, err
	}

	var roleResponses []dtos.RoleResponse
	for _, role := range roles {
		roleResponses = append(roleResponses, dtos.MakeRoleResponse(role))
	}

	return roleResponses, nil
}

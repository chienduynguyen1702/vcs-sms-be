package controllers

import (
	"net/http"

	"github.com/chienduynguyen1702/vcs-sms-be/dtos"
	"github.com/chienduynguyen1702/vcs-sms-be/services"
	"github.com/gin-gonic/gin"
)

type OrganizationController struct {
	organizationService *services.OrganizationService
}

func NewOrganizationController(organizationService *services.OrganizationService) *OrganizationController {
	return &OrganizationController{organizationService}
}

// GetOrganization godoc
// @Summary Get organization
// @Description Get organization
// @Tags Organization
// @Accept  json
// @Produce  json
// @Success 200 {object} string
// @Router /api/v1/organizations/ [get]
func (oc *OrganizationController) GetOrganization(ctx *gin.Context) {
	orgID, exist := ctx.Get("orgID")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, dtos.ErrorResponse("Failed to get organizationID in context"))
		return
	}
	organizationDTOResponse, err := oc.organizationService.GetOrganizationByID(orgID.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, dtos.SuccessResponse("Get organization successfully", organizationDTOResponse))
}

// UpdateOrganization godoc
// @Summary Update organization
// @Description Update organization
// @Tags Organization
// @Accept  json
// @Produce  json
// @Param UpdateOrganizationBodyRequest body dtos.UpdateOrganizationRequest true "Update Organization Request"
// @Success 200 {object} string
// @Router /api/v1/organizations/ [put]
func (oc *OrganizationController) UpdateOrganization(ctx *gin.Context) {
	orgID, exist := ctx.Get("orgID")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, dtos.ErrorResponse("Failed to get organizationID in context"))
		return
	}
	updateOrganizationRequest := &dtos.UpdateOrganizationRequest{}
	if err := ctx.ShouldBindJSON(updateOrganizationRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse(err.Error()))
		return
	}
	updateOrgDToResponse, err := oc.organizationService.UpdateOrganization(updateOrganizationRequest, orgID.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, dtos.SuccessResponse("Update organization successfully", updateOrgDToResponse))
}

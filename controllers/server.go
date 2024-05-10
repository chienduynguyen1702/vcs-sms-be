package controllers

import (
	"net/http"

	"github.com/chienduynguyen1702/vcs-sms-be/dtos"
	"github.com/chienduynguyen1702/vcs-sms-be/services"
	"github.com/gin-gonic/gin"
)

type ServerController struct {
	serverService *services.ServerService
}

func NewServerController(serverService *services.ServerService) *ServerController {
	return &ServerController{serverService}
}

// CreateServer godoc
// @Summary Create a new server
// @Description Create a new server
// @Tags Server
// @Accept  json
// @Produce  json
// @Param CreateServerBodyRequest body dtos.CreateServerRequest true "Create Server Request"
// @Success 200 {object} string
// @Router /api/v1/servers [post]
func (sc *ServerController) CreateServer(ctx *gin.Context) {
	newServer := &dtos.CreateServerRequest{}
	if err := ctx.ShouldBindJSON(newServer); err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse(err.Error()))
		return
	}
	orgID, exist := ctx.Get("orgID")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, dtos.ErrorResponse("Failed to get organizationID in context"))
		return
	}
	err := sc.serverService.CreateServer(newServer, orgID.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, dtos.SuccessResponse(
		"Create server successfully",
		dtos.ServerResponse{
			Name: newServer.Name,
			IP:   newServer.IP,
		},
	))
}

// GetServers godoc
// @Summary Get all servers
// @Description Get all servers
// @Tags Server
// @Accept  json
// @Produce  json
// @Success 200 {object} string
// @Router /api/v1/servers [get]
func (sc *ServerController) GetServers(ctx *gin.Context) {
	orgID, exist := ctx.Get("orgID")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, dtos.ErrorResponse("Failed to get organizationID in context"))
		return
	}
	servers, err := sc.serverService.GetServers(orgID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, dtos.SuccessResponse("Get all servers successfully", servers))
}

// GetServerByID godoc
// @Summary Get server by ID
// @Description Get server by ID
// @Tags Server
// @Accept  json
// @Produce  json
// @Param id path string true "Server ID"
// @Success 200 {object} string
// @Router /api/v1/servers/{id} [get]
func (sc *ServerController) GetServerByID(ctx *gin.Context) {
	id := ctx.Param("id")
	server, err := sc.serverService.GetServerByID(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, dtos.SuccessResponse("Get server by ID successfully", server))
}

// UpdateServer godoc
// @Summary Update server
// @Description Update server
// @Tags Server
// @Accept  json
// @Produce  json
// @Param id path string true "Server ID"
// @Param UpdateServerBodyRequest body dtos.UpdateServerRequest true "Update Server Request"
// @Success 200 {object} string
// @Router /api/v1/servers/{id} [put]
func (sc *ServerController) UpdateServer(ctx *gin.Context) {
	id := ctx.Param("id")
	updateServer := &dtos.UpdateServerRequest{}
	if err := ctx.ShouldBindJSON(updateServer); err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse(err.Error()))
		return
	}
	serverResponse, err := sc.serverService.UpdateServer(id, updateServer)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, dtos.SuccessResponse("Update server successfully", serverResponse))
}

// DeleteServer godoc
// @Summary Delete server
// @Description Delete server
// @Tags Server
// @Accept  json
// @Produce  json
// @Param id path string true "Server ID"
// @Success 200 {object} string
// @Router /api/v1/servers/{id} [delete]
func (sc *ServerController) DeleteServer(ctx *gin.Context) {
	id := ctx.Param("id")
	err := sc.serverService.DeleteServer(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, dtos.SuccessResponse("Delete server successfully", nil))
}

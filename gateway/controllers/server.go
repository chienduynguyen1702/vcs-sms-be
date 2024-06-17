package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/chienduynguyen1702/vcs-sms-be/dtos"
	"github.com/chienduynguyen1702/vcs-sms-be/services"
	"github.com/chienduynguyen1702/vcs-sms-be/utilities"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
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
// @Param page 		query int false "Page"
// @Param limit 	query int false "Limit"
// @Param search	query string false "Search"
// @Success 200 {object} string
// @Router /api/v1/servers [get]
func (sc *ServerController) GetServers(ctx *gin.Context) {
	page := ctx.Query("page")
	limit := ctx.Query("limit")
	search := ctx.Query("search")
	// parse page and limit to int
	pageInt, limitInt := utilities.ParsePageAndLimit(page, limit)
	orgId, exist := ctx.Get("orgID")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, dtos.ErrorResponse("Failed to get organizationID in context"))
		return
	}
	orgIDStr := orgId.(string)
	count, servers, err := sc.serverService.GetServers(orgIDStr, search, pageInt, limitInt)
	if err != nil {
		fmt.Println("err", err.Error())
		ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, dtos.SuccessResponse("Get all servers successfully", dtos.PaginationResponse(count, servers)))
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

// GetArchivedServer godoc
// @Summary GetArchivedServer server
// @Description Archive server
// @Tags Server
// @Accept  json
// @Produce  json
// @Success 200 {object} string
// @Router /api/v1/servers/archived [get]
func (sc *ServerController) GetArchivedServer(ctx *gin.Context) {
	orgID, exist := ctx.Get("orgID")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, dtos.ErrorResponse("Failed to get organizationID in context"))
		return
	}
	servers, err := sc.serverService.GetArchivedServers(orgID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, dtos.SuccessResponse("Get all archived servers successfully", servers))
}

// ArchiveServer godoc
// @Summary Archive server
// @Description Archive server
// @Tags Server
// @Accept  json
// @Produce  json
// @Param id path string true "Server ID"
// @Success 200 {object} string
// @Router /api/v1/servers/{id}/archive [patch]
func (sc *ServerController) ArchiveServer(ctx *gin.Context) {
	// get server id from param
	id := ctx.Param("id")

	// get admin id from context to set as archiver
	adminID, exist := ctx.Get("userID")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, dtos.ErrorResponse("Failed to get userID in context"))
		return
	}
	adminIDStr := adminID.(string)
	adminIDUint, err := utilities.ParseStringToUint(adminIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse(err.Error()))
		return
	}
	err = sc.serverService.ArchiveServer(id, adminIDUint)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, dtos.SuccessResponse("Archive server successfully", nil))
}

// UnarchiveServer godoc
// @Summary Archive server
// @Description Archive server
// @Tags Server
// @Accept  json
// @Produce  json
// @Param id path string true "Server ID"
// @Success 200 {object} string
// @Router /api/v1/servers/{id}/restore [patch]
func (sc *ServerController) Restore(ctx *gin.Context) {
	// get server id from param
	id := ctx.Param("id")

	err := sc.serverService.RestoreServer(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, dtos.SuccessResponse("Restored server successfully", nil))
}

// DownloadTemplate godoc
// @Summary Download template
// @Description Send to client file template from ./files/server_list_template.xlsx
// @Tags Server
// @Accept  json
// @Produce  json
// @Success 200 {object} string
// @Router /api/v1/servers/download-template [get]
func (sc *ServerController) DownloadTemplate(ctx *gin.Context) {
	ctx.File("./files/server_list_template.xlsx")
}

// UploadServerList godoc
// @Summary Upload server list
// @Description Upload server list from client .xlsx file
// @Tags Server
// @Accept  json
// @Produce  json
// @Param file formData file true "Server list file"
// @Success 200 {object} string
// @Router /api/v1/servers/upload [post]
func (sc *ServerController) UploadServerList(ctx *gin.Context) {
	uploadFile, err := ctx.FormFile("uploadFile")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse(err.Error()))
		return
	}
	orgID, exist := ctx.Get("orgID")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, dtos.ErrorResponse("Failed to get organizationID in context"))
		return
	}

	// Save the file to disk
	err = ctx.SaveUploadedFile(uploadFile, "files/"+uploadFile.Filename)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse(err.Error()))
		return
	}
	// Load data from the Excel file
	xlsxFile, err := excelize.OpenFile("files/" + uploadFile.Filename)
	if err != nil {
		// err = os.Remove("files/" + uploadFile.Filename)
		ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse(err.Error()))
		return
	}
	// remove the file after reading
	err = os.Remove("files/" + uploadFile.Filename)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse(err.Error()))
		return
	}

	// Get sheet name
	sheetName := xlsxFile.GetSheetName(0)
	// Get the values from the Excel sheet
	rows, err := xlsxFile.GetRows(sheetName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse(err.Error()))
		return
	}
	var serverList []dtos.CreateServerRequest
	for i, row := range rows {
		log.Println("row", row)
		if i == 0 {
			continue
		}
		des := ""
		if len(row) < 3 {
			des = ""
		} else {
			des = row[2]
		}
		server := dtos.CreateServerRequest{
			IP:          row[0],
			Name:        row[1],
			Description: des,
			Status:      "Unknown",
		}
		serverList = append(serverList, server)
	}
	// Print the data
	// fmt.Println(serverList)
	if len(serverList) == 0 {
		ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse("No data in the file"))
		return
	}

	updatedCount, createdCount, err := sc.serverService.UploadServerList(serverList, orgID.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, dtos.SuccessResponse(fmt.Sprintf("uploaded %d servers, created %d servers.", updatedCount, createdCount), nil))
}

func (sc *ServerController) FlushCache(ctx *gin.Context) {
	// fmt.Println("flush cache in controller")
	sc.serverService.FlushCache()
}

// SendReportByMail godoc
// @Summary Send report by mail
// @Description Send report by mail
// @Tags Server
// @Accept  json
// @Produce  json
// @Param SendMailRequest body dtos.SendMailRequest true "Send Report by Mail Request Date is YYYY-MM-DDThh:mm:ss.000Z"
// @Success 200 {object} string
// @Router /api/v1/servers/send-report [post]
func (sc *ServerController) SendReportByMail(ctx *gin.Context) {
	sendMailRequest := &dtos.SendMailRequest{}
	if err := ctx.ShouldBindJSON(sendMailRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse(err.Error()))
		return
	}

	err := sc.serverService.SendReportByMail(sendMailRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, dtos.SuccessResponse("Send report by mail successfully", nil))
}

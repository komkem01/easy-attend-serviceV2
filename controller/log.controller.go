package controller

import (
	"easy-attend-service/models"
	"easy-attend-service/requests"
	"easy-attend-service/response"
	"easy-attend-service/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LogController struct {
	logService *services.LogService
}

func NewLogController() *LogController {
	return &LogController{
		logService: services.NewLogService(),
	}
}

func (lc *LogController) GetAllLogs(c *gin.Context) {
	logs, err := lc.logService.GetAllLogs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to fetch logs", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Logs retrieved successfully", logs))
}

func (lc *LogController) GetLogByID(c *gin.Context) {
	id := c.Param("id")
	// Convert id from string to uint
	var logID uint
	if _, err := fmt.Sscanf(id, "%d", &logID); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid log ID", err.Error()))
		return
	}

	log, err := lc.logService.GetLogByID(logID)
	if err != nil {
		if err.Error() == "log not found" {
			c.JSON(http.StatusNotFound, response.ErrorResponse("Log not found", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to fetch log", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Log retrieved successfully", log))
}

func (lc *LogController) GetLogsByTeacher(c *gin.Context) {
	teacherID := c.Param("teacher_id")
	logs, err := lc.logService.GetLogsByTeacher(teacherID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to fetch logs", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Logs retrieved successfully", logs))
}

func (lc *LogController) GetLogsByAction(c *gin.Context) {
	actionParam := c.Param("action")
	action := models.LogAction(actionParam)
	logs, err := lc.logService.GetLogsByAction(action)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to fetch logs", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Logs retrieved successfully", logs))
}

func (lc *LogController) CreateLog(c *gin.Context) {
	var req requests.LogCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid request data", err.Error()))
		return
	}

	log, err := lc.logService.CreateLog(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to create log", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.SuccessResponse("Log created successfully", log))
}

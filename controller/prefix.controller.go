package controller

import (
	"easy-attend-service/requests"
	"easy-attend-service/response"
	"easy-attend-service/services"
	"easy-attend-service/utils/logger"
	"net/http"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type PrefixController struct {
	prefixService *services.PrefixService
}

func NewPrefixController() *PrefixController {
	return &PrefixController{
		prefixService: services.NewPrefixService(),
	}
}

// GetAllPrefixes godoc
// @Summary Get all prefixes
// @Description Get a list of all prefixes
// @Tags prefixes
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response
// @Failure 500 {object} response.StatusResponse
// @Router /prefixes [get]
func (pc *PrefixController) GetAllPrefixes(c *gin.Context) {
	logger.LogInfo("Fetching all prefixes", logrus.Fields{})

	prefixes, err := pc.prefixService.GetAllPrefixes()
	if err != nil {
		logger.LogError(err, "Failed to fetch prefixes", logrus.Fields{})
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to fetch prefixes", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Prefixes retrieved successfully", prefixes))
}

// GetPrefixByID godoc
// @Summary Get prefix by ID
// @Description Get a specific prefix by its ID
// @Tags prefixes
// @Accept json
// @Produce json
// @Param id path string true "Prefix ID"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response
// @Failure 404 {object} response.StatusResponse
// @Failure 500 {object} response.StatusResponse
// @Router /prefixes/{id} [get]
func (pc *PrefixController) GetPrefixByID(c *gin.Context) {
	id := c.Param("id")

	logger.LogInfo("Fetching prefix by ID", logrus.Fields{
		"prefix_id": id,
	})

	// Convert id from string to uint
	var uintID uint
	if _, err := fmt.Sscanf(id, "%d", &uintID); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid prefix ID", err.Error()))
		return
	}

	prefix, err := pc.prefixService.GetPrefixByID(uintID)
	if err != nil {
		if err.Error() == "prefix not found" {
			c.JSON(http.StatusNotFound, response.ErrorResponse("Prefix not found", err.Error()))
			return
		}
		logger.LogError(err, "Failed to fetch prefix", logrus.Fields{
			"prefix_id": id,
		})
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to fetch prefix", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Prefix retrieved successfully", prefix))
}

// CreatePrefix godoc
// @Summary Create a new prefix
// @Description Create a new prefix
// @Tags prefixes
// @Accept json
// @Produce json
// @Param prefix body requests.PrefixCreateRequest true "Prefix create request"
// @Security ApiKeyAuth
// @Success 201 {object} response.Response
// @Failure 400 {object} response.StatusResponse
// @Failure 500 {object} response.StatusResponse
// @Router /prefixes [post]
func (pc *PrefixController) CreatePrefix(c *gin.Context) {
	var req requests.PrefixCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error(), "Invalid request body"))
		return
	}

	logger.LogInfo("Creating new prefix", logrus.Fields{
		"name": req.Name,
	})

	prefix, err := pc.prefixService.CreatePrefix(&req)
	if err != nil {
		if err.Error() == "prefix with this name already exists" {
			c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error(), "Duplicate prefix name"))
			return
		}
		logger.LogError(err, "Failed to create prefix", logrus.Fields{
			"name": req.Name,
		})
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to create prefix", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.SuccessResponse("Prefix created successfully", prefix))
}

// UpdatePrefix godoc
// @Summary Update a prefix
// @Description Update an existing prefix by ID
// @Tags prefixes
// @Accept json
// @Produce json
// @Param id path string true "Prefix ID"
// @Param prefix body requests.PrefixUpdateRequest true "Prefix update request"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response
// @Failure 400 {object} response.StatusResponse
// @Failure 404 {object} response.StatusResponse
// @Failure 500 {object} response.StatusResponse
// @Router /prefixes/{id} [put]
func (pc *PrefixController) UpdatePrefix(c *gin.Context) {
	id := c.Param("id")
	var req requests.PrefixUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error(), "Invalid request body"))
		return
	}

	// Convert id from string to uint
	var uintID uint
	if _, err := fmt.Sscanf(id, "%d", &uintID); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid prefix ID", err.Error()))
		return
	}

	logger.LogInfo("Updating prefix", logrus.Fields{
		"prefix_id": id,
		"name":      req.Name,
	})

	prefix, err := pc.prefixService.UpdatePrefix(uintID, &req)
	if err != nil {
		if err.Error() == "prefix not found" {
			c.JSON(http.StatusNotFound, response.ErrorResponse("Prefix not found", err.Error()))
			return
		}
		if err.Error() == "prefix with this name already exists" {
			c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error(), "Duplicate prefix name"))
			return
		}
		logger.LogError(err, "Failed to update prefix", logrus.Fields{
			"prefix_id": id,
		})
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to update prefix", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Prefix updated successfully", prefix))
}

// DeletePrefix godoc
// @Summary Delete a prefix
// @Description Delete a prefix by ID (soft delete)
// @Tags prefixes
// @Accept json
// @Produce json
// @Param id path string true "Prefix ID"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response
// @Failure 404 {object} response.StatusResponse
// @Failure 500 {object} response.StatusResponse
// @Router /prefixes/{id} [delete]
func (pc *PrefixController) DeletePrefix(c *gin.Context) {
	id := c.Param("id")

	logger.LogInfo("Deleting prefix", logrus.Fields{
		"prefix_id": id,
	})

	// Convert id from string to uint
	var uintID uint
	if _, err := fmt.Sscanf(id, "%d", &uintID); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid prefix ID", err.Error()))
		return
	}

	err := pc.prefixService.DeletePrefix(uintID)
	if err != nil {
		if err.Error() == "prefix not found" {
			c.JSON(http.StatusNotFound, response.ErrorResponse("Prefix not found", err.Error()))
			return
		}
		logger.LogError(err, "Failed to delete prefix", logrus.Fields{
			"prefix_id": id,
		})
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to delete prefix", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Prefix deleted successfully", nil))
}

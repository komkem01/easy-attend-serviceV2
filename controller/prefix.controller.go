package controller

import (
	"easy-attend-service/requests"
	"easy-attend-service/response"
	"easy-attend-service/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PrefixController struct {
	prefixService *services.PrefixService
}

func NewPrefixController() *PrefixController {
	return &PrefixController{
		prefixService: services.NewPrefixService(),
	}
}

func (pc *PrefixController) GetAllPrefixes(c *gin.Context) {
	prefixes, err := pc.prefixService.GetAllPrefixes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to fetch prefixes", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Prefixes retrieved successfully", prefixes))
}

func (pc *PrefixController) GetPrefixByID(c *gin.Context) {
	id := c.Param("id")
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
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to fetch prefix", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Prefix retrieved successfully", prefix))
}

func (pc *PrefixController) CreatePrefix(c *gin.Context) {
	var req requests.PrefixCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error(), "Invalid request body"))
		return
	}
	prefix, err := pc.prefixService.CreatePrefix(&req)
	if err != nil {
		if err.Error() == "prefix with this name already exists" {
			c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error(), "Duplicate prefix name"))
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to create prefix", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.SuccessResponse("Prefix created successfully", prefix))
}

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
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to update prefix", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Prefix updated successfully", prefix))
}

func (pc *PrefixController) DeletePrefix(c *gin.Context) {
	id := c.Param("id")
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
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to delete prefix", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Prefix deleted successfully", nil))
}

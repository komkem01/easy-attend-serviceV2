package controller

import (
	"easy-attend-service/requests"
	"easy-attend-service/response"
	"easy-attend-service/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GenderController struct {
	genderService *services.GenderService
}

func NewGenderController() *GenderController {
	return &GenderController{
		genderService: services.NewGenderService(),
	}
}

func (gc *GenderController) GetAllGenders(c *gin.Context) {
	genders, err := gc.genderService.GetAllGenders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to fetch genders", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Genders retrieved successfully", genders))
}

func (gc *GenderController) GetGenderByID(c *gin.Context) {
	id := c.Param("id")
	// Convert id from string to uint
	var genderID uint
	if _, err := fmt.Sscanf(id, "%d", &genderID); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid gender ID", err.Error()))
		return
	}

	gender, err := gc.genderService.GetGenderByID(genderID)
	if err != nil {
		if err.Error() == "gender not found" {
			c.JSON(http.StatusNotFound, response.ErrorResponse("Gender not found", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to fetch gender", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Gender retrieved successfully", gender))
}

func (gc *GenderController) CreateGender(c *gin.Context) {
	var req requests.GenderCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error(), "Invalid request body"))
		return
	}
	gender, err := gc.genderService.CreateGender(&req)
	if err != nil {
		if err.Error() == "gender with this name already exists" {
			c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error(), "Duplicate gender name"))
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to create gender", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.SuccessResponse("Gender created successfully", gender))
}

func (gc *GenderController) UpdateGender(c *gin.Context) {
	id := c.Param("id")
	var req requests.GenderUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error(), "Invalid request body"))
		return
	}

	// Convert id from string to uint
	var genderID uint
	if _, err := fmt.Sscanf(id, "%d", &genderID); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid gender ID", err.Error()))
		return
	}
	gender, err := gc.genderService.UpdateGender(genderID, &req)
	if err != nil {
		if err.Error() == "gender not found" {
			c.JSON(http.StatusNotFound, response.ErrorResponse("Gender not found", err.Error()))
			return
		}
		if err.Error() == "gender with this name already exists" {
			c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error(), "Duplicate gender name"))
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to update gender", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Gender updated successfully", gender))
}

func (gc *GenderController) DeleteGender(c *gin.Context) {
	id := c.Param("id")
	// Convert id from string to uint
	var genderID uint
	if _, err := fmt.Sscanf(id, "%d", &genderID); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid gender ID", err.Error()))
		return
	}

	err := gc.genderService.DeleteGender(genderID)
	if err != nil {
		if err.Error() == "gender not found" {
			c.JSON(http.StatusNotFound, response.ErrorResponse("Gender not found", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to delete gender", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Gender deleted successfully", nil))
}

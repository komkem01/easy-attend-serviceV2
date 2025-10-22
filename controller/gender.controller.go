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

type GenderController struct {
	genderService *services.GenderService
}

func NewGenderController() *GenderController {
	return &GenderController{
		genderService: services.NewGenderService(),
	}
}

// GetAllGenders godoc
// @Summary Get all genders
// @Description Get a list of all genders
// @Tags genders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response
// @Failure 500 {object} response.StatusResponse
// @Router /genders [get]
func (gc *GenderController) GetAllGenders(c *gin.Context) {
	logger.LogInfo("Fetching all genders", logrus.Fields{})

	genders, err := gc.genderService.GetAllGenders()
	if err != nil {
		logger.LogError(err, "Failed to fetch genders", logrus.Fields{})
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to fetch genders", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Genders retrieved successfully", genders))
}

// GetGenderByID godoc
// @Summary Get gender by ID
// @Description Get a specific gender by its ID
// @Tags genders
// @Accept json
// @Produce json
// @Param id path string true "Gender ID"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response
// @Failure 404 {object} response.StatusResponse
// @Failure 500 {object} response.StatusResponse
// @Router /genders/{id} [get]
func (gc *GenderController) GetGenderByID(c *gin.Context) {
	id := c.Param("id")

	logger.LogInfo("Fetching gender by ID", logrus.Fields{
		"gender_id": id,
	})

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
		logger.LogError(err, "Failed to fetch gender", logrus.Fields{
			"gender_id": id,
		})
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to fetch gender", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Gender retrieved successfully", gender))
}

// CreateGender godoc
// @Summary Create a new gender
// @Description Create a new gender
// @Tags genders
// @Accept json
// @Produce json
// @Param gender body requests.GenderCreateRequest true "Gender create request"
// @Security ApiKeyAuth
// @Success 201 {object} response.Response
// @Failure 400 {object} response.StatusResponse
// @Failure 500 {object} response.StatusResponse
// @Router /genders [post]
func (gc *GenderController) CreateGender(c *gin.Context) {
	var req requests.GenderCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error(), "Invalid request body"))
		return
	}

	logger.LogInfo("Creating new gender", logrus.Fields{
		"name": req.Name,
	})

	gender, err := gc.genderService.CreateGender(&req)
	if err != nil {
		if err.Error() == "gender with this name already exists" {
			c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error(), "Duplicate gender name"))
			return
		}
		logger.LogError(err, "Failed to create gender", logrus.Fields{
			"name": req.Name,
		})
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to create gender", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.SuccessResponse("Gender created successfully", gender))
}

// UpdateGender godoc
// @Summary Update a gender
// @Description Update an existing gender by ID
// @Tags genders
// @Accept json
// @Produce json
// @Param id path string true "Gender ID"
// @Param gender body requests.GenderUpdateRequest true "Gender update request"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response
// @Failure 400 {object} response.StatusResponse
// @Failure 404 {object} response.StatusResponse
// @Failure 500 {object} response.StatusResponse
// @Router /genders/{id} [put]
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

	logger.LogInfo("Updating gender", logrus.Fields{
		"gender_id": id,
		"name":      req.Name,
	})

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
		logger.LogError(err, "Failed to update gender", logrus.Fields{
			"gender_id": id,
		})
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to update gender", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Gender updated successfully", gender))
}

// DeleteGender godoc
// @Summary Delete a gender
// @Description Delete a gender by ID (soft delete)
// @Tags genders
// @Accept json
// @Produce json
// @Param id path string true "Gender ID"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response
// @Failure 404 {object} response.StatusResponse
// @Failure 500 {object} response.StatusResponse
// @Router /genders/{id} [delete]
func (gc *GenderController) DeleteGender(c *gin.Context) {
	id := c.Param("id")

	logger.LogInfo("Deleting gender", logrus.Fields{
		"gender_id": id,
	})

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
		logger.LogError(err, "Failed to delete gender", logrus.Fields{
			"gender_id": id,
		})
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to delete gender", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Gender deleted successfully", nil))
}

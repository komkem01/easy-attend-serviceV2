package controller

import (
	"easy-attend-service/requests"
	"easy-attend-service/response"
	"easy-attend-service/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SchoolController struct {
	schoolService *services.SchoolService
}

func NewSchoolController() *SchoolController {
	return &SchoolController{
		schoolService: services.NewSchoolService(),
	}
}

func (sc *SchoolController) GetAllSchools(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	schools, total, err := sc.schoolService.GetAllSchools(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to get schools", err.Error()))
		return
	}

	result := gin.H{
		"schools": schools,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Schools retrieved successfully", result))
}

func (sc *SchoolController) GetSchoolByID(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid request", "School ID is required"))
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid request", "School ID must be a valid number"))
		return
	}

	school, err := sc.schoolService.GetSchoolByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse("School not found", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("School retrieved successfully", school))
}

func (sc *SchoolController) CreateSchool(c *gin.Context) {
	var req requests.SchoolCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid request data", err.Error()))
		return
	}

	school, err := sc.schoolService.CreateSchool(req.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Failed to create school", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.SuccessResponse("School created successfully", school))
}

func (sc *SchoolController) UpdateSchool(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid request", "School ID is required"))
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid request", "School ID must be a valid number"))
		return
	}

	var req requests.SchoolUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid request data", err.Error()))
		return
	}

	school, err := sc.schoolService.UpdateSchool(uint(id), req.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Failed to update school", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("School updated successfully", school))
}

func (sc *SchoolController) DeleteSchool(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid request", "School ID is required"))
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid request", "School ID must be a valid number"))
		return
	}

	if err := sc.schoolService.DeleteSchool(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Failed to delete school", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("School deleted successfully", nil))
}

package controller

import (
	"easy-attend-service/requests"
	"easy-attend-service/response"
	"easy-attend-service/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TeacherController struct {
	teacherService *services.TeacherService
}

func NewTeacherController() *TeacherController {
	return &TeacherController{
		teacherService: services.NewTeacherService(),
	}
}

func (tc *TeacherController) GetAllTeachers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	teachers, total, err := tc.teacherService.GetAllTeachers(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to get teachers", err.Error()))
		return
	}

	result := gin.H{
		"teachers": teachers,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Teachers retrieved successfully", result))
}

// GetTeacherInfo gets comprehensive information for the authenticated teacher
func (tc *TeacherController) GetTeacherInfo(c *gin.Context) {
	// Get teacher ID from JWT context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse("Unauthorized", "User ID not found in token"))
		return
	}

	teacherIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse("Unauthorized", "Invalid user ID format"))
		return
	}

	teacherID, err := strconv.ParseUint(teacherIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid teacher ID", "Teacher ID must be a valid number"))
		return
	}

	info, err := tc.teacherService.GetTeacherInfo(uint(teacherID))
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse("Teacher info not found", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Teacher info retrieved successfully", info))
}

func (tc *TeacherController) GetTeacherByID(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid request", "Teacher ID is required"))
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid teacher ID", "Teacher ID must be a positive integer"))
		return
	}

	teacher, err := tc.teacherService.GetTeacherByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse("Teacher not found", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Teacher retrieved successfully", teacher))
}

func (tc *TeacherController) CreateTeacher(c *gin.Context) {
	var req requests.TeacherCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid request data", err.Error()))
		return
	}

	teacher, err := tc.teacherService.CreateTeacher(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Failed to create teacher", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.SuccessResponse("Teacher created successfully", teacher))
}

func (tc *TeacherController) UpdateTeacher(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid request", "Teacher ID is required"))
		return
	}

	var req requests.TeacherUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid request data", err.Error()))
		return
	}

	teacher, err := tc.teacherService.UpdateTeacher(id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Failed to update teacher", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Teacher updated successfully", teacher))
}

func (tc *TeacherController) DeleteTeacher(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid request", "Teacher ID is required"))
		return
	}

	if err := tc.teacherService.DeleteTeacher(id); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Failed to delete teacher", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Teacher deleted successfully", nil))
}

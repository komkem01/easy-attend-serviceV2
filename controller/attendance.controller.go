package controller

import (
	"easy-attend-service/requests"
	"easy-attend-service/response"
	"easy-attend-service/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AttendanceController struct {
	attendanceService *services.AttendanceService
}

func NewAttendanceController() *AttendanceController {
	return &AttendanceController{
		attendanceService: services.NewAttendanceService(),
	}
}

func (ac *AttendanceController) GetAllAttendances(c *gin.Context) {
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

	attendances, err := ac.attendanceService.GetAttendancesByTeacher(uint(teacherID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to fetch attendances", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Attendances retrieved successfully", attendances))
}

func (ac *AttendanceController) GetAttendanceByID(c *gin.Context) {
	idStr := c.Param("id")

	// Convert string to uint
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid attendance ID", "ID must be a valid number"))
		return
	}

	attendance, err := ac.attendanceService.GetAttendanceByID(uint(id))
	if err != nil {
		if err.Error() == "attendance not found" {
			c.JSON(http.StatusNotFound, response.ErrorResponse("Attendance not found", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to fetch attendance", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Attendance retrieved successfully", attendance))
}

func (ac *AttendanceController) GetAttendancesByClassroom(c *gin.Context) {
	classroomIDStr := c.Param("classroom_id")

	// Convert string to uint
	classroomID, err := strconv.ParseUint(classroomIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid classroom ID", "ID must be a valid number"))
		return
	}

	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}

	attendances, total, err := ac.attendanceService.GetAttendancesByClassroom(uint(classroomID), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to fetch attendances", err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Attendances retrieved successfully",
		"data":    attendances,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (ac *AttendanceController) GetAttendancesByStudent(c *gin.Context) {
	studentIDStr := c.Param("student_id")

	// Convert string to uint
	studentID, err := strconv.ParseUint(studentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid student ID", "ID must be a valid number"))
		return
	}

	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}

	attendances, total, err := ac.attendanceService.GetAttendancesByStudent(uint(studentID), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to fetch attendances", err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Attendances retrieved successfully",
		"data":    attendances,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (ac *AttendanceController) CreateAttendance(c *gin.Context) {
	var req requests.AttendanceCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid request data", err.Error()))
		return
	}

	attendance, err := ac.attendanceService.CreateAttendance(&req)
	if err != nil {
		if err.Error() == "attendance for this student on this date already exists" {
			c.JSON(http.StatusBadRequest, response.ErrorResponse("Attendance already exists", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to create attendance", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.SuccessResponse("Attendance created successfully", attendance))
}

func (ac *AttendanceController) UpdateAttendance(c *gin.Context) {
	idStr := c.Param("id")

	// Convert string to uint
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid attendance ID", "ID must be a valid number"))
		return
	}

	var req requests.AttendanceUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid request data", err.Error()))
		return
	}

	attendance, err := ac.attendanceService.UpdateAttendance(uint(id), &req)
	if err != nil {
		if err.Error() == "attendance not found" {
			c.JSON(http.StatusNotFound, response.ErrorResponse("Attendance not found", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to update attendance", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Attendance updated successfully", attendance))
}

func (ac *AttendanceController) DeleteAttendance(c *gin.Context) {
	idStr := c.Param("id")

	// Convert string to uint
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid attendance ID", "ID must be a valid number"))
		return
	}

	err = ac.attendanceService.DeleteAttendance(uint(id))
	if err != nil {
		if err.Error() == "attendance not found" {
			c.JSON(http.StatusNotFound, response.ErrorResponse("Attendance not found", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to delete attendance", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Attendance deleted successfully", nil))
}

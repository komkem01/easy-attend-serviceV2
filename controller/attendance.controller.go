package controller

import (
	"easy-attend-service/requests"
	"easy-attend-service/response"
	"easy-attend-service/services"
	"easy-attend-service/utils/logger"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
	logger.LogInfo("Fetching all attendances", logrus.Fields{})

	attendances, err := ac.attendanceService.GetAllAttendances()
	if err != nil {
		logger.LogError(err, "Failed to fetch attendances", logrus.Fields{})
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

	logger.LogInfo("Fetching attendance by ID", logrus.Fields{
		"attendance_id": idStr,
	})

	attendance, err := ac.attendanceService.GetAttendanceByID(uint(id))
	if err != nil {
		if err.Error() == "attendance not found" {
			c.JSON(http.StatusNotFound, response.ErrorResponse("Attendance not found", err.Error()))
			return
		}
		logger.LogError(err, "Failed to fetch attendance", logrus.Fields{
			"attendance_id": idStr,
		})
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

	logger.LogInfo("Fetching attendances by classroom", logrus.Fields{
		"classroom_id": classroomIDStr,
	})

	attendances, err := ac.attendanceService.GetAttendancesByClassroom(uint(classroomID))
	if err != nil {
		logger.LogError(err, "Failed to fetch attendances by classroom", logrus.Fields{
			"classroom_id": classroomIDStr,
		})
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to fetch attendances", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Attendances retrieved successfully", attendances))
}

func (ac *AttendanceController) GetAttendancesByStudent(c *gin.Context) {
	studentIDStr := c.Param("student_id")

	// Convert string to uint
	studentID, err := strconv.ParseUint(studentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid student ID", "ID must be a valid number"))
		return
	}

	logger.LogInfo("Fetching attendances by student", logrus.Fields{
		"student_id": studentIDStr,
	})

	attendances, err := ac.attendanceService.GetAttendancesByStudent(uint(studentID))
	if err != nil {
		logger.LogError(err, "Failed to fetch attendances by student", logrus.Fields{
			"student_id": studentIDStr,
		})
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to fetch attendances", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Attendances retrieved successfully", attendances))
}

func (ac *AttendanceController) CreateAttendance(c *gin.Context) {
	var req requests.AttendanceCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid request data", err.Error()))
		return
	}

	logger.LogInfo("Creating new attendance", logrus.Fields{
		"classroom_id": fmt.Sprintf("%d", req.ClassroomID),
		"student_id":   fmt.Sprintf("%d", req.StudentID),
		"session_date": req.SessionDate,
	})

	attendance, err := ac.attendanceService.CreateAttendance(&req)
	if err != nil {
		if err.Error() == "attendance for this student on this date already exists" {
			c.JSON(http.StatusBadRequest, response.ErrorResponse("Attendance already exists", err.Error()))
			return
		}
		logger.LogError(err, "Failed to create attendance", logrus.Fields{
			"classroom_id": fmt.Sprintf("%d", req.ClassroomID),
			"student_id":   fmt.Sprintf("%d", req.StudentID),
		})
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

	logger.LogInfo("Updating attendance", logrus.Fields{
		"attendance_id": idStr,
		"status":        string(req.Status),
	})

	attendance, err := ac.attendanceService.UpdateAttendance(uint(id), &req)
	if err != nil {
		if err.Error() == "attendance not found" {
			c.JSON(http.StatusNotFound, response.ErrorResponse("Attendance not found", err.Error()))
			return
		}
		logger.LogError(err, "Failed to update attendance", logrus.Fields{
			"attendance_id": idStr,
		})
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

	logger.LogInfo("Deleting attendance", logrus.Fields{
		"attendance_id": idStr,
	})

	err = ac.attendanceService.DeleteAttendance(uint(id))
	if err != nil {
		if err.Error() == "attendance not found" {
			c.JSON(http.StatusNotFound, response.ErrorResponse("Attendance not found", err.Error()))
			return
		}
		logger.LogError(err, "Failed to delete attendance", logrus.Fields{
			"attendance_id": idStr,
		})
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to delete attendance", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Attendance deleted successfully", nil))
}

package controller

import (
	"easy-attend-service/requests"
	"easy-attend-service/response"
	"easy-attend-service/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ClassroomController คือคอนโทรลเลอร์สำหรับจัดการห้องเรียน
type ClassroomController struct {
	classroomService *services.ClassroomService
}

// NewClassroomController สร้างอินสแตนซ์ใหม่ของ ClassroomController
func NewClassroomController() *ClassroomController {
	return &ClassroomController{
		classroomService: services.NewClassroomService(),
	}
}

// GetAllClassrooms ดึงข้อมูลห้องเรียนของครูที่ login
func (cc *ClassroomController) GetAllClassrooms(c *gin.Context) {
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

	classrooms, err := cc.classroomService.GetClassroomsByTeacher(uint(teacherID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to fetch classrooms", err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse("Classrooms retrieved successfully", classrooms))
}

// GetClassroomByID ดึงข้อมูลห้องเรียนตาม ID
func (cc *ClassroomController) GetClassroomByID(c *gin.Context) {
	id := c.Param("id")

	// Convert id from string to uint
	classroomID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid classroom ID", err.Error()))
		return
	}

	classroom, err := cc.classroomService.GetClassroomByID(uint(classroomID))
	if err != nil {
		if err.Error() == "classroom not found" {
			c.JSON(http.StatusNotFound, response.ErrorResponse("Classroom not found", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to fetch classroom", err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse("Classroom retrieved successfully", classroom))
}

// CreateClassroom สร้างห้องเรียนใหม่
func (cc *ClassroomController) CreateClassroom(c *gin.Context) {
	var req requests.ClassroomCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid request data", err.Error()))
		return
	}

	classroom, err := cc.classroomService.CreateClassroom(&req)
	if err != nil {
		if err.Error() == "classroom with this name already exists in this school" {
			c.JSON(http.StatusBadRequest, response.ErrorResponse("Classroom name already exists", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to create classroom", err.Error()))
		return
	}
	c.JSON(http.StatusCreated, response.SuccessResponse("Classroom created successfully", classroom))
}

// UpdateClassroom แก้ไขข้อมูลห้องเรียน
func (cc *ClassroomController) UpdateClassroom(c *gin.Context) {
	id := c.Param("id")
	var req requests.ClassroomUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid request data", err.Error()))
		return
	}

	// Convert id from string to uint
	classroomID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid classroom ID", err.Error()))
		return
	}

	classroom, err := cc.classroomService.UpdateClassroom(uint(classroomID), &req)
	if err != nil {
		if err.Error() == "classroom not found" {
			c.JSON(http.StatusNotFound, response.ErrorResponse("Classroom not found", err.Error()))
			return
		}
		if err.Error() == "classroom with this name already exists in this school" {
			c.JSON(http.StatusBadRequest, response.ErrorResponse("Classroom name already exists", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to update classroom", err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse("Classroom updated successfully", classroom))
}

// DeleteClassroom ลบห้องเรียน
func (cc *ClassroomController) DeleteClassroom(c *gin.Context) {
	id := c.Param("id")

	// Convert id from string to uint
	classroomID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid classroom ID", err.Error()))
		return
	}

	err = cc.classroomService.DeleteClassroom(uint(classroomID))
	if err != nil {
		if err.Error() == "classroom not found" {
			c.JSON(http.StatusNotFound, response.ErrorResponse("Classroom not found", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to delete classroom", err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse("Classroom deleted successfully", nil))
}

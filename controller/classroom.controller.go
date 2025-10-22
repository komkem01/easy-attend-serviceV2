package controller

import (
	"easy-attend-service/requests"
	"easy-attend-service/response"
	"easy-attend-service/services"
	"easy-attend-service/utils/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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

// GetAllClassrooms ดึงข้อมูลห้องเรียนทั้งหมด
func (cc *ClassroomController) GetAllClassrooms(c *gin.Context) {
	logger.LogInfo("Fetching all classrooms", nil)

	classrooms, err := cc.classroomService.GetAllClassrooms()
	if err != nil {
		logger.LogError(err, "Failed to fetch classrooms", nil)
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to fetch classrooms", err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse("Classrooms retrieved successfully", classrooms))
}

// GetClassroomByID ดึงข้อมูลห้องเรียนตาม ID
func (cc *ClassroomController) GetClassroomByID(c *gin.Context) {
	id := c.Param("id")
	logger.LogInfo("Fetching classroom by ID", logrus.Fields{"classroom_id": id})

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
		logger.LogError(err, "Failed to fetch classroom", logrus.Fields{"classroom_id": id})
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
	logger.LogInfo("Creating new classroom", logrus.Fields{"name": req.Name})

	classroom, err := cc.classroomService.CreateClassroom(&req)
	if err != nil {
		if err.Error() == "classroom with this name already exists in this school" {
			c.JSON(http.StatusBadRequest, response.ErrorResponse("Classroom name already exists", err.Error()))
			return
		}
		logger.LogError(err, "Failed to create classroom", logrus.Fields{"name": req.Name})
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

	logger.LogInfo("Updating classroom", logrus.Fields{"classroom_id": classroomID, "name": req.Name})

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
		logger.LogError(err, "Failed to update classroom", logrus.Fields{"classroom_id": classroomID})
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to update classroom", err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse("Classroom updated successfully", classroom))
}

// DeleteClassroom ลบห้องเรียน
func (cc *ClassroomController) DeleteClassroom(c *gin.Context) {
	id := c.Param("id")
	logger.LogInfo("Deleting classroom", logrus.Fields{"classroom_id": id})

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
		logger.LogError(err, "Failed to delete classroom", logrus.Fields{"classroom_id": id})
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to delete classroom", err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse("Classroom deleted successfully", nil))
}

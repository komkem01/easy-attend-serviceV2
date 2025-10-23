package controller

import (
	"easy-attend-service/requests"
	"easy-attend-service/response"
	"easy-attend-service/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ClassroomMemberController struct {
	classroomMemberService *services.ClassroomMemberService
}

func NewClassroomMemberController() *ClassroomMemberController {
	return &ClassroomMemberController{
		classroomMemberService: services.NewClassroomMemberService(),
	}
}

func (cmc *ClassroomMemberController) GetAllClassroomMembers(c *gin.Context) {
	members, err := cmc.classroomMemberService.GetAllClassroomMembers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to fetch classroom members", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Classroom members retrieved successfully", members))
}

func (cmc *ClassroomMemberController) GetClassroomMembersByClassroomID(c *gin.Context) {
	classroomIDStr := c.Param("classroom_id")

	// Convert string to uint
	classroomID, err := strconv.ParseUint(classroomIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid classroom ID", "ID must be a valid number"))
		return
	}
	members, err := cmc.classroomMemberService.GetClassroomMembersByClassroomID(uint(classroomID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to fetch classroom members", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Classroom members retrieved successfully", members))
}

func (cmc *ClassroomMemberController) CreateClassroomMember(c *gin.Context) {
	var req requests.ClassroomMemberCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid request data", err.Error()))
		return
	}

	member, err := cmc.classroomMemberService.CreateClassroomMember(&req)
	if err != nil {
		if err.Error() == "either teacher_id or student_id must be provided, but not both" ||
			err.Error() == "member already exists in this classroom" {
			c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid request", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to create classroom member", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.SuccessResponse("Classroom member created successfully", member))
}

func (cmc *ClassroomMemberController) UpdateClassroomMember(c *gin.Context) {
	classroomIDStr := c.Param("classroom_id")
	memberIDStr := c.Param("member_id")

	// Convert string to uint
	classroomID, err := strconv.ParseUint(classroomIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid classroom ID", "ID must be a valid number"))
		return
	}

	memberID, err := strconv.ParseUint(memberIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid member ID", "ID must be a valid number"))
		return
	}

	var req requests.ClassroomMemberUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid request data", err.Error()))
		return
	}
	member, err := cmc.classroomMemberService.UpdateClassroomMember(uint(classroomID), uint(memberID), &req)
	if err != nil {
		if err.Error() == "classroom member not found" {
			c.JSON(http.StatusNotFound, response.ErrorResponse("Classroom member not found", err.Error()))
			return
		}
		if err.Error() == "invalid member ID format" {
			c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid member ID", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to update classroom member", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Classroom member updated successfully", member))
}

func (cmc *ClassroomMemberController) DeleteClassroomMember(c *gin.Context) {
	classroomIDStr := c.Param("classroom_id")
	memberIDStr := c.Param("member_id")

	// Convert string to uint
	classroomID, err := strconv.ParseUint(classroomIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid classroom ID", "ID must be a valid number"))
		return
	}

	memberID, err := strconv.ParseUint(memberIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid member ID", "ID must be a valid number"))
		return
	}
	err = cmc.classroomMemberService.DeleteClassroomMember(uint(classroomID), uint(memberID))
	if err != nil {
		if err.Error() == "classroom member not found" {
			c.JSON(http.StatusNotFound, response.ErrorResponse("Classroom member not found", err.Error()))
			return
		}
		if err.Error() == "invalid member ID format" {
			c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid member ID", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to delete classroom member", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Classroom member deleted successfully", nil))
}

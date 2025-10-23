package controller

import (
	"easy-attend-service/requests"
	"easy-attend-service/response"
	"easy-attend-service/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type StudentController struct {
	studentService *services.StudentService
}

func NewStudentController() *StudentController {
	return &StudentController{
		studentService: services.NewStudentService(),
	}
}

func (sc *StudentController) GetAllStudents(c *gin.Context) {
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

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Get students for this teacher only
	students, total, err := sc.studentService.GetStudentsByTeacherPaginated(uint(teacherID), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to get students", err.Error()))
		return
	}

	result := gin.H{
		"students": students,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Students retrieved successfully", result))
}

func (sc *StudentController) GetStudentByID(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid request", "Student ID is required"))
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid student ID", "Student ID must be a positive integer"))
		return
	}

	student, err := sc.studentService.GetStudentByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse("Student not found", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Student retrieved successfully", student))
}

func (sc *StudentController) CreateStudent(c *gin.Context) {
	type StudentRequest struct {
		Firstname  string `json:"firstname" binding:"required"`
		Lastname   string `json:"lastname" binding:"required"`
		SchoolName string `json:"school_name" binding:"required"`
		StudentNo  string `json:"student_no"` // Optional, will auto-generate if empty
		GenderID   *uint  `json:"gender_id"`
		PrefixID   *uint  `json:"prefix_id"`
	}

	var req StudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid request data", err.Error()))
		return
	}

	var studentNoPtr *string
	if req.StudentNo != "" {
		studentNoPtr = &req.StudentNo
	}
	student, err := sc.studentService.TestCreateStudent(&req.SchoolName, &req.Firstname, &req.Lastname, studentNoPtr, req.GenderID, req.PrefixID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Failed to create student", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.SuccessResponse("Student created successfully", student))
}

func (sc *StudentController) UpdateStudent(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid request", "Student ID is required"))
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid student ID", "Student ID must be a positive integer"))
		return
	}

	var req requests.StudentUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid request data", err.Error()))
		return
	}

	student, err := sc.studentService.UpdateStudent(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Failed to update student", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Student updated successfully", student))
}

func (sc *StudentController) DeleteStudent(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid request", "Student ID is required"))
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid student ID", "Student ID must be a positive integer"))
		return
	}

	if err := sc.studentService.DeleteStudent(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Failed to delete student", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Student deleted successfully", nil))
}

// TestCreateStudent creates a student with auto-generated classroom for testing
func (sc *StudentController) TestCreateStudent(c *gin.Context) {
	type TestStudentRequest struct {
		Firstname  string `json:"firstname" binding:"required"`
		Lastname   string `json:"lastname" binding:"required"`
		SchoolName string `json:"school_name" binding:"required"`
		StudentNo  string `json:"student_no"` // Optional, will auto-generate if empty
		GenderID   *uint  `json:"gender_id"`
		PrefixID   *uint  `json:"prefix_id"`
	}

	var req TestStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid request data", err.Error()))
		return
	}

	var studentNoPtr *string
	if req.StudentNo != "" {
		studentNoPtr = &req.StudentNo
	}
	student, err := sc.studentService.TestCreateStudent(&req.SchoolName, &req.Firstname, &req.Lastname, studentNoPtr, req.GenderID, req.PrefixID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Failed to create student", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Student created successfully", student))
}
